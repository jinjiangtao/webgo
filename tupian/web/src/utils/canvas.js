export const createCanvas = (width, height) => {
  const canvas = document.createElement('canvas')
  canvas.width = width
  canvas.height = height
  return canvas
}

export const getContext = (canvas) => {
  return canvas.getContext('2d')
}

export const clearCanvas = (canvas) => {
  const ctx = getContext(canvas)
  ctx.clearRect(0, 0, canvas.width, canvas.height)
}

export const drawImage = (ctx, image, x, y, width, height, options = {}) => {
  ctx.save()
  if (options.rotation) {
    const cx = x + width / 2
    const cy = y + height / 2
    ctx.translate(cx, cy)
    ctx.rotate((options.rotation * Math.PI) / 180)
    ctx.translate(-cx, -cy)
  }
  if (options.flipH || options.flipV) {
    const cx = x + width / 2
    const cy = y + height / 2
    ctx.translate(cx, cy)
    ctx.scale(options.flipH ? -1 : 1, options.flipV ? -1 : 1)
    ctx.translate(-cx, -cy)
  }
  if (options.filter) {
    ctx.filter = options.filter
  }
  ctx.drawImage(image, x, y, width, height)
  ctx.restore()
}

export const applyFilters = (ctx, filters) => {
  const filterParts = []
  if (filters.grayscale > 0) {
    filterParts.push(`grayscale(${filters.grayscale}%)`)
  }
  if (filters.sepia > 0) {
    filterParts.push(`sepia(${filters.sepia}%)`)
  }
  if (filters.invert > 0) {
    filterParts.push(`invert(${filters.invert}%)`)
  }
  if (filters.brightness !== 100) {
    filterParts.push(`brightness(${filters.brightness}%)`)
  }
  if (filters.contrast !== 100) {
    filterParts.push(`contrast(${filters.contrast}%)`)
  }
  if (filters.saturate !== 100) {
    filterParts.push(`saturate(${filters.saturate}%)`)
  }
  if (filters.blur > 0) {
    filterParts.push(`blur(${filters.blur}px)`)
  }
  if (filters.gaussianBlur > 0) {
    filterParts.push(`blur(${filters.gaussianBlur}px)`)
  }
  if (filterParts.length > 0) {
    ctx.filter = filterParts.join(' ')
  } else {
    ctx.filter = 'none'
  }
  return ctx.filter
}

export const buildFilterString = (filters) => {
  const parts = []
  if (filters.grayscale > 0) parts.push(`grayscale(${filters.grayscale}%)`)
  if (filters.sepia > 0) parts.push(`sepia(${filters.sepia}%)`)
  if (filters.invert > 0) parts.push(`invert(${filters.invert}%)`)
  if (filters.brightness !== 100) parts.push(`brightness(${filters.brightness}%)`)
  if (filters.contrast !== 100) parts.push(`contrast(${filters.contrast}%)`)
  if (filters.saturate !== 100) parts.push(`saturate(${filters.saturate}%)`)
  if (filters.blur > 0) parts.push(`blur(${filters.blur}px)`)
  if (filters.gaussianBlur > 0) parts.push(`blur(${filters.gaussianBlur}px)`)
  return parts.length > 0 ? parts.join(' ') : 'none'
}

export const drawLine = (ctx, x1, y1, x2, y2, color = '#ef4444', lineWidth = 3) => {
  ctx.save()
  ctx.strokeStyle = color
  ctx.lineWidth = lineWidth
  ctx.lineCap = 'round'
  ctx.beginPath()
  ctx.moveTo(x1, y1)
  ctx.lineTo(x2, y2)
  ctx.stroke()
  ctx.restore()
}

export const drawRect = (ctx, x, y, width, height, color = '#ef4444', lineWidth = 3) => {
  ctx.save()
  ctx.strokeStyle = color
  ctx.lineWidth = lineWidth
  ctx.strokeRect(x, y, width, height)
  ctx.restore()
}

export const drawCircle = (ctx, cx, cy, radius, color = '#ef4444', lineWidth = 3) => {
  ctx.save()
  ctx.strokeStyle = color
  ctx.lineWidth = lineWidth
  ctx.beginPath()
  ctx.arc(cx, cy, radius, 0, Math.PI * 2)
  ctx.stroke()
  ctx.restore()
}

export const drawText = (ctx, text, x, y, options = {}) => {
  const {
    color = '#ef4444',
    fontSize = 24,
    fontFamily = 'Arial'
  } = options
  ctx.save()
  ctx.fillStyle = color
  ctx.font = `${fontSize}px ${fontFamily}`
  ctx.textBaseline = 'top'
  ctx.fillText(text, x, y)
  ctx.restore()
}

export const cropCanvas = (canvas, x, y, width, height) => {
  const cropped = createCanvas(width, height)
  const ctx = getContext(cropped)
  ctx.drawImage(canvas, x, y, width, height, 0, 0, width, height)
  return cropped
}

export const rotateCanvas = (canvas, angle) => {
  const ctx = getContext(canvas)
  const radians = (angle * Math.PI) / 180
  const sin = Math.abs(Math.sin(radians))
  const cos = Math.abs(Math.cos(radians))
  const newWidth = canvas.width * cos + canvas.height * sin
  const newHeight = canvas.width * sin + canvas.height * cos
  const rotated = createCanvas(newWidth, newHeight)
  const rctx = getContext(rotated)
  rctx.translate(newWidth / 2, newHeight / 2)
  rctx.rotate(radians)
  rctx.drawImage(canvas, -canvas.width / 2, -canvas.height / 2)
  return rotated
}

export const flipCanvas = (canvas, horizontal = false, vertical = false) => {
  const flipped = createCanvas(canvas.width, canvas.height)
  const ctx = getContext(flipped)
  ctx.translate(
    horizontal ? canvas.width : 0,
    vertical ? canvas.height : 0
  )
  ctx.scale(
    horizontal ? -1 : 1,
    vertical ? -1 : 1
  )
  ctx.drawImage(canvas, 0, 0)
  return flipped
}

export const canvasToBlob = (canvas, format = 'image/png', quality = 0.92) => {
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (blob) {
          resolve(blob)
        } else {
          reject(new Error('Failed to convert canvas to blob'))
        }
      },
      format,
      quality
    )
  })
}

export const downloadCanvas = (canvas, filename, format = 'image/png') => {
  const link = document.createElement('a')
  link.download = filename
  link.href = canvas.toDataURL(format, 0.92)
  link.click()
}

export const resizeCanvas = (canvas, targetWidth, targetHeight) => {
  const resized = createCanvas(targetWidth, targetHeight)
  const ctx = getContext(resized)
  ctx.drawImage(canvas, 0, 0, targetWidth, targetHeight)
  return resized
}

export const scaleCanvas = (canvas, scale) => {
  const newWidth = Math.round(canvas.width * scale)
  const newHeight = Math.round(canvas.height * scale)
  return resizeCanvas(canvas, newWidth, newHeight)
}
