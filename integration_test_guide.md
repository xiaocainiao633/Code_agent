# 前后端集成测试完整指南

## 概述

本指南将帮助你完成 CodeSage 项目的前后端集成测试，验证以下数据流：

```
前端 (Vue3) → Go Backend → Python AI Agent → 返回结果
```

## 前置条件

### 必需软件

- ✅ Python 3.10+ (用于 Python Agent)
- ✅ Go 1.22+ (用于 Go Backend)
- ✅ Node.js 18+ (用于前端)
- ✅ Ollama (用于本地 LLM)

### 检查 Ollama

```bash
# 检查Ollama是否运行
ollama list

# 如果没有llama3.2模型，拉取它
ollama pull llama3.2
```

## 第一步：启动 Python AI Agent

### Windows

```bash
cd backend-ai-agent

# 激活虚拟环境 (如果有)
venv\Scripts\activate

# 安装依赖
pip install -r requirements.txt

# 启动服务
python main.py
```

### Linux/Mac

```bash
cd backend-ai-agent

# 激活虚拟环境 (如果有)
source venv/bin/activate

# 安装依赖
pip install -r requirements.txt

# 启动服务
python main.py
```

**验证启动成功：**

- 控制台显示: `INFO: Application startup complete`
- 访问: http://localhost:8001/docs (FastAPI 文档)
- 访问: http://localhost:8001/api/v1/health

### 快速测试 Python Agent

```bash
# 运行测试脚本
python test_python_agent.py
```

预期输出：

```
✅ 健康检查: 通过
✅ 代码分析: 通过
✅ 代码转换: 通过
✅ 测试生成: 通过

🎉 所有测试通过！Python Agent工作正常。
```

## 第二步：启动 Go Backend

### 编译和运行

```bash
cd backend-go

# 编译
go build -o main.exe ./cmd/server

# 运行
./main.exe
```

或者直接运行：

```bash
cd backend-go
go run cmd/server/main.go
```

**验证启动成功：**

- 控制台显示: `Starting HTTP server on localhost:8080`
- 访问: http://localhost:8080/api/v1/health
- 数据库初始化成功

### 快速测试 Go Backend

**Windows:**

```bash
test_go_backend.bat
```

**Linux/Mac:**

```bash
chmod +x test_go_backend.sh
./test_go_backend.sh
```

预期输出：

```
✅ 健康检查: 通过
✅ 创建分析任务: 通过
✅ 获取任务列表: 通过
✅ 创建转换任务: 通过
✅ 创建测试生成任务: 通过

🎉 所有测试通过！Go Backend工作正常。
```

## 第三步：启动前端

```bash
cd fonteng

# 安装依赖 (首次运行)
npm install

# 启动开发服务器
npm run dev
```

**验证启动成功：**

- 控制台显示: `Local: http://localhost:5173/`
- 浏览器访问: http://localhost:5173

## 第四步：完整集成测试

### 1. 用户注册和登录

1. 打开浏览器访问 http://localhost:5173
2. 点击"注册"
3. 填写用户信息：
   - 用户名: testuser
   - 邮箱: test@example.com
   - 密码: Test123456
4. 注册成功后自动登录

### 2. 测试代码分析功能

1. 导航到"代码分析"页面
2. 在左侧编辑器输入 Python 代码：

```python
def calculate_sum(numbers):
    total = 0
    for num in numbers:
        total += num
    return total

def main():
    numbers = [1, 2, 3, 4, 5]
    result = calculate_sum(numbers)
    print(f"Sum: {result}")
```

3. 点击"开始分析"按钮
4. 验证：
   - ✅ 任务创建成功（显示任务 ID）
   - ✅ 右侧显示分析结果
   - ✅ 显示复杂度分数
   - ✅ 显示依赖列表
   - ✅ 显示安全问题（如果有）
   - ✅ 显示业务逻辑摘要

### 3. 测试代码转换功能

1. 导航到"代码转换"页面
2. 在左侧编辑器输入 Python 2 代码：

```python
print "Hello, World!"

def divide(a, b):
    return a / b

numbers = range(10)
for i in xrange(5):
    print i

if numbers.has_key('test'):
    print "Found"
```

3. 选择转换类型："Python 2 to 3"
4. 点击"开始转换"按钮
5. 验证：
   - ✅ 任务创建成功
   - ✅ 右侧显示转换后的 Python 3 代码
   - ✅ 显示变更列表（print 语句、xrange、has_key 等）
   - ✅ 显示警告和兼容性说明
   - ✅ 可以下载转换后的代码

### 4. 测试测试生成功能

1. 导航到"测试生成"页面
2. 在左侧编辑器输入 Python 代码：

```python
def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

class Calculator:
    def multiply(self, a, b):
        return a * b

    def divide(self, a, b):
        if b == 0:
            raise ValueError("Cannot divide by zero")
        return a / b
```

3. 选择测试框架："pytest"
4. 选择测试类型："单元测试"
5. 设置覆盖率目标：80%
6. 点击"生成测试"按钮
7. 验证：
   - ✅ 任务创建成功
   - ✅ 右侧显示生成的测试代码
   - ✅ 显示测试用例列表
   - ✅ 显示覆盖率估计
   - ✅ 可以下载测试代码

### 5. 验证 Dashboard

1. 导航到"控制台"页面
2. 验证：
   - ✅ 显示最近任务列表
   - ✅ 显示任务状态（待处理、运行中、已完成）
   - ✅ 显示系统统计（总任务数、运行中、已完成、成功率）
   - ✅ 可以点击任务查看详情

## 数据流验证

### 检查 Go 后端日志

```
[INFO] Task created successfully: 测试代码分析 (ID: xxx)
[INFO] Executing analysis task: xxx
[DEBUG] Calling Python Agent API: http://localhost:8001/api/v1/analyze
[INFO] Task xxx completed successfully
```

### 检查 Python Agent 日志

```
INFO: Starting code analysis for python code
INFO: Analysis completed for python code
INFO: 127.0.0.1:xxxxx - "POST /api/v1/analyze HTTP/1.1" 200 OK
```

### 检查数据库

```bash
cd backend-go/data
sqlite3 codesage.db

# 查看任务表
SELECT id, type, status, name, created_at FROM tasks ORDER BY created_at DESC LIMIT 5;

# 查看任务状态分布
SELECT status, COUNT(*) FROM tasks GROUP BY status;
```

## 常见问题排查

### 问题 1：Python Agent 连接失败

**症状：** Go 后端日志显示 "connection refused"

**解决方案：**

1. 确认 Python Agent 正在运行：`curl http://localhost:8001/api/v1/health`
2. 检查端口是否被占用：`netstat -ano | findstr 8001` (Windows)
3. 检查防火墙设置

### 问题 2：Ollama 连接失败

**症状：** Python Agent 日志显示 "Ollama connection failed"

**解决方案：**

1. 确认 Ollama 正在运行：`ollama list`
2. 检查 Ollama 服务：`curl http://localhost:11434/api/tags`
3. 拉取模型：`ollama pull llama3.2`

### 问题 3：数据库锁定

**症状：** Go 后端报错 "database is locked"

**解决方案：**

1. 关闭所有访问数据库的程序
2. 删除数据库锁文件：`rm backend-go/data/codesage.db-shm backend-go/data/codesage.db-wal`
3. 重启 Go 后端

### 问题 4：前端无法连接后端

**症状：** 浏览器控制台显示 "Network Error"

**解决方案：**

1. 确认 Go 后端正在运行：`curl http://localhost:8080/api/v1/health`
2. 检查 CORS 配置
3. 检查前端环境变量：`fonteng/.env`

### 问题 5：WebSocket 连接失败

**症状：** 任务状态不实时更新

**解决方案：**

1. 检查 WebSocket 路由配置
2. 检查浏览器控制台 WebSocket 错误
3. 确认 Go 后端 WebSocket Hub 正常运行

## 性能基准

### 预期响应时间

- 健康检查：< 50ms
- 代码分析：2-5 秒
- 代码转换：3-8 秒
- 测试生成：5-10 秒

### 资源使用

- Python Agent：~500MB RAM
- Go Backend：~50MB RAM
- 前端：~200MB RAM
- Ollama：~2GB RAM (取决于模型)

## 成功标准

✅ **所有服务正常启动**

- Python Agent: http://localhost:8001
- Go Backend: http://localhost:8080
- 前端: http://localhost:5173

✅ **API 端点正常响应**

- 健康检查返回 200
- 任务创建返回 201
- 任务列表返回数据

✅ **数据流完整**

- 前端 → Go → Python Agent → 返回
- 任务状态正确持久化
- 结果正确显示

✅ **实时通信正常**

- WebSocket 连接成功
- 任务进度实时更新
- Dashboard 实时刷新

## 下一步

集成测试通过后，你可以：

1. **添加更多测试用例**

   - 边界条件测试
   - 错误处理测试
   - 并发测试

2. **性能优化**

   - 数据库查询优化
   - 缓存策略
   - 并发处理

3. **功能扩展**

   - 支持更多编程语言
   - 添加代码审查功能
   - 集成 Git 分析

4. **部署准备**
   - Docker 容器化
   - 生产环境配置
   - 监控和日志

## 总结

完成本指南后，你应该已经：

- ✅ 成功启动所有三个服务
- ✅ 验证了完整的数据流
- ✅ 测试了所有核心功能
- ✅ 确认了数据持久化
- ✅ 验证了实时通信

恭喜！你的 CodeSage 项目前后端集成已经完成并正常工作！🎉
