import { defineStore } from 'pinia'

export const useWhiteboardStore = defineStore('whiteboard', {
  state: () => ({
    whiteboardId: '',
    whiteboardName: '',
    background: '#ffffff',
    onlineUsers: [],
    cursors: new Map(),
    zoom: 1,
    panX: 0,
    panY: 0
  }),
  actions: {
    setWhiteboard(id, name, background) {
      this.whiteboardId = id
      this.whiteboardName = name || this.whiteboardName
      if (background !== undefined) {
        this.background = background
      }
    },
    addUser(user) {
      const exists = this.onlineUsers.find(u => u.userId === user.userId)
      if (!exists) {
        this.onlineUsers.push({
          userId: user.userId,
          username: user.username,
          color: user.color
        })
      }
    },
    removeUser(userId) {
      const index = this.onlineUsers.findIndex(u => u.userId === userId)
      if (index !== -1) {
        this.onlineUsers.splice(index, 1)
      }
      this.removeCursor(userId)
    },
    setUsers(users) {
      this.onlineUsers = users.map(u => ({
        userId: u.userId,
        username: u.username,
        color: u.color
      }))
    },
    setCursor(userId, x, y) {
      this.cursors.set(userId, { x, y, userId })
    },
    removeCursor(userId) {
      this.cursors.delete(userId)
    },
    setZoom(zoom) {
      this.zoom = Math.max(0.25, Math.min(3, zoom))
    },
    setPan(x, y) {
      this.panX = x
      this.panY = y
    },
    setBackground(color) {
      this.background = color
    }
  }
})
