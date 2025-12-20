import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Task } from '@/types/task'
import { taskAPI } from '@/services/api'
import { wsManager, TaskProgressHandler, type WebSocketMessage } from '@/services/websocket'

// localStorage key for tasks
const TASKS_STORAGE_KEY = 'codesage_tasks'
const MAX_STORED_TASKS = 100 // 最多存储100个任务

export const useTaskStore = defineStore('task', () => {
  const tasks = ref<Task[]>([])
  const currentTask = ref<Task | null>(null)
  const loading = ref(false)

  // 从localStorage加载任务
  const loadTasksFromStorage = () => {
    try {
      const stored = localStorage.getItem(TASKS_STORAGE_KEY)
      if (stored) {
        const parsedTasks = JSON.parse(stored)
        tasks.value = parsedTasks
      }
    } catch (error) {
      console.error('Failed to load tasks from storage:', error)
    }
  }

  // 保存任务到localStorage
  const saveTasksToStorage = () => {
    try {
      // 只保存最近的任务，避免localStorage过大
      const tasksToSave = tasks.value.slice(0, MAX_STORED_TASKS)
      localStorage.setItem(TASKS_STORAGE_KEY, JSON.stringify(tasksToSave))
    } catch (error) {
      console.error('Failed to save tasks to storage:', error)
    }
  }

  // 初始化时加载任务
  loadTasksFromStorage()

  // 计算属性
  const runningTasks = computed(() =>
    tasks.value.filter((task) => task.status === 'running')
  )

  const completedTasks = computed(() =>
    tasks.value.filter((task) => task.status === 'completed')
  )

  const failedTasks = computed(() =>
    tasks.value.filter((task) => task.status === 'failed')
  )

  // 创建任务进度处理器
  const createTaskProgressHandler = (taskId: string) => {
    return new TaskProgressHandler((updates: Partial<Task>) => {
      updateTask(taskId, updates)
    })
  }

  // 动作
  const createTask = async (type: string, name: string, description: string, params: any) => {
    loading.value = true
    try {
      // 调用API创建任务
      const response = await taskAPI.createTask(type, name, description, params)
      const taskId = response.task_id
      
      // 创建本地任务对象
      const newTask: Task = {
        id: taskId,
        name,
        status: 'pending',
        progress: 0,
        createdAt: new Date().toISOString(),
        type: type as any,
        params,
      }
      
      tasks.value.unshift(newTask)
      
      // 保存到localStorage
      saveTasksToStorage()
      
      // 连接WebSocket监听任务进度
      const progressHandler = createTaskProgressHandler(taskId)
      wsManager.connectToTaskProgress(taskId, (message: WebSocketMessage) => {
        progressHandler.handleProgressMessage(message)
      })
      
      // 连接WebSocket监听Agent思考流
      wsManager.connectToAgentThought(taskId, (message: WebSocketMessage) => {
        progressHandler.handleAgentThoughtMessage(message)
      })
      
      return newTask
    } catch (error) {
      console.error('Failed to create task:', error)
      throw new Error(`创建任务失败: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      loading.value = false
    }
  }

  const updateTask = (taskId: string, updates: Partial<Task>) => {
    const index = tasks.value.findIndex((t) => t.id === taskId)
    if (index !== -1) {
      if (tasks.value[index]) {
        tasks.value[index] = { ...tasks.value[index], ...updates } as Task
      }
    }
    
    // 更新当前任务
    if (currentTask.value && currentTask.value.id === taskId) {
      currentTask.value = { ...currentTask.value, ...updates } as Task
    }
    
    // 保存到localStorage
    saveTasksToStorage()
  }

  const loadTasks = async () => {
    loading.value = true
    try {
      // 调用API获取任务列表
      const taskList = await taskAPI.getTasks()
      tasks.value = taskList
      
      // 为正在运行的任务连接WebSocket
      for (const task of taskList) {
        if (task.status === 'running') {
          const progressHandler = createTaskProgressHandler(task.id)
          wsManager.connectToTaskProgress(task.id, (message: WebSocketMessage) => {
            progressHandler.handleProgressMessage(message)
          })
          wsManager.connectToAgentThought(task.id, (message: WebSocketMessage) => {
            progressHandler.handleAgentThoughtMessage(message)
          })
        }
      }
    } catch (error) {
      console.error('Failed to load tasks:', error)
      throw new Error(`加载任务列表失败: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      loading.value = false
    }
  }

  const getTask = async (taskId: string) => {
    loading.value = true
    try {
      // 调用API获取任务详情
      const task = await taskAPI.getTask(taskId)
      
      // 更新本地任务列表
      updateTask(taskId, task)
      
      // 设置为当前任务
      currentTask.value = task
      
      return task
    } catch (error) {
      console.error(`Failed to get task ${taskId}:`, error)
      throw new Error(`获取任务详情失败: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      loading.value = false
    }
  }

  const getTaskResult = async (taskId: string) => {
    try {
      // 调用API获取任务结果
      const result = await taskAPI.getTaskResult(taskId)
      
      // 更新任务结果
      updateTask(taskId, { result: result.result })
      
      return result
    } catch (error) {
      console.error(`Failed to get task result ${taskId}:`, error)
      throw new Error(`获取任务结果失败: ${error instanceof Error ? error.message : '未知错误'}`)
    }
  }

  const cancelTask = async (taskId: string) => {
    try {
      // 调用API取消任务
      await taskAPI.cancelTask(taskId)
      
      // 更新任务状态
      updateTask(taskId, { status: 'cancelled' as any })
      
      // 断开WebSocket连接
      wsManager.disconnect(taskId)
    } catch (error) {
      console.error(`Failed to cancel task ${taskId}:`, error)
      throw new Error(`取消任务失败: ${error instanceof Error ? error.message : '未知错误'}`)
    }
  }

  const setCurrentTask = (task: Task | null) => {
    currentTask.value = task
  }

  // 清理资源
  const cleanup = () => {
    wsManager.disconnectAll()
  }

  return {
    tasks,
    currentTask,
    loading,
    runningTasks,
    completedTasks,
    failedTasks,
    createTask,
    updateTask,
    loadTasks,
    getTask,
    getTaskResult,
    cancelTask,
    setCurrentTask,
    cleanup,
  }
})