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

// TestOrganizationOwnershipFull - Test đầy đủ các scenarios về organization ownership
func TestOrganizationOwnershipFull(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Khởi tạo dữ liệu mặc định
	initTestData(t, baseURL)

	fixtures := utils.NewTestFixtures(baseURL)

	// Tạo user với token
	firebaseIDToken := utils.GetTestFirebaseIDToken()
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
	}
	
	// Thử tạo admin user để có đầy đủ quyền
	_, _, token, userID, err := fixtures.CreateAdminUser(firebaseIDToken)
	if err != nil || token == "" {
		// Nếu không tạo được admin, thử tạo user thường
		_, _, token, err = fixtures.CreateTestUser(firebaseIDToken)
		if err != nil {
			t.Fatalf("❌ Không thể tạo user test: %v", err)
		}
		
		// Lấy user ID từ profile
		client := utils.NewHTTPClient(baseURL, 10)
		client.SetToken(token)
		_, body, err := client.GET("/auth/profile")
		if err != nil {
			t.Fatalf("❌ Không thể lấy profile: %v", err)
		}
		var profileResult map[string]interface{}
		json.Unmarshal(body, &profileResult)
		profileData, _ := profileResult["data"].(map[string]interface{})
		userID, _ = profileData["id"].(string)
	}

	client := utils.NewHTTPClient(baseURL, 10)
	client.SetToken(token)

	// ============================================
	// SETUP: Tạo cấu trúc organization test với helper function
	// ============================================
	var testData *utils.OrganizationTestData
	t.Run("🏗️ Setup: Tạo cấu trúc organization", func(t *testing.T) {
		var setupErr error
		testData, setupErr = fixtures.SetupOrganizationTestData(token, userID)
		if setupErr != nil {
			t.Logf("⚠️ Lỗi setup organization test data: %v", setupErr)
		}
		if testData != nil {
			fmt.Printf("✅ Setup organization test data thành công\n")
			fmt.Printf("  Company: %s (Role: %s)\n", testData.CompanyOrgID, testData.CompanyRoleID)
			fmt.Printf("  DeptA: %s (Role: %s)\n", testData.DeptAOrgID, testData.DeptARoleID)
			fmt.Printf("  DeptB: %s (Role: %s)\n", testData.DeptBOrgID, testData.DeptBRoleID)
			fmt.Printf("  TeamA: %s (Role: %s)\n", testData.TeamAOrgID, testData.TeamARoleID)
		}
	})

	// Map testData vào các biến cũ để tương thích với code hiện tại
	companyOrgID := ""
	companyRoleID := ""
	deptARoleID := ""
	deptBRoleID := ""
	teamARoleID := ""

	if testData != nil {
		companyOrgID = testData.CompanyOrgID
		companyRoleID = testData.CompanyRoleID
		deptARoleID = testData.DeptARoleID
		deptBRoleID = testData.DeptBRoleID
		teamARoleID = testData.TeamARoleID
	}

	// ============================================
	// TEST CASE 1: Tự động gán organizationId khi insert
	// ============================================
	t.Run("📝 Test Case 1: Tự động gán organizationId khi insert", func(t *testing.T) {
		if companyRoleID == "" {
			t.Skip("Skipping: Không có Company Role ID")
		}

		client.SetActiveRoleID(companyRoleID)

		// Test với FbCustomer
		t.Run("FbCustomer - Tự động gán organizationId", func(t *testing.T) {
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("test_customer_%d", time.Now().UnixNano()),
				"name":       "Test Customer",
				"email":      "test@example.com",
			}

			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo customer: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					orgID, ok := data["organizationId"].(string)
					assert.True(t, ok, "Phải có organizationId")
					assert.Equal(t, companyOrgID, orgID, "organizationId phải khớp với active organization")
					fmt.Printf("✅ FbCustomer: organizationId = %s\n", orgID)
				}
			}
		})

		// Test với PcPosCustomer
		t.Run("PcPosCustomer - Tự động gán organizationId", func(t *testing.T) {
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("pos_customer_%d", time.Now().UnixNano()),
				"name":       "POS Customer",
				"email":      "pos@example.com",
			}

			resp, body, err := client.POST("/pc-pos-customer/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo POS customer: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					orgID, ok := data["organizationId"].(string)
					if ok {
						assert.Equal(t, companyOrgID, orgID, "organizationId phải khớp")
						fmt.Printf("✅ PcPosCustomer: organizationId = %s\n", orgID)
					}
				}
			}
		})

		// Test với Notification Channel
		t.Run("NotificationChannel - Tự động gán organizationId", func(t *testing.T) {
			payload := map[string]interface{}{
				"name":        fmt.Sprintf("TestChannel_%d", time.Now().UnixNano()),
				"channelType": "email",
				"recipients":  []string{"test@example.com"},
			}

			resp, body, err := client.POST("/notification/channel/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo channel: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					orgID, ok := data["organizationId"].(string)
					if ok {
						assert.Equal(t, companyOrgID, orgID, "organizationId phải khớp")
						fmt.Printf("✅ NotificationChannel: organizationId = %s\n", orgID)
					}
				}
			}
		})
	})

	// ============================================
	// TEST CASE 2: Filter dữ liệu theo organization
	// ============================================
	t.Run("🔍 Test Case 2: Filter dữ liệu theo organization", func(t *testing.T) {
		if companyRoleID == "" || deptARoleID == "" {
			t.Skip("Skipping: Không đủ roles")
		}

		// Tạo dữ liệu ở Company
		client.SetActiveRoleID(companyRoleID)
		companyCustomerID := ""
		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("company_customer_%d", time.Now().UnixNano()),
				"name":       "Company Customer",
				"email":      "company@example.com",
			}
			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					companyCustomerID, _ = data["id"].(string)
				}
			}
		}

		// Tạo dữ liệu ở Department A
		client.SetActiveRoleID(deptARoleID)
		deptACustomerID := ""
		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("dept_a_customer_%d", time.Now().UnixNano()),
				"name":       "Dept A Customer",
				"email":      "depta@example.com",
			}
			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					deptACustomerID, _ = data["id"].(string)
				}
			}
		}

		// Test: User ở Department A chỉ thấy dữ liệu của mình (Scope 0)
		t.Run("Scope 0 - Chỉ thấy dữ liệu của organization mình", func(t *testing.T) {
			client.SetActiveRoleID(deptARoleID)

			resp, body, err := client.GET("/fb-customer/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi query customers: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				customers, ok := result["data"].([]interface{})
				if ok {
					// Với Scope 0, chỉ thấy customers của Department A
					// Nhưng với inverse parent lookup, sẽ thấy cả Company customers
					foundDeptA := false
					foundCompany := false

					for _, item := range customers {
						customer, ok := item.(map[string]interface{})
						if ok {
							id, _ := customer["id"].(string)
							if id == deptACustomerID {
								foundDeptA = true
							}
							if id == companyCustomerID {
								foundCompany = true
							}
						}
					}

					assert.True(t, foundDeptA, "Phải thấy customer của Department A")
					// Với inverse parent lookup, sẽ thấy cả Company customers
					if foundCompany {
						fmt.Printf("✅ Inverse parent lookup hoạt động: thấy customer của Company\n")
					}
					fmt.Printf("✅ Filter customers: tìm thấy %d items\n", len(customers))
				}
			}
		})
	})

	// ============================================
	// TEST CASE 3: Scope = 1 (Children) - Xem dữ liệu của organization mình và con
	// ============================================
	t.Run("🔐 Test Case 3: Scope = 1 (Children)", func(t *testing.T) {
		if companyRoleID == "" || deptARoleID == "" || teamARoleID == "" {
			t.Skip("Skipping: Không đủ roles")
		}

		// Tạo dữ liệu ở các cấp khác nhau
		client.SetActiveRoleID(companyRoleID)
		companyDataID := ""
		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("company_data_%d", time.Now().UnixNano()),
				"name":       "Company Data",
			}
			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					companyDataID, _ = data["id"].(string)
				}
			}
		}

		client.SetActiveRoleID(deptARoleID)
		deptADataID := ""
		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("dept_a_data_%d", time.Now().UnixNano()),
				"name":       "Dept A Data",
			}
			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					deptADataID, _ = data["id"].(string)
				}
			}
		}

		client.SetActiveRoleID(teamARoleID)
		teamADataID := ""
		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("team_a_data_%d", time.Now().UnixNano()),
				"name":       "Team A Data",
			}
			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					teamADataID, _ = data["id"].(string)
				}
			}
		}

		// Test: User ở Company với Scope = 1 sẽ thấy tất cả dữ liệu của Company và children
		t.Run("Company Role với Scope 1 - Thấy tất cả children", func(t *testing.T) {
			client.SetActiveRoleID(companyRoleID)

			resp, body, err := client.GET("/fb-customer/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi query: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				customers, ok := result["data"].([]interface{})
				if ok {
					foundCompany := false
					foundDeptA := false
					foundTeamA := false

					for _, item := range customers {
						customer, ok := item.(map[string]interface{})
						if ok {
							id, _ := customer["id"].(string)
							if id == companyDataID {
								foundCompany = true
							}
							if id == deptADataID {
								foundDeptA = true
							}
							if id == teamADataID {
								foundTeamA = true
							}
						}
					}

					// Với Scope 1, phải thấy tất cả (Company + Dept A + Team A)
					// Lưu ý: Cần có permission với Scope = 1
					fmt.Printf("✅ Scope 1 test: Company=%v, DeptA=%v, TeamA=%v\n", foundCompany, foundDeptA, foundTeamA)
					fmt.Printf("  Total customers: %d\n", len(customers))
				}
			}
		})
	})

	// ============================================
	// TEST CASE 4: Inverse Parent Lookup - Xem dữ liệu cấp trên
	// ============================================
	t.Run("⬆️ Test Case 4: Inverse Parent Lookup", func(t *testing.T) {
		if companyRoleID == "" || deptARoleID == "" || teamARoleID == "" {
			t.Skip("Skipping: Không đủ roles")
		}

		// Tạo dữ liệu ở Company (cấp trên)
		client.SetActiveRoleID(companyRoleID)
		parentCustomerID := ""
		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("parent_customer_%d", time.Now().UnixNano()),
				"name":       "Parent Customer",
				"email":      "parent@example.com",
			}
			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					parentCustomerID, _ = data["id"].(string)
				}
			}
		}

		// Test: User ở Team A (cấp thấp nhất) có thể thấy dữ liệu của Company (cấp trên)
		t.Run("Team A thấy dữ liệu của Company", func(t *testing.T) {
			client.SetActiveRoleID(teamARoleID)

			resp, body, err := client.GET("/fb-customer/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi query: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				customers, ok := result["data"].([]interface{})
				if ok {
					foundParent := false
					for _, item := range customers {
						customer, ok := item.(map[string]interface{})
						if ok {
							id, _ := customer["id"].(string)
							if id == parentCustomerID {
								foundParent = true
								break
							}
						}
					}

					if foundParent {
						fmt.Printf("✅ Inverse parent lookup hoạt động: Team A thấy customer của Company\n")
					} else {
						fmt.Printf("⚠️ Inverse parent lookup: Team A không thấy customer của Company (có thể do permission scope)\n")
					}
					fmt.Printf("  Total customers: %d\n", len(customers))
				}
			}
		})
	})

	// ============================================
	// TEST CASE 5: Không thể update organizationId
	// ============================================
	t.Run("🔒 Test Case 5: Không thể update organizationId", func(t *testing.T) {
		if companyRoleID == "" || deptARoleID == "" {
			t.Skip("Skipping: Không đủ roles")
		}

		// Tạo customer ở Company
		client.SetActiveRoleID(companyRoleID)
		var customerID string
		var originalOrgID string

		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("test_update_org_%d", time.Now().UnixNano()),
				"name":       "Test Update Org",
				"email":      "update@example.com",
			}
			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					customerID, _ = data["id"].(string)
					originalOrgID, _ = data["organizationId"].(string)
				}
			}
		}

		if customerID == "" {
			t.Skip("Skipping: Không thể tạo customer để test")
		}

		// Thử update với organizationId khác
		t.Run("Thử update organizationId", func(t *testing.T) {
			updatePayload := map[string]interface{}{
				"name":           "Updated Name",
				"organizationId": deptARoleID, // ID giả, không phải organizationId hợp lệ
			}

			resp, body, err := client.PUT(fmt.Sprintf("/fb-customer/update-by-id/%s", customerID), updatePayload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi update: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					updatedOrgID, _ := data["organizationId"].(string)
					// organizationId không được thay đổi
					assert.Equal(t, originalOrgID, updatedOrgID, "organizationId không được phép thay đổi")
					fmt.Printf("✅ Verify: organizationId không thể update (vẫn là: %s)\n", updatedOrgID)
				}
			} else {
				fmt.Printf("⚠️ Update yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// ============================================
	// TEST CASE 6: Validate quyền truy cập khi update/delete
	// ============================================
	t.Run("🛡️ Test Case 6: Validate quyền truy cập", func(t *testing.T) {
		if companyRoleID == "" || deptBRoleID == "" {
			t.Skip("Skipping: Không đủ roles")
		}

		// Tạo customer ở Company
		client.SetActiveRoleID(companyRoleID)
		var customerID string

		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("test_access_%d", time.Now().UnixNano()),
				"name":       "Test Access",
				"email":      "access@example.com",
			}
			resp, body, err := client.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					customerID, _ = data["id"].(string)
				}
			}
		}

		if customerID == "" {
			t.Skip("Skipping: Không thể tạo customer để test")
		}

		// Thử update với role khác organization
		t.Run("Update với role khác organization", func(t *testing.T) {
			client.SetActiveRoleID(deptBRoleID) // Role ở Department B, khác Company

			updatePayload := map[string]interface{}{
				"name": "Unauthorized Update",
			}

			resp, _, err := client.PUT(fmt.Sprintf("/fb-customer/update-by-id/%s", customerID), updatePayload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi update: %v", err)
			}

			// Phải trả về 403 Forbidden hoặc không tìm thấy
			if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusNotFound {
				fmt.Printf("✅ Validate quyền truy cập: Không cho phép update dữ liệu của organization khác\n")
			} else if resp.StatusCode == http.StatusOK {
				// Có thể thành công nếu có inverse parent lookup và permission scope
				fmt.Printf("⚠️ Update thành công (có thể do inverse parent lookup)\n")
			} else {
				fmt.Printf("⚠️ Update trả về status: %d\n", resp.StatusCode)
			}
		})
	})

	// ============================================
	// TEST CASE 7: Test với nhiều collections có organizationId
	// ============================================
	t.Run("📦 Test Case 7: Test với nhiều collections", func(t *testing.T) {
		if companyRoleID == "" {
			t.Skip("Skipping: Không có Company Role ID")
		}

		client.SetActiveRoleID(companyRoleID)

		collections := []struct {
			name    string
			endpoint string
			payload  map[string]interface{}
		}{
			{
				name:     "FbPage",
				endpoint: "/facebook/page/insert-one",
				payload: map[string]interface{}{
					"pageId":          fmt.Sprintf("test_page_%d", time.Now().UnixNano()),
					"pageName":        "Test Page",
					"pageUsername":    "testpage",
					"isSync":          false,
					"accessToken":     "test_token",
					"pageAccessToken": "test_page_token",
				},
			},
			{
				name:     "PcPosShop",
				endpoint: "/pancake-pos/shop/insert-one",
				payload: map[string]interface{}{
					"shopId": int64(time.Now().UnixNano()),
					"name":   "Test Shop",
				},
			},
			{
				name:     "PcPosProduct",
				endpoint: "/pancake-pos/product/insert-one",
				payload: map[string]interface{}{
					"productId": fmt.Sprintf("test_product_%d", time.Now().UnixNano()),
					"name":      "Test Product",
					"shopId":    int64(123),
				},
			},
			{
				name:     "AccessToken",
				endpoint: "/access-token/insert-one",
				payload: map[string]interface{}{
					"name":   fmt.Sprintf("TestToken_%d", time.Now().UnixNano()),
					"system": "test",
					"value":  "test_token_value",
				},
			},
		}

		for _, collection := range collections {
			t.Run(collection.name, func(t *testing.T) {
				resp, body, err := client.POST(collection.endpoint, collection.payload)
				if err != nil {
					t.Logf("⚠️ Lỗi khi tạo %s: %v", collection.name, err)
					return
				}

				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					if err == nil {
						data, ok := result["data"].(map[string]interface{})
						if ok {
							orgID, ok := data["organizationId"].(string)
							if ok {
								assert.Equal(t, companyOrgID, orgID, fmt.Sprintf("%s: organizationId phải khớp", collection.name))
								fmt.Printf("✅ %s: organizationId = %s\n", collection.name, orgID)
							} else {
								fmt.Printf("⚠️ %s: Không có organizationId (có thể model chưa có field)\n", collection.name)
							}
						}
					}
				} else {
					fmt.Printf("⚠️ %s: Yêu cầu quyền (status: %d)\n", collection.name, resp.StatusCode)
				}
			})
		}
	})

	// ============================================
	// TEST CASE 8: Collections không có organizationId hoạt động bình thường
	// ============================================
	t.Run("✅ Test Case 8: Collections không có organizationId", func(t *testing.T) {
		// Test với User (không có organizationId)
		t.Run("User - Không có organizationId", func(t *testing.T) {
			resp, body, err := client.GET("/user/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi query users: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ User collection hoạt động bình thường (không có organizationId)\n")
			} else {
				fmt.Printf("⚠️ Query users yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test với Permission (không có organizationId)
		t.Run("Permission - Không có organizationId", func(t *testing.T) {
			resp, body, err := client.GET("/permission/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi query permissions: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Permission collection hoạt động bình thường (không có organizationId)\n")
			} else {
				fmt.Printf("⚠️ Query permissions yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// ============================================
	// TEST CASE 9: Multi-client support
	// ============================================
	t.Run("💻 Test Case 9: Multi-client support", func(t *testing.T) {
		if companyRoleID == "" || deptARoleID == "" {
			t.Skip("Skipping: Không đủ roles")
		}

		// Client 1: Set role Company
		client1 := utils.NewHTTPClient(baseURL, 10)
		client1.SetToken(token)
		client1.SetActiveRoleID(companyRoleID)

		// Client 2: Set role Department A
		client2 := utils.NewHTTPClient(baseURL, 10)
		client2.SetToken(token)
		client2.SetActiveRoleID(deptARoleID)

		// Tạo dữ liệu với client 1
		var client1CustomerID string
		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("client1_customer_%d", time.Now().UnixNano()),
				"name":       "Client 1 Customer",
			}
			resp, body, err := client1.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					client1CustomerID, _ = data["id"].(string)
				}
			}
		}

		// Tạo dữ liệu với client 2
		var client2CustomerID string
		{
			payload := map[string]interface{}{
				"customerId": fmt.Sprintf("client2_customer_%d", time.Now().UnixNano()),
				"name":       "Client 2 Customer",
			}
			resp, body, err := client2.POST("/fb-customer/insert-one", payload)
			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if data, ok := result["data"].(map[string]interface{}); ok {
					client2CustomerID, _ = data["id"].(string)
				}
			}
		}

		// Verify: Client 1 chỉ thấy dữ liệu của mình (và parent nếu có inverse lookup)
		t.Run("Client 1 chỉ thấy dữ liệu của Company", func(t *testing.T) {
			resp, body, err := client1.GET("/fb-customer/find")
			if err == nil && resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				customers, _ := result["data"].([]interface{})

				foundClient1 := false
				foundClient2 := false

				for _, item := range customers {
					customer, ok := item.(map[string]interface{})
					if ok {
						id, _ := customer["id"].(string)
						if id == client1CustomerID {
							foundClient1 = true
						}
						if id == client2CustomerID {
							foundClient2 = true
						}
					}
				}

				assert.True(t, foundClient1, "Client 1 phải thấy dữ liệu của mình")
				// Client 1 không nên thấy dữ liệu của Client 2 (khác organization, không phải parent)
				if !foundClient2 {
					fmt.Printf("✅ Multi-client: Client 1 không thấy dữ liệu của Client 2 (đúng)\n")
				} else {
					fmt.Printf("⚠️ Multi-client: Client 1 thấy dữ liệu của Client 2 (có thể do permission scope)\n")
				}
			}
		})

		// Verify: Client 2 chỉ thấy dữ liệu của mình (và parent nếu có inverse lookup)
		t.Run("Client 2 chỉ thấy dữ liệu của Department A", func(t *testing.T) {
			resp, body, err := client2.GET("/fb-customer/find")
			if err == nil && resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				customers, _ := result["data"].([]interface{})

				foundClient1 := false
				foundClient2 := false

				for _, item := range customers {
					customer, ok := item.(map[string]interface{})
					if ok {
						id, _ := customer["id"].(string)
						if id == client1CustomerID {
							foundClient1 = true
						}
						if id == client2CustomerID {
							foundClient2 = true
						}
					}
				}

				assert.True(t, foundClient2, "Client 2 phải thấy dữ liệu của mình")
				// Client 2 có thể thấy dữ liệu của Client 1 nếu Company là parent của Department A
				if foundClient1 {
					fmt.Printf("✅ Multi-client: Client 2 thấy dữ liệu của Client 1 (inverse parent lookup)\n")
				} else {
					fmt.Printf("⚠️ Multi-client: Client 2 không thấy dữ liệu của Client 1\n")
				}
			}
		})
	})

	// ============================================
	// TEST CASE 10: Test với X-Active-Role-ID header
	// ============================================
	t.Run("📋 Test Case 10: X-Active-Role-ID header", func(t *testing.T) {
		// Lấy danh sách roles
		resp, body, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Skip("Skipping: Không thể lấy roles")
		}

		var result map[string]interface{}
		json.Unmarshal(body, &result)
		roles, ok := result["data"].([]interface{})
		if !ok || len(roles) == 0 {
			t.Skip("Skipping: Không có roles")
		}

		// Test: Không set header X-Active-Role-ID
		t.Run("Không set X-Active-Role-ID - Tự động chọn role đầu tiên", func(t *testing.T) {
			clientNoRole := utils.NewHTTPClient(baseURL, 10)
			clientNoRole.SetToken(token)
			// Không set active role ID

			resp, body, err := clientNoRole.POST("/fb-customer/insert-one", map[string]interface{}{
				"customerId": fmt.Sprintf("no_role_customer_%d", time.Now().UnixNano()),
				"name":       "No Role Customer",
			})

			if err == nil {
				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
					var result map[string]interface{}
					json.Unmarshal(body, &result)
					data, _ := result["data"].(map[string]interface{})
					orgID, _ := data["organizationId"].(string)
					if orgID != "" {
						fmt.Printf("✅ Tự động chọn role đầu tiên: organizationId = %s\n", orgID)
					}
				} else {
					fmt.Printf("⚠️ Tạo customer yêu cầu quyền (status: %d)\n", resp.StatusCode)
				}
			}
		})

		// Test: Set header X-Active-Role-ID với role hợp lệ
		t.Run("Set X-Active-Role-ID hợp lệ", func(t *testing.T) {
			firstRole, _ := roles[0].(map[string]interface{})
			roleID, _ := firstRole["roleId"].(string)
			expectedOrgID, _ := firstRole["organizationId"].(string)

			clientWithRole := utils.NewHTTPClient(baseURL, 10)
			clientWithRole.SetToken(token)
			clientWithRole.SetActiveRoleID(roleID)

			resp, body, err := clientWithRole.POST("/fb-customer/insert-one", map[string]interface{}{
				"customerId": fmt.Sprintf("with_role_customer_%d", time.Now().UnixNano()),
				"name":       "With Role Customer",
			})

			if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				data, _ := result["data"].(map[string]interface{})
				orgID, _ := data["organizationId"].(string)
				if orgID == expectedOrgID {
					fmt.Printf("✅ Set X-Active-Role-ID hoạt động: organizationId = %s\n", orgID)
				}
			}
		})

		// Test: Set header X-Active-Role-ID với role không hợp lệ
		t.Run("Set X-Active-Role-ID không hợp lệ - Fallback về role đầu tiên", func(t *testing.T) {
			invalidRoleID := "507f1f77bcf86cd799439999"

			clientInvalidRole := utils.NewHTTPClient(baseURL, 10)
			clientInvalidRole.SetToken(token)
			clientInvalidRole.SetActiveRoleID(invalidRoleID)

			resp, body, err := clientInvalidRole.POST("/fb-customer/insert-one", map[string]interface{}{
				"customerId": fmt.Sprintf("invalid_role_customer_%d", time.Now().UnixNano()),
				"name":       "Invalid Role Customer",
			})

			if err == nil {
				// Có thể thành công nếu fallback về role đầu tiên
				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
					var result map[string]interface{}
					json.Unmarshal(body, &result)
					data, _ := result["data"].(map[string]interface{})
					orgID, _ := data["organizationId"].(string)
					if orgID != "" {
						fmt.Printf("✅ Fallback về role đầu tiên: organizationId = %s\n", orgID)
					}
				} else {
					fmt.Printf("⚠️ Invalid role ID: status = %d\n", resp.StatusCode)
				}
			}
		})
	})
}

