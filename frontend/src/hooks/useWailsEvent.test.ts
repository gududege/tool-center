import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useWailsEvent } from './useWailsEvent'

function createMockRuntime() {
  const handlers: Record<string, Array<(...args: any[]) => void>> = {}
  return {
    EventsOn: vi.fn((event: string, handler: (...args: any[]) => void) => {
      if (!handlers[event]) handlers[event] = []
      handlers[event].push(handler)
    }),
    EventsOff: vi.fn((event: string, handler?: (...args: any[]) => void) => {
      if (handler && handlers[event]) {
        handlers[event] = handlers[event].filter((h) => h !== handler)
      } else if (!handler) {
        delete handlers[event]
      }
    }),
    emit(event: string, data: any) {
      if (handlers[event]) {
        handlers[event].forEach((h) => h(data))
      }
    },
  }
}

describe('useWailsEvent', () => {
  let runtime: ReturnType<typeof createMockRuntime>

  beforeEach(() => {
    runtime = createMockRuntime()
    ;(window as any).runtime = runtime
    delete (window as any).wails
  })

  it('calls handler when event is emitted', () => {
    const handler = vi.fn()

    renderHook(() => useWailsEvent('test:event', handler))

    runtime.emit('test:event', { value: 42 })

    expect(handler).toHaveBeenCalledTimes(1)
    expect(handler).toHaveBeenCalledWith({ value: 42 })
  })

  it('unsubscribes only the specific handler on unmount', () => {
    const handlerA = vi.fn()
    const handlerB = vi.fn()

    const { unmount: unmountA } = renderHook(() =>
      useWailsEvent('test:event', handlerA),
    )
    renderHook(() => useWailsEvent('test:event', handlerB))

    runtime.emit('test:event', {})
    expect(handlerA).toHaveBeenCalledTimes(1)
    expect(handlerB).toHaveBeenCalledTimes(1)

    unmountA()

    runtime.emit('test:event', {})
    expect(handlerA).toHaveBeenCalledTimes(1)
    expect(handlerB).toHaveBeenCalledTimes(2)

    expect(runtime.EventsOff).toHaveBeenCalledWith(
      'test:event',
      expect.any(Function),
    )
  })

  it('does not crash when runtime is not available', () => {
    ;(window as any).runtime = undefined
    const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})

    const handler = vi.fn()
    expect(() =>
      renderHook(() => useWailsEvent('test:event', handler)),
    ).not.toThrow()

    expect(warnSpy).toHaveBeenCalled()
    warnSpy.mockRestore()
  })
})
