package task

// TaskStatus 任务状态
type TaskStatus string

const (
	Created    TaskStatus = "created"
	Queued     TaskStatus = "queued"
	Running    TaskStatus = "running"
	Completed  TaskStatus = "completed"
	Failed     TaskStatus = "failed"
	Cancelling TaskStatus = "cancelling"
	Cancelled  TaskStatus = "cancelled"
	Stopping   TaskStatus = "stopping"
	Stopped    TaskStatus = "stopped"
)

// IsTerminal 判断是否为终态
func (s TaskStatus) IsTerminal() bool {
	return s == Completed || s == Failed || s == Cancelled || s == Stopped
}

// ValidTransitions 返回合法状态转换
func (s TaskStatus) ValidTransitions() []TaskStatus {
	switch s {
	case Created:
		return []TaskStatus{Queued}
	case Queued:
		return []TaskStatus{Running}
	case Running:
		return []TaskStatus{Completed, Failed, Cancelling, Stopping}
	case Cancelling:
		return []TaskStatus{Cancelled, Failed}
	case Stopping:
		return []TaskStatus{Stopped, Failed}
	default:
		return nil
	}
}
