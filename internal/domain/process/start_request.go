package process

// StartRequest 启动进程请求
type StartRequest struct {
	Target           string
	Runner           string
	RunnerPath       string
	RunnerOptions    []string
	Arguments        []string
	WorkingDirectory string
	Environment      []string
}
