/*
第 28 课：项目模块化重构 + 错误处理优化

目录（本课新增，不改仓库里其它章节）：

  28-modular/
  ├── main.go              启动入口：连库、挂路由
  ├── router/              只负责注册路由，不写业务、不碰数据库
  ├── controller/          业务逻辑：调 model，不直接操作 DB
  ├── model/               GORM 模型 + 数据库方法
  ├── middleware/          JWT、管理员权限
  └── utils/               JWT、bcrypt、自定义错误码、统一响应

分层约定：
  请求 → router → middleware → controller → model → MySQL

运行：
  go run ./28-modular

测试：
  # 注册
  curl -s -X POST http://localhost:8088/user/register \
    -H "Content-Type: application/json" \
    -d '{"phone":"13800001111","password":"123456","nickname":"张三"}'

  # 登录（手机号）
  curl -s -X POST http://localhost:8088/user/login \
    -H "Content-Type: application/json" \
    -d '{"phone":"13800000000","password":"admin123"}'

  TOKEN=上一步返回的 token
  curl -s http://localhost:8088/user/list -H "Authorization: Bearer $TOKEN"

内置管理员：phone=13800000000  password=admin123

错误码：
  0      成功
  10001  用户不存在
  10002  密码错误
  10003  手机号已注册
  10004  权限不足
  10005  未登录 / token 无效
  400    参数错误
  500    数据库或系统错误
*/

package main

import (
	"fmt"
	"log"

	"go-study/28-modular/model"
	"go-study/28-modular/router"
)

func main() {
	if err := model.Init(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	if err := model.SeedAdmin(); err != nil {
		log.Fatalf("写入管理员失败: %v", err)
	}

	r := router.Setup()
	addr := ":8088"
	fmt.Println("Day28 模块化 demo: http://localhost" + addr)
	fmt.Println("管理员: 13800000000 / admin123")
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
