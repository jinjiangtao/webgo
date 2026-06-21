@echo off
chcp 65001 >nul
title 企业数据可视化平台 - 开发模式

echo ========================================
echo  企业数据可视化平台 - 开发模式
echo ========================================
echo.

set "PROJECT_DIR=%~dp0"
set "SERVER_DIR=%PROJECT_DIR%server"
set "WEB_DIR=%PROJECT_DIR%web"

echo [1/2] 启动后端服务 (端口: 8080)...
cd /d "%SERVER_DIR%"
start "后端服务" cmd /k "go run ./cmd/main.go"

timeout /t 3 /nobreak >nul

echo [2/2] 启动前端开发服务 (端口: 5173)...
cd /d "%WEB_DIR%"
start "前端服务" cmd /k "npm run dev"

timeout /t 2 /nobreak >nul

echo.
echo ========================================
echo  开发环境已启动！
echo ========================================
echo.
echo 后端API: http://localhost:8080/api
echo 前端地址: http://localhost:5173
echo.
echo 请在浏览器中打开 http://localhost:5173
echo.
pause
