# CodeSage - 智能代码重构与现代化助手

<div align="center">

![CodeSage Logo](https://img.shields.io/badge/CodeSage-v1.0.0-blue)
![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![Vue Version](https://img.shields.io/badge/Vue-3.5+-4FC08D?logo=vue.js)
![Python Version](https://img.shields.io/badge/Python-3.10+-3776AB?logo=python)
![License](https://img.shields.io/badge/license-MIT-green)

**面向遗留系统的本地化、隐私安全的 AI Agent 系统**

[快速开始](#-快速开始) • [功能特性](#-功能特性) • [在线演示](#-在线演示) • [文档](#-文档)

</div>

---

## 📖 项目简介

CodeSage 是一个**完全本地化、隐私安全**的 AI 驱动代码重构助手，专为企业遗留系统现代化设计。它帮助开发者安全、高效地将老旧代码（如 Python 2、VB6、COBOL 等）重构并迁移至现代技术栈。

### 为什么选择 CodeSage？

- 🔒 **100% 本地化**: 所有代码分析与 AI 推理均在本地完成，企业敏感代码永不出内网
- 🤖 **AI 驱动**: 集成本地大模型（Ollama），提供智能代码分析和重构建议
- 🚀 **高性能**: Go 后端高并发处理，支持大规模代码库分析
- 🎨 **现代界面**: Vue 3 + Element Plus 深色主题，Monaco Editor 代码编辑器
- 🔐 **企业级安全**: JWT 认证、密码加密、输入验证等多重安全保障
- ⚡ **实时反馈**: WebSocket 实时通信，任务进度实时更新

---

<img width="2486" height="1407" alt="image" src="https://github.com/user-attachments/assets/6dc3f48d-4b00-4f9b-82db-3c1d43211ad8" />

<img width="2539" height="1393" alt="image" src="https://github.com/user-attachments/assets/e0cd2919-f4b5-45d8-8c8c-7c2d5e858a83" />

## ✨ 功能特性

### 核心功能

#### 🔍 代码分析

- **深度分析**: 复杂度计算、依赖提取、安全检测
- **业务逻辑推断**: AI 理解代码意图，生成业务逻辑文档
- **兼容性检测**: 自动识别 Python 2/3 兼容性问题
- **代码质量评估**: PEP8 规范检查、代码异味检测

#### 🔄 代码转换

- **Python 2 → Python 3**: 15+ 种自动转换规则
- **语法现代化**: 自动应用现代 Python 特性
- **类型注解**: 智能添加类型提示
- **差异对比**: 并排/逐行对比转换前后代码

#### 🧪 测试生成

- **自动生成单元测试**: 支持 pytest/unittest
- **测试覆盖率分析**: 可视化覆盖率报告
- **边界测试**: 自动生成边界条件测试用例
- **影子测试**: 对比原始代码和转换后代码的行为

#### 📊 Git 分析

- **仓库历史分析**: 代码变更追踪、贡献者统计
- **文件历史**: 查看文件的完整修改历史
- **差异对比**: 任意两个提交之间的代码差异
- **分支管理**: 可视化分支结构

#### 👤 用户系统

- **多种登录方式**: 账号登录、游客模式、本地模式
- **用户资料管理**: 头像上传、个人信息编辑
- **任务历史**: 查看所有历史任务和结果
- **设置中心**: 主题、编辑器、自动保存等配置

---

## 🚀 快速开始

### 前置要求

确保你的系统已安装以下软件：

- **Go** 1.22 或更高版本
- **Node.js** 16 或更高版本
- **Python** 3.10 或更高版本
- **Git**
- **Ollama** (可选，用于 AI 功能)

### 一键启动

```bash
# 克隆项目
git clone https://github.com/your-repo/codesage.git
cd codesage

# 启动 Go 后端
cd backend-go
go mod tidy
go run cmd/server/main.go &

# 启动前端
cd ../frontend
npm install
npm run dev &

# 启动 Python AI Agent (可选)
cd ../backend-ai-agent
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python main.py &
```

### 访问应用

- **前端界面**: http://localhost:3000
- **Go 后端**: http://localhost:8082
- **Python AI Agent**: http://localhost:8000

### 首次使用

1. 打开浏览器访问 http://localhost:3000
2. 点击 "Get Started" 按钮
3. 选择登录方式：
   - **注册新账号**: 创建个人账号
   - **游客模式**: 无需注册，快速体验
   - **本地模式**: 仅本地存储，无需后端
4. 开始使用代码分析、转换等功能

📚 **详细指南**: [快速启动文档](quick_start.md)

---

## 🎬 在线演示

### 代码分析示例

```python
# 输入 Python 2 代码
print "Hello World"
def old_function(x, y):
    return x + y
```

**分析结果**:

- 复杂度: 1.2
- 兼容性问题: 2 个
- 建议: 使用 print() 函数，添加类型注解

### 代码转换示例

**转换前** (Python 2):

```python
print "Hello World"
except IOError, e:
    print "Error:", e
```

**转换后** (Python 3):

```python
print("Hello World")
except IOError as e:
    print("Error:", e)
```

---

## 🏗️ 系统架构

### 三层架构设计

```
┌─────────────────────────────────────────────────────────┐
│                    前端层 (Vue 3)                        │
│  Monaco Editor | Element Plus | WebSocket | Pinia       │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/WebSocket
┌────────────────────▼────────────────────────────────────┐
│                 工程层 (Go Backend)                      │
│  任务调度 | 文件处理 | Git 分析 | 用户认证 | WebSocket │
└────────────────────┬────────────────────────────────────┘
                     │ RESTful API
┌────────────────────▼────────────────────────────────────┐
│              AI Agent 层 (Python + Ollama)               │
│  代码分析 | 代码转换 | 测试生成 | LLM 服务 | 向量检索  │
└─────────────────────────────────────────────────────────┘
```

### 技术栈

| 层级         | 技术               | 说明            |
| ------------ | ------------------ | --------------- |
| **前端**     | Vue 3 + TypeScript | 响应式 UI 框架  |
|              | Element Plus       | UI 组件库       |
|              | Monaco Editor      | 代码编辑器      |
|              | Pinia              | 状态管理        |
|              | Axios + WebSocket  | 网络通信        |
| **后端**     | Go 1.22+           | 高性能后端服务  |
|              | SQLite             | 轻量级数据库    |
|              | JWT                | 用户认证        |
|              | Gorilla WebSocket  | 实时通信        |
| **AI Agent** | Python 3.10+       | AI 服务层       |
|              | FastAPI            | 高性能 API 框架 |
|              | Ollama             | 本地大模型      |
|              | ChromaDB           | 向量数据库      |
|              | AST                | 代码解析        |

---

## 📚 文档

### 用户文档

- [快速启动指南](quick_start.md) - 5 分钟快速上手
- [故障排除](troubleshooting.md) - 常见问题解决方案

### 开发文档

- [Go 后端文档](backend-go/README.md) - 后端 API 和架构
- [Python AI Agent 文档](backend-ai-agent/README.md) - AI 服务说明

### 配置文档

- [认证系统配置](authentication_summary.md) - 用户认证详解
- [GitHub OAuth 配置](github_oauth_setup.md) - OAuth 集成指南
- [Kimi API 配置](kimi_api_setup.md) - AI 模型配置

### API 文档

- 详见 [Go 后端文档](backend-go/README.md) - 完整的 API 接口说明

---

## 🔐 安全特性

### 认证与授权

- ✅ JWT Token 认证（7 天有效期）
- ✅ bcrypt 密码加密（成本因子 10）
- ✅ Bearer Token 格式
- ✅ 自动 Token 刷新
- ✅ 会话管理

### 数据安全

- ✅ 参数化查询（防 SQL 注入）
- ✅ 输入验证和清理
- ✅ CORS 白名单
- ✅ 密码不返回前端
- ✅ 本地数据存储

### 隐私保护

- ✅ 代码不上传云端
- ✅ 本地 AI 推理
- ✅ 可选的游客模式
- ✅ 数据自动清理

---

## 📊 项目状态

### 开发进度

| 模块            | 状态    | 完成度 |
| --------------- | ------- | ------ |
| 用户认证系统    | ✅ 完成 | 100%   |
| Go 后端工程层   | ✅ 完成 | 100%   |
| Python AI Agent | ✅ 完成 | 100%   |
| 前端界面        | ✅ 完成 | 95%    |
| 文档            | ✅ 完成 | 100%   |

### 路线图

- [x] 项目基础架构
- [x] 用户认证系统
- [x] 代码分析服务
- [x] 代码转换服务
- [x] 测试生成服务
- [x] Git 分析功能
- [x] WebSocket 实时通信
- [x] 前端界面开发
- [ ] 多语言支持（Java, JavaScript, Go）
- [ ] 多智能体协作
- [ ] 代码审查功能
- [ ] CI/CD 集成

---

## 🛠️ 开发

### 项目结构

```
codesage/
├── backend-go/              # Go 后端服务
│   ├── cmd/server/          # 主程序入口
│   ├── internal/            # 内部模块
│   │   ├── config/          # 配置管理
│   │   ├── database/        # 数据库操作
│   │   ├── handlers/        # HTTP 处理器
│   │   ├── middleware/      # 中间件
│   │   ├── models/          # 数据模型
│   │   ├── services/        # 业务逻辑
│   │   └── websocket/       # WebSocket 管理
│   ├── configs/             # 配置文件
│   └── data/                # SQLite 数据库
│
├── backend-ai-agent/        # Python AI Agent
│   ├── app/                 # 应用代码
│   │   ├── api/             # API 路由
│   │   ├── services/        # AI 服务
│   │   ├── models/          # 数据模型
│   │   └── utils/           # 工具函数
│   ├── main.py              # 主程序
│   └── requirements.txt     # Python 依赖
│
├── frontend/                # Vue 3 前端
│   ├── src/
│   │   ├── components/      # 可复用组件
│   │   ├── views/           # 页面组件
│   │   ├── stores/          # Pinia 状态管理
│   │   ├── services/        # API 服务
│   │   ├── router/          # 路由配置
│   │   └── utils/           # 工具函数
│   ├── public/              # 静态资源
│   └── package.json         # 前端依赖
│
├── quick_start.md           # 快速开始
├── troubleshooting.md       # 故障排除
├── changelog.md             # 更新日志
├── authentication_summary.md # 认证系统配置
├── github_oauth_setup.md    # GitHub OAuth 配置
├── kimi_api_setup.md        # Kimi API 配置
├── integration_test_guide.md # 集成测试指南
└── README.md                # 项目说明
```

### 本地开发

```bash
# 后端开发（热重载）
cd backend-go
go run cmd/server/main.go

# 前端开发（热重载）
cd frontend
npm run dev

# AI Agent 开发
cd backend-ai-agent
python main.py
```

### 构建生产版本

```bash
# 构建 Go 后端
cd backend-go
go build -o bin/server cmd/server/main.go

# 构建前端
cd frontend
npm run build

# 打包 Python Agent
cd backend-ai-agent
pip install -r requirements.txt
```

### 运行测试

```bash
# Go 测试
cd backend-go
go test ./internal/... -v

# 前端测试
cd frontend
npm run test

# Python 测试
cd backend-ai-agent
pytest
```

---

## 🤝 贡献

我们欢迎所有形式的贡献！无论是报告 Bug、提出新功能建议，还是提交代码。

### 贡献流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发规范

- 遵循项目代码风格
- 添加必要的测试
- 更新相关文档
- 提交信息清晰明确

---

## 📝 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

感谢以下开源项目的支持：

- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Go](https://golang.org/) - 高性能编程语言
- [FastAPI](https://fastapi.tiangolo.com/) - 现代 Python Web 框架
- [Ollama](https://ollama.ai/) - 本地大模型运行时
- [Element Plus](https://element-plus.org/) - Vue 3 UI 组件库
- [Monaco Editor](https://microsoft.github.io/monaco-editor/) - VS Code 编辑器核心

---

[⬆ 回到顶部](#codesage---智能代码重构与现代化助手)

</div>
