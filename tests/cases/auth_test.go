package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Cấu trúc để lưu token JWT và HWID
var (
	authToken string
	deviceID  = "test_device_123" // ID thiết bị cố định cho test
)

// TestAuthFlow kiểm tra toàn bộ luồng xác thực
func TestAuthFlow(t *testing.T) {
	// Đợi server khởi động
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:8080/api/v1"

	// Test case 1: Đăng ký tài khoản mới
	t.Run("👤 Đăng ký tài khoản", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "Test@123",
		}

		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API đăng ký: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code phải là 200")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err, "Phải parse được JSON response")

		// Kiểm tra response
		assert.NotNil(t, result["data"], "Phải có thông tin user trong response")
		fmt.Printf("✅ Đăng ký thành công với email: %v\n", payload["email"])
	})

	// Test case 2: Đăng nhập
	t.Run("🔐 Đăng nhập", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":    "test@example.com",
			"password": "Test@123",
			"hwid":     deviceID,
		}

		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API đăng nhập: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code phải là 200")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err, "Phải parse được JSON response")

		// Lưu token để dùng cho các test case sau
		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok, "Phải có data trong response")
		token, ok := data["token"].(string)
		assert.True(t, ok, "Phải có token trong response")
		authToken = token

		fmt.Printf("✅ Đăng nhập thành công và nhận được token\n")
	})

	// Test case 3: Lấy thông tin profile
	t.Run("👤 Lấy thông tin profile", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/auth/profile", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API lấy profile: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code phải là 200")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err, "Phải parse được JSON response")

		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok, "Phải có data trong response")

		// Kiểm tra thông tin profile
		assert.Equal(t, "Test User", data["name"], "Name phải khớp")
		assert.Equal(t, "test@example.com", data["email"], "Email phải khớp")

		fmt.Printf("✅ Lấy thông tin profile thành công\n")
	})

	// Test case 4: Cập nhật profile
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
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code phải là 200")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err, "Phải parse được JSON response")

		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok, "Phải có data trong response")

		// Kiểm tra thông tin đã cập nhật
		assert.Equal(t, "Updated Test User", data["name"], "Tên phải được cập nhật")

		fmt.Printf("✅ Cập nhật profile thành công\n")
	})

	// Test case 5: Đổi mật khẩu
	t.Run("🔑 Đổi mật khẩu", func(t *testing.T) {
		payload := map[string]interface{}{
			"oldPassword": "Test@123",
			"newPassword": "NewTest@123",
		}

		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", baseURL+"/auth/password", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API đổi mật khẩu: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code phải là 200")
		fmt.Printf("✅ Đổi mật khẩu thành công\n")
	})

	// Test case 6: Đăng xuất
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
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code phải là 200")
		fmt.Printf("✅ Đăng xuất thành công\n")
	})
}
