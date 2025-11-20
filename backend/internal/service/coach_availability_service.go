package service

import (
"errors"
"gym-admin/internal/models"
"gym-admin/internal/repository"
)

type CoachAvailabilityService struct {
	repo *repository.CoachAvailabilityRepository
}

func NewCoachAvailabilityService() *CoachAvailabilityService {
	return &CoachAvailabilityService{
		repo: repository.NewCoachAvailabilityRepository(),
	}
}

func (s *CoachAvailabilityService) CreateAvailability(availability *models.CoachAvailability) error {
	return s.repo.Create(availability)
}

func (s *CoachAvailabilityService) GetAvailability(id int64) (*models.CoachAvailability, error) {
	return s.repo.GetByID(id)
}

func (s *CoachAvailabilityService) GetCoachAvailabilities(coachID int64) ([]models.CoachAvailability, error) {
	return s.repo.GetByCoachID(coachID)
}

func (s *CoachAvailabilityService) UpdateAvailability(id int64, updates map[string]interface{}) error {
	availability, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("availability not found")
	}

	if dayOfWeek, ok := updates["day_of_week"].(float64); ok {
		availability.DayOfWeek = int(dayOfWeek)
	}
	if startTime, ok := updates["start_time"].(string); ok {
		availability.StartTime = startTime
	}
	if endTime, ok := updates["end_time"].(string); ok {
		availability.EndTime = endTime
	}
	if isRecurring, ok := updates["is_recurring"].(bool); ok {
		availability.IsRecurring = isRecurring
	}

	return s.repo.Update(availability)
}

func (s *CoachAvailabilityService) DeleteAvailability(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("availability not found")
	}
	return s.repo.Delete(id)
}

func (s *CoachAvailabilityService) DeleteCoachAvailabilities(coachID int64) error {
	return s.repo.DeleteByCoachID(coachID)
}
