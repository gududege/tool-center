/**
 * Dialog API 服务层
 *
 * 调用 Go 后端 SystemApi 的原生文件/目录选择对话框（经由 Wails runtime）。
 * Go 侧用 wails/v2/pkg/runtime 的 OpenDirectoryDialog/OpenFileDialog/SaveFileDialog
 * 唤起操作系统原生对话框。在 Wails 环境外（开发浏览器模式）优雅降级为 prompt。
 */

interface FileFilter {
  displayName: string
  patterns: string[]
}

interface DialogOptions {
  title?: string
  defaultDirectory?: string
}

interface OpenFileOptions extends DialogOptions {
  filters?: FileFilter[]
}

// Go 绑定：window.go.wails.SystemApi（由 wails generate 自动生成）。
// 其返回值是 { path?: string }（DialogResponse），取消选择时 path 为空字符串。
function getSystemApi(): any {
  return (window as any).go?.wails?.SystemApi
}

function isWails(): boolean {
  return typeof getSystemApi()?.OpenFolderDialog === 'function'
}

/** 打开目录选择对话框 */
export async function openDirectoryDialog(options: DialogOptions = {}): Promise<string | null> {
  if (isWails()) {
    const res = await getSystemApi().OpenFolderDialog(options.title ?? 'Select Directory')
    return res?.path || null
  }
  // Fallback for dev/browser mode
  const path = prompt(options.title ?? 'Enter directory path:', options.defaultDirectory ?? '')
  return path || null
}

/** 打开文件选择对话框 */
export async function openFileDialog(options: OpenFileOptions = {}): Promise<string | null> {
  if (isWails()) {
    const filters = options.filters?.map((f) => ({
      displayName: f.displayName,
      patterns: f.patterns,
    })) ?? []
    const res = await getSystemApi().OpenFileDialog(options.title ?? 'Select File', filters)
    return res?.path || null
  }
  // Fallback for dev/browser mode
  const path = prompt(options.title ?? 'Enter file path:', options.defaultDirectory ?? '')
  return path || null
}

/** 打开保存文件对话框 */
export async function saveFileDialog(options: OpenFileOptions = {}): Promise<string | null> {
  if (isWails()) {
    const filters = options.filters?.map((f) => ({
      displayName: f.displayName,
      patterns: f.patterns,
    })) ?? []
    const res = await getSystemApi().SaveFileDialog(options.title ?? 'Save File', filters)
    return res?.path || null
  }
  // Fallback for dev/browser mode
  const path = prompt(options.title ?? 'Enter file path:', options.defaultDirectory ?? '')
  return path || null
}
