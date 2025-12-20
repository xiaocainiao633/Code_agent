import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface ThemeSettings {
  theme: 'dark' | 'light' | 'auto'
  primaryColor: string
  accentColor: string
  fontSize: number
  codeFontFamily: string
  animations: boolean
  glassEffect: boolean
  opacity: number
  compactMode: boolean
  showClock: boolean
}

export const useThemeStore = defineStore('theme', () => {
  // 默认主题设置
  const defaultThemeSettings: ThemeSettings = {
    theme: 'dark',
    primaryColor: '#1e88e5',
    accentColor: '#ff6b35',
    fontSize: 14,
    codeFontFamily: 'Consolas, Monaco, monospace',
    animations: true,
    glassEffect: false,
    opacity: 100,
    compactMode: false,
    showClock: true
  }

  // 当前主题设置
  const themeSettings = ref<ThemeSettings>({ ...defaultThemeSettings })

  // 计算属性 - 当前主题模式
  const currentTheme = computed(() => {
    if (themeSettings.value.theme === 'auto') {
      // 根据系统主题自动切换
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    }
    return themeSettings.value.theme
  })

  // 计算属性 - 是否深色模式
  const isDarkMode = computed(() => currentTheme.value === 'dark')

  // 应用主题到DOM
  const applyTheme = (settings: Partial<ThemeSettings> = {}) => {
    // 合并设置
    const newSettings = { ...themeSettings.value, ...settings }
    themeSettings.value = newSettings

    // 保存到localStorage
    localStorage.setItem('codesage-theme-settings', JSON.stringify(newSettings))

    // 应用主题
    applyThemeToDOM(newSettings)
  }

  // 应用主题到DOM元素
  const applyThemeToDOM = (settings: ThemeSettings) => {
    const root = document.documentElement

    // 设置主题属性
    root.setAttribute('data-theme', settings.theme)

    // 应用CSS变量
    root.style.setProperty('--primary-color', settings.primaryColor)
    root.style.setProperty('--accent-color', settings.accentColor)
    root.style.setProperty('--font-size', `${settings.fontSize}px`)

    // 应用字体
    if (settings.codeFontFamily) {
      root.style.setProperty('--code-font-family', settings.codeFontFamily)
    }

    // 应用动画设置
    if (!settings.animations) {
      root.style.setProperty('--animation-duration', '0s')
    } else {
      root.style.removeProperty('--animation-duration')
    }

    // 应用透明度
    if (settings.opacity !== 100) {
      root.style.setProperty('--window-opacity', `${settings.opacity}%`)
    } else {
      root.style.removeProperty('--window-opacity')
    }

    // 应用紧凑模式
    if (settings.compactMode) {
      root.setAttribute('data-compact-mode', 'true')
    } else {
      root.removeAttribute('data-compact-mode')
    }

    // 应用玻璃效果
    if (settings.glassEffect) {
      root.setAttribute('data-glass-effect', 'true')
    } else {
      root.removeAttribute('data-glass-effect')
    }

    // 触发主题变更事件
    window.dispatchEvent(new CustomEvent('theme-changed', {
      detail: { theme: settings }
    }))
  }

  // 从localStorage加载主题设置
  const loadThemeFromStorage = () => {
    try {
      const saved = localStorage.getItem('codesage-theme-settings')
      if (saved) {
        const settings = JSON.parse(saved) as ThemeSettings
        // 验证设置完整性
        const validSettings = { ...defaultThemeSettings, ...settings }
        themeSettings.value = validSettings
        applyThemeToDOM(validSettings)
      } else {
        // 使用默认主题
        applyThemeToDOM(defaultThemeSettings)
      }
    } catch (error) {
      console.error('Failed to load theme settings:', error)
      // 使用默认主题
      applyThemeToDOM(defaultThemeSettings)
    }
  }

  // 重置为默认主题
  const resetTheme = () => {
    applyTheme(defaultThemeSettings)
  }

  // 监听系统主题变化
  const watchSystemTheme = () => {
    if (themeSettings.value.theme === 'auto') {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
      mediaQuery.addEventListener('change', () => {
        // 主题变化时重新应用
        applyTheme()
      })
    }
  }

  // 初始化主题
  const initTheme = () => {
    loadThemeFromStorage()
    watchSystemTheme()
  }

  // 生成主题样式对象（供动态样式使用）
  const themeStyles = computed(() => {
    const settings = themeSettings.value
    return {
      '--primary-color': settings.primaryColor,
      '--accent-color': settings.accentColor,
      '--font-size': `${settings.fontSize}px`,
      '--code-font-family': settings.codeFontFamily,
      ...(settings.opacity !== 100 && { '--window-opacity': `${settings.opacity}%` }),
      ...(settings.compactMode && { '--compact-mode': 'true' }),
      ...(settings.glassEffect && { '--glass-effect': 'true' }),
      ...(settings.animations && { '--animation-duration': '0s' })
    }
  })

  return {
    themeSettings,
    currentTheme,
    isDarkMode,
    themeStyles,
    applyTheme,
    resetTheme,
    initTheme,
    loadThemeFromStorage
  }
})