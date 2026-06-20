<template>
  <div class="ranking-list">
    <div class="ranking-header">
      <h3>
        <el-icon><Trophy /></el-icon>
        实时票数榜单
      </h3>
      <span class="refresh-tip">
        <el-icon v-if="autoRefresh" class="refresh-icon" :class="{ 'rotating': isRefreshing }">
          <Refresh />
        </el-icon>
        {{ autoRefresh ? '自动刷新中' : '已停止刷新' }}
      </span>
    </div>
    <div class="ranking-items">
      <transition-group name="rank" tag="div">
        <div
          v-for="(item, index) in rankingData"
          :key="item.id"
          class="ranking-item"
          :class="{ 'top-three': index < 3 }"
        >
          <div class="rank-number" :class="`rank-${index + 1}`">
            {{ item.rank }}
          </div>
          <div class="rank-info">
            <div class="rank-name">{{ item.name }}</div>
            <div class="rank-progress">
              <div
                class="progress-bar"
                :style="{ width: `${item.percentage}%` }"
              ></div>
            </div>
          </div>
          <div class="rank-stats">
            <span class="vote-count">{{ item.vote_count }} 票</span>
            <span class="percentage">{{ item.percentage.toFixed(1) }}%</span>
          </div>
        </div>
      </transition-group>
    </div>
    <div class="ranking-footer">
      <span class="total-votes">
        <el-icon><DataAnalysis /></el-icon>
        总票数：{{ totalVotes }} 票
      </span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  rankingData: {
    type: Array,
    default: () => []
  },
  totalVotes: {
    type: Number,
    default: 0
  },
  autoRefresh: {
    type: Boolean,
    default: true
  },
  isRefreshing: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.ranking-list {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.ranking-header {
  padding: 20px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.ranking-header h3 {
  margin: 0;
  font-size: 18px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.refresh-tip {
  font-size: 13px;
  opacity: 0.9;
  display: flex;
  align-items: center;
  gap: 6px;
}

.refresh-icon {
  transition: transform 0.3s;
}

.refresh-icon.rotating {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.ranking-items {
  padding: 16px 24px;
  min-height: 200px;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
  transition: all 0.3s;
}

.ranking-item:last-child {
  border-bottom: none;
}

.ranking-item:hover {
  background: #f5f7fa;
  margin: 0 -24px;
  padding: 12px 24px;
}

.rank-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #e4e7ed;
  color: #606266;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}

.rank-1 {
  background: linear-gradient(135deg, #ffd700 0%, #ffb700 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(255, 215, 0, 0.4);
}

.rank-2 {
  background: linear-gradient(135deg, #c0c4cc 0%, #a8abb2 100%);
  color: #fff;
}

.rank-3 {
  background: linear-gradient(135deg, #e6a23c 0%, #cf9236 100%);
  color: #fff;
}

.rank-info {
  flex: 1;
  min-width: 0;
}

.rank-name {
  font-weight: 500;
  color: #303133;
  margin-bottom: 6px;
  font-size: 14px;
}

.rank-progress {
  height: 8px;
  background: #f0f2f5;
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px;
  transition: width 0.5s ease;
}

.rank-stats {
  text-align: right;
  flex-shrink: 0;
}

.vote-count {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: #667eea;
}

.percentage {
  font-size: 12px;
  color: #909399;
}

.ranking-footer {
  padding: 16px 24px;
  background: #f5f7fa;
  border-top: 1px solid #ebeef5;
}

.total-votes {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.total-votes .el-icon {
  color: #667eea;
}

.rank-move {
  transition: transform 0.5s ease;
}

.rank-enter-active,
.rank-leave-active {
  transition: all 0.5s ease;
}

.rank-enter-from,
.rank-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}
</style>
