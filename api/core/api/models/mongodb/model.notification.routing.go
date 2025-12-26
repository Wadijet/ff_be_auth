package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationRoutingRule - Routing rule định nghĩa: Event nào → Gửi cho teams nào → Qua channels nào
type NotificationRoutingRule struct {
	ID              primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	EventType       string               `json:"eventType" bson:"eventType" index:"single:1"`                    // conversation_unreplied
	OrganizationIDs []primitive.ObjectID `json:"organizationIds" bson:"organizationIds"`                         // Teams nào nhận (có thể nhiều)
	ChannelTypes    []string             `json:"channelTypes,omitempty" bson:"channelTypes,omitempty"`          // Filter channels theo type (optional: email, telegram, webhook)
	IsActive        bool                 `json:"isActive" bson:"isActive" index:"single:1"`
	CreatedAt       int64                `json:"createdAt" bson:"createdAt"`
	UpdatedAt       int64                `json:"updatedAt" bson:"updatedAt"`
}

