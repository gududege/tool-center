package wails

import (
	"fmt"
	"time"

	"github.com/cli-tool-center/tool-center/internal/application/plugin"
	domainPlugin "github.com/cli-tool-center/tool-center/internal/domain/plugin"
)

// PluginApi 插件相关 API
type PluginApi struct {
	service *plugin.Service
}

// NewPluginApi 创建插件 API
func NewPluginApi(service *plugin.Service) *PluginApi {
	return &PluginApi{service: service}
}

// GetPlugins 获取插件摘要列表
func (a *PluginApi) GetPlugins() ([]PluginSummaryDto, error) {
	plugins, err := a.service.List()
	if err != nil {
		return nil, fmt.Errorf("get plugins: %w", err)
	}

	dtos := make([]PluginSummaryDto, 0, len(plugins))
	for _, p := range plugins {
		dtos = append(dtos, pluginToSummary(p))
	}
	return dtos, nil
}

// GetPlugin 获取完整插件定义
func (a *PluginApi) GetPlugin(pluginID string) (*PluginDto, error) {
	p, err := a.service.Get(pluginID)
	if err != nil {
		return nil, fmt.Errorf("get plugin %s: %w", pluginID, err)
	}

	dto := pluginToFull(p)
	return &dto, nil
}

// ReloadPlugins 重新加载所有插件
func (a *PluginApi) ReloadPlugins() (*ReloadPluginsResponse, error) {
	count, err := a.service.Reload()
	if err != nil {
		return nil, fmt.Errorf("reload plugins: %w", err)
	}
	return &ReloadPluginsResponse{Count: count}, nil
}

// ValidateFormData 验证表单数据
func (a *PluginApi) ValidateFormData(pluginID string, data map[string]any) (*ValidationResultDto, error) {
	// 此处可以进行简单的 JSON Schema 验证
	return &ValidationResultDto{Valid: true}, nil
}

// SavePluginHistory 保存当前参数到插件目录的 history.json
func (a *PluginApi) SavePluginHistory(pluginID string, data map[string]any) error {
	if err := a.service.SaveHistory(pluginID, data); err != nil {
		return fmt.Errorf("save plugin history %s: %w", pluginID, err)
	}
	return nil
}

// GetPluginHistory 获取插件参数历史列表
func (a *PluginApi) GetPluginHistory(pluginID string) ([]PluginHistoryEntryDto, error) {
	entries, err := a.service.LoadHistory(pluginID)
	if err != nil {
		return nil, fmt.Errorf("get plugin history %s: %w", pluginID, err)
	}

	dtos := make([]PluginHistoryEntryDto, len(entries))
	for i, e := range entries {
		dtos[i] = PluginHistoryEntryDto{
			Timestamp: e.Timestamp.Format(time.RFC3339),
			Label:     e.Label,
			FormData:  e.FormData,
		}
	}
	return dtos, nil
}

func pluginToSummary(p *domainPlugin.Plugin) PluginSummaryDto {
	return PluginSummaryDto{
		ID:            p.Metadata.ID,
		Name:          p.Metadata.Name,
		NameCn:        p.Metadata.NameCn,
		Description:   p.Metadata.Description,
		DescriptionCn: p.Metadata.DescriptionCn,
		Version:       p.Metadata.Version,
		Icon:          p.Metadata.Icon,
		Navigation: NavigationDto{
			Group:   p.Navigation.Group,
			GroupCn: p.Navigation.GroupCn,
			Order:   p.Navigation.Order,
		},
	}
}

func pluginToFull(p *domainPlugin.Plugin) PluginDto {
	params := make([]ParameterMappingDto, len(p.Execution.Parameters))
	for i, pm := range p.Execution.Parameters {
		params[i] = ParameterMappingDto{
			Field:        pm.Field,
			Kind:         string(pm.Kind),
			Flag:         pm.Flag,
			Style:        pm.Style,
			Separator:    pm.Separator,
			TrueFlag:     pm.TrueFlag,
			FalseFlag:    pm.FalseFlag,
			DefaultValue: pm.DefaultValue,
		}
	}
	return PluginDto{
		Metadata: PluginMetadataDto{
			ID:            p.Metadata.ID,
			Name:          p.Metadata.Name,
			NameCn:        p.Metadata.NameCn,
			Description:   p.Metadata.Description,
			DescriptionCn: p.Metadata.DescriptionCn,
			Version:       p.Metadata.Version,
			Author:        p.Metadata.Author,
			Icon:          p.Metadata.Icon,
		},
		Navigation: NavigationDto{
			Group:   p.Navigation.Group,
			GroupCn: p.Navigation.GroupCn,
			Order:   p.Navigation.Order,
		},
		Form: FormDefinitionDto{
			Schema:   p.Form.Schema,
			UISchema: p.Form.UISchema,
		},
		Execution: ExecutionDefinitionDto{
			Exe:              p.Execution.Executable,
			WorkingDirectory: p.Execution.WorkingDirectory,
			Parameters:       params,
		},
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
