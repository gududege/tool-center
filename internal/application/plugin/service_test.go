package plugin

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/cli-tool-center/tool-center/internal/domain/events"
	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
	"github.com/cli-tool-center/tool-center/internal/domain/settings"
)

// --- Stubs ---

type stubPluginRepo struct {
	plugins map[string]*plugin.Plugin
}

func (s *stubPluginRepo) Get(id string) (*plugin.Plugin, error) {
	p, ok := s.plugins[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}
func (s *stubPluginRepo) List() ([]*plugin.Plugin, error) {
	result := make([]*plugin.Plugin, 0, len(s.plugins))
	for _, p := range s.plugins {
		result = append(result, p)
	}
	return result, nil
}
func (s *stubPluginRepo) Save(p *plugin.Plugin) error {
	s.plugins[p.Metadata.ID] = p
	return nil
}
func (s *stubPluginRepo) Delete(id string) error {
	delete(s.plugins, id)
	return nil
}

type stubLoader struct {
	plugins []*plugin.Plugin
	err     error
}

func (s *stubLoader) Load(dir string) ([]*plugin.Plugin, error) {
	return s.plugins, s.err
}

type testBus struct {
	published []events.DomainEvent
}

func (b *testBus) Publish(event events.DomainEvent) {
	b.published = append(b.published, event)
}
func (b *testBus) Subscribe(topic string, handler events.Handler) func() { return func() {} }
func (b *testBus) PublishAsync(event events.DomainEvent)                 { b.Publish(event) }

// --- Tests ---

func TestService_LoadPlugins(t *testing.T) {
	repo := &stubPluginRepo{plugins: make(map[string]*plugin.Plugin)}
	bus := &testBus{}

	svc := NewService(repo, &stubLoader{
		plugins: []*plugin.Plugin{
			{Metadata: plugin.PluginMetadata{ID: "p1", Name: "Plugin 1"}},
			{Metadata: plugin.PluginMetadata{ID: "p2", Name: "Plugin 2"}},
		},
	}, bus, settings.DefaultSettings())

	if err := svc.LoadPlugins(); err != nil {
		t.Fatalf("LoadPlugins failed: %v", err)
	}

	plugins, _ := svc.List()
	if len(plugins) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(plugins))
	}

	// 验证事件
	found := false
	for _, e := range bus.published {
		if e.Topic() == "plugin.loaded" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected plugin.loaded events")
	}
}

func TestService_GetPlugin(t *testing.T) {
	repo := &stubPluginRepo{plugins: map[string]*plugin.Plugin{
		"p1": {Metadata: plugin.PluginMetadata{ID: "p1", Name: "Plugin 1"}},
	}}
	svc := NewService(repo, nil, &testBus{}, settings.DefaultSettings())

	p, err := svc.Get("p1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if p.Metadata.Name != "Plugin 1" {
		t.Errorf("Name = %s, want 'Plugin 1'", p.Metadata.Name)
	}
}

func TestService_GetPlugin_NotFound(t *testing.T) {
	repo := &stubPluginRepo{plugins: make(map[string]*plugin.Plugin)}
	svc := NewService(repo, nil, &testBus{}, settings.DefaultSettings())

	_, err := svc.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent plugin")
	}
}

func TestService_ListPlugins(t *testing.T) {
	repo := &stubPluginRepo{plugins: map[string]*plugin.Plugin{
		"p1": {Metadata: plugin.PluginMetadata{ID: "p1", Name: "Plugin 1"}},
		"p2": {Metadata: plugin.PluginMetadata{ID: "p2", Name: "Plugin 2"}},
		"p3": {Metadata: plugin.PluginMetadata{ID: "p3", Name: "Plugin 3"}},
	}}
	svc := NewService(repo, nil, &testBus{}, settings.DefaultSettings())

	plugins, err := svc.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(plugins) != 3 {
		t.Errorf("expected 3 plugins, got %d", len(plugins))
	}
}

func TestService_Reload(t *testing.T) {
	repo := &stubPluginRepo{plugins: make(map[string]*plugin.Plugin)}
	bus := &testBus{}

	svc := NewService(repo, &stubLoader{
		plugins: []*plugin.Plugin{
			{Metadata: plugin.PluginMetadata{ID: "p1", Name: "Plugin 1"}},
		},
	}, bus, settings.DefaultSettings())

	// 首次加载
	svc.LoadPlugins()

	// 重新加载应清空已有插件并重新加载
	count, err := svc.Reload()
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected reload count 1, got %d", count)
	}

	plugins, _ := svc.List()
	if len(plugins) != 1 {
		t.Errorf("expected 1 plugin after reload, got %d", len(plugins))
	}
}

func TestService_LoadPlugins_LoaderError(t *testing.T) {
	repo := &stubPluginRepo{plugins: make(map[string]*plugin.Plugin)}
	svc := NewService(repo, &stubLoader{
		err: errors.New("load error"),
	}, &testBus{}, settings.DefaultSettings())

	err := svc.LoadPlugins()
	if err == nil {
		t.Error("expected error when loader fails")
	}
}

func TestService_EmptyList(t *testing.T) {
	repo := &stubPluginRepo{plugins: make(map[string]*plugin.Plugin)}
	svc := NewService(repo, nil, &testBus{}, settings.DefaultSettings())

	plugins, _ := svc.List()
	if len(plugins) != 0 {
		t.Errorf("expected empty list, got %d", len(plugins))
	}
}

func TestService_SaveAndLoadHistory(t *testing.T) {
	dir := t.TempDir()
	repo := &stubPluginRepo{plugins: map[string]*plugin.Plugin{
		"p1": {
			Metadata:  plugin.PluginMetadata{ID: "p1", Name: "Plugin 1"},
			Directory: dir,
			Execution: plugin.ExecutionDefinition{Executable: filepath.Join(dir, "tool.exe")},
		},
	}}
	svc := NewService(repo, nil, &testBus{}, settings.DefaultSettings())

	for i := 1; i <= 3; i++ {
		if err := svc.SaveHistory("p1", map[string]any{"n": i}); err != nil {
			t.Fatalf("SaveHistory failed: %v", err)
		}
	}

	entries, err := svc.LoadHistory("p1")
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}

	// 最新记录在最前（JSON 数字反序列化为 float64）
	if entries[0].FormData["n"].(float64) != 3 {
		t.Errorf("expected newest entry n=3, got %v", entries[0].FormData["n"])
	}

	// 验证文件已写入
	path := filepath.Join(dir, "history.json")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("history file not created: %v", err)
	}
}

func TestService_SaveHistory_RollingMax5(t *testing.T) {
	dir := t.TempDir()
	repo := &stubPluginRepo{plugins: map[string]*plugin.Plugin{
		"p1": {
			Metadata:  plugin.PluginMetadata{ID: "p1", Name: "Plugin 1"},
			Directory: dir,
			Execution: plugin.ExecutionDefinition{Executable: filepath.Join(dir, "tool.exe")},
		},
	}}
	svc := NewService(repo, nil, &testBus{}, settings.DefaultSettings())

	for i := 1; i <= 7; i++ {
		if err := svc.SaveHistory("p1", map[string]any{"n": i}); err != nil {
			t.Fatalf("SaveHistory failed: %v", err)
		}
	}

	entries, err := svc.LoadHistory("p1")
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}
	if len(entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(entries))
	}
	if entries[0].FormData["n"].(float64) != 7 {
		t.Errorf("expected newest entry n=7, got %v", entries[0].FormData["n"])
	}
	if entries[4].FormData["n"].(float64) != 3 {
		t.Errorf("expected oldest kept entry n=3, got %v", entries[4].FormData["n"])
	}
}

func TestService_SaveHistory_PluginNotFound(t *testing.T) {
	repo := &stubPluginRepo{plugins: make(map[string]*plugin.Plugin)}
	svc := NewService(repo, nil, &testBus{}, settings.DefaultSettings())

	err := svc.SaveHistory("missing", map[string]any{})
	if err == nil {
		t.Error("expected error when plugin not found")
	}
}
