<template>
  <div class="task-card" :style="{ '--task-color': task.color || '#6366f1' }">
    <div class="card-header">
      <div class="card-title-wrap">
        <div class="color-dot"></div>
        <div class="title-info">
          <h3 class="task-name">{{ task.name }}</h3>
          <div class="task-meta">
            <span class="cycle-tag">{{ cycleText }}</span>
            <span v-if="task.start_time" class="time-range">{{ task.start_time }} - {{ task.end_time }}</span>
          </div>
        </div>
      </div>
      <div class="card-actions">
        <button class="icon-btn" @click="$emit('edit', task)" title="编辑">✏️</button>
        <button class="icon-btn danger" @click="$emit('delete', task.id)" title="删除">🗑️</button>
      </div>
    </div>

    <div v-if="task.description" class="task-desc">{{ task.description }}</div>

    <div class="countdown-section">
      <div class="countdown-display">
        <div class="time-block">
          <span class="time-num">{{ timeDisplay.hours }}</span>
          <span class="time-label">时</span>
        </div>
        <span class="time-sep">:</span>
        <div class="time-block">
          <span class="time-num">{{ timeDisplay.minutes }}</span>
          <span class="time-label">分</span>
        </div>
        <span class="time-sep">:</span>
        <div class="time-block">
          <span class="time-num">{{ timeDisplay.seconds }}</span>
          <span class="time-label">秒</span>
        </div>
      </div>

      <div class="progress-wrap">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
        </div>
        <span class="progress-text">{{ progressPercent.toFixed(1) }}%</span>
      </div>

      <div class="duration-info">
        <span>目标时长: {{ formatDuration(task.countdown_seconds) }}</span>
        <span>已进行: {{ formatDuration(elapsedSeconds) }}</span>
      </div>
    </div>

    <div v-if="todayRecord" class="status-section completed">
      <div class="status-icon">✓</div>
      <div class="status-info">
        <span class="status-label" :style="{ color: statusColor }">{{ statusText }}</span>
        <span class="status-time">{{ formatCheckTime }}</span>
      </div>
    </div>

    <div class="card-footer">
      <div class="actions">
        <button v-if="!isRunning" class="btn btn-primary" @click="startCountdown">
          {{ elapsedSeconds > 0 ? '继续' : '开始打卡' }}
        </button>
        <button v-else class="btn btn-success" @click="handleCheckIn(false)">
          完成打卡
        </button>
        <button class="btn btn-outline" @click="handleMakeup">补卡</button>
      </div>
      <div class="footer-info">
        <span class="tag" :class="{ active: task.is_active }">
          {{ task.is_active ? '● 进行中' : '○ 已暂停' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { formatDuration, getStatusText, getStatusColor, getCycleText } from '../utils'

const props = defineProps({
  task: { type: Object, required: true },
  todayRecord: { type: Object, default: null }
})

const emit = defineEmits(['check-in', 'makeup', 'edit', 'delete'])

const isRunning = ref(false)
const elapsedSeconds = ref(0)
let timer = null

const isCompleted = computed(() => {
  return elapsedSeconds.value >= props.task.countdown_seconds
})

const timeDisplay = computed(() => {
  const remaining = Math.max(0, props.task.countdown_seconds - elapsedSeconds.value)
  return {
    hours: String(Math.floor(remaining / 3600)).padStart(2, '0'),
    minutes: String(Math.floor((remaining % 3600) / 60)).padStart(2, '0'),
    seconds: String(remaining % 60).padStart(2, '0')
  }
})

const progressPercent = computed(() => {
  if (props.task.countdown_seconds <= 0) return 0
  return Math.min(100, (elapsedSeconds.value / props.task.countdown_seconds) * 100)
})

const cycleText = computed(() => getCycleText(props.task.cycle_type))
const statusText = computed(() => props.todayRecord ? getStatusText(props.todayRecord.status) : '')
const statusColor = computed(() => props.todayRecord ? getStatusColor(props.todayRecord.status) : '')
const formatCheckTime = computed(() => {
  if (!props.todayRecord?.check_in_time) return ''
  const d = new Date(props.todayRecord.check_in_time)
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
})

function startCountdown() {
  isRunning.value = true
  timer = setInterval(() => {
    elapsedSeconds.value++
  }, 1000)
}

function stopCountdown() {
  isRunning.value = false
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

function handleCheckIn(isMakeup) {
  stopCountdown()
  emit('check-in', {
    taskId: props.task.id,
    duration: elapsedSeconds.value,
    isMakeup
  })
}

function handleMakeup() {
  stopCountdown()
  emit('makeup', {
    taskId: props.task.id,
    duration: elapsedSeconds.value
  })
}

watch(() => props.todayRecord, (record) => {
  if (record?.duration) {
    elapsedSeconds.value = record.duration
  }
}, { immediate: true })

onMounted(() => {
  if (props.todayRecord?.duration) {
    elapsedSeconds.value = props.todayRecord.duration
  }
})

onUnmounted(() => {
  stopCountdown()
})
</script>

<style scoped>
.task-card {
  --task-color: #6366f1;
  background: white;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  border: 1px solid #e2e8f0;
  transition: all 0.3s;
  position: relative;
  overflow: hidden;
}

.task-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: var(--task-color);
  opacity: 0.8;
}

.task-card:hover {
  box-shadow: 0 10px 40px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.card-title-wrap {
  display: flex;
  gap: 12px;
  align-items: center;
}

.color-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--task-color);
  flex-shrink: 0;
}

.title-info {
  flex: 1;
}

.task-name {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 4px;
}

.task-meta {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.cycle-tag {
  background: #f1f5f9;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 12px;
  color: #64748b;
}

.time-range {
  font-size: 12px;
  color: #64748b;
}

.card-actions {
  display: flex;
  gap: 4px;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #f8fafc;
  font-size: 16px;
  transition: all 0.2s;
}

.icon-btn:hover {
  background: #f1f5f9;
  transform: scale(1.05);
}

.icon-btn.danger:hover {
  background: #fef2f2;
}

.task-desc {
  color: #64748b;
  font-size: 14px;
  margin-bottom: 20px;
  padding: 12px 16px;
  background: #f8fafc;
  border-radius: 12px;
  line-height: 1.6;
}

.countdown-section {
  background: rgba(99, 102, 241, 0.08);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 20px;
}

.countdown-display {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.time-block {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.time-num {
  font-size: 42px;
  font-weight: 800;
  background: linear-gradient(135deg, var(--task-color), #4f46e5);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  line-height: 1;
  min-width: 60px;
  text-align: center;
}

.time-label {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.time-sep {
  font-size: 36px;
  font-weight: 700;
  color: var(--task-color);
  margin-top: -12px;
}

.progress-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.progress-bar {
  flex: 1;
  height: 10px;
  background: #e2e8f0;
  border-radius: 10px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--task-color), #6366f1);
  border-radius: 10px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--task-color);
  min-width: 60px;
  text-align: right;
}

.duration-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #64748b;
}

.status-section {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 12px;
  margin-bottom: 16px;
}

.status-section.completed {
  background: rgba(16, 185, 129, 0.1);
}

.status-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #10b981;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 20px;
}

.status-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.status-label {
  font-weight: 600;
  font-size: 15px;
}

.status-time {
  font-size: 12px;
  color: #64748b;
}

.card-footer {
  border-top: 1px solid #e2e8f0;
  padding-top: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.btn {
  padding: 10px 20px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}

.btn-success {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
}

.btn-success:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

.btn-outline {
  background: transparent;
  border: 2px solid #e2e8f0;
  color: #1e293b;
}

.btn-outline:hover {
  border-color: #6366f1;
  color: #6366f1;
  background: rgba(99, 102, 241, 0.05);
}

.footer-info {
  display: flex;
  gap: 8px;
}

.tag {
  font-size: 12px;
  padding: 4px 12px;
  border-radius: 20px;
  background: #f1f5f9;
  color: #64748b;
}

.tag.active {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}
</style>
