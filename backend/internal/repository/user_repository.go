package repository

import (
	"gym-admin/internal/models"
	"gym-admin/pkg/database"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.GetDB()}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *UserRepository) List(page, pageSize int, status *int8) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id int64) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) GetStats(userID int64) (*models.UserTrainingStats, error) {
	var stats models.UserTrainingStats
	err := r.db.Where("user_id = ?", userID).First(&stats).Error
	if err == gorm.ErrRecordNotFound {
		// Create new stats if not exists
		stats = models.UserTrainingStats{UserID: userID}
		if err := r.db.Create(&stats).Error; err != nil {
			return nil, err
		}
	}
	return &stats, err
}

func (r *UserRepository) UpdateStats(stats *models.UserTrainingStats) error {
	if stats.ID == 0 {
		// Create new stats if ID is 0
		return r.db.Create(stats).Error
	}
	return r.db.Save(stats).Error
}

// CreateStatusLog creates a status change log
func (r *UserRepository) CreateStatusLog(log *models.UserStatusLog) error {
	return r.db.Create(log).Error
}

// GetStatusLogs gets status change logs for a user
func (r *UserRepository) GetStatusLogs(userID int64, page, pageSize int) ([]models.UserStatusLog, int64, error) {
	var logs []models.UserStatusLog
	var total int64

	query := r.db.Model(&models.UserStatusLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

// GetLatestStatusLog gets the latest status log for a user
func (r *UserRepository) GetLatestStatusLog(userID int64) (*models.UserStatusLog, error) {
	var log models.UserStatusLog
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&log).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &log, err
}

// CountByStatus counts users by status
func (r *UserRepository) CountByStatus(status int8) (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
