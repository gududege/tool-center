import '@testing-library/jest-dom'

// Mock Wails runtime
const mockWailsApi = {
  GetPlugins: vi.fn(),
  GetPlugin: vi.fn(),
  ReloadPlugins: vi.fn(),
  ValidateFormData: vi.fn(),
}

// Set up both possible Wails namespaces for testing
Object.defineProperty(window, 'go', {
  value: {
    main: {
      PluginApi: { ...mockWailsApi },
      TaskApi: {},
      SettingsApi: {},
      SystemApi: {},
    },
    wails: {
      PluginApi: { ...mockWailsApi },
      TaskApi: {},
      SettingsApi: {},
      SystemApi: {},
    },
  },
  writable: true,
  configurable: true,
})
