// nova-launcher is a single-process supervisor that brings up every nova service
// for local development and tears them all down atomically when it exits.
//
// It NOT rely on a fragile PID file, and it makes
// `go run` safe to use. The old script leaked orphan processes (the compiled
// `main` child that `go run` forks kept holding the port after the script was
// killed) — which is exactly how you get ghost servers answering gRPC with
// errors like "unknown service nova.user.v1.UserService".
//
// The fix is the process group. Every service is started with Setpgid, so
// `go run` and the compiled binary it spawns share one process group (pgid ==
// the `go run` pid). Tearing down sends a signal to -<pgid>, which kills the
// whole group — go-run parent and compiled child alike — never leaving an
// orphan on the port.
//
// On top of that it:
//
//  1. starts services in dependency order with a tcp readiness probe;
//  2. forwards Ctrl+C / SIGTERM to every group, escalating to SIGKILL; and
//  3. stops everything (siblings included) the moment any service crashes, so a
//     partial/half-dead fleet can never silently linger again.
//
// Run it from anywhere under the repo:
//
//	cd nova-launcher && go run .
//
// or `make dev` from the repo root.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// service describes one local development process the supervisor manages.
type service struct {
	name string // short name, used in logs
	dir  string // service root, relative to repo root
	port string // tcp port for the (best-effort) readiness probe
	cmd  string // executable; empty means `go`
	args []string
	env  []string
}

// Ordered so dependencies start first: auth/user/post/file, gateway, web, then Vite.
var allServices = []service{
	{name: "nova-auth", dir: "nova-auth", port: "50051"},
	{name: "nova-user", dir: "nova-user", port: "50052"},
	{name: "nova-post", dir: "nova-post", port: "50053"},
	{name: "nova-file", dir: "nova-file", port: "50054"},
	{name: "nova-gateway", dir: "nova-gateway", port: "8080"},
	{name: "nova-web", dir: "nova-web", port: "8081", env: []string{"VITE_PORT=5173"}},
	{name: "nova-web-vite", dir: "nova-web", port: "5173", cmd: "npm", args: []string{"run", "dev"}, env: []string{"VITE_PORT=5173"}},
}

// goRunArgs is the `go run` invocation each service is launched with. Each
// service reads its config from `./config.yaml` in its own directory.
var goRunArgs = []string{"run", "./cmd/server", "-c", "./config.yaml"}

// ANSI colors, indexed by service position so each service gets a stable hue.
var palette = []string{
	"\x1b[36m", // nova-auth    cyan
	"\x1b[35m", // nova-user    magenta
	"\x1b[33m", // nova-post    yellow
	"\x1b[32m", // nova-file    green
	"\x1b[34m", // nova-gateway blue
	"\x1b[95m", // nova-web     bright magenta
	"\x1b[96m", // nova-web-vite bright cyan
}

const (
	colorReset = "\x1b[0m"
	colorDim   = "\x1b[2m"
	colorRed   = "\x1b[31m"
)

func main() {
	var (
		onlyFlag  = flag.String("only", "", "comma-separated subset of services to run (default: all)")
		probeWait = flag.Duration("probe", 6*time.Second, "per-service tcp readiness probe timeout")
	)
	flag.Parse()

	repoRoot, err := findRepoRoot()
	if err != nil {
		fatal(err)
	}

	services, err := selectServices(*onlyFlag)
	if err != nil {
		fatal(err)
	}

	out := newConsole(os.Stdout)
	sv := &supervisor{repoRoot: repoRoot, out: out, probe: *probeWait}

	// Tear everything down if we're interrupted, or if any child exits first.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	exitCode := 0
	select {
	case sig := <-sigCh:
		out.launcher("received %s, shutting down", sig)
		if sig == syscall.SIGINT {
			exitCode = 130
		} else {
			exitCode = 143
		}
	case ex := <-sv.start(services):
		exitCode = 1
		out.launcherErr("%s exited (code=%d); stopping the rest", ex.svc.name, ex.code)
	}

	sv.shutdown()

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

// ---- service selection & repo root ----------------------------------------

func selectServices(only string) ([]service, error) {
	if strings.TrimSpace(only) == "" {
		return allServices, nil
	}
	want := strings.Split(only, ",")
	byName := make(map[string]service, len(allServices))
	for _, s := range allServices {
		byName[s.name] = s
	}
	out := make([]service, 0, len(want))
	for _, raw := range want {
		name := strings.TrimSpace(raw)
		svc, ok := byName[name]
		if !ok {
			return nil, fmt.Errorf("unknown service %q (known: %s)", name, serviceNames(allServices))
		}
		out = append(out, svc)
	}
	return out, nil
}

func serviceNames(svcs []service) string {
	names := make([]string, len(svcs))
	for i, s := range svcs {
		names[i] = s.name
	}
	return strings.Join(names, ", ")
}

// findRepoRoot walks up from cwd until it sees the launcher + a couple of
// services as siblings, so the binary works whether run from the repo root or
// from inside nova-launcher.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for range 8 {
		if isRepoRoot(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("could not locate repo root (expected nova-launcher + nova-gateway siblings)")
}

func isRepoRoot(dir string) bool {
	for _, d := range []string{"nova-launcher", "nova-gateway", "nova-auth"} {
		if info, err := os.Stat(filepath.Join(dir, d)); err != nil || !info.IsDir() {
			return false
		}
	}
	return true
}

// ---- supervisor ------------------------------------------------------------

type supervisor struct {
	repoRoot string
	out      *console
	probe    time.Duration

	mu       sync.Mutex
	procs    []*managedProc
	stopping atomic.Bool
}

// exit carries the first service that stopped, to the main select.
type exit struct {
	svc  service
	code int
}

// start launches services in order (respecting the configured probe) and
// returns a channel that fires when the FIRST service exits. The caller is
// expected to invoke shutdown() afterwards.
func (s *supervisor) start(svcs []service) <-chan exit {
	exitCh := make(chan exit, 1)
	var started []*managedProc

	for idx, svc := range svcs {
		mp, err := s.startOne(svc, idx)
		if err != nil {
			s.out.launcherErr("start %s failed: %v", svc.name, err)
			// Treat a launch failure like a crash: tear down what's already up.
			s.procs = started
			exitCh <- exit{svc: svc, code: -1}
			return exitCh
		}
		started = append(started, mp)
		s.mu.Lock()
		s.procs = started
		s.mu.Unlock()

		if s.probe > 0 {
			if waitReady("127.0.0.1:"+svc.port, s.probe) {
				s.out.svc(svc.name, "ready on :%s", svc.port)
			} else {
				s.out.svcWarn(svc.name, "no response on :%s after %s (continuing)", svc.port, s.probe)
			}
		}
	}
	s.out.launcher("all services up — Ctrl+C to stop")

	// One watcher per process; the first to exit wins.
	for _, mp := range started {
		go func(mp *managedProc) {
			err := mp.cmd.Wait()
			code := exitCodeOf(err)
			if s.stopping.Load() {
				return // expected shutdown, don't report
			}
			select {
			case exitCh <- exit{svc: mp.svc, code: code}:
			default:
			}
		}(mp)
	}
	return exitCh
}

func (s *supervisor) startOne(svc service, idx int) (*managedProc, error) {
	name, args := svc.command()
	cmd := exec.Command(name, args...)
	cmd.Dir = filepath.Join(s.repoRoot, svc.dir)
	if len(svc.env) > 0 {
		cmd.Env = append(os.Environ(), svc.env...)
	}
	// Each child leads its own process group => pgid == pid, so we can signal
	// the whole group (go-run parent + compiled binary, plus any descendants)
	// via syscall.Kill(-pid, ...). This is what makes `go run` orphan-free.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	out := s.out.writerFor(svc.name, idx)
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Start(); err != nil {
		_ = out.Close()
		return nil, err
	}
	s.out.launcher("started %s via `%s` (pid %d)", svc.name, svc.commandLine(), cmd.Process.Pid)
	return &managedProc{svc: svc, cmd: cmd, out: out}, nil
}

func (svc service) command() (string, []string) {
	if svc.cmd != "" {
		return svc.cmd, svc.args
	}
	return "go", goRunArgs
}

func (svc service) commandLine() string {
	name, args := svc.command()
	parts := append([]string{name}, args...)
	return strings.Join(parts, " ")
}

// shutdown sends SIGTERM to every process group in reverse startup order,
// waits briefly, then SIGKILLs anything still alive.
func (s *supervisor) shutdown() {
	s.stopping.Store(true)

	s.mu.Lock()
	procs := make([]*managedProc, len(s.procs))
	copy(procs, s.procs)
	s.mu.Unlock()

	// Stop in reverse so dependents (web, gateway) go down before backends.
	for i, j := 0, len(procs)-1; i < j; i, j = i+1, j-1 {
		procs[i], procs[j] = procs[j], procs[i]
	}

	const graceful = 8 * time.Second
	deadline := time.Now().Add(graceful)

	for _, mp := range procs {
		if mp.cmd.Process == nil {
			continue
		}
		_ = signalGroup(mp.cmd.Process.Pid, syscall.SIGTERM)
	}

	// Wait for each process to report exit (its cmd.Wait goroutine will fire).
	for _, mp := range procs {
		remaining := max(time.Until(deadline), 0)
		if !waitForExit(mp, remaining) {
			s.out.launcherErr("%s did not exit, sending SIGKILL", mp.svc.name)
			_ = signalGroup(mp.cmd.Process.Pid, syscall.SIGKILL)
			_ = waitForExit(mp, 2*time.Second)
		}
		if mp.out != nil {
			_ = mp.out.Close() // flush any partial trailing line
		}
	}
	s.out.launcher("all services stopped")
}

// waitForExit polls whether the process is gone. (cmd.Wait is already running
// in a goroutine; we just need to know when the PID disappears.)
func waitForExit(mp *managedProc, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if mp.cmd.ProcessState != nil {
			return true
		}
		// Sending signal 0 checks process existence without killing it.
		if err := mp.cmd.Process.Signal(syscall.Signal(0)); err != nil {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

type managedProc struct {
	svc service
	cmd *exec.Cmd
	out io.WriteCloser
}

// ---- helpers ---------------------------------------------------------------

func signalGroup(pid int, sig syscall.Signal) error {
	// Negative pid => signal the whole process group.
	return syscall.Kill(-pid, sig)
}

func exitCodeOf(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return -1
}

func waitReady(addr string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 400*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return true
		}
		time.Sleep(150 * time.Millisecond)
	}
	return false
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%snova-launcher:%s %v\n", colorRed, colorReset, err)
	os.Exit(1)
}

// ---- console (prefixed, line-buffered, colorized) -------------------------

type console struct {
	mu sync.Mutex
	w  io.Writer
}

func newConsole(w io.Writer) *console { return &console{w: w} }

// writerFor returns a line-buffered writer that tags every line of a service's
// stdout/stderr with its name and color.
func (c *console) writerFor(name string, idx int) io.WriteCloser {
	color := palette[idx%len(palette)]
	return &lineWriter{c: c, color: color, tag: fmt.Sprintf("%-12s", "["+name+"]")}
}

func (c *console) launcher(format string, args ...any) {
	c.emit(colorDim, "[launcher] ", fmt.Sprintf(format, args...))
}

func (c *console) launcherErr(format string, args ...any) {
	c.emit(colorRed, "[launcher] ", fmt.Sprintf(format, args...))
}

func (c *console) svc(name string, format string, args ...any) {
	c.emit(svcColor(name), fmt.Sprintf("%-12s", "["+name+"]"), fmt.Sprintf(format, args...))
}

func (c *console) svcWarn(name string, format string, args ...any) {
	c.emit(colorRed, fmt.Sprintf("%-12s", "["+name+"]"), fmt.Sprintf(format, args...))
}

func (c *console) emit(color, tag, body string) {
	ts := time.Now().Format("15:04:05")
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Fprintf(c.w, "%s%s%s %s%s│%s %s\n", color, tag, colorReset, colorDim, ts, colorReset, body)
}

func svcColor(name string) string {
	for i, s := range allServices {
		if s.name == name {
			return palette[i%len(palette)]
		}
	}
	return ""
}

// lineWriter buffers writes until a newline, then emits each line through the
// console with the service's color/tag. Partial trailing data is flushed on Close.
type lineWriter struct {
	c     *console
	color string
	tag   string
	buf   []byte
}

func (lw *lineWriter) Write(p []byte) (int, error) {
	lw.buf = append(lw.buf, p...)
	for {
		i := bytes.IndexByte(lw.buf, '\n')
		if i < 0 {
			break
		}
		line := string(bytes.TrimRight(lw.buf[:i], "\r"))
		lw.buf = lw.buf[i+1:]
		lw.c.emit(lw.color, lw.tag, line)
	}
	return len(p), nil
}

func (lw *lineWriter) Close() error {
	if len(lw.buf) > 0 {
		lw.c.emit(lw.color, lw.tag, string(lw.buf))
		lw.buf = nil
	}
	return nil
}
