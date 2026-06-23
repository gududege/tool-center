package wails

import (
	"testing"
	"time"
)

func TestPluginSummaryDto(t *testing.T) {
	dto := PluginSummaryDto{
		ID:          "test-plugin",
		Name:        "Test Plugin",
		Description: "A test plugin",
		Version:     "1.0.0",
		Navigation: NavigationDto{
			Group: []string{"Tools", "Test"},
			Order: 100,
		},
	}

	if dto.ID != "test-plugin" {
		t.Errorf("ID = %s, want test-plugin", dto.ID)
	}
	if dto.Name != "Test Plugin" {
		t.Errorf("Name = %s, want Test Plugin", dto.Name)
	}
	if len(dto.Navigation.Group) != 2 {
		t.Errorf("expected 2 navigation groups, got %d", len(dto.Navigation.Group))
	}
}

func TestTaskDto_Serialization(t *testing.T) {
	dto := TaskDto{
		ID:        "task-1",
		PluginID:  "test-plugin",
		Status:    "running",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if dto.Status != "running" {
		t.Errorf("Status = %s, want running", dto.Status)
	}
	if dto.PluginID != "test-plugin" {
		t.Errorf("PluginID = %s, want test-plugin", dto.PluginID)
	}
}

func TestOutputEventDto(t *testing.T) {
	dto := OutputEventDto{
		TaskID:    "task-1",
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     "info",
		Source:    "stdout",
		Message:   "hello world",
	}

	if dto.Message != "hello world" {
		t.Errorf("Message = %s, want 'hello world'", dto.Message)
	}
	if dto.Source != "stdout" {
		t.Errorf("Source = %s, want stdout", dto.Source)
	}
}

func TestRunPluginRequest(t *testing.T) {
	req := RunPluginRequest{
		PluginID: "test-plugin",
		FormData: map[string]any{
			"project":   "demo.ap20",
			"overwrite": true,
		},
	}

	if req.PluginID != "test-plugin" {
		t.Errorf("PluginID = %s, want test-plugin", req.PluginID)
	}
	if req.FormData["project"] != "demo.ap20" {
		t.Errorf("FormData.project = %s, want demo.ap20", req.FormData["project"])
	}
}

func TestRunPluginResponse(t *testing.T) {
	resp := RunPluginResponse{TaskID: "task-1"}
	if resp.TaskID != "task-1" {
		t.Errorf("TaskID = %s, want task-1", resp.TaskID)
	}
}

func TestCancelTaskResponse(t *testing.T) {
	resp := CancelTaskResponse{Success: true}
	if !resp.Success {
		t.Error("expected Success to be true")
	}
}

func TestApiErrorDto(t *testing.T) {
	errDto := ApiErrorDto{
		Code:    "PLUGIN_NOT_FOUND",
		Message: "Plugin not found",
		Details: "Plugin 'test' does not exist",
	}

	if errDto.Code != "PLUGIN_NOT_FOUND" {
		t.Errorf("Code = %s, want PLUGIN_NOT_FOUND", errDto.Code)
	}
	if errDto.Message == "" {
		t.Error("Message should not be empty")
	}
}

func TestValidationResultDto(t *testing.T) {
	result := ValidationResultDto{
		Valid: true,
	}
	if !result.Valid {
		t.Error("expected Valid to be true")
	}
}

func TestValidationResultDto_WithErrors(t *testing.T) {
	result := ValidationResultDto{
		Valid: false,
		Errors: []ValidationErrorDto{
			{Path: "/project", Message: "is required"},
		},
	}
	if result.Valid {
		t.Error("expected Valid to be false")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Path != "/project" {
		t.Errorf("Error.Path = %s, want /project", result.Errors[0].Path)
	}
}

func TestSystemInfoDto(t *testing.T) {
	info := SystemInfoDto{
		AppVersion: "1.0.0",
		GoVersion:  "go1.22",
		Os:         "windows",
	}

	if info.AppVersion != "1.0.0" {
		t.Errorf("AppVersion = %s, want 1.0.0", info.AppVersion)
	}
	if info.Os != "windows" {
		t.Errorf("Os = %s, want windows", info.Os)
	}
}

func TestSettingsDto(t *testing.T) {
	dto := SettingsDto{
		Theme:           "dark",
		PluginDirectory: "./plugins",
	}
	if dto.Theme != "dark" {
		t.Errorf("Theme = %s, want dark", dto.Theme)
	}
}

func TestEventTypes(t *testing.T) {
	created := TaskCreatedEvent{TaskID: "t1", PluginID: "p1"}
	if created.TaskID != "t1" || created.PluginID != "p1" {
		t.Error("TaskCreatedEvent fields mismatch")
	}

	failed := TaskFailedEvent{TaskID: "t1", Error: "error message"}
	if failed.Error != "error message" {
		t.Errorf("TaskFailedEvent.Error = %s, want 'error message'", failed.Error)
	}

	changed := TaskStatusChangedEvent{TaskID: "t1", OldStatus: "running", NewStatus: "completed"}
	if changed.NewStatus != "completed" {
		t.Errorf("expected newStatus 'completed', got '%s'", changed.NewStatus)
	}
}

func TestFileFilterDto(t *testing.T) {
	filter := FileFilterDto{
		DisplayName: "Executable",
		Patterns:    []string{"*.exe"},
	}
	if len(filter.Patterns) != 1 || filter.Patterns[0] != "*.exe" {
		t.Errorf("Patterns = %v, want [*.exe]", filter.Patterns)
	}
}

func TestDialogResponse(t *testing.T) {
	resp := DialogResponse{Path: "C:\\test\\file.txt"}
	if resp.Path != "C:\\test\\file.txt" {
		t.Errorf("Path = %s, want C:\\test\\file.txt", resp.Path)
	}

	resp2 := DialogResponse{}
	if resp2.Path != "" {
		t.Error("expected empty path for cancelled dialog")
	}
}
