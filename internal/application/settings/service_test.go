package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cli-tool-center/tool-center/internal/domain/settings"
)

func TestService_GetSettings(t *testing.T) {
	repo := NewRepository(filepath.Join(t.TempDir(), "settings.json"))
	svc := NewService(repo)

	s, err := svc.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings failed: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil settings")
	}
	if s.Theme != "dark" {
		t.Errorf("Theme = %s, want dark", s.Theme)
	}
}

func TestService_SaveSettings(t *testing.T) {
	repo := NewRepository(filepath.Join(t.TempDir(), "settings.json"))
	svc := NewService(repo)

	s := &settings.Settings{
		Theme:           "light",
		PluginDirectory: "./custom-plugins",
		MaxOutputLines:  5000,
		MaxTaskHistory:  50,
	}

	if err := svc.SaveSettings(s); err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	got, _ := svc.GetSettings()
	if got.Theme != "light" {
		t.Errorf("Theme = %s, want light", got.Theme)
	}
	if got.PluginDirectory != "./custom-plugins" {
		t.Errorf("PluginDirectory = %s, want ./custom-plugins", got.PluginDirectory)
	}
}

func TestService_RoundTrip(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "settings.json")
	repo := NewRepository(tmpFile)
	svc := NewService(repo)

	s1, _ := svc.GetSettings()
	s1.Theme = "light"
	s1.MaxOutputLines = 2000
	svc.SaveSettings(s1)

	s2, _ := svc.GetSettings()
	if s2.Theme != "light" {
		t.Errorf("Theme = %s, want light", s2.Theme)
	}
	if s2.MaxOutputLines != 2000 {
		t.Errorf("MaxOutputLines = %d, want 2000", s2.MaxOutputLines)
	}
}

func TestService_FilePersistence(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "settings.json")

	// 第一次：保存
	repo1 := NewRepository(tmpFile)
	svc1 := NewService(repo1)
	svc1.SaveSettings(&settings.Settings{
		Theme:           "light",
		PluginDirectory: "./my-plugins",
	})
	data1, _ := os.ReadFile(tmpFile)
	if len(data1) == 0 {
		t.Fatal("settings file should not be empty after save")
	}

	// 第二次：新建实例加载
	repo2 := NewRepository(tmpFile)
	repo2.Load()
	svc2 := NewService(repo2)
	loaded, _ := svc2.GetSettings()
	if loaded.Theme != "light" {
		t.Errorf("Theme = %s, want light", loaded.Theme)
	}
	if loaded.PluginDirectory != "./my-plugins" {
		t.Errorf("PluginDirectory = %s, want ./my-plugins", loaded.PluginDirectory)
	}
}
