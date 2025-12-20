package services

import (
	"baihu/internal/constant"
	"baihu/internal/database"
	"baihu/internal/models"
)

type SettingsService struct{}

func NewSettingsService() *SettingsService {
	return &SettingsService{}
}

// InitSettings 初始化默认设置
func (s *SettingsService) InitSettings() error {
	for section, keys := range constant.DefaultSettings {
		for key, value := range keys {
			var count int64
			database.DB.Model(&models.Setting{}).
				Where("section = ? AND key = ?", section, key).
				Count(&count)
			if count == 0 {
				setting := &models.Setting{
					Section: section,
					Key:     key,
					Value:   value,
				}
				if err := database.DB.Create(setting).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Get 获取单个设置
func (s *SettingsService) Get(section, key string) string {
	var setting models.Setting
	if err := database.DB.Where("section = ? AND key = ?", section, key).First(&setting).Error; err != nil {
		if def, ok := constant.DefaultSettings[section][key]; ok {
			return def
		}
		return ""
	}
	return setting.Value
}

// Set 设置单个值
func (s *SettingsService) Set(section, key, value string) error {
	var setting models.Setting
	result := database.DB.Where("section = ? AND key = ?", section, key).First(&setting)
	if result.Error != nil {
		setting = models.Setting{Section: section, Key: key, Value: value}
		return database.DB.Create(&setting).Error
	}
	return database.DB.Model(&setting).Update("value", value).Error
}

// GetSection 获取整个 section 的设置
func (s *SettingsService) GetSection(section string) map[string]string {
	var settings []models.Setting
	database.DB.Where("section = ?", section).Find(&settings)

	result := make(map[string]string)
	// 先填充默认值
	if defaults, ok := constant.DefaultSettings[section]; ok {
		for k, v := range defaults {
			result[k] = v
		}
	}
	// 覆盖数据库值
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result
}

// SetSection 批量设置
func (s *SettingsService) SetSection(section string, values map[string]string) error {
	for key, value := range values {
		if err := s.Set(section, key, value); err != nil {
			return err
		}
	}
	return nil
}
