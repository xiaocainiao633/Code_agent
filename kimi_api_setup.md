# Kimi API 配置指南

## 🎯 概述

Kimi（月之暗面）是由 Moonshot AI 提供的优秀大语言模型，特别适合中文语境下的代码分析和转换任务。本指南将帮助您快速配置 Kimi API 以启用 CodeSage 的 AI 功能。

## 📋 前提条件

- 已安装并运行 CodeSage 项目
- 拥有有效的 Kimi API 密钥
- 基本的命令行操作知识

## 🔑 获取 Kimi API 密钥

### 步骤 1: 注册 Moonshot AI 账户

1. 访问 [Moonshot AI 开放平台](https://platform.moonshot.cn/)
2. 点击 "注册" 创建新账户
3. 完成邮箱验证和实名认证

### 步骤 2: 创建 API 密钥

1. 登录 Moonshot AI 控制台
2. 导航到 "API 密钥管理" 页面
3. 点击 "创建 API 密钥"
4. 为密钥命名（例如：`codesage-prod`）
5. 复制生成的 API 密钥（⚠️ 只显示一次，请妥善保存）

### 步骤 3: 充值账户（如需要）

1. 导航到 "账户充值" 页面
2. 选择合适的充值金额
3. Kimi 提供优惠的价格：
   - moonshot-v1-8k: ¥0.012/1K tokens
   - moonshot-v1-32k: ¥0.024/1K tokens
   - moonshot-v1-128k: ¥0.048/1K tokens

## ⚙️ 配置 CodeSage 项目

### 步骤 1: 定位配置文件

找到 Python AI Agent 的配置文件：

```bash
cd backend-ai-agent
```

### 步骤 2: 创建环境配置文件

如果还没有 `.env` 文件，复制示例配置：

```bash
cp .env.example .env
```

### 步骤 3: 编辑配置文件

使用您喜欢的文本编辑器打开 `.env` 文件：

```bash
# Windows
notepad .env

# macOS/Linux
nano .env
```

### 步骤 4: 添加 Kimi 配置

在 `.env` 文件中添加或修改以下配置：

```env
# LLM 提供商配置
LLM_PROVIDER=kimi

# Kimi 配置
KIMI_API_KEY=your_api_key_here  # 替换为您的实际 API 密钥
KIMI_MODEL=moonshot-v1-8k
KIMI_BASE_URL=https://api.moonshot.cn/v1

# 通用 LLM 配置
LLM_MAX_TOKENS=2000
LLM_TEMPERATURE=0.1
LLM_TIMEOUT=30
LLM_RETRY_COUNT=3
```

### 步骤 5: 完整的配置示例

您的 `.env` 文件应该类似于：

```env
# LLM 提供商配置
LLM_PROVIDER=kimi

# Kimi 配置
KIMI_API_KEY=sk-1234567890abcdef1234567890abcdef
KIMI_MODEL=moonshot-v1-8k
KIMI_BASE_URL=https://api.moonshot.cn/v1

# FastAPI 配置
API_HOST=0.0.0.0
API_PORT=8000
API_DEBUG=true

# ChromaDB 配置
CHROMA_PERSIST_DIRECTORY=./chroma_db
CHROMA_COLLECTION_NAME=code_embeddings

# 日志配置
LOG_LEVEL=INFO
LOG_FILE=./logs/app.log

# 安全配置
MAX_CODE_SIZE=1048576  # 1MB
ALLOWED_FILE_EXTENSIONS=.py,.js,.java,.cpp,.c
```

## 🚀 重启服务

### 步骤 1: 重启 Python AI Agent

```bash
# 停止当前服务（如果正在运行）
Ctrl+C

# 重新启动
python main.py
```

### 步骤 2: 验证配置

检查日志输出，确认 Kimi 配置已加载：

```
INFO:     Started server process [12345]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     Uvicorn running on http://0.0.0.0:8000
```

## 🧪 测试配置

### 测试 1: 健康检查

```bash
curl http://localhost:8000/api/v1/health
```

预期响应：

```json
{
  "status": "healthy",
  "service": "CodeSage AI Agent",
  "version": "1.0.0"
}
```

### 测试 2: 代码分析功能

1. 打开前端界面：http://localhost:3001
2. 导航到 "代码分析" 页面
3. 上传一个 Python 文件
4. 点击 "开始分析"
5. 观察分析结果

## 💡 Kimi 模型选择指南

### 模型对比

| 模型             | 上下文长度     | 价格             | 适用场景                 |
| ---------------- | -------------- | ---------------- | ------------------------ |
| moonshot-v1-8k   | 8,192 tokens   | ¥0.012/1K tokens | 普通代码分析、短代码转换 |
| moonshot-v1-32k  | 32,768 tokens  | ¥0.024/1K tokens | 长代码分析、复杂转换任务 |
| moonshot-v1-128k | 131,072 tokens | ¥0.048/1K tokens | 大型项目分析、批量处理   |

### 推荐配置

- **开发环境**: `moonshot-v1-8k` - 成本效益最佳
- **生产环境**: `moonshot-v1-32k` - 平衡性能和成本
- **企业级**: `moonshot-v1-128k` - 处理大型代码库

## 🎯 最佳实践

### 1. API 密钥安全

- 不要将 API 密钥提交到代码仓库
- 使用环境变量而非硬编码
- 定期轮换 API 密钥

### 2. 成本控制

- 监控 API 使用量
- 设置合理的 `LLM_MAX_TOKENS` 限制
- 根据任务复杂度选择合适的模型

### 3. 性能优化

- 调整 `LLM_TIMEOUT` 以适应网络条件
- 设置适当的重试次数
- 监控响应时间

### 4. 中文优化

Kimi 在中文语境下表现优异，特别适合：

- 中文代码注释分析
- 中文技术文档处理
- 中文错误信息解析

## 🔧 故障排除

### 问题 1: API 连接失败

**症状**: 日志显示连接错误
**解决**:

- 检查 API 密钥是否正确
- 确认网络连接正常
- 验证 `KIMI_BASE_URL` 设置

### 问题 2: 响应超时

**症状**: 任务长时间无响应
**解决**:

- 增加 `LLM_TIMEOUT` 值
- 检查 `LLM_MAX_TOKENS` 设置
- 验证网络延迟

### 问题 3: 配额不足

**症状**: API 返回配额错误
**解决**:

- 检查账户余额
- 升级服务套餐
- 优化请求频率

## 📊 价格参考

| 模型             | 输入价格         | 输出价格         |
| ---------------- | ---------------- | ---------------- |
| moonshot-v1-8k   | ¥0.012/1K tokens | ¥0.012/1K tokens |
| moonshot-v1-32k  | ¥0.024/1K tokens | ¥0.024/1K tokens |
| moonshot-v1-128k | ¥0.048/1K tokens | ¥0.048/1K tokens |

## 🆘 获取帮助

- **Moonshot AI 文档**: [https://platform.moonshot.cn/docs](https://platform.moonshot.cn/docs)
- **技术支持**: support@moonshot.cn
- **社区论坛**: [Moonshot AI Community](https://community.moonshot.cn)

## ✅ 验证成功

配置成功后，您应该能够：

1. ✅ 在 CodeSage 界面中正常使用 AI 功能
2. ✅ 代码分析返回详细的结果
3. ✅ 代码转换功能正常工作
4. ✅ 测试生成功能可用
5. ✅ 实时任务状态更新

## 🌟 Kimi 特色功能

- **中文理解**: 对中文代码注释和技术文档理解更准确
- **代码推理**: 在代码逻辑分析方面表现优异
- **多语言支持**: 支持多种编程语言的分析和转换
- **实时响应**: 快速的 API 响应时间

恭喜！您已成功配置 Kimi API，现在可以享受 CodeSage 的完整 AI 功能了！Kimi 将为您提供高质量的代码分析和转换服务。
