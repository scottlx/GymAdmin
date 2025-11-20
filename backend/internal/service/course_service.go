package service

import (
	"errors"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
	"time"
)

type CourseService struct {
	repo        *repository.CourseRepository
	leaveRepo   *repository.CoachLeaveRepository
	availRepo   *repository.CoachAvailabilityRepository
	bookingRepo *repository.BookingRepository
}

func NewCourseService() *CourseService {
	return &CourseService{
		repo:        repository.NewCourseRepository(),
		leaveRepo:   repository.NewCoachLeaveRepository(),
		availRepo:   repository.NewCoachAvailabilityRepository(),
		bookingRepo: repository.NewBookingRepository(),
	}
}

func (s *CourseService) CreateCourse(course *models.Course) error {
	if course.Status == 0 {
		course.Status = 1 // Default status: available
	}
	if course.CurrentCount == 0 {
		course.CurrentCount = 0
	}

	// Check for conflicts
	if err := s.checkForConflicts(course.CoachID, course.StartTime, course.EndTime, 0); err != nil {
		return err
	}

	return s.repo.Create(course)
}

func (s *CourseService) GetCourse(id int64) (*models.Course, error) {
	return s.repo.GetByID(id)
}

func (s *CourseService) ListCourses(page, pageSize int, status *int8, coachID *int64, courseType *int8) ([]models.Course, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize, status, coachID, courseType)
}

func (s *CourseService) UpdateCourse(id int64, updates map[string]interface{}) error {
	course, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("course not found")
	}

	originalStartTime := course.StartTime
	originalEndTime := course.EndTime
	originalCoachID := course.CoachID

	// Update fields
	if coachID, ok := updates["coach_id"].(float64); ok {
		course.CoachID = int64(coachID)
	}
	if courseName, ok := updates["course_name"].(string); ok {
		course.CourseName = courseName
	}
	if courseType, ok := updates["course_type"].(float64); ok {
		course.CourseType = int8(courseType)
	}
	if startTime, ok := updates["start_time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			course.StartTime = t
		}
	}
	if endTime, ok := updates["end_time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			course.EndTime = t
		}
	}
	if maxCapacity, ok := updates["max_capacity"].(float64); ok {
		course.MaxCapacity = int(maxCapacity)
	}
	if currentCount, ok := updates["current_count"].(float64); ok {
		course.CurrentCount = int(currentCount)
	}
	if price, ok := updates["price"].(float64); ok {
		course.Price = price
	}
	if status, ok := updates["status"].(float64); ok {
		course.Status = int8(status)
	}
	if description, ok := updates["description"].(string); ok {
		course.Description = description
	}
	if remark, ok := updates["remark"].(string); ok {
		course.Remark = remark
	}

	// Check for conflicts only if time or coach changed
	if originalCoachID != course.CoachID || !originalStartTime.Equal(course.StartTime) || !originalEndTime.Equal(course.EndTime) {
		if err := s.checkForConflicts(course.CoachID, course.StartTime, course.EndTime, id); err != nil {
			return err
		}
	}

	return s.repo.Update(course)
}

func (s *CourseService) DeleteCourse(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("course not found")
	}
	return s.repo.Delete(id)
}

func (s *CourseService) GetCoursesByCoachID(coachID int64) ([]models.Course, error) {
	return s.repo.GetByCoachID(coachID)
}

// checkForConflicts checks for scheduling conflicts for a coach
func (s *CourseService) checkForConflicts(coachID int64, startTime, endTime time.Time, currentCourseID int64) error {
	// Check for conflicts with other courses
	conflictingCourses, err := s.repo.GetConflictingCourses(coachID, startTime, endTime, currentCourseID)
	if err != nil {
		return err
	}
	if len(conflictingCourses) > 0 {
		return errors.New("coach has conflicting courses scheduled")
	}

	// Check for conflicts with coach leaves
	conflictingLeaves, err := s.leaveRepo.GetConflictingLeaves(coachID, startTime, endTime)
	if err != nil {
		return err
	}
	if len(conflictingLeaves) > 0 {
		return errors.New("coach is on leave during this time")
	}

	// Check for conflicts with coach availability (simplified check for now)
	// This would typically involve more complex logic to match specific availability slots
	// For now, we assume if there are no leaves or other courses, the coach is available.
	// A more robust implementation would check if the course time falls within defined availability.

	return nil
}
