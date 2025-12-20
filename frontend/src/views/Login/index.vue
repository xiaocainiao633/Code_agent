<template>
  <div class="login-container">
    <!-- 世界地图点阵背景动画 -->
    <div class="world-map-background">
      <canvas ref="mapCanvas" class="map-canvas"></canvas>
      <div class="gradient-overlay"></div>
      
      <!-- 浮动粒子 -->
      <div class="particles-container">
        <div class="particle" v-for="i in 30" :key="i"></div>
      </div>
      
      <!-- 光晕效果 -->
      <div class="glow-effects">
        <div class="glow glow-1"></div>
        <div class="glow glow-2"></div>
        <div class="glow glow-3"></div>
      </div>
    </div>
    
    <!-- 主标题和Get Started按钮 -->
    <div class="hero-section" v-if="!showLoginCard">
      <div class="hero-content">
        <div class="hero-logo">
          <el-icon size="80" class="logo-icon-large">
            <Cpu />
          </el-icon>
        </div>
        <h1 class="hero-title">
          <span class="gradient-text">CodeSage</span>
        </h1>
        <p class="hero-subtitle">
          Deep Research <span class="highlight">at Your Fingertips</span>
        </p>
        <p class="hero-description">
          Meet CodeSage, your personal Deep Research assistant. With powerful tools like<br>
          code analysis, refactoring, and AI-powered insights.
        </p>
        <div class="hero-actions">
          <button class="get-started-btn" @click="showLoginCard = true">
            <span class="btn-text">Get Started</span>
            <el-icon class="btn-icon"><Right /></el-icon>
          </button>
          <button class="learn-more-btn" @click="handleLearnMore">
            <el-icon class="btn-icon"><Search /></el-icon>
            <span class="btn-text">Learn More</span>
          </button>
        </div>
      </div>
    </div>
    
    <!-- 登录卡片 -->
    <transition name="card-slide" mode="out-in">
      <div class="login-card-wrapper" v-if="showLoginCard" key="login">
        <div class="login-card">
      <div class="login-header">
        <div class="logo-section">
          <div class="logo-icon">
            <el-icon size="48" color="var(--primary-color)">
              <Cpu />
            </el-icon>
          </div>
          <h1 class="logo-text">
            <span class="neon-text">CodeSage</span>
          </h1>
          <p class="logo-subtitle">智能代码重构助手</p>
        </div>
        
        <!-- 标签页切换 -->
        <div class="auth-tabs">
          <button 
            class="auth-tab" 
            :class="{ active: activeTab === 'login' }"
            @click="activeTab = 'login'"
          >
            <span class="tab-text">登录</span>
          </button>
          <button 
            class="auth-tab" 
            :class="{ active: activeTab === 'register' }"
            @click="activeTab = 'register'"
          >
            <span class="tab-text">注册</span>
          </button>
          <div class="tab-indicator" :class="{ 'to-register': activeTab === 'register' }"></div>
        </div>
      </div>
      
      <!-- 登录表单 -->
      <el-form
        v-if="activeTab === 'login'"
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
            class="login-input"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Key"
            class="login-input"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          >
            <span class="button-text">登录</span>
            <el-icon class="button-icon" v-if="!loading">
              <Right />
            </el-icon>
          </el-button>
        </el-form-item>
        
        <div class="form-footer">
          <el-checkbox v-model="rememberMe" class="remember-checkbox">
            记住我
          </el-checkbox>
          <el-button text type="primary" class="forgot-password-btn" @click="handleForgotPassword">
            忘记密码？
          </el-button>
        </div>
      </el-form>
      
      <!-- 注册表单 -->
      <el-form
        v-if="activeTab === 'register'"
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        class="login-form register-form"
        @submit.prevent="handleRegister"
      >
        <el-form-item prop="username">
          <el-input
            v-model="registerForm.username"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
            class="login-input"
          />
        </el-form-item>
        
        <el-form-item prop="email">
          <el-input
            v-model="registerForm.email"
            placeholder="邮箱"
            size="large"
            :prefix-icon="Search"
            class="login-input"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Key"
            class="login-input"
            show-password
          />
        </el-form-item>
        
        <el-form-item prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="确认密码"
            size="large"
            :prefix-icon="Key"
            class="login-input"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleRegister"
          >
            <span class="button-text">注册</span>
            <el-icon class="button-icon" v-if="!loading">
              <Right />
            </el-icon>
          </el-button>
        </el-form-item>
      </el-form>
      
      <!-- 社交登录和快速访问（只在登录标签页显示） -->
      <template v-if="activeTab === 'login'">
      <div class="divider">
        <span class="divider-text">或使用社交账号登录</span>
      </div>
      
      <!-- 社交登录 -->
      <div class="social-login">
        <button class="social-btn github-btn" @click="handleGithubLogin">
          <svg class="social-icon" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
          <span>GitHub 登录</span>
        </button>
      </div>
      
      <div class="divider">
        <span class="divider-text">或快速访问</span>
      </div>
      
      <div class="alternative-login">
        <el-button
          class="alt-button"
          :icon="Key"
          @click="handleGuestLogin"
        >
          游客访问
        </el-button>
        <el-button
          class="alt-button"
          :icon="Setting"
          @click="handleLocalMode"
        >
          本地模式
        </el-button>
      </div>
      </template>
      
      <!-- 返回按钮 -->
      <div class="back-to-home">
        <el-button text type="info" @click="showLoginCard = false" class="back-btn">
          <el-icon><Delete /></el-icon>
          <span>返回首页</span>
        </el-button>
      </div>
      </div>
      </div>
    </transition>
    
    <!-- 装饰线条 -->
    <div class="decoration-lines">
      <div class="line line-1"></div>
      <div class="line line-2"></div>
      <div class="line line-3"></div>
    </div>

    <!-- 忘记密码对话框 -->
    <el-dialog
      v-model="showForgotPasswordDialog"
      title="重置密码"
      width="400px"
      :close-on-click-modal="false"
    >
      <!-- 步骤 1: 输入邮箱 -->
      <el-form v-if="forgotPasswordStep === 1">
        <el-form-item label="邮箱">
          <el-input
            v-model="forgotPasswordForm.email"
            placeholder="请输入注册邮箱"
            type="email"
          />
        </el-form-item>
      </el-form>

      <!-- 步骤 2: 输入验证码 -->
      <el-form v-if="forgotPasswordStep === 2">
        <el-form-item label="验证码">
          <el-input
            v-model="forgotPasswordForm.code"
            placeholder="请输入6位验证码"
            maxlength="6"
          />
        </el-form-item>
      </el-form>

      <!-- 步骤 3: 设置新密码 -->
      <el-form v-if="forgotPasswordStep === 3">
        <el-form-item label="新密码">
          <el-input
            v-model="forgotPasswordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input
            v-model="forgotPasswordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span v-if="forgotPasswordStep === 1">
          <el-button @click="showForgotPasswordDialog = false">取消</el-button>
          <el-button type="primary" @click="handleSendResetCode" :loading="loading">
            发送验证码
          </el-button>
        </span>
        <span v-if="forgotPasswordStep === 2">
          <el-button @click="forgotPasswordStep = 1">上一步</el-button>
          <el-button type="primary" @click="handleVerifyResetCode" :loading="loading">
            验证
          </el-button>
        </span>
        <span v-if="forgotPasswordStep === 3">
          <el-button @click="forgotPasswordStep = 2">上一步</el-button>
          <el-button type="primary" @click="handleResetPassword" :loading="loading">
            重置密码
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- GitHub 登录对话框 -->
    <el-dialog
      v-model="showGithubDialog"
      title="GitHub 登录"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-alert
        title="提示"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      >
        请输入您已绑定的 GitHub 账号信息。如果还未绑定，请先注册账号后在设置中绑定 GitHub。
      </el-alert>
      <el-form>
        <el-form-item label="GitHub ID">
          <el-input
            v-model="githubForm.githubId"
            placeholder="请输入您的 GitHub ID"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showGithubDialog = false">取消</el-button>
        <el-button type="primary" @click="handleConfirmGithubLogin" :loading="loading">
          登录
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Key, Right, Cpu, Setting, Search, Delete } from '@element-plus/icons-vue'
import { useAuthStore } from '../../stores/auth'
import { authAPI } from '../../services/api'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const loginFormRef = ref<FormInstance>()
const registerFormRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)
const mapCanvas = ref<HTMLCanvasElement>()
const showLoginCard = ref(false)
const activeTab = ref<'login' | 'register'>('login')

const loginForm = reactive({
  username: '',
  password: ''
})

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

// 世界地图点阵动画
let animationFrameId: number
let ctx: CanvasRenderingContext2D | null = null
const dots: Array<{ x: number; y: number; baseY: number; speed: number; size: number; opacity: number }> = []

const initWorldMapAnimation = () => {
  if (!mapCanvas.value) return
  
  const canvas = mapCanvas.value
  const resizeCanvas = () => {
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
  }
  resizeCanvas()
  window.addEventListener('resize', resizeCanvas)
  
  ctx = canvas.getContext('2d')
  if (!ctx) return
  
  // 创建世界地图点阵（类似世界地图形状）
  const createDots = () => {
    dots.length = 0
    const cols = Math.floor(canvas.width / 15)
    const rows = Math.floor(canvas.height / 15)
    
    for (let i = 0; i < cols; i++) {
      for (let j = 0; j < rows; j++) {
        // 使用数学函数模拟大陆形状
        const x = i * 15 + 7.5
        const y = j * 15 + 7.5
        
        // 模拟地图分布：中间区域点较多，边缘较少
        const centerX = canvas.width / 2
        const centerY = canvas.height / 2
        const distanceFromCenter = Math.sqrt(
          Math.pow((x - centerX) / canvas.width, 2) + 
          Math.pow((y - centerY) / canvas.height, 2)
        )
        
        // 使用波浪函数创建大陆轮廓效果
        const wave1 = Math.sin(x / 80) * Math.cos(y / 80)
        const wave2 = Math.sin(x / 120 + Math.PI / 4) * Math.cos(y / 120)
        const isLand = (wave1 + wave2) > -0.3 && distanceFromCenter < 0.7
        
        if (isLand && Math.random() > 0.3) {
          dots.push({
            x,
            y,
            baseY: y,
            speed: 0.2 + Math.random() * 0.5,
            size: 1 + Math.random() * 2,
            opacity: 0.3 + Math.random() * 0.7
          })
        }
      }
    }
  }
  
  createDots()
  
  // 动画循环
  let time = 0
  const animate = () => {
    if (!ctx || !canvas) return
    
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    time += 0.01
    
    // 绘制点阵
    dots.forEach(dot => {
      if (!ctx) return
      
      // 波动效果
      const wave = Math.sin(time + dot.x / 100) * 3
      const currentY = dot.baseY + wave
      
      // 颜色渐变：从蓝色到紫色
      const hue = 200 + (dot.x / canvas.width) * 80
      const brightness = 50 + Math.sin(time + dot.x / 50) * 20
      
      ctx.fillStyle = `hsla(${hue}, 80%, ${brightness}%, ${dot.opacity})`
      ctx.beginPath()
      ctx.arc(dot.x, currentY, dot.size, 0, Math.PI * 2)
      ctx.fill()
      
      // 添加微光效果
      if (Math.random() > 0.98) {
        ctx.shadowBlur = 10
        ctx.shadowColor = `hsla(${hue}, 100%, 70%, 0.8)`
        ctx.fill()
        ctx.shadowBlur = 0
      }
    })
    
    // 连接附近的点
    dots.forEach((dot1, i) => {
      if (!ctx) return
      
      dots.slice(i + 1, i + 20).forEach(dot2 => {
        if (!ctx) return
        
        const dx = dot1.x - dot2.x
        const dy = (dot1.baseY + Math.sin(time + dot1.x / 100) * 3) - 
                   (dot2.baseY + Math.sin(time + dot2.x / 100) * 3)
        const distance = Math.sqrt(dx * dx + dy * dy)
        
        if (distance < 80) {
          const opacity = (1 - distance / 80) * 0.2
          const hue = 200 + ((dot1.x + dot2.x) / 2 / canvas.width) * 80
          
          ctx.strokeStyle = `hsla(${hue}, 70%, 60%, ${opacity})`
          ctx.lineWidth = 0.5
          ctx.beginPath()
          ctx.moveTo(dot1.x, dot1.baseY + Math.sin(time + dot1.x / 100) * 3)
          ctx.lineTo(dot2.x, dot2.baseY + Math.sin(time + dot2.x / 100) * 3)
          ctx.stroke()
        }
      })
    })
    
    animationFrameId = requestAnimationFrame(animate)
  }
  
  animate()
}

// 组件加载时检查是否已登录
onMounted(() => {
  // 如果已经登录，直接跳转到功能页面
  if (authStore.isAuthenticated) {
    const redirect = route.query.redirect as string || '/'
    router.push(redirect)
  }
  
  // 初始化背景动画
  setTimeout(initWorldMapAnimation, 100)
})

onUnmounted(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
})

const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

const registerRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule: any, value: any, callback: any) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 登录处理函数
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    // 验证表单
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 调用后端登录API
    const result = await authStore.loginWithBackend(
      loginForm.username,
      loginForm.password,
      rememberMe.value
    )
    
    loading.value = false
    
    if (result.success) {
      ElMessage.success('登录成功！')
      
      // 获取重定向路径，默认跳转到仪表板
      const redirect = route.query.redirect as string || '/'
      router.push(redirect)
    } else {
      ElMessage.error(result.error || '登录失败')
    }
  } catch (error) {
    loading.value = false
    // 表单验证失败
    ElMessage.warning('请正确填写表单')
  }
}

// 游客登录
const handleGuestLogin = () => {
  ElMessage.info('游客模式：部分功能可能受限')
  
  const guestInfo = {
    username: '游客',
    token: 'guest-token-' + Date.now(),
    role: 'guest' as const
  }
  
  authStore.login(guestInfo)
  
  const redirect = route.query.redirect as string || '/'
  router.push(redirect)
}

// 本地模式
const handleLocalMode = () => {
  ElMessage.success('本地模式：所有数据在本地处理')
  
  const localInfo = {
    username: '本地用户',
    token: 'local-token-' + Date.now(),
    role: 'local' as const
  }
  
  authStore.login(localInfo)
  router.push('/')
}

// 忘记密码对话框
const showForgotPasswordDialog = ref(false)
const forgotPasswordStep = ref(1) // 1: 输入邮箱, 2: 输入验证码, 3: 设置新密码
const forgotPasswordForm = reactive({
  email: '',
  code: '',
  newPassword: '',
  confirmPassword: ''
})

// 忘记密码
const handleForgotPassword = () => {
  showForgotPasswordDialog.value = true
  forgotPasswordStep.value = 1
  forgotPasswordForm.email = ''
  forgotPasswordForm.code = ''
  forgotPasswordForm.newPassword = ''
  forgotPasswordForm.confirmPassword = ''
}

// 发送重置码
const handleSendResetCode = async () => {
  if (!forgotPasswordForm.email) {
    ElMessage.warning('请输入邮箱')
    return
  }

  loading.value = true
  try {
    const response = await authAPI.forgotPassword(forgotPasswordForm.email)
    ElMessage.success('重置码已发送到您的邮箱')
    console.log('重置码（演示）:', response.code) // 演示模式
    forgotPasswordStep.value = 2
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '发送失败')
  } finally {
    loading.value = false
  }
}

// 验证重置码
const handleVerifyResetCode = async () => {
  if (!forgotPasswordForm.code) {
    ElMessage.warning('请输入验证码')
    return
  }

  loading.value = true
  try {
    await authAPI.verifyResetCode(forgotPasswordForm.email, forgotPasswordForm.code)
    ElMessage.success('验证码正确')
    forgotPasswordStep.value = 3
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '验证码错误')
  } finally {
    loading.value = false
  }
}

// 重置密码
const handleResetPassword = async () => {
  if (!forgotPasswordForm.newPassword || !forgotPasswordForm.confirmPassword) {
    ElMessage.warning('请输入新密码')
    return
  }

  if (forgotPasswordForm.newPassword !== forgotPasswordForm.confirmPassword) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }

  loading.value = true
  try {
    await authAPI.resetPassword(
      forgotPasswordForm.email,
      forgotPasswordForm.code,
      forgotPasswordForm.newPassword
    )
    ElMessage.success('密码重置成功！请登录')
    showForgotPasswordDialog.value = false
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '重置失败')
  } finally {
    loading.value = false
  }
}

// GitHub 登录对话框
const showGithubDialog = ref(false)
const githubForm = reactive({
  githubId: '',
  githubUsername: ''
})

// GitHub 登录
const handleGithubLogin = () => {
  showGithubDialog.value = true
  githubForm.githubId = ''
  githubForm.githubUsername = ''
  ElMessage.info('请输入您的 GitHub 信息进行登录')
}

// 确认 GitHub 登录
const handleConfirmGithubLogin = async () => {
  if (!githubForm.githubId) {
    ElMessage.warning('请输入 GitHub ID')
    return
  }

  loading.value = true
  try {
    const response = await authAPI.githubLogin(githubForm.githubId)
    
    const userInfo = {
      id: response.user.id,
      username: response.user.username,
      email: response.user.email,
      token: response.token,
      role: response.user.role || 'user',
      avatar: response.user.avatar,
      rememberMe: false,
    }
    
    authStore.login(userInfo)
    ElMessage.success('GitHub 登录成功！')
    showGithubDialog.value = false
    
    const redirect = route.query.redirect as string || '/'
    router.push(redirect)
  } catch (error: any) {
    if (error.response?.data?.error?.includes('未绑定')) {
      ElMessage.error('该 GitHub 账号未绑定，请先注册并绑定')
    } else {
      ElMessage.error(error.response?.data?.error || 'GitHub 登录失败')
    }
  } finally {
    loading.value = false
  }
}

// Learn More 按钮
const handleLearnMore = () => {
  ElMessage.info('了解更多功能...')
  // 可以跳转到介绍页面或显示功能说明
}

// 注册处理函数
const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  try {
    await registerFormRef.value.validate()
    
    loading.value = true
    
    console.log('开始注册，用户名:', registerForm.username)
    
    // 调用后端注册API
    const result = await authStore.registerWithBackend(
      registerForm.username,
      registerForm.email,
      registerForm.password,
      registerForm.confirmPassword
    )
    
    console.log('注册结果:', result)
    
    loading.value = false
    
    if (result.success) {
      ElMessage.success('注册成功！请登录')
      // 切换到登录标签页
      activeTab.value = 'login'
      // 清空注册表单
      registerForm.username = ''
      registerForm.email = ''
      registerForm.password = ''
      registerForm.confirmPassword = ''
    } else {
      console.error('注册失败:', result.error)
      ElMessage.error(result.error || '注册失败')
    }
  } catch (error) {
    loading.value = false
    console.error('注册异常:', error)
    ElMessage.warning('请正确填写表单')
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(ellipse at center, #0f0f23 0%, #050510 100%);
  position: relative;
  overflow: hidden;
}

/* 世界地图点阵背景 */
.world-map-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.map-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

/* 渐变遮罩 */
.gradient-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(
    ellipse at 30% 50%,
    rgba(66, 133, 244, 0.05) 0%,
    rgba(15, 15, 35, 0) 50%
  ),
  radial-gradient(
    ellipse at 70% 50%,
    rgba(156, 39, 176, 0.05) 0%,
    rgba(15, 15, 35, 0) 50%
  );
  pointer-events: none;
}

/* 浮动粒子 */
.particles-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 2px;
  height: 2px;
  background: radial-gradient(circle, rgba(100, 180, 255, 0.8) 0%, transparent 70%);
  border-radius: 50%;
  animation: float-particle 15s infinite ease-in-out;
}

.particle:nth-child(odd) {
  animation-duration: 20s;
}

.particle:nth-child(3n) {
  animation-duration: 25s;
  background: radial-gradient(circle, rgba(156, 100, 255, 0.6) 0%, transparent 70%);
}

.particle:nth-child(1) { top: 10%; left: 20%; animation-delay: 0s; }
.particle:nth-child(2) { top: 20%; left: 40%; animation-delay: 2s; }
.particle:nth-child(3) { top: 30%; left: 60%; animation-delay: 4s; }
.particle:nth-child(4) { top: 40%; left: 80%; animation-delay: 6s; }
.particle:nth-child(5) { top: 50%; left: 10%; animation-delay: 8s; }
.particle:nth-child(6) { top: 60%; left: 30%; animation-delay: 10s; }
.particle:nth-child(7) { top: 70%; left: 50%; animation-delay: 12s; }
.particle:nth-child(8) { top: 80%; left: 70%; animation-delay: 14s; }
.particle:nth-child(9) { top: 15%; left: 85%; animation-delay: 1s; }
.particle:nth-child(10) { top: 25%; left: 15%; animation-delay: 3s; }
.particle:nth-child(11) { top: 35%; left: 35%; animation-delay: 5s; }
.particle:nth-child(12) { top: 45%; left: 55%; animation-delay: 7s; }
.particle:nth-child(13) { top: 55%; left: 75%; animation-delay: 9s; }
.particle:nth-child(14) { top: 65%; left: 25%; animation-delay: 11s; }
.particle:nth-child(15) { top: 75%; left: 45%; animation-delay: 13s; }
.particle:nth-child(16) { top: 85%; left: 65%; animation-delay: 15s; }
.particle:nth-child(17) { top: 12%; left: 52%; animation-delay: 2.5s; }
.particle:nth-child(18) { top: 22%; left: 72%; animation-delay: 4.5s; }
.particle:nth-child(19) { top: 32%; left: 12%; animation-delay: 6.5s; }
.particle:nth-child(20) { top: 42%; left: 32%; animation-delay: 8.5s; }
.particle:nth-child(21) { top: 52%; left: 52%; animation-delay: 10.5s; }
.particle:nth-child(22) { top: 62%; left: 72%; animation-delay: 12.5s; }
.particle:nth-child(23) { top: 72%; left: 22%; animation-delay: 14.5s; }
.particle:nth-child(24) { top: 82%; left: 42%; animation-delay: 16.5s; }
.particle:nth-child(25) { top: 18%; left: 62%; animation-delay: 1.5s; }
.particle:nth-child(26) { top: 28%; left: 82%; animation-delay: 3.5s; }
.particle:nth-child(27) { top: 38%; left: 22%; animation-delay: 5.5s; }
.particle:nth-child(28) { top: 48%; left: 42%; animation-delay: 7.5s; }
.particle:nth-child(29) { top: 58%; left: 62%; animation-delay: 9.5s; }
.particle:nth-child(30) { top: 68%; left: 82%; animation-delay: 11.5s; }

@keyframes float-particle {
  0%, 100% {
    transform: translateY(0) translateX(0) scale(1);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  50% {
    transform: translateY(-100px) translateX(50px) scale(1.5);
    opacity: 0.8;
  }
  90% {
    opacity: 0.5;
  }
}

/* 光晕效果 */
.glow-effects {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.15;
  animation: glow-pulse 8s infinite ease-in-out;
}

.glow-1 {
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, #4285f4 0%, transparent 70%);
  top: -200px;
  left: -200px;
  animation-delay: 0s;
}

.glow-2 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, #9c27b0 0%, transparent 70%);
  bottom: -150px;
  right: -150px;
  animation-delay: 2s;
}

.glow-3 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, #00bcd4 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: 4s;
}

@keyframes glow-pulse {
  0%, 100% {
    opacity: 0.15;
    transform: scale(1);
  }
  50% {
    opacity: 0.25;
    transform: scale(1.1);
  }
}

/* Hero Section - Get Started 区域 */
.hero-section {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.hero-content {
  text-align: center;
  max-width: 800px;
  width: 100%;
  animation: hero-fade-in 1s ease-out;
}

@keyframes hero-fade-in {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.hero-logo {
  margin-bottom: 30px;
  animation: logo-pulse 3s infinite ease-in-out;
}

.logo-icon-large {
  font-size: 80px;
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  filter: drop-shadow(0 0 30px rgba(66, 133, 244, 0.5));
}

@keyframes logo-pulse {
  0%, 100% {
    transform: scale(1);
    filter: drop-shadow(0 0 30px rgba(66, 133, 244, 0.5));
  }
  50% {
    transform: scale(1.05);
    filter: drop-shadow(0 0 40px rgba(156, 39, 176, 0.7));
  }
}

.hero-title {
  font-size: 72px;
  font-weight: 800;
  margin: 0 0 20px 0;
  letter-spacing: -2px;
}

.gradient-text {
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 50%, #00bcd4 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  background-size: 200% 200%;
  animation: gradient-flow 5s ease infinite;
}

@keyframes gradient-flow {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.hero-subtitle {
  font-size: 32px;
  font-weight: 300;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 16px 0;
}

.highlight {
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 600;
}

.hero-description {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.6);
  line-height: 1.8;
  margin: 0 0 50px 0;
  max-width: 700px;
  margin-left: auto;
  margin-right: auto;
}

.hero-actions {
  display: flex;
  gap: 20px;
  justify-content: center;
  align-items: center;
}

/* Get Started 按钮 */
.get-started-btn {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 12px;
  padding: 18px 40px;
  font-size: 18px;
  font-weight: 600;
  color: #ffffff;
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 100%);
  border: none;
  border-radius: 50px;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 10px 40px rgba(66, 133, 244, 0.4);
}

.get-started-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.5s ease;
}

.get-started-btn:hover::before {
  left: 100%;
}

.get-started-btn:hover {
  transform: translateY(-3px) scale(1.05);
  box-shadow: 0 15px 50px rgba(156, 39, 176, 0.5);
}

.get-started-btn:active {
  transform: translateY(-1px) scale(1.02);
}

.btn-text {
  position: relative;
  z-index: 1;
}

.btn-icon {
  position: relative;
  z-index: 1;
  transition: transform 0.3s ease;
}

.get-started-btn:hover .btn-icon {
  transform: translateX(5px);
}

/* Learn More 按钮 */
.learn-more-btn {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 18px 32px;
  font-size: 16px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 50px;
  cursor: pointer;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.learn-more-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(66, 133, 244, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(66, 133, 244, 0.2);
}

/* 登录卡片动画 */
.card-slide-enter-active {
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.card-slide-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.6, 1);
}

.card-slide-enter-from {
  opacity: 0;
  transform: scale(0.85) translateY(60px);
}

.card-slide-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(-30px);
}

.card-slide-enter-to,
.card-slide-leave-from {
  opacity: 1;
  transform: scale(1) translateY(0);
}

/* 登录卡片包装器 */
.login-card-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  z-index: 15;
}

/* 登录卡片 */
.login-card {
  background: rgba(15, 15, 35, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(66, 133, 244, 0.3);
  border-radius: 24px;
  padding: 50px 45px;
  width: 100%;
  max-width: 460px;
  position: relative;
  box-shadow: 
    0 25px 60px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(66, 133, 244, 0.1) inset,
    0 0 40px rgba(66, 133, 244, 0.2);
}

.login-card::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(45deg, 
    #4285f4 0%, 
    #9c27b0 50%, 
    #00bcd4 100%);
  border-radius: 24px;
  z-index: -1;
  opacity: 0;
  filter: blur(20px);
  animation: glow-border 3s infinite ease-in-out;
}

@keyframes glow-border {
  0%, 100% {
    opacity: 0.3;
  }
  50% {
    opacity: 0.6;
  }
}

/* 登录头部 */
.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.logo-icon {
  animation: logo-pulse 3s infinite;
}

@keyframes logo-pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.logo-text {
  font-size: 32px;
  font-weight: bold;
  margin: 0;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: text-glow 4s infinite;
}

@keyframes text-glow {
  0%, 100% {
    text-shadow: 0 0 5px var(--primary-color);
  }
  50% {
    text-shadow: 0 0 20px var(--primary-color), 0 0 30px var(--primary-color);
  }
}

.logo-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

/* 标签页切换 */
.auth-tabs {
  position: relative;
  display: flex;
  gap: 8px;
  margin-top: 32px;
  padding: 6px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.auth-tab {
  position: relative;
  flex: 1;
  padding: 12px 24px;
  background: transparent;
  border: none;
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.5);
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 2;
}

.auth-tab:hover {
  color: rgba(255, 255, 255, 0.8);
  background: rgba(255, 255, 255, 0.05);
}

.auth-tab.active {
  color: rgba(255, 255, 255, 1);
}

.auth-tab .tab-text {
  position: relative;
  z-index: 1;
}

/* 活动指示器 */
.tab-indicator {
  position: absolute;
  top: 6px;
  left: 6px;
  width: calc(50% - 10px);
  height: calc(100% - 12px);
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 100%);
  border-radius: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1;
  box-shadow: 
    0 4px 12px rgba(66, 133, 244, 0.3),
    0 0 0 1px rgba(66, 133, 244, 0.2) inset;
}

.tab-indicator::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(135deg, #4285f4 0%, #9c27b0 100%);
  border-radius: 12px;
  filter: blur(8px);
  opacity: 0.6;
  z-index: -1;
}

.tab-indicator.to-register {
  transform: translateX(calc(100% + 8px));
}

/* 注册表单特殊样式 */
.register-form {
  animation: form-slide-in 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes form-slide-in {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* 登录表单 */
.login-form {
  margin-bottom: 24px;
}

.login-input {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.login-input:hover {
  border-color: var(--primary-color);
  background: rgba(255, 255, 255, 0.08);
}

.login-input :deep(.el-input__wrapper) {
  background: transparent;
  box-shadow: none;
}

.login-button {
  width: 100%;
  background: var(--gradient-primary);
  border: none;
  border-radius: 12px;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(30, 136, 229, 0.3);
}

.button-text {
  flex: 1;
  text-align: center;
}

.button-icon {
  transition: transform 0.3s ease;
}

.login-button:hover .button-icon {
  transform: translateX(4px);
}

/* 表单底部 */
.form-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.remember-checkbox {
  color: var(--text-secondary);
}

.remember-checkbox :deep(.el-checkbox__label) {
  color: var(--text-secondary);
}

.forgot-password-btn {
  color: var(--primary-color) !important;
  font-size: 14px;
}

.forgot-password-btn:hover {
  color: var(--secondary-color) !important;
}

/* 分割线 */
.divider {
  position: relative;
  text-align: center;
  margin: 24px 0;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--border-color), transparent);
}

.divider-text {
  background: var(--bg-card);
  padding: 0 16px;
  color: var(--text-secondary);
  font-size: 14px;
  position: relative;
  z-index: 1;
}

/* 社交登录按钮 */
.social-login {
  display: flex;
  gap: 12px;
  margin: 20px 0;
}

.social-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 14px 20px;
  font-size: 15px;
  font-weight: 500;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  color: #ffffff;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.social-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
}

.github-btn {
  border-color: rgba(36, 41, 46, 0.4);
  background: linear-gradient(135deg, rgba(36, 41, 46, 0.4), rgba(36, 41, 46, 0.2));
}

.github-btn:hover {
  border-color: rgba(36, 41, 46, 0.6);
  background: linear-gradient(135deg, rgba(36, 41, 46, 0.6), rgba(36, 41, 46, 0.4));
  box-shadow: 0 10px 30px rgba(36, 41, 46, 0.5);
}

.social-icon {
  width: 20px;
  height: 20px;
}

/* 返回按钮 */
.back-to-home {
  margin-top: 25px;
  text-align: center;
}

.back-btn {
  color: rgba(255, 255, 255, 0.5) !important;
  transition: all 0.3s ease;
  font-size: 14px;
}

.back-btn:hover {
  color: rgba(66, 133, 244, 0.9) !important;
  transform: translateX(-3px);
}

/* 替代登录 */
.alternative-login {
  display: flex;
  gap: 12px;
}

.alt-button {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  height: 40px;
  color: var(--text-secondary);
  transition: all 0.3s ease;
}

.alt-button:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: var(--primary-color);
  color: var(--text-primary);
}

/* 装饰线条 */
.decoration-lines {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 5;
}

.line {
  position: absolute;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(66, 133, 244, 0.4) 50%,
    transparent 100%
  );
  opacity: 0;
  animation: line-scan 8s infinite ease-in-out;
}

.line-1 {
  top: 20%;
  left: -100%;
  width: 100%;
  height: 1px;
  animation-delay: 0s;
}

.line-2 {
  top: 50%;
  left: -100%;
  width: 100%;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(156, 39, 176, 0.4) 50%,
    transparent 100%
  );
  animation-delay: 2.5s;
}

.line-3 {
  top: 80%;
  left: -100%;
  width: 100%;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(0, 188, 212, 0.4) 50%,
    transparent 100%
  );
  animation-delay: 5s;
}

@keyframes line-scan {
  0% {
    left: -100%;
    opacity: 0;
  }
  20% {
    opacity: 1;
  }
  50% {
    left: 100%;
    opacity: 1;
  }
  51% {
    opacity: 0;
  }
  100% {
    left: 100%;
    opacity: 0;
  }
}

/* 修复浏览器自动填充导致的背景色问题 */
:deep(.el-input__inner:-webkit-autofill),
:deep(.el-input__inner:-webkit-autofill:hover),
:deep(.el-input__inner:-webkit-autofill:focus),
:deep(.el-input__inner:-webkit-autofill:active) {
  -webkit-box-shadow: 0 0 0 1000px rgba(30, 30, 50, 0.8) inset !important;
  -webkit-text-fill-color: #ffffff !important;
  transition: background-color 5000s ease-in-out 0s;
  caret-color: #ffffff !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-card {
    margin: 20px;
    padding: 30px 20px;
  }
  
  .logo-text {
    font-size: 28px;
  }
  
  .alternative-login {
    flex-direction: column;
  }
  
  .particle {
    display: none;
  }
}
</style>