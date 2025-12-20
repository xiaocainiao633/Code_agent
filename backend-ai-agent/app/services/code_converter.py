"""Code conversion service for transforming source code."""

import ast
import re
from typing import Dict, Any, List, Optional
from pathlib import Path

from app.utils.logger import get_logger
from app.models.schemas import CodeLanguage, ConversionType

logger = get_logger(__name__)

class CodeConverter:
    """Service for converting source code between versions and formats."""
    
    def __init__(self):
        """Initialize the code converter."""
        self.conversion_strategies = {
            (CodeLanguage.PYTHON, ConversionType.PYTHON_2_TO_3): self.convert_python2_to_python3,
            (CodeLanguage.PYTHON, ConversionType.MODERNIZATION): self._modernize_python,
            (CodeLanguage.JAVASCRIPT, ConversionType.MODERNIZATION): self._modernize_javascript,
        }
    
    async def convert(self, code: str, language: CodeLanguage, conversion_type: ConversionType,
                     target_version: Optional[str] = None, options: Optional[Dict[str, Any]] = None,
                     filename: Optional[str] = None) -> Dict[str, Any]:
        """Convert source code based on the specified type.
        
        Args:
            code: Source code to convert
            language: Programming language
            conversion_type: Type of conversion
            target_version: Target version (if applicable)
            options: Conversion options
            filename: Original filename
            
        Returns:
            Conversion results
        """
        
        logger.info(f"Starting {conversion_type.value} conversion for {language.value}")
        
        try:
            # Get the appropriate conversion strategy
            strategy = self.conversion_strategies.get((language, conversion_type))
            
            if not strategy:
                # Use generic conversion or raise error
                return await self._generic_convert(code, language, conversion_type, options)
            
            # Apply the conversion strategy
            result = await strategy(code, filename, options)
            
            logger.info(f"{conversion_type.value} conversion completed for {language.value}")
            return result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Code conversion failed: {error_msg}", exc_info=True)
            return {
                "converted_code": code,
                "changes_made": [],
                "warnings": [f"Conversion failed: {str(e)}"],
                "errors": [str(e)]
            }
    
    async def convert_python2_to_python3(self, code: str, filename: Optional[str] = None, 
                                       options: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Convert Python 2 code to Python 3.
        
        Args:
            code: Python 2 source code
            filename: Original filename
            options: Conversion options
            
        Returns:
            Conversion results
        """
        
        logger.info("Starting Python 2 to 3 conversion")
        
        try:
            changes_made = []
            converted_code = code
            
            # Apply Python 2 to 3 conversions
            conversions = [
                self._convert_print_statements,
                self._convert_xrange_to_range,
                self._convert_has_key_to_in,
                self._convert_unicode_to_str,
                self._convert_basestring_to_str,
                self._convert_iteritems_to_items,
                self._convert_iterkeys_to_keys,
                self._convert_itervalues_to_values,
                self._convert_ne_to_not_equal,
                self._convert_exception_syntax,
                self._convert_imports,
                self._convert_string_formatting,
                self._convert_division,
                self._convert_zip_map_filter,
            ]
            
            for conversion_func in conversions:
                try:
                    converted_code, changes = conversion_func(converted_code)
                    changes_made.extend(changes)
                except Exception as e:
                    # Escape curly braces in error message to avoid loguru format issues
                    error_msg = str(e).replace('{', '{{').replace('}', '}}')
                    logger.warning(f"Conversion step {conversion_func.__name__} failed: {error_msg}")
            
            # Validate the converted code
            try:
                ast.parse(converted_code)
            except SyntaxError as e:
                # Escape curly braces in error message to avoid loguru format issues
                error_msg = str(e).replace('{', '{{').replace('}', '}}')
                logger.error(f"Converted code has syntax errors: {error_msg}")
                return {
                    "converted_code": code,
                    "changes_made": [],
                    "warnings": [],
                    "errors": [f"Conversion resulted in invalid syntax: {str(e)}"]
                }
            
            logger.info(f"Python 2 to 3 conversion completed with {len(changes_made)} changes")
            
            return {
                "converted_code": converted_code,
                "changes_made": changes_made,
                "warnings": [],
                "errors": []
            }
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python 2 to 3 conversion failed: {error_msg}", exc_info=True)
            return {
                "converted_code": code,
                "changes_made": [],
                "warnings": [],
                "errors": [str(e)]
            }
    
    def _convert_print_statements(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert Python 2 print statements to print() function calls."""
        
        changes = []
        
        # Pattern for simple print statements
        pattern = r'print\s+(.+?)(?:\s*#.*)?$'
        
        def replace_print(match):
            content = match.group(1).strip()
            # Handle print >>file, content syntax
            if content.startswith('>>'):
                file_part, content_part = content.split(',', 1)
                file_var = file_part.replace('>>', '').strip()
                result = f"print({content_part.strip()}, file={file_var})"
            else:
                result = f"print({content})"
            
            changes.append({
                "type": "print_statement",
                "line": code[:match.start()].count('\n') + 1,
                "original": match.group(0),
                "replacement": result
            })
            
            return result
        
        converted = re.sub(pattern, replace_print, code, flags=re.MULTILINE)
        return converted, changes
    
    def _convert_xrange_to_range(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert xrange() to range()."""
        
        changes = []
        pattern = r'\bxrange\s*\('
        
        def replace_xrange(match):
            changes.append({
                "type": "xrange_to_range",
                "line": code[:match.start()].count('\n') + 1,
                "original": "xrange",
                "replacement": "range"
            })
            return "range("
        
        converted = re.sub(pattern, replace_xrange, code)
        return converted, changes
    
    def _convert_has_key_to_in(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert dict.has_key() to 'in' operator."""
        
        changes = []
        pattern = r'(\w+)\.has_key\s*\(\s*([^)]+)\s*\)'
        
        def replace_has_key(match):
            dict_name = match.group(1)
            key = match.group(2)
            result = f"{key} in {dict_name}"
            
            changes.append({
                "type": "has_key_to_in",
                "line": code[:match.start()].count('\n') + 1,
                "original": match.group(0),
                "replacement": result
            })
            
            return result
        
        converted = re.sub(pattern, replace_has_key, code)
        return converted, changes
    
    def _convert_unicode_to_str(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert unicode() to str()."""
        
        changes = []
        pattern = r'\bunicode\s*\('
        
        def replace_unicode(match):
            changes.append({
                "type": "unicode_to_str",
                "line": code[:match.start()].count('\n') + 1,
                "original": "unicode",
                "replacement": "str"
            })
            return "str("
        
        converted = re.sub(pattern, replace_unicode, code)
        return converted, changes
    
    def _convert_basestring_to_str(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert basestring to str."""
        
        changes = []
        pattern = r'\bbasestring\b'
        
        def replace_basestring(match):
            changes.append({
                "type": "basestring_to_str",
                "line": code[:match.start()].count('\n') + 1,
                "original": "basestring",
                "replacement": "str"
            })
            return "str"
        
        converted = re.sub(pattern, replace_basestring, code)
        return converted, changes
    
    def _convert_iteritems_to_items(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert dict.iteritems() to dict.items()."""
        
        changes = []
        pattern = r'(\w+)\.iteritems\s*\(\s*\)'
        
        def replace_iteritems(match):
            dict_name = match.group(1)
            result = f"list({dict_name}.items())"
            
            changes.append({
                "type": "iteritems_to_items",
                "line": code[:match.start()].count('\n') + 1,
                "original": match.group(0),
                "replacement": result
            })
            
            return result
        
        converted = re.sub(pattern, replace_iteritems, code)
        return converted, changes
    
    def _convert_iterkeys_to_keys(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert dict.iterkeys() to dict.keys()."""
        
        changes = []
        pattern = r'(\w+)\.iterkeys\s*\(\s*\)'
        
        def replace_iterkeys(match):
            dict_name = match.group(1)
            result = f"list({dict_name}.keys())"
            
            changes.append({
                "type": "iterkeys_to_keys",
                "line": code[:match.start()].count('\n') + 1,
                "original": match.group(0),
                "replacement": result
            })
            
            return result
        
        converted = re.sub(pattern, replace_iterkeys, code)
        return converted, changes
    
    def _convert_itervalues_to_values(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert dict.itervalues() to dict.values()."""
        
        changes = []
        pattern = r'(\w+)\.itervalues\s*\(\s*\)'
        
        def replace_itervalues(match):
            dict_name = match.group(1)
            result = f"list({dict_name}.values())"
            
            changes.append({
                "type": "itervalues_to_values",
                "line": code[:match.start()].count('\n') + 1,
                "original": match.group(0),
                "replacement": result
            })
            
            return result
        
        converted = re.sub(pattern, replace_itervalues, code)
        return converted, changes
    
    def _convert_ne_to_not_equal(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert <> to != operator."""
        
        changes = []
        pattern = r'<>'
        
        def replace_ne(match):
            changes.append({
                "type": "ne_to_not_equal",
                "line": code[:match.start()].count('\n') + 1,
                "original": "<>",
                "replacement": "!="
            })
            return "!="
        
        converted = re.sub(pattern, replace_ne, code)
        return converted, changes
    
    def _convert_exception_syntax(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Convert exception syntax."""
        
        changes = []
        
        # Convert except Exception, e: to except Exception as e:
        pattern = r'except\s+(\w+)\s*,\s*(\w+)\s*:'
        
        def replace_except(match):
            exception_type = match.group(1)
            variable = match.group(2)
            result = f"except {exception_type} as {variable}:"
            
            changes.append({
                "type": "except_syntax",
                "line": code[:match.start()].count('\n') + 1,
                "original": match.group(0),
                "replacement": result
            })
            
            return result
        
        converted = re.sub(pattern, replace_except, code)
        return converted, changes
    
    def _convert_imports(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Update import statements."""
        
        changes = []
        
        # Add from __future__ import for Python 2 compatibility if needed
        if "from __future__" not in code:
            # Check if we made significant changes
            if any(change["type"] in ["print_statement", "division"] for change in []):
                future_import = "from __future__ import print_function, division\n"
                converted = future_import + code
                changes.append({
                    "type": "future_imports",
                    "line": 1,
                    "original": "",
                    "replacement": future_import.strip()
                })
                return converted, changes
        
        return code, changes
    
    def _convert_string_formatting(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Update string formatting."""
        
        changes = []
        
        # Convert % formatting to .format() where appropriate
        # This is a simplified implementation
        pattern = r'%s'
        
        if re.search(pattern, code):
            changes.append({
                "type": "string_formatting",
                "line": 0,
                "original": "%s formatting",
                "replacement": "Consider using .format() or f-strings"
            })
        
        return code, changes
    
    def _convert_division(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Handle division changes."""
        
        changes = []
        
        # In Python 3, / is true division, // is floor division
        # Add a warning about division behavior changes
        if re.search(r'/\s*\d', code):
            changes.append({
                "type": "division",
                "line": 0,
                "original": "division operator",
                "replacement": "Note: / is true division in Python 3"
            })
        
        return code, changes
    
    def _convert_zip_map_filter(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Handle zip, map, filter changes."""
        
        changes = []
        
        # In Python 3, these return iterators, not lists
        patterns = [r'\bzip\s*\(', r'\bmap\s*\(', r'\bfilter\s*\(']
        
        for pattern in patterns:
            if re.search(pattern, code):
                func_name = pattern.replace(r'\b', '').replace(r'\s*\(', '')
                changes.append({
                    "type": f"{func_name}_iterator",
                    "line": 0,
                    "original": f"{func_name}()",
                    "replacement": f"Note: {func_name}() returns iterator in Python 3"
                })
        
        return code, changes
    
    async def _modernize_python(self, code: str, filename: Optional[str] = None, 
                              options: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Modernize Python code.
        
        Args:
            code: Python source code
            filename: Original filename
            options: Modernization options
            
        Returns:
            Modernization results
        """
        
        logger.info("Starting Python modernization")
        
        try:
            changes_made = []
            modernized_code = code
            
            # Apply modernization transformations
            modernizations = [
                self._add_type_hints,
                self._modernize_string_formatting,
                self._modernize_collections,
                self._modernize_exceptions,
                self._suggest_f_strings,
            ]
            
            for modernization_func in modernizations:
                try:
                    modernized_code, changes = modernization_func(modernized_code)
                    changes_made.extend(changes)
                except Exception as e:
                    # Escape curly braces in error message to avoid loguru format issues
                    error_msg = str(e).replace('{', '{{').replace('}', '}}')
                    logger.warning(f"Modernization step {modernization_func.__name__} failed: {error_msg}")
            
            logger.info(f"Python modernization completed with {len(changes_made)} changes")
            
            return {
                "modernized_code": modernized_code,
                "changes_made": changes_made,
                "warnings": [],
                "errors": []
            }
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python modernization failed: {error_msg}", exc_info=True)
            return {
                "modernized_code": code,
                "changes_made": [],
                "warnings": [],
                "errors": [str(e)]
            }
    
    def _add_type_hints(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Add type hints to Python code."""
        
        changes = []
        
        # This is a simplified implementation
        # In a real system, you would use more sophisticated analysis
        if "def " in code and "->" not in code:
            changes.append({
                "type": "type_hints",
                "line": 0,
                "original": "function definitions",
                "replacement": "Consider adding type hints to function definitions"
            })
        
        return code, changes
    
    def _modernize_string_formatting(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Modernize string formatting."""
        
        changes = []
        
        # Suggest f-strings for Python 3.6+
        if '%' in code or '.format(' in code:
            changes.append({
                "type": "string_formatting",
                "line": 0,
                "original": "old string formatting",
                "replacement": "Consider using f-strings for better readability"
            })
        
        return code, changes
    
    def _modernize_collections(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Modernize collections usage."""
        
        changes = []
        
        # Suggest using collections.abc for abstract base classes
        if "from collections import" in code:
            changes.append({
                "type": "collections",
                "line": 0,
                "original": "collections imports",
                "replacement": "Consider using collections.abc for abstract base classes"
            })
        
        return code, changes
    
    def _modernize_exceptions(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Modernize exception handling."""
        
        changes = []
        
        # Suggest using specific exception types
        if "except:" in code:
            changes.append({
                "type": "exception_handling",
                "line": 0,
                "original": "bare except",
                "replacement": "Use specific exception types instead of bare except:"
            })
        
        return code, changes
    
    def _suggest_f_strings(self, code: str) -> tuple[str, List[Dict[str, Any]]]:
        """Suggest f-string usage."""
        
        changes = []
        
        # Look for string concatenation or .format() usage
        if '+' in code and '"' in code:
            changes.append({
                "type": "f_strings",
                "line": 0,
                "original": "string concatenation",
                "replacement": "Consider using f-strings for string formatting"
            })
        
        return code, changes
    
    async def _modernize_javascript(self, code: str, filename: Optional[str] = None, 
                                  options: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Modernize JavaScript code.
        
        Args:
            code: JavaScript source code
            filename: Original filename
            options: Modernization options
            
        Returns:
            Modernization results
        """
        
        logger.info("Starting JavaScript modernization")
        
        # Placeholder implementation
        changes = [{
            "type": "javascript_modernization",
            "line": 0,
            "original": "JavaScript code",
            "replacement": "Consider using ES6+ features like arrow functions, const/let, template literals"
        }]
        
        return {
            "modernized_code": code,
            "changes_made": changes,
            "warnings": [],
            "errors": []
        }
    
    async def _generic_convert(self, code: str, language: CodeLanguage, 
                             conversion_type: ConversionType, options: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Generic conversion for unsupported language/type combinations.
        
        Args:
            code: Source code
            language: Programming language
            conversion_type: Type of conversion
            options: Conversion options
            
        Returns:
            Generic conversion results
        """
        
        logger.warning(f"No specific conversion strategy for {language.value} + {conversion_type.value}")
        
        return {
            "converted_code": code,
            "changes_made": [{
                "type": "no_conversion",
                "line": 0,
                "original": "original code",
                "replacement": "No specific conversion available"
            }],
            "warnings": [f"No specific conversion strategy for {language.value} + {conversion_type.value}"],
            "errors": []
        }
    
    async def preview_conversion(self, code: str, language: CodeLanguage, 
                               conversion_type: ConversionType, options: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Generate a preview of what changes would be made.
        
        Args:
            code: Source code
            language: Programming language
            conversion_type: Type of conversion
            options: Conversion options
            
        Returns:
            Preview results
        """
        
        logger.info(f"Generating conversion preview for {conversion_type.value}")
        
        try:
            # Perform a lightweight preview
            if conversion_type == ConversionType.PYTHON_2_TO_3 and language == CodeLanguage.PYTHON:
                # Just detect issues without full conversion
                issues = self._detect_python2_issues(code, {})
                
                preview_changes = []
                for issue in issues[:5]:  # Limit preview to first 5 issues
                    preview_changes.append({
                        "type": "preview",
                        "line": issue.get("line", 0),
                        "description": issue.get("message", ""),
                        "suggestion": issue.get("suggestion", "")
                    })
                
                return {
                    "preview_code": code,
                    "preview_changes": preview_changes,
                    "warnings": [],
                    "errors": [],
                    "compatibility_notes": [f"Found {len(issues)} Python 2 compatibility issues"]
                }
            
            else:
                # Generic preview
                return {
                    "preview_code": code,
                    "preview_changes": [{
                        "type": "preview",
                        "line": 0,
                        "description": f"Preview not available for {conversion_type.value}",
                        "suggestion": "Run full conversion to see changes"
                    }],
                    "warnings": [],
                    "errors": [],
                    "compatibility_notes": []
                }
                
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Preview generation failed: {error_msg}", exc_info=True)
            return {
                "preview_code": code,
                "preview_changes": [],
                "warnings": [],
                "errors": [str(e)],
                "compatibility_notes": []
            }