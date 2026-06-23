/**
 * Settings API 服务层
 */

import type { SettingsDto } from '../types'

const backend = () => (window as any).go.wails.SettingsApi

export const settingsApi = {
  async getSettings(): Promise<SettingsDto> {
    return backend().GetSettings()
  },

  async saveSettings(settings: SettingsDto): Promise<void> {
    return backend().SaveSettings(settings)
  },
}
