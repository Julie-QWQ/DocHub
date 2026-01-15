@echo off
chcp 65001 > nul
REM Study-UPC Docker Services Start Script
REM Purpose: Start only infrastructure services (PostgreSQL, Redis, MinIO)

setlocal enabledelayedexpansion

echo.
echo ========================================
echo   Starting Docker Services
echo ========================================
echo.

cd %~dp0..

REM Check if Docker is running
echo [INFO] Checking Docker...
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker is not running
    pause
    exit /b 1
)
echo [SUCCESS] Docker is running

REM Check .env file
echo [INFO] Checking configuration...
if not exist "docker\.env" (
    if exist "docker\.env.example" (
        copy docker\.env.example docker\.env >nul 2>&1
        echo [SUCCESS] .env file created
    )
)

REM Detect docker compose command
docker compose version >nul 2>&1
if %errorlevel% equ 0 (
    set COMPOSE_CMD=docker compose
) else (
    docker-compose version >nul 2>&1
    if %errorlevel% equ 0 (
        set COMPOSE_CMD=docker-compose
    ) else (
        echo [ERROR] Docker Compose not found
        pause
        exit /b 1
    )
)

REM Start services
echo.
echo [INFO] Starting infrastructure services...
cd docker
!COMPOSE_CMD! up -d postgres redis minio
cd ..

if %errorlevel% neq 0 (
    echo [ERROR] Failed to start services
    pause
    exit /b 1
)

REM Wait and check services
timeout /t 3 /nobreak > nul

echo.
echo [INFO] Service Status:
cd docker
!COMPOSE_CMD! ps postgres redis minio
cd ..

echo.
echo ========================================
echo   Services Started Successfully!
echo ========================================
echo.
echo PostgreSQL:   localhost:5432
echo Redis:        localhost:6379
echo MinIO API:    http://localhost:9000
echo MinIO Console: http://localhost:9001
echo   Username: minioadmin
echo   Password: minioadmin
echo.
echo ========================================
echo.

pause
