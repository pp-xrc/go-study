# JWT 入门速通（教学版）

> 配合本课 `26-JWT/main.go` 学习。  
> 目标：搞清 JWT 是什么、长什么样、登录后怎么用它保护接口。

库文档：https://github.com/golang-jwt/jwt

---

## 1. JWT 是什么？

**JWT**（JSON Web Token）是一种**令牌**：登录成功后，服务端发给客户端的一串字符串。

之后每次请求带上这串 token，服务端就能认出「你是谁」，**不用再查 session、也不用把登录状态存在服务器**。

一句话：

> JWT = 一张带签名的「身份证」，里面写着用户 ID 和过期时间。

| 对比 | Session | JWT |
|------|---------|-----|
| 登录状态存在哪 | 服务器（内存 / Redis） | 客户端（token 本身） |
| 服务端要不要记 token | 要 | 不要（无状态） |
| 多台服务器 | 要共享 session | 只要密钥相同就能校验 |

本课用的就是无状态 JWT：服务端只负责**签发**和**验签**。

---

## 2. 长什么样？

JWT 是三段用 `.` 拼起来的字符串：

```text
Header.Payload.Signature
xxxxx.yyyyy.zzzzz
```

| 部分 | 含义 | 通俗理解 |
|------|------|----------|
| Header | 算法（如 HS256）、类型 JWT | 「这张证件用什么章」 |
| Payload | 载荷：user_id、过期时间等 | 「证件上写的信息」 |
| Signature | 用密钥对前两段做签名 | 「防伪造的钢印」 |

Payload **不是加密**，只是 Base64 编码。任何人拿到 token 都能解码看到内容，所以：

- 不要把密码放进 JWT
- 靠 **Signature + 密钥** 保证没被篡改

本课 Payload 里主要是：`user_id` + 过期时间（24 小时）。

---

## 3. 典型登录流程

```text
1. 前端 POST /login（账号密码，本课简化为 user_id）
2. 服务端校验通过 → GenerateToken(userID) → 返回 JWT
3. 前端把 token 存到 localStorage 或 cookie
4. 之后请求带请求头：
     Authorization: Bearer <token>
5. 中间件 VerifyToken → 取出 userID → 放行
6. 没 token / 伪造 / 过期 → 401
```

对应本课接口：

| 接口 | 要不要登录 |
|------|------------|
| `POST /login` | 否（用来换 token） |
| `GET /ping` | 否 |
| `GET /user/list` | 是 |
| `GET /user/:id` | 是 |

---

## 4. 本课两个核心函数

```go
GenerateToken(userID int64) (string, error)  // 签发
VerifyToken(tokenString string) (int64, error) // 校验，返回 userID
```

- **签发**：把 `user_id`、过期时间写进 Payload，用密钥 HS256 签名  
- **校验**：算法必须是 HMAC、签名对得上、没过期，才返回 `userID`

Gin 中间件做的事：

1. 读 `Authorization` 头  
2. 拆出 `Bearer ` 后面的 token  
3. 调用 `VerifyToken`  
4. 成功则 `c.Set("userID", id)`，失败则 401 并 `Abort`

---

## 5. JWT 用在哪？

日常开发里最常见的几种：

### 5.1 登录鉴权（本课重点）

用户登录后拿 token，访问 `/user/list` 等私有接口。  
适合前后端分离：Vue / React / 小程序 / App 都把 token 放请求头。

### 5.2 API 保护

内部管理后台、个人中心、下单接口：没带有效 JWT 一律拒绝。

### 5.3 单点 / 多端登录（简化理解）

同一套密钥，Web、App、小程序都能用同一套「登录换 token」逻辑。  
（真正的单点登录还会涉及刷新令牌、退出黑名单等，本课不展开。）

### 5.4 微服务之间传身份

网关验完 JWT，把 `user_id` 转给后面的服务，后面不用再查登录态。

---

## 6. 和本课代码怎么对上

| 概念 | 代码 |
|------|------|
| 密钥 | `jwtSecret`（正式项目用环境变量） |
| 24 小时过期 | `tokenTTL = 24 * time.Hour` |
| 载荷 | `CustomClaims{UserID, RegisteredClaims}` |
| 签发 | `GenerateToken` |
| 验签 | `VerifyToken` |
| 从请求头取 token | `JWTAuthMiddleware` |
| 受保护路由 | `auth.Use(JWTAuthMiddleware())` |

本地试跑：

```bash
go run ./26-JWT

curl -s -X POST http://localhost:8086/login \
  -H "Content-Type: application/json" \
  -d '{"user_id":1}'

curl -s http://localhost:8086/user/list \
  -H "Authorization: Bearer <上一步返回的 token>"
```

---

## 7. 使用注意（入门够用）

1. **密钥保密**：泄露后别人能自己签发「合法」token。  
2. **Payload 别放敏感信息**：能被解码。  
3. **一定要设过期时间**：本课 24 小时；越敏感的业务越短。  
4. **HTTPS**：生产环境必须，否则 token 可能被中间人偷走。  
5. **退出登录**：纯 JWT 无法立刻作废（服务端不存它）。常见做法是等过期，或维护黑名单 / 换 refresh token。教学阶段知道这一点即可。

---

## 8. 一句话总结

**JWT 是带签名的登录凭证：登录时签发，之后请求用 `Authorization: Bearer` 带上，中间件验签后放行。**  
先把「签发 → 存储 → 带头发请求 → 中间件校验」这条链路跑通，再学刷新令牌和权限角色。
