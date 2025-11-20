package repository

import (
"gym-admin/internal/models"
"gym-admin/pkg/database"

"gorm.io/gorm"
)

type VoucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository() *VoucherRepository {
	return &VoucherRepository{db: database.GetDB()}
}

func (r *VoucherRepository) Create(voucher *models.VoucherRecord) error {
	return r.db.Create(voucher).Error
}

func (r *VoucherRepository) GetByID(id int64) (*models.VoucherRecord, error) {
	var voucher models.VoucherRecord
	err := r.db.First(&voucher, id).Error
	return &voucher, err
}

func (r *VoucherRepository) GetByCode(code string) (*models.VoucherRecord, error) {
	var voucher models.VoucherRecord
	err := r.db.Where("voucher_code = ?", code).First(&voucher).Error
	return &voucher, err
}

func (r *VoucherRepository) List(page, pageSize int, status *int8, platform *int8) ([]models.VoucherRecord, int64, error) {
	var vouchers []models.VoucherRecord
	var total int64

	query := r.db.Model(&models.VoucherRecord{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if platform != nil {
		query = query.Where("platform = ?", *platform)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&vouchers).Error
	return vouchers, total, err
}

func (r *VoucherRepository) Update(voucher *models.VoucherRecord) error {
	return r.db.Save(voucher).Error
}

func (r *VoucherRepository) Delete(id int64) error {
	return r.db.Delete(&models.VoucherRecord{}, id).Error
}
