package controller

import (
	"gym-admin/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CheckInController struct {
	service *service.CheckInService
}

func NewCheckInController() *CheckInController {
	return &CheckInController{
		service: service.NewCheckInService(),
	}
}

// CheckIn handles user check-in
// @Summary User check-in
// @Tags CheckIn
// @Accept json
// @Produce json
// @Param body body CheckInRequest true "Check-in request"
// @Success 200 {object} Response
// @Router /api/v1/checkins [post]
func (c *CheckInController) CheckIn(ctx *gin.Context) {
	var req CheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkIn, err := c.service.CheckIn(req.UserID, req.CheckInType, req.DeviceID, req.CardID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "check-in successful",
		"data":    checkIn,
	})
}

// GetTodayCheckIn gets today's check-in record for a user
// @Summary Get today's check-in
// @Tags CheckIn
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} Response
// @Router /api/v1/users/:id/checkin/today [get]
func (c *CheckInController) GetTodayCheckIn(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	checkIn, err := c.service.GetTodayCheckIn(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "no check-in record found for today",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    checkIn,
	})
}

// ListCheckIns lists check-in records
// @Summary List check-in records
// @Tags CheckIn
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param user_id query int false "User ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} Response
// @Router /api/v1/checkins [get]
func (c *CheckInController) ListCheckIns(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	var userID *int64
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err == nil {
			userID = &id
		}
	}

	var startDate, endDate *time.Time
	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	checkIns, total, err := c.service.ListCheckIns(page, pageSize, userID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":      checkIns,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Request structs
type CheckInRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	CheckInType int8   `json:"check_in_type" binding:"required,oneof=1 2 3"` // 1-人脸识别，2-刷卡，3-手动签到
	DeviceID    string `json:"device_id"`
	CardID      *int64 `json:"card_id"`
}

// GetUserStats gets user training statistics
func (c *CheckInController) GetUserStats(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	stats, err := c.service.GetUserStats(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetDetailedStats gets detailed training statistics
func (c *CheckInController) GetDetailedStats(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	year, _ := strconv.Atoi(ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	stats, err := c.service.GetDetailedStats(userID, year)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetCheckInCalendar gets check-in calendar data
func (c *CheckInController) GetCheckInCalendar(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	year, _ := strconv.Atoi(ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(ctx.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	if month < 1 || month > 12 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	calendar, err := c.service.GetCheckInCalendar(userID, year, month)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    calendar,
	})
}

// RecalculateStats manually recalculates user statistics
func (c *CheckInController) RecalculateStats(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := c.service.RecalculateUserStats(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "statistics recalculated successfully",
	})
}
