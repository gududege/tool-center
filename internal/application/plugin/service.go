package plugin

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/cli-tool-center/tool-center/internal/domain/events"
	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
	"github.com/cli-tool-center/tool-center/internal/domain/settings"
)

//go:generate mockgen -source=service.go -destination=mock_service_test.go -package=plugin

// PluginLoader 插件加载器接口
type PluginLoader interface {
	Load(dir string) ([]*plugin.Plugin, error)
}

// ParameterBuilder 参数构建器接口
type ParameterBuilder interface {
	Build(mappings []plugin.ParameterMapping, formData map[string]any) ([]string, error)
}

// Service 插件应用服务
type Service struct {
	repo     plugin.PluginRepository
	loader   PluginLoader
	eventBus events.Bus
	cfg      *settings.Settings
}

// NewService 创建插件服务
func NewService(repo plugin.PluginRepository, loader PluginLoader, eventBus events.Bus, cfg *settings.Settings) *Service {
	return &Service{
		repo:     repo,
		loader:   loader,
		eventBus: eventBus,
		cfg:      cfg,
	}
}

// LoadPlugins 扫描并加载所有插件
func (s *Service) LoadPlugins() error {
	dir := s.cfg.PluginDirectory
	plugins, err := s.loader.Load(dir)
	if err != nil {
		return fmt.Errorf("load plugins from %s: %w", dir, err)
	}

	for _, p := range plugins {
		if err := s.repo.Save(p); err != nil {
			slog.Warn("save plugin failed", "id", p.Metadata.ID, "error", err)
			continue
		}
		s.eventBus.Publish(events.PluginLoaded{PluginID: p.Metadata.ID})
		slog.Info("plugin loaded", "id", p.Metadata.ID, "name", p.Metadata.Name)
	}

	slog.Info("plugins loaded", "count", len(plugins))
	return nil
}

// Reload 重新加载所有插件
func (s *Service) Reload() (int, error) {
	// 清空已有插件
	existing, err := s.repo.List()
	if err == nil {
		for _, p := range existing {
			_ = s.repo.Delete(p.Metadata.ID)
			s.eventBus.Publish(events.PluginUnloaded{PluginID: p.Metadata.ID})
		}
	}

	if err := s.LoadPlugins(); err != nil {
		return 0, err
	}

	all, _ := s.repo.List()
	for _, p := range all {
		s.eventBus.Publish(events.PluginReloaded{PluginID: p.Metadata.ID})
	}

	return len(all), nil
}

// Get 获取单个插件
func (s *Service) Get(id string) (*plugin.Plugin, error) {
	return s.repo.Get(id)
}

// List 获取所有插件列表
func (s *Service) List() ([]*plugin.Plugin, error) {
	return s.repo.List()
}

// HistoryEntry 插件参数历史记录条目。
type HistoryEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	Label     string         `json:"label"`
	FormData  map[string]any `json:"formData"`
}

const historyFileName = "history.json"
const maxHistoryEntries = 5

// SaveHistory 把当前表单参数保存到插件目录的 history.json，最多保留 5 条。
func (s *Service) SaveHistory(pluginID string, formData map[string]any) error {
	p, err := s.repo.Get(pluginID)
	if err != nil {
		return fmt.Errorf("get plugin %s: %w", pluginID, err)
	}

	entries, err := s.readHistory(p.Directory)
	if err != nil {
		return err
	}

	label := fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), filepath.Base(p.Execution.Executable))
	newEntry := HistoryEntry{
		Timestamp: time.Now(),
		Label:     label,
		FormData:  formData,
	}

	// prepend 新记录
	entries = append([]HistoryEntry{newEntry}, entries...)
	if len(entries) > maxHistoryEntries {
		entries = entries[:maxHistoryEntries]
	}

	return s.writeHistory(p.Directory, entries)
}

// LoadHistory 读取插件目录的 history.json。
func (s *Service) LoadHistory(pluginID string) ([]HistoryEntry, error) {
	p, err := s.repo.Get(pluginID)
	if err != nil {
		return nil, fmt.Errorf("get plugin %s: %w", pluginID, err)
	}
	return s.readHistory(p.Directory)
}

func (s *Service) readHistory(dir string) ([]HistoryEntry, error) {
	path := filepath.Join(dir, historyFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []HistoryEntry{}, nil
		}
		return nil, fmt.Errorf("read history file %s: %w", path, err)
	}

	var entries []HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parse history file %s: %w", path, err)
	}
	return entries, nil
}

func (s *Service) writeHistory(dir string, entries []HistoryEntry) error {
	path := filepath.Join(dir, historyFileName)
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write history file %s: %w", path, err)
	}
	return nil
}
