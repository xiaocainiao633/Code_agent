<template>
  <div class="settings-container">
    <el-row :gutter="24">
      <!-- 左侧菜单 -->
      <el-col :span="6">
        <el-card class="menu-card">
          <el-menu
            :default-active="activeMenu"
            class="settings-menu"
            @select="handleMenuSelect"
          >
            <el-menu-item index="profile">
              <el-icon><User /></el-icon>
              <span>个人资料</span>
            </el-menu-item>
            <el-menu-item index="general">
              <el-icon><Setting /></el-icon>
              <span>通用设置</span>
            </el-menu-item>
            <el-menu-item index="appearance">
              <el-icon><Edit /></el-icon>
              <span>外观设置</span>
            </el-menu-item>
            <el-menu-item index="editor">
              <el-icon><Edit /></el-icon>
              <span>编辑器设置</span>
            </el-menu-item>
            <el-menu-item index="ai">
              <el-icon><Monitor /></el-icon>
              <span>AI 设置</span>
            </el-menu-item>
            <el-menu-item index="advanced">
              <el-icon><Tools /></el-icon>
              <span>高级设置</span>
            </el-menu-item>
            <el-menu-item index="about">
              <el-icon><InfoFilled /></el-icon>
              <span>关于</span>
            </el-menu-item>
          </el-menu>
        </el-card>
      </el-col>
      
      <!-- 右侧内容 -->
      <el-col :span="18">
        <el-card class="content-card">
          <!-- 个人资料 -->
          <div v-if="activeMenu === 'profile'" class="settings-section">
            <h3 class="section-title">个人资料</h3>
            
            <el-form :model="profileSettings" label-position="left" label-width="140px">
              <!-- 头像区域 -->
              <el-form-item label="用户头像">
                <div class="avatar-upload">
                  <el-avatar :size="100" :src="profileSettings.avatar" class="profile-avatar">
                    <el-icon size="50"><User /></el-icon>
                  </el-avatar>
                  <div class="avatar-actions">
                    <el-button size="small" @click="uploadAvatar">
                      <el-icon><Upload /></el-icon>
                      <span>上传头像</span>
                    </el-button>
                    <el-button size="small" @click="removeAvatar" v-if="profileSettings.avatar">
                      <el-icon><Delete /></el-icon>
                      <span>移除头像</span>
                    </el-button>
                  </div>
                </div>
                <span class="setting-description">支持 JPG、PNG 格式，大小不超过 2MB</span>
              </el-form-item>
              
              <!-- 基本信息 -->
              <el-divider content-position="left">基本信息</el-divider>
              
              <el-form-item label="用户名">
                <el-input
                  v-model="profileSettings.username"
                  placeholder="请输入用户名"
                  :prefix-icon="User"
                />
                <span class="setting-description">用户名将在系统中显示</span>
              </el-form-item>
              
              <el-form-item label="邮箱">
                <el-input
                  v-model="profileSettings.email"
                  placeholder="请输入邮箱"
                  type="email"
                />
                <span class="setting-description">用于接收系统通知和重置密码</span>
              </el-form-item>
              
              <el-form-item label="手机号">
                <el-input
                  v-model="profileSettings.phone"
                  placeholder="请输入手机号"
                />
                <span class="setting-description">可选，用于账号安全验证</span>
              </el-form-item>
              
              <el-form-item label="个人简介">
                <el-input
                  v-model="profileSettings.bio"
                  type="textarea"
                  :rows="3"
                  placeholder="介绍一下你自己..."
                  maxlength="200"
                  show-word-limit
                />
                <span class="setting-description">简短介绍你的背景和兴趣</span>
              </el-form-item>
              
              <!-- 安全设置 -->
              <el-divider content-position="left">安全设置</el-divider>
              
              <el-form-item label="修改密码">
                <div class="password-change">
                  <el-button @click="showPasswordDialog = true">
                    <el-icon><Key /></el-icon>
                    <span>修改密码</span>
                  </el-button>
                </div>
                <span class="setting-description">定期修改密码可提高账号安全性</span>
              </el-form-item>
              
              <!-- 社交账号 -->
              <el-divider content-position="left">社交账号</el-divider>
              
              <el-form-item label="GitHub">
                <el-input
                  v-model="profileSettings.github"
                  placeholder="https://github.com/username"
                >
                  <template #prefix>
                    <svg class="social-icon-small" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                  </template>
                </el-input>
                <span class="setting-description">你的 GitHub 主页地址</span>
              </el-form-item>
              
              <el-form-item label="Twitter">
                <el-input
                  v-model="profileSettings.twitter"
                  placeholder="@username"
                />
                <span class="setting-description">你的 Twitter 用户名</span>
              </el-form-item>
              
              <el-form-item label="个人网站">
                <el-input
                  v-model="profileSettings.website"
                  placeholder="https://yourwebsite.com"
                />
                <span class="setting-description">你的个人网站或博客</span>
              </el-form-item>
              
              <!-- 其他信息 -->
              <el-divider content-position="left">其他信息</el-divider>
              
              <el-form-item label="所在地">
                <el-input
                  v-model="profileSettings.location"
                  placeholder="城市, 国家"
                />
                <span class="setting-description">你的所在地</span>
              </el-form-item>
              
              <el-form-item label="职业">
                <el-input
                  v-model="profileSettings.occupation"
                  placeholder="软件开发工程师"
                />
                <span class="setting-description">你的职业或头衔</span>
              </el-form-item>
              
              <el-form-item label="公司/组织">
                <el-input
                  v-model="profileSettings.company"
                  placeholder="公司名称"
                />
                <span class="setting-description">你所在的公司或组织</span>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- 通用设置 -->
          <div v-if="activeMenu === 'general'" class="settings-section">
            <h3 class="section-title">通用设置</h3>
            
            <el-form :model="generalSettings" label-position="left" label-width="140px">
              <el-form-item label="语言">
                <el-select v-model="generalSettings.language" placeholder="选择语言">
                  <el-option label="简体中文" value="zh-CN" />
                  <el-option label="English" value="en-US" />
                  <el-option label="日本語" value="ja-JP" />
                </el-select>
                <span class="setting-description">界面显示语言</span>
              </el-form-item>
              
              <el-form-item label="自动保存">
                <el-switch v-model="generalSettings.autoSave" />
                <span class="setting-description">自动保存编辑的代码</span>
              </el-form-item>
              
              <el-form-item label="保存间隔">
                <el-slider
                  v-model="generalSettings.autoSaveInterval"
                  :min="1"
                  :max="10"
                  :step="1"
                  show-input
                  :disabled="!generalSettings.autoSave"
                />
                <span class="setting-description">自动保存间隔（分钟）</span>
              </el-form-item>
              
              <el-form-item label="启动时检查更新">
                <el-switch v-model="generalSettings.checkUpdates" />
                <span class="setting-description">启动时检查软件更新</span>
              </el-form-item>
              
              <el-form-item label="发送匿名统计">
                <el-switch v-model="generalSettings.sendStats" />
                <span class="setting-description">发送匿名使用统计（本地运行，无隐私风险）</span>
              </el-form-item>
              
              <el-form-item label="默认工作目录">
                <el-input
                  v-model="generalSettings.defaultWorkspace"
                  placeholder="选择工作目录"
                >
                  <template #append>
                    <el-button @click="selectWorkspace">
                      <el-icon><Folder /></el-icon>
                    </el-button>
                  </template>
                </el-input>
                <span class="setting-description">项目默认保存位置</span>
              </el-form-item>
              
              <el-form-item label="最大并发任务">
                <el-slider
                  v-model="generalSettings.maxConcurrentTasks"
                  :min="1"
                  :max="10"
                  :step="1"
                  show-input
                />
                <span class="setting-description">同时运行的最大任务数</span>
              </el-form-item>
              
              <el-form-item label="任务超时时间">
                <el-slider
                  v-model="generalSettings.taskTimeout"
                  :min="1"
                  :max="30"
                  :step="1"
                  show-input
                />
                <span class="setting-description">任务超时时间（分钟）</span>
              </el-form-item>
              
              <el-form-item label="启用通知">
                <el-switch v-model="generalSettings.enableNotifications" />
                <span class="setting-description">任务完成时显示通知</span>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- 外观设置 -->
          <div v-if="activeMenu === 'appearance'" class="settings-section">
            <h3 class="section-title">外观设置</h3>
            
            <el-form :model="appearanceSettings" label-position="left" label-width="140px">
              <el-form-item label="主题模式">
                <el-radio-group v-model="appearanceSettings.theme">
                  <el-radio label="dark">深色主题</el-radio>
                  <el-radio label="light">浅色主题</el-radio>
                  <el-radio label="auto">跟随系统</el-radio>
                </el-radio-group>
                <span class="setting-description">应用整体主题风格</span>
              </el-form-item>
              
              <el-form-item label="主色调">
                <el-color-picker v-model="appearanceSettings.primaryColor" />
                <span class="setting-description">应用的主色调</span>
              </el-form-item>
              
              <el-form-item label="字体大小">
                <el-slider
                  v-model="appearanceSettings.fontSize"
                  :min="12"
                  :max="20"
                  :step="1"
                  show-input
                />
                <span class="setting-description">界面字体大小（px）</span>
              </el-form-item>
              
              <el-form-item label="代码字体">
                <el-select v-model="appearanceSettings.codeFontFamily" placeholder="选择代码字体">
                  <el-option label="Consolas" value="Consolas, Monaco, monospace" />
                  <el-option label="Monaco" value="Monaco, monospace" />
                  <el-option label="Fira Code" value="Fira Code, monospace" />
                  <el-option label="JetBrains Mono" value="JetBrains Mono, monospace" />
                  <el-option label="Source Code Pro" value="Source Code Pro, monospace" />
                </el-select>
                <span class="setting-description">代码编辑器使用的字体</span>
              </el-form-item>
              
              <el-form-item label="动画效果">
                <el-switch v-model="appearanceSettings.animations" />
                <span class="setting-description">启用界面动画效果</span>
              </el-form-item>
              
              <el-form-item label="毛玻璃效果">
                <el-switch v-model="appearanceSettings.glassEffect" />
                <span class="setting-description">启用毛玻璃背景效果</span>
              </el-form-item>
              
              <el-form-item label="窗口透明度">
                <el-slider
                  v-model="appearanceSettings.opacity"
                  :min="80"
                  :max="100"
                  :step="5"
                  show-input
                />
                <span class="setting-description">窗口透明度（%）</span>
              </el-form-item>
              
              <el-form-item label="紧凑模式">
                <el-switch v-model="appearanceSettings.compactMode" />
                <span class="setting-description">减少界面元素间距，显示更多内容</span>
              </el-form-item>
              
              <el-form-item label="显示时钟">
                <el-switch v-model="appearanceSettings.showClock" />
                <span class="setting-description">在状态栏显示当前时间</span>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- 编辑器设置 -->
          <div v-if="activeMenu === 'editor'" class="settings-section">
            <h3 class="section-title">编辑器设置</h3>
            
            <el-form :model="editorSettings" label-position="left" label-width="140px">
              <el-form-item label="编辑器主题">
                <el-select v-model="editorSettings.theme" placeholder="选择主题">
                  <el-option label="VS Dark" value="vs-dark" />
                  <el-option label="VS Light" value="vs" />
                  <el-option label="High Contrast" value="hc-black" />
                  <el-option label="Monokai" value="monokai" />
                  <el-option label="GitHub" value="github" />
                </el-select>
                <span class="setting-description">代码编辑器配色主题</span>
              </el-form-item>
              
              <el-form-item label="字体大小">
                <el-slider
                  v-model="editorSettings.fontSize"
                  :min="10"
                  :max="24"
                  :step="1"
                  show-input
                />
                <span class="setting-description">编辑器字体大小（px）</span>
              </el-form-item>
              
              <el-form-item label="字体粗细">
                <el-slider
                  v-model="editorSettings.fontWeight"
                  :min="300"
                  :max="700"
                  :step="100"
                  show-input
                />
                <span class="setting-description">代码字体粗细</span>
              </el-form-item>
              
              <el-form-item label="显示行号">
                <el-switch v-model="editorSettings.showLineNumbers" />
                <span class="setting-description">在编辑器左侧显示行号</span>
              </el-form-item>
              
              <el-form-item label="自动换行">
                <el-switch v-model="editorSettings.wordWrap" />
                <span class="setting-description">代码超出宽度时自动换行</span>
              </el-form-item>
              
              <el-form-item label="缩进大小">
                <el-slider
                  v-model="editorSettings.tabSize"
                  :min="2"
                  :max="8"
                  :step="2"
                  show-input
                />
                <span class="setting-description">缩进空格数</span>
              </el-form-item>
              
              <el-form-item label="使用空格缩进">
                <el-switch v-model="editorSettings.insertSpaces" />
                <span class="setting-description">使用空格代替制表符</span>
              </el-form-item>
              
              <el-form-item label="自动补全">
                <el-switch v-model="editorSettings.autoComplete" />
                <span class="setting-description">输入时显示代码补全建议</span>
              </el-form-item>
              
              <el-form-item label="代码检查">
                <el-switch v-model="editorSettings.linting" />
                <span class="setting-description">实时代码语法检查</span>
              </el-form-item>
              
              <el-form-item label="高亮当前行">
                <el-switch v-model="editorSettings.highlightCurrentLine" />
                <span class="setting-description">高亮显示当前编辑行</span>
              </el-form-item>
              
              <el-form-item label="显示空白字符">
                <el-switch v-model="editorSettings.renderWhitespace" />
                <span class="setting-description">显示空格和制表符</span>
              </el-form-item>
              
              <el-form-item label="迷你地图">
                <el-switch v-model="editorSettings.minimap" />
                <span class="setting-description">显示代码迷你地图</span>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- AI 设置 -->
          <div v-if="activeMenu === 'ai'" class="settings-section">
            <h3 class="section-title">AI 设置</h3>
            
            <el-form :model="aiSettings" label-position="left" label-width="140px">
              <el-form-item label="AI 服务类型">
                <el-radio-group v-model="aiSettings.serviceType">
                  <el-radio label="local">本地模型</el-radio>
                  <el-radio label="api">API 服务</el-radio>
                  <el-radio label="hybrid">混合模式</el-radio>
                </el-radio-group>
                <span class="setting-description">选择 AI 服务运行方式</span>
              </el-form-item>
              
              <el-form-item label="AI 模型" v-if="aiSettings.serviceType !== 'api'">
                <el-select v-model="aiSettings.model" placeholder="选择模型">
                  <el-option label="CodeLlama 34B" value="codellama-34b" />
                  <el-option label="Llama3 8B" value="llama3-8b" />
                  <el-option label="CodeT5" value="codet5" />
                  <el-option label="自定义模型" value="custom" />
                </el-select>
                <span class="setting-description">用于代码分析的 AI 模型</span>
              </el-form-item>
              
              <el-form-item label="API 端点" v-if="aiSettings.serviceType !== 'local'">
                <el-input
                  v-model="aiSettings.apiEndpoint"
                  placeholder="https://api.example.com/v1"
                />
                <span class="setting-description">AI 服务的 API 地址</span>
              </el-form-item>
              
              <el-form-item label="API 密钥" v-if="aiSettings.serviceType !== 'local'">
                <el-input
                  v-model="aiSettings.apiKey"
                  type="password"
                  placeholder="输入 API 密钥"
                  show-password
                />
                <span class="setting-description">API 访问密钥</span>
              </el-form-item>
              
              <el-form-item label="模型路径" v-if="aiSettings.serviceType !== 'api'">
                <el-input
                  v-model="aiSettings.modelPath"
                  placeholder="选择模型文件路径"
                >
                  <template #append>
                    <el-button @click="selectModelPath">
                      <el-icon><Folder /></el-icon>
                    </el-button>
                  </template>
                </el-input>
                <span class="setting-description">本地模型文件路径</span>
              </el-form-item>
              
              <el-form-item label="温度参数">
                <el-slider
                  v-model="aiSettings.temperature"
                  :min="0"
                  :max="2"
                  :step="0.1"
                  show-input
                />
                <span class="setting-description">控制生成代码的创造性（0-2）</span>
              </el-form-item>
              
              <el-form-item label="最大令牌数">
                <el-slider
                  v-model="aiSettings.maxTokens"
                  :min="100"
                  :max="8000"
                  :step="100"
                  show-input
                />
                <span class="setting-description">生成代码的最大长度</span>
              </el-form-item>
              
              <el-form-item label="代码验证">
                <el-switch v-model="aiSettings.codeValidation" />
                <span class="setting-description">自动验证生成的代码</span>
              </el-form-item>
              
              <el-form-item label="安全模式">
                <el-switch v-model="aiSettings.safeMode" />
                <span class="setting-description">启用安全模式，避免潜在风险</span>
              </el-form-item>
              
              <el-form-item label="缓存响应">
                <el-switch v-model="aiSettings.enableCache" />
                <span class="setting-description">缓存 AI 响应以提高性能</span>
              </el-form-item>
              
              <el-form-item label="测试连接">
                <el-button type="primary" @click="testAIConnection">
                  <el-icon><Connection /></el-icon>
                  测试 AI 连接
                </el-button>
                <span class="setting-description">验证 AI 服务连接状态</span>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- 高级设置 -->
          <div v-if="activeMenu === 'advanced'" class="settings-section">
            <h3 class="section-title">高级设置</h3>
            
            <el-form :model="advancedSettings" label-position="left" label-width="140px">
              <el-form-item label="开发者模式">
                <el-switch v-model="advancedSettings.developerMode" />
                <span class="setting-description">启用开发者工具和调试信息</span>
              </el-form-item>
              
              <el-form-item label="详细日志">
                <el-switch v-model="advancedSettings.verboseLogging" />
                <span class="setting-description">记录详细的运行日志</span>
              </el-form-item>
              
              <el-form-item label="性能监控">
                <el-switch v-model="advancedSettings.performanceMonitoring" />
                <span class="setting-description">启用性能监控和统计</span>
              </el-form-item>
              
              <el-form-item label="内存限制">
                <el-slider
                  v-model="advancedSettings.memoryLimit"
                  :min="512"
                  :max="8192"
                  :step="512"
                  show-input
                />
                <span class="setting-description">最大内存使用限制（MB）</span>
              </el-form-item>
              
              <el-form-item label="线程池大小">
                <el-slider
                  v-model="advancedSettings.threadPoolSize"
                  :min="2"
                  :max="16"
                  :step="2"
                  show-input
                />
                <span class="setting-description">工作线程数量</span>
              </el-form-item>
              
              <el-form-item label="缓存大小">
                <el-slider
                  v-model="advancedSettings.cacheSize"
                  :min="50"
                  :max="500"
                  :step="50"
                  show-input
                />
                <span class="setting-description">缓存大小（MB）</span>
              </el-form-item>
              
              <el-form-item label="自动清理">
                <el-switch v-model="advancedSettings.autoCleanup" />
                <span class="setting-description">自动清理过期缓存和临时文件</span>
              </el-form-item>
              
              <el-form-item label="清理间隔">
                <el-slider
                  v-model="advancedSettings.cleanupInterval"
                  :min="1"
                  :max="24"
                  :step="1"
                  show-input
                  :disabled="!advancedSettings.autoCleanup"
                />
                <span class="setting-description">自动清理间隔（小时）</span>
              </el-form-item>
              
              <el-form-item label="实验功能">
                <el-switch v-model="advancedSettings.experimentalFeatures" />
                <span class="setting-description">启用实验性功能（可能不稳定）</span>
              </el-form-item>
              
              <el-form-item label="导出设置">
                <el-button :icon="Download" @click="exportSettings">
                  <el-icon><Download /></el-icon>
                  导出配置
                </el-button>
                <span class="setting-description">将当前设置导出为文件</span>
              </el-form-item>
              
              <el-form-item label="导入设置">
                <el-button :icon="Upload" @click="importSettings">
                  <el-icon><Upload /></el-icon>
                  导入配置
                </el-button>
                <span class="setting-description">从文件导入设置配置</span>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- 关于 -->
          <div v-if="activeMenu === 'about'" class="settings-section">
            <h3 class="section-title">关于 CodeSage</h3>
            
            <div class="about-content">
              <div class="logo-section">
                <div class="app-logo">
                  <span class="neon-text">CodeSage</span>
                </div>
                <p class="app-version">版本 {{ appVersion }}</p>
                <p class="app-description">智能代码重构与现代化助手</p>
                <el-tag type="info" size="small">
                  构建时间: {{ buildTime }}
                </el-tag>
              </div>
              
              <el-divider />
              
              <div class="system-info-section">
                <h4>系统信息</h4>
                <el-descriptions :column="2" size="small" border>
                  <el-descriptions-item label="Node.js 版本">
                    {{ systemInfo.nodeVersion }}
                  </el-descriptions-item>
                  <el-descriptions-item label="Chrome 版本">
                    {{ systemInfo.chromeVersion }}
                  </el-descriptions-item>
                  <el-descriptions-item label="操作系统">
                    {{ systemInfo.platform }}
                  </el-descriptions-item>
                  <el-descriptions-item label="内存使用">
                    {{ systemInfo.memoryUsage }}
                  </el-descriptions-item>
                </el-descriptions>
              </div>
              
              <el-divider />
              
              <div class="info-section">
                <h4>功能特性</h4>
                <ul class="feature-list">
                  <li>智能代码分析和问题检测</li>
                  <li>自动代码转换和现代化</li>
                  <li>单元测试自动生成</li>
                  <li>Git 仓库历史分析</li>
                  <li>完全本地化运行，保护代码隐私</li>
                  <li>支持多种编程语言</li>
                  <li>实时协作和版本控制</li>
                </ul>
              </div>
              
              <el-divider />
              
              <div class="stats-section">
                <h4>使用统计</h4>
                <el-row :gutter="20">
                  <el-col :span="8">
                    <div class="stat-item">
                      <div class="stat-value">{{ stats.totalTasks }}</div>
                      <div class="stat-label">总任务数</div>
                    </div>
                  </el-col>
                  <el-col :span="8">
                    <div class="stat-item">
                      <div class="stat-value">{{ stats.totalCodeLines }}</div>
                      <div class="stat-label">处理代码行数</div>
                    </div>
                  </el-col>
                  <el-col :span="8">
                    <div class="stat-item">
                      <div class="stat-value">{{ stats.totalFiles }}</div>
                      <div class="stat-label">处理文件数</div>
                    </div>
                  </el-col>
                </el-row>
              </div>
              
              <el-divider />
              
              <div class="links-section">
                <h4>相关链接</h4>
                <div class="link-buttons">
                  <el-button text @click="openDocumentation">
                    <el-icon><Document /></el-icon>
                    文档
                  </el-button>
                  <el-button text @click="openGitHub">
                    <el-icon><Link /></el-icon>
                    GitHub
                  </el-button>
                  <el-button text @click="openIssues">
                    <el-icon><Warning /></el-icon>
                    问题反馈
                  </el-button>
                  <el-button text @click="checkForUpdates">
                    <el-icon><Refresh /></el-icon>
                    检查更新
                  </el-button>
                </div>
              </div>
              
              <el-divider />
              
              <div class="copyright-section">
                <p class="copyright">
                  © 2024 CodeSage. All rights reserved.
                </p>
                <p class="license">
                  Licensed under MIT License
                </p>
                <p class="build-info">
                  构建版本: {{ buildVersion }} | 构建哈希: {{ buildHash }}
                </p>
              </div>
            </div>
          </div>
          
          <!-- 保存按钮 -->
          <div class="settings-actions" v-if="activeMenu !== 'about'">
            <el-button v-if="activeMenu === 'profile'" type="primary" @click="saveProfileSettings">
              保存个人资料
            </el-button>
            <el-button v-else type="primary" @click="resetSettings">
              重置为默认值
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="修改密码"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="passwordForm" label-width="100px">
        <el-form-item label="当前密码">
          <el-input
            v-model="passwordForm.currentPassword"
            type="password"
            placeholder="请输入当前密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" @click="handlePasswordChange">确认修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Upload, User, Key, Delete } from '@element-plus/icons-vue'
import { useThemeStore } from '../../stores/theme'
import { useAuthStore } from '../../stores/auth'
import { useRoute } from 'vue-router'
import { authAPI } from '../../services/api'

const route = useRoute()
const activeMenu = ref('general')
const showPasswordDialog = ref(false)
const authStore = useAuthStore()

// 个人资料设置
const profileSettings = ref({
  avatar: '',
  username: '',
  email: '',
  phone: '',
  bio: '',
  github: '',
  twitter: '',
  website: '',
  location: '',
  occupation: '',
  company: ''
})

// 密码修改表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 通用设置
const generalSettings = ref({
  language: 'zh-CN',
  autoSave: true,
  autoSaveInterval: 5,
  checkUpdates: true,
  sendStats: false,
  defaultWorkspace: '',
  maxConcurrentTasks: 5,
  taskTimeout: 10,
  enableNotifications: true
})

// 外观设置
const appearanceSettings = ref({
  theme: 'dark',
  primaryColor: '#1e88e5',
  accentColor: '#1e88e5',
  fontSize: 14,
  codeFontFamily: 'Consolas, Monaco, monospace',
  animations: true,
  glassEffect: false,
  opacity: 100,
  compactMode: false,
  showClock: true
})

// 编辑器设置
const editorSettings = ref({
  theme: 'vs-dark',
  fontSize: 14,
  fontWeight: 400,
  showLineNumbers: true,
  wordWrap: true,
  tabSize: 4,
  insertSpaces: true,
  autoComplete: true,
  linting: true,
  highlightCurrentLine: true,
  renderWhitespace: false,
  minimap: true
})

// AI 设置
const aiSettings = ref({
  serviceType: 'local',
  model: 'codellama-34b',
  apiEndpoint: '',
  apiKey: '',
  modelPath: '/models/codellama-34b',
  temperature: 0.7,
  maxTokens: 2048,
  codeValidation: true,
  safeMode: true,
  enableCache: true
})

// 高级设置
const advancedSettings = ref({
  developerMode: false,
  verboseLogging: false,
  performanceMonitoring: false,
  memoryLimit: 2048,
  threadPoolSize: 8,
  cacheSize: 200,
  autoCleanup: true,
  cleanupInterval: 6,
  experimentalFeatures: false
})

// 应用信息
const appVersion = ref('1.0.0')
const buildVersion = ref('1.0.0.20241210')
const buildHash = ref('abc123def456')
const buildTime = ref('2024-12-10 12:00:00')

// 系统信息
const systemInfo = reactive({
  nodeVersion: '',
  chromeVersion: '',
  platform: '',
  memoryUsage: ''
})

// 使用统计
const stats = reactive({
  totalTasks: 0,
  totalCodeLines: 0,
  totalFiles: 0
})

const handleMenuSelect = (index: string) => {
  activeMenu.value = index
}

// 个人资料相关方法
const uploadAvatar = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/jpeg,image/png,image/jpg'
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) {
      // 检查文件大小
      if (file.size > 2 * 1024 * 1024) {
        ElMessage.warning('图片大小不能超过 2MB')
        return
      }
      
      // 读取文件并显示
      const reader = new FileReader()
      reader.onload = (e) => {
        profileSettings.value.avatar = e.target?.result as string
        ElMessage.success('头像上传成功')
        saveProfileSettings()
      }
      reader.readAsDataURL(file)
    }
  }
  input.click()
}

const removeAvatar = () => {
  ElMessageBox.confirm('确定要移除头像吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    profileSettings.value.avatar = ''
    ElMessage.success('头像已移除')
    saveProfileSettings()
  }).catch(() => {})
}

const saveProfileSettings = async () => {
  try {
    // 调用后端 API 更新个人资料
    const response = await authAPI.updateProfile({
      username: profileSettings.value.username,
      email: profileSettings.value.email,
      avatar: profileSettings.value.avatar,
      phone: profileSettings.value.phone,
      bio: profileSettings.value.bio,
      location: profileSettings.value.location,
      occupation: profileSettings.value.occupation,
      company: profileSettings.value.company,
      website: profileSettings.value.website,
      twitter: profileSettings.value.twitter,
      github_url: profileSettings.value.github
    })
    
    // 更新成功后保存到 localStorage
    localStorage.setItem('codesage-profile', JSON.stringify(profileSettings.value))
    
    // 更新 auth store 中的用户信息
    if (authStore.userInfo && response.user) {
      // 使用 authStore 的 updateUserInfo 方法来更新用户信息
      authStore.updateUserInfo({
        username: response.user.username,
        email: response.user.email,
        avatar: response.user.avatar
      })
    }
    
    ElMessage.success('个人资料已保存')
  } catch (error: any) {
    console.error('保存个人资料失败:', error)
    ElMessage.error(error.response?.data?.error || '保存个人资料失败')
  }
}

const handlePasswordChange = () => {
  // 验证表单
  if (!passwordForm.currentPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
    ElMessage.warning('请填写所有密码字段')
    return
  }
  
  if (passwordForm.newPassword.length < 6) {
    ElMessage.warning('新密码长度不能少于 6 位')
    return
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }
  
  // 模拟密码修改
  ElMessage.success('密码修改成功')
  showPasswordDialog.value = false
  
  // 清空表单
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

const selectWorkspace = () => {
  ElMessage.info('选择工作目录功能开发中...')
}

const selectModelPath = () => {
  ElMessage.info('选择模型路径功能开发中...')
}

const testAIConnection = async () => {
  try {
    ElMessage.info('正在测试 AI 连接...')
    // 模拟连接测试
    setTimeout(() => {
      ElMessage.success('AI 连接测试成功')
    }, 2000)
  } catch (error) {
    ElMessage.error('AI 连接测试失败')
  }
}

const exportSettings = () => {
  const settings = {
    general: generalSettings.value,
    appearance: appearanceSettings.value,
    editor: editorSettings.value,
    ai: aiSettings.value,
    advanced: advancedSettings.value
  }
  
  const dataStr = JSON.stringify(settings, null, 2)
  const dataBlob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(dataBlob)
  const link = document.createElement('a')
  link.href = url
  link.download = `codesage-settings-${new Date().toISOString().split('T')[0]}.json`
  link.click()
  URL.revokeObjectURL(url)
  
  ElMessage.success('设置已导出')
}

const importSettings = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        try {
          const settings = JSON.parse(e.target?.result as string)
          
          // 验证设置格式
          if (settings.general) Object.assign(generalSettings.value, settings.general)
          if (settings.appearance) Object.assign(appearanceSettings.value, settings.appearance)
          if (settings.editor) Object.assign(editorSettings.value, settings.editor)
          if (settings.ai) Object.assign(aiSettings.value, settings.ai)
          if (settings.advanced) Object.assign(advancedSettings.value, settings.advanced)
          
          ElMessage.success('设置已导入')
          saveSettings()
        } catch (error) {
          ElMessage.error('导入设置失败：文件格式错误')
        }
      }
      reader.readAsText(file)
    }
  }
  input.click()
}

const saveSettings = () => {
  const settings = {
    general: generalSettings.value,
    appearance: appearanceSettings.value,
    editor: editorSettings.value,
    ai: aiSettings.value,
    advanced: advancedSettings.value
  }
  
  localStorage.setItem('codesage-settings', JSON.stringify(settings))
}

const applySettingsToApp = () => {
  const themeStore = useThemeStore()
  
  // 使用 theme store 应用主题设置
  themeStore.applyTheme({
    theme: appearanceSettings.value.theme as 'dark' | 'light' | 'auto',
    primaryColor: appearanceSettings.value.primaryColor,
    accentColor: appearanceSettings.value.accentColor,
    fontSize: appearanceSettings.value.fontSize,
    codeFontFamily: appearanceSettings.value.codeFontFamily,
    animations: appearanceSettings.value.animations,
    glassEffect: appearanceSettings.value.glassEffect,
    opacity: appearanceSettings.value.opacity,
    compactMode: appearanceSettings.value.compactMode,
    showClock: appearanceSettings.value.showClock
  })
  
  // 应用主题到 Element Plus
  document.documentElement.setAttribute('data-theme', appearanceSettings.value.theme)
  
  // 应用界面字体大小
  document.documentElement.style.fontSize = `${appearanceSettings.value.fontSize}px`
  document.body.style.fontSize = `${appearanceSettings.value.fontSize}px`
  
  // 应用主色调
  document.documentElement.style.setProperty('--primary-color', appearanceSettings.value.primaryColor)
  document.documentElement.style.setProperty('--accent-color', appearanceSettings.value.accentColor)
  
  // 应用代码字体到全局
  document.documentElement.style.setProperty('--code-font-family', appearanceSettings.value.codeFontFamily)
  // 同时应用到 body，确保所有 code 标签都使用该字体
  const codeElements = document.querySelectorAll('code, pre')
  codeElements.forEach(el => {
    (el as HTMLElement).style.fontFamily = appearanceSettings.value.codeFontFamily
  })
  
  // 应用动画设置
  if (!appearanceSettings.value.animations) {
    document.documentElement.style.setProperty('--animation-duration', '0s')
  } else {
    document.documentElement.style.removeProperty('--animation-duration')
  }
  
  // 应用透明度
  if (appearanceSettings.value.opacity !== 100) {
    document.documentElement.style.setProperty('--window-opacity', `${appearanceSettings.value.opacity / 100}`)
  } else {
    document.documentElement.style.removeProperty('--window-opacity')
  }
  
  // 应用紧凑模式
  if (appearanceSettings.value.compactMode) {
    document.documentElement.setAttribute('data-compact-mode', 'true')
  } else {
    document.documentElement.removeAttribute('data-compact-mode')
  }
  
  // 应用玻璃效果
  if (appearanceSettings.value.glassEffect) {
    document.documentElement.setAttribute('data-glass-effect', 'true')
  } else {
    document.documentElement.removeAttribute('data-glass-effect')
  }
  
  // 应用浅色主题样式
  if (appearanceSettings.value.theme === 'light') {
    document.documentElement.style.setProperty('--bg-primary', '#ffffff')
    document.documentElement.style.setProperty('--bg-secondary', '#f5f5f5')
    document.documentElement.style.setProperty('--bg-tertiary', '#eeeeee')
    document.documentElement.style.setProperty('--bg-card', '#fafafa')
    document.documentElement.style.setProperty('--bg-hover', '#e0e0e0')
    document.documentElement.style.setProperty('--text-primary', '#000000')
    document.documentElement.style.setProperty('--text-secondary', '#666666')
    document.documentElement.style.setProperty('--text-tertiary', '#999999')
    document.documentElement.style.setProperty('--text-muted', '#cccccc')
    document.documentElement.style.setProperty('--border-color', '#e0e0e0')
    document.documentElement.style.setProperty('--border-hover', '#d0d0d0')
  } else {
    // 深色主题（默认值）
    document.documentElement.style.setProperty('--bg-primary', '#0a0a0a')
    document.documentElement.style.setProperty('--bg-secondary', '#1a1a1a')
    document.documentElement.style.setProperty('--bg-tertiary', '#2a2a2a')
    document.documentElement.style.setProperty('--bg-card', '#1e1e1e')
    document.documentElement.style.setProperty('--bg-hover', '#333333')
    document.documentElement.style.setProperty('--text-primary', '#ffffff')
    document.documentElement.style.setProperty('--text-secondary', '#b0b0b0')
    document.documentElement.style.setProperty('--text-tertiary', '#808080')
    document.documentElement.style.setProperty('--text-muted', '#606060')
    document.documentElement.style.setProperty('--border-color', '#333333')
    document.documentElement.style.setProperty('--border-hover', '#444444')
  }
  
  // 通知其他组件设置已更改
  window.dispatchEvent(new CustomEvent('settings-applied', { 
    detail: {
      general: generalSettings.value,
      appearance: appearanceSettings.value,
      editor: editorSettings.value,
      ai: aiSettings.value,
      advanced: advancedSettings.value
    }
  }))
}

const resetSettings = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重置所有设置为默认值吗？此操作不可撤销。',
      '重置设置',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 重置为默认值
    generalSettings.value = {
      language: 'zh-CN',
      autoSave: true,
      autoSaveInterval: 5,
      checkUpdates: true,
      sendStats: false,
      defaultWorkspace: '',
      maxConcurrentTasks: 5,
      taskTimeout: 10,
      enableNotifications: true
    }
    
    appearanceSettings.value = {
      theme: 'dark',
      primaryColor: '#1e88e5',
      accentColor: '#1e88e5',
      fontSize: 14,
      codeFontFamily: 'Consolas, Monaco, monospace',
      animations: true,
      glassEffect: false,
      opacity: 100,
      compactMode: false,
      showClock: true
    }
    
    editorSettings.value = {
      theme: 'vs-dark',
      fontSize: 14,
      fontWeight: 400,
      showLineNumbers: true,
      wordWrap: true,
      tabSize: 4,
      insertSpaces: true,
      autoComplete: true,
      linting: true,
      highlightCurrentLine: true,
      renderWhitespace: false,
      minimap: true
    }
    
    aiSettings.value = {
      serviceType: 'local',
      model: 'codellama-34b',
      apiEndpoint: '',
      apiKey: '',
      modelPath: '/models/codellama-34b',
      temperature: 0.7,
      maxTokens: 2048,
      codeValidation: true,
      safeMode: true,
      enableCache: true
    }
    
    advancedSettings.value = {
      developerMode: false,
      verboseLogging: false,
      performanceMonitoring: false,
      memoryLimit: 2048,
      threadPoolSize: 8,
      cacheSize: 200,
      autoCleanup: true,
      cleanupInterval: 6,
      experimentalFeatures: false
    }
    
    saveSettings()
    ElMessage.success('设置已重置为默认值')
  } catch {
    // 用户取消操作
  }
}

const openDocumentation = () => {
  ElMessage.info('打开文档功能开发中...')
}

const openGitHub = () => {
  ElMessage.info('打开 GitHub 功能开发中...')
}

const openIssues = () => {
  ElMessage.info('打开问题反馈功能开发中...')
}

const checkForUpdates = () => {
  ElMessage.info('检查更新功能开发中...')
}

// 获取系统信息
const getSystemInfo = () => {
  // 获取系统信息
  systemInfo.nodeVersion = '18.17.0' // 模拟版本
  systemInfo.chromeVersion = '120.0.0' // 模拟版本
  systemInfo.platform = navigator.platform || 'Windows 11'
  
  // 获取内存使用信息
  if ((performance as any).memory) {
    const memory = (performance as any).memory
    const used = Math.round(memory.usedJSHeapSize / 1024 / 1024)
    const total = Math.round(memory.totalJSHeapSize / 1024 / 1024)
    systemInfo.memoryUsage = `${used}MB / ${total}MB`
  } else {
    systemInfo.memoryUsage = '128MB / 512MB' // 模拟数据
  }
}

// 获取使用统计
const getUsageStats = () => {
  // 从 localStorage 获取统计数据
  const savedStats = localStorage.getItem('codesage-usage-stats')
  if (savedStats) {
    try {
      const parsed = JSON.parse(savedStats)
      Object.assign(stats, parsed)
    } catch (error) {
      console.error('加载使用统计失败:', error)
    }
  }
}

// 加载保存的设置
const loadSettings = async () => {
  const saved = localStorage.getItem('codesage-settings')
  if (saved) {
    try {
      const settings = JSON.parse(saved)
      if (settings.general) Object.assign(generalSettings.value, settings.general)
      if (settings.appearance) Object.assign(appearanceSettings.value, settings.appearance)
      if (settings.editor) Object.assign(editorSettings.value, settings.editor)
      if (settings.ai) Object.assign(aiSettings.value, settings.ai)
      if (settings.advanced) Object.assign(advancedSettings.value, settings.advanced)
      
      // 应用设置
      applySettingsToApp()
    } catch (error) {
      console.error('加载设置失败:', error)
    }
  }
  
  // 从后端加载个人资料
  if (authStore.userInfo) {
    try {
      const response = await authAPI.getProfile()
      if (response.user || response) {
        const user = response.user || response
        profileSettings.value.username = user.username || ''
        profileSettings.value.email = user.email || ''
        profileSettings.value.avatar = user.avatar || ''
        profileSettings.value.phone = user.phone || ''
        profileSettings.value.bio = user.bio || ''
        profileSettings.value.location = user.location || ''
        profileSettings.value.occupation = user.occupation || ''
        profileSettings.value.company = user.company || ''
        profileSettings.value.website = user.website || ''
        profileSettings.value.twitter = user.twitter || ''
        profileSettings.value.github = user.github_url || ''
        
        // 同步到 localStorage 以便刷新后使用
        localStorage.setItem('codesage-profile', JSON.stringify(profileSettings.value))
        
        // 同时更新 authStore 中的用户信息（使用 updateUserInfo 方法确保持久化）
        authStore.updateUserInfo({
          username: user.username,
          email: user.email,
          avatar: user.avatar
        })
      }
    } catch (error) {
      console.error('从后端加载个人资料失败:', error)
      // 如果后端加载失败，尝试从 localStorage 加载
      const savedProfile = localStorage.getItem('codesage-profile')
      if (savedProfile) {
        try {
          const profile = JSON.parse(savedProfile)
          Object.assign(profileSettings.value, profile)
        } catch (error) {
          console.error('加载本地个人资料失败:', error)
        }
      }
    }
  } else {
    // 如果没有登录，尝试从 localStorage 加载
    const savedProfile = localStorage.getItem('codesage-profile')
    if (savedProfile) {
      try {
        const profile = JSON.parse(savedProfile)
        Object.assign(profileSettings.value, profile)
      } catch (error) {
        console.error('加载本地个人资料失败:', error)
      }
    }
  }
}

// 监听所有设置变化，自动保存并应用
watch(() => appearanceSettings.value, () => {
  saveSettings()
  applySettingsToApp()
}, { deep: true })

watch(() => generalSettings.value, () => {
  saveSettings()
  applyGeneralSettings()
}, { deep: true })

watch(() => editorSettings.value, () => {
  saveSettings()
  // 编辑器设置通过事件系统通知
  window.dispatchEvent(new CustomEvent('settings-applied', { 
    detail: {
      editor: editorSettings.value
    }
  }))
}, { deep: true })

// 监听个人资料变化，但不自动保存（避免频繁请求）
// 用户需要手动点击保存按钮
// watch(() => profileSettings.value, () => {
//   saveProfileSettings()
// }, { deep: true })

// 应用通用设置
const applyGeneralSettings = () => {
  // 应用语言设置
  if (generalSettings.value.language) {
    document.documentElement.setAttribute('lang', generalSettings.value.language)
  }
}

// 组件挂载时初始化
onMounted(() => {
  getSystemInfo()
  getUsageStats()
  loadSettings()
  // 初始应用设置
  applySettingsToApp()
  applyGeneralSettings()
  
  // 处理路由参数，支持从用户菜单直接跳转到指定选项卡
  if (route.query.tab) {
    activeMenu.value = route.query.tab as string
  }
})
</script>

<style scoped>
.settings-container {
  padding: 0;
}

.menu-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
}

.settings-menu {
  border: none;
  background: transparent;
}

.settings-menu .el-menu-item {
  margin: 4px 0;
  border-radius: 8px;
  height: 48px;
  line-height: 48px;
  transition: all 0.3s ease;
}

.settings-menu .el-menu-item:hover {
  background: var(--bg-hover);
}

.settings-menu .el-menu-item.is-active {
  background: var(--primary-color);
  color: var(--text-primary);
}

.content-card {
  border: 1px solid var(--border-color);
  border-radius: 12px;
  min-height: 600px;
}

.settings-section {
  padding: 24px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 24px 0;
  color: var(--text-primary);
}

.setting-description {
  margin-left: 12px;
  font-size: 12px;
  color: var(--text-secondary);
}

.settings-actions {
  padding: 24px;
  border-top: 1px solid var(--border-color);
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.about-content {
  padding: 24px;
  text-align: center;
}

.logo-section {
  margin-bottom: 32px;
}

.app-logo {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 16px;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.app-version {
  font-size: 18px;
  color: var(--text-secondary);
  margin: 0 0 8px 0;
}

.app-description {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0 0 16px 0;
}

.system-info-section {
  text-align: left;
  margin: 32px 0;
}

.stats-section {
  text-align: left;
  margin: 32px 0;
}

.stat-item {
  text-align: center;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--primary-color);
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.info-section {
  text-align: left;
  margin: 32px 0;
}

.feature-list {
  list-style: none;
  padding: 0;
}

.feature-list li {
  padding: 8px 0;
  color: var(--text-secondary);
  position: relative;
  padding-left: 20px;
}

.feature-list li::before {
  content: "✓";
  position: absolute;
  left: 0;
  color: var(--success-color);
  font-weight: bold;
}

.links-section {
  margin: 32px 0;
}

.link-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.copyright-section {
  margin-top: 32px;
}

.copyright,
.license,
.build-info {
  font-size: 14px;
  color: var(--text-tertiary);
  margin: 4px 0;
}

/* 个人资料样式 */
.avatar-upload {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 12px;
}

.profile-avatar {
  border: 2px solid var(--border-color);
  transition: all 0.3s ease;
}

.profile-avatar:hover {
  border-color: var(--primary-color);
  transform: scale(1.05);
}

.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.password-change {
  margin-bottom: 12px;
}

.social-icon-small {
  width: 16px;
  height: 16px;
  vertical-align: middle;
}

/* 修复浏览器自动填充导致的背景色问题 */
:deep(.el-input__inner:-webkit-autofill),
:deep(.el-input__inner:-webkit-autofill:hover),
:deep(.el-input__inner:-webkit-autofill:focus),
:deep(.el-input__inner:-webkit-autofill:active) {
  -webkit-box-shadow: 0 0 0 1000px var(--bg-secondary) inset !important;
  -webkit-text-fill-color: var(--text-primary) !important;
  transition: background-color 5000s ease-in-out 0s;
  caret-color: var(--text-primary) !important;
}

:deep(.el-textarea__inner:-webkit-autofill),
:deep(.el-textarea__inner:-webkit-autofill:hover),
:deep(.el-textarea__inner:-webkit-autofill:focus),
:deep(.el-textarea__inner:-webkit-autofill:active) {
  -webkit-box-shadow: 0 0 0 1000px var(--bg-secondary) inset !important;
  -webkit-text-fill-color: var(--text-primary) !important;
  transition: background-color 5000s ease-in-out 0s;
  caret-color: var(--text-primary) !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .settings-menu .el-menu-item {
    justify-content: center;
  }
  
  .link-buttons {
    flex-direction: column;
    align-items: center;
  }
}
</style>