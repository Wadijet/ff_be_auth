package notification

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"time"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/notification/channels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Processor xử lý queue items
type Processor struct {
	queueService    *services.NotificationQueueService
	historyService  *services.NotificationHistoryService
	channelService  *services.NotificationChannelService
	senderService   *services.NotificationSenderService
	templateService *Template
	orgService      *services.OrganizationService
	baseURL         string
}

// NewProcessor tạo mới Processor
func NewProcessor(baseURL string) (*Processor, error) {
	queueService, err := services.NewNotificationQueueService()
	if err != nil {
		return nil, fmt.Errorf("failed to create queue service: %w", err)
	}

	historyService, err := services.NewNotificationHistoryService()
	if err != nil {
		return nil, fmt.Errorf("failed to create history service: %w", err)
	}

	channelService, err := services.NewNotificationChannelService()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel service: %w", err)
	}

	senderService, err := services.NewNotificationSenderService()
	if err != nil {
		return nil, fmt.Errorf("failed to create sender service: %w", err)
	}

	templateService, err := NewTemplate()
	if err != nil {
		return nil, fmt.Errorf("failed to create template service: %w", err)
	}

	orgService, err := services.NewOrganizationService()
	if err != nil {
		return nil, fmt.Errorf("failed to create organization service: %w", err)
	}

	return &Processor{
		queueService:    queueService,
		historyService:  historyService,
		channelService:  channelService,
		senderService:   senderService,
		templateService: templateService,
		orgService:      orgService,
		baseURL:         baseURL,
	}, nil
}

// ProcessQueueItem xử lý một queue item
func (p *Processor) ProcessQueueItem(ctx context.Context, item *models.NotificationQueueItem) error {
	// 1. Lấy channel
	channel, err := p.channelService.FindOneById(ctx, item.ChannelID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return fmt.Errorf("channel not found: %w", err)
		}
		return fmt.Errorf("failed to find channel: %w", err)
	}

	// 2. Tìm sender (với fallback)
	sender, err := p.findSender(ctx, &channel, item.OwnerOrganizationID)
	if err != nil {
		return fmt.Errorf("sender not found: %w", err)
	}

	// 3. Tìm template
	template, err := p.templateService.FindTemplate(ctx, item.EventType, channel.ChannelType, item.OwnerOrganizationID)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	// 4. Render template
	rendered, err := p.templateService.Render(template, item.Payload)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// 5. Tạo tracking token cho history
	token, err := generateTrackingToken()
	if err != nil {
		return fmt.Errorf("failed to generate tracking token: %w", err)
	}

	// 6. Tạo history record (trước khi gửi)
	historyID := primitive.NewObjectID()
	history := &models.NotificationHistory{
		ID:                 historyID,
		QueueItemID:        item.ID,
		EventType:          item.EventType,
		OwnerOrganizationID: item.OwnerOrganizationID, // Phân quyền dữ liệu
		ChannelID:          item.ChannelID,
		ChannelType:    channel.ChannelType,
		Recipient:      item.Recipient,
		Status:         "pending",
		Content:        rendered.Content,
		RetryCount:     item.RetryCount,
		CreatedAt:      time.Now().Unix(),
	}

	// Initialize CTAClicks array
	history.CTAClicks = make([]models.CTAClick, len(rendered.CTAs))
	for i, cta := range rendered.CTAs {
		history.CTAClicks[i] = models.CTAClick{
			CTAIndex:   i,
			Label:      cta.Label,
			ClickCount: 0,
		}
	}

	// 7. Thay thế CTA URLs bằng tracking URLs
	for i := range rendered.CTAs {
		originalURL := rendered.CTAs[i].Action
		trackingURL := fmt.Sprintf("%s/api/v1/notification/track/%s/%d?token=%s&url=%s",
			p.baseURL, historyID.Hex(), i, token, base64.URLEncoding.EncodeToString([]byte(originalURL)))
		rendered.CTAs[i].Action = trackingURL
	}

	// 8. Gửi notification (với fallback sender)
	err = p.sendNotificationWithFallback(ctx, sender, &channel, item.Recipient, rendered, historyID.Hex())
	if err != nil {
		// Update history với error
		history.Status = "failed"
		history.Error = err.Error()
		history.SentAt = nil
	} else {
		// Update history với success
		history.Status = "sent"
		now := time.Now().Unix()
		history.SentAt = &now
	}

	// 9. Lưu history
	_, err = p.historyService.InsertOne(ctx, *history)
	if err != nil {
		return fmt.Errorf("failed to save history: %w", err)
	}

	// 10. Update queue item status
	sendErr := err
	if sendErr != nil {
		// Retry logic
		if item.RetryCount < item.MaxRetries {
			item.RetryCount++
			item.Status = "pending"
			// Exponential backoff: 2^retryCount seconds
			backoffSeconds := int64(math.Pow(2, float64(item.RetryCount)))
			nextRetryAt := time.Now().Unix() + backoffSeconds
			item.NextRetryAt = &nextRetryAt
			item.UpdatedAt = time.Now().Unix()
			updateData := services.UpdateData{
				Set: map[string]interface{}{
					"status":      item.Status,
					"retryCount":  item.RetryCount,
					"nextRetryAt": item.NextRetryAt,
					"updatedAt":   item.UpdatedAt,
				},
			}
			_, err = p.queueService.UpdateOne(ctx, bson.M{"_id": item.ID}, updateData, nil)
		} else {
			item.Status = "failed"
			item.UpdatedAt = time.Now().Unix()
			updateData := services.UpdateData{
				Set: map[string]interface{}{
					"status":    item.Status,
					"updatedAt": item.UpdatedAt,
				},
			}
			_, err = p.queueService.UpdateOne(ctx, bson.M{"_id": item.ID}, updateData, nil)
		}
	} else {
		item.Status = "completed"
		item.UpdatedAt = time.Now().Unix()
		updateData := services.UpdateData{
			Set: map[string]interface{}{
				"status":    item.Status,
				"updatedAt": item.UpdatedAt,
			},
		}
		_, err = p.queueService.UpdateOne(ctx, bson.M{"_id": item.ID}, updateData, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to update queue item: %w", err)
	}
	if sendErr != nil {
		return sendErr
	}
	return nil
}

// findSender tìm sender với fallback logic
func (p *Processor) findSender(ctx context.Context, channel *models.NotificationChannel, orgID primitive.ObjectID) (*models.NotificationChannelSender, error) {
	// 1. Thử các senders trong channel.SenderIDs (nếu có)
	if len(channel.SenderIDs) > 0 {
		for _, senderID := range channel.SenderIDs {
			sender, err := p.senderService.FindOneById(ctx, senderID)
			if err == nil && sender.IsActive && sender.ChannelType == channel.ChannelType {
				return &sender, nil
			}
		}
	}

	// 2. Fallback: Tìm sender theo inheritance (team's company -> group -> global)
	org, err := p.orgService.FindOneById(ctx, orgID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to find organization: %w", err)
	}

	// Tìm company của team
	var companyID *primitive.ObjectID
	if org.ParentID != nil {
		parent, err := p.orgService.FindOneById(ctx, *org.ParentID)
		if err == nil {
			// Tìm company (có thể là parent hoặc parent của parent)
			current := parent
			for current.Type != models.OrganizationTypeCompany {
				if current.ParentID != nil {
					parent, err := p.orgService.FindOneById(ctx, *current.ParentID)
					if err != nil {
						break
					}
					current = parent
				} else {
					break
				}
			}
			if current.Type == models.OrganizationTypeCompany {
				companyID = &current.ID
			}
		}
	}

	// Tìm group (parent của company)
	var groupID *primitive.ObjectID
	if companyID != nil {
		company, err := p.orgService.FindOneById(ctx, *companyID)
		if err == nil && company.ParentID != nil {
			groupID = company.ParentID
		}
	}

	// Thử tìm sender theo thứ tự: company -> group -> global
	searchOrder := []*primitive.ObjectID{companyID, groupID, nil}
	for _, searchOrgID := range searchOrder {
		sender, err := p.findSenderByOrganization(ctx, channel.ChannelType, searchOrgID)
		if err == nil && sender != nil {
			return sender, nil
		}
	}

	return nil, fmt.Errorf("no active sender found for channel type %s", channel.ChannelType)
}

// findSenderByOrganization tìm sender theo organization (null = global)
func (p *Processor) findSenderByOrganization(ctx context.Context, channelType string, orgID *primitive.ObjectID) (*models.NotificationChannelSender, error) {
	filter := bson.M{
		"channelType": channelType,
		"isActive":    true,
	}
	if orgID != nil {
		filter["ownerOrganizationId"] = *orgID
	} else {
		filter["ownerOrganizationId"] = nil
	}

	sender, err := p.senderService.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, nil // Not found, but not an error
		}
		return nil, err
	}
	return &sender, nil
}

// sendNotificationWithFallback gửi notification với fallback sender
func (p *Processor) sendNotificationWithFallback(ctx context.Context, sender *models.NotificationChannelSender, channel *models.NotificationChannel, recipient string, rendered *channels.RenderedTemplate, historyID string) error {
	// Thử gửi với sender hiện tại
	err := p.sendNotification(ctx, sender, channel, recipient, rendered, historyID)
	if err == nil {
		return nil
	}

	// Nếu fail và có fallback senders, thử các sender khác
	if len(channel.SenderIDs) > 1 {
		for _, senderID := range channel.SenderIDs[1:] {
			fallbackSender, err := p.senderService.FindOneById(ctx, senderID)
			if err == nil && fallbackSender.IsActive && fallbackSender.ChannelType == channel.ChannelType {
				err = p.sendNotification(ctx, &fallbackSender, channel, recipient, rendered, historyID)
				if err == nil {
					return nil
				}
			}
		}
	}

	return err
}

// sendNotification gửi notification qua channel tương ứng
func (p *Processor) sendNotification(ctx context.Context, sender *models.NotificationChannelSender, channel *models.NotificationChannel, recipient string, rendered *channels.RenderedTemplate, historyID string) error {
	switch channel.ChannelType {
	case "email":
		return channels.SendEmail(ctx, sender, channel, recipient, rendered, historyID, p.baseURL)
	case "telegram":
		return channels.SendTelegram(ctx, sender, channel, recipient, rendered, historyID, p.baseURL)
	case "webhook":
		return channels.SendWebhook(ctx, channel, rendered, historyID, p.baseURL)
	default:
		return fmt.Errorf("unsupported channel type: %s", channel.ChannelType)
	}
}

// generateTrackingToken tạo token để tracking
func generateTrackingToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Start bắt đầu background worker để xử lý queue
func (p *Processor) Start(ctx context.Context) {
	// Polling interval: 5 giây
	interval := 5 * time.Second
	// Batch size: xử lý tối đa 10 items mỗi lần
	batchSize := 10

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Lấy items từ queue
			items, err := p.queueService.FindPending(ctx, batchSize)
			if err != nil {
				// Log error nhưng tiếp tục
				continue
			}

			if len(items) == 0 {
				continue
			}

			// Xử lý từng item
			for _, item := range items {
				// Update status to processing
				ids := []interface{}{item.ID}
				err = p.queueService.UpdateStatus(ctx, ids, "processing")
				if err != nil {
					continue
				}

				// Process item
				err = p.ProcessQueueItem(ctx, &item)
				if err != nil {
					// Log error nhưng tiếp tục với item khác
					continue
				}
			}
		}
	}
}

