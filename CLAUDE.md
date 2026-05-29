# CLAUDE.md

## 项目概述

goc-quickstart 是一个基于 gRPC 的博客微服务项目，包含认证、用户和文章三个服务。

## 项目结构

```
goc-quickstart
├── apis/                    # Proto 定义和生成代码
│   ├── proto/blog/          # Proto 文件 (auth, post, user)
│   └── gen/                 # 生成的 Go 代码
├── auth-service/             # 认证服务 (gRPC)
├── user-service/             # 用户服务 (gRPC)
├── post-service/             # 文章服务 (gRPC)
├── api-gateway/              # HTTP 网关
├── apidoc-server/          # API 文档服务
└── README.md
```

## Proto 管理

- Proto 文件统一在 `apis/proto/blog/` 管理
- 各服务通过 `buf generate` 生成代码后复制到服务的 `gen/` 目录
- package 命名: `goc.blog.{service}.api.v1`
- service 命名: `{Service}Service` (如 `AuthService`, `PostService`)

## 服务命名规范

| 服务 | 目录 | gRPC 端口 |
|------|------|-----------|
| auth-service | `github.com/miiy/goc-quickstart/auth-service` | 50051 |
| user-service | `github.com/miiy/goc-quickstart/user-service` | 50053 |
| post-service | `github.com/miiy/goc-quickstart/post-service` | 50052 |
| api-gateway | `github.com/miiy/goc-quickstart/api-gateway` | 8080 |

## Service 层规范

- Struct: `XxxService` (如 `AuthService`, `PostService`)
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

- 各服务 `configs/default.yaml` 不包含 HTTP 配置（gateway 统一处理）
- api-gateway 暴露 HTTP 入口，各微服务仅暴露 gRPC
- redis 配置不包含 username（仅 password）

## Wire DI

- `wire.go` 定义依赖注入
- `wire_gen.go` 通过 `wire ./internal/di/` 命令生成，不要手动编辑

## 数据库

- Migration 文件放在 `{service}/deploy/migrations/`
- 表命名: 复数形式 (如 `users`, `posts`)
- 使用 gorm.Model 包含 `id`, `created_at`, `updated_at`, `deleted_at`

## 常见操作

```bash
# 生成 proto 代码
cd apis && buf generate

# 复制 proto 文件到各服务
cp apis/gen/go/blog/{service}/v1/* {service}-service/gen/go/blog/{service}/v1/

# 生成 wire 依赖注入
cd {service}-service && wire ./internal/di/

# 构建服务
cd {service}-service && go build ./...
```
