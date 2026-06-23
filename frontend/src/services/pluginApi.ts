/**
 * Plugin API 服务层
 *
 * 封装 Wails 后端 PluginApi 的所有调用。
 * 组件不得直接调用 window.runtime 或 window.go。
 */

import type {
  PluginSummaryDto,
  PluginDto,
  ReloadPluginsResponse,
  ValidationResultDto,
  PluginHistoryEntryDto,
} from '../types'

// Wails v2 绑定：window.go.<package>.<struct>.<method>
const backend = () => (window as any).go.wails.PluginApi

export const pluginApi = {
  async getPlugins(): Promise<PluginSummaryDto[]> {
    return backend().GetPlugins()
  },

  async getPlugin(pluginId: string): Promise<PluginDto> {
    return backend().GetPlugin(pluginId)
  },

  async reloadPlugins(): Promise<ReloadPluginsResponse> {
    return backend().ReloadPlugins()
  },

  async validateFormData(
    pluginId: string,
    data: Record<string, unknown>,
  ): Promise<ValidationResultDto> {
    return backend().ValidateFormData(pluginId, data)
  },

  async getHistory(pluginId: string): Promise<PluginHistoryEntryDto[]> {
    return backend().GetPluginHistory(pluginId)
  },

  async saveHistory(
    pluginId: string,
    data: Record<string, unknown>,
  ): Promise<void> {
    return backend().SavePluginHistory(pluginId, data)
  },
}
