<script setup>
defineProps({
  stats: { type: Object, default: () => ({}) },
  autoRefresh: { type: Boolean, default: false }
})
const emit = defineEmits(['toggle-refresh', 'add-device', 'manage-types', 'view-logs', 'view-alarms'])
</script>

<template>
  <header class="topbar">
    <div class="brand">
      <div class="brand-icon">
        <svg width="22" height="22" viewBox="0 0 32 32" fill="none">
          <path d="M16 4 L26 10 L26 22 L16 28 L6 22 L6 10 Z" stroke="var(--accent)" stroke-width="1.5"/>
          <circle cx="16" cy="16" r="3" fill="var(--online)"/>
          <circle cx="16" cy="16" r="6" stroke="var(--online)" stroke-width="1" opacity="0.4"/>
        </svg>
      </div>
      <div class="brand-text">
        <h1>设备运维监控</h1>
        <span class="brand-sub">DEVICE OPS CENTER</span>
      </div>
    </div>

    <div class="stat-group">
      <div class="stat-chip total">
        <span class="stat-num">{{ stats.total ?? '-' }}</span>
        <span class="stat-label">设备总数</span>
      </div>
      <div class="stat-chip online">
        <span class="stat-num">{{ stats.online ?? '-' }}</span>
        <span class="stat-label">在线</span>
      </div>
      <div class="stat-chip offline">
        <span class="stat-num">{{ stats.offline ?? '-' }}</span>
        <span class="stat-label">离线</span>
      </div>
      <div class="stat-chip abnormal">
        <span class="stat-num">{{ stats.abnormal ?? '-' }}</span>
        <span class="stat-label">异常</span>
      </div>
      <button class="stat-chip alarm-btn" :class="{ active: stats.alarms > 0 }" @click="emit('view-alarms')">
        <span class="alarm-dot" v-if="stats.alarms > 0"></span>
        <span class="stat-num">{{ stats.alarms ?? '-' }}</span>
        <span class="stat-label">活动告警</span>
      </button>
    </div>

    <div class="actions">
      <button class="btn-toggle" :class="{ on: autoRefresh }" @click="emit('toggle-refresh')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23 4 23 10 17 10"/>
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
        </svg>
        <span>{{ autoRefresh ? '自动刷新' : '已暂停' }}</span>
      </button>
      <button class="btn-ghost" @click="emit('view-logs')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
          <line x1="8" y1="13" x2="16" y2="13"/>
          <line x1="8" y1="17" x2="16" y2="17"/>
        </svg>
        日志
      </button>
      <button class="btn-ghost" @click="emit('manage-types')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="8" y1="6" x2="21" y2="6"/>
          <line x1="8" y1="12" x2="21" y2="12"/>
          <line x1="8" y1="18" x2="21" y2="18"/>
          <line x1="3" y1="6" x2="3.01" y2="6"/>
          <line x1="3" y1="12" x2="3.01" y2="12"/>
          <line x1="3" y1="18" x2="3.01" y2="18"/>
        </svg>
        类型
      </button>
      <button class="btn-primary" @click="emit('add-device')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        录入设备
      </button>
    </div>
  </header>
</template>

<style scoped>
.topbar {
  position: sticky;
  top: 0;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  padding: 14px 24px;
  background: rgba(7, 11, 17, 0.85);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border);
}
.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}
.brand-icon {
  display: flex;
  filter: drop-shadow(0 0 6px var(--accent-glow));
}
.brand-text h1 {
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0.02em;
}
.brand-sub {
  font-family: var(--font-mono);
  font-size: 9px;
  color: var(--accent);
  letter-spacing: 0.2em;
}
.stat-group {
  display: flex;
  gap: 8px;
}
.stat-chip {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px 14px;
  border-radius: var(--radius);
  background: var(--bg-surface);
  border: 1px solid var(--border);
  min-width: 64px;
}
.stat-num {
  font-family: var(--font-mono);
  font-size: 18px;
  font-weight: 700;
  line-height: 1.2;
}
.stat-label {
  font-size: 10px;
  color: var(--text-muted);
  letter-spacing: 0.03em;
}
.stat-chip.total .stat-num { color: var(--text-primary); }
.stat-chip.online .stat-num { color: var(--online); }
.stat-chip.offline .stat-num { color: var(--offline); }
.stat-chip.abnormal .stat-num { color: var(--abnormal); }
.alarm-btn {
  cursor: pointer;
  transition: all 0.2s;
}
.alarm-btn.active {
  border-color: var(--abnormal);
  background: var(--abnormal-bg);
}
.alarm-btn.active .stat-num { color: var(--abnormal); }
.alarm-dot {
  width: 6px;
  height: 6px;
  background: var(--abnormal);
  border-radius: 50%;
  position: absolute;
  margin-left: 32px;
  margin-top: -4px;
  animation: blink 1s infinite;
}

.actions {
  display: flex;
  gap: 8px;
  align-items: center;
}
.btn-toggle, .btn-ghost, .btn-primary {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: var(--radius);
  font-size: 13px;
  font-weight: 500;
  transition: all 0.15s;
  white-space: nowrap;
}
.btn-toggle {
  border: 1px solid var(--border);
  background: var(--bg-surface);
  color: var(--text-secondary);
}
.btn-toggle.on {
  border-color: var(--accent-dim);
  color: var(--accent);
  background: rgba(34, 211, 238, 0.08);
}
.btn-toggle svg { transition: transform 0.3s; }
.btn-toggle.on svg { animation: spin 3s linear infinite; }
.btn-ghost {
  border: 1px solid var(--border);
  background: var(--bg-surface);
  color: var(--text-secondary);
}
.btn-ghost:hover {
  border-color: var(--border-bright);
  color: var(--text-primary);
}
.btn-primary {
  background: var(--accent);
  color: #04141a;
  font-weight: 600;
}
.btn-primary:hover {
  background: #67e8f9;
  box-shadow: 0 0 16px var(--accent-glow);
}

@media (max-width: 1100px) {
  .stat-group { display: none; }
}
@media (max-width: 768px) {
  .topbar { padding: 12px 16px; gap: 10px; flex-wrap: wrap; }
  .brand-text h1 { font-size: 15px; }
  .btn-ghost span, .btn-toggle span { display: none; }
}
</style>
