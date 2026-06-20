<script setup>
defineProps({
  toasts: { type: Array, default: () => [] }
})
const emit = defineEmits(['acknowledge', 'close'])
</script>

<template>
  <div class="toast-container">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      class="toast"
      :class="toast.level"
    >
      <div class="toast-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
          <line x1="12" y1="9" x2="12" y2="13"/>
          <line x1="12" y1="17" x2="12.01" y2="17"/>
        </svg>
      </div>
      <div class="toast-body">
        <div class="toast-title">
          <span class="toast-tag">{{ toast.level === 'critical' ? '严重告警' : '警告' }}</span>
          {{ toast.device_name }}
        </div>
        <div class="toast-msg">{{ toast.message }}</div>
        <div class="toast-time">{{ toast.created_at }}</div>
      </div>
      <div class="toast-actions">
        <button class="toast-btn ack" @click="emit('acknowledge', toast)">确认</button>
        <button class="toast-btn close" @click="emit('close', toast)">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toast-container {
  position: fixed;
  top: 80px;
  right: 24px;
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 380px;
  pointer-events: none;
}
.toast {
  display: flex;
  gap: 12px;
  padding: 14px;
  background: var(--bg-elevated);
  border: 1px solid var(--abnormal);
  border-left: 4px solid var(--abnormal);
  border-radius: var(--radius);
  box-shadow: 0 8px 32px rgba(0,0,0,0.4), 0 0 24px var(--abnormal-glow);
  animation: slide-in-right 0.4s cubic-bezier(0.16, 1, 0.3, 1);
  pointer-events: auto;
}
.toast.warning {
  border-color: var(--warning);
  box-shadow: 0 8px 32px rgba(0,0,0,0.4), 0 0 24px rgba(251,191,36,0.2);
}
.toast-icon {
  color: var(--abnormal);
  flex-shrink: 0;
  animation: blink 1.5s infinite;
}
.toast.warning .toast-icon { color: var(--warning); }
.toast-body {
  flex: 1;
  min-width: 0;
}
.toast-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
}
.toast-tag {
  font-size: 10px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 3px;
  background: var(--abnormal-bg);
  color: var(--abnormal);
  letter-spacing: 0.05em;
}
.toast.warning .toast-tag {
  background: var(--warning-bg);
  color: var(--warning);
}
.toast-msg {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.4;
}
.toast-time {
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
  margin-top: 6px;
}
.toast-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-end;
}
.toast-btn {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 4px;
  transition: all 0.15s;
}
.toast-btn.ack {
  background: var(--abnormal);
  color: #fff;
}
.toast-btn.ack:hover { background: #fca5a5; }
.toast-btn.close {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}
.toast-btn.close:hover { color: var(--text-primary); }
</style>
