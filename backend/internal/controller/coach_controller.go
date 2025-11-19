package controller

import (
	"gym-admin/internal/models"
	"gym-admin/internal/service"
	"gym-admin/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CoachController struct {
	service *service.CoachService
}

func NewCoachController() *CoachController {
	return &CoachController{
		service: service.NewCoachService(),
	}
}

func (ctrl *CoachController) CreateCoach(c *gin.Context) {
	var coach models.Coach
	if err := c.ShouldBindJSON(&coach); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.CreateCoach(&coach); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, coach)
}

func (ctrl *CoachController) GetCoach(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid coach ID")
		return
	}

	coach, err := ctrl.service.GetCoach(id)
	if err != nil {
		response.NotFound(c, "Coach not found")
		return
	}

	response.Success(c, coach)
}

func (ctrl *CoachController) ListCoaches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.ParseInt(statusStr, 10, 8)
		statusVal := int8(s)
		status = &statusVal
	}

	coaches, total, err := ctrl.service.ListCoaches(page, pageSize, status)
	if err != nil {
		response.InternalServerError(c, "Failed to get coaches")
		return
	}

	response.Success(c, gin.H{
		"list":      coaches,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *CoachController) UpdateCoach(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid coach ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UpdateCoach(id, updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Coach updated successfully", nil)
}

func (ctrl *CoachController) DeleteCoach(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid coach ID")
		return
	}

	if err := ctrl.service.DeleteCoach(id); err != nil {
		response.InternalServerError(c, "Failed to delete coach")
		return
	}

	response.SuccessWithMessage(c, "Coach deleted successfully", nil)
}
