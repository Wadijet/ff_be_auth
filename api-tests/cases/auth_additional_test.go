package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestAuthAdditionalCases kiểm tra thêm các tình huống phụ của module auth
func TestAuthAdditionalCases(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Lấy Firebase ID token từ environment variable
	firebaseIDToken := os.Getenv("TEST_FIREBASE_ID_TOKEN")
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN environment variable not set")
	}

	// 1) Đăng nhập Firebase với token không hợp lệ -> 401
	wrongPayload := map[string]interface{}{
		"idToken": "invalid_firebase_token",
		"hwid":    deviceID,
	}
	_, wrongStatus := postJSON(baseURL+"/auth/login/firebase", wrongPayload, nil, t)
	assert.Equal(t, http.StatusUnauthorized, wrongStatus, "Firebase login với token không hợp lệ phải 401")

	// 2) Đăng nhập Firebase đúng để lấy token
	loginPayload := map[string]interface{}{
		"idToken": firebaseIDToken,
		"hwid":    deviceID,
	}
	loginBody, loginStatus := postJSON(baseURL+"/auth/login/firebase", loginPayload, nil, t)
	assert.Equalf(t, http.StatusOK, loginStatus, "Firebase login đúng phải 200. Body: %s", string(loginBody))

	var loginResult map[string]interface{}
	err := json.Unmarshal(loginBody, &loginResult)
	assert.NoError(t, err, "Parse login response")
	loginData, ok := loginResult["data"].(map[string]interface{})
	assert.True(t, ok, "Login response phải có data")
	token, ok := loginData["token"].(string)
	assert.True(t, ok, "Login response phải có token")

	// 4) Lấy roles với token (không yêu cầu permission cụ thể)
	rolesBody, rolesStatus := getWithAuth(baseURL+"/auth/roles", token, t)
	assert.Equalf(t, http.StatusOK, rolesStatus, "Lấy roles phải 200. Body: %s", string(rolesBody))
	var rolesResult map[string]interface{}
	err = json.Unmarshal(rolesBody, &rolesResult)
	assert.NoError(t, err, "Parse roles response")
	if val, exists := rolesResult["data"]; exists && val != nil {
		_, ok = val.([]interface{})
		assert.True(t, ok, "Roles response data phải là array (có thể rỗng)")
	}

	// 5) Gọi profile không token -> 401
	resp, err := http.Get(baseURL + "/auth/profile")
	assert.NoError(t, err, "Call profile không token không được lỗi network")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Profile không token phải 401")
}

// postJSON gửi POST kèm JSON và trả body + status
func postJSON(url string, payload map[string]interface{}, token *string, t *testing.T) ([]byte, int) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err, "POST %s không được lỗi network", url)
	defer resp.Body.Close()
	respBody, _ := readBody(resp)
	return respBody, resp.StatusCode
}

// getWithAuth gửi GET kèm Bearer token
func getWithAuth(url, token string, t *testing.T) ([]byte, int) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err, "GET %s không được lỗi network", url)
	defer resp.Body.Close()
	respBody, _ := readBody(resp)
	return respBody, resp.StatusCode
}
