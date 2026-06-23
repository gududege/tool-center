package task

import (
	"testing"
)

func TestTaskStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		status TaskStatus
		want   bool
	}{
		{TaskCreated, false},
		{TaskQueued, false},
		{TaskRunning, false},
		{TaskCompleted, true},
		{TaskFailed, true},
		{TaskCancelled, true},
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
		{TaskCreated, TaskQueued, true},
		{TaskQueued, TaskRunning, true},
		{TaskRunning, TaskCompleted, true},
		{TaskRunning, TaskFailed, true},
		{TaskRunning, TaskCancelled, true},

		// 非法转换
		{TaskCreated, TaskRunning, false},
		{TaskCreated, TaskCompleted, false},
		{TaskQueued, TaskCompleted, false},
		{TaskCompleted, TaskRunning, false},
		{TaskFailed, TaskRunning, false},
		{TaskCancelled, TaskRunning, false},
		{TaskCompleted, TaskFailed, false},

		// 相同状态
		{TaskCreated, TaskCreated, false},
		{TaskRunning, TaskRunning, false},
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
		task := &Task{Status: TaskCreated}
		if !task.CanTransitionTo(TaskQueued) {
			t.Error("expected Created -> Queued to be valid")
		}
	})

	t.Run("invalid transition returns false", func(t *testing.T) {
		task := &Task{Status: TaskCreated}
		if task.CanTransitionTo(TaskRunning) {
			t.Error("expected Created -> Running to be invalid")
		}
	})

	t.Run("completed task cannot transition", func(t *testing.T) {
		task := &Task{Status: TaskCompleted}
		if task.CanTransitionTo(TaskRunning) {
			t.Error("expected Completed -> Running to be invalid")
		}
	})

	t.Run("cancelled task cannot transition", func(t *testing.T) {
		task := &Task{Status: TaskCancelled}
		if task.CanTransitionTo(TaskRunning) {
			t.Error("expected Cancelled -> Running to be invalid")
		}
	})
}
