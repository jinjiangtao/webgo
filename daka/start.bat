@echo off
chcp 65001 >nul
echo ============================================
echo   自律打卡系统 - 一键启动脚本 (Windows)
echo ============================================
echo.

echo [1/3] 启动 Go 后端服务...
cd /d "%~dp0server"
if not exist "go.sum" (
    echo 正在下载 Go 依赖...
    go mod tidy
)
start "Daka-Server" cmd /k "go run main.go"

echo 等待后端服务启动...
timeout /t 5 /nobreak >nul

echo.
echo [2/3] 安装前端依赖并启动 Vue3 开发服务...
cd /d "%~dp0web"
if not exist "node_modules" (
    echo 正在安装前端依赖，请稍候...
    call npm install
)

echo.
echo [3/3] 启动前端服务...
start "Daka-Web" cmd /k "npm run dev"

echo.
echo ============================================
echo   服务启动中，请稍候...
echo   后端: http://localhost:8080
echo   前端: http://localhost:5173
echo ============================================
echo.
pause
