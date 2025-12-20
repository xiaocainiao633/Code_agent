"""Pydantic models for API requests and responses."""
from typing import Optional, List, Dict, Any
from pydantic import BaseModel, Field, validator
from enum import Enum
from datetime import datetime

# 代码转换服务，支持转换的语言
class CodeLanguage(str, Enum):
    """Supported programming languages."""
    PYTHON = "python"
    JAVASCRIPT = "javascript"
    JAVA = "java"
    CPP = "cpp"
    C = "c"

# 代码转换类型
class ConversionType(str, Enum):
    """Types of code conversion."""
    PYTHON_2_TO_3 = "python_2_to_3"
    MODERNIZATION = "modernization"
    SECURITY_FIX = "security_fix"

# 代码转换请求
class CodeAnalysisRequest(BaseModel):
    """Request model for code analysis."""
    
    code: str = Field(..., description="Source code to analyze")
    language: CodeLanguage = Field(..., description="Programming language")
    filename: Optional[str] = Field(None, description="Original filename")
    context: Optional[Dict[str, Any]] = Field(None, description="Additional context")
    
    @validator('code')
    def validate_code_size(cls, v):
        """Validate code size."""
        from app.utils.config_basic import get_settings
        settings = get_settings()
        if len(v.encode('utf-8')) > settings.max_code_size:
            raise ValueError(f"Code size exceeds maximum allowed size of {settings.max_code_size} bytes")
        return v

class CodeConversionRequest(BaseModel):
    """Request model for code conversion."""
    
    code: str = Field(..., description="Source code to convert")
    language: CodeLanguage = Field(..., description="Programming language")
    conversion_type: ConversionType = Field(..., description="Type of conversion")
    target_version: Optional[str] = Field(None, description="Target version")
    options: Optional[Dict[str, Any]] = Field(None, description="Conversion options")
    filename: Optional[str] = Field(None, description="Original filename")
    

class CodeAnalysisResponse(BaseModel):
    """Response model for code analysis."""
    
    analysis_id: str = Field(..., description="Unique analysis identifier")
    language: CodeLanguage = Field(..., description="Detected language")
    complexity_score: float = Field(..., description="Code complexity score (0-1)")
    dependencies: List[str] = Field(default_factory=list, description="External dependencies")
    security_issues: List[Dict[str, Any]] = Field(default_factory=list, description="Security issues found")
    compatibility_issues: List[Dict[str, Any]] = Field(default_factory=list, description="Compatibility issues")
    business_logic_summary: Optional[str] = Field(None, description="Business logic summary")
    recommendations: List[str] = Field(default_factory=list, description="Improvement recommendations")
    timestamp: datetime = Field(default_factory=datetime.utcnow, description="Analysis timestamp")

class CodeConversionResponse(BaseModel):
    """Response model for code conversion."""
    
    conversion_id: str = Field(..., description="Unique conversion identifier")
    original_code: str = Field(..., description="Original source code")
    converted_code: str = Field(..., description="Converted source code")
    language: CodeLanguage = Field(..., description="Target language")
    conversion_type: ConversionType = Field(..., description="Type of conversion")
    changes_made: List[Dict[str, Any]] = Field(default_factory=list, description="List of changes made")
    warnings: List[str] = Field(default_factory=list, description="Conversion warnings")
    errors: List[str] = Field(default_factory=list, description="Conversion errors")
    compatibility_notes: List[str] = Field(default_factory=list, description="Compatibility notes")
    test_suggestions: List[str] = Field(default_factory=list, description="Test suggestions")
    timestamp: datetime = Field(default_factory=datetime.utcnow, description="Conversion timestamp")

class TestGenerationRequest(BaseModel):
    """Request model for test generation."""
    
    code: str = Field(..., description="Source code to generate tests for")
    language: CodeLanguage = Field(..., description="Programming language")
    test_framework: Optional[str] = Field(None, description="Preferred test framework")
    test_type: Optional[str] = Field(None, description="Type of tests to generate")
    coverage_target: Optional[float] = Field(None, description="Target coverage percentage")

class TestGenerationResponse(BaseModel):
    """Response model for test generation."""
    
    test_id: str = Field(..., description="Unique test generation identifier")
    generated_tests: str = Field(..., description="Generated test code")
    test_framework: str = Field(..., description="Test framework used")
    coverage_estimate: float = Field(..., description="Estimated coverage percentage")
    test_cases: List[Dict[str, Any]] = Field(default_factory=list, description="Test case descriptions")
    timestamp: datetime = Field(default_factory=datetime.utcnow, description="Generation timestamp")

class HealthResponse(BaseModel):
    """Health check response model."""
    
    status: str = Field(..., description="Service status")
    version: str = Field(..., description="Service version")
    timestamp: datetime = Field(default_factory=datetime.utcnow, description="Health check timestamp")
    dependencies: Dict[str, str] = Field(default_factory=dict, description="Dependency statuses")

class ErrorResponse(BaseModel):
    """Error response model."""
    
    error: str = Field(..., description="Error message")
    detail: Optional[str] = Field(None, description="Error details")
    timestamp: datetime = Field(default_factory=datetime.utcnow, description="Error timestamp")
    request_id: Optional[str] = Field(None, description="Request identifier for debugging")
    
    class Config:
        json_encoders = {
            datetime: lambda v: v.isoformat()
        }