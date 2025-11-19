package models

import (
	"time"

	"gorm.io/gorm"
)

type VoucherRecord struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	VoucherCode  string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"voucher_code"`
	Platform     int8           `gorm:"type:tinyint;not null;index" json:"platform"` // 1-美团，2-抖音
	UserID       *int64         `gorm:"index" json:"user_id"`
	CardTypeID   *int64         `json:"card_type_id"`
	Status       int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-未核销，2-已核销，3-已过期
	VerifiedAt   *time.Time     `json:"verified_at"`
	VerifiedBy   *int64         `json:"verified_by"` // 核销操作员ID
	ExpireAt     *time.Time     `gorm:"index" json:"expire_at"`
	PlatformData string         `gorm:"type:text" json:"platform_data"` // 平台原始数据JSON
	Remark       string         `gorm:"type:text" json:"remark"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (VoucherRecord) TableName() string {
	return "voucher_records"
}
