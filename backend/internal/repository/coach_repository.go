package repository

import (
	"gym-admin/internal/models"
	"gym-admin/pkg/database"

	"gorm.io/gorm"
)

type CoachRepository struct {
	db *gorm.DB
}

func NewCoachRepository() *CoachRepository {
	return &CoachRepository{db: database.GetDB()}
}

func (r *CoachRepository) Create(coach *models.Coach) error {
	return r.db.Create(coach).Error
}

func (r *CoachRepository) GetByID(id int64) (*models.Coach, error) {
	var coach models.Coach
	err := r.db.First(&coach, id).Error
	return &coach, err
}

func (r *CoachRepository) List(page, pageSize int, status *int8) ([]models.Coach, int64, error) {
	var coaches []models.Coach
	var total int64

	query := r.db.Model(&models.Coach{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&coaches).Error
	return coaches, total, err
}

func (r *CoachRepository) Update(coach *models.Coach) error {
	return r.db.Save(coach).Error
}

func (r *CoachRepository) Delete(id int64) error {
	return r.db.Delete(&models.Coach{}, id).Error
}
