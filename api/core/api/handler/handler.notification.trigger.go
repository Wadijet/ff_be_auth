package handler

import (
	"fmt"
	"time"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/notification"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationTriggerHandler xử lý việc trigger notification
type NotificationTriggerHandler struct {
	router *notification.Router
	queue  *notification.Queue
}

// NewNotificationTriggerHandler tạo mới NotificationTriggerHandler
func NewNotificationTriggerHandler() (*NotificationTriggerHandler, error) {
	router, err := notification.NewRouter()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification router: %w", err)
	}

	queue, err := notification.NewQueue()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification queue: %w", err)
	}

	return &NotificationTriggerHandler{
		router: router,
		queue:  queue,
	}, nil
}

// TriggerNotificationRequest là request body để trigger notification
type TriggerNotificationRequest struct {
	EventType string                 `json:"eventType" validate:"required"`
	Payload   map[string]interface{} `json:"payload" validate:"required"`
}

// HandleTriggerNotification xử lý request trigger notification
func (h *NotificationTriggerHandler) HandleTriggerNotification(c fiber.Ctx) error {
	return SafeHandlerWrapper(c, func() error {
		var req TriggerNotificationRequest
		if err := c.Bind().Body(&req); err != nil {
			c.Status(common.StatusBadRequest).JSON(fiber.Map{
				"code":    common.ErrCodeValidationFormat.Code,
				"message": fmt.Sprintf("Dữ liệu gửi lên không đúng định dạng JSON. Chi tiết: %v", err),
				"status":  "error",
			})
			return nil
		}

		// Validate
		if req.EventType == "" {
			c.Status(common.StatusBadRequest).JSON(fiber.Map{
				"code":    common.ErrCodeValidationFormat.Code,
				"message": "eventType không được để trống",
				"status":  "error",
			})
			return nil
		}

		if req.Payload == nil {
			req.Payload = make(map[string]interface{})
		}

		// Tìm routes cho eventType
		routes, err := h.router.FindRoutes(c.Context(), req.EventType)
		if err != nil {
			c.Status(common.StatusInternalServerError).JSON(fiber.Map{
				"code":    common.ErrCodeBusinessOperation.Code,
				"message": fmt.Sprintf("Không thể tìm routes cho eventType '%s': %v", req.EventType, err),
				"status":  "error",
			})
			return nil
		}

		if len(routes) == 0 {
			c.JSON(map[string]interface{}{
				"message":   "Không có routing rule nào cho eventType này",
				"eventType": req.EventType,
				"queued":    0,
			})
			return nil
		}

		// Tạo queue items cho mỗi route
		queueItems := make([]*models.NotificationQueueItem, 0)
		channelService, err := services.NewNotificationChannelService()
		if err != nil {
			c.Status(common.StatusInternalServerError).JSON(fiber.Map{
				"code":    common.ErrCodeBusinessOperation.Code,
				"message": fmt.Sprintf("Không thể tạo channel service: %v", err),
				"status":  "error",
			})
			return nil
		}

		for _, route := range routes {
			// Lấy channel để biết recipients
			channel, err := channelService.FindOneById(c.Context(), route.ChannelID)
			if err != nil {
				// Log error nhưng tiếp tục với route khác
				continue
			}

			// Xác định recipients dựa trên channel type
			var recipients []string
			switch channel.ChannelType {
			case "email":
				recipients = channel.Recipients
			case "telegram":
				recipients = channel.ChatIDs
			case "webhook":
				if channel.WebhookURL != "" {
					recipients = []string{channel.WebhookURL}
				}
			default:
				continue
			}

			// Tạo queue item cho mỗi recipient
			for _, recipient := range recipients {
				queueItems = append(queueItems, &models.NotificationQueueItem{
					ID:             primitive.NewObjectID(),
					EventType:      req.EventType,
					OrganizationID: route.OrganizationID,
					ChannelID:      route.ChannelID,
					Recipient:      recipient,
					Payload:        req.Payload,
					Status:         "pending",
					RetryCount:     0,
					MaxRetries:     3,
					CreatedAt:      time.Now().Unix(),
					UpdatedAt:      time.Now().Unix(),
				})
			}
		}

		// Enqueue items
		if len(queueItems) > 0 {
			err = h.queue.Enqueue(c.Context(), queueItems)
			if err != nil {
				c.Status(common.StatusInternalServerError).JSON(fiber.Map{
					"code":    common.ErrCodeBusinessOperation.Code,
					"message": fmt.Sprintf("Không thể thêm items vào queue: %v", err),
					"status":  "error",
				})
				return nil
			}
		}

		c.JSON(map[string]interface{}{
			"message":   "Notification đã được thêm vào queue",
			"eventType": req.EventType,
			"queued":    len(queueItems),
		})
		return nil
	})
}

// SafeHandlerWrapper wrapper để xử lý errors
func SafeHandlerWrapper(c fiber.Ctx, fn func() error) error {
	if err := fn(); err != nil {
		return err
	}
	return nil
}

