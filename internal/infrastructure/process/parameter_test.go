package process

import (
	"testing"

	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
)

func TestParameterBuilder_Argument(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "file", Kind: plugin.ArgumentKind},
	}
	formData := map[string]any{"file": "project.ap20"}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 1 || args[0] != "project.ap20" {
		t.Errorf("expected [project.ap20], got %v", args)
	}
}

func TestParameterBuilder_ArgumentArray(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "files", Kind: plugin.ArgumentArrayKind},
	}
	formData := map[string]any{"files": []any{"a.txt", "b.txt"}}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 2 || args[0] != "a.txt" || args[1] != "b.txt" {
		t.Errorf("expected [a.txt b.txt], got %v", args)
	}
}

func TestParameterBuilder_Option(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "project", Kind: plugin.OptionKind, Flag: "--project"},
	}
	formData := map[string]any{"project": "demo.ap20"}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 2 || args[0] != "--project" || args[1] != "demo.ap20" {
		t.Errorf("expected [--project demo.ap20], got %v", args)
	}
}

func TestParameterBuilder_Option_SkipEmpty(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "project", Kind: plugin.OptionKind, Flag: "--project"},
	}
	formData := map[string]any{} // no value provided

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 0 {
		t.Errorf("expected empty args for missing value, got %v", args)
	}
}

func TestParameterBuilder_Switch_True(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "overwrite", Kind: plugin.SwitchKind, Flag: "--overwrite"},
	}
	formData := map[string]any{"overwrite": true}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 1 || args[0] != "--overwrite" {
		t.Errorf("expected [--overwrite], got %v", args)
	}
}

func TestParameterBuilder_Switch_False(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "overwrite", Kind: plugin.SwitchKind, Flag: "--overwrite"},
	}
	formData := map[string]any{"overwrite": false}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 0 {
		t.Errorf("expected empty args for false switch, got %v", args)
	}
}

func TestParameterBuilder_BoolOption(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "flag", Kind: plugin.BoolOptionKind, Flag: "--flag"},
	}
	formData := map[string]any{"flag": true}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 2 || args[0] != "--flag" || args[1] != "true" {
		t.Errorf("expected [--flag true], got %v", args)
	}
}

func TestParameterBuilder_DualSwitch_True(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{
			Field:     "overwrite",
			Kind:      plugin.DualSwitchKind,
			TrueFlag:  "--overwrite",
			FalseFlag: "--no-overwrite",
		},
	}
	formData := map[string]any{"overwrite": true}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 1 || args[0] != "--overwrite" {
		t.Errorf("expected [--overwrite], got %v", args)
	}
}

func TestParameterBuilder_DualSwitch_False(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{
			Field:     "overwrite",
			Kind:      plugin.DualSwitchKind,
			TrueFlag:  "--overwrite",
			FalseFlag: "--no-overwrite",
		},
	}
	formData := map[string]any{"overwrite": false}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 1 || args[0] != "--no-overwrite" {
		t.Errorf("expected [--no-overwrite], got %v", args)
	}
}

func TestParameterBuilder_OptionArray_Repeat(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{
			Field: "tags",
			Kind:  plugin.OptionArrayKind,
			Flag:  "--tag",
			Style: "repeat",
		},
	}
	formData := map[string]any{"tags": []any{"PLC1", "PLC2", "PLC3"}}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 6 {
		t.Fatalf("expected 6 args for 3 tags repeat, got %v", args)
	}
	expected := []string{"--tag", "PLC1", "--tag", "PLC2", "--tag", "PLC3"}
	for i, v := range expected {
		if args[i] != v {
			t.Errorf("args[%d] = %s, want %s", i, args[i], v)
		}
	}
}

func TestParameterBuilder_OptionArray_Join(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{
			Field:     "tags",
			Kind:      plugin.OptionArrayKind,
			Flag:      "--tags",
			Style:     "join",
			Separator: ",",
		},
	}
	formData := map[string]any{"tags": []any{"PLC1", "PLC2", "PLC3"}}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 2 || args[0] != "--tags" || args[1] != "PLC1,PLC2,PLC3" {
		t.Errorf("expected [--tags PLC1,PLC2,PLC3], got %v", args)
	}
}

func TestParameterBuilder_OptionArray_Equals(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{
			Field:     "tags",
			Kind:      plugin.OptionArrayKind,
			Flag:      "--tags",
			Style:     "equals",
			Separator: ",",
		},
	}
	formData := map[string]any{"tags": []any{"PLC1", "PLC2"}}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 1 || args[0] != "--tags=PLC1,PLC2" {
		t.Errorf("expected [--tags=PLC1,PLC2], got %v", args)
	}
}

func TestParameterBuilder_DefaultValue(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{
			Field:        "project",
			Kind:         plugin.OptionKind,
			Flag:         "--project",
			DefaultValue: "default.ap20",
		},
	}
	formData := map[string]any{} // 未提供值

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 2 || args[0] != "--project" || args[1] != "default.ap20" {
		t.Errorf("expected [--project default.ap20], got %v", args)
	}
}

func TestParameterBuilder_MixedParameters(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "file", Kind: plugin.ArgumentKind},
		{Field: "project", Kind: plugin.OptionKind, Flag: "--project"},
		{Field: "overwrite", Kind: plugin.SwitchKind, Flag: "--overwrite"},
	}
	formData := map[string]any{
		"file":      "output.txt",
		"project":   "demo.ap20",
		"overwrite": true,
	}

	args, err := b.Build(mappings, formData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 顺序应保持定义顺序
	expected := []string{"output.txt", "--project", "demo.ap20", "--overwrite"}
	if len(args) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, args)
	}
	for i, v := range expected {
		if args[i] != v {
			t.Errorf("args[%d] = %s, want %s", i, args[i], v)
		}
	}
}

func TestParameterBuilder_UnknownKind(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "x", Kind: "unknown-kind"},
	}
	formData := map[string]any{"x": "value"}

	_, err := b.Build(mappings, formData)
	if err == nil {
		t.Error("expected error for unknown parameter kind")
	}
}

func TestParameterBuilder_ArgumentArray_InvalidType(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "files", Kind: plugin.ArgumentArrayKind},
	}
	// 传入字符串而非数组
	formData := map[string]any{"files": "not-an-array"}

	_, err := b.Build(mappings, formData)
	if err == nil {
		t.Error("expected error for non-array value")
	}
}

func TestParameterBuilder_EmptyMappings(t *testing.T) {
	b := NewParameterBuilder()
	args, err := b.Build(nil, map[string]any{"x": "y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected empty args, got %v", args)
	}
}

func TestParameterBuilder_EmptyFormData(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "project", Kind: plugin.OptionKind, Flag: "--project"},
	}
	args, err := b.Build(mappings, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected empty args for nil formData, got %v", args)
	}
}

func TestParameterBuilder_Switch_StringTrue(t *testing.T) {
	b := NewParameterBuilder()
	mappings := []plugin.ParameterMapping{
		{Field: "verbose", Kind: plugin.SwitchKind, Flag: "--verbose"},
	}

	// 字符串 "true" 也应被视为 true
	args, err := b.Build(mappings, map[string]any{"verbose": "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 1 || args[0] != "--verbose" {
		t.Errorf("expected [--verbose], got %v", args)
	}

	// 字符串 "false" 应被视为 false
	args, err = b.Build(mappings, map[string]any{"verbose": "false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected empty args, got %v", args)
	}
}
