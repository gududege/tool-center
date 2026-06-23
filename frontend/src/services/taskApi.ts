/**
 * Task API 服务层
 *
 * 封装 Wails 后端 TaskApi 的所有调用。
 */

import type {
  TaskDto,
  OutputEventDto,
  RunPluginRequest,
  RunPluginResponse,
  CancelTaskResponse,
  DeleteTaskResponse,
} from '../types'

const backend = () => (window as any).go.wails.TaskApi

export const taskApi = {
  async runPlugin(request: RunPluginRequest): Promise<RunPluginResponse> {
    return backend().RunPlugin(request)
  },

  async runPluginRaw(
    pluginId: string,
    formData: Record<string, unknown>,
  ): Promise<string> {
    return backend().RunPlugin({ pluginId, formData }).then((r: RunPluginResponse) => r.taskId)
  },

  async getTask(taskId: string): Promise<TaskDto> {
    return backend().GetTask(taskId)
  },

  async getTasks(): Promise<TaskDto[]> {
    return backend().GetTasks()
  },

  async cancelTask(taskId: string): Promise<CancelTaskResponse> {
    return backend().CancelTask(taskId)
  },

  async deleteTask(taskId: string): Promise<DeleteTaskResponse> {
    return backend().DeleteTask(taskId)
  },

  async getTaskOutput(taskId: string): Promise<OutputEventDto[]> {
    return backend().GetTaskOutput(taskId)
  },
}
