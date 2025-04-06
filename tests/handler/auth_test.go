package handler_test

import (
	"bytes"
	"encoding/json"
	models "meta_commerce/app/models/mongodb"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	baseURL = "http://localhost:3000/api"
)

// waitForServer đợi server khởi động
func waitForServer(t *testing.T) {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		_, err := http.Get(baseURL + "/health")
		if err == nil {
			return
		}
		t.Logf("Waiting for server to start... (attempt %d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	t.Fatal("Server did not start in time")
}

func TestHandleLogin(t *testing.T) {
	waitForServer(t)

	tests := []struct {
		name       string
		payload    models.UserLoginInput
		wantStatus int
	}{
		{
			name: "Login thành công",
			payload: models.UserLoginInput{
				Email:    "test@example.com",
				Password: "password123",
				Hwid:     "test-device-id",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Login thất bại - thiếu email",
			payload: models.UserLoginInput{
				Password: "password123",
				Hwid:     "test-device-id",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Login thất bại - thiếu password",
			payload: models.UserLoginInput{
				Email: "test@example.com",
				Hwid:  "test-device-id",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Login thất bại - thiếu hwid",
			payload: models.UserLoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Chuyển payload thành JSON
			body, _ := json.Marshal(tt.payload)

			// Tạo request test
			req, err := http.NewRequest(
				"POST",
				baseURL+"/auth/login",
				bytes.NewBuffer(body),
			)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			// Thực thi request
			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Kiểm tra kết quả
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestHandleRegister(t *testing.T) {
	waitForServer(t)

	tests := []struct {
		name       string
		payload    models.UserCreateInput
		wantStatus int
	}{
		{
			name: "Đăng ký thành công",
			payload: models.UserCreateInput{
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Đăng ký thất bại - thiếu email",
			payload: models.UserCreateInput{
				Password: "password123",
				Name:     "New User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Đăng ký thất bại - thiếu password",
			payload: models.UserCreateInput{
				Email: "newuser@example.com",
				Name:  "New User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Đăng ký thất bại - thiếu name",
			payload: models.UserCreateInput{
				Email:    "newuser@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)

			req, err := http.NewRequest(
				"POST",
				baseURL+"/auth/register",
				bytes.NewBuffer(body),
			)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestHandleLogout(t *testing.T) {
	waitForServer(t)

	tests := []struct {
		name       string
		token      string
		payload    models.UserLogoutInput
		wantStatus int
	}{
		{
			name:  "Đăng xuất thành công",
			token: "valid.token",
			payload: models.UserLogoutInput{
				Hwid: "test-device-id",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Đăng xuất thất bại - không có token",
			payload: models.UserLogoutInput{
				Hwid: "test-device-id",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Đăng xuất thất bại - thiếu hwid",
			token:      "valid.token",
			payload:    models.UserLogoutInput{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)

			req, err := http.NewRequest(
				"POST",
				baseURL+"/auth/logout",
				bytes.NewBuffer(body),
			)
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
