package repository

import (
	"gym-admin/internal/models"
	"gym-admin/pkg/database"

	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository() *CardRepository {
	return &CardRepository{db: database.GetDB()}
}

func (r *CardRepository) Create(card *models.MembershipCard) error {
	return r.db.Create(card).Error
}

func (r *CardRepository) GetByID(id int64) (*models.MembershipCard, error) {
	var card models.MembershipCard
	err := r.db.First(&card, id).Error
	return &card, err
}

func (r *CardRepository) GetByCardNo(cardNo string) (*models.MembershipCard, error) {
	var card models.MembershipCard
	err := r.db.Where("card_no = ?", cardNo).First(&card).Error
	return &card, err
}

func (r *CardRepository) List(page, pageSize int, status *int8, userID *int64) ([]models.MembershipCard, int64, error) {
	var cards []models.MembershipCard
	var total int64

	query := r.db.Model(&models.MembershipCard{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch records with pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&cards).Error

	// Initialize empty slice if nil
	if cards == nil {
		cards = []models.MembershipCard{}
	}

	return cards, total, err
}

func (r *CardRepository) Update(card *models.MembershipCard) error {
	return r.db.Save(card).Error
}

func (r *CardRepository) Delete(id int64) error {
	return r.db.Delete(&models.MembershipCard{}, id).Error
}

func (r *CardRepository) GetByUserID(userID int64) ([]models.MembershipCard, error) {
	var cards []models.MembershipCard
	err := r.db.Where("user_id = ?", userID).Find(&cards).Error
	return cards, err
}

// GetCardType gets a card type by ID
func (r *CardRepository) GetCardType(cardTypeID int64) (*models.CardType, error) {
	var cardType models.CardType
	err := r.db.First(&cardType, cardTypeID).Error
	return &cardType, err
}

// GetUserByID gets a user by ID
func (r *CardRepository) GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	return &user, err
}

// CreateOperation creates a card operation record
func (r *CardRepository) CreateOperation(operation *models.CardOperation) error {
	return r.db.Create(operation).Error
}

// GetOperationsByCardID gets all operations for a card
func (r *CardRepository) GetOperationsByCardID(cardID int64) ([]models.CardOperation, error) {
	var operations []models.CardOperation
	err := r.db.Where("card_id = ?", cardID).Order("created_at DESC").Find(&operations).Error
	return operations, err
}

// GetExpiredCards gets all expired cards that are still marked as active
func (r *CardRepository) GetExpiredCards() ([]models.MembershipCard, error) {
	var cards []models.MembershipCard
	err := r.db.Where("status = 1 AND end_date < CURDATE()").Find(&cards).Error
	return cards, err
}

// GetExpiringCards gets cards expiring within specified days
func (r *CardRepository) GetExpiringCards(days int) ([]models.MembershipCard, error) {
	var cards []models.MembershipCard
	err := r.db.Where("status = 1 AND end_date >= CURDATE() AND end_date <= DATE_ADD(CURDATE(), INTERVAL ? DAY)", days).Find(&cards).Error
	return cards, err
}

// UpdateCardStatus updates card status
func (r *CardRepository) UpdateCardStatus(cardID int64, status int8) error {
	return r.db.Model(&models.MembershipCard{}).Where("id = ?", cardID).Update("status", status).Error
}

// BatchUpdateExpiredCards updates all expired cards to expired status
func (r *CardRepository) BatchUpdateExpiredCards() (int64, error) {
	result := r.db.Model(&models.MembershipCard{}).
		Where("status = 1 AND end_date < CURDATE()").
		Update("status", 2)
	return result.RowsAffected, result.Error
}
