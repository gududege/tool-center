package wails

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/cli-tool-center/tool-center/internal/application/task"
	domainTask "github.com/cli-tool-center/tool-center/internal/domain/task"
)

// TaskApi 任务相关 API
type TaskApi struct {
	service *task.Service
}

// NewTaskApi 创建任务 API
func NewTaskApi(service *task.Service) *TaskApi {
	return &TaskApi{service: service}
}

// RunPlugin 执行插件
func (a *TaskApi) RunPlugin(req RunPluginRequest) (*RunPluginResponse, error) {
	slog.Info("RunPlugin called", "pluginId", req.PluginID, "formData", req.FormData)
	taskID, err := a.service.RunPlugin(req.PluginID, req.FormData)
	if err != nil {
		slog.Error("RunPlugin failed", "pluginId", req.PluginID, "error", err)
		return nil, fmt.Errorf("run plugin: %w", err)
	}
	slog.Info("RunPlugin started", "pluginId", req.PluginID, "taskId", taskID)
	return &RunPluginResponse{TaskID: taskID}, nil
}

// GetTask 获取单个任务
func (a *TaskApi) GetTask(taskID string) (*TaskDto, error) {
	t, err := a.service.GetTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	dto := taskToDto(t)
	return &dto, nil
}

// GetTasks 获取所有任务
func (a *TaskApi) GetTasks() ([]TaskDto, error) {
	tasks, err := a.service.ListTasks()
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}

	dtos := make([]TaskDto, 0, len(tasks))
	for _, t := range tasks {
		dtos = append(dtos, taskToDto(t))
	}
	return dtos, nil
}

// CancelTask 取消任务
func (a *TaskApi) CancelTask(taskID string) (*CancelTaskResponse, error) {
	if err := a.service.CancelTask(taskID); err != nil {
		return &CancelTaskResponse{Success: false}, fmt.Errorf("cancel task: %w", err)
	}
	return &CancelTaskResponse{Success: true}, nil
}

// DeleteTask 删除任务
func (a *TaskApi) DeleteTask(taskID string) (*DeleteTaskResponse, error) {
	if err := a.service.DeleteTask(taskID); err != nil {
		return &DeleteTaskResponse{Success: false}, fmt.Errorf("delete task: %w", err)
	}
	return &DeleteTaskResponse{Success: true}, nil
}

// GetTaskOutput 获取任务输出
func (a *TaskApi) GetTaskOutput(taskID string) ([]OutputEventDto, error) {
	events, err := a.service.GetTaskOutput(taskID)
	if err != nil {
		return nil, fmt.Errorf("get task output: %w", err)
	}

	dtos := make([]OutputEventDto, 0, len(events))
	for _, e := range events {
		dtos = append(dtos, OutputEventDto{
			TaskID:    e.TaskID,
			Timestamp: e.Timestamp.Format(time.RFC3339Nano),
			Level:     string(e.Level),
			Source:    string(e.Source),
			Message:   e.Message,
		})
	}
	return dtos, nil
}

func taskToDto(t *domainTask.Task) TaskDto {
	dto := TaskDto{
		ID:        t.ID,
		PluginID:  t.PluginID,
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
	}
	if t.StartedAt != nil {
		dto.StartedAt = t.StartedAt.Format(time.RFC3339)
	}
	if t.EndedAt != nil {
		dto.EndedAt = t.EndedAt.Format(time.RFC3339)
	}
	return dto
}
