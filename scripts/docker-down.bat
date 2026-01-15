@echo off
chcp 65001 > nul
REM Study-UPC Docker Services Stop Script
REM Purpose: Stop infrastructure services (PostgreSQL, Redis, MinIO)

setlocal enabledelayedexpansion

echo.
echo ========================================
echo   Stopping Docker Services
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

REM Stop services
echo [INFO] Stopping infrastructure services...
cd docker
!COMPOSE_CMD! down
cd ..

if %errorlevel% neq 0 (
    echo [ERROR] Failed to stop services
    pause
    exit /b 1
)

echo.
echo ========================================
echo   Services Stopped Successfully!
echo ========================================
echo.

pause
