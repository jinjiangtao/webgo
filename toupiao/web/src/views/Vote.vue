<template>
  <div class="vote page-container">
    <div class="vote-header">
      <el-button @click="goBack" :icon="Back">
        返回列表
      </el-button>
      <el-button @click="goRecords" :icon="List">
        查看记录
      </el-button>
    </div>

    <div v-loading="loading" class="vote-content">
      <el-card class="activity-info card-shadow" v-if="activity">
        <div class="activity-header">
          <h2>{{ activity.title }}</h2>
          <div class="activity-tags">
            <el-tag
              :type="activity.status === 1 ? 'success' : 'info'"
              size="large"
            >
              {{ activity.status === 1 ? '进行中' : '已关闭' }}
            </el-tag>
            <el-tag
              :type="activity.vote_type === 'single' ? '' : 'warning'"
              size="large"
            >
              {{ activity.vote_type === 'single' ? '单选' : `多选(最多${activity.max_choices}项)` }}
            </el-tag>
          </div>
        </div>
        <p class="activity-desc" v-if="activity.description">
          {{ activity.description }}
        </p>
        <div class="activity-meta">
          <span>
            <el-icon><Clock /></el-icon>
            开始时间：{{ formatDate(activity.start_time) }}
          </span>
          <span>
            <el-icon><Calendar /></el-icon>
            结束时间：{{ formatDate(activity.end_time) }}
          </span>
          <span>
            <el-icon><DataAnalysis /></el-icon>
            总票数：{{ stats.total_votes || 0 }} 票
          </span>
        </div>
      </el-card>

      <el-row :gutter="20" class="vote-main">
        <el-col :xs="24" :md="12">
          <el-card class="vote-options card-shadow">
            <h3 class="section-title">
              <el-icon><Select /></el-icon>
              请选择您要投票的选项
            </h3>

            <div v-if="activity && activity.vote_type === 'single'" class="options-group">
              <el-radio-group v-model="selectedOptions" class="radio-group">
                <el-radio
                  v-for="option in activity.options"
                  :key="option.id"
                  :value="option.id"
                  class="option-item"
                >
                  <span class="option-name">{{ option.name }}</span>
                </el-radio>
              </el-radio-group>
            </div>

            <div v-else-if="activity && activity.vote_type === 'multiple'" class="options-group">
              <el-checkbox-group v-model="selectedOptions" class="checkbox-group">
                <el-checkbox
                  v-for="option in activity.options"
                  :key="option.id"
                  :value="option.id"
                  class="option-item"
                >
                  <span class="option-name">{{ option.name }}</span>
                </el-checkbox>
              </el-checkbox-group>
              <p class="select-tip">
                （最多选择 {{ activity.max_choices }} 个选项，已选择 {{ selectedOptions.length }} 个）
              </p>
            </div>

            <el-alert
              v-if="!canVote"
              title="当前无法投票"
              :type="voteAlertType"
              :description="voteAlertMessage"
              show-icon
              class="vote-alert"
            >
            </el-alert>

            <div class="vote-actions">
              <el-button
                type="primary"
                class="btn-primary vote-btn"
                size="large"
                @click="handleVote"
                :loading="submitting"
                :disabled="!canSubmit"
              >
                <el-icon><Check /></el-icon>
                提交投票
              </el-button>
              <el-button @click="refreshStats" :loading="isRefreshing">
                <el-icon><Refresh /></el-icon>
                刷新数据
              </el-button>
              <el-switch
                v-model="autoRefresh"
                active-text="自动刷新"
                inactive-text="手动刷新"
              >
              </el-switch>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="24" :md="12">
          <RankingList
            :ranking-data="stats.ranking || []"
            :total-votes="stats.total_votes || 0"
            :auto-refresh="autoRefresh"
            :is-refreshing="isRefreshing"
          />
        </el-col>
      </el-row>

      <div class="chart-section">
        <VoteChart :chart-data="stats.chart_data || []" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Back, List, Check, Refresh, Clock, Calendar, DataAnalysis, Select } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import RankingList from '../components/RankingList.vue'
import VoteChart from '../components/VoteChart.vue'
import { getActivityDetail } from '../api/activity'
import { getVoteStats } from '../api/stats'
import { submitVote } from '../api/vote'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const isRefreshing = ref(false)
const autoRefresh = ref(true)
const activity = ref(null)
const stats = reactive({
  total_votes: 0,
  ranking: [],
  chart_data: []
})
const selectedOptions = ref([])

let refreshTimer = null

const canVote = computed(() => {
  if (!activity.value) return false
  if (activity.value.status !== 1) return false
  const now = new Date()
  const startTime = new Date(activity.value.start_time)
  const endTime = new Date(activity.value.end_time)
  return now >= startTime && now <= endTime
})

const canSubmit = computed(() => {
  if (!canVote.value) return false
  if (submitting.value) return false
  if (selectedOptions.value.length === 0) return false
  if (activity.value && activity.value.vote_type === 'multiple') {
    return selectedOptions.value.length <= activity.value.max_choices
  }
  return true
})

const voteAlertType = computed(() => {
  if (!activity.value) return 'info'
  if (activity.value.status !== 1) return 'info'
  const now = new Date()
  const startTime = new Date(activity.value.start_time)
  const endTime = new Date(activity.value.end_time)
  if (now < startTime) return 'warning'
  if (now > endTime) return 'info'
  return 'info'
})

const voteAlertMessage = computed(() => {
  if (!activity.value) return ''
  if (activity.value.status !== 1) return '该活动已被管理员关闭'
  const now = new Date()
  const startTime = new Date(activity.value.start_time)
  const endTime = new Date(activity.value.end_time)
  if (now < startTime) return `活动尚未开始，开始时间：${formatDate(activity.value.start_time)}`
  if (now > endTime) return `活动已结束，结束时间：${formatDate(activity.value.end_time)}`
  return ''
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

const fetchActivity = async () => {
  loading.value = true
  try {
    const data = await getActivityDetail(route.params.id)
    activity.value = data
  } catch (error) {
    ElMessage.error('获取活动信息失败')
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const data = await getVoteStats(route.params.id)
    stats.total_votes = data.total_votes
    stats.ranking = data.ranking
    stats.chart_data = data.chart_data
  } catch (error) {
    // 静默处理，避免频繁报错
  }
}

const refreshStats = async () => {
  isRefreshing.value = true
  await fetchStats()
  setTimeout(() => {
    isRefreshing.value = false
  }, 500)
}

const handleVote = async () => {
  if (!canSubmit.value) return

  const optionIds = Array.isArray(selectedOptions.value)
    ? selectedOptions.value
    : [selectedOptions.value]

  const selectedNames = optionIds.map(id => {
    const opt = activity.value.options.find(o => o.id === id)
    return opt ? opt.name : ''
  }).filter(Boolean)

  try {
    await ElMessageBox.confirm(
      `您确定要为「${selectedNames.join('、')}」投票吗？`,
      '投票确认',
      {
        confirmButtonText: '确认投票',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    submitting.value = true
    await submitVote(activity.value.id, optionIds)

    ElMessage.success('投票成功！')
    selectedOptions.value = []
    await refreshStats()
  } catch (error) {
    if (error !== 'cancel') {
      // 错误已在request拦截器中处理
    }
  } finally {
    submitting.value = false
  }
}

const goBack = () => {
  router.push('/')
}

const goRecords = () => {
  router.push(`/records/${route.params.id}`)
}

const startAutoRefresh = () => {
  stopAutoRefresh()
  refreshTimer = setInterval(() => {
    if (autoRefresh.value && !isRefreshing.value) {
      refreshStats()
    }
  }, 5000)
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

watch(autoRefresh, (val) => {
  if (val) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
})

onMounted(() => {
  fetchActivity()
  fetchStats()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.vote {
  padding: 0 20px;
}

.vote-header {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.vote-content {
  min-height: 400px;
}

.activity-info {
  margin-bottom: 20px;
}

.activity-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.activity-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.activity-tags {
  display: flex;
  gap: 8px;
}

.activity-desc {
  color: #606266;
  line-height: 1.8;
  margin-bottom: 16px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.activity-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

.activity-meta span {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
  font-size: 14px;
}

.activity-meta .el-icon {
  color: #667eea;
}

.vote-main {
  margin-bottom: 20px;
}

.vote-options {
  height: 100%;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #303133;
}

.section-title .el-icon {
  color: #667eea;
}

.options-group {
  margin-bottom: 20px;
}

.radio-group,
.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.option-item {
  padding: 16px 20px;
  background: #f5f7fa;
  border-radius: 8px;
  transition: all 0.3s;
  font-size: 15px;
}

.option-item:hover {
  background: #e8ecf1;
  transform: translateX(4px);
}

.option-name {
  font-weight: 500;
}

.select-tip {
  margin-top: 12px;
  color: #909399;
  font-size: 13px;
}

.vote-alert {
  margin-bottom: 20px;
}

.vote-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

.vote-btn {
  padding: 12px 40px;
  font-size: 16px;
}

.chart-section {
  margin-top: 20px;
}
</style>
