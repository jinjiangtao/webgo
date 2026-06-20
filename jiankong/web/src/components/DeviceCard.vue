<script setup>
import { computed } from 'vue'

const props = defineProps({
  device: { type: Object, required: true }
})
const emit = defineEmits(['edit', 'delete', 'view'])

const statusConfig = {
  online: { label: '在线', color: 'var(--online)' },
  offline: { label: '离线', color: 'var(--offline)' },
  abnormal: { label: '异常', color: 'var(--abnormal)' }
}

const cfg = computed(() => statusConfig[props.device.status] || statusConfig.offline)

const cpuOver = computed(() => props.device.cpu_usage >= props.device.threshold_cpu)
const tempOver = computed(() => props.device.temperature >= props.device.threshold_temp)

function fmt(n) {
  return Number(n).toFixed(1)
}
function barWidth(n) {
  return Math.min(100, Math.max(0, n)) + '%'
}
function barColor(n, threshold) {
  if (n >= threshold) return 'var(--abnormal)'
  if (n >= threshold * 0.8) return 'var(--warning)'
  return 'var(--accent)'
}
</script>

<template>
  <div class="card" :class="device.status" @click="emit('view', device)">
    <div class="card-header">
      <div class="dev-title">
        <span class="status-dot" :class="device.status"></span>
        <span class="dev-name">{{ device.name }}</span>
      </div>
      <span class="dev-status" :style="{ color: cfg.color, borderColor: cfg.color }">{{ cfg.label }}</span>
    </div>

    <div class="card-meta">
      <span class="badge">{{ device.type_name || '未分类' }}</span>
      <span class="meta-code">{{ device.device_code }}</span>
    </div>

    <div class="card-location">
      <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
        <circle cx="12" cy="10" r="3"/>
      </svg>
      <span>{{ device.location || '未设置' }}</span>
    </div>

    <div class="metrics">
      <div class="metric">
        <div class="metric-head">
          <span class="metric-label">CPU</span>
          <span class="metric-value" :class="{ over: cpuOver }">{{ fmt(device.cpu_usage) }}%</span>
        </div>
        <div class="bar-track">
          <div class="bar-fill" :style="{ width: barWidth(device.cpu_usage), background: barColor(device.cpu_usage, device.threshold_cpu) }"></div>
        </div>
      </div>
      <div class="metric">
        <div class="metric-head">
          <span class="metric-label">内存</span>
          <span class="metric-value">{{ fmt(device.memory_usage) }}%</span>
        </div>
        <div class="bar-track">
          <div class="bar-fill" :style="{ width: barWidth(device.memory_usage), background: barColor(device.memory_usage, 95) }"></div>
        </div>
      </div>
      <div class="metric">
        <div class="metric-head">
          <span class="metric-label">温度</span>
          <span class="metric-value" :class="{ over: tempOver }">{{ fmt(device.temperature) }}°C</span>
        </div>
        <div class="bar-track">
          <div class="bar-fill" :style="{ width: barWidth(device.temperature), background: barColor(device.temperature, device.threshold_temp) }"></div>
        </div>
      </div>
    </div>

    <div class="card-footer">
      <span class="heartbeat">
        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 16 14"/>
        </svg>
        {{ device.last_heartbeat || '无心跳' }}
      </span>
      <div class="card-actions">
        <button class="icon-btn" title="编辑" @click.stop="emit('edit', device)">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
          </svg>
        </button>
        <button class="icon-btn danger" title="删除" @click.stop="emit('delete', device)">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"/>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: var(--bg-surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: all 0.25s ease;
  position: relative;
  overflow: hidden;
  cursor: pointer;
}
.card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--border-bright);
  transition: background 0.25s;
}
.card.online::before { background: var(--online); }
.card.offline::before { background: var(--offline); opacity: 0.5; }
.card.abnormal::before { background: var(--abnormal); }

.card.online:hover {
  border-color: var(--online);
  box-shadow: 0 0 20px var(--online-glow);
}
.card.offline:hover {
  border-color: var(--border-bright);
}
.card.abnormal {
  border-color: rgba(248, 113, 113, 0.3);
  animation: pulse-abnormal 2s infinite;
}
.card.abnormal:hover {
  border-color: var(--abnormal);
  box-shadow: 0 0 24px var(--abnormal-glow);
}
.card.offline {
  opacity: 0.65;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}
.dev-title {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  position: relative;
}
.status-dot.online { background: var(--online); box-shadow: 0 0 8px var(--online-glow); }
.status-dot.offline { background: var(--offline); }
.status-dot.abnormal { background: var(--abnormal); box-shadow: 0 0 10px var(--abnormal-glow); animation: blink 1s infinite; }
.dev-name {
  font-weight: 600;
  font-size: 15px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.dev-status {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 3px;
  border: 1px solid;
  letter-spacing: 0.05em;
  flex-shrink: 0;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}
.badge {
  font-size: 11px;
  padding: 2px 8px;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: 3px;
  color: var(--text-secondary);
}
.meta-code {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}
.card-location {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: var(--text-muted);
}

.metrics {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 2px;
}
.metric-head {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}
.metric-label {
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.metric-value {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
}
.metric-value.over { color: var(--abnormal); }
.bar-track {
  height: 4px;
  background: var(--bg-base);
  border-radius: 2px;
  overflow: hidden;
}
.bar-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.6s ease, background 0.3s;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 8px;
  border-top: 1px solid var(--border);
}
.heartbeat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: var(--font-mono);
  font-size: 10px;
  color: var(--text-muted);
}
.card-actions {
  display: flex;
  gap: 4px;
}
.icon-btn {
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  color: var(--text-muted);
  transition: all 0.15s;
}
.icon-btn:hover {
  background: var(--bg-elevated);
  color: var(--accent);
}
.icon-btn.danger:hover {
  color: var(--abnormal);
}
</style>
