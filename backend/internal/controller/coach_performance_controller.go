package controller

import (
	"gym-admin/internal/service"
	"gym-admin/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CoachPerformanceController struct {
	service *service.CoachPerformanceService
}

func NewCoachPerformanceController() *CoachPerformanceController {
	return &CoachPerformanceController{
		service: service.NewCoachPerformanceService(),
	}
}

// GetCoachPerformance retrieves the performance data for a specific coach.
func (ctrl *CoachPerformanceController) GetCoachPerformance(c *gin.Context) {
	coachID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid coach ID")
		return
	}

	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	performance, err := ctrl.service.GetCoachPerformance(coachID, year, month)
	if err != nil {
		response.InternalServerError(c, "Failed to get coach performance")
		return
	}

	response.Success(c, performance)
}

// UpdateCoachPerformance updates the performance data for a specific coach.
func (ctrl *CoachPerformanceController) UpdateCoachPerformance(c *gin.Context) {
	coachID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid coach ID")
		return
	}

	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	if err := ctrl.service.UpdateCoachPerformance(coachID, year, month); err != nil {
		response.InternalServerError(c, "Failed to update coach performance")
		return
	}

	response.SuccessWithMessage(c, "Coach performance updated successfully", nil)
}

// GetCoachRankings gets coach rankings based on different metrics
func (ctrl *CoachPerformanceController) GetCoachRankings(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort_by", "income") // income, hours, rating

	rankings, err := ctrl.service.GetCoachRankings(year, month, limit, sortBy)
	if err != nil {
		response.InternalServerError(c, "Failed to get coach rankings")
		return
	}

	response.Success(c, rankings)
}
