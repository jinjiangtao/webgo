<template>
  <input
    v-if="modelValue"
    ref="inputRef"
    class="canvas-text-input"
    v-model="text"
    :style="inputStyle"
    @keyup.enter="handleConfirm"
    @keyup.esc="handleCancel"
    @blur="handleCancel"
    @mousedown.stop
  />
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  x: {
    type: Number,
    default: 0
  },
  y: {
    type: Number,
    default: 0
  },
  fontSize: {
    type: Number,
    default: 16
  },
  fontFamily: {
    type: String,
    default: 'Arial'
  },
  color: {
    type: String,
    default: '#000000'
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

const emit = defineEmits([
  'update:modelValue',
  'confirm',
  'cancel'
])

const inputRef = ref(null)
const text = ref('')

const inputStyle = computed(() => {
  const screenX = props.x * props.zoom + props.panX
  const screenY = props.y * props.zoom + props.panY
  return {
    left: screenX + 'px',
    top: screenY + 'px',
    fontSize: props.fontSize * props.zoom + 'px',
    fontFamily: props.fontFamily,
    color: props.color,
    width: props.fontSize * props.zoom * 15 + 'px'
  }
})

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      text.value = ''
      nextTick(() => {
        if (inputRef.value) {
          inputRef.value.focus()
        }
      })
    }
  }
)

function handleConfirm() {
  const val = text.value.trim()
  if (val) {
    emit('confirm', val)
  } else {
    emit('cancel')
  }
  emit('update:modelValue', false)
}

function handleCancel() {
  emit('cancel')
  emit('update:modelValue', false)
}
</script>

<style scoped>
.canvas-text-input {
  position: absolute;
  z-index: 20;
  background: transparent;
  border: 1px dashed #409eff;
  outline: none;
  padding: 0 4px;
  min-height: 24px;
  line-height: 1.2;
  transform: translateY(-100%);
  caret-color: auto;
  border-radius: 2px;
}

.canvas-text-input:focus {
  border-style: solid;
  border-color: #409eff;
  background: rgba(255, 255, 255, 0.6);
}
</style>
