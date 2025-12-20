"""Code analysis endpoints."""

import uuid
from typing import Dict, Any, List
from fastapi import APIRouter, HTTPException

from app.utils.logger import get_logger
from app.models.schemas import (
    CodeAnalysisRequest, 
    CodeAnalysisResponse,
    CodeLanguage
)
from app.services.code_analyzer import CodeAnalyzer
from app.services.llm_service import LLMService

router = APIRouter()
logger = get_logger(__name__)

# Initialize services
code_analyzer = CodeAnalyzer()
llm_service = LLMService()

@router.post("/analyze", response_model=CodeAnalysisResponse)
async def analyze_code(request: CodeAnalysisRequest):
    """Analyze source code for complexity, dependencies, and issues.
    
    This endpoint performs comprehensive code analysis including:
    - Syntax validation
    - Complexity analysis
    - Dependency detection
    - Security issue identification
    - Business logic inference
    """
    
    analysis_id = str(uuid.uuid4())
    logger.info(f"Starting code analysis {analysis_id} for {request.language.value}")
    
    try:
        # Perform basic code analysis
        analysis_result = await code_analyzer.analyze(
            code=request.code,
            language=request.language,
            filename=request.filename,
            context=request.context
        )
        
        # Use LLM for business logic inference and recommendations
        llm_analysis = await llm_service.analyze_business_logic(
            code=request.code,
            language=request.language,
            context=request.context
        )
        
        # Combine results
        response = CodeAnalysisResponse(
            analysis_id=analysis_id,
            language=request.language,
            complexity_score=analysis_result.get("complexity_score", 0.0),
            dependencies=analysis_result.get("dependencies", []),
            security_issues=analysis_result.get("security_issues", []),
            compatibility_issues=analysis_result.get("compatibility_issues", []),
            business_logic_summary=llm_analysis.get("business_logic_summary"),
            recommendations=llm_analysis.get("recommendations", []),
        )
        
        logger.info(f"Code analysis {analysis_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Code analysis {analysis_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Analysis failed: {str(e)}")

@router.post("/analyze/python2", response_model=CodeAnalysisResponse)
async def analyze_python2_code(request: CodeAnalysisRequest):
    """Specialized analysis for Python 2 code.
    
    This endpoint specifically analyzes Python 2 code for:
    - Python 2 to 3 compatibility issues
    - Deprecated syntax and libraries
    - Modernization opportunities
    """
    
    # Ensure language is Python
    if request.language != CodeLanguage.PYTHON:
        raise HTTPException(
            status_code=400, 
            detail="This endpoint is specifically for Python code analysis"
        )
    
    analysis_id = str(uuid.uuid4())
    logger.info(f"Starting Python 2 analysis {analysis_id}")
    
    try:
        # Perform Python 2 specific analysis
        py2_analysis = await code_analyzer.analyze_python2_specific(
            code=request.code,
            filename=request.filename,
            context=request.context
        )
        
        # Use LLM for Python 2 specific insights
        llm_analysis = await llm_service.analyze_python2_migration(
            code=request.code,
            context=request.context
        )
        
        response = CodeAnalysisResponse(
            analysis_id=analysis_id,
            language=CodeLanguage.PYTHON,
            complexity_score=py2_analysis.get("complexity_score", 0.0),
            dependencies=py2_analysis.get("dependencies", []),
            security_issues=py2_analysis.get("security_issues", []),
            compatibility_issues=py2_analysis.get("python3_issues", []),
            business_logic_summary=llm_analysis.get("business_logic_summary"),
            recommendations=llm_analysis.get("migration_recommendations", []),
        )
        
        logger.info(f"Python 2 analysis {analysis_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Python 2 analysis {analysis_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Python 2 analysis failed: {str(e)}")

@router.get("/analyze/{analysis_id}", response_model=CodeAnalysisResponse)
async def get_analysis_result(analysis_id: str):
    """Retrieve a previous analysis result.
    
    Note: This is a placeholder implementation. In a real system,
    you would store analysis results in a database or cache.
    """
    
    logger.info(f"Retrieving analysis result {analysis_id}")
    
    # For now, return a not found error
    # In a real implementation, this would fetch from database/cache
    raise HTTPException(
        status_code=404,
        detail=f"Analysis result {analysis_id} not found. Results are not currently persisted."
    )

@router.post("/analyze/batch", response_model=List[CodeAnalysisResponse])
async def batch_analyze_code(requests: List[CodeAnalysisRequest]):
    """Analyze multiple code snippets in batch.
    
    This endpoint allows analyzing multiple code files in a single request,
    which is more efficient for large codebases.
    """
    
    batch_id = str(uuid.uuid4())
    logger.info(f"Starting batch analysis {batch_id} with {len(requests)} items")
    
    results = []
    
    for i, request in enumerate(requests):
        try:
            # Reuse the single analysis logic
            result = await analyze_code(request)
            results.append(result)
            logger.debug(f"Batch analysis {batch_id}: completed item {i+1}/{len(requests)}")
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Batch analysis {batch_id}: failed on item {i+1}: {error_msg}")
            # Continue with other items even if one fails
            continue
    
    logger.info(f"Batch analysis {batch_id} completed: {len(results)}/{len(requests)} successful")
    
    if not results:
        raise HTTPException(
            status_code=500,
            detail="All analyses in the batch failed"
        )
    
    return results