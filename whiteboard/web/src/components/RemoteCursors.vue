<template>
  <div class="remote-cursors-layer">
    <div
      v-for="cursor in cursorList"
      :key="cursor.userId"
      class="remote-cursor"
      :style="getCursorStyle(cursor)"
    >
      <div
        class="cursor-arrow"
        :style="{ borderColor: cursor.color }"
      />
      <div
        v-if="cursor.username"
        class="cursor-label"
        :style="{ backgroundColor: cursor.color }"
      >
        {{ cursor.username }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  cursors: {
    type: Array,
    default: () => []
  },
  zoom: {
    type: Number,
    default: 1
  },
  panX: {
    type: Number,
    default: 0
  },
  panY: {
    type: Number,
    default: 0
  }
})

const cursorList = computed(() => {
  return props.cursors
})

function getCursorStyle(cursor) {
  const screenX = cursor.x * props.zoom + props.panX
  const screenY = cursor.y * props.zoom + props.panY
  return {
    transform: `translate(${screenX}px, ${screenY}px)`
  }
}
</script>

<style scoped>
.remote-cursors-layer {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
  z-index: 10;
}

.remote-cursor {
  position: absolute;
  top: 0;
  left: 0;
  transition: transform 0.1s ease-out;
  will-change: transform;
}

.cursor-arrow {
  width: 0;
  height: 0;
  border-left: 8px solid transparent;
  border-right: 8px solid transparent;
  border-bottom: 12px solid;
  transform: rotate(-45deg);
  transform-origin: 0 0;
  filter: drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.3));
}

.cursor-label {
  position: absolute;
  top: 12px;
  left: 8px;
  padding: 2px 6px;
  border-radius: 3px;
  color: #ffffff;
  font-size: 12px;
  line-height: 1.4;
  white-space: nowrap;
  font-weight: 500;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  text-shadow: 0 0 1px rgba(0, 0, 0, 0.3);
}
</style>
