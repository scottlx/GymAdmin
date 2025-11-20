package controller

import (
	"gym-admin/internal/service"
	"gym-admin/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	service *service.NotificationService
}

func NewNotificationController() *NotificationController {
	return &NotificationController{
		service: service.NewNotificationService(),
	}
}

// ListNotifications lists user notifications
func (ctrl *NotificationController) ListNotifications(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if userID == 0 {
		response.BadRequest(c, "user_id is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var isRead *int8
	if isReadStr := c.Query("is_read"); isReadStr != "" {
		val, _ := strconv.ParseInt(isReadStr, 10, 8)
		isReadVal := int8(val)
		isRead = &isReadVal
	}

	notifications, total, err := ctrl.service.ListNotifications(userID, isRead, page, pageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      notifications,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetNotification gets a notification by ID
func (ctrl *NotificationController) GetNotification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid notification ID")
		return
	}

	notification, err := ctrl.service.GetNotification(id)
	if err != nil {
		response.NotFound(c, "Notification not found")
		return
	}

	response.Success(c, notification)
}

// MarkAsRead marks a notification as read
func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid notification ID")
		return
	}

	if err := ctrl.service.MarkAsRead(id); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Notification marked as read", nil)
}

// MarkAllAsRead marks all notifications as read for a user
func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	var req struct {
		UserID int64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.MarkAllAsRead(req.UserID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "All notifications marked as read", nil)
}

// DeleteNotification deletes a notification
func (ctrl *NotificationController) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid notification ID")
		return
	}

	if err := ctrl.service.DeleteNotification(id); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Notification deleted successfully", nil)
}

// GetUnreadCount gets unread notification count
func (ctrl *NotificationController) GetUnreadCount(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if userID == 0 {
		response.BadRequest(c, "user_id is required")
		return
	}

	count, err := ctrl.service.GetUnreadCount(userID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"count": count,
	})
}
