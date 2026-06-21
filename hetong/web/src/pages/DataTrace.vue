<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import SideNav from '../components/SideNav.vue'
import Toolbar from '../components/Toolbar.vue'
import { getTraceData, exportExcel } from '../api/data'
import type { RawRecord } from '../types'

const loading = ref(false)
const traceLoading = ref(false)
const rawData = ref<RawRecord[]>([])
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  time: '',
  region: '',
  business: '',
  orderNo: '',
  page: 1,
  pageSize: 20,
})

const tracePath = ref<string[]>([])

const timeOptions = [
  { value: '', label: '全部时间' },
  { value: '2025', label: '2025年' },
  { value: '2025-Q1', label: '2025年Q1' },
  { value: '2025-Q2', label: '2025年Q2' },
  { value: '2025-01', label: '2025年1月' },
  { value: '2025-02', label: '2025年2月' },
]

const regionOptions = [
  { value: '', label: '全部区域' },
  { value: '华东', label: '华东区' },
  { value: '华南', label: '华南区' },
  { value: '华北', label: '华北区' },
  { value: '华中', label: '华中区' },
  { value: '西南', label: '西南区' },
]

const businessOptions = [
  { value: '', label: '全部业务' },
  { value: '电商零售', label: '电商零售' },
  { value: '企业服务', label: '企业服务' },
  { value: '线下门店', label: '线下门店' },
  { value: '跨境贸易', label: '跨境贸易' },
]

const handleSearch = async () => {
  loading.value = true
  try {
    const res = await getTraceData({
      ...searchForm,
      page: currentPage.value,
      pageSize: pageSize.value,
    })
    if (res.code === 200) {
      rawData.value = res.data.records || []
      totalCount.value = res.data.total || 0
      buildTracePath()
    }
  } catch (e) {
    ElMessage.error('查询失败')
  } finally {
    loading.value = false
  }
}

const buildTracePath = () => {
  const path: string[] = ['全局视图']
  if (searchForm.time) path.push(`时间: ${searchForm.time}`)
  if (searchForm.region) path.push(`区域: ${searchForm.region}`)
  if (searchForm.business) path.push(`业务: ${searchForm.business}`)
  if (searchForm.orderNo) path.push(`订单: ${searchForm.orderNo}`)
  tracePath.value = path
}

const handleTrace = async (record: RawRecord) => {
  traceLoading.value = true
  try {
    await ElMessageBox.alert(
      `
      <div style="font-family: 'JetBrains Mono', monospace; line-height: 2;">
        <h4 style="color: var(--neon-blue); margin-bottom: 16px;">📊 数据溯源信息</h4>
        <div><strong>订单编号:</strong> ${record.orderNo}</div>
        <div><strong>时间维度:</strong> ${record.dimensions.time}</div>
        <div><strong>区域维度:</strong> ${record.dimensions.region}</div>
        <div><strong>业务维度:</strong> ${record.dimensions.business}</div>
        <hr style="border-color: var(--border-color); margin: 12px 0;">
        <div><strong>销售额:</strong> ¥${record.metrics.sales.toLocaleString()}</div>
        <div><strong>订单量:</strong> ${record.metrics.orders}</div>
        <div><strong>用户数:</strong> ${record.metrics.users}</div>
        <div><strong>客单价:</strong> ¥${record.metrics.amount.toFixed(2)}</div>
        <hr style="border-color: var(--border-color); margin: 12px 0;">
        <div style="color: var(--text-muted); font-size: 12px;">
          <div><strong>数据来源:</strong> 交易系统</div>
          <div><strong>入库时间:</strong> ${record.createdAt}</div>
          <div><strong>数据校验:</strong> ✓ 已通过</div>
          <div><strong>溯源链路:</strong> ${tracePath.value.join(' → ')}</div>
        </div>
      </div>
      `,
      '数据溯源详情',
      {
        dangerouslyUseHTMLString: true,
        confirmButtonText: '关闭',
        customClass: 'trace-dialog',
      }
    )
  } catch (e) {}
  finally {
    traceLoading.value = false
  }
}

const handleReset = () => {
  searchForm.time = ''
  searchForm.region = ''
  searchForm.business = ''
  searchForm.orderNo = ''
  searchForm.page = 1
  searchForm.pageSize = 20
  currentPage.value = 1
  handleSearch()
}

const handleExport = async () => {
  try {
    const blob = await exportExcel({
      drillPath: [],
      metrics: ['sales', 'orders', 'users', 'amount'],
      filters: {
        time: searchForm.time,
        region: searchForm.region,
        business: searchForm.business,
      },
      format: 'xlsx',
    })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `数据溯源报表_${new Date().toISOString().slice(0, 10)}.xlsx`
    link.click()
    URL.revokeObjectURL(url)
    ElMessage.success('溯源报表导出成功！')
  } catch (e) {
    ElMessage.error('导出失败')
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  handleSearch()
}

onMounted(() => {
  document.documentElement.classList.add('dark')
  handleSearch()
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
              <h2 class="page-title">🔍 数据溯源</h2>
              <p class="page-desc">从宏观汇总到单条明细，完整追踪数据链路</p>
            </div>
            <button class="export-btn" @click="handleExport">
              <span>📥</span>
              <span>导出溯源报表</span>
            </button>
          </div>

          <div class="trace-path glass glass-hover" v-if="tracePath.length > 1">
            <span class="path-label">📍 溯源路径:</span>
            <div class="path-items">
              <template v-for="(item, index) in tracePath" :key="index">
                <span class="path-item">{{ item }}</span>
                <span class="path-arrow" v-if="index < tracePath.length - 1">→</span>
              </template>
            </div>
          </div>

          <div class="filter-section glass glass-hover">
            <div class="filter-grid">
              <div class="filter-item">
                <label class="filter-label">时间维度</label>
                <select v-model="searchForm.time" class="filter-select">
                  <option v-for="opt in timeOptions" :key="opt.value" :value="opt.value">
                    {{ opt.label }}
                  </option>
                </select>
              </div>

              <div class="filter-item">
                <label class="filter-label">区域维度</label>
                <select v-model="searchForm.region" class="filter-select">
                  <option v-for="opt in regionOptions" :key="opt.value" :value="opt.value">
                    {{ opt.label }}
                  </option>
                </select>
              </div>

              <div class="filter-item">
                <label class="filter-label">业务维度</label>
                <select v-model="searchForm.business" class="filter-select">
                  <option v-for="opt in businessOptions" :key="opt.value" :value="opt.value">
                    {{ opt.label }}
                  </option>
                </select>
              </div>

              <div class="filter-item">
                <label class="filter-label">订单编号</label>
                <input
                  v-model="searchForm.orderNo"
                  type="text"
                  class="filter-input"
                  placeholder="输入订单编号精确查询"
                />
              </div>
            </div>

            <div class="filter-actions">
              <button class="btn-reset" @click="handleReset">重置</button>
              <button class="btn-search" @click="handleSearch" :disabled="loading">
                <span v-if="loading" class="loading-spin">🔄</span>
                <span v-else>🔍</span>
                查询
              </button>
            </div>
          </div>

          <div class="stats-row">
            <div class="stat-card glass">
              <span class="stat-icon">📄</span>
              <div class="stat-info">
                <span class="stat-value">{{ totalCount.toLocaleString() }}</span>
                <span class="stat-label">总记录数</span>
              </div>
            </div>
            <div class="stat-card glass">
              <span class="stat-icon">📊</span>
              <div class="stat-info">
                <span class="stat-value">{{ rawData.length }}</span>
                <span class="stat-label">当前页</span>
              </div>
            </div>
            <div class="stat-card glass">
              <span class="stat-icon">🏷️</span>
              <div class="stat-info">
                <span class="stat-value">{{ new Set(rawData.map(r => r.dimensions.business)).size }}</span>
                <span class="stat-label">业务类型</span>
              </div>
            </div>
            <div class="stat-card glass">
              <span class="stat-icon">📍</span>
              <div class="stat-info">
                <span class="stat-value">{{ new Set(rawData.map(r => r.dimensions.region)).size }}</span>
                <span class="stat-label">覆盖区域</span>
              </div>
            </div>
          </div>

          <div class="data-table-section glass glass-hover" v-loading="loading">
            <div class="table-header">
              <h3 class="table-title">📋 原始明细数据</h3>
              <span class="table-count">共 {{ totalCount.toLocaleString() }} 条记录</span>
            </div>

            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>订单编号</th>
                    <th>时间</th>
                    <th>区域</th>
                    <th>业务类型</th>
                    <th class="text-right">销售额</th>
                    <th class="text-right">订单量</th>
                    <th class="text-right">用户数</th>
                    <th class="text-right">客单价</th>
                    <th class="text-center">入库时间</th>
                    <th class="text-center">操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="record in rawData" :key="record.id" class="data-row">
                    <td class="order-no">{{ record.orderNo }}</td>
                    <td>{{ record.dimensions.time }}</td>
                    <td>{{ record.dimensions.region }}</td>
                    <td>
                      <span class="biz-tag">{{ record.dimensions.business }}</span>
                    </td>
                    <td class="text-right sales">¥{{ record.metrics.sales.toLocaleString() }}</td>
                    <td class="text-right">{{ record.metrics.orders }}</td>
                    <td class="text-right">{{ record.metrics.users }}</td>
                    <td class="text-right">¥{{ record.metrics.amount.toFixed(2) }}</td>
                    <td class="text-center">{{ record.createdAt }}</td>
                    <td class="text-center">
                      <button
                        class="trace-btn"
                        :disabled="traceLoading"
                        @click="handleTrace(record)"
                      >
                        🔍 溯源
                      </button>
                    </td>
                  </tr>
                  <tr v-if="rawData.length === 0 && !loading">
                    <td colspan="10" class="empty-row">
                      <div class="empty-content">
                        <span class="empty-icon">📭</span>
                        <span class="empty-text">暂无数据，请调整筛选条件</span>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="pagination" v-if="totalCount > 0">
              <button
                class="page-btn"
                :disabled="currentPage === 1"
                @click="handlePageChange(currentPage - 1)"
              >
                ← 上一页
              </button>
              <span class="page-info">
                第 {{ currentPage }} 页 / 共 {{ Math.ceil(totalCount / pageSize) }} 页
              </span>
              <button
                class="page-btn"
                :disabled="currentPage >= Math.ceil(totalCount / pageSize)"
                @click="handlePageChange(currentPage + 1)"
              >
                下一页 →
              </button>
            </div>
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
  max-width: 1600px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
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

.export-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.3), rgba(6, 182, 212, 0.3));
  border: 1px solid var(--neon-blue);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.export-btn:hover {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.5), rgba(6, 182, 212, 0.5));
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.4);
}

.trace-path {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  margin-bottom: 20px;
}

.path-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  font-family: 'JetBrains Mono', monospace;
  white-space: nowrap;
}

.path-items {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.path-item {
  padding: 4px 12px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 4px;
  font-size: 12px;
  color: var(--neon-blue);
  font-family: 'JetBrains Mono', monospace;
}

.path-arrow {
  color: var(--text-muted);
  font-weight: bold;
}

.filter-section {
  padding: 20px;
  margin-bottom: 20px;
}

.filter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-label {
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}

.filter-select,
.filter-input {
  padding: 10px 14px;
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 13px;
  font-family: 'JetBrains Mono', monospace;
  outline: none;
  transition: all 0.2s ease;
}

.filter-select:hover,
.filter-select:focus,
.filter-input:hover,
.filter-input:focus {
  border-color: var(--neon-blue);
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.2);
}

.filter-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-reset {
  padding: 10px 20px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.btn-reset:hover {
  border-color: var(--neon-red);
  color: var(--neon-red);
}

.btn-search {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 24px;
  background: linear-gradient(135deg, var(--neon-blue), var(--neon-cyan));
  border: none;
  border-radius: 6px;
  color: white;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.btn-search:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(59, 130, 246, 0.4);
}

.btn-search:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spin {
  display: inline-block;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
}

.stat-icon {
  font-size: 32px;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  font-family: 'JetBrains Mono', monospace;
  text-shadow: 0 0 10px rgba(59, 130, 246, 0.3);
}

.stat-label {
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
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

.table-count {
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
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
  padding: 14px 16px;
  text-align: left;
  color: var(--text-primary);
  font-weight: 600;
  border-bottom: 2px solid var(--border-color);
  white-space: nowrap;
}

.data-table td {
  padding: 12px 16px;
  border-bottom: 1px solid rgba(59, 130, 246, 0.1);
  color: var(--text-secondary);
  white-space: nowrap;
}

.data-table tbody tr:hover {
  background: rgba(59, 130, 246, 0.08);
}

.text-right {
  text-align: right;
}

.text-center {
  text-align: center;
}

.order-no {
  color: var(--neon-blue);
  font-weight: 500;
  cursor: pointer;
}

.order-no:hover {
  text-decoration: underline;
}

.sales {
  color: var(--neon-green);
  font-weight: 600;
}

.biz-tag {
  display: inline-block;
  padding: 2px 8px;
  background: rgba(139, 92, 246, 0.15);
  border: 1px solid rgba(139, 92, 246, 0.4);
  border-radius: 4px;
  color: #A78BFA;
  font-size: 11px;
}

.trace-btn {
  padding: 4px 12px;
  background: rgba(59, 130, 246, 0.15);
  border: 1px solid var(--neon-blue);
  border-radius: 4px;
  color: var(--neon-blue);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.trace-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.3);
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.3);
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

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 24px;
  padding-top: 20px;
  margin-top: 16px;
  border-top: 1px solid var(--border-color);
}

.page-btn {
  padding: 8px 16px;
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.page-btn:hover:not(:disabled) {
  border-color: var(--neon-blue);
  color: var(--neon-blue);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  font-size: 13px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}
</style>
