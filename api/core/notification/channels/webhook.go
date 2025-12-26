package channels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	models "meta_commerce/core/api/models/mongodb"
)

// SendWebhook gửi webhook
func SendWebhook(ctx context.Context, channel *models.NotificationChannel, template *RenderedTemplate, historyID string, baseURL string) error {
	// Format CTAs thành JSON
	actions := []map[string]interface{}{}
	for _, cta := range template.CTAs {
		actions = append(actions, map[string]interface{}{
			"label": cta.Label,
			"url":   cta.Action, // Đã có tracking URL
			"style": cta.Style,
		})
	}

	payload := map[string]interface{}{
		"content":   template.Content,
		"timestamp": time.Now().Unix(),
		"actions":   actions, // CTAs với tracking URLs
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", channel.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range channel.WebhookHeaders {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

