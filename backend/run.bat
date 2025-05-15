@echo off
echo Starting BPV-Matcher backend...

REM Set environment variables (customize as needed)
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=postgres
set DB_NAME=bpvmatcher
set JWT_SECRET=development-secret-key-change-in-production
set SERVER_PORT=8080
set ENVIRONMENT=development
set ALLOWED_ORIGINS=*

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: Go is not installed or not in PATH
    exit /b 1
)

REM Run go mod tidy to ensure dependencies
echo Installing dependencies...
go mod tidy

REM Build and run the application
echo Building and running application...
go run cmd/backend/main.go

pause 