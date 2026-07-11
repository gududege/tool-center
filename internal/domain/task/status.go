package task

// TaskStatus 任务状态
type TaskStatus string

const (
	Created    TaskStatus = "created"
	Queued     TaskStatus = "queued"
	Running    TaskStatus = "running"
	Completed  TaskStatus = "completed"
	Failed     TaskStatus = "failed"
	Cancelled  TaskStatus = "cancelled"
)

// IsTerminal 判断是否为终态
func (s TaskStatus) IsTerminal() bool {
	return s == Completed || s == Failed || s == Cancelled
}

// ValidTransitions 返回合法状态转换
func (s TaskStatus) ValidTransitions() []TaskStatus {
	switch s {
	case Created:
		return []TaskStatus{Queued}
	case Queued:
		return []TaskStatus{Running}
	case Running:
		return []TaskStatus{Completed, Failed, Cancelled}
	default:
		return nil
	}
}
