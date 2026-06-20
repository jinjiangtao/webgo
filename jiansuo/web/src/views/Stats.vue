<template>
  <div class="stats-page">
    <div class="stats-overview">
      <div class="stat-card primary">
        <div class="stat-icon">
          <el-icon :size="32"><Search /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ statistics.total_searches || 0 }}</div>
          <div class="stat-label">累计搜索次数</div>
        </div>
      </div>
      <div class="stat-card success">
        <div class="stat-icon">
          <el-icon :size="32"><DataLine /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ statistics.today_searches || 0 }}</div>
          <div class="stat-label">今日搜索次数</div>
        </div>
      </div>
      <div class="stat-card warning">
        <div class="stat-icon">
          <el-icon :size="32"><Key /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ statistics.unique_keywords || 0 }}</div>
          <div class="stat-label">搜索词数</div>
        </div>
      </div>
      <div class="stat-card info">
        <div class="stat-icon">
          <el-icon :size="32"><Document /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ statistics.total_keywords || 0 }}</div>
          <div class="stat-label">关键词库总量</div>
        </div>
      </div>
      <div class="stat-card danger">
        <div class="stat-icon">
          <el-icon :size="32"><Menu /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ statistics.total_categories || 0 }}</div>
          <div class="stat-label">分类数量</div>
        </div>
      </div>
    </div>

    <div class="charts-row">
      <div class="chart-card hot-chart">
        <div class="card-header">
          <h3>
            <el-icon><TrendCharts /></el-icon>
            热门搜索词 Top 10
          </h3>
        </div>
        <div v-if="hotKeywords.length" class="hot-list">
          <div
            v-for="(item, idx) in hotKeywords"
            :key="item.keyword"
            class="hot-item"
          >
            <div class="rank" :class="rankClass(idx)">{{ idx + 1 }}</div>
            <div class="hot-content">
              <div class="hot-keyword" @click="goSearch(item.keyword)">
                {{ item.keyword }}
              </div>
              <div class="hot-bar">
                <div
                  class="hot-bar-fill"
                  :style="{ width: getBarWidth(item.search_count) + '%' }"
                  :class="rankClass(idx)"
                ></div>
              </div>
            </div>
            <div class="hot-count">{{ item.search_count }} 次</div>
          </div>
        </div>
        <el-empty v-else description="暂无搜索数据" :image-size="100" />
      </div>

      <div class="chart-card trend-chart">
        <div class="card-header">
          <h3>
            <el-icon><Calendar /></el-icon>
            近30天搜索趋势
          </h3>
        </div>
        <div v-if="dailyTrend.length" class="trend-chart-body">
          <svg class="trend-svg" viewBox="0 0 700 260" preserveAspectRatio="xMidYMid meet">
            <g class="grid-lines">
              <line v-for="i in 5" :key="i"
                :x1="40" :x2="680"
                :y1="20 + (i - 1) * 45" :y2="20 + (i - 1) * 45"
                stroke="#f0f2f5" stroke-width="1"
              />
            </g>
            <g class="y-labels">
              <text v-for="(lbl, i) in yLabels" :key="i"
                x="35" :y="25 + i * 45"
                text-anchor="end" fill="#909399" font-size="11"
              >{{ lbl }}</text>
            </g>
            <polyline
              class="trend-line"
              fill="none"
              stroke="#409EFF"
              stroke-width="2.5"
              :points="linePoints"
            />
            <polygon
              class="trend-area"
              fill="url(#gradient)"
              :points="areaPoints"
            />
            <defs>
              <linearGradient id="gradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="#409EFF" stop-opacity="0.3" />
                <stop offset="100%" stop-color="#409EFF" stop-opacity="0" />
              </linearGradient>
            </defs>
            <g class="data-points">
              <circle
                v-for="(p, i) in pointCoords"
                :key="i"
                :cx="p.x"
                :cy="p.y"
                r="4"
                fill="#fff"
                stroke="#409EFF"
                stroke-width="2"
              >
                <title>{{ dailyTrend[i].date }}: {{ dailyTrend[i].count }} 次搜索</title>
              </circle>
            </g>
            <g class="x-labels">
              <text
                v-for="(lbl, i) in xLabels"
                :key="i"
                :x="60 + i * ((680 - 60) / Math.max(xLabels.length - 1, 1))"
                y="252"
                text-anchor="middle"
                fill="#909399"
                font-size="10"
              >{{ lbl }}</text>
            </g>
          </svg>
        </div>
        <el-empty v-else description="暂无搜索数据" :image-size="100" />
      </div>
    </div>

    <div class="chart-card logs-card">
      <div class="card-header">
        <h3>
          <el-icon><List /></el-icon>
          搜索日志
        </h3>
        <div class="header-filters">
          <el-input
            v-model="logKeyword"
            placeholder="搜索关键词"
            size="small"
            style="width: 180px"
            clearable
            @keyup.enter="loadLogs"
            @clear="loadLogs"
          />
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            size="small"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="margin-left: 12px; width: 260px"
          />
          <el-button type="primary" size="small" @click="loadLogs" style="margin-left: 12px">
            查询
          </el-button>
        </div>
      </div>
      <el-table :data="logList" v-loading="logLoading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column label="搜索关键词" min-width="220">
          <template #default="{ row }">
            <el-tag
              type="primary"
              effect="light"
              size="small"
              style="max-width: 280px; overflow: hidden; text-overflow: ellipsis"
              @click="goSearch(row.keyword)"
              class="clickable-tag"
            >
              {{ row.keyword || '(空)' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="result_count" label="结果数" width="100" align="center">
          <template #default="{ row }">
            <span :class="row.result_count === 0 ? 'no-result' : ''">
              {{ row.result_count }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="user_ip" label="IP" width="140" align="center" />
        <el-table-column label="Session" width="180" align="center">
          <template #default="{ row }">
            <span class="session-text">{{ (row.session_id || '').slice(0, 20) }}...</span>
          </template>
        </el-table-column>
        <el-table-column label="搜索时间" width="180" align="center">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          v-model:current-page="logPage"
          v-model:page-size="logPageSize"
          :total="logTotal"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadLogs"
          @current-change="loadLogs"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  getStatistics,
  listSearchLogs
} from '@/api'

const router = useRouter()

const statistics = reactive({
  total_searches: 0,
  today_searches: 0,
  unique_keywords: 0,
  total_keywords: 0,
  total_categories: 0,
  hot_keywords: [],
  daily_trend: []
})

const hotKeywords = ref([])
const dailyTrend = ref([])

const logKeyword = ref('')
const dateRange = ref([])
const logPage = ref(1)
const logPageSize = ref(20)
const logTotal = ref(0)
const logList = ref([])
const logLoading = ref(false)

const maxCount = computed(() => {
  if (!dailyTrend.value.length) return 1
  return Math.max(...dailyTrend.value.map(d => d.count), 1)
})

const yLabels = computed(() => {
  const max = maxCount.value
  const step = Math.ceil(max / 4) || 1
  return [max, Math.ceil(max * 0.75), Math.ceil(max * 0.5), Math.ceil(max * 0.25), 0]
})

const chartData = computed(() => {
  const trend = [...dailyTrend.value].reverse()
  const len = trend.length
  const max = Math.max(...trend.map(d => d.count), 1)
  const xStart = 60
  const xEnd = 680
  const yTop = 20
  const yBottom = 220
  const points = trend.map((d, i) => {
    const x = len <= 1 ? (xStart + xEnd) / 2 : xStart + (xEnd - xStart) * (i / (len - 1))
    const y = yBottom - (yBottom - yTop) * (d.count / max)
    return { x, y }
  })
  return { trend, points }
})

const linePoints = computed(() => {
  return chartData.value.points.map(p => `${p.x},${p.y}`).join(' ')
})

const areaPoints = computed(() => {
  const pts = chartData.value.points
  if (!pts.length) return ''
  const firstX = pts[0].x
  const lastX = pts[pts.length - 1].x
  const topLine = pts.map(p => `${p.x},${p.y}`).join(' ')
  return `${firstX},220 ${topLine} ${lastX},220`
})

const pointCoords = computed(() => chartData.value.points)

const xLabels = computed(() => {
  const trend = chartData.value.trend
  if (!trend.length) return []
  const len = trend.length
  const labelCount = Math.min(len, 7)
  const indices = []
  for (let i = 0; i < labelCount; i++) {
    indices.push(Math.round(i * (len - 1) / Math.max(labelCount - 1, 1)))
  }
  return indices.map(i => {
    const d = trend[i].date
    return d.slice(5)
  })
})

const rankClass = (idx) => {
  if (idx === 0) return 'rank-1'
  if (idx === 1) return 'rank-2'
  if (idx === 2) return 'rank-3'
  return 'rank-n'
}

const getBarWidth = (count) => {
  if (!hotKeywords.value.length) return 0
  const max = Math.max(...hotKeywords.value.map(k => k.search_count), 1)
  return (count / max) * 100
}

const formatDate = (d) => {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN')
}

const goSearch = (keyword) => {
  router.push({ path: '/search', query: { q: keyword } })
}

const loadStats = async () => {
  try {
    const res = await getStatistics()
    const d = res.data || {}
    Object.assign(statistics, d)
    hotKeywords.value = d.hot_keywords || []
    dailyTrend.value = d.daily_trend || []
  } catch (e) {}
}

const loadLogs = async () => {
  logLoading.value = true
  try {
    const params = {
      page: logPage.value,
      page_size: logPageSize.value
    }
    if (logKeyword.value) params.keyword = logKeyword.value
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const res = await listSearchLogs(params)
    logList.value = res.data.list || []
    logTotal.value = res.data.total || 0
  } catch (e) {
  } finally {
    logLoading.value = false
  }
}

onMounted(() => {
  loadStats()
  loadLogs()
})
</script>

<style scoped>
.stats-page {
  padding-bottom: 40px;
}
.stats-overview {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}
.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  transition: all 0.25s;
  position: relative;
  overflow: hidden;
}
.stat-card::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
}
.stat-card.primary::after { background: #409EFF; }
.stat-card.success::after { background: #67C23A; }
.stat-card.warning::after { background: #E6A23C; }
.stat-card.info::after { background: #909399; }
.stat-card.danger::after { background: #F56C6C; }
.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}
.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-card.primary .stat-icon { background: #ecf5ff; color: #409EFF; }
.stat-card.success .stat-icon { background: #f0f9eb; color: #67C23A; }
.stat-card.warning .stat-icon { background: #fdf6ec; color: #E6A23C; }
.stat-card.info .stat-icon { background: #f4f4f5; color: #909399; }
.stat-card.danger .stat-icon { background: #fef0f0; color: #F56C6C; }
.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #303133;
  line-height: 1.2;
}
.stat-label {
  color: #909399;
  font-size: 13px;
  margin-top: 4px;
}
.charts-row {
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: 20px;
  margin-bottom: 24px;
}
.chart-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}
.card-header h3 {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}
.header-filters {
  display: flex;
  align-items: center;
}
.hot-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.hot-item {
  display: flex;
  align-items: center;
  gap: 12px;
}
.rank {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
  background: #f4f4f5;
  color: #909399;
  flex-shrink: 0;
}
.rank.rank-1 { background: linear-gradient(135deg, #F56C6C, #ff7e7e); color: #fff; }
.rank.rank-2 { background: linear-gradient(135deg, #E6A23C, #ffc96e); color: #fff; }
.rank.rank-3 { background: linear-gradient(135deg, #67C23A, #8ed86a); color: #fff; }
.hot-content {
  flex: 1;
  min-width: 0;
}
.hot-keyword {
  color: #303133;
  font-weight: 500;
  margin-bottom: 6px;
  cursor: pointer;
  transition: color 0.2s;
}
.hot-keyword:hover {
  color: #409EFF;
}
.hot-bar {
  height: 6px;
  background: #f0f2f5;
  border-radius: 3px;
  overflow: hidden;
}
.hot-bar-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.4s ease;
  background: linear-gradient(90deg, #c0c4cc, #909399);
}
.hot-bar-fill.rank-1 { background: linear-gradient(90deg, #F56C6C, #ff7e7e); }
.hot-bar-fill.rank-2 { background: linear-gradient(90deg, #E6A23C, #ffc96e); }
.hot-bar-fill.rank-3 { background: linear-gradient(90deg, #67C23A, #8ed86a); }
.hot-count {
  color: #909399;
  font-size: 13px;
  flex-shrink: 0;
}
.trend-chart-body {
  width: 100%;
}
.trend-svg {
  width: 100%;
  height: auto;
}
.logs-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}
.no-result {
  color: #F56C6C;
  font-weight: 500;
}
.session-text {
  font-family: monospace;
  color: #606266;
  font-size: 12px;
}
.clickable-tag {
  cursor: pointer;
}
.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
@media (max-width: 1200px) {
  .stats-overview {
    grid-template-columns: repeat(3, 1fr);
  }
}
@media (max-width: 900px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
  .stats-overview {
    grid-template-columns: repeat(2, 1fr);
  }
}
@media (max-width: 600px) {
  .stats-overview {
    grid-template-columns: 1fr;
  }
}
</style>
