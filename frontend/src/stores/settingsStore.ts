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
    } catch {
      set({ loading: false })
    }
  },

  saveSettings: async (settings: SettingsDto) => {
    await settingsApi.saveSettings(settings)
    set({ settings })
  },

  updateLayout: (patch) => {
    const current = get().settings
    if (!current) return
    const updated = { ...current, ...patch }
    set({ settings: updated })
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
