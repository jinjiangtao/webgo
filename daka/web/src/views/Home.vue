<template>
  <div class="home-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">打卡中心</h1>
        <p class="page-subtitle">今天是 {{ todayText }}，加油坚持每一天 💪</p>
      </div>
      <button class="btn-create" @click="showCreateForm = true">
        <span class="plus-icon">+</span> 新建任务
      </button>
    </div>

    <div v-if="store.loading" class="loading-wrap">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="store.tasks.length === 0" class="empty-state">
      <div class="empty-icon">📋</div>
      <h3>还没有打卡任务</h3>
      <p>创建你的第一个自律任务，开启好习惯养成之旅！</p>
      <button class="btn-create-lg" @click="showCreateForm = true">
        <span class="plus-icon">+</span> 创建任务
      </button>
    </div>

    <div v-else>
      <div class="filter-bar">
        <div class="filter-tabs">
          <button
            v-for="tab in filterTabs"
            :key="tab.value"
            class="filter-tab"
            :class="{ active: activeFilter === tab.value }"
            @click="activeFilter = tab.value"
          >
            {{ tab.label }}
            <span class="tab-count">{{ getFilterCount(tab.value) }}</span>
          </button>
        </div>
      </div>

      <div v-if="filteredTasks.length === 0" class="empty-tasks">
        暂无此类任务
      </div>

      <div class="tasks-grid">
        <TaskCard
          v-for="item in filteredTasks"
          :key="item.ID || item.id"
          :task="normalizeTask(item)"
          :today-record="item.today_record"
          @check-in="handleCheckIn"
          @makeup="handleMakeup"
          @edit="handleEditTask"
          @delete="handleDeleteTask"
        />
      </div>

      <div class="mini-calendar">
        <h3 class="section-title">近期打卡概览</h3>
        <CalendarView />
      </div>
    </div>

    <TaskForm
      v-if="showCreateForm"
      :task="editingTask"
      @close="closeForm"
      @submit="handleSubmitTask"
    />

    <MakeupDialog
      v-if="showMakeupDialog"
      :task="makeupTask"
      :default-duration="makeupDuration"
      @close="showMakeupDialog = false"
      @submit="handleMakeupSubmit"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useStore } from '../stores'
import TaskCard from '../components/TaskCard.vue'
import TaskForm from '../components/TaskForm.vue'
import CalendarView from '../components/CalendarView.vue'
import MakeupDialog from '../components/MakeupDialog.vue'

const store = useStore()

const showCreateForm = ref(false)
const editingTask = ref(null)
const showMakeupDialog = ref(false)
const makeupTask = ref(null)
const makeupDuration = ref(0)
const activeFilter = ref('all')

const filterTabs = [
  { label: '全部', value: 'all' },
  { label: '进行中', value: 'active' },
  { label: '已完成今日', value: 'done' },
  { label: '待打卡', value: 'pending' }
]

const todayText = computed(() => {
  const now = new Date()
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return `${now.getFullYear()}年${now.getMonth() + 1}月${now.getDate()}日 ${weekdays[now.getDay()]}`
})

const filteredTasks = computed(() => {
  const tasks = store.tasks
  switch (activeFilter.value) {
    case 'active':
      return tasks.filter(t => t.IsActive || t.is_active)
    case 'done':
      return tasks.filter(t => t.today_record)
    case 'pending':
      return tasks.filter(t => !t.today_record && (t.IsActive || t.is_active))
    default:
      return tasks
  }
})

function getFilterCount(value) {
  const tasks = store.tasks
  switch (value) {
    case 'active':
      return tasks.filter(t => t.IsActive || t.is_active).length
    case 'done':
      return tasks.filter(t => t.today_record).length
    case 'pending':
      return tasks.filter(t => !t.today_record && (t.IsActive || t.is_active)).length
    default:
      return tasks.length
  }
}

function normalizeTask(item) {
  return {
    id: item.ID || item.id,
    name: item.Name || item.name,
    description: item.Description || item.description,
    cycle_type: item.CycleType || item.cycle_type,
    countdown_seconds: item.CountdownSeconds || item.countdown_seconds,
    start_time: item.StartTime || item.start_time,
    end_time: item.EndTime || item.end_time,
    color: item.Color || item.color,
    is_active: item.IsActive !== undefined ? item.IsActive : item.is_active
  }
}

function closeForm() {
  showCreateForm.value = false
  editingTask.value = null
}

function handleEditTask(task) {
  editingTask.value = task
  showCreateForm.value = true
}

async function handleSubmitTask(data) {
  try {
    if (editingTask.value) {
      await store.updateTask(editingTask.value.ID || editingTask.value.id, data)
    } else {
      await store.createTask(data)
    }
    closeForm()
  } catch (e) {
    alert('保存失败: ' + (e.response?.data?.error || e.message))
  }
}

function handleDeleteTask(id) {
  if (confirm('确定要删除这个任务吗？相关的打卡记录也会被删除。')) {
    store.deleteTask(id)
  }
}

function handleCheckIn({ taskId, duration }) {
  store.checkIn({
    task_id: taskId,
    is_makeup: false,
    duration
  }).catch(e => {
    alert('打卡失败: ' + (e.response?.data?.error || e.message))
  })
}

function handleMakeup({ taskId, duration }) {
  makeupTask.value = store.tasks.find(t => (t.ID || t.id) === taskId)
  makeupDuration.value = duration
  showMakeupDialog.value = true
}

async function handleMakeupSubmit({ taskId, date, duration, note }) {
  try {
    await store.checkIn({
      task_id: taskId,
      is_makeup: true,
      duration,
      note,
      record_date: date
    })
    showMakeupDialog.value = false
    await store.fetchTasks()
  } catch (e) {
    alert('补卡失败: ' + (e.response?.data?.error || e.message))
  }
}

onMounted(() => {
  store.fetchTasks()
})
</script>

<style scoped>
.home-page {
  position: relative;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 28px;
  flex-wrap: wrap;
  gap: 16px;
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

.btn-create {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  border-radius: 14px;
  font-weight: 600;
  font-size: 15px;
  transition: all 0.2s;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.3);
}

.btn-create:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
}

.plus-icon {
  font-size: 18px;
  font-weight: 700;
}

.loading-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: #64748b;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  background: white;
  border-radius: 24px;
  padding: 60px 40px;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  border: 1px solid #e2e8f0;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 8px;
}

.empty-state p {
  color: #64748b;
  margin-bottom: 28px;
}

.btn-create-lg {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 14px 32px;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
  border-radius: 14px;
  font-weight: 600;
  font-size: 16px;
  transition: all 0.2s;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.3);
}

.btn-create-lg:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
}

.filter-bar {
  margin-bottom: 20px;
}

.filter-tabs {
  display: flex;
  gap: 8px;
  background: white;
  padding: 6px;
  border-radius: 14px;
  display: inline-flex;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
  border: 1px solid #e2e8f0;
}

.filter-tab {
  padding: 8px 18px;
  border-radius: 10px;
  font-weight: 500;
  font-size: 14px;
  color: #64748b;
  background: transparent;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.filter-tab:hover {
  background: #f8fafc;
  color: #334155;
}

.filter-tab.active {
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  color: white;
}

.tab-count {
  background: rgba(255,255,255,0.2);
  padding: 2px 8px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.filter-tab:not(.active) .tab-count {
  background: #f1f5f9;
  color: #64748b;
}

.empty-tasks {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
  background: white;
  border-radius: 16px;
  border: 1px dashed #e2e8f0;
}

.tasks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.mini-calendar {
  margin-top: 40px;
}

.section-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 16px;
}
</style>
