/**
 * Monaco Editor 全局配置管理工具
 * 用于在所有编辑器实例中应用统一的配置
 */

export interface EditorSettings {
  theme: string
  fontSize: number
  fontWeight: number
  showLineNumbers: boolean
  wordWrap: boolean
  tabSize: number
  insertSpaces: boolean
  autoComplete: boolean
  linting: boolean
  highlightCurrentLine: boolean
  renderWhitespace: boolean
  minimap: boolean
}

export const defaultEditorSettings: EditorSettings = {
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

/**
 * 将设置转换为 Monaco Editor 配置选项
 */
export function toMonacoOptions(settings: Partial<EditorSettings> = {}) {
  const config = { ...defaultEditorSettings, ...settings }
  
  return {
    theme: config.theme,
    fontSize: config.fontSize,
    fontWeight: config.fontWeight.toString(),
    lineNumbers: config.showLineNumbers ? 'on' : 'off',
    wordWrap: config.wordWrap ? 'on' : 'off',
    tabSize: config.tabSize,
    insertSpaces: config.insertSpaces,
    quickSuggestions: config.autoComplete,
    renderLineHighlight: config.highlightCurrentLine ? 'all' : 'none',
    renderWhitespace: config.renderWhitespace ? 'all' : 'none',
    minimap: {
      enabled: config.minimap
    },
    automaticLayout: true,
    scrollBeyondLastLine: false,
    roundedSelection: false
  }
}

/**
 * 从 localStorage 加载编辑器设置
 */
export function loadEditorSettings(): EditorSettings {
  try {
    const saved = localStorage.getItem('codesage-settings')
    if (saved) {
      const settings = JSON.parse(saved)
      if (settings.editor) {
        return { ...defaultEditorSettings, ...settings.editor }
      }
    }
  } catch (error) {
    console.error('Failed to load editor settings:', error)
  }
  return defaultEditorSettings
}

/**
 * 更新编辑器实例的配置
 */
export function updateEditorOptions(editor: any, settings: Partial<EditorSettings>) {
  if (!editor) return
  
  const options = toMonacoOptions(settings)
  editor.updateOptions(options)
}

/**
 * 监听全局设置变更事件
 */
export function listenToSettingsChanges(callback: (settings: EditorSettings) => void) {
  const handleSettingsApplied = (event: CustomEvent) => {
    if (event.detail?.editor) {
      callback(event.detail.editor)
    }
  }
  
  window.addEventListener('settings-applied', handleSettingsApplied as EventListener)
  
  // 返回清理函数
  return () => {
    window.removeEventListener('settings-applied', handleSettingsApplied as EventListener)
  }
}
