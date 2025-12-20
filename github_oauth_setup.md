# GitHub OAuth 配置指南

## 概述

本文档说明如何配置 GitHub OAuth 登录功能。

---

## 功能说明

### 1. 忘记密码功能

用户可以通过邮箱重置密码：

1. 点击"忘记密码？"
2. 输入注册邮箱
3. 接收 6 位数字验证码（演示模式直接显示）
4. 输入验证码验证
5. 设置新密码

**注意**: 当前为演示模式，验证码直接返回给前端。生产环境需要集成邮件服务。

### 2. GitHub 登录功能

用户需要先绑定 GitHub 账号才能使用 GitHub 登录：

1. 注册普通账号
2. 登录后在设置中绑定 GitHub ID
3. 下次可以直接使用 GitHub 登录

---

## 配置步骤

### 步骤 1: 创建 GitHub OAuth App

1. 访问 [GitHub Developer Settings](https://github.com/settings/developers)
2. 点击 "New OAuth App"
3. 填写信息：
   - **Application name**: CodeSage
   - **Homepage URL**: `http://localhost:3000`
   - **Authorization callback URL**: `http://localhost:3000/auth/github/callback`
4. 点击 "Register application"
5. 记录 **Client ID** 和 **Client Secret**

### 步骤 2: 配置后端

编辑 `backend-go/configs/config.yaml`:

```yaml
github_oauth:
  client_id: "your_github_client_id"
  client_secret: "your_github_client_secret"
  redirect_url: "http://localhost:3000/auth/github/callback"
```

### 步骤 3: 配置邮件服务（可选）

如果要启用邮件发送功能，编辑配置文件：

```yaml
email:
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  smtp_user: "your_email@gmail.com"
  smtp_password: "your_app_password"
  from_email: "noreply@codesage.dev"
  from_name: "CodeSage"
```

**Gmail 配置**:

1. 启用两步验证
2. 生成应用专用密码
3. 使用应用专用密码作为 `smtp_password`

---

## 使用流程

### 忘记密码流程

```
用户点击"忘记密码？"
    ↓
输入邮箱
    ↓
后端生成6位验证码
    ↓
（演示模式）直接返回验证码
（生产模式）发送邮件
    ↓
用户输入验证码
    ↓
验证通过
    ↓
设置新密码
    ↓
密码重置成功
```

### GitHub 登录流程

```
新用户注册
    ↓
登录系统
    ↓
进入设置页面
    ↓
绑定 GitHub ID
    ↓
退出登录
    ↓
点击 GitHub 登录
    ↓
输入 GitHub ID
    ↓
登录成功
```

---

## API 端点

### 忘记密码相关

| 方法 | 端点                             | 说明       |
| ---- | -------------------------------- | ---------- |
| POST | `/api/v1/auth/forgot-password`   | 发送重置码 |
| POST | `/api/v1/auth/verify-reset-code` | 验证重置码 |
| POST | `/api/v1/auth/reset-password`    | 重置密码   |

### GitHub 相关

| 方法 | 端点                        | 说明        | 需要认证 |
| ---- | --------------------------- | ----------- | -------- |
| POST | `/api/v1/auth/github/login` | GitHub 登录 | 否       |
| POST | `/api/v1/auth/github/bind`  | 绑定 GitHub | 是       |

---

## 数据库字段

新增字段：

```sql
github_id TEXT UNIQUE           -- GitHub 用户 ID
github_username TEXT            -- GitHub 用户名
reset_token TEXT                -- 密码重置令牌
reset_token_expires DATETIME    -- 令牌过期时间
```

---

## 测试

### 测试忘记密码

```bash
# 1. 发送重置码
curl -X POST http://localhost:8082/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# 2. 验证重置码
curl -X POST http://localhost:8082/api/v1/auth/verify-reset-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","code":"123456"}'

# 3. 重置密码
curl -X POST http://localhost:8082/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","code":"123456","newPassword":"newpass123"}'
```

### 测试 GitHub 登录

```bash
# 1. 绑定 GitHub（需要 token）
curl -X POST http://localhost:8082/api/v1/auth/github/bind \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"github_id":"12345","github_username":"testuser","email":"test@example.com","avatar":""}'

# 2. GitHub 登录
curl -X POST http://localhost:8082/api/v1/auth/github/login \
  -H "Content-Type: application/json" \
  -d '{"github_id":"12345"}'
```

---

## 安全建议

### 1. 密码重置

- ✅ 验证码 15 分钟过期
- ✅ 验证码只能使用一次
- ✅ 使用随机生成的 6 位数字
- ⚠️ 生产环境应通过邮件发送，不要返回给前端

### 2. GitHub 绑定

- ✅ GitHub ID 唯一性检查
- ✅ 防止重复绑定
- ✅ 需要登录才能绑定
- ⚠️ 生产环境应验证 GitHub OAuth token

### 3. 邮件安全

- ⚠️ 使用应用专用密码，不要使用账号密码
- ⚠️ 启用 TLS/SSL 加密
- ⚠️ 限制发送频率，防止滥用

---

## 生产环境部署

### 1. 环境变量

```bash
export GITHUB_CLIENT_ID="your_client_id"
export GITHUB_CLIENT_SECRET="your_client_secret"
export SMTP_USER="your_email@gmail.com"
export SMTP_PASSWORD="your_app_password"
```

### 2. 配置文件

生产环境配置应该从环境变量读取：

```go
viper.BindEnv("github_oauth.client_id", "GITHUB_CLIENT_ID")
viper.BindEnv("github_oauth.client_secret", "GITHUB_CLIENT_SECRET")
viper.BindEnv("email.smtp_user", "SMTP_USER")
viper.BindEnv("email.smtp_password", "SMTP_PASSWORD")
```

### 3. HTTPS

生产环境必须使用 HTTPS：

```yaml
github_oauth:
  redirect_url: "https://yourdomain.com/auth/github/callback"
```

---

## 故障排除

### 问题 1: 验证码无效

**原因**: 验证码已过期（15 分钟）

**解决**: 重新发送验证码

### 问题 2: GitHub 账号未绑定

**原因**: 用户未绑定 GitHub ID

**解决**:

1. 先注册普通账号
2. 登录后在设置中绑定 GitHub

### 问题 3: 邮件发送失败

**原因**: SMTP 配置错误

**解决**:

1. 检查 SMTP 服务器地址和端口
2. 检查用户名和密码
3. 确认启用了应用专用密码（Gmail）

---

## 未来改进

### 短期

- [ ] 集成真实的邮件服务
- [ ] 实现完整的 GitHub OAuth 流程
- [ ] 添加验证码发送频率限制

### 中期

- [ ] 支持更多第三方登录（Google, Twitter）
- [ ] 邮箱验证功能
- [ ] 手机号验证

### 长期

- [ ] 双因素认证 (2FA)
- [ ] 生物识别登录
- [ ] 单点登录 (SSO)

---

**文档版本**: 1.0.0  
**最后更新**: 2025-12-17  
**维护者**: CodeSage 开发团队
