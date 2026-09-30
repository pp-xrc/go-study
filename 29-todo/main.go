/*
第 29 课：完整小项目实战 —— ToDo List

两个模块：User（登录拿 JWT）+ Todo（全部接口要登录，且只能操作自己的待办）

表 todo：id / uid(外键→user.id) / title / content / status / create_time / update_time

接口（均需 Authorization: Bearer <token>）：
  POST   /todo         新建
  GET    /todo/list    分页列表  ?page=1&size=10
  GET    /todo/:id     详情
  PUT    /todo/:id     更新标题/内容/状态
  DELETE /todo/:id     删除

运行（开发）：
  go run ./29-todo

编译部署（在仓库根目录）：
  终端 1 API：
    mkdir -p 29-todo/bin
    go build -o 29-todo/bin/todo ./29-todo
    ./29-todo/bin/todo
  终端 2 前端（Nginx 托管 web/，并反代 /user /todo）：
    nginx -p "$PWD/29-todo" -c deploy/nginx.local.conf
  浏览器打开 http://localhost:8090
  停止 Nginx：nginx -p "$PWD/29-todo" -c deploy/nginx.local.conf -s stop
  分享给朋友（需已安装 cloudflared）：
    cloudflared tunnel --url http://localhost:8090

测试：
  # 1) 登录（内置账号）
  curl -s -X POST http://localhost:8089/user/login \
    -H "Content-Type: application/json" \
    -d '{"phone":"13800000000","password":"admin123"}'

  TOKEN=粘贴 token

  # 2) 新建
  curl -s -X POST http://localhost:8089/todo \
    -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d '{"title":"学 Go","content":"把 Day29 跑通","status":0}'

  # 3) 列表 / 详情 / 更新 / 删除
  curl -s "http://localhost:8089/todo/list?page=1&size=10" -H "Authorization: Bearer $TOKEN"
  curl -s http://localhost:8089/todo/1 -H "Authorization: Bearer $TOKEN"
  curl -s -X PUT http://localhost:8089/todo/1 \
    -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d '{"status":1}'
  curl -s -X DELETE http://localhost:8089/todo/1 -H "Authorization: Bearer $TOKEN"
*/

package main

import (
	"fmt"
	"log"

	"go-study/29-todo/model"
	"go-study/29-todo/router"
)

func main() {
	if err := model.Init(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	if err := model.SeedDemoUser(); err != nil {
		log.Fatalf("写入演示用户失败: %v", err)
	}

	r := router.Setup()
	addr := ":8089"
	fmt.Println("Day29 API:  http://localhost" + addr)
	fmt.Println("前端（同源反代）: http://localhost:8090")
	fmt.Println("对外分享: cloudflared tunnel --url http://localhost:8090")
	fmt.Println("演示账号: 13800000000 / admin123")
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
