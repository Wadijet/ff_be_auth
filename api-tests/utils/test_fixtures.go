package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
		"hwid":     "test_device_123",
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

// GetRootOrganizationID lấy Organization Root ID
func (tf *TestFixtures) GetRootOrganizationID(token string) (string, error) {
	tf.client.SetToken(token)

	// Tìm Organization Root (Code: ROOT_GROUP)
	// URL encode filter parameter
	filter := `{"code":"ROOT_GROUP"}`
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
		"name":           name,
		"describe":       describe,
		"organizationId": organizationID, // BẮT BUỘC
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
			"hwid":     "test_device_123",
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
