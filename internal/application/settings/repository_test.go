package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cli-tool-center/tool-center/internal/domain/settings"
)

func TestRepository_DefaultSettings(t *testing.T) {
	// 文件不存在时使用默认值
	tmpFile := filepath.Join(t.TempDir(), "settings.json")
	repo := NewRepository(tmpFile)

	if err := repo.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	s, err := repo.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if s.Theme != "dark" {
		t.Errorf("expected default theme 'dark', got '%s'", s.Theme)
	}
	if s.PluginDirectory != "./plugins" {
		t.Errorf("expected default plugin directory './plugins', got '%s'", s.PluginDirectory)
	}
}

func TestRepository_SaveAndLoad(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "settings.json")
	repo := NewRepository(tmpFile)

	s := &settings.Settings{
		Theme:           "light",
		PluginDirectory: "./my-plugins",
		MaxOutputLines:  5000,
		MaxTaskHistory:  50,
	}

	if err := repo.Save(s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 新的 repository 实例应该能正确加载
	repo2 := NewRepository(tmpFile)
	if err := repo2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	loaded, err := repo2.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if loaded.Theme != "light" {
		t.Errorf("Theme = %s, want light", loaded.Theme)
	}
	if loaded.PluginDirectory != "./my-plugins" {
		t.Errorf("PluginDirectory = %s, want ./my-plugins", loaded.PluginDirectory)
	}
	if loaded.MaxOutputLines != 5000 {
		t.Errorf("MaxOutputLines = %d, want 5000", loaded.MaxOutputLines)
	}
}

func TestRepository_Persistence(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "settings.json")
	repo := NewRepository(tmpFile)

	s := &settings.Settings{
		Theme:           "dark",
		PluginDirectory: "./plugins",
		MaxOutputLines:  10000,
		MaxTaskHistory:  100,
		AutoReloadPlugins: true,
	}
	repo.Save(s)

	// 文件应该实际被写入
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("settings file is empty")
	}
}

func TestRepository_ConcurrentAccess(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "settings.json")
	repo := NewRepository(tmpFile)

	// 并发读写
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			s, _ := repo.Get()
			s.MaxTaskHistory++
			repo.Save(s)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
