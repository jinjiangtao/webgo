<template>
  <div class="dept-container">
    <el-card class="dept-card" shadow="never">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="24" :md="8" :lg="7" :xl="6">
          <div class="tree-panel">
            <div class="panel-header">
              <span class="panel-title">
                <el-icon><OfficeBuilding /></el-icon>
                部门列表
              </span>
              <el-button type="primary" :icon="Plus" size="small" @click="handleAdd(null)">
                新增
              </el-button>
            </div>

            <el-input
              v-model="searchKeyword"
              class="search-input"
              placeholder="搜索部门名称"
              clearable
              :prefix-icon="Search"
            />

            <el-tree
              ref="treeRef"
              class="dept-tree"
              :data="treeData"
              :props="treeProps"
              node-key="id"
              highlight-current
              default-expand-all
              draggable
              :expand-on-click-node="false"
              :filter-node-method="filterNode"
              @node-click="handleNodeClick"
              @node-drop="handleNodeDrop"
            >
              <template #default="{ node, data }">
                <div class="tree-node">
                  <span class="node-label">{{ data.name }}</span>
                  <span v-if="data.status === 0" class="node-tag">
                    <el-tag size="small" type="info" effect="plain">禁用</el-tag>
                  </span>
                </div>
              </template>
            </el-tree>
          </div>
        </el-col>

        <el-col :xs="24" :sm="24" :md="16" :lg="17" :xl="18">
          <div class="detail-panel">
            <div v-if="!currentDept" class="empty-state">
              <el-empty description="请从左侧选择一个部门查看详情" />
            </div>

            <div v-else class="detail-content">
              <div class="detail-header">
                <div class="dept-title">
                  <el-icon class="title-icon"><OfficeBuilding /></el-icon>
                  <span>{{ currentDept.name }}</span>
                  <el-tag
                    :type="currentDept.status === 1 ? 'success' : 'info'"
                    size="default"
                    effect="light"
                    style="margin-left: 12px;"
                  >
                    {{ currentDept.status === 1 ? '启用' : '禁用' }}
                  </el-tag>
                </div>
                <div class="action-btns">
                  <el-button type="primary" :icon="Plus" @click="handleAdd(currentDept)">
                    新增子部门
                  </el-button>
                  <el-button type="warning" :icon="Edit" @click="handleEdit(currentDept)">
                    编辑
                  </el-button>
                  <el-button type="danger" :icon="Delete" @click="handleDelete(currentDept)">
                    删除
                  </el-button>
                </div>
              </div>

              <el-descriptions :column="2" border class="dept-desc">
                <el-descriptions-item label="部门ID">
                  {{ currentDept.id }}
                </el-descriptions-item>
                <el-descriptions-item label="部门名称">
                  {{ currentDept.name }}
                </el-descriptions-item>
                <el-descriptions-item label="父级部门">
                  {{ getParentName(currentDept.parentId) || '顶级部门' }}
                </el-descriptions-item>
                <el-descriptions-item label="排序">
                  {{ currentDept.sort }}
                </el-descriptions-item>
                <el-descriptions-item label="状态">
                  <el-tag :type="currentDept.status === 1 ? 'success' : 'info'" effect="light">
                    {{ currentDept.status === 1 ? '启用' : '禁用' }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="创建时间">
                  {{ currentDept.createdAt || '-' }}
                </el-descriptions-item>
              </el-descriptions>

              <div class="sub-dept-section" v-if="currentDept.children && currentDept.children.length > 0">
                <div class="section-title">
                  <el-icon><Folder /></el-icon>
                  子部门列表 ({{ currentDept.children.length }})
                </div>
                <el-table :data="currentDept.children" stripe border style="width: 100%">
                  <el-table-column prop="id" label="ID" width="80" />
                  <el-table-column prop="name" label="部门名称" min-width="160" />
                  <el-table-column prop="sort" label="排序" width="100" />
                  <el-table-column label="状态" width="100">
                    <template #default="{ row }">
                      <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small" effect="light">
                        {{ row.status === 1 ? '启用' : '禁用' }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="createdAt" label="创建时间" min-width="160" />
                  <el-table-column label="操作" width="220" fixed="right">
                    <template #default="{ row }">
                      <el-button type="primary" link size="small" @click="handleNodeClick(row)">
                        查看
                      </el-button>
                      <el-button type="warning" link size="small" @click="handleEdit(row)">
                        编辑
                      </el-button>
                      <el-button type="danger" link size="small" @click="handleDelete(row)">
                        删除
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="520px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="90px"
        @submit.prevent
      >
        <el-form-item label="部门名称" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入部门名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="父级部门" prop="parentId">
          <el-tree-select
            v-model="formData.parentId"
            :data="treeData"
            :props="treeProps"
            :default-expand-all="true"
            check-strictly
            :node-key="'id'"
            placeholder="请选择父级部门（不选则为顶级部门）"
            style="width: 100%"
            clearable
            filterable
          />
        </el-form-item>

        <el-form-item label="排序" prop="sort">
          <el-input-number
            v-model="formData.sort"
            :min="0"
            :max="9999"
            controls-position="right"
            placeholder="数字越小越靠前"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  Search,
  OfficeBuilding,
  Folder
} from '@element-plus/icons-vue'
import {
  getDepartmentTree,
  getDepartmentList,
  createDepartment,
  updateDepartment,
  deleteDepartment,
  moveDepartment
} from '@/api/department'

const treeRef = ref(null)
const formRef = ref(null)

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('add')
const searchKeyword = ref('')
const treeData = ref([])
const currentDept = ref(null)
const flatList = ref([])

const formData = reactive({
  id: null,
  name: '',
  parentId: null,
  sort: 0,
  status: 1
})

const treeProps = {
  children: 'children',
  label: 'name'
}

const formRules = {
  name: [
    { required: true, message: '请输入部门名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

const dialogTitle = computed(() => {
  return dialogType.value === 'add' ? '新增部门' : '编辑部门'
})

function flattenTree(list, result = []) {
  list.forEach(item => {
    result.push(item)
    if (item.children && item.children.length) {
      flattenTree(item.children, result)
    }
  })
  return result
}

function getParentName(parentId) {
  if (!parentId) return ''
  const parent = flatList.value.find(item => item.id === parentId)
  return parent ? parent.name : ''
}

function filterNode(value, data) {
  if (!value) return true
  return data.name && data.name.includes(value)
}

watch(searchKeyword, val => {
  treeRef.value && treeRef.value.filter(val)
})

async function loadTree() {
  loading.value = true
  try {
    const res = await getDepartmentTree()
    treeData.value = res.data || res || []
    flatList.value = flattenTree(JSON.parse(JSON.stringify(treeData.value)))
  } catch (e) {
    console.error(e)
    ElMessage.error('加载部门列表失败')
  } finally {
    loading.value = false
  }
}

function handleNodeClick(data) {
  currentDept.value = data
}

function handleAdd(parent) {
  dialogType.value = 'add'
  Object.assign(formData, {
    id: null,
    name: '',
    parentId: parent ? parent.id : null,
    sort: 0,
    status: 1
  })
  dialogVisible.value = true
  nextTick(() => {
    formRef.value && formRef.value.clearValidate()
  })
}

function handleEdit(row) {
  dialogType.value = 'edit'
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    parentId: row.parentId || null,
    sort: row.sort ?? 0,
    status: row.status ?? 1
  })
  dialogVisible.value = true
  nextTick(() => {
    formRef.value && formRef.value.clearValidate()
  })
}

function handleDelete(row) {
  const hasChildren = row.children && row.children.length > 0
  const tip = hasChildren
    ? `部门「${row.name}」存在子部门，删除后子部门也将被删除，是否继续？`
    : `确定要删除部门「${row.name}」吗？`

  ElMessageBox.confirm(tip, '删除确认', {
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
    type: 'warning',
    confirmButtonClass: 'el-button--danger'
  }).then(async () => {
    try {
      await deleteDepartment(row.id)
      ElMessage.success('删除成功')
      if (currentDept.value && currentDept.value.id === row.id) {
        currentDept.value = null
      }
      await loadTree()
    } catch (e) {
      console.error(e)
    }
  }).catch(() => {})
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async valid => {
    if (!valid) return

    submitLoading.value = true
    try {
      const payload = {
        name: formData.name,
        parentId: formData.parentId || null,
        sort: Number(formData.sort),
        status: Number(formData.status)
      }

      if (dialogType.value === 'add') {
        await createDepartment(payload)
        ElMessage.success('新增成功')
      } else {
        await updateDepartment(formData.id, payload)
        ElMessage.success('编辑成功')
      }

      dialogVisible.value = false
      await loadTree()

      if (dialogType.value === 'edit' && currentDept.value) {
        await nextTick()
        const updated = flatList.value.find(item => item.id === formData.id)
        if (updated) {
          currentDept.value = updated
          treeRef.value && treeRef.value.setCurrentKey(updated.id)
        }
      }
    } catch (e) {
      console.error(e)
    } finally {
      submitLoading.value = false
    }
  })
}

async function handleNodeDrop(draggingNode, dropNode, dropType) {
  if (dropType === 'inner') {
    ElMessage.warning('暂不支持嵌套到其他部门中间，请使用"之前"或"之后"方式拖动')
    await loadTree()
    return
  }

  const draggingId = draggingNode.data.id
  const targetId = dropNode.data.id

  if (draggingId === targetId) return

  try {
    await moveDepartment({
      id: draggingId,
      targetId: targetId,
      position: dropType
    })
    ElMessage.success('移动成功')
    await loadTree()
  } catch (e) {
    console.error(e)
    ElMessage.error('移动失败')
    await loadTree()
  }
}

onMounted(() => {
  loadTree()
})
</script>

<style lang="scss" scoped>
.dept-container {
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
}

.dept-card {
  border: none;
  border-radius: 8px;
  height: 100%;

  :deep(.el-card__body) {
    padding: 0;
    height: 100%;
  }
}

.tree-panel,
.detail-panel {
  padding: 20px;
  height: calc(100vh - 180px);
  min-height: 600px;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}

.tree-panel {
  border-right: 1px solid #ebeef5;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.search-input {
  margin-bottom: 16px;
}

.dept-tree {
  flex: 1;
  overflow-y: auto;
  padding: 4px;
  background: #fafafa;
  border-radius: 6px;
  border: 1px solid #ebeef5;

  :deep(.el-tree-node__content) {
    height: 36px;
    border-radius: 4px;

    &:hover {
      background-color: #ecf5ff;
    }
  }

  :deep(.is-current > .el-tree-node__content) {
    background-color: #409eff;
    color: #fff;

    .tree-node .node-label {
      color: #fff;
    }
  }
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex: 1;
  padding-right: 8px;
}

.node-label {
  font-size: 14px;
}

.node-tag {
  flex-shrink: 0;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.detail-content {
  flex: 1;
  overflow-y: auto;
}

.detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 20px;
  margin-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
  flex-wrap: wrap;
  gap: 12px;
}

.dept-title {
  display: flex;
  align-items: center;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.title-icon {
  color: #409eff;
  margin-right: 8px;
  font-size: 24px;
}

.action-btns {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.dept-desc {
  margin-bottom: 28px;
}

.sub-dept-section {
  .section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 14px;
  }
}

@media (max-width: 992px) {
  .tree-panel,
  .detail-panel {
    height: auto;
    min-height: 400px;
  }

  .tree-panel {
    border-right: none;
    border-bottom: 1px solid #ebeef5;
  }
}

@media (max-width: 480px) {
  .dept-container {
    padding: 12px;
  }

  .tree-panel,
  .detail-panel {
    padding: 14px;
    min-height: 320px;
  }

  .detail-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .action-btns {
    width: 100%;

    .el-button {
      flex: 1;
    }
  }
}
</style>
