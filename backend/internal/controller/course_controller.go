package controller

import (
	"gym-admin/internal/models"
	"gym-admin/internal/service"
	"gym-admin/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CourseController struct {
	service *service.CourseService
}

func NewCourseController() *CourseController {
	return &CourseController{
		service: service.NewCourseService(),
	}
}

func (ctrl *CourseController) CreateCourse(c *gin.Context) {
	var req struct {
		CoachID     int64   `json:"coach_id" binding:"required"`
		CourseName  string  `json:"course_name" binding:"required"`
		CourseType  int8    `json:"course_type" binding:"required"`
		StartTime   string  `json:"start_time" binding:"required"`
		EndTime     string  `json:"end_time" binding:"required"`
		MaxCapacity int     `json:"max_capacity" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		Status      int8    `json:"status"`
		Description string  `json:"description"`
		Remark      string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// Parse times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		response.BadRequest(c, "Invalid start_time format, expected RFC3339")
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		response.BadRequest(c, "Invalid end_time format, expected RFC3339")
		return
	}

	course := &models.Course{
		CoachID:     req.CoachID,
		CourseName:  req.CourseName,
		CourseType:  req.CourseType,
		StartTime:   startTime,
		EndTime:     endTime,
		MaxCapacity: req.MaxCapacity,
		Price:       req.Price,
		Status:      req.Status,
		Description: req.Description,
		Remark:      req.Remark,
	}

	if err := ctrl.service.CreateCourse(course); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, course)
}

func (ctrl *CourseController) GetCourse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid course ID")
		return
	}

	course, err := ctrl.service.GetCourse(id)
	if err != nil {
		response.NotFound(c, "Course not found")
		return
	}

	response.Success(c, course)
}

func (ctrl *CourseController) ListCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.ParseInt(statusStr, 10, 8)
		statusVal := int8(s)
		status = &statusVal
	}

	var coachID *int64
	if coachIDStr := c.Query("coach_id"); coachIDStr != "" {
		cid, _ := strconv.ParseInt(coachIDStr, 10, 64)
		coachID = &cid
	}

	var courseType *int8
	if courseTypeStr := c.Query("course_type"); courseTypeStr != "" {
		ct, _ := strconv.ParseInt(courseTypeStr, 10, 8)
		courseTypeVal := int8(ct)
		courseType = &courseTypeVal
	}

	courses, total, err := ctrl.service.ListCourses(page, pageSize, status, coachID, courseType)
	if err != nil {
		response.InternalServerError(c, "Failed to get courses")
		return
	}

	response.Success(c, gin.H{
		"list":      courses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *CourseController) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid course ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.UpdateCourse(id, updates); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Course updated successfully", nil)
}

func (ctrl *CourseController) DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid course ID")
		return
	}

	if err := ctrl.service.DeleteCourse(id); err != nil {
		response.InternalServerError(c, "Failed to delete course")
		return
	}

	response.SuccessWithMessage(c, "Course deleted successfully", nil)
}
