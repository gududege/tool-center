package repositories

import (
	"fmt"
	"sync"

	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
)

// MemoryPluginRepository 基于内存的插件仓储实现
type MemoryPluginRepository struct {
	mu      sync.RWMutex
	plugins map[string]*plugin.Plugin
}

// NewMemoryPluginRepository 创建内存插件仓储
func NewMemoryPluginRepository() *MemoryPluginRepository {
	return &MemoryPluginRepository{
		plugins: make(map[string]*plugin.Plugin),
	}
}

// Get 根据ID获取插件
func (r *MemoryPluginRepository) Get(id string) (*plugin.Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.plugins[id]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", id)
	}
	return p, nil
}

// List 获取所有插件
func (r *MemoryPluginRepository) List() ([]*plugin.Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*plugin.Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		result = append(result, p)
	}
	return result, nil
}

// Save 保存插件
func (r *MemoryPluginRepository) Save(p *plugin.Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.plugins[p.Metadata.ID] = p
	return nil
}

// Delete 删除插件
func (r *MemoryPluginRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.plugins, id)
	return nil
}
