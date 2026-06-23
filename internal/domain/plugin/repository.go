package plugin

// PluginRepository 插件仓储接口
type PluginRepository interface {
	Get(id string) (*Plugin, error)
	List() ([]*Plugin, error)
	Save(plugin *Plugin) error
	Delete(id string) error
}
