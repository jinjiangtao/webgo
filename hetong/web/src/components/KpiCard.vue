<script setup lang="ts">
import { computed } from 'vue'
import { useDashboardStore } from '../stores/dashboard'

const props = defineProps<{
  title: string
  metric: string
  icon: string
  color: string
}>()

const store = useDashboardStore()

const total = computed(() => {
  return store.records.reduce((sum, r) => sum + (r.metrics[props.metric] || 0), 0)
})

const comparison = computed(() => {
  if (store.records.length === 0) return null
  const totalYoY = store.records.reduce((sum, r) => sum + (r.comparison?.yoyPercent || 0), 0)
  return totalYoY / store.records.length
})

const anomalyCount = computed(() => {
  return store.records.filter(r => r.anomaly?.isAnomaly).length
})

const formatNumber = (num: number) => {
  if (num >= 100000000) {
    return (num / 100000000).toFixed(2) + '亿'
  } else if (num >= 10000) {
    return (num / 10000).toFixed(2) + '万'
  }
  return num.toLocaleString()
}
</script>

<template>
  <div class="kpi-card" :style="{ '--card-color': color }">
    <div class="kpi-header">
      <span class="kpi-icon">{{ icon }}</span>
      <span class="kpi-title">{{ title }}</span>
    </div>
    <div class="kpi-value">{{ formatNumber(total) }}</div>
    <div class="kpi-footer">
      <span
        class="kpi-comparison"
        :class="comparison >= 0 ? 'positive' : 'negative'"
        v-if="comparison !== null"
      >
        {{ comparison >= 0 ? '↑' : '↓' }} {{ Math.abs(comparison).toFixed(2) }}% 同比
      </span>
      <span class="kpi-anomaly" v-if="anomalyCount > 0">
        ⚠ {{ anomalyCount }} 个异常
      </span>
    </div>
    <div class="kpi-glow"></div>
  </div>
</template>

<style scoped>
.kpi-card {
  position: relative;
  padding: 20px;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--card-color);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
  cursor: pointer;
}

.kpi-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 40px var(--card-color);
}

.kpi-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.kpi-icon {
  font-size: 24px;
}

.kpi-title {
  font-size: 14px;
  color: #94A3B8;
  font-family: 'JetBrains Mono', monospace;
}

.kpi-value {
  font-size: 32px;
  font-weight: bold;
  color: #E2E8F0;
  font-family: 'JetBrains Mono', monospace;
  margin-bottom: 12px;
  text-shadow: 0 0 20px var(--card-color);
}

.kpi-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.kpi-comparison {
  font-size: 13px;
  font-family: 'JetBrains Mono', monospace;
}

.kpi-comparison.positive {
  color: #10B981;
}

.kpi-comparison.negative {
  color: #EF4444;
}

.kpi-anomaly {
  font-size: 13px;
  color: #F59E0B;
  font-family: 'JetBrains Mono', monospace;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.kpi-glow {
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, var(--card-color) 0%, transparent 70%);
  opacity: 0.1;
  pointer-events: none;
}
</style>
