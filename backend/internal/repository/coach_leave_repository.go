package repository

import (
	"gym-admin/internal/models"
	"gym-admin/pkg/database"
	"time"

	"gorm.io/gorm"
)

type CoachLeaveRepository struct {
	db *gorm.DB
}

func NewCoachLeaveRepository() *CoachLeaveRepository {
	return &CoachLeaveRepository{db: database.GetDB()}
}

func (r *CoachLeaveRepository) Create(leave *models.CoachLeave) error {
	return r.db.Create(leave).Error
}

func (r *CoachLeaveRepository) GetByID(id int64) (*models.CoachLeave, error) {
	var leave models.CoachLeave
	err := r.db.First(&leave, id).Error
	return &leave, err
}

func (r *CoachLeaveRepository) GetByCoachID(coachID int64) ([]models.CoachLeave, error) {
	var leaves []models.CoachLeave
	err := r.db.Where("coach_id = ?", coachID).Find(&leaves).Error
	return leaves, err
}

func (r *CoachLeaveRepository) Update(leave *models.CoachLeave) error {
	return r.db.Save(leave).Error
}

func (r *CoachLeaveRepository) Delete(id int64) error {
	return r.db.Delete(&models.CoachLeave{}, id).Error
}

func (r *CoachLeaveRepository) GetConflictingLeaves(coachID int64, startTime, endTime time.Time) ([]models.CoachLeave, error) {
	var leaves []models.CoachLeave
	err := r.db.Where("coach_id = ? AND status = ? AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?) OR (end_time > ? AND end_time <= ?))",
		coachID, models.LeaveStatusApproved, endTime, startTime, startTime, endTime, startTime, endTime).Find(&leaves).Error
	return leaves, err
}
