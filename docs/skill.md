# Plugin Authoring Skill

**Project:** CLI Tool Center  
**Version:** 1.0  
**Last Updated:** 2026-06-23

---

## 1. 目录结构模板

```text
plugins/<plugin-id>/
├── plugin.json
├── schema.json
├── uischema.json
├── <executable>
└── README.md
```

---

## 2. plugin.json 速查

```json
{
  "$schema": "../../docs/plugin.schema.json",
  "schemaVersion": "1.0",
  "id": "<plugin-id>",
  "name": "<English Name>",
  "name_cn": "<中文名>",
  "description": "...",
  "description_cn": "...",
  "version": "1.0.0",
  "author": "...",
  "homepage": "...",
  "navigation": {
    "group": ["Group", "Subgroup"],
    "group_cn": ["分组", "子分组"],
    "order": 1
  },
  "form": {
    "schema": "schema.json",
    "uischema": "uischema.json"
  },
  "execution": {
    "exe": "<executable>",
    "parameters": []
  }
}
```

---

## 3. Parameter Mapping 速查

```json
{ "field": "x", "kind": "argument" }
{ "field": "x", "kind": "argument-array" }
{ "field": "x", "kind": "option", "flag": "--x" }
{ "field": "x", "kind": "option-array", "flag": "--x", "style": "repeat" }
{ "field": "x", "kind": "option-array", "flag": "--x", "style": "join", "separator": "," }
{ "field": "x", "kind": "option-array", "flag": "--x", "style": "equals", "separator": "," }
{ "field": "x", "kind": "switch", "flag": "--x" }
{ "field": "x", "kind": "bool-option", "flag": "--x" }
{ "field": "x", "kind": "dual-switch", "trueFlag": "--x", "falseFlag": "--no-x" }
```

---

## 4. Widget 速查

### 4.1 本地自定义（优先使用）

| 效果 | uischema Control options |
| --- | --- |
| 目录选择 | `{ "folderPicker": true }` |
| 文件选择 | `{ "filePicker": true }` |
| 保存文件 | `{ "saveFilePicker": true }` |

### 4.2 RJSF 内置 widget

| 效果 | options |
| --- | --- |
| 多行文本 | `{ "widget": "textarea" }` |
| 单选按钮 | `{ "widget": "radio" }` |
| 下拉选择 | `{ "widget": "select" }` |
| 隐藏字段 | `{ "widget": "hidden" }` |
| 范围滑块 | `{ "widget": "range" }` |
| 数字步进 | `{ "widget": "updown" }` |
| 密码输入 | `{ "inputType": "password" }` |

### 4.3 Format 自动控件

| Format | 控件 |
| --- | --- |
| `email` | email |
| `uri` / `url` | url |
| `date` | date |
| `date-time` | datetime-local |
| `time` | time |
| `color` | color |
| `password` | password |

---

## 5. schema.json 字段速查

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["fieldA"],
  "properties": {
    "fieldA": {
      "type": "string",
      "title": "Field A",
      "title_cn": "字段 A",
      "description": "...",
      "description_cn": "...",
      "default": "default value"
    },
    "fieldB": {
      "type": "integer",
      "minimum": 0,
      "maximum": 100,
      "default": 50
    },
    "fieldC": {
      "type": "boolean",
      "default": true
    },
    "fieldD": {
      "type": "string",
      "enum": ["a", "b", "c"],
      "default": "a"
    },
    "fieldE": {
      "type": "array",
      "items": { "type": "string", "enum": ["x", "y", "z"] },
      "uniqueItems": true
    }
  }
}
```

---

## 6. uischema.json 模板

```json
{
  "type": "VerticalLayout",
  "elements": [
    {
      "type": "Group",
      "label": "Group Name",
      "label_cn": "分组名称",
      "elements": [
        { "type": "Control", "scope": "#/properties/fieldA" },
        { "type": "Control", "scope": "#/properties/fieldB", "options": { "widget": "radio" } }
      ]
    }
  ]
}
```

---

## 7. 常用命令

```bash
# 验证 Go 后端
go test ./...

# 验证前端
cd frontend
npm run test
npm run build

# 重新生成 Wails 绑定（修改 Go API 后）
wails generate module
```

---

## 8. 避坑清单

- `id` 必须全局唯一，且只能使用 `a-zA-Z0-9._-`。
- `schema.json` 必须使用 `$schema: https://json-schema.org/draft/2020-12/schema`。
- `uischema.json` 使用 JSON Forms 格式，`scope` 为 `#/properties/<field>`。
- 本地选择器使用 `options.folderPicker` / `options.filePicker` / `options.saveFilePicker`，不是 format。
- 通用 widget 使用 `options.widget`，不是 `options.format`。
- `switch` 在 false 时不输出任何参数；需要显式 false 时用 `bool-option` 或 `dual-switch`。
- 数组默认渲染为复选框组；需要下拉多选用 `options.widget: select`。
- 表单校验通过后才会保存历史，保存历史不依赖运行结果。

---

## 9. 参考

- [Plugin Authoring Guide](./usage.md)
- [Plugin Specification](./Plugin%20Specification.md)
- [RJSF Widgets](https://rjsf-team.github.io/react-jsonschema-form/docs/usage/widgets)
