# CLAUDE.md

## 项目概述

nova 是一个基于 gRPC 的博客微服务项目，包含认证、用户、文章和文件四个服务。

## 协作原则

- 如果用户提出的方案不符合业内最佳实践、项目既有架构边界、可维护性或安全性要求，必须明确指出问题、风险和推荐替代方案，避免在错误方向上继续深入。
- 不要为了迎合需求直接实现明显不合理的方案；如果最终仍按用户方案执行，需要先说明取舍和后续风险。
- 优先遵循本项目已经形成的边界：HTTP 入口在 `nova-gateway`，业务服务只暴露 gRPC，通用能力沉淀到 sibling repo `../goc`。

## 项目结构

```
goc-quickstart
├── nova-proto/                    # Proto 定义和生成代码
│   ├── proto/nova/                # Proto 文件 (auth, post, user, file)
│   └── gen/                       # 生成的 Go/OpenAPI 代码
├── nova-auth/                     # 认证服务 (gRPC)
├── nova-user/                     # 用户服务 (gRPC)
├── nova-post/                     # 文章服务 (gRPC)
├── nova-file/                     # 文件服务 (gRPC)
├── nova-gateway/                  # HTTP 网关
├── nova-web/                      # Web 前端服务
├── nova-apidoc/                   # API 文档服务
└── README.md
```

## 通用库依赖

- 本项目通过 `replace github.com/miiy/goc => ../../goc` 使用 sibling repo `/Users/mac/IdeaProjects/goc`。
- 修改认证、Gin middleware、gRPC server/interceptor、mTLS helper 等通用能力时，优先在 `../goc` 内收敛，再在 quickstart 中接入。
- `goc/grpc/gateway` 只表示 grpc-gateway reverse proxy 相关能力，不应承载认证、metadata 提取或 mTLS helper。
- gRPC 客户端 mTLS helper 放在 `goc/grpc/credentials`。
- gRPC metadata 认证提取放在 `goc/grpc/interceptor/auth`。

## Proto 管理

- Proto 文件统一在 `nova-proto/proto/nova/` 管理
- 各服务通过 `buf generate` 生成代码后复制到服务的 `gen/` 目录
- package 命名: `nova.{service}.v1`
- service 命名: `{Service}Service` (如 `AuthService`, `PostService`, `FileService`)
- 保留 `protoc-gen-openapiv2` 用于生成 OpenAPI 文档；不要启用 `protoc-gen-grpc-gateway` 生成 `*.pb.gw.go`，运行时 HTTP 网关由 Gin 手写 handler 承担。

## API Gateway 与认证边界

- `nova-gateway` 使用 Gin 手写 HTTP 路由和 handler，不使用 grpc-gateway generated mux 作为主链路。
- grpc-gateway 不负责认证；认证由 `nova-gateway` 的 Gin middleware 完成。
- JWT middleware 需要完成两步：
  - 校验 JWT 签名、issuer、过期时间等 token claims。
  - 通过 nova-auth 的 `GetAuthenticatedUser` 做当前用户二次校验，只允许 active 用户继续访问。
- 不要把原始 access token 继续传给下游业务服务作为主认证手段；gateway 是认证边界，下游接收已验证后的身份结果。
- Gin 认证成功后，通过 gRPC metadata 向下游传递最小身份信息：
  - `x-auth-user-id`
  - `x-auth-username`
- 下游服务通过 `goc/grpc/interceptor/auth.MetadataAuthFunc` 读取 metadata 并注入 `goc/auth.AuthenticatedUser` 到 `context`。
- 下游 gRPC auth interceptor 只匹配受保护 RPC；公开 RPC 不应要求认证 metadata。
- 下游业务代码需要当前用户时，从 `goc/auth.ExtractAuthenticatedUser(ctx)` 获取，不要直接解析 metadata。
- metadata 只传稳定、必要、非敏感的身份字段；不要传 email、phone、完整用户对象或原始 token。

## 路由保护约定

- nova-post：
  - 公开：`ListPosts`, `GetPost`
  - 保护：`CreatePost`, `UpdatePost`, `DeletePost`
- nova-user：
  - 当前 user 路由都挂在 nova-gateway protected group 下，对应 gRPC 方法需要 metadata。
- nova-auth：
  - 公开：注册、登录、字段检查、短信登录、`GetAuthenticatedUser`、`RefreshToken`（refresh token 本身是凭证；access 过期后必须可达，安全由 opaque refresh + rotation + reuse detection 保证）
  - 保护：`Logout`, `ChangePassword`
- nova-file：
  - 保护：`CreateFile`（所有文件上传需要认证）

## 服务命名规范

| 服务         | 目录                                            | gRPC 端口 |
| ------------ | ----------------------------------------------- | --------- |
| nova-auth    | `github.com/miiy/goc-quickstart/nova-auth`    | 50051     |
| nova-user    | `github.com/miiy/goc-quickstart/nova-user`    | 50052     |
| nova-post    | `github.com/miiy/goc-quickstart/nova-post`    | 50053     |
| nova-file    | `github.com/miiy/goc-quickstart/nova-file`    | 50054     |
| nova-gateway | `github.com/miiy/goc-quickstart/nova-gateway` | 8080      |
| nova-web     | `github.com/miiy/goc-quickstart/nova-web`     | 8081      |

## Service 层规范

- Struct: `XxxService` (如 `AuthService`, `PostService`, `FileService`)
- 构造函数: `NewXxxServiceServer` (返回 `pb.XxxServiceServer`)
- Getter: `XxxService()` (如 `AuthService()`)
- 接口: 使用 proto 生成的 `pb.XxxServiceServer`

## Repository 层规范

- 接口: `XxxRepository`
- 方法顺序: First → List → Create → Update → Delete
- First 支持可选 columns 参数
- List 支持 filter + 分页

## 错误处理

- 使用 `google.golang.org/grpc/codes` + `status.Error()`
- 常见模式:
  ```go
  if err != nil {
      if errors.Is(err, gorm.ErrRecordNotFound) {
          return nil, status.Error(codes.NotFound, ErrXxxNotFound.Error())
      }
      s.logger.Error("repo.First", zap.Error(err))
      return nil, status.Error(codes.Internal, err.Error())
  }
  ```

## 配置

- 各服务从目录根部的 `config.yaml` 读取配置，并提交 `config.yaml.example` 作为模板；实际 `config.yaml` 为本地运行配置，不提交。
- 本地 mTLS 证书和私钥放在各服务的 `configs/x509/`，`configs/` 整体忽略，不提交证书材料。
- nova-gateway 暴露 HTTP 入口，各微服务仅暴露 gRPC
- redis 配置不包含 username（仅 password）

## Wire DI

- `wire.go` 定义依赖注入
- `wire_gen.go` 通过 `wire ./internal/di/` 命令生成，不要手动编辑

## 数据库

- Migration 文件放在 `{service}/deploy/migrations/`
- 表命名: 复数形式 (如 `users`, `posts`, `files`)
- 使用 gorm.Model 包含 `id`, `created_at`, `updated_at`, `deleted_at`

## 常见操作

使用根目录 Makefile 统一管理：

```bash
# Proto 相关
make proto              # 生成 proto 代码 + 复制到各服务
make proto-generate     # 仅生成 proto 代码
make proto-copy         # 仅复制生成代码到各服务

# 构建
make build              # 构建所有服务（二进制命名为服务名，如 bin/nova-auth）
make wire               # 生成所有服务的 Wire DI

# 测试
make test               # 运行所有服务测试
make lint               # 运行 lint 检查

# 本地开发：单进程 supervisor 拉起全部服务（go run，Ctrl+C 原子回收，无孤儿）
make dev                       # 启动全部
make dev ONLY=nova-auth,nova-user   # 只启动子集

# 单个服务运行
make dev-auth           # 运行 nova-auth
make dev-user           # 运行 nova-user
make dev-post           # 运行 nova-post
make dev-file           # 运行 nova-file
make dev-gateway        # 运行 nova-gateway
make dev-web            # 运行 nova-web
make dev-apidoc         # 运行 nova-apidoc

# Docker
make docker-up          # 启动 docker-compose
make docker-down        # 停止 docker-compose
make docker-build       # 构建 docker 镜像
```

各服务内部也有独立 Makefile，支持 `make build/run/test/fmt/clean`。
