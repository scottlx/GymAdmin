package service

import (
"errors"
"gym-admin/internal/models"
"gym-admin/internal/repository"
"time"
)

type CoachLeaveService struct {
	repo *repository.CoachLeaveRepository
}

func NewCoachLeaveService() *CoachLeaveService {
	return &CoachLeaveService{
		repo: repository.NewCoachLeaveRepository(),
	}
}

func (s *CoachLeaveService) RequestLeave(leave *models.CoachLeave) error {
	leave.Status = models.LeaveStatusPending
	return s.repo.Create(leave)
}

func (s *CoachLeaveService) GetLeave(id int64) (*models.CoachLeave, error) {
	return s.repo.GetByID(id)
}

func (s *CoachLeaveService) GetCoachLeaves(coachID int64) ([]models.CoachLeave, error) {
	return s.repo.GetByCoachID(coachID)
}

func (s *CoachLeaveService) ApproveLeave(id int64, approverID int64) error {
	leave, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("leave not found")
	}

	leave.Status = models.LeaveStatusApproved
	leave.ApproverID = &approverID
	return s.repo.Update(leave)
}

func (s *CoachLeaveService) RejectLeave(id int64, approverID int64) error {
	leave, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("leave not found")
	}

	leave.Status = models.LeaveStatusRejected
	leave.ApproverID = &approverID
	return s.repo.Update(leave)
}

func (s *CoachLeaveService) DeleteLeave(id int64) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("leave not found")
	}
	return s.repo.Delete(id)
}

func (s *CoachLeaveService) IsCoachOnLeave(coachID int64, t time.Time) (bool, error) {
	leaves, err := s.repo.GetConflictingLeaves(coachID, t, t)
	if err != nil {
		return false, err
	}
	return len(leaves) > 0, nil
}
