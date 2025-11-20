package controller

import (
	"gym-admin/internal/models"
	"gym-admin/internal/service"
	"gym-admin/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CoachCertificationController struct {
	service *service.CoachCertificationService
}

func NewCoachCertificationController() *CoachCertificationController {
	return &CoachCertificationController{
		service: service.NewCoachCertificationService(),
	}
}

// CreateCertification creates a new coach certification
func (ctrl *CoachCertificationController) CreateCertification(c *gin.Context) {
	var certification models.CoachCertification
	if err := c.ShouldBindJSON(&certification); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.CreateCertification(&certification); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, certification)
}

// GetCertification gets a certification by ID
func (ctrl *CoachCertificationController) GetCertification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid certification ID")
		return
	}

	certification, err := ctrl.service.GetCertification(id)
	if err != nil {
		response.NotFound(c, "Certification not found")
		return
	}

	response.Success(c, certification)
}

// GetCoachCertifications gets all certifications for a coach
func (ctrl *CoachCertificationController) GetCoachCertifications(c *gin.Context) {
	coachID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid coach ID")
		return
	}

	certifications, err := ctrl.service.GetCoachCertifications(coachID)
	if err != nil {
		response.InternalServerError(c, "Failed to get coach certifications")
		return
	}

	response.Success(c, certifications)
}

// ListCertifications lists certifications with pagination
func (ctrl *CoachCertificationController) ListCertifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.ParseInt(statusStr, 10, 8)
		statusVal := int8(s)
		status = &statusVal
	}

	certifications, total, err := ctrl.service.ListCertifications(page, pageSize, status)
	if err != nil {
		response.InternalServerError(c, "Failed to get certifications")
		return
	}

	response.Success(c, gin.H{
		"list":      certifications,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateCertification updates a certification
func (ctrl *CoachCertificationController) UpdateCertification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid certification ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UpdateCertification(id, updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Certification updated successfully", nil)
}

// DeleteCertification deletes a certification
func (ctrl *CoachCertificationController) DeleteCertification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid certification ID")
		return
	}

	if err := ctrl.service.DeleteCertification(id); err != nil {
		response.InternalServerError(c, "Failed to delete certification")
		return
	}

	response.SuccessWithMessage(c, "Certification deleted successfully", nil)
}

// ApproveCertification approves a certification
func (ctrl *CoachCertificationController) ApproveCertification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid certification ID")
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Get reviewer ID from context (assuming it's set by auth middleware)
	reviewerID := c.GetInt64("userID")

	if err := ctrl.service.ApproveCertification(id, reviewerID, req.Notes); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Certification approved successfully", nil)
}

// RejectCertification rejects a certification
func (ctrl *CoachCertificationController) RejectCertification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid certification ID")
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Get reviewer ID from context (assuming it's set by auth middleware)
	reviewerID := c.GetInt64("userID")

	if err := ctrl.service.RejectCertification(id, reviewerID, req.Notes); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Certification rejected successfully", nil)
}
