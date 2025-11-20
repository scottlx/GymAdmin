package models

import (
"time"

"gorm.io/gorm"
)

// CoachPerformance represents a coach's monthly performance statistics
type CoachPerformance struct {
ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
CoachID       int64          `gorm:"index;not null" json:"coach_id"`
Year          int            `gorm:"index;not null" json:"year"`
Month         int            `gorm:"index;not null" json:"month"`
TotalCourses  int            `gorm:"default:0" json:"total_courses"`
TotalHours    float64        `gorm:"type:decimal(10,2);default:0.0" json:"total_hours"`
TotalIncome   float64        `gorm:"type:decimal(10,2);default:0.0" json:"total_income"`
AverageRating float64        `gorm:"type:decimal(3,2);default:0.0" json:"average_rating"`
CreatedAt     time.Time      `json:"created_at"`
UpdatedAt     time.Time      `json:"updated_at"`
DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CoachPerformance) TableName() string {
return "coach_performance"
}
