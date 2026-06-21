<template>
  <div class="calendar-card">
    <div class="calendar-header">
      <button class="nav-btn" @click="prevMonth">←</button>
      <div class="month-title">{{ year }}年 {{ month }}月</div>
      <button class="nav-btn" @click="nextMonth">→</button>
    </div>

    <div class="weekdays">
      <div v-for="d in weekdays" :key="d" class="weekday">{{ d }}</div>
    </div>

    <div class="calendar-grid">
      <div
        v-for="(cell, idx) in calendarCells"
        :key="idx"
        class="day-cell"
        :class="{
          'other-month': !cell.isCurrentMonth,
          'today': cell.isToday,
          'has-records': cell.dayRecord && cell.dayRecord.records && cell.dayRecord.records.length > 0,
          'all-done': cell.dayRecord && (cell.dayRecord.completed_count || 0) >= (cell.dayRecord.task_count || 0) && (cell.dayRecord.task_count || 0) > 0
        }"
        @click="handleDayClick(cell)"
      >
        <div class="day-num">{{ cell.day }}</div>
        <div v-if="cell.dayRecord && cell.dayRecord.records && cell.dayRecord.records.length > 0" class="day-indicators">
          <div class="indicator-row">
            <span
              v-for="r in cell.dayRecord.records.slice(0, 3)"
              :key="r.id || r.ID"
              class="record-dot"
              :style="{ background: getRecordColor(r) }"
              :title="getRecordTaskName(r) + ' - ' + getStatusText(r.status || r.Status)"
            ></span>
            <span v-if="cell.dayRecord.records.length > 3" class="more-dot">+{{ cell.dayRecord.records.length - 3 }}</span>
          </div>
          <div class="progress-mini">
            <div
              class="progress-mini-fill"
              :style="{
                width: (cell.dayRecord.task_count || 0) > 0
                  ? ((cell.dayRecord.completed_count || 0) / cell.dayRecord.task_count * 100) + '%'
                  : '0%'
              }"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <div class="calendar-legend">
      <span class="legend-item"><span class="legend-dot success"></span> 完成</span>
      <span class="legend-item"><span class="legend-dot warning"></span> 迟到</span>
      <span class="legend-item"><span class="legend-dot purple"></span> 补卡</span>
      <span class="legend-item"><span class="legend-dot danger"></span> 缺勤</span>
    </div>

    <div v-if="selectedDay && selectedDay.dayRecord" class="day-detail">
      <div class="detail-header">
        <strong>{{ selectedDay.date }} 打卡详情</strong>
      </div>
      <div v-if="!selectedDay.dayRecord.records || selectedDay.dayRecord.records.length === 0" class="empty-detail">
        当日暂无打卡记录
      </div>
      <div v-else class="record-list">
        <div v-for="r in selectedDay.dayRecord.records" :key="r.id || r.ID" class="record-item">
          <div class="record-color" :style="{ background: getRecordTaskColor(r) }"></div>
          <div class="record-info">
            <div class="record-task">{{ getRecordTaskName(r) }}</div>
            <div class="record-meta">
              <span class="status-tag" :style="{ background: getStatusBg(r.status || r.Status), color: getRecordColor(r) }">
                {{ getStatusText(r.status || r.Status) }}
              </span>
              <span v-if="r.check_in_time || r.CheckInTime" class="record-time">
                {{ formatTime(r.check_in_time || r.CheckInTime) }}
              </span>
              <span v-if="r.duration || r.Duration" class="record-duration">
                {{ formatDuration(r.duration || r.Duration) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { getStatusText, getStatusColor, formatDuration } from '../utils'
import { useStore } from '../stores'

const store = useStore()

const weekdays = ['日', '一', '二', '三', '四', '五', '六']

const year = ref(new Date().getFullYear())
const month = ref(new Date().getMonth() + 1)
const calendarData = ref(null)
const selectedDay = ref(null)

const calendarCells = computed(() => {
  if (!calendarData.value || !Array.isArray(calendarData.value.calendar)) {
    return buildEmptyCalendar()
  }

  const firstDay = new Date(year.value, month.value - 1, 1)
  const startWeekday = firstDay.getDay()
  const daysInMonth = new Date(year.value, month.value, 0).getDate()
  const prevMonthDays = new Date(year.value, month.value - 1, 0).getDate()

  const cells = []
  const recordMap = {}
  calendarData.value.calendar.forEach(d => {
    if (d && d.date) {
      recordMap[d.date] = d
    }
  })

  const todayStr = new Date().toISOString().slice(0, 10)

  for (let i = startWeekday - 1; i >= 0; i--) {
    const day = prevMonthDays - i
    const prevMonth = month.value === 1 ? 12 : month.value - 1
    const prevYear = month.value === 1 ? year.value - 1 : year.value
    const dateStr = `${prevYear}-${String(prevMonth).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    cells.push({
      day,
      date: dateStr,
      isCurrentMonth: false,
      isToday: dateStr === todayStr,
      dayRecord: recordMap[dateStr] || null
    })
  }

  for (let day = 1; day <= daysInMonth; day++) {
    const dateStr = `${year.value}-${String(month.value).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    cells.push({
      day,
      date: dateStr,
      isCurrentMonth: true,
      isToday: dateStr === todayStr,
      dayRecord: recordMap[dateStr] || null
    })
  }

  let nextDayCounter = 1
  while (cells.length % 7 !== 0) {
    const day = nextDayCounter++
    const nextMonth = month.value === 12 ? 1 : month.value + 1
    const nextYear = month.value === 12 ? year.value + 1 : year.value
    const dateStr = `${nextYear}-${String(nextMonth).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    cells.push({
      day,
      date: dateStr,
      isCurrentMonth: false,
      isToday: dateStr === todayStr,
      dayRecord: recordMap[dateStr] || null
    })
  }

  return cells
})

function buildEmptyCalendar() {
  const firstDay = new Date(year.value, month.value - 1, 1)
  const startWeekday = firstDay.getDay()
  const daysInMonth = new Date(year.value, month.value, 0).getDate()
  const prevMonthDays = new Date(year.value, month.value - 1, 0).getDate()

  const cells = []
  const todayStr = new Date().toISOString().slice(0, 10)

  for (let i = startWeekday - 1; i >= 0; i--) {
    const day = prevMonthDays - i
    const prevMonth = month.value === 1 ? 12 : month.value - 1
    const prevYear = month.value === 1 ? year.value - 1 : year.value
    const dateStr = `${prevYear}-${String(prevMonth).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    cells.push({ day, date: dateStr, isCurrentMonth: false, isToday: dateStr === todayStr, dayRecord: null })
  }

  for (let day = 1; day <= daysInMonth; day++) {
    const dateStr = `${year.value}-${String(month.value).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    cells.push({ day, date: dateStr, isCurrentMonth: true, isToday: dateStr === todayStr, dayRecord: null })
  }

  let nextDayCounter = 1
  while (cells.length % 7 !== 0) {
    const day = nextDayCounter++
    const nextMonth = month.value === 12 ? 1 : month.value + 1
    const nextYear = month.value === 12 ? year.value + 1 : year.value
    const dateStr = `${nextYear}-${String(nextMonth).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    cells.push({ day, date: dateStr, isCurrentMonth: false, isToday: dateStr === todayStr, dayRecord: null })
  }

  return cells
}

async function loadCalendar() {
  try {
    const data = await store.fetchCalendar({
      year: year.value,
      month: month.value
    })
    calendarData.value = data || { calendar: [] }
  } catch (e) {
    console.error('Load calendar failed:', e)
    calendarData.value = { calendar: [] }
  }
}

function prevMonth() {
  if (month.value === 1) {
    month.value = 12
    year.value--
  } else {
    month.value--
  }
  selectedDay.value = null
}

function nextMonth() {
  if (month.value === 12) {
    month.value = 1
    year.value++
  } else {
    month.value++
  }
  selectedDay.value = null
}

function getRecordColor(record) {
  const status = record.status || record.Status
  if (status === 'absent') return '#ef4444'
  return getStatusColor(status)
}

function getRecordTaskName(record) {
  if (record.task && record.task.name) return record.task.name
  if (record.Task && record.Task.Name) return record.Task.Name
  return '打卡任务'
}

function getRecordTaskColor(record) {
  if (record.task && record.task.color) return record.task.color
  if (record.Task && record.Task.Color) return record.Task.Color
  return '#6366f1'
}

function getStatusBg(status) {
  const colors = {
    'on_time': 'rgba(16, 185, 129, 0.12)',
    'checked_in': 'rgba(16, 185, 129, 0.12)',
    'late': 'rgba(245, 158, 11, 0.12)',
    'makeup': 'rgba(139, 92, 246, 0.12)',
    'absent': 'rgba(239, 68, 68, 0.12)'
  }
  return colors[status] || 'rgba(100, 116, 139, 0.12)'
}

function formatTime(timeStr) {
  if (!timeStr) return ''
  const d = new Date(timeStr)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

function handleDayClick(cell) {
  selectedDay.value = cell
}

watch([year, month], loadCalendar)
onMounted(loadCalendar)
</script>

<style scoped>
.calendar-card {
  background: white;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  border: 1px solid #e2e8f0;
}

.calendar-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  margin-bottom: 20px;
}

.month-title {
  font-size: 18px;
  font-weight: 700;
  min-width: 140px;
  text-align: center;
}

.nav-btn {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: #f8fafc;
  font-size: 20px;
  font-weight: 600;
  color: #475569;
  transition: all 0.2s;
}

.nav-btn:hover {
  background: #6366f1;
  color: white;
  transform: scale(1.05);
}

.weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 6px;
  margin-bottom: 8px;
}

.weekday {
  text-align: center;
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  padding: 8px 0;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 6px;
}

.day-cell {
  aspect-ratio: 1;
  border-radius: 12px;
  padding: 8px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  border: 2px solid transparent;
  position: relative;
  background: #f8fafc;
  min-height: 80px;
}

.day-cell:hover {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.05);
}

.day-cell.other-month {
  opacity: 0.4;
}

.day-cell.today {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.08);
}

.day-cell.today .day-num {
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 4px;
}

.day-cell.all-done {
  background: rgba(16, 185, 129, 0.08);
}

.day-num {
  font-size: 14px;
  font-weight: 600;
  color: #334155;
  margin-bottom: 4px;
}

.day-indicators {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  gap: 4px;
}

.indicator-row {
  display: flex;
  gap: 3px;
  flex-wrap: wrap;
}

.record-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.more-dot {
  font-size: 10px;
  color: #64748b;
  font-weight: 600;
}

.progress-mini {
  height: 4px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-mini-fill {
  height: 100%;
  background: linear-gradient(90deg, #10b981, #059669);
  border-radius: 4px;
  transition: width 0.3s;
}

.calendar-legend {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #64748b;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.legend-dot.success { background: #10b981; }
.legend-dot.warning { background: #f59e0b; }
.legend-dot.purple { background: #8b5cf6; }
.legend-dot.danger { background: #ef4444; }

.day-detail {
  margin-top: 20px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 12px;
}

.detail-header {
  margin-bottom: 12px;
  color: #334155;
}

.empty-detail {
  color: #94a3b8;
  font-size: 14px;
  text-align: center;
  padding: 12px;
}

.record-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.record-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  background: white;
  border-radius: 10px;
}

.record-color {
  width: 4px;
  height: 32px;
  border-radius: 4px;
  flex-shrink: 0;
}

.record-info {
  flex: 1;
  min-width: 0;
}

.record-task {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 2px;
}

.record-meta {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.status-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 20px;
  font-weight: 600;
}

.record-time,
.record-duration {
  font-size: 12px;
  color: #64748b;
}
</style>
