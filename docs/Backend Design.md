# Backend Design

Project: CLI Tool Center

Document: Backend Design

Version: 1.0

Status: Draft

Last Updated: 2026-06-21

---

# 1. Purpose

本文档定义 Backend 实现架构。

包括：

* 包结构
* 服务分层
* Repository
* Process Runner
* Plugin Loader
* Parameter Builder
* Event Bus
* Wails Adapter

目标：

```text
高内聚
低耦合
可测试
可扩展
```

---

# 2. Backend Architecture

## Layered Architecture

```text
┌─────────────────────────────┐
│ Wails Adapter Layer         │
├─────────────────────────────┤
│ Application Service Layer   │
├─────────────────────────────┤
│ Domain Layer                │
├─────────────────────────────┤
│ Infrastructure Layer        │
└─────────────────────────────┘
```

---

# 3. Dependency Rules

允许：

```text
Adapter
 ↓

Application
 ↓

Domain
 ↓

Infrastructure
```

禁止：

```text
Domain
 ↓
Wails

Domain
 ↓
React

PluginManager
 ↓
TaskManager
```

通过接口解耦。

---

# 4. Project Structure

```text
/

├── main.go
├── wails.json
├── go.mod
├── internal/
│   ├── domain/
│   │   ├── plugin/
│   │   ├── task/
│   │   ├── process/
│   │   ├── output/
│   │   ├── session/
│   │   ├── settings/
│   │   └── events/
│   ├── application/
│   │   ├── plugin/
│   │   ├── task/
│   │   ├── settings/
│   │   └── dialog/
│   ├── infrastructure/
│   │   ├── filesystem/
│   │   ├── process/
│   │   ├── events/
│   │   └── repositories/
│   ├── adapter/
│   │   └── wails/
│   └── bootstrap/
├── plugins/
├── frontend/
└── docs/
```

---

# 5. Domain Layer

禁止依赖：

```text
Wails

os/exec

文件系统

数据库
```

---

## Structure

```text
domain/

├── plugin/
├── task/
├── process/
├── output/
├── session/
├── settings/
└── events/
```

---

# 6. Plugin Domain

```text
plugin/

├── model.go
├── navigation.go
├── parameter.go
└── repository.go
```

---

## Repository Interface

```go
type PluginRepository interface {
    Get(id string) (*Plugin, error)

    List() ([]*Plugin, error)

    Save(plugin *Plugin) error

    Delete(id string) error
}
```

---

# 7. Task Domain

```text
task/

├── model.go
├── status.go
└── repository.go
```

---

## Repository Interface

```go
type TaskRepository interface {
    Get(id string) (*Task, error)

    List() ([]*Task, error)

    Save(task *Task) error

    Delete(id string) error
}
```

---

# 8. Application Layer

负责：

```text
业务流程

事务协调

领域对象协作
```

---

## Structure

```text
application/

├── plugin/
├── task/
├── process/
├── output/
├── dialog/
└── settings/
```

---

# 9. Plugin Service

```go
type PluginService struct {
    plugins PluginRepository

    loader PluginLoader

    validator PluginValidator
}
```

---

## Responsibilities

```text
加载插件

重载插件

校验插件

菜单生成
```

---

## API

```go
func (s *PluginService) Load()

func (s *PluginService) Reload()

func (s *PluginService) Get()

func (s *PluginService) List()
```

---

# 10. Task Service

核心服务。

---

```go
type TaskService struct {
    tasks TaskRepository

    processes ProcessRunner

    plugins PluginRepository

    events EventBus
}
```

---

## Responsibilities

```text
创建任务

执行任务

取消任务

状态更新
```

---

## API

```go
RunPlugin()

GetTask()

ListTasks()

CancelTask()

DeleteTask()
```

---

# 11. RunPlugin Flow

```text
Frontend

 ↓

RunPlugin()

 ↓

Load Plugin

 ↓

Build Arguments

 ↓

Create Task

 ↓

Start Process

 ↓

Update Task

 ↓

Return TaskID
```

---

# 12. Process Service

```go
type ProcessService struct {
    runner ProcessRunner
}
```

---

## Responsibilities

```text
启动进程

终止进程

监控进程

收集输出
```

---

# 13. Infrastructure Layer

负责：

```text
文件系统

CLI进程

事件实现

Schema验证
```

---

## Structure

```text
infrastructure/

├── filesystem/
├── process/
├── events/
├── schema/
└── repositories/
```

---

# 14. Plugin Loader

```go
type PluginLoader interface {
    Load(dir string) ([]*Plugin, error)
}
```

---

## Default Implementation

```text
filesystem/pluginloader
```

---

### Flow

```text
Scan

 ↓

Read plugin.json

 ↓

Validate

 ↓

Resolve Paths

 ↓

Build Plugin
```

---

# 15. Plugin Validator

```go
type PluginValidator interface {
    Validate(
        plugin *Plugin,
    ) error
}
```

---

## Validation Rules

```text
ID唯一

Exe存在

Schema存在

Navigation合法
```

---

# 16. Parameter Builder

核心组件。

---

## Purpose

负责：

```text
FormData

 ↓

CLI Arguments
```

转换。

---

## Interface

```go
type ParameterBuilder interface {
    Build(
        mappings []ParameterMapping,
        formData map[string]any,
    ) ([]string, error)
}
```

---

## Example

Input:

```json
{
  "project": "demo.ap20",
  "overwrite": true
}
```

Output:

```text
--project
demo.ap20
--overwrite
```

---

# 17. Process Runner

核心执行器。

---

## Interface

```go
type ProcessRunner interface {
    Start(
        request StartProcessRequest,
    ) (*ProcessInstance, error)

    Stop(
        processID string,
    ) error
}
```

---

# 18. StartProcessRequest

```go
type StartProcessRequest struct {
    Executable string

    Arguments []string

    WorkingDirectory string

    Environment []string
}
```

---

# 19. Process Lifecycle

```text
Created

 ↓

Starting

 ↓

Running

 ↓

Exited
```

异常：

```text
Running

 ↓

Killed
```

---

# 20. Default Runner

实现：

```go
exec.CommandContext()
```

---

## Cancellation

```go
ctx,cancel :=
    context.WithCancel(...)
```

保存：

```go
task.CancelFunc
```

---

# 21. Windows Process Tree

推荐实现：

```text
CREATE_NEW_PROCESS_GROUP
```

---

取消时：

```text
Kill Process Tree
```

避免子进程泄漏。

---

# 22. Output Collector

负责：

```text
stdout

stderr
```

采集。

---

## Interface

```go
type OutputCollector interface {
    Collect(
        process *ProcessInstance,
    )
}
```

---

## Output Flow

```text
stdout

 ↓

Scanner

 ↓

OutputEvent

 ↓

EventBus

 ↓

Frontend
```

---

# 23. Event Bus

系统核心。

---

## Interface

```go
type EventBus interface {
    Publish(
        event DomainEvent,
    )

    Subscribe(
        topic string,
    )
}
```

---

# 24. Event Topics

```text
plugin.loaded

plugin.reloaded

task.created

task.started

task.completed

task.failed

task.cancelled

output.received
```

---

# 25. Wails Event Adapter

桥接：

```text
Domain Event

 ↓

Wails Event

 ↓

Frontend
```

---

## Example

```go
events.Publish(
    TaskStarted{},
)
```

转换：

```go
runtime.EventsEmit(
    ctx,
    "task:started",
    payload,
)
```

---

# 26. Output Buffer

每个任务维护：

```go
type OutputBuffer struct {
    TaskID string

    Events []OutputEvent
}
```

---

## Retention

配置：

```go
MaxOutputLines
```

超出：

```text
FIFO Eviction
```

---

# 27. Repositories

V1全部采用内存实现。

---

## PluginRepository

```go
MemoryPluginRepository
```

---

## TaskRepository

```go
MemoryTaskRepository
```

---

## SessionRepository

```go
MemorySessionRepository
```

---

# 28. Future Persistence

V2支持：

```text
SQLite
```

实现：

```go
SqliteTaskRepository

SqlitePluginRepository
```

无需修改业务层。

---

# 29. Settings Service

```go
type SettingsService struct {
    repository SettingsRepository
}
```

---

## Storage

```text
settings.json
```

---

## Example

```json
{
  "theme": "dark",

  "pluginDirectory": "./plugins"
}
```

---

# 30. Bootstrap Layer

负责依赖注入。

---

## Structure

```text
bootstrap/

├── repositories.go
├── services.go
├── events.go
└── app.go
```

---

# 31. Startup Sequence

```text
Application Start

 ↓

Load Settings

 ↓

Create EventBus

 ↓

Create Repositories

 ↓

Create Services

 ↓

Load Plugins

 ↓

Ready
```

---

# 32. Wails Adapter Layer

对外暴露API。

---

## Structure

```text
adapter/wails/

├── plugin_api.go
├── task_api.go
├── dialog_api.go
├── settings_api.go
└── system_api.go
```

---

# 33. Plugin API

```go
func (a *PluginApi)
    GetPlugins()

func (a *PluginApi)
    GetPlugin(id string)

func (a *PluginApi)
    ReloadPlugins()
```

---

# 34. Task API

```go
func (a *TaskApi)
    RunPlugin()

func (a *TaskApi)
    GetTask()

func (a *TaskApi)
    GetTasks()

func (a *TaskApi)
    CancelTask()
```

---

# 35. Error Handling

统一错误模型。

---

## Domain Error

```go
var (
    ErrPluginNotFound

    ErrTaskNotFound

    ErrSchemaInvalid

    ErrProcessFailed
)
```

---

## API Error

统一转换：

```go
ApiErrorDto
```

---

# 36. Logging

实际使用：

```text
slog
```

---

## Categories

```text
plugin

task

process

eventbus

system
```

---

# 37. Testing Strategy

## Unit Test

覆盖：

```text
ParameterBuilder

PluginLoader

TaskService

ProcessRunner
```

---

## Integration Test

覆盖：

```text
Plugin → Task

Task → Process

Process → Output
```

---

# 38. Concurrency Model

Task支持并发执行。

---

保护对象：

```go
PluginRegistry

TaskRepository

ProcessRegistry

OutputBuffer
```

---

推荐：

```go
sync.RWMutex
```

---

# 39. Future Extensions

V2：

```text
SQLite

Plugin Marketplace

Auto Update

Task History
```

---

V3：

```text
Workflow Engine

Remote Agent

Distributed Execution

Cluster Scheduler
```

---

# 40. Summary

核心服务：

```text
PluginService

TaskService

ProcessService

SettingsService
```

核心基础设施：

```text
PluginLoader

ParameterBuilder

ProcessRunner

OutputCollector

EventBus
```

核心原则：

```text
Domain First

Service Oriented

Plugin Driven

Task Driven

Event Driven
```

Backend 的唯一职责是：

```text
Plugin
 ↓

Task
 ↓

Process
 ↓

Output
```

将插件描述转换为可执行任务，并通过事件实时反馈执行结果。
