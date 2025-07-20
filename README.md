# 个人博客系统后端
这是一个使用Go语言、Gin框架和GORM框架开发的个人博客系统后端，支持用户认证、文章管理和评论功能。
***
## 功能特性
***
- 用户注册与登录
- JWT用户认证与授权
- 文章管理功能的增删改查 (CRUD)
- 评论功能
- 错误处理与日志记录
- 统一的错误处理
- RESTful API 设计
## 技术栈
***
- 语言: Go 1.24
- Web框架: Gin
- ORM框架: GORM 
- 数据库: MySQL 9.2.0
- 认证: JWT
- 密码加密:bcrypt
## 项目结构
***
```markdown
go-blog/
├── base/            # 公共基础包
├── config/          # 配置管理
├── web/             # 控制器层
├── dao/             # 数据库连接层
├── middleware/      # 中间件
├── service/         # 业务层
├── router/          # 路由配置
├── env/             # 环境参数配置
├── go.mod           # 依赖管理
├── main.go          # 程序入口
└── README.md        # 项目说明
```
## 数据库表结构
***
### users表
  - id - 主键
  - username - 用户名（唯一）
  - password - 加密密码
  - email - 邮箱（唯一）
  - created_at - 创建时间
  - updated_at - 更新时间
  - deleted_at - 删除时间
### posts表
- id - 主键
- title - 文章标题
- content - 文章内容
- user_id - 作者ID（外键）
- created_at - 创建时间
- updated_at - 更新时间
- deleted_at - 删除时间
### comments表
- id - 主键
- content - 评论内容
- user_id - 评论者ID（外键）
- post_id - 文章ID（外键）
- created_at - 创建时间
- updated_at - 更新时间
- deleted_at - 删除时间
## 安装与运行
***
### 环境要求
- Go 1.20+
- Mysql 5.7+
## 安装步骤
``` 
# 克隆项目（如果从git获取）
git clone <repository-url>
cd go-blog
# 安装依赖
go mod tidy
```
## 数据库配置
1.创建 MySQL 数据库
``` sql
CREATE DATABASE `blog` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```
2.修改数据库连接配置（已在代码中配置）
```sql
Dsn = "root:root@admin123@tcp(127.0.0.1:3306)/blog?charset=utf8&parseTime=True&loc=Local"
```
## 配置说明
系统支持通过指定环境(dev/prod)进行不同配置：
- DSN: 数据库连接字符串
- JwtSecretKey: JWT 签名密钥
- Port: 服务器端口

如果不设置指定环境，将使用默认环境dev配置。
## 运行项目
```sql
# 直接运行
go run main.go

# 或者编译后运行
go build -o go-blog -tags="dev" #指定环境配置
./go-blog
```
服务器将在 http://localhost:8080 启动
## 错误处理
系统提供统一的错误处理机制，所有 API 返回格式如下：
- 成功响应
```
{
    "code": 0,
    "message": "操作成功",
    "data": {}
}
```
- 错误响应
```
{
    "code": 1,
    "message": "错误信息",
}
```
## 安全特性
- 密码使用 bcrypt 加密存储
- JWT Token 用于用户认证
- API 权限验证
- 用户只能操作自己的文章
- CORS 支持
## API接口
***
### 用户认证与授权
- 用户注册
```
POST /v1/user/register
Content-Type: application/json

{
    "username" : "用户名",
    "password" : "密码",
    "email" : "邮箱"
}
```
- 用户登录
```
POST /v1/user/login
Content-Type: application/json

{
    "username":"用户名",
    "password":"密码"
}
```
### 文章管理功能
- 创建文章 （需要认证）
```
POST /v1/post/create
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
    "title":"文章标题",
    "content":"文章内容"
}
```
- 获取文章
```
GET /v1/post/find/0  # 所有文章
GET /v1/post/find/{id} # 单个文章
```
- 更新文章 需要认证，只能更新自己的文章）
```
PUT /v1/post/update/{id}
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
    "title":"文章标题",
    "content":"文章内容"
}
```
- 删除文章（需要认证，只能删除自己的文章）
```
DELETE /v1/post/delete/{id}
Authorization: Bearer {jwt_token}
```
### 评论功能
- 创建评论（需要认证）
```
POST /v1/comment/create
Authorization: Bearer {jwt_token}
Content-Type: application/json
{
    "postId":文章ID,
    "content":"评论内容"
}
```
- 获取评论
```
GET /v1/comment/find?postId=1
```
## 测试用例
***
### 使用curl进行测试
- 用户注册
```
curl -X POST http://localhost:8080/v1/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "password": "123abc",
    "email": "test@example.com"
  }'
```
- 用户登录
```
curl -X POST http://localhost:8080/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "password": "123abc",
  }'
```
- 创建文章（需要替换 {token} 为登录获取的 JWT）
```
curl -X POST http://localhost:8080/v1/post/create \
  -H Authorization: Bearer {jwt_token} \
  -H "Content-Type: application/json" \
  -d '{
    "title":"user3titile",
    "content":"user3content"
}'
```
- 获取所有文章
```
curl -X GET http://localhost:8080/v1/post/find/0
```
- 获取指定文章
```
curl -X GET http://localhost:8080/v1/post/find/9
```
- 修改文章（需要替换 {token} 为登录获取的 JWT）
```
curl -X PUT http://localhost:8080/v1/update/8 \
  -H Authorization: Bearer {jwt_token} \
  -H "Content-Type: application/json" \
  -d '{
    "title":"user3titile",
    "content":"user3content"
}'
```
- 删除指定文章 （需要替换 {token} 为登录获取的 JWT）
```
curl -X DELETE http://localhost:8080/v1/post/delete/1 \
     -H Authorization: Bearer {jwt_token} \
```
- 创建文章（需要替换 {token} 为登录获取的 JWT）
```
curl -X POST http://localhost:8080/v1/comment/create \
  -H Authorization: Bearer {jwt_token} \
  -H "Content-Type: application/json" \
  -d '{
    "postId":7,
    "content":"测试中文abc123"
}'
```
- 获取文章评论
```
curl -X GET http://localhost:8080/v1/comment/find?postId=1
```
