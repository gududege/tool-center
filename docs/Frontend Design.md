# Frontend Design

**Project:** CLI Tool Center

**Document:** Frontend Design

**Version:** 1.0

**Status:** Draft

**Last Updated:** 2026-06-21

---

# 1. Purpose

本文档定义前端架构设计。

包括：

* 页面结构
* 状态管理
* 组件设计
* Tab系统
* Plugin菜单
* Dynamic Form
* Task管理
* Output Viewer
* Wails事件集成

目标：

```text
Plugin First

Schema First

Task Driven

Event Driven
```

---

# 2. Technology Stack

## Core

```text
React 19

TypeScript

Vite

Wails v3
```

---

## UI

```text
TailwindCSS

shadcn/ui

lucide-react
```

---

## Form

```text
react-jsonschema-form

@rjsf/validator-ajv8 + Ajv2020
```

---

## State

```text
Zustand
```

---

## Terminal Output

日志输出面板自行实现滚动。

---

# 3. Frontend Architecture

```text
App
 │
 ├── Sidebar
 │
 ├── TabManager
 │
 ├── TaskPanel
 │
 └── OutputPanel
```

---

采用：

```text
Feature Based Architecture
```

---

目录：

```text
src/

├── app
├── layouts
├── features
├── components
├── services
├── stores
├── hooks
└── types
```

---

# 4. Layout Design

整体布局：

```text
┌────────────────────────────────────┐
│ Toolbar                            │
├──────────┬─────────────────────────┤
│ Sidebar  │                         │
│          │        Tabs             │
│          │                         │
│          │                         │
├──────────┴─────────────────────────┤
│ Task Panel                         │
├────────────────────────────────────┤
│ Output Panel                       │
└────────────────────────────────────┘
```

---

# 5. MainLayout

```tsx
<MainLayout>
    <Sidebar />

    <Workspace />

    <TaskPanel />

    <OutputPanel />
</MainLayout>
```

---

职责：

```text
布局

面板管理

窗口状态
```

---

# 6. Sidebar

显示插件菜单。

---

数据来源：

```text
GetPlugins()
```

---

结构：

```text
TIA
 ├ Export
 ├ Import
 └ Compare

Tools
 ├ Diff
 └ Convert
```

---

组件：

```tsx
<PluginTree />
```

---

# 7. Plugin Tree Model

```ts
interface MenuNode {
    id: string

    name: string

    children: MenuNode[]

    pluginId?: string
}
```

---

# 8. Plugin Click Flow

```text
Click Plugin

 ↓

GetPlugin()

 ↓

Open Tab

 ↓

Load Schema

 ↓

Render Form
```

---

# 9. Tab System

核心模块。

---

一个插件：

```text
可打开多个Tab
```

例如：

```text
Export #1

Export #2

Export #3
```

---

互不影响。

---

# 10. Tab Model

```ts
interface Tab {
    id: string

    pluginId: string

    title: string

    dirty: boolean
}
```

---

# 11. Tab Store

```ts
interface TabStore {
    tabs: Tab[]

    activeTabId?: string

    addTab()

    closeTab()

    activateTab()
}
```

---

# 12. Workspace

显示当前Tab。

---

```tsx
<Workspace>
    <DynamicFormTab />
</Workspace>
```

---

# 13. Dynamic Form

核心组件。

---

职责：

```text
读取Schema

生成表单

收集数据

提交任务
```

---

组件：

```tsx
<DynamicForm />
```

---

# 14. Dynamic Form Flow

```text
Plugin

 ↓

schema.json

 ↓

RJSF

 ↓

Generated Form

 ↓

FormData
```

---

# 15. Dynamic Form Component

```tsx
<DynamicForm
    schema={schema}
    uiSchema={uiSchema}
    formData={formData}
/>
```

---

# 16. Supported Fields

JSON Schema 基础类型：

```text
string
number
integer
boolean
array
object
```

以及由 `enum` / `anyOf` / `oneOf` 描述的枚举、多选和单选。

---

# 17. Widgets

表单基于 **react-jsonschema-form (RJSF)** + **shadcn/ui** 主题渲染，并在 Shadow DOM 中隔离样式。

Validator 使用 **Ajv2020**，与插件 schema 的 `$schema: https://json-schema.org/draft/2020-12/schema` 保持一致。

## 17.1 RJSF / shadcn 内置 Widgets

| Widget | 适用 Schema | 说明 |
| --- | --- | --- |
| `checkbox` | boolean | 复选框，boolean 默认值 |
| `radio` | boolean / enum string | 单选按钮组 |
| `select` | boolean / enum string / enum array | 下拉选择 |
| `textarea` | string | 多行文本 |
| `hidden` | string / number / integer / boolean | 隐藏字段 |
| `range` | number / integer | 范围滑块 |
| `updown` | number / integer | 数字步进器 |
| `color` | string | 颜色选择器 |
| `password` | string | 密码输入框 |

## 17.2 Format 自动映射

以下 `format` 会自动映射到对应 input type：

| Format | 控件 |
| --- | --- |
| `email` | email |
| `uri` / `url` | url |
| `date` | date |
| `date-time` | datetime-local |
| `time` | time |
| `color` | color |
| `password` | password |
| `data-url` | file |

## 17.3 本地自定义 Widgets

| Widget | 说明 |
| --- | --- |
| `folderPicker` | 目录选择对话框，返回路径字符串 |
| `filePicker` | 文件选择对话框，返回路径字符串 |
| `saveFilePicker` | 保存文件对话框，返回路径字符串 |

在 `uischema.json` 中通过 `Control.options` 启用：

```json
{
  "type": "Control",
  "scope": "#/properties/indir",
  "options": { "folderPicker": true }
}
```

通用 widget 覆盖：

```json
{
  "type": "Control",
  "scope": "#/properties/mode",
  "options": { "widget": "radio" }
}
```

`inputType` 可用于改变 `<input>` 类型：

```json
{
  "type": "Control",
  "scope": "#/properties/secret",
  "options": { "inputType": "password" }
}
```

---

# 20. Form Actions

底部工具栏：

```text
Run          # 校验表单并启动 CLI
Reset        # 清空表单数据
Preview      # 切换命令行实时预览
Export       # 导出当前表单数据为 JSON
Import       # 从 JSON 文件加载表单数据
History      # 加载最近保存的历史参数
```

---

# 21. Run Flow

```text
Click Run
  ↓
RJSF/Ajv Validate
  ↓
Save current params to plugin history.json (max 5)
  ↓
taskApi.runPlugin()
  ↓
Task Created
  ↓
Process Started
  ↓
Output Events
```

历史保存仅依赖校验通过，不依赖运行结果。

---

# 22. Command Preview

执行前预览：

```text
export.exe

--project demo.ap20

--overwrite
```

---

便于调试。

---

# 23. Task Panel

显示任务列表。

---

位置：

```text
底部
```

---

组件：

```tsx
<TaskPanel />
```

---

# 24. Task Item

显示：

```text
Task ID

Plugin

Status

Duration
```

---

状态：

```text
Created

Queued

Running

Completed

Failed

Cancelled
```

---

# 25. Task Store

```ts
interface TaskStore {
    tasks: TaskDto[]

    addTask()

    updateTask()

    removeTask()
}
```

---

# 26. Task Click

点击任务：

```text
Open Output Tab
```

---

# 27. Output Panel

实时日志窗口。

---

组件：

```tsx
<OutputViewer />
```

---

支持：

```text
Auto Scroll

Search

Copy

Export
```

---

# 28. Output Model

```ts
interface OutputLine {
    timestamp: string

    source: string

    level: string

    message: string
}
```

---

# 29. Output Store

```ts
interface OutputStore {
    outputs: Map<string, OutputLine[]>
}
```

---

key：

```text
taskId
```

---

# 30. Large Output Handling

日志可能：

```text
100,000+

1,000,000+
```

---

必须：

```text
Virtualized Rendering
```

---

推荐：

```text
react-virtuoso
```

---

# 31. Output Event Flow

```text
stdout

 ↓

Go

 ↓

OutputEvent

 ↓

Wails Event

 ↓

Output Store

 ↓

Output Viewer
```

---

# 32. Wails Events

统一监听。

---

封装：

```ts
useWailsEvent()
```

---

# 33. Event Registration

```ts
useWailsEvent(
    "task:started",
    handler
)
```

---

# 34. Supported Events

```text
plugin:reloaded

task:created

task:started

task:completed

task:failed

task:cancelled

output:append
```

---

# 35. Zustand Stores

推荐拆分：

```text
pluginStore

tabStore

taskStore

outputStore

settingsStore
```

---

避免：

```text
Global Mega Store
```

---

# 36. Plugin Store

```ts
interface PluginStore {
    plugins: PluginSummaryDto[]
}
```

---

# 37. Tab Store

```ts
interface TabStore {
    tabs: Tab[]

    activeTabId?: string
}
```

---

# 38. Task Store

```ts
interface TaskStore {
    tasks: TaskDto[]
}
```

---

# 39. Output Store

```ts
interface OutputStore {
    outputs:
        Record<
            string,
            OutputLine[]
        >
}
```

---

# 40. Query Strategy

使用：

```text
TanStack Query
```

---

缓存：

```text
Plugins

Settings
```

---

不缓存：

```text
Output

Running Tasks
```

通过事件驱动。

---

# 41. Theme System

支持：

```text
Light

Dark

System
```

---

实现：

```text
next-themes
```

即使不是Next.js也可使用。

---

# 42. Error Boundary

必须实现：

```tsx
<ErrorBoundary>
```

---

避免：

```text
Schema错误导致整个UI崩溃
```

---

# 43. Loading States

统一组件：

```tsx
<LoadingState />

<EmptyState />

<ErrorState />
```

---

# 44. Accessibility

要求：

```text
Keyboard Navigation

Focus Management

ARIA Labels
```

---

兼容：

```text
Windows

Linux

macOS
```

---

# 45. Frontend Service Layer

禁止：

```ts
window.go...
```

出现在组件中。

---

统一：

```ts
services/
```

---

例如：

```ts
pluginApi.ts

taskApi.ts

settingsApi.ts
```

---

# 46. Plugin API

```ts
getPlugins()

getPlugin()

reloadPlugins()
```

---

# 47. Task API

```ts
runPlugin()

cancelTask()

getTasks()

getTask()
```

---

# 48. Component Hierarchy

```text
App

 ├ Sidebar

 │   └ PluginTree

 ├ Workspace

 │   └ DynamicForm

 ├ TaskPanel

 │   └ TaskList

 └ OutputPanel

     └ OutputViewer
```

---

# 49. Recommended shadcn Components

建议使用：

```text
Button

Input

Textarea

Select

Combobox

Tabs

Sheet

Dialog

Tooltip

ScrollArea

Accordion

TreeView(自定义)
```

---

# 50. Future Extensions

V2：

```text
Workflow Editor

Plugin Marketplace

Plugin Update

Saved Profiles
```

---

V3：

```text
Remote Agent

Distributed Tasks

Multi Host Execution

Cluster Dashboard
```

---

# 51. Summary

核心页面：

```text
Sidebar

Workspace

TaskPanel

OutputPanel
```

核心状态：

```text
Plugin

Tab

Task

Output
```

核心技术：

```text
React

TypeScript

shadcn/ui

react-jsonschema-form

Zustand

TanStack Query
```

核心思想：

```text
Plugin First

Schema First

Task Driven

Event Driven
```

数据流最终统一为：

```text
Plugin

 ↓

Schema

 ↓

Form

 ↓

Task

 ↓

Process

 ↓

Output

 ↓

Event

 ↓

UI
```

对于你的项目，我会再做两处升级：

1. **OutputPanel 不要固定在底部**，改成和 VSCode 一样的可停靠 Panel（Bottom / Right / Hidden），否则日志多的时候会严重挤占表单空间。
2. **Task 与 Tab 完全解耦**，推荐采用：

```text
Plugin Tab
    ↓ Run
Task

Task
    ↓ Open
Output Tab
```

而不是：

```text
Plugin Tab
    ↓
Task Tab
    ↓
Output Tab
```

这样一个表单 Tab 可以启动多个任务，一个任务也可以单独查看日志，更符合你的 CLI Launcher 场景。
