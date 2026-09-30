# MySQL 基础速通（教学版）

> 配合本课 `go_db` / `user` 表学习。只讲最常用的部分，够写 Go 增删改查即可。

---

## 1. 什么是 SQL？

**SQL**（Structured Query Language）是跟数据库说话的语言。

| 你想做的事 | 对应操作 | 常见语句 |
|-----------|---------|---------|
| 查数据 | 读 | `SELECT` |
| 加数据 | 增 | `INSERT` |
| 改数据 | 改 | `UPDATE` |
| 删数据 | 删 | `DELETE` |
| 建库/建表 | 结构 | `CREATE` |

进 MySQL 命令行：

```bash
mysql -u root -h 127.0.0.1
```

选中本课数据库：

```sql
USE go_db;
```

---

## 2. 库与表（先有架子，再放数据）

### 2.1 数据库（Database）

一个库像一个「文件夹」，里面放很多表。

```sql
-- 查看有哪些库
SHOW DATABASES;

-- 创建库（utf8mb4 支持中文和 emoji）
CREATE DATABASE go_db DEFAULT CHARACTER SET utf8mb4;

-- 使用某个库
USE go_db;

-- 删除库（危险！数据全没）
-- DROP DATABASE go_db;
```

### 2.2 表（Table）

表像 Excel：行 = 一条记录，列 = 字段。

本课的 `user` 表结构：

| 字段 | 含义 | 类型（简化理解） |
|------|------|------------------|
| id | 主键，自动编号 | 整数，自增 |
| name | 姓名 | 字符串 |
| age | 年龄 | 整数 |
| phone | 手机号 | 字符串 |
| password | 密码 | 字符串 |
| create_time | 创建时间 | 日期时间 |

```sql
CREATE TABLE user (
  id          INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(64)  NOT NULL,
  age         INT UNSIGNED NOT NULL DEFAULT 0,
  phone       VARCHAR(20)  NOT NULL DEFAULT '',
  password    VARCHAR(255) NOT NULL,
  create_time DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

几个关键字：

- `PRIMARY KEY`：主键，唯一标识一行（通常用 id）
- `AUTO_INCREMENT`：插入时 id 自动 +1
- `NOT NULL`：不能为空
- `DEFAULT`：不写时用默认值
- `VARCHAR(n)`：最长 n 个字符的字符串

常用查看命令：

```sql
SHOW TABLES;          -- 看有哪些表
DESCRIBE user;        -- 看表结构（可简写 DESC user）
SHOW CREATE TABLE user; -- 看建表完整 SQL
```

---

## 3. 增：INSERT

```sql
INSERT INTO user (name, age, phone, password)
VALUES ('张三', 20, '13800001111', '123456');
```

一次插多行：

```sql
INSERT INTO user (name, age, phone, password) VALUES
  ('李四', 25, '13800002222', '123456'),
  ('王五', 30, '13800003333', '123456');
```

说明：

- 没写的 `id`、`create_time` 会用默认值（自增 / 当前时间）
- 列名顺序和 `VALUES` 里的值顺序必须一一对应

---

## 4. 查：SELECT（最常用）

### 4.1 查全部

```sql
SELECT * FROM user;
```

`*` 表示所有列。学习阶段可以，正式项目更推荐写出具体列名。

### 4.2 查指定列

```sql
SELECT id, name, age, phone FROM user;
```

### 4.3 条件过滤：WHERE

```sql
-- 等于
SELECT * FROM user WHERE id = 1;

-- 比较
SELECT * FROM user WHERE age >= 25;

-- 模糊匹配：% 表示任意多个字符
SELECT * FROM user WHERE name LIKE '张%';

-- 多条件
SELECT * FROM user WHERE age >= 20 AND phone != '';
SELECT * FROM user WHERE age < 18 OR age > 60;
```

### 4.4 排序与条数

```sql
-- 按 id 升序（默认 ASC）；DESC 是降序
SELECT * FROM user ORDER BY id ASC;
SELECT * FROM user ORDER BY age DESC;

-- 只要前 2 条
SELECT * FROM user ORDER BY id LIMIT 2;

-- 分页：跳过前 2 条，再取 2 条（第 2 页，每页 2 条）
SELECT * FROM user ORDER BY id LIMIT 2 OFFSET 2;
```

### 4.5 计数

```sql
SELECT COUNT(*) FROM user;
SELECT COUNT(*) FROM user WHERE age >= 25;
```

---

## 5. 改：UPDATE

```sql
UPDATE user SET age = 21 WHERE id = 1;
```

改多个字段：

```sql
UPDATE user
SET nickname = '小张', age = 22
WHERE id = 1;
```

**非常重要：**

```sql
-- ❌ 没有 WHERE 会改整张表所有行！
UPDATE user SET age = 0;
```

写 `UPDATE` / `DELETE` 时，先想清楚 `WHERE` 条件。

---

## 6. 删：DELETE

```sql
-- 删某一行
DELETE FROM user WHERE id = 3;

-- ❌ 没有 WHERE 会删光整张表
-- DELETE FROM user;
```

和「删表」的区别：

| 语句 | 效果 |
|------|------|
| `DELETE FROM user WHERE ...` | 删数据，表还在 |
| `TRUNCATE TABLE user` | 清空整表数据，表结构还在 |
| `DROP TABLE user` | 表直接删掉 |

---

## 7. 一条完整小练习

```sql
USE go_db;

-- 1. 插入
INSERT INTO user (name, age, phone, password)
VALUES ('赵六', 28, '13800004444', 'abc123');

-- 2. 查询
SELECT id, name, age FROM user WHERE name = '赵六';

-- 3. 更新
UPDATE user SET age = 29 WHERE name = '赵六';

-- 4. 再查确认
SELECT * FROM user WHERE name = '赵六';

-- 5. 删除
DELETE FROM user WHERE name = '赵六';
```

---

## 8. 和 Go 代码的对应关系

本课 `main.go` 里常见写法：

| SQL | Go（database/sql） |
|-----|---------------------|
| 连接并确认 | `sql.Open` + `db.Ping()` |
| `INSERT` | `db.Exec("INSERT ...", 参数...)` |
| `SELECT` 多行 | `db.Query` → `rows.Next` → `rows.Scan` |
| `SELECT` 单行 | `db.QueryRow(...).Scan(...)` |
| `UPDATE` / `DELETE` | `db.Exec(...)` |

占位符 `?`（防 SQL 注入）：

```go
db.Exec(
  "INSERT INTO user (name, age, phone, password) VALUES (?, ?, ?, ?)",
  "张三", 20, "13800001111", "123456",
)
```

不要自己用字符串拼接用户输入拼 SQL。

---

## 9. 容易踩的坑

1. **忘了 `USE go_db`**：语句可能打到别的库，或报「表不存在」。
2. **`UPDATE`/`DELETE` 忘了 `WHERE`**：全表被改/删。
3. **字符串要用引号**：`'张三'` 对，`张三` 错。
4. **中文乱码**：建库用 `utf8mb4`；Go DSN 加 `charset=utf8mb4`。
5. **时间字段**：Go DSN 加 `parseTime=True`，才能把 `DATETIME` 扫进 `time.Time`。

推荐本课 DSN：

```text
root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local
```

---

## 10. 速查卡片

```sql
SHOW DATABASES;
USE go_db;
SHOW TABLES;
DESC user;

INSERT INTO user (name, age, phone, password) VALUES (...);
SELECT * FROM user WHERE ... ORDER BY id LIMIT 10;
UPDATE user SET age = 21 WHERE id = 1;
DELETE FROM user WHERE id = 1;
```

记这四句就够入门：**INSERT / SELECT / UPDATE / DELETE**。  
结构相关记两句：**CREATE DATABASE / CREATE TABLE**。
