/// <reference types="vite/client" />

interface Window {
  go?: {
    main?: {
      PluginApi?: any
      TaskApi?: any
      SettingsApi?: any
      SystemApi?: any
    }
  }
  runtime?: {
    EventsOn?: (eventName: string, callback: (data?: any) => void) => void
    EventsOff?: (eventName: string) => void
    EventsEmit?: (eventName: string, data?: any) => void
  }
  wails?: {
    runtime?: {
      EventsOn?: (eventName: string, callback: (data?: any) => void) => void
      EventsOff?: (eventName: string) => void
      EventsEmit?: (eventName: string, data?: any) => void
    }
  }
}
