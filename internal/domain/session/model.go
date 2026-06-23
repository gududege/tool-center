package session

import "time"

// FormSession 表单会话（Tab）
type FormSession struct {
	ID        string
	PluginID  string
	FormData  map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SessionRegistry 会话注册表
type SessionRegistry struct {
	Sessions map[string]*FormSession
}

// NewSessionRegistry 创建会话注册表
func NewSessionRegistry() *SessionRegistry {
	return &SessionRegistry{
		Sessions: make(map[string]*FormSession),
	}
}
