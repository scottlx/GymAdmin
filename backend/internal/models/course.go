package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CoachID       int64          `gorm:"index;not null" json:"coach_id"`
	CourseName    string         `gorm:"type:varchar(100);not null" json:"course_name"`
	CourseType    int8           `gorm:"type:tinyint;not null" json:"course_type"` // 1-私教课，2-团课
	StartTime     time.Time      `gorm:"not null;index" json:"start_time"`
	EndTime       time.Time      `gorm:"not null" json:"end_time"`
	MaxCapacity   int            `gorm:"default:1" json:"max_capacity"`
	CurrentCount  int            `gorm:"default:0" json:"current_count"`
	Price         float64        `gorm:"type:decimal(10,2)" json:"price"`
	Status        int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-可预约，2-已满员，3-已取消，4-已完成
	AverageRating float64        `gorm:"type:decimal(3,2);default:0.0" json:"average_rating"`
	RatingCount   int            `gorm:"default:0" json:"rating_count"`
	Description   string         `gorm:"type:text" json:"description"`
	Remark        string         `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Course) TableName() string {
	return "courses"
}

const (
	CourseStatusAvailable = 1
	CourseStatusFull      = 2
	CourseStatusCancelled = 3
	CourseStatusCompleted = 4
)

type Booking struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64          `gorm:"index;not null" json:"user_id"`
	CourseID    int64          `gorm:"index;not null" json:"course_id"`
	Status      int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-已预约，2-已取消，3-已完成，4-缺席
	HasRated    bool           `gorm:"default:false" json:"has_rated"`
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

const (
	BookingStatusBooked    = 1
	BookingStatusCancelled = 2
	BookingStatusCompleted = 3
	BookingStatusAbsent    = 4
)

// CoachAvailability represents a coach's available time slots
type CoachAvailability struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CoachID      int64          `gorm:"index;not null" json:"coach_id"`
	DayOfWeek    int            `gorm:"type:tinyint;not null" json:"day_of_week"`   // 0 = Sunday, 1 = Monday, ...
	StartTime    string         `gorm:"type:varchar(5);not null" json:"start_time"` // "HH:MM"
	EndTime      string         `gorm:"type:varchar(5);not null" json:"end_time"`   // "HH:MM"
	IsRecurring  bool           `gorm:"default:true" json:"is_recurring"`           // Is this a recurring availability?
	SpecificDate *time.Time     `gorm:"type:date" json:"specific_date"`             // For non-recurring availability
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CoachAvailability) TableName() string {
	return "coach_availability"
}

// CoachLeave represents a coach's leave/vacation time
type CoachLeave struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CoachID    int64          `gorm:"index;not null" json:"coach_id"`
	StartTime  time.Time      `gorm:"not null;index"`
	EndTime    time.Time      `gorm:"not null;index"`
	Reason     string         `gorm:"type:varchar(255)" json:"reason"`
	Status     int8           `gorm:"type:tinyint;default:1"` // 1-Pending, 2-Approved, 3-Rejected
	ApproverID *int64         `json:"approver_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CoachLeave) TableName() string {
	return "coach_leaves"
}

const (
	LeaveStatusPending  = 1
	LeaveStatusApproved = 2
	LeaveStatusRejected = 3
)

// CourseRating represents a user's rating for a course
type CourseRating struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64          `gorm:"index;not null" json:"user_id"`
	CourseID  int64          `gorm:"index;not null" json:"course_id"`
	BookingID int64          `gorm:"uniqueIndex;not null" json:"booking_id"`
	Rating    int8           `gorm:"type:tinyint;not null"` // 1-5
	Comment   string         `gorm:"type:text" json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CourseRating) TableName() string {
	return "course_ratings"
}
