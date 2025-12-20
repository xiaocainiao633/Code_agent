# CodeSage AI Agent Backend

面向遗留系统的智能重构与现代化助手 - AI Agent 层后端服务

## 项目概述

CodeSage AI Agent Backend 是一个基于 FastAPI 和本地大语言模型（Ollama）的智能代码重构服务。它能够：

- **代码分析**: 深度分析代码结构、复杂度和潜在问题
- **智能转换**: 将 Python 2 代码安全转换为 Python 3
- **测试生成**: 自动生成单元测试和影子测试
- **业务逻辑推断**: 利用 AI 理解代码的业务含义
- **安全检测**: 识别代码中的安全漏洞和兼容性问题

## 技术栈

- **后端框架**: FastAPI + Uvicorn
- **AI 模型**: Ollama + llama3.2 (本地部署)
- **代码分析**: Python AST, libcst
- **数据存储**: ChromaDB (向量数据库)
- **日志**: Loguru
- **配置**: Pydantic Settings

## 快速开始

### 环境要求

- Python 3.10+
- Ollama (已安装 llama3.2 模型)
- 8GB+ RAM (推荐 16GB)

### 安装步骤

1. **克隆项目**

   ```bash
   git clone <repository-url>
   cd backend-ai-agent
   ```

2. **创建虚拟环境**

   ```bash
   python -m venv venv
   source venv/bin/activate  # Linux/Mac
   # 或
   venv\Scripts\activate  # Windows
   ```

3. **安装依赖**

   ```bash
   pip install -r requirements.txt
   ```

4. **配置环境**

   ```bash
   cp .env.example .env
   # 编辑 .env 文件，配置你的 Ollama 设置
   ```

5. **启动服务**

   ```bash
   # Linux/Mac
   ./scripts/start.sh

   # Windows
   scripts\start.bat
   ```

### 手动启动

```bash
python main.py
```

服务将在 `http://localhost:8000` 启动

## API 文档

启动服务后，可以访问：

- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc
- **健康检查**: http://localhost:8000/health

## 核心功能

### 1. 代码分析 (`/api/v1/analyze`)

分析代码的复杂度、依赖关系、安全问题和兼容性：

```bash
curl -X POST "http://localhost:8000/api/v1/analyze" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "print \"Hello, World!\"",
    "language": "python",
    "filename": "example.py"
  }'
```

### 2. Python 2 转 3 (`/api/v1/convert/python2-to-3`)

将 Python 2 代码转换为 Python 3：

```bash
curl -X POST "http://localhost:8000/api/v1/convert/python2-to-3" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "print \"Hello, World!\"",
    "language": "python",
    "conversion_type": "python_2_to_3"
  }'
```

### 3. 测试生成 (`/api/v1/generate-tests`)

为代码生成单元测试：

```bash
curl -X POST "http://localhost:8000/api/v1/generate-tests" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "def add(a, b): return a + b",
    "language": "python",
    "test_framework": "pytest"
  }'
```

### 4. 健康检查 (`/api/v1/health`)

检查服务状态和依赖项：

```bash
curl "http://localhost:8000/api/v1/health"
```

## 项目结构

```
backend-ai-agent/
├── app/                    # 应用代码
│   ├── api/               # API 路由
│   ├── models/            # 数据模型
│   ├── services/          # 业务逻辑服务
│   └── utils/             # 工具函数
├── scripts/               # 启动脚本
├── tests/                 # 测试文件
├── docs/                  # 文档
├── main.py               # 主应用文件
├── requirements.txt      # 依赖列表
└── .env.example         # 环境配置模板
```

## 配置说明

### 环境变量

| 变量名                     | 说明              | 默认值                   |
| -------------------------- | ----------------- | ------------------------ |
| `OLLAMA_HOST`              | Ollama 服务地址   | `http://localhost:11434` |
| `OLLAMA_MODEL`             | 使用的模型名称    | `llama3.2`               |
| `API_HOST`                 | API 监听地址      | `0.0.0.0`                |
| `API_PORT`                 | API 监听端口      | `8000`                   |
| `CHROMA_PERSIST_DIRECTORY` | ChromaDB 存储路径 | `./chroma_db`            |
| `MAX_CODE_SIZE`            | 最大代码大小      | `1048576` (1MB)          |

## 开发指南

### 添加新的语言支持

1. 在 `app/models/schemas.py` 中添加新的语言类型
2. 在 `app/services/code_analyzer.py` 中实现解析逻辑
3. 在 `app/services/code_converter.py` 中实现转换逻辑
4. 在 `app/services/test_generator.py` 中实现测试生成逻辑

### 添加新的转换类型

1. 在 `app/models/schemas.py` 中添加新的转换类型
2. 在 `app/services/code_converter.py` 中实现转换逻辑
3. 更新 API 路由以支持新的转换类型

## 故障排除

### Ollama 连接问题

确保 Ollama 服务正在运行：

```bash
ollama serve
```

检查模型是否已安装：

```bash
ollama list
```

如果 llama3.2 未安装，请运行：

```bash
ollama pull llama3.2
```

### 内存不足

如果遇到内存不足的问题：

1. 减小批处理大小
2. 降低并发请求数
3. 使用更小的模型

## 性能优化

- 使用异步处理提高并发性能
- 实现请求缓存机制
- 添加请求队列管理
- 使用连接池管理数据库连接

## 安全考虑

- 所有代码处理都在本地进行，不会发送到外部服务
- 实现代码大小限制防止 DoS 攻击
- 输入验证和清理
- 错误信息脱敏

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- 项目维护者: CodeSage Team
- 邮箱: team@codesage.ai
- 项目主页: https://github.com/codesage-ai/backend

---

**注意**: 这是一个正在开发中的项目，API 可能会发生变化。请关注版本更新和变更日志。
