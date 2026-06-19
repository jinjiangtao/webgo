<template>
  <div class="annotation-panel">
    <div class="tool-section">
      <div class="section-title">绘图工具</div>
      <div class="tool-buttons">
        <button
          :class="['tool-btn', { active: currentTool === 'none' }]"
          @click="$emit('select-tool', 'none')"
        >
          <span class="tool-icon">🖱️</span>
          <span class="tool-name">选择</span>
        </button>
        <button
          :class="['tool-btn', { active: currentTool === 'line' }]"
          @click="$emit('select-tool', 'line')"
        >
          <span class="tool-icon">📏</span>
          <span class="tool-name">线条</span>
        </button>
        <button
          :class="['tool-btn', { active: currentTool === 'rect' }]"
          @click="$emit('select-tool', 'rect')"
        >
          <span class="tool-icon">⬜</span>
          <span class="tool-name">矩形</span>
        </button>
        <button
          :class="['tool-btn', { active: currentTool === 'circle' }]"
          @click="$emit('select-tool', 'circle')"
        >
          <span class="tool-icon">⭕</span>
          <span class="tool-name">圆形</span>
        </button>
        <button
          :class="['tool-btn', { active: currentTool === 'text' }]"
          @click="$emit('select-tool', 'text')"
        >
          <span class="tool-icon">📝</span>
          <span class="tool-name">文字</span>
        </button>
      </div>
    </div>

    <div class="settings-section">
      <div class="section-title">样式设置</div>

      <div class="setting-row">
        <label class="setting-label">颜色</label>
        <input
          type="color"
          :value="settings.color"
          @input="$emit('update-settings', 'color', $event.target.value)"
        />
      </div>

      <div class="setting-row">
        <label class="setting-label">线宽</label>
        <div class="setting-control">
          <input
            type="range"
            :value="settings.lineWidth"
            min="1"
            max="20"
            @input="$emit('update-settings', 'lineWidth', Number($event.target.value))"
          />
          <span class="setting-value">{{ settings.lineWidth }}px</span>
        </div>
      </div>
    </div>

    <div v-if="currentTool === 'text'" class="text-section">
      <div class="section-title">文字设置</div>

      <div class="setting-row">
        <label class="setting-label">文字内容</label>
        <input
          type="text"
          class="text-input"
          :value="settings.text"
          placeholder="请输入文字..."
          @input="$emit('update-settings', 'text', $event.target.value)"
        />
      </div>

      <div class="setting-row">
        <label class="setting-label">字体</label>
        <select
          :value="settings.fontFamily"
          @change="$emit('update-settings', 'fontFamily', $event.target.value)"
        >
          <option value="Arial">Arial</option>
          <option value="Helvetica">Helvetica</option>
          <option value="Times New Roman">Times New Roman</option>
          <option value="Georgia">Georgia</option>
          <option value="Verdana">Verdana</option>
          <option value="Courier New">Courier New</option>
          <option value="Microsoft YaHei">微软雅黑</option>
          <option value="SimSun">宋体</option>
          <option value="SimHei">黑体</option>
        </select>
      </div>

      <div class="setting-row">
        <label class="setting-label">字号</label>
        <div class="setting-control">
          <input
            type="range"
            :value="settings.fontSize"
            min="10"
            max="120"
            @input="$emit('update-settings', 'fontSize', Number($event.target.value))"
          />
          <span class="setting-value">{{ settings.fontSize }}px</span>
        </div>
      </div>
    </div>

    <div class="hint-section" v-if="currentTool !== 'none'">
      <div class="hint-box">
        <span class="hint-icon">💡</span>
        <span class="hint-text">{{ currentToolHint }}</span>
      </div>
    </div>

    <button class="clear-btn" @click="$emit('clear')">
      🗑️ 清除所有标注
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  currentTool: {
    type: String,
    default: 'none'
  },
  settings: {
    type: Object,
    default: () => ({})
  }
})

defineEmits(['select-tool', 'update-settings', 'clear'])

const currentToolHint = computed(() => {
  const hints = {
    line: '在画布上拖拽绘制线条',
    rect: '在画布上拖拽绘制矩形',
    circle: '在画布上拖拽绘制圆形',
    text: '在画布上点击添加文字标注'
  }
  return hints[props.currentTool] || ''
})
</script>

<style scoped>
.annotation-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #aaa;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.tool-buttons {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.tool-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px 6px;
  background: #1a1a2e;
  border: 1px solid #2a2a4a;
  border-radius: 8px;
  transition: all 0.2s;
}

.tool-btn:hover {
  border-color: #3a3a5a;
  background: #1e1e36;
}

.tool-btn.active {
  border-color: #6366f1;
  background: rgba(99, 102, 241, 0.15);
}

.tool-icon {
  font-size: 22px;
}

.tool-name {
  font-size: 12px;
  color: #aaa;
}

.tool-btn.active .tool-name {
  color: #6366f1;
}

.settings-section,
.text-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.setting-label {
  font-size: 13px;
  color: #ccc;
  min-width: 60px;
}

.setting-control {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
}

.setting-control input[type="range"] {
  flex: 1;
}

.setting-value {
  font-size: 12px;
  color: #6366f1;
  font-weight: 500;
  min-width: 48px;
  text-align: right;
}

.text-input {
  flex: 1;
  padding: 8px 10px;
  background: #1a1a2e;
  border: 1px solid #2a2a4a;
  border-radius: 6px;
  color: #e0e0e0;
  font-size: 13px;
  outline: none;
  transition: border-color 0.2s;
}

.text-input:focus {
  border-color: #6366f1;
}

.hint-section .hint-box {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 12px;
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 8px;
}

.hint-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.hint-text {
  font-size: 12px;
  color: #a5b4fc;
  line-height: 1.5;
}

.clear-btn {
  padding: 10px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  color: #ef4444;
  font-size: 13px;
  transition: all 0.2s;
  margin-top: auto;
}

.clear-btn:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.5);
}
</style>
