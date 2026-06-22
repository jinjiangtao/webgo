<template>
  <el-dialog
    :model-value="modelValue"
    :title="dialogTitle"
    width="600px"
    @update:model-value="(val) => emit('update:modelValue', val)"
    @close="handleClose"
  >
    <template v-if="mode === 'open'">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="白板列表" name="whiteboards">
          <el-table
            :data="whiteboardList"
            stripe
            style="width: 100%"
            v-loading="loading"
          >
            <el-table-column prop="name" label="名称" min-width="180" />
            <el-table-column prop="updatedAt" label="更新时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.updatedAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button
                  type="primary"
                  size="small"
                  @click="handleOpen(row)"
                >
                  打开
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click="handleDelete(row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <div v-if="whiteboardList.length === 0 && !loading" class="empty-tip">
            暂无白板，快去创建一个吧！
          </div>
        </el-tab-pane>

        <el-tab-pane label="快照列表" name="snapshots">
          <el-table
            :data="snapshotList"
            stripe
            style="width: 100%"
            v-loading="loading"
          >
            <el-table-column prop="name" label="名称" min-width="180" />
            <el-table-column prop="createdAt" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button
                  type="primary"
                  size="small"
                  @click="handleLoadSnapshot(row)"
                >
                  加载
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <div v-if="snapshotList.length === 0 && !loading" class="empty-tip">
            暂无快照
          </div>
        </el-tab-pane>
      </el-tabs>
    </template>

    <template v-else-if="mode === 'save'">
      <el-form ref="formRef" :model="saveForm" :rules="saveRules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input
            v-model="saveForm.name"
            placeholder="请输入快照名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSaveSnapshot">
            保存
          </el-button>
          <el-button @click="handleClose">取消</el-button>
        </el-form-item>
      </el-form>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listWhiteboards,
  deleteWhiteboard,
  listSnapshots,
  createSnapshot
} from '../utils/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  mode: {
    type: String,
    default: 'open',
    validator: (v) => ['open', 'save'].includes(v)
  },
  whiteboardId: {
    type: String,
    default: ''
  },
  snapshotOperations: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits([
  'update:modelValue',
  'open',
  'delete',
  'save',
  'loadSnapshot'
])

const activeTab = ref('whiteboards')
const loading = ref(false)
const whiteboardList = ref([])
const snapshotList = ref([])

const formRef = ref(null)
const saveForm = ref({
  name: ''
})
const saveRules = {
  name: [
    { required: true, message: '请输入快照名称', trigger: 'blur' }
  ]
}

const dialogTitle = computed(() => {
  if (props.mode === 'open') {
    return activeTab.value === 'whiteboards' ? '打开白板' : '加载快照'
  }
  return '保存快照'
})

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      if (props.mode === 'open') {
        activeTab.value = 'whiteboards'
        loadWhiteboards()
      } else {
        saveForm.value.name = ''
      }
    }
  }
)

watch(activeTab, (val) => {
  if (val === 'snapshots' && props.whiteboardId) {
    loadSnapshots()
  }
})

function formatTime(timeStr) {
  if (!timeStr) return '-'
  const d = new Date(timeStr)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}`
}

async function loadWhiteboards() {
  loading.value = true
  try {
    const res = await listWhiteboards()
    whiteboardList.value = res?.data || res || []
  } catch (e) {
    ElMessage.error('加载白板列表失败: ' + e.message)
    whiteboardList.value = []
  } finally {
    loading.value = false
  }
}

async function loadSnapshots() {
  loading.value = true
  try {
    const res = await listSnapshots(props.whiteboardId)
    snapshotList.value = res?.data || res || []
  } catch (e) {
    ElMessage.error('加载快照列表失败: ' + e.message)
    snapshotList.value = []
  } finally {
    loading.value = false
  }
}

function handleOpen(row) {
  emit('open', row)
  emit('update:modelValue', false)
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      `确定删除白板「${row.name}」吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await deleteWhiteboard(row.id)
    ElMessage.success('删除成功')
    emit('delete', row)
    loadWhiteboards()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败: ' + e.message)
    }
  }
}

function handleLoadSnapshot(row) {
  emit('loadSnapshot', row)
  emit('update:modelValue', false)
}

async function handleSaveSnapshot() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch (e) {
    return
  }

  if (!props.whiteboardId) {
    ElMessage.error('请先创建或打开一个白板')
    return
  }

  loading.value = true
  try {
    const res = await createSnapshot(
      props.whiteboardId,
      saveForm.value.name,
      props.snapshotOperations
    )
    ElMessage.success('保存成功')
    emit('save', res?.data || res)
    emit('update:modelValue', false)
  } catch (e) {
    ElMessage.error('保存失败: ' + e.message)
  } finally {
    loading.value = false
  }
}

function handleClose() {
  emit('update:modelValue', false)
}
</script>

<style scoped>
.empty-tip {
  padding: 40px 0;
  text-align: center;
  color: #909399;
  font-size: 14px;
}
</style>
