<template>
  <div class="vote-records page-container">
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><List /></el-icon>
        投票记录
      </h2>
      <el-button @click="goBack" :icon="Back">
        返回列表
      </el-button>
    </div>

    <el-card class="activity-card card-shadow" v-if="activity">
      <div class="activity-brief">
        <div>
          <h3>{{ activity.title }}</h3>
          <p class="activity-desc" v-if="activity.description">
            {{ activity.description }}
          </p>
        </div>
        <div class="activity-stats">
          <div class="stat-item">
            <span class="stat-label">总票数</span>
            <span class="stat-value">{{ total }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">选项数</span>
            <span class="stat-value">{{ activity.options?.length || 0 }}</span>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="table-card card-shadow" v-loading="loading">
      <div class="table-toolbar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索选项名称或IP..."
          clearable
          style="width: 300px"
          :prefix-icon="Search"
          @keyup.enter="fetchRecords"
        >
        </el-input>
        <el-button @click="fetchRecords" :icon="Refresh">
          刷新
        </el-button>
      </div>

      <el-table :data="filteredRecords" stripe style="width: 100%">
        <el-table-column
          prop="id"
          label="ID"
          width="80"
          align="center"
        />
        <el-table-column
          prop="option_name"
          label="选项名称"
          min-width="150"
        >
          <template #default="{ row }">
            <el-tag size="small">{{ row.option_name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="ip"
          label="IP地址"
          width="140"
        >
          <template #default="{ row }">
            <span class="ip-text">{{ maskIP(row.ip) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="cookie"
          label="Cookie"
          min-width="200"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <span class="cookie-text">{{ row.cookie || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="created_at"
          label="投票时间"
          width="180"
        >
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>

      <el-empty
        v-if="!loading && records.length === 0"
        description="暂无投票记录"
        class="empty-state"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Back, List, Search, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getActivityDetail } from '../api/activity'
import { getVoteRecords } from '../api/stats'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const activity = ref(null)
const records = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchKeyword = ref('')

const filteredRecords = computed(() => {
  if (!searchKeyword.value) return records.value
  const keyword = searchKeyword.value.toLowerCase()
  return records.value.filter(record =>
    record.option_name.toLowerCase().includes(keyword) ||
    record.ip.toLowerCase().includes(keyword)
  )
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const maskIP = (ip) => {
  if (!ip) return '-'
  const parts = ip.split('.')
  if (parts.length === 4) {
    return `${parts[0]}.${parts[1]}.*.*`
  }
  return ip
}

const fetchActivity = async () => {
  try {
    const data = await getActivityDetail(route.params.id)
    activity.value = data
  } catch (error) {
    ElMessage.error('获取活动信息失败')
  }
}

const fetchRecords = async () => {
  loading.value = true
  try {
    const data = await getVoteRecords(route.params.id, page.value, pageSize.value)
    records.value = data.records
    total.value = data.total
  } catch (error) {
    ElMessage.error('获取投票记录失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (val) => {
  page.value = val
  fetchRecords()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  fetchRecords()
}

const goBack = () => {
  router.push('/')
}

watch(() => route.params.id, () => {
  page.value = 1
  fetchActivity()
  fetchRecords()
})

onMounted(() => {
  fetchActivity()
  fetchRecords()
})
</script>

<style scoped>
.vote-records {
  padding: 0 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
}

.page-title .el-icon {
  color: #667eea;
}

.activity-card {
  margin-bottom: 20px;
}

.activity-brief {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  flex-wrap: wrap;
}

.activity-brief h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.activity-desc {
  color: #606266;
  margin: 0;
  max-width: 600px;
}

.activity-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  text-align: center;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: #fff;
}

.stat-label {
  display: block;
  font-size: 13px;
  opacity: 0.9;
  margin-bottom: 4px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 600;
}

.table-card {
  padding: 20px;
}

.table-toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.ip-text {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  color: #606266;
}

.cookie-text {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 12px;
  color: #909399;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.empty-state {
  padding: 60px 0;
}
</style>
