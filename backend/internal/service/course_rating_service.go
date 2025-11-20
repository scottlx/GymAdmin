package service

import (
	"errors"
	"fmt"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
)

type CourseRatingService struct {
	ratingRepo  *repository.CourseRatingRepository
	courseRepo  *repository.CourseRepository
	bookingRepo *repository.BookingRepository
}

func NewCourseRatingService() *CourseRatingService {
	return &CourseRatingService{
		ratingRepo:  repository.NewCourseRatingRepository(),
		courseRepo:  repository.NewCourseRepository(),
		bookingRepo: repository.NewBookingRepository(),
	}
}

func (s *CourseRatingService) CreateRating(rating *models.CourseRating) error {
	// Check if booking exists and is completed
	booking, err := s.bookingRepo.GetByID(rating.BookingID)
	if err != nil {
		return errors.New("booking not found")
	}
	if booking.Status != models.BookingStatusCompleted {
		return errors.New("course not completed yet")
	}
	if booking.UserID != rating.UserID {
		return errors.New("user does not match booking")
	}
	if booking.HasRated {
		return errors.New("booking already rated")
	}

	rating.CourseID = booking.CourseID
	if err := s.ratingRepo.Create(rating); err != nil {
		return err
	}

	// Update course average rating
	if err := s.updateCourseRating(booking.CourseID); err != nil {
		// Log error, but don't fail the request
		fmt.Println("Failed to update course rating:", err)
	}

	// Mark booking as rated
	booking.HasRated = true
	return s.bookingRepo.Update(booking)
}

func (s *CourseRatingService) GetRating(id int64) (*models.CourseRating, error) {
	return s.ratingRepo.GetByID(id)
}

func (s *CourseRatingService) GetCourseRatings(courseID int64) ([]models.CourseRating, error) {
	return s.ratingRepo.GetByCourseID(courseID)
}

func (s *CourseRatingService) updateCourseRating(courseID int64) error {
	avgRating, count, err := s.ratingRepo.GetAverageRating(courseID)
	if err != nil {
		return err
	}

	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return err
	}

	course.AverageRating = avgRating
	course.RatingCount = int(count)
	return s.courseRepo.Update(course)
}
