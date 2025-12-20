<template>
  <div class="code-conversion-container">
    <el-row :gutter="24">
      <!-- 左侧原始代码 -->
      <el-col :span="12">
        <el-card class="code-card">
          <template #header>
            <div class="card-header">
              <h3 class="card-title">原始代码</h3>
              <el-upload
                action="#"
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleOriginalFileChange"
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
              ref="originalEditorContainer"
              class="monaco-editor-container"
              style="height: 100%"
            ></div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧转换后代码 -->
      <el-col :span="12">
        <el-card class="code-card">
          <template #header>
            <div class="card-header">
              <h3 class="card-title">转换后代码</h3>
              <div class="header-actions">
                <el-button
                  type="success"
                  :icon="Download"
                  size="small"
                  @click="downloadConverted"
                >
                  下载代码
                </el-button>
                <el-button
                  :icon="CopyDocument"
                  size="small"
                  @click="copyConverted"
                >
                  复制代码
                </el-button>
              </div>
            </div>
          </template>

          <div class="editor-container">
            <div
              ref="convertedEditorContainer"
              class="monaco-editor-container"
              style="height: 100%"
            ></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 转换配置和差异对比 -->
    <el-row :gutter="24" class="bottom-section">
      <el-col :span="8">
        <el-card class="config-card">
          <template #header>
            <h3 class="card-title">转换配置</h3>
          </template>

          <el-form :model="conversionConfig" label-position="top">
            <el-form-item label="目标语言">
              <el-select
                v-model="conversionConfig.targetLanguage"
                placeholder="选择目标语言"
              >
                <el-option label="Python 3" value="python3" />
                <el-option label="JavaScript (ES6+)" value="javascript" />
                <el-option label="TypeScript" value="typescript" />
                <el-option label="Go" value="go" />
              </el-select>
            </el-form-item>

            <el-form-item label="转换规则">
              <el-checkbox-group v-model="conversionConfig.rules">
                <el-checkbox label="syntax">语法现代化</el-checkbox>
                <el-checkbox label="types">类型注解</el-checkbox>
                <el-checkbox label="async">异步转换</el-checkbox>
                <el-checkbox label="optimize">性能优化</el-checkbox>
              </el-checkbox-group>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                :icon="Refresh"
                @click="handleConvertClick"
                :loading="converting"
                size="large"
                style="width: 100%"
              >
                开始转换
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card class="diff-card">
          <template #header>
            <div class="card-header">
              <h3 class="card-title">差异对比</h3>
              <el-radio-group v-model="diffViewMode" size="small">
                <el-radio-button label="side-by-side">并排对比</el-radio-button>
                <el-radio-button label="line-by-line">逐行对比</el-radio-button>
              </el-radio-group>
            </div>
          </template>

          <div class="diff-content" v-if="diffHtml" v-html="diffHtml"></div>
          <el-empty v-else description="请先进行代码转换" />
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
  Refresh,
} from "@element-plus/icons-vue";
import * as monaco from "monaco-editor";
import { createTwoFilesPatch } from "diff";
import * as Diff2Html from "diff2html";
import "diff2html/bundles/css/diff2html.min.css";
import { useTaskStore } from "@/stores/task";
import {
  loadEditorSettings,
  toMonacoOptions,
  updateEditorOptions,
  listenToSettingsChanges,
  type EditorSettings,
} from "../../utils/editorConfig";

const originalEditorContainer = ref<HTMLElement>();
const convertedEditorContainer = ref<HTMLElement>();
let originalEditor: any;
let convertedEditor: any;
let saveTimer: any = null;
const taskStore = useTaskStore();
const currentTaskId = ref<string | null>(null);

// 代码自动保存的localStorage key
const CODE_STORAGE_KEY = "code_conversion_original_code";
const TARGET_LANGUAGE_STORAGE_KEY = "code_conversion_target_language";
const CONVERTED_CODE_STORAGE_KEY = "code_conversion_converted_code";
const DIFF_HTML_STORAGE_KEY = "code_conversion_diff_html";

// 从localStorage加载保存的代码（如果启用了自动保存）
const loadSavedCode = () => {
  const allSettings = JSON.parse(
    localStorage.getItem("codesage-settings") || "{}"
  );
  const settings = allSettings.general || {};
  if (settings.autoSave) {
    const savedCode = localStorage.getItem(CODE_STORAGE_KEY);
    const savedLanguage = localStorage.getItem(TARGET_LANGUAGE_STORAGE_KEY);
    const savedConverted = localStorage.getItem(CONVERTED_CODE_STORAGE_KEY);
    const savedDiff = localStorage.getItem(DIFF_HTML_STORAGE_KEY);

    if (savedCode) {
      originalCode.value = savedCode;
    }
    if (savedLanguage) {
      conversionConfig.value.targetLanguage = savedLanguage;
    }
    if (savedConverted) {
      convertedCode.value = savedConverted;
    }
    if (savedDiff) {
      diffHtml.value = savedDiff;
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
    localStorage.setItem(CODE_STORAGE_KEY, originalCode.value);
    localStorage.setItem(
      TARGET_LANGUAGE_STORAGE_KEY,
      conversionConfig.value.targetLanguage
    );
    if (convertedCode.value) {
      localStorage.setItem(CONVERTED_CODE_STORAGE_KEY, convertedCode.value);
    }
    if (diffHtml.value) {
      localStorage.setItem(DIFF_HTML_STORAGE_KEY, diffHtml.value);
    }
  }
};

const originalCode = ref(`# Python 2 代码示例
def old_function(x, y):
    result = x + y
    print "结果:", result
    return result

class OldClass:
    def __init__(self):
        self.data = []
    
    def add_item(self, item):
        self.data.append(item)
    
    def get_length(self):
        return len(self.data)

# 使用老式异常处理
try:
    file = open('test.txt', 'r')
    content = file.read()
    file.close()
except IOError, e:
    print "文件读取错误:", e
`);

const convertedCode = ref("");
const converting = ref(false);
const diffHtml = ref("");
const diffViewMode = ref("side-by-side");

const conversionConfig = ref({
  targetLanguage: "python3",
  rules: ["syntax", "types"],
});

onMounted(() => {
  console.log(
    "DEBUG: CodeConversion component mounted - JavaScript is working!"
  );

  // 加载保存的代码
  loadSavedCode();

  // 延迟初始化编辑器，确保DOM完全加载
  setTimeout(() => {
    initEditors();
  }, 100);

  // 添加全局点击监听器用于调试
  document.addEventListener("click", (e) => {
    const target = e.target as HTMLElement;
    console.log("DEBUG: Global click event:", target.tagName, target.className);
  });
});

onUnmounted(() => {
  // 清理定时器
  if (saveTimer) {
    clearTimeout(saveTimer);
  }

  // 保存代码到localStorage
  saveCodeToStorage();

  if (originalEditor) {
    // 清理设置监听器
    if ((originalEditor as any)._settingsCleanup) {
      (originalEditor as any)._settingsCleanup();
    }
    originalEditor.dispose();
  }
  if (convertedEditor) convertedEditor.dispose();

  // 断开WebSocket连接
  if (currentTaskId.value) {
    taskStore.cancelTask(currentTaskId.value);
  }
});

const initEditors = () => {
  if (!originalEditorContainer.value || !convertedEditorContainer.value) {
    console.error("Editor containers not found");
    return;
  }

  // 确保容器有合适的高度
  originalEditorContainer.value.style.height = "400px";
  convertedEditorContainer.value.style.height = "400px";

  // 加载全局编辑器设置
  const editorSettings = loadEditorSettings();
  const options = toMonacoOptions(editorSettings);

  // 原始代码编辑器
  originalEditor = (monaco as any).editor.create(
    originalEditorContainer.value,
    {
      value: originalCode.value,
      language: "python",
      readOnly: false,
      ...options,
    }
  );

  // 转换后代码编辑器
  convertedEditor = (monaco as any).editor.create(
    convertedEditorContainer.value,
    {
      value: convertedCode.value || "",
      language: "python",
      readOnly: true,
      ...options,
    }
  );

  originalEditor.onDidChangeModelContent(() => {
    originalCode.value = originalEditor.getValue();
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
    updateEditorOptions(originalEditor, settings);
    updateEditorOptions(convertedEditor, settings);
  });

  // 保存清理函数
  (originalEditor as any)._settingsCleanup = cleanup;
};

const handleOriginalFileChange = async (file: any) => {
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
        if (originalEditor) {
          originalEditor.setValue(content);
        }
        originalCode.value = content;

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

// 根据目标语言更新代码注释
const updateCodeComments = (code: string, targetLanguage: string): string => {
  // 获取语言的注释样式和名称
  const languageInfo: Record<string, { name: string; commentStart: string; commentEnd?: string }> = {
    python3: { name: "Python 3", commentStart: "#" },
    javascript: { name: "JavaScript (ES6+)", commentStart: "//" },
    typescript: { name: "TypeScript", commentStart: "//" },
    go: { name: "Go", commentStart: "//" },
    java: { name: "Java", commentStart: "//" },
    rust: { name: "Rust", commentStart: "//" },
    cpp: { name: "C++", commentStart: "//" },
    c: { name: "C", commentStart: "//" },
  };

  const info = languageInfo[targetLanguage] || { name: targetLanguage, commentStart: "#" };
  
  // 替换第一行的语言标识注释
  const lines = code.split('\n');
  if (lines.length > 0 && lines[0]) {
    // 检查第一行是否是注释
    const firstLine = lines[0].trim();
    if (firstLine.startsWith('#') || firstLine.startsWith('//')) {
      // 替换为目标语言的注释
      lines[0] = `${info.commentStart} ${info.name} 代码示例`;
    }
  }
  
  return lines.join('\n');
};

const convertCode = async () => {
  console.log("=== DEBUG: Starting code conversion ===");
  console.log("Original code:", originalCode.value);
  console.log("Conversion config:", conversionConfig.value);

  if (!originalCode.value.trim()) {
    console.log("DEBUG: No code to convert");
    ElMessage.warning("请先输入或上传原始代码");
    return;
  }
  console.log("DEBUG: Code validation passed");

  converting.value = true;
  try {
    console.log("Creating conversion task...");
    // 使用任务存储创建转换任务
    const task = await taskStore.createTask(
      "convert",
      "代码转换任务",
      "将代码从Python 2转换到Python 3",
      {
        code: originalCode.value,
        language: "python",
        conversion_type: "python_2_to_3", // 修复转换类型
        from_version: "python2",
        to_version: "python3",
        options: {
          modernize: conversionConfig.value.rules.includes("syntax"),
          add_type_hints: conversionConfig.value.rules.includes("types"),
          optimize: conversionConfig.value.rules.includes("optimize"),
        },
      }
    );

    console.log("Task created:", task);
    currentTaskId.value = task.id;
    ElMessage.success("代码转换任务已创建，正在处理...");

    // 等待任务完成并获取结果
    await waitForTaskCompletion(task.id);
  } catch (error) {
    console.error("Conversion error:", error);
    ElMessage.error(
      `转换失败: ${error instanceof Error ? error.message : "未知错误"}`
    );
  } finally {
    converting.value = false;
  }
};

const handleConvertClick = () => {
  console.log("DEBUG: Convert button clicked!");
  convertCode();
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
        let conversionData = null
        
        if (result.result) {
          // 尝试多种可能的数据结构
          if (result.result.result) {
            conversionData = result.result.result
          } else if (result.result.converted_code || result.result.converted) {
            // 直接是转换结果
            conversionData = result.result
          } else {
            conversionData = result.result
          }
        }

        console.log("Conversion data:", conversionData);

        if (conversionData && (conversionData.converted_code || conversionData.converted)) {
          // 支持两种字段名：converted_code 和 converted
          let code = conversionData.converted_code || conversionData.converted;
          
          // 根据目标语言更新代码注释
          code = updateCodeComments(code, conversionConfig.value.targetLanguage);
          
          convertedCode.value = code;

          if (convertedEditor) {
            convertedEditor.setValue(convertedCode.value);
          }

          generateDiff();
          saveCodeToStorage(); // 保存转换结果
          ElMessage.success("代码转换完成");
          return;
        } else {
          console.error("Invalid result structure:", result);
          ElMessage.error("转换结果格式错误");
          return;
        }
      } else if (task.status === "failed") {
        ElMessage.error(`转换失败: ${task.error || "未知错误"}`);
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

  ElMessage.error("转换超时，请稍后重试");
};

const generateDiff = () => {
  if (!originalCode.value || !convertedCode.value) return;

  // 根据目标语言确定文件扩展名
  const getFileExtension = (lang: string) => {
    const extensionMap: Record<string, string> = {
      python3: "py",
      javascript: "js",
      typescript: "ts",
      go: "go",
      java: "java",
      rust: "rs",
      cpp: "cpp",
      c: "c",
    };
    return extensionMap[lang] || "txt";
  };

  const targetExt = getFileExtension(conversionConfig.value.targetLanguage);
  const targetLangName = conversionConfig.value.targetLanguage === "python3" 
    ? "Python 3" 
    : conversionConfig.value.targetLanguage;

  const patch = createTwoFilesPatch(
    "original.py",
    `converted.${targetExt}`,
    originalCode.value,
    convertedCode.value,
    "原始代码",
    `转换后代码 (${targetLangName})`,
    { context: 10 }
  );

  const diffJson = Diff2Html.parse(patch);
  diffHtml.value = Diff2Html.html(diffJson, {
    outputFormat: diffViewMode.value as any,
    drawFileList: false,
    matching: "lines",
    renderNothingWhenEmpty: false,
  });
};

const downloadConverted = () => {
  if (!convertedCode.value) {
    ElMessage.warning("没有可下载的转换代码");
    return;
  }

  const blob = new Blob([convertedCode.value], { type: "text/plain" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = "converted_code.py";
  link.click();
  URL.revokeObjectURL(url);
  ElMessage.success("代码下载成功");
};

const copyConverted = () => {
  if (!convertedCode.value) {
    ElMessage.warning("没有可复制的转换代码");
    return;
  }

  navigator.clipboard.writeText(convertedCode.value).then(() => {
    ElMessage.success("代码已复制到剪贴板");
  });
};
</script>

<style scoped>
.code-conversion-container {
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

.bottom-section {
  margin-top: 24px;
}

.config-card,
.diff-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
  height: 400px;
  display: flex;
  flex-direction: column;
}

.diff-content {
  flex: 1;
  overflow: auto;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
}

/* diff2html 样式覆盖 */
.diff-content :deep(.d2h-wrapper) {
  font-family: "Consolas", "Monaco", "Courier New", monospace;
}

.diff-content :deep(.d2h-file-wrapper) {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 16px;
  overflow: hidden;
  background: var(--bg-primary);
}

.diff-content :deep(.d2h-file-header) {
  background: var(--bg-tertiary) !important;
  color: var(--text-primary) !important;
  border: none !important;
  border-bottom: 1px solid var(--border-color) !important;
  padding: 12px 16px;
  font-size: 14px;
}

.diff-content :deep(.d2h-file-name) {
  font-weight: 600;
}

.diff-content :deep(.d2h-code-wrapper) {
  overflow-x: auto;
  position: relative;
}

.diff-content :deep(.d2h-diff-table) {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  line-height: 1.5;
}

.diff-content :deep(.d2h-diff-tbody) {
  background: var(--bg-primary);
}

.diff-content :deep(tr) {
  display: table-row;
}

.diff-content :deep(td) {
  display: table-cell;
  padding: 0;
  vertical-align: top;
  border: none;
}

.diff-content :deep(.d2h-code-line) {
  background: var(--bg-primary) !important;
  color: var(--text-primary) !important;
  padding: 0;
  vertical-align: top;
}

.diff-content :deep(.d2h-code-linenumber),
.diff-content :deep(.d2h-code-side-linenumber) {
  background: var(--bg-tertiary) !important;
  color: var(--text-secondary) !important;
  border-right: 1px solid var(--border-color) !important;
  padding: 0 8px !important;
  text-align: right;
  width: 45px;
  min-width: 45px;
  user-select: none;
  vertical-align: top;
  font-size: 12px;
  white-space: nowrap;
  position: static !important;
}

.diff-content :deep(.d2h-code-line-prefix) {
  padding: 0 8px !important;
  text-align: center;
  width: 25px;
  min-width: 25px;
  user-select: none;
  font-weight: bold;
  vertical-align: top;
  white-space: nowrap;
}

.diff-content :deep(.d2h-code-line-ctn) {
  color: var(--text-primary) !important;
  padding: 0 10px !important;
  white-space: pre;
  word-wrap: normal;
  vertical-align: top;
}

.diff-content :deep(.d2h-del) {
  background: rgba(239, 83, 80, 0.1) !important;
}

.diff-content :deep(.d2h-del .d2h-code-linenumber),
.diff-content :deep(.d2h-del .d2h-code-side-linenumber) {
  background: rgba(239, 83, 80, 0.2) !important;
  border-right-color: rgba(239, 83, 80, 0.4) !important;
}

.diff-content :deep(.d2h-del .d2h-code-line-prefix) {
  background: rgba(239, 83, 80, 0.15) !important;
  color: #ef5350 !important;
}

.diff-content :deep(.d2h-del .d2h-code-line-ctn) {
  background: rgba(239, 83, 80, 0.08) !important;
}

.diff-content :deep(.d2h-ins) {
  background: rgba(102, 187, 106, 0.1) !important;
}

.diff-content :deep(.d2h-ins .d2h-code-linenumber),
.diff-content :deep(.d2h-ins .d2h-code-side-linenumber) {
  background: rgba(102, 187, 106, 0.2) !important;
  border-right-color: rgba(102, 187, 106, 0.4) !important;
}

.diff-content :deep(.d2h-ins .d2h-code-line-prefix) {
  background: rgba(102, 187, 106, 0.15) !important;
  color: #66bb6a !important;
}

.diff-content :deep(.d2h-ins .d2h-code-line-ctn) {
  background: rgba(102, 187, 106, 0.08) !important;
}

.diff-content :deep(.d2h-info) {
  background: var(--bg-tertiary) !important;
  color: var(--text-secondary) !important;
  border-top: 1px solid var(--border-color) !important;
  border-bottom: 1px solid var(--border-color) !important;
}

.diff-content :deep(.d2h-emptyplaceholder) {
  background: var(--bg-tertiary) !important;
  border-right: 1px solid var(--border-color) !important;
  width: 45px;
  min-width: 45px;
}

.diff-content :deep(.d2h-moved-tag) {
  display: none;
}

/* 并排对比模式 */
.diff-content :deep(.d2h-file-side-diff) {
  width: 100%;
}

.diff-content :deep(.d2h-file-side-diff .d2h-code-side-line) {
  width: 50%;
}

/* 逐行对比模式 */
.diff-content :deep(.d2h-file-diff) {
  width: 100%;
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
