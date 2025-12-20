"""Code analysis service for parsing and analyzing source code."""

import ast
import re
from typing import Dict, Any, List, Optional
from pathlib import Path

from app.utils.logger import get_logger
from app.models.schemas import CodeLanguage

logger = get_logger(__name__)

class CodeAnalyzer:
    """Service for analyzing source code."""
    
    def __init__(self):
        """Initialize the code analyzer."""
        self.language_parsers = {
            CodeLanguage.PYTHON: self._parse_python,
            CodeLanguage.JAVASCRIPT: self._parse_javascript,
            CodeLanguage.JAVA: self._parse_java,
            CodeLanguage.CPP: self._parse_cpp,
            CodeLanguage.C: self._parse_c,
        }
    
    async def analyze(self, code: str, language: CodeLanguage, 
                     filename: Optional[str] = None, context: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Analyze source code.
        
        Args:
            code: Source code to analyze
            language: Programming language
            filename: Original filename
            context: Additional context
            
        Returns:
            Analysis results
        """
        
        logger.info(f"Analyzing {language.value} code")
        
        try:
            # Get the appropriate parser
            parser = self.language_parsers.get(language)
            if not parser:
                raise ValueError(f"Language {language.value} not supported")
            
            # Parse the code
            parse_result = parser(code, filename)
            
            # Perform general analysis
            analysis_result = {
                "complexity_score": self._calculate_complexity(parse_result),
                "dependencies": self._extract_dependencies(parse_result, language),
                "security_issues": self._detect_security_issues(parse_result, language),
                "compatibility_issues": self._detect_compatibility_issues(parse_result, language),
                "code_metrics": self._calculate_metrics(parse_result, language),
            }
            
            logger.info(f"Analysis completed for {language.value} code")
            return analysis_result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Code analysis failed: {error_msg}", exc_info=True)
            return {
                "complexity_score": 0.0,
                "dependencies": [],
                "security_issues": [],
                "compatibility_issues": [],
                "code_metrics": {},
                "error": str(e)
            }
    
    async def analyze_python2_specific(self, code: str, filename: Optional[str] = None, 
                                     context: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Analyze Python 2 code for Python 3 migration issues.
        
        Args:
            code: Python 2 source code
            filename: Original filename
            context: Additional context
            
        Returns:
            Python 2 specific analysis results
        """
        
        logger.info("Analyzing Python 2 code for migration issues")
        
        try:
            # Parse Python 2 code
            parse_result = self._parse_python(code, filename)
            
            # Detect Python 2 specific issues
            python2_issues = self._detect_python2_issues(code, parse_result)
            
            analysis_result = {
                "complexity_score": self._calculate_complexity(parse_result),
                "dependencies": self._extract_dependencies(parse_result, CodeLanguage.PYTHON),
                "security_issues": self._detect_security_issues(parse_result, CodeLanguage.PYTHON),
                "python3_issues": python2_issues,
                "code_metrics": self._calculate_metrics(parse_result, CodeLanguage.PYTHON),
            }
            
            logger.info("Python 2 analysis completed")
            return analysis_result
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python 2 analysis failed: {error_msg}", exc_info=True)
            return {
                "complexity_score": 0.0,
                "dependencies": [],
                "security_issues": [],
                "python3_issues": [],
                "code_metrics": {},
                "error": str(e)
            }
    
    def _parse_python(self, code: str, filename: Optional[str] = None) -> Dict[str, Any]:
        """Parse Python code using AST.
        
        Args:
            code: Python source code
            filename: Original filename
            
        Returns:
            Parsed AST information
        """
        
        try:
            # First try to parse as Python 3
            tree = ast.parse(code, filename=filename or '<unknown>')
            
            # Extract information from AST
            result = {
                "ast": tree,
                "functions": [],
                "classes": [],
                "imports": [],
                "variables": [],
                "complexity_indicators": []
            }
            
            # Walk through the AST
            for node in ast.walk(tree):
                if isinstance(node, ast.FunctionDef):
                    result["functions"].append({
                        "name": node.name,
                        "lineno": node.lineno,
                        "args": [arg.arg for arg in node.args.args],
                        "decorators": [d.id if isinstance(d, ast.Name) else str(d) for d in node.decorator_list]
                    })
                elif isinstance(node, ast.ClassDef):
                    result["classes"].append({
                        "name": node.name,
                        "lineno": node.lineno,
                        "methods": [n.name for n in node.body if isinstance(n, ast.FunctionDef)]
                    })
                elif isinstance(node, ast.Import):
                    for alias in node.names:
                        result["imports"].append(alias.name)
                elif isinstance(node, ast.ImportFrom):
                    module = node.module or ""
                    for alias in node.names:
                        result["imports"].append(f"{module}.{alias.name}")
                elif isinstance(node, ast.Name) and isinstance(node.ctx, ast.Store):
                    result["variables"].append(node.id)
            
            return result
            
        except SyntaxError as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python syntax error: {error_msg}")
            # For Python 2 code, try to parse with more lenient approach
            logger.warning("Python 2 syntax detected, using compatibility mode")
            return self._parse_python2_compat(code, filename)
    
    def _parse_javascript(self, code: str, filename: Optional[str] = None) -> Dict[str, Any]:
        """Parse JavaScript code (placeholder implementation).
        
        Args:
            code: JavaScript source code
            filename: Original filename
            
        Returns:
            Parsed information
        """
        
        # This would use tree-sitter or similar for real JavaScript parsing
        logger.warning("JavaScript parsing not fully implemented, using basic regex analysis")
        
        result = {
            "functions": [],
            "classes": [],
            "imports": [],
            "variables": [],
            "complexity_indicators": []
        }
        
        # Basic regex-based analysis
        function_pattern = r'function\s+(\w+)\s*\('
        arrow_pattern = r'const\s+(\w+)\s*=\s*\([^)]*\)\s*=>'
        class_pattern = r'class\s+(\w+)'
        
        result["functions"] = re.findall(function_pattern, code)
        result["functions"].extend(re.findall(arrow_pattern, code))
        result["classes"] = re.findall(class_pattern, code)
        
        return result
    
    def _parse_java(self, code: str, filename: Optional[str] = None) -> Dict[str, Any]:
        """Parse Java code (placeholder implementation)."""
        logger.warning("Java parsing not fully implemented")
        return {"functions": [], "classes": [], "imports": [], "variables": [], "complexity_indicators": []}
    
    def _parse_cpp(self, code: str, filename: Optional[str] = None) -> Dict[str, Any]:
        """Parse C++ code (placeholder implementation)."""
        logger.warning("C++ parsing not fully implemented")
        return {"functions": [], "classes": [], "imports": [], "variables": [], "complexity_indicators": []}
    
    def _parse_c(self, code: str, filename: Optional[str] = None) -> Dict[str, Any]:
        """Parse C code (placeholder implementation)."""
        logger.warning("C parsing not fully implemented")
        return {"functions": [], "classes": [], "imports": [], "variables": [], "complexity_indicators": []}
    
    def _calculate_complexity(self, parse_result: Dict[str, Any]) -> float:
        """Calculate code complexity score (0-1).
        
        Args:
            parse_result: Parsed code information
            
        Returns:
            Complexity score
        """
        
        complexity_score = 0.0
        
        # Base complexity on function count
        functions = parse_result.get("functions", [])
        complexity_score += min(len(functions) * 0.1, 0.3)
        
        # Complexity based on nesting depth (for Python AST)
        if "ast" in parse_result and parse_result["ast"]:
            max_depth = self._calculate_max_nesting_depth(parse_result["ast"])
            complexity_score += min(max_depth * 0.05, 0.3)
        
        # Complexity based on import count
        imports = parse_result.get("imports", [])
        complexity_score += min(len(imports) * 0.02, 0.2)
        
        # Complexity based on class count
        classes = parse_result.get("classes", [])
        complexity_score += min(len(classes) * 0.05, 0.2)
        
        return min(complexity_score, 1.0)
    
    def _calculate_max_nesting_depth(self, ast_node) -> int:
        """Calculate maximum nesting depth in AST.
        
        Args:
            ast_node: AST node
            
        Returns:
            Maximum nesting depth
        """
        
        max_depth = 0
        
        for node in ast.walk(ast_node):
            if isinstance(node, (ast.If, ast.For, ast.While, ast.With, ast.Try)):
                depth = self._get_nesting_depth(node)
                max_depth = max(max_depth, depth)
        
        return max_depth
    
    def _get_nesting_depth(self, node) -> int:
        """Get nesting depth of a specific node.
        
        Args:
            node: AST node
            
        Returns:
            Nesting depth
        """
        
        depth = 0
        current = node
        
        while hasattr(current, 'parent') and current.parent:
            if isinstance(current.parent, (ast.If, ast.For, ast.While, ast.With, ast.Try)):
                depth += 1
            current = current.parent
        
        return depth
    
    def _extract_dependencies(self, parse_result: Dict[str, Any], language: CodeLanguage) -> List[str]:
        """Extract external dependencies.
        
        Args:
            parse_result: Parsed code information
            language: Programming language
            
        Returns:
            List of dependencies
        """
        
        imports = parse_result.get("imports", [])
        
        # Filter out standard library modules (basic implementation)
        if language == CodeLanguage.PYTHON:
            stdlib_modules = {
                'os', 'sys', 'json', 'datetime', 'math', 'random', 're', 
                'collections', 'itertools', 'functools', 'operator', 'string',
                'urllib', 'http', 'subprocess', 'pathlib', 'typing'
            }
            return [imp for imp in imports if imp.split('.')[0] not in stdlib_modules]
        
        return imports
    
    def _detect_security_issues(self, parse_result: Dict[str, Any], language: CodeLanguage) -> List[Dict[str, Any]]:
        """Detect potential security issues.
        
        Args:
            parse_result: Parsed code information
            language: Programming language
            
        Returns:
            List of security issues
        """
        
        issues = []
        
        if language == CodeLanguage.PYTHON and "ast" in parse_result:
            issues.extend(self._detect_python_security_issues(parse_result["ast"]))
        
        return issues
    
    def _detect_python_security_issues(self, ast_node) -> List[Dict[str, Any]]:
        """Detect Python-specific security issues.
        
        Args:
            ast_node: Python AST node
            
        Returns:
            List of security issues
        """
        
        issues = []
        
        for node in ast.walk(ast_node):
            # Detect exec() usage
            if isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == 'exec':
                issues.append({
                    "type": "code_injection",
                    "severity": "high",
                    "message": "Use of exec() can lead to code injection vulnerabilities",
                    "line": getattr(node, 'lineno', 0)
                })
            
            # Detect eval() usage
            if isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == 'eval':
                issues.append({
                    "type": "code_injection",
                    "severity": "high",
                    "message": "Use of eval() can lead to code injection vulnerabilities",
                    "line": getattr(node, 'lineno', 0)
                })
            
            # Detect SQL injection patterns (basic)
            if isinstance(node, ast.Call) and isinstance(node.func, ast.Attribute):
                if node.func.attr in ['execute', 'executemany']:
                    issues.append({
                        "type": "sql_injection",
                        "severity": "medium",
                        "message": "Potential SQL injection vulnerability. Use parameterized queries.",
                        "line": getattr(node, 'lineno', 0)
                    })
        
        return issues
    
    def _detect_compatibility_issues(self, parse_result: Dict[str, Any], language: CodeLanguage) -> List[Dict[str, Any]]:
        """Detect compatibility issues.
        
        Args:
            parse_result: Parsed code information
            language: Programming language
            
        Returns:
            List of compatibility issues
        """
        
        issues = []
        
        if language == CodeLanguage.PYTHON and "ast" in parse_result:
            issues.extend(self._detect_python_compatibility_issues(parse_result["ast"]))
        
        return issues
    
    def _detect_python_compatibility_issues(self, ast_node) -> List[Dict[str, Any]]:
        """Detect Python-specific compatibility issues.
        
        Args:
            ast_node: Python AST node
            
        Returns:
            List of compatibility issues
        """
        
        issues = []
        
        for node in ast.walk(ast_node):
            # Detect print statements (Python 2 style)
            if isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == 'print':
                # Check if it's a print statement (not function call)
                if hasattr(node, 'kwargs') and not node.keywords:
                    issues.append({
                        "type": "python2_print",
                        "severity": "low",
                        "message": "Print statement detected. Use print() function for Python 3 compatibility.",
                        "line": getattr(node, 'lineno', 0)
                    })
        
        return issues
    
    def _detect_python2_issues(self, code: str, parse_result: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Detect Python 2 specific issues for migration.
        
        Args:
            code: Python source code
            parse_result: Parsed code information
            
        Returns:
            List of Python 2 to 3 migration issues
        """
        
        issues = []
        
        # Check for Python 2 specific patterns
        python2_patterns = [
            (r'print\s+[^(]', "Print statement (should be print() function)"),
            (r'xrange\s*\(', "xrange() (should be range() in Python 3)"),
            (r'\.has_key\s*\(', "dict.has_key() (use 'in' operator)"),
            (r'unicode\s*\(', "unicode() (use str() in Python 3)"),
            (r'basestring', "basestring (use str in Python 3)"),
            (r'__nonzero__', "__nonzero__ (use __bool__ in Python 3)"),
            (r'\.iteritems\s*\(\)', "dict.iteritems() (use .items() in Python 3)"),
            (r'\.iterkeys\s*\(\)', "dict.iterkeys() (use .keys() in Python 3)"),
            (r'\.itervalues\s*\(\)', "dict.itervalues() (use .values() in Python 3)"),
        ]
        
        for pattern, description in python2_patterns:
            matches = re.finditer(pattern, code, re.MULTILINE)
            for match in matches:
                issues.append({
                    "type": "python2_compatibility",
                    "severity": "high",
                    "message": f"Python 2 pattern: {description}",
                    "line": code[:match.start()].count('\n') + 1,
                    "pattern": pattern,
                    "suggestion": self._get_python2_suggestion(pattern)
                })
        
        return issues
    
    def _get_python2_suggestion(self, pattern: str) -> str:
        """Get migration suggestion for Python 2 pattern.
        
        Args:
            pattern: Detected pattern
            
        Returns:
            Migration suggestion
        """
        
        suggestions = {
            r'print\s+[^(]': "Use print() function: print('hello')",
            r'xrange\s*\(': "Use range() function: range(10)",
            r'\.has_key\s*\(': "Use 'in' operator: 'key' in dict",
            r'unicode\s*\(': "Use str() function: str('text')",
            r'basestring': "Use str type: isinstance(obj, str)",
            r'__nonzero__': "Use __bool__ method: def __bool__(self):",
            r'\.iteritems\s*\(\)': "Use .items() method: dict.items()",
            r'\.iterkeys\s*\(\)': "Use .keys() method: dict.keys()",
            r'\.itervalues\s*\(\)': "Use .values() method: dict.values()",
        }
        
        return suggestions.get(pattern, "Update for Python 3 compatibility")
    
    def _calculate_metrics(self, parse_result: Dict[str, Any], language: CodeLanguage) -> Dict[str, Any]:
        """Calculate code metrics.
        
        Args:
            parse_result: Parsed code information
            language: Programming language
            
        Returns:
            Code metrics
        """
        
        metrics = {
            "function_count": len(parse_result.get("functions", [])),
            "class_count": len(parse_result.get("classes", [])),
            "import_count": len(parse_result.get("imports", [])),
            "variable_count": len(parse_result.get("variables", [])),
        }
        
        return metrics
    
    def _parse_python2_compat(self, code: str, filename: Optional[str] = None) -> Dict[str, Any]:
        """Parse Python 2 code with compatibility mode.
        
        Args:
            code: Python 2 source code
            filename: Original filename
            
        Returns:
            Parsed information with Python 2 compatibility
        """
        
        logger.info("Using Python 2 compatibility parsing mode")
        
        result = {
            "ast": None,
            "functions": [],
            "classes": [],
            "imports": [],
            "variables": [],
            "complexity_indicators": [],
            "python2_compat": True
        }
        
        try:
            # Basic regex-based parsing for Python 2 code
            # Function definitions
            function_pattern = r'def\s+(\w+)\s*\('
            result["functions"] = [{"name": name, "lineno": 0, "args": [], "decorators": []}
                                 for name in re.findall(function_pattern, code)]
            
            # Class definitions
            class_pattern = r'class\s+(\w+)'
            result["classes"] = [{"name": name, "lineno": 0, "methods": []}
                               for name in re.findall(class_pattern, code)]
            
            # Import statements
            import_pattern = r'import\s+(\w+(?:\s*,\s*\w+)*)'
            from_pattern = r'from\s+(\w+)\s+import\s+(\w+(?:\s*,\s*\w+)*)'
            
            for match in re.finditer(import_pattern, code):
                modules = match.group(1).split(',')
                result["imports"].extend([mod.strip() for mod in modules])
            
            for match in re.finditer(from_pattern, code):
                module = match.group(1)
                imports = match.group(2).split(',')
                result["imports"].extend([f"{module}.{imp.strip()}" for imp in imports])
            
            # Variables (basic detection)
            assignment_pattern = r'(\w+)\s*='
            result["variables"] = list(set(re.findall(assignment_pattern, code)))
            
            logger.info("Python 2 compatibility parsing completed")
            
        except Exception as e:
            # Escape curly braces in error message to avoid loguru format issues
            error_msg = str(e).replace('{', '{{').replace('}', '}}')
            logger.error(f"Python 2 compatibility parsing failed: {error_msg}")
            # Return basic result even if parsing fails
            pass
        
        return result