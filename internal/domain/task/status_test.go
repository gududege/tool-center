package task

import (
	"testing"
)

func TestTaskStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		status TaskStatus
		want   bool
	}{
		{Created, false},
		{Queued, false},
		{Running, false},
		{Cancelling, false},
		{Stopping, false},
		{Completed, true},
		{Failed, true},
		{Cancelled, true},
		{Stopped, true},
	}

	for _, tt := range tests {
		got := tt.status.IsTerminal()
		if got != tt.want {
			t.Errorf("TaskStatus(%s).IsTerminal() = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestTaskStatus_ValidTransitions(t *testing.T) {
	tests := []struct {
		from TaskStatus
		to   TaskStatus
		want bool
	}{
		// 合法转换
		{Created, Queued, true},
		{Queued, Running, true},
		{Running, Completed, true},
		{Running, Failed, true},
		{Running, Cancelling, true},
		{Cancelling, Cancelled, true},
		{Running, Stopping, true},
		{Stopping, Stopped, true},

		// 非法转换
		{Created, Running, false},
		{Created, Completed, false},
		{Queued, Completed, false},
		{Completed, Running, false},
		{Failed, Running, false},
		{Cancelled, Running, false},
		{Stopped, Running, false},
		{Completed, Failed, false},

		// 相同状态
		{Created, Created, false},
		{Running, Running, false},
	}

	for _, tt := range tests {
		valid := false
		for _, s := range tt.from.ValidTransitions() {
			if s == tt.to {
				valid = true
				break
			}
		}
		if valid != tt.want {
			t.Errorf("transition %s -> %s: got %v, want %v", tt.from, tt.to, valid, tt.want)
		}
	}
}

func TestTask_CanTransitionTo(t *testing.T) {
	t.Run("successful transition", func(t *testing.T) {
		task := &Task{Status: Created}
		if !task.CanTransitionTo(Queued) {
			t.Error("expected Created -> Queued to be valid")
		}
	})

	t.Run("invalid transition returns false", func(t *testing.T) {
		task := &Task{Status: Created}
		if task.CanTransitionTo(Running) {
			t.Error("expected Created -> Running to be invalid")
		}
	})

	t.Run("completed task cannot transition", func(t *testing.T) {
		task := &Task{Status: Completed}
		if task.CanTransitionTo(Running) {
			t.Error("expected Completed -> Running to be invalid")
		}
	})

	t.Run("cancelled task cannot transition", func(t *testing.T) {
		task := &Task{Status: Cancelled}
		if task.CanTransitionTo(Running) {
			t.Error("expected Cancelled -> Running to be invalid")
		}
	})
}
