# CodeSage Go Backend

CodeSage Go Backend 是一个高性能的工程层服务，负责处理前端请求、任务调度、用户认证和与 Python AI Agent 的通信。

## 技术栈

- **语言**: Go 1.22+
- **数据库**: SQLite 3
- **认证**: JWT (JSON Web Token)
- **密码加密**: bcrypt
- **WebSocket**: 实时通信支持

## 功能特性

### 1. 用户认证系统

- 用户注册和登录
- JWT token 认证
- 密码加密存储（bcrypt）
- 用户资料管理
- 密码修改功能

### 2. 任务管理系统

- 异步任务创建和调度
- 任务状态跟踪
- 任务结果存储
- 支持多种任务类型：
  - 代码分析（analysis）
  - 代码转换（convert）
  - 测试生成（test）
  - Git 操作（git_clone, git_analyze, git_history, git_diff）
  - 批处理（batch）

### 3. WebSocket 实时通信

- Agent 思考流推送
- 任务进度实时更新
- 自动重连机制

### 4. 文件管理

- 文件上传
- 文件列表查询
- 批量文件处理

## 项目结构

```
backend-go/
├── cmd/
│   └── server/          # 服务入口
├── internal/
│   ├── config/          # 配置管理
│   ├── database/        # 数据库操作
│   ├── handlers/        # HTTP 处理器
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   ├── services/        # 业务逻辑
│   ├── utils/           # 工具函数
│   └── websocket/       # WebSocket 管理
├── configs/
│   └── config.yaml      # 配置文件
├── data/                # 数据库文件
├── logs/                # 日志文件
└── scripts/             # 脚本文件
```

## 快速开始

### 1. 安装依赖

```bash
cd backend-go
go mod download
```

### 2. 配置

编辑 `configs/config.yaml`：

```yaml
server:
  port: "8082"
  host: "0.0.0.0"

python_agent:
  host: "localhost"
  port: "8000"
  timeout: 30s
  retry_count: 3

task_scheduler:
  max_concurrent_tasks: 10
  task_timeout: 10m
  result_retention: 24h

websocket:
  ping_interval: 30s
  pong_timeout: 60s
  max_message_size: 512000
```

### 3. 运行服务

```bash
# 开发模式
go run cmd/server/main.go

# 或使用编译后的二进制文件
./bin/server
```

### 4. 构建

```bash
go build -o bin/server cmd/server/main.go
```

## API 文档

### 基础信息

- **Base URL**: `http://localhost:8082`
- **API Version**: v1

### 认证端点

#### 用户注册

```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "confirmPassword": "password123"
}
```

#### 用户登录

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

#### 获取用户资料

```
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

#### 更新用户资料

```
PUT /api/v1/auth/profile/update
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "newusername",
  "email": "newemail@example.com",
  "avatar": "https://example.com/avatar.jpg"
}
```

#### 修改密码

```
POST /api/v1/auth/password/change
Authorization: Bearer <token>
Content-Type: application/json

{
  "oldPassword": "oldpassword123",
  "newPassword": "newpassword123"
}
```

### 任务管理端点

#### 创建任务

```
POST /api/v1/tasks
Content-Type: application/json

{
  "type": "analysis",
  "name": "Python Code Analysis",
  "description": "Analyze Python code",
  "params": {
    "code": "print('Hello World')",
    "language": "python"
  }
}
```

#### 获取任务列表

```
GET /api/v1/tasks
```

#### 获取任务详情

```
GET /api/v1/tasks/{taskId}
```

#### 获取任务结果

```
GET /api/v1/tasks/{taskId}/result
```

#### 取消任务

```
DELETE /api/v1/tasks/{taskId}
```

### WebSocket 端点

#### Agent 思考流

```
ws://localhost:8082/ws/agent/{taskId}
```

#### 任务进度

```
ws://localhost:8082/ws/progress/{taskId}
```

### 健康检查

```
GET /api/v1/health
GET /api/v1/health/detailed
```

## 安全特性

1. **密码加密**: 使用 bcrypt 算法加密存储密码
2. **JWT 认证**: 使用 JWT token 进行用户认证（7 天有效期）
3. **输入验证**: 正则表达式验证用户输入
4. **SQL 注入防护**: 使用参数化查询
5. **CORS 支持**: 配置跨域资源共享

## 数据库结构

### users 表

```sql
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'user',
  avatar TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### tasks 表

```sql
CREATE TABLE tasks (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  status TEXT NOT NULL,
  name TEXT,
  description TEXT,
  params TEXT,
  result TEXT,
  error TEXT,
  progress INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  started_at DATETIME,
  completed_at DATETIME
);
```

## 错误代码

| HTTP 状态码 | 说明           |
| ----------- | -------------- |
| 200         | 成功           |
| 201         | 创建成功       |
| 400         | 请求参数错误   |
| 401         | 未授权         |
| 404         | 资源不存在     |
| 405         | 方法不允许     |
| 500         | 服务器内部错误 |
| 503         | 服务不可用     |

## 开发

### 运行测试

```bash
go test ./internal/... -v
```

### 代码格式化

```bash
go fmt ./...
```

### 代码检查

```bash
go vet ./...
```

## 生产部署建议

1. **环境变量**: 使用环境变量存储敏感配置（JWT 密钥等）
2. **HTTPS**: 生产环境必须使用 HTTPS
3. **日志**: 配置日志轮转和持久化
4. **监控**: 添加性能监控和告警
5. **限流**: 实现 API 限流机制
6. **备份**: 定期备份数据库

## 贡献

欢迎提交 Issue 和 Pull Request！

---
