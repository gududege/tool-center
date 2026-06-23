# Architecture Design Document

**Project:** CLI Tool Center

**Version:** 1.0

**Status:** Draft

**Last Updated:** 2026-06-21

---

# 1. Introduction

## 1.1 Purpose

CLI Tool Center 是一个基于 Wails 的插件化桌面应用平台，用于统一管理、配置和执行各种命令行工具（CLI）。

系统通过 JSON Schema 驱动动态表单生成，通过插件机制实现工具扩展，通过任务系统管理 CLI 生命周期，并提供统一的日志查看和执行体验。

本架构文档用于定义系统整体架构、模块职责、数据流、系统边界和扩展策略。

---

## 1.2 Goals

系统目标：

### G1 插件化

支持新增工具而无需修改主程序代码。

### G2 Schema驱动

支持使用 JSON Schema 自动生成参数配置界面。

### G3 多任务

支持多个 CLI 同时执行。

### G4 解耦

实现：

* UI 与 CLI 解耦
* Plugin 与 Runtime 解耦
* Task 与 Tab 解耦

### G5 可扩展

为未来功能预留扩展能力：

* 插件市场
* 自动更新
* 任务队列
* 工作流
* 远程执行

---

## 1.3 Non Goals

V1 不包含：

* 云同步
* 插件签名验证
* 远程 Agent
* Workflow Engine
* Web 版本
* Linux/macOS 支持优化

---

# 2. Technology Stack

## Backend

| Technology   | Purpose           |
| ------------ | ----------------- |
| Go           | Core Runtime      |
| Wails        | Desktop Framework |
| os/exec      | CLI Execution     |
| context      | Task Control      |
| slog         | Logging           |

---

## Frontend

| Technology            | Purpose           |
| --------------------- | ----------------- |
| React                 | UI Framework      |
| TypeScript            | Type Safety       |
| Vite                  | Build Tool        |
| TailwindCSS           | Styling           |
| shadcn/ui             | UI Components     |
| react-jsonschema-form | Dynamic Form      |
| @rjsf/validator-ajv8  | Schema Validation |
| Ajv2020               | JSON Schema 2020-12 |
| Zustand               | Client State      |

---

# 3. High Level Architecture

## 3.1 Logical View

```text
┌────────────────────────────────────────────┐
│ React Frontend                             │
│                                            │
│ Navigation Tree                            │
│ Dynamic Forms                              │
│ Tabs                                       │
│ Task Center                                │
│ Output Viewer                              │
└───────────────────┬────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│ Wails Bridge                               │
└───────────────────┬────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│ Backend Services                           │
│                                            │
│ PluginManager                              │
│ SchemaManager                              │
│ TaskManager                                │
│ ProcessRunner                              │
│ EventBus                                   │
│ SettingsManager                            │
└───────────────────┬────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│ External CLI Processes                     │
└────────────────────────────────────────────┘
```

---

## 3.2 Layer Responsibilities

### Presentation Layer

负责：

* 用户交互
* 表单渲染
* 任务展示
* 日志展示

不负责：

* 参数构造
* 进程管理

---

### Application Layer

负责：

* 插件管理
* 任务管理
* 生命周期管理

不负责：

* UI

---

### Infrastructure Layer

负责：

* CLI启动
* 文件访问
* 事件传递

---

# 4. Core Architectural Principles

## 4.1 Plugin First

所有业务能力均通过插件实现。

主程序仅提供：

* Runtime
* UI Framework
* Task Engine

插件负责：

* 表单定义
* 参数映射
* CLI执行

---

## 4.2 Schema First

界面由 Schema 驱动生成。

禁止：

```text
手写业务配置界面
```

允许：

```text
JSON Schema
 ↓
Dynamic Form
```

---

## 4.3 Task First

所有执行均抽象为 Task。

统一管理：

* 状态
* 日志
* 生命周期

---

## 4.4 Event Driven

模块间通信采用事件机制。

禁止：

```text
Frontend
 ↓
直接操作Process
```

采用：

```text
Frontend
 ↓
TaskManager
 ↓
ProcessRunner
```

---

# 5. Plugin Architecture

## 5.1 Plugin Directory

```text
plugins/

├─ export/
│  ├─ plugin.json
│  ├─ schema.json
│  ├─ uischema.json
│  ├─ export.exe
│  └─ icon.svg

├─ aml-import/
│  ├─ plugin.json
│  └─ aml.exe

└─ docgen/
   ├─ plugin.json
   └─ docgen.exe
```

---

## 5.2 Plugin Lifecycle

```text
Startup
 ↓
Scan
 ↓
Load
 ↓
Validate
 ↓
Register
 ↓
Ready
```

---

## 5.3 Plugin Runtime

插件运行时包含：

```text
Metadata
Navigation
Form Definition
Execution Definition
```

---

# 6. Frontend Architecture

## 6.1 UI Layout

```text
┌─────────────────────────────────────────┐
│ Header                                  │
├──────────────┬──────────────────────────┤
│ Navigation   │ Tabs                     │
│              │                          │
│              │ Dynamic Form             │
│              │                          │
├──────────────┴──────────────────────────┤
│ Output Viewer                           │
└─────────────────────────────────────────┘
```

---

## 6.2 Navigation Tree

根据插件自动生成：

```text
Project
 ├ Export
 ├ Import

Tools
 ├ Compare
 ├ Analyze
```

---

## 6.3 Dynamic Form

采用：

```text
react-jsonschema-form
```

生成界面。

流程：

```text
Schema
 +
UISchema
 ↓
RJSF
 ↓
React Components
```

---

## 6.4 Widget Registry

支持 RJSF 内置控件和本地自定义控件。

RJSF 内置：

```text
checkbox / radio / select
textarea / hidden / range / updown
color / password
```

本地自定义：

```text
folderPicker
filePicker
saveFilePicker
```

未来扩展：

```text
plc-selector
block-selector
tag-selector
```

---

## 6.5 Tab Model

重要原则：

```text
Tab ≠ Task
```

Tab 仅代表：

```text
Form Session
```

支持：

```text
Export #1
Export #2
Export #3
```

多个实例。

---

# 7. Backend Architecture

## 7.1 PluginManager

职责：

* 扫描插件
* 加载插件
* 校验配置
* 构建菜单

接口：

```go
LoadPlugins()

GetPlugin()

GetPlugins()
```

---

## 7.2 SchemaManager

职责：

* 加载 Schema
* 缓存 Schema
* 校验 Schema

接口：

```go
LoadSchema()

ValidateSchema()
```

---

## 7.3 TaskManager

职责：

* 创建任务
* 管理状态
* 管理生命周期

状态机：

```text
Created
 ↓
Queued
 ↓
Running
 ↓
Completed

Running
 ↓
Failed

Running
 ↓
Cancelled
```

---

## 7.4 ProcessRunner

职责：

* 参数构造
* 启动CLI
* 输出采集
* 终止CLI

接口：

```go
Run()

Stop()

BuildArgs()
```

---

## 7.5 EventBus

职责：

* Output Event
* Task Event
* Plugin Event

统一事件分发。

---

# 8. Task Architecture

## 8.1 Why Task Exists

CLI 生命周期通常长于 UI 生命周期。

例如：

```text
Export Project
  ↓
关闭Tab
  ↓
继续运行
```

因此：

```text
Task
```

必须独立于：

```text
Tab
```

---

## 8.2 Task Model

```go
type Task struct {
    ID string

    PluginID string

    Status TaskStatus

    CreatedAt time.Time

    StartedAt time.Time

    EndedAt *time.Time

    Cancel context.CancelFunc
}
```

---

## 8.3 Task Center

统一查看：

```text
Running
Completed
Failed
Cancelled
```

任务。

---

# 9. Process Architecture

## 9.1 Execution Flow

```text
Form Data
 ↓
Parameter Mapping
 ↓
Build Args
 ↓
Create Process
 ↓
Capture Output
 ↓
Update Task
```

---

## 9.2 Cancellation

采用：

```go
context.WithCancel()
```

管理生命周期。

禁止：

```go
Process.Kill()
```

作为默认实现。

---

## 9.3 Process Tree

Windows 下预留：

```text
CREATE_NEW_PROCESS_GROUP
```

支持终止整个进程树。

---

# 10. Output Architecture

## 10.1 Unified Event Model

```go
type OutputEvent struct {
    TaskID string

    Timestamp time.Time

    Source string

    Level string

    Message string
}
```

---

## 10.2 Sources

支持：

```text
stdout
stderr
system
```

统一处理。

---

## 10.3 Benefits

未来支持：

* 搜索
* 导出
* 过滤
* JSON日志

无需修改核心架构。

---

# 11. Data Flow

## Plugin Load

```text
Application Start
 ↓
PluginManager
 ↓
Scan Plugins
 ↓
Validate
 ↓
Register
 ↓
Build Navigation Tree
 ↓
Frontend Render
```

---

## Execute CLI

```text
User
 ↓
Submit Form
 ↓
RunPlugin
 ↓
TaskManager
 ↓
Create Task
 ↓
ProcessRunner
 ↓
Start CLI
 ↓
Output Events
 ↓
Frontend
```

---

## Cancel Task

```text
User
 ↓
Stop
 ↓
TaskManager
 ↓
Cancel Context
 ↓
Process Exit
 ↓
Task Updated
```

---

# 12. Scalability Strategy

## V1

支持：

* Plugin Discovery
* Dynamic Forms
* Task System
* Output Viewer

---

## V2

支持：

* Plugin Marketplace
* Plugin Update
* Task History
* Favorites

---

## V3

支持：

* Workflow
* Batch Execution
* Remote Agents
* Distributed Runtime

---

# 13. Architecture Decisions

## ADR-001

使用 Wails。

原因：

* Go原生
* 打包简单
* CLI集成方便

---

## ADR-002

使用 React JSON Schema Form。

原因：

* 适合参数配置场景
* 学习成本低
* 社区成熟

---

## ADR-003

使用 Plugin Runtime。

原因：

* 新增工具无需重新编译

---

## ADR-004

Task 与 Tab 解耦。

原因：

* 支持后台运行
* 支持任务中心

---

## ADR-005

采用 Event Driven Architecture。

原因：

* 模块低耦合
* 易于扩展

---

# 14. Conclusion

CLI Tool Center 采用 Plugin + Schema + Task 的核心架构。

核心设计理念：

```text
Plugin First
Schema First
Task First
Event Driven
```

该架构能够满足当前 TIA Portal 工具中心需求，并为未来插件生态、自动化工作流和远程执行能力提供稳定的扩展基础。
