package controller

import (
"gym-admin/internal/models"
"gym-admin/internal/service"
"gym-admin/pkg/response"
"strconv"

"github.com/gin-gonic/gin"
)

type UserFaceController struct {
	service *service.UserFaceService
}

func NewUserFaceController() *UserFaceController {
	return &UserFaceController{
		service: service.NewUserFaceService(),
	}
}

// RegisterFace registers a new face for a user
func (ctrl *UserFaceController) RegisterFace(c *gin.Context) {
	var face models.UserFace
	if err := c.ShouldBindJSON(&face); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.RegisterFace(&face); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, face)
}

// GetFace gets face by ID
func (ctrl *UserFaceController) GetFace(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid face ID")
		return
	}

	face, err := ctrl.service.GetFace(id)
	if err != nil {
		response.NotFound(c, "Face not found")
		return
	}

	response.Success(c, face)
}

// GetUserFaces gets all faces for a user
func (ctrl *UserFaceController) GetUserFaces(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	faces, err := ctrl.service.GetUserFaces(userID)
	if err != nil {
		response.InternalServerError(c, "Failed to get user faces")
		return
	}

	response.Success(c, faces)
}

// GetMainFace gets the main face for a user
func (ctrl *UserFaceController) GetMainFace(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	face, err := ctrl.service.GetMainFace(userID)
	if err != nil {
		response.NotFound(c, "Main face not found")
		return
	}

	response.Success(c, face)
}

// ListFaces lists faces with pagination
func (ctrl *UserFaceController) ListFaces(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var userID *int64
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, _ := strconv.ParseInt(userIDStr, 10, 64)
		userID = &id
	}

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.ParseInt(statusStr, 10, 8)
		statusVal := int8(s)
		status = &statusVal
	}

	faces, total, err := ctrl.service.ListFaces(page, pageSize, userID, status)
	if err != nil {
		response.InternalServerError(c, "Failed to get faces")
		return
	}

	response.Success(c, gin.H{
"list":      faces,
"total":     total,
"page":      page,
"page_size": pageSize,
})
}

// UpdateFace updates face information
func (ctrl *UserFaceController) UpdateFace(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid face ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UpdateFace(id, updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Face updated successfully", nil)
}

// DeleteFace deletes a face
func (ctrl *UserFaceController) DeleteFace(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid face ID")
		return
	}

	if err := ctrl.service.DeleteFace(id); err != nil {
		response.InternalServerError(c, "Failed to delete face")
		return
	}

	response.SuccessWithMessage(c, "Face deleted successfully", nil)
}

// SetMainFace sets a face as the main face for a user
func (ctrl *UserFaceController) SetMainFace(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	faceID, err := strconv.ParseInt(c.Param("face_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid face ID")
		return
	}

	if err := ctrl.service.SetMainFace(userID, faceID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Main face set successfully", nil)
}

// DeleteUserFaces deletes all faces for a user
func (ctrl *UserFaceController) DeleteUserFaces(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	if err := ctrl.service.DeleteUserFaces(userID); err != nil {
		response.InternalServerError(c, "Failed to delete user faces")
		return
	}

	response.SuccessWithMessage(c, "User faces deleted successfully", nil)
}

// GetFaceStatistics gets face statistics
func (ctrl *UserFaceController) GetFaceStatistics(c *gin.Context) {
	stats, err := ctrl.service.GetFaceStatistics()
	if err != nil {
		response.InternalServerError(c, "Failed to get face statistics")
		return
	}

	response.Success(c, stats)
}
