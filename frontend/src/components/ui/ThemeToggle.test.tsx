import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { ThemeToggle } from './ThemeToggle'
import * as settingsStoreModule from '../../stores/settingsStore'
import * as TranslationContextModule from '../../contexts/TranslationContext'
import type { SettingsDto } from '../../types'

vi.mock('../../stores/settingsStore', () => ({
  useSettingsStore: vi.fn(),
}))

vi.mock('../../contexts/TranslationContext', () => ({
  useT: vi.fn(),
}))

function mockT(key: string) {
  const map: Record<string, string> = {
    'theme.light': 'Light',
    'theme.dark': 'Dark',
    'theme.system': 'System',
  }
  return map[key] ?? key
}

function createMockSettings(overrides: Partial<SettingsDto> = {}): SettingsDto {
  return {
    theme: 'dark',
    pluginDirectory: './plugins',
    maxOutputLines: 10000,
    maxTaskHistory: 100,
    autoReloadPlugins: false,
    language: 'en',
    windowWidth: 1200,
    windowHeight: 800,
    sidebarCollapsed: false,
    sidebarWidth: 280,
    panelPosition: 'bottom',
    panelHeight: 200,
    panelWidth: 400,
    ...overrides,
  }
}

describe('ThemeToggle', () => {
  const saveSettings = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(TranslationContextModule.useT as any).mockReturnValue({ t: mockT })
    ;(settingsStoreModule.useSettingsStore as any).mockReturnValue({
      settings: null,
      saveSettings,
    })
  })

  it('renders current theme label', () => {
    const settings = createMockSettings({ theme: 'dark' })
    ;(settingsStoreModule.useSettingsStore as any).mockReturnValue({
      settings,
      saveSettings,
    })

    render(<ThemeToggle />)
    expect(screen.getByText('Dark')).toBeDefined()
  })

  it('opens dropdown on button click', () => {
    const settings = createMockSettings({ theme: 'dark' })
    ;(settingsStoreModule.useSettingsStore as any).mockReturnValue({
      settings,
      saveSettings,
    })

    render(<ThemeToggle />)

    const button = screen.getByTitle('Dark')
    fireEvent.click(button)

    expect(screen.getByText('Light')).toBeDefined()
    expect(screen.getByText('System')).toBeDefined()
  })

  it('selects a new theme and saves', () => {
    const settings = createMockSettings({ theme: 'dark', language: 'en' })
    ;(settingsStoreModule.useSettingsStore as any).mockReturnValue({
      settings,
      saveSettings,
    })

    render(<ThemeToggle />)

    fireEvent.click(screen.getByTitle('Dark'))
    fireEvent.click(screen.getByText('Light'))

    expect(saveSettings).toHaveBeenCalledWith({
      ...settings,
      theme: 'light',
    })
  })

  it('does not crash when settings is null', () => {
    ;(settingsStoreModule.useSettingsStore as any).mockReturnValue({
      settings: null,
      saveSettings,
    })

    render(<ThemeToggle />)

    const button = screen.getByTitle('Dark')
    fireEvent.click(button)

    fireEvent.click(screen.getByText('Light'))

    expect(saveSettings).not.toHaveBeenCalled()
  })

  it('closes dropdown on outside click', () => {
    const settings = createMockSettings({ theme: 'dark' })
    ;(settingsStoreModule.useSettingsStore as any).mockReturnValue({
      settings,
      saveSettings,
    })

    render(<ThemeToggle />)

    fireEvent.click(screen.getByTitle('Dark'))
    expect(screen.queryByText('Light')).toBeTruthy()

    fireEvent.click(document.body)
    expect(screen.queryByText('Light')).toBeNull()
  })
})
