import { useEffect } from 'react'
import { useSettingsStore } from '../stores/settingsStore'

export function useTheme() {
  const theme = useSettingsStore((s) => s.settings?.theme ?? 'dark')

  useEffect(() => {
    const root = document.documentElement
    applyTheme(theme)
    if (theme !== 'system') return

    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    const handler = (e: MediaQueryListEvent) => applySystemTheme(e.matches)
    mq.addEventListener('change', handler)
    return () => mq.removeEventListener('change', handler)
  }, [theme])
}

function applyTheme(theme: string) {
  const root = document.documentElement
  if (theme === 'dark') {
    root.classList.add('dark')
  } else if (theme === 'light') {
    root.classList.remove('dark')
  } else {
    applySystemTheme(window.matchMedia('(prefers-color-scheme: dark)').matches)
  }
}

function applySystemTheme(prefersDark: boolean) {
  document.documentElement.classList.toggle('dark', prefersDark)
}
