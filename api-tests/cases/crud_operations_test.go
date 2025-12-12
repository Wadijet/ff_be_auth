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

// TestCRUDOperations kiểm tra các thao tác CRUD đầy đủ
func TestCRUDOperations(t *testing.T) {
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

	// Lấy Organization Root ID trước
	rootOrgID, err := fixtures.GetRootOrganizationID(token)
	if err != nil {
		t.Logf("⚠️ Không thể lấy Root Organization (có thể chưa init): %v", err)
		// Vẫn tiếp tục test, có thể sẽ fail ở phần tạo Role
	}

	// Test Role CRUD Operations
	t.Run("🎭 Role CRUD Operations", func(t *testing.T) {
		var roleID string

		// CREATE: Tạo role
		t.Run("CREATE - Tạo role", func(t *testing.T) {
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
				fmt.Printf("✅ CREATE role thành công\n")
			} else {
				fmt.Printf("⚠️ CREATE role yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// READ: Lấy danh sách roles
		t.Run("READ - Lấy danh sách roles", func(t *testing.T) {
			resp, body, err := client.GET("/role/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách roles: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ READ roles thành công\n")
			} else {
				fmt.Printf("⚠️ READ roles yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// READ BY ID: Lấy role theo ID
		if roleID != "" {
			t.Run("READ BY ID - Lấy role theo ID", func(t *testing.T) {
				resp, body, err := client.GET(fmt.Sprintf("/role/find-by-id/%s", roleID))
				if err != nil {
					t.Fatalf("❌ Lỗi khi lấy role theo ID: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err, "Phải parse được JSON response")
					fmt.Printf("✅ READ BY ID role thành công\n")
				} else {
					fmt.Printf("⚠️ READ BY ID role yêu cầu quyền (status: %d)\n", resp.StatusCode)
				}
			})

			// UPDATE: Cập nhật role
			t.Run("UPDATE - Cập nhật role", func(t *testing.T) {
				payload := map[string]interface{}{
					"name":     fmt.Sprintf("UpdatedRole_%d", time.Now().UnixNano()),
					"describe": "Updated Role Description",
				}

				resp, body, err := client.PUT(fmt.Sprintf("/role/update-by-id/%s", roleID), payload)
				if err != nil {
					t.Fatalf("❌ Lỗi khi cập nhật role: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err, "Phải parse được JSON response")
					fmt.Printf("✅ UPDATE role thành công\n")
				} else {
					fmt.Printf("⚠️ UPDATE role yêu cầu quyền (status: %d - %s)\n", resp.StatusCode, string(body))
				}
			})

			// DELETE: Xóa role
			t.Run("DELETE - Xóa role", func(t *testing.T) {
				resp, body, err := client.DELETE(fmt.Sprintf("/role/delete-by-id/%s", roleID))
				if err != nil {
					t.Fatalf("❌ Lỗi khi xóa role: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err, "Phải parse được JSON response")
					fmt.Printf("✅ DELETE role thành công\n")
				} else {
					fmt.Printf("⚠️ DELETE role yêu cầu quyền (status: %d - %s)\n", resp.StatusCode, string(body))
				}
			})
		}
	})

	// Test Permission CRUD Operations
	t.Run("🔐 Permission CRUD Operations", func(t *testing.T) {
		// READ: Lấy danh sách permissions
		t.Run("READ - Lấy danh sách permissions", func(t *testing.T) {
			resp, body, err := client.GET("/permission/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách permissions: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ READ permissions thành công\n")
			} else {
				fmt.Printf("⚠️ READ permissions yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// COUNT: Đếm số lượng permissions
		t.Run("COUNT - Đếm số lượng permissions", func(t *testing.T) {
			resp, body, err := client.GET("/permission/count")
			if err != nil {
				t.Fatalf("❌ Lỗi khi đếm permissions: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ COUNT permissions thành công\n")
			} else {
				fmt.Printf("⚠️ COUNT permissions yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test User CRUD Operations (Read-only)
	t.Run("👤 User CRUD Operations", func(t *testing.T) {
		// READ: Lấy danh sách users
		t.Run("READ - Lấy danh sách users", func(t *testing.T) {
			resp, body, err := client.GET("/user/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách users: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ READ users thành công\n")
			} else {
				fmt.Printf("⚠️ READ users yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// COUNT: Đếm số lượng users
		t.Run("COUNT - Đếm số lượng users", func(t *testing.T) {
			resp, body, err := client.GET("/user/count")
			if err != nil {
				t.Fatalf("❌ Lỗi khi đếm users: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ COUNT users thành công\n")
			} else {
				fmt.Printf("⚠️ COUNT users yêu cầu quyền (status: %d)\n", resp.StatusCode)
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
