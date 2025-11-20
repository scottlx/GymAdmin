package repository

import (
"gym-admin/internal/models"
"gym-admin/pkg/database"

"gorm.io/gorm"
)

type FaceDeviceRepository struct {
	db *gorm.DB
}

func NewFaceDeviceRepository() *FaceDeviceRepository {
	return &FaceDeviceRepository{db: database.GetDB()}
}

func (r *FaceDeviceRepository) Create(device *models.FaceDevice) error {
	return r.db.Create(device).Error
}

func (r *FaceDeviceRepository) GetByID(id int64) (*models.FaceDevice, error) {
	var device models.FaceDevice
	err := r.db.First(&device, id).Error
	return &device, err
}

func (r *FaceDeviceRepository) GetByDeviceNo(deviceNo string) (*models.FaceDevice, error) {
	var device models.FaceDevice
	err := r.db.Where("device_no = ?", deviceNo).First(&device).Error
	return &device, err
}

func (r *FaceDeviceRepository) List(page, pageSize int, status *int8, deviceType *int8) ([]models.FaceDevice, int64, error) {
	var devices []models.FaceDevice
	var total int64

	query := r.db.Model(&models.FaceDevice{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if deviceType != nil {
		query = query.Where("device_type = ?", *deviceType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&devices).Error
	return devices, total, err
}

func (r *FaceDeviceRepository) Update(device *models.FaceDevice) error {
	return r.db.Save(device).Error
}

func (r *FaceDeviceRepository) Delete(id int64) error {
	return r.db.Delete(&models.FaceDevice{}, id).Error
}

func (r *FaceDeviceRepository) UpdateStatus(id int64, status int8) error {
	return r.db.Model(&models.FaceDevice{}).Where("id = ?", id).Update("status", status).Error
}

func (r *FaceDeviceRepository) UpdateLastOnline(id int64) error {
	return r.db.Model(&models.FaceDevice{}).Where("id = ?", id).Update("last_online", gorm.Expr("NOW()")).Error
}

func (r *FaceDeviceRepository) CountByStatus(status int8) (int64, error) {
	var count int64
	err := r.db.Model(&models.FaceDevice{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *FaceDeviceRepository) GetOnlineDevices() ([]models.FaceDevice, error) {
	var devices []models.FaceDevice
	err := r.db.Where("status = ?", models.DeviceStatusOnline).Find(&devices).Error
	return devices, err
}
