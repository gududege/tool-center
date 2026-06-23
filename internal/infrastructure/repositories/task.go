package repositories

import (
	"fmt"
	"sync"

	"github.com/cli-tool-center/tool-center/internal/domain/task"
)

// MemoryTaskRepository 基于内存的任务仓储实现
type MemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*task.Task
}

// NewMemoryTaskRepository 创建内存任务仓储
func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make(map[string]*task.Task),
	}
}

// Get 根据ID获取任务
func (r *MemoryTaskRepository) Get(id string) (*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task %s not found", id)
	}
	return t, nil
}

// List 获取所有任务
func (r *MemoryTaskRepository) List() ([]*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*task.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result, nil
}

// Save 保存任务
func (r *MemoryTaskRepository) Save(t *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[t.ID] = t
	return nil
}

// Delete 删除任务
func (r *MemoryTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tasks, id)
	return nil
}
