package bootstrap

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/cli-tool-center/tool-center/internal/adapter/wails"
	appPlugin "github.com/cli-tool-center/tool-center/internal/application/plugin"
	appSettings "github.com/cli-tool-center/tool-center/internal/application/settings"
	appTask "github.com/cli-tool-center/tool-center/internal/application/task"
	"github.com/cli-tool-center/tool-center/internal/domain/settings"
	"github.com/cli-tool-center/tool-center/internal/infrastructure/events"
	"github.com/cli-tool-center/tool-center/internal/infrastructure/filesystem"
	processInfra "github.com/cli-tool-center/tool-center/internal/infrastructure/process"
	"github.com/cli-tool-center/tool-center/internal/infrastructure/repositories"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用容器，持有所有依赖
type App struct {
	// Settings
	Settings     *settings.Settings
	SettingsRepo *appSettings.Repository
	SettingsSvc  *appSettings.Service
	SettingsApi  *wails.SettingsApi

	// Window size watcher
	WindowSizeStop   chan struct{}
	LastWindowWidth  int
	LastWindowHeight int

	// Event
	EventBus    *events.InMemoryBus
	EventBridge *wails.EventBridge

	// Plugin
	PluginRepo *repositories.MemoryPluginRepository
	PluginSvc  *appPlugin.Service
	PluginApi  *wails.PluginApi

	// Task
	TaskRepo *repositories.MemoryTaskRepository
	TaskSvc  *appTask.Service
	TaskApi  *wails.TaskApi

	// Process
	ProcessRunner    *processInfra.Runner
	ParameterBuilder *processInfra.ParameterBuilder

	// System
	SystemApi *wails.SystemApi
}

// Initialize 初始化所有依赖
func Initialize() (*App, error) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	slog.Info("initializing application...")

	// 1. Settings
	settingsPath := "settings.json"
	settingsRepo := appSettings.NewRepository(settingsPath)
	if err := settingsRepo.Load(); err != nil {
		slog.Warn("load settings", "error", err)
	}
	settingsSvc := appSettings.NewService(settingsRepo)
	settingsApi := wails.NewSettingsApi(settingsSvc)
	cfg, _ := settingsSvc.GetSettings()

	// 2. Event Bus
	eventBus := events.NewInMemoryBus()

	// 3. Repositories
	pluginRepo := repositories.NewMemoryPluginRepository()
	taskRepo := repositories.NewMemoryTaskRepository()

	// 4. Infrastructure
	pluginLoader := filesystem.NewLoader()
	paramBuilder := processInfra.NewParameterBuilder()
	processRunner := processInfra.NewRunner()

	// 5. Application Services
	pluginSvc := appPlugin.NewService(pluginRepo, pluginLoader, eventBus, cfg)
	taskSvc := appTask.NewService(taskRepo, pluginRepo, processRunner, paramBuilder, eventBus, cfg.MaxOutputLines)

	// 6. Wails APIs
	pluginApi := wails.NewPluginApi(pluginSvc)
	taskApi := wails.NewTaskApi(taskSvc)
	systemApi := wails.NewSystemApi()

	// 7. Event Bridge
	eventBridge := wails.NewEventBridge(eventBus)

	// 8. Load plugins
	slog.Info("loading plugins...")
	if err := pluginSvc.LoadPlugins(); err != nil {
		slog.Error("load plugins", "error", err)
	}

	slog.Info("application initialized")
	return &App{
		Settings:         cfg,
		SettingsRepo:     settingsRepo,
		SettingsSvc:      settingsSvc,
		SettingsApi:      settingsApi,
		EventBus:         eventBus,
		EventBridge:      eventBridge,
		PluginRepo:       pluginRepo,
		PluginSvc:        pluginSvc,
		PluginApi:        pluginApi,
		TaskRepo:         taskRepo,
		TaskSvc:          taskSvc,
		TaskApi:          taskApi,
		ProcessRunner:    processRunner,
		ParameterBuilder: paramBuilder,
		SystemApi:        systemApi,
	}, nil
}

// StartWindowSizeWatcher starts a goroutine that polls the exact window size
// every 2 seconds and persists it to settings. It skips maximized windows so
// the next launch does not open in a maximized state.
func (a *App) StartWindowSizeWatcher(ctx context.Context) {
	a.WindowSizeStop = make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.saveWindowSize(ctx)
			case <-a.WindowSizeStop:
				return
			}
		}
	}()
}

// StopWindowSizeWatcher stops the polling goroutine and saves the final size.
func (a *App) StopWindowSizeWatcher(ctx context.Context) {
	if a.WindowSizeStop != nil {
		close(a.WindowSizeStop)
		a.WindowSizeStop = nil
	}
	a.saveWindowSize(ctx)
}

func (a *App) saveWindowSize(ctx context.Context) {
	w, h := wailsRuntime.WindowGetSize(ctx)
	if w <= 0 || h <= 0 {
		return
	}
	if wailsRuntime.WindowIsMaximised(ctx) {
		return
	}
	if w == a.LastWindowWidth && h == a.LastWindowHeight {
		return
	}

	settings, err := a.SettingsSvc.GetSettings()
	if err != nil {
		slog.Error("get settings for window size", "error", err)
		return
	}

	settings.WindowWidth = w
	settings.WindowHeight = h
	if err := a.SettingsSvc.SaveSettings(settings); err != nil {
		slog.Error("save window size", "error", err)
		return
	}

	a.LastWindowWidth = w
	a.LastWindowHeight = h
	slog.Debug("window size saved", "width", w, "height", h)
}
