package repositories

import (
	"testing"

	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
)

func TestMemoryPluginRepository_SaveAndGet(t *testing.T) {
	repo := NewMemoryPluginRepository()

	p := &plugin.Plugin{
		Metadata: plugin.PluginMetadata{
			ID:   "test-plugin",
			Name: "Test Plugin",
		},
	}

	if err := repo.Save(p); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.Get("test-plugin")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Metadata.ID != "test-plugin" {
		t.Errorf("got ID = %s, want test-plugin", got.Metadata.ID)
	}
}

func TestMemoryPluginRepository_GetNotFound(t *testing.T) {
	repo := NewMemoryPluginRepository()
	_, err := repo.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent plugin")
	}
}

func TestMemoryPluginRepository_List(t *testing.T) {
	repo := NewMemoryPluginRepository()

	ids := []string{"a", "b", "c"}
	for _, id := range ids {
		repo.Save(&plugin.Plugin{
			Metadata: plugin.PluginMetadata{ID: id, Name: id},
		})
	}

	plugins, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(plugins) != 3 {
		t.Errorf("expected 3 plugins, got %d", len(plugins))
	}
}

func TestMemoryPluginRepository_Delete(t *testing.T) {
	repo := NewMemoryPluginRepository()

	repo.Save(&plugin.Plugin{
		Metadata: plugin.PluginMetadata{ID: "to-delete", Name: "To Delete"},
	})

	if err := repo.Delete("to-delete"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.Get("to-delete")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestMemoryPluginRepository_Concurrency(t *testing.T) {
	repo := NewMemoryPluginRepository()

	// 并发读写
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			id := string(rune('A' + n))
			repo.Save(&plugin.Plugin{
				Metadata: plugin.PluginMetadata{ID: id, Name: id},
			})
			repo.List()
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	plugins, _ := repo.List()
	if len(plugins) != 10 {
		t.Errorf("expected 10 plugins after concurrent writes, got %d", len(plugins))
	}
}
