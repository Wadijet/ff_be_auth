package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationQueueItem - Queue item để xử lý
type NotificationQueueItem struct {
	ID                 primitive.ObjectID            `json:"id,omitempty" bson:"_id,omitempty"`
	EventType          string                        `json:"eventType" bson:"eventType" index:"single:1"`
	OwnerOrganizationID primitive.ObjectID            `json:"ownerOrganizationId" bson:"ownerOrganizationId" index:"single:1"` // Tổ chức sở hữu dữ liệu (phân quyền)
	ChannelID      primitive.ObjectID            `json:"channelId" bson:"channelId" index:"single:1"`
	Recipient      string                        `json:"recipient" bson:"recipient"` // Email, chatId, webhook URL
	Payload        map[string]interface{}        `json:"payload" bson:"payload"`

	Status      string `json:"status" bson:"status" index:"single:1"` // pending, processing, completed, failed
	RetryCount  int    `json:"retryCount" bson:"retryCount"`
	MaxRetries  int    `json:"maxRetries" bson:"maxRetries"` // Mặc định: 3
	NextRetryAt *int64 `json:"nextRetryAt,omitempty" bson:"nextRetryAt,omitempty" index:"single:1"`

	Error     string `json:"error,omitempty" bson:"error,omitempty"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64  `json:"updatedAt" bson:"updatedAt"`
}

