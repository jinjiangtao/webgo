<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useDashboardStore } from '../stores/dashboard'

const store = useDashboardStore()

const timeFilter = ref('')
const regionFilter = ref('')
const businessFilter = ref('')
const metricSelector = ref<string[]>(['sales', 'orders', 'users'])

const metricOptions = [
  { value: 'sales', label: '销售额' },
  { value: 'orders', label: '订单量' },
  { value: 'users', label: '用户数' },
  { value: 'amount', label: '客单价' },
]

const timeOptions = computed(() => [
  { value: '', label: '全部时间' },
  ...store.dimensions.timeOptions.map(d => ({ value: d.value, label: d.label })),
])

const regionOptions = computed(() => [
  { value: '', label: '全部区域' },
  ...store.dimensions.regionOptions.map(d => ({ value: d.value, label: d.label })),
])

const businessOptions = computed(() => [
  { value: '', label: '全部业务' },
  ...store.dimensions.businessOptions.map(d => ({ value: d.value, label: d.label })),
])

watch([timeFilter, regionFilter, businessFilter], () => {
  const filters: Record<string, string> = {}
  if (timeFilter.value) filters.time = timeFilter.value
  if (regionFilter.value) filters.region = regionFilter.value
  if (businessFilter.value) filters.business = businessFilter.value
  store.updateFilters(filters)
})

watch(metricSelector, () => {
  store.updateMetrics(metricSelector.value)
}, { deep: true })

const handleReset = () => {
  timeFilter.value = ''
  regionFilter.value = ''
  businessFilter.value = ''
  metricSelector.value = ['sales', 'orders', 'users']
}
</script>

<template>
  <div class="filter-panel">
    <div class="filter-header">
      <span class="filter-title">🔍 维度筛选</span>
      <button class="reset-btn" @click="handleReset">重置</button>
    </div>

    <div class="filter-content">
      <div class="filter-item">
        <label class="filter-label">时间维度</label>
        <select v-model="timeFilter" class="filter-select">
          <option v-for="opt in timeOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <div class="filter-item">
        <label class="filter-label">区域维度</label>
        <select v-model="regionFilter" class="filter-select">
          <option v-for="opt in regionOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <div class="filter-item">
        <label class="filter-label">业务维度</label>
        <select v-model="businessFilter" class="filter-select">
          <option v-for="opt in businessOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <div class="filter-item metrics">
        <label class="filter-label">指标选择</label>
        <div class="metric-checkboxes">
          <label
            v-for="opt in metricOptions"
            :key="opt.value"
            class="metric-checkbox"
          >
            <input
              type="checkbox"
              :value="opt.value"
              v-model="metricSelector"
            />
            <span>{{ opt.label }}</span>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.filter-panel {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.8), rgba(15, 23, 42, 0.9));
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  backdrop-filter: blur(10px);
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
}

.filter-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  font-family: 'JetBrains Mono', monospace;
}

.reset-btn {
  padding: 4px 12px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 12px;
  font-family: 'JetBrains Mono', monospace;
  transition: all 0.2s ease;
}

.reset-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  border-color: var(--neon-red);
  color: var(--neon-red);
}

.filter-content {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  align-items: flex-start;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 180px;
}

.filter-item.metrics {
  min-width: 300px;
  flex: 1;
}

.filter-label {
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}

.filter-select {
  padding: 8px 12px;
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 13px;
  font-family: 'JetBrains Mono', monospace;
  cursor: pointer;
  transition: all 0.2s ease;
  outline: none;
}

.filter-select:hover,
.filter-select:focus {
  border-color: var(--neon-blue);
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.2);
}

.metric-checkboxes {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.metric-checkbox {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  font-family: 'JetBrains Mono', monospace;
  transition: color 0.2s ease;
}

.metric-checkbox:hover {
  color: var(--text-primary);
}

.metric-checkbox input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: var(--neon-blue);
  cursor: pointer;
}

.metric-checkbox input[type="checkbox"]:checked + span {
  color: var(--neon-blue);
}
</style>
