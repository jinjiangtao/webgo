<template>
  <div class="admin page-container">
    <h2 class="page-title">
      <el-icon><Setting /></el-icon>
      管理后台
    </h2>

    <el-row :gutter="20" class="stats-row" v-loading="loading">
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
            <el-icon><Files /></el-icon>
          </div>
          <div class="stat-content">
            <div class="label">活动总数</div>
            <div class="value">{{ dashboard.activity_count || 0 }}</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            <el-icon><TrendCharts /></el-icon>
          </div>
          <div class="stat-content">
            <div class="label">总投票数</div>
            <div class="value">{{ dashboard.total_votes || 0 }}</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe 0%, #43e97b 100%);">
            <el-icon><Sunny /></el-icon>
          </div>
          <div class="stat-content">
            <div class="label">今日投票</div>
            <div class="value">{{ dashboard.today_votes || 0 }}</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-content">
            <div class="label">活跃活动</div>
            <div class="value">{{ activeCount }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="content-row">
      <el-col :span="24">
        <el-card class="activities-card card-shadow">
          <div class="card-header">
            <h3>
              <el-icon><List /></el-icon>
              活动管理
            </h3>
            <div class="header-actions">
              <el-button type="primary" class="btn-primary" @click="goCreate">
                <el-icon><Plus /></el-icon>
                新建活动
              </el-button>
              <el-button @click="refreshData" :icon="Refresh" :loading="loading">
                刷新
              </el-button>
            </div>
          </div>

          <el-table
            :data="activities"
            v-loading="loading"
            stripe
            style="width: 100%"
          >
            <el-table-column
              prop="id"
              label="ID"
              width="80"
              align="center"
            />
            <el-table-column
              prop="title"
              label="活动标题"
              min-width="200"
              show-overflow-tooltip
            />
            <el-table-column
              label="投票类型"
              width="120"
              align="center"
            >
              <template #default="{ row }">
                <el-tag :type="row.vote_type === 'single' ? '' : 'warning'" size="small">
                  {{ row.vote_type === 'single' ? '单选' : '多选' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              label="选项数"
              width="100"
              align="center"
            >
              <template #default="{ row }">
                {{ row.options?.length || 0 }}
              </template>
            </el-table-column>
            <el-table-column
              label="总票数"
              width="100"
              align="center"
            >
              <template #default="{ row }">
                <span class="vote-count">{{ getTotalVotes(row) }}</span>
              </template>
            </el-table-column>
            <el-table-column
              label="状态"
              width="100"
              align="center"
            >
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? '进行中' : '已关闭' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column
              label="开始时间"
              width="180"
            >
              <template #default="{ row }">
                {{ formatDate(row.start_time) }}
              </template>
            </el-table-column>
            <el-table-column
              label="结束时间"
              width="180"
            >
              <template #default="{ row }">
                {{ formatDate(row.end_time) }}
              </template>
            </el-table-column>
            <el-table-column
              label="操作"
              width="260"
              fixed="right"
              align="center"
            >
              <template #default="{ row }">
                <el-button
                  type="primary"
                  size="small"
                  :icon="View"
                  @click="goVote(row.id)"
                >
                  投票
                </el-button>
                <el-button
                  :type="row.status === 1 ? 'warning' : 'success'"
                  size="small"
                  :icon="row.status === 1 ? VideoPause : VideoPlay"
                  @click="toggleStatus(row)"
                  :loading="row.toggling"
                >
                  {{ row.status === 1 ? '关闭' : '开启' }}
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  :icon="Delete"
                  @click="handleDelete(row)"
                  :loading="row.deleting"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-empty
            v-if="!loading && activities.length === 0"
            description="暂无活动数据"
            class="empty-state"
          />
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="recent-row" v-if="dashboard.recent_activities && dashboard.recent_activities.length > 0">
      <el-col :span="24">
        <el-card class="recent-card card-shadow">
          <div class="card-header">
            <h3>
              <el-icon><Clock /></el-icon>
              最近活动
            </h3>
          </div>
          <div class="recent-list">
            <div
              v-for="item in dashboard.recent_activities"
              :key="item.id"
              class="recent-item"
            >
              <div class="recent-info">
                <h4>{{ item.title }}</h4>
                <div class="recent-meta">
                  <el-tag :type="item.status === 1 ? 'success' : 'info'" size="small">
                    {{ item.status === 1 ? '进行中' : '已关闭' }}
                  </el-tag>
                  <span>{{ item.vote_type === 'single' ? '单选' : '多选' }}</span>
                  <span>总票数：{{ item.total_votes }}</span>
                  <span>创建于：{{ formatDate(item.created_at) }}</span>
                </div>
              </div>
              <el-button type="primary" size="small" @click="goVote(item.id)">
                查看详情
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Setting, Files, TrendCharts, Sunny, User, List, Plus, Refresh,
  View, VideoPause, VideoPlay, Delete, Clock
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getActivityList, toggleActivityStatus, deleteActivity, getDashboardStats } from '../api/activity'

const router = useRouter()

const loading = ref(false)
const activities = ref([])
const dashboard = reactive({
  activity_count: 0,
  total_votes: 0,
  today_votes: 0,
  recent_activities: []
})

const activeCount = computed(() => {
  return activities.value.filter(a => a.status === 1).length
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getTotalVotes = (activity) => {
  if (!activity.options) return 0
  return activity.options.reduce((sum, opt) => sum + (opt.vote_count || 0), 0)
}

const fetchDashboard = async () => {
  try {
    const data = await getDashboardStats()
    dashboard.activity_count = data.activity_count
    dashboard.total_votes = data.total_votes
    dashboard.today_votes = data.today_votes
    dashboard.recent_activities = data.recent_activities
  } catch (error) {
    // 静默处理
  }
}

const fetchActivities = async () => {
  loading.value = true
  try {
    const data = await getActivityList()
    activities.value = data
  } catch (error) {
    ElMessage.error('获取活动列表失败')
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  fetchDashboard()
  fetchActivities()
}

const toggleStatus = async (row) => {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '开启' : '关闭'

  try {
    await ElMessageBox.confirm(
      `确定要${action}活动「${row.title}」吗？`,
      `${action}确认`,
      {
        confirmButtonText: `确认${action}`,
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    row.toggling = true
    await toggleActivityStatus(row.id, newStatus)
    row.status = newStatus
    ElMessage.success(`活动${action}成功`)
    fetchDashboard()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`${action}失败`)
    }
  } finally {
    row.toggling = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除活动「${row.title}」吗？此操作将同时删除所有相关的投票记录，且不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'error',
        confirmButtonClass: 'el-button--danger'
      }
    )

    row.deleting = true
    await deleteActivity(row.id)
    ElMessage.success('活动删除成功')
    refreshData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  } finally {
    row.deleting = false
  }
}

const goCreate = () => {
  router.push('/create')
}

const goVote = (id) => {
  router.push(`/vote/${id}`)
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.admin {
  padding: 0 20px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.page-title .el-icon {
  color: #667eea;
}

.stats-row {
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
  margin-bottom: 20px;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 28px;
  flex-shrink: 0;
}

.stat-content .label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 4px;
}

.stat-content .value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.content-row {
  margin-bottom: 24px;
}

.activities-card {
  padding: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #303133;
}

.card-header h3 .el-icon {
  color: #667eea;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.vote-count {
  font-weight: 600;
  color: #667eea;
}

.empty-state {
  padding: 60px 0;
}

.recent-card {
  padding: 24px;
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recent-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #f5f7fa;
  border-radius: 8px;
  transition: all 0.3s;
}

.recent-item:hover {
  background: #e8ecf1;
}

.recent-info h4 {
  margin: 0 0 8px 0;
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.recent-meta {
  display: flex;
  gap: 16px;
  align-items: center;
  font-size: 13px;
  color: #909399;
  flex-wrap: wrap;
}

.recent-meta .el-tag {
  margin-right: 0;
}
</style>
