package models

import (
	"baihu/internal/constant"
)

// Setting 系统设置
type Setting struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Section string `json:"section" gorm:"size:50;not null;index:idx_section_key"`
	Key     string `json:"key" gorm:"column:setting_key;size:100;not null;index:idx_section_key"`
	Value   string `json:"value" gorm:"type:text"`
}

func (Setting) TableName() string {
	return constant.TablePrefix + "settings"
}
