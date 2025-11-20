package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification 通知消息
type Notification struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64          `gorm:"index;not null" json:"user_id"`
	Type      int8           `gorm:"type:tinyint;not null;index" json:"type"` // 1-会员卡到期, 2-课程提醒, 3-系统通知
	Title     string         `gorm:"type:varchar(100);not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	RelatedID *int64         `json:"related_id"` // 关联ID（如会员卡ID）
	IsRead    int8           `gorm:"type:tinyint;default:0;index" json:"is_read"`
	ReadAt    *time.Time     `json:"read_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}
