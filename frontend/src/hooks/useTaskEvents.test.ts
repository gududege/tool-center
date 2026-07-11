import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook } from '@testing-library/react'
import { useTaskEvents } from './useTaskEvents'
import { useTaskStore } from '../stores/taskStore'
import { useOutputStore } from '../stores/outputStore'
import { usePluginStore } from '../stores/pluginStore'
import { useTabStore } from '../stores/tabStore'
import type {
  TaskCreatedEvent,
  TaskStartedEvent,
  TaskCompletedEvent,
  TaskFailedEvent,
  TaskCancelledEvent,
  TaskStatusChangedEvent,
  OutputEventDto,
  PluginsReloadedEvent,
} from '../types'

vi.mock('../stores/taskStore', () => ({
  useTaskStore: vi.fn(),
}))

vi.mock('../stores/outputStore', () => ({
  useOutputStore: vi.fn(),
}))

vi.mock('../stores/pluginStore', () => ({
  usePluginStore: vi.fn(),
}))

vi.mock('../stores/tabStore', () => ({
  useTabStore: vi.fn(() => ({ tabs: [], activeTabId: null })),
}))

function createMockRuntime() {
  const handlers: Record<string, Array<(...args: any[]) => void>> = {}
  return {
    EventsOn: vi.fn((event: string, handler: (...args: any[]) => void) => {
      if (!handlers[event]) handlers[event] = []
      handlers[event].push(handler)
    }),
    EventsOff: vi.fn(),
    emit(event: string, data: any) {
      if (handlers[event]) {
        handlers[event].forEach((h) => h(data))
      }
    },
  }
}

function setupStores() {
  const addTask = vi.fn()
  const updateTask = vi.fn()
  const removeTask = vi.fn()
  const appendOutput = vi.fn()
  const fetchPlugins = vi.fn()

  ;(useTaskStore as any).mockReturnValue({ addTask, updateTask, removeTask })
  ;(useOutputStore as any).mockReturnValue({ appendOutput })
  ;(usePluginStore as any).mockReturnValue({ fetchPlugins })

  return { addTask, updateTask, removeTask, appendOutput, fetchPlugins }
}

describe('useTaskEvents', () => {
  let runtime: ReturnType<typeof createMockRuntime>
  let stores: ReturnType<typeof setupStores>

  beforeEach(() => {
    vi.clearAllMocks()
    runtime = createMockRuntime()
    ;(window as any).runtime = runtime
    ;(window as any).wails = undefined
    stores = setupStores()
  })

  it('dispatches task:created to addTask', () => {
    renderHook(() => useTaskEvents())

    const event: TaskCreatedEvent = { taskId: 'task-1', pluginId: 'p1' }
    runtime.emit('task:created', event)

    expect(stores.addTask).toHaveBeenCalledTimes(1)
    const call = stores.addTask.mock.calls[0][0]
    expect(call.id).toBe('task-1')
    expect(call.pluginId).toBe('p1')
    expect(call.status).toBe('created')
  })

  it('dispatches task:failed to updateTask with failure status', () => {
    renderHook(() => useTaskEvents())

    const event: TaskFailedEvent = { taskId: 'task-1', error: 'something went wrong' }
    runtime.emit('task:failed', event)

    expect(stores.updateTask).toHaveBeenCalledWith('task-1', {
      status: 'failed',
      endedAt: expect.any(String),
    })
  })

  it('dispatches output:append to appendOutput', () => {
    renderHook(() => useTaskEvents())

    const event: OutputEventDto = {
      taskId: 'task-1',
      timestamp: '2024-01-01T00:00:00Z',
      level: 'info',
      source: 'stdout',
      message: 'hello world',
    }
    runtime.emit('output:append', event)

    expect(stores.appendOutput).toHaveBeenCalledWith('task-1', event)
  })

  it('dispatches plugin:reloaded to fetchPlugins', () => {
    renderHook(() => useTaskEvents())

    const event: PluginsReloadedEvent = { count: 3 }
    runtime.emit('plugin:reloaded', event)

    expect(stores.fetchPlugins).toHaveBeenCalledTimes(1)
  })
})
