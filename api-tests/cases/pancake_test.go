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

// TestPancakeAPIs kiểm tra các API Pancake integration
func TestPancakeAPIs(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

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

	// Test Pancake Order APIs
	t.Run("🥞 Pancake Order APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách orders
		t.Run("Lấy danh sách orders", func(t *testing.T) {
			resp, body, err := client.GET("/pancake/order/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách orders: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách orders thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách orders yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Count orders
		t.Run("Đếm số lượng orders", func(t *testing.T) {
			resp, body, err := client.GET("/pancake/order/count")
			if err != nil {
				t.Fatalf("❌ Lỗi khi đếm orders: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Đếm orders thành công\n")
			} else {
				fmt.Printf("⚠️ Đếm orders yêu cầu quyền (status: %d)\n", resp.StatusCode)
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
