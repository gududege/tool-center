package task

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/cli-tool-center/tool-center/internal/domain/events"
	"github.com/cli-tool-center/tool-center/internal/domain/output"
	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
	domainProcess "github.com/cli-tool-center/tool-center/internal/domain/process"
	domainTask "github.com/cli-tool-center/tool-center/internal/domain/task"
	"github.com/google/uuid"
)

// ProcessRunner 进程运行器接口
type ProcessRunner interface {
	Start(processID string, req domainProcess.StartRequest) error
	Wait(processID string) (int, error)
	Stop(processID string) error
	Remove(processID string)
	SetOutputHandler(handler func(taskID string, line string, source string))
}

// ParameterBuilder 参数构建器接口
type ParameterBuilder interface {
	Build(mappings []plugin.ParameterMapping, formData map[string]any) ([]string, error)
}

// PluginRepository 插件仓储接口（应用层定义）
type PluginRepository interface {
	Get(id string) (*plugin.Plugin, error)
}

// Service 任务应用服务
type Service struct {
	tasks        domainTask.Repository
	plugins      PluginRepository
	runner       ProcessRunner
	paramBuilder ParameterBuilder
	eventBus     events.Bus
	outputsMu    sync.RWMutex
	outputs      map[string]*output.OutputBuffer
	maxOutput    int
	cancels      sync.Map
}

// NewService 创建任务服务
func NewService(
	tasks domainTask.Repository,
	plugins PluginRepository,
	runner ProcessRunner,
	paramBuilder ParameterBuilder,
	eventBus events.Bus,
	maxOutput int,
) *Service {
	s := &Service{
		tasks:        tasks,
		plugins:      plugins,
		runner:       runner,
		paramBuilder: paramBuilder,
		eventBus:     eventBus,
		outputs:      make(map[string]*output.OutputBuffer),
		maxOutput:    maxOutput,
	}

	// 设置输出处理器
	runner.SetOutputHandler(s.handleOutput)

	return s
}

// RunPlugin 执行插件，返回任务ID
func (s *Service) RunPlugin(pluginID string, formData map[string]any) (string, error) {
	p, err := s.plugins.Get(pluginID)
	if err != nil {
		return "", fmt.Errorf("plugin %s not found: %w", pluginID, err)
	}

	// 构建参数
	args, err := s.paramBuilder.Build(p.Execution.Parameters, formData)
	if err != nil {
		return "", fmt.Errorf("build arguments: %w", err)
	}

	slog.Info("run plugin", "pluginId", pluginID, "exe", p.Execution.Executable, "args", args)

	// 创建任务
	task := &domainTask.Task{
		ID:        generateID(),
		PluginID:  pluginID,
		Status:    domainTask.Created,
		CreatedAt: time.Now(),
	}
	if err := s.tasks.Save(task); err != nil {
		return "", fmt.Errorf("save task: %w", err)
	}

	// 初始化输出缓冲区
	s.outputsMu.Lock()
	s.outputs[task.ID] = &output.OutputBuffer{
		TaskID: task.ID,
		Events: make([]output.OutputEvent, 0),
		MaxLen: s.maxOutput,
	}
	s.outputsMu.Unlock()

	// 发送创建事件
	s.eventBus.Publish(events.TaskCreated{
		TaskID:   task.ID,
		PluginID: pluginID,
	})

	// 异步执行
	go s.executeTask(task.ID, p, args)

	return task.ID, nil
}

func (s *Service) executeTask(taskID string, p *plugin.Plugin, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancels.Store(taskID, cancel)
	defer s.cancels.Delete(taskID)
	defer cancel()

	s.transitionTaskStatus(taskID, domainTask.Queued)
	s.transitionTaskStatus(taskID, domainTask.Running)

	_ = ctx // placeholder for future context propagation to runner

	s.eventBus.Publish(events.TaskStarted{TaskID: taskID})
	s.appendSystemOutput(taskID, "Task started")

	env := make([]string, 0)
	for k, v := range p.Execution.Environment {
		env = append(env, k+"="+v)
	}

	err := s.runner.Start(taskID, domainProcess.StartRequest{
		Executable:       p.Execution.Executable,
		Arguments:        args,
		WorkingDirectory: p.Execution.WorkingDirectory,
		Environment:      env,
	})
	if err != nil {
		s.saveTaskFinalState(taskID, domainTask.Failed, err.Error(), nil)
		s.appendSystemOutput(taskID, fmt.Sprintf("Task failed: %v", err))
		s.eventBus.Publish(events.TaskFailed{TaskID: taskID, Error: err.Error()})
		return
	}

	exitCode, waitErr := s.runner.Wait(taskID)

	if waitErr != nil {
		t, err := s.tasks.Get(taskID)
		if err != nil {
			slog.Warn("executeTask: get task failed", "id", taskID, "error", err)
			s.runner.Remove(taskID)
			return
		}
		if t.Status.IsTerminal() {
			if t.Status == domainTask.Cancelled {
				s.appendSystemOutput(taskID, "Task cancelled")
				s.eventBus.Publish(events.TaskCancelled{TaskID: taskID})
			}
			s.runner.Remove(taskID)
			return
		}

		s.saveTaskFinalState(taskID, domainTask.Failed, waitErr.Error(), nil)
		s.appendSystemOutput(taskID, fmt.Sprintf("Task failed: %v", waitErr))
		s.eventBus.Publish(events.TaskFailed{TaskID: taskID, Error: waitErr.Error()})
		if err := s.runner.Stop(taskID); err != nil {
			slog.Warn("executeTask: stop process failed", "id", taskID, "error", err)
		}
		s.runner.Remove(taskID)
		return
	}

	s.saveTaskFinalState(taskID, domainTask.Completed, "", func(t *domainTask.Task) {
		t.ExitCode = &exitCode
	})
	s.appendSystemOutput(taskID, fmt.Sprintf("Task completed with exit code %d", exitCode))
	s.eventBus.Publish(events.TaskCompleted{TaskID: taskID, ExitCode: exitCode})
	s.runner.Remove(taskID)
}

// transitionTaskStatus 原子更新任务状态
func (s *Service) transitionTaskStatus(taskID string, status domainTask.TaskStatus) {
	t, err := s.tasks.Get(taskID)
	if err != nil {
		slog.Warn("transition task status: task not found", "id", taskID, "status", status)
		return
	}
	t.Status = status
	if status == domainTask.Running {
		now := time.Now()
		t.StartedAt = &now
	}
	if err := s.tasks.Save(t); err != nil {
		slog.Warn("transition task status: save failed", "id", taskID, "status", status, "error", err)
	}
}

// saveTaskFinalState 保存任务终态，若已被 CancelTask 设为终态则跳过。
// apply 可选，用于在 Save 前设置附加字段（如 ExitCode）。
func (s *Service) saveTaskFinalState(taskID string, status domainTask.TaskStatus, errMsg string, apply func(t *domainTask.Task)) {
	t, err := s.tasks.Get(taskID)
	if err != nil {
		slog.Warn("save task final state: task not found", "id", taskID)
		return
	}
	if t.Status.IsTerminal() {
		return
	}
	t.Status = status
	t.ErrorMessage = errMsg
	now := time.Now()
	t.EndedAt = &now
	if apply != nil {
		apply(t)
	}
	if err := s.tasks.Save(t); err != nil {
		slog.Warn("save task final state: save failed", "id", taskID, "error", err)
	}
}

// handleOutput 处理进程输出
func (s *Service) handleOutput(taskID string, line string, source string) {
	sourceType := output.StdoutSource
	if source == "stderr" {
		sourceType = output.StderrSource
	}

	event := output.OutputEvent{
		TaskID:    taskID,
		Timestamp: time.Now(),
		Source:    sourceType,
		Level:     output.InfoLevel,
		Message:   line,
	}

	if buf := s.getOutputBuffer(taskID); buf != nil {
		buf.Append(event)
	}

	s.eventBus.Publish(events.OutputReceived{
		TaskID: taskID,
		Event:  event,
	})
}

// appendSystemOutput 追加系统日志
func (s *Service) appendSystemOutput(taskID string, message string) {
	event := output.OutputEvent{
		TaskID:    taskID,
		Timestamp: time.Now(),
		Source:    output.SystemSource,
		Level:     output.InfoLevel,
		Message:   message,
	}

	if buf := s.getOutputBuffer(taskID); buf != nil {
		buf.Append(event)
	}

	s.eventBus.Publish(events.OutputReceived{
		TaskID: taskID,
		Event:  event,
	})
}

// getOutputBuffer 线程安全获取输出缓冲区
func (s *Service) getOutputBuffer(taskID string) *output.OutputBuffer {
	s.outputsMu.RLock()
	defer s.outputsMu.RUnlock()
	return s.outputs[taskID]
}

// GetTask 获取单个任务
func (s *Service) GetTask(id string) (*domainTask.Task, error) {
	return s.tasks.Get(id)
}

// ListTasks 获取所有任务
func (s *Service) ListTasks() ([]*domainTask.Task, error) {
	return s.tasks.List()
}

// CancelTask 取消任务
func (s *Service) CancelTask(id string) error {
	task, err := s.tasks.Get(id)
	if err != nil {
		return fmt.Errorf("task %s not found", id)
	}

	if task.Status.IsTerminal() {
		return fmt.Errorf("task %s already finished with status %s", id, task.Status)
	}

	oldStatus := task.Status
	task.Status = domainTask.Cancelled
	if err := s.tasks.Save(task); err != nil {
		slog.Warn("cancel task: save failed", "id", id, "error", err)
	}

	if cancel, ok := s.cancels.Load(id); ok {
		cancel.(context.CancelFunc)()
	}

	if err := s.runner.Stop(id); err != nil {
		slog.Warn("cancel task: stop process failed", "id", id, "error", err)
	}

	s.appendSystemOutput(id, "Task cancelled by user")

	s.eventBus.Publish(events.TaskStatusChanged{
		TaskID:    id,
		OldStatus: string(oldStatus),
		NewStatus: string(task.Status),
	})

	return nil
}

// DeleteTask 删除任务
func (s *Service) DeleteTask(id string) error {
	if err := s.tasks.Delete(id); err != nil {
		slog.Warn("delete task: failed", "id", id, "error", err)
	}
	s.outputsMu.Lock()
	delete(s.outputs, id)
	s.outputsMu.Unlock()
	return nil
}

// GetTaskOutput 获取任务输出
func (s *Service) GetTaskOutput(taskID string) ([]output.OutputEvent, error) {
	buf := s.getOutputBuffer(taskID)
	if buf == nil {
		return nil, fmt.Errorf("output for task %s not found", taskID)
	}
	return buf.GetEvents(), nil
}

// generateID 生成简单 ID
func generateID() string {
	return "task-" + uuid.New().String()
}
