# Day29 ToDo List（前后端分离）

用户模块 + 待办 CRUD。**只能操作自己的待办**（`todo.uid` = JWT 里的 `user_id`）。

浏览器只访问 **8090**。Nginx 托管静态页，并把 `/user`、`/todo` 反代到本机 Gin。

```text
浏览器  →  http://localhost:8090   Nginx
              ├─ /、HTML/CSS/JS     →  web/
              └─ /user /todo        →  127.0.0.1:8089  Gin
```

Gin 仍监听 8089，但页面不再直接跨域请求它（`api.js` 里 `API_BASE` 为空）。

---

## 前置条件

- Go 1.26+
- 本机 MySQL 已启动，库名 `go_db`
- Nginx（`brew install nginx`）
- 对外分享：Cloudflare Tunnel（`brew install cloudflared`）
- 命令在**仓库根目录** `go-study/` 执行

```bash
mysql -u root -h 127.0.0.1 -e "CREATE DATABASE IF NOT EXISTS go_db DEFAULT CHARACTER SET utf8mb4;"
```

---

## 怎么跑

**终端 1：后端 API**

```bash
mkdir -p 29-todo/bin
go build -o 29-todo/bin/todo ./29-todo
./29-todo/bin/todo
```

开发也可用：`go run ./29-todo`  
直连 API（curl 调试）：http://localhost:8089

**终端 2：前端 Nginx**

```bash
mkdir -p 29-todo/logs 29-todo/tmp
nginx -p "$PWD/29-todo" -c deploy/nginx.local.conf
```

本机打开：**http://localhost:8090**

改完配置后：

```bash
nginx -p "$PWD/29-todo" -c deploy/nginx.local.conf -s reload
```

停止 Nginx：

```bash
nginx -p "$PWD/29-todo" -c deploy/nginx.local.conf -s stop
```

演示账号：`13800000000` / `admin123`

---

## 分享给外部（Cloudflare Tunnel）

本机 8090 / 8089 仍只监听在你电脑上。`cloudflared` 打一条到 Cloudflare 的隧道，给你一个 **https 公网地址**，朋友打开即可。不需要改路由器、也不用注册也能试（临时域名）。

先保证终端 1、2 已经在跑，然后：

```bash
brew install cloudflared   # 装过可跳过
cloudflared tunnel --url http://localhost:8090
```

终端里会出现类似：

```text
https://xxxx-xxxx.trycloudflare.com
```

把这个 URL 发给朋友。关掉该终端或 Ctrl+C，公网地址立刻失效。

注意：

- 临时隧道每次启动域名都不同；电脑休眠/断网朋友会打不开
- 演示账号会一起暴露，测完停掉 tunnel 和 API
- **不要**把 MySQL 3306 暴露出去
- 这是教学快速通道。固定域名需要 Cloudflare 账号 + Named Tunnel（`cloudflared tunnel login`）

---

## 目录

```text
29-todo/
├── main.go              API 入口，监听 8089
├── router/              只注册 /user /todo
├── controller/          业务（不直接碰 DB）
├── model/               User / Todo + GORM
├── middleware/          JWT、CORS
├── utils/               密码、JWT、错误码、统一响应
├── web/                 前端静态文件（由 Nginx 托管）
├── deploy/nginx.local.conf
└── bin/                 编译产物（不入库）
```

---

## 数据表

`todo`（启动时 AutoMigrate）：id / uid / title / content / status(0 未完成 1 已完成) / create_time / update_time

---

## 接口

统一响应：`{ "code": 0, "message": "...", "data": ... }`  
登录后请求头：`Authorization: Bearer <token>`  
页面请求前缀：当前站点（经 Nginx）；curl 可直连 `http://localhost:8089`

| 方法 | 路径 | 登录 | 说明 |
|------|------|------|------|
| POST | `/user/register` | 否 | `{phone, password, nickname}` |
| POST | `/user/login` | 否 | `{phone, password}` → `data.token` |
| GET | `/user/me` | 是 | 当前用户 |
| POST | `/todo` | 是 | 新建，uid 从 token 取 |
| GET | `/todo/list?page=1&size=10` | 是 | 分页，只返回自己的 |
| GET | `/todo/:id` | 是 | 详情 |
| PUT | `/todo/:id` | 是 | 更新 |
| DELETE | `/todo/:id` | 是 | 删除 |

错误码：0 成功；10001 用户不存在；10002 密码错误；10003 手机号已注册；10005 未登录；10006 待办不存在。

---

## 前端说明

- 页面和接口都走 8090（见 `web/assets/api.js` 的 `API_BASE = ""`）
- token 存在 `localStorage`（`todo_token`）
- `code === 10005` 清 token 并跳回登录页 `/`
