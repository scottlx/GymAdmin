package service

import (
	"errors"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
)

type CardTypeService struct {
	repo *repository.CardTypeRepository
}

func NewCardTypeService() *CardTypeService {
	return &CardTypeService{
		repo: repository.NewCardTypeRepository(),
	}
}

// CreateCardType creates a new card type
func (s *CardTypeService) CreateCardType(cardType *models.CardType) error {
	// Validate card type
	if cardType.TypeName == "" {
		return errors.New("type name is required")
	}
	if cardType.TypeCode == "" {
		return errors.New("type code is required")
	}
	if cardType.DurationValue <= 0 {
		return errors.New("duration value must be positive")
	}
	if cardType.Price <= 0 {
		return errors.New("price must be positive")
	}

	// Check if type code already exists
	existing, _ := s.repo.GetByTypeCode(cardType.TypeCode)
	if existing != nil {
		return errors.New("type code already exists")
	}

	// Set default values
	if cardType.Status == 0 {
		cardType.Status = 1 // Default: enabled
	}

	return s.repo.Create(cardType)
}

// GetCardType gets a card type by ID
func (s *CardTypeService) GetCardType(id int64) (*models.CardType, error) {
	return s.repo.GetByID(id)
}

// GetCardTypeByCode gets a card type by type code
func (s *CardTypeService) GetCardTypeByCode(typeCode string) (*models.CardType, error) {
	return s.repo.GetByTypeCode(typeCode)
}

// ListCardTypes lists all card types
func (s *CardTypeService) ListCardTypes(status *int8) ([]models.CardType, error) {
	return s.repo.List(status)
}

// UpdateCardType updates a card type
func (s *CardTypeService) UpdateCardType(id int64, updates map[string]interface{}) error {
	cardType, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("card type not found")
	}

	// Update fields
	if typeName, ok := updates["type_name"].(string); ok {
		cardType.TypeName = typeName
	}
	if durationType, ok := updates["duration_type"].(float64); ok {
		cardType.DurationType = int8(durationType)
	}
	if durationValue, ok := updates["duration_value"].(float64); ok {
		cardType.DurationValue = int(durationValue)
	}
	if price, ok := updates["price"].(float64); ok {
		cardType.Price = price
	}
	if originalPrice, ok := updates["original_price"].(float64); ok {
		cardType.OriginalPrice = originalPrice
	}
	if description, ok := updates["description"].(string); ok {
		cardType.Description = description
	}
	if benefits, ok := updates["benefits"].(string); ok {
		cardType.Benefits = benefits
	}
	if canFreeze, ok := updates["can_freeze"].(float64); ok {
		cardType.CanFreeze = int8(canFreeze)
	}
	if maxFreezeTimes, ok := updates["max_freeze_times"].(float64); ok {
		cardType.MaxFreezeTimes = int(maxFreezeTimes)
	}
	if maxFreezeDays, ok := updates["max_freeze_days"].(float64); ok {
		cardType.MaxFreezeDays = int(maxFreezeDays)
	}
	if canTransfer, ok := updates["can_transfer"].(float64); ok {
		cardType.CanTransfer = int8(canTransfer)
	}
	if transferFee, ok := updates["transfer_fee"].(float64); ok {
		cardType.TransferFee = transferFee
	}
	if status, ok := updates["status"].(float64); ok {
		cardType.Status = int8(status)
	}
	if sortOrder, ok := updates["sort_order"].(float64); ok {
		cardType.SortOrder = int(sortOrder)
	}

	return s.repo.Update(cardType)
}

// DeleteCardType deletes a card type
func (s *CardTypeService) DeleteCardType(id int64) error {
	// Check if any cards are using this type
	count, err := s.repo.CountCardsByType(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete card type that is in use")
	}

	return s.repo.Delete(id)
}

// EnableCardType enables a card type
func (s *CardTypeService) EnableCardType(id int64) error {
	cardType, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("card type not found")
	}

	cardType.Status = 1
	return s.repo.Update(cardType)
}

// DisableCardType disables a card type
func (s *CardTypeService) DisableCardType(id int64) error {
	cardType, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("card type not found")
	}

	cardType.Status = 2
	return s.repo.Update(cardType)
}

// UpdateSortOrder updates the sort order of card types
func (s *CardTypeService) UpdateSortOrder(orders map[int64]int) error {
	for id, order := range orders {
		cardType, err := s.repo.GetByID(id)
		if err != nil {
			continue
		}
		cardType.SortOrder = order
		if err := s.repo.Update(cardType); err != nil {
			return err
		}
	}
	return nil
}
