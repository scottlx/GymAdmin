package models

import (
	"time"

	"gorm.io/gorm"
)

type Coach struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CoachNo        string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"coach_no"`
	Name           string         `gorm:"type:varchar(50);not null" json:"name"`
	Gender         int8           `gorm:"type:tinyint" json:"gender"`
	Phone          string         `gorm:"type:varchar(11);uniqueIndex;not null" json:"phone"`
	Email          string         `gorm:"type:varchar(100)" json:"email"`
	AvatarURL      string         `gorm:"type:varchar(255)" json:"avatar_url"`
	Specialties    string         `gorm:"type:text" json:"specialties"`    // JSON格式
	Certifications string         `gorm:"type:text" json:"certifications"` // JSON格式
	Experience     int            `gorm:"default:0" json:"experience"`     // 工作年限
	Introduction   string         `gorm:"type:text" json:"introduction"`
	HourlyRate     float64        `gorm:"type:decimal(10,2)" json:"hourly_rate"`
	Status         int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-在职，2-离职
	HireDate       *time.Time     `gorm:"type:date" json:"hire_date"`
	Remark         string         `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Coach) TableName() string {
	return "coaches"
}
