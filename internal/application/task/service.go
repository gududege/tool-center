package task

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/cli-tool-center/tool-center/internal/domain/events"
	"github.com/cli-tool-center/tool-center/internal/domain/output"
	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
	domainProcess "github.com/cli-tool-center/tool-center/internal/domain/process"
	domainTask "github.com/cli-tool-center/tool-center/internal/domain/task"
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
	outputs      map[string]*output.OutputBuffer
	maxOutput    int
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
		Status:    domainTask.TaskCreated,
		CreatedAt: time.Now(),
	}
	if err := s.tasks.Save(task); err != nil {
		return "", fmt.Errorf("save task: %w", err)
	}

	// 初始化输出缓冲区
	s.outputs[task.ID] = &output.OutputBuffer{
		TaskID: task.ID,
		Events: make([]output.OutputEvent, 0),
		MaxLen: s.maxOutput,
	}

	// 发送创建事件
	s.eventBus.Publish(events.TaskCreated{
		TaskID:   task.ID,
		PluginID: pluginID,
	})

	// 异步执行
	go s.executeTask(task, p, args)

	return task.ID, nil
}

func (s *Service) executeTask(task *domainTask.Task, p *plugin.Plugin, args []string) {
	// 更新状态为 queued
	task.Status = domainTask.TaskQueued
	_ = s.tasks.Save(task)

	// 更新为 running
	now := time.Now()
	task.Status = domainTask.TaskRunning
	task.StartedAt = &now
	_ = s.tasks.Save(task)
	s.eventBus.Publish(events.TaskStarted{TaskID: task.ID})

	s.appendSystemOutput(task.ID, "Task started")

	// 准备环境变量
	env := make([]string, 0)
	for k, v := range p.Execution.Environment {
		env = append(env, k+"="+v)
	}

	// 启动进程
	err := s.runner.Start(task.ID, domainProcess.StartRequest{
		Executable:       p.Execution.Executable,
		Arguments:        args,
		WorkingDirectory: p.Execution.WorkingDirectory,
		Environment:      env,
	})
	if err != nil {
		task.Status = domainTask.TaskFailed
		task.ErrorMessage = err.Error()
		now := time.Now()
		task.EndedAt = &now
		_ = s.tasks.Save(task)
		s.appendSystemOutput(task.ID, fmt.Sprintf("Task failed: %v", err))
		s.eventBus.Publish(events.TaskFailed{TaskID: task.ID, Error: err.Error()})
		return
	}

	// 等待结束
	exitCode, waitErr := s.runner.Wait(task.ID)

	if waitErr != nil {
		// 检查是否是取消导致的
		t, _ := s.tasks.Get(task.ID)
		if t != nil && t.Status == domainTask.TaskCancelled {
			s.appendSystemOutput(task.ID, "Task cancelled")
			s.eventBus.Publish(events.TaskCancelled{TaskID: task.ID})
			_ = s.runner.Stop(task.ID)
			s.runner.Remove(task.ID)
			return
		}

		task.Status = domainTask.TaskFailed
		task.ErrorMessage = waitErr.Error()
		now := time.Now()
		task.EndedAt = &now
		_ = s.tasks.Save(task)
		s.appendSystemOutput(task.ID, fmt.Sprintf("Task failed: %v", waitErr))
		s.eventBus.Publish(events.TaskFailed{TaskID: task.ID, Error: waitErr.Error()})
		_ = s.runner.Stop(task.ID)
		s.runner.Remove(task.ID)
		return
	}

	// 正常完成
	task.Status = domainTask.TaskCompleted
	task.ExitCode = &exitCode
	now = time.Now()
	task.EndedAt = &now
	_ = s.tasks.Save(task)
	s.appendSystemOutput(task.ID, fmt.Sprintf("Task completed with exit code %d", exitCode))
	s.eventBus.Publish(events.TaskCompleted{TaskID: task.ID, ExitCode: exitCode})
	s.runner.Remove(task.ID)
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

	// 追加到缓冲区
	if buf, ok := s.outputs[taskID]; ok {
		buf.Append(event)
	}

	// 发布事件
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

	if buf, ok := s.outputs[taskID]; ok {
		buf.Append(event)
	}

	s.eventBus.Publish(events.OutputReceived{
		TaskID: taskID,
		Event:  event,
	})
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
	task.Status = domainTask.TaskCancelled
	_ = s.tasks.Save(task)

	// 停止进程
	_ = s.runner.Stop(id)

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
	_ = s.tasks.Delete(id)
	delete(s.outputs, id)
	return nil
}

// GetTaskOutput 获取任务输出
func (s *Service) GetTaskOutput(taskID string) ([]output.OutputEvent, error) {
	buf, ok := s.outputs[taskID]
	if !ok {
		return nil, fmt.Errorf("output for task %s not found", taskID)
	}
	return buf.GetEvents(), nil
}

// generateID 生成简单 ID
func generateID() string {
	return fmt.Sprintf("task-%d", time.Now().UnixNano())
}
