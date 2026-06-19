<template>
  <div class="toolbar">
    <div class="tool-group">
      <button class="tool-btn" title="撤销" @click="$emit('undo')">
        <span>↩️</span>
      </button>
      <button class="tool-btn" title="重做" @click="$emit('redo')">
        <span>↪️</span>
      </button>
    </div>

    <div class="divider"></div>

    <div class="tool-group">
      <button class="tool-btn" title="裁剪" @click="$emit('crop')">
        <span>✂️</span>
      </button>
    </div>

    <div class="divider"></div>

    <div class="tool-group scale-group">
      <span class="tool-label">缩放</span>
      <input
        type="range"
        :value="scale"
        min="10"
        max="300"
        step="1"
        @input="$emit('scale', Number($event.target.value))"
      />
      <span class="scale-value">{{ scale }}%</span>
    </div>

    <div class="divider"></div>

    <div class="tool-group rotate-group">
      <button class="tool-btn" title="向左旋转90°" @click="$emit('rotate', 'left')">
        <span>↺</span>
      </button>
      <div class="custom-rotate">
        <input
          type="number"
          :value="rotation"
          min="0"
          max="359"
          step="1"
          @input="$emit('rotate', Number($event.target.value))"
        />
        <span class="unit">°</span>
      </div>
      <button class="tool-btn" title="向右旋转90°" @click="$emit('rotate', 'right')">
        <span>↻</span>
      </button>
    </div>

    <div class="divider"></div>

    <div class="tool-group">
      <button class="tool-btn" title="水平翻转" @click="$emit('flip', 'h')">
        <span>⇋</span>
      </button>
      <button class="tool-btn" title="垂直翻转" @click="$emit('flip', 'v')">
        <span>⇵</span>
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  scale: {
    type: Number,
    default: 100
  },
  rotation: {
    type: Number,
    default: 0
  }
})

defineEmits(['scale', 'rotate', 'flip', 'crop', 'undo', 'redo'])
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: #16162a;
  border-bottom: 1px solid #2a2a4a;
  flex-wrap: wrap;
}

.tool-group {
  display: flex;
  align-items: center;
  gap: 6px;
}

.divider {
  width: 1px;
  height: 28px;
  background: #2a2a4a;
  margin: 0 4px;
}

.tool-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #2a2a4a;
  border-radius: 8px;
  transition: all 0.2s;
  font-size: 18px;
}

.tool-btn:hover {
  background: #3a3a5a;
  transform: translateY(-1px);
}

.tool-btn:active {
  transform: translateY(0);
}

.tool-label {
  font-size: 13px;
  color: #888;
  margin-right: 4px;
}

.scale-group {
  min-width: 180px;
}

.scale-group input[type="range"] {
  flex: 1;
  min-width: 100px;
}

.scale-value {
  font-size: 12px;
  color: #6366f1;
  font-weight: 500;
  min-width: 42px;
  text-align: right;
}

.rotate-group .custom-rotate {
  display: flex;
  align-items: center;
  background: #2a2a4a;
  border-radius: 6px;
  padding: 4px 8px;
}

.rotate-group input[type="number"] {
  width: 50px;
  background: transparent;
  border: none;
  color: #e0e0e0;
  font-size: 13px;
  text-align: center;
  outline: none;
  -moz-appearance: textfield;
}

.rotate-group input[type="number"]::-webkit-outer-spin-button,
.rotate-group input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.unit {
  font-size: 12px;
  color: #888;
}
</style>
