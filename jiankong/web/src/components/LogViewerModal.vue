<script setup>
import { ref, watch, computed } from 'vue'
import { api } from '../api'

const props = defineProps({
  show: { type: Boolean, default: false },
  devices: { type: Array, default: () => [] }
})
const emit = defineEmits(['close'])

const logs = ref([])
const loading = ref(false)
const filterDevice = ref('')
const filterLevel = ref('')
const autoScroll = ref(true)

const levelConfig = {
  info: { label: '信息', color: 'var(--accent)' },
  warning: { label: '警告', color: 'var(--warning)' },
  error: { label: '错误', color: 'var(--abnormal)' }
}

const filtered = computed(() => logs.value.filter(l => {
  if (filterDevice.value && l.device_id != filterDevice.value) return false
  if (filterLevel.value && l.level !== filterLevel.value) return false
  return true
}))

async function load() {
  loading.value = true
  try {
    logs.value = await api.getLogs({ limit: 500 })
  } finally {
    loading.value = false
  }
}
watch(() => props.show, (v) => { if (v) load() })

async function clearAll() {
  if (!confirm('确定清空所有日志？此操作不可恢复。')) return
  await api.clearLogs()
  await load()
}
</script>

<template>
  <div v-if="show" class="modal-overlay" @click.self="emit('close')">
    <div class="modal modal-wide">
      <div class="modal-head">
        <h2>设备运行日志</h2>
        <div class="head-actions">
          <button class="btn-refresh" @click="load" :disabled="loading">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
            </svg>
          </button>
          <button class="btn-clear" @click="clearAll">清空</button>
          <button class="modal-close" @click="emit('close')">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
      </div>
      <div class="modal-body">
        <div class="log-filters">
          <select v-model="filterDevice">
            <option value="">全部设备</option>
            <option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }}</option>
          </select>
          <div class="level-tabs">
            <button :class="{ active: !filterLevel }" @click="filterLevel = ''">全部</button>
            <button :class="{ active: filterLevel === 'info' }" @click="filterLevel = 'info'">信息</button>
            <button :class="{ active: filterLevel === 'warning' }" @click="filterLevel = 'warning'">警告</button>
            <button :class="{ active: filterLevel === 'error' }" @click="filterLevel = 'error'">错误</button>
          </div>
        </div>
        <div class="log-list">
          <div v-if="filtered.length === 0 && !loading" class="empty">暂无日志记录</div>
          <div v-for="log in filtered" :key="log.id" class="log-row" :class="log.level">
            <span class="log-time">{{ log.created_at }}</span>
            <span class="log-level" :style="{ color: levelConfig[log.level]?.color }">{{ levelConfig[log.level]?.label }}</span>
            <span class="log-device">{{ log.device_name || '系统' }}</span>
            <span class="log-msg">{{ log.message }}</span>
          </div>
        </div>
        <div class="log-count">共 {{ filtered.length }} 条记录</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  z-index: 200;
  display: flex; align-items: center; justify-content: center;
  padding: 20px; animation: fade-in 0.2s;
}
.modal {
  background: var(--bg-surface);
  border: 1px solid var(--border-bright);
  border-radius: var(--radius-lg);
  width: 100%;
  animation: scale-in 0.25s;
  box-shadow: 0 24px 60px rgba(0,0,0,0.5);
}
.modal-wide { max-width: 860px; }
.modal-head {
  display: flex; justify-content: space-between; align-items: center;
  padding: 18px 22px; border-bottom: 1px solid var(--border);
}
.modal-head h2 { font-size: 16px; font-weight: 600; }
.head-actions { display: flex; gap: 8px; align-items: center; }
.btn-refresh, .btn-clear, .modal-close {
  display: flex; align-items: center; justify-content: center;
  border-radius: 6px; color: var(--text-muted);
}
.btn-refresh { width: 32px; height: 32px; }
.btn-refresh:disabled { opacity: 0.5; }
.btn-refresh:hover { color: var(--accent); }
.btn-clear { padding: 6px 12px; font-size: 12px; border: 1px solid var(--border); }
.btn-clear:hover { color: var(--abnormal); border-color: var(--abnormal); }
.modal-close { width: 32px; height: 32px; }
.modal-close:hover { background: var(--bg-elevated); color: var(--text-primary); }
.modal-body { padding: 22px; display: flex; flex-direction: column; gap: 14px; }
.log-filters { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
.log-filters select {
  background: var(--bg-base); border: 1px solid var(--border);
  border-radius: var(--radius); padding: 7px 10px; font-size: 13px; outline: none;
}
.level-tabs { display: flex; gap: 4px; }
.level-tabs button {
  padding: 6px 12px; font-size: 12px; border-radius: var(--radius);
  background: var(--bg-base); border: 1px solid var(--border); color: var(--text-secondary);
}
.level-tabs button.active { color: var(--accent); border-color: var(--accent-dim); }
.log-list {
  max-height: 440px; overflow-y: auto; background: var(--bg-base);
  border: 1px solid var(--border); border-radius: var(--radius); padding: 6px;
}
.log-row {
  display: grid;
  grid-template-columns: 140px 50px 130px 1fr;
  gap: 12px; padding: 7px 10px; border-bottom: 1px solid rgba(30,45,61,0.5);
  font-size: 12px; align-items: baseline;
}
.log-row:last-child { border-bottom: none; }
.log-row.error { background: rgba(248,113,113,0.04); }
.log-row.warning { background: rgba(251,191,36,0.04); }
.log-time { font-family: var(--font-mono); color: var(--text-muted); }
.log-level { font-weight: 600; text-align: center; }
.log-device { color: var(--text-secondary); }
.log-msg { color: var(--text-primary); }
.empty { text-align: center; color: var(--text-muted); padding: 40px; }
.log-count { font-size: 12px; color: var(--text-muted); text-align: right; }
</style>
