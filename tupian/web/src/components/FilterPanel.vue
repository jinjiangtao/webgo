<template>
  <div class="filter-panel">
    <div class="preset-filters">
      <button
        v-for="preset in presets"
        :key="preset.name"
        :class="['preset-btn', { active: isPresetActive(preset) }]"
        @click="applyPreset(preset)"
      >
        <span class="preset-icon">{{ preset.icon }}</span>
        <span class="preset-name">{{ preset.label }}</span>
      </button>
    </div>

    <div class="filter-sliders">
      <div class="filter-item">
        <div class="filter-header">
          <span class="filter-name">灰度</span>
          <span class="filter-value">{{ filters.grayscale }}%</span>
        </div>
        <input
          type="range"
          :value="filters.grayscale"
          min="0"
          max="100"
          @input="$emit('update', 'grayscale', Number($event.target.value))"
        />
      </div>

      <div class="filter-item">
        <div class="filter-header">
          <span class="filter-name">复古</span>
          <span class="filter-value">{{ filters.sepia }}%</span>
        </div>
        <input
          type="range"
          :value="filters.sepia"
          min="0"
          max="100"
          @input="$emit('update', 'sepia', Number($event.target.value))"
        />
      </div>

      <div class="filter-item">
        <div class="filter-header">
          <span class="filter-name">反相</span>
          <span class="filter-value">{{ filters.invert }}%</span>
        </div>
        <input
          type="range"
          :value="filters.invert"
          min="0"
          max="100"
          @input="$emit('update', 'invert', Number($event.target.value))"
        />
      </div>

      <div class="filter-item">
        <div class="filter-header">
          <span class="filter-name">亮度</span>
          <span class="filter-value">{{ filters.brightness }}%</span>
        </div>
        <input
          type="range"
          :value="filters.brightness"
          min="0"
          max="200"
          @input="$emit('update', 'brightness', Number($event.target.value))"
        />
      </div>

      <div class="filter-item">
        <div class="filter-header">
          <span class="filter-name">对比度</span>
          <span class="filter-value">{{ filters.contrast }}%</span>
        </div>
        <input
          type="range"
          :value="filters.contrast"
          min="0"
          max="200"
          @input="$emit('update', 'contrast', Number($event.target.value))"
        />
      </div>

      <div class="filter-item">
        <div class="filter-header">
          <span class="filter-name">饱和度</span>
          <span class="filter-value">{{ filters.saturate }}%</span>
        </div>
        <input
          type="range"
          :value="filters.saturate"
          min="0"
          max="200"
          @input="$emit('update', 'saturate', Number($event.target.value))"
        />
      </div>

      <div class="filter-item">
        <div class="filter-header">
          <span class="filter-name">模糊</span>
          <span class="filter-value">{{ filters.blur }}px</span>
        </div>
        <input
          type="range"
          :value="filters.blur"
          min="0"
          max="20"
          step="0.5"
          @input="$emit('update', 'blur', Number($event.target.value))"
        />
      </div>

      <div class="filter-item">
        <div class="filter-header">
          <span class="filter-name">高斯模糊</span>
          <span class="filter-value">{{ filters.gaussianBlur }}px</span>
        </div>
        <input
          type="range"
          :value="filters.gaussianBlur"
          min="0"
          max="20"
          step="0.5"
          @input="$emit('update', 'gaussianBlur', Number($event.target.value))"
        />
      </div>
    </div>

    <button class="reset-btn" @click="$emit('reset')">
      🔄 重置所有滤镜
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  filters: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update', 'reset'])

const presets = [
  { name: 'original', label: '原图', icon: '🖼️', values: { grayscale: 0, sepia: 0, invert: 0, brightness: 100, contrast: 100, saturate: 100, blur: 0, gaussianBlur: 0 } },
  { name: 'grayscale', label: '黑白', icon: '⬛', values: { grayscale: 100, sepia: 0, invert: 0, brightness: 100, contrast: 100, saturate: 100, blur: 0, gaussianBlur: 0 } },
  { name: 'sepia', label: '复古', icon: '📜', values: { grayscale: 0, sepia: 80, invert: 0, brightness: 105, contrast: 95, saturate: 100, blur: 0, gaussianBlur: 0 } },
  { name: 'invert', label: '反色', icon: '🔄', values: { grayscale: 0, sepia: 0, invert: 100, brightness: 100, contrast: 100, saturate: 100, blur: 0, gaussianBlur: 0 } },
  { name: 'bright', label: '明亮', icon: '☀️', values: { grayscale: 0, sepia: 0, invert: 0, brightness: 130, contrast: 110, saturate: 120, blur: 0, gaussianBlur: 0 } },
  { name: 'dark', label: '暗调', icon: '🌙', values: { grayscale: 0, sepia: 0, invert: 0, brightness: 75, contrast: 120, saturate: 85, blur: 0, gaussianBlur: 0 } },
  { name: 'vivid', label: '鲜艳', icon: '🎨', values: { grayscale: 0, sepia: 0, invert: 0, brightness: 110, contrast: 130, saturate: 160, blur: 0, gaussianBlur: 0 } },
  { name: 'soft', label: '柔和', icon: '☁️', values: { grayscale: 0, sepia: 10, invert: 0, brightness: 105, contrast: 90, saturate: 85, blur: 1, gaussianBlur: 0 } }
]

const isPresetActive = (preset) => {
  return Object.keys(preset.values).every((key) => props.filters[key] === preset.values[key])
}

const applyPreset = (preset) => {
  Object.keys(preset.values).forEach((key) => {
    emit('update', key, preset.values[key])
  })
}
</script>

<style scoped>
.filter-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.preset-filters {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.preset-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 6px;
  background: #1a1a2e;
  border: 1px solid #2a2a4a;
  border-radius: 8px;
  transition: all 0.2s;
}

.preset-btn:hover {
  border-color: #3a3a5a;
  background: #1e1e36;
}

.preset-btn.active {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
}

.preset-icon {
  font-size: 20px;
}

.preset-name {
  font-size: 11px;
  color: #aaa;
}

.preset-btn.active .preset-name {
  color: #6366f1;
}

.filter-sliders {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-name {
  font-size: 13px;
  color: #ccc;
}

.filter-value {
  font-size: 12px;
  color: #6366f1;
  font-weight: 500;
  min-width: 48px;
  text-align: right;
}

.reset-btn {
  padding: 10px;
  background: #2a2a4a;
  border: 1px solid #3a3a5a;
  border-radius: 8px;
  color: #ccc;
  font-size: 13px;
  transition: all 0.2s;
}

.reset-btn:hover {
  background: #3a3a5a;
  border-color: #4a4a6a;
}
</style>
