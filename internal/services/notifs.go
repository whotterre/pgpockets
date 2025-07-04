package services

import (
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"

	"github.com/google/uuid"
)

type NotificationService interface {
	GetNotificationCount(userID uuid.UUID, includeRead bool) (int64, error)
	GetUnreadNotificationCount(userID uuid.UUID) (int64, error)

	GetNotificationByID(id uuid.UUID, userID uuid.UUID) (*models.Notification, error)
	GetNotificationsByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error)
	GetUnreadNotificationsByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error)

	MarkNotificationAsRead(id uuid.UUID, userID uuid.UUID) error
	MarkNotificationAsUnread(id uuid.UUID, userID uuid.UUID) error
	MarkAllNotificationsAsRead(userID uuid.UUID) (int64, error)

	DeleteNotification(id uuid.UUID, userID uuid.UUID) error
	DeleteAllNotifications(userID uuid.UUID) (int64, error)
	DeleteAllReadNotifications(userID uuid.UUID) (int64, error)
}

type notificationService struct {
	repo repositories.NotifRepo
}

func NewNotificationService(repo repositories.NotifRepo) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) GetNotificationCount(userID uuid.UUID, includeRead bool) (int64, error) {
	return s.repo.GetNotificationCount(userID, includeRead)
}

func (s *notificationService) GetUnreadNotificationCount(userID uuid.UUID) (int64, error) {
	return s.repo.GetUnreadNotificationCount(userID)
}

func (s *notificationService) GetNotificationByID(id uuid.UUID, userID uuid.UUID) (*models.Notification, error) {
	return s.repo.GetByID(id, userID)
}

func (s *notificationService) GetNotificationsByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

func (s *notificationService) GetUnreadNotificationsByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	return s.repo.GetUnreadByUserID(userID, limit, offset)
}

func (s *notificationService) MarkNotificationAsRead(id uuid.UUID, userID uuid.UUID) error {
	return s.repo.UpdateReadStatus(id, userID, true)
}

func (s *notificationService) MarkNotificationAsUnread(id uuid.UUID, userID uuid.UUID) error {
	return s.repo.UpdateReadStatus(id, userID, false)
}

func (s *notificationService) MarkAllNotificationsAsRead(userID uuid.UUID) (int64, error) {
	return s.repo.MarkAllAsReadByUserID(userID)
}

func (s *notificationService) DeleteNotification(id uuid.UUID, userID uuid.UUID) error {
	return s.repo.DeleteByID(id, userID)
}

func (s *notificationService) DeleteAllNotifications(userID uuid.UUID) (int64, error) {
	return s.repo.DeleteAllByUserID(userID)
}

func (s *notificationService) DeleteAllReadNotifications(userID uuid.UUID) (int64, error) {
	return s.repo.DeleteAllReadByUserID(userID)
}
