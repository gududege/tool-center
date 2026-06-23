# Plugin Authoring Guide

**Project:** CLI Tool Center  
**Version:** 1.0  
**Last Updated:** 2026-06-23

---

## 1. 快速开始

一个最小插件只需要：

```text
plugins/my-tool/
├── plugin.json
├── schema.json
└── my-tool.exe
```

1. 在 `plugins/` 下新建一个目录，目录名即为插件 ID。
2. 创建 `plugin.json` 描述插件。
3. 创建 `schema.json` 描述表单。
4. 放入可执行文件。
5. 重启 CLI Tool Center，左侧导航树会自动出现该插件。

---

## 2. plugin.json 详解

```json
{
  "$schema": "../../docs/plugin.schema.json",
  "schemaVersion": "1.0",
  "id": "my-tool",
  "name": "My Tool",
  "name_cn": "我的工具",
  "description": "A short description",
  "description_cn": "简短描述",
  "version": "1.0.0",
  "author": "Your Name",
  "homepage": "https://github.com/you/my-tool",
  "icon": "./icon.svg",
  "navigation": {
    "group": ["Tools", "My Tool"],
    "group_cn": ["工具", "我的工具"],
    "order": 10
  },
  "form": {
    "schema": "schema.json",
    "uischema": "uischema.json"
  },
  "execution": {
    "exe": "my-tool.exe",
    "workingDirectory": ".",
    "environment": { "LOG_LEVEL": "info" },
    "parameters": [
      { "field": "input", "kind": "option", "flag": "--input" },
      { "field": "output", "kind": "option", "flag": "--output" },
      { "field": "verbose", "kind": "switch", "flag": "--verbose" }
    ]
  }
}
```

### 2.1 字段速查

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `id` | 是 | 唯一标识，仅含 `a-zA-Z0-9._-` |
| `name` | 是 | 英文名称 |
| `name_cn` | 否 | 中文名称 |
| `description` / `description_cn` | 否 | 描述 |
| `navigation.group` | 是 | 菜单路径数组，支持多级 |
| `navigation.group_cn` | 否 | 中文菜单路径 |
| `navigation.order` | 否 | 排序，越小越靠前 |
| `form.schema` | 是 | schema.json 路径 |
| `form.uischema` | 否 | uischema.json 路径 |
| `execution.exe` | 是 | 可执行文件路径 |
| `execution.workingDirectory` | 否 | 工作目录 |
| `execution.environment` | 否 | 环境变量键值对 |
| `execution.parameters` | 否 | 表单字段到 CLI 参数映射 |

---

## 3. schema.json 详解

CLI Tool Center 使用 **JSON Schema draft 2020-12** 描述表单。

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["host", "port"],
  "properties": {
    "host": {
      "type": "string",
      "title": "Host",
      "title_cn": "主机",
      "description": "Target host",
      "description_cn": "目标主机"
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

### 3.1 本地化

- `title_cn`：中文标题
- `description_cn`：中文描述
- `anyOf` / `oneOf` 子项也可使用 `title_cn`

### 3.2 字段类型与默认控件

| Schema Type | 默认控件 |
| --- | --- |
| `string` | 文本输入 |
| `string` + `enum` | 下拉选择 |
| `integer` / `number` | 数字输入 |
| `boolean` | 复选框 |
| `array` + `enum` items | 复选框组 |
| `object` | 嵌套对象卡片 |

### 3.3 常用 format

```json
{
  "emailField": { "type": "string", "format": "email" },
  "urlField": { "type": "string", "format": "uri" },
  "dateField": { "type": "string", "format": "date" },
  "dateTimeField": { "type": "string", "format": "date-time" },
  "timeField": { "type": "string", "format": "time" },
  "colorField": { "type": "string", "format": "color" },
  "passwordField": { "type": "string", "format": "password" }
}
```

---

## 4. uischema.json 详解

`uischema.json` 使用 **JSON Forms** 风格描述布局，运行时会自动展平为 RJSF 的 `uiSchema`。

### 4.1 基本结构

```json
{
  "type": "VerticalLayout",
  "elements": [
    {
      "type": "Group",
      "label": "Basic",
      "label_cn": "基础",
      "elements": [
        { "type": "Control", "scope": "#/properties/host" },
        { "type": "Control", "scope": "#/properties/port" }
      ]
    }
  ]
}
```

### 4.2 布局元素

| 元素 | 说明 |
| --- | --- |
| `VerticalLayout` | 垂直排列子元素 |
| `HorizontalLayout` | 水平排列子元素 |
| `Group` | 带边框卡片分组 |
| `Categorization` / `Category` | 分组容器 |
| `Control` | 单个字段，scope 为 `#/properties/<field>` |

### 4.3 Control Options

```json
{
  "type": "Control",
  "scope": "#/properties/indir",
  "options": { "folderPicker": true }
}
```

| Option | 值 | 说明 |
| --- | --- | --- |
| `folderPicker` | `true` | 目录选择器 |
| `filePicker` | `true` | 文件选择器 |
| `saveFilePicker` | `true` | 保存文件选择器 |
| `widget` | `"radio"`, `"select"`, `"textarea"`, `"range"`, `"hidden"`, ... | RJSF widget 名称 |
| `inputType` | `"password"`, `"number"`, ... | HTML input type |

### 4.4 支持的 widget 列表

**RJSF / shadcn 内置：**

| Widget | 适用类型 | 说明 |
| --- | --- | --- |
| `checkbox` | boolean | 复选框（默认） |
| `radio` | boolean / enum string | 单选按钮组 |
| `select` | boolean / enum string / enum array | 下拉选择 |
| `textarea` | string | 多行文本 |
| `hidden` | string / number / integer / boolean | 隐藏字段 |
| `range` | number / integer | 滑块 |
| `updown` | number / integer | 数字步进器 |
| `color` | string | 颜色选择 |
| `password` | string | 密码输入 |

**本地自定义：**

| Widget | 说明 |
| --- | --- |
| `folderPicker` | 目录选择对话框 |
| `filePicker` | 文件选择对话框 |
| `saveFilePicker` | 保存文件对话框 |

### 4.5 完整 widget 示例

```json
{
  "type": "VerticalLayout",
  "elements": [
    {
      "type": "Group",
      "label": "Inputs",
      "label_cn": "输入",
      "elements": [
        { "type": "Control", "scope": "#/properties/text" },
        { "type": "Control", "scope": "#/properties/notes", "options": { "widget": "textarea" } },
        { "type": "Control", "scope": "#/properties/secret", "options": { "inputType": "password" } }
      ]
    },
    {
      "type": "Group",
      "label": "Pickers",
      "label_cn": "选择器",
      "elements": [
        { "type": "Control", "scope": "#/properties/indir", "options": { "folderPicker": true } },
        { "type": "Control", "scope": "#/properties/infile", "options": { "filePicker": true } },
        { "type": "Control", "scope": "#/properties/outfile", "options": { "saveFilePicker": true } }
      ]
    },
    {
      "type": "Group",
      "label": "Choices",
      "label_cn": "选项",
      "elements": [
        { "type": "Control", "scope": "#/properties/mode", "options": { "widget": "radio" } },
        { "type": "Control", "scope": "#/properties/tags", "options": { "widget": "select" } }
      ]
    }
  ]
}
```

---

## 5. Parameter Mapping

`execution.parameters` 定义表单数据如何转换为 CLI 参数。

### 5.1 参数类型

| Kind | 字段值 | 输出 |
| --- | --- | --- |
| `argument` | `"value"` | `value` |
| `argument-array` | `["a", "b"]` | `a b` |
| `option` | `"value"` | `--flag value` |
| `switch` | `true` | `--flag` |
| `switch` | `false` | （无输出） |
| `bool-option` | `true` | `--flag true` |
| `bool-option` | `false` | `--flag false` |
| `dual-switch` | `true` | `--trueFlag` |
| `dual-switch` | `false` | `--falseFlag` |

### 5.2 option-array 风格

```json
{ "field": "tags", "kind": "option-array", "flag": "--tag", "style": "repeat" }
```

输出：

```bash
--tag a --tag b --tag c
```

```json
{ "field": "tags", "kind": "option-array", "flag": "--tags", "style": "join", "separator": "," }
```

输出：

```bash
--tags a,b,c
```

```json
{ "field": "tags", "kind": "option-array", "flag": "--tags", "style": "equals", "separator": "," }
```

输出：

```bash
--tags=a,b,c
```

---

## 6. 完整示例：文件拷贝工具

### plugin.json

```json
{
  "$schema": "../../docs/plugin.schema.json",
  "schemaVersion": "1.0",
  "id": "file-copy",
  "name": "File Copy",
  "name_cn": "文件拷贝",
  "description": "Copy files from source to destination",
  "description_cn": "将文件从源路径拷贝到目标路径",
  "version": "1.0.0",
  "author": "CLI Tool Center",
  "navigation": {
    "group": ["Tools", "File"],
    "group_cn": ["工具", "文件"],
    "order": 1
  },
  "form": {
    "schema": "schema.json",
    "uischema": "uischema.json"
  },
  "execution": {
    "exe": "file-copy.exe",
    "parameters": [
      { "field": "source", "kind": "option", "flag": "--source" },
      { "field": "target", "kind": "option", "flag": "--target" },
      { "field": "overwrite", "kind": "switch", "flag": "--overwrite" }
    ]
  }
}
```

### schema.json

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["source", "target"],
  "properties": {
    "source": {
      "type": "string",
      "title": "Source File",
      "title_cn": "源文件"
    },
    "target": {
      "type": "string",
      "title": "Target Directory",
      "title_cn": "目标目录"
    },
    "overwrite": {
      "type": "boolean",
      "title": "Overwrite existing files",
      "title_cn": "覆盖已存在文件",
      "default": false
    }
  }
}
```

### uischema.json

```json
{
  "type": "VerticalLayout",
  "elements": [
    {
      "type": "Group",
      "label": "Paths",
      "label_cn": "路径",
      "elements": [
        { "type": "Control", "scope": "#/properties/source", "options": { "filePicker": true } },
        { "type": "Control", "scope": "#/properties/target", "options": { "folderPicker": true } }
      ]
    },
    {
      "type": "Group",
      "label": "Options",
      "label_cn": "选项",
      "elements": [
        { "type": "Control", "scope": "#/properties/overwrite" }
      ]
    }
  ]
}
```

---

## 7. 调试技巧

1. **表单不显示**：检查 `plugin.json` 的 `form.schema` 路径是否正确，以及 `schema.json` 是否为合法 JSON。
2. **widget 不生效**：确认 `uischema.json` 中 `options.widget` 或 `options.folderPicker` 等配置正确，并检查插件是否已重新加载。
3. **参数不对**：使用底部 **Preview** 按钮实时查看生成的命令行。
4. **校验失败**：确保 `schema.json` 声明 `$schema: https://json-schema.org/draft/2020-12/schema`。
5. **历史加载**：点击 Run 且校验通过后，当前参数会自动保存到插件目录的 `history.json`，最多保留 5 条。

---

## 8. 参考

- [Plugin Specification](./Plugin%20Specification.md)
- [Frontend Design](./Frontend%20Design.md)
- [RJSF Widgets 官方文档](https://rjsf-team.github.io/react-jsonschema-form/docs/usage/widgets)
- [JSON Schema 2020-12](https://json-schema.org/draft/2020-12/schema)
