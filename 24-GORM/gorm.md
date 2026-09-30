# GORM 入门速通（教学版）

> 衔接上一课：你已经会用 `database/sql` 手写 SQL。  
> 本课学 **GORM**：用 Go 结构体操作数据库，少写 SQL，够做基本 CRUD。

官方文档：https://gorm.io/zh_CN/docs/

---

## 1. GORM 是什么？

**GORM** 是 Go 里最常用的 **ORM**（Object-Relational Mapping，对象关系映射）框架。

一句话理解：

| 以前（database/sql） | 现在（GORM） |
|---------------------|--------------|
| 自己写 `SELECT ...` | 写 Go 代码：`db.Find(&users)` |
| `Scan` 一行行填字段 | 自动填进结构体 |
| 表结构和代码分开维护 | 结构体 ≈ 表结构 |

**不是替代 MySQL**，而是帮你少写样板 SQL。底层还是连 MySQL，SQL 基础仍然要会。

---

## 2. 和 database/sql 的对比

查用户列表：

**database/sql（上节课）：**

```go
rows, _ := db.Query("SELECT id, name, age FROM user")
for rows.Next() {
    var u User
    rows.Scan(&u.ID, &u.Name, &u.Age)
    // ...
}
```

**GORM：**

```go
var users []User
db.Find(&users) // 自动 SELECT，并填进切片
```

增删改查都能用类似「方法调用」的方式完成。

---

## 3. 安装

在项目根目录（有 `go.mod` 的地方）：

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

- `gorm.io/gorm`：GORM 本体  
- `gorm.io/driver/mysql`：MySQL 驱动适配（内部仍会用到 MySQL 驱动）

若下载超时，先设代理：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

---

## 4. 连接数据库

DSN 和上节课一样：

```text
root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local
```

```go
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	fmt.Println("✅ GORM 连接成功")
	_ = db
}
```

要点：

- `gorm.Open` 成功 ≈ 能连上库  
- 有密码：`root:密码@tcp(127.0.0.1:3306)/go_db?...`

---

## 5. 模型（Model）= 表结构

GORM 用 **结构体** 描述一张表。字段名默认转成蛇形列名：

| Go 字段 | 默认列名 |
|---------|----------|
| `Name` | `name` |
| `CreateTime` | `create_time` |
| `ID` | `id`（常当作主键） |

对应本课已有的 `user` 表：

```go
type User struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"size:64;not null"`
	Age        uint      `gorm:"not null;default:0"`
	Phone      string    `gorm:"size:20;not null;default:''"`
	Password   string    `gorm:"size:255;not null"`
	CreateTime time.Time `gorm:"autoCreateTime"`
}

// 表名默认是结构体名的复数 users；我们表叫 user，需要指定：
func (User) TableName() string {
	return "user"
}
```

常用 tag：

| tag | 含义 |
|-----|------|
| `primaryKey` | 主键 |
| `size:64` | VARCHAR(64) |
| `not null` | 非空 |
| `default:0` | 默认值 |
| `autoCreateTime` | 创建时自动填时间 |
| `column:xxx` | 自定义列名 |
| `-` | 忽略该字段 |

### AutoMigrate（可选）

让 GORM 根据结构体**自动建表/补字段**（适合学习、小项目）：

```go
db.AutoMigrate(&User{})
```

注意：

- 一般**不会删列**，比较安全  
- 正式项目复杂迁移常用专门工具；本课有现成 `user` 表时，能连上即可，不强制 AutoMigrate

---

## 6. CRUD 四件套

下面都假设已有 `db *gorm.DB` 和 `User` 模型。

### 6.1 增 Create

```go
u := User{
	Name:     "张三",
	Age:      20,
	Phone:    "13800001111",
	Password: "123456",
}
result := db.Create(&u)
// result.Error 看有没有错
// u.ID 会被自动填上自增主键
```

对应 SQL 大致：

```sql
INSERT INTO user (name, age, phone, password, create_time) VALUES (...);
```

### 6.2 查 Query

```go
// 查全部
var users []User
db.Find(&users)

// 按主键查一条
var u User
db.First(&u, 1) // id = 1

// 条件查询
db.Where("age >= ?", 25).Find(&users)
db.Where("name = ?", "张三").First(&u)

// 排序 + 限制条数
db.Order("id desc").Limit(10).Find(&users)
```

| GORM | 近似 SQL |
|------|----------|
| `Find` | `SELECT * FROM user` |
| `First` | `SELECT ... LIMIT 1`（找不到会报 `ErrRecordNotFound`） |
| `Where` | `WHERE ...` |
| `Order` | `ORDER BY ...` |
| `Limit` | `LIMIT ...` |

`?` 占位符仍然要用，防止 SQL 注入。

### 6.3 改 Update

```go
// 更新单个字段
db.Model(&User{}).Where("id = ?", 1).Update("age", 21)

// 更新多个字段
db.Model(&User{}).Where("id = ?", 1).Updates(User{
	Age:   22,
	Phone: "13900001111",
})
// 注意：Updates 用结构体时，默认**不会**更新零值字段（如 age=0）
// 若要更新零值，用 map：
db.Model(&User{}).Where("id = ?", 1).Updates(map[string]interface{}{
	"age": 0,
})
```

### 6.4 删 Delete

```go
// 按主键删
db.Delete(&User{}, 1)

// 按条件删
db.Where("name = ?", "李四").Delete(&User{})
```

若模型里加了 `gorm.DeletedAt`，`Delete` 会变成**软删除**（只标记删除，不真删行）。本课 `user` 表没有该字段，就是真删除。

---

## 7. 错误处理（必养成习惯）

```go
result := db.Create(&u)
if result.Error != nil {
	log.Println("创建失败:", result.Error)
	return
}
fmt.Println("影响行数:", result.RowsAffected)

err := db.First(&u, 999).Error
if errors.Is(err, gorm.ErrRecordNotFound) {
	fmt.Println("没有这条记录")
} else if err != nil {
	log.Println("查询失败:", err)
}
```

不要忽略 `Error`，否则出问题不好排查。

---

## 8. 调试：看 GORM 实际执行了什么 SQL

学习阶段强烈建议打开日志：

```go
import "gorm.io/gorm/logger"

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	Logger: logger.Default.LogMode(logger.Info),
})
```

控制台会打印类似：

```text
INSERT INTO `user` ...
SELECT * FROM `user` WHERE age >= 25
```

方便对照自己写的方法和真实 SQL。

---

## 9. 完整最小示例（串起来）

```go
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"size:64;not null"`
	Age        uint      `gorm:"not null;default:0"`
	Phone      string    `gorm:"size:20"`
	Password   string    `gorm:"size:255;not null"`
	CreateTime time.Time `gorm:"autoCreateTime"`
}

func (User) TableName() string { return "user" }

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 增
	u := User{Name: "赵六", Age: 28, Phone: "13800004444", Password: "abc"}
	db.Create(&u)
	fmt.Println("新建 ID:", u.ID)

	// 查
	var list []User
	db.Where("age >= ?", 20).Order("id").Find(&list)
	for _, item := range list {
		fmt.Printf("%d %s %d\n", item.ID, item.Name, item.Age)
	}

	// 改
	db.Model(&User{}).Where("id = ?", u.ID).Update("age", 29)

	// 删（演示后可注释掉，避免误删）
	// db.Delete(&User{}, u.ID)
}
```

运行（有 `main.go` 后）：

```bash
go run ./24-GORM
```

---

## 10. 学习路线（建议顺序）

1. **会连接**：`gorm.Open` + 打印成功  
2. **会建模型**：结构体 + `TableName` + 常用 tag  
3. **会 CRUD**：`Create` / `Find`/`First`/`Where` / `Update`/`Updates` / `Delete`  
4. **会看 SQL 日志**：对照 GORM 和原生 SQL  
5. **再进阶**（下一批再学即可）：  
   - 关联（一对多、多对一）  
   - 事务 `db.Transaction`  
   - 软删除 `DeletedAt`  
   - 钩子（创建前/更新前）

---

## 11. 和上节课对照速查

| 你想做的事 | SQL | database/sql | GORM |
|-----------|-----|--------------|------|
| 连接 | — | `sql.Open` + `Ping` | `gorm.Open` |
| 插入 | `INSERT` | `db.Exec` | `db.Create` |
| 查多条 | `SELECT` | `Query` + `Scan` | `db.Find` |
| 查一条 | `SELECT ... LIMIT 1` | `QueryRow` | `db.First` |
| 条件 | `WHERE` | SQL 字符串 | `db.Where` |
| 更新 | `UPDATE` | `db.Exec` | `db.Update` / `Updates` |
| 删除 | `DELETE` | `db.Exec` | `db.Delete` |

---

## 12. 常见坑

1. **表名不对**：结构体 `User` 默认表名是 `users`，本课表是 `user`，务必写 `TableName()`。  
2. **`Updates` 不更新 0 / 空字符串**：结构体零值会被忽略，改用 `map[string]interface{}`。  
3. **忘了 `parseTime=True`**：时间字段扫不进 `time.Time`。  
4. **忽略 `Error`**：查不到、唯一键冲突都看不出来。  
5. **以为 GORM 不用学 SQL**：复杂查询、排查性能，还是要会看 SQL。

---

## 13. 一句话总结

**GORM = 用 Go 结构体和链式方法，代替大部分手写 SQL 的增删改查。**  
先把「连接 → 模型 → CRUD → 看日志」跑通，再学关联和事务。
