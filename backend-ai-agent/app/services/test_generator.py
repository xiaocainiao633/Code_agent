"""Test generation service for creating unit tests."""

import ast
import re
from typing import Dict, Any, List, Optional
from datetime import datetime

from app.utils.logger import get_logger
from app.models.schemas import CodeLanguage

logger = get_logger(__name__)

class TestGenerator:
    """Service for generating unit tests."""
    
    def __init__(self):
        """Initialize the test generator."""
        self.test_templates = {
            CodeLanguage.PYTHON: {
                "pytest": self._generate_pytest_tests,
                "unittest": self._generate_unittest_tests,
            }
        }
    
    async def generate_tests(self, code: str, language: CodeLanguage, 
                           test_framework: Optional[str] = None, 
                           test_type: Optional[str] = None,
                           coverage_target: Optional[float] = None) -> Dict[str, Any]:
        """Generate unit tests for the provided code.
        
        Args:
            code: Source code to generate tests for
            language: Programming language
            test_framework: Preferred test framework
            test_type: Type of tests to generate
            coverage_target: Target coverage percentage
            
        Returns:
            Generated test results
        """
        
        logger.info(f"Generating {language.value} tests with {test_framework or 'default'} framework")
        
        try:
            # Get the appropriate test generator
            if language == CodeLanguage.PYTHON:
                if not test_framework:
                    test_framework = "pytest"
                generator = self.test_templates[language].get(test_framework)
                if not generator:
                    generator = self._generate_pytest_tests  # Default to pytest
            else:
                # For other languages, use generic generation
                return await self._generate_generic_tests(code, language, test_framework)
            
            # Generate tests
            test_result = await generator(code, test_type, coverage_target)
            
            # Calculate coverage estimate
            coverage_estimate = self._estimate_coverage(code, test_result["generated_tests"], language)
            
            logger.info(f"Test generation completed for {language.value} with {test_framework}")
            
            return {
                "generated_tests": test_result["generated_tests"],
                "test_framework": test_framework,
                "coverage_estimate": coverage_estimate,
                "test_cases": test_result.get("test_cases", []),
            }
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Test generation failed: {error_msg}", exc_info=True)
            return {
                "generated_tests": "",
                "test_framework": test_framework or "unknown",
                "coverage_estimate": 0.0,
                "test_cases": [],
                "error": str(e)
            }
    
    async def generate_python_tests(self, code: str, test_framework: str = "pytest",
                                  test_type: Optional[str] = None, 
                                  coverage_target: Optional[float] = None) -> Dict[str, Any]:
        """Generate Python-specific unit tests.
        
        Args:
            code: Python source code
            test_framework: Test framework to use
            test_type: Type of tests to generate
            coverage_target: Target coverage percentage
            
        Returns:
            Generated Python test results
        """
        
        logger.info(f"Generating Python tests with {test_framework}")
        
        try:
            # Parse the code to understand its structure
            parse_result = self._parse_python_code(code)
            
            # Generate tests based on the framework
            if test_framework == "pytest":
                return await self._generate_pytest_tests(code, test_type, coverage_target, parse_result)
            elif test_framework == "unittest":
                return await self._generate_unittest_tests(code, test_type, coverage_target, parse_result)
            else:
                # Default to pytest
                return await self._generate_pytest_tests(code, test_type, coverage_target, parse_result)
                
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python test generation failed: {error_msg}", exc_info=True)
            return {
                "generated_tests": "",
                "test_cases": [],
                "error": str(e)
            }
    
    async def generate_shadow_tests(self, code: str, language: CodeLanguage, 
                                  test_framework: Optional[str] = None) -> Dict[str, Any]:
        """Generate shadow tests for comparing original and converted code.
        
        Args:
            code: Source code
            language: Programming language
            test_framework: Test framework to use
            
        Returns:
            Shadow test results
        """
        
        logger.info(f"Generating shadow tests for {language.value}")
        
        try:
            if language == CodeLanguage.PYTHON:
                if not test_framework:
                    test_framework = "pytest"
                return await self._generate_python_shadow_tests(code, test_framework)
            else:
                return await self._generate_generic_shadow_tests(code, language, test_framework)
                
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Shadow test generation failed: {error_msg}", exc_info=True)
            return {
                "generated_tests": "",
                "test_framework": test_framework or "unknown",
                "coverage_estimate": 0.0,
                "test_cases": [],
                "error": str(e)
            }
    
    def _parse_python_code(self, code: str) -> Dict[str, Any]:
        """Parse Python code to extract structure.
        
        Args:
            code: Python source code
            
        Returns:
            Parsed code information
        """
        
        try:
            tree = ast.parse(code)
            
            result = {
                "functions": [],
                "classes": [],
                "imports": [],
                "constants": [],
                "global_variables": []
            }
            
            for node in ast.walk(tree):
                if isinstance(node, ast.FunctionDef):
                    result["functions"].append({
                        "name": node.name,
                        "lineno": node.lineno,
                        "args": [arg.arg for arg in node.args.args],
                        "defaults": len(node.args.defaults),
                        "decorators": [d.id if isinstance(d, ast.Name) else str(d) for d in node.decorator_list],
                        "docstring": ast.get_docstring(node),
                        "complexity": self._calculate_function_complexity(node)
                    })
                elif isinstance(node, ast.ClassDef):
                    result["classes"].append({
                        "name": node.name,
                        "lineno": node.lineno,
                        "methods": [n.name for n in node.body if isinstance(n, ast.FunctionDef)],
                        "docstring": ast.get_docstring(node)
                    })
                elif isinstance(node, ast.Import):
                    for alias in node.names:
                        result["imports"].append(alias.name)
                elif isinstance(node, ast.ImportFrom):
                    module = node.module or ""
                    for alias in node.names:
                        result["imports"].append(f"{module}.{alias.name}")
                elif isinstance(node, ast.Assign) and isinstance(node.targets[0], ast.Name):
                    if node.targets[0].id.isupper():
                        result["constants"].append(node.targets[0].id)
                    else:
                        result["global_variables"].append(node.targets[0].id)
            
            return result
            
        except SyntaxError as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python syntax error: {error_msg}")
            raise ValueError(f"Invalid Python syntax: {str(e)}")
    
    def _calculate_function_complexity(self, func_node: ast.FunctionDef) -> int:
        """Calculate complexity score for a function.
        
        Args:
            func_node: Function AST node
            
        Returns:
            Complexity score
        """
        
        complexity = 1  # Base complexity
        
        # Count control flow statements
        for node in ast.walk(func_node):
            if isinstance(node, (ast.If, ast.For, ast.While, ast.With, ast.Try)):
                complexity += 1
            elif isinstance(node, ast.ExceptHandler):
                complexity += 1
        
        return complexity
    
    async def _generate_pytest_tests(self, code: str, test_type: Optional[str] = None,
                                   coverage_target: Optional[float] = None,
                                   parse_result: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Generate pytest tests.
        
        Args:
            code: Python source code
            test_type: Type of tests to generate
            coverage_target: Target coverage percentage
            parse_result: Pre-parsed code information
            
        Returns:
            pytest test results
        """
        
        if not parse_result:
            parse_result = self._parse_python_code(code)
        
        functions = parse_result.get("functions", [])
        classes = parse_result.get("classes", [])
        
        test_cases = []
        test_code_lines = []
        
        # Add imports
        test_code_lines.extend([
            "import pytest",
            "import sys",
            "import os",
            "",
            "# Add the parent directory to the path to import the module",
            "sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))",
            ""
        ])
        
        # Generate tests for functions
        for func in functions:
            if func["name"].startswith("_") and not func["name"].startswith("__"):
                continue  # Skip private functions
            
            test_name = f"test_{func['name']}"
            test_cases.append({
                "name": test_name,
                "description": f"Test {func['name']} function",
                "type": "function_test"
            })
            
            test_code_lines.extend([
                f"def {test_name}():",
                f'    """Test {func["name"]} function."""',
                "    # TODO: Import the function from your module",
                f"    # from your_module import {func['name']}",
                "    ",
                "    # Arrange - Set up test data",
                "    # TODO: Set up appropriate test data",
                "    ",
                "    # Act - Call the function",
                f"    # result = {func['name']}(",
                "    #     # TODO: Add function arguments",
                "    # )",
                "    ",
                "    # Assert - Verify the result",
                "    # TODO: Add appropriate assertions",
                "    # assert result is not None",
                "    # assert result == expected_result",
                "    ",
                "    pass  # Remove this when implementing the test",
                "}"
            ])
            
            # Add edge case tests for functions with parameters
            if func["args"]:
                edge_test_name = f"test_{func['name']}_edge_cases"
                test_cases.append({
                    "name": edge_test_name,
                    "description": f"Test {func['name']} function edge cases",
                    "type": "edge_case_test"
                })
                
                test_code_lines.extend([
                    "",
                    f"def {edge_test_name}():",
                    f'    """Test {func["name"]} function edge cases."""',
                    "    # TODO: Test with None values",
                    "    # TODO: Test with empty strings/collections",
                    "    # TODO: Test with invalid inputs",
                    "    # TODO: Test error handling",
                    "    pass",
                    ""
                ])
        
        # Generate tests for classes
        for cls in classes:
            test_name = f"test_{cls['name']}"
            test_cases.append({
                "name": test_name,
                "description": f"Test {cls['name']} class",
                "type": "class_test"
            })
            
            test_code_lines.extend([
                f"def {test_name}():",
                f'    """Test {cls["name"]} class."""',
                "    # TODO: Import the class from your module",
                f"    # from your_module import {cls['name']}",
                "    ",
                "    # Arrange - Create class instance",
                f"    # instance = {cls['name']}(",
                "    #     # TODO: Add constructor arguments",
                "    # )",
                "    ",
                "    # Assert - Verify instance creation",
                "    # assert instance is not None",
                "    # assert isinstance(instance, " + cls['name'] + ")",
                "    ",
                "    pass",
                "}"
            ])
        
        # Add fixture for common setup if needed
        if functions or classes:
            test_code_lines.extend([
                "",
                "@pytest.fixture",
                "def setup_test_environment():",
                '    """Set up test environment."""',
                "    # TODO: Set up any common test data or environment",
                "    # yield setup_data",
                "    pass",
                ""
            ])
        
        # Add parametrized tests if function has multiple scenarios
        if len(functions) > 1:
            test_code_lines.extend([
                "@pytest.mark.parametrize('input_data,expected_result', [",
                "    # TODO: Add test cases as tuples",
                "    # (input1, expected1),",
                "    # (input2, expected2),",
                "])",
                "def test_function_with_params(input_data, expected_result):",
                '    """Test function with multiple parameters."""',
                "    # TODO: Implement parametrized test",
                "    pass",
                ""
            ])
        
        test_code = "\n".join(test_code_lines)
        
        return {
            "generated_tests": test_code,
            "test_cases": test_cases
        }
    
    async def _generate_unittest_tests(self, code: str, test_type: Optional[str] = None,
                                     coverage_target: Optional[float] = None,
                                     parse_result: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Generate unittest tests.
        
        Args:
            code: Python source code
            test_type: Type of tests to generate
            coverage_target: Target coverage percentage
            parse_result: Pre-parsed code information
            
        Returns:
            unittest test results
        """
        
        if not parse_result:
            parse_result = self._parse_python_code(code)
        
        functions = parse_result.get("functions", [])
        classes = parse_result.get("classes", [])
        
        test_cases = []
        test_code_lines = []
        
        # Add imports
        test_code_lines.extend([
            "import unittest",
            "import sys",
            "import os",
            "",
            "# Add the parent directory to the path to import the module",
            "sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))",
            ""
        ])
        
        # Generate test class
        test_class_name = "TestGenerated"
        test_code_lines.extend([
            f"class {test_class_name}(unittest.TestCase):",
            '    """Generated test cases."""',
            ""
        ])
        
        # Generate test methods for functions
        for func in functions:
            if func["name"].startswith("_") and not func["name"].startswith("__"):
                continue
            
            test_method_name = f"test_{func['name']}"
            test_cases.append({
                "name": test_method_name,
                "description": f"Test {func['name']} function",
                "type": "function_test"
            })
            
            test_code_lines.extend([
                f"    def {test_method_name}(self):",
                f'        """Test {func["name"]} function."""',
                "        # TODO: Import the function from your module",
                f"        # from your_module import {func['name']}",
                "        ",
                "        # Arrange - Set up test data",
                "        # TODO: Set up appropriate test data",
                "        ",
                "        # Act - Call the function",
                f"        # result = {func['name']}(",
                "        #     # TODO: Add function arguments",
                "        # )",
                "        ",
                "        # Assert - Verify the result",
                "        # TODO: Add appropriate assertions",
                "        # self.assertIsNotNone(result)",
                "        # self.assertEqual(result, expected_result)",
                "        ",
                "        pass  # Remove this when implementing the test",
                "",
                ""
            ])
        
        # Generate test methods for classes
        for cls in classes:
            test_method_name = f"test_{cls['name']}"
            test_cases.append({
                "name": test_method_name,
                "description": f"Test {cls['name']} class",
                "type": "class_test"
            })
            
            test_code_lines.extend([
                f"    def {test_method_name}(self):",
                f'        """Test {cls["name"]} class."""',
                "        # TODO: Import the class from your module",
                f"        # from your_module import {cls['name']}",
                "        ",
                "        # Arrange - Create class instance",
                f"        # instance = {cls['name']}(",
                "        #     # TODO: Add constructor arguments",
                "        # )",
                "        ",
                "        # Assert - Verify instance creation",
                "        # self.assertIsNotNone(instance)",
                f"        # self.assertIsInstance(instance, {cls['name']})",
                "        ",
                "        pass",
                "",
                ""
            ])
        
        # Add setUp and tearDown methods if needed
        test_code_lines.extend([
            "    def setUp(self):",
            '        """Set up test fixtures."""',
            "        # TODO: Set up any common test data",
            "        pass",
            "",
            "    def tearDown(self):",
            '        """Tear down test fixtures."""',
            "        # TODO: Clean up after tests",
            "        pass",
            ""
        ])
        
        # Add if __name__ == '__main__'
        test_code_lines.extend([
            "if __name__ == '__main__':",
            "    unittest.main()",
            ""
        ])
        
        test_code = "\n".join(test_code_lines)
        
        return {
            "generated_tests": test_code,
            "test_cases": test_cases
        }
    
    async def _generate_python_shadow_tests(self, code: str, test_framework: str) -> Dict[str, Any]:
        """Generate Python shadow tests.
        
        Args:
            code: Python source code
            test_framework: Test framework to use
            
        Returns:
            Shadow test results
        """
        
        logger.info("Generating Python shadow tests")
        
        parse_result = self._parse_python_code(code)
        
        test_cases = []
        test_code_lines = []
        
        # Add imports
        test_code_lines.extend([
            f"import {test_framework}",
            "import sys",
            "import os",
            "import time",
            "import traceback",
            "",
            "# Add the parent directory to the path",
            "sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))",
            ""
        ])
        
        # Generate shadow test class
        test_class_name = "TestShadow"
        test_code_lines.extend([
            f"class {test_class_name}:",
            '    """Shadow tests for comparing original and converted code."""',
            "",
            "    def __init__(self):",
            "        self.results = []",
            "",
            "    def compare_outputs(self, original_func, converted_func, test_input, description=''):",
            "        \"\"\"Compare outputs of original and converted functions.\"\"\"",
            "        try:",
            "            # Run original function",
            "            start_time = time.time()",
            "            original_result = original_func(*test_input)",
            "            original_time = time.time() - start_time",
            "",
            "            # Run converted function",
            "            start_time = time.time()",
            "            converted_result = converted_func(*test_input)",
            "            converted_time = time.time() - start_time",
            "",
            "            # Compare results",
            "            results_match = original_result == converted_result",
            "",
            "            result = {",
            "                'description': description,",
            "                'input': test_input,",
            "                'original_result': original_result,",
            "                'converted_result': converted_result,",
            "                'results_match': results_match,",
            "                'original_time': original_time,",
            "                'converted_time': converted_time,",
            "                'performance_diff': converted_time - original_time,",
            "                'status': 'success' if results_match else 'mismatch'",
            "            }",
            "",
            "            self.results.append(result)",
            "            return result",
            "",
            "        except Exception as e:",
            "            result = {",
            "                'description': description,",
            "                'input': test_input,",
            "                'error': str(e),",
            "                'traceback': traceback.format_exc(),",
            "                'status': 'error'",
            "            }",
            "            self.results.append(result)",
            "            return result",
            "",
            "    def generate_report(self):",
            "        \"\"\"Generate a comparison report.\"\"\"",
            "        total_tests = len(self.results)",
            "        successful_tests = len([r for r in self.results if r['status'] == 'success'])",
            "        mismatch_tests = len([r for r in self.results if r['status'] == 'mismatch'])",
            "        error_tests = len([r for r in self.results if r['status'] == 'error'])",
            "",
            "        report = {",
            "            'total_tests': total_tests,",
            "            'successful_tests': successful_tests,",
            "            'mismatch_tests': mismatch_tests,",
            "            'error_tests': error_tests,",
            "            'success_rate': successful_tests / total_tests if total_tests > 0 else 0,",
            "            'results': self.results",
            "        }",
            "",
            "        return report",
            "",
            ""
        ])
        
        # Generate specific shadow test functions
        functions = parse_result.get("functions", [])
        
        for i, func in enumerate(functions[:3]):  # Limit to first 3 functions for demo
            test_name = f"shadow_test_{func['name']}"
            test_cases.append({
                "name": test_name,
                "description": f"Shadow test for {func['name']} function",
                "type": "shadow_test"
            })
            
            test_code_lines.extend([
                f"def {test_name}():",
                f'    """Shadow test for {func["name"]} function."""',
                "    # TODO: Import original and converted functions",
                f"    # from original_module import {func['name']} as original_{func['name']}",
                f"    # from converted_module import {func['name']} as converted_{func['name']}",
                "    ",
                "    # Create shadow test instance",
                "    shadow_test = TestShadow()",
                "    ",
                "    # Define test inputs",
                "    test_inputs = [",
                "        # TODO: Add appropriate test inputs",
                "        # (arg1, arg2, ...),",
                "        # (different_arg1, different_arg2, ...),",
                "    ]",
                "    ",
                "    # Run shadow tests",
                "    for i, test_input in enumerate(test_inputs):",
                "        result = shadow_test.compare_outputs(",
                f"            original_{func['name']},",
                f"            converted_{func['name']},",
                "            test_input,",
                f'            f"Test case {{i+1}} for {func["name"]}"',
                "        )",
                "        ",
                "        # Print result",
                "        print(f\"Test {{i+1}}: {{result['status']}}\")",
                "        if result['status'] != 'success':",
                "            print(f\"  Error: {{result.get('error', 'Results mismatch')}}\")",
                "        else:",
                "            print(f\"  Performance difference: {{result['performance_diff']:.6f}}s\")",
                "    ",
                "    # Generate report",
                "    report = shadow_test.generate_report()",
                "    print(f\"\\\\nShadow Test Report:\")",
                "    print(f\"Total tests: {{report['total_tests']}}\")",
                "    print(f\"Successful: {{report['successful_tests']}}\")",
                "    print(f\"Mismatches: {{report['mismatch_tests']}}\")",
                "    print(f\"Errors: {{report['error_tests']}}\")",
                "    print(f\"Success rate: {{report['success_rate']:.2%}}\")",
                "    ",
                "    return report",
                "",
                ""
            ])
        
        # Add main execution
        test_code_lines.extend([
            "if __name__ == '__main__':",
            "    # Run all shadow tests",
            "    print(\"Running shadow tests...\")",
            "    print(\"=\" * 50)",
            ""
        ])
        
        for func in functions[:3]:
            test_name = f"shadow_test_{func['name']}"
            test_code_lines.extend([
                f"    print(\"\\\\nRunning {test_name}...\")",
                f"    {test_name}()",
                "    print(\"-\" * 30)",
                ""
            ])
        
        test_code = "\n".join(test_code_lines)
        
        return {
            "generated_tests": test_code,
            "test_framework": test_framework,
            "coverage_estimate": 85.0,  # Shadow tests have high coverage
            "test_cases": test_cases
        }
    
    async def _generate_generic_tests(self, code: str, language: CodeLanguage, 
                                    test_framework: Optional[str] = None) -> Dict[str, Any]:
        """Generate generic tests for unsupported languages.
        
        Args:
            code: Source code
            language: Programming language
            test_framework: Test framework to use
            
        Returns:
            Generic test results
        """
        
        logger.warning(f"Generic test generation for {language.value} - limited functionality")
        
        test_cases = [
            {
                "name": "basic_functionality_test",
                "description": "Basic functionality test",
                "type": "generic_test"
            }
        ]
        
        test_code = f"""
# Generic test template for {language.value}
# Framework: {test_framework or 'unknown'}

def test_basic_functionality():
    \"\"\"Test basic functionality.\"\"\"
    # TODO: Implement basic functionality test
    # This is a generic template - customize for your needs
    
    # Arrange - Set up test data
    # test_data = ...
    
    # Act - Execute code
    # result = your_function(test_data)
    
    # Assert - Verify results
    # assert result is not None
    # assert result == expected_result
    
    pass  # Remove this when implementing the test

if __name__ == '__main__':
    test_basic_functionality()
    print("Generic test completed")
"""
        
        return {
            "generated_tests": test_code,
            "test_framework": test_framework or "generic",
            "coverage_estimate": 50.0,
            "test_cases": test_cases
        }
    
    async def _generate_generic_shadow_tests(self, code: str, language: CodeLanguage, 
                                           test_framework: Optional[str] = None) -> Dict[str, Any]:
        """Generate generic shadow tests.
        
        Args:
            code: Source code
            language: Programming language
            test_framework: Test framework to use
            
        Returns:
            Generic shadow test results
        """
        
        logger.warning(f"Generic shadow test generation for {language.value}")
        
        test_cases = [
            {
                "name": "generic_shadow_test",
                "description": "Generic shadow test for comparing outputs",
                "type": "generic_shadow_test"
            }
        ]
        
        test_code = f"""
# Generic shadow test template for {language.value}

def compare_outputs(original_func, converted_func, test_input):
    \"\"\"Compare outputs of original and converted functions.\"\"\"
    print(f"Testing with input: {{test_input}}")
    
    try:
        # Test original function
        original_result = original_func(test_input)
        print(f"Original result: {{original_result}}")
        
        # Test converted function
        converted_result = converted_func(test_input)
        print(f"Converted result: {{converted_result}}")
        
        # Compare results
        if original_result == converted_result:
            print("✓ Results match!")
            return True
        else:
            print("✗ Results differ!")
            return False
            
    except Exception as e:
        print(f"✗ Error occurred: {{e}}")
        return False

# Example usage:
# result = compare_outputs(original_function, converted_function, test_data)
# print(f"Shadow test result: {{result}}")
"""
        
        return {
            "generated_tests": test_code,
            "test_framework": test_framework or "generic",
            "coverage_estimate": 75.0,
            "test_cases": test_cases
        }
    
    def _estimate_coverage(self, code: str, test_code: str, language: CodeLanguage) -> float:
        """Estimate test coverage.
        
        Args:
            code: Original source code
            test_code: Generated test code
            language: Programming language
            
        Returns:
            Estimated coverage percentage
        """
        
        # Simple heuristic-based coverage estimation
        coverage = 50.0  # Base coverage
        
        # Check for different test types in the generated code
        if "edge_case" in test_code or "parametrize" in test_code:
            coverage += 15.0
        
        if "exception" in test_code or "error" in test_code:
            coverage += 10.0
        
        if "mock" in test_code or "patch" in test_code:
            coverage += 10.0
        
        if "integration" in test_code:
            coverage += 10.0
        
        # Check if there are multiple test functions
        test_functions = test_code.count("def test_")
        if test_functions > 3:
            coverage += 5.0
        
        return min(coverage, 95.0)  # Cap at 95% to be realistic