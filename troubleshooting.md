# 故障排除指南

## 注册功能问题诊断

### 问题：注册后显示"注册失败"

#### 步骤 1: 检查后端是否启动

```bash
# 检查后端是否在运行
curl http://localhost:8082/api/v1/health
```

**预期响应**:

```json
{ "status": "healthy", "service": "go-backend", "timestamp": "..." }
```

如果无响应，启动后端：

```bash
cd backend-go
go run cmd/server/main.go
```

#### 步骤 2: 测试注册 API

**Windows**:

```bash
cd backend-go
test_register.bat
```

**Linux/Mac**:

```bash
cd backend-go
chmod +x test_register.sh
./test_register.sh
```

**预期响应**:

```json
{
  "message": "注册成功",
  "user": {
    "id": 1,
    "username": "testuser123",
    "email": "test123@example.com",
    "role": "user",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

#### 步骤 3: 检查浏览器控制台

1. 打开浏览器开发者工具 (F12)
2. 切换到 Console 标签
3. 尝试注册
4. 查看控制台输出

**查找以下信息**:

- `开始注册，用户名: xxx`
- `Registration response: {...}`
- `注册结果: {...}`

#### 步骤 4: 检查网络请求

1. 打开浏览器开发者工具 (F12)
2. 切换到 Network 标签
3. 尝试注册
4. 查找 `register` 请求

**检查项**:

- 请求 URL: `http://localhost:8082/api/v1/auth/register`
- 请求方法: POST
- 状态码: 201 (成功) 或 400 (失败)
- 响应内容

#### 步骤 5: 检查 CORS 错误

如果看到 CORS 错误，检查：

1. 后端配置文件 `backend-go/configs/config.yaml`:

```yaml
cors:
  allowed_origins:
    - "http://localhost:3000" # 确保包含前端地址
```

2. 重启后端服务

#### 步骤 6: 检查数据库

```bash
# 查看数据库文件是否存在
ls backend-go/data/codesage.db

# 使用 SQLite 查看用户表
sqlite3 backend-go/data/codesage.db "SELECT * FROM users;"
```

### 常见错误及解决方案

#### 错误 1: "用户名或邮箱已存在"

**原因**: 该用户名或邮箱已被注册

**解决**: 使用不同的用户名或邮箱

#### 错误 2: "用户名格式不正确"

**原因**: 用户名不符合规则

**要求**:

- 3-20 个字符
- 只允许字母、数字、下划线

#### 错误 3: "邮箱格式不正确"

**原因**: 邮箱格式无效

**要求**: 标准邮箱格式 (例如: user@example.com)

#### 错误 4: "密码格式不正确"

**原因**: 密码不符合规则

**要求**: 6-20 个字符

#### 错误 5: "两次输入的密码不一致"

**原因**: 密码和确认密码不匹配

**解决**: 确保两次输入的密码完全相同

#### 错误 6: 网络错误

**可能原因**:

1. 后端未启动
2. 端口被占用
3. 防火墙阻止

**解决**:

1. 启动后端服务
2. 修改端口配置
3. 关闭防火墙或添加例外

### 调试技巧

#### 1. 启用详细日志

后端日志位置: `backend-go/logs/go-backend.log`

查看实时日志:

```bash
tail -f backend-go/logs/go-backend.log
```

#### 2. 使用 curl 测试

```bash
# 测试注册
curl -X POST http://localhost:8082/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456","confirmPassword":"123456"}' \
  -v
```

#### 3. 检查前端环境变量

文件: `fonteng/.env`

```env
VITE_API_BASE_URL=http://localhost:8082
```

确保 URL 正确，重启前端服务。

### 重置数据库

如果需要重新开始：

```bash
# 删除数据库文件
rm backend-go/data/codesage.db

# 重启后端服务（会自动创建新数据库）
cd backend-go
go run cmd/server/main.go
```

### 获取帮助

如果问题仍未解决：

1. 查看完整日志
2. 检查浏览器控制台
3. 检查网络请求详情
4. 提供错误信息截图

---

**最后更新**: 2025-12-17
