"""LLM service for AI-powered code analysis and generation."""

import json
from typing import Dict, Any, List, Optional
from ollama import Client as OllamaClient

from app.utils.logger import get_logger
from app.utils.config_basic import get_settings
from app.models.schemas import CodeLanguage, ConversionType

logger = get_logger(__name__)
settings = get_settings()

class LLMService:
    """Service for interacting with Large Language Models via Ollama."""
    
    def __init__(self):
        """Initialize the LLM service."""
        self.client = OllamaClient(host=settings.ollama_host)
        self.model = settings.ollama_model
        
    async def analyze_business_logic(self, code: str, language: CodeLanguage, context: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Analyze business logic in the code using LLM.
        
        Args:
            code: Source code to analyze
            language: Programming language
            context: Additional context
            
        Returns:
            Business logic analysis results
        """
        
        prompt = f"""
        Analyze the following {language.value} code for business logic and functionality.
        
        Code:
        ```{language.value}
        {code}
        ```
        
        Please provide:
        1. A summary of the business logic and purpose
        2. Key functions and their roles
        3. Data flow and processing steps
        4. Potential business rules or constraints
        
        Format your response as JSON with the following structure:
        {{
            "business_logic_summary": "Brief summary of what this code does",
            "key_functions": ["list of important functions and their purposes"],
            "data_flow": "Description of how data flows through the code",
            "business_rules": ["list of business rules or constraints"],
            "recommendations": ["list of improvement recommendations"]
        }}
        """
        
        try:
            response = self.client.generate(
                model=self.model,
                prompt=prompt,
                format="json",
                options={
                    "temperature": 0.3,
                    "top_p": 0.9
                }
            )
            
            result = json.loads(response['response'])
            logger.info(f"Business logic analysis completed for {language.value} code")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Business logic analysis failed: {error_msg}")
            return {
                "business_logic_summary": "Analysis failed",
                "key_functions": [],
                "data_flow": "",
                "business_rules": [],
                "recommendations": []
            }
    
    async def analyze_python2_migration(self, code: str, context: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Analyze Python 2 code for migration to Python 3.
        
        Args:
            code: Python 2 source code
            context: Additional context
            
        Returns:
            Python 2 to 3 migration analysis
        """
        
        prompt = f"""
        Analyze the following Python 2 code for migration to Python 3.
        
        Code:
        ```python
        {code}
        ```
        
        Please provide:
        1. A summary of the business logic
        2. Python 2 specific issues that need migration
        3. Migration complexity assessment
        4. Specific recommendations for Python 3 migration
        
        Format your response as JSON with the following structure:
        {{
            "business_logic_summary": "Brief summary of what this code does",
            "python2_issues": ["list of Python 2 specific issues"],
            "migration_complexity": "low|medium|high",
            "migration_recommendations": ["list of specific migration recommendations"],
            "critical_changes": ["list of critical changes needed"]
        }}
        """
        
        try:
            response = self.client.generate(
                model=self.model,
                prompt=prompt,
                format="json",
                options={
                    "temperature": 0.3,
                    "top_p": 0.9
                }
            )
            
            result = json.loads(response['response'])
            logger.info("Python 2 migration analysis completed")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python 2 migration analysis failed: {error_msg}")
            return {
                "business_logic_summary": "Analysis failed",
                "python2_issues": [],
                "migration_complexity": "unknown",
                "migration_recommendations": [],
                "critical_changes": []
            }
    
    async def validate_conversion(self, original_code: str, converted_code: str, 
                                conversion_type: ConversionType, language: CodeLanguage) -> Dict[str, Any]:
        """Validate code conversion using LLM.
        
        Args:
            original_code: Original source code
            converted_code: Converted source code
            conversion_type: Type of conversion performed
            language: Programming language
            
        Returns:
            Validation results
        """
        
        prompt = f"""
        Validate the following code conversion from {conversion_type.value}.
        
        Original {language.value} code:
        ```{language.value}
        {original_code}
        ```
        
        Converted {language.value} code:
        ```{language.value}
        {converted_code}
        ```
        
        Please analyze and provide:
        1. Validation of functional equivalence
        2. Any potential issues or regressions
        3. Compatibility notes
        4. Test suggestions to verify the conversion
        
        Format your response as JSON with the following structure:
        {{
            "validation_status": "valid|warnings|errors",
            "warnings": ["list of warnings"],
            "errors": ["list of errors"],
            "compatibility_notes": ["list of compatibility notes"],
            "test_suggestions": ["list of test suggestions"]
        }}
        """
        
        try:
            response = self.client.generate(
                model=self.model,
                prompt=prompt,
                format="json",
                options={
                    "temperature": 0.3,
                    "top_p": 0.9
                }
            )
            
            result = json.loads(response['response'])
            logger.info(f"Conversion validation completed for {conversion_type.value}")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Conversion validation failed: {error_msg}")
            return {
                "validation_status": "unknown",
                "warnings": ["Validation failed"],
                "errors": [],
                "compatibility_notes": [],
                "test_suggestions": []
            }
    
    async def validate_python3_conversion(self, original_code: str, converted_code: str) -> Dict[str, Any]:
        """Specifically validate Python 2 to 3 conversion.
        
        Args:
            original_code: Original Python 2 code
            converted_code: Converted Python 3 code
            
        Returns:
            Python 3 conversion validation results
        """
        
        prompt = f"""
        Validate this Python 2 to Python 3 conversion.
        
        Original Python 2 code:
        ```python
        {original_code}
        ```
        
        Converted Python 3 code:
        ```python
        {converted_code}
        ```
        
        Please analyze:
        1. Correctness of Python 2 to 3 syntax changes
        2. Proper handling of Python 2 specific features
        3. Modern Python 3 best practices
        4. Potential runtime issues
        
        Format your response as JSON with the following structure:
        {{
            "validation_status": "valid|warnings|errors",
            "python3_compliance": "excellent|good|needs_improvement",
            "warnings": ["list of warnings"],
            "errors": ["list of errors"],
            "compatibility_notes": ["list of Python 3 compatibility notes"],
            "improvement_suggestions": ["list of improvement suggestions"]
        }}
        """
        
        try:
            response = self.client.generate(
                model=self.model,
                prompt=prompt,
                format="json",
                options={
                    "temperature": 0.3,
                    "top_p": 0.9
                }
            )
            
            result = json.loads(response['response'])
            logger.info("Python 3 conversion validation completed")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python 3 conversion validation failed: {error_msg}")
            return {
                "validation_status": "unknown",
                "python3_compliance": "unknown",
                "warnings": ["Validation failed"],
                "errors": [],
                "compatibility_notes": [],
                "improvement_suggestions": []
            }
    
    async def suggest_modernizations(self, code: str, language: CodeLanguage) -> Dict[str, Any]:
        """Suggest code modernizations using LLM.
        
        Args:
            code: Source code to modernize
            language: Programming language
            
        Returns:
            Modernization suggestions
        """
        
        prompt = f"""
        Suggest modernizations for the following {language.value} code.
        
        Code:
        ```{language.value}
        {code}
        ```
        
        Please suggest:
        1. Modern syntax improvements
        2. Updated libraries or frameworks
        3. Better coding practices
        4. Performance improvements
        
        Format your response as JSON with the following structure:
        {{
            "suggested_changes": ["list of suggested changes"],
            "warnings": ["list of warnings"],
            "errors": [],
            "compatibility_notes": ["list of compatibility notes"],
            "test_suggestions": ["list of test suggestions"]
        }}
        """
        
        try:
            response = self.client.generate(
                model=self.model,
                prompt=prompt,
                format="json",
                options={
                    "temperature": 0.4,
                    "top_p": 0.9
                }
            )
            
            result = json.loads(response['response'])
            logger.info(f"Modernization suggestions generated for {language.value}")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Modernization suggestions failed: {error_msg}")
            return {
                "suggested_changes": [],
                "warnings": ["Suggestion generation failed"],
                "errors": [],
                "compatibility_notes": [],
                "test_suggestions": []
            }
    
    async def suggest_additional_tests(self, code: str, language: CodeLanguage, 
                                     generated_tests: str, test_framework: str) -> Dict[str, Any]:
        """Suggest additional test cases using LLM.
        
        Args:
            code: Source code
            language: Programming language
            generated_tests: Already generated tests
            test_framework: Test framework being used
            
        Returns:
            Additional test suggestions
        """
        
        prompt = f"""
        Suggest additional test cases for the following {language.value} code.
        
        Original code:
        ```{language.value}
        {code}
        ```
        
        Already generated tests ({test_framework}):
        ```{language.value}
        {generated_tests}
        ```
        
        Please suggest:
        1. Additional edge cases
        2. Error handling tests
        3. Integration tests
        4. Performance tests (if applicable)
        
        Format your response as JSON with the following structure:
        {{
            "additional_test_cases": [
                {{
                    "name": "test_name",
                    "description": "What this test covers",
                    "implementation": "test implementation code"
                }}
            ]
        }}
        """
        
        try:
            response = self.client.generate(
                model=self.model,
                prompt=prompt,
                format="json",
                options={
                    "temperature": 0.4,
                    "top_p": 0.9
                }
            )
            
            result = json.loads(response['response'])
            logger.info(f"Additional test suggestions generated for {language.value}")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Additional test suggestions failed: {error_msg}")
            return {
                "additional_test_cases": []
            }
    
    async def suggest_shadow_test_scenarios(self, code: str, language: CodeLanguage) -> Dict[str, Any]:
        """Suggest shadow test scenarios for comparing original and converted code.
        
        Args:
            code: Source code
            language: Programming language
            
        Returns:
            Shadow test scenarios
        """
        
        prompt = f"""
        Suggest shadow test scenarios for the following {language.value} code.
        
        Shadow tests are used to compare the behavior of original code and converted code
        to ensure functional equivalence.
        
        Code:
        ```{language.value}
        {code}
        ```
        
        Please suggest:
        1. Input scenarios that test key functionality
        2. Edge cases that might reveal conversion issues
        3. Performance comparison scenarios
        4. Integration scenarios
        
        Format your response as JSON with the following structure:
        {{
            "shadow_scenarios": [
                {{
                    "name": "scenario_name",
                    "description": "What this scenario tests",
                    "input_data": "sample input data",
                    "expected_behavior": "expected behavior description"
                }}
            ]
        }}
        """
        
        try:
            response = self.client.generate(
                model=self.model,
                prompt=prompt,
                format="json",
                options={
                    "temperature": 0.4,
                    "top_p": 0.9
                }
            )
            
            result = json.loads(response['response'])
            logger.info(f"Shadow test scenarios generated for {language.value}")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Shadow test scenarios failed: {error_msg}")
            return {
                "shadow_scenarios": []
            }
    
    async def improve_python_tests(self, code: str, generated_tests: str, test_framework: str) -> Dict[str, Any]:
        """Improve Python tests using LLM.
        
        Args:
            code: Source code
            generated_tests: Generated test code
            test_framework: Test framework being used
            
        Returns:
            Improved test suggestions
        """
        
        prompt = f"""
        Improve the following Python tests for better coverage and quality.
        
        Original code:
        ```python
        {code}
        ```
        
        Generated tests ({test_framework}):
        ```python
        {generated_tests}
        ```
        
        Please suggest improvements for:
        1. Better test organization
        2. More comprehensive assertions
        3. Proper use of {test_framework} features
        4. Test data management
        5. Mocking strategies (if applicable)
        
        Format your response as JSON with the following structure:
        {{
            "improved_test_cases": [
                {{
                    "name": "improved_test_name",
                    "description": "What this improved test covers",
                    "implementation": "improved test implementation"
                }}
            ]
        }}
        """
        
        try:
            response = self.client.generate(
                model=self.model,
                prompt=prompt,
                format="json",
                options={
                    "temperature": 0.3,
                    "top_p": 0.9
                }
            )
            
            result = json.loads(response['response'])
            logger.info(f"Python test improvements generated for {test_framework}")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python test improvements failed: {error_msg}")
            return {
                "improved_test_cases": []
            }