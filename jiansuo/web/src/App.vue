<template>
  <div class="layout">
    <el-container class="main-container">
      <el-header class="app-header">
        <div class="header-content">
          <div class="logo" @click="$router.push('/')">
            <el-icon :size="28" color="#409EFF"><Search /></el-icon>
            <span class="title">智能关键词检索系统</span>
          </div>
          <el-menu
            mode="horizontal"
            :default-active="$route.path"
            router
            class="nav-menu"
            background-color="transparent"
            text-color="#303133"
            active-text-color="#409EFF"
          >
            <el-menu-item index="/search">
              <el-icon><Search /></el-icon>
              <span>智能搜索</span>
            </el-menu-item>
            <el-menu-item index="/admin">
              <el-icon><Management /></el-icon>
              <span>关键词管理</span>
            </el-menu-item>
            <el-menu-item index="/stats">
              <el-icon><DataAnalysis /></el-icon>
              <span>统计分析</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-header>
      <el-main class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  let sessionId = localStorage.getItem('session_id')
  if (!sessionId) {
    sessionId = 'sess_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
    localStorage.setItem('session_id', sessionId)
  }
})
</script>

<style scoped>
.layout {
  min-height: 100vh;
}
.main-container {
  min-height: 100vh;
}
.app-header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 0;
  height: 64px;
  position: sticky;
  top: 0;
  z-index: 100;
}
.header-content {
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 24px;
}
.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  margin-right: 40px;
}
.logo .title {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, #409EFF, #67C23A);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.nav-menu {
  flex: 1;
  border-bottom: none;
}
.app-main {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
