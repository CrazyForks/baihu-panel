package services

import (
	"baihu/internal/logger"
)

const (
	InitSection = "system"
	InitKey     = "initialized"
	InitValue   = "true"
)

type InitService struct {
	settingsService *SettingsService
	userService     *UserService
}

func NewInitService(settingsService *SettingsService, userService *UserService) *InitService {
	return &InitService{
		settingsService: settingsService,
		userService:     userService,
	}
}

// Initialize 执行初始化，如果已初始化则跳过
func (s *InitService) Initialize() {
	//if s.IsInitialized() {
	//	logger.Info("系统已初始化，跳过")
	//	return
	//}

	logger.Info("开始初始化系统...")

	// 创建管理员账号
	s.createAdminUser()

	// 初始化默认设置（每次启动都检查）
	if err := s.settingsService.InitSettings(); err != nil {
		logger.Warnf("初始化设置失败: %v", err)
	}

	//// 标记为已初始化
	//s.settingsService.Set(InitSection, InitKey, InitValue)
	//logger.Info("系统初始化完成")
}

// IsInitialized 检查是否已初始化
func (s *InitService) IsInitialized() bool {
	return s.settingsService.Get(InitSection, InitKey) == InitValue
}

// createAdminUser 创建管理员账号
func (s *InitService) createAdminUser() {
	existingUser := s.userService.GetUserByUsername("admin")
	if existingUser != nil {
		logger.Info("管理员账号已存在，跳过创建")
		return
	}

	s.userService.CreateUser("admin", "123456", "admin@local", "admin")
	logger.Info("管理员账号创建成功: admin / 123456")
}
