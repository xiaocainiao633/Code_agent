<template>
  <div class="test-generation-container">
    <el-row :gutter="24">
      <!-- 左侧源代码 -->
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
                <el-button type="primary" :icon="Upload" size="small"
                  >上传文件</el-button
                >
              </el-upload>
            </div>
          </template>

          <div class="editor-container">
            <div
              ref="sourceEditorContainer"
              class="monaco-editor-container"
            ></div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧生成的测试代码 -->
      <el-col :span="12">
        <el-card class="code-card">
          <template #header>
            <div class="card-header">
              <h3 class="card-title">生成的测试代码</h3>
              <div class="header-actions">
                <el-button
                  type="success"
                  :icon="Download"
                  size="small"
                  @click="downloadTests"
                >
                  下载测试
                </el-button>
                <el-button :icon="CopyDocument" size="small" @click="copyTests">
                  复制代码
                </el-button>
              </div>
            </div>
          </template>

          <div class="editor-container">
            <div
              ref="testEditorContainer"
              class="monaco-editor-container"
            ></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 测试配置 -->
    <el-row :gutter="24" class="config-section">
      <el-col :span="24">
        <el-card class="config-card">
          <template #header>
            <h3 class="card-title">测试配置</h3>
          </template>

          <el-form :model="testConfig" label-position="left" inline>
            <el-form-item label="测试框架">
              <el-select
                v-model="testConfig.framework"
                placeholder="选择测试框架"
                style="width: 150px"
              >
                <el-option label="pytest" value="pytest" />
                <el-option label="unittest" value="unittest" />
                <el-option label="Jest" value="jest" />
                <el-option label="Mocha" value="mocha" />
              </el-select>
            </el-form-item>

            <el-form-item label="测试类型">
              <el-checkbox-group v-model="testConfig.testTypes">
                <el-checkbox label="unit">单元测试</el-checkbox>
                <el-checkbox label="integration">集成测试</el-checkbox>
                <el-checkbox label="functional">功能测试</el-checkbox>
              </el-checkbox-group>
            </el-form-item>

            <el-form-item label="覆盖率目标">
              <el-slider
                v-model="testConfig.coverageTarget"
                :min="50"
                :max="100"
                :step="10"
                show-input
                style="width: 200px"
              />
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                :icon="MagicStick"
                @click="generateTests"
                :loading="generating"
                size="large"
              >
                生成测试
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>

    <!-- 测试覆盖率 -->
    <el-row :gutter="24" class="coverage-section">
      <el-col :span="12">
        <el-card class="coverage-card">
          <template #header>
            <h3 class="card-title">测试覆盖率</h3>
          </template>

          <div class="coverage-content" v-if="coverageData">
            <el-progress
              type="dashboard"
              :percentage="coverageData.total"
              :color="getCoverageColor(coverageData.total)"
            >
              <template #default="{ percentage }">
                <span class="percentage-value">{{ percentage }}%</span>
                <span class="percentage-label">总覆盖率</span>
              </template>
            </el-progress>

            <div class="coverage-details">
              <div class="coverage-item">
                <span class="coverage-label">语句覆盖</span>
                <el-progress
                  :percentage="coverageData.statements"
                  :stroke-width="8"
                />
              </div>
              <div class="coverage-item">
                <span class="coverage-label">分支覆盖</span>
                <el-progress
                  :percentage="coverageData.branches"
                  :stroke-width="8"
                />
              </div>
              <div class="coverage-item">
                <span class="coverage-label">函数覆盖</span>
                <el-progress
                  :percentage="coverageData.functions"
                  :stroke-width="8"
                />
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无覆盖率数据" />
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="test-cases-card">
          <template #header>
            <h3 class="card-title">测试用例</h3>
          </template>

          <div class="test-cases-content" v-if="testCases.length > 0">
            <el-timeline>
              <el-timeline-item
                v-for="(testCase, index) in testCases"
                :key="index"
                :timestamp="`用例 ${index + 1}`"
                :type="getTestCaseType(testCase.status)"
              >
                <div class="test-case-item">
                  <h4 class="test-case-title">{{ testCase.name }}</h4>
                  <p class="test-case-description">
                    {{ testCase.description }}
                  </p>
                  <el-tag size="small" :type="getTestCaseType(testCase.status)">
                    {{ getTestCaseStatus(testCase.status) }}
                  </el-tag>
                </div>
              </el-timeline-item>
            </el-timeline>
          </div>
          <el-empty v-else description="暂无测试用例" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { ElMessage } from "element-plus";
import {
  Upload,
  Download,
  CopyDocument,
  MagicStick,
} from "@element-plus/icons-vue";
import * as monaco from "monaco-editor";
import { useTaskStore } from "@/stores/task";
import {
  loadEditorSettings,
  toMonacoOptions,
  updateEditorOptions,
  listenToSettingsChanges,
  type EditorSettings,
} from "../../utils/editorConfig";

const sourceEditorContainer = ref<HTMLElement>();
const testEditorContainer = ref<HTMLElement>();
let sourceEditor: any;
let testEditor: any;
let saveTimer: any = null;
const taskStore = useTaskStore();
const currentTaskId = ref<string | null>(null);

// 代码自动保存的localStorage key
const CODE_STORAGE_KEY = "test_generation_source_code";
const FRAMEWORK_STORAGE_KEY = "test_generation_framework";
const TEST_CODE_STORAGE_KEY = "test_generation_test_code";
const COVERAGE_DATA_STORAGE_KEY = "test_generation_coverage_data";
const TEST_CASES_STORAGE_KEY = "test_generation_test_cases";

// 从localStorage加载保存的代码（如果启用了自动保存）
const loadSavedCode = () => {
  const allSettings = JSON.parse(
    localStorage.getItem("codesage-settings") || "{}"
  );
  const settings = allSettings.general || {};
  if (settings.autoSave) {
    const savedCode = localStorage.getItem(CODE_STORAGE_KEY);
    const savedFramework = localStorage.getItem(FRAMEWORK_STORAGE_KEY);
    const savedTestCode = localStorage.getItem(TEST_CODE_STORAGE_KEY);
    const savedCoverage = localStorage.getItem(COVERAGE_DATA_STORAGE_KEY);
    const savedTestCases = localStorage.getItem(TEST_CASES_STORAGE_KEY);

    if (savedCode) {
      sourceCode.value = savedCode;
    }
    if (savedFramework) {
      testConfig.value.framework = savedFramework;
    }
    if (savedTestCode) {
      testCode.value = savedTestCode;
    }
    if (savedCoverage) {
      try {
        coverageData.value = JSON.parse(savedCoverage);
      } catch (error) {
        console.error('Failed to parse saved coverage:', error);
      }
    }
    if (savedTestCases) {
      try {
        testCases.value = JSON.parse(savedTestCases);
      } catch (error) {
        console.error('Failed to parse saved test cases:', error);
      }
    }
  }
};

// 保存代码到localStorage（如果启用了自动保存）
const saveCodeToStorage = () => {
  const allSettings = JSON.parse(
    localStorage.getItem("codesage-settings") || "{}"
  );
  const settings = allSettings.general || {};
  if (settings.autoSave) {
    localStorage.setItem(CODE_STORAGE_KEY, sourceCode.value);
    localStorage.setItem(FRAMEWORK_STORAGE_KEY, testConfig.value.framework);
    if (testCode.value) {
      localStorage.setItem(TEST_CODE_STORAGE_KEY, testCode.value);
    }
    if (coverageData.value) {
      localStorage.setItem(COVERAGE_DATA_STORAGE_KEY, JSON.stringify(coverageData.value));
    }
    if (testCases.value && testCases.value.length > 0) {
      localStorage.setItem(TEST_CASES_STORAGE_KEY, JSON.stringify(testCases.value));
    }
  }
};

const sourceCode = ref(`def calculate_total(items):
    total = 0
    for item in items:
        total += item
    return total

def divide_numbers(a, b):
    if b == 0:
        raise ValueError("Cannot divide by zero")
    return a / b

class ShoppingCart:
    def __init__(self):
        self.items = []
    
    def add_item(self, item):
        self.items.append(item)
    
    def get_total(self):
        return sum(self.items)
`);

const testCode = ref("");
const generating = ref(false);

const testConfig = ref({
  framework: "pytest",
  testTypes: ["unit"],
  coverageTarget: 80,
});

const coverageData = ref<{
  total: number;
  statements: number;
  branches: number;
  functions: number;
} | null>(null);

const testCases = ref<
  Array<{
    name: string;
    description: string;
    status: "pending" | "passed" | "failed";
  }>
>([]);

onMounted(() => {
  // 加载保存的代码
  loadSavedCode();

  initEditors();
});

onUnmounted(() => {
  // 清理定时器
  if (saveTimer) {
    clearTimeout(saveTimer);
  }

  // 保存代码到localStorage
  saveCodeToStorage();

  if (sourceEditor) {
    // 清理设置监听器
    if ((sourceEditor as any)._settingsCleanup) {
      (sourceEditor as any)._settingsCleanup();
    }
    sourceEditor.dispose();
  }
  if (testEditor) testEditor.dispose();

  // 断开WebSocket连接
  if (currentTaskId.value) {
    taskStore.cancelTask(currentTaskId.value);
  }
});

const initEditors = () => {
  if (!sourceEditorContainer.value || !testEditorContainer.value) {
    console.error("Editor containers not found");
    return;
  }

  // 确保容器有合适的高度
  sourceEditorContainer.value.style.height = "400px";
  testEditorContainer.value.style.height = "400px";

  // 加载全局编辑器设置
  const editorSettings = loadEditorSettings();
  const options = toMonacoOptions(editorSettings);

  // 源代码编辑器
  sourceEditor = (monaco as any).editor.create(sourceEditorContainer.value, {
    value: sourceCode.value,
    language: "python",
    readOnly: false,
    ...options,
  });

  // 测试代码编辑器
  testEditor = (monaco as any).editor.create(testEditorContainer.value, {
    value: testCode.value || "",
    language: "python",
    readOnly: true,
    ...options,
  });

  sourceEditor.onDidChangeModelContent(() => {
    sourceCode.value = sourceEditor.getValue();
    // 延迟保存，避免频繁写入localStorage
    if (saveTimer) {
      clearTimeout(saveTimer);
    }
    saveTimer = setTimeout(() => {
      saveCodeToStorage();
    }, 1000); // 1秒后保存
  });

  // 监听全局设置变更
  const cleanup = listenToSettingsChanges((settings: EditorSettings) => {
    updateEditorOptions(sourceEditor, settings);
    updateEditorOptions(testEditor, settings);
  });

  // 保存清理函数
  (sourceEditor as any)._settingsCleanup = cleanup;
};

const handleFileChange = async (file: any) => {
  try {
    const rawFile = file.raw;

    // 验证文件类型
    const allowedExtensions = [
      ".py",
      ".js",
      ".java",
      ".go",
      ".rs",
      ".ts",
      ".cpp",
      ".c",
      ".jsx",
      ".tsx",
    ];
    const fileName = rawFile.name.toLowerCase();
    const isValidType = allowedExtensions.some((ext) => fileName.endsWith(ext));

    if (!isValidType) {
      ElMessage.error(
        `不支持的文件类型。支持的类型：${allowedExtensions.join(", ")}`
      );
      return;
    }

    // 验证文件大小（限制为5MB）
    const maxSize = 5 * 1024 * 1024; // 5MB
    if (rawFile.size > maxSize) {
      ElMessage.error("文件大小不能超过5MB");
      return;
    }

    // 读取文件内容
    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      if (content) {
        if (sourceEditor) {
          sourceEditor.setValue(content);
        }
        sourceCode.value = content;

        // 保存到localStorage
        saveCodeToStorage();

        ElMessage.success(`文件加载成功: ${rawFile.name}`);
      }
    };

    reader.onerror = () => {
      ElMessage.error("文件读取失败");
    };

    reader.readAsText(rawFile);
  } catch (error) {
    console.error("File handling error:", error);
    ElMessage.error(
      `文件处理失败: ${error instanceof Error ? error.message : "未知错误"}`
    );
  }
};

const generateTests = async () => {
  if (!sourceCode.value.trim()) {
    ElMessage.warning("请先输入或上传源代码");
    return;
  }

  generating.value = true;
  try {
    // 使用任务存储创建测试生成任务
    const task = await taskStore.createTask(
      "test",
      "测试生成任务",
      "为代码生成单元测试和集成测试",
      {
        code: sourceCode.value,
        language: "python",
        test_type: testConfig.value.testTypes.includes("unit")
          ? "unit"
          : "integration",
        framework: testConfig.value.framework,
        coverage_target: testConfig.value.coverageTarget,
      }
    );

    currentTaskId.value = task.id;
    ElMessage.success("测试生成任务已创建，正在处理...");

    // 等待任务完成并获取结果
    await waitForTaskCompletion(task.id);
  } catch (error) {
    ElMessage.error(
      `生成失败: ${error instanceof Error ? error.message : "未知错误"}`
    );
  } finally {
    generating.value = false;
  }
};

// 等待任务完成
const waitForTaskCompletion = async (taskId: string) => {
  const maxAttempts = 60; // 最多等待5分钟
  let attempts = 0;

  while (attempts < maxAttempts) {
    try {
      // 获取任务详情
      const task = await taskStore.getTask(taskId);

      if (task.status === "completed") {
        // 获取任务结果
        const result = await taskStore.getTaskResult(taskId);
        console.log("Task result received:", result);

        // 处理不同的结果格式
        let testData = null
        
        if (result.result) {
          // 尝试多种可能的数据结构
          if (result.result.result) {
            testData = result.result.result
          } else if (result.result.generated_tests || result.result.test_code) {
            // 直接是测试结果
            testData = result.result
          } else {
            testData = result.result
          }
        }

        console.log("Test data:", testData);

        // 支持多种字段名：generated_tests 或 test_code
        const testCodeValue = testData?.generated_tests || testData?.test_code;
        
        if (testData && testCodeValue) {
          testCode.value = testCodeValue;

          if (testEditor) {
            testEditor.setValue(testCode.value);
          }

          // 设置测试用例
          if (testData.test_cases && Array.isArray(testData.test_cases)) {
            testCases.value = testData.test_cases.map((tc: any) => ({
              name: tc.name || tc.test_name || "未知测试",
              description: tc.description || tc.test_description || "",
              status: tc.status || "pending",
            }));
          }

          // 设置覆盖率数据
          if (testData.coverage) {
            coverageData.value = {
              total: testData.coverage.total || 0,
              statements: testData.coverage.statements || 0,
              branches: testData.coverage.branches || 0,
              functions: testData.coverage.functions || 0,
            };
          } else if (testData.coverage_estimate !== undefined) {
            // 如果只有覆盖率估计值
            coverageData.value = {
              total: testData.coverage_estimate,
              statements: testData.coverage_estimate,
              branches: testData.coverage_estimate,
              functions: testData.coverage_estimate,
            };
          }

          saveCodeToStorage(); // 保存测试结果
          ElMessage.success("测试生成完成");
          return;
        } else {
          console.error("Invalid result structure:", result);
          ElMessage.error("测试结果格式错误");
          return;
        }
      } else if (task.status === "failed") {
        ElMessage.error(`生成失败: ${task.error || "未知错误"}`);
        return;
      }

      // 等待1秒后重试
      await new Promise((resolve) => setTimeout(resolve, 5000));
      attempts++;
    } catch (error) {
      console.error("Error checking task status:", error);
      await new Promise((resolve) => setTimeout(resolve, 5000));
      attempts++;
    }
  }

  ElMessage.error("生成超时，请稍后重试");
};

const downloadTests = () => {
  if (!testCode.value) {
    ElMessage.warning("没有可下载的测试代码");
    return;
  }

  const blob = new Blob([testCode.value], { type: "text/plain" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = "test_generated.py";
  link.click();
  URL.revokeObjectURL(url);
  ElMessage.success("测试代码下载成功");
};

const copyTests = () => {
  if (!testCode.value) {
    ElMessage.warning("没有可复制的测试代码");
    return;
  }

  navigator.clipboard.writeText(testCode.value).then(() => {
    ElMessage.success("测试代码已复制到剪贴板");
  });
};

const getCoverageColor = (percentage: number) => {
  if (percentage >= 80) return "#66bb6a";
  if (percentage >= 60) return "#ffb74d";
  return "#ef5350";
};

const getTestCaseType = (status: string) => {
  const typeMap: Record<string, string> = {
    passed: "success",
    failed: "danger",
    pending: "info",
  };
  return typeMap[status] || "info";
};

const getTestCaseStatus = (status: string) => {
  const statusMap: Record<string, string> = {
    passed: "通过",
    failed: "失败",
    pending: "待执行",
  };
  return statusMap[status] || status;
};
</script>

<style scoped>
.test-generation-container {
  padding: 0;
}

.code-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
  height: 500px;
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

.header-actions {
  display: flex;
  gap: 8px;
}

.editor-container {
  flex: 1;
  min-height: 400px;
}

.monaco-editor-container {
  width: 100%;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}

.config-section,
.coverage-section {
  margin-top: 24px;
}

.config-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.coverage-card,
.test-cases-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
  height: 500px;
  display: flex;
  flex-direction: column;
}

.coverage-card :deep(.el-card__body),
.test-cases-card :deep(.el-card__body) {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.coverage-content {
  text-align: center;
  padding: 20px;
  overflow-y: auto;
  flex: 1;
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

.coverage-details {
  margin-top: 24px;
}

.coverage-item {
  margin-bottom: 16px;
}

.coverage-label {
  display: block;
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.test-cases-content {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.test-case-item {
  padding: 12px;
  background: var(--bg-secondary);
  border-radius: 8px;
}

.test-case-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--text-primary);
}

.test-case-description {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 8px 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .el-col {
    margin-bottom: 16px;
  }

  .code-card {
    height: 400px;
  }
}
</style>
