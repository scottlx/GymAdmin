package controller

import (
	"gym-admin/internal/models"
	"gym-admin/internal/service"
	"gym-admin/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CardController struct {
	service *service.CardService
}

func NewCardController() *CardController {
	return &CardController{
		service: service.NewCardService(),
	}
}

func (ctrl *CardController) CreateCard(c *gin.Context) {
	var req struct {
		UserID        int64   `json:"user_id" binding:"required"`
		CardTypeID    int64   `json:"card_type_id" binding:"required"`
		StartDate     string  `json:"start_date" binding:"required"`
		EndDate       string  `json:"end_date" binding:"required"`
		PurchasePrice float64 `json:"purchase_price" binding:"required"`
		Status        int8    `json:"status"`
		Remark        string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format, expected YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format, expected YYYY-MM-DD")
		return
	}

	card := &models.MembershipCard{
		UserID:        req.UserID,
		CardTypeID:    req.CardTypeID,
		StartDate:     startDate,
		EndDate:       endDate,
		PurchasePrice: req.PurchasePrice,
		Status:        req.Status,
		Remark:        req.Remark,
	}

	if err := ctrl.service.CreateCard(card); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, card)
}

func (ctrl *CardController) GetCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid card ID")
		return
	}

	card, err := ctrl.service.GetCard(id)
	if err != nil {
		response.NotFound(c, "Card not found")
		return
	}

	response.Success(c, card)
}

func (ctrl *CardController) ListCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.ParseInt(statusStr, 10, 8)
		statusVal := int8(s)
		status = &statusVal
	}

	var userID *int64
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		uid, _ := strconv.ParseInt(userIDStr, 10, 64)
		userID = &uid
	}

	cards, total, err := ctrl.service.ListCards(page, pageSize, status, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Failed to get cards: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 确保返回空数组而不是nil
	if cards == nil {
		cards = []models.MembershipCard{}
	}

	response.Success(c, gin.H{
		"list":      cards,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *CardController) UpdateCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid card ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UpdateCard(id, updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Card updated successfully", nil)
}

func (ctrl *CardController) DeleteCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid card ID")
		return
	}

	if err := ctrl.service.DeleteCard(id); err != nil {
		response.InternalServerError(c, "Failed to delete card")
		return
	}

	response.SuccessWithMessage(c, "Card deleted successfully", nil)
}

// RenewCard renews a membership card
func (ctrl *CardController) RenewCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid card ID")
		return
	}

	var req struct {
		Months     int     `json:"months" binding:"required,min=1"`
		Amount     float64 `json:"amount" binding:"required,min=0"`
		OperatorID int64   `json:"operator_id" binding:"required"`
		Remark     string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.RenewCard(id, req.Months, req.Amount, req.OperatorID, req.Remark); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Card renewed successfully", nil)
}

// FreezeCard freezes a membership card
func (ctrl *CardController) FreezeCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid card ID")
		return
	}

	var req struct {
		FreezeDays int    `json:"freeze_days" binding:"required,min=1"`
		OperatorID int64  `json:"operator_id" binding:"required"`
		Remark     string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.FreezeCard(id, req.FreezeDays, req.OperatorID, req.Remark); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Card frozen successfully", nil)
}

// UnfreezeCard unfreezes a membership card
func (ctrl *CardController) UnfreezeCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid card ID")
		return
	}

	var req struct {
		OperatorID int64  `json:"operator_id" binding:"required"`
		Remark     string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UnfreezeCard(id, req.OperatorID, req.Remark); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Card unfrozen successfully", nil)
}

// TransferCard transfers a membership card to another user
func (ctrl *CardController) TransferCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid card ID")
		return
	}

	var req struct {
		ToUserID    int64   `json:"to_user_id" binding:"required"`
		TransferFee float64 `json:"transfer_fee" binding:"min=0"`
		OperatorID  int64   `json:"operator_id" binding:"required"`
		Remark      string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.TransferCard(id, req.ToUserID, req.TransferFee, req.OperatorID, req.Remark); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Card transferred successfully", nil)
}

// GetCardOperations gets operation history of a card
func (ctrl *CardController) GetCardOperations(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid card ID")
		return
	}

	operations, err := ctrl.service.GetCardOperations(id)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, operations)
}
