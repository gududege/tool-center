package task

import (
	"context"
	"time"
)

// Task 任务领域聚合根
type Task struct {
	ID           string
	PluginID     string
	SessionID    string
	Status       TaskStatus
	CreatedAt    time.Time
	StartedAt    *time.Time
	EndedAt      *time.Time
	ExitCode     *int
	ErrorMessage string
	ProcessID    *string
	Cancel       context.CancelFunc
}

// CanTransitionTo 判断是否可以转换到目标状态
func (t *Task) CanTransitionTo(target TaskStatus) bool {
	if t.Status == target {
		return false
	}
	for _, s := range t.Status.ValidTransitions() {
		if s == target {
			return true
		}
	}
	return false
}
