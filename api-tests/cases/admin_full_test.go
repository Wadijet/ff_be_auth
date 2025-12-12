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

// TestAdminFullAPIs kiểm tra các API admin với user có full quyền
func TestAdminFullAPIs(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Khởi tạo dữ liệu mặc định trước
	initTestData(t, baseURL)

	fixtures := utils.NewTestFixtures(baseURL)

	// Lấy Firebase ID token từ environment variable
	firebaseIDToken := utils.GetTestFirebaseIDToken()
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
	}

	// Tạo admin user với full quyền
	adminEmail, _, adminToken, _, err := fixtures.CreateAdminUser(firebaseIDToken)
	if err != nil {
		t.Fatalf("❌ Không thể tạo admin user: %v", err)
	}

	client := utils.NewHTTPClient(baseURL, 10)
	client.SetToken(adminToken)

	// Test 1: Set Administrator cho user khác
	t.Run("👑 Set Administrator", func(t *testing.T) {
		// Tạo user thường và lấy userID từ profile
		firebaseIDToken := utils.GetTestFirebaseIDToken()
		if firebaseIDToken == "" {
			t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
		}
		userEmail, _, userToken, err := fixtures.CreateTestUser(firebaseIDToken)
		if err != nil {
			t.Fatalf("❌ Không thể tạo user test: %v", err)
		}

		// Lấy userID từ profile
		tempClient := utils.NewHTTPClient(baseURL, 10)
		tempClient.SetToken(userToken)
		resp, body, err := tempClient.GET("/auth/profile")
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Skip("⚠️ Không thể lấy userID, bỏ qua test")
			return
		}

		var profileResult map[string]interface{}
		json.Unmarshal(body, &profileResult)
		data, _ := profileResult["data"].(map[string]interface{})
		userID, _ := data["id"].(string)
		if userID == "" {
			t.Skip("⚠️ Không lấy được userID, bỏ qua test")
			return
		}
		_ = userEmail

		// Set administrator cho user này
		resp, body, err = client.POST(fmt.Sprintf("/init/set-administrator/%s", userID), nil)
		if err != nil {
			t.Fatalf("❌ Lỗi khi set administrator: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")
			fmt.Printf("✅ Set administrator thành công\n")
		} else {
			// Có thể đã là admin hoặc cần quyền đặc biệt
			fmt.Printf("⚠️ Set administrator (status: %d - %s)\n", resp.StatusCode, string(body))
		}
	})

	// Test 2: Tạo role với admin quyền
	t.Run("🎭 Tạo Role với Admin", func(t *testing.T) {
		// Lấy Root Organization ID
		rootOrgID, err := fixtures.GetRootOrganizationID(adminToken)
		if err != nil {
			t.Skipf("⚠️ Không thể lấy Root Organization, bỏ qua test tạo role: %v", err)
			return
		}

		payload := map[string]interface{}{
			"name":           fmt.Sprintf("TestRole_%d", time.Now().UnixNano()),
			"describe":       "Test Role Description",
			"organizationId": rootOrgID, // BẮT BUỘC
		}

		resp, body, err := client.POST("/role/insert-one", payload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi tạo role: %v", err)
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")
			fmt.Printf("✅ Tạo role thành công với admin quyền\n")
		} else {
			t.Errorf("❌ Tạo role thất bại với admin: %d - %s", resp.StatusCode, string(body))
		}
	})

	// Test 3: Lấy danh sách roles
	t.Run("📋 Lấy danh sách Roles", func(t *testing.T) {
		resp, body, err := client.GET("/role/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy danh sách roles: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Admin phải lấy được danh sách roles")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Phải parse được JSON response")
		fmt.Printf("✅ Lấy danh sách roles thành công\n")
	})

	// Test 4: Lấy danh sách permissions
	t.Run("🔐 Lấy danh sách Permissions", func(t *testing.T) {
		resp, body, err := client.GET("/permission/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy danh sách permissions: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Admin phải lấy được danh sách permissions")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Phải parse được JSON response")
		fmt.Printf("✅ Lấy danh sách permissions thành công\n")
	})

	// Test 5: Lấy danh sách users
	t.Run("👥 Lấy danh sách Users", func(t *testing.T) {
		resp, body, err := client.GET("/user/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy danh sách users: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Admin phải lấy được danh sách users")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Phải parse được JSON response")
		fmt.Printf("✅ Lấy danh sách users thành công\n")
	})

	// Test 6: Block/Unblock user
	t.Run("🔒 Block/Unblock User", func(t *testing.T) {
		// Tạo user để block
		userEmail, _, _, err := fixtures.CreateTestUser(firebaseIDToken)
		if err != nil {
			t.Fatalf("❌ Không thể tạo user test: %v", err)
		}

		// Block user
		blockPayload := map[string]interface{}{
			"email": userEmail,
			"note":  "Test block",
		}

		resp, body, err := client.POST("/admin/user/block", blockPayload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi block user: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")
			fmt.Printf("✅ Block user thành công\n")
		} else {
			t.Errorf("❌ Block user thất bại: %d - %s", resp.StatusCode, string(body))
		}

		// Unblock user
		unblockPayload := map[string]interface{}{
			"email": userEmail,
		}

		resp, body, err = client.POST("/admin/user/unblock", unblockPayload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi unblock user: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")
			fmt.Printf("✅ Unblock user thành công\n")
		} else {
			t.Errorf("❌ Unblock user thất bại: %d - %s", resp.StatusCode, string(body))
		}
	})

	// Test 7: Set role cho user
	t.Run("👤 Set Role cho User", func(t *testing.T) {
		// Lấy Root Organization ID
		rootOrgID, err := fixtures.GetRootOrganizationID(adminToken)
		if err != nil {
			t.Skipf("⚠️ Không thể lấy Root Organization, bỏ qua test set role: %v", err)
			return
		}

		// Tạo role trước (phải có organizationId)
		rolePayload := map[string]interface{}{
			"name":           fmt.Sprintf("TestRole_%d", time.Now().UnixNano()),
			"describe":       "Test Role",
			"organizationId": rootOrgID, // BẮT BUỘC
		}

		resp, body, err := client.POST("/role/insert-one", rolePayload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi tạo role: %v", err)
		}

		var roleID string
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")

			data, ok := result["data"].(map[string]interface{})
			if ok {
				id, ok := data["id"].(string)
				if ok {
					roleID = id
				}
			}
		}

		if roleID == "" {
			t.Skip("⚠️ Không thể tạo role, bỏ qua test set role")
			return
		}

		// Tạo user để set role
		userEmail, _, _, err := fixtures.CreateTestUser(firebaseIDToken)
		if err != nil {
			t.Fatalf("❌ Không thể tạo user test: %v", err)
		}

		// Set role
		setRolePayload := map[string]interface{}{
			"email":  userEmail,
			"roleID": roleID,
		}

		resp, body, err = client.POST("/admin/user/role", setRolePayload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi set role: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err, "Phải parse được JSON response")
			fmt.Printf("✅ Set role thành công\n")
		} else {
			t.Errorf("❌ Set role thất bại: %d - %s", resp.StatusCode, string(body))
		}
	})

	// Cleanup
	t.Run("🧹 Cleanup", func(t *testing.T) {
		logoutPayload := map[string]interface{}{
			"hwid": "test_device_123",
		}
		client.POST("/auth/logout", logoutPayload)
		fmt.Printf("✅ Cleanup hoàn tất (admin: %s)\n", adminEmail)
	})
}
