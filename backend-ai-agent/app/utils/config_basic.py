"""Basic configuration management without pydantic-settings."""

import os
from pathlib import Path
from typing import Optional, List
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

class BasicSettings:
    """Basic application settings without pydantic complexity."""
    
    def __init__(self):
        # Ollama configuration
        self.ollama_host = os.getenv("OLLAMA_HOST", "http://localhost:11434")
        self.ollama_model = os.getenv("OLLAMA_MODEL", "llama3.2")
        
        # FastAPI configuration
        self.api_host = os.getenv("API_HOST", "0.0.0.0")
        self.api_port = int(os.getenv("API_PORT", "8000"))
        self.api_debug = os.getenv("API_DEBUG", "true").lower() == "true"
        
        # ChromaDB configuration
        self.chroma_persist_directory = os.getenv("CHROMA_PERSIST_DIRECTORY", "./chroma_db")
        self.chroma_collection_name = os.getenv("CHROMA_COLLECTION_NAME", "code_embeddings")
        
        # Logging configuration
        self.log_level = os.getenv("LOG_LEVEL", "INFO")
        self.log_file = os.getenv("LOG_FILE", "./logs/app.log")
        
        # Security configuration
        self.max_code_size = int(os.getenv("MAX_CODE_SIZE", "1048576"))  # 1MB
        
        # Parse allowed file extensions
        extensions_str = os.getenv("ALLOWED_FILE_EXTENSIONS", ".py,.js,.java,.cpp,.c")
        self.allowed_file_extensions = [ext.strip() for ext in extensions_str.split(",") if ext.strip()]
    
    def __repr__(self):
        return f"Settings(ollama_host={self.ollama_host}, api_port={self.api_port})"

# Global settings instance
settings = BasicSettings()

def get_settings():
    """Get application settings.
    
    Returns:
        Settings instance
    """
    return settings

def create_directories() -> None:
    """Create necessary directories for the application."""
    # Create logs directory
    if settings.log_file:
        log_path = Path(settings.log_file)
        log_path.parent.mkdir(parents=True, exist_ok=True)
    
    # Create ChromaDB directory
    chroma_path = Path(settings.chroma_persist_directory)
    chroma_path.mkdir(parents=True, exist_ok=True)