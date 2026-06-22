import { ref, reactive, computed, onUnmounted } from 'vue'
import { CanvasEngine, TOOL_TYPES } from '../utils/canvas/engine'
import { HistoryManager } from '../utils/canvas/history'

export function useCanvas() {
  const canvasRef = ref(null)
  const engine = ref(null)
  const history = ref(new HistoryManager(200))
  const shapes = ref([])
  const isDrawing = ref(false)
  const isTextPending = ref(false)
  const textPosition = ref({ x: 0, y: 0 })
  const previewShape = ref(null)

  const config = reactive({
    tool: TOOL_TYPES.PEN,
    strokeWidth: 2,
    strokeColor: '#000000',
    fillColor: 'transparent',
    opacity: 1,
    fontSize: 16,
    fontFamily: 'Arial',
    backgroundColor: '#ffffff',
    scale: 1,
    offsetX: 0,
    offsetY: 0
  })

  const canUndo = computed(() => history.value.canUndo())
  const canRedo = computed(() => history.value.canRedo())

  function initCanvas(canvas) {
    if (!canvas) return
    canvasRef.value = canvas
    engine.value = new CanvasEngine(canvas)
    _applyConfigToEngine()
    renderAll()
  }

  function _applyConfigToEngine() {
    if (!engine.value) return
    engine.value.setTool(config.tool)
    engine.value.setStrokeWidth(config.strokeWidth)
    engine.value.setStrokeColor(config.strokeColor)
    engine.value.setFillColor(config.fillColor)
    engine.value.setOpacity(config.opacity)
    engine.value.setFontSize(config.fontSize)
    engine.value.setFontFamily(config.fontFamily)
  }

  function setTool(tool) {
    config.tool = tool
    engine.value?.setTool(tool)
  }

  function setStrokeWidth(width) {
    config.strokeWidth = width
    engine.value?.setStrokeWidth(width)
  }

  function setStrokeColor(color) {
    config.strokeColor = color
    engine.value?.setStrokeColor(color)
  }

  function setFillColor(color) {
    config.fillColor = color
    engine.value?.setFillColor(color)
  }

  function setOpacity(opacity) {
    config.opacity = opacity
    engine.value?.setOpacity(opacity)
  }

  function setFontSize(size) {
    config.fontSize = size
    engine.value?.setFontSize(size)
  }

  function setFontFamily(family) {
    config.fontFamily = family
    engine.value?.setFontFamily(family)
  }

  function setBackground(color) {
    config.backgroundColor = color
    engine.value?.setBackground(color)
  }

  function applyZoom(scale, offsetX, offsetY) {
    config.scale = scale
    config.offsetX = offsetX
    config.offsetY = offsetY
    engine.value?.applyZoom(scale, offsetX, offsetY)
    renderAll()
  }

  function _getCanvasCoords(e) {
    if (!canvasRef.value) return { x: 0, y: 0 }
    const rect = canvasRef.value.getBoundingClientRect()
    const clientX = e.touches ? e.touches[0].clientX : e.clientX
    const clientY = e.touches ? e.touches[0].clientY : e.clientY
    const x = (clientX - rect.left - config.offsetX) / config.scale
    const y = (clientY - rect.top - config.offsetY) / config.scale
    return { x, y }
  }

  function handleStart(e) {
    if (!engine.value) return
    e.preventDefault?.()
    const { x, y } = _getCanvasCoords(e)
    const shape = engine.value.startDraw(x, y)

    if (config.tool === TOOL_TYPES.TEXT && shape) {
      isTextPending.value = true
      textPosition.value = { x, y }
      previewShape.value = null
      return
    }

    if (shape) {
      isDrawing.value = true
      previewShape.value = shape
    }
  }

  function handleMove(e) {
    if (!engine.value || !isDrawing.value) return
    e.preventDefault?.()
    const { x, y } = _getCanvasCoords(e)
    const shape = engine.value.moveDraw(x, y)

    if (shape && shape.type !== TOOL_TYPES.PEN && shape.type !== TOOL_TYPES.ERASER) {
      previewShape.value = { ...shape }
      renderAll()
      engine.value.renderShape(shape)
    } else if (shape) {
      previewShape.value = shape
    }
  }

  function handleEnd(e) {
    if (!engine.value) return
    e.preventDefault?.()
    const { x, y } = _getCanvasCoords(e)

    if (isTextPending.value) {
      return
    }

    if (!isDrawing.value) {
      previewShape.value = null
      return
    }

    const shape = engine.value.endDraw(x, y)
    if (shape) {
      history.value.push(shape)
      shapes.value = history.value.getAllShapes()
    }

    isDrawing.value = false
    previewShape.value = null
    renderAll()
  }

  function confirmText(text) {
    if (!engine.value || !isTextPending.value) return
    const shape = engine.value.setText(text)
    if (shape && shape.text) {
      history.value.push(shape)
      shapes.value = history.value.getAllShapes()
      renderAll()
    }
    isTextPending.value = false
    previewShape.value = null
  }

  function cancelText() {
    engine.value?.cancelText()
    isTextPending.value = false
    previewShape.value = null
  }

  function undo() {
    if (!history.value.canUndo()) return
    history.value.undo()
    shapes.value = history.value.getAllShapes()
    renderAll()
  }

  function redo() {
    if (!history.value.canRedo()) return
    history.value.redo()
    shapes.value = history.value.getAllShapes()
    renderAll()
  }

  function renderAll() {
    if (!engine.value) return
    engine.value.renderShapes(shapes.value)
  }

  function renderShape(shape) {
    engine.value?.renderShape(shape)
  }

  function clearCanvas() {
    history.value.clear()
    shapes.value = []
    engine.value?.clear()
    engine.value?.setBackground(config.backgroundColor)
  }

  function setShapes(newShapes) {
    history.value.setShapes(newShapes)
    shapes.value = history.value.getAllShapes()
    renderAll()
  }

  function getAllShapes() {
    return history.value.getAllShapes()
  }

  function resize(width, height) {
    engine.value?.resize(width, height)
    renderAll()
  }

  onUnmounted(() => {
    engine.value = null
    canvasRef.value = null
  })

  return {
    canvasRef,
    engine,
    config,
    shapes,
    isDrawing,
    isTextPending,
    textPosition,
    previewShape,
    canUndo,
    canRedo,
    TOOL_TYPES,
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
    renderShape,
    clearCanvas,
    setShapes,
    getAllShapes,
    resize
  }
}

export default useCanvas
