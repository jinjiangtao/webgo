#!/bin/bash
echo "============================================"
echo "  自律打卡系统 - 一键启动脚本 (Mac/Linux)"
echo "============================================"
echo ""

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "[1/3] 启动 Go 后端服务..."
cd "$SCRIPT_DIR/server"
if [ ! -f "go.sum" ]; then
    echo "正在下载 Go 依赖..."
    go mod tidy
fi
gnome-terminal --title="Daka-Server" -- go run main.go 2>/dev/null || \
osascript -e 'tell app "Terminal" to do script "cd '"$SCRIPT_DIR/server"' && go run main.go"' 2>/dev/null || \
(go run main.go &)

sleep 5

echo ""
echo "[2/3] 安装前端依赖并启动 Vue3 开发服务..."
cd "$SCRIPT_DIR/web"
if [ ! -d "node_modules" ]; then
    echo "正在安装前端依赖，请稍候..."
    npm install
fi

echo ""
echo "[3/3] 启动前端服务..."
gnome-terminal --title="Daka-Web" -- npm run dev 2>/dev/null || \
osascript -e 'tell app "Terminal" to do script "cd '"$SCRIPT_DIR/web"' && npm run dev"' 2>/dev/null || \
(npm run dev &)

echo ""
echo "============================================"
echo "  服务启动中，请稍候..."
echo "  后端: http://localhost:8080"
echo "  前端: http://localhost:5173"
echo "============================================"
