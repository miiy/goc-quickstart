# nova-launcher

本地开发用的单进程 supervisor：用 `go run` 一次性拉起 nova 全部服务，并在它退出时
把所有子进程**原子地**一起回收。

## 为什么能做到零孤儿

关键在**进程组**：每个服务以 `Setpgid` 启动，`go run` 父进程和它编译出的 `main`
子进程落在**同一个进程组**（pgid == `go run` 的 pid）。回收时向 `-<pgid>` 发信号，
父进程和子进程一起被带走，绝不留孤儿。所以开发阶段无需预编译，全部走 `go run`，
改完代码下次启动自动重新编译。

除此之外它还：

- **顺序启动 + 就绪探测**：按 auth→user→post→file→gateway→web 顺序启动，每步做 TCP
   就绪探测，避免 gateway 在后端未就绪时连接。
- **崩溃即全停**：任一服务意外退出，supervisor 立刻把其余服务一并停掉，杜绝「半死」的
   进程组合悄悄残留。
- **优雅回收**：Ctrl+C / SIGTERM 先向每个进程组发 `SIGTERM`，超时则升级为 `SIGKILL`。

## 用法

从仓库任意位置：

```bash
# 方式一：根目录 Makefile
make dev                            # 拉起全部
make dev ONLY=nova-auth,nova-user   # 只起子集

# 方式二：直接运行
cd nova-launcher && go run .
cd nova-launcher && go run . -only nova-auth,nova-gateway
```

`Ctrl+C`（SIGINT）或 `SIGTERM` 触发回收。

## 配置

端口与启动顺序硬编码在 `main.go` 的 `allServices` 中，与本仓库各服务 `config.yaml`
保持一致（auth 50051 / user 50052 / post 50053 / file 50054 / gateway 8080 / web 8081）。
