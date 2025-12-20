"""Health check endpoints."""

from datetime import datetime
from typing import Dict, Any
from fastapi import APIRouter, HTTPException
from ollama import Client as OllamaClient

from app.utils.logger import get_logger
from app.models.schemas import HealthResponse
from app.utils.config_basic import get_settings

router = APIRouter()
logger = get_logger(__name__)
settings = get_settings()

@router.get("/health", response_model=HealthResponse)
async def health_check():
    """Comprehensive health check endpoint.
    
    Checks the status of all external dependencies:
    - Ollama service
    - ChromaDB (if configured)
    """
    
    dependencies: Dict[str, str] = {}
    
    # Check Ollama
    try:
        ollama_client = OllamaClient(host=settings.ollama_host)
        # Try to list models to check connectivity
        models = ollama_client.list()
        dependencies["ollama"] = "healthy"
        logger.debug(f"Ollama health check passed, found {len(models.get('models', []))} models")
    except Exception as e:
        dependencies["ollama"] = f"unhealthy: {str(e)}"
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Ollama health check failed: {error_msg}")
    
    # Check ChromaDB (placeholder - implement when ChromaDB is integrated)
    dependencies["chromadb"] = "pending"
    
    # Determine overall health
    unhealthy_deps = [k for k, v in dependencies.items() if v != "healthy" and v != "pending"]
    
    if unhealthy_deps:
        logger.warning(f"Unhealthy dependencies: {unhealthy_deps}")
        # Still return 200 but indicate issues in response
    
    return HealthResponse(
        status="healthy" if not unhealthy_deps else "degraded",
        version="0.1.0",
        timestamp=datetime.utcnow(),
        dependencies=dependencies
    )

@router.get("/health/ollama", response_model=Dict[str, str])
async def ollama_health():
    """Specific Ollama health check."""
    
    try:
        ollama_client = OllamaClient(host=settings.ollama_host)
        models = ollama_client.list()
        
        return {
            "status": "healthy",
            "host": settings.ollama_host,
            "model_count": str(len(models.get('models', []))),
            "target_model": settings.ollama_model
        }
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Ollama health check failed: {error_msg}")
        raise HTTPException(
            status_code=503,
            detail=f"Ollama service unavailable: {str(e)}"
        )

@router.get("/health/detailed", response_model=Dict[str, Any])
async def detailed_health():
    """Detailed health information."""
    
    health_info = {
        "service": {
            "name": "CodeSage AI Agent Backend",
            "version": "0.1.0",
            "timestamp": datetime.utcnow().isoformat(),
            "uptime": "Not implemented"  # Could be implemented with a startup timestamp
        },
        "configuration": {
            "api_host": settings.api_host,
            "api_port": settings.api_port,
            "ollama_host": settings.ollama_host,
            "ollama_model": settings.ollama_model,
            "chroma_collection": settings.chroma_collection_name,
            "max_code_size": settings.max_code_size
        },
        "dependencies": {}
    }
    
    # Check Ollama with more detail
    try:
        ollama_client = OllamaClient(host=settings.ollama_host)
        models = ollama_client.list()
        
        health_info["dependencies"]["ollama"] = {
            "status": "healthy",
            "host": settings.ollama_host,
            "total_models": len(models.get('models', [])),
            "target_model": settings.ollama_model,
            "models_available": [model.get('name', 'unknown') for model in models.get('models', [])]
        }
    except Exception as e:
        health_info["dependencies"]["ollama"] = {
            "status": "unhealthy",
            "error": str(e),
            "host": settings.ollama_host
        }
    
    # ChromaDB status (placeholder)
    health_info["dependencies"]["chromadb"] = {
        "status": "pending",
        "collection": settings.chroma_collection_name,
        "persist_directory": settings.chroma_persist_directory
    }
    
    return health_info