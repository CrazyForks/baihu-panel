package services

import (
	"baihu/internal/database"
	"baihu/internal/models"
)

type SettingsService struct{}

func NewSettingsService() *SettingsService {
	return &SettingsService{}
}

// Get 获取设置值
func (s *SettingsService) Get(section, key string) string {
	var setting models.Setting
	if err := database.DB.Where("section = ? AND key = ?", section, key).First(&setting).Error; err != nil {
		return ""
	}
	return setting.Value
}

// GetWithDefault 获取设置值，如果不存在则返回默认值
func (s *SettingsService) GetWithDefault(section, key, defaultValue string) string {
	value := s.Get(section, key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Set 设置值
func (s *SettingsService) Set(section, key, value string) error {
	var setting models.Setting
	result := database.DB.Where("section = ? AND key = ?", section, key).First(&setting)
	if result.Error != nil {
		// 不存在则创建
		setting = models.Setting{
			Section: section,
			Key:     key,
			Value:   value,
		}
		return database.DB.Create(&setting).Error
	}
	// 存在则更新
	return database.DB.Model(&setting).Update("value", value).Error
}

// GetBySection 获取某个 section 下的所有设置
func (s *SettingsService) GetBySection(section string) map[string]string {
	var settings []models.Setting
	database.DB.Where("section = ?", section).Find(&settings)

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result
}

// Delete 删除设置
func (s *SettingsService) Delete(section, key string) error {
	return database.DB.Where("section = ? AND key = ?", section, key).Delete(&models.Setting{}).Error
}

// DeleteBySection 删除某个 section 下的所有设置
func (s *SettingsService) DeleteBySection(section string) error {
	return database.DB.Where("section = ?", section).Delete(&models.Setting{}).Error
}
