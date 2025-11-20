package repository

import (
	"gym-admin/internal/models"

	"gorm.io/gorm"
)

type CardTypeRepository struct {
	db *gorm.DB
}

func NewCardTypeRepository() *CardTypeRepository {
	return &CardTypeRepository{}
}

func (r *CardTypeRepository) Create(cardType *models.CardType) error {
	return r.db.Create(cardType).Error
}

func (r *CardTypeRepository) GetByID(id int64) (*models.CardType, error) {
	var cardType models.CardType
	err := r.db.First(&cardType, id).Error
	return &cardType, err
}

func (r *CardTypeRepository) GetByTypeCode(typeCode string) (*models.CardType, error) {
	var cardType models.CardType
	err := r.db.Where("type_code = ?", typeCode).First(&cardType).Error
	if err != nil {
		return nil, err
	}
	return &cardType, nil
}

func (r *CardTypeRepository) List(status *int8) ([]models.CardType, error) {
	var cardTypes []models.CardType
	query := r.db.Model(&models.CardType{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("sort_order ASC, created_at DESC").Find(&cardTypes).Error
	return cardTypes, err
}

func (r *CardTypeRepository) Update(cardType *models.CardType) error {
	return r.db.Save(cardType).Error
}

func (r *CardTypeRepository) Delete(id int64) error {
	return r.db.Delete(&models.CardType{}, id).Error
}

func (r *CardTypeRepository) CountCardsByType(cardTypeID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.MembershipCard{}).
		Where("card_type_id = ?", cardTypeID).
		Count(&count).Error
	return count, err
}
