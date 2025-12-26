package notification

import (
	"context"
	"fmt"

	"meta_commerce/core/api/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Route đại diện cho một route từ routing rule
type Route struct {
	OrganizationID primitive.ObjectID
	ChannelID      primitive.ObjectID
}

// Router xử lý việc tìm routing rules và tạo routes
type Router struct {
	routingService *services.NotificationRoutingService
	channelService *services.NotificationChannelService
}

// NewRouter tạo mới Router
func NewRouter() (*Router, error) {
	routingService, err := services.NewNotificationRoutingService()
	if err != nil {
		return nil, fmt.Errorf("failed to create routing service: %w", err)
	}

	channelService, err := services.NewNotificationChannelService()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel service: %w", err)
	}

	return &Router{
		routingService: routingService,
		channelService: channelService,
	}, nil
}

// FindRoutes tìm tất cả routes cho một eventType
func (r *Router) FindRoutes(ctx context.Context, eventType string) ([]Route, error) {
	// Tìm tất cả rules cho eventType
	rules, err := r.routingService.FindByEventType(ctx, eventType)
	if err != nil {
		return nil, err
	}

	routes := []Route{}

	// Với mỗi rule
	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		// Với mỗi team trong rule
		for _, orgID := range rule.OrganizationIDs {
			// Lấy TẤT CẢ channels của team (filter theo ChannelTypes nếu có)
			channels, err := r.channelService.FindByOrganizationID(ctx, orgID, rule.ChannelTypes)
			if err != nil {
				// Log error nhưng tiếp tục với team khác
				continue
			}

			// Với mỗi channel của team
			for _, channel := range channels {
				if !channel.IsActive {
					continue
				}
				routes = append(routes, Route{
					OrganizationID: orgID,
					ChannelID:      channel.ID,
				})
			}
		}
	}

	return routes, nil
}

