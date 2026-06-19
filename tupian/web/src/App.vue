<template>
  <div class="app-container">
    <header class="app-header">
      <h1 class="app-title">在线图片编辑器</h1>
      <div class="header-actions">
        <label class="upload-btn">
          <span>📁 上传图片</span>
          <input type="file" accept="image/*" multiple @change="handleUpload" hidden />
        </label>
        <button class="action-btn" @click="exportImage('png')">⬇️ 导出PNG</button>
        <button class="action-btn" @click="exportImage('jpeg')">⬇️ 导出JPG</button>
      </div>
    </header>

    <div class="main-content">
      <aside class="left-panel">
        <LayerPanel
          :layers="layers"
          :activeLayerId="activeLayerId"
          @select="selectLayer"
          @reorder="reorderLayers"
          @toggle="toggleLayerVisibility"
          @remove="removeLayer"
        />
      </aside>

      <main class="editor-area">
        <Toolbar
          :scale="scale"
          :rotation="rotation"
          @scale="handleScale"
          @rotate="handleRotate"
          @flip="handleFlip"
          @crop="startCrop"
          @undo="undo"
          @redo="redo"
        />
        <ImageEditor
          ref="editorRef"
          :layers="layers"
          :activeLayerId="activeLayerId"
          :scale="scale"
          :rotation="rotation"
          :flipH="flipH"
          :flipV="flipV"
          :filters="filters"
          :annotations="annotations"
          :isCropping="isCropping"
          :currentTool="currentTool"
          :annotationSettings="annotationSettings"
          @crop-complete="applyCrop"
          @add-annotation="addAnnotation"
          @update-annotation="updateAnnotation"
        />
      </main>

      <aside class="right-panel">
        <div class="panel-tabs">
          <button
            :class="['tab-btn', { active: activeTab === 'filter' }]"
            @click="activeTab = 'filter'"
          >
            🎨 滤镜
          </button>
          <button
            :class="['tab-btn', { active: activeTab === 'annotation' }]"
            @click="activeTab = 'annotation'"
          >
            ✏️ 标注
          </button>
        </div>
        <div class="panel-content">
          <FilterPanel
            v-if="activeTab === 'filter'"
            :filters="filters"
            @update="updateFilters"
            @reset="resetFilters"
          />
          <AnnotationPanel
            v-if="activeTab === 'annotation'"
            :currentTool="currentTool"
            :settings="annotationSettings"
            @select-tool="selectTool"
            @update-settings="updateAnnotationSettings"
            @clear="clearAnnotations"
          />
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import ImageEditor from './components/ImageEditor.vue'
import Toolbar from './components/Toolbar.vue'
import LayerPanel from './components/LayerPanel.vue'
import FilterPanel from './components/FilterPanel.vue'
import AnnotationPanel from './components/AnnotationPanel.vue'

const editorRef = ref(null)

const layers = ref([])
const activeLayerId = ref(null)
const scale = ref(100)
const rotation = ref(0)
const flipH = ref(false)
const flipV = ref(false)
const activeTab = ref('filter')
const isCropping = ref(false)
const currentTool = ref('none')

const filters = reactive({
  grayscale: 0,
  sepia: 0,
  invert: 0,
  brightness: 100,
  contrast: 100,
  saturate: 100,
  blur: 0,
  gaussianBlur: 0
})

const annotationSettings = reactive({
  color: '#ef4444',
  lineWidth: 3,
  fontSize: 24,
  fontFamily: 'Arial',
  text: ''
})

const annotations = ref([])

const history = ref([])
const historyIndex = ref(-1)

let layerCounter = 0

const saveHistory = () => {
  const state = JSON.stringify({
    layers: layers.value,
    annotations: annotations.value,
    filters: { ...filters }
  })
  history.value = history.value.slice(0, historyIndex.value + 1)
  history.value.push(state)
  historyIndex.value = history.value.length - 1
  if (history.value.length > 50) {
    history.value.shift()
    historyIndex.value--
  }
}

const undo = () => {
  if (historyIndex.value > 0) {
    historyIndex.value--
    const state = JSON.parse(history.value[historyIndex.value])
    layers.value = state.layers
    annotations.value = state.annotations
    Object.assign(filters, state.filters)
  }
}

const redo = () => {
  if (historyIndex.value < history.value.length - 1) {
    historyIndex.value++
    const state = JSON.parse(history.value[historyIndex.value])
    layers.value = state.layers
    annotations.value = state.annotations
    Object.assign(filters, state.filters)
  }
}

const handleUpload = (e) => {
  const files = Array.from(e.target.files)
  files.forEach((file) => {
    const reader = new FileReader()
    reader.onload = (event) => {
      const img = new Image()
      img.onload = () => {
        const layer = {
          id: ++layerCounter,
          name: file.name,
          image: img,
          visible: true,
          x: 0,
          y: 0,
          width: img.width,
          height: img.height
        }
        layers.value.push(layer)
        if (!activeLayerId.value) {
          activeLayerId.value = layer.id
        }
        saveHistory()
      }
      img.src = event.target.result
    }
    reader.readAsDataURL(file)
  })
  e.target.value = ''
}

const selectLayer = (id) => {
  activeLayerId.value = id
}

const reorderLayers = (newLayers) => {
  layers.value = newLayers
  saveHistory()
}

const toggleLayerVisibility = (id) => {
  const layer = layers.value.find((l) => l.id === id)
  if (layer) {
    layer.visible = !layer.visible
    saveHistory()
  }
}

const removeLayer = (id) => {
  const idx = layers.value.findIndex((l) => l.id === id)
  if (idx > -1) {
    layers.value.splice(idx, 1)
    if (activeLayerId.value === id) {
      activeLayerId.value = layers.value.length > 0 ? layers.value[layers.value.length - 1].id : null
    }
    saveHistory()
  }
}

const handleScale = (val) => {
  scale.value = val
}

const handleRotate = (val) => {
  if (val === 'left') {
    rotation.value = (rotation.value - 90 + 360) % 360
  } else if (val === 'right') {
    rotation.value = (rotation.value + 90) % 360
  } else {
    rotation.value = val
  }
}

const handleFlip = (dir) => {
  if (dir === 'h') {
    flipH.value = !flipH.value
  } else {
    flipV.value = !flipV.value
  }
}

const startCrop = () => {
  isCropping.value = !isCropping.value
}

const applyCrop = (cropData) => {
  if (editorRef.value) {
    editorRef.value.applyCrop(cropData)
    saveHistory()
  }
  isCropping.value = false
}

const updateFilters = (key, val) => {
  filters[key] = val
}

const resetFilters = () => {
  filters.grayscale = 0
  filters.sepia = 0
  filters.invert = 0
  filters.brightness = 100
  filters.contrast = 100
  filters.saturate = 100
  filters.blur = 0
  filters.gaussianBlur = 0
}

const selectTool = (tool) => {
  currentTool.value = tool
  isCropping.value = false
}

const updateAnnotationSettings = (key, val) => {
  annotationSettings[key] = val
}

const addAnnotation = (annotation) => {
  annotations.value.push(annotation)
  saveHistory()
}

const updateAnnotation = (annotation) => {
  const idx = annotations.value.findIndex((a) => a.id === annotation.id)
  if (idx > -1) {
    annotations.value[idx] = annotation
  }
}

const clearAnnotations = () => {
  annotations.value = []
  saveHistory()
}

const exportImage = (format) => {
  if (editorRef.value) {
    editorRef.value.exportImage(format)
  }
}
</script>

<style scoped>
.app-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1a1a2e;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  background: #16162a;
  border-bottom: 1px solid #2a2a4a;
}

.app-title {
  font-size: 20px;
  font-weight: 600;
  color: #e0e0e0;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.upload-btn,
.action-btn {
  padding: 8px 16px;
  background: #6366f1;
  color: white;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.2s;
}

.upload-btn:hover,
.action-btn:hover {
  background: #5856d6;
  transform: translateY(-1px);
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.left-panel,
.right-panel {
  width: 280px;
  background: #16162a;
  border-right: 1px solid #2a2a4a;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.right-panel {
  border-right: none;
  border-left: 1px solid #2a2a4a;
}

.editor-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-tabs {
  display: flex;
  border-bottom: 1px solid #2a2a4a;
}

.tab-btn {
  flex: 1;
  padding: 14px;
  font-size: 14px;
  color: #888;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;
}

.tab-btn.active {
  color: #6366f1;
  border-bottom-color: #6366f1;
  background: rgba(99, 102, 241, 0.05);
}

.tab-btn:hover {
  color: #aaa;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}
</style>
