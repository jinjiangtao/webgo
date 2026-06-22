<template>
  <div class="whiteboard-view" ref="viewRef" @paste="handlePaste">
    <div class="view-header">
      <Toolbar
        @undo="handleUndo"
        @redo="handleRedo"
        @clear="handleClear"
        @save="handleSaveClick"
        @new="handleNew"
        @open="handleOpenClick"
      />
    </div>

    <div class="view-body">
      <div class="canvas-container" ref="containerRef">
        <canvas
          ref="canvasEl"
          class="whiteboard-canvas"
          @mousedown="handleMouseDown"
          @mousemove="handleMouseMove"
          @mouseup="handleMouseUp"
          @mouseleave="handleMouseUp"
          @wheel="handleWheel"
          @touchstart.prevent="handleMouseDown"
          @touchmove.prevent="handleMouseMove"
          @touchend.prevent="handleMouseUp"
        />
        <RemoteCursors
          :cursors="cursorsArray"
          :zoom="whiteboardStore.zoom"
          :pan-x="whiteboardStore.panX"
          :pan-y="whiteboardStore.panY"
        />
        <TextInput
          v-model:model-value="showTextInput"
          :x="textPosition.x"
          :y="textPosition.y"
          :font-size="toolStore.fontSize"
          :font-family="toolStore.fontFamily"
          :color="toolStore.strokeColor"
          :zoom="whiteboardStore.zoom"
          :pan-x="whiteboardStore.panX"
          :pan-y="whiteboardStore.panY"
          @confirm="handleTextConfirm"
          @cancel="handleTextCancel"
        />
      </div>

      <UserPanel
        @rename="handleRename"
        @change-background="handleBackgroundChange"
        @change-zoom="handleZoomChange"
      />
    </div>

    <WhiteboardDialog
      v-model:model-value="openDialogVisible"
      mode="open"
      :whiteboard-id="whiteboardStore.whiteboardId"
      @open="handleOpenWhiteboard"
      @delete="handleDeleteWhiteboard"
      @load-snapshot="handleLoadSnapshot"
    />

    <WhiteboardDialog
      v-model:model-value="saveDialogVisible"
      mode="save"
      :whiteboard-id="whiteboardStore.whiteboardId"
      :snapshot-operations="pendingOperations"
      @save="handleSaveSnapshot"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import { v4 as uuidv4 } from 'uuid'
import Toolbar from '../components/Toolbar.vue'
import UserPanel from '../components/UserPanel.vue'
import RemoteCursors from '../components/RemoteCursors.vue'
import TextInput from '../components/TextInput.vue'
import WhiteboardDialog from '../components/WhiteboardDialog.vue'
import { useToolStore, TOOL_TYPES } from '../stores/tool'
import { useWhiteboardStore } from '../stores/whiteboard'
import { useCanvas } from '../composables/useCanvas'
import { getWSClient } from '../utils/websocket'
import {
  createWhiteboard,
  getWhiteboard,
  updateWhiteboard,
  saveOperations
} from '../utils/api'

const viewRef = ref(null)
const containerRef = ref(null)
const canvasEl = ref(null)

const toolStore = useToolStore()
const whiteboardStore = useWhiteboardStore()

const { cursors, zoom, panX, panY, background } = storeToRefs(whiteboardStore)

const ws = getWSClient()

const {
  engine,
  config,
  shapes,
  isDrawing,
  isTextPending,
  textPosition,
  TOOL_TYPES: CANVAS_TOOL_TYPES,
  initCanvas,
  setTool,
  setStrokeWidth,
  setStrokeColor,
  setFillColor,
  setOpacity,
  setFontSize,
  setFontFamily,
  setBackground,
  applyZoom,
  handleStart,
  handleMove,
  handleEnd,
  confirmText,
  cancelText,
  undo,
  redo,
  renderAll,
  clearCanvas,
  setShapes,
  getAllShapes,
  resize
} = useCanvas()

const showTextInput = ref(false)
const openDialogVisible = ref(false)
const saveDialogVisible = ref(false)

const isPanning = ref(false)
const panStartPoint = ref({ x: 0, y: 0 })
const panStartOffset = ref({ x: 0, y: 0 })

const pendingOperations = ref([])

const userId = ref(localStorage.getItem('wb_user_id') || uuidv4())
const username = ref(localStorage.getItem('wb_username') || ('用户' + Math.floor(Math.random() * 10000)))
const userColor = ref(localStorage.getItem('wb_user_color') || getRandomColor())

localStorage.setItem('wb_user_id', userId.value)
localStorage.setItem('wb_user_name', username.value)
localStorage.setItem('wb_user_color', userColor.value)

function getRandomColor() {
  const colors = ['#f56c6c', '#e6a23c', '#67c23a', '#409eff', '#909399', '#9b59b6', '#1abc9c', '#e91e63']
  return colors[Math.floor(Math.random() * colors.length)]
}

const cursorsArray = computed(() => {
  const arr = []
  for (const [uid, cursor] of cursors.value.entries()) {
    if (uid === userId.value) continue
    const user = whiteboardStore.onlineUsers.find(u => u.userId === uid)
    arr.push({
      userId: uid,
      x: cursor.x,
      y: cursor.y,
      color: user?.color || userColor.value,
      username: user?.username || ''
    })
  }
  return arr
})

function syncToolToCanvas() {
  setTool(toolStore.currentTool)
  setStrokeWidth(toolStore.strokeWidth)
  setStrokeColor(toolStore.strokeColor)
  setFillColor(toolStore.fillColor)
  setOpacity(toolStore.opacity / 100)
  setFontSize(toolStore.fontSize)
  setFontFamily(toolStore.fontFamily)
}

function syncWhiteboardToCanvas() {
  applyZoom(whiteboardStore.zoom, whiteboardStore.panX, whiteboardStore.panY)
  setBackground(whiteboardStore.background)
}

watch(
  () => toolStore.$state,
  () => {
    syncToolToCanvas()
  },
  { deep: true }
)

watch(
  [zoom, panX, panY],
  () => {
    applyZoom(zoom.value, panX.value, panY.value)
  }
)

watch(background, (val) => {
  setBackground(val)
  if (whiteboardStore.whiteboardId) {
    updateWhiteboard(whiteboardStore.whiteboardId, undefined, val).catch(() => {})
  }
})

watch(isTextPending, (val) => {
  if (val) {
    nextTick(() => {
      showTextInput.value = true
    })
  } else {
    showTextInput.value = false
  }
})

function getCanvasCoords(e) {
  if (!canvasEl.value) return { x: 0, y: 0 }
  const rect = canvasEl.value.getBoundingClientRect()
  const clientX = e.touches ? e.touches[0].clientX : e.clientX
  const clientY = e.touches ? e.touches[0].clientY : e.clientY
  const x = (clientX - rect.left - panX.value) / zoom.value
  const y = (clientY - rect.top - panY.value) / zoom.value
  return { x, y, clientX, clientY }
}

function sendCursorMessage(x, y) {
  ws.send('cursor', {
    userId: userId.value,
    whiteboardId: whiteboardStore.whiteboardId,
    x,
    y
  })
}

let lastCursorSendTime = 0
function throttleSendCursor(x, y) {
  const now = Date.now()
  if (now - lastCursorSendTime < 50) return
  lastCursorSendTime = now
  sendCursorMessage(x, y)
}

function handleMouseDown(e) {
  if (!engine.value) return

  const { x, y, clientX, clientY } = getCanvasCoords(e)

  if (e.button === 1 || (e.button === 0 && e.altKey) || toolStore.currentTool === 'pan') {
    isPanning.value = true
    panStartPoint.value = { x: clientX, y: clientY }
    panStartOffset.value = { x: panX.value, y: panY.value }
    return
  }

  sendCursorMessage(x, y)
  handleStart(e)

  if (toolStore.currentTool === TOOL_TYPES.TEXT) {
    return
  }

  if (isDrawing.value) {
    ws.send('draw', {
      userId: userId.value,
      whiteboardId: whiteboardStore.whiteboardId,
      phase: 'start',
      shape: { ...shapes.value }
    })
  }
}

function handleMouseMove(e) {
  if (!engine.value) return

  const { x, y, clientX, clientY } = getCanvasCoords(e)
  throttleSendCursor(x, y)

  if (isPanning.value) {
    const dx = clientX - panStartPoint.value.x
    const dy = clientY - panStartPoint.value.y
    whiteboardStore.setPan(panStartOffset.value.x + dx, panStartOffset.value.y + dy)
    return
  }

  if (!isDrawing.value) return

  handleMove(e)

  if (isDrawing.value) {
    ws.send('draw', {
      userId: userId.value,
      whiteboardId: whiteboardStore.whiteboardId,
      phase: 'move',
      shape: { ...shapes.value }
    })
  }
}

function handleMouseUp(e) {
  if (isPanning.value) {
    isPanning.value = false
    return
  }

  if (!engine.value) return

  const { x, y } = getCanvasCoords(e)
  sendCursorMessage(x, y)

  const prevLength = shapes.value.length
  handleEnd(e)

  if (!isTextPending.value && shapes.value.length > prevLength) {
    const newShape = shapes.value[shapes.value.length - 1]
    pushOperation({
      type: 'add',
      shape: newShape
    })

    ws.send('draw', {
      userId: userId.value,
      whiteboardId: whiteboardStore.whiteboardId,
      phase: 'end',
      shape: newShape
    })
  }
}

function handleWheel(e) {
  e.preventDefault()
  const { x, y } = getCanvasCoords(e)
  const delta = e.deltaY > 0 ? -0.1 : 0.1
  const newZoom = Math.max(0.25, Math.min(3, zoom.value + delta))

  const newPanX = panX.value - x * (newZoom - zoom.value)
  const newPanY = panY.value - y * (newZoom - zoom.value)

  whiteboardStore.setZoom(newZoom)
  whiteboardStore.setPan(newPanX, newPanY)
}

function handleTextConfirm(text) {
  const prevLength = shapes.value.length
  confirmText(text)
  if (shapes.value.length > prevLength) {
    const newShape = shapes.value[shapes.value.length - 1]
    pushOperation({
      type: 'add',
      shape: newShape
    })
    ws.send('draw', {
      userId: userId.value,
      whiteboardId: whiteboardStore.whiteboardId,
      phase: 'end',
      shape: newShape
    })
  }
}

function handleTextCancel() {
  cancelText()
}

function handleUndo() {
  undo()
  pushOperation({ type: 'undo' })
  ws.send('draw', {
    userId: userId.value,
    whiteboardId: whiteboardStore.whiteboardId,
    phase: 'undo',
    shapes: getAllShapes()
  })
}

function handleRedo() {
  redo()
  pushOperation({ type: 'redo' })
  ws.send('draw', {
    userId: userId.value,
    whiteboardId: whiteboardStore.whiteboardId,
    phase: 'redo',
    shapes: getAllShapes()
  })
}

function handleClear() {
  clearCanvas()
  pushOperation({ type: 'clear' })
  ws.send('clear', {
    userId: userId.value,
    whiteboardId: whiteboardStore.whiteboardId
  })
}

function handleSaveClick() {
  if (!whiteboardStore.whiteboardId) {
    ElMessage.warning('请先创建或打开一个白板')
    return
  }
  flushPendingOperations().then(() => {
    saveDialogVisible.value = true
  })
}

function handleOpenClick() {
  openDialogVisible.value = true
}

async function handleNew() {
  try {
    await flushPendingOperations()
    const wb = await createWhiteboard('未命名白板', '#ffffff')
    const data = wb?.data || wb
    initWhiteboard(data.id, data.name, data.background || '#ffffff', [])
    joinWhiteboard(data.id)
    ElMessage.success('已创建新白板')
  } catch (e) {
    ElMessage.error('创建白板失败: ' + e.message)
  }
}

function handleRename(name) {
  if (!whiteboardStore.whiteboardId) return
  updateWhiteboard(whiteboardStore.whiteboardId, name).catch((e) => {
    ElMessage.error('重命名失败: ' + e.message)
  })
}

function handleBackgroundChange(color) {
  if (!whiteboardStore.whiteboardId) return
  updateWhiteboard(whiteboardStore.whiteboardId, undefined, color).catch((e) => {
    ElMessage.error('修改背景失败: ' + e.message)
  })
}

function handleZoomChange(val) {
  // zoom is handled by watcher
}

function handleOpenWhiteboard(row) {
  flushPendingOperations()
    .then(() => getWhiteboard(row.id))
    .then((res) => {
      const data = res?.data || res
      const ops = data.operations || []
      const shapesList = ops
        .filter(o => o.type === 'add')
        .map(o => o.shape)
        .filter(Boolean)
      initWhiteboard(row.id, data.name || row.name, data.background || '#ffffff', shapesList)
      joinWhiteboard(row.id)
    })
    .catch((e) => {
      ElMessage.error('打开白板失败: ' + e.message)
    })
}

function handleDeleteWhiteboard() {
  // handled in dialog
}

function handleLoadSnapshot(row) {
  const ops = row.operations || []
  const shapesList = ops
    .filter(o => o.type === 'add')
    .map(o => o.shape)
    .filter(Boolean)
  setShapes(shapesList)
  renderAll()
  ElMessage.success('快照已加载')
}

function handleSaveSnapshot() {
  ElMessage.success('快照保存成功')
}

function pushOperation(op) {
  pendingOperations.value.push({
    userId: userId.value,
    type: op.type,
    data: op.shape || op.shapes || null,
    timestamp: Date.now()
  })
  if (pendingOperations.value.length >= 50) {
    flushPendingOperations()
  }
}

async function flushPendingOperations() {
  if (!whiteboardStore.whiteboardId || pendingOperations.value.length === 0) {
    return
  }
  const ops = pendingOperations.value.splice(0)
  try {
    await saveOperations(whiteboardStore.whiteboardId, ops)
  } catch (e) {
    pendingOperations.value.unshift(...ops)
    console.error('保存操作失败', e)
  }
}

function handlePaste(e) {
  if (!engine.value) return
  const items = e.clipboardData?.items
  if (!items) return

  for (const item of items) {
    if (item.type.startsWith('image/')) {
      e.preventDefault()
      const file = item.getAsFile()
      if (!file) continue
      const reader = new FileReader()
      reader.onload = (ev) => {
        const base64 = ev.target.result
        const { x, y } = getCanvasCoords({
          clientX: window.innerWidth / 2,
          clientY: window.innerHeight / 2,
          touches: null
        })
        addImageShape(base64, x, y)
      }
      reader.readAsDataURL(file)
      break
    }
  }
}

function addImageShape(src, x, y) {
  const img = new Image()
  img.onload = () => {
    let w = img.width
    let h = img.height
    const maxSize = 400
    if (w > maxSize || h > maxSize) {
      const ratio = Math.min(maxSize / w, maxSize / h)
      w *= ratio
      h *= ratio
    }
    const shape = {
      id: 'shape_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9),
      type: 'image',
      src,
      x,
      y,
      width: w,
      height: h,
      opacity: 1,
      strokeWidth: 0,
      strokeColor: 'transparent',
      fillColor: 'transparent'
    }
    const ctx = engine.value.ctx
    ctx.save()
    ctx.globalAlpha = 1
    ctx.drawImage(img, x, y, w, h)
    ctx.restore()

    const historyMgr = (engine._currentShape ? null : null) // skip
    // Use setShapes path
    const currentShapes = getAllShapes()
    currentShapes.push(shape)
    setShapes(currentShapes)

    pushOperation({ type: 'add', shape })
    ws.send('draw', {
      userId: userId.value,
      whiteboardId: whiteboardStore.whiteboardId,
      phase: 'end',
      shape
    })

    // Patch engine to render image
    patchEngineRenderImage()
  }
  img.src = src
}

function patchEngineRenderImage() {
  if (!engine.value) return
  const orig = engine.value.renderShape.bind(engine.value)
  if (engine.value._renderImagePatched) return
  engine.value._renderImagePatched = true
  engine.value.renderShape = function (shape) {
    if (!shape) return
    if (shape.type === 'image') {
      const ctx = this.ctx
      ctx.save()
      ctx.globalAlpha = shape.opacity ?? 1
      const img = new Image()
      img.onload = () => {
        this.renderShapes(getAllShapes())
      }
      if (!shape._imgLoaded) {
        img.src = shape.src
        shape._imgLoaded = true
      }
      ctx.drawImage(img, shape.x, shape.y, shape.width, shape.height)
      ctx.restore()
      return
    }
    orig(shape)
  }
}

function joinWhiteboard(id) {
  if (!id) return
  ws.send('join', {
    userId: userId.value,
    username: username.value,
    color: userColor.value,
    whiteboardId: id
  })
}

function setupWSHandlers() {
  ws.on('users', (payload) => {
    const users = payload?.users || payload || []
    whiteboardStore.setUsers(users)
  })

  ws.on('join', (payload) => {
    whiteboardStore.addUser({
      userId: payload.userId,
      username: payload.username,
      color: payload.color
    })
  })

  ws.on('leave', (payload) => {
    whiteboardStore.removeUser(payload.userId)
  })

  ws.on('cursor', (payload) => {
    if (payload.userId === userId.value) return
    whiteboardStore.setCursor(payload.userId, payload.x, payload.y)
  })

  ws.on('draw', (payload) => {
    if (payload.userId === userId.value) return
    if (!payload.shape && !payload.shapes) return

    if (payload.phase === 'undo' || payload.phase === 'redo') {
      setShapes(payload.shapes || [])
      renderAll()
    } else if (payload.phase === 'end') {
      patchEngineRenderImage()
      const current = getAllShapes()
      current.push(payload.shape)
      setShapes(current)
      renderAll()
    }
  })

  ws.on('clear', (payload) => {
    if (payload.userId === userId.value) return
    clearCanvas()
  })
}

function initWhiteboard(id, name, bg, shapesList) {
  whiteboardStore.setWhiteboard(id, name, bg)
  setShapes(shapesList || [])
  renderAll()
  pendingOperations.value = []
}

function setupCanvasSize() {
  if (!canvasEl.value || !containerRef.value) return
  const { width, height } = containerRef.value.getBoundingClientRect()
  resize(width, height)
}

let resizeObserver = null

onMounted(async () => {
  nextTick(() => {
    initCanvas(canvasEl.value)
    syncToolToCanvas()
    syncWhiteboardToCanvas()
    setupCanvasSize()
    patchEngineRenderImage()
  })

  resizeObserver = new ResizeObserver(() => {
    setupCanvasSize()
  })
  if (containerRef.value) {
    resizeObserver.observe(containerRef.value)
  }

  setupWSHandlers()

  const wsProto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${wsProto}//${window.location.host}/ws`
  try {
    await ws.connect(wsUrl)
  } catch (e) {
    console.warn('WS 连接失败，将以离线模式运行', e)
  }

  const id = whiteboardStore.whiteboardId
  if (id) {
    joinWhiteboard(id)
  } else {
    try {
      const wb = await createWhiteboard('未命名白板', '#ffffff')
      const data = wb?.data || wb
      initWhiteboard(data.id, data.name, data.background || '#ffffff', [])
      joinWhiteboard(data.id)
    } catch (e) {
      console.warn('创建默认白板失败，以空画布运行', e)
    }
  }
})

onUnmounted(() => {
  if (resizeObserver && containerRef.value) {
    resizeObserver.unobserve(containerRef.value)
  }
  flushPendingOperations().finally(() => {
    if (whiteboardStore.whiteboardId) {
      ws.send('leave', {
        userId: userId.value,
        whiteboardId: whiteboardStore.whiteboardId
      })
    }
  })
})
</script>

<style scoped>
.whiteboard-view {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f0f2f5;
}

.view-header {
  flex-shrink: 0;
  z-index: 100;
}

.view-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.canvas-container {
  flex: 1;
  position: relative;
  overflow: hidden;
  background: #e8eaed;
}

.whiteboard-canvas {
  display: block;
  width: 100%;
  height: 100%;
  cursor: crosshair;
}
</style>
