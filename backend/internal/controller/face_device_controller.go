package controller

import (
"gym-admin/internal/models"
"gym-admin/internal/service"
"gym-admin/pkg/response"
"strconv"

"github.com/gin-gonic/gin"
)

type FaceDeviceController struct {
	service *service.FaceDeviceService
}

func NewFaceDeviceController() *FaceDeviceController {
	return &FaceDeviceController{
		service: service.NewFaceDeviceService(),
	}
}

// CreateDevice creates a new face recognition device
func (ctrl *FaceDeviceController) CreateDevice(c *gin.Context) {
	var device models.FaceDevice
	if err := c.ShouldBindJSON(&device); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.CreateDevice(&device); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, device)
}

// GetDevice gets device by ID
func (ctrl *FaceDeviceController) GetDevice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	device, err := ctrl.service.GetDevice(id)
	if err != nil {
		response.NotFound(c, "Device not found")
		return
	}

	response.Success(c, device)
}

// ListDevices lists devices with pagination
func (ctrl *FaceDeviceController) ListDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.ParseInt(statusStr, 10, 8)
		statusVal := int8(s)
		status = &statusVal
	}

	var deviceType *int8
	if typeStr := c.Query("device_type"); typeStr != "" {
		t, _ := strconv.ParseInt(typeStr, 10, 8)
		typeVal := int8(t)
		deviceType = &typeVal
	}

	devices, total, err := ctrl.service.ListDevices(page, pageSize, status, deviceType)
	if err != nil {
		response.InternalServerError(c, "Failed to get devices")
		return
	}

	response.Success(c, gin.H{
"list":      devices,
"total":     total,
"page":      page,
"page_size": pageSize,
})
}

// UpdateDevice updates device information
func (ctrl *FaceDeviceController) UpdateDevice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UpdateDevice(id, updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Device updated successfully", nil)
}

// DeleteDevice deletes a device
func (ctrl *FaceDeviceController) DeleteDevice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	if err := ctrl.service.DeleteDevice(id); err != nil {
		response.InternalServerError(c, "Failed to delete device")
		return
	}

	response.SuccessWithMessage(c, "Device deleted successfully", nil)
}

// EnableDevice enables a device
func (ctrl *FaceDeviceController) EnableDevice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	if err := ctrl.service.EnableDevice(id); err != nil {
		response.InternalServerError(c, "Failed to enable device")
		return
	}

	response.SuccessWithMessage(c, "Device enabled successfully", nil)
}

// DisableDevice disables a device
func (ctrl *FaceDeviceController) DisableDevice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	if err := ctrl.service.DisableDevice(id); err != nil {
		response.InternalServerError(c, "Failed to disable device")
		return
	}

	response.SuccessWithMessage(c, "Device disabled successfully", nil)
}

// GetDeviceStatusSummary gets device status summary
func (ctrl *FaceDeviceController) GetDeviceStatusSummary(c *gin.Context) {
	summary, err := ctrl.service.GetDeviceStatusSummary()
	if err != nil {
		response.InternalServerError(c, "Failed to get device status summary")
		return
	}

	response.Success(c, summary)
}

// UpdateDeviceHeartbeat updates device last online time
func (ctrl *FaceDeviceController) UpdateDeviceHeartbeat(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	if err := ctrl.service.UpdateLastOnline(id); err != nil {
		response.InternalServerError(c, "Failed to update device heartbeat")
		return
	}

	response.SuccessWithMessage(c, "Device heartbeat updated", nil)
}
