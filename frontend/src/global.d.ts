// Global type declarations for the project

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module '@/stores/auth' {
  export interface UserInfo {
    username: string
    token: string
    role: 'admin' | 'user' | 'guest' | 'local'
    rememberMe?: boolean
  }

  export const useAuthStore: () => {
    userInfo: import('vue').Ref<UserInfo | null>
    isAuthenticated: import('vue').ComputedRef<boolean>
    isGuest: import('vue').ComputedRef<boolean>
    isLocal: import('vue').ComputedRef<boolean>
    isAdmin: import('vue').ComputedRef<boolean>
    login: (info: UserInfo) => void
    logout: () => void
    loadUserInfo: () => void
    updateUserInfo: (updates: Partial<UserInfo>) => void
    clearStorage: () => void
  }
}

declare module '@/stores/task' {
  export interface Task {
    id: string
    name: string
    status: 'pending' | 'running' | 'completed' | 'failed'
    progress: number
    result?: any
    error?: string
    createdAt: string
    updatedAt?: string
    type: 'analysis' | 'conversion' | 'test_generation' | 'git_analysis'
    filePath?: string
    agentThoughts?: AgentThought[]
  }

  export interface AgentThought {
    message: string
    timestamp: number
    type: 'thought' | 'action' | 'result'
  }

  export interface CodeAnalysisResult {
    issues: CodeIssue[]
    metrics: CodeMetrics
    suggestions: CodeSuggestion[]
  }

  export interface CodeIssue {
    line: number
    column: number
    severity: 'error' | 'warning' | 'info'
    message: string
    rule: string
  }

  export interface CodeMetrics {
    lines: number
    complexity: number
    maintainability: number
    testCoverage?: number
  }

  export interface CodeSuggestion {
    line: number
    original: string
    suggested: string
    reason: string
  }

  export const useTaskStore: () => {
    tasks: import('vue').Ref<Task[]>
    currentTask: import('vue').Ref<Task | null>
    loading: import('vue').Ref<boolean>
    runningTasks: import('vue').ComputedRef<Task[]>
    completedTasks: import('vue').ComputedRef<Task[]>
    createTask: (taskData: Partial<Task>) => Promise<Task>
    updateTask: (taskId: string, updates: Partial<Task>) => void
    loadTasks: () => Promise<void>
  }
}

// Monaco Editor types
declare module 'monaco-editor' {
  export interface IStandaloneCodeEditor {
    getValue(): string
    setValue(value: string): void
    dispose(): void
    onDidChangeModelContent(callback: () => void): void
  }

  export function create(domElement: HTMLElement, options: any): IStandaloneCodeEditor
}

// Element Plus icons
declare module '@element-plus/icons-vue' {
  export const House: any
  export const Search: any
  export const Refresh: any
  export const DocumentChecked: any
  export const Coin: any
  export const Setting: any
  export const Plus: any
  export const User: any
  export const SwitchButton: any
  export const ArrowDown: any
  export const Upload: any
  export const Delete: any
  export const Download: any
  export const CopyDocument: any
  export const MagicStick: any
  export const Link: any
  export const Right: any
  export const Key: any
  export const Cpu: any
  export const ArrowRight: any
  export const Document: any
  export const TrendCharts: any
  export const Clock: any
  export const CircleCheck: any
}

// ECharts
declare module 'echarts' {
  export function init(dom: HTMLElement): any
}

// Diff2Html
declare module 'diff2html' {
  export function parse(diff: string): any
  export function html(json: any, options?: any): string
}

// Diff
declare module 'diff' {
  export function createTwoFilesPatch(
    oldFileName: string,
    newFileName: string,
    oldStr: string,
    newStr: string,
    oldHeader?: string,
    newHeader?: string,
    options?: any
  ): string
}