<script setup>
import { ref, watch } from 'vue'
import { api } from '../api'

const props = defineProps({ show: { type: Boolean, default: false } })
const emit = defineEmits(['close'])

const alarms = ref([])
const tab = ref('active')

async function load() {
  alarms.value = await api.getAlarms(tab.value)
}
watch(() => props.show, (v) => { if (v) load() })
watch(tab, () => load())

async function ack(id) {
  await api.acknowledgeAlarm(id)
  await load()
}

function levelText(l) { return l === 'critical' ? '严重' : '警告' }
</script>

<template>
  <div v-if="show" class="modal-overlay" @click.self="emit('close')">
    <div class="modal modal-wide">
      <div class="modal-head">
        <h2>告警记录</h2>
        <button class="modal-close" @click="emit('close')">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
      <div class="modal-body">
        <div class="tabs">
          <button :class="{ active: tab === 'active' }" @click="tab = 'active'">活动中</button>
          <button :class="{ active: tab === 'acknowledged' }" @click="tab = 'acknowledged'">已确认</button>
          <button :class="{ active: tab === 'resolved' }" @click="tab = 'resolved'">已恢复</button>
          <button :class="{ active: tab === '' }" @click="tab = ''">全部</button>
        </div>
        <div class="alarm-list">
          <div v-if="alarms.length === 0" class="empty">暂无告警记录</div>
          <div v-for="a in alarms" :key="a.id" class="alarm-row" :class="a.level">
            <div class="alarm-left">
              <span class="alarm-level" :class="a.level">{{ levelText(a.level) }}</span>
              <div class="alarm-info">
                <div class="alarm-dev">{{ a.device_name }}</div>
                <div class="alarm-msg">{{ a.message }}</div>
                <div class="alarm-time">{{ a.created_at }}<span v-if="a.acknowledged_at"> · 确认于 {{ a.acknowledged_at }}</span></div>
              </div>
            </div>
            <button v-if="a.status === 'active'" class="btn-ack" @click="ack(a.id)">确认处理</button>
            <span v-else class="status-badge" :class="a.status">{{ a.status === 'acknowledged' ? '已确认' : '已恢复' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.6); backdrop-filter: blur(4px);
  z-index: 200; display: flex; align-items: center; justify-content: center;
  padding: 20px; animation: fade-in 0.2s;
}
.modal {
  background: var(--bg-surface); border: 1px solid var(--border-bright);
  border-radius: var(--radius-lg); width: 100%; animation: scale-in 0.25s;
  box-shadow: 0 24px 60px rgba(0,0,0,0.5);
}
.modal-wide { max-width: 720px; }
.modal-head {
  display: flex; justify-content: space-between; align-items: center;
  padding: 18px 22px; border-bottom: 1px solid var(--border);
}
.modal-head h2 { font-size: 16px; font-weight: 600; }
.modal-close {
  width: 32px; height: 32px; display: flex; align-items: center; justify-content: center;
  border-radius: 6px; color: var(--text-muted);
}
.modal-close:hover { background: var(--bg-elevated); color: var(--text-primary); }
.modal-body { padding: 22px; display: flex; flex-direction: column; gap: 14px; }
.tabs { display: flex; gap: 4px; }
.tabs button {
  padding: 7px 14px; font-size: 12px; border-radius: var(--radius);
  background: var(--bg-base); border: 1px solid var(--border); color: var(--text-secondary);
}
.tabs button.active { color: var(--accent); border-color: var(--accent-dim); }
.alarm-list { display: flex; flex-direction: column; gap: 8px; max-height: 460px; overflow-y: auto; }
.alarm-row {
  display: flex; justify-content: space-between; align-items: center; gap: 12px;
  padding: 14px; background: var(--bg-base); border: 1px solid var(--border);
  border-left: 3px solid var(--border-bright); border-radius: var(--radius);
}
.alarm-row.critical { border-left-color: var(--abnormal); }
.alarm-row.warning { border-left-color: var(--warning); }
.alarm-left { display: flex; gap: 12px; align-items: flex-start; min-width: 0; flex: 1; }
.alarm-level {
  font-size: 10px; font-weight: 700; padding: 3px 8px; border-radius: 3px;
  letter-spacing: 0.05em; flex-shrink: 0; margin-top: 1px;
}
.alarm-level.critical { background: var(--abnormal-bg); color: var(--abnormal); }
.alarm-level.warning { background: var(--warning-bg); color: var(--warning); }
.alarm-info { min-width: 0; }
.alarm-dev { font-size: 14px; font-weight: 600; margin-bottom: 3px; }
.alarm-msg { font-size: 12px; color: var(--text-secondary); margin-bottom: 3px; }
.alarm-time { font-family: var(--font-mono); font-size: 10px; color: var(--text-muted); }
.btn-ack {
  padding: 7px 14px; font-size: 12px; border-radius: var(--radius);
  background: var(--accent); color: #04141a; font-weight: 600; white-space: nowrap; flex-shrink: 0;
}
.btn-ack:hover { background: #67e8f9; }
.status-badge {
  font-size: 11px; padding: 4px 10px; border-radius: 3px; flex-shrink: 0;
  background: var(--bg-elevated); color: var(--text-muted);
}
.status-badge.resolved { color: var(--online); background: var(--online-bg); }
.empty { text-align: center; color: var(--text-muted); padding: 40px; }
</style>
