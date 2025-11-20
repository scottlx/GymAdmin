package service

import (
	"errors"
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
	"time"
)

type CoachCertificationService struct {
	certRepo  *repository.CoachCertificationRepository
	coachRepo *repository.CoachRepository
}

func NewCoachCertificationService() *CoachCertificationService {
	return &CoachCertificationService{
		certRepo:  repository.NewCoachCertificationRepository(),
		coachRepo: repository.NewCoachRepository(),
	}
}

// CreateCertification creates a new coach certification
func (s *CoachCertificationService) CreateCertification(certification *models.CoachCertification) error {
	// Verify coach exists
	_, err := s.coachRepo.GetByID(certification.CoachID)
	if err != nil {
		return errors.New("coach not found")
	}

	// Set default status
	certification.Status = models.CertificationStatusPending

	if err := s.certRepo.Create(certification); err != nil {
		return err
	}

	// Update coach's certification status to "pending"
	return s.updateCoachCertificationStatus(certification.CoachID)
}

// GetCertification gets a certification by ID
func (s *CoachCertificationService) GetCertification(id int64) (*models.CoachCertification, error) {
	return s.certRepo.GetByID(id)
}

// GetCoachCertifications gets all certifications for a coach
func (s *CoachCertificationService) GetCoachCertifications(coachID int64) ([]models.CoachCertification, error) {
	return s.certRepo.GetByCoachID(coachID)
}

// ListCertifications lists certifications with pagination
func (s *CoachCertificationService) ListCertifications(page, pageSize int, status *int8) ([]models.CoachCertification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.certRepo.List(page, pageSize, status)
}

// UpdateCertification updates a certification
func (s *CoachCertificationService) UpdateCertification(id int64, updates map[string]interface{}) error {
	certification, err := s.certRepo.GetByID(id)
	if err != nil {
		return errors.New("certification not found")
	}

	// Update fields
	if certificationName, ok := updates["certification_name"].(string); ok {
		certification.CertificationName = certificationName
	}
	if issuingOrg, ok := updates["issuing_org"].(string); ok {
		certification.IssuingOrg = issuingOrg
	}
	if issueDate, ok := updates["issue_date"].(string); ok {
		if t, err := time.Parse(time.RFC3339, issueDate); err == nil {
			certification.IssueDate = &t
		}
	}
	if expiryDate, ok := updates["expiry_date"].(string); ok {
		if t, err := time.Parse(time.RFC3339, expiryDate); err == nil {
			certification.ExpiryDate = &t
		}
	}
	if fileURL, ok := updates["file_url"].(string); ok {
		certification.FileURL = fileURL
	}

	return s.certRepo.Update(certification)
}

// DeleteCertification deletes a certification
func (s *CoachCertificationService) DeleteCertification(id int64) error {
	certification, err := s.certRepo.GetByID(id)
	if err != nil {
		return errors.New("certification not found")
	}

	if err := s.certRepo.Delete(id); err != nil {
		return err
	}

	// Update coach's certification status
	return s.updateCoachCertificationStatus(certification.CoachID)
}

// ApproveCertification approves a certification
func (s *CoachCertificationService) ApproveCertification(id int64, reviewerID int64, notes string) error {
	certification, err := s.certRepo.GetByID(id)
	if err != nil {
		return errors.New("certification not found")
	}

	if err := s.certRepo.Approve(id, reviewerID, notes); err != nil {
		return err
	}

	// Update coach's certification status
	return s.updateCoachCertificationStatus(certification.CoachID)
}

// RejectCertification rejects a certification
func (s *CoachCertificationService) RejectCertification(id int64, reviewerID int64, notes string) error {
	certification, err := s.certRepo.GetByID(id)
	if err != nil {
		return errors.New("certification not found")
	}

	if err := s.certRepo.Reject(id, reviewerID, notes); err != nil {
		return err
	}

	// Update coach's certification status
	return s.updateCoachCertificationStatus(certification.CoachID)
}

// updateCoachCertificationStatus updates the certification status of a coach based on their certifications
func (s *CoachCertificationService) updateCoachCertificationStatus(coachID int64) error {
	certifications, err := s.certRepo.GetByCoachID(coachID)
	if err != nil {
		return err
	}

	// Determine coach's certification status
	var coachStatus int8 = models.CoachCertificationStatusUnverified // Default to unverified

	for _, cert := range certifications {
		if cert.Status == models.CertificationStatusApproved {
			coachStatus = models.CoachCertificationStatusVerified
			break
		} else if cert.Status == models.CertificationStatusPending {
			coachStatus = models.CoachCertificationStatusPending
		}
	}

	// Update coach's certification status
	coach, err := s.coachRepo.GetByID(coachID)
	if err != nil {
		return err
	}

	coach.CertificationStatus = coachStatus
	return s.coachRepo.Update(coach)
}
