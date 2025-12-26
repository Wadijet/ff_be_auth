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

// SendTelegram gửi telegram message
func SendTelegram(ctx context.Context, sender *models.NotificationChannelSender, channel *models.NotificationChannel, chatID string, template *RenderedTemplate, historyID string, baseURL string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", sender.BotToken)

	// Format CTAs thành inline keyboard
	inlineKeyboard := [][]map[string]interface{}{}
	row := []map[string]interface{}{}
	for _, cta := range template.CTAs {
		button := map[string]interface{}{
			"text": cta.Label,
			"url":  cta.Action, // Đã có tracking URL
		}
		row = append(row, button)
		if len(row) >= 3 { // Tối đa 3 buttons/row
			inlineKeyboard = append(inlineKeyboard, row)
			row = []map[string]interface{}{}
		}
	}
	if len(row) > 0 {
		inlineKeyboard = append(inlineKeyboard, row)
	}

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    template.Content,
	}

	if len(inlineKeyboard) > 0 {
		keyboard := map[string]interface{}{
			"inline_keyboard": inlineKeyboard,
		}
		payload["reply_markup"] = keyboard
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}

