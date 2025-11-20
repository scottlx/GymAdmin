package service

import (
"errors"
"fmt"
"gym-admin/internal/models"
"gym-admin/internal/repository"
"time"
)

type UserFaceService struct {
	repo     *repository.UserFaceRepository
	userRepo *repository.UserRepository
}

func NewUserFaceService() *UserFaceService {
	return &UserFaceService{
		repo:     repository.NewUserFaceRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserFaceService) RegisterFace(face *models.UserFace) error {
	// Verify user exists
	_, err := s.userRepo.GetByID(face.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Generate face ID if not provided
	if face.FaceID == "" {
		face.FaceID = s.generateFaceID(face.UserID)
	}

	// Set default values
	face.Status = models.FaceStatusNormal
	face.RegisteredAt = time.Now()

	// Check if this is the first face for the user
	count, err := s.repo.CountByUser(face.UserID)
	if err != nil {
		return err
	}
	if count == 0 {
		face.IsMain = true // First face is main by default
	}

	return s.repo.Create(face)
}

func (s *UserFaceService) GetFace(id int64) (*models.UserFace, error) {
	return s.repo.GetByID(id)
}

func (s *UserFaceService) GetFaceByFaceID(faceID string) (*models.UserFace, error) {
	return s.repo.GetByFaceID(faceID)
}

func (s *UserFaceService) GetUserFaces(userID int64) ([]models.UserFace, error) {
	return s.repo.GetByUserID(userID)
}

func (s *UserFaceService) GetMainFace(userID int64) (*models.UserFace, error) {
	return s.repo.GetMainFace(userID)
}

func (s *UserFaceService) ListFaces(page, pageSize int, userID *int64, status *int8) ([]models.UserFace, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize, userID, status)
}

func (s *UserFaceService) UpdateFace(id int64, updates map[string]interface{}) error {
	face, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("face not found")
	}

	// Update fields
	if faceImageURL, ok := updates["face_image_url"].(string); ok {
		face.FaceImageURL = faceImageURL
	}
	if quality, ok := updates["quality"].(float64); ok {
		face.Quality = quality
	}
	if remark, ok := updates["remark"].(string); ok {
		face.Remark = remark
	}

	return s.repo.Update(face)
}

func (s *UserFaceService) DeleteFace(id int64) error {
	face, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("face not found")
	}

	// If deleting main face, set another face as main
	if face.IsMain {
		faces, err := s.repo.GetByUserID(face.UserID)
		if err == nil && len(faces) > 1 {
			// Find another face to set as main
			for _, f := range faces {
				if f.ID != id {
					s.repo.SetMainFace(face.UserID, f.ID)
					break
				}
			}
		}
	}

	return s.repo.Delete(id)
}

func (s *UserFaceService) SetMainFace(userID, faceID int64) error {
	// Verify face exists and belongs to user
	face, err := s.repo.GetByID(faceID)
	if err != nil {
		return errors.New("face not found")
	}
	if face.UserID != userID {
		return errors.New("face does not belong to this user")
	}

	return s.repo.SetMainFace(userID, faceID)
}

func (s *UserFaceService) DeleteUserFaces(userID int64) error {
	return s.repo.DeleteByUserID(userID)
}

func (s *UserFaceService) GetFaceStatistics() (map[string]interface{}, error) {
	totalCount, err := s.repo.GetTotalCount()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_faces": totalCount,
	}, nil
}

// generateFaceID generates a unique face ID
func (s *UserFaceService) generateFaceID(userID int64) string {
	return fmt.Sprintf("FACE_%d_%d", userID, time.Now().Unix())
}
