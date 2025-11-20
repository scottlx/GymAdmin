package service

import (
"errors"
"gym-admin/internal/models"
"gym-admin/internal/repository"
"time"
)

type BookingService struct {
	repo       *repository.BookingRepository
	courseRepo *repository.CourseRepository
}

func NewBookingService() *BookingService {
	return &BookingService{
		repo:       repository.NewBookingRepository(),
		courseRepo: repository.NewCourseRepository(),
	}
}

func (s *BookingService) CreateBooking(booking *models.Booking) error {
	course, err := s.courseRepo.GetByID(booking.CourseID)
	if err != nil {
		return errors.New("course not found")
	}

	if course.CurrentCount >= course.MaxCapacity {
		return errors.New("course is full")
	}

	// Check for user's time conflicts with other bookings
if err := s.checkUserBookingConflicts(booking.UserID, course.StartTime, course.EndTime, 0); err != nil {
return err
}

booking.BookedAt = time.Now()
booking.Status = models.BookingStatusBooked

if err := s.repo.Create(booking); err != nil {
return err
}

// Update course current count
course.CurrentCount++
if course.CurrentCount >= course.MaxCapacity {
course.Status = models.CourseStatusFull
}
return s.courseRepo.Update(course)
}

func (s *BookingService) GetBooking(id int64) (*models.Booking, error) {
return s.repo.GetByID(id)
}

func (s *BookingService) ListBookings(page, pageSize int, userID *int64, courseID *int64, status *int8) ([]models.Booking, int64, error) {
if page < 1 {
page = 1
}
if pageSize < 1 || pageSize > 100 {
pageSize = 10
}
return s.repo.List(page, pageSize, userID, courseID, status)
}

func (s *BookingService) UpdateBooking(id int64, updates map[string]interface{}) error {
booking, err := s.repo.GetByID(id)
if err != nil {
return errors.New("booking not found")
}

// Update fields
if status, ok := updates["status"].(float64); ok {
booking.Status = int8(status)
}
if remark, ok := updates["remark"].(string); ok {
booking.Remark = remark
}

return s.repo.Update(booking)
}

func (s *BookingService) CancelBooking(id int64) error {
booking, err := s.repo.GetByID(id)
if err != nil {
return errors.New("booking not found")
}

if booking.Status == models.BookingStatusCancelled {
return errors.New("booking already cancelled")
}

booking.Status = models.BookingStatusCancelled
now := time.Now()
booking.CancelledAt = &now

if err := s.repo.Update(booking); err != nil {
return err
}

// Update course current count
course, err := s.courseRepo.GetByID(booking.CourseID)
if err == nil && course.CurrentCount > 0 {
course.CurrentCount--
if course.CurrentCount < course.MaxCapacity {
course.Status = models.CourseStatusAvailable
}
s.courseRepo.Update(course)
}

return nil
}

func (s *BookingService) CompleteBooking(id int64) error {
booking, err := s.repo.GetByID(id)
if err != nil {
return errors.New("booking not found")
}

if booking.Status == models.BookingStatusCompleted {
return errors.New("booking already completed")
}

booking.Status = models.BookingStatusCompleted
now := time.Now()
booking.CheckedInAt = &now

return s.repo.Update(booking)
}

func (s *BookingService) checkUserBookingConflicts(userID int64, startTime, endTime time.Time, excludeBookingID int64) error {
conflictingBookings, err := s.repo.GetConflictingBookings(userID, startTime, endTime, excludeBookingID)
if err != nil {
return err
}
if len(conflictingBookings) > 0 {
return errors.New("user has conflicting bookings")
}
return nil
}
