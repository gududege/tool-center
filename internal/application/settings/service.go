package settings

import (
	"github.com/cli-tool-center/tool-center/internal/domain/settings"
)

// Service 配置应用服务
type Service struct {
	repo *Repository
}

// NewService 创建配置服务
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetSettings 获取配置
func (s *Service) GetSettings() (*settings.Settings, error) {
	return s.repo.Get()
}

// SaveSettings 保存配置
func (s *Service) SaveSettings(set *settings.Settings) error {
	return s.repo.Save(set)
}
