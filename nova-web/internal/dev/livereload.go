package dev

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/miiy/goc/gin"
)

type LiveReload struct {
	mu         sync.Mutex
	clients    map[chan struct{}]struct{}
	watchRoots []string
	lastMod    time.Time
}

func NewLiveReload(watchRoots ...string) *LiveReload {
	if len(watchRoots) == 0 {
		watchRoots = []string{"dist"}
	}

	return &LiveReload{
		clients:    make(map[chan struct{}]struct{}),
		watchRoots: append([]string(nil), watchRoots...),
	}
}

func (lr *LiveReload) RegisterRouter(r gin.IRouter) {
	go lr.watch()
	r.GET("/__dev/livereload", lr.events)
}

func (lr *LiveReload) events(c *gin.Context) {
	ch := make(chan struct{}, 1)
	lr.add(ch)
	defer lr.remove(ch)

	w := c.Writer
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	fmt.Fprint(w, ": connected\n\n")
	flusher.Flush()

	notify := c.Request.Context().Done()
	for {
		select {
		case <-notify:
			return
		case <-ch:
			fmt.Fprint(w, "event: reload\ndata: now\n\n")
			flusher.Flush()
		}
	}
}

func (lr *LiveReload) add(ch chan struct{}) {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	lr.clients[ch] = struct{}{}
}

func (lr *LiveReload) remove(ch chan struct{}) {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	delete(lr.clients, ch)
	close(ch)
}

func (lr *LiveReload) watch() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		modTime, err := newestModTime(lr.watchRoots...)
		if err != nil {
			continue
		}
		if modTime.IsZero() {
			continue
		}
		if lr.lastMod.IsZero() {
			lr.lastMod = modTime
			continue
		}
		if !modTime.After(lr.lastMod) {
			continue
		}
		lr.lastMod = modTime
		lr.broadcast()
	}
}

func (lr *LiveReload) broadcast() {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	for ch := range lr.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func newestModTime(roots ...string) (time.Time, error) {
	var newest time.Time
	for _, root := range roots {
		if root == "" {
			continue
		}
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if info.ModTime().After(newest) {
				newest = info.ModTime()
			}
			return nil
		})
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return newest, err
		}
	}
	return newest, nil
}
