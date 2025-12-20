import type { Task } from '@/types/task'

// WebSocket消息类型
export const WebSocketMessageType = {
  TASK_PROGRESS: 'task_progress' as const,
  AGENT_THOUGHT: 'agent_thought' as const,
  SYSTEM: 'system' as const,
  PING: 'ping' as const,
  PONG: 'pong' as const,
}

export type WebSocketMessageType = typeof WebSocketMessageType[keyof typeof WebSocketMessageType]

// WebSocket消息接口
export interface WebSocketMessage {
  type: WebSocketMessageType
  task_id?: string
  timestamp: string 
  data: any
}

// WebSocket连接管理器
export class WebSocketManager {
  private connections: Map<string, WebSocket> = new Map()
  private messageHandlers: Map<string, (message: WebSocketMessage) => void> = new Map()
  private reconnectAttempts: Map<string, number> = new Map()
  private maxReconnectAttempts = 3
  private reconnectDelay = 1000 // 1秒

  // 连接到任务进度WebSocket
  connectToTaskProgress(taskId: string, onMessage: (message: WebSocketMessage) => void): void {
    const wsBaseUrl = import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8082'
    const url = `${wsBaseUrl}/ws/progress/${taskId}`
    this.connect(taskId, url, onMessage)
  }

  // 连接到Agent思考流WebSocket
  connectToAgentThought(taskId: string, onMessage: (message: WebSocketMessage) => void): void {
    const wsBaseUrl = import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8082'
    const url = `${wsBaseUrl}/ws/agent/${taskId}`
    this.connect(taskId, url, onMessage)
  }

  // 建立WebSocket连接
  private connect(taskId: string, url: string, onMessage: (message: WebSocketMessage) => void): void {
    try {
      console.log(`Connecting to WebSocket: ${url}`)
      const ws = new WebSocket(url)

      // 存储消息处理器
      this.messageHandlers.set(taskId, onMessage)

      ws.onopen = () => {
        console.log(`WebSocket connected: ${url}`)
        this.connections.set(taskId, ws)
        this.reconnectAttempts.set(taskId, 0)
        
        // 发送连接成功消息
        onMessage({
          type: WebSocketMessageType.SYSTEM,
          task_id: taskId,
          timestamp: new Date().toISOString(),
          data: { message: 'WebSocket connected successfully' }
        })
      }

      ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          console.log(`WebSocket message received:`, message)
          
          // 处理ping/pong
          if (message.type === WebSocketMessageType.PING) {
            this.sendPong(taskId)
            return
          }

          // 调用消息处理器
          const handler = this.messageHandlers.get(taskId)
          if (handler) {
            handler(message)
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      ws.onerror = (error) => {
        console.error(`WebSocket error: ${url}`, error)
      }

      ws.onclose = () => {
        console.log(`WebSocket disconnected: ${url}`)
        this.connections.delete(taskId)
        
        // 尝试重连
        this.attemptReconnect(taskId, url, onMessage)
      }

    } catch (error) {
      console.error(`Failed to create WebSocket connection: ${url}`, error)
    }
  }

  // 发送Pong响应
  private sendPong(taskId: string): void {
    const ws = this.connections.get(taskId)
    if (ws && ws.readyState === WebSocket.OPEN) {
      const pongMessage: WebSocketMessage = {
        type: WebSocketMessageType.PONG,
        task_id: taskId,
        timestamp: new Date().toISOString(),
        data: {}
      }
      ws.send(JSON.stringify(pongMessage))
    }
  }

  // 尝试重连
  private attemptReconnect(taskId: string, url: string, onMessage: (message: WebSocketMessage) => void): void {
    const attempts = this.reconnectAttempts.get(taskId) || 0
    
    if (attempts < this.maxReconnectAttempts) {
      console.log(`Attempting to reconnect WebSocket (${attempts + 1}/${this.maxReconnectAttempts}): ${url}`)
      this.reconnectAttempts.set(taskId, attempts + 1)
      
      setTimeout(() => {
        this.connect(taskId, url, onMessage)
      }, this.reconnectDelay * (attempts + 1)) // 递增延迟
    } else {
      console.error(`Max reconnection attempts reached for WebSocket: ${url}`)
      
      // 发送连接失败消息
      onMessage({
        type: WebSocketMessageType.SYSTEM,
        task_id: taskId,
        timestamp: new Date().toISOString(),
        data: { message: 'WebSocket connection failed after max retry attempts' }
      })
    }
  }

  // 断开连接
  disconnect(taskId: string): void {
    const ws = this.connections.get(taskId)
    if (ws) {
      console.log(`Disconnecting WebSocket: ${taskId}`)
      ws.close()
      this.connections.delete(taskId)
      this.messageHandlers.delete(taskId)
      this.reconnectAttempts.delete(taskId)
    }
  }

  // 断开所有连接
  disconnectAll(): void {
    console.log('Disconnecting all WebSocket connections')
    for (const taskId of this.connections.keys()) {
      this.disconnect(taskId)
    }
  }

  // 检查连接状态
  isConnected(taskId: string): boolean {
    const ws = this.connections.get(taskId)
    return ws !== undefined && ws.readyState === WebSocket.OPEN
  }

  // 获取连接数
  getConnectionCount(): number {
    return this.connections.size
  }
}

// 创建全局WebSocket管理器实例
export const wsManager = new WebSocketManager()

// 任务进度处理器
export class TaskProgressHandler {
  private taskUpdateCallback: (task: Partial<Task>) => void

  constructor(taskUpdateCallback: (task: Partial<Task>) => void) {
    this.taskUpdateCallback = taskUpdateCallback
  }

  // 处理任务进度消息
  handleProgressMessage(message: WebSocketMessage): void {
    if (message.type === WebSocketMessageType.TASK_PROGRESS) {
      const { progress, status, message: msg } = message.data
      
      this.taskUpdateCallback({
        id: message.task_id!,
        progress,
        status,
        // 可以根据需要添加更多字段
      })
    }
  }

  // 处理Agent思考消息
  handleAgentThoughtMessage(message: WebSocketMessage): void {
    if (message.type === WebSocketMessageType.AGENT_THOUGHT) {
      const { thought, step } = message.data
      
      // 可以在这里处理Agent思考流，比如更新UI显示
      console.log(`Agent Thought [${step}]: ${thought}`)
      
      // 可以触发回调来更新UI
      this.taskUpdateCallback({
        id: message.task_id!,
        agentThoughts: [{
          message: thought,
          timestamp: Date.now(),
          type: step === 'result' ? 'result' : 'thought'
        }]
      })
    }
  }
}

export default wsManager