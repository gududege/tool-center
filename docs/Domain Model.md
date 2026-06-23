# Domain Model

Project: CLI Tool Center

Version: 1.0

Status: Draft

Last Updated: 2026-06-21

---

# 1. Overview

本系统采用：

```text
Plugin Driven
Task Driven
Event Driven
```

架构。

核心目标：

```text
Plugin
 ↓

Form Session
 ↓

Task
 ↓

Process
 ↓

Output
```

整个系统围绕两个核心聚合构建：

```text
Plugin Aggregate
Task Aggregate
```

---

# 2. Domain Context

## Plugin Context

负责：

```text
插件发现
插件加载
插件校验
菜单生成
Schema管理
```

---

## Runtime Context

负责：

```text
表单会话
任务管理
进程管理
日志管理
```

---

## System Context

负责：

```text
配置
主题
插件目录
系统状态
```

---

# 3. Aggregate Overview

```text
Plugin Aggregate

Plugin
 ├ Metadata
 ├ Navigation
 ├ FormDefinition
 └ ExecutionDefinition
```

---

```text
Task Aggregate

Task
 ├ ProcessInstance
 ├ OutputBuffer
 └ OutputEvents
```

---

# 4. Plugin Aggregate

Plugin 是整个系统的扩展单元。

---

## Plugin

```go
type Plugin struct {
    Metadata PluginMetadata

    Navigation Navigation

    Form FormDefinition

    Execution ExecutionDefinition
}
```

---

## Responsibilities

```text
描述工具能力

描述菜单位置

描述表单

描述参数映射

描述执行行为
```

---

## Invariants

```text
ID必须唯一

Execution.Exe必须存在

Schema必须合法

Navigation不能为空
```

---

# 5. PluginMetadata

```go
type PluginMetadata struct {
    ID string

    Name string

    Description string

    Version string

    Author string

    Homepage string

    Icon string
}
```

---

## Example

```go
PluginMetadata{
    ID: "tia-export",

    Name: "TIA Export",

    Version: "1.0.0",
}
```

---

# 6. Navigation

菜单定义。

```go
type Navigation struct {
    Group []string

    Order int
}
```

---

## Example

```json
{
  "group": [
    "TIA",
    "Export"
  ]
}
```

生成：

```text
TIA
 └ Export
```

---

# 7. FormDefinition

```go
type FormDefinition struct {
    SchemaPath string

    UISchemaPath string

    Schema map[string]any

    UISchema map[string]any
}
```

---

## Responsibilities

```text
加载Schema

缓存Schema

提供前端表单定义
```

---

# 8. ExecutionDefinition

```go
type ExecutionDefinition struct {
    Exe string

    WorkingDirectory string

    Environment map[string]string

    Parameters []ParameterMapping
}
```

---

# 9. ParameterMapping

用于：

```text
FormData
 ↓
CLI Arguments
```

转换。

---

```go
type ParameterMapping struct {
    Field string

    Kind ParameterKind

    Flag string

    Style string

    Separator string

    TrueFlag string

    FalseFlag string

    DefaultValue any
}
```

---

# 10. ParameterKind

```go
type ParameterKind string
```

---

```go
const (
    ArgumentKind ParameterKind = "argument"

    ArgumentArrayKind ParameterKind = "argument-array"

    OptionKind ParameterKind = "option"

    OptionArrayKind ParameterKind = "option-array"

    SwitchKind ParameterKind = "switch"

    BoolOptionKind ParameterKind = "bool-option"

    DualSwitchKind ParameterKind = "dual-switch"
)
```

---

# 11. Menu Tree

运行时菜单结构。

---

## MenuNode

```go
type MenuNode struct {
    ID string

    Name string

    Order int

    Parent *MenuNode

    Children []*MenuNode

    PluginID *string
}
```

---

## Example

```text
TIA
 ├ Export
 ├ Import
 └ Compare
```

---

# 12. Plugin Registry

运行时插件仓库。

---

```go
type PluginRegistry struct {
    Plugins map[string]*Plugin
}
```

---

## Responsibilities

```text
查询插件

注册插件

卸载插件

重载插件
```

---

# 13. Form Session

Form Session ≠ Task

Form Session = 一个打开的Tab

---

## FormSession

```go
type FormSession struct {
    ID string

    PluginID string

    FormData map[string]any

    CreatedAt time.Time

    UpdatedAt time.Time
}
```

---

## Example

```text
Export #1

Export #2

Export #3
```

三个独立Session。

---

# 14. Session Registry

```go
type SessionRegistry struct {
    Sessions map[string]*FormSession
}
```

---

# 15. Task Aggregate

Task 是运行时核心对象。

---

## 生命周期

```text
Created
 ↓

Queued
 ↓

Running
 ↓

Completed
```

异常：

```text
Running
 ↓
Failed

Running
 ↓
Cancelled
```

---

# 16. Task

```go
type Task struct {
    ID string

    PluginID string

    SessionID string

    Status TaskStatus

    CreatedAt time.Time

    StartedAt *time.Time

    EndedAt *time.Time

    ExitCode *int

    ErrorMessage string

    ProcessID *string
}
```

---

## Responsibilities

```text
描述执行状态

关联Process

关联输出

关联插件
```

---

# 17. TaskStatus

```go
type TaskStatus string
```

---

```go
const (
    TaskCreated TaskStatus = "created"

    TaskQueued TaskStatus = "queued"

    TaskRunning TaskStatus = "running"

    TaskCompleted TaskStatus = "completed"

    TaskFailed TaskStatus = "failed"

    TaskCancelled TaskStatus = "cancelled"
)
```

---

# 18. Task State Machine

```text
Created
 ↓

Queued
 ↓

Running
 ├ Completed
 ├ Failed
 └ Cancelled
```

---

## Illegal Transitions

禁止：

```text
Completed → Running

Failed → Running

Cancelled → Running
```

---

# 19. Task Repository

```go
type TaskRepository struct {
    Tasks map[string]*Task
}
```

---

# 20. Process Instance

实际CLI进程。

---

```go
type ProcessInstance struct {
    ID string

    TaskID string

    PID int

    Command string

    Arguments []string

    WorkingDirectory string

    Environment []string

    StartedAt time.Time
}
```

---

## Example

```text
export.exe
    --project demo.ap20
    --overwrite
```

---

# 21. Process Registry

```go
type ProcessRegistry struct {
    Processes map[string]*ProcessInstance
}
```

---

# 22. Output Event

统一日志模型。

---

```go
type OutputEvent struct {
    ID string

    TaskID string

    Timestamp time.Time

    Source OutputSource

    Level LogLevel

    Message string
}
```

---

## Example

```json
{
  "taskId": "task1",

  "source": "stdout",

  "message": "Export started"
}
```

---

# 23. OutputSource

```go
type OutputSource string
```

---

```go
const (
    StdoutSource OutputSource = "stdout"

    StderrSource OutputSource = "stderr"

    SystemSource OutputSource = "system"
)
```

---

# 24. LogLevel

```go
type LogLevel string
```

---

```go
const (
    TraceLevel LogLevel = "trace"

    DebugLevel LogLevel = "debug"

    InfoLevel LogLevel = "info"

    WarnLevel LogLevel = "warn"

    ErrorLevel LogLevel = "error"
)
```

---

# 25. Output Buffer

任务日志缓存。

---

```go
type OutputBuffer struct {
    TaskID string

    Events []OutputEvent
}
```

---

## Rules

```text
Append Only

不允许修改历史记录
```

---

# 26. Settings Aggregate

系统配置。

---

```go
type Settings struct {
    Theme string

    PluginDirectory string

    MaxOutputLines int

    MaxTaskHistory int

    AutoReloadPlugins bool
}
```

---

# 27. Runtime Context

运行时状态。

---

```go
type RuntimeState struct {
    StartedAt time.Time

    LoadedPlugins int

    RunningTasks int

    ActiveSessions int
}
```

---

# 28. Domain Events

所有状态变化统一事件化。

---

## DomainEvent

```go
type DomainEvent interface {
    Topic() string
}
```

---

# 29. Plugin Events

## PluginLoaded

```go
type PluginLoaded struct {
    PluginID string
}
```

---

## PluginUnloaded

```go
type PluginUnloaded struct {
    PluginID string
}
```

---

## PluginReloaded

```go
type PluginReloaded struct {
    PluginID string
}
```

---

# 30. Task Events

## TaskCreated

```go
type TaskCreated struct {
    TaskID string
}
```

---

## TaskStarted

```go
type TaskStarted struct {
    TaskID string
}
```

---

## TaskCompleted

```go
type TaskCompleted struct {
    TaskID string

    ExitCode int
}
```

---

## TaskFailed

```go
type TaskFailed struct {
    TaskID string

    Error string
}
```

---

## TaskCancelled

```go
type TaskCancelled struct {
    TaskID string
}
```

---

# 31. Output Events

## OutputReceived

```go
type OutputReceived struct {
    TaskID string

    Event OutputEvent
}
```

---

# 32. Event Bus

```go
type EventBus interface {
    Publish(event DomainEvent)

    Subscribe(
        topic string,
    )
}
```

---

# 33. Aggregate Boundaries

## Plugin Aggregate

```text
Plugin
 ├ Metadata
 ├ Navigation
 ├ FormDefinition
 └ ExecutionDefinition
```

修改Plugin不会影响Task。

---

## Task Aggregate

```text
Task
 ├ ProcessInstance
 ├ OutputBuffer
 └ OutputEvents
```

Task拥有完整执行生命周期。

---

# 34. Domain Invariants

## Plugin

```text
Plugin.ID 唯一

Plugin.Name 非空

Exe存在

Schema合法
```

---

## Task

```text
Running状态最多一个Process

Completed不可恢复

Cancelled不可恢复
```

---

## Output

```text
OutputEvent不可修改

仅允许Append
```

---

# 35. Future Extensions

V2：

```text
PluginDependency

PluginMarketplace

TaskQueue

TaskScheduler
```

---

V3：

```text
Workflow

RemoteAgent

DistributedExecution

ClusterTask
```

---

# 36. Summary

核心聚合：

```text
Plugin
Task
```

核心实体：

```text
FormSession
ProcessInstance
OutputEvent
```

核心值对象：

```text
Navigation
PluginMetadata
ParameterMapping
```

核心仓库：

```text
PluginRegistry
SessionRegistry
TaskRepository
ProcessRegistry
```

核心事件：

```text
PluginLoaded

TaskCreated

TaskStarted

TaskCompleted

OutputReceived
```

整个 Runtime 的所有服务均围绕这些领域对象构建。
