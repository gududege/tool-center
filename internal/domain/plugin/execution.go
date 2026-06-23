package plugin

// ExecutionDefinition 插件执行定义
type ExecutionDefinition struct {
	Executable       string
	WorkingDirectory string
	Environment      map[string]string
	Parameters       []ParameterMapping
}
