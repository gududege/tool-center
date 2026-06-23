package filesystem

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoader_Load_ValidPlugin(t *testing.T) {
	dir := t.TempDir()

	// 创建插件目录
	pluginDir := filepath.Join(dir, "test-plugin")
	mkDir(t, pluginDir)

	// 创建一个假的 exe
	touchFile(t, filepath.Join(pluginDir, "tool.exe"))

	// 写入 plugin.json
	writeJSON(t, filepath.Join(pluginDir, "plugin.json"), map[string]any{
		"schemaVersion": "1.0",
		"id":            "test-plugin",
		"name":          "Test Plugin",
		"description":   "A test plugin",
		"version":       "1.0.0",
		"author":        "Tester",
		"navigation": map[string]any{
			"group": []string{"Tools", "Test"},
			"order": 100,
		},
		"execution": map[string]any{
			"exe": "./tool.exe",
		},
	})

	loader := NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(plugins))
	}

	p := plugins[0]
	if p.Metadata.ID != "test-plugin" {
		t.Errorf("ID = %s, want test-plugin", p.Metadata.ID)
	}
	if p.Metadata.Name != "Test Plugin" {
		t.Errorf("Name = %s, want Test Plugin", p.Metadata.Name)
	}
	if p.Metadata.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", p.Metadata.Version)
	}
	if p.Execution.Executable != filepath.Join(pluginDir, "tool.exe") {
		t.Errorf("Executable = %s, want %s", p.Execution.Executable, filepath.Join(pluginDir, "tool.exe"))
	}
	if len(p.Navigation.Group) != 2 || p.Navigation.Group[0] != "Tools" || p.Navigation.Group[1] != "Test" {
		t.Errorf("Navigation.Group = %v, want [Tools Test]", p.Navigation.Group)
	}
	if p.Navigation.Order != 100 {
		t.Errorf("Navigation.Order = %d, want 100", p.Navigation.Order)
	}
}

func TestLoader_Load_WithSchema(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "schema-plugin")
	mkDir(t, pluginDir)
	touchFile(t, filepath.Join(pluginDir, "tool.exe"))

	// 写入 schema.json
	writeJSON(t, filepath.Join(pluginDir, "schema.json"), map[string]any{
		"type":       "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	})

	// 写入 uischema.json
	writeJSON(t, filepath.Join(pluginDir, "uischema.json"), map[string]any{
		"name": map[string]any{"ui:autofocus": true},
	})

	// 写入 plugin.json
	writeJSON(t, filepath.Join(pluginDir, "plugin.json"), map[string]any{
		"schemaVersion": "1.0",
		"id":            "schema-plugin",
		"name":          "Schema Plugin",
		"navigation":    map[string]any{"group": []string{"Tools"}, "order": 1},
		"form": map[string]any{
			"schema":   "./schema.json",
			"uischema": "./uischema.json",
		},
		"execution": map[string]any{
			"exe": "./tool.exe",
		},
	})

	loader := NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(plugins))
	}

	p := plugins[0]
	if len(p.Form.Schema) == 0 {
		t.Fatal("expected Schema to be loaded")
	}

	// RawMessage 透传原始字节，解析后 type 应为 object
	var schema map[string]any
	if err := json.Unmarshal(p.Form.Schema, &schema); err != nil {
		t.Fatalf("schema is not valid JSON: %v", err)
	}
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want 'object'", schema["type"])
	}

	// 关键：透传 RawMessage 后，序列化 DTO 时 properties 的字段顺序必须与源文件一致，
	// 而不是按字母序重排（map[string]any + json.Marshal 会字母序，导致前端表单字段顺序错乱）。
	// 用 json.Decoder 的 Token 流直接从原始字节里提取 properties 下的字段顺序（不经 map）。
	if order := propertyOrder(string(p.Form.Schema)); order == nil {
		t.Fatal("could not parse property order from schema bytes")
	} else {
		want := []string{"name"}
		if len(order) != len(want) || order[0] != want[0] {
			t.Errorf("property order = %v, want %v", order, want)
		}
	}

	if len(p.Form.UISchema) == 0 {
		t.Fatal("expected UISchema to be loaded")
	}
}

// propertyOrder 从 schema JSON 字节里按出现顺序提取 properties 下的字段名，
// 不经过 map（避免字母序），用于验证 RawMessage 透传保留了文件定义顺序。
func propertyOrder(schemaJSON string) []string {
	dec := json.NewDecoder(strings.NewReader(schemaJSON))
	tok, err := dec.Token()
	if err != nil {
		return nil
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil
	}
	for {
		key, err := dec.Token()
		if err != nil {
			return nil
		}
		if d, ok := key.(json.Delim); ok && d == '}' {
			return nil
		}
		keyStr, _ := key.(string)
		if keyStr == "properties" {
			// 进入 properties 对象
			t, err := dec.Token()
			if err != nil {
				return nil
			}
			d, ok := t.(json.Delim)
			if !ok || d != '{' {
				return nil
			}
			var order []string
			for {
				name, err := dec.Token()
				if err != nil {
					return nil
				}
				if d, ok := name.(json.Delim); ok && d == '}' {
					return order
				}
				ns, _ := name.(string)
				order = append(order, ns)
				// 跳过该字段的值
				var v any
				if err := dec.Decode(&v); err != nil {
					return nil
				}
			}
		}
		// 跳过非 properties 字段的值
		var v any
		if err := dec.Decode(&v); err != nil {
			return nil
		}
	}
}

func TestLoader_Load_WithEnvironment(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "env-plugin")
	mkDir(t, pluginDir)
	touchFile(t, filepath.Join(pluginDir, "tool.exe"))

	writeJSON(t, filepath.Join(pluginDir, "plugin.json"), map[string]any{
		"schemaVersion": "1.0",
		"id":            "env-plugin",
		"name":          "Env Plugin",
		"navigation":    map[string]any{"group": []string{"Tools"}, "order": 1},
		"execution": map[string]any{
			"exe": "./tool.exe",
			"environment": map[string]string{
				"LOG_LEVEL": "debug",
				"MODE":      "test",
			},
		},
	})

	loader := NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := plugins[0]
	if p.Execution.Environment["LOG_LEVEL"] != "debug" {
		t.Errorf("LOG_LEVEL = %s, want debug", p.Execution.Environment["LOG_LEVEL"])
	}
	if p.Execution.Environment["MODE"] != "test" {
		t.Errorf("MODE = %s, want test", p.Execution.Environment["MODE"])
	}
}

func TestLoader_Load_WithParameters(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "param-plugin")
	mkDir(t, pluginDir)
	touchFile(t, filepath.Join(pluginDir, "tool.exe"))

	writeJSON(t, filepath.Join(pluginDir, "plugin.json"), map[string]any{
		"schemaVersion": "1.0",
		"id":            "param-plugin",
		"name":          "Param Plugin",
		"navigation":    map[string]any{"group": []string{"Tools"}, "order": 1},
		"execution": map[string]any{
			"exe": "./tool.exe",
			"parameters": []map[string]any{
				{"field": "project", "kind": "option", "flag": "--project"},
				{"field": "overwrite", "kind": "switch", "flag": "--overwrite"},
				{"field": "files", "kind": "argument-array"},
			},
		},
	})

	loader := NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := plugins[0]
	if len(p.Execution.Parameters) != 3 {
		t.Fatalf("expected 3 parameters, got %d", len(p.Execution.Parameters))
	}
	if p.Execution.Parameters[0].Field != "project" || p.Execution.Parameters[0].Flag != "--project" {
		t.Errorf("param[0] = %+v, want field=project flag=--project", p.Execution.Parameters[0])
	}
	if p.Execution.Parameters[1].Field != "overwrite" {
		t.Errorf("param[1].Field = %s, want overwrite", p.Execution.Parameters[1].Field)
	}
	if p.Execution.Parameters[2].Kind != "argument-array" {
		t.Errorf("param[2].Kind = %s, want argument-array", p.Execution.Parameters[2].Kind)
	}
}

func TestLoader_Load_MissingPluginJSON(t *testing.T) {
	loader := NewLoader()
	// 空目录
	dir := t.TempDir()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(plugins))
	}
}

func TestLoader_Load_InvalidPluginJSON(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "broken")
	mkDir(t, pluginDir)
	touchFile(t, filepath.Join(pluginDir, "tool.exe"))

	// 写入无效 JSON
	if err := os.WriteFile(filepath.Join(pluginDir, "plugin.json"), []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("should skip broken plugin, got error: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins after skip, got %d", len(plugins))
	}
}

func TestLoader_Load_MissingRequiredFields(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "incomplete")
	mkDir(t, pluginDir)
	touchFile(t, filepath.Join(pluginDir, "tool.exe"))

	// 缺少 id
	writeJSON(t, filepath.Join(pluginDir, "plugin.json"), map[string]any{
		"schemaVersion": "1.0",
		"name":          "No ID",
		"navigation":    map[string]any{"group": []string{"Tools"}, "order": 1},
		"execution":     map[string]any{"exe": "./tool.exe"},
	})

	loader := NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("should skip invalid plugin: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(plugins))
	}
}

func TestLoader_Load_MultiplePlugins(t *testing.T) {
	dir := t.TempDir()

	// 创建两个插件
	for _, id := range []string{"plugin-a", "plugin-b"} {
		pDir := filepath.Join(dir, id)
		mkDir(t, pDir)
		touchFile(t, filepath.Join(pDir, "tool.exe"))
		writeJSON(t, filepath.Join(pDir, "plugin.json"), map[string]any{
			"schemaVersion": "1.0",
			"id":            id,
			"name":          id,
			"navigation":    map[string]any{"group": []string{"Tools"}, "order": 1},
			"execution":     map[string]any{"exe": "./tool.exe"},
		})
	}

	loader := NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plugins) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(plugins))
	}
}

func TestLoader_Load_SkipNonDirectories(t *testing.T) {
	dir := t.TempDir()

	// 创建一个文件而不是目录
	touchFile(t, filepath.Join(dir, "file.txt"))

	pluginDir := filepath.Join(dir, "valid-plugin")
	mkDir(t, pluginDir)
	touchFile(t, filepath.Join(pluginDir, "tool.exe"))
	writeJSON(t, filepath.Join(pluginDir, "plugin.json"), map[string]any{
		"schemaVersion": "1.0",
		"id":            "valid-plugin",
		"name":          "Valid",
		"navigation":    map[string]any{"group": []string{"Tools"}, "order": 1},
		"execution":     map[string]any{"exe": "./tool.exe"},
	})

	loader := NewLoader()
	plugins, err := loader.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plugins) != 1 {
		t.Errorf("expected 1 plugin (skip file.txt), got %d", len(plugins))
	}
}

// helpers
func mkDir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}

func touchFile(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func writeJSON(t *testing.T, path string, data map[string]any) {
	t.Helper()
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		t.Fatal(err)
	}
}
