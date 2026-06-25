# Cookie Session Data Flow

`nova-web` 使用 Redis-backed session。浏览器只保存 session id，完整 session 数据保存在 Redis。

## 存储位置

```text
Browser Cookie
  COOKIE_SESSION = securecookie(session.ID)

Redis
  session_{session.ID} = session.Values
```

Redis key 示例：

```text
session_JUQ57UX6FEVPK2R6SPL4CXEW3JJXMVPLSYHY734NLOGKHNFL6NVA
```

规则：

- `COOKIE_SESSION`：cookie 名，来自 `nova-web/config.yaml` 的 `session.name`。
- `session_`：Redis key 前缀，来自底层 `boj/redistore` 默认值。
- `session.ID`：底层生成的随机 ID，规则是 `32 bytes random -> base32 -> trim "=" padding`。

## 首次打开页面

以打开 `/login` 为例：

```text
Browser -> GET /login
nova-web session middleware -> 创建空 session 对象
csrf.Token(...) -> 写入 _csrf_token
session.Save() -> Redis 写入 session_{session.ID}
Response -> Set-Cookie: COOKIE_SESSION=securecookie(session.ID)
```

此时 Redis 里主要有：

```text
_csrf_token
```

## 登录成功

登录表单提交时：

```text
Browser -> POST /login + COOKIE_SESSION + _csrf
nova-web -> 根据 COOKIE_SESSION 解出 session.ID
nova-web -> 从 Redis 读取 session_{session.ID}
CSRF middleware -> 校验 _csrf
nova-web -> 调用 nova-gateway /api/v1/auth/login
nova-auth -> 返回 access_token + refresh_token
nova-web SaveLogin -> 写入新的登录 session
Response -> Set-Cookie: COOKIE_SESSION=securecookie(new session.ID)
```

登录后的 Redis session value 主要有：

```text
goc.auth.user
access_token
access_expires_at
refresh_token
```

## 普通请求

登录后的页面请求：

```text
Browser -> Request + COOKIE_SESSION
nova-web -> 解出 session.ID
nova-web -> 读取 Redis session_{session.ID}
auth.RefreshSessionToken -> 注入 access_token 到 request context
sessionauth.LoadSessionUser -> 从 session.Values 注入 current user
Handler -> 使用当前用户和 access_token 调用后端
```

## Token 刷新

当 access token 快过期时：

```text
Browser -> Request + COOKIE_SESSION
nova-web -> 读取 Redis session
auth.RefreshSessionToken -> 发现 access_token 快过期
auth.RefreshSessionToken -> 使用 refresh_token 调用 /api/v1/auth/token/refresh
nova-auth -> 返回新的 access_token + refresh_token
nova-web -> 更新同一个 session 的 token 字段
Redis -> 覆盖 session_{session.ID} 的 value，并刷新 TTL
```

更新的字段：

```text
access_token
access_expires_at
refresh_token
```

## 退出或失效

退出登录：

```text
Browser -> POST /auth/logout
nova-web -> 调用 nova-auth logout，撤销 token
nova-web -> 清空 session
Redis -> 删除 session_{session.ID}
Response -> Set-Cookie: COOKIE_SESSION=; Max-Age=-1
```

refresh token 无效或被撤销时：

```text
webauth.RefreshSessionToken -> refresh 失败且返回 401
nova-web -> 清空 session
Redis -> 删除 session_{session.ID}
Browser -> 后续需要重新登录
```

超过 `session.maxAge` 时：

```text
Redis -> session_{session.ID} 自动过期
Browser cookie -> 到期后不再发送
```
