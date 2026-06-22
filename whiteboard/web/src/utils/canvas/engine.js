const TOOL_TYPES = {
  PEN: 'pen',
  LINE: 'line',
  RECT: 'rect',
  CIRCLE: 'circle',
  ELLIPSE: 'ellipse',
  TRIANGLE: 'triangle',
  TEXT: 'text',
  ERASER: 'eraser'
}

class CanvasEngine {
  constructor(canvas) {
    this.canvas = canvas
    this.ctx = canvas.getContext('2d')
    this.tool = TOOL_TYPES.PEN
    this.strokeWidth = 2
    this.strokeColor = '#000000'
    this.fillColor = 'transparent'
    this.opacity = 1
    this.fontSize = 16
    this.fontFamily = 'Arial'
    this.backgroundColor = '#ffffff'
    this.scale = 1
    this.offsetX = 0
    this.offsetY = 0
    this._currentShape = null
    this._isDrawing = false
    this._textPending = false
  }

  _generateId() {
    return 'shape_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
  }

  setTool(tool) {
    this.tool = tool
  }

  setStrokeWidth(width) {
    this.strokeWidth = width
  }

  setStrokeColor(color) {
    this.strokeColor = color
  }

  setFillColor(color) {
    this.fillColor = color
  }

  setOpacity(opacity) {
    this.opacity = Math.max(0, Math.min(1, opacity))
  }

  setFontSize(size) {
    this.fontSize = size
  }

  setFontFamily(family) {
    this.fontFamily = family
  }

  setBackground(color) {
    this.backgroundColor = color
    this._paintBackground()
  }

  _paintBackground() {
    const prevTransform = this.ctx.getTransform()
    this.ctx.setTransform(1, 0, 0, 1, 0, 0)
    this.ctx.save()
    this.ctx.fillStyle = this.backgroundColor
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height)
    this.ctx.restore()
    this.ctx.setTransform(prevTransform)
  }

  applyZoom(scale, offsetX, offsetY) {
    this.scale = scale
    this.offsetX = offsetX
    this.offsetY = offsetY
    this.ctx.setTransform(scale, 0, 0, scale, offsetX, offsetY)
  }

  _applyStyles(ctx = this.ctx) {
    ctx.lineWidth = this.strokeWidth
    ctx.strokeStyle = this.strokeColor
    ctx.fillStyle = this.fillColor
    ctx.globalAlpha = this.opacity
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'
    ctx.font = `${this.fontSize}px ${this.fontFamily}`
  }

  _applyShapeStyles(shape, ctx = this.ctx) {
    ctx.lineWidth = shape.strokeWidth
    ctx.strokeStyle = shape.strokeColor
    ctx.fillStyle = shape.fillColor || 'transparent'
    ctx.globalAlpha = shape.opacity
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'
    if (shape.fontSize && shape.fontFamily) {
      ctx.font = `${shape.fontSize}px ${shape.fontFamily}`
    }
  }

  startDraw(x, y) {
    if (this._textPending) return null

    this._isDrawing = true
    const baseShape = {
      id: this._generateId(),
      type: this.tool,
      points: [{ x, y }],
      strokeWidth: this.strokeWidth,
      strokeColor: this.strokeColor,
      fillColor: this.fillColor,
      opacity: this.opacity
    }

    if (this.tool === TOOL_TYPES.TEXT) {
      baseShape.fontSize = this.fontSize
      baseShape.fontFamily = this.fontFamily
      baseShape.text = ''
      this._currentShape = baseShape
      this._textPending = true
      this._isDrawing = false
      return this._currentShape
    }

    if (this.tool === TOOL_TYPES.ERASER) {
      baseShape.strokeColor = this.backgroundColor
      baseShape.strokeWidth = this.strokeWidth * 3
    }

    this._currentShape = baseShape
    return this._currentShape
  }

  moveDraw(x, y) {
    if (!this._isDrawing || !this._currentShape) return null

    const shape = this._currentShape
    const lastPoint = shape.points[shape.points.length - 1]

    switch (shape.type) {
      case TOOL_TYPES.PEN:
      case TOOL_TYPES.ERASER:
        shape.points.push({ x, y })
        this._drawPathSegment(lastPoint, { x, y }, shape)
        break
      case TOOL_TYPES.LINE:
        shape.points[1] = { x, y }
        break
      case TOOL_TYPES.RECT:
      case TOOL_TYPES.CIRCLE:
      case TOOL_TYPES.ELLIPSE:
      case TOOL_TYPES.TRIANGLE:
        shape.points[1] = { x, y }
        break
      default:
        break
    }

    return shape
  }

  endDraw(x, y) {
    if (!this._currentShape) return null

    const shape = this._currentShape

    if (shape.type === TOOL_TYPES.TEXT) {
      this._currentShape = null
      return shape
    }

    if (!this._isDrawing) {
      this._currentShape = null
      return null
    }

    if (shape.points.length < 2 && shape.type !== TOOL_TYPES.PEN && shape.type !== TOOL_TYPES.ERASER) {
      shape.points.push({ x: x + 0.1, y: y + 0.1 })
    }

    if (shape.type === TOOL_TYPES.PEN || shape.type === TOOL_TYPES.ERASER) {
      if (shape.points.length < 2) {
        shape.points.push({ x: x + 0.1, y: y + 0.1 })
      }
    }

    this._isDrawing = false
    this._currentShape = null
    return shape
  }

  setText(text) {
    if (!this._textPending || !this._currentShape) return null
    this._currentShape.text = text
    const shape = this._currentShape
    this._currentShape = null
    this._textPending = false
    return shape
  }

  cancelText() {
    this._textPending = false
    this._currentShape = null
  }

  isTextPending() {
    return this._textPending
  }

  _drawPathSegment(from, to, shape) {
    const ctx = this.ctx
    this._applyShapeStyles(shape, ctx)
    ctx.beginPath()
    ctx.moveTo(from.x, from.y)
    ctx.lineTo(to.x, to.y)
    ctx.stroke()
  }

  renderShape(shape) {
    if (!shape) return
    const ctx = this.ctx
    ctx.save()
    this._applyShapeStyles(shape, ctx)

    switch (shape.type) {
      case TOOL_TYPES.PEN:
      case TOOL_TYPES.ERASER:
        this._renderPen(shape, ctx)
        break
      case TOOL_TYPES.LINE:
        this._renderLine(shape, ctx)
        break
      case TOOL_TYPES.RECT:
        this._renderRect(shape, ctx)
        break
      case TOOL_TYPES.CIRCLE:
        this._renderCircle(shape, ctx)
        break
      case TOOL_TYPES.ELLIPSE:
        this._renderEllipse(shape, ctx)
        break
      case TOOL_TYPES.TRIANGLE:
        this._renderTriangle(shape, ctx)
        break
      case TOOL_TYPES.TEXT:
        this._renderText(shape, ctx)
        break
      default:
        break
    }

    ctx.restore()
  }

  _renderPen(shape, ctx) {
    if (shape.points.length < 2) return
    ctx.beginPath()
    ctx.moveTo(shape.points[0].x, shape.points[0].y)
    for (let i = 1; i < shape.points.length; i++) {
      ctx.lineTo(shape.points[i].x, shape.points[i].y)
    }
    ctx.stroke()
  }

  _renderLine(shape, ctx) {
    if (shape.points.length < 2) return
    ctx.beginPath()
    ctx.moveTo(shape.points[0].x, shape.points[0].y)
    ctx.lineTo(shape.points[1].x, shape.points[1].y)
    ctx.stroke()
  }

  _renderRect(shape, ctx) {
    if (shape.points.length < 2) return
    const p1 = shape.points[0]
    const p2 = shape.points[1]
    const x = Math.min(p1.x, p2.x)
    const y = Math.min(p1.y, p2.y)
    const w = Math.abs(p2.x - p1.x)
    const h = Math.abs(p2.y - p1.y)
    ctx.beginPath()
    ctx.rect(x, y, w, h)
    if (shape.fillColor && shape.fillColor !== 'transparent') {
      ctx.fill()
    }
    ctx.stroke()
  }

  _renderCircle(shape, ctx) {
    if (shape.points.length < 2) return
    const p1 = shape.points[0]
    const p2 = shape.points[1]
    const cx = (p1.x + p2.x) / 2
    const cy = (p1.y + p2.y) / 2
    const r = Math.sqrt(Math.pow(p2.x - p1.x, 2) + Math.pow(p2.y - p1.y, 2)) / 2
    ctx.beginPath()
    ctx.arc(cx, cy, r, 0, Math.PI * 2)
    if (shape.fillColor && shape.fillColor !== 'transparent') {
      ctx.fill()
    }
    ctx.stroke()
  }

  _renderEllipse(shape, ctx) {
    if (shape.points.length < 2) return
    const p1 = shape.points[0]
    const p2 = shape.points[1]
    const cx = (p1.x + p2.x) / 2
    const cy = (p1.y + p2.y) / 2
    const rx = Math.abs(p2.x - p1.x) / 2
    const ry = Math.abs(p2.y - p1.y) / 2
    ctx.beginPath()
    ctx.ellipse(cx, cy, rx, ry, 0, 0, Math.PI * 2)
    if (shape.fillColor && shape.fillColor !== 'transparent') {
      ctx.fill()
    }
    ctx.stroke()
  }

  _renderTriangle(shape, ctx) {
    if (shape.points.length < 2) return
    const p1 = shape.points[0]
    const p2 = shape.points[1]
    const leftX = Math.min(p1.x, p2.x)
    const rightX = Math.max(p1.x, p2.x)
    const bottomY = Math.max(p1.y, p2.y)
    const topY = Math.min(p1.y, p2.y)
    const apexX = (leftX + rightX) / 2
    ctx.beginPath()
    ctx.moveTo(apexX, topY)
    ctx.lineTo(rightX, bottomY)
    ctx.lineTo(leftX, bottomY)
    ctx.closePath()
    if (shape.fillColor && shape.fillColor !== 'transparent') {
      ctx.fill()
    }
    ctx.stroke()
  }

  _renderText(shape, ctx) {
    if (!shape.points.length || !shape.text) return
    ctx.fillStyle = shape.strokeColor
    ctx.globalAlpha = shape.opacity
    ctx.fillText(shape.text, shape.points[0].x, shape.points[0].y)
  }

  renderShapes(shapes) {
    this.clear()
    this._paintBackground()
    shapes.forEach(shape => this.renderShape(shape))
  }

  clear() {
    const prevTransform = this.ctx.getTransform()
    this.ctx.setTransform(1, 0, 0, 1, 0, 0)
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
    this.ctx.setTransform(prevTransform)
  }

  resize(width, height) {
    this.canvas.width = width
    this.canvas.height = height
    this.applyZoom(this.scale, this.offsetX, this.offsetY)
    this._paintBackground()
  }
}

export { CanvasEngine, TOOL_TYPES }
export default CanvasEngine
