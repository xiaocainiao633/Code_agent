@echo off
REM CodeSage AI Agent Backend Startup Script for Windows

echo Starting CodeSage AI Agent Backend...

REM Check if Python is available
python --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Python is not installed or not in PATH
    exit /b 1
)

REM Check if pip is available
pip --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: pip is not installed or not in PATH
    exit /b 1
)

REM Create virtual environment if it doesn't exist
if not exist "venv" (
    echo Creating virtual environment...
    python -m venv venv
)

REM Activate virtual environment
echo Activating virtual environment...
call venv\Scripts\activate.bat

REM Install dependencies
echo Installing dependencies...
pip install -r requirements.txt

REM Create necessary directories
echo Creating necessary directories...
if not exist "logs" mkdir logs
if not exist "chroma_db" mkdir chroma_db

REM Copy environment file if .env doesn't exist
if not exist ".env" (
    echo Creating .env file from template...
    copy .env.example .env
    echo Please review and update the .env file with your configuration
)

REM Check Ollama connection
echo Checking Ollama connection...
powershell -Command "try { Invoke-RestMethod -Uri 'http://localhost:11434/api/tags' -TimeoutSec 5 | Out-Null; Write-Host '✓ Ollama is running' } catch { Write-Host '⚠ Warning: Ollama is not responding at http://localhost:11434'; Write-Host 'Please make sure Ollama is installed and running with llama3.2 model' }"

REM Start the application
echo Starting FastAPI application...
echo API will be available at: http://localhost:8000
echo API documentation will be available at: http://localhost:8000/docs
echo Press Ctrl+C to stop the server
echo.

REM Run the application
python main.py

REM Deactivate virtual environment when done
deactivate