import { create } from 'zustand'
import type { TaskDto } from '../types'
import { taskApi } from '../services/taskApi'

interface TaskState {
  tasks: TaskDto[]
  loading: boolean

  fetchTasks: () => Promise<void>
  addTask: (task: TaskDto) => void
  updateTask: (taskId: string, updates: Partial<TaskDto>) => void
  removeTask: (taskId: string) => void
  refreshTask: (taskId: string) => Promise<void>
}

export const useTaskStore = create<TaskState>((set, get) => ({
  tasks: [],
  loading: false,

  fetchTasks: async () => {
    set({ loading: true })
    try {
      const tasks = await taskApi.getTasks()
      set({ tasks, loading: false })
    } catch {
      set({ loading: false })
    }
  },

  addTask: (task: TaskDto) => {
    set((state) => ({
      tasks: [task, ...state.tasks],
    }))
  },

  updateTask: (taskId: string, updates: Partial<TaskDto>) => {
    set((state) => ({
      tasks: state.tasks.map((t) =>
        t.id === taskId ? { ...t, ...updates } : t,
      ),
    }))
  },

  removeTask: (taskId: string) => {
    set((state) => ({
      tasks: state.tasks.filter((t) => t.id !== taskId),
    }))
  },

  refreshTask: async (taskId: string) => {
    try {
      const task = await taskApi.getTask(taskId)
      // Task may be deleted (404)
      if (!task) {
        get().removeTask(taskId)
        return
      }
      get().updateTask(taskId, task)
    } catch {
      get().removeTask(taskId)
    }
  },
}))
