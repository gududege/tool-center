package wails

import (
	"github.com/cli-tool-center/tool-center/internal/application/settings"
	domainSettings "github.com/cli-tool-center/tool-center/internal/domain/settings"
)

// SettingsApi 配置相关 API
type SettingsApi struct {
	service *settings.Service
}

// NewSettingsApi 创建配置 API
func NewSettingsApi(service *settings.Service) *SettingsApi {
	return &SettingsApi{service: service}
}

// GetSettings 获取配置
func (a *SettingsApi) GetSettings() (*SettingsDto, error) {
	s, err := a.service.GetSettings()
	if err != nil {
		return nil, err
	}
	return &SettingsDto{
		Theme:            s.Theme,
		FormTheme:        s.FormTheme,
		PluginDirectory:  s.PluginDirectory,
		Language:         s.Language,
		SidebarCollapsed: s.SidebarCollapsed,
		SidebarSize:      s.SidebarSize,
		BottomPanelSize:  s.BottomPanelSize,
		BottomTab:        s.BottomTab,
	}, nil
}

// SaveSettings 保存配置
func (a *SettingsApi) SaveSettings(dto SettingsDto) error {
	current, err := a.service.GetSettings()
	if err != nil {
		current = domainSettings.DefaultSettings()
	}
	s := &domainSettings.Settings{
		Theme:             dto.Theme,
		FormTheme:         dto.FormTheme,
		PluginDirectory:   dto.PluginDirectory,
		MaxOutputLines:    current.MaxOutputLines,
		MaxTaskHistory:    current.MaxTaskHistory,
		AutoReloadPlugins: current.AutoReloadPlugins,
		Language:          dto.Language,
		SidebarCollapsed:  dto.SidebarCollapsed,
		SidebarSize:       dto.SidebarSize,
		BottomPanelSize:   dto.BottomPanelSize,
		BottomTab:         dto.BottomTab,
		// WindowWidth/WindowHeight 不在 DTO 里（仅退出时由 main.go OnShutdown
		// 经 SettingsSvc.SaveSettings 直存）。前端走本方法保存时必须回填 current
		// 的窗口尺寸，否则会把 OnShutdown 存好的值清零，导致下次启动回退默认尺寸。
		WindowWidth:  current.WindowWidth,
		WindowHeight: current.WindowHeight,
	}
	return a.service.SaveSettings(s)
}
