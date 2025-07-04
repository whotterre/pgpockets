package handlers

import (
	"pgpockets/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationHandlers struct {
	notifService services.NotificationService
	logger       *zap.Logger
}

func NewNotificationHandlers(
	notifService services.NotificationService,
	logger *zap.Logger,
) *NotificationHandlers {
	return &NotificationHandlers{
		notifService: notifService,
		logger:       logger,
	}
}

func (h *NotificationHandlers) GetNotificationCount(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	includeRead := c.Query("includeRead") == "true"
	count, err := h.notifService.GetNotificationCount(userID, includeRead)
	if err != nil {
		h.logger.Error("Failed to get notification count", 
			zap.Error(err),
			zap.String("userID", userID.String()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get notification count",
		})
	}
	h.logger.Info("Successfully got notifications count",
		zap.String("userID", userID.String()),
		zap.Int64("count", count),
	)
	return c.JSON(fiber.Map{
		"count": count,
	})
}

func (h *NotificationHandlers) GetUnreadNotificationCount(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	count, err := h.notifService.GetUnreadNotificationCount(userID)
	if err != nil {
		h.logger.Error("Failed to get unread notification count",
			zap.Error(err),
			zap.String("userID", userID.String()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread notification count",
		})
	}
	h.logger.Info("Successfully got unread notification count",
		zap.String("userID", userID.String()),
		zap.Int64("count", count),
	)
	return c.JSON(fiber.Map{
		"count": count,
	})
}

func (h *NotificationHandlers) GetNotifications(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	limit, offset := getPaginationParams(c)
	notifications, err := h.notifService.GetNotificationsByUserID(userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get notifications",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get notifications",
		})
	}
	h.logger.Info("Successfully got notifications",
		zap.String("userID", userID.String()),
		zap.Int("count", len(notifications)),
	)
	return c.JSON(notifications)
}

func (h *NotificationHandlers) GetUnreadNotifications(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	limit, offset := getPaginationParams(c)
	notifications, err := h.notifService.GetUnreadNotificationsByUserID(userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get unread notifications",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread notifications",
		})
	}
	h.logger.Info("Successfully got unread notifications",
		zap.String("userID", userID.String()),
		zap.Int("count", len(notifications)),
	)
	return c.JSON(notifications)
}

func (h *NotificationHandlers) GetNotification(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Error("Invalid notification ID format",
			zap.Error(err),
			zap.String("input", c.Params("id")),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID",
		})
	}

	notification, err := h.notifService.GetNotificationByID(id, userID)
	if err != nil {
		h.logger.Error("Failed to get notification",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.String("notificationID", id.String()),
		)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}
	h.logger.Info("Successfully retrieved notification",
		zap.String("userID", userID.String()),
		zap.String("notificationID", id.String()),
	)
	return c.JSON(notification)
}

func (h *NotificationHandlers) MarkAsRead(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Error("Invalid notification ID format",
			zap.Error(err),
			zap.String("input", c.Params("id")),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID",
		})
	}

	if err := h.notifService.MarkNotificationAsRead(id, userID); err != nil {
		h.logger.Error("Failed to mark notification as read",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.String("notificationID", id.String()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notification as read",
		})
	}
	h.logger.Info("Successfully marked notification as read",
		zap.String("userID", userID.String()),
		zap.String("notificationID", id.String()),
	)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *NotificationHandlers) MarkAsUnread(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Error("Invalid notification ID format",
			zap.Error(err),
			zap.String("input", c.Params("id")),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID",
		})
	}

	if err := h.notifService.MarkNotificationAsUnread(id, userID); err != nil {
		h.logger.Error("Failed to mark notification as unread",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.String("notificationID", id.String()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notification as unread",
		})
	}
	h.logger.Info("Successfully marked notification as unread",
		zap.String("userID", userID.String()),
		zap.String("notificationID", id.String()),
	)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *NotificationHandlers) MarkAllAsRead(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	count, err := h.notifService.MarkAllNotificationsAsRead(userID)
	if err != nil {
		h.logger.Error("Failed to mark all notifications as read",
			zap.Error(err),
			zap.String("userID", userID.String()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark all notifications as read",
		})
	}
	h.logger.Info("Successfully marked all notifications as read",
		zap.String("userID", userID.String()),
		zap.Int64("count", count),
	)
	return c.JSON(fiber.Map{
		"count": count,
	})
}

func (h *NotificationHandlers) DeleteNotification(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		h.logger.Error("Invalid notification ID format",
			zap.Error(err),
			zap.String("input", c.Params("id")),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID",
		})
	}

	if err := h.notifService.DeleteNotification(id, userID); err != nil {
		h.logger.Error("Failed to delete notification",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.String("notificationID", id.String()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete notification",
		})
	}
	h.logger.Info("Successfully deleted notification",
		zap.String("userID", userID.String()),
		zap.String("notificationID", id.String()),
	)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *NotificationHandlers) DeleteAllNotifications(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	count, err := h.notifService.DeleteAllNotifications(userID)
	if err != nil {
		h.logger.Error("Failed to delete all notifications",
			zap.Error(err),
			zap.String("userID", userID.String()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete all notifications",
		})
	}
	h.logger.Info("Successfully deleted all notifications",
		zap.String("userID", userID.String()),
		zap.Int64("count", count),
	)
	return c.JSON(fiber.Map{
		"count": count,
	})
}

func (h *NotificationHandlers) DeleteAllReadNotifications(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.logger.Error("Invalid user", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user",
		})
	}

	count, err := h.notifService.DeleteAllReadNotifications(userID)
	if err != nil {
		h.logger.Error("Failed to delete read notifications",
			zap.Error(err),
			zap.String("userID", userID.String()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete read notifications",
		})
	}
	h.logger.Info("Successfully deleted read notifications",
		zap.String("userID", userID.String()),
		zap.Int64("count", count),
	)
	return c.JSON(fiber.Map{
		"count": count,
	})
}

// Helper functions
func getUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "invalid user")
	}
	return uuid.Parse(userIDStr)
}

func getPaginationParams(c *fiber.Ctx) (limit, offset int) {
	limit, _ = strconv.Atoi(c.Query("limit", "10"))
	offset, _ = strconv.Atoi(c.Query("offset", "0"))
	return limit, offset
}