package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationChannel - Cấu hình kênh nhận (recipients) cho team
type NotificationChannel struct {
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID   `json:"organizationId" bson:"organizationId" index:"single:1"` // Team ID
	ChannelType    string               `json:"channelType" bson:"channelType" index:"single:1"`        // email, telegram, webhook
	Name           string               `json:"name" bson:"name" index:"single:1"`
	IsActive       bool                 `json:"isActive" bson:"isActive" index:"single:1"`

	// Sender configs (dự phòng - thứ tự ưu tiên)
	SenderIDs []primitive.ObjectID `json:"senderIds,omitempty" bson:"senderIds,omitempty"` // Mảng sender IDs (thứ tự ưu tiên), null/empty = dùng inheritance

	// Recipients (kênh nhận)
	// Email recipients
	Recipients []string `json:"recipients,omitempty" bson:"recipients,omitempty"` // Email addresses

	// Telegram recipients
	ChatIDs []string `json:"chatIds,omitempty" bson:"chatIds,omitempty"` // Telegram chat IDs

	// Webhook recipients
	WebhookURL     string            `json:"webhookUrl,omitempty" bson:"webhookUrl,omitempty"`         // Webhook URL (chỉ 1 URL)
	WebhookHeaders map[string]string `json:"webhookHeaders,omitempty" bson:"webhookHeaders,omitempty"` // Webhook headers

	CreatedAt int64 `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt"`
}

