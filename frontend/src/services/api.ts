import axios from 'axios'
import type { Task, FileInfo } from '@/types/task'

// API基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ? `${import.meta.env.VITE_API_BASE_URL}/api/v1` : '/api/v1'

// 创建axios实例
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30秒超时
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`)
    
    // 添加认证token
    const token = localStorage.getItem('auth_token') || sessionStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    return config
  },
  (error) => {
    console.error('API Request Error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    console.log(`API Response: ${response.status} ${response.config.url}`)
    return response
  },
  (error) => {
    console.error('API Response Error:', error)
    if (error.response) {
      // 服务器响应错误
      const message = error.response.data?.error || error.response.statusText
      console.error(`API Error ${error.response.status}: ${message}`)
    } else if (error.request) {
      // 请求发送失败
      console.error('API Request failed - no response received')
    } else {
      // 请求配置错误
      console.error('API Request configuration error:', error.message)
    }
    return Promise.reject(error)
  }
)

// 任务相关API
export const taskAPI = {
  // 创建任务
  createTask: async (type: string, name: string, description: string, params: any) => {
    const response = await api.post('/tasks', {
      type,
      name,
      description,
      params,
    })
    return response.data
  },

  // 获取任务列表
  getTasks: async (): Promise<Task[]> => {
    const response = await api.get('/tasks')
    return response.data.tasks || []
  },

  // 获取单个任务
  getTask: async (taskId: string): Promise<Task> => {
    const response = await api.get(`/tasks/${taskId}`)
    return response.data.task
  },

  // 获取任务结果
  getTaskResult: async (taskId: string) => {
    const response = await api.get(`/tasks/${taskId}/result`)
    return response.data
  },

  // 取消任务
  cancelTask: async (taskId: string) => {
    const response = await api.delete(`/tasks/${taskId}`)
    return response.data
  },
}

// 文件相关API
export const fileAPI = {
  // 上传文件
  uploadFile: async (file: File): Promise<FileInfo> => {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post('/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  },

  // 获取文件列表
  getFiles: async (): Promise<FileInfo[]> => {
    const response = await api.get('/files/list')
    return response.data.files || []
  },

  // 删除文件
  deleteFile: async (fileId: string) => {
    const response = await api.delete(`/files/${fileId}`)
    return response.data
  },

  // 批量处理文件
  batchProcessFiles: async (fileIds: string[], operation: string, params: any) => {
    const response = await api.post('/files/batch', {
      file_ids: fileIds,
      operation,
      params,
    })
    return response.data
  },
}

// Git相关API
export const gitAPI = {
  // 克隆Git仓库
  cloneRepository: async (remoteUrl: string, targetPath: string) => {
    const response = await api.post('/git/clone', {
      remote_url: remoteUrl,
      target_path: targetPath,
    })
    return response.data
  },

  // 分析Git仓库
  analyzeRepository: async (repoPath: string, remoteUrl?: string, cloneIfNotExists = true) => {
    const response = await api.post('/git/analyze', {
      repo_path: repoPath,
      remote_url: remoteUrl,
      clone_if_not_exists: cloneIfNotExists,
    })
    return response.data
  },

  // 获取文件历史
  getFileHistory: async (repoPath: string, filePath: string) => {
    const response = await api.post(`/git/history/${filePath}`, {
      repo_path: repoPath,
      file_path: filePath,
    })
    return response.data
  },

  // 获取文件差异
  getFileDiff: async (repoPath: string, filePath: string, fromCommit: string, toCommit: string) => {
    const response = await api.post('/git/diff', {
      repo_path: repoPath,
      file_path: filePath,
      from_commit: fromCommit,
      to_commit: toCommit,
    })
    return response.data
  },
}

// 健康检查API
export const healthAPI = {
  // 基础健康检查
  checkHealth: async () => {
    const response = await api.get('/health')
    return response.data
  },

  // 详细健康检查
  checkDetailedHealth: async () => {
    const response = await api.get('/health/detailed')
    return response.data
  },
}

// 认证相关API
export const authAPI = {
  // 用户注册
  register: async (username: string, email: string, password: string, confirmPassword: string) => {
    const response = await api.post('/auth/register', {
      username,
      email,
      password,
      confirmPassword,
    })
    return response.data
  },

  // 用户登录
  login: async (username: string, password: string) => {
    const response = await api.post('/auth/login', {
      username,
      password,
    })
    return response.data
  },

  // 获取用户资料
  getProfile: async () => {
    const response = await api.get('/auth/profile')
    return response.data
  },

  // 更新用户资料
  updateProfile: async (data: { 
    username?: string
    email?: string
    avatar?: string
    phone?: string
    bio?: string
    location?: string
    occupation?: string
    company?: string
    website?: string
    twitter?: string
    github_url?: string
  }) => {
    const response = await api.put('/auth/profile/update', data)
    return response.data
  },

  // 修改密码
  changePassword: async (oldPassword: string, newPassword: string) => {
    const response = await api.post('/auth/password/change', {
      oldPassword,
      newPassword,
    })
    return response.data
  },

  // 根据ID获取用户
  getUserById: async (userId: number) => {
    const response = await api.get(`/users/${userId}`)
    return response.data
  },

  // 忘记密码
  forgotPassword: async (email: string) => {
    const response = await api.post('/auth/forgot-password', { email })
    return response.data
  },

  // 验证重置码
  verifyResetCode: async (email: string, code: string) => {
    const response = await api.post('/auth/verify-reset-code', { email, code })
    return response.data
  },

  // 重置密码
  resetPassword: async (email: string, code: string, newPassword: string) => {
    const response = await api.post('/auth/reset-password', {
      email,
      code,
      newPassword,
    })
    return response.data
  },

  // GitHub 登录
  githubLogin: async (githubId: string) => {
    const response = await api.post('/auth/github/login', { github_id: githubId })
    return response.data
  },

  // 绑定 GitHub
  bindGithub: async (githubId: string, githubUsername: string, email: string, avatar: string) => {
    const response = await api.post('/auth/github/bind', {
      github_id: githubId,
      github_username: githubUsername,
      email,
      avatar,
    })
    return response.data
  },
}

export default api