<template>
  <el-card class="activity-card card-shadow" :body-style="{ padding: '0' }">
    <div class="card-header">
      <div class="card-title">
        <h3>{{ activity.title }}</h3>
        <el-tag
          :type="activity.status === 1 ? 'success' : 'info'"
          size="small"
        >
          {{ activity.status === 1 ? '进行中' : '已关闭' }}
        </el-tag>
      </div>
      <div class="card-type">
        <el-tag :type="activity.vote_type === 'single' ? '' : 'warning'" size="small">
          {{ activity.vote_type === 'single' ? '单选' : '多选(最多' + activity.max_choices + '项)' }}
        </el-tag>
      </div>
    </div>
    <div class="card-body">
      <p class="description" v-if="activity.description">
        {{ activity.description }}
      </p>
      <div class="meta">
        <div class="meta-item">
          <el-icon><Clock /></el-icon>
          <span>开始时间：{{ formatDate(activity.start_time) }}</span>
        </div>
        <div class="meta-item">
          <el-icon><Calendar /></el-icon>
          <span>结束时间：{{ formatDate(activity.end_time) }}</span>
        </div>
        <div class="meta-item">
          <el-icon><Tickets /></el-icon>
          <span>选项数量：{{ activity.options?.length || 0 }} 个</span>
        </div>
        <div class="meta-item">
          <el-icon><TrendCharts /></el-icon>
          <span>总票数：{{ totalVotes }} 票</span>
        </div>
      </div>
    </div>
    <div class="card-footer">
      <el-button type="primary" class="btn-primary" @click="handleVote">
        <el-icon><Vote /></el-icon>
        去投票
      </el-button>
      <el-button @click="handleRecords">
        <el-icon><List /></el-icon>
        查看记录
      </el-button>
    </div>
  </el-card>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  activity: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const totalVotes = computed(() => {
  if (!props.activity.options) return 0
  return props.activity.options.reduce((sum, opt) => sum + (opt.vote_count || 0), 0)
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

const handleVote = () => {
  router.push(`/vote/${props.activity.id}`)
}

const handleRecords = () => {
  router.push(`/records/${props.activity.id}`)
}
</script>

<style scoped>
.activity-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  transition: all 0.3s;
}

.card-header {
  padding: 20px 24px;
  border-bottom: 1px solid #ebeef5;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8ecf1 100%);
}

.card-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.card-title h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  flex: 1;
  margin-right: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-type {
  margin-top: 8px;
}

.card-body {
  flex: 1;
  padding: 20px 24px;
}

.description {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 16px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 48px;
}

.meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #909399;
  width: calc(50% - 6px);
}

.meta-item .el-icon {
  color: #667eea;
}

.card-footer {
  padding: 16px 24px;
  border-top: 1px solid #ebeef5;
  display: flex;
  gap: 12px;
}

.card-footer .el-button {
  flex: 1;
}
</style>
