# go-study — 前端转 Go 全栈学习笔记

按章节拆分的 Go 学习仓库。每个编号目录可独立 `go run`，从语法到 Gin、MySQL、GORM、JWT，最后串成 ToDo 项目。

更适合复习的总纲见：[学习总纲.md](./学习总纲.md)  
Go + Gin + MySQL 知识体系（一张图串起来）：[Go-Gin-MySQL知识体系.md](./Go-Gin-MySQL知识体系.md)

---

## 环境要求

- Go 1.26.1（见 `go.mod`），检查：`go version`
- 第 23 课起需要本机 MySQL（库名 `go_db`）
- 国内拉依赖建议：`go env -w GOPROXY=https://goproxy.cn,direct`

## 如何运行

在**项目根目录**执行，每次只跑一章：

```bash
go run ./01-hello
go run ./16-gin-intro
go run ./27-Gin-GORM-JWT
```

Web 课启动后用 curl / Postman 测，端口写在各章 `main.go` 头部注释里。

---

## 学习路线

```
语法基础 01–12（指针在 07b-pointers）
    → 并发与标准库 13–15
    → Gin Web 16–22
    → MySQL / GORM 23–25
    → JWT 鉴权 26
    → 整合 / 分层 / ToDo 27–29
```

建议按编号学。后面章节会用到前面的概念，但各章代码互不影响，可单独反复跑。

---

## 各目录知识点

### 一、Go 语法基础

| 目录 | 知识点 |
|------|--------|
| [01-hello](./01-hello) | `package main`、`fmt` 打印、`const` / `var`、`if-else`、`Sprintf` |
| [02-basic-types](./02-basic-types) | 基本类型、零值、`:=` 短声明、`iota`、`%T` 看类型、`rune` / `[]byte` |
| [03-operators](./03-operators) | 算术运算、`switch`、多返回值 `(结果, error)`、`nil` |
| [04-loop](./04-loop) | `for` 三种写法、`range`、`for range n`、`break` / `continue` |
| [05-arr-slice](./05-arr-slice) | 数组 vs 切片、`make` / `append` / `len` / `cap`、切片截取 |
| [06-map-string](./06-map-string) | `map` 增删改查、`delete`、comma-ok、`strings` 常用函数 |
| [07-functions](./07-functions) | 值传递、可变参数 `...`、函数作参数/返回值、`defer` |
| [07b-pointers](./07b-pointers) | 指针声明、`&` 取址、`*` 解引用、值传参 vs 指针传参 |
| [08-scope](./08-scope) | 函数与变量作用域（包级 / 函数内 / `:=` 影子） |
| [09-struct](./09-struct) | 结构体定义、嵌套、方法、值接收者 vs 指针接收者 |
| [10-struct-error](./10-struct-error) | 结构体 + map + 切片；`error`；`panic` / `recover` |
| [11-interface](./11-interface) | 接口、多态、空接口 `any`、类型断言、`switch` 判类型 |
| [12-package-module](./12-package-module) | `package` / `import` / 导出规则（首字母大小写）、`go.mod` |

### 二、并发与标准库

| 目录 | 知识点 |
|------|--------|
| [13-goroutine-channel](./13-goroutine-channel) | `go` 启动协程、无缓冲/有缓冲 channel、Worker Pool |
| [14-sync-lock](./14-sync-lock) | `WaitGroup`、竞态、`Mutex`、`defer Unlock` |
| [15-standard-library](./15-standard-library) | `strconv`、`time`、`encoding/json`、`os` 文件读写 |

### 三、Gin Web

| 目录 | 知识点 |
|------|--------|
| [16-gin-intro](./16-gin-intro) | Gin 入门、`gin.Default()`、第一个 GET 接口 |
| [17-gin-router](./17-gin-router) | GET/POST/PUT/DELETE、路径参数、路由分组 |
| [18-gin-data-bind](./18-gin-data-bind) | `Param` / `Query`、JSON 绑定 `ShouldBindJSON`、参数校验 |
| [19-gin-response](./19-gin-response) | 统一返回 `{code,message,data}`，见 [response](./response) |
| [20-gin-middleware](./20-gin-middleware) | 中间件、洋葱模型、`c.Next()` / `c.Abort()`、日志 / CORS |
| [21-gin-upload](./21-gin-upload) | 单/多文件上传、保存到 [upload](./upload)、`Redirect` 301/302 |
| [22-gin-advance](./22-gin-advance) | 用户模块综合（内存 map）：注册登录列表详情更新 |

### 四、数据库与 ORM

| 目录 | 知识点 |
|------|--------|
| [23-MySQL](./23-MySQL) | `database/sql` + MySQL 驱动、DSN、`Ping`、`Query`/`Scan`/`Exec`；SQL 笔记 [mysql.md](./23-MySQL/mysql.md) |
| [24-GORM](./24-GORM) | GORM 连接、模型 tag、`TableName`、CRUD；笔记 [gorm.md](./24-GORM/gorm.md) |
| [25-GORM-advance](./25-GORM-advance) | `Where`/`Or`、`Limit`/`Offset` 分页、`Order`、`Count` |

本机库：`go_db`，表：`user`。DSN 示例：

```text
root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local
```

### 五、鉴权与整合

| 目录 | 知识点 |
|------|--------|
| [26-JWT](./26-JWT) | JWT 三段结构、签发/校验、`Authorization: Bearer`、鉴权中间件；笔记 [jwt.md](./26-JWT/jwt.md) |
| [27-Gin-GORM-JWT](./27-Gin-GORM-JWT) | Gin + GORM + JWT：bcrypt、注册登录、JWT 护路由、admin/user 权限 |
| [28-modular](./28-modular) | 模块化：router / controller / model / middleware / utils；自定义错误码 |
| [29-todo](./29-todo) | ToDo：user + todo 关联，JWT 下只能操作自己的待办；Nginx + cloudflared |

---

## 公共包

| 路径 | 作用 |
|------|------|
| [response](./response) | 统一 API 响应：`OK` / `BadRequest` / `NotFound` / `Fail` |
| [upload](./upload) | 第 21 课上传文件存放目录 |

整个仓库共用根目录 `go.mod`（模块名 `go-study`）。

---

## 复习建议

1. 语法课改完立刻 `go run ./xx` 看输出。  
2. Web 课对着 `main.go` 顶部 curl 示例测正常/异常。  
3. 数据库课先确认 `mysql -u root -h 127.0.0.1` 能连上。  
4. 总复习按 [学习总纲.md](./学习总纲.md) 的「必背清单」过一遍即可。  
5. 要把 Gin / MySQL / JWT 串成体系，读 [Go-Gin-MySQL知识体系.md](./Go-Gin-MySQL知识体系.md)。
