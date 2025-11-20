package controller

import (
	"gym-admin/internal/models"
	"gym-admin/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardTypeController struct {
	service *service.CardTypeService
}

func NewCardTypeController() *CardTypeController {
	return &CardTypeController{
		service: service.NewCardTypeService(),
	}
}

// CreateCardType creates a new card type
// @Summary Create card type
// @Tags CardType
// @Accept json
// @Produce json
// @Param body body CreateCardTypeRequest true "Card type info"
// @Success 200 {object} Response
// @Router /api/v1/card-types [post]
func (c *CardTypeController) CreateCardType(ctx *gin.Context) {
	var req CreateCardTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cardType := &models.CardType{
		TypeName:       req.TypeName,
		TypeCode:       req.TypeCode,
		DurationType:   req.DurationType,
		DurationValue:  req.DurationValue,
		Price:          req.Price,
		OriginalPrice:  req.OriginalPrice,
		Description:    req.Description,
		Benefits:       req.Benefits,
		CanFreeze:      req.CanFreeze,
		MaxFreezeTimes: req.MaxFreezeTimes,
		MaxFreezeDays:  req.MaxFreezeDays,
		CanTransfer:    req.CanTransfer,
		TransferFee:    req.TransferFee,
		SortOrder:      req.SortOrder,
	}

	if err := c.service.CreateCardType(cardType); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "card type created successfully",
		"data":    cardType,
	})
}

// GetCardType gets a card type by ID
// @Summary Get card type
// @Tags CardType
// @Produce json
// @Param id path int true "Card Type ID"
// @Success 200 {object} Response
// @Router /api/v1/card-types/:id [get]
func (c *CardTypeController) GetCardType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid card type ID"})
		return
	}

	cardType, err := c.service.GetCardType(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "card type not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    cardType,
	})
}

// ListCardTypes lists all card types
// @Summary List card types
// @Tags CardType
// @Produce json
// @Param status query int false "Status filter"
// @Success 200 {object} Response
// @Router /api/v1/card-types [get]
func (c *CardTypeController) ListCardTypes(ctx *gin.Context) {
	var status *int8
	if statusStr := ctx.Query("status"); statusStr != "" {
		s, err := strconv.ParseInt(statusStr, 10, 8)
		if err == nil {
			status8 := int8(s)
			status = &status8
		}
	}

	cardTypes, err := c.service.ListCardTypes(status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    cardTypes,
	})
}

// UpdateCardType updates a card type
// @Summary Update card type
// @Tags CardType
// @Accept json
// @Produce json
// @Param id path int true "Card Type ID"
// @Param body body UpdateCardTypeRequest true "Update info"
// @Success 200 {object} Response
// @Router /api/v1/card-types/:id [put]
func (c *CardTypeController) UpdateCardType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid card type ID"})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateCardType(id, updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "card type updated successfully",
	})
}

// DeleteCardType deletes a card type
// @Summary Delete card type
// @Tags CardType
// @Produce json
// @Param id path int true "Card Type ID"
// @Success 200 {object} Response
// @Router /api/v1/card-types/:id [delete]
func (c *CardTypeController) DeleteCardType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid card type ID"})
		return
	}

	if err := c.service.DeleteCardType(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "card type deleted successfully",
	})
}

// EnableCardType enables a card type
// @Summary Enable card type
// @Tags CardType
// @Produce json
// @Param id path int true "Card Type ID"
// @Success 200 {object} Response
// @Router /api/v1/card-types/:id/enable [post]
func (c *CardTypeController) EnableCardType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid card type ID"})
		return
	}

	if err := c.service.EnableCardType(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "card type enabled successfully",
	})
}

// DisableCardType disables a card type
// @Summary Disable card type
// @Tags CardType
// @Produce json
// @Param id path int true "Card Type ID"
// @Success 200 {object} Response
// @Router /api/v1/card-types/:id/disable [post]
func (c *CardTypeController) DisableCardType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid card type ID"})
		return
	}

	if err := c.service.DisableCardType(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "card type disabled successfully",
	})
}

// UpdateSortOrder updates the sort order of card types
// @Summary Update sort order
// @Tags CardType
// @Accept json
// @Produce json
// @Param body body UpdateSortOrderRequest true "Sort order info"
// @Success 200 {object} Response
// @Router /api/v1/card-types/sort-order [post]
func (c *CardTypeController) UpdateSortOrder(ctx *gin.Context) {
	var req UpdateSortOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.UpdateSortOrder(req.Orders); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "sort order updated successfully",
	})
}

// Request structs
type CreateCardTypeRequest struct {
	TypeName       string  `json:"type_name" binding:"required"`
	TypeCode       string  `json:"type_code" binding:"required"`
	DurationType   int8    `json:"duration_type" binding:"required,oneof=1 2 3 4 5"` // 1-天卡，2-月卡，3-季卡，4-年卡，5-次卡
	DurationValue  int     `json:"duration_value" binding:"required,min=1"`
	Price          float64 `json:"price" binding:"required,min=0"`
	OriginalPrice  float64 `json:"original_price"`
	Description    string  `json:"description"`
	Benefits       string  `json:"benefits"` // JSON格式
	CanFreeze      int8    `json:"can_freeze" binding:"oneof=0 1"`
	MaxFreezeTimes int     `json:"max_freeze_times"`
	MaxFreezeDays  int     `json:"max_freeze_days"`
	CanTransfer    int8    `json:"can_transfer" binding:"oneof=0 1"`
	TransferFee    float64 `json:"transfer_fee"`
	SortOrder      int     `json:"sort_order"`
}

type UpdateCardTypeRequest struct {
	TypeName       string  `json:"type_name"`
	DurationType   int8    `json:"duration_type"`
	DurationValue  int     `json:"duration_value"`
	Price          float64 `json:"price"`
	OriginalPrice  float64 `json:"original_price"`
	Description    string  `json:"description"`
	Benefits       string  `json:"benefits"`
	CanFreeze      int8    `json:"can_freeze"`
	MaxFreezeTimes int     `json:"max_freeze_times"`
	MaxFreezeDays  int     `json:"max_freeze_days"`
	CanTransfer    int8    `json:"can_transfer"`
	TransferFee    float64 `json:"transfer_fee"`
	SortOrder      int     `json:"sort_order"`
}

type UpdateSortOrderRequest struct {
	Orders map[int64]int `json:"orders" binding:"required"`
}
