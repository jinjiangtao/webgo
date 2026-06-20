<script setup>
const props = defineProps({
  deviceTypes: { type: Array, default: () => [] },
  filters: { type: Object, required: true },
  resultCount: { type: Number, default: 0 }
})
const emit = defineEmits(['update:filters'])

const statusOptions = [
  { value: '', label: '全部状态' },
  { value: 'online', label: '在线' },
  { value: 'offline', label: '离线' },
  { value: 'abnormal', label: '异常' }
]

function update(key, value) {
  emit('update:filters', { ...props.filters, [key]: value })
}
</script>

<template>
  <div class="filter-bar">
    <div class="filter-left">
      <div class="filter-group">
        <label>设备类型</label>
        <select :value="filters.device_type" @change="update('device_type', $event.target.value)">
          <option value="">全部类型</option>
          <option v-for="t in deviceTypes" :key="t.id" :value="t.id">{{ t.name }}</option>
        </select>
      </div>
      <div class="filter-group">
        <label>运行状态</label>
        <div class="status-tabs">
          <button
            v-for="opt in statusOptions"
            :key="opt.value"
            class="status-tab"
            :class="[opt.value, { active: filters.status === opt.value }]"
            @click="update('status', opt.value)"
          >{{ opt.label }}</button>
        </div>
      </div>
      <div class="filter-group search-group">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input
          type="text"
          placeholder="搜索设备名称 / 编号 / 位置"
          :value="filters.keyword"
          @input="update('keyword', $event.target.value)"
        />
      </div>
    </div>
    <div class="filter-result">
      共 <span class="count">{{ resultCount }}</span> 台设备
    </div>
  </div>
</template>

<style scoped>
.filter-bar {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 24px;
  border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
}
.filter-left {
  display: flex;
  gap: 18px;
  flex-wrap: wrap;
  align-items: flex-end;
}
.filter-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.filter-group label {
  font-size: 10px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}
select {
  background: var(--bg-surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 7px 28px 7px 10px;
  font-size: 13px;
  outline: none;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%2394a8bd' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 8px center;
  cursor: pointer;
}
select:focus { border-color: var(--accent-dim); }

.status-tabs {
  display: flex;
  background: var(--bg-surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}
.status-tab {
  padding: 7px 14px;
  font-size: 12px;
  color: var(--text-secondary);
  border-right: 1px solid var(--border);
  transition: all 0.15s;
}
.status-tab:last-child { border-right: none; }
.status-tab:hover { background: var(--bg-elevated); }
.status-tab.active { font-weight: 600; }
.status-tab.active.online { color: var(--online); background: var(--online-bg); }
.status-tab.active.offline { color: var(--offline); background: var(--offline-bg); }
.status-tab.active.abnormal { color: var(--abnormal); background: var(--abnormal-bg); }
.status-tab.active:not(.online):not(.offline):not(.abnormal) { color: var(--accent); background: rgba(34,211,238,0.08); }

.search-group {
  position: relative;
}
.search-group svg {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  pointer-events: none;
}
.search-group input {
  background: var(--bg-surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 7px 10px 7px 32px;
  font-size: 13px;
  width: 240px;
  outline: none;
}
.search-group input:focus { border-color: var(--accent-dim); }
.search-group input::placeholder { color: var(--text-muted); }

.filter-result {
  font-size: 12px;
  color: var(--text-muted);
}
.filter-result .count {
  font-family: var(--font-mono);
  color: var(--accent);
  font-weight: 700;
  font-size: 14px;
}

@media (max-width: 768px) {
  .filter-bar { padding: 12px 16px; }
  .search-group input { width: 100%; }
  .filter-left { gap: 12px; }
}
</style>
