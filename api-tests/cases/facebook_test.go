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

// TestFacebookAPIs kiểm tra các API Facebook integration
func TestFacebookAPIs(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Khởi tạo dữ liệu mặc định trước
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

	// Test AccessToken APIs
	t.Run("🔑 AccessToken APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách access tokens
		t.Run("Lấy danh sách access tokens", func(t *testing.T) {
			resp, body, err := client.GET("/access-token/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách access tokens: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách access tokens thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách access tokens yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

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

	// Test Facebook Page APIs
	t.Run("📄 Facebook Page APIs", func(t *testing.T) {
		var pageID string

		// Test 1: Lấy danh sách pages
		t.Run("Lấy danh sách pages", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/page/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách pages: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách pages thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách pages yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Tạo page mới (nếu có quyền)
		t.Run("Tạo page mới", func(t *testing.T) {
			payload := map[string]interface{}{
				"pageId":          fmt.Sprintf("test_page_%d", time.Now().UnixNano()),
				"pageName":        "Test Page",
				"pageUsername":    "testpage",
				"isSync":          false,
				"accessToken":     "test_token",
				"pageAccessToken": "test_page_token",
			}

			resp, body, err := client.POST("/facebook/page/insert-one", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi tạo page: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				data, ok := result["data"].(map[string]interface{})
				if ok {
					pageID, _ = data["id"].(string)
					// Verify organizationId đã được tự động gán
					orgID, ok := data["organizationId"].(string)
					if ok {
						fmt.Printf("✅ Tạo page thành công với organizationId: %s\n", orgID)
					} else {
						fmt.Printf("✅ Tạo page thành công: %s\n", pageID)
					}
				}
			} else {
				fmt.Printf("⚠️ Tạo page yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 3: Find by page ID (endpoint đặc biệt)
		if pageID != "" {
			t.Run("Find by page ID", func(t *testing.T) {
				testPageID := "test_page_123"
				resp, body, err := client.GET(fmt.Sprintf("/facebook/page/find-by-page-id/%s", testPageID))
				if err != nil {
					t.Fatalf("❌ Lỗi khi tìm page by ID: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					var result map[string]interface{}
					err = json.Unmarshal(body, &result)
					assert.NoError(t, err)
					fmt.Printf("✅ Find by page ID thành công\n")
				} else {
					fmt.Printf("⚠️ Find by page ID yêu cầu quyền hoặc không tìm thấy (status: %d)\n", resp.StatusCode)
				}
			})
		}

		// Test 4: Update token (endpoint đặc biệt)
		if pageID != "" {
			t.Run("Update token", func(t *testing.T) {
				payload := map[string]interface{}{
					"accessToken":     "updated_token",
					"pageAccessToken": "updated_page_token",
				}

				resp, _, err := client.PUT("/facebook/page/update-token", payload)
				if err != nil {
					t.Fatalf("❌ Lỗi khi update token: %v", err)
				}

				if resp.StatusCode == http.StatusOK {
					fmt.Printf("✅ Update token thành công\n")
				} else {
					fmt.Printf("⚠️ Update token yêu cầu quyền (status: %d)\n", resp.StatusCode)
				}
			})
		}
	})

	// Test Facebook Post APIs
	t.Run("📝 Facebook Post APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách posts
		t.Run("Lấy danh sách posts", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/post/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách posts: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách posts thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách posts yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Find by post ID (endpoint đặc biệt)
		t.Run("Find by post ID", func(t *testing.T) {
			testPostID := "test_post_123"
			resp, body, err := client.GET(fmt.Sprintf("/facebook/post/find-by-post-id/%s", testPostID))
			if err != nil {
				t.Fatalf("❌ Lỗi khi tìm post by ID: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Find by post ID thành công\n")
			} else {
				fmt.Printf("⚠️ Find by post ID yêu cầu quyền hoặc không tìm thấy (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test Facebook Conversation APIs
	t.Run("💬 Facebook Conversation APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách conversations
		t.Run("Lấy danh sách conversations", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/conversation/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách conversations: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách conversations thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách conversations yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Lấy conversations sắp xếp theo API update
		t.Run("Lấy conversations sắp xếp theo API update", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/conversation/sort-by-api-update")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy conversations sort: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy conversations sort thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy conversations sort yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test Facebook Message APIs
	t.Run("📨 Facebook Message APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách messages
		t.Run("Lấy danh sách messages", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/message/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách messages: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Phải parse được JSON response")
				fmt.Printf("✅ Lấy danh sách messages thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy danh sách messages yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Upsert messages (endpoint đặc biệt)
		t.Run("Upsert messages", func(t *testing.T) {
			payload := map[string]interface{}{
				"conversationId": "test_conv_123",
				"pageId":         "test_page_123",
				"pageUsername":   "testpage",
				"customerId":     "test_customer_123",
				"panCakeData": map[string]interface{}{
					"id":             "test_conv_123",
					"conversation_id": "test_conv_123",
					"messages": []interface{}{
						map[string]interface{}{
							"id":         "msg_1",
							"message":     "Test message",
							"inserted_at": "2024-01-01T00:00:00.000000",
						},
					},
				},
				"hasMore": false,
			}

			resp, body, err := client.POST("/facebook/message/upsert-messages", payload)
			if err != nil {
				t.Fatalf("❌ Lỗi khi upsert messages: %v", err)
			}

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Upsert messages thành công\n")
			} else {
				fmt.Printf("⚠️ Upsert messages yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})
	})

	// Test Facebook Message Item APIs
	t.Run("📨 Facebook Message Item APIs", func(t *testing.T) {
		// Test 1: Lấy danh sách message items
		t.Run("Lấy danh sách message items", func(t *testing.T) {
			resp, body, err := client.GET("/facebook/message-item/find")
			if err != nil {
				t.Fatalf("❌ Lỗi khi lấy danh sách message items: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Lấy danh sách message items thành công\n")
			} else {
				fmt.Printf("⚠️ Lấy message items yêu cầu quyền (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 2: Find by conversation ID (endpoint đặc biệt)
		t.Run("Find by conversation ID", func(t *testing.T) {
			testConvID := "test_conv_123"
			resp, body, err := client.GET(fmt.Sprintf("/facebook/message-item/find-by-conversation/%s", testConvID))
			if err != nil {
				t.Fatalf("❌ Lỗi khi tìm message items by conversation: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Find by conversation ID thành công\n")
			} else {
				fmt.Printf("⚠️ Find by conversation ID yêu cầu quyền hoặc không tìm thấy (status: %d)\n", resp.StatusCode)
			}
		})

		// Test 3: Find by message ID (endpoint đặc biệt)
		t.Run("Find by message ID", func(t *testing.T) {
			testMsgID := "test_msg_123"
			resp, body, err := client.GET(fmt.Sprintf("/facebook/message-item/find-by-message-id/%s", testMsgID))
			if err != nil {
				t.Fatalf("❌ Lỗi khi tìm message item by message ID: %v", err)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)
				fmt.Printf("✅ Find by message ID thành công\n")
			} else {
				fmt.Printf("⚠️ Find by message ID yêu cầu quyền hoặc không tìm thấy (status: %d)\n", resp.StatusCode)
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
