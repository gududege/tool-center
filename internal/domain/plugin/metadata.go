package plugin

// PluginMetadata 插件元信息
type PluginMetadata struct {
	ID          string
	Name        string
	NameCn      string // 中文名称
	Description string
	DescriptionCn string // 中文描述
	Version     string
	Author      string
	Homepage    string
	Icon        string
}
