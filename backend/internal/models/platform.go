package models

import (
"time"

"gorm.io/gorm"
)

// ThirdPartyPlatform represents a third-party platform configuration
type ThirdPartyPlatform struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Platform    int8           `gorm:"type:tinyint;uniqueIndex;not null" json:"platform"` // 1-美团，2-抖音
	AppName     string         `gorm:"type:varchar(100);not null" json:"app_name"`
	AppID       string         `gorm:"type:varchar(100)" json:"app_id"`
	AppSecret   string         `gorm:"type:varchar(255)" json:"app_secret"`
	CallbackURL string         `gorm:"type:varchar(500)" json:"callback_url"`
	Config      string         `gorm:"type:text" json:"config"` // JSON格式的额外配置
	Status      int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-启用，2-禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ThirdPartyPlatform) TableName() string {
	return "third_party_platforms"
}
