<template>
  <div class="stats-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">数据统计</h1>
        <p class="page-subtitle">查看你的自律数据，见证每一步成长 📈</p>
      </div>
    </div>

    <div class="stats-layout">
      <div class="stats-main">
        <StatsPanel />
      </div>
      <div class="stats-side">
        <CalendarView />
      </div>
    </div>

    <div class="task-stats-section">
      <h2 class="section-title">📝 各任务统计</h2>
      <div v-if="!store.tasks.length" class="empty-tip">暂无任务数据</div>
      <div v-else class="task-stats-grid">
        <div
          v-for="task in store.tasks"
          :key="task.ID || task.id"
          class="task-stat-card"
        >
          <div class="task-stat-header">
            <div class="color-bar" :style="{ background: task.Color || task.color || '#6366f1' }"></div>
            <div class="task-stat-info">
              <h4>{{ task.Name || task.name }}</h4>
              <span>{{ getCycleText(task.CycleType || task.cycle_type) }}</span>
            </div>
          </div>
          <div class="task-stat-body" v-if="taskStats[task.ID || task.id]">
            <div class="rate-circle">
              <svg viewBox="0 0 36 36" class="circle-progress">
                <circle cx="18" cy="18" r="15.9155" fill="none" stroke="#e2e8f0" stroke-width="3"/>
                <circle
                  cx="18" cy="18" r="15.9155" fill="none"
                  :stroke="task.Color || task.color || '#6366f1'"
                  stroke-width="3"
                  stroke-linecap="round"
                  :stroke-dasharray="`${taskStats[task.ID || task.id].completion_rate || 0}, 100`"
                  transform="rotate(-90 18 18)"
                />
              </svg>
              <span class="rate-text">{{ (taskStats[task.ID || task.id].completion_rate || 0).toFixed(0) }}%</span>
            </div>
            <div class="stat-items">
              <div class="stat-item">
                <span class="label">完成天数</span>
                <span class="value">{{ (taskStats[task.ID || task.id].on_time || 0) + (taskStats[task.ID || task.id].late || 0) + (taskStats[task.ID || task.id].makeup || 0) }}</span>
              </div>
              <div class="stat-item">
                <span class="label">迟到</span>
                <span class="value warning">{{ taskStats[task.ID || task.id].late || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="label">补卡</span>
                <span class="value purple">{{ taskStats[task.ID || task.id].makeup || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="label">累计时长</span>
                <span class="value primary">{{ formatDuration(taskStats[task.ID || task.id].total_duration || 0) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useStore } from '../stores'
import StatsPanel from '../components/StatsPanel.vue'
import CalendarView from '../components/CalendarView.vue'
import { formatDuration, getCycleText } from '../utils'

const store = useStore()
const taskStats = ref({})

async function loadAllTaskStats() {
  const year = new Date().getFullYear()
  const month = new Date().getMonth() + 1
  for (const task of store.tasks) {
    const id = task.ID || task.id
    try {
      const data = await store.fetchTaskStats(id, { year, month })
      taskStats.value[id] = data
    } catch (e) {
      console.error('Failed to load task stats:', e)
    }
  }
}

watch(() => store.tasks, () => {
  if (store.tasks.length) {
    loadAllTaskStats()
  }
}, { immediate: true })

onMounted(async () => {
  await store.fetchTasks()
})
</script>

<style scoped>
.stats-page {
  position: relative;
}

.page-header {
  margin-bottom: 28px;
}

.page-title {
  font-size: 28px;
  font-weight: 800;
  margin-bottom: 6px;
  background: linear-gradient(135deg, #1e293b, #475569);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.page-subtitle {
  color: #64748b;
  font-size: 15px;
}

.stats-layout {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 24px;
  margin-bottom: 40px;
}

@media (max-width: 960px) {
  .stats-layout {
    grid-template-columns: 1fr;
  }
}

.task-stats-section {
  margin-top: 8px;
}

.section-title {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 20px;
}

.empty-tip {
  background: white;
  padding: 40px;
  border-radius: 16px;
  text-align: center;
  color: #94a3b8;
  border: 1px dashed #e2e8f0;
}

.task-stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.task-stat-card {
  background: white;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  border: 1px solid #e2e8f0;
  transition: all 0.3s;
}

.task-stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 40px rgba(0,0,0,0.1);
}

.task-stat-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f1f5f9;
}

.color-bar {
  width: 6px;
  height: 40px;
  border-radius: 6px;
}

.task-stat-info h4 {
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 2px;
}

.task-stat-info span {
  font-size: 12px;
  color: #64748b;
  background: #f1f5f9;
  padding: 2px 10px;
  border-radius: 20px;
}

.task-stat-body {
  display: flex;
  gap: 24px;
  align-items: center;
}

.rate-circle {
  position: relative;
  width: 96px;
  height: 96px;
  flex-shrink: 0;
}

.circle-progress {
  width: 100%;
  height: 100%;
}

.rate-text {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 800;
  color: #1e293b;
}

.stat-items {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-item .label {
  font-size: 11px;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-item .value {
  font-size: 16px;
  font-weight: 700;
  color: #1e293b;
}

.stat-item .value.warning { color: #f59e0b; }
.stat-item .value.purple { color: #8b5cf6; }
.stat-item .value.primary { color: #6366f1; font-size: 13px; }
</style>
