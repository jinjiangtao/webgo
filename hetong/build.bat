@echo off
chcp 65001 >nul
echo ========================================
echo  企业数据可视化平台 - 一键编译部署
echo ========================================
echo.

echo [1/4] 检查环境...
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [错误] 未检测到 Go 环境，请先安装 Go 1.22+
    pause
    exit /b 1
)

where node >nul 2>nul
if %errorlevel% neq 0 (
    echo [错误] 未检测到 Node.js 环境，请先安装 Node.js 18+
    pause
    exit /b 1
)

for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
for /f "tokens=*" %%i in ('node --version') do set NODE_VERSION=%%i
echo Go 版本: %GO_VERSION%
echo Node 版本: %NODE_VERSION%
echo.

echo [2/4] 编译后端服务...
cd /d "%~dp0server"
echo 正在下载依赖...
go mod download
echo 正在编译...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o dashboard-server.exe ./cmd/main.go
if %errorlevel% neq 0 (
    echo [错误] 后端编译失败
    pause
    exit /b 1
)
echo 后端编译成功！
echo.

echo [3/4] 构建前端应用...
cd /d "%~dp0web"
echo 正在安装依赖...
call npm install
if %errorlevel% neq 0 (
    echo [警告] npm install 失败，尝试使用国内镜像...
    call npm install --registry=https://registry.npmmirror.com
)
echo 正在构建...
call npm run build
if %errorlevel% neq 0 (
    echo [错误] 前端构建失败
    pause
    exit /b 1
)
echo 前端构建成功！
echo.

echo [4/4] 准备部署目录...
cd /d "%~dp0"
if not exist "dist" mkdir dist
if not exist "dist\server" mkdir dist\server
if not exist "dist\web" mkdir dist\web

copy /y "server\dashboard-server.exe" "dist\server\" >nul
if exist "server\data.db" copy /y "server\data.db" "dist\server\" >nul
xcopy /e /i /y "web\dist" "dist\web" >nul

echo.
echo ========================================
echo  编译部署完成！
echo ========================================
echo.
echo 部署目录: %~dp0dist
echo   - 后端: dist\server\dashboard-server.exe
echo   - 前端: dist\web\
echo.
echo 启动方式:
echo   双击运行 start.bat 启动服务
echo   或手动执行:
echo     1. cd dist\server && dashboard-server.exe
echo     2. 浏览器访问: http://localhost:8080
echo.
pause
