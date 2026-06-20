<template>
  <div id="app">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1 class="logo" @click="goHome">
            <el-icon><Box /></el-icon>
            <span>投票系统</span>
          </h1>
          <el-menu
            :default-active="activeMenu"
            mode="horizontal"
            router
            class="nav-menu"
          >
            <el-menu-item index="/">
              <el-icon><HomeFilled /></el-icon>
              <span>首页</span>
            </el-menu-item>
            <el-menu-item index="/create">
              <el-icon><Plus /></el-icon>
              <span>创建活动</span>
            </el-menu-item>
            <el-menu-item index="/admin">
              <el-icon><Setting /></el-icon>
              <span>管理后台</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
      <el-footer class="footer">
        <p>© 2024 投票系统 - 让每一票都有价值</p>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const activeMenu = computed(() => route.path)

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 0;
  height: 70px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.header-content {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #fff;
  font-size: 22px;
  font-weight: 600;
  margin: 0;
  cursor: pointer;
  transition: transform 0.3s;
}

.logo:hover {
  transform: scale(1.05);
}

.logo .el-icon {
  font-size: 28px;
}

.nav-menu {
  background: transparent;
  border: none;
}

.nav-menu .el-menu-item {
  color: rgba(255, 255, 255, 0.9);
  border-bottom: 2px solid transparent;
  transition: all 0.3s;
}

.nav-menu .el-menu-item:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.nav-menu .el-menu-item.is-active {
  color: #fff;
  border-bottom-color: #ffd700;
}

.main {
  background: #f5f7fa;
  min-height: calc(100vh - 140px);
  padding: 30px 20px;
}

.footer {
  background: #2c3e50;
  color: #95a5a6;
  text-align: center;
  padding: 20px 0;
  height: 70px;
  line-height: 30px;
}

.footer p {
  margin: 0;
}
</style>
