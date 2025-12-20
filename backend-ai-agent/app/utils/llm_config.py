"""LLM配置管理"""

import os
from typing import Optional, Dict, Any
from pydantic import BaseModel, Field

class LLMConfig(BaseModel):
    """LLM配置"""
    
    # 提供商配置
    provider: str = Field(default="ollama", env="LLM_PROVIDER")  # ollama, openai, anthropic, deepseek, kimi
    
    # Ollama配置
    ollama_host: str = Field(default="http://localhost:11434", env="OLLAMA_HOST")
    ollama_model: str = Field(default="llama3.2", env="OLLAMA_MODEL")
    
    # OpenAI配置
    openai_api_key: Optional[str] = Field(default=None, env="OPENAI_API_KEY")
    openai_model: str = Field(default="gpt-3.5-turbo", env="OPENAI_MODEL")
    openai_base_url: Optional[str] = Field(default=None, env="OPENAI_BASE_URL")
    
    # Anthropic配置
    anthropic_api_key: Optional[str] = Field(default=None, env="ANTHROPIC_API_KEY")
    anthropic_model: str = Field(default="claude-3-haiku-20240307", env="ANTHROPIC_MODEL")
    
    # DeepSeek配置
    deepseek_api_key: Optional[str] = Field(default=None, env="DEEPSEEK_API_KEY")
    deepseek_model: str = Field(default="deepseek-chat", env="DEEPSEEK_MODEL")
    deepseek_base_url: str = Field(default="https://api.deepseek.com", env="DEEPSEEK_BASE_URL")
    
    # Kimi配置
    kimi_api_key: Optional[str] = Field(default=None, env="KIMI_API_KEY")
    kimi_model: str = Field(default="moonshot-v1-8k", env="KIMI_MODEL")
    kimi_base_url: str = Field(default="https://api.moonshot.cn/v1", env="KIMI_BASE_URL")
    
    # 通用配置
    max_tokens: int = Field(default=2000, env="LLM_MAX_TOKENS")
    temperature: float = Field(default=0.1, env="LLM_TEMPERATURE")
    timeout: int = Field(default=30, env="LLM_TIMEOUT")
    retry_count: int = Field(default=3, env="LLM_RETRY_COUNT")

def get_llm_config() -> LLMConfig:
    """获取LLM配置"""
    return LLMConfig()

def get_available_providers() -> Dict[str, Dict[str, Any]]:
    """获取可用的LLM提供商"""
    return {
        "ollama": {
            "name": "Ollama (本地)",
            "description": "本地运行的开源模型",
            "requires_api_key": False,
            "models": ["llama3.2", "llama2", "codellama", "mistral"]
        },
        "openai": {
            "name": "OpenAI",
            "description": "GPT系列模型",
            "requires_api_key": True,
            "models": ["gpt-3.5-turbo", "gpt-4", "gpt-4-turbo"]
        },
        "anthropic": {
            "name": "Anthropic",
            "description": "Claude系列模型",
            "requires_api_key": True,
            "models": ["claude-3-haiku-20240307", "claude-3-sonnet-20240229"]
        },
        "deepseek": {
            "name": "DeepSeek",
            "description": "DeepSeek系列模型",
            "requires_api_key": True,
            "models": ["deepseek-chat", "deepseek-coder"]
        },
        "kimi": {
            "name": "Kimi (月之暗面)",
            "description": "Moonshot AI 系列模型，中文表现优秀",
            "requires_api_key": True,
            "models": ["moonshot-v1-8k", "moonshot-v1-32k", "moonshot-v1-128k"]
        }
    }