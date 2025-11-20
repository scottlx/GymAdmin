package repository

import (
	"gym-admin/internal/models"
	"gym-admin/pkg/database"

	"gorm.io/gorm"
)

type CoachCertificationRepository struct {
	db *gorm.DB
}

func NewCoachCertificationRepository() *CoachCertificationRepository {
	return &CoachCertificationRepository{db: database.GetDB()}
}

func (r *CoachCertificationRepository) Create(certification *models.CoachCertification) error {
	return r.db.Create(certification).Error
}

func (r *CoachCertificationRepository) GetByID(id int64) (*models.CoachCertification, error) {
	var certification models.CoachCertification
	err := r.db.First(&certification, id).Error
	return &certification, err
}

func (r *CoachCertificationRepository) GetByCoachID(coachID int64) ([]models.CoachCertification, error) {
	var certifications []models.CoachCertification
	err := r.db.Where("coach_id = ?", coachID).Find(&certifications).Error
	return certifications, err
}

func (r *CoachCertificationRepository) List(page, pageSize int, status *int8) ([]models.CoachCertification, int64, error) {
	var certifications []models.CoachCertification
	var total int64

	query := r.db.Model(&models.CoachCertification{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&certifications).Error
	return certifications, total, err
}

func (r *CoachCertificationRepository) Update(certification *models.CoachCertification) error {
	return r.db.Save(certification).Error
}

func (r *CoachCertificationRepository) Delete(id int64) error {
	return r.db.Delete(&models.CoachCertification{}, id).Error
}

func (r *CoachCertificationRepository) Approve(id int64, reviewerID int64, notes string) error {
	return r.db.Model(&models.CoachCertification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.CertificationStatusApproved,
		"reviewer_id":  reviewerID,
		"reviewed_at":  gorm.Expr("NOW()"),
		"review_notes": notes,
		"updated_at":   gorm.Expr("NOW()"),
	}).Error
}

func (r *CoachCertificationRepository) Reject(id int64, reviewerID int64, notes string) error {
	return r.db.Model(&models.CoachCertification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.CertificationStatusRejected,
		"reviewer_id":  reviewerID,
		"reviewed_at":  gorm.Expr("NOW()"),
		"review_notes": notes,
		"updated_at":   gorm.Expr("NOW()"),
	}).Error
}
