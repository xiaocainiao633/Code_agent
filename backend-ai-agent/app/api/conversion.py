"""Code conversion endpoints."""

import uuid
from typing import List
from fastapi import APIRouter, HTTPException

from app.utils.logger import get_logger
from app.models.schemas import (
    CodeConversionRequest, 
    CodeConversionResponse,
    ConversionType,
    CodeLanguage
)
from app.services.code_converter import CodeConverter
from app.services.llm_service import LLMService

router = APIRouter()
logger = get_logger(__name__)

# Initialize services
code_converter = CodeConverter()
llm_service = LLMService()

def clean_string_list(data: list) -> List[str]:
    """Convert list items to strings, handling dict objects.
    
    Args:
        data: List that may contain strings or dicts
        
    Returns:
        List of strings
    """
    result = []
    for item in data:
        if isinstance(item, str):
            result.append(item)
        elif isinstance(item, dict):
            # If it's a dict, try to extract a 'message' field or convert to string
            if 'message' in item:
                result.append(item['message'])
            else:
                result.append(str(item))
        else:
            result.append(str(item))
    return result

@router.post("/convert", response_model=CodeConversionResponse)
async def convert_code(request: CodeConversionRequest):
    """Convert source code based on the specified conversion type.
    
    This endpoint performs code conversion including:
    - Python 2 to Python 3 conversion
    - Code modernization
    - Security fixes
    """
    
    conversion_id = str(uuid.uuid4())
    logger.info(f"Starting code conversion {conversion_id}: {request.conversion_type.value}")
    
    try:
        # Perform code conversion
        conversion_result = await code_converter.convert(
            code=request.code,
            language=request.language,
            conversion_type=request.conversion_type,
            target_version=request.target_version,
            options=request.options,
            filename=request.filename
        )
        
        # Use LLM for additional improvements and validation
        llm_validation = await llm_service.validate_conversion(
            original_code=request.code,
            converted_code=conversion_result.get("converted_code", ""),
            conversion_type=request.conversion_type,
            language=request.language
        )
        
        # Clean LLM validation results to ensure all list items are strings
        llm_warnings = clean_string_list(llm_validation.get("warnings", []))
        llm_errors = clean_string_list(llm_validation.get("errors", []))
        llm_compatibility = clean_string_list(llm_validation.get("compatibility_notes", []))
        llm_tests = clean_string_list(llm_validation.get("test_suggestions", []))
        
        # Combine results
        response = CodeConversionResponse(
            conversion_id=conversion_id,
            original_code=request.code,
            converted_code=conversion_result.get("converted_code", ""),
            language=request.language,
            conversion_type=request.conversion_type,
            changes_made=conversion_result.get("changes_made", []),
            warnings=conversion_result.get("warnings", []) + llm_warnings,
            errors=conversion_result.get("errors", []) + llm_errors,
            compatibility_notes=llm_compatibility,
            test_suggestions=llm_tests
        )
        
        logger.info(f"Code conversion {conversion_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Code conversion {conversion_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Conversion failed: {str(e)}")

@router.post("/convert/python2-to-3", response_model=CodeConversionResponse)
async def convert_python2_to_python3(request: CodeConversionRequest):
    """Specialized Python 2 to Python 3 conversion.
    
    This endpoint specifically handles Python 2 to Python 3 conversion with:
    - Automatic syntax updates
    - Library migration
    - Compatibility improvements
    - Best practices application
    """
    
    # Validate request
    if request.language != CodeLanguage.PYTHON:
        raise HTTPException(
            status_code=400,
            detail="This endpoint is specifically for Python code conversion"
        )
    
    if request.conversion_type != ConversionType.PYTHON_2_TO_3:
        request.conversion_type = ConversionType.PYTHON_2_TO_3
    
    conversion_id = str(uuid.uuid4())
    logger.info(f"Starting Python 2 to 3 conversion {conversion_id}")
    
    try:
        # Perform Python 2 to 3 specific conversion
        py2_conversion = await code_converter.convert_python2_to_python3(
            code=request.code,
            filename=request.filename,
            options=request.options
        )
        
        # Use LLM for Python 3 best practices and validation
        llm_validation = await llm_service.validate_python3_conversion(
            original_code=request.code,
            converted_code=py2_conversion.get("converted_code", "")
        )
        
        # Clean LLM validation results
        llm_warnings = clean_string_list(llm_validation.get("warnings", []))
        llm_errors = clean_string_list(llm_validation.get("errors", []))
        llm_compatibility = clean_string_list(llm_validation.get("compatibility_notes", []))
        llm_tests = clean_string_list(llm_validation.get("test_suggestions", []))
        
        response = CodeConversionResponse(
            conversion_id=conversion_id,
            original_code=request.code,
            converted_code=py2_conversion.get("converted_code", ""),
            language=CodeLanguage.PYTHON,
            conversion_type=ConversionType.PYTHON_2_TO_3,
            changes_made=py2_conversion.get("changes_made", []),
            warnings=py2_conversion.get("warnings", []) + llm_warnings,
            errors=py2_conversion.get("errors", []) + llm_errors,
            compatibility_notes=llm_compatibility,
            test_suggestions=llm_tests
        )
        
        logger.info(f"Python 2 to 3 conversion {conversion_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Python 2 to 3 conversion {conversion_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Python 2 to 3 conversion failed: {str(e)}")

@router.post("/convert/modernize", response_model=CodeConversionResponse)
async def modernize_code(request: CodeConversionRequest):
    """Modernize code to use current best practices.
    
    This endpoint modernizes code by:
    - Updating to modern syntax
    - Replacing deprecated functions
    - Improving code structure
    - Adding type hints (for Python)
    """
    
    # Set conversion type to modernization
    request.conversion_type = ConversionType.MODERNIZATION
    
    conversion_id = str(uuid.uuid4())
    logger.info(f"Starting code modernization {conversion_id} for {request.language.value}")
    
    try:
        # Perform modernization
        modernization_result = await code_converter.modernize(
            code=request.code,
            language=request.language,
            filename=request.filename,
            options=request.options
        )
        
        # Use LLM for additional modernization suggestions
        llm_suggestions = await llm_service.suggest_modernizations(
            code=request.code,
            language=request.language
        )
        
        # Clean LLM suggestions results
        llm_changes = clean_string_list(llm_suggestions.get("suggested_changes", []))
        llm_warnings = clean_string_list(llm_suggestions.get("warnings", []))
        llm_errors = clean_string_list(llm_suggestions.get("errors", []))
        llm_compatibility = clean_string_list(llm_suggestions.get("compatibility_notes", []))
        llm_tests = clean_string_list(llm_suggestions.get("test_suggestions", []))
        
        response = CodeConversionResponse(
            conversion_id=conversion_id,
            original_code=request.code,
            converted_code=modernization_result.get("modernized_code", ""),
            language=request.language,
            conversion_type=ConversionType.MODERNIZATION,
            changes_made=modernization_result.get("changes_made", []) + llm_changes,
            warnings=modernization_result.get("warnings", []) + llm_warnings,
            errors=modernization_result.get("errors", []) + llm_errors,
            compatibility_notes=llm_compatibility,
            test_suggestions=llm_tests
        )
        
        logger.info(f"Code modernization {conversion_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Code modernization {conversion_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Modernization failed: {str(e)}")

@router.get("/convert/{conversion_id}", response_model=CodeConversionResponse)
async def get_conversion_result(conversion_id: str):
    """Retrieve a previous conversion result.
    
    Note: This is a placeholder implementation. In a real system,
    you would store conversion results in a database or cache.
    """
    
    logger.info(f"Retrieving conversion result {conversion_id}")
    
    # For now, return a not found error
    # In a real implementation, this would fetch from database/cache
    raise HTTPException(
        status_code=404,
        detail=f"Conversion result {conversion_id} not found. Results are not currently persisted."
    )

@router.post("/convert/preview", response_model=CodeConversionResponse)
async def preview_conversion(request: CodeConversionRequest):
    """Preview code conversion without full processing.
    
    This endpoint provides a quick preview of what changes would be made
    without performing the full conversion process.
    """
    
    preview_id = str(uuid.uuid4())
    logger.info(f"Starting conversion preview {preview_id}: {request.conversion_type.value}")
    
    try:
        # Generate preview (simplified conversion)
        preview_result = await code_converter.preview_conversion(
            code=request.code,
            language=request.language,
            conversion_type=request.conversion_type,
            options=request.options
        )
        
        response = CodeConversionResponse(
            conversion_id=preview_id,
            original_code=request.code,
            converted_code=preview_result.get("preview_code", ""),
            language=request.language,
            conversion_type=request.conversion_type,
            changes_made=preview_result.get("preview_changes", []),
            warnings=preview_result.get("warnings", []),
            errors=preview_result.get("errors", []),
            compatibility_notes=preview_result.get("compatibility_notes", []),
            test_suggestions=[]
        )
        
        logger.info(f"Conversion preview {preview_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Conversion preview {preview_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Preview failed: {str(e)}")