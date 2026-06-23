package process

import "time"

// ProcessInstance CLI进程实例
type ProcessInstance struct {
	ID               string
	TaskID           string
	PID              int
	Command          string
	Arguments        []string
	WorkingDirectory string
	Environment      []string
	StartedAt        time.Time
}
