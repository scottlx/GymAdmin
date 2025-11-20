package repository

import (
"gym-admin/internal/models"
"gym-admin/pkg/database"
"time"

"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{db: database.GetDB()}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) GetByID(id int64) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.First(&booking, id).Error
	return &booking, err
}

func (r *BookingRepository) List(page, pageSize int, userID *int64, courseID *int64, status *int8) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if courseID != nil {
		query = query.Where("course_id = ?", *courseID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("booked_at DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *BookingRepository) Delete(id int64) error {
	return r.db.Delete(&models.Booking{}, id).Error
}

func (r *BookingRepository) GetByUserID(userID int64) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) GetByCourseID(courseID int64) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("course_id = ?", courseID).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) GetConflictingBookings(userID int64, startTime, endTime time.Time, excludeBookingID int64) ([]models.Booking, error) {
	var bookings []models.Booking
	query := r.db.Joins("JOIN courses ON courses.id = bookings.course_id").
		Where("bookings.user_id = ? AND bookings.status = ? AND ((courses.start_time < ? AND courses.end_time > ?) OR (courses.start_time >= ? AND courses.start_time < ?) OR (courses.end_time > ? AND courses.end_time <= ?))",
userID, models.BookingStatusBooked, endTime, startTime, startTime, endTime, startTime, endTime)

	if excludeBookingID > 0 {
		query = query.Where("bookings.id <> ?", excludeBookingID)
	}

	err := query.Find(&bookings).Error
	return bookings, err
}
