package events

import (
	"testing"
)

func TestPluginEvents_Topics(t *testing.T) {
	tests := []struct {
		event DomainEvent
		topic string
	}{
		{PluginLoaded{PluginID: "test"}, "plugin.loaded"},
		{PluginUnloaded{PluginID: "test"}, "plugin.unloaded"},
		{PluginReloaded{PluginID: "test"}, "plugin.reloaded"},
	}

	for _, tt := range tests {
		if got := tt.event.Topic(); got != tt.topic {
			t.Errorf("event.Topic() = %s, want %s", got, tt.topic)
		}
	}
}

func TestTaskEvents_Topics(t *testing.T) {
	tests := []struct {
		event DomainEvent
		topic string
	}{
		{TaskCreated{TaskID: "t1", PluginID: "p1"}, "task.created"},
		{TaskStarted{TaskID: "t1"}, "task.started"},
		{TaskCompleted{TaskID: "t1", ExitCode: 0}, "task.completed"},
		{TaskFailed{TaskID: "t1", Error: "err"}, "task.failed"},
		{TaskCancelled{TaskID: "t1"}, "task.cancelled"},
		{TaskStatusChanged{TaskID: "t1", OldStatus: "running", NewStatus: "completed"}, "task.status-changed"},
	}

	for _, tt := range tests {
		if got := tt.event.Topic(); got != tt.topic {
			t.Errorf("%T.Topic() = %s, want %s", tt.event, got, tt.topic)
		}
	}
}

func TestOutputEvent_Topic(t *testing.T) {
	e := OutputReceived{TaskID: "t1"}
	if got := e.Topic(); got != "output.received" {
		t.Errorf("OutputReceived.Topic() = %s, want %s", got, "output.received")
	}
}

func TestPluginLoaded_Fields(t *testing.T) {
	e := PluginLoaded{PluginID: "tia-export"}
	if e.PluginID != "tia-export" {
		t.Errorf("PluginLoaded.PluginID = %s, want tia-export", e.PluginID)
	}
}

func TestTaskFailed_Fields(t *testing.T) {
	e := TaskFailed{TaskID: "task-1", Error: "something went wrong"}
	if e.TaskID != "task-1" || e.Error != "something went wrong" {
		t.Errorf("TaskFailed fields mismatch: %+v", e)
	}
}
