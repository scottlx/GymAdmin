package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID               int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserNo           string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"user_no"`
	Name             string         `gorm:"type:varchar(50);not null" json:"name"`
	Gender           int8           `gorm:"type:tinyint" json:"gender"` // 1-男，2-女
	Birthday         *time.Time     `gorm:"type:date" json:"birthday"`
	IDCard           string         `gorm:"type:varchar(18)" json:"id_card"`
	Phone            string         `gorm:"type:varchar(11);uniqueIndex;not null" json:"phone"`
	Email            string         `gorm:"type:varchar(100)" json:"email"`
	AvatarURL        string         `gorm:"type:varchar(255)" json:"avatar_url"`
	Address          string         `gorm:"type:varchar(255)" json:"address"`
	EmergencyContact string         `gorm:"type:varchar(50)" json:"emergency_contact"`
	EmergencyPhone   string         `gorm:"type:varchar(11)" json:"emergency_phone"`
	HealthStatus     string         `gorm:"type:text" json:"health_status"`
	TrainingGoal     string         `gorm:"type:text" json:"training_goal"`
	Source           int8           `gorm:"type:tinyint;default:1" json:"source"`       // 1-前台录入，2-小程序注册，3-美团，4-抖音
	Status           int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-正常，2-冻结，3-黑名单
	Remark           string         `gorm:"type:text" json:"remark"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type UserTrainingStats struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64      `gorm:"uniqueIndex;not null" json:"user_id"`
	TotalDays       int        `gorm:"default:0" json:"total_days"`
	TotalTimes      int        `gorm:"default:0" json:"total_times"`
	ContinuousDays  int        `gorm:"default:0" json:"continuous_days"`
	LastCheckInDate *time.Time `gorm:"type:date" json:"last_check_in_date"`
	MonthTimes      int        `gorm:"default:0" json:"month_times"`
	YearTimes       int        `gorm:"default:0" json:"year_times"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (UserTrainingStats) TableName() string {
	return "user_training_stats"
}

// UserStatusLog records user status changes
type UserStatusLog struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64          `gorm:"index;not null" json:"user_id"`
	OldStatus  int8           `gorm:"type:tinyint" json:"old_status"` // 1-正常，2-冻结，3-黑名单
	NewStatus  int8           `gorm:"type:tinyint" json:"new_status"` // 1-正常，2-冻结，3-黑名单
	Reason     string         `gorm:"type:text" json:"reason"`
	OperatorID *int64         `gorm:"index" json:"operator_id"` // 操作人ID
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserStatusLog) TableName() string {
	return "user_status_logs"
}

// Status constants
const (
	UserStatusActive    = 1 // 正常
	UserStatusFrozen    = 2 // 冻结
	UserStatusBlacklist = 3 // 黑名单
)

// GetStatusText returns the text representation of status
func GetUserStatusText(status int8) string {
	switch status {
	case UserStatusActive:
		return "正常"
	case UserStatusFrozen:
		return "冻结"
	case UserStatusBlacklist:
		return "黑名单"
	default:
		return "未知"
	}
}
