import { create } from 'zustand'
import type { PluginSummaryDto, PluginDto } from '../types'
import { pluginApi } from '../services/pluginApi'

interface PluginState {
  plugins: PluginSummaryDto[]
  pluginCache: Map<string, PluginDto>
  loading: boolean
  error: string | null

  fetchPlugins: () => Promise<void>
  fetchPlugin: (id: string) => Promise<PluginDto>
  reloadPlugins: () => Promise<number>
}

export const usePluginStore = create<PluginState>((set, get) => ({
  plugins: [],
  pluginCache: new Map(),
  loading: false,
  error: null,

  fetchPlugins: async () => {
    set({ loading: true, error: null })
    try {
      const plugins = await pluginApi.getPlugins()
      set({ plugins, loading: false })
    } catch (err: any) {
      set({ error: err?.message ?? 'Failed to load plugins', loading: false })
    }
  },

  fetchPlugin: async (id: string) => {
    const cached = get().pluginCache.get(id)
    if (cached) return cached

    const plugin = await pluginApi.getPlugin(id)
    set((state) => {
      const next = new Map(state.pluginCache)
      next.set(id, plugin)
      return { pluginCache: next }
    })
    return plugin
  },

  reloadPlugins: async () => {
    try {
      // 清除缓存，确保重新加载后获取最新数据
      set({ pluginCache: new Map() })
      const resp = await pluginApi.reloadPlugins()
      await get().fetchPlugins()
      return resp.count
    } catch (err: any) {
      set({ error: err?.message ?? 'Failed to reload plugins' })
      return 0
    }
  },
}))
