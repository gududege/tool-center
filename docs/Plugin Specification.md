# Plugin Specification

**Project:** CLI Tool Center  
**Version:** 1.0  
**Status:** Active  
**Last Updated:** 2026-06-23

---

## 1. 概述

CLI Tool Center 的插件是声明式的：一个插件目录包含若干 JSON 文件和一个可执行文件，主程序负责加载、渲染表单并执行 CLI。

插件负责声明：

- 元数据（名称、描述、图标）
- 导航位置
- 表单 Schema / UISchema
- CLI 参数映射
- 执行方式

主程序负责：

- 插件发现与加载
- 表单渲染（RJSF + shadcn/ui）
- 参数构造
- 任务与进程管理
- 输出展示

---

## 2. 目录结构

```text
plugins/
└─ <plugin-id>/
   ├─ plugin.json       # 必填，插件主配置
   ├─ schema.json       # 必填，表单 JSON Schema
   ├─ uischema.json     # 可选，表单布局
   ├─ <executable>      # 必填，CLI 可执行文件
   ├─ icon.svg          # 可选，图标
   └─ README.md         # 可选，说明文档
```

---

## 3. plugin.json

### 3.1 完整示例

```json
{
  "$schema": "../../docs/plugin.schema.json",
  "schemaVersion": "1.0",
  "id": "tia-export",
  "name": "TIA Portal Export",
  "name_cn": "TIA Portal 导出",
  "description": "Export content from Siemens TIA Portal projects",
  "description_cn": "从西门子 TIA Portal 项目导出内容",
  "version": "1.0.0",
  "author": "gududege",
  "homepage": "https://github.com/gududege/tia-exporter",
  "icon": "./icon.svg",
  "navigation": {
    "group": ["Industrial Automation", "TIA Portal"],
    "group_cn": ["工业自动化", "TIA Portal"],
    "order": 1
  },
  "form": {
    "schema": "schema.json",
    "uischema": "uischema.json"
  },
  "execution": {
    "exe": "tia-export.exe",
    "workingDirectory": ".",
    "environment": {
      "LOG_LEVEL": "info"
    },
    "parameters": [
      { "field": "indir",  "kind": "option", "flag": "--indir" },
      { "field": "outdir", "kind": "option", "flag": "--outdir" },
      { "field": "keepFolderStructure", "kind": "switch", "flag": "--keep-folder-structure" }
    ]
  },
  "output": {
    "type": "text"
  },
  "capabilities": ["file-input", "file-output"]
}
```

### 3.2 字段说明

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `schemaVersion` | 是 | 固定 `"1.0"` |
| `id` | 是 | 唯一标识，`^[a-zA-Z0-9._-]+$` |
| `name` | 是 | 英文显示名称 |
| `name_cn` | 否 | 中文显示名称 |
| `description` | 否 | 英文描述 |
| `description_cn` | 否 | 中文描述 |
| `version` | 否 | 插件版本 |
| `author` | 否 | 作者 |
| `homepage` | 否 | 主页 |
| `icon` | 否 | 图标路径 |
| `navigation` | 是 | 导航配置 |
| `form` | 否 | 表单配置 |
| `execution` | 是 | 执行配置 |
| `output` | 否 | 输出类型提示 |
| `capabilities` | 否 | 能力标签 |

### 3.3 navigation

```json
{
  "navigation": {
    "group": ["Network", "TCPING"],
    "group_cn": ["网络", "TCPING"],
    "order": 1
  }
}
```

- `group`：菜单路径，支持多级嵌套。
- `group_cn`：中文菜单路径（可选）。
- `order`：同组内排序，数值越小越靠前。

### 3.4 form

```json
{
  "form": {
    "schema": "schema.json",
    "uischema": "uischema.json"
  }
}
```

- `schema`：JSON Schema 文件路径，必填。
- `uischema`：布局文件路径，可选；缺失时使用 RJSF 默认布局。

### 3.5 execution

```json
{
  "execution": {
    "exe": "tcping.exe",
    "workingDirectory": ".",
    "environment": { "KEY": "value" },
    "parameters": []
  }
}
```

- `exe`：可执行文件路径，相对插件目录或绝对路径。
- `workingDirectory`：工作目录，可选。
- `environment`：额外环境变量。
- `parameters`：表单字段到 CLI 参数的映射，见第 4 节。

---

## 4. Parameter Mapping

JSON Schema 负责 UI、校验和默认值，不直接决定 CLI 参数格式，因此需要独立的 `parameters` 映射。

### 4.1 参数类型

| Kind | 说明 | 示例 |
| --- | --- | --- |
| `argument` | 位置参数 | `tool.exe value` |
| `argument-array` | 多个位置参数 | `tool.exe a b c` |
| `option` | `--flag value` | `--project demo.ap20` |
| `option-array` | 数组型选项 | 见 4.2 |
| `switch` | 为 true 时输出 flag | `--overwrite` |
| `bool-option` | 输出 `flag true` / `flag false` | `--overwrite true` |
| `dual-switch` | true/false 各对应一个 flag | `--overwrite` / `--no-overwrite` |

### 4.2 option-array 风格

```json
{ "field": "tags", "kind": "option-array", "flag": "--tag", "style": "repeat" }
```

输出：

```bash
--tag PLC1 --tag PLC2 --tag PLC3
```

```json
{ "field": "tags", "kind": "option-array", "flag": "--tags", "style": "join", "separator": "," }
```

输出：

```bash
--tags PLC1,PLC2,PLC3
```

```json
{ "field": "tags", "kind": "option-array", "flag": "--tags", "style": "equals", "separator": "," }
```

输出：

```bash
--tags=PLC1,PLC2,PLC3
```

### 4.3 参数顺序

默认按 `parameters` 数组的定义顺序输出。

---

## 5. schema.json

采用 JSON Schema draft 2020-12。前端使用 RJSF（react-jsonschema-form）+ Ajv2020 渲染和校验。

### 5.1 最小示例

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["host", "port"],
  "properties": {
    "host": {
      "type": "string",
      "title": "Hostname or IP",
      "title_cn": "主机名或 IP",
      "description": "Target hostname"
    },
    "port": {
      "type": "integer",
      "title": "Port",
      "minimum": 1,
      "maximum": 65535,
      "default": 80
    }
  }
}
```

### 5.2 本地化字段

- `title_cn`：中文标题（前端会根据当前语言替换 `title`）。
- `description_cn`：中文描述。
- `anyOf` / `oneOf` 子项中的 `title_cn` 同样会被本地化。

### 5.3 支持的 JSON Schema 类型

- `string`：文本输入
- `integer`：整数
- `number`：浮点数
- `boolean`：布尔值
- `array`：数组
- `object`：嵌套对象
- `enum` / `anyOf` / `oneOf`：枚举/多选

### 5.4 内置 format

RJSF 会根据 `format` 自动选择控件：

| Format | 默认控件 |
| --- | --- |
| `email` | email 输入 |
| `uri` / `url` | url 输入 |
| `date` | 日期选择器 |
| `date-time` | 日期时间选择器 |
| `time` | 时间选择器 |
| `color` | 颜色选择器 |
| `password` | 密码输入 |
| `data-url` | 文件选择（本地会被 folderPicker / saveFilePicker 覆盖） |

---

## 6. uischema.json

CLI Tool Center 使用 **JSON Forms 风格** 的布局格式，在运行时自动展平为 RJSF 的 `uiSchema`。

### 6.1 布局元素

```json
{
  "type": "VerticalLayout",
  "elements": [
    {
      "type": "Group",
      "label": "Required Paths",
      "label_cn": "必需路径",
      "elements": [
        { "type": "Control", "scope": "#/properties/indir" },
        { "type": "Control", "scope": "#/properties/outdir" }
      ]
    }
  ]
}
```

- `VerticalLayout` / `HorizontalLayout`：布局容器。
- `Group`：带边框卡片分组，会收集组内字段并在前端渲染为卡片。
- `Categorization` / `Category`：分组容器（当前按普通容器处理）。
- `Control`：对应一个 schema 字段，`scope` 格式为 `#/properties/<field>`。

### 6.2 Control options

`Control` 的 `options` 支持以下字段：

| Option | 说明 | 对应 RJSF uiSchema |
| --- | --- | --- |
| `folderPicker: true` | 目录选择器 | `ui:widget: folderPicker` |
| `filePicker: true` | 文件选择器 | `ui:widget: filePicker` |
| `saveFilePicker: true` | 保存文件选择器 | `ui:widget: saveFilePicker` |
| `widget: "<name>"` | RJSF widget 名称 | `ui:widget: <name>` |
| `inputType: "<type>"` | HTML input type | `ui:inputType: <type>` |

### 6.3 支持的 widget 名称

**RJSF 内置 / shadcn 主题提供：**

| Widget | 适用类型 | 说明 |
| --- | --- | --- |
| `checkbox` | boolean | 复选框（默认） |
| `radio` | boolean / string(enum) | 单选按钮组 |
| `select` | boolean / string(enum) / array(enum) | 下拉选择 |
| `textarea` | string | 多行文本 |
| `hidden` | string / number / integer / boolean | 隐藏字段 |
| `range` | number / integer | 滑块 |
| `updown` | number / integer | 数字步进输入 |
| `color` | string | 颜色选择器（也可通过 `format: color`） |
| `password` | string | 密码输入（也可通过 `format: password`） |

**本地自定义 widget：**

| Widget | 适用类型 | 说明 |
| --- | --- | --- |
| `folderPicker` | string | 目录选择，返回路径字符串 |
| `filePicker` | string | 文件选择，返回路径字符串 |
| `saveFilePicker` | string | 保存文件对话框，返回路径字符串 |

### 6.4 完整 widget 示例

```json
{
  "type": "VerticalLayout",
  "elements": [
    {
      "type": "Group",
      "label": "Paths",
      "elements": [
        {
          "type": "Control",
          "scope": "#/properties/indir",
          "options": { "folderPicker": true }
        },
        {
          "type": "Control",
          "scope": "#/properties/outfile",
          "options": { "saveFilePicker": true }
        }
      ]
    },
    {
      "type": "Group",
      "label": "Options",
      "elements": [
        {
          "type": "Control",
          "scope": "#/properties/mode",
          "options": { "widget": "radio" }
        },
        {
          "type": "Control",
          "scope": "#/properties/notes",
          "options": { "widget": "textarea" }
        },
        {
          "type": "Control",
          "scope": "#/properties/level",
          "options": { "widget": "range" }
        }
      ]
    }
  ]
}
```

---

## 7. 生命周期

```text
Application Start
  ↓
Scan plugins/*
  ↓
Read plugin.json
  ↓
Validate against plugin.schema.json
  ↓
Resolve executable / schema / uischema paths
  ↓
Register Plugin
  ↓
Build Navigation Tree
  ↓
Ready
```

单个插件加载失败不影响其他插件。

---

## 8. 校验规则

- `id` 必须唯一。
- `name`、`navigation`、`execution` 必填。
- `execution.exe` 必须存在且可执行。
- `form.schema` 若定义则文件必须存在。
- `schema.json` 必须是合法 JSON Schema draft 2020-12。

---

## 9. 历史记录

运行时会在插件目录下自动生成 `history.json`，保存最近 5 次成功校验的表单参数，供用户快速加载。

---

## 10. 参考

- [React JSON Schema Form 官方文档](https://rjsf-team.github.io/react-jsonschema-form/docs/)
- [RJSF Widgets 列表](https://rjsf-team.github.io/react-jsonschema-form/docs/usage/widgets)
- [JSON Schema 2020-12](https://json-schema.org/draft/2020-12/schema)
- [JSON Forms](https://jsonforms.io/)
