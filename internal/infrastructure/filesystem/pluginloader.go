package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
)

// PluginConfig plugin.json 根结构
type PluginConfig struct {
	SchemaVersion string            `json:"schemaVersion"`
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	NameCn        string            `json:"name_cn,omitempty"`
	Description   string            `json:"description,omitempty"`
	DescriptionCn string            `json:"description_cn,omitempty"`
	Version       string            `json:"version,omitempty"`
	Author        string            `json:"author,omitempty"`
	Homepage      string            `json:"homepage,omitempty"`
	Icon          string            `json:"icon,omitempty"`
	Navigation    *NavigationConfig `json:"navigation"`
	Form          *FormConfig       `json:"form,omitempty"`
	Execution     *ExecutionConfig  `json:"execution"`
}

type NavigationConfig struct {
	Group   []string `json:"group"`
	GroupCn []string `json:"group_cn,omitempty"`
	Order   int      `json:"order"`
}

type FormConfig struct {
	Schema   string `json:"schema"`
	UISchema string `json:"uischema,omitempty"`
}

type ExecutionConfig struct {
	Exe              string            `json:"exe"`
	WorkingDirectory string            `json:"workingDirectory,omitempty"`
	Environment      map[string]string `json:"environment,omitempty"`
	Parameters       []ParameterConfig `json:"parameters,omitempty"`
}

type ParameterConfig struct {
	Field        string `json:"field"`
	Kind         string `json:"kind"`
	Flag         string `json:"flag,omitempty"`
	Style        string `json:"style,omitempty"`
	Separator    string `json:"separator,omitempty"`
	TrueFlag     string `json:"trueFlag,omitempty"`
	FalseFlag    string `json:"falseFlag,omitempty"`
	DefaultValue any    `json:"defaultValue,omitempty"`
}

// Loader 文件系统插件加载器
type Loader struct{}

// NewLoader 创建插件加载器
func NewLoader() *Loader {
	return &Loader{}
}

// Load 扫描并加载 plugins 目录下的所有插件
func (l *Loader) Load(dir string) ([]*plugin.Plugin, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read plugins dir %s: %w", dir, err)
	}

	var plugins []*plugin.Plugin
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		p, err := l.loadPlugin(filepath.Join(dir, entry.Name()))
		if err != nil {
			// 单个插件加载失败不影响其他插件
			fmt.Fprintf(os.Stderr, "warn: skip plugin %s: %v\n", entry.Name(), err)
			continue
		}
		plugins = append(plugins, p)
	}
	return plugins, nil
}

// loadPlugin 加载单个插件
func (l *Loader) loadPlugin(pluginDir string) (*plugin.Plugin, error) {
	configPath := filepath.Join(pluginDir, "plugin.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read plugin.json: %w", err)
	}

	var cfg PluginConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse plugin.json: %w", err)
	}

	if cfg.ID == "" || cfg.Name == "" {
		return nil, fmt.Errorf("plugin id and name are required")
	}

	p := &plugin.Plugin{
		Metadata: plugin.PluginMetadata{
			ID:            cfg.ID,
			Name:          cfg.Name,
			NameCn:        cfg.NameCn,
			Description:   cfg.Description,
			DescriptionCn: cfg.DescriptionCn,
			Version:       cfg.Version,
			Author:        cfg.Author,
			Homepage:      cfg.Homepage,
		},
		Directory: pluginDir,
	}

	// Navigation
	if cfg.Navigation != nil {
		p.Navigation = plugin.Navigation{
			Group:   cfg.Navigation.Group,
			GroupCn: cfg.Navigation.GroupCn,
			Order:   cfg.Navigation.Order,
		}
	} else {
		p.Navigation = plugin.Navigation{
			Group: []string{cfg.Name},
			Order: 999,
		}
	}

	// Icon
	if cfg.Icon != "" {
		p.Metadata.Icon = l.resolvePath(pluginDir, cfg.Icon)
	}

	// Form
	if cfg.Form != nil {
		p.Form.SchemaPath = l.resolvePath(pluginDir, cfg.Form.Schema)
		p.Form.UISchemaPath = l.resolvePath(pluginDir, cfg.Form.UISchema)

		// 加载 schema.json
		// 直接保存原始字节（json.RawMessage），避免 map 序列化时按字母序重排字段。
		if p.Form.SchemaPath != "" {
			schemaData, err := os.ReadFile(p.Form.SchemaPath)
			if err != nil {
				return nil, fmt.Errorf("read schema: %w", err)
			}
			if !json.Valid(schemaData) {
				return nil, fmt.Errorf("parse schema: invalid JSON")
			}
			p.Form.Schema = json.RawMessage(schemaData)
		}

		// 加载 uischema.json
		if p.Form.UISchemaPath != "" {
			uiData, err := os.ReadFile(p.Form.UISchemaPath)
			if err == nil && json.Valid(uiData) {
				p.Form.UISchema = json.RawMessage(uiData)
			}
		}
	}

	// Execution
	if cfg.Execution != nil {
		p.Execution.Executable = l.resolvePath(pluginDir, cfg.Execution.Exe)
		p.Execution.WorkingDirectory = cfg.Execution.WorkingDirectory
		p.Execution.Environment = cfg.Execution.Environment

		for _, param := range cfg.Execution.Parameters {
			mapping := plugin.ParameterMapping{
				Field:        param.Field,
				Kind:         plugin.ParameterKind(param.Kind),
				Flag:         param.Flag,
				Style:        param.Style,
				Separator:    param.Separator,
				TrueFlag:     param.TrueFlag,
				FalseFlag:    param.FalseFlag,
				DefaultValue: param.DefaultValue,
			}
			p.Execution.Parameters = append(p.Execution.Parameters, mapping)
		}
	}

	if p.Execution.Executable == "" {
		return nil, fmt.Errorf("plugin %s: executable is required", cfg.ID)
	}

	return p, nil
}

// resolvePath 解析路径（支持相对路径）
func (l *Loader) resolvePath(baseDir, path string) string {
	if path == "" {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}
