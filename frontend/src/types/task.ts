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

export interface TaskCreateData {
  type: string
  name: string
  description?: string
  params?: any
}

export interface AgentThought {
  message: string
  timestamp: number
  type: 'thought' | 'action' | 'result'
}

export interface FileInfo {
  id: string
  filename: string
  size: number
  uploadedAt: string
}

export interface TaskCreateRequest {
  type: string
  name: string
  description?: string
  params?: any
}

export interface TaskCreateResponse {
  task_id: string
  message: string
}

export interface TaskListResponse {
  tasks: Task[]
  total: number
}

export interface TaskResponse {
  task: Task
}

export interface TaskResultResponse {
  task_id: string
  result: any
  status: string
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

export interface ConversionResult {
  original: string
  converted: string
  diff: string
  warnings: string[]
}

export interface TestGenerationResult {
  testCode: string
  testFramework: string
  coverage: number
  testCases: TestCase[]
}

export interface TestCase {
  name: string
  description: string
  code: string
}