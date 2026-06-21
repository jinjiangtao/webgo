<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import SideNav from '../components/SideNav.vue'
import Toolbar from '../components/Toolbar.vue'
import { getSnapshotList, getSnapshotDetail } from '../api/data'
import type { SnapshotItem, DrillLevel } from '../types'

const router = useRouter()
const loading = ref(false)
const snapshots = ref<SnapshotItem[]>([])

const loadSnapshots = async () => {
  loading.value = true
  try {
    const res = await getSnapshotList()
    if (res.code === 200) {
      snapshots.value = res.data.sort((a, b) =>
        new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      )
    }
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleView = async (snapshot: SnapshotItem) => {
  try {
    const res = await getSnapshotDetail(snapshot.id)
    if (res.code === 200) {
      const data = res.data
      await ElMessageBox.alert(
        `
        <div style="font-family: 'JetBrains Mono', monospace; line-height: 2;">
          <h4 style="color: var(--neon-blue); margin-bottom: 16px;">📸 ${data.name}</h4>
          <div><strong>描述:</strong> ${data.description}</div>
          <div><strong>创建时间:</strong> ${data.createdAt}</div>
          <hr style="border-color: var(--border-color); margin: 12px 0;">
          <div><strong>钻取路径:</strong></div>
          <div style="background: rgba(30, 41, 59, 0.5); padding: 8px 12px; border-radius: 4px; margin: 8px 0;">
            ${data.drillPath.map((d: any) => d.label).join(' → ') || '全局视图'}
          </div>
          <hr style="border-color: var(--border-color); margin: 12px 0;">
          <div><strong>数据记录:</strong> ${data.records?.length || 0} 条</div>
        </div>
        `,
        '快照详情',
        {
          dangerouslyUseHTMLString: true,
          confirmButtonText: '关闭',
          customClass: 'snapshot-dialog',
        }
      )
    }
  } catch (e) {}
}

const handleRestore = async (snapshot: SnapshotItem) => {
  try {
    await ElMessageBox.confirm(
      `确定要恢复快照 "${snapshot.name}" 吗？\n这将覆盖当前的钻取路径和筛选条件。`,
      '恢复快照',
      {
        confirmButtonText: '恢复',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const res = await getSnapshotDetail(snapshot.id)
    if (res.code === 200) {
      sessionStorage.setItem('snapshot_restore', JSON.stringify(res.data))
      router.push('/dashboard')
      ElMessage.success('快照恢复成功！')
    }
  } catch (e) {}
}

const handleDelete = async (snapshot: SnapshotItem) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除快照 "${snapshot.name}" 吗？此操作不可撤销。`,
      '删除快照',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'error',
      }
    )
    snapshots.value = snapshots.value.filter(s => s.id !== snapshot.id)
    ElMessage.success('删除成功')
  } catch (e) {}
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const getTimeAgo = (dateStr: string) => {
  const now = Date.now()
  const date = new Date(dateStr).getTime()
  const diff = now - date

  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  return formatDate(dateStr)
}

onMounted(() => {
  document.documentElement.classList.add('dark')
  loadSnapshots()
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
              <h2 class="page-title">📸 快照管理</h2>
              <p class="page-desc">保存和管理数据分析快照，一键恢复历史视图</p>
            </div>
            <button class="refresh-btn" @click="loadSnapshots" :disabled="loading">
              <span :class="{ 'animate-spin': loading }">🔄</span>
              刷新
            </button>
          </div>

          <div class="empty-state" v-if="snapshots.length === 0 && !loading">
            <span class="empty-icon">📸</span>
            <h3 class="empty-title">暂无快照</h3>
            <p class="empty-desc">
              在数据大屏中点击"快照"按钮，即可保存当前的数据分析视图
            </p>
            <button class="go-btn" @click="router.push('/dashboard')">
              前往数据大屏 →
            </button>
          </div>

          <div v-else class="snapshot-grid" v-loading="loading">
            <div
              v-for="snapshot in snapshots"
              :key="snapshot.id"
              class="snapshot-card glass glass-hover"
            >
              <div class="snapshot-header">
                <span class="snapshot-icon">📸</span>
                <span class="snapshot-time">{{ getTimeAgo(snapshot.createdAt) }}</span>
              </div>

              <h3 class="snapshot-name">{{ snapshot.name }}</h3>
              <p class="snapshot-desc">{{ snapshot.description }}</p>

              <div class="snapshot-path">
                <span class="path-label">📍</span>
                <span class="path-text">
                  {{ (snapshot.drillPath as DrillLevel[])?.map((d) => d.label).join(' → ') || '全局视图' }}
                </span>
              </div>

              <div class="snapshot-meta">
                <span class="meta-item">
                  <span>📊</span>
                  {{ snapshot.recordCount || 0 }} 条数据
                </span>
                <span class="meta-item">
                  <span>🕐</span>
                  {{ formatDate(snapshot.createdAt) }}
                </span>
              </div>

              <div class="snapshot-actions">
                <button class="action-btn view" @click="handleView(snapshot)">
                  👁️ 查看
                </button>
                <button class="action-btn restore" @click="handleRestore(snapshot)">
                  ↩️ 恢复
                </button>
                <button class="action-btn delete" @click="handleDelete(snapshot)">
                  🗑️ 删除
                </button>
              </div>
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
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
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

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.refresh-btn:hover:not(:disabled) {
  border-color: var(--neon-blue);
  color: var(--neon-blue);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
  opacity: 0.6;
}

.empty-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 12px 0;
  font-family: 'JetBrains Mono', monospace;
}

.empty-desc {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0 0 24px 0;
  max-width: 400px;
  line-height: 1.6;
  font-family: 'JetBrains Mono', monospace;
}

.go-btn {
  padding: 12px 28px;
  background: linear-gradient(135deg, var(--neon-blue), var(--neon-cyan));
  border: none;
  border-radius: 6px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.go-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(59, 130, 246, 0.4);
}

.snapshot-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}

.snapshot-card {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.snapshot-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.snapshot-icon {
  font-size: 28px;
}

.snapshot-time {
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}

.snapshot-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  font-family: 'JetBrains Mono', monospace;
  line-height: 1.4;
}

.snapshot-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0;
  line-height: 1.5;
  min-height: 40px;
}

.snapshot-path {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 12px;
  background: rgba(30, 41, 59, 0.6);
  border-radius: 6px;
  border: 1px solid var(--border-color);
}

.path-label {
  font-size: 14px;
  flex-shrink: 0;
}

.path-text {
  font-size: 12px;
  color: var(--text-secondary);
  font-family: 'JetBrains Mono', monospace;
  line-height: 1.5;
  word-break: break-all;
}

.snapshot-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}

.snapshot-actions {
  display: flex;
  gap: 8px;
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 12px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'JetBrains Mono', monospace;
}

.action-btn.view {
  color: var(--neon-blue);
}

.action-btn.view:hover {
  background: rgba(59, 130, 246, 0.15);
  border-color: var(--neon-blue);
}

.action-btn.restore {
  color: var(--neon-green);
}

.action-btn.restore:hover {
  background: rgba(16, 185, 129, 0.15);
  border-color: var(--neon-green);
}

.action-btn.delete {
  color: var(--neon-red);
}

.action-btn.delete:hover {
  background: rgba(239, 68, 68, 0.15);
  border-color: var(--neon-red);
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
