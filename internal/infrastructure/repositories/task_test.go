package repositories

import (
	"testing"

	"github.com/cli-tool-center/tool-center/internal/domain/task"
)

func TestMemoryTaskRepository_SaveAndGet(t *testing.T) {
	repo := NewMemoryTaskRepository()

	tsk := &task.Task{
		ID:       "task-1",
		PluginID: "test-plugin",
		Status:   task.TaskCreated,
	}

	if err := repo.Save(tsk); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.Get("task-1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.ID != "task-1" || got.PluginID != "test-plugin" {
		t.Errorf("got Task = %+v, want ID=task-1 PluginID=test-plugin", got)
	}
}

func TestMemoryTaskRepository_GetNotFound(t *testing.T) {
	repo := NewMemoryTaskRepository()
	_, err := repo.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestMemoryTaskRepository_List(t *testing.T) {
	repo := NewMemoryTaskRepository()

	for i := 0; i < 5; i++ {
		repo.Save(&task.Task{
			ID:       string(rune('A' + i)),
			PluginID: "p1",
			Status:   task.TaskCreated,
		})
	}

	tasks, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(tasks) != 5 {
		t.Errorf("expected 5 tasks, got %d", len(tasks))
	}
}

func TestMemoryTaskRepository_Delete(t *testing.T) {
	repo := NewMemoryTaskRepository()

	repo.Save(&task.Task{ID: "to-delete", PluginID: "p1"})
	if err := repo.Delete("to-delete"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.Get("to-delete")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestMemoryTaskRepository_Update(t *testing.T) {
	repo := NewMemoryTaskRepository()

	repo.Save(&task.Task{ID: "t1", PluginID: "p1", Status: task.TaskCreated})

	// 更新
	repo.Save(&task.Task{ID: "t1", PluginID: "p1", Status: task.TaskRunning})

	got, _ := repo.Get("t1")
	if got.Status != task.TaskRunning {
		t.Errorf("expected status 'running', got '%s'", got.Status)
	}
}
