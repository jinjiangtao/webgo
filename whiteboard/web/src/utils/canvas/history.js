class HistoryManager {
  constructor(maxHistory = 100) {
    this._history = []
    this._index = -1
    this._maxHistory = maxHistory
  }

  push(shape) {
    if (!shape) return
    if (this._index < this._history.length - 1) {
      this._history = this._history.slice(0, this._index + 1)
    }
    this._history.push(shape)
    if (this._history.length > this._maxHistory) {
      this._history.shift()
    } else {
      this._index++
    }
  }

  undo() {
    if (!this.canUndo()) return null
    this._index--
    return this.getAllShapes()
  }

  redo() {
    if (!this.canRedo()) return null
    this._index++
    return this.getAllShapes()
  }

  canUndo() {
    return this._index >= 0
  }

  canRedo() {
    return this._index < this._history.length - 1
  }

  getAllShapes() {
    return this._history.slice(0, this._index + 1)
  }

  setShapes(shapes) {
    this._history = shapes ? [...shapes] : []
    this._index = this._history.length - 1
  }

  clear() {
    this._history = []
    this._index = -1
  }

  getHistoryLength() {
    return this._history.length
  }

  getIndex() {
    return this._index
  }
}

export { HistoryManager }
export default HistoryManager
