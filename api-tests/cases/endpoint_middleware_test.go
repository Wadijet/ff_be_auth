package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"ff_be_auth_tests/utils"
)

// EndpointInfo chứa thông tin về một endpoint để test
type EndpointInfo struct {
	Path           string   // Đường dẫn endpoint
	Method         string   // HTTP method (GET, POST, PUT, DELETE)
	RequiresAuth   bool     // Có yêu cầu authentication không
	RequiresPerm   bool     // Có yêu cầu permission không (cần X-Active-Role-ID)
	Permission     string   // Permission cần thiết (nếu có)
	IsPublic       bool     // Endpoint công khai (không cần auth)
	Description    string   // Mô tả endpoint
	TestData       interface{} // Dữ liệu test (cho POST/PUT)
	ExpectedStatus int      // Status code mong đợi (0 = không kiểm tra)
}

// TestAllEndpointsMiddleware kiểm tra middleware có được gọi cho tất cả các endpoint
func TestAllEndpointsMiddleware(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Lấy token nếu cần
	firebaseIDToken := getTestFirebaseIDToken(t)
	var token string
	var activeRoleID string

	if firebaseIDToken != "" {
		client := utils.NewHTTPClient(baseURL, 10)
		payload := map[string]interface{}{
			"idToken": firebaseIDToken,
			"hwid":    "test_device_all_endpoints",
		}

		resp, respBody, err := client.POST("/auth/login/firebase", payload)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var loginResult map[string]interface{}
			json.Unmarshal(respBody, &loginResult)
			if data, ok := loginResult["data"].(map[string]interface{}); ok {
				if t, ok := data["token"].(string); ok {
					token = t
				}
			}

			// Lấy roles để có activeRoleID
			if token != "" {
				client.SetToken(token)
				resp, respBody, err = client.GET("/auth/roles")
				if err == nil && resp.StatusCode == http.StatusOK {
					defer resp.Body.Close()
					var rolesResult map[string]interface{}
					json.Unmarshal(respBody, &rolesResult)
					if rolesData, ok := rolesResult["data"].([]interface{}); ok && len(rolesData) > 0 {
						if firstRole, ok := rolesData[0].(map[string]interface{}); ok {
							if roleID, ok := firstRole["id"].(string); ok {
								activeRoleID = roleID
							}
						}
					}
				}
			}
		}
	}

	// Định nghĩa tất cả các endpoint cần test
	endpoints := getAllEndpoints()

	// Thống kê
	stats := struct {
		Total          int
		Tested         int
		Passed         int
		Failed         int
		Skipped        int
		MiddlewareOK   int
		MiddlewareFail int
	}{
		Total: len(endpoints),
	}

	fmt.Printf("\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("🧪 TEST TẤT CẢ ENDPOINT - KIỂM TRA MIDDLEWARE\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("📊 Tổng số endpoint: %d\n", stats.Total)
	fmt.Printf("🔑 Token: %s\n", func() string {
		if token != "" {
			return "✅ Có"
		}
		return "❌ Không có"
	}())
	fmt.Printf("👤 Active Role ID: %s\n", func() string {
		if activeRoleID != "" {
			return "✅ " + activeRoleID[:8] + "..."
		}
		return "❌ Không có"
	}())
	fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

	// Test từng endpoint
	for i, endpoint := range endpoints {
		stats.Tested++
		endpointNum := i + 1

		t.Run(fmt.Sprintf("[%d/%d] %s %s", endpointNum, stats.Total, endpoint.Method, endpoint.Path), func(t *testing.T) {
			// Tạo client
			client := utils.NewHTTPClient(baseURL, 10)

			// Set token nếu endpoint yêu cầu auth
			if endpoint.RequiresAuth && token != "" {
				client.SetToken(token)
			}

			// Set active role ID nếu endpoint yêu cầu permission
			if endpoint.RequiresPerm && activeRoleID != "" {
				client.SetActiveRoleID(activeRoleID)
			}

			// Thực hiện request
			var resp *http.Response
			var respBody []byte
			var err error

			switch endpoint.Method {
			case "GET":
				resp, respBody, err = client.GET(endpoint.Path)
			case "POST":
				resp, respBody, err = client.POST(endpoint.Path, endpoint.TestData)
			case "PUT":
				resp, respBody, err = client.PUT(endpoint.Path, endpoint.TestData)
			case "DELETE":
				resp, respBody, err = client.DELETE(endpoint.Path)
			default:
				t.Skipf("⚠️ Method %s chưa được hỗ trợ", endpoint.Method)
				stats.Skipped++
				return
			}

			if err != nil {
				t.Logf("❌ Lỗi khi gọi API: %v", err)
				stats.Failed++
				return
			}
			defer resp.Body.Close()

			// Kiểm tra middleware toàn cục
			middlewareOK := true
			middlewareIssues := []string{}

			// 1. Kiểm tra Request ID Middleware
			requestID := resp.Header.Get("X-Request-ID")
			if requestID == "" {
				middlewareOK = false
				middlewareIssues = append(middlewareIssues, "❌ Thiếu X-Request-ID header (Request ID Middleware)")
			}

			// 2. Kiểm tra Security Headers Middleware
			securityHeaders := map[string]string{
				"X-Content-Type-Options": "nosniff",
				"X-Frame-Options":        "DENY",
				"X-XSS-Protection":        "1; mode=block",
			}
			for header, expectedValue := range securityHeaders {
				actualValue := resp.Header.Get(header)
				if actualValue == "" {
					middlewareOK = false
					middlewareIssues = append(middlewareIssues, fmt.Sprintf("❌ Thiếu %s header", header))
				} else if actualValue != expectedValue {
					middlewareOK = false
					middlewareIssues = append(middlewareIssues, fmt.Sprintf("⚠️ %s không đúng: %s (mong đợi: %s)", header, actualValue, expectedValue))
				}
			}

			// 3. Kiểm tra AuthMiddleware (nếu endpoint yêu cầu auth)
			if endpoint.RequiresAuth {
				// Nếu không có token, phải trả về 401
				if token == "" {
					if resp.StatusCode != http.StatusUnauthorized {
						middlewareOK = false
						middlewareIssues = append(middlewareIssues, fmt.Sprintf("❌ AuthMiddleware không từ chối request không có token (Status: %d, mong đợi: 401)", resp.StatusCode))
					} else {
						middlewareIssues = append(middlewareIssues, "✅ AuthMiddleware đã từ chối request không có token")
					}
				} else {
					// Có token, kiểm tra response
					var result map[string]interface{}
					json.Unmarshal(respBody, &result)
					message, _ := result["message"].(string)

					// Nếu endpoint yêu cầu permission và thiếu X-Active-Role-ID
					if endpoint.RequiresPerm && activeRoleID == "" {
						if resp.StatusCode == http.StatusBadRequest {
							if strings.Contains(strings.ToLower(message), "x-active-role-id") || strings.Contains(strings.ToLower(message), "role") {
								middlewareIssues = append(middlewareIssues, "✅ AuthMiddleware đã từ chối request thiếu X-Active-Role-ID")
							} else {
								middlewareOK = false
								middlewareIssues = append(middlewareIssues, fmt.Sprintf("⚠️ AuthMiddleware từ chối nhưng message không rõ: %s", message))
							}
						} else if resp.StatusCode == http.StatusOK {
							middlewareOK = false
							middlewareIssues = append(middlewareIssues, "❌ AuthMiddleware KHÔNG từ chối request thiếu X-Active-Role-ID (Status: 200)")
						}
					}

					// Nếu có token và có X-Active-Role-ID, kiểm tra xem có bị từ chối do permission không
					if endpoint.RequiresPerm && activeRoleID != "" {
						if resp.StatusCode == http.StatusForbidden {
							middlewareIssues = append(middlewareIssues, "✅ AuthMiddleware đã check permission (403 - có thể user không có quyền)")
						} else if resp.StatusCode == http.StatusOK {
							middlewareIssues = append(middlewareIssues, "✅ AuthMiddleware đã cho phép (200 - user có quyền)")
						}
					}
				}
			}

			// 4. Kiểm tra OrganizationContextMiddleware (nếu endpoint yêu cầu permission)
			if endpoint.RequiresPerm && token != "" && activeRoleID != "" {
				// Middleware này set context, không có header riêng để kiểm tra
				// Nhưng nếu request thành công (200) hoặc bị từ chối do permission (403),
				// có nghĩa là middleware đã chạy
				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusForbidden {
					middlewareIssues = append(middlewareIssues, "✅ OrganizationContextMiddleware đã set context")
				}
			}

			// In kết quả
			statusIcon := "✅"
			if !middlewareOK {
				statusIcon = "❌"
				stats.Failed++
				stats.MiddlewareFail++
			} else {
				stats.Passed++
				stats.MiddlewareOK++
			}

			fmt.Printf("%s [%d/%d] %s %s\n", statusIcon, endpointNum, stats.Total, endpoint.Method, endpoint.Path)
			if endpoint.Description != "" {
				fmt.Printf("   📝 %s\n", endpoint.Description)
			}
			fmt.Printf("   📊 Status: %d", resp.StatusCode)
			if endpoint.ExpectedStatus > 0 && resp.StatusCode != endpoint.ExpectedStatus {
				fmt.Printf(" (mong đợi: %d)", endpoint.ExpectedStatus)
			}
			fmt.Printf("\n")

			if requestID != "" {
				fmt.Printf("   🆔 Request ID: %s\n", requestID)
			}

			// In các vấn đề về middleware
			if len(middlewareIssues) > 0 {
				for _, issue := range middlewareIssues {
					fmt.Printf("   %s\n", issue)
				}
			}

			// In response message nếu có
			var result map[string]interface{}
			if json.Unmarshal(respBody, &result) == nil {
				if message, ok := result["message"].(string); ok && message != "" {
					if len(message) > 100 {
						message = message[:100] + "..."
					}
					fmt.Printf("   💬 Message: %s\n", message)
				}
			}

			fmt.Printf("\n")

			// Assert middleware OK
			if !middlewareOK {
				t.Errorf("❌ Middleware không hoạt động đúng cho endpoint %s %s", endpoint.Method, endpoint.Path)
			}
		})
	}

	// In tổng kết
	fmt.Printf("\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("📊 TỔNG KẾT\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("📈 Tổng số endpoint: %d\n", stats.Total)
	fmt.Printf("✅ Đã test: %d\n", stats.Tested)
	fmt.Printf("✅ Passed: %d\n", stats.Passed)
	fmt.Printf("❌ Failed: %d\n", stats.Failed)
	fmt.Printf("⏭️  Skipped: %d\n", stats.Skipped)
	fmt.Printf("🔒 Middleware OK: %d\n", stats.MiddlewareOK)
	fmt.Printf("⚠️  Middleware Fail: %d\n", stats.MiddlewareFail)
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
}

// getAllEndpoints trả về danh sách tất cả các endpoint cần test
func getAllEndpoints() []EndpointInfo {
	return []EndpointInfo{
		// System Routes
		{
			Path:        "/system/health",
			Method:      "GET",
			RequiresAuth: false,
			IsPublic:    true,
			Description: "Health check endpoint",
		},

		// Init Routes (chỉ khi chưa có admin)
		{
			Path:        "/init/status",
			Method:      "GET",
			RequiresAuth: false,
			IsPublic:    true,
			Description: "Kiểm tra trạng thái init",
		},

		// Auth Routes
		{
			Path:        "/auth/login/firebase",
			Method:      "POST",
			RequiresAuth: false,
			IsPublic:    true,
			Description: "Đăng nhập với Firebase",
			TestData: map[string]interface{}{
				"idToken": "test_token",
				"hwid":    "test_device",
			},
		},
		{
			Path:        "/auth/profile",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: false,
			Description: "Lấy profile user",
		},
		{
			Path:        "/auth/roles",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: false,
			Description: "Lấy danh sách roles của user",
		},
		{
			Path:        "/auth/logout",
			Method:      "POST",
			RequiresAuth: true,
			RequiresPerm: false,
			Description: "Đăng xuất",
		},

		// RBAC Routes - User
		{
			Path:        "/user/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "User.Read",
			Description: "Tìm tất cả users",
		},
		{
			Path:        "/user/find-one",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "User.Read",
			Description: "Tìm một user",
		},
		{
			Path:        "/user/count",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "User.Read",
			Description: "Đếm số lượng users",
		},

		// RBAC Routes - Permission
		{
			Path:        "/permission/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "Permission.Read",
			Description: "Tìm tất cả permissions",
		},
		{
			Path:        "/permission/count",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "Permission.Read",
			Description: "Đếm số lượng permissions",
		},

		// RBAC Routes - Role
		{
			Path:        "/role/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "Role.Read",
			Description: "Tìm tất cả roles",
		},
		{
			Path:        "/role/count",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "Role.Read",
			Description: "Đếm số lượng roles",
		},

		// RBAC Routes - Organization
		{
			Path:        "/organization/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "Organization.Read",
			Description: "Tìm tất cả organizations",
		},
		{
			Path:        "/organization/count",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "Organization.Read",
			Description: "Đếm số lượng organizations",
		},

		// RBAC Routes - RolePermission
		{
			Path:        "/role-permission/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "RolePermission.Read",
			Description: "Tìm tất cả role permissions",
		},

		// RBAC Routes - UserRole
		{
			Path:        "/user-role/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "UserRole.Read",
			Description: "Tìm tất cả user roles",
		},

		// RBAC Routes - Agent
		{
			Path:        "/agent/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "Agent.Read",
			Description: "Tìm tất cả agents",
		},

		// Facebook Routes - Access Token
		{
			Path:        "/access-token/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "AccessToken.Read",
			Description: "Tìm tất cả access tokens",
		},

		// Facebook Routes - Page
		{
			Path:        "/facebook/page/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "FbPage.Read",
			Description: "Tìm tất cả Facebook pages",
		},

		// Facebook Routes - Post
		{
			Path:        "/facebook/post/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "FbPost.Read",
			Description: "Tìm tất cả Facebook posts",
		},

		// Facebook Routes - Conversation
		{
			Path:        "/facebook/conversation/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "FbConversation.Read",
			Description: "Tìm tất cả Facebook conversations",
		},

		// Facebook Routes - Message
		{
			Path:        "/facebook/message/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "FbMessage.Read",
			Description: "Tìm tất cả Facebook messages",
		},

		// Notification Routes - Sender
		{
			Path:        "/notification/sender/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "NotificationSender.Read",
			Description: "Tìm tất cả notification senders",
		},

		// Notification Routes - Channel
		{
			Path:        "/notification/channel/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "NotificationChannel.Read",
			Description: "Tìm tất cả notification channels",
		},

		// Notification Routes - Template
		{
			Path:        "/notification/template/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "NotificationTemplate.Read",
			Description: "Tìm tất cả notification templates",
		},

		// Notification Routes - Routing
		{
			Path:        "/notification/routing/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "NotificationRouting.Read",
			Description: "Tìm tất cả notification routing rules",
		},

		// Notification Routes - History
		{
			Path:        "/notification/history/find",
			Method:      "GET",
			RequiresAuth: true,
			RequiresPerm: true,
			Permission:  "NotificationHistory.Read",
			Description: "Tìm tất cả notification history",
		},

		// Notification Routes - Tracking (public)
		{
			Path:        "/notification/track/open/test_history_id",
			Method:      "GET",
			RequiresAuth: false,
			IsPublic:    true,
			Description: "Track notification open (public endpoint)",
		},
	}
}
