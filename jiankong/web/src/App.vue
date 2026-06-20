<script setup>
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue'
import { api } from './api'
import StatBar from './components/StatBar.vue'
import FilterBar from './components/FilterBar.vue'
import DeviceCard from './components/DeviceCard.vue'
import AlarmToast from './components/AlarmToast.vue'
import DeviceFormModal from './components/DeviceFormModal.vue'
import DeviceTypeModal from './components/DeviceTypeModal.vue'
import LogViewerModal from './components/LogViewerModal.vue'
import AlarmListModal from './components/AlarmListModal.vue'

const stats = ref({})
const devices = ref([])
const deviceTypes = ref([])
const filters = reactive({ device_type: '', status: '', keyword: '' })
const autoRefresh = ref(true)
const alarmToasts = ref([])
const knownAlarmIds = new Set()
const loading = ref(false)

const showDeviceForm = ref(false)
const showDeviceType = ref(false)
const showLogs = ref(false)
const showAlarms = ref(false)
const editingDevice = ref(null)

let refreshTimer = null
let keywordTimer = null

async function fetchStats() {
  try { stats.value = await api.getStats() } catch {}
}
async function fetchDevices() {
  loading.value = true
  try {
    devices.value = await api.getDevices({
      device_type: filters.device_type,
      status: filters.status,
      keyword: filters.keyword
    })
  } catch {} finally { loading.value = false }
}
async function fetchTypes() {
  try { deviceTypes.value = await api.getDeviceTypes() } catch {}
}

async function checkAlarms() {
  try {
    const active = await api.getAlarms('active')
    const newOnes = active.filter(a => !knownAlarmIds.has(a.id))
    newOnes.forEach(a => {
      knownAlarmIds.add(a.id)
      alarmToasts.value.push(a)
      setTimeout(() => dismissToast(a), 12000)
    })
  } catch {}
}

function dismissToast(toast) {
  const i = alarmToasts.value.indexOf(toast)
  if (i > -1) alarmToasts.value.splice(i, 1)
}

async function refreshAll() {
  await Promise.all([fetchStats(), fetchDevices()])
  if (autoRefresh.value) checkAlarms()
}

function startPolling() {
  stopPolling()
  refreshTimer = setInterval(refreshAll, 4000)
}
function stopPolling() {
  if (refreshTimer) { clearInterval(refreshTimer); refreshTimer = null }
}

function toggleRefresh() {
  autoRefresh.value = !autoRefresh.value
}

watch(autoRefresh, (v) => {
  if (v) { refreshAll(); startPolling() } else { stopPolling() }
})

watch(() => filters.device_type, refreshAll)
watch(() => filters.status, refreshAll)
watch(() => filters.keyword, () => {
  clearTimeout(keywordTimer)
  keywordTimer = setTimeout(fetchDevices, 300)
})

function handleAdd() { editingDevice.value = null; showDeviceForm.value = true }
function handleEdit(d) { editingDevice.value = d; showDeviceForm.value = true }

async function handleSubmitDevice(data) {
  try {
    if (editingDevice.value) {
      await api.updateDevice(editingDevice.value.id, data)
    } else {
      await api.createDevice(data)
    }
    showDeviceForm.value = false
    await refreshAll()
  } catch (e) {
    alert(e.message)
  }
}

async function handleDelete(d) {
  if (!confirm(`确定删除设备「${d.name}」？相关告警与日志将一并清除。`)) return
  try {
    await api.deleteDevice(d.id)
    await refreshAll()
  } catch (e) { alert(e.message) }
}

async function handleAckToast(toast) {
  try {
    await api.acknowledgeAlarm(toast.id)
    dismissToast(toast)
    await fetchStats()
  } catch {}
}

onMounted(async () => {
  await Promise.all([fetchStats(), fetchDevices(), fetchTypes()])
  const active = await api.getAlarms('active')
  active.forEach(a => knownAlarmIds.add(a.id))
  startPolling()
})
onUnmounted(stopPolling)
</script>

<template>
  <div class="app">
    <StatBar
      :stats="stats"
      :auto-refresh="autoRefresh"
      @toggle-refresh="toggleRefresh"
      @add-device="handleAdd"
      @manage-types="showDeviceType = true"
      @view-logs="showLogs = true"
      @view-alarms="showAlarms = true"
    />

    <FilterBar
      :device-types="deviceTypes"
      :filters="filters"
      :result-count="devices.length"
      @update:filters="(v) => { Object.assign(filters, v) }"
    />

    <main class="main">
      <div v-if="loading && devices.length === 0" class="loading">
        <div class="spinner"></div>
        <span>加载设备数据中...</span>
      </div>
      <div v-else-if="devices.length === 0" class="empty-state">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <rect x="2" y="3" width="20" height="14" rx="2"/>
          <line x1="8" y1="21" x2="16" y2="21"/>
          <line x1="12" y1="17" x2="12" y2="21"/>
        </svg>
        <p>未找到匹配的设备</p>
        <span>请调整筛选条件或录入新设备</span>
      </div>
      <div v-else class="grid">
        <DeviceCard
          v-for="d in devices"
          :key="d.id"
          :device="d"
          @edit="handleEdit"
          @delete="handleDelete"
          @view="handleEdit"
        />
      </div>
    </main>

    <AlarmToast
      :toasts="alarmToasts"
      @acknowledge="handleAckToast"
      @close="dismissToast"
    />

    <DeviceFormModal
      :show="showDeviceForm"
      :device="editingDevice"
      :device-types="deviceTypes"
      @close="showDeviceForm = false"
      @submit="handleSubmitDevice"
    />
    <DeviceTypeModal
      :show="showDeviceType"
      @close="showDeviceType = false"
      @changed="fetchTypes"
    />
    <LogViewerModal
      :show="showLogs"
      :devices="devices"
      @close="showLogs = false"
    />
    <AlarmListModal
      :show="showAlarms"
      @close="showAlarms = false"
    />
  </div>
</template>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.main {
  flex: 1;
  padding: 20px 24px 40px;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  animation: fade-in 0.3s;
}
.loading, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 80px 20px;
  color: var(--text-muted);
}
.spinner {
  width: 32px;
  height: 32px;
  border: 2px solid var(--border-bright);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
.empty-state svg { color: var(--border-bright); }
.empty-state p { font-size: 15px; color: var(--text-secondary); }
.empty-state span { font-size: 12px; }

@media (max-width: 768px) {
  .main { padding: 16px; }
  .grid { grid-template-columns: 1fr; }
}
</style>
