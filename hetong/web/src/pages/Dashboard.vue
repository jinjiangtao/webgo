<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useDashboardStore } from '../stores/dashboard'
import SideNav from '../components/SideNav.vue'
import Toolbar from '../components/Toolbar.vue'
import FilterPanel from '../components/FilterPanel.vue'
import DrillBreadcrumb from '../components/DrillBreadcrumb.vue'
import KpiCard from '../components/KpiCard.vue'
import BarChart from '../components/charts/BarChart.vue'
import LineChart from '../components/charts/LineChart.vue'
import PieChart from '../components/charts/PieChart.vue'
import { useTheme } from '../composables/useTheme'

const store = useDashboardStore()
const { isDark } = useTheme()

const kpiConfig = [
  { title: '总销售额', metric: 'sales', icon: '💰', color: 'rgba(59, 130, 246, 0.5)' },
  { title: '订单总量', metric: 'orders', icon: '📦', color: 'rgba(16, 185, 129, 0.5)' },
  { title: '用户总数', metric: 'users', icon: '👥', color: 'rgba(139, 92, 246, 0.5)' },
  { title: '客单价', metric: 'amount', icon: '💵', color: 'rgba(245, 158, 11, 0.5)' },
]

const displayMetrics = computed(() => {
  return store.metrics.slice(0, 4)
})

const activeKpiConfig = computed(() => {
  return kpiConfig.filter(k => store.metrics.includes(k.metric))
})

const getNextDimension = computed(() => {
  const usedDims = new Set(store.drillPath.map(d => d.dimension))
  const allDims: ('time' | 'region' | 'business')[] = ['time', 'region', 'business']
  return allDims.find(d => !usedDims.has(d)) || 'time'
})

onMounted(() => {
  document.documentElement.classList.add('dark')
  store.loadDimensions()
  store.loadData()
})
</script>

<template>
  <div class="dashboard-container" :class="{ dark: isDark }">
    <SideNav />

    <div class="main-content">
      <Toolbar />

      <div class="content-area bg-grid">
        <div class="content-wrapper">
          <FilterPanel />

          <DrillBreadcrumb />

          <div v-loading="store.loading" class="loading-wrapper">
            <div class="kpi-grid">
              <KpiCard
                v-for="kpi in activeKpiConfig"
                :key="kpi.metric"
                :title="kpi.title"
                :metric="kpi.metric"
                :icon="kpi.icon"
                :color="kpi.color"
              />
            </div>

            <div class="charts-grid">
              <div class="chart-card large">
                <BarChart
                  title="销售趋势分析"
                  :metric="store.metrics[0] || 'sales'"
                  :dimension="getNextDimension"
                />
              </div>

              <div class="chart-card medium">
                <LineChart
                  title="同比环比趋势"
                  :metric="store.metrics[0] || 'sales'"
                  :dimension="getNextDimension"
                />
              </div>

              <div class="chart-card small">
                <PieChart
                  title="业务分布"
                  :metric="store.metrics[0] || 'sales'"
                  dimension="business"
                />
              </div>

              <div class="chart-card small">
                <PieChart
                  title="区域分布"
                  :metric="store.metrics[0] || 'sales'"
                  dimension="region"
                />
              </div>
            </div>

            <div class="data-table-section glass glass-hover">
              <div class="table-header">
                <h3 class="table-title">📊 数据明细</h3>
                <div class="table-stats">
                  <span class="stat-item">共 {{ store.records.length }} 条记录</span>
                  <span class="stat-item anomaly" v-if="store.records.some(r => r.anomaly?.isAnomaly)">
                    ⚠ {{ store.records.filter(r => r.anomaly?.isAnomaly).length }} 条异常
                  </span>
                </div>
              </div>

              <div class="table-container">
                <table class="data-table">
                  <thead>
                    <tr>
                      <th>维度</th>
                      <th v-for="m in store.metrics" :key="m" class="text-right">
                        {{ m === 'sales' ? '销售额' : m === 'orders' ? '订单量' : m === 'users' ? '用户数' : '客单价' }}
                      </th>
                      <th class="text-right">同比</th>
                      <th class="text-right">环比</th>
                      <th class="text-center">状态</th>
                      <th class="text-center">操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="(record, index) in store.records"
                      :key="index"
                      :class="{ 'anomaly-row': record.anomaly?.isAnomaly }"
                    >
                      <td class="dimension-cell">
                        <span class="dimension-icon">
                          {{ getNextDimension === 'time' ? '📅' : getNextDimension === 'region' ? '📍' : '💼' }}
                        </span>
                        <span class="dimension-value">
                          {{ record.dimensions[getNextDimension as any] || '总计' }}
                        </span>
                      </td>
                      <td v-for="m in store.metrics" :key="m" class="text-right metric-cell">
                        {{ (record.metrics[m] || 0).toLocaleString() }}
                      </td>
                      <td class="text-right">
                        <span
                          :class="(record.comparison?.yoyPercent || 0) >= 0 ? 'positive' : 'negative'"
                        >
                          {{ (record.comparison?.yoyPercent || 0) >= 0 ? '↑' : '↓' }}
                          {{ Math.abs(record.comparison?.yoyPercent || 0).toFixed(2) }}%
                        </span>
                      </td>
                      <td class="text-right">
                        <span
                          :class="(record.comparison?.momPercent || 0) >= 0 ? 'positive' : 'negative'"
                        >
                          {{ (record.comparison?.momPercent || 0) >= 0 ? '↑' : '↓' }}
                          {{ Math.abs(record.comparison?.momPercent || 0).toFixed(2) }}%
                        </span>
                      </td>
                      <td class="text-center">
                        <span
                          v-if="record.anomaly?.isAnomaly"
                          class="anomaly-tag"
                          :title="record.anomaly.reason"
                        >
                          ⚠ 异常
                        </span>
                        <span v-else class="normal-tag">✓ 正常</span>
                      </td>
                      <td class="text-center">
                        <button
                          v-if="record.canDrillDown"
                          class="drill-btn"
                          @click="store.drillDown(
                            getNextDimension as any,
                            record.dimensions[getNextDimension as any] || '',
                            record.dimensions[getNextDimension as any] || ''
                          )"
                        >
                          下钻 →
                        </button>
                        <span v-else class="no-drill">-</span>
                      </td>
                    </tr>
                    <tr v-if="store.records.length === 0">
                      <td :colspan="store.metrics.length + 5" class="empty-row">
                        <div class="empty-content">
                          <span class="empty-icon">📭</span>
                          <span class="empty-text">暂无数据</span>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background: var(--bg-primary);
  transition: background 0.3s ease;
}

.main-content {
  margin-left: 240px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.content-wrapper {
  max-width: 1800px;
  margin: 0 auto;
}

.loading-wrapper {
  min-height: 400px;
}

.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.chart-card {
  height: 380px;
}

.chart-card.large {
  grid-column: span 3;
}

.chart-card.medium {
  grid-column: span 3;
}

.chart-card.small {
  grid-column: span 3;
  height: 320px;
}

@media (max-width: 1200px) {
  .charts-grid {
    grid-template-columns: repeat(1, 1fr);
  }

  .chart-card.large,
  .chart-card.medium,
  .chart-card.small {
    grid-column: span 1;
  }
}

.data-table-section {
  padding: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
}

.table-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  font-family: 'JetBrains Mono', monospace;
}

.table-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}

.stat-item.anomaly {
  color: var(--neon-amber);
  animation: pulse 2s infinite;
}

.table-container {
  overflow-x: auto;
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  font-family: 'JetBrains Mono', monospace;
}

.data-table th {
  background: rgba(30, 41, 59, 0.9);
  padding: 12px 16px;
  text-align: left;
  color: var(--text-primary);
  font-weight: 600;
  border-bottom: 2px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 10;
}

.data-table td {
  padding: 12px 16px;
  border-bottom: 1px solid rgba(59, 130, 246, 0.1);
  color: var(--text-secondary);
  transition: background 0.2s ease;
}

.data-table tbody tr:hover {
  background: rgba(59, 130, 246, 0.08);
}

.data-table tbody tr.anomaly-row {
  background: rgba(239, 68, 68, 0.05);
}

.data-table tbody tr.anomaly-row:hover {
  background: rgba(239, 68, 68, 0.1);
}

.text-right {
  text-align: right;
}

.text-center {
  text-align: center;
}

.dimension-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dimension-icon {
  font-size: 16px;
}

.dimension-value {
  color: var(--text-primary);
  font-weight: 500;
}

.metric-cell {
  color: var(--text-primary);
  font-weight: 500;
}

.positive {
  color: var(--neon-green);
  font-weight: 500;
}

.negative {
  color: var(--neon-red);
  font-weight: 500;
}

.anomaly-tag {
  display: inline-block;
  padding: 2px 8px;
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid var(--neon-red);
  border-radius: 4px;
  color: var(--neon-red);
  font-size: 11px;
  font-weight: 500;
  animation: pulse 2s infinite;
  cursor: help;
}

.normal-tag {
  display: inline-block;
  padding: 2px 8px;
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid var(--neon-green);
  border-radius: 4px;
  color: var(--neon-green);
  font-size: 11px;
  font-weight: 500;
}

.drill-btn {
  padding: 4px 12px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(6, 182, 212, 0.2));
  border: 1px solid var(--neon-blue);
  border-radius: 4px;
  color: var(--neon-blue);
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.drill-btn:hover {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.4), rgba(6, 182, 212, 0.4));
  box-shadow: 0 0 15px rgba(59, 130, 246, 0.3);
}

.no-drill {
  color: var(--text-muted);
}

.empty-row {
  text-align: center;
  padding: 60px 20px !important;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.empty-icon {
  font-size: 48px;
  opacity: 0.5;
}

.empty-text {
  font-size: 14px;
  color: var(--text-muted);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}
</style>
