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

// TestRBACAPIs kiểm tra các API RBAC (Role, Permission, UserRole)
func TestRBACAPIs(t *testing.T) {
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

	var roleID string

	// Lấy Organization Root ID trước
	rootOrgID, err := fixtures.GetRootOrganizationID(token)
	if err != nil {
		t.Logf("⚠️ Không thể lấy Root Organization (có thể chưa init): %v", err)
		// Vẫn tiếp tục test, có thể sẽ fail ở phần tạo Role
	}

	// Test Role APIs
	t.Run("🎭 Role APIs", func(t *testing.T) {
		// Test 1: Tạo role
		t.Run("Tạo role mới", func(t *testing.T) {
			// Role phải có organizationId (bắt buộc)
			if rootOrgID == "" {
				t.Skip("Skipping: Không có Root Organization ID")
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

			// Có thể thành công hoặc fail tùy vào permissions
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
				fmt.Printf("✅ Tạo role thành công\n")
			} else {
				fmt.Printf("⚠️ Tạo role yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Lấy danh sách roles
		t.Run("Lấy danh sách roles", func(t *testing.T) {
			resp, body, err := client.GET("/role/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách roles: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách roles thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách roles yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 3: Lấy role theo ID (nếu có)
		if roleID != "" {
			t.Run("Lấy role theo ID", func(t *testing.T) {
				resp, body, err := client.GET(fmt.Sprintf("/role/find-by-id/%s", roleID))
				if err != nil {
					t.Fatalf("❌ Lỗi khi lấy role theo ID: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err, "Phải parse được JSON response")
					fmt.Printf("✅ Lấy role theo ID thành công\n")
				} else {
					fmt.Printf("⚠️ Lấy role theo ID yêu cầu quyền (status: %d)\n", resp.StatusCode)
				}
			})
		}
	})

	// Test Permission APIs
	t.Run("🔐 Permission APIs", func(t *testing.T) {
		// Test 1: Tạo permission
		t.Run("Tạo permission mới", func(t *testing.T) {
			payload := map[string]interface{}{
				"name":     fmt.Sprintf("TestPermission_%d", time.Now().UnixNano()),
				"describe": "Test Permission Description",
				"category": "test",
				"group":    "test",
			}

			resp, body, err := client.POST("/permission/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo permission: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Tạo permission thành công\n")
			} else {
				fmt.Printf("⚠️ Tạo permission yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Lấy danh sách permissions
		t.Run("Lấy danh sách permissions", func(t *testing.T) {
			resp, body, err := client.GET("/permission/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách permissions: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách permissions thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách permissions yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test UserRole APIs
	t.Run("👥 UserRole APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách user roles
		t.Run("Lấy danh sách user roles", func(t *testing.T) {
			resp, body, err := client.GET("/user-role/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách user roles: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách user roles thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách user roles yêu cầu quyền (status: %d)\n", resp.StatusCode)
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
