package settings

// Settings 系统配置
type Settings struct {
	Theme             string
	// FormTheme 是 RJSF shadcn 的子主题 id（default/caffeine/claude/clean-slate/
	// amethyst-haze/neo-brutalism/pastel-dreams/soft-pop/twitter/vercel）。仅存储
	// 字符串，明暗模式跟随 Theme 字段。
	FormTheme         string
	PluginDirectory   string
	MaxOutputLines    int
	MaxTaskHistory    int
	AutoReloadPlugins bool
	Language          string
	SidebarCollapsed  bool
	SidebarSize       int
	BottomPanelSize   int
	BottomTab         string
	WindowWidth       int
	WindowHeight      int
}

// RuntimeState 运行时状态
type RuntimeState struct {
	StartedAt      string
	LoadedPlugins  int
	RunningTasks   int
	ActiveSessions int
}

// DefaultSettings 返回默认配置
func DefaultSettings() *Settings {
	return &Settings{
		Theme:             "dark",
		FormTheme:         "default",
		PluginDirectory:   "./plugins",
		MaxOutputLines:    10000,
		MaxTaskHistory:    100,
		AutoReloadPlugins: false,
		Language:          "en",
		SidebarCollapsed:  false,
		SidebarSize:       18,
		BottomPanelSize:   35,
		BottomTab:         "tasks",
		WindowWidth:       1200,
		WindowHeight:      800,
	}
}
