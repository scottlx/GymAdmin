package service

import (
	"errors"
	"fmt"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
	"time"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

func (s *UserService) CreateUser(user *models.User) error {
	// Generate user number
	user.UserNo = s.generateUserNo()
	user.Status = 1 // Default status: active

	return s.repo.Create(user)
}

func (s *UserService) GetUser(id int64) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetUserByPhone(phone string) (*models.User, error) {
	return s.repo.GetByPhone(phone)
}

func (s *UserService) ListUsers(page, pageSize int, status *int8) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize, status)
}

func (s *UserService) UpdateUser(id int64, updates map[string]interface{}) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	// Update fields
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	if gender, ok := updates["gender"].(float64); ok {
		user.Gender = int8(gender)
	}
	if phone, ok := updates["phone"].(string); ok {
		user.Phone = phone
	}
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if avatarURL, ok := updates["avatar_url"].(string); ok {
		user.AvatarURL = avatarURL
	}
	if address, ok := updates["address"].(string); ok {
		user.Address = address
	}
	if emergencyContact, ok := updates["emergency_contact"].(string); ok {
		user.EmergencyContact = emergencyContact
	}
	if emergencyPhone, ok := updates["emergency_phone"].(string); ok {
		user.EmergencyPhone = emergencyPhone
	}
	if healthStatus, ok := updates["health_status"].(string); ok {
		user.HealthStatus = healthStatus
	}
	if trainingGoal, ok := updates["training_goal"].(string); ok {
		user.TrainingGoal = trainingGoal
	}
	if status, ok := updates["status"].(float64); ok {
		user.Status = int8(status)
	}
	if remark, ok := updates["remark"].(string); ok {
		user.Remark = remark
	}

	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(id)
}

func (s *UserService) GetUserStats(userID int64) (*models.UserTrainingStats, error) {
	return s.repo.GetStats(userID)
}

// generateUserNo generates a unique user number
func (s *UserService) generateUserNo() string {
	return fmt.Sprintf("U%s%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)
}
