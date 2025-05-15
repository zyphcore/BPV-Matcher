@echo off
echo Starting BPV-Matcher with Docker...

REM Check if Docker is installed
where docker >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: Docker is not installed or not in PATH
    exit /b 1
)

REM Check if Docker Compose is installed
where docker-compose >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: Docker Compose is not installed or not in PATH
    exit /b 1
)

REM Build and start services
echo Building and starting services...
docker-compose up --build -d

echo.
echo BPV-Matcher is now running!
echo API available at: http://localhost:8080
echo.
echo Press any key to stop the services...
pause

REM Stop services
echo Stopping services...
docker-compose down

pause 