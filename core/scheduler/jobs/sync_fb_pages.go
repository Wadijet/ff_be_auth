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
	scheduler "meta_commerce/core/scheduler"

	"go.mongodb.org/mongo-driver/bson"
)

// SyncFbPagesJob là job đồng bộ danh sách Facebook pages từ Pages.fm
// Kế thừa BaseJob để dùng chung name, schedule và các hàm mặc định
type SyncFbPagesJob struct {
	*scheduler.BaseJob
	accessTokenService *services.AccessTokenService
	fbPageService      *services.FbPageService
}

// NewSyncFbPagesJob tạo mới một instance của SyncFbPagesJob
func NewSyncFbPagesJob(schedule string) (*SyncFbPagesJob, error) {
	// Khởi tạo các services cần thiết
	accessTokenService, err := services.NewAccessTokenService()
	if err != nil {
		return nil, fmt.Errorf("failed to create access token service: %v", err)
	}

	fbPageService, err := services.NewFbPageService()
	if err != nil {
		return nil, fmt.Errorf("failed to create fb page service: %v", err)
	}

	GetJobLogger().Info("SyncFbPagesJob created successfully")
	base := scheduler.NewBaseJob("sync_fb_pages", schedule)
	return &SyncFbPagesJob{
		BaseJob:            base,
		accessTokenService: accessTokenService,
		fbPageService:      fbPageService,
	}, nil
}

// Execute thực thi job đồng bộ Facebook pages (chuẩn interface Job)
func (j *SyncFbPagesJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	GetJobLogger().WithField("job", "sync_fb_pages").Info("Bắt đầu thực thi job")

	// Lấy access token từ hệ thống với system là "pancake"
	filter := bson.M{"system": "pancake"}
	accessToken, err := j.accessTokenService.FindOne(ctx, filter, nil)
	if err != nil {
		GetJobLogger().WithError(err).Error("Không lấy được access token pancake")
		return fmt.Errorf("không lấy được access token pancake: %v", err)
	}
	GetJobLogger().WithField("token_id", accessToken.ID.Hex()).Info("Đã tìm thấy access token pancake")

	// Tạo request đến API Pages.fm
	url := fmt.Sprintf("https://pages.fm/api/v1/pages?access_token=%s", accessToken.Value)
	GetJobLogger().WithField("url", url).Debug("Gọi API Pages.fm")

	resp, err := http.Get(url)
	if err != nil {
		GetJobLogger().WithError(err).Error("Gọi API Pages.fm thất bại")
		return fmt.Errorf("gọi API pages.fm thất bại: %v", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		GetJobLogger().WithError(err).Error("Đọc response body thất bại")
		return fmt.Errorf("đọc response body thất bại: %v", err)
	}

	// Kiểm tra status code
	if resp.StatusCode != http.StatusOK {
		GetJobLogger().WithFields(map[string]interface{}{
			"status_code": resp.StatusCode,
			"body":        string(body),
		}).Error("Pages.fm API trả về lỗi")
		return fmt.Errorf("pages.fm API trả về lỗi: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response JSON
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		GetJobLogger().WithError(err).Error("Parse response JSON thất bại")
		return fmt.Errorf("parse response JSON thất bại: %v", err)
	}

	// Kiểm tra response success
	success, ok := response["success"].(bool)
	if !ok || !success {
		GetJobLogger().Error("Pages.fm API trả về trạng thái thành công không hợp lệ hoặc false")
		return fmt.Errorf("pages.fm API trả về trạng thái thành công không hợp lệ hoặc false")
	}

	// Lấy danh sách pages từ categorized.activated
	categorized, ok := response["categorized"].(map[string]interface{})
	if !ok {
		GetJobLogger().Error("Response không hợp lệ: thiếu trường categorized")
		return fmt.Errorf("response không hợp lệ: thiếu trường categorized")
	}

	activated, ok := categorized["activated"].([]interface{})
	if !ok {
		GetJobLogger().Error("Response không hợp lệ: thiếu hoặc sai trường activated")
		return fmt.Errorf("response không hợp lệ: thiếu hoặc sai trường activated")
	}

	GetJobLogger().WithField("page_count", len(activated)).Info("Tìm thấy các Facebook page cần đồng bộ")

	// Cập nhật hoặc tạo mới các Facebook pages
	successCount := 0
	for i, page := range activated {
		pageData := page.(map[string]interface{})
		GetJobLogger().WithFields(map[string]interface{}{
			"index":     i + 1,
			"page_id":   pageData["id"],
			"page_name": pageData["name"],
		}).Debug("Đang xử lý page")

		input := &models.FbPageCreateInput{
			AccessToken: accessToken.Value,
			PanCakeData: pageData,
		}

		_, err := j.fbPageService.ReviceData(ctx, input)
		if err != nil {
			GetJobLogger().WithError(err).WithFields(map[string]interface{}{
				"page_id":   pageData["id"],
				"page_name": pageData["name"],
			}).Error("Cập nhật/tạo mới fb page thất bại")
			return fmt.Errorf("cập nhật/tạo mới fb page thất bại: %v", err)
		}
		successCount++
	}

	duration := time.Since(startTime)
	GetJobLogger().WithFields(map[string]interface{}{
		"success_count": successCount,
		"total_count":   len(activated),
		"duration_ms":   duration.Milliseconds(),
	}).Info("Thực thi job hoàn tất")

	return nil
}
