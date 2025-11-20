package repository

import (
"gym-admin/internal/models"
"gym-admin/pkg/database"

"gorm.io/gorm"
)

type CoachAvailabilityRepository struct {
	db *gorm.DB
}

func NewCoachAvailabilityRepository() *CoachAvailabilityRepository {
	return &CoachAvailabilityRepository{db: database.GetDB()}
}

func (r *CoachAvailabilityRepository) Create(availability *models.CoachAvailability) error {
	return r.db.Create(availability).Error
}

func (r *CoachAvailabilityRepository) GetByID(id int64) (*models.CoachAvailability, error) {
	var availability models.CoachAvailability
	err := r.db.First(&availability, id).Error
	return &availability, err
}

func (r *CoachAvailabilityRepository) GetByCoachID(coachID int64) ([]models.CoachAvailability, error) {
	var availabilities []models.CoachAvailability
	err := r.db.Where("coach_id = ?", coachID).Find(&availabilities).Error
	return availabilities, err
}

func (r *CoachAvailabilityRepository) Update(availability *models.CoachAvailability) error {
	return r.db.Save(availability).Error
}

func (r *CoachAvailabilityRepository) Delete(id int64) error {
	return r.db.Delete(&models.CoachAvailability{}, id).Error
}

func (r *CoachAvailabilityRepository) DeleteByCoachID(coachID int64) error {
	return r.db.Where("coach_id = ?", coachID).Delete(&models.CoachAvailability{}).Error
}
