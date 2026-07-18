package process

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	domainProcess "github.com/cli-tool-center/tool-center/internal/domain/process"
)

type Runner struct {
	mu            sync.Mutex
	processes     map[string]*runningProcess
	outputHandler func(taskID string, line string, source string)
	runnerPaths   map[string]string
	runnerOptions map[string]string
}

type runningProcess struct {
	cmd       *exec.Cmd
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	killTimer *time.Timer // Cancel 兜底 Kill 定时器，nil 表示未启用
}

func NewRunner(runnerPaths, runnerOptions map[string]string) *Runner {
	return &Runner{
		processes:     make(map[string]*runningProcess),
		runnerPaths:   runnerPaths,
		runnerOptions: runnerOptions,
	}
}

func (r *Runner) SetOutputHandler(handler func(taskID string, line string, source string)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.outputHandler = handler
}

func (r *Runner) getOutputHandler() func(taskID string, line string, source string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.outputHandler
}

func (r *Runner) Start(processID string, req domainProcess.StartRequest) error {
	_, cancel := context.WithCancel(context.Background())

	var cmd *exec.Cmd
	if req.Runner == "binary" || req.Runner == "" {
		cmd = exec.Command(req.Target, req.Arguments...)
	} else {
		opts := req.RunnerOptions
		if len(opts) == 0 {
			opts = strings.Fields(r.runnerOptions[req.Runner])
		}
		allArgs := make([]string, 0, len(opts)+1+len(req.Arguments))
		allArgs = append(allArgs, opts...)
		allArgs = append(allArgs, req.Target)
		allArgs = append(allArgs, req.Arguments...)
		cmd = exec.Command(req.RunnerPath, allArgs...)
	}
	setHideWindow(cmd)

	if req.WorkingDirectory != "" {
		cmd.Dir = req.WorkingDirectory
	}
	if len(req.Environment) > 0 {
		cmd.Env = req.Environment
	}

	rp := &runningProcess{
		cmd:    cmd,
		cancel: cancel,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("create stdout pipe: %w", err)
	}

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

	rp.wg.Add(2)

	go func() {
		defer rp.wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if handler := r.getOutputHandler(); handler != nil {
				handler(processID, scanner.Text(), "stdout")
			}
		}
	}()

	go func() {
		defer rp.wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if handler := r.getOutputHandler(); handler != nil {
				handler(processID, scanner.Text(), "stderr")
			}
		}
	}()

	return nil
}

func (r *Runner) Wait(processID string) (int, error) {
	r.mu.Lock()
	rp, ok := r.processes[processID]
	r.mu.Unlock()
	if !ok {
		return -1, fmt.Errorf("process %s not found", processID)
	}

	err := rp.cmd.Wait()
	rp.wg.Wait()
	if rp.killTimer != nil {
		rp.killTimer.Stop()
	}
	rp.cancel()

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		return -1, err
	}
	return exitCode, nil
}

// Cancel 优雅取消：发 Ctrl+C → 等 5min → Kill（兜底）
func (r *Runner) Cancel(processID string) error {
	r.mu.Lock()
	rp, ok := r.processes[processID]
	r.mu.Unlock()
	if !ok {
		return fmt.Errorf("process %s not found", processID)
	}

	if sendInterrupt(rp.cmd.Process) {
		// 信号已发出 → 等 5 分钟让进程优雅退出再 Kill
		// 使用 AfterFunc 以便进程提前退出时能通过 Stop 取消
		rp.killTimer = time.AfterFunc(5*time.Minute, func() {
			rp.cmd.Process.Kill()
		})
	} else {
		// 进程无控制台（PowerShell with HideWindow），无法发信号 → 直接 Kill
		rp.cmd.Process.Kill()
	}

	rp.cancel()
	return nil
}

// Stop 强制停止：直接 Kill
func (r *Runner) Stop(processID string) error {
	r.mu.Lock()
	rp, ok := r.processes[processID]
	r.mu.Unlock()
	if !ok {
		return fmt.Errorf("process %s not found", processID)
	}

	rp.cmd.Process.Kill()
	if rp.killTimer != nil {
		rp.killTimer.Stop()
	}
	rp.cancel()
	return nil
}

func (r *Runner) IsRunning(processID string) bool {
	r.mu.Lock()
	_, ok := r.processes[processID]
	r.mu.Unlock()
	return ok
}

func (r *Runner) Remove(processID string) {
	r.mu.Lock()
	rp, ok := r.processes[processID]
	if ok {
		if rp.killTimer != nil {
			rp.killTimer.Stop()
		}
		rp.cancel()
		delete(r.processes, processID)
	}
	r.mu.Unlock()
}

func (r *Runner) UpdateRunnerConfig(paths, options map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runnerPaths = paths
	r.runnerOptions = options
}
