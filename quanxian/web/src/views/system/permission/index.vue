<template>
  <div class="permission-container">
    <el-tabs v-model="activeTab" class="permission-tabs">
      <el-tab-pane
        v-if="userStore.hasPermission('system:permission:query') || userStore.hasPermission('system:permission:config')"
        label="权限分配"
        name="assign"
      >
        <el-card class="assign-card" shadow="never">
          <div class="assign-header">
            <el-form :inline="true">
              <el-form-item label="选择角色">
                <el-select
                  v-model="selectedRoleId"
                  placeholder="请选择角色"
                  style="width: 260px"
                  :disabled="!userStore.hasPermission('system:permission:config')"
                  @change="handleRoleChange"
                >
                  <el-option
                    v-for="role in roleList"
                    :key="role.id"
                    :label="role.name"
                    :value="role.id"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-tooltip content="展开全部" placement="top">
                  <el-button :icon="Expand" circle @click="expandAll(true)" />
                </el-tooltip>
                <el-tooltip content="收起全部" placement="top">
                  <el-button :icon="Fold" circle @click="expandAll(false)" />
                </el-tooltip>
              </el-form-item>
            </el-form>
            <div class="stats-info">
              <el-statistic title="已选菜单" :value="selectedMenuCount" />
              <el-divider direction="vertical" />
              <el-statistic title="已选按钮" :value="selectedButtonCount" />
            </div>
          </div>

          <el-divider class="header-divider" />

          <div class="assign-tree-wrapper">
            <el-alert
              v-if="!selectedRoleId"
              title="请先选择一个角色，然后配置该角色的菜单和按钮权限"
              type="info"
              :closable="false"
              show-icon
              style="margin-bottom: 16px"
            />
            <el-tree
              ref="permissionTreeRef"
              v-loading="treeLoading"
              :data="permissionTree"
              :props="{ label: 'label', children: 'children' }"
              show-checkbox
              node-key="id"
              :default-expand-all="true"
              :check-strictly="false"
              :disabled="!userStore.hasPermission('system:permission:config') || !selectedRoleId"
              @check="handleCheckChange"
            >
              <template #default="{ node, data }">
                <span class="tree-node">
                  <el-icon v-if="data.type === 'menu'" class="node-icon menu-icon">
                    <component :is="getMenuIcon(data)" />
                  </el-icon>
                  <el-icon v-else class="node-icon button-icon">
                    <Coin />
                  </el-icon>
                  <span class="node-label">{{ data.label }}</span>
                  <el-tag
                    v-if="data.type === 'button'"
                    size="small"
                    type="warning"
                    effect="plain"
                    class="node-tag"
                  >
                    按钮
                  </el-tag>
                  <el-tag
                    v-else
                    size="small"
                    type="primary"
                    effect="plain"
                    class="node-tag"
                  >
                    菜单
                  </el-tag>
                </span>
              </template>
            </el-tree>
          </div>

          <el-divider class="footer-divider" />

          <div class="assign-footer">
            <el-alert
              title="说明：勾选菜单项会自动勾选其下的所有按钮；取消勾选父菜单会自动取消勾选所有子项。按钮节点ID > 100000 为按钮偏移量。"
              type="warning"
              :closable="false"
              show-icon
            />
            <div class="footer-actions">
              <el-button
                :disabled="!selectedRoleId || isSaving"
                @click="handleResetChecked"
              >
                <el-icon><RefreshLeft /></el-icon>
                重置
              </el-button>
              <el-button
                v-if="userStore.hasPermission('system:permission:config')"
                type="primary"
                :disabled="!selectedRoleId || isSaving"
                :loading="isSaving"
                @click="handleSavePermissions"
              >
                <el-icon><Check /></el-icon>
                批量保存
              </el-button>
            </div>
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane
        v-if="userStore.hasPermission('system:permission:query')"
        label="权限预览"
        name="preview"
      >
        <el-card class="preview-card" shadow="never">
          <div class="preview-header">
            <el-radio-group v-model="previewMode" @change="handlePreviewModeChange">
              <el-radio-button value="role">按角色预览</el-radio-button>
              <el-radio-button value="user">按用户预览</el-radio-button>
            </el-radio-group>
            <el-select
              v-if="previewMode === 'role'"
              v-model="previewRoleId"
              placeholder="请选择角色"
              style="width: 260px; margin-left: 16px"
              filterable
              @change="loadPreviewData"
            >
              <el-option
                v-for="role in roleList"
                :key="role.id"
                :label="role.name"
                :value="role.id"
              />
            </el-select>
            <el-select
              v-else
              v-model="previewUserId"
              placeholder="请选择用户"
              style="width: 260px; margin-left: 16px"
              filterable
              @change="loadPreviewData"
            >
              <el-option
                v-for="user in userList"
                :key="user.id"
                :label="`${user.nickname || user.username} (${user.username})`"
                :value="user.id"
              />
            </el-select>
          </div>

          <el-divider class="header-divider" />

          <div v-loading="previewLoading" class="preview-body">
            <template v-if="previewData">
              <div class="preview-sidebar">
                <div class="sidebar-title">
                  <el-icon><Menu /></el-icon>
                  <span>侧边菜单预览</span>
                  <el-tag size="small" type="success" effect="light" style="margin-left: auto">
                    {{ previewMenuCount }} 个菜单
                  </el-tag>
                </div>
                <div class="sidebar-menu">
                  <el-scrollbar height="100%">
                    <template v-if="previewData.menus && previewData.menus.length > 0">
                      <div
                        v-for="menu in previewData.menus"
                        :key="menu.id"
                        class="menu-group"
                      >
                        <div
                          class="menu-item parent-menu"
                          @click="toggleMenuExpand(menu.id)"
                        >
                          <el-icon class="menu-icon">
                            <component :is="getMenuIcon(menu)" />
                          </el-icon>
                          <span class="menu-text">{{ menu.name }}</span>
                          <el-icon
                            v-if="menu.children && menu.children.length > 0"
                            class="expand-icon"
                            :class="{ expanded: expandedMenuIds.includes(menu.id) }"
                          >
                            <ArrowRight />
                          </el-icon>
                        </div>
                        <div
                          v-if="menu.children && menu.children.length > 0 && expandedMenuIds.includes(menu.id)"
                          class="submenu-list"
                        >
                          <div
                            v-for="child in menu.children"
                            :key="child.id"
                            class="menu-item child-menu"
                            :class="{ active: activePreviewMenu === child.id }"
                            @click="selectPreviewMenu(child)"
                          >
                            <span class="menu-dot" />
                            <span class="menu-text">{{ child.name }}</span>
                            <el-tag
                              v-if="getMenuButtons(child.path).length > 0"
                              size="small"
                              type="primary"
                              effect="light"
                            >
                              {{ getMenuButtons(child.path).length }}
                            </el-tag>
                          </div>
                        </div>
                        <div
                          v-if="(!menu.children || menu.children.length === 0)"
                          class="submenu-list"
                        >
                          <div
                            class="menu-item child-menu"
                            :class="{ active: activePreviewMenu === menu.id }"
                            @click="selectPreviewMenu(menu)"
                          >
                            <span class="menu-dot" />
                            <span class="menu-text">{{ menu.name }}</span>
                            <el-tag
                              v-if="getMenuButtons(menu.path).length > 0"
                              size="small"
                              type="primary"
                              effect="light"
                            >
                              {{ getMenuButtons(menu.path).length }}
                            </el-tag>
                          </div>
                        </div>
                      </div>
                    </template>
                    <el-empty v-else description="暂无可见菜单" :image-size="80" />
                  </el-scrollbar>
                </div>
              </div>

              <div class="preview-content">
                <div class="content-header">
                  <div class="title-section">
                    <el-icon class="header-icon"><Operation /></el-icon>
                    <span class="header-title">
                      {{ currentPreviewMenu ? currentPreviewMenu.name : '请选择左侧菜单' }}
                    </span>
                  </div>
                  <div v-if="currentPreviewMenu" class="path-section">
                    <el-tag size="small" type="info" effect="plain">
                      路径: {{ currentPreviewMenu.path || '-' }}
                    </el-tag>
                  </div>
                </div>

                <div class="content-body">
                  <template v-if="currentPreviewMenu">
                    <div class="buttons-section">
                      <div class="section-title">
                        <el-icon><Coin /></el-icon>
                        <span>可用按钮权限</span>
                        <el-tag
                          size="small"
                          type="success"
                          effect="light"
                          style="margin-left: auto"
                        >
                          共 {{ getMenuButtons(currentPreviewMenu.path).length }} 个
                        </el-tag>
                      </div>
                      <div class="buttons-grid">
                        <template v-if="getMenuButtons(currentPreviewMenu.path).length > 0">
                          <div
                            v-for="(code, idx) in getMenuButtons(currentPreviewMenu.path)"
                            :key="idx"
                            class="button-item"
                          >
                            <el-icon class="btn-icon"><Select /></el-icon>
                            <el-tooltip :content="code" placement="top" show-after="200">
                              <el-tag
                                type="primary"
                                effect="dark"
                                round
                                class="btn-tag"
                              >
                                {{ formatButtonCode(code) }}
                              </el-tag>
                            </el-tooltip>
                            <span class="btn-code">{{ code }}</span>
                          </div>
                        </template>
                        <el-empty
                          v-else
                          description="该菜单下暂无可用按钮"
                          :image-size="60"
                        />
                      </div>
                    </div>

                    <div class="mock-section">
                      <div class="section-title">
                        <el-icon><Picture /></el-icon>
                        <span>页面操作区模拟</span>
                      </div>
                      <div class="mock-toolbar">
                        <el-button
                          v-for="(code, idx) in getMenuButtons(currentPreviewMenu.path)"
                          :key="'mock-' + idx"
                          :type="getMockButtonType(code)"
                          size="small"
                        >
                          <el-icon>
                            <component :is="getMockButtonIcon(code)" />
                          </el-icon>
                          {{ getMockButtonLabel(code) }}
                        </el-button>
                        <el-button
                          v-if="getMenuButtons(currentPreviewMenu.path).length === 0"
                          disabled
                        >
                          无可用操作
                        </el-button>
                      </div>
                      <div class="mock-table">
                        <el-table :data="mockTableData" size="small" border stripe>
                          <el-table-column prop="id" label="ID" width="60" align="center" />
                          <el-table-column prop="name" label="名称" />
                          <el-table-column prop="status" label="状态" width="80" align="center">
                            <template #default>
                              <el-tag size="small" type="success">正常</el-tag>
                            </template>
                          </el-table-column>
                          <el-table-column
                            v-if="hasActionButtons(currentPreviewMenu.path)"
                            label="操作"
                            width="180"
                            align="center"
                          >
                            <template #default>
                              <el-button
                                v-if="hasButtonCode(currentPreviewMenu.path, 'edit')"
                                link
                                type="primary"
                                size="small"
                              >
                                编辑
                              </el-button>
                              <el-button
                                v-if="hasButtonCode(currentPreviewMenu.path, 'delete')"
                                link
                                type="danger"
                                size="small"
                              >
                                删除
                              </el-button>
                              <el-button
                                v-if="hasButtonCode(currentPreviewMenu.path, 'view')"
                                link
                                type="info"
                                size="small"
                              >
                                查看
                              </el-button>
                            </template>
                          </el-table-column>
                        </el-table>
                      </div>
                    </div>
                  </template>
                  <el-empty
                    v-else
                    description="请从左侧菜单选择要预览的菜单项"
                    :image-size="100"
                  >
                    <template #description>
                      <p>选择左侧菜单后可预览：</p>
                      <p>1. 该菜单下可操作的按钮列表</p>
                      <p>2. 按钮编码及实际渲染效果</p>
                    </template>
                  </el-empty>
                </div>
              </div>
            </template>
            <el-empty
              v-else
              :description="previewMode === 'role' ? '请选择要预览的角色' : '请选择要预览的用户'"
              :image-size="100"
            />
          </div>

          <el-divider v-if="previewData" class="footer-divider" />

          <div v-if="previewData" class="preview-footer">
            <div class="footer-tags">
              <span class="tags-label">关联角色：</span>
              <el-tag
                v-for="role in previewData.roles"
                :key="role.id"
                type="success"
                effect="light"
                style="margin-right: 6px"
              >
                {{ role.name }}
              </el-tag>
              <el-tag
                v-if="previewData.user"
                type="warning"
                effect="light"
              >
                用户：{{ previewData.user.nickname || previewData.user.username }}
              </el-tag>
            </div>
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Expand,
  Fold,
  RefreshLeft,
  Check,
  Menu,
  ArrowRight,
  Operation,
  Coin,
  Picture,
  Select,
  Document,
  Setting,
  User,
  UserFilled,
  Lock,
  DataBoard,
  HomeFilled,
  OfficeBuilding,
  Search,
  Plus,
  Edit,
  Delete,
  Download,
  Upload,
  Refresh,
  Key
} from '@element-plus/icons-vue'
import { getRoleList } from '@/api/role'
import { getUserList } from '@/api/user'
import {
  getPermissionsTree,
  getPermissionPreview,
  batchAssignPermissions
} from '@/api/permission'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const activeTab = ref('assign')

const roleList = ref([])
const userList = ref([])

const permissionTree = ref([])
const permissionTreeRef = ref(null)
const selectedRoleId = ref(null)
const treeLoading = ref(false)
const isSaving = ref(false)
const currentCheckedKeys = ref([])
const originalCheckedKeys = ref([])

const previewMode = ref('role')
const previewRoleId = ref(null)
const previewUserId = ref(null)
const previewLoading = ref(false)
const previewData = ref(null)
const expandedMenuIds = ref([])
const activePreviewMenu = ref(null)
const currentPreviewMenu = ref(null)

const mockTableData = [
  { id: 1, name: '示例数据一', status: 1 },
  { id: 2, name: '示例数据二', status: 1 },
  { id: 3, name: '示例数据三', status: 1 }
]

const selectedMenuCount = computed(() => {
  return currentCheckedKeys.value.filter(id => id <= 100000).length
})

const selectedButtonCount = computed(() => {
  return currentCheckedKeys.value.filter(id => id > 100000).length
})

const previewMenuCount = computed(() => {
  if (!previewData.value || !previewData.value.menus) return 0
  let count = 0
  const walk = (menus) => {
    for (const m of menus) {
      count++
      if (m.children && m.children.length > 0) walk(m.children)
    }
  }
  walk(previewData.value.menus)
  return count
})

function getMenuIcon(data) {
  const iconMap = {
    'dashboard': DataBoard,
    '首页': HomeFilled,
    '系统管理': Setting,
    '用户管理': User,
    '用户': User,
    '角色管理': UserFilled,
    '角色': UserFilled,
    '权限管理': Lock,
    '权限': Lock,
    '菜单管理': Menu,
    '菜单': Menu,
    '部门管理': OfficeBuilding,
    '部门': OfficeBuilding
  }
  const key = data.name || data.label || ''
  return iconMap[key] || Document
}

function formatButtonCode(code) {
  if (!code) return ''
  const parts = code.split(':')
  if (parts.length > 0) {
    const last = parts[parts.length - 1]
    const map = {
      'query': '查询',
      'add': '新增',
      'edit': '编辑',
      'delete': '删除',
      'view': '查看',
      'export': '导出',
      'import': '导入',
      'config': '配置',
      'reset': '重置'
    }
    return map[last] || last.toUpperCase()
  }
  return code
}

function getMockButtonType(code) {
  if (!code) return 'default'
  if (code.includes('add')) return 'primary'
  if (code.includes('delete') || code.includes('remove')) return 'danger'
  if (code.includes('edit') || code.includes('update') || code.includes('config')) return 'warning'
  if (code.includes('export') || code.includes('import')) return 'success'
  return 'default'
}

function getMockButtonIcon(code) {
  if (!code) return Document
  if (code.includes('add')) return Plus
  if (code.includes('delete') || code.includes('remove')) return Delete
  if (code.includes('edit') || code.includes('update')) return Edit
  if (code.includes('query') || code.includes('view') || code.includes('search')) return Search
  if (code.includes('export')) return Download
  if (code.includes('import')) return Upload
  if (code.includes('reset')) return Refresh
  if (code.includes('config')) return Setting
  if (code.includes('reset') || code.includes('password')) return Key
  return Document
}

function getMockButtonLabel(code) {
  if (!code) return '按钮'
  const parts = code.split(':')
  const last = parts[parts.length - 1]
  const map = {
    'query': '查询',
    'add': '新增',
    'edit': '编辑',
    'delete': '删除',
    'view': '查看',
    'export': '导出',
    'import': '导入',
    'config': '配置',
    'reset': '重置',
    'assign': '分配'
  }
  return map[last] || last.charAt(0).toUpperCase() + last.slice(1)
}

function getMenuButtons(path) {
  if (!previewData.value || !previewData.value.buttons || !path) return []
  return previewData.value.buttons[path] || []
}

function hasButtonCode(path, keyword) {
  const btns = getMenuButtons(path)
  return btns.some(b => b.includes(keyword))
}

function hasActionButtons(path) {
  const btns = getMenuButtons(path)
  return btns.some(b => ['edit', 'delete', 'view'].some(k => b.includes(k)))
}

async function loadRoleList() {
  try {
    const res = await getRoleList({ page: 1, pageSize: 1000 })
    roleList.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function loadUserList() {
  try {
    const res = await getUserList({ page: 1, pageSize: 1000 })
    userList.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function loadPermissionTree() {
  try {
    const res = await getPermissionsTree()
    permissionTree.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

function expandAll(val) {
  if (!permissionTreeRef.value) return
  const store = permissionTreeRef.value.store
  for (const [, node] of Object.entries(store.nodesMap)) {
    node.expanded = val
  }
}

async function handleRoleChange(roleId) {
  if (!roleId) {
    currentCheckedKeys.value = []
    originalCheckedKeys.value = []
    if (permissionTreeRef.value) {
      permissionTreeRef.value.setCheckedKeys([])
    }
    return
  }
  treeLoading.value = true
  try {
    const res = await getPermissionPreview({ roleId })
    const data = res.data || {}
    const keys = []
    const collectMenuIds = (menus) => {
      for (const m of menus || []) {
        keys.push(m.id)
        if (m.children && m.children.length > 0) collectMenuIds(m.children)
      }
    }
    collectMenuIds(data.menus)
    const btnCodeToIdMap = {}
    const collectButtonIds = (nodes) => {
      for (const n of nodes || []) {
        if (n.type === 'button' && n.id > 100000) {
          btnCodeToIdMap[n.label] = n.id
        }
        if (n.children && n.children.length > 0) collectButtonIds(n.children)
      }
    }
    collectButtonIds(permissionTree.value)
    const btnLabelToCodeMap = {}
    const collectButtonLabels = (menus) => {
      for (const m of menus || []) {
        if (m.buttons) {
          for (const b of m.buttons) {
            btnLabelToCodeMap[b.name] = b.code
          }
        }
        if (m.children && m.children.length > 0) collectButtonLabels(m.children)
      }
    }
    collectButtonLabels(data.menus)
    for (const [label, code] of Object.entries(btnLabelToCodeMap)) {
      if (btnCodeToIdMap[label]) {
        keys.push(btnCodeToIdMap[label])
      }
    }
    const allBtnCodes = new Set()
    if (data.buttons) {
      for (const [, codes] of Object.entries(data.buttons)) {
        for (const c of codes) allBtnCodes.add(c)
      }
    }
    const findBtnIds = (nodes) => {
      for (const n of nodes || []) {
        if (n.type === 'button' && n.id > 100000) {
          const btnCode = btnLabelToCodeMap[n.label]
          if (btnCode && allBtnCodes.has(btnCode) && !keys.includes(n.id)) {
            keys.push(n.id)
          }
        }
        if (n.children && n.children.length > 0) findBtnIds(n.children)
      }
    }
    findBtnIds(permissionTree.value)
    currentCheckedKeys.value = [...keys]
    originalCheckedKeys.value = [...keys]
    await nextTick()
    if (permissionTreeRef.value) {
      permissionTreeRef.value.setCheckedKeys(keys)
    }
  } catch (e) {
    console.error(e)
  } finally {
    treeLoading.value = false
  }
}

function handleCheckChange(_data, checkedInfo) {
  currentCheckedKeys.value = (checkedInfo.checkedKeys || []).map(k => Number(k))
}

function handleResetChecked() {
  if (!selectedRoleId.value) return
  currentCheckedKeys.value = [...originalCheckedKeys.value]
  if (permissionTreeRef.value) {
    permissionTreeRef.value.setCheckedKeys(originalCheckedKeys.value)
  }
  ElMessage.info('已重置为上次保存的状态')
}

async function handleSavePermissions() {
  if (!selectedRoleId.value) {
    ElMessage.warning('请先选择角色')
    return
  }
  const allChecked = currentCheckedKeys.value.map(k => Number(k))
  const menuIds = allChecked.filter(id => id <= 100000)
  const buttonIds = allChecked.filter(id => id > 100000)
  try {
    await ElMessageBox.confirm(
      `确认保存角色权限配置？\n菜单：${menuIds.length} 项，按钮：${buttonIds.length} 项`,
      '保存确认',
      {
        confirmButtonText: '确定保存',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch (e) {
    return
  }
  isSaving.value = true
  try {
    await batchAssignPermissions({
      roleId: selectedRoleId.value,
      menuIds,
      buttonIds
    })
    ElMessage.success('权限配置保存成功')
    originalCheckedKeys.value = [...currentCheckedKeys.value]
  } catch (e) {
    console.error(e)
    ElMessage.error('保存失败，请稍后重试')
  } finally {
    isSaving.value = false
  }
}

function handlePreviewModeChange() {
  previewRoleId.value = null
  previewUserId.value = null
  previewData.value = null
  currentPreviewMenu.value = null
  activePreviewMenu.value = null
  expandedMenuIds.value = []
}

async function loadPreviewData() {
  const id = previewMode.value === 'role' ? previewRoleId.value : previewUserId.value
  if (!id) {
    previewData.value = null
    return
  }
  previewLoading.value = true
  try {
    const params = previewMode.value === 'role'
      ? { roleId: id }
      : { userId: id }
    const res = await getPermissionPreview(params)
    previewData.value = res.data || null
    const firstMenuIds = []
    const walkFirst = (menus) => {
      for (const m of menus || []) {
        firstMenuIds.push(m.id)
        if (m.children && m.children.length > 0) return
      }
    }
    walkFirst(previewData.value?.menus || [])
    expandedMenuIds.value = firstMenuIds
    const findFirstLeaf = (menus) => {
      for (const m of menus || []) {
        if (m.children && m.children.length > 0) {
          const found = findFirstLeaf(m.children)
          if (found) return found
        } else {
          return m
        }
      }
      return null
    }
    const firstLeaf = findFirstLeaf(previewData.value?.menus || [])
    if (firstLeaf) {
      selectPreviewMenu(firstLeaf)
    }
  } catch (e) {
    console.error(e)
  } finally {
    previewLoading.value = false
  }
}

function toggleMenuExpand(menuId) {
  const idx = expandedMenuIds.value.indexOf(menuId)
  if (idx >= 0) {
    expandedMenuIds.value.splice(idx, 1)
  } else {
    expandedMenuIds.value.push(menuId)
  }
}

function selectPreviewMenu(menu) {
  currentPreviewMenu.value = menu
  activePreviewMenu.value = menu.id
}

onMounted(async () => {
  await Promise.all([
    loadRoleList(),
    loadUserList(),
    loadPermissionTree()
  ])
})
</script>

<style lang="scss" scoped>
.permission-container {
  padding: 16px;

  :deep(.el-tabs__item) {
    font-size: 15px;
    height: 44px;
    line-height: 44px;
  }

  :deep(.el-tabs__active-bar) {
    height: 3px;
  }
}

.assign-card,
.preview-card {
  border-radius: 8px;

  :deep(.el-card__body) {
    padding: 20px;
  }
}

.assign-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;

  .stats-info {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 16px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 8px;
    color: #fff;

    :deep(.el-statistic__head) {
      color: rgba(255, 255, 255, 0.85);
      font-size: 12px;
    }

    :deep(.el-statistic__content) {
      color: #fff;
      font-size: 22px;
      font-weight: 600;
    }

    :deep(.el-divider) {
      border-color: rgba(255, 255, 255, 0.3);
    }
  }
}

.header-divider {
  margin: 16px 0;
}

.footer-divider {
  margin: 20px 0 16px;
}

.assign-tree-wrapper {
  min-height: 300px;
  max-height: 520px;
  overflow: auto;
  padding: 8px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  background: var(--el-fill-color-lighter);

  :deep(.el-tree-node__content) {
    height: 36px;
    border-radius: 4px;
    margin-bottom: 2px;

    &:hover {
      background-color: var(--el-color-primary-light-9);
    }
  }

  :deep(.el-tree-node.is-current > .el-tree-node__content) {
    background-color: var(--el-color-primary-light-8);
  }

  .tree-node {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;

    .node-icon {
      width: 18px;
      height: 18px;

      &.menu-icon {
        color: #409eff;
      }

      &.button-icon {
        color: #e6a23c;
      }
    }

    .node-label {
      font-weight: 500;
    }

    .node-tag {
      transform: scale(0.85);
    }
  }
}

.assign-footer {
  display: flex;
  flex-direction: column;
  gap: 16px;

  .footer-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}

.preview-header {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.preview-body {
  display: flex;
  gap: 16px;
  min-height: 500px;
}

.preview-sidebar {
  width: 280px;
  flex-shrink: 0;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  background: linear-gradient(180deg, #1f2937 0%, #111827 100%);
  display: flex;
  flex-direction: column;
  overflow: hidden;

  .sidebar-title {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 14px 16px;
    color: #fff;
    font-weight: 600;
    font-size: 14px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .sidebar-menu {
    flex: 1;
    padding: 8px;
    overflow: hidden;

    :deep(.el-scrollbar__wrap) {
      overflow-x: hidden;
    }

    .menu-group {
      margin-bottom: 2px;
    }

    .menu-item {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 10px 12px;
      border-radius: 6px;
      cursor: pointer;
      color: rgba(255, 255, 255, 0.8);
      font-size: 14px;
      transition: all 0.2s;

      &:hover {
        background: rgba(255, 255, 255, 0.1);
        color: #fff;
      }

      &.parent-menu {
        font-weight: 600;
        color: rgba(255, 255, 255, 0.9);

        .menu-icon {
          color: #60a5fa;
        }
      }

      &.child-menu {
        color: rgba(255, 255, 255, 0.7);

        &.active {
          background: linear-gradient(90deg, #3b82f6 0%, #6366f1 100%);
          color: #fff;
          font-weight: 500;
        }
      }

      .menu-icon {
        width: 16px;
        height: 16px;
      }

      .menu-dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.3);
        margin-left: 4px;
      }

      .menu-text {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .expand-icon {
        transition: transform 0.2s;
        font-size: 12px;

        &.expanded {
          transform: rotate(90deg);
        }
      }
    }

    .submenu-list {
      padding-left: 12px;
      margin-top: 2px;
    }
  }
}

.preview-content {
  flex: 1;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--el-bg-color);

  .content-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 20px;
    background: linear-gradient(90deg, #f0f9ff 0%, #f5f3ff 100%);
    border-bottom: 1px solid var(--el-border-color-lighter);

    .title-section {
      display: flex;
      align-items: center;
      gap: 10px;

      .header-icon {
        color: #409eff;
        font-size: 18px;
      }

      .header-title {
        font-size: 16px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }
    }
  }

  .content-body {
    flex: 1;
    padding: 20px;
    overflow: auto;
  }
}

.buttons-section,
.mock-section {
  margin-bottom: 24px;

  .section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 14px;
    padding-bottom: 10px;
    border-bottom: 2px solid var(--el-border-color-lighter);
    font-weight: 600;
    font-size: 15px;
    color: var(--el-text-color-primary);

    :deep(.el-icon) {
      color: #6366f1;
    }
  }
}

.buttons-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
  padding: 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;

  .button-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 14px 10px;
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 8px;
    transition: all 0.2s;

    &:hover {
      border-color: var(--el-color-primary);
      box-shadow: 0 4px 12px rgba(64, 158, 255, 0.12);
      transform: translateY(-1px);
    }

    .btn-icon {
      color: var(--el-color-primary);
      font-size: 18px;
    }

    .btn-tag {
      font-size: 12px;
      min-width: 60px;
    }

    .btn-code {
      font-size: 11px;
      color: var(--el-text-color-secondary);
      font-family: 'Courier New', monospace;
      background: var(--el-fill-color-lighter);
      padding: 2px 6px;
      border-radius: 3px;
    }
  }
}

.mock-toolbar {
  padding: 12px 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px 8px 0 0;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.mock-table {
  :deep(.el-table) {
    border-radius: 0 0 8px 8px;
  }
}

.preview-footer {
  .footer-tags {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 4px;

    .tags-label {
      font-size: 13px;
      color: var(--el-text-color-secondary);
      margin-right: 4px;
    }
  }
}
</style>
