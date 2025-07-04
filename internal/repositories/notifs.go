package repositories

import (
	"pgpockets/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotifRepo interface {
    GetNotificationCount(userID uuid.UUID, includeRead bool) (int64, error)
    GetUnreadNotificationCount(userID uuid.UUID) (int64, error)

    GetByID(id uuid.UUID, userID uuid.UUID) (*models.Notification, error)
    GetByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error)
    GetUnreadByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error)

    UpdateReadStatus(id uuid.UUID, userID uuid.UUID, isRead bool) error
    MarkAllAsReadByUserID(userID uuid.UUID) (int64, error)

    DeleteByID(id uuid.UUID, userID uuid.UUID) error
    DeleteAllByUserID(userID uuid.UUID) (int64, error)
    DeleteAllReadByUserID(userID uuid.UUID) (int64, error)
}

type notifRepo struct {
    db *gorm.DB
}

func NewNotifRepo(db *gorm.DB) NotifRepo {
    return &notifRepo{db: db}
}

func (r *notifRepo) GetNotificationCount(userID uuid.UUID, includeRead bool) (int64, error) {
    var count int64
    query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)
    if !includeRead {
        query = query.Where("is_read = ?", false)
    }
    if err := query.Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}

func (r *notifRepo) GetUnreadNotificationCount(userID uuid.UUID) (int64, error) {
    return r.GetNotificationCount(userID, false)
}

func (r *notifRepo) GetByID(id uuid.UUID, userID uuid.UUID) (*models.Notification, error) {
    var notification models.Notification
    if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&notification).Error; err != nil {
        return nil, err
    }
    return &notification, nil
}

func (r *notifRepo) GetByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
    var notifications []models.Notification
    if err := r.db.Where("user_id = ?", userID).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&notifications).Error; err != nil {
        return nil, err
    }
    return notifications, nil
}

func (r *notifRepo) GetUnreadByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
    var notifications []models.Notification
    if err := r.db.Where("user_id = ? AND is_read = ?", userID, false).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&notifications).Error; err != nil {
        return nil, err
    }
    return notifications, nil
}

func (r *notifRepo) UpdateReadStatus(id uuid.UUID, userID uuid.UUID, isRead bool) error {
    result := r.db.Model(&models.Notification{}).
	Where("id = ? AND user_id = ?", id, userID).
	Update("is_read", isRead)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}

func (r *notifRepo) MarkAllAsReadByUserID(userID uuid.UUID) (int64, error) {
    result := r.db.Model(&models.Notification{}).
	Where("user_id = ? AND is_read = ?", userID, false).
	Update("is_read", true)
    return result.RowsAffected, result.Error
}

func (r *notifRepo) DeleteByID(id uuid.UUID, userID uuid.UUID) error {
    result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{})
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}

func (r *notifRepo) DeleteAllByUserID(userID uuid.UUID) (int64, error) {
    result := r.db.Where("user_id = ?", userID).Delete(&models.Notification{})
    return result.RowsAffected, result.Error
}

func (r *notifRepo) DeleteAllReadByUserID(userID uuid.UUID) (int64, error) {
    result := r.db.Where("user_id = ? AND is_read = ?", userID, true).Delete(&models.Notification{})
    return result.RowsAffected, result.Error
}