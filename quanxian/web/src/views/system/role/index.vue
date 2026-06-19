<template>
  <div class="role-container">
    <el-card class="search-card" shadow="never">
      <el-form :model="searchForm" :inline="true">
        <el-form-item label="角色名称">
          <el-input v-model="searchForm.name" placeholder="请输入角色名称" clearable />
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
            v-if="userStore.hasPermission('system:role:add')"
            type="primary"
            @click="handleAdd"
          >
            <el-icon><Plus /></el-icon>
            新增角色
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card" shadow="never">
      <el-table :data="tableData" v-loading="loading" border stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="角色名称" min-width="140" />
        <el-table-column prop="code" label="角色编码" min-width="140" />
        <el-table-column prop="remark" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="100" align="center" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              :disabled="!userStore.hasPermission('system:role:edit')"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" min-width="180" />
        <el-table-column label="操作" width="320" fixed="right" align="center">
          <template #default="{ row }">
            <el-button
              v-if="userStore.hasPermission('system:role:edit')"
              link
              type="primary"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="userStore.hasPermission('system:permission:config')"
              link
              type="warning"
              @click="handlePermission(row)"
            >
              分配权限
            </el-button>
            <el-button
              v-if="userStore.hasPermission('system:role:delete')"
              link
              type="danger"
              :disabled="row.code === 'admin'"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑角色' : '新增角色'"
      width="560px"
      destroy-on-close
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="90px"
      >
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入角色名称" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入角色编码（如 system:admin）" :disabled="isEdit" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="描述" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入描述" maxlength="255" show-word-limit />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number
            v-model="formData.sort"
            :min="0"
            :max="9999"
            controls-position="right"
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
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="permissionVisible"
      :title="`分配权限 - ${currentRole?.name || ''}`"
      width="720px"
      destroy-on-close
      @close="resetPermissionForm"
    >
      <div class="permission-container" v-loading="permissionLoading">
        <div class="permission-hint">
          <el-icon color="#409eff"><InfoFilled /></el-icon>
          <span>勾选菜单会自动联动子菜单，按钮可单独勾选配置</span>
        </div>
        <el-scrollbar height="460px" class="permission-scrollbar">
          <el-tree
            ref="menuTreeRef"
            :data="menuTreeWithButtons"
            :props="treeProps"
            node-key="_nodeKey"
            show-checkbox
            :check-strictly="false"
            default-expand-all
            :expand-on-click-node="false"
          >
            <template #default="{ node, data }">
              <div class="tree-node-content">
                <el-icon v-if="data._type === 'menu'" class="node-icon" color="#409eff">
                  <component :is="data.icon || 'Folder'" />
                </el-icon>
                <el-icon v-else class="node-icon" color="#67C23A">
                  <Coin />
                </el-icon>
                <span class="node-label">{{ data._label }}</span>
                <el-tag v-if="data._type === 'button'" size="small" type="success" effect="plain" class="node-tag">
                  按钮
                </el-tag>
              </div>
            </template>
          </el-tree>
        </el-scrollbar>
      </div>
      <template #footer>
        <el-button @click="permissionVisible = false">取消</el-button>
        <el-button type="primary" :loading="permissionSubmitLoading" @click="handlePermissionSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  InfoFilled,
  Coin,
  Folder
} from '@element-plus/icons-vue'
import {
  getRoleList,
  getRole,
  createRole,
  updateRole,
  deleteRole,
  getRoleMenus,
  getRoleButtons,
  bindRoleMenus,
  bindRoleButtons
} from '@/api/role'
import { getMenuTree } from '@/api/menu'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const loading = ref(false)
const submitLoading = ref(false)
const permissionLoading = ref(false)
const permissionSubmitLoading = ref(false)
const dialogVisible = ref(false)
const permissionVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const menuTreeRef = ref(null)
const currentRoleId = ref(null)
const currentRole = ref(null)

const searchForm = reactive({
  name: '',
  status: null
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const tableData = ref([])
const menuTree = ref([])
const menuTreeWithButtons = ref([])

const formData = reactive({
  name: '',
  code: '',
  remark: '',
  sort: 0,
  status: 1
})

const treeProps = {
  children: '_children',
  label: '_label'
}

const formRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

async function loadData() {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    }
    const res = await getRoleList(params)
    tableData.value = res.data || []
    pagination.total = res.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadMenuTree() {
  try {
    const res = await getMenuTree()
    menuTree.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

function buildMenuTreeWithButtons(menus) {
  const result = []
  for (const menu of menus) {
    const menuNode = {
      _nodeKey: `menu_${menu.id}`,
      _type: 'menu',
      _id: menu.id,
      _label: menu.name,
      icon: menu.icon,
      _children: []
    }
    if (menu.children && menu.children.length > 0) {
      menuNode._children = buildMenuTreeWithButtons(menu.children)
    }
    if (menu.buttons && menu.buttons.length > 0) {
      for (const btn of menu.buttons) {
        menuNode._children.push({
          _nodeKey: `button_${btn.id}`,
          _type: 'button',
          _id: btn.id,
          _label: btn.name,
          _children: []
        })
      }
    }
    result.push(menuNode)
  }
  return result
}

function handleSearch() {
  pagination.page = 1
  loadData()
}

function handleReset() {
  searchForm.name = ''
  searchForm.status = null
  pagination.page = 1
  loadData()
}

function resetForm() {
  formData.name = ''
  formData.code = ''
  formData.remark = ''
  formData.sort = 0
  formData.status = 1
  isEdit.value = false
  currentRoleId.value = null
  formRef.value?.resetFields()
}

function handleAdd() {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

async function handleEdit(row) {
  isEdit.value = true
  currentRoleId.value = row.id
  try {
    const res = await getRole(row.id)
    const role = res.data || {}
    formData.name = role.name || ''
    formData.code = role.code || ''
    formData.remark = role.remark || ''
    formData.sort = role.sort ?? 0
    formData.status = role.status ?? 1
    dialogVisible.value = true
  } catch (e) {
    console.error(e)
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch (e) {
    return
  }
  submitLoading.value = true
  try {
    const data = {
      name: formData.name,
      code: formData.code,
      remark: formData.remark,
      sort: Number(formData.sort),
      status: Number(formData.status)
    }
    if (!isEdit.value) {
      await createRole(data)
      ElMessage.success('新增角色成功')
    } else {
      await updateRole(currentRoleId.value, data)
      ElMessage.success('编辑角色成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (e) {
    console.error(e)
  } finally {
    submitLoading.value = false
  }
}

async function handleStatusChange(row) {
  try {
    await updateRole(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch (e) {
    row.status = row.status === 1 ? 0 : 1
    console.error(e)
  }
}

async function handleDelete(row) {
  if (row.code === 'admin') {
    ElMessage.warning('超级管理员不可删除')
    return
  }
  try {
    await ElMessageBox.confirm(
      `确定要删除角色"${row.name}"吗？`,
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
    await deleteRole(row.id)
    ElMessage.success('删除角色成功')
    if (tableData.value.length === 1 && pagination.page > 1) {
      pagination.page--
    }
    loadData()
  } catch (e) {
    console.error(e)
  }
}

async function handlePermission(row) {
  currentRole.value = row
  currentRoleId.value = row.id
  permissionVisible.value = true
  permissionLoading.value = true
  menuTreeWithButtons.value = buildMenuTreeWithButtons(JSON.parse(JSON.stringify(menuTree.value)))
  try {
    const [menuRes, buttonRes] = await Promise.all([
      getRoleMenus(row.id),
      getRoleButtons(row.id)
    ])
    const menuIds = menuRes.data || []
    const buttonIds = buttonRes.data || []
    const checkedKeys = [
      ...menuIds.map(id => `menu_${id}`),
      ...buttonIds.map(id => `button_${id}`)
    ]
    await nextTick()
    menuTreeRef.value?.setCheckedKeys(checkedKeys)
  } catch (e) {
    console.error(e)
  } finally {
    permissionLoading.value = false
  }
}

function resetPermissionForm() {
  currentRole.value = null
  currentRoleId.value = null
  menuTreeWithButtons.value = []
  menuTreeRef.value?.setCheckedKeys([])
}

async function handlePermissionSubmit() {
  if (!menuTreeRef.value || !currentRoleId.value) return
  const checkedKeys = menuTreeRef.value.getCheckedKeys(true)
  const halfCheckedKeys = menuTreeRef.value.getHalfCheckedKeys()
  const allMenuKeys = [...checkedKeys, ...halfCheckedKeys]
  const menuIds = []
  const buttonIds = []
  for (const key of allMenuKeys) {
    if (key.startsWith('menu_')) {
      menuIds.push(Number(key.replace('menu_', '')))
    }
  }
  for (const key of checkedKeys) {
    if (key.startsWith('button_')) {
      buttonIds.push(Number(key.replace('button_', '')))
    }
  }
  permissionSubmitLoading.value = true
  try {
    await Promise.all([
      bindRoleMenus(currentRoleId.value, { menuIds }),
      bindRoleButtons(currentRoleId.value, { buttonIds })
    ])
    ElMessage.success('权限分配成功')
    permissionVisible.value = false
  } catch (e) {
    console.error(e)
  } finally {
    permissionSubmitLoading.value = false
  }
}

onMounted(() => {
  loadData()
  loadMenuTree()
})
</script>

<style lang="scss" scoped>
.role-container {
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

  .pagination-wrapper {
    margin-top: 16px;
    display: flex;
    justify-content: flex-end;
  }
}

.permission-container {
  .permission-hint {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 10px 12px;
    margin-bottom: 12px;
    background: #ecf5ff;
    border-radius: 6px;
    font-size: 13px;
    color: #409eff;
  }

  .permission-scrollbar {
    border: 1px solid #ebeef5;
    border-radius: 6px;
    padding: 8px;
    background: #fafafa;

    :deep(.el-tree) {
      background: transparent;

      .el-tree-node__content {
        height: 36px;
        border-radius: 4px;

        &:hover {
          background-color: #ecf5ff;
        }
      }

      .is-leaf {
        display: inline-block;
        width: 0 !important;
        height: 0 !important;
      }
    }
  }

  .tree-node-content {
    display: flex;
    align-items: center;
    gap: 6px;
    flex: 1;
    padding-right: 8px;

    .node-icon {
      font-size: 14px;
    }

    .node-label {
      font-size: 14px;
      color: #303133;
    }

    .node-tag {
      margin-left: auto;
    }
  }
}
</style>
