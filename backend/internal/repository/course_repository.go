package repository

import (
	"gym-admin/internal/models"
	"gym-admin/pkg/database"

	"time"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository() *CourseRepository {
	return &CourseRepository{db: database.GetDB()}
}

func (r *CourseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *CourseRepository) GetByID(id int64) (*models.Course, error) {
	var course models.Course
	err := r.db.First(&course, id).Error
	return &course, err
}

func (r *CourseRepository) List(page, pageSize int, status *int8, coachID *int64, courseType *int8) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	query := r.db.Model(&models.Course{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if coachID != nil {
		query = query.Where("coach_id = ?", *coachID)
	}
	if courseType != nil {
		query = query.Where("course_type = ?", *courseType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("start_time DESC").Find(&courses).Error
	return courses, total, err
}

func (r *CourseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

func (r *CourseRepository) Delete(id int64) error {
	return r.db.Delete(&models.Course{}, id).Error
}

func (r *CourseRepository) GetByCoachID(coachID int64) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Where("coach_id = ?", coachID).Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) GetConflictingCourses(coachID int64, startTime, endTime time.Time, excludeCourseID int64) ([]models.Course, error) {
	var courses []models.Course
	query := r.db.Where("coach_id = ? AND status IN (?, ?) AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?) OR (end_time > ? AND end_time <= ?))",
		coachID, models.CourseStatusAvailable, models.CourseStatusFull, endTime, startTime, startTime, endTime, startTime, endTime)

	if excludeCourseID > 0 {
		query = query.Where("id <> ?", excludeCourseID)
	}

	err := query.Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) GetCompletedCoursesByCoachAndMonth(coachID int64, year, month int) ([]models.Course, error) {
	var courses []models.Course
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	err := r.db.Where("coach_id = ? AND status = ? AND start_time BETWEEN ? AND ?", coachID, models.CourseStatusCompleted, firstDay, lastDay).Find(&courses).Error
	return courses, err
}
