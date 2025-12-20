import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login/index.vue'),
    meta: { title: '登录', requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard/index.vue'),
        meta: { title: '控制台' },
      },
      {
        path: '/analysis',
        name: 'CodeAnalysis',
        component: () => import('@/views/CodeAnalysis/index.vue'),
        meta: { title: '代码分析' },
      },
      {
        path: '/conversion',
        name: 'CodeConversion',
        component: () => import('@/views/CodeConversion/index.vue'),
        meta: { title: '代码转换' },
      },
      {
        path: '/test-generation',
        name: 'TestGeneration',
        component: () => import('@/views/TestGeneration/index.vue'),
        meta: { title: '测试生成' },
      },
      {
        path: '/git-analysis',
        name: 'GitAnalysis',
        component: () => import('@/views/GitAnalysis/index.vue'),
        meta: { title: 'Git分析' },
      },
      {
        path: '/settings',
        name: 'Settings',
        component: () => import('@/views/Settings/index.vue'),
        meta: { title: '设置' },
      },
      {
        path: '/tasks',
        name: 'TaskList',
        component: () => import('@/views/TaskList/index.vue'),
        meta: { title: '任务列表' },
      },
    ]
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, _from, next) => {
  document.title = `${to.meta.title || 'CodeSage'} - CodeSage`
  
  const authStore = useAuthStore()
  
  // 检查是否需要认证
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // 保存原始目标路径
    next(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    // 已登录用户访问登录页，重定向到首页
    next('/')
  } else {
    next()
  }
})

export default router