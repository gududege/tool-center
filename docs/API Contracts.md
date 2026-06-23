# API Contracts

**Project:** CLI Tool Center

**Document:** API Contracts

**Version:** 1.0

**Status:** Draft

**Last Updated:** 2026-06-21

---

# 1. Introduction

## 1.1 Purpose

本文档定义：

```text
React Frontend
      ↕
Wails Bridge
      ↕
Go Backend
```

之间的所有接口契约。

包括：

* Request Models
* Response Models
* Event Models
* Error Models
* DTO Definitions

---

# 2. Design Principles

## API-001

Frontend 不访问内部对象。

禁止：

```typescript
Task
Plugin
Process
```

直接暴露。

---

## API-002

所有接口返回 DTO。

```typescript
PluginDto
TaskDto
OutputEventDto
```

---

## API-003

所有事件均为只追加（Append Only）。

禁止：

```typescript
UpdateOutputLine()
```

允许：

```typescript
AppendOutputEvent()
```

---

## API-004

所有接口异步。

```typescript
Promise<T>
```

---

# 3. Service Overview

## Backend Services

```text
PluginService

TaskService

DialogService

SettingsService

SystemService
```

---

# 4. Plugin Service

## GetPlugins

获取插件列表。

### React

```typescript
await GetPlugins()
```

---

### Request

无

---

### Response

```typescript
type PluginSummaryDto = {
    id: string

    name: string

    name_cn?: string

    description?: string

    description_cn?: string

    version?: string

    icon?: string

    navigation: NavigationDto
}
```

---

### Example

```json
[
  {
    "id": "tia-export",
    "name": "TIA Export",

    "navigation": {
      "group": [
        "TIA",
        "Export"
      ],
      "order": 100
    }
  }
]
```

---

## GetPlugin

获取完整插件定义。

---

### Request

```typescript
type GetPluginRequest = {
    pluginId: string
}
```

---

### Response

```typescript
type PluginDto = {
    metadata: PluginMetadataDto

    navigation: NavigationDto

    form: FormDefinitionDto

    execution: ExecutionDefinitionDto
}
```

---

## ReloadPlugins

重新扫描插件目录。

---

### Response

```typescript
type ReloadPluginsResponse = {
    count: number
}
```

---

## ValidateFormData

校验表单数据是否符合插件 schema。

### Request

```typescript
type ValidateFormDataRequest = {
    pluginId: string

    data: Record<string, unknown>
}
```

### Response

```typescript
type ValidationResultDto = {
    valid: boolean

    errors: ValidationErrorDto[]
}
```

---

## SavePluginHistory

保存当前表单参数到插件历史（最多 5 条滚动保存）。

### Request

```typescript
type SavePluginHistoryRequest = {
    pluginId: string

    formData: Record<string, unknown>
}
```

### Response

无返回值，出错时抛错。

---

## GetPluginHistory

获取插件已保存的历史参数列表。

### Response

```typescript
type PluginHistoryEntryDto = {
    timestamp: string

    label: string

    formData: Record<string, unknown>
}
```

---

# 5. Form Service

## GetFormDefinition

获取 Schema 和 UISchema。

---

### Request

```typescript
type GetFormDefinitionRequest = {
    pluginId: string
}
```

---

### Response

```typescript
type FormDefinitionDto = {
    schema: any

    uiSchema?: any
}
```

---

## ValidateFormData

本地调试接口。

通常由 AJV 完成。

---

### Request

```typescript
type ValidateFormDataRequest = {
    pluginId: string

    data: Record<string, unknown>
}
```

---

### Response

```typescript
type ValidationResultDto = {
    valid: boolean

    errors: ValidationErrorDto[]
}
```

---

# 6. Task Service

## RunPlugin

执行插件。

---

### Request

```typescript
type RunPluginRequest = {
    pluginId: string

    formData: Record<string, unknown>
}
```

---

### Response

```typescript
type RunPluginResponse = {
    taskId: string
}
```

---

### Example

```json
{
  "pluginId": "tia-export",

  "formData": {
    "project": "demo.ap20",

    "overwrite": true
  }
}
```

---

## GetTask

获取单个任务。

---

### Request

```typescript
type GetTaskRequest = {
    taskId: string
}
```

---

### Response

```typescript
type TaskDto = {
    id: string

    pluginId: string

    status: TaskStatus

    createdAt: string

    startedAt?: string

    endedAt?: string
}
```

---

## GetTasks

获取任务列表。

---

### Response

```typescript
type TaskDto[]
```

---

## CancelTask

取消任务。

---

### Request

```typescript
type CancelTaskRequest = {
    taskId: string
}
```

---

### Response

```typescript
type CancelTaskResponse = {
    success: boolean
}
```

---

## DeleteTask

删除历史任务。

---

### Request

```typescript
type DeleteTaskRequest = {
    taskId: string
}
```

---

### Response

```typescript
type DeleteTaskResponse = {
    success: boolean
}
```

---

# 7. Output Service

## GetTaskOutput

获取历史输出。

---

### Request

```typescript
type GetTaskOutputRequest = {
    taskId: string
}
```

---

### Response

```typescript
type OutputEventDto[] 
```

---

## ExportOutput

导出日志。

---

### Request

```typescript
type ExportOutputRequest = {
    taskId: string

    path: string
}
```

---

### Response

```typescript
type ExportOutputResponse = {
    success: boolean
}
```

---

# 8. Dialog Service

## OpenFileDialog

打开文件选择器。

---

### Request

```typescript
type OpenFileDialogRequest = {
    title?: string

    filters?: FileFilterDto[]
}
```

---

### Response

```typescript
type OpenFileDialogResponse = {
    path?: string
}
```

---

## OpenFolderDialog

---

### Response

```typescript
type OpenFolderDialogResponse = {
    path?: string
}
```

---

## SaveFileDialog

---

### Response

```typescript
type SaveFileDialogResponse = {
    path?: string
}
```

---

# 9. Settings Service

## GetSettings

---

### Response

```typescript
type SettingsDto = {
    theme: string

    formTheme?: string

    pluginDirectory: string

    language?: string

    sidebarCollapsed: boolean

    sidebarSize: number

    bottomPanelSize: number

    bottomTab: string
}
```

---

## SaveSettings

---

### Request

```typescript
type SaveSettingsRequest = {
    settings: SettingsDto
}
```

---

### Response

```typescript
type SaveSettingsResponse = {
    success: boolean
}
```

---

# 10. System Service

## GetSystemInfo

---

### Response

```typescript
type SystemInfoDto = {
    appVersion: string

    buildTime: string

    goVersion: string

    os: string
}
```

---

# 11. DTO Definitions

## PluginMetadataDto

```typescript
type PluginMetadataDto = {
    id: string

    name: string

    name_cn?: string

    description?: string

    description_cn?: string

    version?: string

    author?: string

    icon?: string
}
```

---

## NavigationDto

```typescript
type NavigationDto = {
    group: string[]

    group_cn?: string[]

    order: number
}
```

---

## ExecutionDefinitionDto

```typescript
type ExecutionDefinitionDto = {
    exe: string

    workingDirectory?: string

    parameters?: ParameterMappingDto[]
}

type ParameterMappingDto = {
    field: string

    kind: string

    flag?: string

    style?: string

    separator?: string

    trueFlag?: string

    falseFlag?: string

    defaultValue?: any
}
```

---

## TaskStatus

```typescript
type TaskStatus =
    | "created"
    | "queued"
    | "running"
    | "completed"
    | "failed"
    | "cancelled"
```

---

## TaskDto

```typescript
type TaskDto = {
    id: string

    pluginId: string

    status: TaskStatus

    createdAt: string

    startedAt?: string

    endedAt?: string
}
```

---

## OutputEventDto

```typescript
type OutputEventDto = {
    taskId: string

    timestamp: string

    level: string

    source: OutputSource

    message: string
}
```

---

## OutputSource

```typescript
type OutputSource =
    | "stdout"
    | "stderr"
    | "system"
```

---

## ValidationErrorDto

```typescript
type ValidationErrorDto = {
    path: string

    message: string
}
```

---

# 12. Event Contracts

Wails Events 用于实时通知。

---

# 13. Plugin Events

## plugins:reloaded

### Payload

```typescript
type PluginsReloadedEvent = {
    count: number
}
```

---

# 14. Task Events

## task:created

### Payload

```typescript
type TaskCreatedEvent = {
    taskId: string

    pluginId: string
}
```

---

## task:started

### Payload

```typescript
type TaskStartedEvent = {
    taskId: string
}
```

---

## task:completed

### Payload

```typescript
type TaskCompletedEvent = {
    taskId: string
}
```

---

## task:failed

### Payload

```typescript
type TaskFailedEvent = {
    taskId: string

    error: string
}
```

---

## task:cancelled

### Payload

```typescript
type TaskCancelledEvent = {
    taskId: string
}
```

---

## task:status-changed

统一状态变化事件。

### Payload

```typescript
type TaskStatusChangedEvent = {
    taskId: string

    oldStatus: TaskStatus

    newStatus: TaskStatus
}
```

---

# 15. Output Events

## output:append

实时日志事件。

### Payload

```typescript
type OutputEventDto
```

---

### Example

```json
{
  "taskId": "task-001",

  "timestamp": "2026-06-21T12:00:00Z",

  "level": "info",

  "source": "stdout",

  "message": "Export started..."
}
```

---

# 16. Error Model

## ErrorResponse

所有接口统一返回：

```typescript
type ApiErrorDto = {
    code: string

    message: string

    details?: string
}
```

---

## Standard Error Codes

```text
PLUGIN_NOT_FOUND

PLUGIN_LOAD_FAILED

PLUGIN_VALIDATION_FAILED

TASK_NOT_FOUND

TASK_ALREADY_FINISHED

PROCESS_START_FAILED

PROCESS_CANCEL_FAILED

SCHEMA_NOT_FOUND

INVALID_FORM_DATA

INTERNAL_ERROR
```

---

# 17. Backend Interface Definition

## PluginService

```go
type PluginService interface {
    GetPlugins() ([]PluginSummaryDto, error)

    GetPlugin(id string) (*PluginDto, error)

    ReloadPlugins() (*ReloadPluginsResponse, error)

    ValidateFormData(pluginID string, data map[string]any) (*ValidationResultDto, error)

    SavePluginHistory(pluginID string, data map[string]any) error

    GetPluginHistory(pluginID string) ([]PluginHistoryEntryDto, error)
}
```

---

## TaskService

```go
type TaskService interface {
    RunPlugin(
        pluginID string,
        formData map[string]any,
    ) (*RunPluginResponse, error)

    GetTask(
        taskID string,
    ) (*TaskDto, error)

    GetTasks() ([]TaskDto, error)

    CancelTask(taskID string) (*CancelTaskResponse, error)

    DeleteTask(taskID string) (*DeleteTaskResponse, error)

    GetTaskOutput(taskID string) ([]OutputEventDto, error)
}
```

---

## OutputService

```go
type OutputService interface {
    GetTaskOutput(
        taskID string,
    ) ([]OutputEventDto, error)
}
```

---

# 18. Frontend Generated SDK

推荐自动生成：

```typescript
export const api = {
    getPlugins(),
    getPlugin(),
    runPlugin(),
    getTask(),
    getTasks(),
    cancelTask()
}
```

React 仅调用 SDK。

禁止直接调用：

```typescript
window.go....
```

---

# 19. Recommended Event Flow

## Execute Plugin

```text
RunPlugin
    ↓
task:created
    ↓
task:started
    ↓
output:append
    ↓
output:append
    ↓
task:completed
```

---

## Failed Plugin

```text
RunPlugin
    ↓
task:created
    ↓
task:started
    ↓
output:append
    ↓
task:failed
```

---

## Cancel Plugin

```text
RunPlugin
    ↓
task:created
    ↓
task:started
    ↓
CancelTask
    ↓
task:cancelled
```

---

# 20. Version Compatibility

API Version：

```text
v1
```

规则：

* 新增字段允许
* 删除字段禁止
* 修改字段类型禁止
* Event Name 修改禁止

保证插件和前端长期兼容。
