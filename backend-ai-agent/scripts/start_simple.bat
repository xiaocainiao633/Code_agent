@echo off
REM Simple startup script for CodeSage AI Agent Backend

echo.
echo ========================================
echo CodeSage AI Agent Backend Startup
echo ========================================
echo.

REM Check if Python is available
python --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Python is not installed or not in PATH
    echo Please install Python 3.10+ and try again
    pause
    exit /b 1
)

REM Check if we're in the right directory
if not exist "main.py" (
    echo ERROR: main.py not found
    echo Please run this script from the backend-ai-agent directory
    pause
    exit /b 1
)

REM Install missing dependencies
echo Checking dependencies...
pip install pydantic-settings==2.1.0 >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Failed to install pydantic-settings
    echo Please run: pip install pydantic-settings==2.1.0
    pause
    exit /b 1
)

REM Create necessary directories
echo Creating directories...
if not exist "logs" mkdir logs
if not exist "chroma_db" mkdir chroma_db

REM Copy environment file if .env doesn't exist
if not exist ".env" (
    echo Creating .env file from template...
    if exist ".env.example" (
        copy .env.example .env >nul
        echo Please review and update the .env file with your configuration
    ) else (
        echo WARNING: .env.example not found, creating basic .env file
        echo OLLAMA_HOST=http://localhost:11434 > .env
        echo OLLAMA_MODEL=llama3.2 >> .env
        echo API_HOST=0.0.0.0 >> .env
        echo API_PORT=8000 >> .env
    )
)

REM Test basic imports
echo Testing Python environment...
python -c "import fastapi, pydantic, ollama, loguru" >nul 2>&1
if %errorlevel% neq 0 (
    echo WARNING: Some dependencies may be missing
    echo Installing basic dependencies...
    pip install fastapi uvicorn pydantic ollama loguru python-dotenv >nul 2>&1
)

REM Check Ollama connection
echo Checking Ollama connection...
powershell -Command "try { $response = Invoke-RestMethod -Uri 'http://localhost:11434/api/tags' -TimeoutSec 3; Write-Host '✓ Ollama is running' } catch { Write-Host '⚠ Warning: Ollama is not responding at http://localhost:11434'; Write-Host 'Please make sure Ollama is installed and running with llama3.2 model' }"

echo.
echo Starting FastAPI application...
echo API will be available at: http://localhost:8000
echo API documentation: http://localhost:8000/docs
echo Press Ctrl+C to stop the server
echo ========================================
echo.

REM Run the application
python main.py

echo.
echo ========================================
echo Server stopped
echo ========================================
pause