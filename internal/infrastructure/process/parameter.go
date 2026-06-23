package process

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/cli-tool-center/tool-center/internal/domain/plugin"
)

// ParameterBuilder 参数构建器，将 FormData 转换为 CLI 参数
type ParameterBuilder struct{}

// NewParameterBuilder 创建参数构建器
func NewParameterBuilder() *ParameterBuilder {
	return &ParameterBuilder{}
}

// Build 根据参数映射和表单数据生成 CLI 参数列表
func (b *ParameterBuilder) Build(mappings []plugin.ParameterMapping, formData map[string]any) ([]string, error) {
	var args []string

	for _, m := range mappings {
		val, exists := formData[m.Field]
		if !exists {
			// 使用默认值
			if m.DefaultValue != nil {
				val = m.DefaultValue
			} else {
				continue
			}
		}

		built, err := b.buildParam(m, val)
		if err != nil {
			return nil, fmt.Errorf("build param %s: %w", m.Field, err)
		}
		args = append(args, built...)
	}
	slog.Debug("built arguments", "args", args)
	return args, nil
}

func (b *ParameterBuilder) buildParam(m plugin.ParameterMapping, val any) ([]string, error) {
	switch m.Kind {
	case plugin.ArgumentKind:
		return b.buildArgument(m, val)
	case plugin.ArgumentArrayKind:
		return b.buildArgumentArray(m, val)
	case plugin.OptionKind:
		return b.buildOption(m, val)
	case plugin.OptionArrayKind:
		return b.buildOptionArray(m, val)
	case plugin.SwitchKind:
		return b.buildSwitch(m, val)
	case plugin.BoolOptionKind:
		return b.buildBoolOption(m, val)
	case plugin.DualSwitchKind:
		return b.buildDualSwitch(m, val)
	default:
		return nil, fmt.Errorf("unknown parameter kind: %s", m.Kind)
	}
}

func (b *ParameterBuilder) buildArgument(m plugin.ParameterMapping, val any) ([]string, error) {
	s := fmt.Sprint(val)
	if s == "" {
		return nil, nil
	}
	return []string{s}, nil
}

func (b *ParameterBuilder) buildArgumentArray(m plugin.ParameterMapping, val any) ([]string, error) {
	arr, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("expected array for argument-array field %s", m.Field)
	}
	var args []string
	for _, v := range arr {
		if s := fmt.Sprint(v); s != "" {
			args = append(args, s)
		}
	}
	return args, nil
}

func (b *ParameterBuilder) buildOption(m plugin.ParameterMapping, val any) ([]string, error) {
	s := fmt.Sprint(val)
	if s == "" {
		return nil, nil
	}
	return []string{m.Flag, s}, nil
}

func (b *ParameterBuilder) buildOptionArray(m plugin.ParameterMapping, val any) ([]string, error) {
	arr, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("expected array for option-array field %s", m.Field)
	}
	var args []string
	switch m.Style {
	case "repeat":
		for _, v := range arr {
			if s := fmt.Sprint(v); s != "" {
				args = append(args, m.Flag, s)
			}
		}
	case "join":
		var parts []string
		for _, v := range arr {
			parts = append(parts, fmt.Sprint(v))
		}
		sep := m.Separator
		if sep == "" {
			sep = ","
		}
		args = append(args, m.Flag, strings.Join(parts, sep))
	case "equals":
		var parts []string
		for _, v := range arr {
			parts = append(parts, fmt.Sprint(v))
		}
		sep := m.Separator
		if sep == "" {
			sep = ","
		}
		args = append(args, fmt.Sprintf("%s=%s", m.Flag, strings.Join(parts, sep)))
	default:
		// 默认 repeat
		for _, v := range arr {
			if s := fmt.Sprint(v); s != "" {
				args = append(args, m.Flag, s)
			}
		}
	}
	return args, nil
}

func (b *ParameterBuilder) buildSwitch(m plugin.ParameterMapping, val any) ([]string, error) {
	flagVal, ok := val.(bool)
	if !ok {
		// 尝试字符串
		s := fmt.Sprint(val)
		flagVal = s == "true" || s == "1"
	}
	if !flagVal {
		return nil, nil
	}
	return []string{m.Flag}, nil
}

func (b *ParameterBuilder) buildBoolOption(m plugin.ParameterMapping, val any) ([]string, error) {
	flagVal, ok := val.(bool)
	if !ok {
		s := fmt.Sprint(val)
		flagVal = s == "true" || s == "1"
	}
	return []string{m.Flag, fmt.Sprint(flagVal)}, nil
}

func (b *ParameterBuilder) buildDualSwitch(m plugin.ParameterMapping, val any) ([]string, error) {
	flagVal, ok := val.(bool)
	if !ok {
		s := fmt.Sprint(val)
		flagVal = s == "true" || s == "1"
	}
	if flagVal {
		return []string{m.TrueFlag}, nil
	}
	if m.FalseFlag != "" {
		return []string{m.FalseFlag}, nil
	}
	return nil, nil
}
