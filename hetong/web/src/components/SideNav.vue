<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const menuItems = [
  { path: '/dashboard', name: '数据大屏', icon: '📊' },
  { path: '/reports', name: '报表中心', icon: '📋' },
  { path: '/trace', name: '数据溯源', icon: '🔍' },
  { path: '/snapshots', name: '快照管理', icon: '📸' },
]

const isActive = (path: string) => route.path === path

const navigate = (path: string) => {
  router.push(path)
}

const currentTime = computed(() => {
  return new Date().toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
})
</script>

<template>
  <aside class="side-nav">
    <div class="nav-header">
      <div class="logo">
        <span class="logo-icon">📈</span>
        <span class="logo-text">DataViz</span>
      </div>
      <div class="system-time">{{ currentTime }}</div>
    </div>

    <nav class="nav-menu">
      <div
        v-for="item in menuItems"
        :key="item.path"
        class="nav-item"
        :class="{ active: isActive(item.path) }"
        @click="navigate(item.path)"
      >
        <span class="nav-icon">{{ item.icon }}</span>
        <span class="nav-label">{{ item.name }}</span>
        <div class="nav-indicator" v-if="isActive(item.path)"></div>
      </div>
    </nav>

    <div class="nav-footer">
      <div class="status-indicator">
        <span class="status-dot"></span>
        <span class="status-text">系统正常</span>
      </div>
      <div class="version">v1.0.0</div>
    </div>
  </aside>
</template>

<style scoped>
.side-nav {
  width: 240px;
  height: 100vh;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.98), rgba(30, 41, 59, 0.95));
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(10px);
}

.nav-header {
  padding: 24px 20px;
  border-bottom: 1px solid var(--border-color);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.logo-icon {
  font-size: 28px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--neon-blue), var(--neon-cyan));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-family: 'JetBrains Mono', monospace;
}

.system-time {
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}

.nav-menu {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  color: var(--text-secondary);
}

.nav-item:hover {
  background: rgba(59, 130, 246, 0.1);
  color: var(--text-primary);
}

.nav-item.active {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(6, 182, 212, 0.2));
  color: var(--text-primary);
}

.nav-icon {
  font-size: 18px;
}

.nav-label {
  font-size: 14px;
  font-weight: 500;
}

.nav-indicator {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: linear-gradient(180deg, var(--neon-blue), var(--neon-cyan));
  border-radius: 2px 0 0 2px;
  box-shadow: 0 0 10px var(--neon-blue);
}

.nav-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  background: var(--neon-green);
  border-radius: 50%;
  animation: pulse 2s infinite;
  box-shadow: 0 0 10px var(--neon-green);
}

.status-text {
  font-size: 12px;
  color: var(--neon-green);
  font-family: 'JetBrains Mono', monospace;
}

.version {
  font-size: 11px;
  color: var(--text-muted);
  font-family: 'JetBrains Mono', monospace;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
