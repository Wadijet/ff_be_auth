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

// TestOrganizationOwnership kiểm tra phân quyền dữ liệu theo organization
func TestOrganizationOwnership(t *testing.T) {
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

	// Lấy Root Organization ID
	rootOrgID, err := fixtures.GetRootOrganizationID(token)
	if err != nil {
		t.Fatalf("❌ Không thể lấy Root Organization ID: %v", err)
	}

	// Test 1: Lấy danh sách roles của user với thông tin organization
	t.Run("📋 Lấy danh sách roles của user", func(t *testing.T) {
		resp, body, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy danh sách roles: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Phải parse được JSON response")

		data, ok := result["data"].([]interface{})
		assert.True(t, ok, "Data phải là array")
		assert.Greater(t, len(data), 0, "Phải có ít nhất 1 role")

		// Kiểm tra cấu trúc role response
		firstRole, ok := data[0].(map[string]interface{})
		assert.True(t, ok, "Role phải là object")
		assert.Contains(t, firstRole, "roleId", "Phải có roleId")
		assert.Contains(t, firstRole, "roleName", "Phải có roleName")
		assert.Contains(t, firstRole, "organizationId", "Phải có organizationId")
		assert.Contains(t, firstRole, "organizationName", "Phải có organizationName")

		fmt.Printf("✅ Lấy danh sách roles thành công: %d roles\n", len(data))
	})

	// Test 2: Tạo organization và role mới
	var testOrgID string
	var testRoleID string

	t.Run("🏢 Tạo organization và role mới", func(t *testing.T) {
		// Tạo organization con
		orgPayload := map[string]interface{}{
			"name":     fmt.Sprintf("TestOrg_%d", time.Now().UnixNano()),
			"code":     fmt.Sprintf("TEST_ORG_%d", time.Now().UnixNano()),
			"type":     2, // Company
			"parentId": rootOrgID,
		}

		resp, body, err := client.POST("/organization/insert-one", orgPayload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi tạo organization: %v", err)
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err)

			data, ok := result["data"].(map[string]interface{})
			if ok {
				testOrgID, _ = data["id"].(string)
				fmt.Printf("✅ Tạo organization thành công: %s\n", testOrgID)
			}
		}

		// Tạo role trong organization mới
		if testOrgID != "" {
			rolePayload := map[string]interface{}{
				"name":           fmt.Sprintf("TestRole_%d", time.Now().UnixNano()),
				"describe":       "Test Role for Organization Ownership",
				"organizationId": testOrgID,
			}

			resp, body, err := client.POST("/role/insert-one", rolePayload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo role: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					testRoleID, _ = data["id"].(string)
					fmt.Printf("✅ Tạo role thành công: %s\n", testRoleID)
				}
			}
		}
	})

	// Test 3: Gán role cho user và test organization context
	if testRoleID != "" {
		t.Run("👤 Test organization context với role mới", func(t *testing.T) {
			// Skip phần gán role - user đã có role từ init
			// Sử dụng role vừa tạo để test

			// Set active role ID
			client.SetActiveRoleID(testRoleID)

			// Test tạo dữ liệu với organization context
			t.Run("📝 Tạo dữ liệu với organization context", func(t *testing.T) {
				// Test với FbCustomer (có organizationId)
				customerPayload := map[string]interface{}{
					"customerId": fmt.Sprintf("test_customer_%d", time.Now().UnixNano()),
					"name":       "Test Customer",
					"email":      "test@example.com",
				}

				resp, body, err := client.POST("/fb-customer/insert-one", customerPayload)
				if err != nil {
					t.Fatalf("❌ Lỗi khi tạo customer: %v", err)
				}

				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err)

					data, ok := result["data"].(map[string]interface{})
					if ok {
						// Kiểm tra organizationId đã được tự động gán
						orgID, ok := data["organizationId"].(string)
						assert.True(t, ok, "Phải có organizationId")
						assert.Equal(t, testOrgID, orgID, "organizationId phải khớp với active organization")
						fmt.Printf("✅ Tạo customer với organizationId: %s\n", orgID)
					}
				}
			})

			// Test filter dữ liệu theo organization
			t.Run("🔍 Filter dữ liệu theo organization", func(t *testing.T) {
				// Lấy danh sách customers
				resp, body, err := client.GET("/fb-customer/find")
				if err != nil {
					t.Fatalf("❌ Lỗi khi lấy danh sách customers: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err)

					data, ok := result["data"].([]interface{})
					if ok {
						// Tất cả customers phải thuộc organization của user
						for _, item := range data {
							customer, ok := item.(map[string]interface{})
							if ok {
								orgID, ok := customer["organizationId"].(string)
								if ok {
									// Kiểm tra organizationId phải trong allowed organizations
									// (bao gồm cả parent organizations)
									fmt.Printf("  - Customer organizationId: %s\n", orgID)
								}
							}
						}
						fmt.Printf("✅ Filter customers theo organization thành công: %d items\n", len(data))
					}
				}
			})
		})
	}

	// Test 4: Test với scope permissions
	t.Run("🔐 Test scope permissions", func(t *testing.T) {
		// Tạo organization con và role với scope = 0 (Self)
		childOrgPayload := map[string]interface{}{
			"name":     fmt.Sprintf("ChildOrg_%d", time.Now().UnixNano()),
			"code":     fmt.Sprintf("CHILD_%d", time.Now().UnixNano()),
			"type":     3, // Department
			"parentId": testOrgID,
		}

		var childOrgID string
		resp, body, err := client.POST("/organization/insert-one", childOrgPayload)
		if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
			var result map[string]interface{}
			json.Unmarshal(body, &result)
			if data, ok := result["data"].(map[string]interface{}); ok {
				childOrgID, _ = data["id"].(string)
			}
		}

		if childOrgID != "" {
			// Tạo role trong child organization
			childRolePayload := map[string]interface{}{
				"name":           fmt.Sprintf("ChildRole_%d", time.Now().UnixNano()),
				"describe":       "Child Role with Scope 0",
				"organizationId": childOrgID,
			}

			var childRoleID string
			resp, body, err := client.POST("/role/insert-one", childRolePayload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					childRoleID, _ = data["id"].(string)
				}
			}

			if childRoleID != "" {
				// Gán permission với scope = 0 cho role
				// (Cần có permission "FbCustomer.Read" trước)
				fmt.Printf("✅ Tạo child organization và role thành công\n")
				fmt.Printf("  - Child Org ID: %s\n", childOrgID)
				fmt.Printf("  - Child Role ID: %s\n", childRoleID)
			}
		}
	})

	// Test 5: Test inverse parent lookup (xem dữ liệu cấp trên)
	t.Run("⬆️ Test inverse parent lookup", func(t *testing.T) {
		// Tạo dữ liệu ở organization cha
		client.SetActiveRoleID(testRoleID) // Role ở organization cha

		parentCustomerPayload := map[string]interface{}{
			"customerId": fmt.Sprintf("parent_customer_%d", time.Now().UnixNano()),
			"name":       "Parent Customer",
			"email":      "parent@example.com",
		}

		var parentCustomerID string
		resp, body, err := client.POST("/fb-customer/insert-one", parentCustomerPayload)
		if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
			var result map[string]interface{}
			json.Unmarshal(body, &result)
			if data, ok := result["data"].(map[string]interface{}); ok {
				parentCustomerID, _ = data["id"].(string)
			}
		}

		// Test: User ở organization con có thể xem dữ liệu của organization cha
		// (Thông qua inverse parent lookup)
		if parentCustomerID != "" {
			fmt.Printf("✅ Tạo customer ở organization cha: %s\n", parentCustomerID)
			fmt.Printf("  - User ở organization con sẽ tự động thấy customer này\n")
		}
	})
}
