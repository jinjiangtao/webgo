<template>
  <div class="dashboard-container">
    <el-row :gutter="20" class="mb-20">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon user-icon">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.userCount }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon role-icon">
              <el-icon><Avatar /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.roleCount }}</div>
              <div class="stat-label">角色数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon menu-icon">
              <el-icon><Menu /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.menuCount }}</div>
              <div class="stat-label">菜单数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon dept-icon">
              <el-icon><OfficeBuilding /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.deptCount }}</div>
              <div class="stat-label">部门数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>系统功能介绍</span>
            </div>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="(item, index) in features"
              :key="index"
              :timestamp="item.title"
              :color="item.color"
              placement="top"
            >
              <p>{{ item.desc }}</p>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>快速导航</span>
            </div>
          </template>
          <el-row :gutter="12">
            <el-col :span="8" v-for="nav in quickNav" :key="nav.path">
              <div class="nav-item" @click="goPage(nav.path)">
                <el-icon class="nav-icon" :size="28" :color="nav.color">
                  <component :is="nav.icon" />
                </el-icon>
                <div class="nav-text">{{ nav.title }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getUserList } from '@/api/user'
import { getRoleList } from '@/api/role'
import { getMenuList } from '@/api/menu'
import { getDepartmentList } from '@/api/department'

const router = useRouter()

const stats = ref({
  userCount: 0,
  roleCount: 0,
  menuCount: 0,
  deptCount: 0
})

const features = [
  { title: '多级树形权限', desc: '支持部门、角色、菜单、按钮四级联动权限配置', color: '#409EFF' },
  { title: '树形可视化', desc: '以树形结构展示层级关系，支持拖拽调整', color: '#67C23A' },
  { title: '联动勾选授权', desc: '节点勾选联动，支持批量权限配置', color: '#E6A23C' },
  { title: '动态菜单渲染', desc: '根据用户角色动态渲染侧边菜单', color: '#F56C6C' },
  { title: '按钮级权限', desc: '支持按钮级精细化权限控制', color: '#909399' },
  { title: '权限预览', desc: '实时预览权限配置效果', color: '#409EFF' }
]

const quickNav = [
  { path: '/system/user', title: '用户管理', icon: 'User', color: '#409EFF' },
  { path: '/system/role', title: '角色管理', icon: 'Avatar', color: '#67C23A' },
  { path: '/system/menu', title: '菜单管理', icon: 'Menu', color: '#E6A23C' },
  { path: '/system/dept', title: '部门管理', icon: 'OfficeBuilding', color: '#F56C6C' },
  { path: '/system/permission', title: '权限配置', icon: 'Key', color: '#909399' },
  { path: '/dashboard', title: '系统首页', icon: 'HomeFilled', color: '#409EFF' }
]

async function loadStats() {
  try {
    const [users, roles, menus, depts] = await Promise.all([
      getUserList({ page: 1, pageSize: 1 }),
      getRoleList(),
      getMenuList(),
      getDepartmentList()
    ])
    stats.value.userCount = users.total || 0
    stats.value.roleCount = roles.total || 0
    stats.value.menuCount = menus.total || 0
    stats.value.deptCount = depts.total || 0
  } catch (e) {
    console.error(e)
  }
}

function goPage(path) {
  router.push(path)
}

onMounted(() => {
  loadStats()
})
</script>

<style lang="scss" scoped>
.dashboard-container {
  .stat-card {
    border: none;
    border-radius: 8px;
    overflow: hidden;

    :deep(.el-card__body) {
      padding: 24px;
    }
  }

  .stat-content {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-size: 28px;
  }

  .user-icon { background: linear-gradient(135deg, #667eea, #764ba2); }
  .role-icon { background: linear-gradient(135deg, #11998e, #38ef7d); }
  .menu-icon { background: linear-gradient(135deg, #f093fb, #f5576c); }
  .dept-icon { background: linear-gradient(135deg, #4facfe, #00f2fe); }

  .stat-info {
    flex: 1;

    .stat-value {
      font-size: 28px;
      font-weight: 700;
      color: #303133;
      line-height: 1.2;
    }

    .stat-label {
      font-size: 14px;
      color: #909399;
      margin-top: 4px;
    }
  }

  .card-header {
    font-weight: 600;
    font-size: 16px;
  }

  .nav-item {
    padding: 20px 12px;
    text-align: center;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;
    margin-bottom: 12px;

    &:hover {
      background: #f5f7fa;
      transform: translateY(-2px);
    }

    .nav-icon {
      margin-bottom: 8px;
    }

    .nav-text {
      font-size: 13px;
      color: #606266;
    }
  }
}
</style>
