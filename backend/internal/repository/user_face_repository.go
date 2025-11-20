package repository

import (
"gym-admin/internal/models"
"gym-admin/pkg/database"

"gorm.io/gorm"
)

type UserFaceRepository struct {
	db *gorm.DB
}

func NewUserFaceRepository() *UserFaceRepository {
	return &UserFaceRepository{db: database.GetDB()}
}

func (r *UserFaceRepository) Create(face *models.UserFace) error {
	return r.db.Create(face).Error
}

func (r *UserFaceRepository) GetByID(id int64) (*models.UserFace, error) {
	var face models.UserFace
	err := r.db.First(&face, id).Error
	return &face, err
}

func (r *UserFaceRepository) GetByFaceID(faceID string) (*models.UserFace, error) {
	var face models.UserFace
	err := r.db.Where("face_id = ?", faceID).First(&face).Error
	return &face, err
}

func (r *UserFaceRepository) GetByUserID(userID int64) ([]models.UserFace, error) {
	var faces []models.UserFace
	err := r.db.Where("user_id = ? AND status = ?", userID, models.FaceStatusNormal).Find(&faces).Error
	return faces, err
}

func (r *UserFaceRepository) GetMainFace(userID int64) (*models.UserFace, error) {
	var face models.UserFace
	err := r.db.Where("user_id = ? AND is_main = ? AND status = ?", userID, true, models.FaceStatusNormal).First(&face).Error
	if err == gorm.ErrRecordNotFound {
		// If no main face, get the first one
		err = r.db.Where("user_id = ? AND status = ?", userID, models.FaceStatusNormal).Order("created_at ASC").First(&face).Error
	}
	return &face, err
}

func (r *UserFaceRepository) List(page, pageSize int, userID *int64, status *int8) ([]models.UserFace, int64, error) {
	var faces []models.UserFace
	var total int64

	query := r.db.Model(&models.UserFace{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&faces).Error
	return faces, total, err
}

func (r *UserFaceRepository) Update(face *models.UserFace) error {
	return r.db.Save(face).Error
}

func (r *UserFaceRepository) Delete(id int64) error {
	return r.db.Delete(&models.UserFace{}, id).Error
}

func (r *UserFaceRepository) UpdateStatus(id int64, status int8) error {
	return r.db.Model(&models.UserFace{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserFaceRepository) SetMainFace(userID, faceID int64) error {
	// First, unset all main faces for this user
	if err := r.db.Model(&models.UserFace{}).Where("user_id = ?", userID).Update("is_main", false).Error; err != nil {
		return err
	}
	// Then set the specified face as main
	return r.db.Model(&models.UserFace{}).Where("id = ?", faceID).Update("is_main", true).Error
}

func (r *UserFaceRepository) CountByUser(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.UserFace{}).Where("user_id = ? AND status = ?", userID, models.FaceStatusNormal).Count(&count).Error
	return count, err
}

func (r *UserFaceRepository) DeleteByUserID(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserFace{}).Error
}

func (r *UserFaceRepository) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.UserFace{}).Where("status = ?", models.FaceStatusNormal).Count(&count).Error
	return count, err
}
