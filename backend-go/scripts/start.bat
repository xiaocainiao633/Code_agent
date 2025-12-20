@echo off
REM CodeSage Go Backend Startup Script for Windows

echo.
echo ========================================
echo CodeSage Go Backend Startup
echo ========================================
echo.

REM 检查Go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go 1.22+ and try again
    pause
    exit /b 1
)

REM 检查是否在项目目录
if not exist "go.mod" (
    echo ERROR: go.mod not found
    echo Please run this script from the backend-go directory
    pause
    exit /b 1
)

REM 创建必要目录
echo Creating directories...
if not exist "logs" mkdir logs
if not exist "temp" mkdir temp
if not exist "bin" mkdir bin

REM 安装依赖
echo Installing dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo ERROR: Failed to install dependencies
    pause
    exit /b 1
)

REM 构建项目
echo Building project...
go build -o bin/server cmd/server/main.go
if %errorlevel% neq 0 (
    echo ERROR: Build failed
    pause
    exit /b 1
)

echo.
echo Starting Go backend server...
echo Server will be available at: http://localhost:8080
echo Health check: http://localhost:8080/api/v1/health
echo Press Ctrl+C to stop the server
echo ========================================
echo.

REM 运行服务器
bin\server

echo.
echo ========================================
echo Server stopped
echo ========================================
pause