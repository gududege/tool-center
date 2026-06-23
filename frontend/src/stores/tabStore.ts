import { create } from 'zustand'
import type { Tab } from '../types'

interface TabState {
  tabs: Tab[]
  activeTabId: string | null

  addTab: (tab: Tab) => string
  closeTab: (tabId: string) => void
  activateTab: (tabId: string) => void
  setTabDirty: (tabId: string, dirty: boolean) => void
  setTabTaskId: (tabId: string, taskId: string) => void
  updateTabTitle: (tabId: string, title: string) => void
}

let tabCounter = 0

export const useTabStore = create<TabState>((set, get) => ({
  tabs: [],
  activeTabId: null,

  addTab: (tab: Tab) => {
    const id = tab.id || `tab-${++tabCounter}`
    const newTab = { ...tab, id }
    set((state) => {
      // 如果同 plugin 已有未命名的 tab，复用它
      const existing = state.tabs.find(
        (t) => t.pluginId === tab.pluginId && t.title === tab.title && !t.taskId,
      )
      if (existing) {
        return { activeTabId: existing.id }
      }
      return {
        tabs: [...state.tabs, newTab],
        activeTabId: id,
      }
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
