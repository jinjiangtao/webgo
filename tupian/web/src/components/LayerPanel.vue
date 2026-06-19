<template>
  <div class="layer-panel">
    <div class="panel-header">
      <span class="panel-title">图层</span>
      <span class="layer-count">{{ layers.length }}</span>
    </div>

    <div class="layer-list" ref="listRef">
      <div
        v-for="(layer, index) in displayLayers"
        :key="layer.id"
        :class="['layer-item', { active: layer.id === activeLayerId }]"
        draggable="true"
        @dragstart="handleDragStart($event, index)"
        @dragover="handleDragOver"
        @drop="handleDrop($event, index)"
        @dragend="handleDragEnd"
        @click="$emit('select', layer.id)"
      >
        <div class="layer-thumb">
          <div v-if="layer.image" class="thumb-placeholder">🖼️</div>
          <div v-else class="thumb-placeholder empty">⬜</div>
        </div>

        <div class="layer-info">
          <div class="layer-name" :title="layer.name">{{ layer.name }}</div>
          <div class="layer-size" v-if="layer.width && layer.height">
            {{ layer.width }} × {{ layer.height }}
          </div>
        </div>

        <div class="layer-actions">
          <button
            :class="['action-icon', { active: layer.visible }]"
            title="显示/隐藏"
            @click.stop="$emit('toggle', layer.id)"
          >
            {{ layer.visible ? '👁️' : '🚫' }}
          </button>
          <button
            class="action-icon"
            title="删除图层"
            @click.stop="$emit('remove', layer.id)"
          >
            🗑️
          </button>
        </div>

        <div class="drag-handle" title="拖拽排序">⋮⋮</div>
      </div>

      <div v-if="layers.length === 0" class="empty-layers">
        <div class="empty-icon">📋</div>
        <div class="empty-text">暂无图层</div>
        <div class="empty-hint">上传图片以添加图层</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  layers: {
    type: Array,
    default: () => []
  },
  activeLayerId: {
    type: [Number, String],
    default: null
  }
})

const emit = defineEmits(['select', 'reorder', 'toggle', 'remove'])

const listRef = ref(null)
const dragIndex = ref(null)
const dragOverIndex = ref(null)

const displayLayers = computed(() => {
  return [...props.layers].reverse()
})

const handleDragStart = (e, index) => {
  dragIndex.value = index
  e.dataTransfer.effectAllowed = 'move'
  e.target.style.opacity = '0.5'
}

const handleDragOver = (e) => {
  e.preventDefault()
  e.dataTransfer.dropEffect = 'move'
}

const handleDrop = (e, dropIndex) => {
  e.preventDefault()
  if (dragIndex.value === null || dragIndex.value === dropIndex) return

  const reversed = [...props.layers].reverse()
  const originalItems = [...props.layers]
  const fromReversed = reversed[dragIndex.value]
  const toReversed = reversed[dropIndex]
  const fromIdx = originalItems.findIndex((l) => l.id === fromReversed.id)
  const toIdx = originalItems.findIndex((l) => l.id === toReversed.id)

  const newLayers = [...originalItems]
  const [removed] = newLayers.splice(fromIdx, 1)
  newLayers.splice(toIdx, 0, removed)

  emit('reorder', newLayers)
}

const handleDragEnd = (e) => {
  e.target.style.opacity = '1'
  dragIndex.value = null
  dragOverIndex.value = null
}
</script>

<style scoped>
.layer-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid #2a2a4a;
}

.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #e0e0e0;
}

.layer-count {
  font-size: 12px;
  color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
  padding: 2px 8px;
  border-radius: 10px;
}

.layer-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.layer-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  background: #1a1a2e;
  border: 1px solid #2a2a4a;
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.layer-item:hover {
  border-color: #3a3a5a;
  background: #1e1e36;
}

.layer-item.active {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.layer-thumb {
  width: 40px;
  height: 40px;
  background: #2a2a4a;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
}

.thumb-placeholder {
  font-size: 20px;
}

.thumb-placeholder.empty {
  opacity: 0.5;
}

.layer-info {
  flex: 1;
  min-width: 0;
}

.layer-name {
  font-size: 13px;
  color: #e0e0e0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.layer-size {
  font-size: 11px;
  color: #666;
  margin-top: 2px;
}

.layer-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.layer-item:hover .layer-actions {
  opacity: 1;
}

.action-icon {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s;
}

.action-icon:hover {
  background: #3a3a5a;
}

.drag-handle {
  color: #444;
  font-size: 12px;
  cursor: grab;
  padding: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.layer-item:hover .drag-handle {
  opacity: 1;
}

.drag-handle:active {
  cursor: grabbing;
}

.empty-layers {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #555;
}

.empty-layers .empty-icon {
  font-size: 40px;
  margin-bottom: 12px;
  opacity: 0.5;
}

.empty-layers .empty-text {
  font-size: 14px;
  margin-bottom: 4px;
}

.empty-layers .empty-hint {
  font-size: 12px;
  opacity: 0.7;
}
</style>
