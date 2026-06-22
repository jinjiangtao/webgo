<template>
  <div class="user-panel">
    <div class="panel-section">
      <div class="section-label">白板名称</div>
      <div class="whiteboard-name">
        <el-input
          v-model="editableName"
          size="small"
          @blur="handleNameBlur"
          @keyup.enter="handleNameBlur"
          :placeholder="'未命名白板'"
        />
      </div>
    </div>

    <el-divider class="panel-divider" />

    <div class="panel-section">
      <div class="section-label">在线用户 ({{ onlineUsers.length }})</div>
      <div class="user-list">
        <div
          v-for="user in onlineUsers"
          :key="user.userId"
          class="user-item"
          :title="user.username"
        >
          <div
            class="user-avatar"
            :style="{ backgroundColor: user.color }"
          >
            {{ user.username ? user.username.charAt(0).toUpperCase() : '?' }}
          </div>
          <span class="user-name">{{ user.username }}</span>
        </div>
        <div v-if="onlineUsers.length === 0" class="empty-text">暂无在线用户</div>
      </div>
    </div>

    <el-divider class="panel-divider" />

    <div class="panel-section">
      <div class="section-label">背景颜色</div>
      <div class="background-list">
        <div
          v-for="color in backgroundPresets"
          :key="color"
          class="background-item"
          :class="{ active: background === color }"
          :style="{ backgroundColor: color }"
          :title="color"
          @click="handleBackgroundChange(color)"
        />
      </div>
    </div>

    <el-divider class="panel-divider" />

    <div class="panel-section">
      <div class="section-label">
        <span>缩放</span>
        <span class="zoom-value">{{ (zoom * 100).toFixed(0) }}%</span>
      </div>
      <el-slider
        v-model="zoom"
        :min="0.25"
        :max="3"
        :step="0.25"
        :show-tooltip="false"
        @change="handleZoomChange"
      />
      <div class="zoom-buttons">
        <el-button
          size="small"
          :icon="Minus"
          circle
          @click="adjustZoom(-0.25)"
        />
        <el-button
          size="small"
          @click="handleZoomReset"
        >重置</el-button>
        <el-button
          size="small"
          :icon="Plus"
          circle
          @click="adjustZoom(0.25)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useWhiteboardStore } from '../stores/whiteboard'
import { Minus, Plus } from '@element-plus/icons-vue'

const emit = defineEmits(['rename', 'changeBackground', 'changeZoom'])

const whiteboardStore = useWhiteboardStore()
const { whiteboardName, onlineUsers, background, zoom } = storeToRefs(whiteboardStore)
const { setBackground, setZoom, setWhiteboard } = whiteboardStore

const editableName = ref(whiteboardName.value)

watch(whiteboardName, (val) => {
  editableName.value = val
})

const backgroundPresets = [
  '#ffffff',
  '#f5f5f5',
  '#fff9e6',
  '#e6f7ff',
  '#f6ffed',
  '#fff0f6',
  '#1f1f1f'
]

const handleNameBlur = () => {
  if (editableName.value !== whiteboardName.value) {
    setWhiteboard(whiteboardStore.whiteboardId, editableName.value)
    emit('rename', editableName.value)
  }
}

const handleBackgroundChange = (color) => {
  setBackground(color)
  emit('changeBackground', color)
}

const handleZoomChange = (val) => {
  setZoom(val)
  emit('changeZoom', val)
}

const adjustZoom = (delta) => {
  const newZoom = Math.max(0.25, Math.min(3, zoom.value + delta))
  setZoom(newZoom)
  emit('changeZoom', newZoom)
}

const handleZoomReset = () => {
  setZoom(1)
  emit('changeZoom', 1)
}
</script>

<style scoped>
.user-panel {
  width: 240px;
  background: #ffffff;
  border-left: 1px solid #e4e7ed;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
}

.panel-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-label {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-divider {
  margin: 12px 0;
}

.whiteboard-name {
  width: 100%;
}

.user-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 200px;
  overflow-y: auto;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 8px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.user-item:hover {
  background-color: #f5f7fa;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-name {
  font-size: 13px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty-text {
  font-size: 13px;
  color: #c0c4cc;
  text-align: center;
  padding: 12px 0;
}

.background-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.background-item {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.background-item:hover {
  transform: scale(1.1);
}

.background-item.active {
  border-color: #409eff;
}

.zoom-value {
  font-size: 12px;
  color: #909399;
  font-weight: 400;
}

.zoom-buttons {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}
</style>
