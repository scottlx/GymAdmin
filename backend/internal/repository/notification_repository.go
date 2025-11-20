package repository

import (
	"gym-admin/internal/models"
	"gym-admin/pkg/database"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{db: database.GetDB()}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetByID(id int64) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, id).Error
	return &notification, err
}

func (r *NotificationRepository) List(userID int64, isRead *int8, page, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notifications).Error

	if notifications == nil {
		notifications = []models.Notification{}
	}

	return notifications, total, err
}

func (r *NotificationRepository) MarkAsRead(id int64) error {
	now := gorm.Expr("NOW()")
	return r.db.Model(&models.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_read": 1,
		"read_at": now,
	}).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID int64) error {
	now := gorm.Expr("NOW()")
	return r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = 0", userID).Updates(map[string]interface{}{
		"is_read": 1,
		"read_at": now,
	}).Error
}

func (r *NotificationRepository) Delete(id int64) error {
	return r.db.Delete(&models.Notification{}, id).Error
}

func (r *NotificationRepository) GetUnreadCount(userID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = 0", userID).Count(&count).Error
	return count, err
}
