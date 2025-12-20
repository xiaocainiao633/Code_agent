<template>
  <div id="app">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// 初始化用户认证状态
onMounted(() => {
  authStore.loadUserInfo()
  
  // 检查当前路由是否需要认证
  const currentRoute = router.currentRoute.value
  if (currentRoute.meta.requiresAuth && !authStore.isAuthenticated) {
    router.push('/login')
  }
})
</script>

<style>
#app {
  height: 100vh;
  margin: 0;
  padding: 0;
}

/* 全局过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 全局样式重置 */
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}
</style>

<style scoped>
.app-container {
  height: 100vh;
  background: var(--bg-primary);
}

.sidebar {
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  padding: 20px 0;
}

.logo-container {
  text-align: center;
  padding: 0 20px 30px;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 20px;
}

.logo {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 8px;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.tagline {
  font-size: 12px;
  color: var(--text-tertiary);
  margin: 0;
}

.nav-menu {
  border: none;
  background: transparent !important;
}

.nav-menu .el-menu-item {
  margin: 4px 12px;
  border-radius: 8px;
  height: 48px;
  line-height: 48px;
  transition: all 0.3s ease;
}

.nav-menu .el-menu-item:hover {
  background: var(--bg-hover) !important;
  transform: translateX(4px);
}

.nav-menu .el-menu-item.is-active {
  background: var(--primary-color) !important;
  box-shadow: var(--shadow-sm);
}

.header {
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  padding: 0 24px;
  display: flex;
  align-items: center;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.user-button {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 8px 12px;
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.user-button:hover {
  background: rgba(255, 255, 255, 0.15) !important;
  border-color: var(--primary-color) !important;
  color: var(--text-primary) !important;
}

.user-button :deep(.el-button__text) {
  color: var(--text-primary) !important;
}

.username {
  margin: 0 4px;
  font-size: 14px;
}

.dropdown-icon {
  font-size: 12px;
  transition: transform 0.3s ease;
}

.el-dropdown:hover .dropdown-icon {
  transform: rotate(180deg);
}

.main-content {
  background: var(--bg-primary);
  padding: 24px;
  overflow-y: auto;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 200px !important;
  }
  
  .page-title {
    font-size: 20px;
  }
}
</style>
