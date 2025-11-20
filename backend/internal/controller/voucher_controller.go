package controller

import (
"gym-admin/internal/service"
"gym-admin/pkg/response"
"strconv"

"github.com/gin-gonic/gin"
)

type VoucherController struct {
	service *service.VoucherService
}

func NewVoucherController() *VoucherController {
	return &VoucherController{
		service: service.NewVoucherService(),
	}
}

// VerifyVoucher verifies a voucher code
func (ctrl *VoucherController) VerifyVoucher(c *gin.Context) {
	var req struct {
		VoucherCode string `json:"voucher_code" binding:"required"`
		Platform    int8   `json:"platform" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Get operator ID from context (assuming it's set by auth middleware)
operatorID := c.GetInt64("userID")

voucher, err := ctrl.service.VerifyVoucher(req.VoucherCode, req.Platform, operatorID)
if err != nil {
response.Error(c, 400, err.Error())
return
}

response.Success(c, voucher)
}

// ListVouchers lists vouchers
func (ctrl *VoucherController) ListVouchers(c *gin.Context) {
page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

var status *int8
if statusStr := c.Query("status"); statusStr != "" {
s, _ := strconv.ParseInt(statusStr, 10, 8)
statusVal := int8(s)
status = &statusVal
}

var platform *int8
if platformStr := c.Query("platform"); platformStr != "" {
p, _ := strconv.ParseInt(platformStr, 10, 8)
platformVal := int8(p)
platform = &platformVal
}

vouchers, total, err := ctrl.service.ListVouchers(page, pageSize, status, platform)
if err != nil {
response.InternalServerError(c, "Failed to get vouchers")
return
}

response.Success(c, gin.H{
"list":      vouchers,
"total":     total,
"page":      page,
"page_size": pageSize,
})
}

// GetVoucher gets a voucher by ID
func (ctrl *VoucherController) GetVoucher(c *gin.Context) {
id, err := strconv.ParseInt(c.Param("id"), 10, 64)
if err != nil {
response.BadRequest(c, "Invalid voucher ID")
return
}

voucher, err := ctrl.service.GetVoucher(id)
if err != nil {
response.NotFound(c, "Voucher not found")
return
}

response.Success(c, voucher)
}
