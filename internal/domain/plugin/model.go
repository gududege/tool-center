package plugin

// Plugin 领域聚合根——是整个系统的扩展单元
type Plugin struct {
	Metadata   PluginMetadata
	Navigation Navigation
	Form       FormDefinition
	Execution  ExecutionDefinition
	// Directory 是插件在文件系统中的目录，用于读写 history.json 等插件本地文件。
	Directory string
}
