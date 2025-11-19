package models

import (
	"time"

	"gorm.io/gorm"
)

type CheckIn struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64          `gorm:"index;not null" json:"user_id"`
	CardID      *int64         `gorm:"index" json:"card_id"`
	CheckInType int8           `gorm:"type:tinyint;not null" json:"check_in_type"` // 1-人脸识别，2-刷卡，3-手动签到
	CheckInTime time.Time      `gorm:"not null;index" json:"check_in_time"`
	DeviceID    string         `gorm:"type:varchar(50)" json:"device_id"`
	Remark      string         `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CheckIn) TableName() string {
	return "check_ins"
}

type FaceRecord struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64          `gorm:"index;not null" json:"user_id"`
	FaceImageURL string         `gorm:"type:varchar(255);not null" json:"face_image_url"`
	FaceFeature  string         `gorm:"type:text" json:"face_feature"`              // 人脸特征向量
	Status       int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-启用，2-停用
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FaceRecord) TableName() string {
	return "face_records"
}
