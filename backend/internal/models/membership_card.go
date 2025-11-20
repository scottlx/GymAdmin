package models

import (
	"time"

	"gorm.io/gorm"
)

type CardType struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeName       string    `gorm:"type:varchar(50);not null" json:"type_name"`
	TypeCode       string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"type_code"`
	DurationType   int8      `gorm:"type:tinyint;not null" json:"duration_type"` // 1-天卡，2-月卡，3-季卡，4-年卡，5-次卡
	DurationValue  int       `gorm:"not null" json:"duration_value"`
	Price          float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice  float64   `gorm:"type:decimal(10,2)" json:"original_price"`
	Description    string    `gorm:"type:text" json:"description"`
	Benefits       string    `gorm:"type:text" json:"benefits"` // JSON格式
	CanFreeze      int8      `gorm:"type:tinyint;default:1" json:"can_freeze"`
	MaxFreezeTimes int       `gorm:"default:0" json:"max_freeze_times"`
	MaxFreezeDays  int       `gorm:"default:0" json:"max_freeze_days"`
	CanTransfer    int8      `gorm:"type:tinyint;default:1" json:"can_transfer"`
	TransferFee    float64   `gorm:"type:decimal(10,2);default:0" json:"transfer_fee"`
	Status         int8      `gorm:"type:tinyint;default:1;index" json:"status"` // 1-启用，2-停用
	SortOrder      int       `gorm:"default:0;index" json:"sort_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (CardType) TableName() string {
	return "card_types"
}

type MembershipCard struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CardNo         string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"card_no"`
	UserID         int64          `gorm:"index;not null" json:"user_id"`
	CardTypeID     int64          `gorm:"not null" json:"card_type_id"`
	Status         int8           `gorm:"type:tinyint;default:1;index" json:"status"` // 1-正常，2-已过期，3-已冻结，4-已转出，5-已退卡
	StartDate      time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate        time.Time      `gorm:"type:date;not null;index" json:"end_date"`
	RemainingTimes *int           `json:"remaining_times"`
	TotalTimes     *int           `json:"total_times"`
	FreezeTimes    int            `gorm:"default:0" json:"freeze_times"`
	FreezeDays     int            `gorm:"default:0" json:"freeze_days"`
	IsFrozen       int8           `gorm:"type:tinyint;default:0" json:"is_frozen"`
	FrozenAt       *time.Time     `json:"frozen_at"`
	Source         int8           `gorm:"type:tinyint;default:1" json:"source"` // 1-前台办理，2-小程序购买，3-美团，4-抖音
	PurchasePrice  float64        `gorm:"type:decimal(10,2);not null" json:"purchase_price"`
	OperatorID     *int64         `json:"operator_id"`
	Remark         string         `gorm:"type:text" json:"remark"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (MembershipCard) TableName() string {
	return "membership_cards"
}

// CardOperation 会员卡操作记录
type CardOperation struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CardID        int64     `gorm:"index;not null" json:"card_id"`
	OperationType int8      `gorm:"type:tinyint;not null;index" json:"operation_type"` // 1-续费，2-冻结，3-解冻，4-转卡，5-退卡
	OperatorID    int64     `gorm:"not null" json:"operator_id"`
	Amount        float64   `gorm:"type:decimal(10,2)" json:"amount"`
	OldEndDate    time.Time `gorm:"type:date" json:"old_end_date"`
	NewEndDate    time.Time `gorm:"type:date" json:"new_end_date"`
	FreezeDays    int       `json:"freeze_days"`
	TransferToID  *int64    `json:"transfer_to_id"` // 转卡目标用户ID
	Remark        string    `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
}

func (CardOperation) TableName() string {
	return "card_operations"
}
