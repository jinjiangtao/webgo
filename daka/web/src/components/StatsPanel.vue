<template>
  <div class="stats-overview">
    <div class="stats-header">
      <h2 class="section-title">📊 本月打卡概览</h2>
      <div class="month-selector">
        <button class="nav-btn-sm" @click="prevMonth">←</button>
        <span class="month-text">{{ year }}年{{ month }}月</span>
        <button class="nav-btn-sm" @click="nextMonth">→</button>
      </div>
    </div>

    <div class="stats-grid" v-if="stats">
      <div class="stat-card primary">
        <div class="stat-icon">🎯</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.completion_rate.toFixed(1) }}%</div>
          <div class="stat-label">完成率</div>
        </div>
        <div class="stat-bar">
          <div class="stat-bar-fill" :style="{ width: stats.completion_rate + '%' }"></div>
        </div>
      </div>

      <div class="stat-card success">
        <div class="stat-icon">✅</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.total_check_ins }}</div>
          <div class="stat-label">已完成打卡</div>
        </div>
      </div>

      <div class="stat-card warning">
        <div class="stat-icon">🔥</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.current_streak }} <small>天</small></div>
          <div class="stat-label">连续打卡</div>
        </div>
      </div>

      <div class="stat-card purple">
        <div class="stat-icon">🏆</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.max_streak }} <small>天</small></div>
          <div class="stat-label">最长连续</div>
        </div>
      </div>

      <div class="stat-card blue">
        <div class="stat-icon">⏱️</div>
        <div class="stat-content">
          <div class="stat-value">{{ formatDuration(stats.total_duration) }}</div>
          <div class="stat-label">累计时长</div>
        </div>
      </div>

      <div class="stat-card danger">
        <div class="stat-icon">⚠️</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.absent_count }}</div>
          <div class="stat-label">缺勤次数</div>
        </div>
      </div>
    </div>

    <div class="detail-stats" v-if="stats">
      <h3 class="sub-title">打卡分布</h3>
      <div class="chart-card">
        <div class="bar-chart">
          <div class="bar-item">
            <div class="bar-label">准时</div>
            <div class="bar-track">
              <div class="bar-fill success" :style="{ width: getBarWidth(stats.on_time_count) }"></div>
            </div>
            <div class="bar-value">{{ stats.on_time_count }}</div>
          </div>
          <div class="bar-item">
            <div class="bar-label">迟到</div>
            <div class="bar-track">
              <div class="bar-fill warning" :style="{ width: getBarWidth(stats.late_count) }"></div>
            </div>
            <div class="bar-value">{{ stats.late_count }}</div>
          </div>
          <div class="bar-item">
            <div class="bar-label">补卡</div>
            <div class="bar-track">
              <div class="bar-fill purple" :style="{ width: getBarWidth(stats.makeup_count) }"></div>
            </div>
            <div class="bar-value">{{ stats.makeup_count }}</div>
          </div>
          <div class="bar-item">
            <div class="bar-label">缺勤</div>
            <div class="bar-track">
              <div class="bar-fill danger" :style="{ width: getBarWidth(stats.absent_count) }"></div>
            </div>
            <div class="bar-value">{{ stats.absent_count }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { formatDuration } from '../utils'
import { useStore } from '../stores'

const store = useStore()

const year = ref(new Date().getFullYear())
const month = ref(new Date().getMonth() + 1)
const stats = ref(null)
const maxBarCount = ref(1)

async function loadStats() {
  stats.value = await store.fetchMonthly({
    year: year.value,
    month: month.value
  })
  if (stats.value) {
    maxBarCount.value = Math.max(
      stats.value.on_time_count || 0,
      stats.value.late_count || 0,
      stats.value.makeup_count || 0,
      stats.value.absent_count || 0,
      1
    )
  }
}

function getBarWidth(count) {
  return (count / maxBarCount.value * 100) + '%'
}

function prevMonth() {
  if (month.value === 1) {
    month.value = 12
    year.value--
  } else {
    month.value--
  }
}

function nextMonth() {
  if (month.value === 12) {
    month.value = 1
    year.value++
  } else {
    month.value++
  }
}

watch([year, month], loadStats)
onMounted(loadStats)
</script>

<style scoped>
.stats-overview {
  background: white;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  border: 1px solid #e2e8f0;
}

.stats-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 700;
}

.month-selector {
  display: flex;
  align-items: center;
  gap: 12px;
}

.month-text {
  font-weight: 600;
  color: #475569;
}

.nav-btn-sm {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #f8fafc;
  font-weight: 600;
  color: #475569;
  transition: all 0.2s;
}

.nav-btn-sm:hover {
  background: #6366f1;
  color: white;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
  margin-bottom: 28px;
}

.stat-card {
  border-radius: 16px;
  padding: 20px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.08);
}

.stat-card.primary {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(99, 102, 241, 0.05));
}
.stat-card.success {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(16, 185, 129, 0.05));
}
.stat-card.warning {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.1), rgba(245, 158, 11, 0.05));
}
.stat-card.purple {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.1), rgba(139, 92, 246, 0.05));
}
.stat-card.blue {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1), rgba(59, 130, 246, 0.05));
}
.stat-card.danger {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.1), rgba(239, 68, 68, 0.05));
}

.stat-icon {
  font-size: 28px;
  margin-bottom: 8px;
}

.stat-content {
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 800;
  color: #1e293b;
  line-height: 1.2;
}

.stat-value small {
  font-size: 14px;
  font-weight: 500;
  color: #64748b;
}

.stat-label {
  font-size: 13px;
  color: #64748b;
  margin-top: 2px;
  font-weight: 500;
}

.stat-bar {
  height: 6px;
  background: rgba(255,255,255,0.5);
  border-radius: 6px;
  overflow: hidden;
  margin-top: 12px;
}

.stat-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #6366f1, #4f46e5);
  border-radius: 6px;
  transition: width 0.5s ease;
}

.detail-stats {
  border-top: 1px solid #e2e8f0;
  padding-top: 24px;
}

.sub-title {
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 16px;
}

.chart-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 20px;
}

.bar-chart {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.bar-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.bar-label {
  width: 50px;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
}

.bar-track {
  flex: 1;
  height: 24px;
  background: white;
  border-radius: 12px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 12px;
  transition: width 0.5s ease;
  min-width: 4px;
}

.bar-fill.success { background: linear-gradient(90deg, #10b981, #059669); }
.bar-fill.warning { background: linear-gradient(90deg, #f59e0b, #d97706); }
.bar-fill.purple { background: linear-gradient(90deg, #8b5cf6, #7c3aed); }
.bar-fill.danger { background: linear-gradient(90deg, #ef4444, #dc2626); }

.bar-value {
  width: 50px;
  text-align: right;
  font-weight: 700;
  color: #1e293b;
}
</style>
