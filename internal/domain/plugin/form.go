package plugin

import "encoding/json"

// FormDefinition 表单定义
//
// Schema/UISchema 以 json.RawMessage 保存原始 JSON 字节，而非 map[string]any。
// 原因：Go 的 encoding/json 序列化 map 时会按 key 字母序排序，导致前端拿到的
// 表单字段顺序与 schema.json 定义不一致。使用 RawMessage 透传原始字节，
// 前端 JSON.parse 保留对象键的插入顺序，RJSF 即按文件定义顺序渲染字段。
type FormDefinition struct {
	SchemaPath   string
	UISchemaPath string
	Schema       json.RawMessage
	UISchema     json.RawMessage
}
