package service

import (
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
)

type CoachPerformanceService struct {
	repo       *repository.CoachPerformanceRepository
	courseRepo *repository.CourseRepository
	ratingRepo *repository.CourseRatingRepository
}

func NewCoachPerformanceService() *CoachPerformanceService {
	return &CoachPerformanceService{
		repo:       repository.NewCoachPerformanceRepository(),
		courseRepo: repository.NewCourseRepository(),
		ratingRepo: repository.NewCourseRatingRepository(),
	}
}

// UpdateCoachPerformance updates a coach's performance for a given month
func (s *CoachPerformanceService) UpdateCoachPerformance(coachID int64, year, month int) error {
	// Get all completed courses for the coach in the given month
	courses, err := s.courseRepo.GetCompletedCoursesByCoachAndMonth(coachID, year, month)
	if err != nil {
		return err
	}

	var totalCourses int
	var totalHours float64
	var totalIncome float64
	var totalRating float64
	var ratedCourses int

	for _, course := range courses {
		totalCourses++
		totalHours += course.EndTime.Sub(course.StartTime).Hours()
		totalIncome += course.Price

		if course.RatingCount > 0 {
			totalRating += course.AverageRating * float64(course.RatingCount)
			ratedCourses += course.RatingCount
		}
	}

	avgRating := 0.0
	if ratedCourses > 0 {
		avgRating = totalRating / float64(ratedCourses)
	}

	performance, err := s.repo.GetByCoachAndMonth(coachID, year, month)
	if err != nil {
		// Create new record if not found
		performance = &models.CoachPerformance{
			CoachID: coachID,
			Year:    year,
			Month:   month,
		}
		performance.TotalCourses = totalCourses
		performance.TotalHours = totalHours
		performance.TotalIncome = totalIncome
		performance.AverageRating = avgRating
		return s.repo.Create(performance)
	}

	// Update existing record
	performance.TotalCourses = totalCourses
	performance.TotalHours = totalHours
	performance.TotalIncome = totalIncome
	performance.AverageRating = avgRating
	return s.repo.Update(performance)
}

// GetCoachPerformance gets a coach's performance for a given month
func (s *CoachPerformanceService) GetCoachPerformance(coachID int64, year, month int) (*models.CoachPerformance, error) {
	return s.repo.GetByCoachAndMonth(coachID, year, month)
}

// GetCoachRankings gets coach rankings based on different metrics
func (s *CoachPerformanceService) GetCoachRankings(year, month, limit int, sortBy string) ([]models.CoachPerformance, error) {
	switch sortBy {
	case "income":
		return s.repo.GetTopCoachesByIncome(year, month, limit)
	case "hours":
		return s.repo.GetTopCoachesByHours(year, month, limit)
	case "rating":
		return s.repo.GetTopCoachesByRating(year, month, limit)
	default:
		return s.repo.GetTopCoachesByIncome(year, month, limit) // Default sort by income
	}
}
