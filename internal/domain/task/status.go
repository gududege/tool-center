package task

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskCreated    TaskStatus = "created"
	TaskQueued     TaskStatus = "queued"
	TaskRunning    TaskStatus = "running"
	TaskCompleted  TaskStatus = "completed"
	TaskFailed     TaskStatus = "failed"
	TaskCancelled  TaskStatus = "cancelled"
)

// IsTerminal 判断是否为终态
func (s TaskStatus) IsTerminal() bool {
	return s == TaskCompleted || s == TaskFailed || s == TaskCancelled
}

// ValidTransitions 返回合法状态转换
func (s TaskStatus) ValidTransitions() []TaskStatus {
	switch s {
	case TaskCreated:
		return []TaskStatus{TaskQueued}
	case TaskQueued:
		return []TaskStatus{TaskRunning}
	case TaskRunning:
		return []TaskStatus{TaskCompleted, TaskFailed, TaskCancelled}
	default:
		return nil
	}
}
