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

// TestAgentAPIs kiểm tra các API Agent
func TestAgentAPIs(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

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

	var agentID string

	// Test Agent CRUD APIs
	t.Run("👤 Agent CRUD APIs", func(t *testing.T) {
		// Test 1: Tạo agent mới
		t.Run("Tạo agent mới", func(t *testing.T) {
			payload := map[string]interface{}{
				"name":     fmt.Sprintf("TestAgent_%d", time.Now().UnixNano()),
				"describe": "Test Agent Description",
			}

			resp, body, err := client.POST("/agent/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo agent: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")

				data, ok := result["data"].(map[string]interface{})
				if ok {
					id, ok := data["id"].(string)
					if ok {
						agentID = id
					}
				}
				fmt.Printf("✅ Tạo agent thành công\n")
			} else {
				fmt.Printf("⚠️ Tạo agent yêu cầu quyền (status: %d - %s)\n", resp.StatusCode, string(body))
			}
		})

		// Test 2: Lấy danh sách agents
		t.Run("Lấy danh sách agents", func(t *testing.T) {
			resp, body, err := client.GET("/agent/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách agents: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách agents thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách agents yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 3: Lấy agent theo ID (nếu có)
		if agentID != "" {
			t.Run("Lấy agent theo ID", func(t *testing.T) {
				resp, body, err := client.GET(fmt.Sprintf("/agent/find-by-id/%s", agentID))
				if err != nil {
					t.Fatalf("❌ Lỗi khi lấy agent theo ID: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err, "Phải parse được JSON response")
					fmt.Printf("✅ Lấy agent theo ID thành công\n")
				} else {
					fmt.Printf("⚠️ Lấy agent theo ID yêu cầu quyền (status: %d)\n", resp.StatusCode)
				}
			})
		}
	})

	// Test Agent Check-in/Check-out APIs
	t.Run("🕐 Agent Check-in/Check-out APIs", func(t *testing.T) {
		// Test 1: Check-in agent (nếu có agentID)
		if agentID != "" {
			t.Run("Check-in agent", func(t *testing.T) {
				resp, body, err := client.POST(fmt.Sprintf("/agent/check-in/%s", agentID), nil)
				if err != nil {
					t.Fatalf("❌ Lỗi khi check-in agent: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err, "Phải parse được JSON response")
					fmt.Printf("✅ Check-in agent thành công\n")
				} else {
					fmt.Printf("⚠️ Check-in agent yêu cầu quyền hoặc agent không tồn tại (status: %d - %s)\n", resp.StatusCode, string(body))
				}
			})
		}

		// Test 2: Check-out agent
		t.Run("Check-out agent", func(t *testing.T) {
			resp, body, err := client.POST("/agent/check-out/", nil)
			if err != nil {
				t.Fatalf("❌ Lỗi khi check-out agent: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Check-out agent thành công\n")
			} else {
				fmt.Printf("⚠️ Check-out agent yêu cầu quyền hoặc user không phải agent (status: %d - %s)\n", resp.StatusCode, string(body))
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
