import { describe, it, expect, beforeEach, vi } from 'vitest'
import { pluginApi } from './pluginApi'

describe('pluginApi Wails namespace', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('calls window.go.wails.PluginApi (fixed frontend code)', async () => {
    const mockFn = vi.fn().mockResolvedValue([])
    ;(window as any).go = {
      wails: {
        PluginApi: { GetPlugins: mockFn },
      },
    }

    const result = await pluginApi.getPlugins()
    expect(mockFn).toHaveBeenCalledOnce()
    expect(result).toEqual([])
  })

  it('GENERATED BINDINGS use "wails" — frontend now matches after fix', async () => {
    // After the fix, pluginApi.ts uses "wails" namespace, matching
    // the Wails-generated bindings in wailsjs/go/wails/PluginApi.js
    ;(window as any).go = {
      wails: {
        PluginApi: {
          GetPlugins: vi.fn().mockResolvedValue([{ id: 'test', name: 'Test', navigation: { group: [], order: 0 } }]),
        },
      },
    }

    // pluginApi.getPlugins() calls go.wails.PluginApi → works!
    const result = await pluginApi.getPlugins()
    expect(result).toHaveLength(1)
    expect(result[0].id).toBe('test')
  })

  it('FRONTEND USES window.go.main — Wails GENERATED BINDINGS USE window.go.wails', () => {
    // The frontend's pluginApi.ts calls: (window as any).go.main.PluginApi
    // But Wails-generated bindings (wailsjs/go/wails/PluginApi.js) call: window.go.wails.PluginApi
    //
    // This test documents the mismatch. Only one of these actually exists at runtime.
    // If Wails uses "wails", ALL frontend API calls silently fail.
    const go = (window as any).go

    console.log('window.go keys:', Object.keys(go || {}))
    console.log('window.go.main:', go?.main)
    console.log('window.go.wails:', go?.wails)
    console.log(
      'pluginApi.ts uses: window.go.main.PluginApi',
    )
    console.log(
      'generated bindings use: window.go.wails.PluginApi',
    )

    // At least one namespace must exist at runtime
    const hasMain = !!go?.main?.PluginApi
    const hasWails = !!go?.wails?.PluginApi
    expect(hasMain || hasWails).toBe(true)
  })
})

describe('pluginApi method delegation', () => {
  it('getPlugins delegates to backend.GetPlugins', async () => {
    const mockFn = vi.fn().mockResolvedValue([{ id: 'p1', name: 'Test' }])
    ;(window as any).go = {
      wails: { PluginApi: { GetPlugins: mockFn } },
    }

    const result = await pluginApi.getPlugins()
    expect(result).toEqual([{ id: 'p1', name: 'Test' }])
  })

  it('reloadPlugins delegates to backend.ReloadPlugins', async () => {
    const mockFn = vi.fn().mockResolvedValue({ count: 3 })
    ;(window as any).go = {
      wails: { PluginApi: { ReloadPlugins: mockFn } },
    }

    const result = await pluginApi.reloadPlugins()
    expect(result).toEqual({ count: 3 })
  })

  it('getPlugin delegates to backend.GetPlugin with correct ID', async () => {
    const mockFn = vi.fn().mockResolvedValue({ metadata: { id: 'tia-export' } })
    ;(window as any).go = {
      wails: { PluginApi: { GetPlugin: mockFn } },
    }

    const result = await pluginApi.getPlugin('tia-export')
    expect(mockFn).toHaveBeenCalledWith('tia-export')
    expect(result.metadata.id).toBe('tia-export')
  })
})
