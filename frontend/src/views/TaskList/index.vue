<template>
  <div class="task-list-container">
    <el-card class="task-list-card">
      <template #header>
        <div class="card-header">
          <h3 class="card-title">任务列表</h3>
          <div class="header-actions">
            <el-button :icon="Refresh" @click="refreshTasks" :loading="loading">
              刷新
            </el-button>
            <el-button :icon="Delete" @click="clearAllTasks" type="danger">
              清空全部
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选器 -->
      <div class="filters">
        <el-radio-group v-model="filterStatus" @change="handleFilterChange">
          <el-radio-button label="all">全部 ({{ tasks.length }})</el-radio-button>
          <el-radio-button label="running">运行中 ({{ runningCount }})</el-radio-button>
          <el-radio-button label="completed">已完成 ({{ completedCount }})</el-radio-button>
          <el-radio-button label="failed">失败 ({{ failedCount }})</el-radio-button>
        </el-radio-group>

        <el-input
          v-model="searchQuery"
          placeholder="搜索任务名称..."
          :prefix-icon="Search"
          clearable
          style="width: 300px; margin-left: 16px"
        />
      </div>

      <!-- 任务表格 -->
      <el-table
        :data="filteredTasks"
        style="width: 100%; margin-top: 20px"
        v-loading="loading"
      >
        <el-table-column prop="name" label="任务名称" min-width="200">
          <template #default="{ row }">
            <div class="task-name">
              <el-icon :size="16" style="margin-right: 8px">
                <component :is="getTaskIcon(row.type)" />
              </el-icon>
              {{ row.name }}
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getTaskTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="progress" label="进度" width="150">
          <template #default="{ row }">
            <el-progress
              :percentage="row.progress || 0"
              :status="getProgressStatus(row.status)"
              :stroke-width="8"
            />
          </template>
        </el-table-column>

        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              text
              type="primary"
              @click="viewTaskDetail(row)"
              :icon="Document"
            >
              查看
            </el-button>
            <el-button
              text
              type="danger"
              @click="deleteTask(row)"
              :icon="Delete"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 空状态 -->
      <el-empty v-if="filteredTasks.length === 0 && !loading" description="暂无任务记录" />
    </el-card>

    <!-- 任务详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="任务详情"
      width="800px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedTask" class="task-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务ID">
            {{ selectedTask.id }}
          </el-descriptions-item>
          <el-descriptions-item label="任务名称">
            {{ selectedTask.name }}
          </el-descriptions-item>
          <el-descriptions-item label="任务类型">
            {{ getTaskTypeText(selectedTask.type) }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(selectedTask.status)">
              {{ getStatusText(selectedTask.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="进度">
            {{ selectedTask.progress || 0 }}%
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(selectedTask.createdAt) }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="selectedTask.error" class="error-section">
          <h4>错误信息</h4>
          <el-alert type="error" :closable="false">
            {{ selectedTask.error }}
          </el-alert>
        </div>

        <div v-if="selectedTask.result" class="result-section">
          <h4>任务结果</h4>
          <pre class="result-content">{{ JSON.stringify(selectedTask.result, null, 2) }}</pre>
        </div>
      </div>

      <template #footer>
        <el-button @click="showDetailDialog = false">关闭</el-button>
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
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Delete,
  DocumentChecked,
  Coin,
  Document
} from '@element-plus/icons-vue'
import { useTaskStore } from '@/stores/task'
import type { Task } from '@/types/task'

const router = useRouter()
const taskStore = useTaskStore()

const loading = ref(false)
const filterStatus = ref('all')
const searchQuery = ref('')
const showDetailDialog = ref(false)
const selectedTask = ref<Task | null>(null)

const tasks = computed(() => taskStore.tasks)

const runningCount = computed(() =>
  tasks.value.filter((t: Task) => t.status === 'running').length
)

const completedCount = computed(() =>
  tasks.value.filter((t: Task) => t.status === 'completed').length
)

const failedCount = computed(() =>
  tasks.value.filter((t: Task) => t.status === 'failed').length
)

const filteredTasks = computed(() => {
  let result = tasks.value

  // 按状态筛选
  if (filterStatus.value !== 'all') {
    result = result.filter((t: Task) => t.status === filterStatus.value)
  }

  // 按名称搜索
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter((t: Task) => t.name.toLowerCase().includes(query))
  }

  return result
})

onMounted(() => {
  refreshTasks()
})

const refreshTasks = async () => {
  loading.value = true
  try {
    // 任务已经从localStorage加载，这里只是刷新显示
    await new Promise((resolve) => setTimeout(resolve, 500))
  } catch (error) {
    ElMessage.error('刷新任务列表失败')
  } finally {
    loading.value = false
  }
}

const handleFilterChange = () => {
  // 筛选变化时的处理
}

const viewTaskDetail = (task: Task) => {
  selectedTask.value = task
  showDetailDialog.value = true
}

const deleteTask = async (task: Task) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除任务"${task.name}"吗？此操作不可撤销。`,
      '删除任务',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const index = taskStore.tasks.findIndex((t: Task) => t.id === task.id)
    if (index !== -1) {
      taskStore.tasks.splice(index, 1)
      // 保存到localStorage
      localStorage.setItem('codesage_tasks', JSON.stringify(taskStore.tasks))
      ElMessage.success('任务已删除')
    }
  } catch {
    // 用户取消
  }
}

const clearAllTasks = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清空所有任务吗？此操作不可撤销。',
      '清空任务',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    taskStore.tasks.length = 0
    localStorage.setItem('codesage_tasks', JSON.stringify([]))
    ElMessage.success('所有任务已清空')
  } catch {
    // 用户取消
  }
}

const goToAnalysis = () => {
  showDetailDialog.value = false
  router.push('/analysis')
}

const goToConversion = () => {
  showDetailDialog.value = false
  router.push('/conversion')
}

const goToTestGeneration = () => {
  showDetailDialog.value = false
  router.push('/test-generation')
}

const getTaskIcon = (type: string) => {
  const iconMap: Record<string, any> = {
    analysis: Search,
    convert: Refresh,
    test: DocumentChecked,
    git: Coin
  }
  return iconMap[type] || Search
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

const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    pending: 'info',
    running: 'warning',
    completed: 'success',
    failed: 'danger',
    cancelled: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待处理',
    running: '运行中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

const getProgressStatus = (status: string) => {
  if (status === 'completed') return 'success'
  if (status === 'failed') return 'exception'
  return undefined
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.task-list-container {
  padding: 0;
}

.task-list-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  gap: 12px;
}

.filters {
  display: flex;
  align-items: center;
  margin-top: 20px;
}

.task-name {
  display: flex;
  align-items: center;
  color: var(--text-primary);
}

.task-detail {
  padding: 20px 0;
}

.error-section,
.result-section {
  margin-top: 24px;
}

.error-section h4,
.result-section h4 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  color: var(--text-primary);
}

.result-content {
  background: var(--bg-secondary);
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  color: var(--text-primary);
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .filters {
    flex-direction: column;
    align-items: stretch;
  }

  .filters .el-input {
    margin-left: 0;
    margin-top: 12px;
    width: 100% !important;
  }
}
</style>
