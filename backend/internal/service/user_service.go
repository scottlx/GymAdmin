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

// ChangeUserStatus changes user status and logs the change
func (s *UserService) ChangeUserStatus(userID int64, newStatus int8, reason string, operatorID *int64) error {
	// Validate status
	if newStatus < models.UserStatusActive || newStatus > models.UserStatusBlacklist {
		return errors.New("invalid status value")
	}

	// Get user
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if status is already the same
	if user.Status == newStatus {
		return errors.New("user is already in this status")
	}

	oldStatus := user.Status

	// Update user status
	user.Status = newStatus
	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Create status log
	log := &models.UserStatusLog{
		UserID:     userID,
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		Reason:     reason,
		OperatorID: operatorID,
	}
	if err := s.repo.CreateStatusLog(log); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

// FreezeUser freezes a user account
func (s *UserService) FreezeUser(userID int64, reason string, operatorID *int64) error {
	return s.ChangeUserStatus(userID, models.UserStatusFrozen, reason, operatorID)
}

// UnfreezeUser unfreezes a user account (set to active)
func (s *UserService) UnfreezeUser(userID int64, reason string, operatorID *int64) error {
	return s.ChangeUserStatus(userID, models.UserStatusActive, reason, operatorID)
}

// AddToBlacklist adds a user to blacklist
func (s *UserService) AddToBlacklist(userID int64, reason string, operatorID *int64) error {
	return s.ChangeUserStatus(userID, models.UserStatusBlacklist, reason, operatorID)
}

// RemoveFromBlacklist removes a user from blacklist (set to active)
func (s *UserService) RemoveFromBlacklist(userID int64, reason string, operatorID *int64) error {
	return s.ChangeUserStatus(userID, models.UserStatusActive, reason, operatorID)
}

// GetStatusLogs gets status change logs for a user
func (s *UserService) GetStatusLogs(userID int64, page, pageSize int) ([]models.UserStatusLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.GetStatusLogs(userID, page, pageSize)
}

// GetUserStatusSummary gets a summary of user counts by status
func (s *UserService) GetUserStatusSummary() (map[string]interface{}, error) {
	activeCount, err := s.repo.CountByStatus(models.UserStatusActive)
	if err != nil {
		return nil, err
	}

	frozenCount, err := s.repo.CountByStatus(models.UserStatusFrozen)
	if err != nil {
		return nil, err
	}

	blacklistCount, err := s.repo.CountByStatus(models.UserStatusBlacklist)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"active":    activeCount,
		"frozen":    frozenCount,
		"blacklist": blacklistCount,
		"total":     activeCount + frozenCount + blacklistCount,
	}, nil
}

// BatchFreezeUsers freezes multiple users
func (s *UserService) BatchFreezeUsers(userIDs []int64, reason string, operatorID *int64) (map[string]interface{}, error) {
	successCount := 0
	failedCount := 0
	errors := make([]string, 0)

	for _, userID := range userIDs {
		if err := s.FreezeUser(userID, reason, operatorID); err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("User %d: %s", userID, err.Error()))
		} else {
			successCount++
		}
	}

	return map[string]interface{}{
		"success_count": successCount,
		"failed_count":  failedCount,
		"errors":        errors,
	}, nil
}

// BatchUnfreezeUsers unfreezes multiple users
func (s *UserService) BatchUnfreezeUsers(userIDs []int64, reason string, operatorID *int64) (map[string]interface{}, error) {
	successCount := 0
	failedCount := 0
	errors := make([]string, 0)

	for _, userID := range userIDs {
		if err := s.UnfreezeUser(userID, reason, operatorID); err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("User %d: %s", userID, err.Error()))
		} else {
			successCount++
		}
	}

	return map[string]interface{}{
		"success_count": successCount,
		"failed_count":  failedCount,
		"errors":        errors,
	}, nil
}

// generateUserNo generates a unique user number
func (s *UserService) generateUserNo() string {
	return fmt.Sprintf("U%s%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)
}
