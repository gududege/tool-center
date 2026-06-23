package process

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"

	domainProcess "github.com/cli-tool-center/tool-center/internal/domain/process"
)

// Runner 进程运行器
type Runner struct {
	mu            sync.Mutex
	processes     map[string]*runningProcess
	outputHandler func(taskID string, line string, source string)
}

type runningProcess struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
	taskID string
}

// NewRunner 创建进程运行器
func NewRunner() *Runner {
	return &Runner{
		processes: make(map[string]*runningProcess),
	}
}

// SetOutputHandler 设置输出处理器
func (r *Runner) SetOutputHandler(handler func(taskID string, line string, source string)) {
	r.outputHandler = handler
}

// Start 启动进程
func (r *Runner) Start(processID string, req domainProcess.StartRequest) error {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, req.Executable, req.Arguments...)

	if req.WorkingDirectory != "" {
		cmd.Dir = req.WorkingDirectory
	}
	if len(req.Environment) > 0 {
		cmd.Env = req.Environment
	}

	rp := &runningProcess{
		cmd:    cmd,
		cancel: cancel,
		taskID: processID,
	}

	// stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("create stdout pipe: %w", err)
	}

	// stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("start process: %w", err)
	}

	r.mu.Lock()
	r.processes[processID] = rp
	r.mu.Unlock()

	// 采集 stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if r.outputHandler != nil {
				r.outputHandler(processID, scanner.Text(), "stdout")
			}
		}
	}()

	// 采集 stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if r.outputHandler != nil {
				r.outputHandler(processID, scanner.Text(), "stderr")
			}
		}
	}()

	return nil
}

// Wait 等待进程结束并返回退出码
func (r *Runner) Wait(processID string) (int, error) {
	r.mu.Lock()
	rp, ok := r.processes[processID]
	r.mu.Unlock()
	if !ok {
		return -1, fmt.Errorf("process %s not found", processID)
	}

	err := rp.cmd.Wait()
	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		// context canceled 等
		return -1, err
	}
	return exitCode, nil
}

// Stop 停止进程
func (r *Runner) Stop(processID string) error {
	r.mu.Lock()
	rp, ok := r.processes[processID]
	r.mu.Unlock()
	if !ok {
		return fmt.Errorf("process %s not found", processID)
	}

	rp.cancel()
	return nil
}

// IsRunning 检查进程是否还在运行
func (r *Runner) IsRunning(processID string) bool {
	r.mu.Lock()
	_, ok := r.processes[processID]
	r.mu.Unlock()
	return ok
}

// Remove 移除进程记录
func (r *Runner) Remove(processID string) {
	r.mu.Lock()
	delete(r.processes, processID)
	r.mu.Unlock()
}
