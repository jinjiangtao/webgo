<template>
  <div class="menu-container">
    <el-card class="search-card" shadow="never">
      <el-form :model="searchForm" :inline="true">
        <el-form-item label="菜单名称">
          <el-input v-model="searchForm.name" placeholder="请输入菜单名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
        <el-form-item class="ml-auto">
          <el-button
            v-if="userStore.hasPermission('system:menu:add')"
            type="primary"
            @click="handleAdd(null)"
          >
            <el-icon><Plus /></el-icon>
            新增顶级菜单
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card" shadow="never">
      <el-table
        :data="tableData"
        v-loading="loading"
        border
        stripe
        style="width: 100%"
        row-key="id"
        :expand-row-keys="expandRowKeys"
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="expand-content" v-if="row.buttons && row.buttons.length > 0">
              <div class="expand-title">
                <el-icon><Promotion /></el-icon>
                按钮列表
              </div>
              <el-table :data="row.buttons" size="small" border stripe style="width: 100%; margin-left: 40px;">
                <el-table-column prop="name" label="按钮名称" min-width="120" />
                <el-table-column prop="code" label="按钮编码" min-width="160" />
                <el-table-column prop="sort" label="排序" width="80" align="center" />
                <el-table-column label="状态" width="80" align="center">
                  <template #default="{ row: btn }">
                    <el-tag :type="btn.status === 1 ? 'success' : 'info'" size="small" effect="light">
                      {{ btn.status === 1 ? '启用' : '禁用' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="160" align="center">
                  <template #default="{ row: btn }">
                    <el-button link type="primary" size="small" @click="handleEditButton(btn)">编辑</el-button>
                    <el-button link type="danger" size="small" @click="handleDeleteButton(btn, row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
            <el-empty v-else description="暂无按钮" :image-size="60" />
          </template>
        </el-table-column>
        <el-table-column label="菜单名称" min-width="180">
          <template #default="{ row }">
            <div class="menu-name-cell">
              <el-icon v-if="row.icon" class="menu-icon"><component :is="row.icon" /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="图标" width="80" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.icon" class="table-icon"><component :is="row.icon" /></el-icon>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由地址" min-width="180" show-overflow-tooltip />
        <el-table-column prop="component" label="组件路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column label="显示状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.visible === 1 ? 'success' : 'info'" size="small" effect="light">
              {{ row.visible === 1 ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="启用状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              :disabled="!userStore.hasPermission('system:menu:edit')"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" min-width="160" />
        <el-table-column label="操作" width="320" fixed="right" align="center">
          <template #default="{ row }">
            <el-button
              v-if="userStore.hasPermission('system:menu:add')"
              link
              type="primary"
              @click="handleAdd(row)"
            >
              新增子菜单
            </el-button>
            <el-button
              v-if="userStore.hasPermission('system:menu:edit')"
              link
              type="warning"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="userStore.hasPermission('system:menu:delete')"
              link
              type="danger"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
            <el-button
              link
              type="success"
              @click="handleManageButton(row)"
            >
              管理按钮
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="menuDialogVisible"
      :title="isMenuEdit ? '编辑菜单' : (menuFormData.parentId ? '新增子菜单' : '新增顶级菜单')"
      width="620px"
      destroy-on-close
      @close="resetMenuForm"
    >
      <el-form
        ref="menuFormRef"
        :model="menuFormData"
        :rules="menuFormRules"
        label-width="100px"
      >
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="menuFormData.name" placeholder="请输入菜单名称" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="父级菜单" prop="parentId">
          <el-tree-select
            v-model="menuFormData.parentId"
            :data="menuTreeForSelect"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="不选则为顶级菜单"
            check-strictly
            clearable
            filterable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="路由地址" prop="path">
          <el-input v-model="menuFormData.path" placeholder="请输入路由地址，如：/system/user" />
        </el-form-item>
        <el-form-item label="组件路径" prop="component">
          <el-input v-model="menuFormData.component" placeholder="请输入组件路径，如：system/user/index" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="menuFormData.icon" placeholder="请输入图标名称，如：Setting、User" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number
            v-model="menuFormData.sort"
            :min="0"
            :max="9999"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="显示状态" prop="visible">
          <el-radio-group v-model="menuFormData.visible">
            <el-radio :value="1">显示</el-radio>
            <el-radio :value="0">隐藏</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="启用状态" prop="status">
          <el-radio-group v-model="menuFormData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="menuDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="menuSubmitLoading" @click="handleMenuSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="buttonDialogVisible"
      :title="`按钮管理 - ${currentMenu?.name || ''}`"
      width="800px"
      destroy-on-close
      @close="resetButtonState"
    >
      <div class="button-toolbar">
        <el-button
          type="primary"
          :icon="Plus"
          @click="handleAddButton"
        >
          新增按钮
        </el-button>
      </div>
      <el-table :data="buttonList" v-loading="buttonLoading" border stripe style="width: 100%">
        <el-table-column prop="name" label="按钮名称" min-width="140" />
        <el-table-column prop="code" label="按钮编码" min-width="180" />
        <el-table-column prop="sort" label="排序" width="100" align="center" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" effect="light">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" min-width="160" />
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEditButton(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDeleteButton(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog
      v-model="buttonFormDialogVisible"
      :title="isButtonEdit ? '编辑按钮' : '新增按钮'"
      width="520px"
      destroy-on-close
      @close="resetButtonForm"
    >
      <el-form
        ref="buttonFormRef"
        :model="buttonFormData"
        :rules="buttonFormRules"
        label-width="100px"
      >
        <el-form-item label="按钮名称" prop="name">
          <el-input v-model="buttonFormData.name" placeholder="请输入按钮名称" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="按钮编码" prop="code">
          <el-input v-model="buttonFormData.code" placeholder="请输入按钮编码，如：system:user:add" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number
            v-model="buttonFormData.sort"
            :min="0"
            :max="9999"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="buttonFormData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="buttonFormDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="buttonSubmitLoading" @click="handleButtonSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Promotion } from '@element-plus/icons-vue'
import {
  getMenuTree,
  getMenuList,
  createMenu,
  updateMenu,
  deleteMenu
} from '@/api/menu'
import {
  getButtonList,
  createButton,
  updateButton,
  deleteButton
} from '@/api/button'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const loading = ref(false)
const menuSubmitLoading = ref(false)
const buttonSubmitLoading = ref(false)
const buttonLoading = ref(false)

const menuDialogVisible = ref(false)
const buttonDialogVisible = ref(false)
const buttonFormDialogVisible = ref(false)

const isMenuEdit = ref(false)
const isButtonEdit = ref(false)

const menuFormRef = ref(null)
const buttonFormRef = ref(null)

const expandRowKeys = ref([])
const currentMenu = ref(null)
const currentMenuId = ref(null)
const buttonList = ref([])

const searchForm = reactive({
  name: '',
  status: null
})

const tableData = ref([])
const flatMenuList = ref([])

const menuFormData = reactive({
  id: null,
  name: '',
  path: '',
  component: '',
  icon: '',
  parentId: null,
  sort: 0,
  visible: 1,
  status: 1
})

const buttonFormData = reactive({
  id: null,
  name: '',
  code: '',
  menuId: null,
  sort: 0,
  status: 1
})

const menuFormRules = {
  name: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' }
  ],
  visible: [
    { required: true, message: '请选择显示状态', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择启用状态', trigger: 'change' }
  ]
}

const buttonFormRules = {
  name: [
    { required: true, message: '请输入按钮名称', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入按钮编码', trigger: 'blur' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

const menuTreeForSelect = computed(() => {
  const excludeId = isMenuEdit.value ? menuFormData.id : null
  const filterTree = (list) => {
    return list
      .filter(item => item.id !== excludeId)
      .map(item => ({
        ...item,
        children: item.children ? filterTree(item.children) : []
      }))
  }
  return filterTree(tableData.value)
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

function filterMenuTree(list, keyword, status) {
  return list
    .filter(item => {
      const matchName = !keyword || (item.name && item.name.includes(keyword))
      const matchStatus = status === null || status === undefined || item.status === status
      return matchName || matchStatus
    })
    .map(item => {
      const children = item.children ? filterMenuTree(item.children, keyword, status) : []
      if (children.length > 0) {
        return { ...item, children }
      }
      const matchName = !keyword || (item.name && item.name.includes(keyword))
      const matchStatus = status === null || status === undefined || item.status === status
      if (matchName && matchStatus) {
        return { ...item, children: [] }
      }
      return null
    })
    .filter(item => item !== null)
}

async function loadData() {
  loading.value = true
  try {
    const res = await getMenuTree()
    const treeData = res.data || res || []
    flatMenuList.value = flattenTree(JSON.parse(JSON.stringify(treeData)))
    tableData.value = treeData
  } catch (e) {
    console.error(e)
    ElMessage.error('加载菜单列表失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  loadData()
}

function handleReset() {
  searchForm.name = ''
  searchForm.status = null
  loadData()
}

function resetMenuForm() {
  menuFormData.id = null
  menuFormData.name = ''
  menuFormData.path = ''
  menuFormData.component = ''
  menuFormData.icon = ''
  menuFormData.parentId = null
  menuFormData.sort = 0
  menuFormData.visible = 1
  menuFormData.status = 1
  isMenuEdit.value = false
  nextTick(() => {
    menuFormRef.value?.resetFields()
  })
}

function handleAdd(parent) {
  isMenuEdit.value = false
  resetMenuForm()
  if (parent) {
    menuFormData.parentId = parent.id
  }
  menuDialogVisible.value = true
}

function handleEdit(row) {
  isMenuEdit.value = true
  menuFormData.id = row.id
  menuFormData.name = row.name
  menuFormData.path = row.path || ''
  menuFormData.component = row.component || ''
  menuFormData.icon = row.icon || ''
  menuFormData.parentId = row.parentId || null
  menuFormData.sort = row.sort ?? 0
  menuFormData.visible = row.visible ?? 1
  menuFormData.status = row.status ?? 1
  menuDialogVisible.value = true
  nextTick(() => {
    menuFormRef.value?.clearValidate()
  })
}

async function handleMenuSubmit() {
  if (!menuFormRef.value) return
  try {
    await menuFormRef.value.validate()
  } catch (e) {
    return
  }
  menuSubmitLoading.value = true
  try {
    const data = {
      name: menuFormData.name,
      path: menuFormData.path,
      component: menuFormData.component,
      icon: menuFormData.icon,
      parentId: menuFormData.parentId || null,
      sort: Number(menuFormData.sort),
      visible: Number(menuFormData.visible),
      status: Number(menuFormData.status)
    }
    if (!isMenuEdit.value) {
      await createMenu(data)
      ElMessage.success('新增菜单成功')
    } else {
      await updateMenu(menuFormData.id, data)
      ElMessage.success('编辑菜单成功')
    }
    menuDialogVisible.value = false
    loadData()
  } catch (e) {
    console.error(e)
  } finally {
    menuSubmitLoading.value = false
  }
}

async function handleStatusChange(row) {
  try {
    await updateMenu(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch (e) {
    row.status = row.status === 1 ? 0 : 1
    console.error(e)
  }
}

function hasChildren(row) {
  return row.children && row.children.length > 0
}

async function handleDelete(row) {
  if (hasChildren(row)) {
    ElMessage.warning('该菜单下存在子菜单，不可删除')
    return
  }
  try {
    await ElMessageBox.confirm(
      `确定要删除菜单"${row.name}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch (e) {
    return
  }
  try {
    await deleteMenu(row.id)
    ElMessage.success('删除菜单成功')
    loadData()
  } catch (e) {
    console.error(e)
  }
}

function resetButtonState() {
  currentMenu.value = null
  currentMenuId.value = null
  buttonList.value = []
}

function resetButtonForm() {
  buttonFormData.id = null
  buttonFormData.name = ''
  buttonFormData.code = ''
  buttonFormData.menuId = null
  buttonFormData.sort = 0
  buttonFormData.status = 1
  isButtonEdit.value = false
  nextTick(() => {
    buttonFormRef.value?.resetFields()
  })
}

async function loadButtonList(menuId) {
  buttonLoading.value = true
  try {
    const res = await getButtonList({ menuId, page: 1, pageSize: 1000 })
    buttonList.value = res.data || res || []
  } catch (e) {
    console.error(e)
    ElMessage.error('加载按钮列表失败')
  } finally {
    buttonLoading.value = false
  }
}

async function handleManageButton(row) {
  currentMenu.value = row
  currentMenuId.value = row.id
  buttonDialogVisible.value = true
  await loadButtonList(row.id)
}

function handleAddButton() {
  isButtonEdit.value = false
  resetButtonForm()
  buttonFormData.menuId = currentMenuId.value
  buttonFormDialogVisible.value = true
}

function handleEditButton(row) {
  isButtonEdit.value = true
  buttonFormData.id = row.id
  buttonFormData.name = row.name
  buttonFormData.code = row.code
  buttonFormData.menuId = row.menuId || currentMenuId.value
  buttonFormData.sort = row.sort ?? 0
  buttonFormData.status = row.status ?? 1
  buttonFormDialogVisible.value = true
  nextTick(() => {
    buttonFormRef.value?.clearValidate()
  })
}

async function handleButtonSubmit() {
  if (!buttonFormRef.value) return
  try {
    await buttonFormRef.value.validate()
  } catch (e) {
    return
  }
  buttonSubmitLoading.value = true
  try {
    const data = {
      name: buttonFormData.name,
      code: buttonFormData.code,
      menuId: buttonFormData.menuId,
      sort: Number(buttonFormData.sort),
      status: Number(buttonFormData.status)
    }
    if (!isButtonEdit.value) {
      await createButton(data)
      ElMessage.success('新增按钮成功')
    } else {
      await updateButton(buttonFormData.id, data)
      ElMessage.success('编辑按钮成功')
    }
    buttonFormDialogVisible.value = false
    loadButtonList(currentMenuId.value)
  } catch (e) {
    console.error(e)
  } finally {
    buttonSubmitLoading.value = false
  }
}

async function handleDeleteButton(row, parentMenu) {
  try {
    await ElMessageBox.confirm(
      `确定要删除按钮"${row.name}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch (e) {
    return
  }
  try {
    await deleteButton(row.id)
    ElMessage.success('删除按钮成功')
    if (buttonDialogVisible.value) {
      loadButtonList(currentMenuId.value)
    } else if (parentMenu) {
      if (parentMenu.buttons) {
        const idx = parentMenu.buttons.findIndex(b => b.id === row.id)
        if (idx > -1) {
          parentMenu.buttons.splice(idx, 1)
        }
      }
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
.menu-container {
  padding: 16px;

  .search-card {
    margin-bottom: 16px;
    border-radius: 8px;

    :deep(.el-card__body) {
      padding: 16px 16px 0;
    }

    .ml-auto {
      margin-left: auto;
    }
  }

  .table-card {
    border-radius: 8px;

    :deep(.el-card__body) {
      padding: 16px;
    }
  }

  .menu-name-cell {
    display: flex;
    align-items: center;
    gap: 6px;

    .menu-icon {
      color: #409eff;
      font-size: 16px;
    }
  }

  .table-icon {
    color: #606266;
    font-size: 18px;
  }

  .expand-content {
    padding: 12px 0;

    .expand-title {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 12px;
      margin-left: 40px;
      color: #409eff;
    }
  }

  .button-toolbar {
    margin-bottom: 16px;
  }
}
</style>
