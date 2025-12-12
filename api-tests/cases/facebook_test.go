package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"ff_be_auth_tests/utils"

	"github.com/stretchr/testify/assert"
)

// TestFacebookAPIs kiểm tra các API Facebook integration
func TestFacebookAPIs(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Khởi tạo dữ liệu mặc định trước
	initTestData(t, baseURL)

	fixtures := utils.NewTestFixtures(baseURL)

	// Tạo user với token
	firebaseIDToken := utils.GetTestFirebaseIDToken()
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
	}
	_, _, token, err := fixtures.CreateTestUser(firebaseIDToken)
	if err != nil {
		t.Fatalf("❌ Không thể tạo user test: %v", err)
	}

	client := utils.NewHTTPClient(baseURL, 10)
	client.SetToken(token)

	// Test AccessToken APIs
	t.Run("🔑 AccessToken APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách access tokens
		t.Run("Lấy danh sách access tokens", func(t *testing.T) {
			resp, body, err := client.GET("/access-token/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách access tokens: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách access tokens thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách access tokens yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test Facebook Page APIs
	t.Run("📄 Facebook Page APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách pages
		t.Run("Lấy danh sách pages", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/page/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách pages: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách pages thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách pages yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test Facebook Post APIs
	t.Run("📝 Facebook Post APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách posts
		t.Run("Lấy danh sách posts", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/post/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách posts: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách posts thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách posts yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test Facebook Conversation APIs
	t.Run("💬 Facebook Conversation APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách conversations
		t.Run("Lấy danh sách conversations", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/conversation/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách conversations: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách conversations thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách conversations yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Lấy conversations sắp xếp theo API update
		t.Run("Lấy conversations sắp xếp theo API update", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/conversation/sort-by-api-update")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy conversations sort: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy conversations sort thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy conversations sort yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test Facebook Message APIs
	t.Run("📨 Facebook Message APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách messages
		t.Run("Lấy danh sách messages", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/message/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách messages: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách messages thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách messages yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Cleanup
	t.Run("🧹 Cleanup", func(t *testing.T) {
		logoutPayload := map[string]interface{}{
			"hwid": "test_device_123",
		}
		client.POST("/auth/logout", logoutPayload)
		fmt.Printf("✅ Cleanup hoàn tất\n")
	})
}
