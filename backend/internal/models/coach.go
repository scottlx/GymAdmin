package models

import (
	"time"

	"gorm.io/gorm"
)

type Coach struct {
	ID                  int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CoachNo             string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"coach_no"`
	Name                string         `gorm:"type:varchar(50);not null" json:"name"`
	Gender              int8           `gorm:"type:tinyint" json:"gender"`
	Phone               string         `gorm:"type:varchar(11);uniqueIndex;not null" json:"phone"`
	Email               string         `gorm:"type:varchar(100)" json:"email"`
	AvatarURL           string         `gorm:"type:varchar(255)" json:"avatar_url"`
	Specialties         string         `gorm:"type:text" json:"specialties"`    // JSON格式
	Certifications      string         `gorm:"type:text" json:"certifications"` // JSON格式
	Experience          int            `gorm:"default:0" json:"experience"`     // 工作年限
	Introduction        string         `gorm:"type:text" json:"introduction"`
	HourlyRate          float64        `gorm:"type:decimal(10,2)" json:"hourly_rate"`
	Status              int8           `gorm:"type:tinyint;default:1;index" json:"status"`               // 1-在职，2-离职
	CertificationStatus int8           `gorm:"type:tinyint;default:1;index" json:"certification_status"` // 1-未认证, 2-待审核, 3-已认证
	HireDate            *time.Time     `gorm:"type:date" json:"hire_date"`
	Remark              string         `gorm:"type:text" json:"remark"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Coach) TableName() string {
	return "coaches"
}

const (
	CoachCertificationStatusUnverified = 1
	CoachCertificationStatusPending    = 2
	CoachCertificationStatusVerified   = 3
)

// CoachCertification represents a coach's certification
type CoachCertification struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CoachID           int64          `gorm:"index;not null" json:"coach_id"`
	CertificationName string         `gorm:"type:varchar(100);not null" json:"certification_name"`
	IssuingOrg        string         `gorm:"type:varchar(100)" json:"issuing_org"`
	IssueDate         *time.Time     `gorm:"type:date" json:"issue_date"`
	ExpiryDate        *time.Time     `gorm:"type:date" json:"expiry_date"`
	FileURL           string         `gorm:"type:varchar(500)" json:"file_url"`
	Status            int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-待审核, 2-已通过, 3-未通过
	ReviewerID        *int64         `json:"reviewer_id"`
	ReviewedAt        *time.Time     `json:"reviewed_at"`
	ReviewNotes       string         `gorm:"type:text" json:"review_notes"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CoachCertification) TableName() string {
	return "coach_certifications"
}

const (
	CertificationStatusPending  = 1
	CertificationStatusApproved = 2
	CertificationStatusRejected = 3
)
