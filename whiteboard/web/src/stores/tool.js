import { defineStore } from 'pinia'

export const TOOL_TYPES = {
  PEN: 'pen',
  LINE: 'line',
  RECT: 'rect',
  CIRCLE: 'circle',
  ELLIPSE: 'ellipse',
  TRIANGLE: 'triangle',
  TEXT: 'text',
  ERASER: 'eraser'
}

export const useToolStore = defineStore('tool', {
  state: () => ({
    currentTool: TOOL_TYPES.PEN,
    strokeWidth: 3,
    strokeColor: '#000000',
    fillColor: 'transparent',
    opacity: 100,
    fontSize: 16,
    fontFamily: 'Arial'
  }),
  actions: {
    setCurrentTool(tool) {
      this.currentTool = tool
    },
    setStrokeWidth(width) {
      this.strokeWidth = Math.max(1, Math.min(30, width))
    },
    setStrokeColor(color) {
      this.strokeColor = color
    },
    setFillColor(color) {
      this.fillColor = color
    },
    setOpacity(opacity) {
      this.opacity = Math.max(0, Math.min(100, opacity))
    },
    setFontSize(size) {
      this.fontSize = Math.max(12, Math.min(72, size))
    },
    setFontFamily(family) {
      this.fontFamily = family
    }
  }
})
