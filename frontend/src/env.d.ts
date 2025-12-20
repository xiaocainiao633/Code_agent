/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module '@/types/task' {
  export interface Task {
    id: string
    name: string
    status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled'
    progress: number
    result?: any
    error?: string
    createdAt: string
    updatedAt?: string
    startedAt?: string
    completedAt?: string
    type: 'analysis' | 'convert' | 'test' | 'git_clone' | 'git_analyze' | 'git_history' | 'git_diff' | 'batch'
    description?: string
    params?: any
    filePath?: string
    agentThoughts?: AgentThought[]
  }

  export interface FileInfo {
    id: string
    filename: string
    size: number
    uploadedAt: string
  }

  export interface AgentThought {
    message: string
    timestamp: number
    type: 'thought' | 'action' | 'result'
  }
}

declare module '@/stores/task' {
  export const useTaskStore: any
}