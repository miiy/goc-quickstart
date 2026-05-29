# goc-quickstart

基于 gRPC 的博客微服务项目，采用 Go 语言开发。

## 项目架构

```
┌─────────────┐     ┌─────────────┐
│    Web      │────▶│ API Gateway │
│  (前端)     │     │  (HTTP网关) │
└─────────────┘     └──────┬──────┘
                           │ gRPC
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
   │Auth Service │  │Post Service │  │User Service │
   │  认证服务   │  │  文章服务   │  │  用户服务   │
   └─────────────┘  └─────────────┘  └─────────────┘
          │                │                │
          └────────────────┴────────────────┘
                           │
                    ┌──────┴──────┐
                    │   MySQL     │
                    │   Redis     │
                    └─────────────┘
```

## 项目结构

```
goc-quickstart
├── apis/                    # Proto 定义和生成代码
│   ├── proto/blog/          # Proto 文件
│   └── gen/                 # 生成的 Go/OpenAPI 代码
├── auth-service/            # 认证服务 (gRPC :50051)
├── user-service/            # 用户服务 (gRPC :50053)
├── post-service/            # 文章服务 (gRPC :50052)
├── api-gateway/             # HTTP 网关 (:8080)
├── web/                     # Web 前端服务 (:8081)
├── apidoc-server/           # API 文档服务
└── README.md
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| api-gateway | 8080 | HTTP 入口，路由转发 |
| web | 8081 | Web 前端页面 |
| auth-service | 50051 | 认证/登录/注册 |
| post-service | 50052 | 文章 CRUD |
| user-service | 50053 | 用户信息管理 |

## 技术栈

- **语言**: Go 1.26
- **RPC**: gRPC + gRPC-Gateway
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **依赖注入**: Wire
- **配置**: Viper
- **日志**: Zap
- **Proto**: Buf

## 快速开始

### 前置要求

- Go 1.26+
- MySQL 8.0+
- Redis 7.0+
- Buf CLI (用于 proto 生成)

### 安装依赖

```bash
# 安装 buf
go install github.com/bufbuild/buf/cmd/buf@latest

# 安装 wire
go install github.com/google/wire/cmd/wire@latest
```

### 生成代码

```bash
# 生成 proto 代码
cd apis && buf generate

# 复制生成代码到各服务
cp -r apis/gen/go/blog/auth/v1 auth-service/gen/go/blog/auth/v1/
cp -r apis/gen/go/blog/post/v1 post-service/gen/go/blog/post/v1/
cp -r apis/gen/go/blog/user/v1 user-service/gen/go/blog/user/v1/
cp -r apis/gen/go/blog/* api-gateway/gen/go/blog/

# 生成 wire 依赖注入
cd auth-service && wire ./internal/di/
cd post-service && wire ./internal/di/
cd user-service && wire ./internal/di/
```

### 数据库初始化

```sql
CREATE DATABASE goc CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

各服务的 migration 文件位于 `{service}/deploy/migrations/` 目录。

### 启动服务

```bash
# 启动 auth-service
cd auth-service && go run cmd/server/main.go

# 启动 user-service
cd user-service && go run cmd/server/main.go

# 启动 post-service
cd post-service && go run cmd/server/main.go

# 启动 api-gateway
cd api-gateway && go run cmd/main.go

# 启动 web
cd web && go run cmd/server/main.go
```

### 访问服务

- Web 界面: http://localhost:8081
- API 网关: http://localhost:8080
- API 文档: 查看 `apis/gen/openapiv2/` 目录

## 配置说明

各服务配置文件位于 `{service}/configs/default.yaml`：

```yaml
app:
  name: blog
  env: local
  debug: true

database:
  driver: mysql
  host: 127.0.0.1
  port: 3306
  database: goc
  username: root
  password: 123456

redis:
  addrs:
    - 127.0.0.1:6379
  password:
  db: 0

server:
  grpc:
    addr: 0.0.0.0:50051

jwt:
  secret: 123456
  expiresIn: 7200
```

## API 示例

### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","username":"testuser","password":"password123","password_confirmation":"password123"}'
```

### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### 获取文章列表

```bash
curl http://localhost:8080/api/v1/posts
```

## 开发指南

### 添加新服务

1. 在 `apis/proto/blog/` 下创建 proto 文件
2. 运行 `buf generate` 生成代码
3. 创建服务目录，实现 `pb.XxxServiceServer` 接口
4. 添加 wire 依赖注入
5. 在 api-gateway 注册服务路由

### 代码规范

- Service 层: `internal/service/xxx_service.go`
- Repository 层: `internal/repository/xxx_repository.go`
- Entity: `internal/entity/xxx.go`
- 使用 `status.Error()` 返回 gRPC 错误

## License

MIT
