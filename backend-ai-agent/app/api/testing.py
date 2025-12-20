"""Test generation endpoints."""

import uuid
from typing import List, Dict, Any
from fastapi import APIRouter, HTTPException

from app.utils.logger import get_logger
from app.models.schemas import (
    TestGenerationRequest,
    TestGenerationResponse,
    CodeLanguage
)
from app.services.test_generator import TestGenerator
from app.services.llm_service import LLMService

router = APIRouter()
logger = get_logger(__name__)

# Initialize services
test_generator = TestGenerator()
llm_service = LLMService()

@router.post("/generate-tests", response_model=TestGenerationResponse)
async def generate_tests(request: TestGenerationRequest):
    """Generate unit tests for the provided code.
    
    This endpoint generates comprehensive unit tests including:
    - Basic functionality tests
    - Edge case tests
    - Error handling tests
    - Integration tests (if applicable)
    """
    
    test_id = str(uuid.uuid4())
    logger.info(f"Starting test generation {test_id} for {request.language.value}")
    
    try:
        # Generate tests using the test generator
        test_result = await test_generator.generate_tests(
            code=request.code,
            language=request.language,
            test_framework=request.test_framework,
            test_type=request.test_type,
            coverage_target=request.coverage_target
        )
        
        # Use LLM for additional test suggestions and improvements
        llm_suggestions = await llm_service.suggest_additional_tests(
            code=request.code,
            language=request.language,
            generated_tests=test_result.get("generated_tests", ""),
            test_framework=test_result.get("test_framework", "")
        )
        
        # Combine results
        response = TestGenerationResponse(
            test_id=test_id,
            generated_tests=test_result.get("generated_tests", ""),
            test_framework=test_result.get("test_framework", ""),
            coverage_estimate=test_result.get("coverage_estimate", 0.0),
            test_cases=test_result.get("test_cases", []) + llm_suggestions.get("additional_test_cases", [])
        )
        
        logger.info(f"Test generation {test_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Test generation {test_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Test generation failed: {str(e)}")

@router.post("/generate-tests/python", response_model=TestGenerationResponse)
async def generate_python_tests(request: TestGenerationRequest):
    """Generate Python-specific unit tests.
    
    This endpoint generates Python unit tests with:
    - pytest framework (default)
    - unittest framework (if specified)
    - Type hints and modern Python features
    - Mock objects and fixtures
    """
    
    # Ensure language is Python
    if request.language != CodeLanguage.PYTHON:
        raise HTTPException(
            status_code=400,
            detail="This endpoint is specifically for Python test generation"
        )
    
    # Default to pytest if no framework specified
    if not request.test_framework:
        request.test_framework = "pytest"
    
    test_id = str(uuid.uuid4())
    logger.info(f"Starting Python test generation {test_id} with {request.test_framework}")
    
    try:
        # Generate Python-specific tests
        python_tests = await test_generator.generate_python_tests(
            code=request.code,
            test_framework=request.test_framework,
            test_type=request.test_type,
            coverage_target=request.coverage_target
        )
        
        # Use LLM for Python-specific test improvements
        llm_improvements = await llm_service.improve_python_tests(
            code=request.code,
            generated_tests=python_tests.get("generated_tests", ""),
            test_framework=request.test_framework
        )
        
        response = TestGenerationResponse(
            test_id=test_id,
            generated_tests=python_tests.get("generated_tests", ""),
            test_framework=python_tests.get("test_framework", ""),
            coverage_estimate=python_tests.get("coverage_estimate", 0.0),
            test_cases=python_tests.get("test_cases", []) + llm_improvements.get("improved_test_cases", [])
        )
        
        logger.info(f"Python test generation {test_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Python test generation {test_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Python test generation failed: {str(e)}")

@router.post("/generate-shadow-tests", response_model=TestGenerationResponse)
async def generate_shadow_tests(request: TestGenerationRequest):
    """Generate shadow tests for comparing original and converted code.
    
    This endpoint generates shadow tests that can be used to:
    - Compare behavior between original and converted code
    - Ensure functional equivalence after conversion
    - Validate that conversions didn't introduce regressions
    """
    
    test_id = str(uuid.uuid4())
    logger.info(f"Starting shadow test generation {test_id} for {request.language.value}")
    
    try:
        # Generate shadow tests
        shadow_tests = await test_generator.generate_shadow_tests(
            code=request.code,
            language=request.language,
            test_framework=request.test_framework
        )
        
        # Use LLM for additional shadow test scenarios
        llm_scenarios = await llm_service.suggest_shadow_test_scenarios(
            code=request.code,
            language=request.language
        )
        
        response = TestGenerationResponse(
            test_id=test_id,
            generated_tests=shadow_tests.get("generated_tests", ""),
            test_framework=shadow_tests.get("test_framework", ""),
            coverage_estimate=shadow_tests.get("coverage_estimate", 0.0),
            test_cases=shadow_tests.get("test_cases", []) + llm_scenarios.get("shadow_scenarios", [])
        )
        
        logger.info(f"Shadow test generation {test_id} completed successfully")
        return response
        
    except Exception as e:
        # Escape curly braces in error message to avoid loguru format issues
        error_msg = str(e).replace('{', '{{').replace('}', '}}')
        logger.error(f"Shadow test generation {test_id} failed: {error_msg}", exc_info=True)
        raise HTTPException(status_code=500, detail=f"Shadow test generation failed: {str(e)}")

@router.get("/test-frameworks", response_model=List[Dict[str, str]])
async def get_test_frameworks():
    """Get available test frameworks for different languages.
    
    Returns a list of supported test frameworks organized by language.
    """
    
    frameworks = {
        CodeLanguage.PYTHON: [
            {"name": "pytest", "description": "Modern Python testing framework"},
            {"name": "unittest", "description": "Python standard library testing framework"},
            {"name": "nose2", "description": "Successor to nose testing framework"}
        ],
        CodeLanguage.JAVASCRIPT: [
            {"name": "jest", "description": "JavaScript testing framework"},
            {"name": "mocha", "description": "Feature-rich JavaScript test framework"},
            {"name": "jasmine", "description": "Behavior-driven development framework"}
        ],
        CodeLanguage.JAVA: [
            {"name": "junit", "description": "Java unit testing framework"},
            {"name": "testng", "description": "Java testing framework inspired by JUnit"},
            {"name": "spock", "description": "Testing and specification framework for Java"}
        ]
    }
    
    # Return frameworks for all languages or format as needed
    result = []
    for language, framework_list in frameworks.items():
        for framework in framework_list:
            result.append({
                "language": language.value,
                "framework": framework["name"],
                "description": framework["description"]
            })
    
    return result

@router.get("/test-templates/{language}", response_model=Dict[str, str])
async def get_test_templates(language: CodeLanguage):
    """Get test templates for the specified language.
    
    Returns basic test templates that can be used as starting points
    for test generation.
    """
    
    templates = {
        CodeLanguage.PYTHON: {
            "pytest": '''
import pytest
from your_module import your_function

def test_your_function():
    """Test your function with basic cases."""
    # Arrange
    input_data = "test_input"
    expected_output = "expected_output"
    
    # Act
    result = your_function(input_data)
    
    # Assert
    assert result == expected_output

def test_your_function_edge_cases():
    """Test edge cases and error conditions."""
    # Test with None input
    with pytest.raises(ValueError):
        your_function(None)
    
    # Test with empty input
    result = your_function("")
    assert result == "expected_empty_result"
''',
            "unittest": '''
import unittest
from your_module import your_function

class TestYourFunction(unittest.TestCase):
    """Test cases for your_function."""
    
    def test_basic_functionality(self):
        """Test your function with basic cases."""
        # Arrange
        input_data = "test_input"
        expected_output = "expected_output"
        
        # Act
        result = your_function(input_data)
        
        # Assert
        self.assertEqual(result, expected_output)
    
    def test_edge_cases(self):
        """Test edge cases and error conditions."""
        # Test with None input
        with self.assertRaises(ValueError):
            your_function(None)
        
        # Test with empty input
        result = your_function("")
        self.assertEqual(result, "expected_empty_result")

if __name__ == '__main__':
    unittest.main()
'''
        }
    }
    
    return templates.get(language, {})