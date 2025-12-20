/**
 * 自动保存工具
 * 用于自动保存编辑器内容到 localStorage
 */

export interface AutoSaveOptions {
  key: string
  interval: number // 分钟
  enabled: boolean
}

export class AutoSave {
  private timerId: number | null = null
  private getValue: () => string
  private options: AutoSaveOptions

  constructor(getValue: () => string, options: AutoSaveOptions) {
    this.getValue = getValue
    this.options = options
  }

  /**
   * 启动自动保存
   */
  start() {
    if (!this.options.enabled) return
    
    this.stop() // 先停止已存在的定时器
    
    // 设置定时器
    this.timerId = window.setInterval(() => {
      this.save()
    }, this.options.interval * 60 * 1000) // 转换为毫秒
  }

  /**
   * 停止自动保存
   */
  stop() {
    if (this.timerId !== null) {
      clearInterval(this.timerId)
      this.timerId = null
    }
  }

  /**
   * 立即保存
   */
  save() {
    try {
      const content = this.getValue()
      if (content) {
        localStorage.setItem(this.options.key, content)
        console.log(`Auto-saved to ${this.options.key}`)
      }
    } catch (error) {
      console.error('Auto-save failed:', error)
    }
  }

  /**
   * 加载保存的内容
   */
  load(): string | null {
    try {
      return localStorage.getItem(this.options.key)
    } catch (error) {
      console.error('Auto-load failed:', error)
      return null
    }
  }

  /**
   * 更新选项
   */
  updateOptions(options: Partial<AutoSaveOptions>) {
    this.options = { ...this.options, ...options }
    
    // 如果更改了enabled或interval，需要重新启动
    if (this.options.enabled) {
      this.start()
    } else {
      this.stop()
    }
  }
}

/**
 * 从设置中加载自动保存配置
 */
export function loadAutoSaveSettings() {
  try {
    const saved = localStorage.getItem('codesage-settings')
    if (saved) {
      const settings = JSON.parse(saved)
      if (settings.general) {
        return {
          enabled: settings.general.autoSave !== false,
          interval: settings.general.autoSaveInterval || 5
        }
      }
    }
  } catch (error) {
    console.error('Failed to load auto-save settings:', error)
  }
  
  return {
    enabled: true,
    interval: 5
  }
}

/**
 * 监听自动保存设置变更
 */
export function listenToAutoSaveChanges(callback: (options: { enabled: boolean; interval: number }) => void) {
  const handleSettingsApplied = (event: CustomEvent) => {
    if (event.detail?.general) {
      callback({
        enabled: event.detail.general.autoSave !== false,
        interval: event.detail.general.autoSaveInterval || 5
      })
    }
  }
  
  window.addEventListener('settings-applied', handleSettingsApplied as EventListener)
  
  // 返回清理函数
  return () => {
    window.removeEventListener('settings-applied', handleSettingsApplied as EventListener)
  }
}
