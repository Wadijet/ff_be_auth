package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BaseURL = "http://localhost:8080/api/v1"

// CallAPI gọi API với method, path và body cho trước
func CallAPI(method string, path string, body interface{}, token string) (map[string]interface{}, error) {
	// Tạo request body
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
	}

	// Tạo request
	var req *http.Request
	if body != nil {
		req, err = http.NewRequest(method, fmt.Sprintf("%s%s", BaseURL, path), bytes.NewBuffer(reqBody))
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("%s%s", BaseURL, path), nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Gửi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Đọc response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse response JSON
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result, nil
}

// RegisterUser gọi API đăng ký user
func RegisterUser(userData map[string]interface{}) (map[string]interface{}, error) {
	return CallAPI("POST", "/users/register", userData, "")
}

// LoginUser gọi API đăng nhập
func LoginUser(loginData map[string]interface{}) (map[string]interface{}, error) {
	return CallAPI("POST", "/users/login", loginData, "")
}

// SetAdministrator gọi API set quyền Administrator
func SetAdministrator(userId string, token string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/init/setadmin/%s", userId)
	return CallAPI("POST", path, map[string]interface{}{}, token)
}

// LogoutUser gọi API đăng xuất
func LogoutUser(token string) (map[string]interface{}, error) {
	return CallAPI("POST", "/users/logout", TestLogoutData, token)
}
