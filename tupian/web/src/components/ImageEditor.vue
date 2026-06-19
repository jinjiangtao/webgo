<template>
  <div class="editor-container" ref="containerRef">
    <div class="canvas-wrapper" ref="wrapperRef">
      <canvas
        ref="mainCanvas"
        class="main-canvas"
        @mousedown="handleMouseDown"
        @mousemove="handleMouseMove"
        @mouseup="handleMouseUp"
        @mouseleave="handleMouseUp"
      ></canvas>
      <div
        v-if="isCropping"
        class="crop-overlay"
        ref="cropOverlay"
        @mousedown="startCrop"
      >
        <div
          v-if="cropArea"
          class="crop-box"
          :style="cropBoxStyle"
          ref="cropBox"
        >
          <div class="crop-handle nw" @mousedown.stop="startResize('nw', $event)"></div>
          <div class="crop-handle ne" @mousedown.stop="startResize('ne', $event)"></div>
          <div class="crop-handle sw" @mousedown.stop="startResize('sw', $event)"></div>
          <div class="crop-handle se" @mousedown.stop="startResize('se', $event)"></div>
          <div class="crop-handle n" @mousedown.stop="startResize('n', $event)"></div>
          <div class="crop-handle s" @mousedown.stop="startResize('s', $event)"></div>
          <div class="crop-handle w" @mousedown.stop="startResize('w', $event)"></div>
          <div class="crop-handle e" @mousedown.stop="startResize('e', $event)"></div>
          <div class="crop-actions">
            <button class="crop-btn apply" @click.stop="confirmCrop">✓ 确认</button>
            <button class="crop-btn cancel" @click.stop="cancelCrop">✕ 取消</button>
          </div>
        </div>
      </div>
      <div v-if="layers.length === 0" class="empty-state">
        <div class="empty-icon">🖼️</div>
        <div class="empty-text">请上传图片开始编辑</div>
        <div class="empty-hint">点击左上角"上传图片"按钮，或拖拽图片到此处</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { buildFilterString, clearCanvas, downloadCanvas } from '../utils/canvas.js'

const props = defineProps({
  layers: {
    type: Array,
    default: () => []
  },
  activeLayerId: {
    type: [Number, String],
    default: null
  },
  scale: {
    type: Number,
    default: 100
  },
  rotation: {
    type: Number,
    default: 0
  },
  flipH: {
    type: Boolean,
    default: false
  },
  flipV: {
    type: Boolean,
    default: false
  },
  filters: {
    type: Object,
    default: () => ({})
  },
  annotations: {
    type: Array,
    default: () => []
  },
  isCropping: {
    type: Boolean,
    default: false
  },
  currentTool: {
    type: String,
    default: 'none'
  },
  annotationSettings: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['crop-complete', 'add-annotation', 'update-annotation'])

const containerRef = ref(null)
const wrapperRef = ref(null)
const mainCanvas = ref(null)
const cropOverlay = ref(null)
const cropBox = ref(null)

const cropArea = ref(null)
const isDrawing = ref(false)
const isResizing = ref(false)
const resizeHandle = ref(null)
const startPos = ref({ x: 0, y: 0 })
const currentAnnotation = ref(null)

let annotationCounter = 0

const cropBoxStyle = computed(() => {
  if (!cropArea.value) return {}
  return {
    left: `${cropArea.value.x}px`,
    top: `${cropArea.value.y}px`,
    width: `${cropArea.value.width}px`,
    height: `${cropArea.value.height}px`
  }
})

const getCanvasSize = () => {
  let maxWidth = 0
  let maxHeight = 0
  props.layers.forEach((layer) => {
    if (layer.width > maxWidth) maxWidth = layer.width
    if (layer.height > maxHeight) maxHeight = layer.height
  })
  return { width: maxWidth || 800, height: maxHeight || 600 }
}

const render = () => {
  if (!mainCanvas.value) return
  const canvas = mainCanvas.value
  const ctx = canvas.getContext('2d')
  const size = getCanvasSize()
  canvas.width = size.width
  canvas.height = size.height
  clearCanvas(canvas)

  const filterStr = buildFilterString(props.filters)

  props.layers.forEach((layer) => {
    if (!layer.visible || !layer.image) return
    ctx.save()
    const cx = layer.x + layer.width / 2
    const cy = layer.y + layer.height / 2
    ctx.translate(cx, cy)
    if (props.rotation) {
      ctx.rotate((props.rotation * Math.PI) / 180)
    }
    ctx.scale(props.flipH ? -1 : 1, props.flipV ? -1 : 1)
    ctx.translate(-cx, -cy)
    if (filterStr !== 'none') {
      ctx.filter = filterStr
    }
    const s = props.scale / 100
    const w = layer.width * s
    const h = layer.height * s
    const dx = layer.x + (layer.width - w) / 2
    const dy = layer.y + (layer.height - h) / 2
    ctx.drawImage(layer.image, dx, dy, w, h)
    ctx.restore()
  })

  props.annotations.forEach((ann) => {
    drawAnnotation(ctx, ann)
  })

  if (currentAnnotation.value) {
    drawAnnotation(ctx, currentAnnotation.value)
  }
}

const drawAnnotation = (ctx, ann) => {
  ctx.save()
  if (ann.type === 'line') {
    ctx.strokeStyle = ann.color
    ctx.lineWidth = ann.lineWidth
    ctx.lineCap = 'round'
    ctx.beginPath()
    ctx.moveTo(ann.x1, ann.y1)
    ctx.lineTo(ann.x2, ann.y2)
    ctx.stroke()
  } else if (ann.type === 'rect') {
    ctx.strokeStyle = ann.color
    ctx.lineWidth = ann.lineWidth
    const x = Math.min(ann.x1, ann.x2)
    const y = Math.min(ann.y1, ann.y2)
    const w = Math.abs(ann.x2 - ann.x1)
    const h = Math.abs(ann.y2 - ann.y1)
    ctx.strokeRect(x, y, w, h)
  } else if (ann.type === 'circle') {
    ctx.strokeStyle = ann.color
    ctx.lineWidth = ann.lineWidth
    const cx = (ann.x1 + ann.x2) / 2
    const cy = (ann.y1 + ann.y2) / 2
    const rx = Math.abs(ann.x2 - ann.x1) / 2
    const ry = Math.abs(ann.y2 - ann.y1) / 2
    const r = Math.max(rx, ry)
    ctx.beginPath()
    ctx.arc(cx, cy, r, 0, Math.PI * 2)
    ctx.stroke()
  } else if (ann.type === 'text') {
    ctx.fillStyle = ann.color
    ctx.font = `${ann.fontSize}px ${ann.fontFamily}`
    ctx.textBaseline = 'top'
    ctx.fillText(ann.text, ann.x, ann.y)
  }
  ctx.restore()
}

const getCanvasCoords = (e) => {
  const canvas = mainCanvas.value
  const rect = canvas.getBoundingClientRect()
  const scaleX = canvas.width / rect.width
  const scaleY = canvas.height / rect.height
  return {
    x: (e.clientX - rect.left) * scaleX,
    y: (e.clientY - rect.top) * scaleY
  }
}

const handleMouseDown = (e) => {
  if (props.currentTool === 'none' && !props.isCropping) return
  if (props.isCropping) return

  const coords = getCanvasCoords(e)
  isDrawing.value = true
  startPos.value = coords

  if (props.currentTool === 'text') {
    const text = props.annotationSettings.text || '请输入文字'
    currentAnnotation.value = {
      id: ++annotationCounter,
      type: 'text',
      x: coords.x,
      y: coords.y,
      text,
      color: props.annotationSettings.color,
      fontSize: props.annotationSettings.fontSize,
      fontFamily: props.annotationSettings.fontFamily
    }
    emit('add-annotation', { ...currentAnnotation.value })
    currentAnnotation.value = null
    isDrawing.value = false
    render()
    return
  }

  if (props.currentTool === 'line') {
    currentAnnotation.value = {
      id: ++annotationCounter,
      type: 'line',
      x1: coords.x,
      y1: coords.y,
      x2: coords.x,
      y2: coords.y,
      color: props.annotationSettings.color,
      lineWidth: props.annotationSettings.lineWidth
    }
  } else if (props.currentTool === 'rect') {
    currentAnnotation.value = {
      id: ++annotationCounter,
      type: 'rect',
      x1: coords.x,
      y1: coords.y,
      x2: coords.x,
      y2: coords.y,
      color: props.annotationSettings.color,
      lineWidth: props.annotationSettings.lineWidth
    }
  } else if (props.currentTool === 'circle') {
    currentAnnotation.value = {
      id: ++annotationCounter,
      type: 'circle',
      x1: coords.x,
      y1: coords.y,
      x2: coords.x,
      y2: coords.y,
      color: props.annotationSettings.color,
      lineWidth: props.annotationSettings.lineWidth
    }
  }
}

const handleMouseMove = (e) => {
  if (!isDrawing.value || !currentAnnotation.value) return
  const coords = getCanvasCoords(e)
  if (currentAnnotation.value.type === 'line' ||
      currentAnnotation.value.type === 'rect' ||
      currentAnnotation.value.type === 'circle') {
    currentAnnotation.value.x2 = coords.x
    currentAnnotation.value.y2 = coords.y
    render()
  }
}

const handleMouseUp = () => {
  if (!isDrawing.value) return
  if (currentAnnotation.value) {
    emit('add-annotation', { ...currentAnnotation.value })
  }
  isDrawing.value = false
  currentAnnotation.value = null
  render()
}

const startCrop = (e) => {
  if (cropArea.value) return
  const rect = cropOverlay.value.getBoundingClientRect()
  cropArea.value = {
    x: e.clientX - rect.left,
    y: e.clientY - rect.top,
    width: 0,
    height: 0
  }
  isDrawing.value = true
  startPos.value = { x: e.clientX - rect.left, y: e.clientY - rect.top }

  const moveHandler = (ev) => {
    if (!isDrawing.value) return
    const cx = ev.clientX - rect.left
    const cy = ev.clientY - rect.top
    cropArea.value.x = Math.min(startPos.value.x, cx)
    cropArea.value.y = Math.min(startPos.value.y, cy)
    cropArea.value.width = Math.abs(cx - startPos.value.x)
    cropArea.value.height = Math.abs(cy - startPos.value.y)
  }

  const upHandler = () => {
    isDrawing.value = false
    document.removeEventListener('mousemove', moveHandler)
    document.removeEventListener('mouseup', upHandler)
  }

  document.addEventListener('mousemove', moveHandler)
  document.addEventListener('mouseup', upHandler)
}

const startResize = (handle, e) => {
  isResizing.value = true
  resizeHandle.value = handle
  startPos.value = { x: e.clientX, y: e.clientY }
  const startCrop = { ...cropArea.value }
  const rect = cropOverlay.value.getBoundingClientRect()

  const moveHandler = (ev) => {
    if (!isResizing.value) return
    const dx = ev.clientX - startPos.value.x
    const dy = ev.clientY - startPos.value.y
    let { x, y, width, height } = startCrop

    if (handle.includes('e')) width = Math.max(10, startCrop.width + dx)
    if (handle.includes('w')) {
      x = startCrop.x + dx
      width = Math.max(10, startCrop.width - dx)
    }
    if (handle.includes('s')) height = Math.max(10, startCrop.height + dy)
    if (handle.includes('n')) {
      y = startCrop.y + dy
      height = Math.max(10, startCrop.height - dy)
    }

    cropArea.value = { x, y, width, height }
  }

  const upHandler = () => {
    isResizing.value = false
    resizeHandle.value = null
    document.removeEventListener('mousemove', moveHandler)
    document.removeEventListener('mouseup', upHandler)
  }

  document.addEventListener('mousemove', moveHandler)
  document.addEventListener('mouseup', upHandler)
}

const confirmCrop = () => {
  if (!cropArea.value || !mainCanvas.value) return
  const canvas = mainCanvas.value
  const overlay = cropOverlay.value
  const orect = overlay.getBoundingClientRect()
  const crect = canvas.getBoundingClientRect()
  const scaleX = canvas.width / crect.width
  const scaleY = canvas.height / crect.height
  const offsetX = (crect.left - orect.left) * scaleX
  const offsetY = (crect.top - orect.top) * scaleY

  emit('crop-complete', {
    x: Math.max(0, (cropArea.value.x * scaleX) - offsetX),
    y: Math.max(0, (cropArea.value.y * scaleY) - offsetY),
    width: cropArea.value.width * scaleX,
    height: cropArea.value.height * scaleY
  })
  cropArea.value = null
}

const cancelCrop = () => {
  cropArea.value = null
}

const applyCrop = (cropData) => {
  if (!mainCanvas.value) return
  const canvas = mainCanvas.value
  const ctx = canvas.getContext('2d')
  const tempCanvas = document.createElement('canvas')
  tempCanvas.width = cropData.width
  tempCanvas.height = cropData.height
  const tempCtx = tempCanvas.getContext('2d')
  tempCtx.drawImage(
    canvas,
    cropData.x,
    cropData.y,
    cropData.width,
    cropData.height,
    0,
    0,
    cropData.width,
    cropData.height
  )
  canvas.width = cropData.width
  canvas.height = cropData.height
  ctx.drawImage(tempCanvas, 0, 0)

  props.layers.forEach((layer) => {
    if (layer.image) {
      layer.x = Math.max(0, layer.x - cropData.x)
      layer.y = Math.max(0, layer.y - cropData.y)
    }
  })
  render()
}

const exportImage = (format) => {
  if (!mainCanvas.value) return
  const ext = format === 'jpeg' ? 'jpg' : 'png'
  const mimeType = format === 'jpeg' ? 'image/jpeg' : 'image/png'
  downloadCanvas(mainCanvas.value, `edited-image.${ext}`, mimeType)
}

defineExpose({ applyCrop, exportImage, render })

watch(
  () => [props.layers, props.scale, props.rotation, props.flipH, props.flipV, props.filters, props.annotations],
  () => {
    nextTick(() => render())
  },
  { deep: true }
)

onMounted(() => {
  render()
})
</script>

<style scoped>
.editor-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;
  background:
    linear-gradient(45deg, #1e1e36 25%, transparent 25%),
    linear-gradient(-45deg, #1e1e36 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #1e1e36 75%),
    linear-gradient(-45deg, transparent 75%, #1e1e36 75%);
  background-size: 20px 20px;
  background-position: 0 0, 0 10px, 10px -10px, -10px 0px;
  background-color: #16162a;
  padding: 24px;
}

.canvas-wrapper {
  position: relative;
  display: inline-block;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
  border-radius: 8px;
  overflow: hidden;
}

.main-canvas {
  display: block;
  max-width: 100%;
  max-height: 100%;
  cursor: crosshair;
}

.crop-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  cursor: crosshair;
}

.crop-box {
  position: absolute;
  border: 2px dashed #6366f1;
  box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.3);
  cursor: move;
}

.crop-handle {
  position: absolute;
  width: 12px;
  height: 12px;
  background: #6366f1;
  border: 2px solid white;
  border-radius: 50%;
  z-index: 10;
}

.crop-handle.nw { top: -6px; left: -6px; cursor: nwse-resize; }
.crop-handle.ne { top: -6px; right: -6px; cursor: nesw-resize; }
.crop-handle.sw { bottom: -6px; left: -6px; cursor: nesw-resize; }
.crop-handle.se { bottom: -6px; right: -6px; cursor: nwse-resize; }
.crop-handle.n { top: -6px; left: 50%; transform: translateX(-50%); cursor: ns-resize; }
.crop-handle.s { bottom: -6px; left: 50%; transform: translateX(-50%); cursor: ns-resize; }
.crop-handle.w { top: 50%; left: -6px; transform: translateY(-50%); cursor: ew-resize; }
.crop-handle.e { top: 50%; right: -6px; transform: translateY(-50%); cursor: ew-resize; }

.crop-actions {
  position: absolute;
  top: -44px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
  white-space: nowrap;
}

.crop-btn {
  padding: 6px 14px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
}

.crop-btn.apply {
  background: #22c55e;
  color: white;
}

.crop-btn.apply:hover {
  background: #16a34a;
}

.crop-btn.cancel {
  background: #ef4444;
  color: white;
}

.crop-btn.cancel:hover {
  background: #dc2626;
}

.empty-state {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: #666;
  pointer-events: none;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-text {
  font-size: 18px;
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 14px;
  opacity: 0.7;
}
</style>
