package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GetTestFirebaseIDToken lấy Firebase ID token từ environment variable
// Lưu ý: Test cần có Firebase ID token hợp lệ từ Firebase test project
// Có thể set qua environment variable: TEST_FIREBASE_ID_TOKEN
func GetTestFirebaseIDToken() string {
	return os.Getenv("TEST_FIREBASE_ID_TOKEN")
}

// TestFixtures chứa các helper để setup test data
type TestFixtures struct {
	client  *HTTPClient
	baseURL string
}

// NewTestFixtures tạo mới TestFixtures
func NewTestFixtures(baseURL string) *TestFixtures {
	return &TestFixtures{
		client:  NewHTTPClient(baseURL, 10),
		baseURL: baseURL,
	}
}

// CreateTestUser tạo user test và trả về email, firebaseUID, token
// Lưu ý: Cần cung cấp Firebase ID token hợp lệ từ Firebase test project
// Firebase ID token có thể lấy từ environment variable TEST_FIREBASE_ID_TOKEN
// hoặc tạo bằng Firebase Admin SDK trong test setup
func (tf *TestFixtures) CreateTestUser(firebaseIDToken string) (email, firebaseUID, token string, err error) {
	if firebaseIDToken == "" {
		return "", "", "", fmt.Errorf("firebase ID token là bắt buộc cho test")
	}

	// Đăng nhập bằng Firebase để tạo/lấy user
	loginPayload := map[string]interface{}{
		"idToken": firebaseIDToken,
		"hwid":    "test_device_123",
	}

	resp, body, err := tf.client.POST("/auth/login/firebase", loginPayload)
	if err != nil {
		return "", "", "", fmt.Errorf("lỗi đăng nhập Firebase: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("đăng nhập Firebase thất bại: %d - %s", resp.StatusCode, string(body))
	}

	// Parse token từ response
	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", "", "", fmt.Errorf("lỗi parse response: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", "", "", fmt.Errorf("không có data trong response")
	}

	token, ok = data["token"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("không có token trong response")
	}

	// Lấy email và firebaseUID từ response
	email, _ = data["email"].(string)
	firebaseUID, _ = data["firebaseUid"].(string)

	return email, firebaseUID, token, nil
}

// CreateTestUserDirect tạo user trực tiếp trong database (bypass Firebase) - CHỈ DÙNG CHO TEST
// Tạo user với email và FirebaseUID giả để test nhanh hơn
// Lưu ý: User này sẽ không thể login qua Firebase, chỉ dùng để test database operations
// ⚠️ KHÔNG KHUYẾN NGHỊ: Hệ thống yêu cầu Firebase authentication, không thể bypass
func (tf *TestFixtures) CreateTestUserDirect(email, name string) (userID, token string, err error) {
	// Hệ thống yêu cầu Firebase authentication, không thể tạo user trực tiếp
	// Sử dụng CreateTestUser() với Firebase ID token thay thế
	return "", "", fmt.Errorf("không thể tạo user trực tiếp - cần Firebase authentication. Sử dụng CreateTestUser() với Firebase ID token")
}

// GetRootOrganizationID lấy Organization Root ID
func (tf *TestFixtures) GetRootOrganizationID(token string) (string, error) {
	tf.client.SetToken(token)

	// Tìm Organization System (Code: SYSTEM)
	// URL encode filter parameter
	filter := `{"code":"SYSTEM"}`
	encodedFilter := url.QueryEscape(filter)
	resp, body, err := tf.client.GET(fmt.Sprintf("/organization/find?filter=%s", encodedFilter))
	if err != nil {
		return "", fmt.Errorf("lỗi lấy root organization: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("lấy root organization thất bại: %d - %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("lỗi parse response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok || len(data) == 0 {
		return "", fmt.Errorf("không tìm thấy root organization")
	}

	firstOrg, ok := data[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("không parse được organization data")
	}

	id, ok := firstOrg["id"].(string)
	if !ok {
		return "", fmt.Errorf("không có id trong organization response")
	}

	return id, nil
}

// CreateTestRole tạo role test và trả về role ID
// Role phải có organizationId (bắt buộc)
func (tf *TestFixtures) CreateTestRole(token, name, describe, organizationID string) (string, error) {
	tf.client.SetToken(token)

	// Nếu không có organizationID, lấy Root Organization
	if organizationID == "" {
		rootOrgID, err := tf.GetRootOrganizationID(token)
		if err != nil {
			return "", fmt.Errorf("lỗi lấy root organization: %v", err)
		}
		organizationID = rootOrgID
	}

	payload := map[string]interface{}{
		"name":                name,
		"describe":            describe,
		"ownerOrganizationId": organizationID, // BẮT BUỘC - Phân quyền dữ liệu
	}

	resp, body, err := tf.client.POST("/role/insert-one", payload)
	if err != nil {
		return "", fmt.Errorf("lỗi tạo role: %v", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("tạo role thất bại: %d - %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("lỗi parse response: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("không có data trong response")
	}

	id, ok := data["id"].(string)
	if !ok {
		return "", fmt.Errorf("không có id trong response")
	}

	return id, nil
}

// CreateTestPermission tạo permission test và trả về permission ID
func (tf *TestFixtures) CreateTestPermission(token, name, describe, category, group string) (string, error) {
	tf.client.SetToken(token)

	payload := map[string]interface{}{
		"name":     name,
		"describe": describe,
		"category": category,
		"group":    group,
	}

	resp, body, err := tf.client.POST("/permission/insert-one", payload)
	if err != nil {
		return "", fmt.Errorf("lỗi tạo permission: %v", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("tạo permission thất bại: %d - %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("lỗi parse response: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("không có data trong response")
	}

	id, ok := data["id"].(string)
	if !ok {
		return "", fmt.Errorf("không có id trong response")
	}

	return id, nil
}

// CreateAdminUser tạo user và set làm administrator với full quyền
// Trả về userID để có thể dùng cho các test khác
// Lưu ý: Cần cung cấp Firebase ID token hợp lệ
func (tf *TestFixtures) CreateAdminUser(firebaseIDToken string) (email, firebaseUID, token, userID string, err error) {
	// Tạo user thường trước
	email, firebaseUID, token, err = tf.CreateTestUser(firebaseIDToken)
	if err != nil {
		return "", "", "", "", fmt.Errorf("lỗi tạo user: %v", err)
	}

	// Lấy user ID từ profile
	tf.client.SetToken(token)
	resp, body, err := tf.client.GET("/auth/profile")
	if err != nil {
		return "", "", "", "", fmt.Errorf("lỗi lấy profile: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", fmt.Errorf("lấy profile thất bại: %d - %s", resp.StatusCode, string(body))
	}

	var profileResult map[string]interface{}
	if err = json.Unmarshal(body, &profileResult); err != nil {
		return "", "", "", "", fmt.Errorf("lỗi parse profile: %v", err)
	}

	data, ok := profileResult["data"].(map[string]interface{})
	if !ok {
		return "", "", "", "", fmt.Errorf("không có data trong profile response")
	}

	userID, ok = data["id"].(string)
	if !ok {
		return "", "", "", "", fmt.Errorf("không có id trong profile response")
	}

	// Set administrator - API này yêu cầu permission "Init.SetAdmin"
	// Thử với token hiện tại (có thể thành công nếu là lần đầu init hoặc đã có quyền)
	resp, body, err = tf.client.POST(fmt.Sprintf("/init/set-administrator/%s", userID), nil)
	if err != nil {
		return "", "", "", "", fmt.Errorf("lỗi set administrator: %v", err)
	}

	// Nếu thành công, đăng nhập lại bằng Firebase để refresh token với permissions mới
	if resp.StatusCode == http.StatusOK {
		loginPayload := map[string]interface{}{
			"idToken": firebaseIDToken,
			"hwid":    "test_device_123",
		}

		// Tạo client mới không có token để login
		loginClient := NewHTTPClient(tf.baseURL, 10)
		resp, body, err = loginClient.POST("/auth/login/firebase", loginPayload)
		if err != nil {
			return "", "", "", "", fmt.Errorf("lỗi đăng nhập lại: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			return "", "", "", "", fmt.Errorf("đăng nhập lại thất bại: %d - %s", resp.StatusCode, string(body))
		}

		var loginResult map[string]interface{}
		if err = json.Unmarshal(body, &loginResult); err != nil {
			return "", "", "", "", fmt.Errorf("lỗi parse login response: %v", err)
		}

		loginData, ok := loginResult["data"].(map[string]interface{})
		if !ok {
			return "", "", "", "", fmt.Errorf("không có data trong login response")
		}

		newToken, ok := loginData["token"].(string)
		if !ok {
			return "", "", "", "", fmt.Errorf("không có token trong login response")
		}

		return email, firebaseUID, newToken, userID, nil
	}

	// Nếu fail (403 - không có quyền), vẫn trả về token hiện tại
	// Test sẽ phải xử lý trường hợp này
	return email, firebaseUID, token, userID, nil
}

// InitData khởi tạo tất cả dữ liệu mặc định của hệ thống
// Bao gồm: Root Organization, Permissions, Roles
// API này chỉ hoạt động khi chưa có admin trong hệ thống
func (tf *TestFixtures) InitData() error {
	// Gọi API init/all để khởi tạo tất cả
	resp, body, err := tf.client.POST("/init/all", nil)
	if err != nil {
		return fmt.Errorf("lỗi gọi init API: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		// Nếu đã có admin, API sẽ không được đăng ký (404) hoặc đã init rồi
		if resp.StatusCode == http.StatusNotFound {
			// Có thể đã có admin, thử kiểm tra status
			return tf.checkInitStatus()
		}
		return fmt.Errorf("init data thất bại: %d - %s", resp.StatusCode, string(body))
	}

	// Parse response để kiểm tra kết quả
	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		// Không parse được cũng không sao, có thể đã init thành công
		return nil
	}

	// Kiểm tra từng phần init
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil // Không có data, có thể đã init rồi
	}

	// Kiểm tra status của từng phần
	if orgStatus, ok := data["organization"].(map[string]interface{}); ok {
		if status, ok := orgStatus["status"].(string); ok && status != "success" {
			return fmt.Errorf("init organization thất bại: %v", orgStatus)
		}
	}

	if permStatus, ok := data["permissions"].(map[string]interface{}); ok {
		if status, ok := permStatus["status"].(string); ok && status != "success" {
			return fmt.Errorf("init permissions thất bại: %v", permStatus)
		}
	}

	if roleStatus, ok := data["roles"].(map[string]interface{}); ok {
		if status, ok := roleStatus["status"].(string); ok && status != "success" {
			return fmt.Errorf("init roles thất bại: %v", roleStatus)
		}
	}

	return nil
}

// checkInitStatus kiểm tra trạng thái init của hệ thống
func (tf *TestFixtures) checkInitStatus() error {
	resp, body, err := tf.client.GET("/init/status")
	if err != nil {
		// Nếu không có endpoint (404), có thể đã có admin rồi
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil // Có thể đã init rồi
		}
		return fmt.Errorf("lỗi kiểm tra init status: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		// Nếu không OK, có thể đã có admin rồi
		return nil
	}

	// Parse response
	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil // Không parse được, có thể đã init rồi
	}

	// Kiểm tra data
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil
	}

	// Kiểm tra các thành phần đã init chưa
	hasOrg, _ := data["hasOrganization"].(bool)
	hasPerm, _ := data["hasPermissions"].(bool)
	hasRole, _ := data["hasRoles"].(bool)

	if !hasOrg || !hasPerm || !hasRole {
		return fmt.Errorf("chưa init đầy đủ: org=%v, perm=%v, role=%v", hasOrg, hasPerm, hasRole)
	}

	return nil
}

// OrganizationTestData chứa dữ liệu test cho organization ownership
type OrganizationTestData struct {
	CompanyOrgID  string
	DeptAOrgID    string
	DeptBOrgID    string
	TeamAOrgID    string
	CompanyRoleID string
	DeptARoleID   string
	DeptBRoleID   string
	TeamARoleID   string
}

// SetupOrganizationTestData tạo đầy đủ dữ liệu test cho organization ownership
// Bao gồm: organization hierarchy, roles, permissions với scope, và gán roles cho user
// Lưu ý: User cần có quyền Organization.Insert và Role.Insert để tạo dữ liệu
// Nếu user không có quyền, function sẽ thử set user làm admin trước
func (tf *TestFixtures) SetupOrganizationTestData(token, userID string) (*OrganizationTestData, error) {
	tf.client.SetToken(token)

	// Lấy Root Organization ID
	rootOrgID, err := tf.GetRootOrganizationID(token)
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy root organization: %v", err)
	}

	data := &OrganizationTestData{}

	fmt.Printf("🔧 Bắt đầu setup organization test data...\n")

	// Thử set user làm admin nếu chưa có quyền (chỉ khi chưa có admin trong hệ thống)
	// API /init/set-administrator chỉ hoạt động khi chưa có admin
	resp, _, _ := tf.client.POST(fmt.Sprintf("/init/set-administrator/%s", userID), nil)
	if resp != nil && resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Đã set user làm administrator để có quyền tạo organization/roles\n")
		// Refresh token để có permissions mới (nhưng không cần thiết vì token đã có trong context)
	}

	// 1. Tạo Company (cấp 2)
	companyPayload := map[string]interface{}{
		"name":     fmt.Sprintf("TestCompany_%d", time.Now().UnixNano()),
		"code":     fmt.Sprintf("COMP_%d", time.Now().UnixNano()),
		"type":     "company", // Company - phải là string
		"parentId": rootOrgID,
	}
	resp, body, err := tf.client.POST("/organization/insert-one", companyPayload)
	if err != nil {
		fmt.Printf("⚠️ Lỗi khi tạo Company: %v\n", err)
	} else if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			if dataMap, ok := result["data"].(map[string]interface{}); ok {
				data.CompanyOrgID, _ = dataMap["id"].(string)
				if data.CompanyOrgID != "" {
					fmt.Printf("✅ Tạo Company thành công: %s\n", data.CompanyOrgID)
				}
			}
		}
	} else {
		fmt.Printf("⚠️ Tạo Company thất bại (status: %d): %s\n", resp.StatusCode, string(body))
	}

	// 2. Tạo Department A (cấp 3)
	if data.CompanyOrgID != "" {
		deptAPayload := map[string]interface{}{
			"name":     fmt.Sprintf("DeptA_%d", time.Now().UnixNano()),
			"code":     fmt.Sprintf("DEPT_A_%d", time.Now().UnixNano()),
			"type":     "department", // Department - phải là string
			"parentId": data.CompanyOrgID,
		}
		resp, body, err := tf.client.POST("/organization/insert-one", deptAPayload)
		if err != nil {
			fmt.Printf("⚠️ Lỗi khi tạo Department A: %v\n", err)
		} else if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err == nil {
				if dataMap, ok := result["data"].(map[string]interface{}); ok {
					data.DeptAOrgID, _ = dataMap["id"].(string)
					if data.DeptAOrgID != "" {
						fmt.Printf("✅ Tạo Department A thành công: %s\n", data.DeptAOrgID)
					}
				}
			}
		} else {
			fmt.Printf("⚠️ Tạo Department A thất bại (status: %d): %s\n", resp.StatusCode, string(body))
		}
	}

	// 3. Tạo Department B (cấp 3)
	if data.CompanyOrgID != "" {
		deptBPayload := map[string]interface{}{
			"name":     fmt.Sprintf("DeptB_%d", time.Now().UnixNano()),
			"code":     fmt.Sprintf("DEPT_B_%d", time.Now().UnixNano()),
			"type":     "department", // Department - phải là string
			"parentId": data.CompanyOrgID,
		}
		resp, body, err := tf.client.POST("/organization/insert-one", deptBPayload)
		if err != nil {
			fmt.Printf("⚠️ Lỗi khi tạo Department B: %v\n", err)
		} else if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err == nil {
				if dataMap, ok := result["data"].(map[string]interface{}); ok {
					data.DeptBOrgID, _ = dataMap["id"].(string)
					if data.DeptBOrgID != "" {
						fmt.Printf("✅ Tạo Department B thành công: %s\n", data.DeptBOrgID)
					}
				}
			}
		} else {
			fmt.Printf("⚠️ Tạo Department B thất bại (status: %d): %s\n", resp.StatusCode, string(body))
		}
	}

	// 4. Tạo Team A (cấp 4) thuộc Department A
	if data.DeptAOrgID != "" {
		teamAPayload := map[string]interface{}{
			"name":     fmt.Sprintf("TeamA_%d", time.Now().UnixNano()),
			"code":     fmt.Sprintf("TEAM_A_%d", time.Now().UnixNano()),
			"type":     "team", // Team - phải là string
			"parentId": data.DeptAOrgID,
		}
		resp, body, err := tf.client.POST("/organization/insert-one", teamAPayload)
		if err != nil {
			fmt.Printf("⚠️ Lỗi khi tạo Team A: %v\n", err)
		} else if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err == nil {
				if dataMap, ok := result["data"].(map[string]interface{}); ok {
					data.TeamAOrgID, _ = dataMap["id"].(string)
					if data.TeamAOrgID != "" {
						fmt.Printf("✅ Tạo Team A thành công: %s\n", data.TeamAOrgID)
					}
				}
			}
		} else {
			fmt.Printf("⚠️ Tạo Team A thất bại (status: %d): %s\n", resp.StatusCode, string(body))
		}
	}

	// 5. Tạo roles cho từng organization
	if data.CompanyOrgID != "" {
		roleID, err := tf.CreateTestRole(token, fmt.Sprintf("CompanyRole_%d", time.Now().UnixNano()), "Company Role", data.CompanyOrgID)
		if err != nil {
			fmt.Printf("⚠️ Lỗi khi tạo Company Role: %v\n", err)
		} else if roleID != "" {
			data.CompanyRoleID = roleID
			fmt.Printf("✅ Tạo Company Role thành công: %s\n", roleID)
		}
	}
	if data.DeptAOrgID != "" {
		roleID, err := tf.CreateTestRole(token, fmt.Sprintf("DeptARole_%d", time.Now().UnixNano()), "Department A Role", data.DeptAOrgID)
		if err != nil {
			fmt.Printf("⚠️ Lỗi khi tạo Department A Role: %v\n", err)
		} else if roleID != "" {
			data.DeptARoleID = roleID
			fmt.Printf("✅ Tạo Department A Role thành công: %s\n", roleID)
		}
	}
	if data.DeptBOrgID != "" {
		roleID, err := tf.CreateTestRole(token, fmt.Sprintf("DeptBRole_%d", time.Now().UnixNano()), "Department B Role", data.DeptBOrgID)
		if err != nil {
			fmt.Printf("⚠️ Lỗi khi tạo Department B Role: %v\n", err)
		} else if roleID != "" {
			data.DeptBRoleID = roleID
			fmt.Printf("✅ Tạo Department B Role thành công: %s\n", roleID)
		}
	}
	if data.TeamAOrgID != "" {
		roleID, err := tf.CreateTestRole(token, fmt.Sprintf("TeamARole_%d", time.Now().UnixNano()), "Team A Role", data.TeamAOrgID)
		if err != nil {
			fmt.Printf("⚠️ Lỗi khi tạo Team A Role: %v\n", err)
		} else if roleID != "" {
			data.TeamARoleID = roleID
			fmt.Printf("✅ Tạo Team A Role thành công: %s\n", roleID)
		}
	}

	// 6. Lấy permissions cần thiết (FbCustomer.*, NotificationChannel.*, AccessToken.*)
	permissionNames := []string{
		"FbCustomer.Insert", "FbCustomer.Read", "FbCustomer.Update", "FbCustomer.Delete",
		"NotificationChannel.Insert", "NotificationChannel.Read", "NotificationChannel.Update", "NotificationChannel.Delete",
		"AccessToken.Insert", "AccessToken.Read", "AccessToken.Update", "AccessToken.Delete",
	}
	permissionIDs := make([]string, 0)

	for _, permName := range permissionNames {
		filter := fmt.Sprintf(`{"name":"%s"}`, permName)
		encodedFilter := url.QueryEscape(filter)
		resp, body, err := tf.client.GET(fmt.Sprintf("/permission/find?filter=%s", encodedFilter))
		if err == nil && resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			json.Unmarshal(body, &result)
			if dataList, ok := result["data"].([]interface{}); ok && len(dataList) > 0 {
				if perm, ok := dataList[0].(map[string]interface{}); ok {
					if id, ok := perm["id"].(string); ok {
						permissionIDs = append(permissionIDs, id)
					}
				}
			}
		}
	}

	// 7. Gán permissions cho roles với Scope 0 hoặc Scope 1
	// Company Role: Scope 1 (xem tất cả children)
	// Dept/Team Roles: Scope 0 (chỉ xem organization mình)
	if data.CompanyRoleID != "" && len(permissionIDs) > 0 {
		tf.assignPermissionsToRole(token, data.CompanyRoleID, permissionIDs, 1) // Scope 1
	}
	if data.DeptARoleID != "" && len(permissionIDs) > 0 {
		tf.assignPermissionsToRole(token, data.DeptARoleID, permissionIDs, 0) // Scope 0
	}
	if data.DeptBRoleID != "" && len(permissionIDs) > 0 {
		tf.assignPermissionsToRole(token, data.DeptBRoleID, permissionIDs, 0) // Scope 0
	}
	if data.TeamARoleID != "" && len(permissionIDs) > 0 {
		tf.assignPermissionsToRole(token, data.TeamARoleID, permissionIDs, 0) // Scope 0
	}

	// 8. Gán tất cả roles cho user
	roleIDs := make([]string, 0)
	if data.CompanyRoleID != "" {
		roleIDs = append(roleIDs, data.CompanyRoleID)
	}
	if data.DeptARoleID != "" {
		roleIDs = append(roleIDs, data.DeptARoleID)
	}
	if data.DeptBRoleID != "" {
		roleIDs = append(roleIDs, data.DeptBRoleID)
	}
	if data.TeamARoleID != "" {
		roleIDs = append(roleIDs, data.TeamARoleID)
	}

	if len(roleIDs) > 0 && userID != "" {
		updatePayload := map[string]interface{}{
			"userID":  userID,
			"roleIDs": roleIDs,
		}
		tf.client.PUT("/user-role/update-user", updatePayload)
	}

	return data, nil
}

// assignPermissionsToRole gán permissions cho role với scope cụ thể
func (tf *TestFixtures) assignPermissionsToRole(token, roleID string, permissionIDs []string, scope byte) {
	tf.client.SetToken(token)

	// Tạo danh sách permissions với scope (format đúng theo DTO)
	permissions := make([]map[string]interface{}, 0)
	for _, permID := range permissionIDs {
		permissions = append(permissions, map[string]interface{}{
			"permissionId": permID, // Chữ i thường, đúng format DTO
			"scope":        scope,
		})
	}

	// Sử dụng API update-role để gán permissions
	updatePayload := map[string]interface{}{
		"roleId":      roleID,
		"permissions": permissions,
	}

	// Thử dùng API update-role (PUT /role-permission/update-role)
	resp, _, _ := tf.client.PUT("/role-permission/update-role", updatePayload)
	if resp != nil && resp.StatusCode == http.StatusOK {
		return
	}

	// Nếu không có API update-role hoặc không có quyền, thử insert từng permission
	for _, permID := range permissionIDs {
		payload := map[string]interface{}{
			"roleId":       roleID,
			"permissionId": permID,
			"scope":        scope,
		}
		tf.client.POST("/role-permission/insert-one", payload)
	}
}
