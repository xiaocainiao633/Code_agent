# CodeSage 认证系统 - 快速启动指南

## 🚀 5 分钟快速启动

### 前置要求

- ✅ Go 1.22+ 已安装
- ✅ Node.js 16+ 已安装
- ✅ Git 已安装

---

## 步骤 1: 安装后端依赖

```bash
cd backend-go
go mod tidy
```

**预期输出**: 依赖包下载完成

---

## 步骤 2: 启动后端服务

```bash
go run cmd/server/main.go
```

**预期输出**:

```
INFO: Starting CodeSage Go Backend Server
INFO: Configuration loaded successfully
INFO: Database initialized successfully at ./data/codesage.db
INFO: Starting HTTP server on 0.0.0.0:8082
```

**验证**: 打开浏览器访问 `http://localhost:8082/api/v1/health`

应该看到:

```json
{
  "status": "healthy",
  "service": "go-backend",
  "timestamp": "2025-12-17T10:00:00Z"
}
```

---

## 步骤 3: 启动前端服务

**新开一个终端窗口**:

```bash
cd fonteng
npm run dev
```

**预期输出**:

```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:3000/
➜  Network: use --host to expose
```

---

## 步骤 4: 测试注册和登录

### 方法 1: 使用前端界面 (推荐)

1. 打开浏览器访问 `http://localhost:3000`
2. 点击 **"Get Started"** 按钮
3. 切换到 **"注册"** 标签页
4. 填写信息:
   - 用户名: `testuser`
   - 邮箱: `test@example.com`
   - 密码: `password123`
   - 确认密码: `password123`
5. 点击 **"注册"** 按钮
6. 看到 "注册成功！请登录" 提示
7. 自动切换到登录页面
8. 输入用户名和密码
9. 点击 **"登录"** 按钮
10. 登录成功，跳转到主界面 🎉

### 方法 2: 使用 curl 命令

**注册**:

```bash
curl -X POST http://localhost:8082/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\",\"confirmPassword\":\"password123\"}"
```

**登录**:

```bash
curl -X POST http://localhost:8082/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"testuser\",\"password\":\"password123\"}"
```

**获取用户资料** (需要替换 YOUR_TOKEN):

```bash
curl -X GET http://localhost:8082/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🎯 功能验证清单

- [ ] 后端服务启动成功
- [ ] 前端服务启动成功
- [ ] 健康检查接口正常
- [ ] 用户注册成功
- [ ] 用户登录成功
- [ ] 登录后跳转到主界面
- [ ] Token 自动保存
- [ ] 刷新页面保持登录状态

---

## 📁 重要文件位置

| 文件       | 路径                                | 说明          |
| ---------- | ----------------------------------- | ------------- |
| 后端主程序 | `backend-go/cmd/server/main.go`     | Go 服务入口   |
| 配置文件   | `backend-go/configs/config.yaml`    | 服务配置      |
| 数据库     | `backend-go/data/codesage.db`       | SQLite 数据库 |
| 前端登录页 | `fonteng/src/views/Login/index.vue` | 登录界面      |
| API 服务   | `fonteng/src/services/api.ts`       | API 封装      |
| 认证状态   | `fonteng/src/stores/auth.ts`        | 状态管理      |

---

## 🔧 常见问题

### Q1: 后端启动失败，提示端口被占用

**解决方案**: 修改配置文件中的端口

```yaml
# backend-go/configs/config.yaml
server:
  port: "8083" # 改为其他端口
```

### Q2: 前端无法连接后端

**解决方案**: 检查前端环境变量

```env
# fonteng/.env
VITE_API_BASE_URL=http://localhost:8082
```

### Q3: 注册时提示"用户名或邮箱已存在"

**解决方案**: 使用不同的用户名或邮箱，或删除数据库重新开始

```bash
rm backend-go/data/codesage.db
```

### Q4: 登录后刷新页面需要重新登录

**解决方案**: 确保勾选了"记住我"选项，或检查浏览器是否禁用了 localStorage

---

## 📚 下一步

现在你已经成功启动了认证系统！接下来可以:

1. 📖 阅读 [API 文档](backend-go/AUTH_API_DOCUMENTATION.md)
2. 🔍 查看 [完整实现总结](AUTHENTICATION_SUMMARY.md)
3. 🚀 查看 [部署指南](AUTHENTICATION_SETUP.md)
4. 💻 开始开发你的功能

---

## 🎨 界面预览

### 登录页面

- 深色主题
- 世界地图点阵背景动画
- 浮动粒子效果
- 登录/注册标签页切换
- 社交登录按钮
- 游客模式和本地模式

### 主界面

- 侧边栏导航
- 用户头像和下拉菜单
- 代码分析、转换、测试生成等功能
- 实时任务状态更新

---

## 🛡️ 安全提示

- ✅ 密码使用 bcrypt 加密
- ✅ JWT token 7 天自动过期
- ✅ 所有 API 请求自动携带 token
- ⚠️ 生产环境请使用 HTTPS
- ⚠️ 修改 JWT 密钥为环境变量

---

## 📞 获取帮助

遇到问题？查看:

- [完整文档](AUTHENTICATION_SETUP.md)
- [API 文档](backend-go/AUTH_API_DOCUMENTATION.md)
- [实现总结](AUTHENTICATION_SUMMARY.md)

---

**祝你使用愉快！** 🎉

---

**版本**: 1.0.0  
**更新日期**: 2025-12-17
