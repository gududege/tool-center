package task

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/cli-tool-center/tool-center/internal/domain/events"
	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
	domainProcess "github.com/cli-tool-center/tool-center/internal/domain/process"
	domainTask "github.com/cli-tool-center/tool-center/internal/domain/task"
)

// --- Stub implementations ---

type stubTaskRepo struct {
	mu    sync.Mutex
	tasks map[string]*domainTask.Task
}

func (s *stubTaskRepo) Get(id string) (*domainTask.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	return t, nil
}

func (s *stubTaskRepo) List() ([]*domainTask.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]*domainTask.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result, nil
}

func (s *stubTaskRepo) Save(task *domainTask.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
	return nil
}

func (s *stubTaskRepo) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
	return nil
}

type stubPluginRepo struct {
	mu      sync.Mutex
	plugins map[string]*plugin.Plugin
}

func (s *stubPluginRepo) Get(id string) (*plugin.Plugin, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.plugins[id]
	if !ok {
		return nil, errors.New("plugin not found")
	}
	return p, nil
}

type stubRunner struct {
	mu      sync.Mutex
	handler func(taskID string, line string, source string)
	started map[string]bool
}

func (s *stubRunner) Start(processID string, req domainProcess.StartRequest) error {
	s.mu.Lock()
	s.started[processID] = true
	s.mu.Unlock()
	if s.handler != nil {
		s.handler(processID, "stub running: "+req.Executable, "stdout")
	}
	return nil
}

func (s *stubRunner) Wait(processID string) (int, error) {
	s.mu.Lock()
	delete(s.started, processID)
	s.mu.Unlock()
	return 0, nil
}

func (s *stubRunner) Stop(processID string) error {
	s.mu.Lock()
	delete(s.started, processID)
	s.mu.Unlock()
	return nil
}

func (s *stubRunner) Remove(processID string) {
	s.mu.Lock()
	delete(s.started, processID)
	s.mu.Unlock()
}

func (s *stubRunner) SetOutputHandler(handler func(taskID string, line string, source string)) {
	s.handler = handler
}

type stubParamBuilder struct{}

func (s *stubParamBuilder) Build(mappings []plugin.ParameterMapping, formData map[string]any) ([]string, error) {
	return []string{"--stub"}, nil
}

type stubFailRunner struct{}

func (s *stubFailRunner) Start(processID string, req domainProcess.StartRequest) error {
	return errors.New("start failed")
}
func (s *stubFailRunner) Wait(processID string) (int, error) { return -1, nil }
func (s *stubFailRunner) Stop(processID string) error        { return nil }
func (s *stubFailRunner) Remove(processID string)            {}
func (s *stubFailRunner) SetOutputHandler(handler func(taskID string, line string, source string)) {}

type eventCollector struct {
	mu     sync.Mutex
	events []events.DomainEvent
}

func (e *eventCollector) Publish(event events.DomainEvent) {
	e.mu.Lock()
	e.events = append(e.events, event)
	e.mu.Unlock()
}

func (e *eventCollector) Subscribe(topic string, handler events.Handler) func() { return func() {} }
func (e *eventCollector) PublishAsync(event events.DomainEvent)                { go e.Publish(event) }

func newStubPlugin() *plugin.Plugin {
	return &plugin.Plugin{
		Metadata: plugin.PluginMetadata{
			ID:   "test-plugin",
			Name: "Test Plugin",
		},
		Navigation: plugin.Navigation{Group: []string{"Test"}, Order: 1},
		Execution: plugin.ExecutionDefinition{
			Executable: "echo",
			Parameters: []plugin.ParameterMapping{
				{Field: "msg", Kind: plugin.ArgumentKind},
			},
		},
	}
}

func waitForTask(t *testing.T, svc *Service, taskID string) {
	t.Helper()
	deadline := time.After(2 * time.Second)
	for {
		task, err := svc.GetTask(taskID)
		if err != nil {
			return // already deleted
		}
		if task.Status.IsTerminal() {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for task %s to complete", taskID)
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// --- Tests ---

func TestService_RunPlugin(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{"test-plugin": newStubPlugin()}},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	taskID, err := svc.RunPlugin("test-plugin", map[string]any{"msg": "hello"})
	if err != nil {
		t.Fatalf("RunPlugin failed: %v", err)
	}
	if taskID == "" {
		t.Error("expected non-empty task ID")
	}

	waitForTask(t, svc, taskID)

	task, err := svc.GetTask(taskID)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if task.PluginID != "test-plugin" {
		t.Errorf("PluginID = %s, want test-plugin", task.PluginID)
	}
	if task.Status != domainTask.TaskCompleted {
		t.Errorf("Status = %s, want %s", task.Status, domainTask.TaskCompleted)
	}
}

func TestService_RunPlugin_PluginNotFound(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: make(map[string]*plugin.Plugin)},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	_, err := svc.RunPlugin("nonexistent", nil)
	if err == nil {
		t.Error("expected error for nonexistent plugin")
	}
}

func TestService_GetTask_NotFound(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{}},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	_, err := svc.GetTask("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestService_ListTasks(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{"test-plugin": newStubPlugin()}},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	tasks, _ := svc.ListTasks()
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks initially, got %d", len(tasks))
	}

	taskID, _ := svc.RunPlugin("test-plugin", nil)
	waitForTask(t, svc, taskID)

	tasks, _ = svc.ListTasks()
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestService_CancelTask(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{"test-plugin": newStubPlugin()}},
		// Use stubRunner so we cancel before it completes
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	taskID, _ := svc.RunPlugin("test-plugin", nil)

	err := svc.CancelTask(taskID)
	if err != nil {
		t.Fatalf("CancelTask failed: %v", err)
	}

	task, _ := svc.GetTask(taskID)
	if task.Status != domainTask.TaskCancelled {
		t.Errorf("Status = %s, want %s", task.Status, domainTask.TaskCancelled)
	}
}

func TestService_CancelTask_AlreadyFinished(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{"test-plugin": newStubPlugin()}},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	taskID, _ := svc.RunPlugin("test-plugin", nil)
	waitForTask(t, svc, taskID)

	// 已完成的任务不能取消
	err := svc.CancelTask(taskID)
	if err == nil {
		t.Error("expected error for already finished task")
	}
}

func TestService_DeleteTask(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{"test-plugin": newStubPlugin()}},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	taskID, _ := svc.RunPlugin("test-plugin", nil)
	waitForTask(t, svc, taskID)

	err := svc.DeleteTask(taskID)
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	_, err = svc.GetTask(taskID)
	if err == nil {
		t.Error("expected task to be deleted")
	}
}

func TestService_GetTaskOutput(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{"test-plugin": newStubPlugin()}},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	taskID, _ := svc.RunPlugin("test-plugin", nil)
	waitForTask(t, svc, taskID)

	output, err := svc.GetTaskOutput(taskID)
	if err != nil {
		t.Fatalf("GetTaskOutput failed: %v", err)
	}
	if len(output) == 0 {
		t.Error("expected some output")
	}
}

func TestService_GetTaskOutput_NotFound(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{}},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	_, err := svc.GetTaskOutput("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task output")
	}
}

func TestService_RunPlugin_ProcessStartFailure(t *testing.T) {
	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{"test-plugin": newStubPlugin()}},
		&stubFailRunner{},
		&stubParamBuilder{},
		&eventCollector{},
		1000,
	)

	taskID, err := svc.RunPlugin("test-plugin", nil)
	if err != nil {
		t.Fatalf("RunPlugin should not fail in sync path: %v", err)
	}

	waitForTask(t, svc, taskID)

	task, _ := svc.GetTask(taskID)
	if task.Status != domainTask.TaskFailed {
		t.Errorf("expected status '%s', got '%s'", domainTask.TaskFailed, task.Status)
	}
}

func TestService_EventPublished(t *testing.T) {
	collector := &eventCollector{}

	svc := NewService(
		&stubTaskRepo{tasks: make(map[string]*domainTask.Task)},
		&stubPluginRepo{plugins: map[string]*plugin.Plugin{"test-plugin": newStubPlugin()}},
		&stubRunner{started: make(map[string]bool)},
		&stubParamBuilder{},
		collector,
		1000,
	)

	taskID, _ := svc.RunPlugin("test-plugin", nil)
	waitForTask(t, svc, taskID)

	// 应该发出 task.created 事件
	found := false
	collector.mu.Lock()
	for _, e := range collector.events {
		if e.Topic() == "task.created" {
			found = true
			break
		}
	}
	collector.mu.Unlock()
	if !found {
		t.Error("expected task.created event to be published")
	}
}
