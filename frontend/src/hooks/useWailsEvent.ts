import { useEffect, useCallback } from 'react'

type WailsEventHandler = (data: any) => void

/**
 * 监听 Wails 事件的 Hook
 *
 * Wails v2 通过 runtime.EventsOn / EventsOff 实现事件
 * 在 WebView 中通过 window.runtime.EventsOn 暴露
 */
export function useWailsEvent(
  eventName: string,
  handler: WailsEventHandler,
  deps: any[] = [],
) {
  const stableHandler = useCallback(handler, deps)

  useEffect(() => {
    const runtime = (window as any).runtime ?? (window as any).wails?.runtime
    if (!runtime?.EventsOn) {
      console.warn('[WailsEvent] runtime.EventsOn not available')
      return
    }

    runtime.EventsOn(eventName, stableHandler)

    return () => {
      if (runtime?.EventsOff) {
        runtime.EventsOff(eventName)
      }
    }
  }, [eventName, stableHandler])
}

/**
 * 触发 Wails 事件（通常前端不需要，留给后端触发）
 */
export function wailsEmit(eventName: string, data?: any) {
  const runtime = (window as any).runtime ?? (window as any).wails?.runtime
  if (runtime?.EventsEmit) {
    runtime.EventsEmit(eventName, data)
  }
}
