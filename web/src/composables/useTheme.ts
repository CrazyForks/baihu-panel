import { ref, watch, onMounted } from 'vue'

export type Theme = 'light' | 'dark' | 'system'

const theme = ref<Theme>('system')

function getSystemTheme(): 'light' | 'dark' {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function applyTheme(t: Theme) {
  const root = document.documentElement
  const effectiveTheme = t === 'system' ? getSystemTheme() : t

  root.classList.remove('light', 'dark')
  root.classList.add(effectiveTheme)
}

export function useTheme() {
  onMounted(() => {
    // 从 localStorage 读取主题
    const saved = localStorage.getItem('theme') as Theme | null
    if (saved && ['light', 'dark', 'system'].includes(saved)) {
      theme.value = saved
    }
    applyTheme(theme.value)

    // 监听系统主题变化
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
      if (theme.value === 'system') {
        applyTheme('system')
      }
    })
  })

  watch(theme, (newTheme) => {
    localStorage.setItem('theme', newTheme)
    applyTheme(newTheme)
  })

  function setTheme(t: Theme) {
    theme.value = t
  }

  return {
    theme,
    setTheme
  }
}
