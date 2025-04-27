package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"

	"go.mongodb.org/mongo-driver/bson"
)

// SyncFbPagesJob là job đồng bộ danh sách Facebook pages từ Pages.fm
type SyncFbPagesJob struct {
	accessTokenService *services.AccessTokenService
	fbPageService      *services.FbPageService
}

// NewSyncFbPagesJob tạo mới một instance của SyncFbPagesJob
func NewSyncFbPagesJob() (*SyncFbPagesJob, error) {
	// Khởi tạo các services cần thiết
	accessTokenService, err := services.NewAccessTokenService()
	if err != nil {
		return nil, fmt.Errorf("failed to create access token service: %v", err)
	}

	fbPageService, err := services.NewFbPageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create fb page service: %v", err)
	}

	jobLogger.Info("SyncFbPagesJob created successfully")
	return &SyncFbPagesJob{
		accessTokenService: accessTokenService,
		fbPageService:      fbPageService,
	}, nil
}

// Name trả về tên của job
func (j *SyncFbPagesJob) Name() string {
	return "sync_fb_pages"
}

// Run thực thi job đồng bộ Facebook pages
func (j *SyncFbPagesJob) Run() error {
	startTime := time.Now()
	jobLogger.WithField("job", "sync_fb_pages").Info("Starting job execution")

	ctx := context.Background()

	// Lấy access token từ hệ thống với system là "pancake"
	filter := bson.M{"system": "pancake"}
	accessToken, err := j.accessTokenService.FindOne(ctx, filter, nil)
	if err != nil {
		jobLogger.WithError(err).Error("Failed to get pancake access token")
		return fmt.Errorf("failed to get pancake access token: %v", err)
	}
	jobLogger.WithField("token_id", accessToken.ID.Hex()).Info("Found pancake access token")

	// Tạo request đến API Pages.fm
	url := fmt.Sprintf("https://pages.fm/api/v1/pages?access_token=%s", accessToken.Value)
	jobLogger.WithField("url", url).Debug("Calling Pages.fm API")

	resp, err := http.Get(url)
	if err != nil {
		jobLogger.WithError(err).Error("Failed to call Pages.fm API")
		return fmt.Errorf("failed to call pages.fm API: %v", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		jobLogger.WithError(err).Error("Failed to read response body")
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Kiểm tra status code
	if resp.StatusCode != http.StatusOK {
		jobLogger.WithFields(map[string]interface{}{
			"status_code": resp.StatusCode,
			"body":        string(body),
		}).Error("Pages.fm API returned error status")
		return fmt.Errorf("pages.fm API returned error status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response JSON
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		jobLogger.WithError(err).Error("Failed to parse response JSON")
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// Kiểm tra response success
	success, ok := response["success"].(bool)
	if !ok || !success {
		jobLogger.Error("Pages.fm API returned invalid or false success status")
		return fmt.Errorf("pages.fm API returned invalid or false success status")
	}

	// Lấy danh sách pages từ categorized.activated
	categorized, ok := response["categorized"].(map[string]interface{})
	if !ok {
		jobLogger.Error("Invalid response format: missing categorized field")
		return fmt.Errorf("invalid response format: missing categorized field")
	}

	activated, ok := categorized["activated"].([]interface{})
	if !ok {
		jobLogger.Error("Invalid response format: missing or invalid activated field")
		return fmt.Errorf("invalid response format: missing or invalid activated field")
	}

	jobLogger.WithField("page_count", len(activated)).Info("Found Facebook pages to sync")

	// Cập nhật hoặc tạo mới các Facebook pages
	successCount := 0
	for i, page := range activated {
		pageData := page.(map[string]interface{})
		jobLogger.WithFields(map[string]interface{}{
			"index":     i + 1,
			"page_id":   pageData["id"],
			"page_name": pageData["name"],
		}).Debug("Processing page")

		input := &models.FbPageCreateInput{
			AccessToken: accessToken.Value,
			PanCakeData: pageData,
		}

		_, err := j.fbPageService.ReviceData(ctx, input)
		if err != nil {
			jobLogger.WithError(err).WithFields(map[string]interface{}{
				"page_id":   pageData["id"],
				"page_name": pageData["name"],
			}).Error("Failed to update/create fb page")
			return fmt.Errorf("failed to update/create fb page: %v", err)
		}
		successCount++
	}

	duration := time.Since(startTime)
	jobLogger.WithFields(map[string]interface{}{
		"success_count": successCount,
		"total_count":   len(activated),
		"duration_ms":   duration.Milliseconds(),
	}).Info("Job execution completed")

	return nil
}

// Schedule trả về lịch chạy của job theo định dạng cron
func (j *SyncFbPagesJob) Schedule() string {
	return "0 */1 * * * *" // Chạy 1 phút một lần, bắt đầu ở giây thứ 0
}
