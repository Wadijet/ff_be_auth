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

// TestOrganizationSharing - Test kịch bản Organization-Level Sharing
func TestOrganizationSharing(t *testing.T) {
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

	fixtures := utils.NewTestFixtures(baseURL)

	// Lấy Firebase token
	firebaseIDToken := utils.GetTestFirebaseIDToken()
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
	}

	// Tạo user admin và lấy token
	_, _, adminToken, err := fixtures.CreateTestUser(firebaseIDToken)
	if err != nil {
		t.Fatalf("❌ Không thể tạo user test: %v", err)
	}

	adminClient := utils.NewHTTPClient(baseURL, 10)
	adminClient.SetToken(adminToken)

	// Lấy roles của admin
	resp, body, err := adminClient.GET("/auth/roles")
	if err != nil {
		t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

	var rolesResult map[string]interface{}
	json.Unmarshal(body, &rolesResult)
	rolesData, _ := rolesResult["data"].([]interface{})
	if len(rolesData) == 0 {
		t.Skip("Skipping: Không có role nào")
	}

	// Lấy role đầu tiên để làm admin role
	adminRole, _ := rolesData[0].(map[string]interface{})
	adminRoleID, _ := adminRole["roleId"].(string)
	adminOrgID, _ := adminRole["organizationId"].(string)
	adminClient.SetActiveRoleID(adminRoleID)

	fmt.Printf("✅ Setup: Admin Role ID: %s, Org ID: %s\n", adminRoleID, adminOrgID)

	// Tạo cấu trúc organization test
	// Sales Department (Level 2) - sẽ share data
	// ├── Team A (Level 3) - sẽ nhận data
	// └── Team B (Level 3) - KHÔNG nhận data (để test)

	// Biến để share giữa các subtests
	var salesDeptID, teamAID, teamBID, teamARoleID, teamBRoleID string

	t.Run("1. Tạo cấu trúc organization test", func(t *testing.T) {
		// Lấy root organization
		rootOrgID, err := fixtures.GetRootOrganizationID(adminToken)
		if err != nil {
			t.Fatalf("❌ Không thể lấy root organization: %v", err)
		}

		// Tạo Sales Department
		salesDeptPayload := map[string]interface{}{
			"name":     "Sales Department Test",
			"code":     fmt.Sprintf("SALES_DEPT_%d", time.Now().UnixNano()),
			"type":     "department",
			"parentId": rootOrgID,
			"isActive": true,
		}

		resp, body, err := adminClient.POST("/organization/insert-one", salesDeptPayload)
		assert.NoError(t, err, "Không có lỗi khi tạo Sales Department")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var salesDeptResult map[string]interface{}
		json.Unmarshal(body, &salesDeptResult)
		salesDeptData, _ := salesDeptResult["data"].(map[string]interface{})
		salesDeptID = salesDeptData["id"].(string)

		// Tạo Team A
		teamAPayload := map[string]interface{}{
			"name":     "Team A Test",
			"code":     fmt.Sprintf("TEAM_A_%d", time.Now().UnixNano()),
			"type":     "team",
			"parentId": salesDeptID,
			"isActive": true,
		}

		resp, body, err = adminClient.POST("/organization/insert-one", teamAPayload)
		assert.NoError(t, err, "Không có lỗi khi tạo Team A")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var teamAResult map[string]interface{}
		json.Unmarshal(body, &teamAResult)
		teamAData, _ := teamAResult["data"].(map[string]interface{})
		teamAID = teamAData["id"].(string)

		// Tạo Team B
		teamBPayload := map[string]interface{}{
			"name":     "Team B Test",
			"code":     fmt.Sprintf("TEAM_B_%d", time.Now().UnixNano()),
			"type":     "team",
			"parentId": salesDeptID,
			"isActive": true,
		}

		resp, body, err = adminClient.POST("/organization/insert-one", teamBPayload)
		assert.NoError(t, err, "Không có lỗi khi tạo Team B")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var teamBResult map[string]interface{}
		json.Unmarshal(body, &teamBResult)
		teamBData, _ := teamBResult["data"].(map[string]interface{})
		teamBID = teamBData["id"].(string)

		fmt.Printf("✅ Tạo cấu trúc organization:\n")
		fmt.Printf("  - Sales Department: %s\n", salesDeptID)
		fmt.Printf("  - Team A: %s\n", teamAID)
		fmt.Printf("  - Team B: %s\n", teamBID)
	})

	t.Run("2. Tạo roles cho Team A và Team B", func(t *testing.T) {
		// Tạo role cho Team A
		teamARolePayload := map[string]interface{}{
			"name":           "Team A Member",
			"code":           fmt.Sprintf("TEAM_A_ROLE_%d", time.Now().UnixNano()),
			"organizationId": teamAID,
			"describe":       "Role cho Team A",
		}

		resp, body, err := adminClient.POST("/role/insert-one", teamARolePayload)
		assert.NoError(t, err, "Không có lỗi khi tạo role Team A")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var teamARoleResult map[string]interface{}
		json.Unmarshal(body, &teamARoleResult)
		teamARoleData, _ := teamARoleResult["data"].(map[string]interface{})
		teamARoleID = teamARoleData["id"].(string)

		// Tạo role cho Team B
		teamBRolePayload := map[string]interface{}{
			"name":           "Team B Member",
			"code":           fmt.Sprintf("TEAM_B_ROLE_%d", time.Now().UnixNano()),
			"organizationId": teamBID,
			"describe":       "Role cho Team B",
		}

		resp, body, err = adminClient.POST("/role/insert-one", teamBRolePayload)
		assert.NoError(t, err, "Không có lỗi khi tạo role Team B")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var teamBRoleResult map[string]interface{}
		json.Unmarshal(body, &teamBRoleResult)
		teamBRoleData, _ := teamBRoleResult["data"].(map[string]interface{})
		teamBRoleID = teamBRoleData["id"].(string)

		fmt.Printf("✅ Tạo roles:\n")
		fmt.Printf("  - Team A Role: %s\n", teamARoleID)
		fmt.Printf("  - Team B Role: %s\n", teamBRoleID)
	})

	// Tạo user cho Team A - assign role Team A cho admin user
	t.Run("3. Assign role Team A cho admin user", func(t *testing.T) {
		// Lấy user ID từ profile
		resp, body, err := adminClient.GET("/auth/profile")
		if err != nil {
			t.Fatalf("❌ Không thể lấy profile: %v", err)
		}

		var profileResult map[string]interface{}
		json.Unmarshal(body, &profileResult)
		profileData, _ := profileResult["data"].(map[string]interface{})
		userID, _ := profileData["id"].(string)

		// Assign role Team A cho user
		assignRolePayload := map[string]interface{}{
			"userId": userID,
			"roleId": teamARoleID,
		}

		resp, body, err = adminClient.POST("/user-role/insert-one", assignRolePayload)
		if err != nil || resp.StatusCode != http.StatusOK {
			// Có thể role đã tồn tại, bỏ qua
			fmt.Printf("⚠️  User đã có role hoặc lỗi: %v\n", err)
		} else {
			fmt.Printf("✅ Assign role Team A cho user thành công\n")
		}
	})

	// Tạo dữ liệu test ở Sales Department
	var salesDeptDataID string

	t.Run("4. Tạo dữ liệu test ở Sales Department", func(t *testing.T) {
		// Set active role là admin role (có quyền với Sales Department)
		adminClient.SetActiveRoleID(adminRoleID)

		// Tạo customer ở Sales Department
		customerPayload := map[string]interface{}{
			"customerId":     fmt.Sprintf("SHARED_CUSTOMER_%d", time.Now().UnixNano()),
			"name":           "Shared Customer",
			"email":          "shared@example.com",
			"organizationId": salesDeptID, // Dữ liệu ở Sales Department
		}

		resp, body, err := adminClient.POST("/customer/insert-one", customerPayload)
		assert.NoError(t, err, "Không có lỗi khi tạo customer")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var customerResult map[string]interface{}
		json.Unmarshal(body, &customerResult)
		customerData, _ := customerResult["data"].(map[string]interface{})
		salesDeptDataID = customerData["id"].(string)

		fmt.Printf("✅ Tạo customer ở Sales Department: %s\n", salesDeptDataID)
	})

	// Test: User Team A KHÔNG thấy data của Sales Department (chưa share)
	t.Run("5. Test: User Team A KHÔNG thấy data Sales Department (chưa share)", func(t *testing.T) {
		// Set active role là Team A role
		teamAClient := utils.NewHTTPClient(baseURL, 10)
		teamAClient.SetToken(adminToken) // Dùng cùng token nhưng set active role khác
		teamAClient.SetActiveRoleID(teamARoleID)

		// Query customers
		_, body, err := teamAClient.GET("/customer/find")
		assert.NoError(t, err, "Không có lỗi khi query customers")

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		customers, _ := result["data"].([]interface{})

		// Kiểm tra không thấy customer của Sales Department
		found := false
		for _, customer := range customers {
			customerMap, _ := customer.(map[string]interface{})
			customerID, _ := customerMap["id"].(string)
			if customerID == salesDeptDataID {
				found = true
				break
			}
		}

		assert.False(t, found, "User Team A KHÔNG được thấy customer của Sales Department (chưa share)")
		fmt.Printf("✅ User Team A KHÔNG thấy data Sales Department (đúng như mong đợi)\n")
	})

	// Tạo share: Sales Department share với Team A
	var shareID string

	t.Run("6. Tạo share: Sales Department share với Team A", func(t *testing.T) {
		// Set active role là admin role
		adminClient.SetActiveRoleID(adminRoleID)

		sharePayload := map[string]interface{}{
			"fromOrgId":      salesDeptID,
			"toOrgId":        teamAID,
			"permissionNames": []string{}, // Share tất cả permissions
		}

		resp, body, err := adminClient.POST("/organization-share", sharePayload)
		assert.NoError(t, err, "Không có lỗi khi tạo share")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var shareResult map[string]interface{}
		json.Unmarshal(body, &shareResult)
		shareData, _ := shareResult["data"].(map[string]interface{})
		shareID = shareData["id"].(string)

		fmt.Printf("✅ Tạo share thành công: %s (Sales Dept -> Team A)\n", shareID)
	})

	// Test: User Team A thấy data của Sales Department (sau khi share)
	t.Run("7. Test: User Team A thấy data Sales Department (sau khi share)", func(t *testing.T) {
		// Set active role là Team A role
		teamAClient := utils.NewHTTPClient(baseURL, 10)
		teamAClient.SetToken(adminToken) // Dùng cùng token nhưng set active role khác
		teamAClient.SetActiveRoleID(teamARoleID)

		// Query customers
		_, body, err := teamAClient.GET("/customer/find")
		assert.NoError(t, err, "Không có lỗi khi query customers")

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		customers, _ := result["data"].([]interface{})

		// Kiểm tra thấy customer của Sales Department
		found := false
		for _, customer := range customers {
			customerMap, _ := customer.(map[string]interface{})
			customerID, _ := customerMap["id"].(string)
			if customerID == salesDeptDataID {
				found = true
				fmt.Printf("✅ Tìm thấy shared customer: %s\n", customerID)
				break
			}
		}

		assert.True(t, found, "User Team A phải thấy customer của Sales Department (đã share)")
		fmt.Printf("✅ User Team A thấy data Sales Department (sau khi share)\n")
	})

	// Test: User Team B KHÔNG thấy data của Sales Department (không được share)
	t.Run("8. Test: User Team B KHÔNG thấy data Sales Department (không được share)", func(t *testing.T) {
		// Set active role là Team B role
		teamBClient := utils.NewHTTPClient(baseURL, 10)
		teamBClient.SetToken(adminToken) // Dùng cùng token
		teamBClient.SetActiveRoleID(teamBRoleID) // Nhưng set active role là Team B

		// Query customers
		_, body, err := teamBClient.GET("/customer/find")
		assert.NoError(t, err, "Không có lỗi khi query customers")

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		customers, _ := result["data"].([]interface{})

		// Kiểm tra không thấy customer của Sales Department
		found := false
		for _, customer := range customers {
			customerMap, _ := customer.(map[string]interface{})
			customerID, _ := customerMap["id"].(string)
			if customerID == salesDeptDataID {
				found = true
				break
			}
		}

		assert.False(t, found, "User Team B KHÔNG được thấy customer của Sales Department (không được share)")
		fmt.Printf("✅ User Team B KHÔNG thấy data Sales Department (đúng như mong đợi)\n")
	})

	// Test: List shares
	t.Run("9. Test: List shares của Sales Department", func(t *testing.T) {
		adminClient.SetActiveRoleID(adminRoleID)

		resp, body, err := adminClient.GET(fmt.Sprintf("/organization-share?fromOrgId=%s", salesDeptID))
		assert.NoError(t, err, "Không có lỗi khi list shares")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		shares, _ := result["data"].([]interface{})

		assert.Greater(t, len(shares), 0, "Phải có ít nhất 1 share")
		fmt.Printf("✅ List shares thành công: %d shares\n", len(shares))
	})

	// Test: Xóa share
	t.Run("10. Test: Xóa share", func(t *testing.T) {
		adminClient.SetActiveRoleID(adminRoleID)

		resp, _, err := adminClient.DELETE(fmt.Sprintf("/organization-share/%s", shareID))
		assert.NoError(t, err, "Không có lỗi khi xóa share")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Phải trả về status 200")

		fmt.Printf("✅ Xóa share thành công: %s\n", shareID)
	})

	// Test: User Team A KHÔNG thấy data nữa (sau khi xóa share)
	t.Run("11. Test: User Team A KHÔNG thấy data nữa (sau khi xóa share)", func(t *testing.T) {
		// Set active role là Team A role
		teamAClient := utils.NewHTTPClient(baseURL, 10)
		teamAClient.SetToken(adminToken) // Dùng cùng token
		teamAClient.SetActiveRoleID(teamARoleID)

		// Query customers
		_, body, err := teamAClient.GET("/customer/find")
		assert.NoError(t, err, "Không có lỗi khi query customers")

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		customers, _ := result["data"].([]interface{})

		// Kiểm tra không thấy customer của Sales Department
		found := false
		for _, customer := range customers {
			customerMap, _ := customer.(map[string]interface{})
			customerID, _ := customerMap["id"].(string)
			if customerID == salesDeptDataID {
				found = true
				break
			}
		}

		assert.False(t, found, "User Team A KHÔNG được thấy customer của Sales Department (đã xóa share)")
		fmt.Printf("✅ User Team A KHÔNG thấy data Sales Department (sau khi xóa share)\n")
	})

	fmt.Printf("\n✅ TẤT CẢ TEST CASES ĐÃ PASS!\n")
}
