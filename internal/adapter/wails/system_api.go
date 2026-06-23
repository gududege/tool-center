package wails

import (
	"context"
	"runtime"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	AppVersion = "1.0.0"
	BuildTime  = "unknown"
)

// SystemApi 系统相关 API
type SystemApi struct {
	ctx context.Context
}

// NewSystemApi 创建系统 API
func NewSystemApi() *SystemApi {
	return &SystemApi{}
}

// SetContext 注入 Wails runtime context（在 OnStartup 中调用，与 EventBridge 同模式）。
// 文件/目录对话框必须经由 wails runtime 的 context 才能唤起原生系统对话框。
func (a *SystemApi) SetContext(ctx context.Context) {
	a.ctx = ctx
}

// GetSystemInfo 获取系统信息
func (a *SystemApi) GetSystemInfo() SystemInfoDto {
	return SystemInfoDto{
		AppVersion: AppVersion,
		BuildTime:  BuildTime,
		GoVersion:  runtime.Version(),
		Os:         runtime.GOOS,
	}
}

// toWailsFilters 把 DTO 的 []FileFilterDto 转成 wails runtime 的 []runtime.FileFilter。
// wails 的 FileFilter.Pattern 是分号分隔的单个字符串（如 "*.txt;*.csv"），而 DTO 用数组。
func toWailsFilters(filters []FileFilterDto) []wailsRuntime.FileFilter {
	out := make([]wailsRuntime.FileFilter, 0, len(filters))
	for _, f := range filters {
		out = append(out, wailsRuntime.FileFilter{
			DisplayName: f.DisplayName,
			Pattern:     strings.Join(f.Patterns, ";"),
		})
	}
	return out
}

// OpenFileDialog 打开文件选择对话框（原生系统对话框）。
// 取消选择时返回空 DialogResponse（path 为空字符串），由前端判空处理。
func (a *SystemApi) OpenFileDialog(title string, filters []FileFilterDto) (*DialogResponse, error) {
	if a.ctx == nil {
		return &DialogResponse{}, nil
	}
	path, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title:   title,
		Filters: toWailsFilters(filters),
	})
	if err != nil {
		return nil, err
	}
	return &DialogResponse{Path: path}, nil
}

// OpenFolderDialog 打开文件夹选择对话框（原生系统对话框）。
func (a *SystemApi) OpenFolderDialog(title string) (*DialogResponse, error) {
	if a.ctx == nil {
		return &DialogResponse{}, nil
	}
	path, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return nil, err
	}
	return &DialogResponse{Path: path}, nil
}

// SaveFileDialog 打开保存文件对话框（原生系统对话框）。
func (a *SystemApi) SaveFileDialog(title string, filters []FileFilterDto) (*DialogResponse, error) {
	if a.ctx == nil {
		return &DialogResponse{}, nil
	}
	path, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title:   title,
		Filters: toWailsFilters(filters),
	})
	if err != nil {
		return nil, err
	}
	return &DialogResponse{Path: path}, nil
}
