"""改进的LLM服务，支持多种提供商"""

import asyncio
import json
import logging
from typing import Dict, Any, Optional, List
import httpx
import openai
from anthropic import Anthropic
from app.utils.llm_config import get_llm_config, get_available_providers

logger = logging.getLogger(__name__)

class LLMService:
    """改进的LLM服务，支持多种提供商"""
    
    def __init__(self):
        self.config = get_llm_config()
        self.providers = get_available_providers()
        
    async def analyze_code(self, code: str, language: str, analysis_type: str = "general") -> Dict[str, Any]:
        """分析代码"""
        try:
            prompt = self._build_analysis_prompt(code, language, analysis_type)
            response = await self._call_llm(prompt)
            return self._parse_analysis_response(response, analysis_type)
        except Exception as e:
            logger.error(f"代码分析失败: {e}")
            return self._get_fallback_analysis(code, language, analysis_type)
    
    async def suggest_improvements(self, code: str, language: str, issues: List[str]) -> Dict[str, Any]:
        """建议改进"""
        try:
            prompt = self._build_improvement_prompt(code, language, issues)
            response = await self._call_llm(prompt)
            return self._parse_improvement_response(response)
        except Exception as e:
            logger.error(f"改进建议失败: {e}")
            return self._get_fallback_improvements(code, language, issues)
    
    async def generate_tests(self, code: str, language: str, framework: str = "pytest") -> Dict[str, Any]:
        """生成测试"""
        try:
            prompt = self._build_test_prompt(code, language, framework)
            response = await self._call_llm(prompt)
            return self._parse_test_response(response, framework)
        except Exception as e:
            logger.error(f"测试生成失败: {e}")
            return self._get_fallback_tests(code, language, framework)
    
    async def convert_code(self, code: str, from_version: str, to_version: str, conversion_type: str) -> Dict[str, Any]:
        """转换代码"""
        try:
            prompt = self._build_conversion_prompt(code, from_version, to_version, conversion_type)
            response = await self._call_llm(prompt)
            return self._parse_conversion_response(response, conversion_type)
        except Exception as e:
            logger.error(f"代码转换失败: {e}")
            return self._get_fallback_conversion(code, from_version, to_version, conversion_type)
    
    async def _call_llm(self, prompt: str) -> str:
        """调用LLM API"""
        if self.config.provider == "ollama":
            return await self._call_ollama(prompt)
        elif self.config.provider == "openai":
            return await self._call_openai(prompt)
        elif self.config.provider == "anthropic":
            return await self._call_anthropic(prompt)
        elif self.config.provider == "deepseek":
            return await self._call_deepseek(prompt)
        elif self.config.provider == "kimi":
            return await self._call_kimi(prompt)
        else:
            raise ValueError(f"不支持的LLM提供商: {self.config.provider}")
    
    async def _call_ollama(self, prompt: str) -> str:
        """调用Ollama API"""
        url = f"{self.config.ollama_host}/api/generate"
        
        async with httpx.AsyncClient(timeout=self.config.timeout) as client:
            response = await client.post(url, json={
                "model": self.config.ollama_model,
                "prompt": prompt,
                "stream": False,
                "options": {
                    "temperature": self.config.temperature,
                    "num_predict": self.config.max_tokens
                }
            })
            
            if response.status_code != 200:
                raise Exception(f"Ollama API错误: {response.status_code}")
            
            result = response.json()
            return result.get("response", "")
    
    async def _call_openai(self, prompt: str) -> str:
        """调用OpenAI API"""
        if not self.config.openai_api_key:
            raise ValueError("OpenAI API密钥未配置")
        
        client = openai.AsyncOpenAI(
            api_key=self.config.openai_api_key,
            base_url=self.config.openai_base_url
        )
        
        response = await client.chat.completions.create(
            model=self.config.openai_model,
            messages=[
                {"role": "system", "content": "你是一个专业的代码分析和转换助手。"},
                {"role": "user", "content": prompt}
            ],
            temperature=self.config.temperature,
            max_tokens=self.config.max_tokens
        )
        
        return response.choices[0].message.content or ""
    
    async def _call_anthropic(self, prompt: str) -> str:
        """调用Anthropic API"""
        if not self.config.anthropic_api_key:
            raise ValueError("Anthropic API密钥未配置")
        
        client = Anthropic(api_key=self.config.anthropic_api_key)
        
        response = await client.messages.create(
            model=self.config.anthropic_model,
            max_tokens=self.config.max_tokens,
            temperature=self.config.temperature,
            messages=[
                {"role": "user", "content": prompt}
            ]
        )
        
        return response.content[0].text if response.content else ""
    
    async def _call_deepseek(self, prompt: str) -> str:
        """调用DeepSeek API"""
        if not self.config.deepseek_api_key:
            raise ValueError("DeepSeek API密钥未配置")
        
        url = f"{self.config.deepseek_base_url}/chat/completions"
        
        async with httpx.AsyncClient(timeout=self.config.timeout) as client:
            response = await client.post(url, json={
                "model": self.config.deepseek_model,
                "messages": [
                    {"role": "system", "content": "你是一个专业的代码分析和转换助手。"},
                    {"role": "user", "content": prompt}
                ],
                "temperature": self.config.temperature,
                "max_tokens": self.config.max_tokens
            }, headers={
                "Authorization": f"Bearer {self.config.deepseek_api_key}",
                "Content-Type": "application/json"
            })
            
            if response.status_code != 200:
                raise Exception(f"DeepSeek API错误: {response.status_code}")
            
            result = response.json()
            return result["choices"][0]["message"]["content"] if result.get("choices") else ""
    
    async def _call_kimi(self, prompt: str) -> str:
        """调用Kimi (Moonshot) API"""
        if not self.config.kimi_api_key:
            raise ValueError("Kimi API密钥未配置")
        
        url = f"{self.config.kimi_base_url}/chat/completions"
        
        async with httpx.AsyncClient(timeout=self.config.timeout) as client:
            response = await client.post(url, json={
                "model": self.config.kimi_model,
                "messages": [
                    {"role": "system", "content": "你是一个专业的代码分析和转换助手。请用中文或英文回复，根据用户输入的语言。"},
                    {"role": "user", "content": prompt}
                ],
                "temperature": self.config.temperature,
                "max_tokens": self.config.max_tokens
            }, headers={
                "Authorization": f"Bearer {self.config.kimi_api_key}",
                "Content-Type": "application/json"
            })
            
            if response.status_code != 200:
                raise Exception(f"Kimi API错误: {response.status_code}")
            
            result = response.json()
            return result["choices"][0]["message"]["content"] if result.get("choices") else ""
    
    def _build_analysis_prompt(self, code: str, language: str, analysis_type: str) -> str:
        """构建分析提示词"""
        base_prompt = f"""请分析以下{language}代码，提供详细的分析报告：

代码：
```{language}
{code}
```

请提供以下分析：
1. 代码质量评估
2. 潜在问题和改进建议
3. 性能优化建议
4. 安全性分析
5. 代码复杂度评估

请以JSON格式返回结果，包含以下字段：
- quality_score: 代码质量评分 (0-100)
- issues: 问题列表，每个问题包含description和severity
- recommendations: 改进建议列表
- complexity_score: 复杂度评分 (0-100)
- security_issues: 安全问题列表
"""
        
        if analysis_type == "python2_migration":
            base_prompt += "\n特别注意Python 2到Python 3的迁移问题。"
        
        return base_prompt
    
    def _build_improvement_prompt(self, code: str, language: str, issues: List[str]) -> str:
        """构建改进提示词"""
        return f"""请改进以下{language}代码，解决指定的问题：

代码：
```{language}
{code}
```

需要解决的问题：
{chr(10).join(f"- {issue}" for issue in issues)}

请提供：
1. 改进后的代码
2. 改进说明
3. 性能影响评估

请以JSON格式返回结果。"""
    
    def _build_test_prompt(self, code: str, language: str, framework: str) -> str:
        """构建测试生成提示词"""
        return f"""请为以下{language}代码生成完整的单元测试：

代码：
```{language}
{code}
```

测试框架：{framework}

请生成：
1. 基础功能测试
2. 边界条件测试
3. 异常处理测试
4. 性能测试（如适用）

请以JSON格式返回结果，包含：
- test_code: 生成的测试代码
- test_cases: 测试用例说明
- coverage_estimate: 覆盖率估计 (0-100)"""
    
    def _build_conversion_prompt(self, code: str, from_version: str, to_version: str, conversion_type: str) -> str:
        """构建转换提示词"""
        return f"""请将以下代码从{from_version}转换到{to_version}：

代码：
```
{code}
```

转换类型：{conversion_type}

请提供：
1. 转换后的代码
2. 转换说明
3. 注意事项
4. 测试建议

请以JSON格式返回结果，包含：
- converted_code: 转换后的代码
- changes: 变更说明
- warnings: 注意事项
- testing_notes: 测试建议"""
    
    def _parse_analysis_response(self, response: str, analysis_type: str) -> Dict[str, Any]:
        """解析分析响应"""
        try:
            # 尝试解析JSON响应
            if "{" in response and "}" in response:
                start = response.find("{")
                end = response.rfind("}") + 1
                json_str = response[start:end]
                return json.loads(json_str)
        except:
            pass
        
        # 回退到默认解析
        return self._get_fallback_analysis("", "", analysis_type)
    
    def _parse_improvement_response(self, response: str) -> Dict[str, Any]:
        """解析改进响应"""
        try:
            if "{" in response and "}" in response:
                start = response.find("{")
                end = response.rfind("}") + 1
                json_str = response[start:end]
                return json.loads(json_str)
        except:
            pass
        
        return {"improved_code": response, "explanation": "AI生成的改进代码"}
    
    def _parse_test_response(self, response: str, framework: str) -> Dict[str, Any]:
        """解析测试响应"""
        try:
            if "{" in response and "}" in response:
                start = response.find("{")
                end = response.rfind("}") + 1
                json_str = response[start:end]
                return json.loads(json_str)
        except:
            pass
        
        return {
            "test_code": response,
            "test_cases": ["基础功能测试", "边界条件测试"],
            "coverage_estimate": 70
        }
    
    def _parse_conversion_response(self, response: str, conversion_type: str) -> Dict[str, Any]:
        """解析转换响应"""
        try:
            if "{" in response and "}" in response:
                start = response.find("{")
                end = response.rfind("}") + 1
                json_str = response[start:end]
                return json.loads(json_str)
        except:
            pass
        
        return {
            "converted_code": response,
            "changes": ["AI生成的转换"],
            "warnings": ["请仔细检查转换结果"],
            "testing_notes": ["运行测试确保功能正常"]
        }
    
    def _get_fallback_analysis(self, code: str, language: str, analysis_type: str) -> Dict[str, Any]:
        """获取回退分析结果"""
        return {
            "quality_score": 75,
            "issues": [
                {"description": "需要进一步分析", "severity": "medium"},
                {"description": "建议使用更详细的分析工具", "severity": "low"}
            ],
            "recommendations": ["使用专业的代码分析工具", "添加单元测试"],
            "complexity_score": 50,
            "security_issues": []
        }
    
    def _get_fallback_improvements(self, code: str, language: str, issues: List[str]) -> Dict[str, Any]:
        """获取回退改进建议"""
        return {
            "improved_code": code,
            "explanation": "由于LLM服务不可用，返回原始代码。请手动检查并改进。",
            "performance_impact": "neutral"
        }
    
    def _get_fallback_tests(self, code: str, language: str, framework: str) -> Dict[str, Any]:
        """获取回退测试代码"""
        return {
            "test_code": f"# {framework} 测试框架\n# 由于LLM服务不可用，请手动编写测试",
            "test_cases": ["手动编写基础测试"],
            "coverage_estimate": 0
        }
    
    def _get_fallback_conversion(self, code: str, from_version: str, to_version: str, conversion_type: str) -> Dict[str, Any]:
        """获取回退转换结果"""
        return {
            "converted_code": f"# 从 {from_version} 到 {to_version} 的转换\n# 由于LLM服务不可用，请手动转换\n{code}",
            "changes": ["手动转换所需"],
            "warnings": ["请仔细检查手动转换结果"],
            "testing_notes": ["全面测试转换后的代码"]
        }