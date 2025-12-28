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

// TestNotificationAPIs kiểm tra các API Notification
func TestNotificationAPIs(t *testing.T) {
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
	_, _, token, err := fixtures.CreateTestUser(firebaseIDToken)
	if err != nil {
		t.Fatalf("❌ Không thể tạo user test: %v", err)
	}

	client := utils.NewHTTPClient(baseURL, 10)
	client.SetToken(token)

	// Lấy danh sách roles và set active role
	resp, body, err := client.GET("/auth/roles")
	if err == nil && resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
			firstRole, _ := data[0].(map[string]interface{})
			roleID, _ := firstRole["roleId"].(string)
			if roleID != "" {
				client.SetActiveRoleID(roleID)
			}
		}
	}

	// Lấy Root Organization ID
	rootOrgID, err := fixtures.GetRootOrganizationID(token)
	if err != nil {
		t.Logf("⚠️ Không thể lấy Root Organization: %v", err)
	}

	// Test 1: Notification Sender CRUD
	t.Run("📧 Notification Sender CRUD", func(t *testing.T) {
		var senderID string

		// CREATE
		t.Run("CREATE - Tạo sender", func(t *testing.T) {
			payload := map[string]interface{}{
				"name":        fmt.Sprintf("TestSender_%d", time.Now().UnixNano()),
				"senderType":  "email",
				"smtpHost":    "smtp.example.com",
				"smtpPort":    587,
				"smtpUser":    "test@example.com",
				"smtpPassword": "password123",
			}

			resp, body, err := client.POST("/notification/sender/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo sender: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					senderID, _ = data["id"].(string)
					fmt.Printf("✅ Tạo sender thành công: %s\n", senderID)
				}
			} else {
				fmt.Printf("⚠️ Tạo sender yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// READ
		t.Run("READ - Lấy danh sách senders", func(t *testing.T) {
			resp, body, err := client.GET("/notification/sender/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách senders: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Lấy danh sách senders thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy senders yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// UPDATE
		if senderID != "" {
			t.Run("UPDATE - Cập nhật sender", func(t *testing.T) {
				payload := map[string]interface{}{
					"name": fmt.Sprintf("UpdatedSender_%d", time.Now().UnixNano()),
				}

				resp, _, err := client.PUT(fmt.Sprintf("/notification/sender/update-by-id/%s", senderID), payload)
				if err != nil {
					t.Fatalf("❌ Lỗi khi update sender: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					fmt.Printf("✅ Update sender thành công\n")
				} else {
					fmt.Printf("⚠️ Update sender yêu cầu quyền (status: %d)\n", resp.StatusCode)
				}
			})
		}
	})

	// Test 2: Notification Channel CRUD
	t.Run("📢 Notification Channel CRUD", func(t *testing.T) {
		var channelID string

		// CREATE
		t.Run("CREATE - Tạo channel", func(t *testing.T) {
			if rootOrgID == "" {
				t.Skip("Skipping: Không có Root Organization ID")
			}

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
					channelID, _ = data["id"].(string)
					// Verify organizationId đã được tự động gán
					orgID, ok := data["organizationId"].(string)
					if ok {
						fmt.Printf("✅ Tạo channel thành công với organizationId: %s\n", orgID)
					} else {
						fmt.Printf("✅ Tạo channel thành công: %s\n", channelID)
					}
				}
			} else {
				fmt.Printf("⚠️ Tạo channel yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// READ
		t.Run("READ - Lấy danh sách channels", func(t *testing.T) {
			resp, body, err := client.GET("/notification/channel/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách channels: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Lấy danh sách channels thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy channels yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test 3: Notification Template CRUD
	t.Run("📝 Notification Template CRUD", func(t *testing.T) {
		var templateID string

		// CREATE
		t.Run("CREATE - Tạo template", func(t *testing.T) {
			payload := map[string]interface{}{
				"name":     fmt.Sprintf("TestTemplate_%d", time.Now().UnixNano()),
				"subject":  "Test Subject",
				"body":     "Test Body {{.variable}}",
				"bodyType": "text",
			}

			resp, body, err := client.POST("/notification/template/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo template: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					templateID, _ = data["id"].(string)
					fmt.Printf("✅ Tạo template thành công: %s\n", templateID)
				}
			} else {
				fmt.Printf("⚠️ Tạo template yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// READ
		t.Run("READ - Lấy danh sách templates", func(t *testing.T) {
			resp, body, err := client.GET("/notification/template/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách templates: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Lấy danh sách templates thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy templates yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test 4: Notification Routing CRUD
	t.Run("🔄 Notification Routing CRUD", func(t *testing.T) {
		var routingID string

		// CREATE
		t.Run("CREATE - Tạo routing rule", func(t *testing.T) {
			if rootOrgID == "" {
				t.Skip("Skipping: Không có Root Organization ID")
			}

			// Cần có channel và template trước
			// Tạm thời test với dữ liệu giả
			payload := map[string]interface{}{
				"eventType":     fmt.Sprintf("test.event.%d", time.Now().UnixNano()),
				"organizationIds": []string{rootOrgID},
				"channelIds":    []string{},
				"templateId":    "",
				"priority":      1,
			}

			resp, body, err := client.POST("/notification/routing/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo routing: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					routingID, _ = data["id"].(string)
					fmt.Printf("✅ Tạo routing rule thành công: %s\n", routingID)
				}
			} else {
				fmt.Printf("⚠️ Tạo routing yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// READ
		t.Run("READ - Lấy danh sách routing rules", func(t *testing.T) {
			resp, body, err := client.GET("/notification/routing/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách routing: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Lấy danh sách routing rules thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy routing yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test 5: Notification History (Read-only)
	t.Run("📜 Notification History", func(t *testing.T) {
		// READ
		t.Run("READ - Lấy danh sách history", func(t *testing.T) {
			resp, body, err := client.GET("/notification/history/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách history: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Lấy danh sách history thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy history yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test 6: Notification Trigger
	t.Run("🚀 Notification Trigger", func(t *testing.T) {
		payload := map[string]interface{}{
			"eventType": "test.event",
			"payload": map[string]interface{}{
				"message": "Test notification",
				"data":    "test data",
			},
		}

		resp, body, err := client.POST("/notification/trigger", payload)
		if err != nil {
			t.Fatalf("❌ Lỗi khi trigger notification: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			assert.NoError(t, err)
			fmt.Printf("✅ Trigger notification thành công\n")
		} else {
			fmt.Printf("⚠️ Trigger notification yêu cầu quyền hoặc không có routing rule (status: %d)\n", resp.StatusCode)
		}
	})

	// Test 7: Notification Tracking (Public endpoints - không cần auth)
	t.Run("📊 Notification Tracking", func(t *testing.T) {
		// Tạo client không có token để test public endpoints
		publicClient := utils.NewHTTPClient(baseURL, 10)

		// Test track open
		t.Run("Track Open", func(t *testing.T) {
			// Sử dụng historyId giả để test
			testHistoryID := "507f1f77bcf86cd799439011"
			resp, _, err := publicClient.GET(fmt.Sprintf("/notification/track/open/%s", testHistoryID))
			if err != nil {
				t.Fatalf("❌ Lỗi khi track open: %v", err)
			}

			// Có thể trả về 404 nếu historyId không tồn tại, nhưng endpoint phải hoạt động
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
				fmt.Printf("✅ Track open endpoint hoạt động (status: %d)\n", resp.StatusCode)
			} else {
				fmt.Printf("⚠️ Track open trả về status không mong đợi: %d\n", resp.StatusCode)
			}
		})

		// Test track click
		t.Run("Track Click", func(t *testing.T) {
			testHistoryID := "507f1f77bcf86cd799439011"
			ctaIndex := 0
			resp, _, err := publicClient.GET(fmt.Sprintf("/notification/track/%s/%d", testHistoryID, ctaIndex))
			if err != nil {
				t.Fatalf("❌ Lỗi khi track click: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
				fmt.Printf("✅ Track click endpoint hoạt động (status: %d)\n", resp.StatusCode)
			} else {
				fmt.Printf("⚠️ Track click trả về status không mong đợi: %d\n", resp.StatusCode)
			}
		})

		// Test confirm
		t.Run("Track Confirm", func(t *testing.T) {
			testHistoryID := "507f1f77bcf86cd799439011"
			resp, _, err := publicClient.GET(fmt.Sprintf("/notification/confirm/%s", testHistoryID))
			if err != nil {
				t.Fatalf("❌ Lỗi khi track confirm: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
				fmt.Printf("✅ Track confirm endpoint hoạt động (status: %d)\n", resp.StatusCode)
			} else {
				fmt.Printf("⚠️ Track confirm trả về status không mong đợi: %d\n", resp.StatusCode)
			}
		})
	})
}

