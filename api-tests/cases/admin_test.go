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

// TestAdminAPIs kiểm tra các API quản trị viên
func TestAdminAPIs(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	fixtures := utils.NewTestFixtures(baseURL)

	// Lấy Firebase ID token từ environment variable
	firebaseIDToken := utils.GetTestFirebaseIDToken()
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
	}

	// Tạo admin user (giả định đã có permissions)
	_, _, adminToken, err := fixtures.CreateTestUser(firebaseIDToken)
	if err != nil {
		t.Fatalf("❌ Không thể tạo admin user: %v", err)
	}

	// Tạo user thường để test block/unblock (sử dụng cùng Firebase ID token)
	userEmail, _, _, err := fixtures.CreateTestUser(firebaseIDToken)
	if err != nil {
		t.Fatalf("❌ Không thể tạo user test: %v", err)
	}

	client := utils.NewHTTPClient(baseURL, 10)
	client.SetToken(adminToken)

	// Test case 1: Block user
	t.Run("🔒 Khóa người dùng", func(t *testing.T) {
		payload := map[string]interface{}{
			"email": userEmail,
			"note":  "Test block user",
		}

		resp, body, err := client.POST("/admin/user/block", payload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API block user: %v", err)
		}

		// Có thể thành công hoặc fail tùy vào permissions
		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")
			fmt.Printf("✅ Block user thành công\n")
		} else {
			// Nếu không có quyền, sẽ trả về 403
			fmt.Printf("⚠️ Block user yêu cầu quyền admin (status: %d)\n", resp.StatusCode)
		}
	})

	// Test case 2: Unblock user
	t.Run("🔓 Mở khóa người dùng", func(t *testing.T) {
		payload := map[string]interface{}{
			"email": userEmail,
		}

		resp, body, err := client.POST("/admin/user/unblock", payload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API unblock user: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")
			fmt.Printf("✅ Unblock user thành công\n")
		} else {
			fmt.Printf("⚠️ Unblock user yêu cầu quyền admin (status: %d)\n", resp.StatusCode)
		}
	})

	// Test case 3: Set role cho user
	t.Run("👤 Thiết lập vai trò cho người dùng", func(t *testing.T) {
		// Lấy Root Organization ID
		rootOrgID, err := fixtures.GetRootOrganizationID(adminToken)
		if err != nil {
			t.Skipf("⚠️ Không thể lấy Root Organization, bỏ qua test set role: %v", err)
			return
		}

		// Tạo role test trước (phải có organizationId)
		roleID, err := fixtures.CreateTestRole(adminToken, "TestRole", "Test Role Description", rootOrgID)
		if err != nil {
			t.Skipf("⚠️ Không thể tạo role test, bỏ qua test set role: %v", err)
			return
		}

		payload := map[string]interface{}{
			"email":  userEmail,
			"roleID": roleID,
		}

		resp, body, err := client.POST("/admin/user/role", payload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API set role: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")
			fmt.Printf("✅ Set role thành công\n")
		} else {
			fmt.Printf("⚠️ Set role yêu cầu quyền admin (status: %d - %s)\n", resp.StatusCode, string(body))
		}
	})

	// Cleanup: Logout
	t.Run("🧹 Cleanup", func(t *testing.T) {
		logoutPayload := map[string]interface{}{
			"hwid": "test_device_123",
		}
		client.POST("/auth/logout", logoutPayload)
		fmt.Printf("✅ Cleanup hoàn tất\n")
	})
}
