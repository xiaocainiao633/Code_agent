<template>
  <div class="git-analysis-container">
    <!-- 仓库输入 -->
    <el-row :gutter="24" class="repo-input-section">
      <el-col :span="24">
        <el-card class="repo-card">
          <template #header>
            <h3 class="card-title">Git 仓库分析</h3>
          </template>
          
          <el-form :model="repoForm" inline>
            <el-form-item label="仓库地址">
              <el-input
                v-model="repoForm.url"
                placeholder="输入 Git 仓库地址"
                style="width: 400px"
              >
                <template #prefix>
                  <el-icon><Link /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item label="分支">
              <el-select 
                v-model="repoForm.branch" 
                placeholder="选择分支" 
                @change="handleBranchChange"
                style="width: 150px"
              >
                <el-option label="main" value="main" />
                <el-option label="master" value="master" />
                <el-option label="develop" value="develop" />
              </el-select>
            </el-form-item>
            
            <el-form-item>
              <el-button 
                type="primary" 
                :icon="Search" 
                @click="analyzeRepository"
                :loading="analyzing"
              >
                分析仓库
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 分析结果 -->
    <el-row :gutter="24" class="results-section" v-if="analysisResults">
      <el-col :span="8">
        <el-card class="stat-card">
          <template #header>
            <h4 class="stat-title">仓库统计</h4>
          </template>
          
          <div class="stat-content">
            <div class="stat-item">
              <span class="stat-label">总提交数</span>
              <span class="stat-value">{{ analysisResults.totalCommits }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">贡献者数</span>
              <span class="stat-value">{{ analysisResults.contributors.length }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">文件数</span>
              <span class="stat-value">{{ analysisResults.totalFiles }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">代码行数</span>
              <span class="stat-value">{{ analysisResults.totalLines }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="16">
        <el-card class="chart-card">
          <template #header>
            <h4 class="stat-title">提交历史趋势</h4>
          </template>
          
          <div class="chart-container">
            <div ref="commitChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 贡献者排行 -->
    <el-row :gutter="24" class="contributors-section" v-if="analysisResults">
      <el-col :span="24">
        <el-card class="contributors-card">
          <template #header>
            <h4 class="stat-title">贡献者排行</h4>
          </template>
          
          <el-table :data="contributorRankings" style="width: 100%">
            <el-table-column prop="rank" label="排名" width="80" />
            <el-table-column prop="name" label="贡献者" />
            <el-table-column prop="commits" label="提交数" width="100" />
            <el-table-column prop="additions" label="新增行数" width="100" />
            <el-table-column prop="deletions" label="删除行数" width="100" />
            <el-table-column label="活跃度" width="200">
              <template #default="{ row }">
                <el-progress :percentage="row.activity" :stroke-width="6" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 文件变更分析 -->
    <el-row :gutter="24" class="files-section" v-if="analysisResults">
      <el-col :span="24">
        <el-card class="files-card">
          <template #header>
            <h4 class="stat-title">文件变更分析</h4>
          </template>
          
          <el-table :data="fileChanges" style="width: 100%">
            <el-table-column prop="fileName" label="文件名" />
            <el-table-column prop="changeCount" label="变更次数" width="120" />
            <el-table-column prop="lastModified" label="最后修改" width="180" />
            <el-table-column prop="contributors" label="贡献者" width="200">
              <template #default="{ row }">
                <el-tag v-for="contributor in row.contributors.slice(0, 3)" :key="contributor" size="small">
                  {{ contributor }}
                </el-tag>
                <el-tag v-if="row.contributors.length > 3" size="small">
                  +{{ row.contributors.length - 3 }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button text type="primary" @click="viewFileHistory(row)">
                  查看历史
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Link } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { useTaskStore } from '@/stores/task'
import { gitAPI } from '../../services/api'

const repoForm = ref({
  url: 'https://github.com/example/repo.git',
  branch: 'main'
})

const analyzing = ref(false)
const analysisResults = ref<{
  totalCommits: number
  contributors: Array<any>
  totalFiles: number
  totalLines: number
  commitHistory: Array<any>
} | null>(null)

const commitChart = ref<HTMLElement>()
const taskStore = useTaskStore()
const currentTaskId = ref<string | null>(null)

onMounted(() => {
  // 初始化图表
  if (commitChart.value) {
    initChart()
  }
})

const initChart = () => {
  const chart = echarts.init(commitChart.value!)
  const option = {
    backgroundColor: 'transparent',
    title: {
      text: '提交历史',
      textStyle: {
        color: '#ffffff'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月'],
      axisLine: {
        lineStyle: {
          color: '#333333'
        }
      },
      axisLabel: {
        color: '#b0b0b0'
      }
    },
    yAxis: {
      type: 'value',
      axisLine: {
        lineStyle: {
          color: '#333333'
        }
      },
      axisLabel: {
        color: '#b0b0b0'
      },
      splitLine: {
        lineStyle: {
          color: '#333333'
        }
      }
    },
    series: [{
      data: [120, 200, 150, 80, 70, 110],
      type: 'bar',
      itemStyle: {
        color: '#1e88e5'
      }
    }]
  }
  
  chart.setOption(option)
}

const analyzeRepository = async () => {
  if (!repoForm.value.url.trim()) {
    ElMessage.warning('请输入仓库地址')
    return
  }

  analyzing.value = true
  try {
    // 第一步：克隆仓库
    ElMessage.info('正在克隆仓库...')
    const cloneResult = await gitAPI.cloneRepository(
      repoForm.value.url,
      `./temp/repo_${Date.now()}`
    )
    
    const cloneTaskId = cloneResult.task_id
    
    // 等待克隆完成
    await waitForTaskCompletion(cloneTaskId, '克隆')
    
    // 第二步：分析仓库
    ElMessage.info('正在分析仓库...')
    const repoPath = `./temp/repo_${Date.now()}`
    const analyzeResult = await gitAPI.analyzeRepository(
      repoPath,
      repoForm.value.url,
      false // 不需要重新克隆
    )
    
    const analyzeTaskId = analyzeResult.task_id
    currentTaskId.value = analyzeTaskId
    
    // 等待分析完成
    await waitForTaskCompletion(analyzeTaskId, '分析')
    
  } catch (error) {
    ElMessage.error(`分析失败: ${error instanceof Error ? error.message : '未知错误'}`)
  } finally {
    analyzing.value = false
  }
}

// 等待任务完成
const waitForTaskCompletion = async (taskId: string, taskType: string) => {
  const maxAttempts = 60 // 最多等待5分钟
  let attempts = 0
  
  while (attempts < maxAttempts) {
    try {
      // 获取任务详情
      const task = await taskStore.getTask(taskId)
      
      if (task.status === 'completed') {
        // 获取任务结果
        const result = await taskStore.getTaskResult(taskId)
        
        if (result.result && result.result.result) {
          // 解析分析结果
          const analysisData = result.result.result
          
          // 更新分析结果
          analysisResults.value = {
            totalCommits: analysisData.total_commits || 0,
            contributors: analysisData.contributors || [],
            totalFiles: analysisData.total_files || 0,
            totalLines: analysisData.total_lines || 0,
            commitHistory: analysisData.commit_history || []
          }
          
          // 更新图表
          if (analysisData.commit_history && commitChart.value) {
            updateCommitChart(analysisData.commit_history)
          }
          
          ElMessage.success(`仓库${taskType}完成`)
          return
        }
      } else if (task.status === 'failed') {
        ElMessage.error(`${taskType}失败: ${task.error || '未知错误'}`)
        return
      }
      
      // 等待1秒后重试
      await new Promise(resolve => setTimeout(resolve, 5000))
      attempts++
      
    } catch (error) {
      console.error(`Error checking ${taskType} task status:`, error)
      await new Promise(resolve => setTimeout(resolve, 5000))
      attempts++
    }
  }
  
  ElMessage.error(`${taskType}超时，请稍后重试`)
}

// 更新提交历史图表
const updateCommitChart = (commitHistory: any[]) => {
  if (!commitChart.value) return
  
  const chart = (echarts as any).getInstanceByDom(commitChart.value) || (echarts as any).init(commitChart.value)
  
  // 处理提交历史数据
  const dates = commitHistory.map(item => item.date || item.month || '未知')
  const counts = commitHistory.map(item => item.count || item.commits || 0)
  
  const option = {
    backgroundColor: 'transparent',
    title: {
      text: '提交历史',
      textStyle: {
        color: '#ffffff'
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: {
        lineStyle: {
          color: '#333333'
        }
      },
      axisLabel: {
        color: '#b0b0b0'
      }
    },
    yAxis: {
      type: 'value',
      axisLine: {
        lineStyle: {
          color: '#333333'
        }
      },
      axisLabel: {
        color: '#b0b0b0'
      },
      splitLine: {
        lineStyle: {
          color: '#333333'
        }
      }
    },
    series: [{
      data: counts,
      type: 'bar',
      itemStyle: {
        color: '#1e88e5'
      }
    }]
  }
  
  chart.setOption(option)
}

const contributorRankings = computed(() => {
  if (!analysisResults.value) return []
  
  return analysisResults.value.contributors.map((contributor, index) => ({
    rank: index + 1,
    name: contributor.name,
    commits: contributor.commits,
    additions: contributor.additions,
    deletions: contributor.deletions,
    activity: Math.round((contributor.commits / analysisResults.value!.totalCommits) * 100)
  }))
})

const fileChanges = computed(() => {
  if (!analysisResults.value) return []
  
  // 从分析结果中获取文件变更数据
  return analysisResults.value.contributors.map((contributor, index) => ({
    fileName: `src/${contributor.name.toLowerCase()}_module.py`, // 模拟文件名
    changeCount: Math.floor(Math.random() * 50) + 10, // 模拟变更次数
    lastModified: new Date(Date.now() - index * 24 * 60 * 60 * 1000).toISOString().split('T')[0], // 模拟日期
    contributors: [contributor.name]
  })).slice(0, 10) // 限制显示数量
})

const viewFileHistory = async (file: any) => {
  if (!currentTaskId.value) {
    ElMessage.warning('请先分析仓库')
    return
  }
  
  try {
    // 获取文件历史
    const result = await gitAPI.getFileHistory(
      `./temp/repo_${Date.now()}`, // 使用当前分析的仓库路径
      file.fileName
    )
    
    const taskId = result.task_id
    ElMessage.info(`正在获取 ${file.fileName} 的历史记录...`)
    
    // 等待任务完成
    await waitForTaskCompletion(taskId, '文件历史')
    
  } catch (error) {
    ElMessage.error(`获取文件历史失败: ${error instanceof Error ? error.message : '未知错误'}`)
  }
}

const handleBranchChange = (branch: string) => {
  ElMessage.success(`已选择分支: ${branch}`)
  // 这里可以添加重新分析的逻辑
  if (analysisResults.value) {
    ElMessage.info('分支已更改，请重新分析仓库以获取最新数据')
  }
}
</script>

<style scoped>
.git-analysis-container {
  padding: 0;
}

.repo-input-section {
  margin-bottom: 24px;
}

.repo-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.results-section {
  margin-bottom: 24px;
}

.stat-card,
.chart-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.stat-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.stat-content {
  padding: 20px 0;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.stat-item:last-child {
  border-bottom: none;
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: var(--primary-color);
}

.chart-container {
  height: 300px;
  padding: 20px;
}

.chart {
  width: 100%;
  height: 100%;
}

.contributors-section,
.files-section {
  margin-bottom: 24px;
}

.contributors-card,
.files-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .el-col {
    margin-bottom: 16px;
  }
  
  .chart-container {
    height: 250px;
  }
}
</style>