package repository

import (
"gym-admin/internal/models"
"gym-admin/pkg/database"

"gorm.io/gorm"
)

type CoachPerformanceRepository struct {
	db *gorm.DB
}

func NewCoachPerformanceRepository() *CoachPerformanceRepository {
	return &CoachPerformanceRepository{db: database.GetDB()}
}

func (r *CoachPerformanceRepository) Create(performance *models.CoachPerformance) error {
	return r.db.Create(performance).Error
}

func (r *CoachPerformanceRepository) GetByCoachAndMonth(coachID int64, year, month int) (*models.CoachPerformance, error) {
	var performance models.CoachPerformance
	err := r.db.Where("coach_id = ? AND year = ? AND month = ?", coachID, year, month).First(&performance).Error
	return &performance, err
}

func (r *CoachPerformanceRepository) Update(performance *models.CoachPerformance) error {
	return r.db.Save(performance).Error
}

func (r *CoachPerformanceRepository) GetTopCoachesByIncome(year, month, limit int) ([]models.CoachPerformance, error) {
	var performances []models.CoachPerformance
	err := r.db.Where("year = ? AND month = ?", year, month).Order("total_income DESC").Limit(limit).Find(&performances).Error
	return performances, err
}

func (r *CoachPerformanceRepository) GetTopCoachesByHours(year, month, limit int) ([]models.CoachPerformance, error) {
	var performances []models.CoachPerformance
	err := r.db.Where("year = ? AND month = ?", year, month).Order("total_hours DESC").Limit(limit).Find(&performances).Error
	return performances, err
}

func (r *CoachPerformanceRepository) GetTopCoachesByRating(year, month, limit int) ([]models.CoachPerformance, error) {
	var performances []models.CoachPerformance
	err := r.db.Where("year = ? AND month = ?", year, month).Order("average_rating DESC").Limit(limit).Find(&performances).Error
	return performances, err
}
