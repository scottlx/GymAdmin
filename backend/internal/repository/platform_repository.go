package repository

import (
"gym-admin/internal/models"
"gym-admin/pkg/database"

"gorm.io/gorm"
)

type PlatformRepository struct {
	db *gorm.DB
}

func NewPlatformRepository() *PlatformRepository {
	return &PlatformRepository{db: database.GetDB()}
}

func (r *PlatformRepository) Create(platform *models.ThirdPartyPlatform) error {
	return r.db.Create(platform).Error
}

func (r *PlatformRepository) GetByID(id int64) (*models.ThirdPartyPlatform, error) {
	var platform models.ThirdPartyPlatform
	err := r.db.First(&platform, id).Error
	return &platform, err
}

func (r *PlatformRepository) GetByPlatform(platform int8) (*models.ThirdPartyPlatform, error) {
	var p models.ThirdPartyPlatform
	err := r.db.Where("platform = ? AND status = ?", platform, 1).First(&p).Error
	return &p, err
}

func (r *PlatformRepository) List() ([]models.ThirdPartyPlatform, error) {
	var platforms []models.ThirdPartyPlatform
	err := r.db.Find(&platforms).Error
	return platforms, err
}

func (r *PlatformRepository) Update(platform *models.ThirdPartyPlatform) error {
	return r.db.Save(platform).Error
}

func (r *PlatformRepository) Delete(id int64) error {
	return r.db.Delete(&models.ThirdPartyPlatform{}, id).Error
}
