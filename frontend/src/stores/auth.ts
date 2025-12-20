import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '@/services/api'

export interface UserInfo {
  id?: number
  username: string
  email?: string
  token: string
  role: 'admin' | 'user' | 'guest' | 'local'
  avatar?: string
  rememberMe?: boolean
}

export const useAuthStore = defineStore('auth', () => {
  const userInfo = ref<UserInfo | null>(null)
  const isAuthenticated = computed(() => !!userInfo.value)
  const isGuest = computed(() => userInfo.value?.role === 'guest')
  const isLocal = computed(() => userInfo.value?.role === 'local')
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

  const login = (info: UserInfo) => {
    userInfo.value = info
    
    // 根据记住我选项选择存储方式
    const storage = info.rememberMe ? localStorage : sessionStorage
    storage.setItem('userInfo', JSON.stringify(info))
    storage.setItem('auth_token', info.token)
  }

  // 后端登录
  const loginWithBackend = async (username: string, password: string, rememberMe: boolean = false) => {
    try {
      const response = await authAPI.login(username, password)
      const userInfo: UserInfo = {
        id: response.user.id,
        username: response.user.username,
        email: response.user.email,
        token: response.token,
        role: response.user.role || 'user',
        avatar: response.user.avatar,
        rememberMe,
      }
      login(userInfo)
      return { success: true, user: userInfo }
    } catch (error: any) {
      console.error('Login failed:', error)
      return { 
        success: false, 
        error: error.response?.data?.error || '登录失败，请检查用户名和密码' 
      }
    }
  }

  // 后端注册
  const registerWithBackend = async (username: string, email: string, password: string, confirmPassword: string) => {
    try {
      const response = await authAPI.register(username, email, password, confirmPassword)
      console.log('Registration response:', response)
      
      // 检查响应是否包含 message 字段
      if (response && (response.message || response.user)) {
        return { success: true, message: response.message || '注册成功' }
      }
      
      return { success: false, error: '注册响应格式错误' }
    } catch (error: any) {
      console.error('Registration failed:', error)
      console.error('Error response:', error.response)
      return { 
        success: false, 
        error: error.response?.data?.error || error.message || '注册失败，请稍后重试' 
      }
    }
  }

  const logout = () => {
    userInfo.value = null
    localStorage.removeItem('userInfo')
    localStorage.removeItem('auth_token')
    sessionStorage.removeItem('userInfo')
    sessionStorage.removeItem('auth_token')
  }

  const loadUserInfo = () => {
    // 优先从 localStorage 加载
    const localInfo = localStorage.getItem('userInfo')
    if (localInfo) {
      userInfo.value = JSON.parse(localInfo)
      return
    }
    
    // 其次从 sessionStorage 加载
    const sessionInfo = sessionStorage.getItem('userInfo')
    if (sessionInfo) {
      userInfo.value = JSON.parse(sessionInfo)
    }
  }

  const updateUserInfo = (updates: Partial<UserInfo>) => {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...updates }
      
      // 更新存储
      const storage = userInfo.value.rememberMe ? localStorage : sessionStorage
      storage.setItem('userInfo', JSON.stringify(userInfo.value))
    }
  }

  const clearStorage = () => {
    localStorage.removeItem('userInfo')
    localStorage.removeItem('auth_token')
    sessionStorage.removeItem('userInfo')
    sessionStorage.removeItem('auth_token')
  }

  // 初始化时加载用户信息
  loadUserInfo()

  return {
    userInfo,
    isAuthenticated,
    isGuest,
    isLocal,
    isAdmin,
    login,
    loginWithBackend,
    registerWithBackend,
    logout,
    loadUserInfo,
    updateUserInfo,
    clearStorage
  }
})