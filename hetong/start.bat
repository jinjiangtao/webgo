@echo off
chcp 65001 >nul
title 企业数据可视化平台 - 服务启动

echo ========================================
echo  企业数据可视化平台
echo ========================================
echo.

set "PROJECT_DIR=%~dp0"
set "SERVER_DIR=%PROJECT_DIR%server"
set "WEB_DIR=%PROJECT_DIR%web\dist"

cd /d "%SERVER_DIR%"

if not exist "dashboard-server.exe" (
    echo [错误] 未找到后端可执行文件，请先运行 build.bat 编译
    echo.
    pause
    exit /b 1
)

echo 正在启动后端服务...
echo 服务地址: http://localhost:8080
echo API 地址: http://localhost:8080/api
echo.
echo 按 Ctrl+C 停止服务
echo ========================================
echo.

start "" http://localhost:8080

dashboard-server.exe

echo.
echo 服务已停止
pause
