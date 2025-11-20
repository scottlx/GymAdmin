package service

import (
	"gym-admin/internal/models"
	"gym-admin/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		repo: repository.NewNotificationRepository(),
	}
}

func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	return s.repo.Create(notification)
}

func (s *NotificationService) GetNotification(id int64) (*models.Notification, error) {
	return s.repo.GetByID(id)
}

func (s *NotificationService) ListNotifications(userID int64, isRead *int8, page, pageSize int) ([]models.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.List(userID, isRead, page, pageSize)
}

func (s *NotificationService) MarkAsRead(id int64) error {
	return s.repo.MarkAsRead(id)
}

func (s *NotificationService) MarkAllAsRead(userID int64) error {
	return s.repo.MarkAllAsRead(userID)
}

func (s *NotificationService) DeleteNotification(id int64) error {
	return s.repo.Delete(id)
}

func (s *NotificationService) GetUnreadCount(userID int64) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}
