import { describe, it, expect, beforeEach, vi } from 'vitest'
import { usePluginStore } from './pluginStore'
import { pluginApi } from '../services/pluginApi'

vi.mock('../services/pluginApi', () => ({
  pluginApi: {
    getPlugins: vi.fn(),
    getPlugin: vi.fn(),
    reloadPlugins: vi.fn(),
    validateFormData: vi.fn(),
  },
}))

const mockPlugins = [
  {
    id: 'tia-export',
    name: 'TIA Portal Export',
    description: 'Export content from Siemens TIA Portal projects',
    version: '1.0.0',
    navigation: { group: ['Industrial Automation', 'TIA Portal'], order: 1 },
  },
  {
    id: 'other-plugin',
    name: 'Other Tool',
    navigation: { group: ['Utilities'], order: 2 },
  },
]

describe('pluginStore', () => {
  beforeEach(() => {
    // Reset store state between tests
    usePluginStore.setState({
      plugins: [],
      pluginCache: new Map(),
      loading: false,
      error: null,
    })
    vi.clearAllMocks()
  })

  it('initial state has empty plugins and no error', () => {
    const state = usePluginStore.getState()
    expect(state.plugins).toEqual([])
    expect(state.loading).toBe(false)
    expect(state.error).toBeNull()
  })

  it('fetchPlugins loads plugins into state', async () => {
    vi.mocked(pluginApi.getPlugins).mockResolvedValue(mockPlugins)

    await usePluginStore.getState().fetchPlugins()

    const state = usePluginStore.getState()
    expect(state.plugins).toHaveLength(2)
    expect(state.plugins[0].id).toBe('tia-export')
    expect(state.plugins[1].id).toBe('other-plugin')
    expect(state.loading).toBe(false)
    expect(state.error).toBeNull()
  })

  it('fetchPlugins sets loading state correctly', async () => {
    vi.mocked(pluginApi.getPlugins).mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve(mockPlugins), 50)),
    )

    const fetchPromise = usePluginStore.getState().fetchPlugins()

    // During fetch, loading should be true
    expect(usePluginStore.getState().loading).toBe(true)

    await fetchPromise

    // After fetch, loading should be false
    expect(usePluginStore.getState().loading).toBe(false)
  })

  it('fetchPlugins handles errors gracefully', async () => {
    vi.mocked(pluginApi.getPlugins).mockRejectedValue(new Error('API timeout'))

    await usePluginStore.getState().fetchPlugins()

    const state = usePluginStore.getState()
    expect(state.plugins).toEqual([])
    expect(state.loading).toBe(false)
    expect(state.error).toBe('API timeout')
  })

  it('fetchPlugins handles empty API response', async () => {
    vi.mocked(pluginApi.getPlugins).mockResolvedValue([])

    await usePluginStore.getState().fetchPlugins()

    const state = usePluginStore.getState()
    expect(state.plugins).toEqual([])
    expect(state.error).toBeNull()
  })

  it('fetchPlugin caches results', async () => {
    const pluginDetail = { metadata: { id: 'tia-export', name: 'TIA Portal Export' } }
    vi.mocked(pluginApi.getPlugin).mockResolvedValue(pluginDetail as any)

    // First call should hit the API
    const result1 = await usePluginStore.getState().fetchPlugin('tia-export')
    expect(result1).toEqual(pluginDetail)
    expect(pluginApi.getPlugin).toHaveBeenCalledTimes(1)

    // Second call should use cache, not hit API
    const result2 = await usePluginStore.getState().fetchPlugin('tia-export')
    expect(result2).toEqual(pluginDetail)
    expect(pluginApi.getPlugin).toHaveBeenCalledTimes(1) // Still 1
  })

  it('reloadPlugins reloads and returns count', async () => {
    vi.mocked(pluginApi.reloadPlugins).mockResolvedValue({ count: 2 })
    vi.mocked(pluginApi.getPlugins).mockResolvedValue(mockPlugins)

    const count = await usePluginStore.getState().reloadPlugins()
    expect(pluginApi.reloadPlugins).toHaveBeenCalledOnce()
    expect(count).toBe(2)
  })

  it('reloadPlugins handles failures gracefully', async () => {
    vi.mocked(pluginApi.reloadPlugins).mockRejectedValue(new Error('Backend offline'))

    const count = await usePluginStore.getState().reloadPlugins()
    expect(count).toBe(0)
    const state = usePluginStore.getState()
    expect(state.error).toBe('Backend offline')
  })
})
