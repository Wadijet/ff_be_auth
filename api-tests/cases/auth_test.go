package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Cấu trúc để lưu token JWT và HWID
var (
	authToken string
	deviceID  = "test_device_123" // ID thiết bị cố định cho test
)

// waitForHealth đợi server sẵn sàng bằng cách ping endpoint health
func waitForHealth(baseURL string, attempts int, delay time.Duration, t *testing.T) {
	for i := 0; i < attempts; i++ {
		resp, err := http.Get(baseURL + "/system/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(delay)
	}
	t.Fatalf("❌ Server chưa sẵn sàng sau %d lần thử", attempts)
}

// readBody đọc body và đóng response
func readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// getTestFirebaseIDToken lấy Firebase ID token từ environment variable
// Lưu ý: Test cần có Firebase ID token hợp lệ từ Firebase test project
// Có thể set qua environment variable: TEST_FIREBASE_ID_TOKEN
func getTestFirebaseIDToken(t *testing.T) string {
	// TODO: Lấy từ environment variable hoặc tạo bằng Firebase Admin SDK
	// Tạm thời skip test nếu không có token
	token := os.Getenv("TEST_FIREBASE_ID_TOKEN")
	if token == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
	}
	return token
}

// TestAuthFlow kiểm tra toàn bộ luồng xác thực với Firebase
func TestAuthFlow(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Lấy Firebase ID token từ environment variable
	firebaseIDToken := getTestFirebaseIDToken(t)

	// Test case 1: Đăng nhập bằng Firebase
	t.Run("🔐 Đăng nhập bằng Firebase", func(t *testing.T) {
		payload := map[string]interface{}{
			"idToken": firebaseIDToken,
			"hwid":    deviceID,
		}

		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/auth/login/firebase", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API đăng nhập Firebase: %v", err)
		}
		respBody, readErr := readBody(resp)
		assert.NoError(t, readErr, "Phải đọc được response body")
		assert.Equalf(t, http.StatusOK, resp.StatusCode, "Status code phải là 200. Body: %s", string(respBody))

		var result map[string]interface{}
		err = json.Unmarshal(respBody, &result)
		assert.NoError(t, err, "Phải parse được JSON response")

		// Lưu token để dùng cho các test case sau
		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok, "Phải có data trong response")
		token, ok := data["token"].(string)
		assert.True(t, ok, "Phải có token trong response")
		authToken = token

		fmt.Printf("✅ Đăng nhập Firebase thành công và nhận được token\n")
	})

	// Test case 2: Lấy thông tin profile
	t.Run("👤 Lấy thông tin profile", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/auth/profile", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API lấy profile: %v", err)
		}
		respBody, readErr := readBody(resp)
		assert.NoError(t, readErr, "Phải đọc được response body")
		assert.Equalf(t, http.StatusOK, resp.StatusCode, "Status code phải là 200. Body: %s", string(respBody))

		var result map[string]interface{}
		err = json.Unmarshal(respBody, &result)
		assert.NoError(t, err, "Phải parse được JSON response")

		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok, "Phải có data trong response")

		// Kiểm tra thông tin profile
		assert.NotNil(t, data["name"], "Phải có name trong profile")
		assert.NotNil(t, data["email"], "Phải có email trong profile")

		fmt.Printf("✅ Lấy thông tin profile thành công\n")
	})

	// Test case 3: Cập nhật profile
	t.Run("✏️ Cập nhật profile", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "Updated Test User",
		}

		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", baseURL+"/auth/profile", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API cập nhật profile: %v", err)
		}
		respBody, readErr := readBody(resp)
		assert.NoError(t, readErr, "Phải đọc được response body")
		assert.Equalf(t, http.StatusOK, resp.StatusCode, "Status code phải là 200. Body: %s", string(respBody))

		var result map[string]interface{}
		err = json.Unmarshal(respBody, &result)
		assert.NoError(t, err, "Phải parse được JSON response")

		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok, "Phải có data trong response")

		// Kiểm tra thông tin đã cập nhật
		assert.Equal(t, "Updated Test User", data["name"], "Tên phải được cập nhật")

		fmt.Printf("✅ Cập nhật profile thành công\n")
	})

	// Test case 4: Đăng xuất
	t.Run("🚪 Đăng xuất", func(t *testing.T) {
		payload := map[string]interface{}{
			"hwid": deviceID,
		}

		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/auth/logout", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API đăng xuất: %v", err)
		}
		respBody, readErr := readBody(resp)
		assert.NoError(t, readErr, "Phải đọc được response body")
		assert.Equalf(t, http.StatusOK, resp.StatusCode, "Status code phải là 200. Body: %s", string(respBody))
		fmt.Printf("✅ Đăng xuất thành công\n")
	})
}
