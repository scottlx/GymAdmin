package service

import (
	"errors"
	"fmt"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
	"time"
)

type CoachService struct {
	repo *repository.CoachRepository
}

func NewCoachService() *CoachService {
	return &CoachService{
		repo: repository.NewCoachRepository(),
	}
}

func (s *CoachService) CreateCoach(coach *models.Coach) error {
	// Generate coach number
	coach.CoachNo = s.generateCoachNo()
	coach.Status = 1 // Default status: active

	return s.repo.Create(coach)
}

func (s *CoachService) GetCoach(id int64) (*models.Coach, error) {
	return s.repo.GetByID(id)
}

func (s *CoachService) ListCoaches(page, pageSize int, status *int8) ([]models.Coach, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize, status)
}

func (s *CoachService) UpdateCoach(id int64, updates map[string]interface{}) error {
	coach, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("coach not found")
	}

	// Update fields
	if name, ok := updates["name"].(string); ok {
		coach.Name = name
	}
	if phone, ok := updates["phone"].(string); ok {
		coach.Phone = phone
	}
	if email, ok := updates["email"].(string); ok {
		coach.Email = email
	}
	if introduction, ok := updates["introduction"].(string); ok {
		coach.Introduction = introduction
	}

	return s.repo.Update(coach)
}

func (s *CoachService) DeleteCoach(id int64) error {
	return s.repo.Delete(id)
}

func (s *CoachService) generateCoachNo() string {
	return fmt.Sprintf("C%d", time.Now().UnixNano()/1000000)
}

func (s *CoachService) GetAllCoaches() ([]models.Coach, error) {
	return s.repo.GetAll()
}
