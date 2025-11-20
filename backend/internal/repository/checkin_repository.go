package repository

import (
	"gym-admin/internal/models"
	"gym-admin/pkg/database"
	"time"

	"gorm.io/gorm"
)

type CheckInRepository struct {
	db *gorm.DB
}

func NewCheckInRepository() *CheckInRepository {
	return &CheckInRepository{db: database.GetDB()}
}

func (r *CheckInRepository) Create(checkIn *models.CheckIn) error {
	return r.db.Create(checkIn).Error
}

func (r *CheckInRepository) GetByID(id int64) (*models.CheckIn, error) {
	var checkIn models.CheckIn
	err := r.db.First(&checkIn, id).Error
	return &checkIn, err
}

func (r *CheckInRepository) GetTodayCheckIn(userID int64, startTime, endTime time.Time) (*models.CheckIn, error) {
	var checkIn models.CheckIn
	err := r.db.Where("user_id = ? AND check_in_time >= ? AND check_in_time < ?",
		userID, startTime, endTime).First(&checkIn).Error
	if err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (r *CheckInRepository) List(page, pageSize int, userID *int64, startDate, endDate *time.Time) ([]models.CheckIn, int64, error) {
	var checkIns []models.CheckIn
	var total int64

	query := r.db.Model(&models.CheckIn{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if startDate != nil {
		query = query.Where("check_in_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("check_in_time <= ?", *endDate)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("check_in_time DESC").Offset(offset).Limit(pageSize).Find(&checkIns).Error

	return checkIns, total, err
}

func (r *CheckInRepository) GetUserCheckInCount(userID int64, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.CheckIn{}).
		Where("user_id = ? AND check_in_time >= ? AND check_in_time <= ?", userID, startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetAllByUserID gets all check-ins for a user ordered by check-in time
func (r *CheckInRepository) GetAllByUserID(userID int64) ([]models.CheckIn, error) {
	var checkIns []models.CheckIn
	err := r.db.Where("user_id = ?", userID).Order("check_in_time ASC").Find(&checkIns).Error
	return checkIns, err
}

// GetCheckInsByDateRange gets check-ins for a user within a date range
func (r *CheckInRepository) GetCheckInsByDateRange(userID int64, startDate, endDate time.Time) ([]models.CheckIn, error) {
	var checkIns []models.CheckIn
	err := r.db.Where("user_id = ? AND check_in_time >= ? AND check_in_time <= ?",
		userID, startDate, endDate).Order("check_in_time ASC").Find(&checkIns).Error
	return checkIns, err
}

// GetMonthlyCheckInDays gets unique check-in days count for a month
func (r *CheckInRepository) GetMonthlyCheckInDays(userID int64, year int, month int) (int, error) {
	var count int64
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	err := r.db.Model(&models.CheckIn{}).
		Select("COUNT(DISTINCT DATE(check_in_time))").
		Where("user_id = ? AND check_in_time >= ? AND check_in_time <= ?", userID, startDate, endDate).
		Count(&count).Error

	return int(count), err
}

// GetYearlyCheckInDays gets unique check-in days count for a year
func (r *CheckInRepository) GetYearlyCheckInDays(userID int64, year int) (int, error) {
	var count int64
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.Local)

	err := r.db.Model(&models.CheckIn{}).
		Select("COUNT(DISTINCT DATE(check_in_time))").
		Where("user_id = ? AND check_in_time >= ? AND check_in_time <= ?", userID, startDate, endDate).
		Count(&count).Error

	return int(count), err
}

// GetDetailedStats retrieves detailed training statistics for a user for a given year
func (r *CheckInRepository) GetDetailedStats(userID int64, year int) (map[string]interface{}, error) {
	var results []struct {
		Month int `json:"month"`
		Count int `json:"count"`
	}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.Local)

	err := r.db.Model(&models.CheckIn{}).
		Select("MONTH(check_in_time) as month, COUNT(*) as count").
		Where("user_id = ? AND check_in_time >= ? AND check_in_time <= ?", userID, startDate, endDate).
		Group("MONTH(check_in_time)").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	monthlyData := make(map[int]int)
	for _, result := range results {
		monthlyData[result.Month] = result.Count
	}

	// Ensure all months are present
	for i := 1; i <= 12; i++ {
		if _, ok := monthlyData[i]; !ok {
			monthlyData[i] = 0
		}
	}

	return map[string]interface{}{"monthly_check_ins": monthlyData}, nil
}

// GetCheckInCalendar retrieves the check-in calendar for a user for a given month
func (r *CheckInRepository) GetCheckInCalendar(userID int64, year int, month int) ([]int, error) {
	var days []int
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	err := r.db.Model(&models.CheckIn{}).
		Select("DISTINCT DAY(check_in_time)").
		Where("user_id = ? AND check_in_time >= ? AND check_in_time <= ?", userID, startDate, endDate).
		Pluck("DAY(check_in_time)", &days).Error

	return days, err
}
