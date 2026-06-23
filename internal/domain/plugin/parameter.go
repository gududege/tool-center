package plugin

// ParameterKind 参数映射类型
type ParameterKind string

const (
	ArgumentKind      ParameterKind = "argument"
	ArgumentArrayKind ParameterKind = "argument-array"
	OptionKind        ParameterKind = "option"
	OptionArrayKind   ParameterKind = "option-array"
	SwitchKind        ParameterKind = "switch"
	BoolOptionKind    ParameterKind = "bool-option"
	DualSwitchKind    ParameterKind = "dual-switch"
)

// ParameterMapping 表单字段到 CLI 参数的映射
type ParameterMapping struct {
	Field        string
	Kind         ParameterKind
	Flag         string
	Style        string // repeat | join | equals
	Separator    string
	TrueFlag     string
	FalseFlag    string
	DefaultValue any
}
