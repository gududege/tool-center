import { create } from 'zustand'
import type { Tab } from '../types'

interface TabState {
  _counter: number
  tabs: Tab[]
  activeTabId: string | null

  addTab: (tab: Tab) => string
  closeTab: (tabId: string) => void
  activateTab: (tabId: string) => void
  setTabDirty: (tabId: string, dirty: boolean) => void
  setTabTaskId: (tabId: string, taskId: string) => void
  updateTabTitle: (tabId: string, title: string) => void
}

export const useTabStore = create<TabState>((set, get) => ({
  _counter: 0,
  tabs: [],
  activeTabId: null,

  addTab: (tab: Tab) => {
    const state = get()
    const id = tab.id || `tab-${state._counter + 1}`
    const newTab = { ...tab, id }

    const existing = state.tabs.find(
      (t) => t.pluginId === tab.pluginId,
    )
    if (existing) {
      set({ activeTabId: existing.id })
      return existing.id
    }

    set({
      _counter: state._counter + (tab.id ? 0 : 1),
      tabs: [...state.tabs, newTab],
      activeTabId: id,
    })
    return id
  },

  closeTab: (tabId: string) => {
    set((state) => {
      const tabs = state.tabs.filter((t) => t.id !== tabId)
      let activeTabId = state.activeTabId
      if (activeTabId === tabId) {
        const idx = state.tabs.findIndex((t) => t.id === tabId)
        activeTabId = tabs[Math.min(idx, tabs.length - 1)]?.id ?? null
      }
      return { tabs, activeTabId }
    })
  },

  activateTab: (tabId: string) => {
    set({ activeTabId: tabId })
  },

  setTabDirty: (tabId: string, dirty: boolean) => {
    set((state) => ({
      tabs: state.tabs.map((t) => (t.id === tabId ? { ...t, dirty } : t)),
    }))
  },

  setTabTaskId: (tabId: string, taskId: string) => {
    set((state) => ({
      tabs: state.tabs.map((t) => (t.id === tabId ? { ...t, taskId } : t)),
    }))
  },

  updateTabTitle: (tabId: string, title: string) => {
    set((state) => ({
      tabs: state.tabs.map((t) => (t.id === tabId ? { ...t, title } : t)),
    }))
  },
}))
