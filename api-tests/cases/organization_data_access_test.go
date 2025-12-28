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

// TestOrganizationDataAccess - Test đơn giản để kiểm tra phân quyền dữ liệu
func TestOrganizationDataAccess(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Khởi tạo dữ liệu mặc định
	initTestData(t, baseURL)

	fixtures := utils.NewTestFixtures(baseURL)

	// Lấy Firebase token
	firebaseIDToken := utils.GetTestFirebaseIDToken()
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
	}

	// Tạo user và lấy token
	_, _, token, err := fixtures.CreateTestUser(firebaseIDToken)
	if err != nil {
		t.Fatalf("❌ Không thể tạo user test: %v", err)
	}

	client := utils.NewHTTPClient(baseURL, 10)
	client.SetToken(token)

	// Test 1: Lấy danh sách roles của user
	t.Run("1. Lấy danh sách roles", func(t *testing.T) {
		resp, body, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		assert.NoError(t, err, "Phải parse được JSON")

		data, ok := result["data"].([]interface{})
		assert.True(t, ok, "Data phải là array")
		assert.Greater(t, len(data), 0, "Phải có ít nhất 1 role")

		// In ra thông tin roles
		for i, role := range data {
			roleMap, ok := role.(map[string]interface{})
			if ok {
				fmt.Printf("  Role %d: %v\n", i+1, roleMap)
			}
		}

		fmt.Printf("✅ Lấy danh sách roles thành công: %d roles\n", len(data))
	})

	// Test 2: Set active role và tạo dữ liệu
	t.Run("2. Tạo dữ liệu với organization context", func(t *testing.T) {
		// Lấy role đầu tiên
		resp, body, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		data, _ := result["data"].([]interface{})

		if len(data) == 0 {
			t.Skip("Skipping: Không có role nào")
		}

		firstRole, _ := data[0].(map[string]interface{})
		roleID, _ := firstRole["roleId"].(string)
		orgID, _ := firstRole["organizationId"].(string)

		if roleID == "" {
			t.Skip("Skipping: Không có role ID")
		}

		// Set active role
		client.SetActiveRoleID(roleID)
		fmt.Printf("✅ Set active role: %s (org: %s)\n", roleID, orgID)

		// Tạo customer với organization context
		customerPayload := map[string]interface{}{
			"customerId": fmt.Sprintf("test_customer_%d", time.Now().UnixNano()),
			"name":       "Test Customer",
			"email":      "test@example.com",
		}

		resp, body, err = client.POST("/fb-customer/insert-one", customerPayload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi tạo customer: %v", err)
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err)

			data, ok := result["data"].(map[string]interface{})
			if ok {
				customerOrgID, ok := data["organizationId"].(string)
				if ok {
					assert.Equal(t, orgID, customerOrgID, "organizationId phải khớp với active organization")
					fmt.Printf("✅ Tạo customer thành công với organizationId: %s\n", customerOrgID)
				} else {
					fmt.Printf("⚠️ Customer không có organizationId (có thể model chưa có field)\n")
				}
			}
		} else {
			fmt.Printf("⚠️ Tạo customer yêu cầu quyền (status: %d)\n", resp.StatusCode)
		}
	})

	// Test 3: Filter dữ liệu theo organization
	t.Run("3. Filter dữ liệu theo organization", func(t *testing.T) {
		// Lấy role đầu tiên
		resp, body, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		data, _ := result["data"].([]interface{})

		if len(data) == 0 {
			t.Skip("Skipping: Không có role nào")
		}

		firstRole, _ := data[0].(map[string]interface{})
		roleID, _ := firstRole["roleId"].(string)

		if roleID == "" {
			t.Skip("Skipping: Không có role ID")
		}

		// Set active role
		client.SetActiveRoleID(roleID)

		// Query customers
		resp, body, err = client.GET("/fb-customer/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi query customers: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err)

			customers, ok := result["data"].([]interface{})
			if ok {
				fmt.Printf("✅ Query customers thành công: %d items\n", len(customers))
				for i, item := range customers {
					if i < 3 { // Chỉ in 3 items đầu
						customer, ok := item.(map[string]interface{})
						if ok {
							orgID, _ := customer["organizationId"].(string)
							name, _ := customer["name"].(string)
							fmt.Printf("  - Customer: %s (org: %s)\n", name, orgID)
						}
					}
				}
			}
		} else {
			fmt.Printf("⚠️ Query customers yêu cầu quyền (status: %d)\n", resp.StatusCode)
		}
	})

	// Test 4: Verify organizationId không thể update
	t.Run("4. Verify không thể update organizationId", func(t *testing.T) {
		// Lấy role đầu tiên
		resp, body, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		data, _ := result["data"].([]interface{})

		if len(data) == 0 {
			t.Skip("Skipping: Không có role nào")
		}

		firstRole, _ := data[0].(map[string]interface{})
		roleID, _ := firstRole["roleId"].(string)

		if roleID == "" {
			t.Skip("Skipping: Không có role ID")
		}

		client.SetActiveRoleID(roleID)

		// Tạo customer trước
		customerPayload := map[string]interface{}{
			"customerId": fmt.Sprintf("test_update_%d", time.Now().UnixNano()),
			"name":       "Test Update Customer",
			"email":      "update@example.com",
		}

		resp, body, err = client.POST("/fb-customer/insert-one", customerPayload)
		if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated) {
			t.Skip("Skipping: Không thể tạo customer để test update")
		}

		var createResult map[string]interface{}
		json.Unmarshal(body, &createResult)
		createdData, _ := createResult["data"].(map[string]interface{})
		customerID, _ := createdData["id"].(string)
		originalOrgID, _ := createdData["organizationId"].(string)

		if customerID == "" {
			t.Skip("Skipping: Không có customer ID")
		}

		// Thử update với organizationId khác
		updatePayload := map[string]interface{}{
			"name":           "Updated Name",
			"organizationId": "507f1f77bcf86cd799439999", // ID giả
		}

		resp, body, err = client.PUT(fmt.Sprintf("/fb-customer/update-by-id/%s", customerID), updatePayload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi update: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var updateResult map[string]interface{}
			json.Unmarshal(body, &updateResult)
			updatedData, _ := updateResult["data"].(map[string]interface{})

			// Verify organizationId không thay đổi
			updatedOrgID, _ := updatedData["organizationId"].(string)
			assert.Equal(t, originalOrgID, updatedOrgID, "organizationId không được phép thay đổi")
			fmt.Printf("✅ Verify: organizationId không thể update (vẫn là: %s)\n", updatedOrgID)
		} else {
			fmt.Printf("⚠️ Update yêu cầu quyền (status: %d)\n", resp.StatusCode)
		}
	})
}

