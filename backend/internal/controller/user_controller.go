package controller

import (
	"gym-admin/internal/models"
	"gym-admin/internal/service"
	"gym-admin/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		service: service.NewUserService(),
	}
}

// CreateUser creates a new user
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.CreateUser(&user); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, user)
}

// GetUser gets user by ID
func (ctrl *UserController) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := ctrl.service.GetUser(id)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, user)
}

// ListUsers lists users with pagination
func (ctrl *UserController) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.ParseInt(statusStr, 10, 8)
		statusVal := int8(s)
		status = &statusVal
	}

	users, total, err := ctrl.service.ListUsers(page, pageSize, status)
	if err != nil {
		response.InternalServerError(c, "Failed to get users")
		return
	}

	response.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateUser updates user information
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UpdateUser(id, updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User updated successfully", nil)
}

// DeleteUser deletes a user
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	if err := ctrl.service.DeleteUser(id); err != nil {
		response.InternalServerError(c, "Failed to delete user")
		return
	}

	response.SuccessWithMessage(c, "User deleted successfully", nil)
}

// GetUserStats gets user training statistics
func (ctrl *UserController) GetUserStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	stats, err := ctrl.service.GetUserStats(id)
	if err != nil {
		response.InternalServerError(c, "Failed to get user stats")
		return
	}

	response.Success(c, stats)
}

// FreezeUser freezes a user account
func (ctrl *UserController) FreezeUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req StatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.FreezeUser(id, req.Reason, req.OperatorID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User frozen successfully", nil)
}

// UnfreezeUser unfreezes a user account
func (ctrl *UserController) UnfreezeUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req StatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UnfreezeUser(id, req.Reason, req.OperatorID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User unfrozen successfully", nil)
}

// AddToBlacklist adds a user to blacklist
func (ctrl *UserController) AddToBlacklist(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req StatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.AddToBlacklist(id, req.Reason, req.OperatorID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User added to blacklist successfully", nil)
}

// RemoveFromBlacklist removes a user from blacklist
func (ctrl *UserController) RemoveFromBlacklist(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req StatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.RemoveFromBlacklist(id, req.Reason, req.OperatorID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User removed from blacklist successfully", nil)
}

// GetStatusLogs gets status change logs for a user
func (ctrl *UserController) GetStatusLogs(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logs, total, err := ctrl.service.GetStatusLogs(id, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to get status logs")
		return
	}

	response.Success(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetUserStatusSummary gets a summary of user counts by status
func (ctrl *UserController) GetUserStatusSummary(c *gin.Context) {
	summary, err := ctrl.service.GetUserStatusSummary()
	if err != nil {
		response.InternalServerError(c, "Failed to get status summary")
		return
	}

	response.Success(c, summary)
}

// BatchFreezeUsers freezes multiple users
func (ctrl *UserController) BatchFreezeUsers(c *gin.Context) {
	var req BatchStatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	result, err := ctrl.service.BatchFreezeUsers(req.UserIDs, req.Reason, req.OperatorID)
	if err != nil {
		response.InternalServerError(c, "Failed to freeze users")
		return
	}

	response.Success(c, result)
}

// BatchUnfreezeUsers unfreezes multiple users
func (ctrl *UserController) BatchUnfreezeUsers(c *gin.Context) {
	var req BatchStatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	result, err := ctrl.service.BatchUnfreezeUsers(req.UserIDs, req.Reason, req.OperatorID)
	if err != nil {
		response.InternalServerError(c, "Failed to unfreeze users")
		return
	}

	response.Success(c, result)
}

// Request structs
type StatusChangeRequest struct {
	Reason     string `json:"reason" binding:"required"`
	OperatorID *int64 `json:"operator_id"`
}

type BatchStatusChangeRequest struct {
	UserIDs    []int64 `json:"user_ids" binding:"required"`
	Reason     string  `json:"reason" binding:"required"`
	OperatorID *int64  `json:"operator_id"`
}
