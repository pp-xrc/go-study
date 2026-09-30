# Go + Gin + MySQL 知识体系

> 面向本仓库 `go-study`（尤其 16–29 课）。  
> 目标：把「语法课里的零件」组装成一张 **Web 后端地图**，复习时先看图，再下钻到对应目录。

按课号扫目录请用：[README.md](./README.md) · [学习总纲.md](./学习总纲.md)  
专题笔记：[23-MySQL/mysql.md](./23-MySQL/mysql.md) · [24-GORM/gorm.md](./24-GORM/gorm.md) · [26-JWT/jwt.md](./26-JWT/jwt.md) · [29-todo/README.md](./29-todo/README.md)

---

## 1. 先建立一张总图

一次完整请求（以第 29 课 ToDo 为例）：

```text
浏览器
  │  GET /app.html、POST /user/login、GET /todo/list?...
  ▼
Nginx :8090
  ├─ 静态页 HTML/CSS/JS  →  web/
  └─ /user /todo         →  反代到 127.0.0.1:8089
  ▼
Gin :8089
  1. 中间件 CORS
  2. 路由匹配（Group + GET/POST/PUT/DELETE）
  3. JWT 中间件（需登录的接口）→ 解析 token → c.Set("user_id")
  4. Controller 绑参数（JSON / Query / Param）
  5. Model / GORM 访问 MySQL（go_db）
  6. 统一 JSON：{ code, message, data }
  ▼
MySQL
  表 user、todo（todo.uid = 当前登录用户）
```

记住三件事：

1. **Gin 不存数据**，只处理 HTTP。
2. **GORM 不替代 SQL**，只是用结构体少写样板。
3. **JWT 不加密 Payload**，只签名；密码用 bcrypt 哈希存在库里。

---

## 2. 知识分层（从上到下）

```text
┌─────────────────────────────────────────┐
│  工程形态：分层、错误码、前后端分离、部署   │  28–29
├─────────────────────────────────────────┤
│  鉴权：bcrypt 存密码 + JWT 认人           │  26–27
├─────────────────────────────────────────┤
│  数据：SQL ↔ GORM CRUD / 分页 / 关联      │  23–25
├─────────────────────────────────────────┤
│  Web：路由、绑定、中间件、统一响应         │  16–22
├─────────────────────────────────────────┤
│  语言：类型、结构体、error、包、strconv    │  01–15
└─────────────────────────────────────────┘
```

越往上越「能做出一个接口」；底层不会，上面全是黑盒。

---

## 3. Go 后端必备（够写接口即可）

Web 课反复出现的语言点，不必一次背完语法全书。

| 点 | 为什么后端天天用 | 课 |
|----|------------------|----|
| `(值, error)`，先判断 `err != nil` | 查库、绑 JSON、签发 JWT 都会失败 | 03、10 |
| 结构体 + **指针接收者** | 模型、改字段、GORM 填结构体 | 09 |
| 结构体 tag | `json:"phone"` 绑请求体；`gorm:"column:phone"` 对表 | 18、24 |
| `package` / 大写导出 | `controller` 调 `model.GetUserByPhone` | 12、28 |
| 切片 / `make` | 列表接口 `[]Todo` | 05 |
| `map[string]any` | GORM `Updates` 才能改零值字段 | 06、24 |
| comma-ok | `v, ok := m[k]`；类型断言；JWT 算法检查 | 06、11、26 |
| `strconv.Atoi` | Query 里的 `page`、`size` 是字符串 | 15、29 |
| `time` 模板 `2006-01-02 15:04:05` | 时间格式化 | 15 |
| `sync.RWMutex` | 第 22 课内存仓库；有了 MySQL 后一般不再用 | 14、22 |

**零值提醒：** `0`、`""`、`false`、`nil` 是类型的默认值。GORM 用结构体 `Updates` **不会更新零值**，要改 `status: 0` 得传 `map`。

---

## 4. HTTP 与 REST（Gin 的工作语言）

接口 = **方法 + 路径 + 请求体/查询串 + 状态码 + JSON**。

| 动词 | 习惯用途 | 本仓库例子 |
|------|----------|------------|
| GET | 读，可重复，无 body | `GET /todo/list`、`GET /todo/:id` |
| POST | 新建 / 登录 | `POST /user/login`、`POST /todo` |
| PUT | 更新整条或若干字段 | `PUT /todo/:id` |
| DELETE | 删除 | `DELETE /todo/:id` |
| OPTIONS | CORS 预检，浏览器自动发 | `middleware/cors.go` 回 **204** |

参数从哪来：

| 位置 | Gin | 例子 |
|------|-----|------|
| 路径 | `c.Param("id")` | `/todo/3` |
| 查询串 | `c.Query` / `c.DefaultQuery` | `?page=1&size=10` |
| JSON body | `c.ShouldBindJSON` | `{ "phone", "password" }` |
| 请求头 | `c.GetHeader` | `Authorization: Bearer ...` |

**路由顺序：** 静态路径写在通配前面。`/todo/list` 必须在 `/todo/:id` 之前，否则 `list` 会被当成 id。

统一响应（本仓库约定）：

```json
{ "code": 0, "message": "ok", "data": { } }
```

- HTTP 状态码：给浏览器 / Nginx 看（200、401、404…）
- `code`：给前端业务看（`0` 成功，`10005` 未登录）

常见 HTTP 状态码见第 8 节。

---

## 5. Gin 核心

### 5.1 最小骨架

```go
r := gin.Default()          // 自带 Logger、Recovery
r.GET("/ping", handler)
r.Run(":8089")              // 监听所有网卡的 8089
```

`gin.Default()` vs `gin.New()`：后者是空引擎，要自己挂中间件。

### 5.2 请求生命周期

```text
请求进入
  → 全局中间件（CORS）
  → 路由组中间件（JWTAuth）
  → Handler
       Bind → 调 model → 写 JSON
  → 中间件从内往外收尾（洋葱模型）
```

| API | 作用 |
|-----|------|
| `c.Next()` | 放行到下一层 |
| `c.Abort()` / `AbortWithStatus` | 拦住，后面 handler 不跑 |
| `c.Set` / `c.Get` | 中间件把 `user_id` 传给 handler |
| `c.JSON(status, obj)` | 写响应 |

### 5.3 中间件该放什么

| 类型 | 做什么 | 本仓库 |
|------|--------|--------|
| CORS | 跨域时告诉浏览器谁能访问 | `29-todo/middleware/cors.go` |
| JWT | 验 token，失败直接 401 | `middleware/auth.go` |
| 日志 | 访问记录 | `gin.Default()` 自带 |
| Recovery | panic 转 500，进程不挂 | `gin.Default()` 自带 |

**CORS 预检：** 跨域 + JSON + `Authorization` 时，浏览器先 OPTIONS。服务端要先写 `Access-Control-Allow-*` 头，再对 OPTIONS 返回 204 并 `Abort`。同源（页面和接口都走 Nginx 8090）时预检通常不会发。

### 5.4 分组

```go
user := r.Group("/user")
todo := r.Group("/todo")
todo.Use(middleware.JWTAuth())  // 整组都要登录
```

登录、注册放在 JWT 外面；`/user/me`、全部 `/todo` 套 JWT。

---

## 6. MySQL 核心

库像文件夹，表像 Excel：行 = 记录，列 = 字段。本课库名 **`go_db`**。

### 6.1 四句 SQL

```sql
INSERT INTO user (phone, password, nickname) VALUES (?, ?, ?);
SELECT id, phone, nickname FROM user WHERE phone = ?;
UPDATE todo SET status = 1 WHERE id = ? AND uid = ?;
DELETE FROM todo WHERE id = ? AND uid = ?;
```

占位符必须用 `?`，**禁止**把用户输入拼进 SQL（防注入）。

### 6.2 表设计要点

| 概念 | 含义 |
|------|------|
| PRIMARY KEY + AUTO_INCREMENT | `id` 自增主键 |
| NOT NULL / DEFAULT | 约束 |
| utf8mb4 | 中文和 emoji |
| DATETIME | 配合 DSN 的 `parseTime=True` 扫进 `time.Time` |
| 外键思路 | `todo.uid` 指向 `user.id`（本课迁移时关掉了自动建 FK，避免旧表类型冲突） |

第 29 课两张表：

```text
user  1 ──────<  todo
 id               id
 phone            uid   ← 只查 / 只改「当前用户」的行
 password(哈希)    title, content, status
```

### 6.3 从 Go 连接

DSN：

```text
用户:密码@tcp(主机:3306)/库名?charset=utf8mb4&parseTime=True&loc=Local
```

本机示例：`root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local`

| API（`database/sql`） | 用途 |
|----------------------|------|
| `sql.Open` | 只解析 DSN，**还不保证连得上** |
| `Ping` | 真正握手，必须做 |
| `Query` + `Scan` | 查多行 |
| `QueryRow` | 查一行 |
| `Exec` | 增删改 |

第 23 课用手写 SQL；第 24 课起交给 GORM，但 `Ping` 仍然要做（见 `29-todo/model/db.go`）。

---

## 7. GORM 核心

ORM = 用结构体描述表，用方法代替手写 SQL。

### 7.1 模型

```go
type Todo struct {
    ID  uint   `gorm:"column:id;primaryKey;autoIncrement"`
    UID uint   `gorm:"column:uid;index;not null"`
    // ...
}
func (Todo) TableName() string { return "todo" }  // 否则默认 todos
```

- 默认表名是结构体名的**蛇形复数**，本课表是 `user` / `todo`，必须 `TableName()`。
- `AutoMigrate` 按模型补列，**不会**随意删列。

### 7.2 CRUD 对照

| 意图 | SQL | GORM |
|------|-----|------|
| 插入 | INSERT | `db.Create(&u)` |
| 多条 | SELECT | `db.Find(&list)` |
| 一条 | SELECT … LIMIT 1 | `db.First(&u, id)` / `Where(...).First` |
| 条件 | WHERE | `Where("uid = ?", uid)`，多个 Where = AND |
| 或 | OR | `Or("x = ?", v)` |
| 更新 | UPDATE | `Update` / `Updates`；结构体跳过零值，用 **map** |
| 删除 | DELETE | `Delete` |
| 分页 | LIMIT / OFFSET | `Limit(size).Offset((page-1)*size)` |
| 总数 | COUNT | `Model(&Todo{}).Count(&total)` |
| 排序 | ORDER BY | `Order("id desc")` |

找不到行：`errors.Is(err, gorm.ErrRecordNotFound)`（封装成 `model.IsNotFound`）。

### 7.3 本课最重要的查询模式

待办必须带 **uid**，防止越权：

```go
db.Where("id = ? AND uid = ?", id, uid).First(&t)
```

别人的 id 当成 404，不告诉前端「这条存在但不是你的」。

---

## 8. 错误处理：HTTP 状态码 × 业务码

见 `29-todo/utils/errors.go` + `response.go`。

| HTTP | 含义 | 业务码例子 |
|------|------|------------|
| 200 | 成功 | `code = 0` |
| 201 | 新建成功 | `code = 0` |
| 204 | 成功无 body | CORS OPTIONS |
| 400 | 参数不合法 | `400` / `CodeBadRequest` |
| 401 | 未认证（没登录 / 密码错） | `10001` 用户不存在、`10002` 密码错、`10005` 未登录 |
| 403 | 已登录但没权限 | `10004`（第 27 课 admin 列表） |
| 404 | 资源不存在 | `10006` 待办不存在 |
| 409 | 冲突 | `10003` 手机号已注册 |
| 500 | 服务器/数据库异常 | `500` |

口诀：**2 成功，4 你错了，5 我错了；401 是谁，403 能不能，404 有没有，409 撞车。**

---

## 9. 认证与鉴权

两个问题分开：

1. **认证 Authentication**：你是谁？（登录 + JWT）
2. **鉴权 Authorization**：你能干什么？（只能改自己的 todo；admin 才能看全量用户）

### 9.1 密码：bcrypt

- 注册：`HashPassword(明文)` → 把哈希写入 `user.password`
- 登录：`CompareHashAndPassword(哈希, 明文)`
- **不能解密**，只能比对
- 前端 Network 里仍能看到明文，那是 TLS 解开后的请求体；路上靠 HTTPS 保护

### 9.2 JWT

```text
Header.Payload.Signature
```

Payload 可被 Base64 解码（别放密码）。密钥只在服务端，用来验签。

登录链路：

```text
POST /user/login { phone, password }
  → 查 user → bcrypt 比对
  → GenerateToken(userID, role)
  → 前端 localStorage 存 token
  → 之后：Authorization: Bearer <token>
  → JWTAuth 中间件 VerifyToken
  → c.Set("user_id", …)
  → Controller 用这个 id 去查 / 改数据
```

token 过期或缺失 → HTTP 401 + `code = 10005` → 前端清 token 跳登录页。

---

## 10. 工程分层（第 28–29 课）

单文件 `main.go` 适合学语法；上项目后按职责拆包：

```text
main.go          启动：连库、Listen
router/          只注册路由和中间件，不写业务
middleware/      CORS、JWT
controller/      绑参数、调 model、写 JSON（不直接拼 SQL）
model/           结构体 + GORM + 表名
utils/           密码、JWT、错误码、统一响应
web/             静态前端（Nginx 托管，不属于 Go 包）
```

依赖方向：**router → controller → model**，不要反向。`controller` 不 `import` 数据库驱动。

---

## 11. 前后端分离与部署（第 29 课）

| 角色 | 端口 | 职责 |
|------|------|------|
| Gin | 8089 | 只提供 JSON API |
| Nginx | 8090 | 静态页 + 把 `/user` `/todo` 反代到 8089 |
| MySQL | 3306 | 只本机，不要映射到公网 |
| cloudflared | — | 把 8090 打成一条 https 临时域名 |

`api.js` 里 `API_BASE = ""`：浏览器只请求**当前网站**同源路径，朋友才不会打到他自己的 `localhost`。

临时分享：

```bash
cloudflared tunnel --url http://localhost:8090
```

不需要 `~/.cloudflared/config.yml`。关掉进程，域名立刻失效。

---

## 12. 一张对照表（默写用）

| 我想… | 用什么 |
|--------|--------|
| 起一个 HTTP 服务 | Gin `Run` |
| 读 `?page=1` | `DefaultQuery` + `strconv.Atoi` |
| 读 JSON body | `ShouldBindJSON` + 结构体 tag |
| 保护接口 | `Group.Use(JWTAuth)` |
| 知道当前用户 | 中间件 `c.Set`，handler 再 `Get` |
| 连 MySQL | DSN + `Ping` |
| 少写 SQL | GORM `Create/First/Find/Updates/Delete` |
| 分页 | `Count` + `Limit` + `Offset` |
| 密码入库 | bcrypt Hash，不要明文、不要可逆加密当密码存储 |
| 登录态 | JWT Bearer，服务端不存 session |
| 统一出错 | `AppError{HTTP, Code, Message}` |
| 给外部访问 | Nginx 同源反代 + cloudflared |

---

## 13. 和课程目录的映射

| 体系章节 | 目录 |
|----------|------|
| Go 必备语法 | `01-hello` … `12-package-module`（指针 `07b-pointers`）、`15-standard-library`（含 json） |
| Gin 路由绑定中间件 | `16-gin-intro` … `21-gin-upload` |
| 内存版用户模块（无 DB） | `22-gin-advance` |
| 手写 SQL | `23-MySQL` |
| GORM CRUD | `24-GORM` |
| 条件 / 分页 / 排序 | `25-GORM-advance` |
| JWT 原理 | `26-JWT` |
| 三件套合在一个 main | `27-Gin-GORM-JWT` |
| 模块化 + 错误码 | `28-modular` |
| 关联数据 + 前后端分离 | `29-todo` |

建议复习路径：`22`（REST 手感）→ `23/24`（数据）→ `26/27`（登录）→ `28/29`（分层和越权）。

---

## 14. 过关标准

能不看文档讲清并动手复现下面这条链，就算这套知识体系过关：

1. 注册时密码 bcrypt 入库。  
2. 登录校验后发 JWT。  
3. 带 `Authorization` 访问 `/todo`。  
4. 列表/详情/改/删都带 `uid`，别人的数据 404。  
5. 响应永远是 `{code, message, data}`。  
6. 浏览器只打 8090，API 由 Nginx 转到 8089。

卡住时：打开对应目录 `main.go` 顶部的 curl，先让接口通，再回头看这一节在总图的哪一层。
