<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import SideNav from '../components/SideNav.vue'
import Toolbar from '../components/Toolbar.vue'
import { exportAggregateExcel, exportExcel } from '../api/data'
import { useDashboardStore } from '../stores/dashboard'

const store = useDashboardStore()

const exportLoading = ref(false)
const reportTemplates = [
  {
    id: 1,
    name: '销售数据分析报表',
    description: '包含销售额、订单量、用户数等核心指标的多维度分析',
    icon: '📊',
    type: 'aggregate',
    metrics: ['sales', 'orders', 'users'],
  },
  {
    id: 2,
    name: '区域运营分析报表',
    description: '按区域维度分析各地区业务表现和市场占比',
    icon: '🗺️',
    type: 'aggregate',
    metrics: ['sales', 'orders'],
  },
  {
    id: 3,
    name: '业务类型分析报表',
    description: '按业务类型分析不同业务线的盈利情况',
    icon: '💼',
    type: 'aggregate',
    metrics: ['sales', 'amount'],
  },
  {
    id: 4,
    name: '时间趋势分析报表',
    description: '按时间维度分析业务同比环比变化趋势',
    icon: '📈',
    type: 'aggregate',
    metrics: ['sales', 'orders', 'users', 'amount'],
  },
  {
    id: 5,
    name: '原始数据明细报表',
    description: '导出包含溯源信息的完整原始明细数据',
    icon: '📋',
    type: 'raw',
    metrics: ['sales', 'orders', 'users', 'amount'],
  },
  {
    id: 6,
    name: '异常数据专项报表',
    description: '仅导出标记为异常的数据记录，用于问题排查',
    icon: '⚠️',
    type: 'raw',
    metrics: ['sales', 'orders', 'users'],
  },
]

const exportHistory = ref([
  { id: 1, name: '销售数据分析报表', time: '2024-01-15 14:30:22', size: '2.4 MB', status: 'success' },
  { id: 2, name: '区域运营分析报表', time: '2024-01-15 11:20:15', size: '1.8 MB', status: 'success' },
  { id: 3, name: '原始数据明细报表', time: '2024-01-14 16:45:33', size: '15.6 MB', status: 'success' },
])

const handleExport = async (template: any) => {
  exportLoading.value = true
  try {
    let res
    let blob: Blob
    if (template.type === 'aggregate') {
      blob = await exportAggregateExcel({
        drillPath: store.drillPath,
        metrics: template.metrics,
        filters: store.filters,
        format: 'xlsx',
      })
    } else {
      blob = await exportExcel({
        drillPath: store.drillPath,
        metrics: template.metrics,
        filters: store.filters,
        format: 'xlsx',
      })
    }
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${template.name}_${new Date().toISOString().slice(0, 10)}.xlsx`
    link.click()
    URL.revokeObjectURL(url)

    exportHistory.value.unshift({
      id: Date.now(),
      name: template.name,
      time: new Date().toLocaleString('zh-CN'),
      size: template.type === 'raw' ? '~10 MB' : '~2 MB',
      status: 'success',
    })

    ElMessage.success(`${template.name} 导出成功！`)
  } catch (e) {
    ElMessage.error('导出失败，请重试')
  } finally {
    exportLoading.value = false
  }
}

onMounted(() => {
  document.documentElement.classList.add('dark')
})
</script>

<template>
  <div class="page-container dark">
    <SideNav />

    <div class="main-content">
      <Toolbar />

      <div class="content-area bg-grid">
        <div class="content-wrapper">
          <div class="page-header">
            <div>
              <h2 class="page-title">📋 报表中心</h2>
              <p class="page-desc">选择报表模板，一键导出专业数据分析报告</p>
            </div>
          </div>

          <div class="section-title">
            <span class="section-icon">📑</span>
            <span class="section-text">报表模板</span>
          </div>

          <div class="template-grid">
            <div
              v-for="template in reportTemplates"
              :key="template.id"
              class="template-card glass glass-hover"
            >
              <div class="template-header">
                <span class="template-icon">{{ template.icon }}</span>
                <span
                  class="template-badge"
                  :class="template.type === 'aggregate' ? 'agg' : 'raw'"
                >
                  {{ template.type === 'aggregate' ? '聚合报表' : '明细报表' }}
                </span>
              </div>
              <h3 class="template-name">{{ template.name }}</h3>
              <p class="template-desc">{{ template.description }}</p>
              <div class="template-metrics">
                <span
                  v-for="m in template.metrics"
                  :key="m"
                  class="metric-tag"
                >
                  {{ m === 'sales' ? '销售额' : m === 'orders' ? '订单量' : m === 'users' ? '用户数' : '客单价' }}
                </span>
              </div>
              <button
                class="export-btn"
                :disabled="exportLoading"
                @click="handleExport(template)"
              >
                <span>📥</span>
                <span>导出报表</span>
              </button>
            </div>
          </div>

          <div class="section-title" style="margin-top: 32px;">
            <span class="section-icon">📜</span>
            <span class="section-text">导出历史</span>
          </div>

          <div class="history-table glass glass-hover">
            <table class="data-table">
              <thead>
                <tr>
                  <th>报表名称</th>
                  <th>导出时间</th>
                  <th>文件大小</th>
                  <th>状态</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in exportHistory" :key="item.id">
                  <td class="name-cell">
                    <span class="name-icon">📄</span>
                    <span class="name-text">{{ item.name }}</span>
                  </td>
                  <td>{{ item.time }}</td>
                  <td>{{ item.size }}</td>
                  <td>
                    <span class="status-tag success">✓ 成功</span>
                  </td>
                  <td>
                    <button class="action-btn">下载</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  min-height: 100vh;
  background: var(--bg-primary);
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
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  font-family: 'JetBrains Mono', monospace;
  text-shadow: 0 0 20px rgba(59, 130, 246, 0.5);
}

.page-desc {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0;
  font-family: 'JetBrains Mono', monospace;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.section-icon {
  font-size: 20px;
}

.section-text {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  font-family: 'JetBrains Mono', monospace;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.template-card {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.template-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.template-icon {
  font-size: 32px;
}

.template-badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 500;
  font-family: 'JetBrains Mono', monospace;
}

.template-badge.agg {
  background: rgba(59, 130, 246, 0.15);
  border: 1px solid var(--neon-blue);
  color: var(--neon-blue);
}

.template-badge.raw {
  background: rgba(139, 92, 246, 0.15);
  border: 1px solid var(--neon-purple);
  color: #8B5CF6;
}

.template-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  font-family: 'JetBrains Mono', monospace;
}

.template-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0;
  line-height: 1.5;
  min-height: 40px;
}

.template-metrics {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.metric-tag {
  padding: 2px 8px;
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 11px;
  color: var(--text-secondary);
  font-family: 'JetBrains Mono', monospace;
}

.export-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.3), rgba(6, 182, 212, 0.3));
  border: 1px solid var(--neon-blue);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
  margin-top: auto;
}

.export-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.5), rgba(6, 182, 212, 0.5));
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.4);
  transform: translateY(-2px);
}

.export-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.history-table {
  padding: 0;
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  font-family: 'JetBrains Mono', monospace;
}

.data-table th {
  background: rgba(30, 41, 59, 0.9);
  padding: 14px 20px;
  text-align: left;
  color: var(--text-primary);
  font-weight: 600;
  border-bottom: 2px solid var(--border-color);
}

.data-table td {
  padding: 14px 20px;
  border-bottom: 1px solid rgba(59, 130, 246, 0.1);
  color: var(--text-secondary);
}

.data-table tbody tr:hover {
  background: rgba(59, 130, 246, 0.08);
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.name-icon {
  font-size: 16px;
}

.name-text {
  color: var(--text-primary);
  font-weight: 500;
}

.status-tag {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-tag.success {
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid var(--neon-green);
  color: var(--neon-green);
}

.action-btn {
  padding: 4px 12px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--neon-blue);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.action-btn:hover {
  background: rgba(59, 130, 246, 0.15);
  border-color: var(--neon-blue);
}
</style>
