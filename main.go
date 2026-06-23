package main

import (
	"context"
	"embed"
	"log"

	"github.com/cli-tool-center/tool-center/internal/bootstrap"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化依赖
	app, err := bootstrap.Initialize()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	// 创建 Wails 应用
	err = wails.Run(&options.App{
		Title:  "CLI Tool Center",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			// 恢复上次的窗口大小，窗口居中
			settings, _ := app.SettingsSvc.GetSettings()
			if settings != nil && settings.WindowWidth > 0 && settings.WindowHeight > 0 {
				wailsRuntime.WindowSetSize(ctx, settings.WindowWidth, settings.WindowHeight)
			}
			wailsRuntime.WindowCenter(ctx)
			// 文件/目录对话框需经由 wails runtime context 才能唤起原生系统对话框。
			app.SystemApi.SetContext(ctx)
			// 设置事件桥接上下文
			app.EventBridge.SetContext(ctx)
			app.EventBridge.Start()
			// 启动窗口尺寸定时保存
			app.StartWindowSizeWatcher(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			// 停止窗口尺寸 watcher（停止前会再保存一次精确尺寸）
			app.StopWindowSizeWatcher(ctx)
			app.EventBridge.Stop()
		},
		Bind: []interface{}{
			// API 层绑定到前端
			app.PluginApi,
			app.TaskApi,
			app.SettingsApi,
			app.SystemApi,
		},
	})

	if err != nil {
		log.Fatalf("application error: %v", err)
	}
}
