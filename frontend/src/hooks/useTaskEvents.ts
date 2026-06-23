import { useEffect } from 'react'
import { useWailsEvent } from './useWailsEvent'
import { useTaskStore } from '../stores/taskStore'
import { useOutputStore } from '../stores/outputStore'
import { useTabStore } from '../stores/tabStore'
import { usePluginStore } from '../stores/pluginStore'
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

/**
 * 全局 Wails 事件监听器
 *
 * 在 App 根组件挂载，监听所有后端事件并分发到对应 Store。
 */
export function useTaskEvents() {
  const { addTask, updateTask, removeTask } = useTaskStore()
  const { appendOutput } = useOutputStore()
  const { fetchPlugins } = usePluginStore()

  // task:created
  useWailsEvent('task:created', (event: TaskCreatedEvent) => {
    addTask({
      id: event.taskId,
      pluginId: event.pluginId,
      status: 'created',
      createdAt: new Date().toISOString(),
    })
  })

  // task:started
  useWailsEvent('task:started', (event: TaskStartedEvent) => {
    updateTask(event.taskId, {
      status: 'running',
      startedAt: new Date().toISOString(),
    })
  })

  // task:completed
  useWailsEvent('task:completed', (event: TaskCompletedEvent) => {
    updateTask(event.taskId, {
      status: 'completed',
      endedAt: new Date().toISOString(),
    })
  })

  // task:failed
  useWailsEvent('task:failed', (event: TaskFailedEvent) => {
    updateTask(event.taskId, {
      status: 'failed',
      endedAt: new Date().toISOString(),
    })
  })

  // task:cancelled
  useWailsEvent('task:cancelled', (event: TaskCancelledEvent) => {
    updateTask(event.taskId, {
      status: 'cancelled',
      endedAt: new Date().toISOString(),
    })
  })

  // task:status-changed
  useWailsEvent('task:status-changed', (event: TaskStatusChangedEvent) => {
    updateTask(event.taskId, {
      status: event.newStatus as any,
    })
  })

  // output:append
  useWailsEvent('output:append', (event: OutputEventDto) => {
    appendOutput(event.taskId, event)
  })

  // plugin:reloaded
  useWailsEvent('plugin:reloaded', (_event: PluginsReloadedEvent) => {
    fetchPlugins()
  })
}
