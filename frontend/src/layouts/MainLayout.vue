<template>
  <el-container class="app-container">
    <!-- 背景装饰效果 -->
    <div class="background-decoration">
      <div class="bg-particles">
        <div class="bg-particle" v-for="i in 15" :key="i"></div>
      </div>
      <div class="bg-glow bg-glow-1"></div>
      <div class="bg-glow bg-glow-2"></div>
      <div class="bg-mesh"></div>
    </div>
    
    <!-- 侧边栏 -->
    <el-aside width="260px" class="sidebar">
      <div class="logo-container">
        <div class="logo-icon-wrapper">
          <el-icon size="36" class="logo-icon-main">
            <Cpu />
          </el-icon>
        </div>
        <h1 class="logo">
          <span class="neon-text">CodeSage</span>
        </h1>
        <p class="tagline">智能代码重构助手</p>
      </div>
      
      <el-menu
        :default-active="$route.path"
        class="nav-menu"
        router
        :collapse="false"
        background-color="transparent"
        text-color="var(--text-secondary)"
        active-text-color="var(--text-primary)"
      >
        <el-menu-item index="/">
          <el-icon><House /></el-icon>
          <span>控制台</span>
        </el-menu-item>
        
        <el-menu-item index="/analysis">
          <el-icon><Search /></el-icon>
          <span>代码分析</span>
        </el-menu-item>
        
        <el-menu-item index="/conversion">
          <el-icon><Refresh /></el-icon>
          <span>代码转换</span>
        </el-menu-item>
        
        <el-menu-item index="/test-generation">
          <el-icon><DocumentChecked /></el-icon>
          <span>测试生成</span>
        </el-menu-item>
        
        <el-menu-item index="/git-analysis">
          <el-icon><Coin /></el-icon>
          <span>Git分析</span>
        </el-menu-item>
        
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>设置</span>
        </el-menu-item>
      </el-menu>
      
      <!-- 侧边栏底部装饰 -->
      <div class="sidebar-footer">
        <div class="version-info">v1.0.0</div>
      </div>
    </el-aside>
    
    <!-- 主内容区 -->
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h2 class="page-title">{{ $route.meta.title }}</h2>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" @click="handleQuickAction">
              快速开始
            </el-button>
            
            <!-- 用户菜单 -->
            <el-dropdown trigger="click" @command="handleUserCommand">
              <el-button text class="user-button">
                <el-avatar 
                  v-if="(authStore.userInfo as any)?.avatar" 
                  :size="32" 
                  :src="(authStore.userInfo as any)?.avatar"
                  class="user-avatar"
                >
                  <el-icon><User /></el-icon>
                </el-avatar>
                <div v-else class="user-avatar-default">
                  <el-icon class="user-icon-default"><User /></el-icon>
                </div>
                <span class="username">{{ (authStore.userInfo as any)?.username ?? '用户' }}</span>
                <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile" :icon="User">
                    个人资料
                  </el-dropdown-item>
                  <el-dropdown-item command="settings" :icon="Setting">
                    账号设置
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout" :icon="SwitchButton">
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
// import { ref } from 'vue'
import { House, Search, Refresh, DocumentChecked, Coin, Setting, Plus, User, SwitchButton, ArrowDown, Cpu } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
// const showUserMenu = ref(false)

const handleQuickAction = () => {
  ElMessage.success('请输入代码进行分析')
  router.push('/analysis')
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要退出登录吗？',
      '退出确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch {
    // 用户取消退出
  }
}

const handleUserCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/settings?tab=profile')
      break
    case 'settings':
      router.push('/settings?tab=general')
      break
    case 'logout':
      handleLogout()
      break
  }
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  background: radial-gradient(ellipse at top, #0f0f23 0%, #050510 100%);
  position: relative;
  overflow: hidden;
}

/* 背景装饰 */
.background-decoration {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
}

/* 背景粒子 */
.bg-particles {
  position: absolute;
  width: 100%;
  height: 100%;
}

.bg-particle {
  position: absolute;
  width: 2px;
  height: 2px;
  background: radial-gradient(circle, rgba(66, 133, 244, 0.6) 0%, transparent 70%);
  border-radius: 50%;
  animation: float-bg-particle 20s infinite ease-in-out;
}

.bg-particle:nth-child(odd) {
  animation-duration: 25s;
  background: radial-gradient(circle, rgba(156, 39, 176, 0.4) 0%, transparent 70%);
}

.bg-particle:nth-child(1) { top: 10%; left: 15%; animation-delay: 0s; }
.bg-particle:nth-child(2) { top: 20%; left: 35%; animation-delay: 2s; }
.bg-particle:nth-child(3) { top: 30%; left: 55%; animation-delay: 4s; }
.bg-particle:nth-child(4) { top: 40%; left: 75%; animation-delay: 6s; }
.bg-particle:nth-child(5) { top: 50%; left: 25%; animation-delay: 8s; }
.bg-particle:nth-child(6) { top: 60%; left: 45%; animation-delay: 10s; }
.bg-particle:nth-child(7) { top: 70%; left: 65%; animation-delay: 12s; }
.bg-particle:nth-child(8) { top: 80%; left: 85%; animation-delay: 14s; }
.bg-particle:nth-child(9) { top: 15%; left: 60%; animation-delay: 1s; }
.bg-particle:nth-child(10) { top: 25%; left: 80%; animation-delay: 3s; }
.bg-particle:nth-child(11) { top: 35%; left: 20%; animation-delay: 5s; }
.bg-particle:nth-child(12) { top: 45%; left: 40%; animation-delay: 7s; }
.bg-particle:nth-child(13) { top: 55%; left: 70%; animation-delay: 9s; }
.bg-particle:nth-child(14) { top: 65%; left: 30%; animation-delay: 11s; }
.bg-particle:nth-child(15) { top: 75%; left: 50%; animation-delay: 13s; }

@keyframes float-bg-particle {
  0%, 100% {
    transform: translateY(0) scale(1);
    opacity: 0.3;
  }
  50% {
    transform: translateY(-50px) scale(1.5);
    opacity: 0.6;
  }
}

/* 背景光晕 */
.bg-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.1;
  animation: glow-pulse 10s infinite ease-in-out;
}

.bg-glow-1 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, #4285f4 0%, transparent 70%);
  top: -150px;
  left: -150px;
}

.bg-glow-2 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, #9c27b0 0%, transparent 70%);
  bottom: -100px;
  right: -100px;
  animation-delay: 5s;
}

@keyframes glow-pulse {
  0%, 100% {
    opacity: 0.1;
    transform: scale(1);
  }
  50% {
    opacity: 0.2;
    transform: scale(1.1);
  }
}

/* 网格背景 */
.bg-mesh {
  position: absolute;
  width: 100%;
  height: 100%;
  background-image: 
    linear-gradient(rgba(66, 133, 244, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(66, 133, 244, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  opacity: 0.5;
}

/* 侧边栏 */
.sidebar {
  background: rgba(15, 15, 35, 0.7);
  backdrop-filter: blur(20px) saturate(150%);
  border-right: 1px solid rgba(66, 133, 244, 0.2);
  padding: 24px 0;
  position: relative;
  z-index: 10;
  box-shadow: 
    4px 0 20px rgba(0, 0, 0, 0.3),
    inset -1px 0 0 rgba(66, 133, 244, 0.1);
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  right: -1px;
  width: 1px;
  height: 100%;
  background: linear-gradient(
    180deg,
    transparent 0%,
    rgba(66, 133, 244, 0.5) 50%,
    transparent 100%
  );
  animation: sidebar-glow 3s infinite ease-in-out;
}

@keyframes sidebar-glow {
  0%, 100% {
    opacity: 0.3;
  }
  50% {
    opacity: 0.8;
  }
}

.logo-container {
  text-align: center;
  padding: 0 24px 28px;
  border-bottom: 1px solid rgba(66, 133, 244, 0.15);
  margin-bottom: 24px;
  position: relative;
}

.logo-icon-wrapper {
  margin-bottom: 12px;
  animation: logo-pulse 3s infinite ease-in-out;
}

.logo-icon-main {
  font-size: 36px;
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  filter: drop-shadow(0 0 20px rgba(66, 133, 244, 0.4));
}

@keyframes logo-pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.logo {
  font-size: 36px;
  font-weight: 900;
  margin: 0 0 8px 0;
  letter-spacing: -1.5px;
  line-height: 1;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Helvetica Neue', Arial, sans-serif;
  position: relative;
}

.neon-text {
  display: inline-block;
  background: linear-gradient(
    135deg, 
    #64b5f6 0%,
    #4285f4 25%,
    #9c27b0 50%,
    #ba68c8 75%,
    #00bcd4 100%
  );
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  background-size: 200% 200%;
  animation: gradient-flow 5s ease infinite;
  filter: brightness(1.2) contrast(1.1);
  text-shadow: none;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.neon-text::before {
  content: 'CodeSage';
  position: absolute;
  top: 0;
  left: 0;
  background: linear-gradient(
    135deg,
    rgba(100, 181, 246, 0.3) 0%,
    rgba(66, 133, 244, 0.3) 25%,
    rgba(156, 39, 176, 0.3) 50%,
    rgba(186, 104, 200, 0.3) 75%,
    rgba(0, 188, 212, 0.3) 100%
  );
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  filter: blur(8px);
  opacity: 0.8;
  z-index: -1;
  animation: gradient-flow 5s ease infinite;
}

@keyframes gradient-flow {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.tagline {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin: 0;
  font-weight: 300;
}

/* 导航菜单 */
.nav-menu {
  border: none;
  background: transparent !important;
  padding: 0 12px;
}

.nav-menu .el-menu-item {
  margin: 6px 0;
  border-radius: 12px;
  height: 50px;
  line-height: 50px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  font-weight: 500;
}

.nav-menu .el-menu-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(66, 133, 244, 0.2),
    transparent
  );
  transition: left 0.5s ease;
}

.nav-menu .el-menu-item:hover::before {
  left: 100%;
}

.nav-menu .el-menu-item:hover {
  background: rgba(66, 133, 244, 0.1) !important;
  transform: translateX(6px);
  box-shadow: 0 4px 12px rgba(66, 133, 244, 0.2);
}

.nav-menu .el-menu-item.is-active {
  background: linear-gradient(135deg, rgba(66, 133, 244, 0.3), rgba(156, 39, 176, 0.2)) !important;
  box-shadow: 
    0 4px 15px rgba(66, 133, 244, 0.3),
    inset 0 0 20px rgba(66, 133, 244, 0.1);
  border-left: 3px solid #4285f4;
}

.nav-menu .el-menu-item.is-active::after {
  content: '';
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 6px;
  height: 6px;
  background: #4285f4;
  border-radius: 50%;
  box-shadow: 0 0 10px #4285f4;
}

/* 侧边栏底部 */
.sidebar-footer {
  position: absolute;
  bottom: 20px;
  left: 0;
  width: 100%;
  text-align: center;
  padding: 0 24px;
}

.version-info {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

/* 顶栏 */
.header {
  background: rgba(15, 15, 35, 0.7);
  backdrop-filter: blur(20px) saturate(150%);
  border-bottom: 1px solid rgba(66, 133, 244, 0.2);
  padding: 0 32px;
  display: flex;
  align-items: center;
  position: relative;
  z-index: 10;
  box-shadow: 
    0 4px 20px rgba(0, 0, 0, 0.2),
    inset 0 -1px 0 rgba(66, 133, 244, 0.1);
}

.header::before {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  width: 100%;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(66, 133, 244, 0.6) 50%,
    transparent 100%
  );
  animation: header-glow 3s infinite ease-in-out;
}

@keyframes header-glow {
  0%, 100% {
    opacity: 0.3;
  }
  50% {
    opacity: 0.8;
  }
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #ffffff 0%, rgba(255, 255, 255, 0.7) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  letter-spacing: -0.5px;
}

.header-actions {
  display: flex;
  gap: 16px;
  align-items: center;
}

/* 快速开始按钮 */
.header-actions .el-button--primary {
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 100%);
  border: none;
  border-radius: 10px;
  padding: 10px 24px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(66, 133, 244, 0.3);
  transition: all 0.3s ease;
}

.header-actions .el-button--primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(66, 133, 244, 0.5);
}

/* 用户按钮 */
.user-button {
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid rgba(66, 133, 244, 0.3) !important;
  border-radius: 12px;
  padding: 8px 16px;
  color: rgba(255, 255, 255, 0.9) !important;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-button:hover {
  background: rgba(66, 133, 244, 0.15) !important;
  border-color: rgba(66, 133, 244, 0.5) !important;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(66, 133, 244, 0.2);
}

.user-button :deep(.el-button__text) {
  color: rgba(255, 255, 255, 0.9) !important;
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-avatar {
  border: 2px solid rgba(66, 133, 244, 0.4);
  transition: all 0.3s ease;
}

.user-button:hover .user-avatar {
  border-color: rgba(66, 133, 244, 0.8);
  box-shadow: 0 0 12px rgba(66, 133, 244, 0.4);
}

.user-avatar-default {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid rgba(66, 133, 244, 0.4);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.user-avatar-default::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, transparent 0%, rgba(255, 255, 255, 0.2) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.user-button:hover .user-avatar-default {
  border-color: rgba(66, 133, 244, 0.8);
  box-shadow: 0 0 12px rgba(66, 133, 244, 0.4);
  transform: scale(1.05);
}

.user-button:hover .user-avatar-default::before {
  opacity: 1;
}

.user-icon-default {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.95);
  position: relative;
  z-index: 1;
}

.username {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}

.dropdown-icon {
  font-size: 12px;
  transition: transform 0.3s ease;
}

.el-dropdown:hover .dropdown-icon {
  transform: rotate(180deg);
}

/* 主内容区 */
.main-content {
  background: transparent;
  padding: 32px;
  overflow-y: auto;
  position: relative;
  z-index: 5;
}

/* 滚动条样式 */
.main-content::-webkit-scrollbar {
  width: 8px;
}

.main-content::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 4px;
}

.main-content::-webkit-scrollbar-thumb {
  background: rgba(66, 133, 244, 0.3);
  border-radius: 4px;
  transition: background 0.3s ease;
}

.main-content::-webkit-scrollbar-thumb:hover {
  background: rgba(66, 133, 244, 0.5);
}

/* 过渡动画增强 */
.fade-enter-active,
.fade-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.98);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.98);
}

/* 下拉菜单样式 */
.el-dropdown-menu {
  background: rgba(15, 15, 35, 0.95) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(66, 133, 244, 0.2) !important;
  border-radius: 12px;
  padding: 8px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5) !important;
}

.el-dropdown-menu__item {
  color: rgba(255, 255, 255, 0.8) !important;
  border-radius: 8px;
  padding: 10px 16px;
  transition: all 0.3s ease;
}

.el-dropdown-menu__item:hover {
  background: rgba(66, 133, 244, 0.15) !important;
  color: rgba(255, 255, 255, 1) !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 220px !important;
  }
  
  .page-title {
    font-size: 22px;
  }
  
  .header-actions {
    gap: 12px;
  }
  
  .bg-particle {
    display: none;
  }
}
</style>