package wails

import (
	"context"
	"log/slog"
	"time"

	"github.com/cli-tool-center/tool-center/internal/domain/events"
	domainOutput "github.com/cli-tool-center/tool-center/internal/domain/output"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// EventBridge 将领域事件桥接到前端 Wails 事件
type EventBridge struct {
	ctx      context.Context
	eventBus events.Bus
	cancels  []func()
}

// NewEventBridge 创建事件桥接器
func NewEventBridge(eventBus events.Bus) *EventBridge {
	return &EventBridge{
		eventBus: eventBus,
	}
}

// SetContext 设置 Wails runtime context（在 OnStartup 中调用）
func (b *EventBridge) SetContext(ctx context.Context) {
	b.ctx = ctx
}

// Start 开始监听领域事件并转发到前端
func (b *EventBridge) Start() {
	b.subscribe("plugin.reloaded", "plugins:reloaded", b.toPluginReloaded)
	b.subscribe("task.created", "task:created", b.toTaskCreated)
	b.subscribe("task.started", "task:started", b.toTaskStarted)
	b.subscribe("task.completed", "task:completed", b.toTaskCompleted)
	b.subscribe("task.failed", "task:failed", b.toTaskFailed)
	b.subscribe("task.cancelled", "task:cancelled", b.toTaskCancelled)
	b.subscribe("task.status-changed", "task:status-changed", b.toTaskStatusChanged)
	b.subscribe("output.received", "output:append", b.toOutputAppend)
}

// Stop 停止事件监听
func (b *EventBridge) Stop() {
	for _, cancel := range b.cancels {
		cancel()
	}
}

func (b *EventBridge) subscribe(domainTopic, frontendTopic string, mapper func(events.DomainEvent) any) {
	cancel := b.eventBus.Subscribe(domainTopic, func(event events.DomainEvent) {
		if b.ctx == nil {
			return
		}
		payload := mapper(event)
		if payload == nil {
			return
		}
		runtime.EventsEmit(b.ctx, frontendTopic, payload)
	})
	b.cancels = append(b.cancels, cancel)
}

func (b *EventBridge) toPluginReloaded(e events.DomainEvent) any {
	event, ok := e.(events.PluginReloaded)
	if !ok {
		return nil
	}
	_ = event
	return PluginsReloadedEvent{Count: 1}
}

func (b *EventBridge) toTaskCreated(e events.DomainEvent) any {
	event, ok := e.(events.TaskCreated)
	if !ok {
		return nil
	}
	return TaskCreatedEvent{TaskID: event.TaskID, PluginID: event.PluginID}
}

func (b *EventBridge) toTaskStarted(e events.DomainEvent) any {
	event, ok := e.(events.TaskStarted)
	if !ok {
		return nil
	}
	return TaskStartedEvent{TaskID: event.TaskID}
}

func (b *EventBridge) toTaskCompleted(e events.DomainEvent) any {
	event, ok := e.(events.TaskCompleted)
	if !ok {
		return nil
	}
	return TaskCompletedEvent{TaskID: event.TaskID}
}

func (b *EventBridge) toTaskFailed(e events.DomainEvent) any {
	event, ok := e.(events.TaskFailed)
	if !ok {
		return nil
	}
	return TaskFailedEvent{TaskID: event.TaskID, Error: event.Error}
}

func (b *EventBridge) toTaskCancelled(e events.DomainEvent) any {
	event, ok := e.(events.TaskCancelled)
	if !ok {
		return nil
	}
	return TaskCancelledEvent{TaskID: event.TaskID}
}

func (b *EventBridge) toTaskStatusChanged(e events.DomainEvent) any {
	event, ok := e.(events.TaskStatusChanged)
	if !ok {
		return nil
	}
	return TaskStatusChangedEvent{
		TaskID:    event.TaskID,
		OldStatus: event.OldStatus,
		NewStatus: event.NewStatus,
	}
}

func (b *EventBridge) toOutputAppend(e events.DomainEvent) any {
	event, ok := e.(events.OutputReceived)
	if !ok {
		return nil
	}
	// 尝试从事件中提取 OutputEvent
	if outputEvent, ok2 := event.Event.(domainOutput.OutputEvent); ok2 {
		return OutputEventDto{
			TaskID:    outputEvent.TaskID,
			Timestamp: outputEvent.Timestamp.Format(time.RFC3339Nano),
			Level:     string(outputEvent.Level),
			Source:    string(outputEvent.Source),
			Message:   outputEvent.Message,
		}
	}
	slog.Warn("unexpected output event type", "event", event)
	return nil
}
