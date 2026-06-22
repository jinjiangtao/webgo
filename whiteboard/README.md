# 多人实时协同在线白板系统

适配线上会议、远程协作、教学演示场景，支持多人实时同步操作的在线白板系统。

## 技术栈

- **后端**: Go 1.21+ / Gin / Gorilla WebSocket / GORM / SQLite
- **前端**: Vue 3.4+ / Vite 5 / Pinia / Element Plus / Canvas API
- **实时通信**: WebSocket (房间隔离广播)
- **数据存储**: SQLite (白板、操作记录、快照)

## 功能特性

### 🎨 绘图工具
- **画笔**: 自由曲线绘制
- **直线**: 两点直线
- **图形**: 矩形、圆形、椭圆、三角形
- **文字**: 自定义文字插入
- **橡皮擦**: 清除笔迹
- **属性**: 自定义粗细(1-30px)、描边颜色、填充颜色、透明度(0-100%)、字号(12-72px)

### 👥 多人协同
- **实时同步**: 毫秒级绘图操作同步
- **用户列表**: 显示当前在线成员
- **远程光标**: 实时显示其他用户光标位置(颜色区分)
- **房间隔离**: 不同白板互不干扰

### 📁 白板管理
- **新建白板**: 创建空白白板
- **保存快照**: 保存当前状态为命名快照
- **打开白板**: 从白板列表或历史快照加载
- **云端存储**: SQLite 持久化存储
- **操作历史**: 撤销/重做 (最多 100 步)

### 🪄 辅助功能
- **缩放平移**: 滚轮缩放(0.25x-3x)、Alt+拖拽平移
- **背景切换**: 7 种预设背景色
- **图片粘贴**: 剪贴板图片直接粘贴到画布
- **一键清空**: 清空画布所有内容

## 目录结构

```
whiteboard/
├── server/                     # Go 后端
│   ├── main.go                 # 主入口
│   ├── data/                   # SQLite 数据目录
│   └── internal/
│       ├── model/              # 数据模型
│       ├── database/           # 数据库层
│       ├── ws/                 # WebSocket 管理器
│       ├── handler/            # HTTP/WS 处理器
│       └── middleware/         # 中间件
├── web/                        # Vue 前端
│   ├── public/
│   ├── src/
│   │   ├── components/         # UI 组件
│   │   ├── views/              # 页面视图
│   │   ├── utils/              # 工具函数 (API/WS/Canvas)
│   │   ├── composables/        # Vue 组合式函数
│   │   ├── stores/             # Pinia 状态管理
│   │   ├── router/             # 路由配置
│   │   ├── App.vue
│   │   └── main.js
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
├── start.bat                   # Windows 一键启动
└── README.md
```

## 快速启动

### 方式一: 一键启动 (Windows)

双击运行 `start.bat` 即可自动安装依赖并启动前后端服务。

### 方式二: 手动启动

#### 启动后端

```bash
cd server
go mod tidy          # 首次运行下载依赖
go run main.go       # 启动服务: http://localhost:8080
```

#### 启动前端

```bash
cd web
npm install          # 首次运行下载依赖
npm run dev          # 启动服务: http://localhost:5173
```

### 访问

打开浏览器访问 **http://localhost:5173** 即可使用。

如需多人协同，在多台设备或多个浏览器窗口中打开同一 URL，所有操作会实时同步。

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/whiteboards` | 创建白板 |
| GET | `/api/whiteboards` | 白板列表 |
| GET | `/api/whiteboards/:id` | 获取白板详情+操作记录 |
| PUT | `/api/whiteboards/:id` | 更新白板信息 |
| DELETE | `/api/whiteboards/:id` | 删除白板 |
| POST | `/api/whiteboards/:id/operations` | 批量保存操作 |
| POST | `/api/whiteboards/:id/clear` | 清空操作记录 |
| POST | `/api/whiteboards/:id/snapshots` | 创建快照 |
| GET | `/api/whiteboards/:id/snapshots` | 快照列表 |
| GET | `/ws` | WebSocket 连接 |

### WebSocket 消息类型

| 类型 | 说明 |
|------|------|
| `join` / `leave` | 用户加入/离开通知 |
| `users` | 在线用户列表推送 |
| `cursor` | 光标位置同步 |
| `draw` | 绘图操作同步 (start/move/end/undo/redo) |
| `clear` | 清空画布 |
| `ping` / `pong` | 心跳检测 |

## 数据库表结构

- **users**: 用户信息 (ID, 用户名, 颜色)
- **whiteboards**: 白板元数据 (ID, 名称, 背景色, 创建/更新时间)
- **operations**: 操作日志 (白板ID, 用户ID, 操作类型, JSON数据, 时间戳)
- **snapshots**: 快照记录 (白板ID, 名称, 序列化操作集合, 创建时间)

## 快捷键提示

| 操作 | 快捷键 |
|------|--------|
| 平移画布 | `Alt + 鼠标拖拽` / `鼠标中键拖拽` |
| 缩放画布 | `鼠标滚轮` |
| 文字确认 | `Enter` |
| 文字取消 | `ESC` / 点击空白处 |
| 粘贴图片 | `Ctrl + V` |
