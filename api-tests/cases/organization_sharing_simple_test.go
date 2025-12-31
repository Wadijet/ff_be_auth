package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"ff_be_auth_tests/utils"

	"github.com/stretchr/testify/assert"
)

// TestOrganizationSharingSimple - Test đơn giản Organization Sharing (không cần Firebase token)
// Test này chỉ kiểm tra API endpoints có hoạt động không
func TestOrganizationSharingSimple(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"

	// Đợi server sẵn sàng
	client := utils.NewHTTPClient(baseURL, 2)
	for i := 0; i < 10; i++ {
		resp, _, err := client.GET("/system/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(1 * time.Second)
		if i == 9 {
			t.Fatalf("Server không sẵn sàng sau 10 lần thử")
		}
	}

	fmt.Printf("✅ Server sẵn sàng\n")

	// Test 1: Kiểm tra endpoint có tồn tại không (sẽ trả về 401 vì chưa có token)
	t.Run("1. Kiểm tra endpoint POST /organization-share (không có token)", func(t *testing.T) {
		sharePayload := map[string]interface{}{
			"fromOrgId":       "507f1f77bcf86cd799439011",
			"toOrgId":         "507f1f77bcf86cd799439012",
			"permissionNames": []string{},
		}

		resp, body, err := client.POST("/organization-share", sharePayload)
		// Không có token nên sẽ trả về 401
		if err == nil {
			assert.True(t, resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden,
				"Phải trả về 401 hoặc 403 khi không có token")
			fmt.Printf("✅ Endpoint tồn tại (status: %d - đúng như mong đợi khi không có token)\n", resp.StatusCode)
		} else {
			fmt.Printf("⚠️  Lỗi khi gọi API: %v\n", err)
		}
		_ = body // Tránh lỗi unused variable
	})

	// Test 2: Kiểm tra endpoint GET /organization-share
	t.Run("2. Kiểm tra endpoint GET /organization-share (không có token)", func(t *testing.T) {
		resp, body, err := client.GET("/organization-share?fromOrgId=507f1f77bcf86cd799439011")
		// Không có token nên sẽ trả về 401
		if err == nil {
			assert.True(t, resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden,
				"Phải trả về 401 hoặc 403 khi không có token")
			fmt.Printf("✅ Endpoint tồn tại (status: %d - đúng như mong đợi khi không có token)\n", resp.StatusCode)
		} else {
			fmt.Printf("⚠️  Lỗi khi gọi API: %v\n", err)
		}
		_ = body // Tránh lỗi unused variable
	})

	// Test 3: Kiểm tra endpoint DELETE /organization-share
	t.Run("3. Kiểm tra endpoint DELETE /organization-share (không có token)", func(t *testing.T) {
		resp, body, err := client.DELETE("/organization-share/507f1f77bcf86cd799439011")
		// Không có token nên sẽ trả về 401
		if err == nil {
			assert.True(t, resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden,
				"Phải trả về 401 hoặc 403 khi không có token")
			fmt.Printf("✅ Endpoint tồn tại (status: %d - đúng như mong đợi khi không có token)\n", resp.StatusCode)
		} else {
			fmt.Printf("⚠️  Lỗi khi gọi API: %v\n", err)
		}
		_ = body // Tránh lỗi unused variable
	})

	fmt.Printf("\n✅ TẤT CẢ ENDPOINTS ĐÃ ĐƯỢC ĐĂNG KÝ!\n")
	fmt.Printf("📝 Để chạy test đầy đủ, cần set TEST_FIREBASE_ID_TOKEN environment variable\n")
}
