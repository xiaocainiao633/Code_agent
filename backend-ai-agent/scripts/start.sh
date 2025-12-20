#!/bin/bash

# CodeSage AI Agent Backend Startup Script

echo "Starting CodeSage AI Agent Backend..."

# Check if Python is available
if ! command -v python3 &> /dev/null; then
    echo "Error: Python 3 is not installed or not in PATH"
    exit 1
fi

# Check if pip is available
if ! command -v pip3 &> /dev/null; then
    echo "Error: pip3 is not installed or not in PATH"
    exit 1
fi

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment
echo "Activating virtual environment..."
source venv/bin/activate

# Install dependencies
echo "Installing dependencies..."
pip install -r requirements.txt

# Create necessary directories
echo "Creating necessary directories..."
mkdir -p logs
mkdir -p chroma_db

# Copy environment file if .env doesn't exist
if [ ! -f ".env" ]; then
    echo "Creating .env file from template..."
    cp .env.example .env
    echo "Please review and update the .env file with your configuration"
fi

# Check Ollama connection
echo "Checking Ollama connection..."
if command -v curl &> /dev/null; then
    if curl -s "http://localhost:11434/api/tags" > /dev/null; then
        echo "✓ Ollama is running"
    else
        echo "⚠ Warning: Ollama is not responding at http://localhost:11434"
        echo "Please make sure Ollama is installed and running with llama3.2 model"
    fi
else
    echo "⚠ Warning: curl not available, skipping Ollama connection check"
fi

# Start the application
echo "Starting FastAPI application..."
echo "API will be available at: http://localhost:8000"
echo "API documentation will be available at: http://localhost:8000/docs"
echo "Press Ctrl+C to stop the server"

# Run the application
python main.py