export class WSClient {
  constructor() {
    this.ws = null
    this.url = ''
    this.handlers = new Map()
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 10
    this.baseReconnectDelay = 1000
    this.reconnectTimer = null
    this.heartbeatInterval = null
    this.heartbeatTimeout = null
    this.heartbeatIntervalMs = 30000
    this.heartbeatTimeoutMs = 10000
    this.isManualClose = false
    this.pendingMessages = []
  }

  connect(url) {
    return new Promise((resolve, reject) => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        resolve()
        return
      }

      this.url = url
      this.isManualClose = false

      try {
        this.ws = new WebSocket(url)

        this.ws.onopen = () => {
          this.reconnectAttempts = 0
          this.flushPendingMessages()
          this.startHeartbeat()
          this.emit('connect')
          resolve()
        }

        this.ws.onmessage = (event) => {
          this.handleMessage(event)
        }

        this.ws.onerror = (error) => {
          this.emit('error', error)
          reject(error)
        }

        this.ws.onclose = (event) => {
          this.stopHeartbeat()
          this.emit('disconnect', event)

          if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.scheduleReconnect()
          }
        }
      } catch (error) {
        this.emit('error', error)
        reject(error)
      }
    })
  }

  disconnect() {
    this.isManualClose = true
    this.stopHeartbeat()
    this.clearReconnectTimer()

    if (this.ws) {
      this.ws.close(1000, 'Manual disconnect')
      this.ws = null
    }
  }

  send(type, payload = {}) {
    const message = JSON.stringify({ type, payload })

    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(message)
    } else {
      this.pendingMessages.push(message)
    }
  }

  on(type, handler) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, [])
    }
    this.handlers.get(type).push(handler)
  }

  off(type, handler) {
    if (!this.handlers.has(type)) return

    if (handler) {
      const handlers = this.handlers.get(type)
      const index = handlers.indexOf(handler)
      if (index !== -1) {
        handlers.splice(index, 1)
      }
    } else {
      this.handlers.delete(type)
    }
  }

  handleMessage(event) {
    try {
      const data = JSON.parse(event.data)
      const { type, payload } = data

      if (type === 'pong') {
        this.resetHeartbeatTimeout()
        return
      }

      this.emit(type, payload)
    } catch (error) {
      this.emit('error', error)
    }
  }

  emit(type, payload) {
    if (this.handlers.has(type)) {
      this.handlers.get(type).forEach((handler) => {
        try {
          handler(payload)
        } catch (error) {
          console.error(`[WSClient] Handler error for "${type}":`, error)
        }
      })
    }
  }

  scheduleReconnect() {
    this.reconnectAttempts++
    const delay = this.getReconnectDelay()

    this.emit('reconnecting', {
      attempt: this.reconnectAttempts,
      delay
    })

    this.reconnectTimer = setTimeout(() => {
      this.connect(this.url).catch((error) => {
        console.error('[WSClient] Reconnect failed:', error)
      })
    }, delay)
  }

  getReconnectDelay() {
    const delay = this.baseReconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
    const maxDelay = 30000
    return Math.min(delay, maxDelay)
  }

  clearReconnectTimer() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
  }

  flushPendingMessages() {
    while (this.pendingMessages.length > 0) {
      const message = this.pendingMessages.shift()
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(message)
      }
    }
  }

  startHeartbeat() {
    this.heartbeatInterval = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.send('ping', { timestamp: Date.now() })

        this.heartbeatTimeout = setTimeout(() => {
          if (this.ws) {
            this.ws.close(1006, 'Heartbeat timeout')
          }
        }, this.heartbeatTimeoutMs)
      }
    }, this.heartbeatIntervalMs)
  }

  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
    this.resetHeartbeatTimeout()
  }

  resetHeartbeatTimeout() {
    if (this.heartbeatTimeout) {
      clearTimeout(this.heartbeatTimeout)
      this.heartbeatTimeout = null
    }
  }

  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN
  }
}

let client = null

export const getWSClient = () => {
  if (!client) {
    client = new WSClient()
  }
  return client
}
