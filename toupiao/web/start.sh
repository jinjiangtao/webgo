#!/bin/bash
echo "========================================"
echo "   投票系统前端启动脚本 (Linux/macOS)"
echo "========================================"
echo ""

echo "[1/3] 检查 Node.js 环境..."
if ! command -v node &> /dev/null; then
    echo "[错误] 未检测到 Node.js，请先安装 Node.js 16+ 版本"
    echo "下载地址: https://nodejs.org/"
    exit 1
fi
echo "[OK] Node.js 环境正常"
echo ""

echo "[2/3] 安装项目依赖..."
if [ ! -d "node_modules" ]; then
    echo "正在安装依赖，请稍候..."
    npm install
    if [ $? -ne 0 ]; then
        echo "[错误] 依赖安装失败，请检查网络连接"
        exit 1
    fi
    echo "[OK] 依赖安装完成"
else
    echo "[跳过] 依赖已存在，如需更新请手动执行 npm install"
fi
echo ""

echo "[3/3] 启动开发服务器..."
echo "访问地址: http://localhost:3000"
echo "按 Ctrl+C 停止服务器"
echo ""
npm run dev
