<template>
  <div class="dashboard-container">
    <!-- 欢迎区域 -->
    <el-row :gutter="24" class="welcome-section">
      <el-col :span="24">
        <el-card class="welcome-card fade-in">
          <div class="welcome-content">
            <div class="welcome-text">
              <h1 class="welcome-title">
                欢迎使用 <span class="text-gradient">CodeSage</span>
              </h1>
              <p class="welcome-description">
                智能代码重构与现代化助手，帮助您安全、高效地将老旧代码迁移至现代技术栈
              </p>
              <div class="welcome-actions">
                <el-button type="primary" size="large" @click="startAnalysis">
                  <el-icon><Search /></el-icon>
                  开始代码分析
                </el-button>
                <el-button size="large" @click="viewDocumentation">
                  <el-icon><Document /></el-icon>
                  查看文档
                </el-button>
              </div>
            </div>
            <div class="welcome-animation">
              <div class="floating-icon">
                <el-icon size="80" color="var(--primary-color)">
                  <Cpu />
                </el-icon>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快速操作 -->
    <el-row :gutter="24" class="quick-actions-section">
      <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="action in quickActions" :key="action.key">
        <el-card class="action-card" @click="handleQuickAction(action)" hoverable>
          <div class="action-content">
            <div class="action-icon">
              <el-icon size="32" :style="{ color: action.color }">
                <component :is="action.icon" />
              </el-icon>
            </div>
            <h3 class="action-title">{{ action.title }}</h3>
            <p class="action-description">{{ action.description }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近任务 -->
    <el-row :gutter="24" class="recent-tasks-section">
      <el-col :span="24">
        <el-card class="tasks-card">
          <template #header>
            <div class="card-header">
              <h3 class="card-title">最近任务</h3>
              <el-button text @click="viewAllTasks">
                查看全部
                <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
          </template>
          
          <div class="tasks-content">
            <el-table
              :data="recentTasks"
              style="width: 100%"
              :show-header="false"
              v-if="recentTasks.length > 0"
            >
              <el-table-column prop="name" label="任务名称" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="getStatusType(row.status)" size="small">
                    {{ getStatusText(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="createdAt" label="创建时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.createdAt) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120">
                <template #default="{ row }">
                  <el-button 
                    type="primary" 
                    size="small"
                    @click="viewTask(row)"
                  >
                    <el-icon style="margin-right: 4px"><Document /></el-icon>
                    查看报告
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            
            <el-empty v-else description="暂无任务记录" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统状态 -->
    <el-row :gutter="24" class="system-status-section">
      <el-col :xs="24" :sm="12" :md="6" v-for="stat in systemStats" :key="stat.key">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" :style="{ backgroundColor: stat.bgColor }">
              <el-icon size="24" :style="{ color: stat.color }">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <h4 class="stat-title">{{ stat.title }}</h4>
              <p class="stat-value">{{ stat.value }}</p>
              <p class="stat-description">{{ stat.description }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 任务报告对话框 -->
    <el-dialog
      v-model="showReportDialog"
      :title="selectedTask ? `${selectedTask.name} - 报告` : '任务报告'"
      width="900px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedTask" class="task-report">
        <!-- 任务基本信息 -->
        <el-descriptions :column="2" border class="task-info">
          <el-descriptions-item label="任务类型">
            <el-tag size="small">{{ getTaskTypeText(selectedTask.type) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="任务状态">
            <el-tag :type="getStatusType(selectedTask.status)" size="small">
              {{ getStatusText(selectedTask.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(selectedTask.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="进度">
            {{ selectedTask.progress || 0 }}%
          </el-descriptions-item>
        </el-descriptions>

        <!-- 错误信息 -->
        <div v-if="selectedTask.error" class="error-section">
          <h4 class="section-title">错误信息</h4>
          <el-alert type="error" :closable="false" show-icon>
            {{ selectedTask.error }}
          </el-alert>
        </div>

        <!-- 任务结果 -->
        <div v-if="selectedTask.result" class="result-section">
          <h4 class="section-title">任务结果</h4>
          
          <!-- 代码分析结果 -->
          <div v-if="selectedTask.type === 'analysis'" class="analysis-result">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="质量评分">
                <el-tag :type="getScoreType(selectedTask.result.quality_score || 0)">
                  {{ selectedTask.result.quality_score || 0 }}分
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="复杂度分数">
                {{ ((selectedTask.result.complexity_score || 0) * 100).toFixed(1) }}%
              </el-descriptions-item>
              <el-descriptions-item label="依赖数量">
                {{ selectedTask.result.dependencies?.length || 0 }}
              </el-descriptions-item>
              <el-descriptions-item label="问题数量">
                {{ selectedTask.result.issues?.length || 0 }}
              </el-descriptions-item>
            </el-descriptions>
            
            <div v-if="selectedTask.result.business_logic_summary" style="margin-top: 16px">
              <strong>业务逻辑摘要：</strong>
              <p class="summary-text">{{ selectedTask.result.business_logic_summary }}</p>
            </div>
          </div>

          <!-- 代码转换结果 -->
          <div v-if="selectedTask.type === 'convert'" class="conversion-result">
            <el-alert type="success" :closable="false" show-icon>
              代码转换已完成，转换后的代码已保存
            </el-alert>
          </div>

          <!-- 测试生成结果 -->
          <div v-if="selectedTask.type === 'test'" class="test-result">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="测试用例数">
                {{ selectedTask.result.test_cases?.length || 0 }}
              </el-descriptions-item>
              <el-descriptions-item label="覆盖率">
                {{ selectedTask.result.coverage?.total || 0 }}%
              </el-descriptions-item>
            </el-descriptions>
          </div>

          <!-- 原始JSON数据（折叠） -->
          <el-collapse style="margin-top: 16px">
            <el-collapse-item title="查看原始数据" name="raw">
              <pre class="result-json">{{ JSON.stringify(selectedTask.result, null, 2) }}</pre>
            </el-collapse-item>
          </el-collapse>
        </div>

        <!-- 无结果提示 -->
        <div v-if="!selectedTask.result && !selectedTask.error && selectedTask.status === 'completed'" class="no-result">
          <el-empty description="暂无结果数据" />
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showReportDialog = false">关闭</el-button>
          <el-button
            v-if="selectedTask && selectedTask.type === 'analysis'"
            type="primary"
            @click="goToAnalysis"
          >
            前往代码分析
          </el-button>
          <el-button
            v-if="selectedTask && selectedTask.type === 'convert'"
            type="primary"
            @click="goToConversion"
          >
            前往代码转换
          </el-button>
          <el-button
            v-if="selectedTask && selectedTask.type === 'test'"
            type="primary"
            @click="goToTestGeneration"
          >
            前往测试生成
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Search,
  Document,
  Cpu,
  ArrowRight,
  Refresh,
  DocumentChecked,
  Coin,
  TrendCharts,
  Clock,
  CircleCheck
} from '@element-plus/icons-vue'
import { useTaskStore } from '@/stores/task'
import type { Task } from '@/types/task'

const router = useRouter()
const taskStore = useTaskStore()

const recentTasks = ref<Task[]>([])
const showReportDialog = ref(false)
const selectedTask = ref<Task | null>(null)
const systemStats = ref([
  {
    key: 'total',
    title: '总任务数',
    value: '0',
    description: '累计处理任务',
    icon: TrendCharts,
    color: '#1e88e5',
    bgColor: 'rgba(30, 136, 229, 0.1)'
  },
  {
    key: 'running',
    title: '运行中',
    value: '0',
    description: '当前运行任务',
    icon: Clock,
    color: '#ffb74d',
    bgColor: 'rgba(255, 183, 77, 0.1)'
  },
  {
    key: 'completed',
    title: '已完成',
    value: '0',
    description: '成功完成任务',
    icon: CircleCheck,
    color: '#66bb6a',
    bgColor: 'rgba(102, 187, 106, 0.1)'
  },
  {
    key: 'success',
    title: '成功率',
    value: '0%',
    description: '任务成功率',
    icon: CircleCheck,
    color: '#26c6da',
    bgColor: 'rgba(38, 198, 218, 0.1)'
  }
])

const quickActions = ref([
  {
    key: 'analysis',
    title: '代码分析',
    description: '分析代码质量和潜在问题',
    icon: Search,
    color: '#1e88e5',
    route: '/analysis'
  },
  {
    key: 'conversion',
    title: '代码转换',
    description: '将代码转换为现代语法',
    icon: Refresh,
    color: '#26c6da',
    route: '/conversion'
  },
  {
    key: 'test',
    title: '测试生成',
    description: '自动生成单元测试',
    icon: DocumentChecked,
    color: '#66bb6a',
    route: '/test-generation'
  },
  {
    key: 'git',
    title: 'Git分析',
    description: '分析Git仓库历史',
    icon: Coin,
    color: '#ffb74d',
    route: '/git-analysis'
  }
])

onMounted(() => {
  loadRecentTasks()
  updateSystemStats()
})

const loadRecentTasks = () => {
  // 从store获取最近任务
  recentTasks.value = taskStore.tasks.slice(0, 5)
}

const updateSystemStats = () => {
  const tasks = taskStore.tasks
  const stats = systemStats.value as any[]
  if (stats && stats.length >= 4) {
    stats[0].value = tasks.length.toString()
    stats[1].value = tasks.filter((t: Task) => t.status === 'running').length.toString()
    stats[2].value = tasks.filter((t: Task) => t.status === 'completed').length.toString()
    
    const completed = tasks.filter((t: Task) => t.status === 'completed').length
    const total = tasks.length
    stats[3].value = total > 0 ? `${Math.round((completed / total) * 100)}%` : '0%'
  }
}

const startAnalysis = () => {
  router.push('/analysis')
}

const viewDocumentation = () => {
  ElMessage.info('文档功能开发中...')
}

const handleQuickAction = (action: any) => {
  router.push(action.route)
}

const viewAllTasks = () => {
  router.push('/tasks')
}

const viewTask = (task: Task) => {
  selectedTask.value = task
  showReportDialog.value = true
}

const goToAnalysis = () => {
  showReportDialog.value = false
  router.push('/analysis')
}

const goToConversion = () => {
  showReportDialog.value = false
  router.push('/conversion')
}

const goToTestGeneration = () => {
  showReportDialog.value = false
  router.push('/test-generation')
}

const getTaskTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    analysis: '代码分析',
    convert: '代码转换',
    test: '测试生成',
    git: 'Git分析'
  }
  return typeMap[type] || type
}

const getScoreType = (score: number) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    'pending': 'info',
    'running': 'warning',
    'completed': 'success',
    'failed': 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    'pending': '待处理',
    'running': '运行中',
    'completed': '已完成',
    'failed': '失败'
  }
  return statusMap[status] || status
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.dashboard-container {
  padding: 0;
}

.welcome-section {
  margin-bottom: 24px;
}

.welcome-card {
  background: linear-gradient(135deg, var(--bg-card) 0%, var(--bg-secondary) 100%);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 40px;
}

.welcome-text {
  flex: 1;
}

.welcome-title {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 16px;
  color: var(--text-primary);
}

.welcome-description {
  font-size: 16px;
  color: var(--text-secondary);
  margin-bottom: 24px;
  line-height: 1.6;
}

.welcome-actions {
  display: flex;
  gap: 16px;
}

.welcome-animation {
  margin-left: 40px;
}

.floating-icon {
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-10px);
  }
}

.quick-actions-section {
  margin-bottom: 24px;
}

.action-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
}

.action-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  border-color: var(--primary-color);
}

.action-content {
  text-align: center;
  padding: 24px;
}

.action-icon {
  margin-bottom: 16px;
}

.action-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.action-description {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.recent-tasks-section {
  margin-bottom: 24px;
}

.tasks-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
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

.tasks-content {
  padding: 20px;
}

.system-status-section {
  margin-bottom: 24px;
}

.stat-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 20px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
}

.stat-info {
  flex: 1;
}

.stat-title {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 4px 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 4px 0;
}

.stat-description {
  font-size: 12px;
  color: var(--text-tertiary);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .welcome-content {
    flex-direction: column;
    text-align: center;
    padding: 24px;
  }
  
  .welcome-animation {
    margin-left: 0;
    margin-top: 24px;
  }
  
  .welcome-title {
    font-size: 28px;
  }
  
  .welcome-actions {
    justify-content: center;
  }
}

/* 任务报告对话框样式 */
.task-report {
  padding: 20px 0;
}

.task-info {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin: 24px 0 12px 0;
  color: var(--text-primary);
}

.error-section {
  margin-top: 24px;
}

.result-section {
  margin-top: 24px;
}

.summary-text {
  margin: 8px 0;
  padding: 12px;
  background: var(--bg-secondary);
  border-radius: 8px;
  color: var(--text-primary);
  line-height: 1.6;
}

.result-json {
  background: var(--bg-secondary);
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  color: var(--text-primary);
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  max-height: 400px;
  overflow-y: auto;
}

.no-result {
  margin-top: 24px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
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
</style>