<script setup lang="ts">
import { useDashboardStore } from '../stores/dashboard'

const store = useDashboardStore()

const handleRollUp = (index: number) => {
  if (index === store.drillPath.length - 1) return
  store.rollUp(index)
}

const handleRollUpOne = () => {
  store.rollUpOne()
}
</script>

<template>
  <div class="drill-breadcrumb">
    <div class="breadcrumb-label">
      <span>📍</span>
      <span class="label-text">钻取路径</span>
    </div>

    <div class="breadcrumb-items">
      <div
        class="breadcrumb-item global"
        @click="handleRollUp(-1)"
      >
        <span class="item-icon">🌐</span>
        <span class="item-label">全局视图</span>
      </div>

      <span class="breadcrumb-separator" v-if="store.drillPath.length > 0">→</span>

      <template v-for="(level, index) in store.drillPath" :key="index">
        <div
          class="breadcrumb-item"
          :class="{ current: index === store.drillPath.length - 1 }"
          @click="handleRollUp(index)"
        >
          <span class="item-icon">{{
            level.dimension === 'time' ? '📅' :
            level.dimension === 'region' ? '📍' : '💼'
          }}</span>
          <span class="item-label">{{ level.label }}</span>
        </div>
        <span class="breadcrumb-separator" v-if="index < store.drillPath.length - 1">→</span>
      </template>
    </div>

    <div class="breadcrumb-actions" v-if="store.drillPath.length > 1">
      <button class="rollup-btn" @click="handleRollUpOne">
        <span>⬆️</span>
        <span>上卷一级</span>
      </button>
    </div>

    <div class="breadcrumb-hint" v-if="store.canDrillDown">
      <span class="hint-icon">💡</span>
      <span class="hint-text">点击图表可继续下钻</span>
    </div>
  </div>
</template>

<style scoped>
.drill-breadcrumb {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 20px;
  background: linear-gradient(90deg, rgba(59, 130, 246, 0.1), rgba(6, 182, 212, 0.05));
  border-radius: 10px;
  border: 1px solid var(--border-color);
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.breadcrumb-label {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-right: 16px;
  border-right: 1px solid var(--border-color);
}

.label-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  font-family: 'JetBrains Mono', monospace;
}

.breadcrumb-items {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  flex: 1;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  background: rgba(30, 41, 59, 0.8);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 13px;
  color: var(--text-secondary);
  border: 1px solid transparent;
}

.breadcrumb-item:hover:not(.current) {
  background: rgba(59, 130, 246, 0.15);
  border-color: var(--neon-blue);
  color: var(--text-primary);
  box-shadow: 0 0 15px rgba(59, 130, 246, 0.2);
  transform: translateY(-1px);
}

.breadcrumb-item.current {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.3), rgba(6, 182, 212, 0.3));
  border: 1px solid var(--neon-blue);
  color: var(--text-primary);
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.3);
  cursor: default;
}

.breadcrumb-item.global {
  background: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.3);
  color: var(--neon-green);
}

.breadcrumb-item.global:hover {
  background: rgba(16, 185, 129, 0.2);
  border-color: var(--neon-green);
  box-shadow: 0 0 15px rgba(16, 185, 129, 0.2);
}

.item-icon {
  font-size: 14px;
}

.item-label {
  font-family: 'JetBrains Mono', monospace;
  font-weight: 500;
}

.breadcrumb-separator {
  color: var(--text-muted);
  font-size: 14px;
  font-weight: bold;
}

.breadcrumb-actions {
  margin-left: auto;
}

.rollup-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(245, 158, 11, 0.15);
  border: 1px solid rgba(245, 158, 11, 0.4);
  border-radius: 6px;
  color: var(--neon-amber);
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 12px;
  font-family: 'JetBrains Mono', monospace;
}

.rollup-btn:hover {
  background: rgba(245, 158, 11, 0.3);
  box-shadow: 0 0 15px rgba(245, 158, 11, 0.3);
}

.breadcrumb-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 6px;
  animation: pulse 2s infinite;
}

.hint-icon {
  font-size: 14px;
}

.hint-text {
  font-size: 12px;
  color: var(--neon-blue);
  font-family: 'JetBrains Mono', monospace;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}
</style>
