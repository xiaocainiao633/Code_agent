"""Main FastAPI application for CodeSage AI Agent Backend."""

import uuid
from contextlib import asynccontextmanager
from typing import AsyncGenerator

from fastapi import FastAPI, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from fastapi.staticfiles import StaticFiles

from app.utils.config_basic import get_settings, create_directories
from app.utils.logger import setup_logging, get_logger
from app.models.schemas import HealthResponse, ErrorResponse
from app.api import analysis, conversion, testing, health

# Initialize logger
logger = get_logger(__name__)

# Get settings
settings = get_settings()

@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    """Application lifespan manager."""
    # Startup
    logger.info("Starting CodeSage AI Agent Backend...")
    
    # Create necessary directories
    create_directories()
    
    # Setup logging
    setup_logging(
        log_level=settings.log_level,
        log_file=settings.log_file,
        enable_console=True
    )
    
    logger.info("Application startup complete")
    logger.info(f"API running on {settings.api_host}:{settings.api_port}")
    logger.info(f"Ollama host: {settings.ollama_host}")
    logger.info(f"Ollama model: {settings.ollama_model}")
    
    yield
    
    # Shutdown
    logger.info("Shutting down CodeSage AI Agent Backend...")
    logger.info("Application shutdown complete")

# Create FastAPI application
app = FastAPI(
    title="CodeSage AI Agent Backend",
    description="AI-powered legacy code modernization assistant",
    version="0.1.0",
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
    lifespan=lifespan
)

# Configure CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Exception handler
@app.exception_handler(HTTPException)
async def http_exception_handler(request: Request, exc: HTTPException) -> JSONResponse:
    """Handle HTTP exceptions."""
    request_id = getattr(request.state, "request_id", str(uuid.uuid4()))
    
    error_response = ErrorResponse(
        error=exc.detail,
        detail=f"HTTP {exc.status_code}",
        request_id=request_id
    )
    
    # Escape curly braces in detail to avoid format string issues
    detail_safe = str(exc.detail).replace('{', '{{').replace('}', '}}')
    logger.error(f"HTTP {exc.status_code}: {detail_safe} (request_id: {request_id})")
    
    # Use model_dump() with mode='json' to properly serialize datetime
    return JSONResponse(
        status_code=exc.status_code,
        content=error_response.model_dump(mode='json')
    )

# General exception handler
@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception) -> JSONResponse:
    """Handle general exceptions."""
    request_id = getattr(request.state, "request_id", str(uuid.uuid4()))
    
    error_response = ErrorResponse(
        error="Internal server error",
        detail=str(exc),
        request_id=request_id
    )
    
    # Use a safe string representation for logging to avoid format string issues
    error_msg = str(exc).replace('{', '{{').replace('}', '}}')
    logger.error(f"Internal error: {error_msg} (request_id: {request_id})", exc_info=True)
    
    # Use model_dump() with mode='json' to properly serialize datetime
    return JSONResponse(
        status_code=500,
        content=error_response.model_dump(mode='json')
    )

# Middleware to add request ID
@app.middleware("http")
async def add_request_id(request: Request, call_next):
    """Add request ID to all requests."""
    request_id = str(uuid.uuid4())
    request.state.request_id = request_id
    
    response = await call_next(request)
    response.headers["X-Request-ID"] = request_id
    
    return response

# Include routers
app.include_router(health.router, prefix="/api/v1", tags=["health"])
app.include_router(analysis.router, prefix="/api/v1", tags=["analysis"])
app.include_router(conversion.router, prefix="/api/v1", tags=["conversion"])
app.include_router(testing.router, prefix="/api/v1", tags=["testing"])

# Root endpoint
@app.get("/", response_model=dict)
async def root():
    """Root endpoint."""
    return {
        "message": "CodeSage AI Agent Backend",
        "version": "0.1.0",
        "docs": "/docs",
        "health": "/api/v1/health"
    }

# Health check endpoint
@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint."""
    from datetime import datetime
    
    return HealthResponse(
        status="healthy",
        version="0.1.0",
        timestamp=datetime.utcnow(),
        dependencies={
            "ollama": "pending",
            "chromadb": "pending"
        }
    )

if __name__ == "__main__":
    import uvicorn
    
    uvicorn.run(
        "main:app",
        host=settings.api_host,
        port=settings.api_port,
        reload=settings.api_debug,
        log_level=settings.log_level.lower()
    )