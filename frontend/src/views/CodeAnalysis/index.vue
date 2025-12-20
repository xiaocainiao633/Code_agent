<template>
  <div class="code-analysis-container">
    <el-row :gutter="24">
      <!-- 左侧代码编辑器 -->
      <el-col :span="12">
        <el-card class="code-card">
          <template #header>
            <div class="card-header">
              <h3 class="card-title">源代码</h3>
              <el-upload
                action="#"
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleFileChange"
                accept=".py,.js,.java,.go,.rs,.ts,.cpp,.c"
              >
                <el-button type="primary" :icon="Upload" size="small">上传文件</el-button>
              </el-upload>
            </div>
          </template>
          
          <div class="editor-container">
            <div ref="editorContainer" class="monaco-editor-container"></div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 右侧分析结果 -->
      <el-col :span="12">
        <el-card class="result-card">
          <template #header>
            <div class="card-header">
              <h3 class="card-title">分析结果</h3>
              <el-button type="success" :icon="Download" size="small" @click="downloadReport">
                下载报告
              </el-button>
            </div>
          </template>
          
          <div class="result-content" v-if="analyzing">
            <el-skeleton :rows="8" animated />
            <div style="text-align: center; margin-top: 20px;">
              <el-text type="info">正在分析代码，请稍候...</el-text>
            </div>
          </div>
          <div class="result-content" v-else-if="analysisResult">
            <!-- 代码质量评分 -->
            <div class="quality-score">
              <el-progress
                type="dashboard"
                :percentage="analysisResult.quality_score || 0"
                :color="getScoreColor(analysisResult.quality_score || 0)"
              >
                <template #default="{ percentage }">
                  <span class="percentage-value">{{ percentage }}</span>
                  <span class="percentage-label">质量评分</span>
                </template>
              </el-progress>
            </div>
            
            <!-- 复杂度信息 -->
            <div class="metrics-section">
              <el-descriptions :column="2" border class="metrics-descriptions">
                <el-descriptions-item 
                  label="复杂度分数"
                  v-if="shouldShowComplexity"
                >
                  {{ (analysisResult.complexity_score * 100).toFixed(1) }}%
                </el-descriptions-item>
                <el-descriptions-item label="依赖数量">
                  {{ analysisResult.dependencies?.length || 0 }}
                </el-descriptions-item>
                <el-descriptions-item 
                  label="问题数量"
                  v-if="shouldShowSecurity || shouldShowSyntax"
                >
                  {{ analysisResult.issues?.length || 0 }}
                </el-descriptions-item>
                <el-descriptions-item label="建议数量">
                  {{ analysisResult.suggestions?.length || 0 }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
            
            <!-- 业务逻辑摘要 -->
            <div class="summary-section" v-if="analysisResult.business_logic_summary">
              <h4 class="section-title">业务逻辑摘要</h4>
              <el-alert
                :title="analysisResult.business_logic_summary"
                type="info"
                :closable="false"
                show-icon
              />
            </div>
            
            <!-- 依赖列表 -->
            <div class="dependencies-section" v-if="analysisResult.dependencies && analysisResult.dependencies.length > 0">
              <h4 class="section-title">外部依赖</h4>
              <el-tag
                v-for="(dep, index) in analysisResult.dependencies"
                :key="index"
                style="margin: 4px"
              >
                {{ dep }}
              </el-tag>
            </div>
            
            <!-- 问题列表 -->
            <div 
              class="issues-section" 
              v-if="(shouldShowSecurity || shouldShowSyntax) && analysisResult.issues && analysisResult.issues.length > 0"
            >
              <h4 class="section-title">发现的问题 ({{ analysisResult.issues.length }})</h4>
              <el-timeline class="issues-timeline">
                <el-timeline-item
                  v-for="(issue, index) in analysisResult.issues"
                  :key="index"
                  :type="getIssueType(issue.severity)"
                  :timestamp="issue.line > 0 ? `行 ${issue.line}` : '全局'"
                >
                  <div class="issue-item">
                    <el-tag :type="getIssueType(issue.severity)" size="small">
                      {{ issue.severity }}
                    </el-tag>
                    <el-tag size="small" style="margin-left: 8px">{{ issue.type }}</el-tag>
                    <p class="issue-message">{{ issue.message }}</p>
                  </div>
                </el-timeline-item>
              </el-timeline>
            </div>
            
            <!-- 建议列表 -->
            <div 
              class="suggestions-section" 
              v-if="shouldShowPerformance && analysisResult.suggestions && analysisResult.suggestions.length > 0"
            >
              <h4 class="section-title">优化建议 ({{ analysisResult.suggestions.length }})</h4>
              <el-collapse class="suggestions-collapse">
                <el-collapse-item
                  v-for="(suggestion, index) in analysisResult.suggestions"
                  :key="index"
                  :title="suggestion.title || `建议 ${index + 1}`"
                  :name="index"
                >
                  <p class="suggestion-description">{{ suggestion.description }}</p>
                  <div v-if="suggestion.code_example" class="code-example">
                    <strong>示例代码：</strong>
                    <pre><code>{{ suggestion.code_example }}</code></pre>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
            
            <!-- 无问题提示 -->
            <div v-if="shouldShowNoIssuesMessage">
              <el-result
                icon="success"
                title="代码质量良好"
                sub-title="未发现明显问题，继续保持！"
              >
                <template #extra>
                  <el-text type="info" size="small">
                    已完成的分析：{{ getCompletedAnalysisTypes() }}
                  </el-text>
                </template>
              </el-result>
            </div>
          </div>
          <el-empty v-else description="请先进行代码分析" />
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 分析配置 -->
    <el-row :gutter="24" class="config-section">
      <el-col :span="24">
        <el-card class="config-card">
          <template #header>
            <h3 class="card-title">分析配置</h3>
          </template>
          
          <el-form :model="analysisConfig" label-position="left" inline>
            <el-form-item label="语言类型">
              <el-select 
                v-model="analysisConfig.language" 
                placeholder="选择语言"
                popper-class="language-select-dropdown"
                style="width: 150px"
              >
                <el-option label="Python" value="python" />
                <el-option label="JavaScript" value="javascript" />
                <el-option label="Java" value="java" />
                <el-option label="Go" value="go" />
                <el-option label="TypeScript" value="typescript" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="分析类型">
              <el-checkbox-group v-model="analysisConfig.analysisTypes">
                <el-checkbox label="syntax">语法检查</el-checkbox>
                <el-checkbox label="complexity">复杂度分析</el-checkbox>
                <el-checkbox label="security">安全检测</el-checkbox>
                <el-checkbox label="performance">性能分析</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                :icon="Search"
                @click="analyzeCode"
                :loading="analyzing"
                size="large"
              >
                开始分析
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, Download, Search } from '@element-plus/icons-vue'
import * as monaco from 'monaco-editor'
import { useTaskStore } from '@/stores/task'
import { loadEditorSettings, toMonacoOptions, updateEditorOptions, listenToSettingsChanges, type EditorSettings } from '../../utils/editorConfig'

const editorContainer = ref<HTMLElement>()
let editor: any
let saveTimer: any = null
const taskStore = useTaskStore()
const currentTaskId = ref<string | null>(null)

// 代码自动保存的localStorage key
const CODE_STORAGE_KEY = 'code_analysis_source_code'
const LANGUAGE_STORAGE_KEY = 'code_analysis_language'
const RESULT_STORAGE_KEY = 'code_analysis_result'

// 从localStorage加载保存的代码（如果启用了自动保存）
const loadSavedCode = () => {
  const allSettings = JSON.parse(localStorage.getItem('codesage-settings') || '{}')
  const settings = allSettings.general || {}
  if (settings.autoSave) {
    const savedCode = localStorage.getItem(CODE_STORAGE_KEY)
    const savedLanguage = localStorage.getItem(LANGUAGE_STORAGE_KEY)
    const savedResult = localStorage.getItem(RESULT_STORAGE_KEY)
    
    if (savedCode) {
      sourceCode.value = savedCode
    }
    if (savedLanguage) {
      analysisConfig.value.language = savedLanguage
    }
    if (savedResult) {
      try {
        analysisResult.value = JSON.parse(savedResult)
      } catch (error) {
        console.error('Failed to parse saved result:', error)
      }
    }
  }
}

const sourceCode = ref(`# Python 代码示例
def calculate_sum(numbers):
    total = 0
    for num in numbers:
        total += num
    return total

def process_data(data):
    result = []
    for item in data:
        if item > 0:
            result.append(item * 2)
    return result

class DataProcessor:
    def __init__(self):
        self.data = []
    
    def add_data(self, value):
        self.data.append(value)
    
    def get_average(self):
        if len(self.data) == 0:
            return 0
        return sum(self.data) / len(self.data)
`)

// 保存代码到localStorage（如果启用了自动保存）
const saveCodeToStorage = () => {
  const allSettings = JSON.parse(localStorage.getItem('codesage-settings') || '{}')
  const settings = allSettings.general || {}
  if (settings.autoSave) {
    localStorage.setItem(CODE_STORAGE_KEY, sourceCode.value)
    localStorage.setItem(LANGUAGE_STORAGE_KEY, analysisConfig.value.language)
    if (analysisResult.value) {
      localStorage.setItem(RESULT_STORAGE_KEY, JSON.stringify(analysisResult.value))
    }
  }
}

const analysisResult = ref<any>(null)
const analyzing = ref(false)

const analysisConfig = ref({
  language: 'python',
  analysisTypes: ['syntax', 'complexity']
})

// 计算属性：根据分析类型判断是否显示对应部分
const shouldShowSyntax = computed(() => analysisConfig.value.analysisTypes.includes('syntax'))
const shouldShowComplexity = computed(() => analysisConfig.value.analysisTypes.includes('complexity'))
const shouldShowSecurity = computed(() => analysisConfig.value.analysisTypes.includes('security'))
const shouldShowPerformance = computed(() => analysisConfig.value.analysisTypes.includes('performance'))

// 计算属性：判断是否显示"无问题"提示
const shouldShowNoIssuesMessage = computed(() => {
  if (!analysisResult.value) return false
  
  const hasIssues = analysisResult.value.issues && analysisResult.value.issues.length > 0
  const hasSuggestions = analysisResult.value.suggestions && analysisResult.value.suggestions.length > 0
  
  // 如果选择了安全或语法检查，但没有问题
  const noIssues = (shouldShowSecurity.value || shouldShowSyntax.value) && !hasIssues
  
  // 如果选择了性能分析，但没有建议
  const noSuggestions = shouldShowPerformance.value && !hasSuggestions
  
  // 只有当所有选中的分析类型都没有发现问题时，才显示"无问题"提示
  return (noIssues || (!shouldShowSecurity.value && !shouldShowSyntax.value)) && 
         (noSuggestions || !shouldShowPerformance.value)
})

// 获取已完成的分析类型
const getCompletedAnalysisTypes = () => {
  const typeNames: Record<string, string> = {
    'syntax': '语法检查',
    'complexity': '复杂度分析',
    'security': '安全检测',
    'performance': '性能分析'
  }
  
  return analysisConfig.value.analysisTypes
    .map(type => typeNames[type])
    .join('、')
}

// 监听语言变化，保存到localStorage
watch(() => analysisConfig.value.language, () => {
  saveCodeToStorage()
  // 更新编辑器语言
  if (editor) {
    const model = editor.getModel()
    if (model) {
      (monaco as any).editor.setModelLanguage(model, analysisConfig.value.language)
    }
  }
})

onMounted(() => {
  // 加载保存的代码
  loadSavedCode()
  
  setTimeout(() => {
    initEditor()
  }, 100)
})

onUnmounted(() => {
  // 清理定时器
  if (saveTimer) {
    clearTimeout(saveTimer)
  }
  
  // 保存代码到localStorage
  saveCodeToStorage()
  
  if (editor) {
    if ((editor as any)._settingsCleanup) {
      ;(editor as any)._settingsCleanup()
    }
    editor.dispose()
  }
  
  if (currentTaskId.value) {
    taskStore.cancelTask(currentTaskId.value)
  }
})

const initEditor = () => {
  if (!editorContainer.value) {
    console.error('Editor container not found')
    return
  }
  
  editorContainer.value.style.height = '500px'
  
  const editorSettings = loadEditorSettings()
  const options = toMonacoOptions(editorSettings)
  
  editor = (monaco as any).editor.create(editorContainer.value, {
    value: sourceCode.value,
    language: 'python',
    readOnly: false,
    ...options
  })

  editor.onDidChangeModelContent(() => {
    sourceCode.value = editor.getValue()
    // 延迟保存，避免频繁写入localStorage
    if (saveTimer) {
      clearTimeout(saveTimer)
    }
    saveTimer = setTimeout(() => {
      saveCodeToStorage()
    }, 1000) // 1秒后保存
  })
  
  const cleanup = listenToSettingsChanges((settings: EditorSettings) => {
    updateEditorOptions(editor, settings)
  })
  
  ;(editor as any)._settingsCleanup = cleanup
}

const handleFileChange = async (file: any) => {
  try {
    const rawFile = file.raw
    
    // 验证文件类型
    const allowedExtensions = ['.py', '.js', '.java', '.go', '.rs', '.ts', '.cpp', '.c', '.jsx', '.tsx']
    const fileName = rawFile.name.toLowerCase()
    const isValidType = allowedExtensions.some(ext => fileName.endsWith(ext))
    
    if (!isValidType) {
      ElMessage.error(`不支持的文件类型。支持的类型：${allowedExtensions.join(', ')}`)
      return
    }
    
    // 验证文件大小（限制为5MB）
    const maxSize = 5 * 1024 * 1024 // 5MB
    if (rawFile.size > maxSize) {
      ElMessage.error('文件大小不能超过5MB')
      return
    }
    
    // 读取文件内容
    const reader = new FileReader()
    reader.onload = (e) => {
      const content = e.target?.result as string
      if (content) {
        if (editor) {
          editor.setValue(content)
        }
        sourceCode.value = content
        
        // 根据文件扩展名自动设置语言
        const ext = fileName.substring(fileName.lastIndexOf('.'))
        const languageMap: Record<string, string> = {
          '.py': 'python',
          '.js': 'javascript',
          '.jsx': 'javascript',
          '.ts': 'typescript',
          '.tsx': 'typescript',
          '.java': 'java',
          '.go': 'go',
          '.rs': 'rust',
          '.cpp': 'cpp',
          '.c': 'c'
        }
        
        const detectedLanguage = languageMap[ext]
        if (detectedLanguage) {
          analysisConfig.value.language = detectedLanguage
          // 更新编辑器语言
          if (editor) {
            const model = editor.getModel()
            if (model) {
              (monaco as any).editor.setModelLanguage(model, detectedLanguage)
            }
          }
        }
        
        // 保存到localStorage
        saveCodeToStorage()
        
        ElMessage.success(`文件加载成功: ${rawFile.name}`)
      }
    }
    
    reader.onerror = () => {
      ElMessage.error('文件读取失败')
    }
    
    reader.readAsText(rawFile)
    
  } catch (error) {
    console.error('File handling error:', error)
    ElMessage.error(`文件处理失败: ${error instanceof Error ? error.message : '未知错误'}`)
  }
}

const analyzeCode = async () => {
  if (!sourceCode.value.trim()) {
    ElMessage.warning('请先输入或上传代码')
    return
  }
  
  if (analysisConfig.value.analysisTypes.length === 0) {
    ElMessage.warning('请至少选择一种分析类型')
    return
  }

  // 清空之前的结果
  analysisResult.value = null
  analyzing.value = true
  
  try {
    console.log('Creating analysis task with config:', analysisConfig.value)
    
    const task = await taskStore.createTask(
      'analysis',
      '代码分析任务',
      `分析${analysisConfig.value.language}代码质量和潜在问题`,
      {
        code: sourceCode.value,
        language: analysisConfig.value.language,
        analysis_types: analysisConfig.value.analysisTypes
      }
    )
    
    currentTaskId.value = task.id
    console.log('Task created:', task.id)
    ElMessage.success('代码分析任务已创建，正在处理...')
    
    await waitForTaskCompletion(task.id)
    
  } catch (error) {
    console.error('Analysis error:', error)
    ElMessage.error(`分析失败: ${error instanceof Error ? error.message : '未知错误'}`)
  } finally {
    analyzing.value = false
  }
}

const waitForTaskCompletion = async (taskId: string) => {
  const maxAttempts = 60
  let attempts = 0
  
  while (attempts < maxAttempts) {
    try {
      const task = await taskStore.getTask(taskId)
      
      if (task.status === 'completed') {
        const result = await taskStore.getTaskResult(taskId)
        
        console.log('Task result:', result)
        
        // 处理不同的结果格式
        let analysisData = null
        
        if (result.result) {
          // 尝试多种可能的数据结构
          if (result.result.result) {
            analysisData = result.result.result
          } else if (result.result.analysis_id) {
            // 直接是分析结果
            analysisData = result.result
          } else {
            analysisData = result.result
          }
        }
        
        if (analysisData) {
          // 转换Python Agent返回的格式到前端期望的格式
          analysisResult.value = {
            quality_score: Math.round((1 - (analysisData.complexity_score || 0)) * 100),
            complexity_score: analysisData.complexity_score || 0,
            dependencies: analysisData.dependencies || [],
            issues: [
              ...(analysisData.security_issues || []).map((issue: any) => ({
                severity: issue.severity || 'warning',
                line: issue.line || 0,
                message: issue.message || issue.type || '安全问题',
                type: issue.type || 'security'
              })),
              ...(analysisData.compatibility_issues || []).map((issue: any) => ({
                severity: issue.severity || 'info',
                line: issue.line || 0,
                message: issue.message || issue.type || '兼容性问题',
                type: issue.type || 'compatibility'
              }))
            ],
            suggestions: (analysisData.recommendations || []).map((rec: string) => ({
              line: 0,
              title: '优化建议',
              description: rec
            })),
            business_logic_summary: analysisData.business_logic_summary || '',
            analysis_id: analysisData.analysis_id || taskId
          }
          
          console.log('Processed analysis result:', analysisResult.value)
          saveCodeToStorage() // 保存结果到localStorage
          ElMessage.success('代码分析完成')
          return
        } else {
          console.error('Invalid result structure:', result)
          ElMessage.error('分析结果格式错误')
          return
        }
      } else if (task.status === 'failed') {
        ElMessage.error(`分析失败: ${task.error || '未知错误'}`)
        return
      }
      
      await new Promise(resolve => setTimeout(resolve, 5000))
      attempts++
      
    } catch (error) {
      console.error('Error checking task status:', error)
      await new Promise(resolve => setTimeout(resolve, 5000))
      attempts++
    }
  }
  
  ElMessage.error('分析超时，请稍后重试')
}

const downloadReport = () => {
  if (!analysisResult.value) {
    ElMessage.warning('没有可下载的分析报告')
    return
  }
  
  const report = JSON.stringify(analysisResult.value, null, 2)
  const blob = new Blob([report], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'analysis_report.json'
  link.click()
  URL.revokeObjectURL(url)
  ElMessage.success('报告下载成功')
}

const getScoreColor = (score: number) => {
  if (score >= 80) return '#66bb6a'
  if (score >= 60) return '#ffb74d'
  return '#ef5350'
}

const getIssueType = (severity: string) => {
  const typeMap: Record<string, string> = {
    'error': 'danger',
    'warning': 'warning',
    'info': 'info'
  }
  return typeMap[severity] || 'info'
}
</script>

<style scoped>
.code-analysis-container {
  padding: 0;
}

.code-card,
.result-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
  height: 600px;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.editor-container {
  flex: 1;
  min-height: 500px;
}

.monaco-editor-container {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}

.result-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.quality-score {
  text-align: center;
  margin-bottom: 32px;
}

.percentage-value {
  display: block;
  font-size: 28px;
  font-weight: bold;
  color: var(--text-primary);
}

.percentage-label {
  display: block;
  font-size: 14px;
  color: var(--text-secondary);
  margin-top: 8px;
}

.metrics-section {
  margin-bottom: 24px;
}

.summary-section {
  margin-bottom: 24px;
}

.dependencies-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin: 24px 0 16px 0;
  color: var(--text-primary);
}

.issue-item {
  padding: 8px;
}

.issue-message {
  margin: 8px 0 0 0;
  color: var(--text-secondary);
  line-height: 1.6;
}

.config-section {
  margin-top: 24px;
}

.config-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

/* 修复Element Plus组件背景色 */
:deep(.el-descriptions) {
  background-color: transparent;
}

:deep(.el-descriptions__body) {
  background-color: var(--bg-card);
}

:deep(.el-descriptions__label) {
  background-color: var(--bg-secondary);
  color: var(--text-secondary);
  font-weight: 500;
}

:deep(.el-descriptions__content) {
  background-color: var(--bg-card);
  color: var(--text-primary);
  font-weight: 600;
}

:deep(.el-alert) {
  background-color: var(--bg-secondary);
  border-color: var(--border-color);
}

:deep(.el-alert__title) {
  color: var(--text-primary);
}

:deep(.el-timeline-item__wrapper) {
  padding-left: 20px;
}

:deep(.el-timeline-item__timestamp) {
  color: var(--text-secondary);
}

:deep(.el-collapse) {
  border-color: var(--border-color);
}

:deep(.el-collapse-item__header) {
  background-color: var(--bg-card);
  color: var(--text-primary);
  border-color: var(--border-color);
}

:deep(.el-collapse-item__wrap) {
  background-color: var(--bg-card);
  border-color: var(--border-color);
}

:deep(.el-collapse-item__content) {
  color: var(--text-secondary);
  background-color: var(--bg-card);
}

.suggestion-description {
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 12px;
}

.code-example {
  margin-top: 12px;
  padding: 12px;
  background-color: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.code-example pre {
  margin: 8px 0 0 0;
  padding: 8px;
  background-color: var(--bg-primary);
  border-radius: 4px;
  overflow-x: auto;
}

.code-example code {
  color: var(--text-primary);
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .el-col {
    margin-bottom: 16px;
  }
  
  .code-card,
  .result-card {
    height: 400px;
  }
}
</style>

<style>
/* 全局样式：修复语言选择下拉框显示问题 */
.language-select-dropdown {
  z-index: 9999 !important;
  max-height: 300px;
}

.language-select-dropdown .el-select-dropdown__item {
  background-color: var(--bg-card);
  color: var(--text-primary);
}

.language-select-dropdown .el-select-dropdown__item:hover {
  background-color: var(--bg-secondary);
}

.language-select-dropdown .el-select-dropdown__item.selected {
  color: var(--primary-color);
  font-weight: 600;
}
</style>
