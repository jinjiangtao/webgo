<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { saveSnapshot, exportAggregateExcel } from '../api/data'
import { useDashboardStore } from '../stores/dashboard'
import { useTheme } from '../composables/useTheme'

const store = useDashboardStore()
const { toggleTheme, isDark } = useTheme()

const snapshotLoading = ref(false)
const exportLoading = ref(false)
const autoRefresh = ref(false)
let refreshTimer: number | null = null

const handleSaveSnapshot = async () => {
  snapshotLoading.value = true
  try {
    await ElMessageBox.prompt('请输入快照名称', '保存快照', {
      confirmButtonText: '保存',
      cancelButtonText: '取消',
      inputPlaceholder: '例如：2024年度Q1运营数据',
      inputValidator: (value) => {
        if (!value) return '快照名称不能为空'
        return true
      },
    }).then(async ({ value }) => {
      await saveSnapshot({
        name: value,
        description: `钻取路径: ${store.pathSummary}`,
        drillPath: store.drillPath,
        filters: store.filters,
        records: store.records,
      })
      ElMessage.success('快照保存成功！')
    }).catch(() => {})
  } catch (e) {
    ElMessage.error('快照保存失败')
  } finally {
    snapshotLoading.value = false
  }
}

const handleExport = async () => {
  exportLoading.value = true
  try {
    const blob = await exportAggregateExcel({
      drillPath: store.drillPath,
      metrics: store.metrics,
      filters: store.filters,
      format: 'xlsx',
    })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `数据分析报表_${new Date().toISOString().slice(0, 10)}.xlsx`
    link.click()
    URL.revokeObjectURL(url)
    ElMessage.success('报表导出成功！')
  } catch (e) {
    ElMessage.error('导出失败')
  } finally {
    exportLoading.value = false
  }
}

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    refreshTimer = window.setInterval(() => {
      store.loadData()
    }, 30000)
    ElMessage.info('已开启自动刷新（每30秒）')
  } else {
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
    ElMessage.info('已关闭自动刷新')
  }
}

const handleRefresh = () => {
  store.loadData()
  store.loadDimensions()
  ElMessage.success('数据已刷新')
}
</script>

<template>
  <div class="toolbar">
    <div class="toolbar-left">
      <h1 class="page-title">企业数据分析平台</h1>
      <span class="page-subtitle">精细化运营 · 智能决策</span>
    </div>

    <div class="toolbar-right">
      <div class="toolbar-actions">
        <button
          class="toolbar-btn"
          :class="{ active: autoRefresh }"
          @click="toggleAutoRefresh"
          title="自动刷新"
        >
          <span :class="{ 'animate-spin': autoRefresh }">🔄</span>
          <span class="btn-text">{{ autoRefresh ? '自动' : '自动' }}</span>
        </button>

        <button class="toolbar-btn" @click="handleRefresh" title="手动刷新">
          <span>⟳</span>
          <span class="btn-text">刷新</span>
        </button>

        <button
          class="toolbar-btn"
          @click="handleSaveSnapshot"
          :disabled="snapshotLoading"
          title="保存快照"
        >
          <span>📸</span>
          <span class="btn-text">快照</span>
        </button>

        <button
          class="toolbar-btn primary"
          @click="handleExport"
          :disabled="exportLoading"
          title="导出报表"
        >
          <span>📥</span>
          <span class="btn-text">导出</span>
        </button>

        <button class="toolbar-btn" @click="toggleTheme" title="切换主题">
          <span>{{ isDark ? '☀️' : '🌙' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: linear-gradient(90deg, rgba(15, 23, 42, 0.95), rgba(30, 41, 59, 0.9));
  border-bottom: 1px solid var(--border-color);
  backdrop-filter: blur(10px);
}

.toolbar-left {
  display: flex;
  align-items: baseline;
  gap: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  font-family: 'JetBrains Mono', monospace;
  text-shadow: 0 0 20px rgba(59, 130, 246, 0.5);
}

.page-subtitle {
  font-size: 13px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 13px;
  font-family: 'JetBrains Mono', monospace;
}

.toolbar-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.15);
  border-color: var(--neon-blue);
  color: var(--text-primary);
  box-shadow: 0 0 15px rgba(59, 130, 246, 0.2);
}

.toolbar-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.toolbar-btn.primary {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.3), rgba(6, 182, 212, 0.3));
  border-color: var(--neon-blue);
  color: var(--text-primary);
}

.toolbar-btn.primary:hover:not(:disabled) {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.5), rgba(6, 182, 212, 0.5));
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.4);
}

.toolbar-btn.active {
  background: rgba(16, 185, 129, 0.2);
  border-color: var(--neon-green);
  color: var(--neon-green);
}

.btn-text {
  font-weight: 500;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
