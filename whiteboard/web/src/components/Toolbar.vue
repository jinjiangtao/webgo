<template>
  <div class="toolbar">
    <div class="toolbar-section">
      <el-button-group>
        <el-tooltip content="画笔" placement="bottom">
          <el-button :type="currentTool === TOOL_TYPES.PEN ? 'primary' : 'default'" @click="setCurrentTool(TOOL_TYPES.PEN)">
            <el-icon><Edit /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="直线" placement="bottom">
          <el-button :type="currentTool === TOOL_TYPES.LINE ? 'primary' : 'default'" @click="setCurrentTool(TOOL_TYPES.LINE)">
            <el-icon><Minus /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="矩形" placement="bottom">
          <el-button :type="currentTool === TOOL_TYPES.RECT ? 'primary' : 'default'" @click="setCurrentTool(TOOL_TYPES.RECT)">
            <el-icon><Grid /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="圆形" placement="bottom">
          <el-button :type="currentTool === TOOL_TYPES.CIRCLE ? 'primary' : 'default'" @click="setCurrentTool(TOOL_TYPES.CIRCLE)">
            <el-icon><Coin /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="椭圆" placement="bottom">
          <el-button :type="currentTool === TOOL_TYPES.ELLIPSE ? 'primary' : 'default'" @click="setCurrentTool(TOOL_TYPES.ELLIPSE)">
            <span style="font-weight:bold;font-style:italic;font-size:16px">O</span>
          </el-button>
        </el-tooltip>
        <el-tooltip content="三角形" placement="bottom">
          <el-button :type="currentTool === TOOL_TYPES.TRIANGLE ? 'primary' : 'default'" @click="setCurrentTool(TOOL_TYPES.TRIANGLE)">
            <el-icon><CaretTop /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="文字" placement="bottom">
          <el-button :type="currentTool === TOOL_TYPES.TEXT ? 'primary' : 'default'" @click="setCurrentTool(TOOL_TYPES.TEXT)">
            <span style="font-weight:bold;font-size:16px">T</span>
          </el-button>
        </el-tooltip>
        <el-tooltip content="橡皮擦" placement="bottom">
          <el-button :type="currentTool === TOOL_TYPES.ERASER ? 'primary' : 'default'" @click="setCurrentTool(TOOL_TYPES.ERASER)">
            <el-icon><Delete /></el-icon>
          </el-button>
        </el-tooltip>
      </el-button-group>
    </div>

    <el-divider direction="vertical" />

    <div class="toolbar-section">
      <div class="toolbar-item">
        <span class="toolbar-label">粗细</span>
        <el-slider v-model="strokeWidth" :min="1" :max="30" :show-tooltip="true" style="width: 100px" @change="setStrokeWidth" />
        <span class="toolbar-value">{{ strokeWidth }}</span>
      </div>
    </div>

    <el-divider direction="vertical" />

    <div class="toolbar-section">
      <div class="toolbar-item">
        <span class="toolbar-label">描边</span>
        <el-color-picker v-model="strokeColor" size="small" @change="setStrokeColor" />
      </div>
      <div class="toolbar-item">
        <span class="toolbar-label">填充</span>
        <el-color-picker v-model="fillColor" size="small" @change="setFillColor" />
      </div>
    </div>

    <el-divider direction="vertical" />

    <div class="toolbar-section">
      <div class="toolbar-item">
        <span class="toolbar-label">透明度</span>
        <el-slider v-model="opacity" :min="0" :max="100" :show-tooltip="true" style="width: 100px" @change="setOpacity" />
        <span class="toolbar-value">{{ opacity }}%</span>
      </div>
    </div>

    <el-divider direction="vertical" />

    <div class="toolbar-section">
      <div class="toolbar-item">
        <span class="toolbar-label">字号</span>
        <el-slider v-model="fontSize" :min="12" :max="72" :step="1" :show-tooltip="true" style="width: 100px" @change="setFontSize" />
        <span class="toolbar-value">{{ fontSize }}px</span>
      </div>
    </div>

    <el-divider direction="vertical" />

    <div class="toolbar-section">
      <el-button-group>
        <el-tooltip content="撤销" placement="bottom">
          <el-button @click="$emit('undo')">
            <el-icon><RefreshLeft /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="重做" placement="bottom">
          <el-button @click="$emit('redo')">
            <el-icon><RefreshRight /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="清空" placement="bottom">
          <el-button type="warning" @click="$emit('clear')">
            <el-icon><DeleteFilled /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="保存" placement="bottom">
          <el-button type="success" @click="$emit('save')">
            <el-icon><DocumentAdd /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="新建" placement="bottom">
          <el-button @click="$emit('new')">
            <el-icon><Plus /></el-icon>
          </el-button>
        </el-tooltip>
        <el-tooltip content="打开" placement="bottom">
          <el-button @click="$emit('open')">
            <el-icon><Folder /></el-icon>
          </el-button>
        </el-tooltip>
      </el-button-group>
    </div>
  </div>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { useToolStore, TOOL_TYPES } from '../stores/tool'
import {
  Edit,
  Minus,
  Grid,
  Coin,
  CaretTop,
  Delete,
  RefreshLeft,
  RefreshRight,
  DeleteFilled,
  DocumentAdd,
  Plus,
  Folder
} from '@element-plus/icons-vue'

defineEmits(['undo', 'redo', 'clear', 'save', 'new', 'open'])

const toolStore = useToolStore()
const {
  currentTool,
  strokeWidth,
  strokeColor,
  fillColor,
  opacity,
  fontSize
} = storeToRefs(toolStore)

const {
  setCurrentTool,
  setStrokeWidth,
  setStrokeColor,
  setFillColor,
  setOpacity,
  setFontSize
} = toolStore
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  background: #ffffff;
  border-bottom: 1px solid #e4e7ed;
  gap: 8px;
  flex-wrap: wrap;
}

.toolbar-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-label {
  font-size: 13px;
  color: #606266;
  white-space: nowrap;
}

.toolbar-value {
  font-size: 12px;
  color: #909399;
  min-width: 40px;
  text-align: right;
}
</style>
