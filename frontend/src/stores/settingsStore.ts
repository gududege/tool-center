import { create } from 'zustand'
import type { SettingsDto } from '../types'
import { settingsApi } from '../services/settingsApi'

interface SettingsState {
  settings: SettingsDto | null
  loading: boolean

  fetchSettings: () => Promise<void>
  saveSettings: (settings: SettingsDto) => Promise<void>
  updateLayout: (patch: Partial<Pick<SettingsDto, 'sidebarCollapsed' | 'sidebarSize' | 'bottomPanelSize' | 'bottomTab'>>) => void
  setLanguage: (lang: string) => Promise<void>
}

export const useSettingsStore = create<SettingsState>((set, get) => ({
  settings: null,
  loading: false,

  fetchSettings: async () => {
    set({ loading: true })
    try {
      const settings = await settingsApi.getSettings()
      set({ settings, loading: false })
      applyTheme(settings.theme)
    } catch {
      set({ loading: false })
    }
  },

  saveSettings: async (settings: SettingsDto) => {
    await settingsApi.saveSettings(settings)
    set({ settings })
    applyTheme(settings.theme)
  },

  updateLayout: (patch) => {
    const current = get().settings
    if (!current) return
    const updated = { ...current, ...patch }
    set({ settings: updated })
    // 异步保存不阻塞 UI
    settingsApi.saveSettings(updated).catch(() => {})
  },

  setLanguage: async (lang: string) => {
    const current = get().settings
    if (!current) return
    const updated = { ...current, language: lang }
    set({ settings: updated })
    await settingsApi.saveSettings(updated)
  },
}))

function applyTheme(theme: string) {
  if (theme === 'dark') {
    document.documentElement.classList.add('dark')
  } else if (theme === 'light') {
    document.documentElement.classList.remove('dark')
  } else {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    document.documentElement.classList.toggle('dark', prefersDark)
  }
}
