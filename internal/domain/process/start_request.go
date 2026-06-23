package process

// StartRequest 启动进程请求
type StartRequest struct {
	Executable       string
	Arguments        []string
	WorkingDirectory string
	Environment      []string
}
