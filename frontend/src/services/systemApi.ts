/**
 * System API 服务层
 */

import type { SystemInfoDto } from '../types'

const backend = () => (window as any).go.wails.SystemApi

export const systemApi = {
  async getSystemInfo(): Promise<SystemInfoDto> {
    return backend().GetSystemInfo()
  },
}
