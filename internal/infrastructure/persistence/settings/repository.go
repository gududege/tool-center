package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	domainSettings "github.com/cli-tool-center/tool-center/internal/domain/settings"
)

// Repository 基于文件的配置仓储实现
type Repository struct {
	mu       sync.RWMutex
	filePath string
	settings *domainSettings.Settings
}

// NewRepository 创建配置仓储
func NewRepository(filePath string) *Repository {
	return &Repository{
		filePath: filePath,
		settings: domainSettings.DefaultSettings(),
	}
}

// Load 从文件加载配置
func (r *Repository) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read settings file: %w", err)
	}

	var s domainSettings.Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("parse settings: %w", err)
	}
	r.settings = &s
	return nil
}

// Get 获取当前配置（返回副本）
func (r *Repository) Get() (*domainSettings.Settings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s := *r.settings
	return &s, nil
}

// Save 保存配置
func (r *Repository) Save(s *domainSettings.Settings) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create settings dir: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("write settings file: %w", err)
	}

	copied := *s
	r.settings = &copied
	return nil
}
