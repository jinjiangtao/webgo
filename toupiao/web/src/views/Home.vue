<template>
  <div class="home page-container">
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><HomeFilled /></el-icon>
        投票活动列表
      </h2>
      <el-button type="primary" class="btn-primary" @click="goCreate">
        <el-icon><Plus /></el-icon>
        创建新活动
      </el-button>
    </div>

    <div class="search-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索活动名称..."
        clearable
        :prefix-icon="Search"
      >
      </el-input>
    </div>

    <el-row :gutter="20" v-loading="loading">
      <el-col
        v-for="activity in filteredActivities"
        :key="activity.id"
        :xs="24"
        :sm="12"
        :md="8"
        :lg="8"
        style="margin-bottom: 20px"
      >
        <ActivityCard :activity="activity" />
      </el-col>
    </el-row>

    <el-empty
      v-if="!loading && filteredActivities.length === 0"
      description="暂无活动数据"
      class="empty-state"
    >
      <el-button type="primary" @click="goCreate">
        <el-icon><Plus /></el-icon>
        立即创建
      </el-button>
    </el-empty>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import ActivityCard from '../components/ActivityCard.vue'
import { getActivityList } from '../api/activity'

const router = useRouter()
const loading = ref(false)
const searchKeyword = ref('')
const activities = ref([])

const filteredActivities = computed(() => {
  if (!searchKeyword.value) return activities.value
  const keyword = searchKeyword.value.toLowerCase()
  return activities.value.filter(act =>
    act.title.toLowerCase().includes(keyword) ||
    (act.description && act.description.toLowerCase().includes(keyword))
  )
})

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

const goCreate = () => {
  router.push('/create')
}

onMounted(() => {
  fetchActivities()
})
</script>

<style scoped>
.home {
  padding: 0 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
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

.search-bar {
  margin-bottom: 24px;
  max-width: 400px;
}

.empty-state {
  padding: 80px 0;
}
</style>
