package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
	"math/rand"
	"time"
)

type VoucherService struct {
	voucherRepo  *repository.VoucherRepository
	platformRepo *repository.PlatformRepository
}

func NewVoucherService() *VoucherService {
	return &VoucherService{
		voucherRepo:  repository.NewVoucherRepository(),
		platformRepo: repository.NewPlatformRepository(),
	}
}

// VerifyVoucher verifies a voucher code with a third-party platform (simulated)
func (s *VoucherService) VerifyVoucher(voucherCode string, platform int8, operatorID int64) (*models.VoucherRecord, error) {
	// Check if voucher already exists and has been verified
	existingVoucher, err := s.voucherRepo.GetByCode(voucherCode)
	if err == nil && existingVoucher.Status == 2 {
		return nil, errors.New("voucher already verified")
	}

	// Get platform configuration
	_, err = s.platformRepo.GetByPlatform(platform)
	if err != nil {
		return nil, errors.New("platform not configured or disabled")
	}

	// Simulate API call to third-party platform
	apiResponse, err := s.simulateThirdPartyAPICall(voucherCode, platform)
	if err != nil {
		return nil, err
	}

	// Process the API response
	voucher, err := s.processVoucherResponse(apiResponse, platform, operatorID)
	if err != nil {
		return nil, err
	}

	// Save the voucher record
	if existingVoucher != nil {
		// Update existing record
		voucher.ID = existingVoucher.ID
		if err := s.voucherRepo.Update(voucher); err != nil {
			return nil, err
		}
	} else {
		// Create new record
		if err := s.voucherRepo.Create(voucher); err != nil {
			return nil, err
		}
	}

	return voucher, nil
}

// simulateThirdPartyAPICall simulates a call to Meituan or Douyin API
func (s *VoucherService) simulateThirdPartyAPICall(voucherCode string, platform int8) (map[string]interface{}, error) {
	fmt.Printf("Simulating API call to platform %d for voucher %s\n", platform, voucherCode)
	time.Sleep(500 * time.Millisecond) // Simulate network latency

	// Simulate different responses based on voucher code
	switch {
	case len(voucherCode) < 8:
		return nil, errors.New("invalid voucher code format")
	case voucherCode == "DOUYIN_EXPIRED_123":
		return map[string]interface{}{"valid": false, "reason": "voucher expired"}, nil
	case voucherCode == "MEITUAN_USED_456":
		return map[string]interface{}{"valid": false, "reason": "voucher already used"}, nil
	default:
		// Simulate a successful response
		rand.Seed(time.Now().UnixNano())
		cardTypeID := int64(rand.Intn(3) + 1) // Random card type ID (1, 2, or 3)
		expireDays := 30 + rand.Intn(60)

		return map[string]interface{}{
			"valid":        true,
			"voucher_code": voucherCode,
			"card_type_id": cardTypeID,
			"expire_at":    time.Now().Add(time.Duration(expireDays) * 24 * time.Hour),
			"order_id":     fmt.Sprintf("ORDER_%s", voucherCode),
			"user_phone":   fmt.Sprintf("138%08d", rand.Intn(100000000)),
		}, nil
	}
}

// processVoucherResponse processes the simulated API response and creates a VoucherRecord
func (s *VoucherService) processVoucherResponse(response map[string]interface{}, platform int8, operatorID int64) (*models.VoucherRecord, error) {
	if !response["valid"].(bool) {
		return nil, fmt.Errorf("voucher verification failed: %s", response["reason"])
	}

	platformData, _ := json.Marshal(response)
	now := time.Now()
	expireAt := response["expire_at"].(time.Time)
	cardTypeID := response["card_type_id"].(int64)

	return &models.VoucherRecord{
		VoucherCode:  response["voucher_code"].(string),
		Platform:     platform,
		CardTypeID:   &cardTypeID,
		Status:       2, // 2-已核销
		VerifiedAt:   &now,
		VerifiedBy:   &operatorID,
		ExpireAt:     &expireAt,
		PlatformData: string(platformData),
	}, nil
}

// ListVouchers lists vouchers with pagination
func (s *VoucherService) ListVouchers(page, pageSize int, status *int8, platform *int8) ([]models.VoucherRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.voucherRepo.List(page, pageSize, status, platform)
}

// GetVoucher gets a voucher by ID
func (s *VoucherService) GetVoucher(id int64) (*models.VoucherRecord, error) {
	return s.voucherRepo.GetByID(id)
}
