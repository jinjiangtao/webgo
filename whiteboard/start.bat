@echo off
chcp 65001 >nul
echo ========================================
echo    多人实时协同在线白板系统
echo    Multi-User Real-Time Whiteboard
echo ========================================
echo.

echo [1/3] 检查后端依赖...
cd server
if not exist "go.sum" (
    echo 首次运行，正在下载 Go 依赖...
    go mod tidy
)
cd ..

echo.
echo [2/3] 检查前端依赖...
cd web
if not exist "node_modules" (
    echo 首次运行，正在下载 npm 依赖...
    call npm install
)
cd ..

echo.
echo [3/3] 启动服务...
echo.

echo 正在启动 Go 后端服务 (端口: 8080)...
start "Go Backend" cmd /k "cd server && go run main.go"

echo 等待后端启动...
timeout /t 3 /nobreak >nul

echo.
echo 正在启动 Vue 前端开发服务器 (端口: 5173)...
start "Vue Frontend" cmd /k "cd web && npm run dev"

echo.
echo ========================================
echo   服务启动中，请在浏览器打开：
echo   http://localhost:5173
echo ========================================
echo.
echo   关闭此窗口不会停止服务
echo   如需停止，请关闭对应的后端/前端窗口
echo.
pause
