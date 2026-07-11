package process

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"

	domainProcess "github.com/cli-tool-center/tool-center/internal/domain/process"
)

type Runner struct {
	mu            sync.Mutex
	processes     map[string]*runningProcess
	outputHandler func(taskID string, line string, source string)
}

type runningProcess struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewRunner() *Runner {
	return &Runner{
		processes: make(map[string]*runningProcess),
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
	rp.cancel()

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		return -1, err
	}
	return exitCode, nil
}

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
		rp.cancel()
		delete(r.processes, processID)
	}
	r.mu.Unlock()
}
