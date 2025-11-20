package controller

import (
"gym-admin/internal/models"
"gym-admin/internal/service"
"gym-admin/pkg/response"
"strconv"

"github.com/gin-gonic/gin"
)

type CourseRatingController struct {
	service *service.CourseRatingService
}

func NewCourseRatingController() *CourseRatingController {
	return &CourseRatingController{
		service: service.NewCourseRatingService(),
	}
}

// CreateRating creates a new rating for a course
func (ctrl *CourseRatingController) CreateRating(c *gin.Context) {
	var rating models.CourseRating
	if err := c.ShouldBindJSON(&rating); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := ctrl.service.CreateRating(&rating); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, rating)
}

// GetRating gets rating by ID
func (ctrl *CourseRatingController) GetRating(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid rating ID")
		return
	}

	rating, err := ctrl.service.GetRating(id)
	if err != nil {
		response.NotFound(c, "Rating not found")
		return
	}

	response.Success(c, rating)
}

// GetCourseRatings gets all ratings for a course
func (ctrl *CourseRatingController) GetCourseRatings(c *gin.Context) {
	courseID, err := strconv.ParseInt(c.Param("course_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid course ID")
		return
	}

	ratings, err := ctrl.service.GetCourseRatings(courseID)
	if err != nil {
		response.InternalServerError(c, "Failed to get course ratings")
		return
	}

	response.Success(c, ratings)
}
