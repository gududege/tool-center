// ============================================================
// CLI Tool Center — TypeScript Types
// 匹配 Go 后端 DTO (internal/adapter/wails/types.go)
// ============================================================

// --- Plugin ---

export interface PluginSummaryDto {
  id: string
  name: string
  name_cn?: string
  description?: string
  description_cn?: string
  version?: string
  icon?: string
  navigation: NavigationDto
}

export interface PluginDto {
  metadata: PluginMetadataDto
  navigation: NavigationDto
  form: FormDefinitionDto
  execution: ExecutionDefinitionDto
}

export interface PluginMetadataDto {
  id: string
  name: string
  name_cn?: string
  description?: string
  description_cn?: string
  version?: string
  author?: string
  icon?: string
}

export interface NavigationDto {
  group: string[]
  group_cn?: string[]
  order: number
}

export interface FormDefinitionDto {
  schema: any
  uiSchema?: any
}

export interface ExecutionDefinitionDto {
  exe: string
  workingDirectory?: string
  parameters?: ParameterMappingDto[]
}

export interface ParameterMappingDto {
  field: string
  kind: string
  flag?: string
  style?: string
  separator?: string
  trueFlag?: string
  falseFlag?: string
  defaultValue?: any
}

// --- Task ---

export type TaskStatus =
  | 'created'
  | 'queued'
  | 'running'
  | 'completed'
  | 'failed'
  | 'cancelled'

export interface TaskDto {
  id: string
  pluginId: string
  status: TaskStatus
  createdAt: string
  startedAt?: string
  endedAt?: string
}

// --- Output ---

export type OutputSource = 'stdout' | 'stderr' | 'system'
export type LogLevel = 'trace' | 'debug' | 'info' | 'warn' | 'error'

export interface OutputEventDto {
  taskId: string
  timestamp: string
  level: string
  source: OutputSource
  message: string
}

// --- Settings ---

export interface SettingsDto {
  theme: string
  formTheme?: string
  pluginDirectory: string
  language?: string
  sidebarCollapsed: boolean
  sidebarSize: number
  bottomPanelSize: number
  bottomTab: string
}

// --- System ---

export interface SystemInfoDto {
  appVersion: string
  buildTime: string
  goVersion: string
  os: string
}

// --- API Request/Response ---

export interface RunPluginRequest {
  pluginId: string
  formData: Record<string, unknown>
}

export interface RunPluginResponse {
  taskId: string
}

export interface CancelTaskResponse {
  success: boolean
}

export interface DeleteTaskResponse {
  success: boolean
}

export interface ReloadPluginsResponse {
  count: number
}

export interface ValidationResultDto {
  valid: boolean
  errors?: ValidationErrorDto[]
}

export interface ValidationErrorDto {
  path: string
  message: string
}

export interface PluginHistoryEntryDto {
  timestamp: string
  label: string
  formData: Record<string, unknown>
}

export interface ApiErrorDto {
  code: string
  message: string
  details?: string
}

export interface FileFilterDto {
  displayName: string
  patterns: string[]
}

export interface DialogResponse {
  path?: string
}

// --- Event Payloads ---

export interface PluginsReloadedEvent {
  count: number
}

export interface TaskCreatedEvent {
  taskId: string
  pluginId: string
}

export interface TaskStartedEvent {
  taskId: string
}

export interface TaskCompletedEvent {
  taskId: string
}

export interface TaskFailedEvent {
  taskId: string
  error: string
}

export interface TaskCancelledEvent {
  taskId: string
}

export interface TaskStatusChangedEvent {
  taskId: string
  oldStatus: TaskStatus
  newStatus: TaskStatus
}

// --- UI types ---

export interface MenuNode {
  id: string
  name: string
  children: MenuNode[]
  pluginId?: string
}

export interface Tab {
  id: string
  pluginId: string
  title: string
  dirty: boolean
  taskId?: string
}

export interface OutputLine {
  timestamp: string
  source: OutputSource
  level: string
  message: string
}

export type PanelPosition = 'bottom' | 'right' | 'hidden'
