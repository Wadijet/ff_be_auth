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

// TestErrorHandling kiểm tra xử lý lỗi và edge cases
func TestErrorHandling(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	client := utils.NewHTTPClient(baseURL, 10)

	// Test 1: Đăng nhập Firebase với token không hợp lệ
	t.Run("❌ Đăng nhập Firebase với token không hợp lệ", func(t *testing.T) {
		payload := map[string]interface{}{
			"idToken": "invalid_firebase_token",
			"hwid":    "test_device_123",
		}

		resp, body, err := client.POST("/auth/login/firebase", payload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi đăng nhập Firebase: %v", err)
		}

		// Phải trả về 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Firebase token không hợp lệ phải trả về 401")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Phải parse được JSON response")
		fmt.Printf("✅ Test Firebase token không hợp lệ thành công (status: %d)\n", resp.StatusCode)
	})

	// Test 2: Đăng nhập Firebase với dữ liệu thiếu
	t.Run("❌ Đăng nhập Firebase với dữ liệu thiếu", func(t *testing.T) {
		// Thiếu idToken
		payload1 := map[string]interface{}{
			"hwid": "test_device_123",
		}

		resp1, body1, err := client.POST("/auth/login/firebase", payload1)
		if err != nil {
			t.Fatalf("❌ Lỗi khi đăng nhập Firebase: %v", err)
		}

		// Phải trả về 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp1.StatusCode, "Thiếu idToken phải trả về 400")
		fmt.Printf("✅ Test thiếu idToken thành công (status: %d - %s)\n", resp1.StatusCode, string(body1))

		// Thiếu hwid
		payload2 := map[string]interface{}{
			"idToken": "test_token",
		}

		resp2, body2, err := client.POST("/auth/login/firebase", payload2)
		if err != nil {
			t.Fatalf("❌ Lỗi khi đăng nhập Firebase: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, resp2.StatusCode, "Thiếu hwid phải trả về 400")
		fmt.Printf("✅ Test thiếu hwid thành công (status: %d - %s)\n", resp2.StatusCode, string(body2))
	})

	// Test 3: Truy cập API cần auth mà không có token
	t.Run("❌ Truy cập API cần auth không có token", func(t *testing.T) {
		// Không set token
		clientNoToken := utils.NewHTTPClient(baseURL, 10)

		resp, body, err := clientNoToken.GET("/auth/profile")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}

		// Phải trả về 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Truy cập không token phải trả về 401")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Phải parse được JSON response")
		fmt.Printf("✅ Test không có token thành công (status: %d)\n", resp.StatusCode)
	})

	// Test 4: Truy cập API với token không hợp lệ
	t.Run("❌ Truy cập API với token không hợp lệ", func(t *testing.T) {
		clientInvalidToken := utils.NewHTTPClient(baseURL, 10)
		clientInvalidToken.SetToken("invalid_token_12345")

		resp, body, err := clientInvalidToken.GET("/auth/profile")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}

		// Phải trả về 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Token không hợp lệ phải trả về 401")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Phải parse được JSON response")
		fmt.Printf("✅ Test token không hợp lệ thành công (status: %d)\n", resp.StatusCode)
	})

	// Test 5: Truy cập API không tồn tại
	t.Run("❌ Truy cập API không tồn tại", func(t *testing.T) {
		resp, _, err := client.GET("/api/v1/nonexistent/endpoint")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}

		// Phải trả về 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "API không tồn tại phải trả về 404")
		fmt.Printf("✅ Test API không tồn tại thành công (status: %d)\n", resp.StatusCode)
	})

}
