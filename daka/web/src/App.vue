<template>
  <div class="app">
    <header class="header">
      <div class="header-inner">
        <div class="logo">
          <span class="logo-icon">⏰</span>
          <span class="logo-text">自律打卡</span>
        </div>
        <nav class="nav">
          <router-link to="/" class="nav-link" :class="{ active: $route.path === '/' }">
            <span>🏠</span> 打卡中心
          </router-link>
          <router-link to="/stats" class="nav-link" :class="{ active: $route.path === '/stats' }">
            <span>📊</span> 数据统计
          </router-link>
        </nav>
      </div>
    </header>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>

<style scoped>
.app {
  min-height: 100vh;
}

.header {
  background: white;
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: var(--shadow);
}

.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  font-size: 28px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav {
  display: flex;
  gap: 8px;
}

.nav-link {
  padding: 8px 16px;
  border-radius: 10px;
  color: var(--text-secondary);
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.nav-link:hover {
  background: #f1f5f9;
  color: var(--text);
}

.nav-link.active {
  background: linear-gradient(135deg, var(--primary), var(--primary-dark));
  color: white;
}

.main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}
</style>
