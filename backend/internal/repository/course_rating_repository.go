package repository

import (
"gym-admin/internal/models"
"gym-admin/pkg/database"

"gorm.io/gorm"
)

type CourseRatingRepository struct {
	db *gorm.DB
}

func NewCourseRatingRepository() *CourseRatingRepository {
	return &CourseRatingRepository{db: database.GetDB()}
}

func (r *CourseRatingRepository) Create(rating *models.CourseRating) error {
	return r.db.Create(rating).Error
}

func (r *CourseRatingRepository) GetByID(id int64) (*models.CourseRating, error) {
	var rating models.CourseRating
	err := r.db.First(&rating, id).Error
	return &rating, err
}

func (r *CourseRatingRepository) GetByBookingID(bookingID int64) (*models.CourseRating, error) {
	var rating models.CourseRating
	err := r.db.Where("booking_id = ?", bookingID).First(&rating).Error
	return &rating, err
}

func (r *CourseRatingRepository) GetByCourseID(courseID int64) ([]models.CourseRating, error) {
	var ratings []models.CourseRating
	err := r.db.Where("course_id = ?", courseID).Find(&ratings).Error
	return ratings, err
}

func (r *CourseRatingRepository) Update(rating *models.CourseRating) error {
	return r.db.Save(rating).Error
}

func (r *CourseRatingRepository) Delete(id int64) error {
	return r.db.Delete(&models.CourseRating{}, id).Error
}

func (r *CourseRatingRepository) GetAverageRating(courseID int64) (float64, int64, error) {
	var avgRating float64
	var count int64

	// Calculate average rating
	err := r.db.Model(&models.CourseRating{}).Where("course_id = ?", courseID).Select("COALESCE(AVG(rating), 0)").Row().Scan(&avgRating)
	if err != nil {
		return 0, 0, err
	}

	// Count total ratings
	err = r.db.Model(&models.CourseRating{}).Where("course_id = ?", courseID).Count(&count).Error
	if err != nil {
		return 0, 0, err
	}

	return avgRating, count, nil
}
