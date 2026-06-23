package wails

import "encoding/json"

// PluginSummaryDto 插件摘要 DTO
type PluginSummaryDto struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	NameCn        string        `json:"name_cn,omitempty"`
	Description   string        `json:"description,omitempty"`
	DescriptionCn string        `json:"description_cn,omitempty"`
	Version       string        `json:"version,omitempty"`
	Icon          string        `json:"icon,omitempty"`
	Navigation    NavigationDto `json:"navigation"`
}

// PluginDto 完整插件 DTO
type PluginDto struct {
	Metadata   PluginMetadataDto      `json:"metadata"`
	Navigation NavigationDto          `json:"navigation"`
	Form       FormDefinitionDto      `json:"form"`
	Execution  ExecutionDefinitionDto `json:"execution"`
}

// PluginMetadataDto 插件元数据 DTO
type PluginMetadataDto struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	NameCn        string `json:"name_cn,omitempty"`
	Description   string `json:"description,omitempty"`
	DescriptionCn string `json:"description_cn,omitempty"`
	Version       string `json:"version,omitempty"`
	Author        string `json:"author,omitempty"`
	Icon          string `json:"icon,omitempty"`
}

// NavigationDto 导航 DTO
type NavigationDto struct {
	Group   []string `json:"group"`
	GroupCn []string `json:"group_cn,omitempty"`
	Order   int      `json:"order"`
}

// FormDefinitionDto 表单定义 DTO
//
// Schema/UISchema 用 json.RawMessage 透传原始 JSON 字节，保留字段定义顺序。
type FormDefinitionDto struct {
	Schema   json.RawMessage `json:"schema,omitempty"`
	UISchema json.RawMessage `json:"uiSchema,omitempty"`
}

// ExecutionDefinitionDto 执行定义 DTO
type ExecutionDefinitionDto struct {
	Exe              string                `json:"exe"`
	WorkingDirectory string                `json:"workingDirectory,omitempty"`
	Parameters       []ParameterMappingDto `json:"parameters,omitempty"`
}

// ParameterMappingDto 参数映射 DTO
type ParameterMappingDto struct {
	Field        string `json:"field"`
	Kind         string `json:"kind"`
	Flag         string `json:"flag,omitempty"`
	Style        string `json:"style,omitempty"`
	Separator    string `json:"separator,omitempty"`
	TrueFlag     string `json:"trueFlag,omitempty"`
	FalseFlag    string `json:"falseFlag,omitempty"`
	DefaultValue any    `json:"defaultValue,omitempty"`
}

// TaskDto 任务 DTO
type TaskDto struct {
	ID        string `json:"id"`
	PluginID  string `json:"pluginId"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	StartedAt string `json:"startedAt,omitempty"`
	EndedAt   string `json:"endedAt,omitempty"`
}

// OutputEventDto 输出事件 DTO
type OutputEventDto struct {
	TaskID    string `json:"taskId"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Source    string `json:"source"`
	Message   string `json:"message"`
}

// SettingsDto 配置 DTO
type SettingsDto struct {
	Theme            string `json:"theme"`
	FormTheme        string `json:"formTheme,omitempty"`
	PluginDirectory  string `json:"pluginDirectory"`
	Language         string `json:"language,omitempty"`
	SidebarCollapsed bool   `json:"sidebarCollapsed"`
	SidebarSize      int    `json:"sidebarSize"`
	BottomPanelSize  int    `json:"bottomPanelSize"`
	BottomTab        string `json:"bottomTab"`
}

// SystemInfoDto 系统信息 DTO
type SystemInfoDto struct {
	AppVersion string `json:"appVersion"`
	BuildTime  string `json:"buildTime"`
	GoVersion  string `json:"goVersion"`
	Os         string `json:"os"`
}

// ApiErrorDto 统一错误 DTO
type ApiErrorDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// RunPluginRequest 执行插件请求
type RunPluginRequest struct {
	PluginID string         `json:"pluginId"`
	FormData map[string]any `json:"formData"`
}

// RunPluginResponse 执行插件响应
type RunPluginResponse struct {
	TaskID string `json:"taskId"`
}

// CancelTaskResponse 取消任务响应
type CancelTaskResponse struct {
	Success bool `json:"success"`
}

// DeleteTaskResponse 删除任务响应
type DeleteTaskResponse struct {
	Success bool `json:"success"`
}

// ReloadPluginsResponse 重载插件响应
type ReloadPluginsResponse struct {
	Count int `json:"count"`
}

// ValidationResultDto 表单验证结果
type ValidationResultDto struct {
	Valid  bool                 `json:"valid"`
	Errors []ValidationErrorDto `json:"errors,omitempty"`
}

// ValidationErrorDto 验证错误
type ValidationErrorDto struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// PluginHistoryEntryDto 插件参数历史条目
type PluginHistoryEntryDto struct {
	Timestamp string         `json:"timestamp"`
	Label     string         `json:"label"`
	FormData  map[string]any `json:"formData"`
}

// FileFilterDto 文件选择过滤器
type FileFilterDto struct {
	DisplayName string   `json:"displayName"`
	Patterns    []string `json:"patterns"`
}

// DialogResponse 对话框响应
type DialogResponse struct {
	Path string `json:"path,omitempty"`
}

// Event contracts
type PluginsReloadedEvent struct {
	Count int `json:"count"`
}

type TaskCreatedEvent struct {
	TaskID   string `json:"taskId"`
	PluginID string `json:"pluginId"`
}

type TaskStartedEvent struct {
	TaskID string `json:"taskId"`
}

type TaskCompletedEvent struct {
	TaskID string `json:"taskId"`
}

type TaskFailedEvent struct {
	TaskID string `json:"taskId"`
	Error  string `json:"error"`
}

type TaskCancelledEvent struct {
	TaskID string `json:"taskId"`
}

type TaskStatusChangedEvent struct {
	TaskID    string `json:"taskId"`
	OldStatus string `json:"oldStatus"`
	NewStatus string `json:"newStatus"`
}
