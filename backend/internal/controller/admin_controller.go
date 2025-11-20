package controller

import (
	"gym-admin/internal/scheduler"
	"gym-admin/internal/service"
	"gym-admin/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	cardService *service.CardService
	scheduler   *scheduler.CardExpiryScheduler
}

func NewAdminController(scheduler *scheduler.CardExpiryScheduler) *AdminController {
	return &AdminController{
		cardService: service.NewCardService(),
		scheduler:   scheduler,
	}
}

// TriggerExpiryCheck manually triggers card expiry check
func (ctrl *AdminController) TriggerExpiryCheck(c *gin.Context) {
	if ctrl.scheduler != nil {
		ctrl.scheduler.RunNow()
		response.SuccessWithMessage(c, "Card expiry check triggered successfully", nil)
	} else {
		response.Error(c, 500, "Scheduler not initialized")
	}
}

// GetSchedulerStatus gets scheduler status
func (ctrl *AdminController) GetSchedulerStatus(c *gin.Context) {
	if ctrl.scheduler != nil {
		status := ctrl.scheduler.GetStatus()
		response.Success(c, gin.H{
			"status": status,
		})
	} else {
		response.Error(c, 500, "Scheduler not initialized")
	}
}

// GetExpiringCards gets cards expiring soon
func (ctrl *AdminController) GetExpiringCards(c *gin.Context) {
	days := 7 // Default 7 days
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	cards, err := ctrl.cardService.GetExpiringCards(days)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"days":  days,
		"count": len(cards),
		"cards": cards,
	})
}

// GetExpiredCards gets expired cards
func (ctrl *AdminController) GetExpiredCards(c *gin.Context) {
	cards, err := ctrl.cardService.GetExpiredCards()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"count": len(cards),
		"cards": cards,
	})
}
