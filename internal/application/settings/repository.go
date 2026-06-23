package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cli-tool-center/tool-center/internal/domain/settings"
)

// Repository 配置仓储
type Repository struct {
	mu       sync.RWMutex
	filePath string
	settings *settings.Settings
}

// NewRepository 创建配置仓储
func NewRepository(filePath string) *Repository {
	return &Repository{
		filePath: filePath,
		settings: settings.DefaultSettings(),
	}
}

// Load 从文件加载配置
func (r *Repository) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，使用默认配置
			return nil
		}
		return fmt.Errorf("read settings file: %w", err)
	}

	var s settings.Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("parse settings: %w", err)
	}
	r.settings = &s
	return nil
}

// Get 获取当前配置
func (r *Repository) Get() (*settings.Settings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 返回副本
	s := *r.settings
	return &s, nil
}

// Save 保存配置
func (r *Repository) Save(s *settings.Settings) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 确保目录存在
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

	r.settings = s
	return nil
}
