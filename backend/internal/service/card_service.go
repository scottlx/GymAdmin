package service

import (
	"errors"
	"fmt"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
	"time"
)

type CardService struct {
	repo *repository.CardRepository
}

func NewCardService() *CardService {
	return &CardService{
		repo: repository.NewCardRepository(),
	}
}

func (s *CardService) CreateCard(card *models.MembershipCard) error {
	// Generate card number
	card.CardNo = s.generateCardNo()
	if card.Status == 0 {
		card.Status = 1 // Default status: active
	}

	return s.repo.Create(card)
}

func (s *CardService) GetCard(id int64) (*models.MembershipCard, error) {
	return s.repo.GetByID(id)
}

func (s *CardService) GetCardByCardNo(cardNo string) (*models.MembershipCard, error) {
	return s.repo.GetByCardNo(cardNo)
}

func (s *CardService) ListCards(page, pageSize int, status *int8, userID *int64) ([]models.MembershipCard, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize, status, userID)
}

func (s *CardService) UpdateCard(id int64, updates map[string]interface{}) error {
	card, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("card not found")
	}

	// Update fields
	if userID, ok := updates["user_id"].(float64); ok {
		card.UserID = int64(userID)
	}
	if cardTypeID, ok := updates["card_type_id"].(float64); ok {
		card.CardTypeID = int64(cardTypeID)
	}
	if status, ok := updates["status"].(float64); ok {
		card.Status = int8(status)
	}
	if startDate, ok := updates["start_date"].(string); ok {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			card.StartDate = t
		}
	}
	if endDate, ok := updates["end_date"].(string); ok {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			card.EndDate = t
		}
	}
	if remainingTimes, ok := updates["remaining_times"].(float64); ok {
		times := int(remainingTimes)
		card.RemainingTimes = &times
	}
	if totalTimes, ok := updates["total_times"].(float64); ok {
		times := int(totalTimes)
		card.TotalTimes = &times
	}
	if purchasePrice, ok := updates["purchase_price"].(float64); ok {
		card.PurchasePrice = purchasePrice
	}
	if remark, ok := updates["remark"].(string); ok {
		card.Remark = remark
	}

	return s.repo.Update(card)
}

func (s *CardService) DeleteCard(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("card not found")
	}
	return s.repo.Delete(id)
}

func (s *CardService) GetCardsByUserID(userID int64) ([]models.MembershipCard, error) {
	return s.repo.GetByUserID(userID)
}

// generateCardNo generates a unique card number
func (s *CardService) generateCardNo() string {
	return fmt.Sprintf("C%s%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)
}

// RenewCard renews a membership card
func (s *CardService) RenewCard(cardID int64, months int, amount float64, operatorID int64, remark string) error {
	card, err := s.repo.GetByID(cardID)
	if err != nil {
		return errors.New("card not found")
	}

	// Check card status
	if card.Status == 4 || card.Status == 5 {
		return errors.New("cannot renew transferred or refunded card")
	}

	// Calculate new end date
	oldEndDate := card.EndDate
	var newEndDate time.Time
	if time.Now().After(card.EndDate) {
		// If expired, start from today
		newEndDate = time.Now().AddDate(0, months, 0)
	} else {
		// If not expired, extend from current end date
		newEndDate = card.EndDate.AddDate(0, months, 0)
	}

	// Update card
	card.EndDate = newEndDate
	if card.Status == 2 { // If expired, set to active
		card.Status = 1
	}

	// Save card and operation record in transaction
	if err := s.repo.Update(card); err != nil {
		return err
	}

	// Record operation
	operation := &models.CardOperation{
		CardID:        cardID,
		OperationType: 1, // Renewal
		OperatorID:    operatorID,
		Amount:        amount,
		OldEndDate:    oldEndDate,
		NewEndDate:    newEndDate,
		Remark:        remark,
	}

	return s.repo.CreateOperation(operation)
}

// FreezeCard freezes a membership card
func (s *CardService) FreezeCard(cardID int64, freezeDays int, operatorID int64, remark string) error {
	card, err := s.repo.GetByID(cardID)
	if err != nil {
		return errors.New("card not found")
	}

	// Check card status
	if card.Status != 1 {
		return errors.New("only active cards can be frozen")
	}

	if card.IsFrozen == 1 {
		return errors.New("card is already frozen")
	}

	// Get card type to check freeze rules
	cardType, err := s.repo.GetCardType(card.CardTypeID)
	if err != nil {
		return errors.New("card type not found")
	}

	if cardType.CanFreeze != 1 {
		return errors.New("this card type cannot be frozen")
	}

	// Check freeze times limit
	if cardType.MaxFreezeTimes > 0 && card.FreezeTimes >= cardType.MaxFreezeTimes {
		return errors.New("freeze times limit exceeded")
	}

	// Check freeze days limit
	if cardType.MaxFreezeDays > 0 && card.FreezeDays+freezeDays > cardType.MaxFreezeDays {
		return errors.New("freeze days limit exceeded")
	}

	// Update card
	now := time.Now()
	card.IsFrozen = 1
	card.FrozenAt = &now
	card.Status = 3 // Frozen
	card.FreezeTimes++
	card.FreezeDays += freezeDays

	// Extend end date
	oldEndDate := card.EndDate
	card.EndDate = card.EndDate.AddDate(0, 0, freezeDays)

	if err := s.repo.Update(card); err != nil {
		return err
	}

	// Record operation
	operation := &models.CardOperation{
		CardID:        cardID,
		OperationType: 2, // Freeze
		OperatorID:    operatorID,
		OldEndDate:    oldEndDate,
		NewEndDate:    card.EndDate,
		FreezeDays:    freezeDays,
		Remark:        remark,
	}

	return s.repo.CreateOperation(operation)
}

// UnfreezeCard unfreezes a membership card
func (s *CardService) UnfreezeCard(cardID int64, operatorID int64, remark string) error {
	card, err := s.repo.GetByID(cardID)
	if err != nil {
		return errors.New("card not found")
	}

	if card.IsFrozen != 1 {
		return errors.New("card is not frozen")
	}

	// Update card
	card.IsFrozen = 0
	card.FrozenAt = nil
	card.Status = 1 // Active

	if err := s.repo.Update(card); err != nil {
		return err
	}

	// Record operation
	operation := &models.CardOperation{
		CardID:        cardID,
		OperationType: 3, // Unfreeze
		OperatorID:    operatorID,
		Remark:        remark,
	}

	return s.repo.CreateOperation(operation)
}

// TransferCard transfers a membership card to another user
func (s *CardService) TransferCard(cardID int64, toUserID int64, transferFee float64, operatorID int64, remark string) error {
	card, err := s.repo.GetByID(cardID)
	if err != nil {
		return errors.New("card not found")
	}

	// Check card status
	if card.Status != 1 {
		return errors.New("only active cards can be transferred")
	}

	if card.IsFrozen == 1 {
		return errors.New("frozen card cannot be transferred")
	}

	// Get card type to check transfer rules
	cardType, err := s.repo.GetCardType(card.CardTypeID)
	if err != nil {
		return errors.New("card type not found")
	}

	if cardType.CanTransfer != 1 {
		return errors.New("this card type cannot be transferred")
	}

	// Check if target user exists
	if _, err := s.repo.GetUserByID(toUserID); err != nil {
		return errors.New("target user not found")
	}

	// Update card
	oldUserID := card.UserID
	card.UserID = toUserID
	card.Status = 4 // Transferred

	if err := s.repo.Update(card); err != nil {
		return err
	}

	// Record operation
	operation := &models.CardOperation{
		CardID:        cardID,
		OperationType: 4, // Transfer
		OperatorID:    operatorID,
		Amount:        transferFee,
		TransferToID:  &toUserID,
		Remark:        fmt.Sprintf("从用户%d转至用户%d, %s", oldUserID, toUserID, remark),
	}

	return s.repo.CreateOperation(operation)
}

// GetCardOperations gets operation history of a card
func (s *CardService) GetCardOperations(cardID int64) ([]models.CardOperation, error) {
	return s.repo.GetOperationsByCardID(cardID)
}

// CheckAndUpdateExpiredCards checks and updates expired cards status
func (s *CardService) CheckAndUpdateExpiredCards() (int64, error) {
	return s.repo.BatchUpdateExpiredCards()
}

// GetExpiringCards gets cards expiring within specified days
func (s *CardService) GetExpiringCards(days int) ([]models.MembershipCard, error) {
	return s.repo.GetExpiringCards(days)
}

// SendExpiryNotifications sends expiry notifications for cards expiring soon
func (s *CardService) SendExpiryNotifications(days int) (int, error) {
	cards, err := s.repo.GetExpiringCards(days)
	if err != nil {
		return 0, err
	}

	notificationRepo := repository.NewNotificationRepository()
	count := 0

	for _, card := range cards {
		// Calculate days until expiry
		daysUntilExpiry := int(card.EndDate.Sub(time.Now()).Hours() / 24)

		notification := &models.Notification{
			UserID:    card.UserID,
			Type:      1, // Card expiry notification
			Title:     "会员卡即将到期提醒",
			Content:   fmt.Sprintf("您的会员卡（卡号：%s）将在%d天后到期，到期日期为%s，请及时续费。", card.CardNo, daysUntilExpiry, card.EndDate.Format("2006-01-02")),
			RelatedID: &card.ID,
		}

		if err := notificationRepo.Create(notification); err == nil {
			count++
		}
	}

	return count, nil
}

// GetExpiredCards gets all expired cards
func (s *CardService) GetExpiredCards() ([]models.MembershipCard, error) {
	return s.repo.GetExpiredCards()
}
