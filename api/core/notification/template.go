package notification

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/notification/channels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Template xử lý việc tìm và render template
type Template struct {
	templateService *services.NotificationTemplateService
}

// NewTemplate tạo mới Template
func NewTemplate() (*Template, error) {
	templateService, err := services.NewNotificationTemplateService()
	if err != nil {
		return nil, fmt.Errorf("failed to create template service: %w", err)
	}

	return &Template{
		templateService: templateService,
	}, nil
}

// FindTemplate tìm template theo EventType, ChannelType, và OrganizationID
// Logic: Tìm team-specific trước, nếu không có → tìm global
func (t *Template) FindTemplate(ctx context.Context, eventType string, channelType string, organizationID primitive.ObjectID) (*models.NotificationTemplate, error) {
	// 1. Tìm team-specific template
	filter := bson.M{
		"eventType":   eventType,
		"channelType": channelType,
		"organizationId": organizationID,
		"isActive":    true,
	}

	template, err := t.templateService.FindOne(ctx, filter, nil)
	if err == nil {
		return &template, nil
	}
	if !errors.Is(err, common.ErrNotFound) {
		return nil, fmt.Errorf("failed to find team-specific template: %w", err)
	}

	// 2. Nếu không có → Tìm global template (organizationId = null)
	filter = bson.M{
		"eventType":      eventType,
		"channelType":    channelType,
		"organizationId": nil,
		"isActive":       true,
	}

	template, err = t.templateService.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, fmt.Errorf("template not found for eventType=%s, channelType=%s", eventType, channelType)
		}
		return nil, fmt.Errorf("failed to find global template: %w", err)
	}

	return &template, nil
}

// RenderedTemplate và RenderedCTA đã được di chuyển vào channels package để tránh import cycle

// Render render template với payload
func (t *Template) Render(template *models.NotificationTemplate, payload map[string]interface{}) (*channels.RenderedTemplate, error) {
	// Render subject
	subject := template.Subject
	for _, variable := range template.Variables {
		value, exists := payload[variable]
		if !exists {
			value = ""
		}
		placeholder := "{{" + variable + "}}"
		subject = strings.ReplaceAll(subject, placeholder, fmt.Sprintf("%v", value))
	}

	// Render content
	content := template.Content
	for _, variable := range template.Variables {
		value, exists := payload[variable]
		if !exists {
			value = ""
		}
		placeholder := "{{" + variable + "}}"
		content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", value))
	}

	// Render CTAs (nếu có)
	renderedCTAs := []channels.RenderedCTA{}
	for _, cta := range template.CTAs {
		renderedCTA := channels.RenderedCTA{
			Label:  cta.Label,
			Action: cta.Action,
			Style:  cta.Style,
		}

		// Render variables trong Action
		for _, variable := range template.Variables {
			value, exists := payload[variable]
			if !exists {
				value = ""
			}
			placeholder := "{{" + variable + "}}"
			renderedCTA.Action = strings.ReplaceAll(renderedCTA.Action, placeholder, fmt.Sprintf("%v", value))
		}

		// Lưu original URL (trước khi thay bằng tracking URL)
		renderedCTA.OriginalURL = renderedCTA.Action

		renderedCTAs = append(renderedCTAs, renderedCTA)
	}

	return &channels.RenderedTemplate{
		Subject: subject,
		Content: content,
		CTAs:    renderedCTAs,
	}, nil
}

