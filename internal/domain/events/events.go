package events

// DomainEvent 领域事件接口
type DomainEvent interface {
	Topic() string
}

// PluginLoaded 插件加载事件
type PluginLoaded struct {
	PluginID string
}

func (e PluginLoaded) Topic() string { return "plugin.loaded" }

// PluginUnloaded 插件卸载事件
type PluginUnloaded struct {
	PluginID string
}

func (e PluginUnloaded) Topic() string { return "plugin.unloaded" }

// PluginReloaded 插件重载事件
type PluginReloaded struct {
	PluginID string
}

func (e PluginReloaded) Topic() string { return "plugin.reloaded" }

// --- Task Events ---

// TaskCreated 任务创建事件
type TaskCreated struct {
	TaskID   string
	PluginID string
}

func (e TaskCreated) Topic() string { return "task.created" }

// TaskStarted 任务开始事件
type TaskStarted struct {
	TaskID string
}

func (e TaskStarted) Topic() string { return "task.started" }

// TaskCompleted 任务完成事件
type TaskCompleted struct {
	TaskID   string
	ExitCode int
}

func (e TaskCompleted) Topic() string { return "task.completed" }

// TaskFailed 任务失败事件
type TaskFailed struct {
	TaskID string
	Error  string
}

func (e TaskFailed) Topic() string { return "task.failed" }

// TaskCancelled 任务取消事件
type TaskCancelled struct {
	TaskID string
}

func (e TaskCancelled) Topic() string { return "task.cancelled" }

// TaskStatusChanged 任务状态变更事件
type TaskStatusChanged struct {
	TaskID    string
	OldStatus string
	NewStatus string
}

func (e TaskStatusChanged) Topic() string { return "task.status-changed" }

// --- Output Events ---

// OutputReceived 输出事件
type OutputReceived struct {
	TaskID string
	Event  interface{} // OutputEvent DTO
}

func (e OutputReceived) Topic() string { return "output.received" }
