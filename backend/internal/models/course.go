package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CoachID      int64          `gorm:"index;not null" json:"coach_id"`
	CourseName   string         `gorm:"type:varchar(100);not null" json:"course_name"`
	CourseType   int8           `gorm:"type:tinyint;not null" json:"course_type"` // 1-私教课，2-团课
	StartTime    time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime      time.Time      `gorm:"not null" json:"end_time"`
	MaxCapacity  int            `gorm:"default:1" json:"max_capacity"`
	CurrentCount int            `gorm:"default:0" json:"current_count"`
	Price        float64        `gorm:"type:decimal(10,2)" json:"price"`
	Status       int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-可预约，2-已满员，3-已取消，4-已完成
	Description  string         `gorm:"type:text" json:"description"`
	Remark       string         `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Course) TableName() string {
	return "courses"
}

type Booking struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64          `gorm:"index;not null" json:"user_id"`
	CourseID    int64          `gorm:"index;not null" json:"course_id"`
	Status      int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-已预约，2-已取消，3-已完成，4-缺席
	BookedAt    time.Time      `gorm:"not null" json:"booked_at"`
	CancelledAt *time.Time     `json:"cancelled_at"`
	CheckedInAt *time.Time     `json:"checked_in_at"`
	Remark      string         `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Booking) TableName() string {
	return "bookings"
}
