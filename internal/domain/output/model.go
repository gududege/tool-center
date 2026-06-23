package output

import "time"

// OutputSource 输出源
type OutputSource string

const (
	StdoutSource OutputSource = "stdout"
	StderrSource OutputSource = "stderr"
	SystemSource OutputSource = "system"
)

// LogLevel 日志级别
type LogLevel string

const (
	TraceLevel LogLevel = "trace"
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

// OutputEvent 统一输出事件模型
type OutputEvent struct {
	ID        string
	TaskID    string
	Timestamp time.Time
	Source    OutputSource
	Level     LogLevel
	Message   string
}

// OutputBuffer 任务日志缓存（只追加）
type OutputBuffer struct {
	TaskID string
	Events []OutputEvent
	MaxLen int
}

// Append 追加事件，超出限制时FIFO淘汰
func (b *OutputBuffer) Append(event OutputEvent) {
	if b.MaxLen > 0 && len(b.Events) >= b.MaxLen {
		b.Events = b.Events[len(b.Events)-b.MaxLen+1:]
	}
	b.Events = append(b.Events, event)
}

// GetEvents 获取所有事件
func (b *OutputBuffer) GetEvents() []OutputEvent {
	result := make([]OutputEvent, len(b.Events))
	copy(result, b.Events)
	return result
}
