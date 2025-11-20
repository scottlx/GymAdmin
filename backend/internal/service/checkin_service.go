package service

import (
	"errors"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
	"time"
)

type CheckInService struct {
	checkinRepo *repository.CheckInRepository
	userRepo    *repository.UserRepository
	cardRepo    *repository.CardRepository
}

func NewCheckInService() *CheckInService {
	return &CheckInService{
		checkinRepo: repository.NewCheckInRepository(),
		userRepo:    repository.NewUserRepository(),
		cardRepo:    repository.NewCardRepository(),
	}
}

// CheckIn performs user check-in
func (s *CheckInService) CheckIn(userID int64, checkInType int8, deviceID string, cardID *int64) (*models.CheckIn, error) {
	// Verify user exists and is active
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Status != 1 {
		return nil, errors.New("user is not active")
	}

	// Check if user has valid membership card
	if cardID != nil {
		card, err := s.cardRepo.GetByID(*cardID)
		if err != nil {
			return nil, errors.New("membership card not found")
		}
		if card.Status != 1 {
			return nil, errors.New("membership card is not active")
		}
		if card.EndDate.Before(time.Now()) {
			return nil, errors.New("membership card has expired")
		}
		if card.IsFrozen == 1 {
			return nil, errors.New("membership card is frozen")
		}
	}

	// Check if already checked in today
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)
	todayEnd := todayStart.Add(24 * time.Hour)

	existingCheckIn, _ := s.checkinRepo.GetTodayCheckIn(userID, todayStart, todayEnd)
	if existingCheckIn != nil {
		return nil, errors.New("already checked in today")
	}

	// Create check-in record
	checkIn := &models.CheckIn{
		UserID:      userID,
		CardID:      cardID,
		CheckInType: checkInType,
		CheckInTime: time.Now(),
		DeviceID:    deviceID,
	}

	if err := s.checkinRepo.Create(checkIn); err != nil {
		return nil, err
	}

	// Update training statistics
	if err := s.UpdateTrainingStats(userID); err != nil {
		// Log error but don't fail the check-in
		// TODO: Add proper logging
	}

	return checkIn, nil
}

// GetTodayCheckIn gets today's check-in record for a user
func (s *CheckInService) GetTodayCheckIn(userID int64) (*models.CheckIn, error) {
	today := time.Now().Format("2006-01-02")
	todayStart, _ := time.Parse("2006-01-02", today)
	todayEnd := todayStart.Add(24 * time.Hour)

	return s.checkinRepo.GetTodayCheckIn(userID, todayStart, todayEnd)
}

// ListCheckIns lists check-in records with pagination
func (s *CheckInService) ListCheckIns(page, pageSize int, userID *int64, startDate, endDate *time.Time) ([]models.CheckIn, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.checkinRepo.List(page, pageSize, userID, startDate, endDate)
}

// UpdateTrainingStats updates user training statistics after check-in
func (s *CheckInService) UpdateTrainingStats(userID int64) error {
	stats, err := s.userRepo.GetStats(userID)
	if err != nil {
		// Create new stats if not exists
		stats = &models.UserTrainingStats{
			UserID:     userID,
			TotalDays:  0,
			TotalTimes: 0,
		}
	}

	// Update statistics
	stats.TotalTimes++

	// Check if this is a new training day
	today := time.Now()
	if stats.LastCheckInDate == nil || !isSameDay(*stats.LastCheckInDate, today) {
		stats.TotalDays++

		// Update continuous days
		if stats.LastCheckInDate != nil && isConsecutiveDay(*stats.LastCheckInDate, today) {
			stats.ContinuousDays++
		} else {
			stats.ContinuousDays = 1
		}
	}

	stats.LastCheckInDate = &today

	// Update month and year times
	if stats.LastCheckInDate == nil || stats.LastCheckInDate.Month() != today.Month() {
		stats.MonthTimes = 1
	} else {
		stats.MonthTimes++
	}

	if stats.LastCheckInDate == nil || stats.LastCheckInDate.Year() != today.Year() {
		stats.YearTimes = 1
	} else {
		stats.YearTimes++
	}

	return s.userRepo.UpdateStats(stats)
}

// Helper functions
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func isConsecutiveDay(t1, t2 time.Time) bool {
	diff := t2.Sub(t1)
	return diff >= 23*time.Hour && diff <= 25*time.Hour
}

// GetUserStats retrieves user training statistics
func (s *CheckInService) GetUserStats(userID int64) (*models.UserTrainingStats, error) {
	return s.userRepo.GetStats(userID)
}

// GetDetailedStats retrieves detailed training statistics for a user
func (s *CheckInService) GetDetailedStats(userID int64, year int) (map[string]interface{}, error) {
	return s.checkinRepo.GetDetailedStats(userID, year)
}

// GetCheckInCalendar retrieves the check-in calendar for a user
func (s *CheckInService) GetCheckInCalendar(userID int64, year int, month int) ([]int, error) {
	return s.checkinRepo.GetCheckInCalendar(userID, year, month)
}

// RecalculateUserStats recalculates user training statistics from scratch
func (s *CheckInService) RecalculateUserStats(userID int64) error {
	// This is a more complex operation, for now, we just reset and update
	// A more complete implementation would re-process all check-ins
	stats, err := s.userRepo.GetStats(userID)
	if err != nil {
		return err
	}

	stats.TotalDays = 0
	stats.TotalTimes = 0
	stats.ContinuousDays = 0
	stats.MonthTimes = 0
	stats.YearTimes = 0
	stats.LastCheckInDate = nil

	checkIns, _, err := s.checkinRepo.List(1, 10000, &userID, nil, nil) // Assuming max 10000 check-ins
	if err != nil {
		return err
	}

	for _, checkIn := range checkIns {
		// Simplified update logic for recalculation
		stats.TotalTimes++
		if stats.LastCheckInDate == nil || !isSameDay(*stats.LastCheckInDate, checkIn.CheckInTime) {
			stats.TotalDays++
		}
		stats.LastCheckInDate = &checkIn.CheckInTime
	}

	return s.userRepo.UpdateStats(stats)
}
