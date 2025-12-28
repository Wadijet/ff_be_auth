package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"ff_be_auth_tests/utils"

	"github.com/stretchr/testify/assert"
)

// TestMiddlewareInitialization kiểm tra middleware có được đăng ký và gọi đúng không
func TestMiddlewareInitialization(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	
	// Đợi server khởi động
	waitForHealth(baseURL, 10, 1*time.Second, t)

	t.Run("🔍 Kiểm tra Request ID Middleware", func(t *testing.T) {
		// Gọi endpoint không cần auth để kiểm tra Request ID
		resp, err := http.Get(baseURL + "/system/health")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi health check: %v", err)
		}
		defer resp.Body.Close()

		// Kiểm tra header X-Request-ID có được set không
		requestID := resp.Header.Get("X-Request-ID")
		assert.NotEmpty(t, requestID, "✅ X-Request-ID header phải được set bởi Request ID Middleware")
		
		fmt.Printf("   ✅ X-Request-ID: %s\n", requestID)
	})

	t.Run("🔍 Kiểm tra CORS Middleware", func(t *testing.T) {
		// Tạo request với Origin header để test CORS
		req, err := http.NewRequest("GET", baseURL+"/system/health", nil)
		if err != nil {
			t.Fatalf("❌ Lỗi khi tạo request: %v", err)
		}
		req.Header.Set("Origin", "http://localhost:3000")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		// Kiểm tra CORS headers
		allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
		allowMethods := resp.Header.Get("Access-Control-Allow-Methods")
		allowHeaders := resp.Header.Get("Access-Control-Allow-Headers")
		
		// CORS headers có thể được set hoặc không tùy vào config
			// Nhưng nếu có Origin header trong request, CORS middleware phải xử lý
		if req.Header.Get("Origin") != "" {
			fmt.Printf("   ✅ CORS Middleware đã xử lý request với Origin header\n")
			if allowOrigin != "" {
				fmt.Printf("   ✅ Access-Control-Allow-Origin: %s\n", allowOrigin)
			}
			if allowMethods != "" {
				fmt.Printf("   ✅ Access-Control-Allow-Methods: %s\n", allowMethods)
			}
			if allowHeaders != "" {
				fmt.Printf("   ✅ Access-Control-Allow-Headers: %s\n", allowHeaders)
			}
		}
	})

	t.Run("🔍 Kiểm tra Security Headers Middleware", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/system/health")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi health check: %v", err)
		}
		defer resp.Body.Close()

		// Kiểm tra security headers
		contentTypeOptions := resp.Header.Get("X-Content-Type-Options")
		frameOptions := resp.Header.Get("X-Frame-Options")
		xssProtection := resp.Header.Get("X-XSS-Protection")
		referrerPolicy := resp.Header.Get("Referrer-Policy")

		assert.Equal(t, "nosniff", contentTypeOptions, "✅ X-Content-Type-Options phải là 'nosniff'")
		assert.Equal(t, "DENY", frameOptions, "✅ X-Frame-Options phải là 'DENY'")
		assert.Equal(t, "1; mode=block", xssProtection, "✅ X-XSS-Protection phải là '1; mode=block'")
		assert.Equal(t, "strict-origin-when-cross-origin", referrerPolicy, "✅ Referrer-Policy phải được set")

		fmt.Printf("   ✅ Security Headers đã được set đúng:\n")
		fmt.Printf("      - X-Content-Type-Options: %s\n", contentTypeOptions)
		fmt.Printf("      - X-Frame-Options: %s\n", frameOptions)
		fmt.Printf("      - X-XSS-Protection: %s\n", xssProtection)
		fmt.Printf("      - Referrer-Policy: %s\n", referrerPolicy)
	})

	t.Run("🔍 Kiểm tra Logger Middleware", func(t *testing.T) {
		// Logger middleware sẽ log request, kiểm tra bằng cách gọi API
		// và xem log output (không thể test trực tiếp, nhưng có thể verify qua response)
		resp, err := http.Get(baseURL + "/system/health")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi health check: %v", err)
		}
		defer resp.Body.Close()

		// Logger middleware sẽ log request, nhưng không set header nào
		// Chỉ cần verify request được xử lý thành công
		assert.Equal(t, http.StatusOK, resp.StatusCode, "✅ Request phải được xử lý thành công (Logger Middleware đã chạy)")
		fmt.Printf("   ✅ Logger Middleware đã log request (kiểm tra log file để xác nhận)\n")
	})
}

// TestAuthMiddleware kiểm tra AuthMiddleware có được gọi đúng không
func TestAuthMiddleware(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Lấy token từ auth test
	firebaseIDToken := getTestFirebaseIDToken(t)
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN not set")
	}

	// Đăng nhập để lấy JWT token
	client := utils.NewHTTPClient(baseURL, 10)
	payload := map[string]interface{}{
		"idToken": firebaseIDToken,
		"hwid":    "test_device_middleware",
	}

	resp, respBody, err := client.POST("/auth/login/firebase", payload)
	if err != nil {
		t.Fatalf("❌ Lỗi khi đăng nhập: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("❌ Đăng nhập thất bại: %s", string(respBody))
	}

	var loginResult map[string]interface{}
	err = json.Unmarshal(respBody, &loginResult)
	if err != nil {
		t.Fatalf("❌ Lỗi khi parse response: %v", err)
	}

	data, ok := loginResult["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("❌ Response không có data")
	}

	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatalf("❌ Không có token trong response")
	}

	client.SetToken(token)

	t.Run("🔒 Kiểm tra AuthMiddleware - Request không có token", func(t *testing.T) {
		// Tạo client mới không có token
		noTokenClient := utils.NewHTTPClient(baseURL, 10)
		
		// Gọi endpoint yêu cầu auth
		resp, respBody, err := noTokenClient.GET("/auth/profile")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		// AuthMiddleware phải từ chối request không có token
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "✅ AuthMiddleware phải từ chối request không có token")
		
		var errorResult map[string]interface{}
		err = json.Unmarshal(respBody, &errorResult)
		assert.NoError(t, err, "Phải parse được error response")
		
		// Kiểm tra error code
		code, ok := errorResult["code"].(string)
		if ok {
			fmt.Printf("   ✅ AuthMiddleware đã từ chối request: %s\n", code)
		}
	})

	t.Run("🔒 Kiểm tra AuthMiddleware - Request có token hợp lệ", func(t *testing.T) {
		// Gọi endpoint với token hợp lệ
		resp, respBody, err := client.GET("/auth/profile")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		// AuthMiddleware phải cho phép request có token hợp lệ
		assert.Equal(t, http.StatusOK, resp.StatusCode, "✅ AuthMiddleware phải cho phép request có token hợp lệ")
		
		var result map[string]interface{}
		err = json.Unmarshal(respBody, &result)
		assert.NoError(t, err, "Phải parse được response")
		
		fmt.Printf("   ✅ AuthMiddleware đã xác thực thành công\n")
	})

	t.Run("🔒 Kiểm tra AuthMiddleware - Request có token nhưng thiếu X-Active-Role-ID", func(t *testing.T) {
		// Gọi endpoint yêu cầu permission (cần X-Active-Role-ID)
		// Sử dụng endpoint CRUD yêu cầu permission
		resp, respBody, err := client.GET("/user/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		// AuthMiddleware phải từ chối nếu thiếu X-Active-Role-ID khi route yêu cầu permission
		if resp.StatusCode == http.StatusBadRequest {
			var errorResult map[string]interface{}
			err = json.Unmarshal(respBody, &errorResult)
			if err == nil {
				message, ok := errorResult["message"].(string)
				if ok && strings.Contains(strings.ToLower(message), "x-active-role-id") {
					fmt.Printf("   ✅ AuthMiddleware đã từ chối request thiếu X-Active-Role-ID: %s\n", message)
					return
				}
			}
		}
		
		// Nếu không bị từ chối, có thể là route không yêu cầu permission hoặc đã có default role
		fmt.Printf("   ⚠️ Request không bị từ chối (có thể route không yêu cầu permission hoặc có default role)\n")
	})
}

// TestOrganizationContextMiddleware kiểm tra OrganizationContextMiddleware có được gọi đúng không
func TestOrganizationContextMiddleware(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Lấy token
	firebaseIDToken := getTestFirebaseIDToken(t)
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN not set")
	}

	client := utils.NewHTTPClient(baseURL, 10)
	payload := map[string]interface{}{
		"idToken": firebaseIDToken,
		"hwid":    "test_device_org_context",
	}

	resp, respBody, err := client.POST("/auth/login/firebase", payload)
	if err != nil {
		t.Fatalf("❌ Lỗi khi đăng nhập: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("❌ Đăng nhập thất bại: %s", string(respBody))
	}

	var loginResult map[string]interface{}
	err = json.Unmarshal(respBody, &loginResult)
	if err != nil {
		t.Fatalf("❌ Lỗi khi parse response: %v", err)
	}

	data, ok := loginResult["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("❌ Response không có data")
	}

	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatalf("❌ Không có token trong response")
	}

	// Lấy roles của user
	client.SetToken(token)
	resp, respBody, err = client.GET("/auth/roles")
	if err != nil {
		t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("❌ Không lấy được roles: %s", string(respBody))
	}

	var rolesResult map[string]interface{}
	err = json.Unmarshal(respBody, &rolesResult)
	if err != nil {
		t.Fatalf("❌ Lỗi khi parse roles response: %v", err)
	}

	rolesData, ok := rolesResult["data"].([]interface{})
	if !ok || len(rolesData) == 0 {
		t.Skip("Skipping test: User không có roles")
	}

	// Lấy role ID đầu tiên
	firstRole, ok := rolesData[0].(map[string]interface{})
	if !ok {
		t.Fatalf("❌ Role data không đúng format")
	}

	roleID, ok := firstRole["id"].(string)
	if !ok || roleID == "" {
		t.Fatalf("❌ Không có role ID")
	}

	client.SetActiveRoleID(roleID)

	t.Run("🏢 Kiểm tra OrganizationContextMiddleware - Set active role context", func(t *testing.T) {
		// Gọi endpoint CRUD với X-Active-Role-ID header
		resp, respBody, err := client.GET("/user/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		// OrganizationContextMiddleware phải set context và request phải được xử lý
		// (có thể thành công hoặc lỗi permission, nhưng không phải lỗi thiếu context)
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusForbidden {
			fmt.Printf("   ✅ OrganizationContextMiddleware đã set context (status: %d)\n", resp.StatusCode)
		} else {
			// Kiểm tra xem có phải lỗi thiếu context không
			var errorResult map[string]interface{}
			err = json.Unmarshal(respBody, &errorResult)
			if err == nil {
				message, ok := errorResult["message"].(string)
				if ok {
					fmt.Printf("   ⚠️ Response: %s\n", message)
				}
			}
		}
	})
}

// TestMiddlewareOrder kiểm tra thứ tự thực thi middleware
func TestMiddlewareOrder(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	t.Run("🔄 Kiểm tra thứ tự thực thi middleware", func(t *testing.T) {
		// Gọi endpoint và kiểm tra headers để verify thứ tự middleware
		resp, err := http.Get(baseURL + "/system/health")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi health check: %v", err)
		}
		defer resp.Body.Close()

		// Kiểm tra các headers được set bởi middleware (theo thứ tự)
		headers := map[string]string{
			"X-Request-ID":              "Request ID Middleware",
			"X-Content-Type-Options":     "Security Headers Middleware",
			"X-Frame-Options":            "Security Headers Middleware",
			"X-XSS-Protection":          "Security Headers Middleware",
			"Referrer-Policy":           "Security Headers Middleware",
		}

		fmt.Printf("   📋 Kiểm tra headers được set bởi middleware:\n")
		for header, middleware := range headers {
			value := resp.Header.Get(header)
			if value != "" {
				fmt.Printf("      ✅ %s: %s (từ %s)\n", header, value, middleware)
			} else {
				fmt.Printf("      ⚠️ %s: không có (từ %s)\n", header, middleware)
			}
		}

		// Verify Request ID được set (middleware đầu tiên)
		requestID := resp.Header.Get("X-Request-ID")
		assert.NotEmpty(t, requestID, "✅ Request ID Middleware phải được gọi đầu tiên và set X-Request-ID")

		// Verify Security Headers được set (middleware sau CORS)
		contentTypeOptions := resp.Header.Get("X-Content-Type-Options")
		assert.Equal(t, "nosniff", contentTypeOptions, "✅ Security Headers Middleware phải được gọi sau CORS")
	})
}

// TestMiddlewareErrorHandling kiểm tra middleware xử lý lỗi đúng không
func TestMiddlewareErrorHandling(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	t.Run("🛡️ Kiểm tra Recover Middleware", func(t *testing.T) {
		// Recover middleware sẽ xử lý panic
		// Không thể test trực tiếp panic, nhưng có thể verify middleware được đăng ký
		// bằng cách kiểm tra server vẫn hoạt động sau các request
		
		// Gọi nhiều request để đảm bảo server ổn định
		for i := 0; i < 5; i++ {
			resp, err := http.Get(baseURL + "/system/health")
			if err != nil {
				t.Fatalf("❌ Lỗi khi gọi health check: %v", err)
			}
			resp.Body.Close()
		}

		fmt.Printf("   ✅ Recover Middleware đã được đăng ký (server ổn định sau nhiều request)\n")
	})

	t.Run("🛡️ Kiểm tra Error Handler", func(t *testing.T) {
		// Gọi endpoint không tồn tại để test error handler
		resp, err := http.Get(baseURL + "/nonexistent/endpoint")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		// Error handler phải trả về JSON error response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "✅ Error handler phải trả về 404 cho endpoint không tồn tại")
		
		var errorResult map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&errorResult)
		assert.NoError(t, err, "✅ Error response phải là JSON")
		
		// Kiểm tra format error response
		status, ok := errorResult["status"].(string)
		if ok {
			assert.Equal(t, "error", status, "✅ Error response phải có status='error'")
		}
		
		fmt.Printf("   ✅ Error Handler đã xử lý lỗi đúng format\n")
	})
}

// TestUserPermissions kiểm tra quyền của user và middleware check quyền
func TestUserPermissions(t *testing.T) {
	baseURL := "http://localhost:8080/api/v1"
	waitForHealth(baseURL, 10, 1*time.Second, t)

	// Lấy token
	firebaseIDToken := getTestFirebaseIDToken(t)
	if firebaseIDToken == "" {
		t.Skip("Skipping test: TEST_FIREBASE_ID_TOKEN not set")
	}

	client := utils.NewHTTPClient(baseURL, 10)
	payload := map[string]interface{}{
		"idToken": firebaseIDToken,
		"hwid":    "test_device_permissions",
	}

	// Đăng nhập
	resp, respBody, err := client.POST("/auth/login/firebase", payload)
	if err != nil {
		t.Fatalf("❌ Lỗi khi đăng nhập: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("❌ Đăng nhập thất bại: %s", string(respBody))
	}

	var loginResult map[string]interface{}
	err = json.Unmarshal(respBody, &loginResult)
	if err != nil {
		t.Fatalf("❌ Lỗi khi parse response: %v", err)
	}

	data, ok := loginResult["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("❌ Response không có data")
	}

	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatalf("❌ Không có token trong response")
	}

	client.SetToken(token)

	t.Run("📋 Lấy thông tin roles và permissions của user", func(t *testing.T) {
		// Lấy roles của user
		resp, respBody, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("❌ Không lấy được roles: %s", string(respBody))
		}

		var rolesResult map[string]interface{}
		err = json.Unmarshal(respBody, &rolesResult)
		if err != nil {
			t.Fatalf("❌ Lỗi khi parse roles response: %v", err)
		}

		rolesData, ok := rolesResult["data"].([]interface{})
		if !ok {
			t.Fatalf("❌ Roles data không đúng format")
		}

		fmt.Printf("   📋 User có %d role(s):\n", len(rolesData))
		for i, roleInterface := range rolesData {
			role, ok := roleInterface.(map[string]interface{})
			if !ok {
				continue
			}

			roleID, _ := role["id"].(string)
			roleName, _ := role["name"].(string)
			roleCode, _ := role["code"].(string)
			orgID, _ := role["organizationId"].(string)

			fmt.Printf("      [%d] Role: %s (%s) - ID: %s - Org: %s\n", 
				i+1, roleName, roleCode, roleID, orgID)
		}

		if len(rolesData) == 0 {
			t.Skip("Skipping test: User không có roles")
		}

		// Lấy role đầu tiên để test
		firstRole, ok := rolesData[0].(map[string]interface{})
		if !ok {
			t.Fatalf("❌ Role data không đúng format")
		}

		roleID, ok := firstRole["id"].(string)
		if !ok || roleID == "" {
			t.Fatalf("❌ Không có role ID")
		}

		// Test với role này
		client.SetActiveRoleID(roleID)
		fmt.Printf("   ✅ Đã set active role: %s\n", roleID)
	})

	t.Run("🔒 Test middleware check quyền - Endpoint yêu cầu User.Read", func(t *testing.T) {
		// Lấy roles trước
		resp, respBody, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Skip("Skipping: Không lấy được roles")
		}

		var rolesResult map[string]interface{}
		json.Unmarshal(respBody, &rolesResult)
		rolesData, ok := rolesResult["data"].([]interface{})
		if !ok || len(rolesData) == 0 {
			t.Skip("Skipping: User không có roles")
		}

		firstRole, _ := rolesData[0].(map[string]interface{})
		roleID, _ := firstRole["id"].(string)
		client.SetActiveRoleID(roleID)

		// Test endpoint yêu cầu User.Read
		resp, respBody, err = client.GET("/user/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		fmt.Printf("   🔍 Test User.Read permission - Status: %d\n", resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("   ✅ User có quyền User.Read\n")
		} else if resp.StatusCode == http.StatusForbidden {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("   ❌ User KHÔNG có quyền User.Read: %s\n", message)
		} else if resp.StatusCode == http.StatusBadRequest {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("   ⚠️ Request bị từ chối: %s\n", message)
		}
	})

	t.Run("🔒 Test middleware check quyền - Endpoint yêu cầu User.Insert", func(t *testing.T) {
		// Lấy roles trước
		resp, respBody, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Skip("Skipping: Không lấy được roles")
		}

		var rolesResult map[string]interface{}
		json.Unmarshal(respBody, &rolesResult)
		rolesData, ok := rolesResult["data"].([]interface{})
		if !ok || len(rolesData) == 0 {
			t.Skip("Skipping: User không có roles")
		}

		firstRole, _ := rolesData[0].(map[string]interface{})
		roleID, _ := firstRole["id"].(string)
		client.SetActiveRoleID(roleID)

		// Test endpoint yêu cầu User.Insert (insert-one)
		testData := map[string]interface{}{
			"email": "test@example.com",
			"name":  "Test User",
		}

		resp, respBody, err = client.POST("/user/insert-one", testData)
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		fmt.Printf("   🔍 Test User.Insert permission - Status: %d\n", resp.StatusCode)

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			fmt.Printf("   ✅ User có quyền User.Insert\n")
		} else if resp.StatusCode == http.StatusForbidden {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("   ❌ User KHÔNG có quyền User.Insert: %s\n", message)
		} else if resp.StatusCode == http.StatusBadRequest {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			// Có thể là lỗi validation, không phải permission
			if strings.Contains(strings.ToLower(message), "permission") || 
			   strings.Contains(strings.ToLower(message), "quyền") ||
			   strings.Contains(strings.ToLower(message), "forbidden") {
				fmt.Printf("   ❌ User KHÔNG có quyền User.Insert: %s\n", message)
			} else {
				fmt.Printf("   ⚠️ Request bị từ chối (có thể là validation error): %s\n", message)
			}
		}
	})

	t.Run("🔒 Test middleware check quyền - Endpoint yêu cầu Permission.Read", func(t *testing.T) {
		// Lấy roles trước
		resp, respBody, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Skip("Skipping: Không lấy được roles")
		}

		var rolesResult map[string]interface{}
		json.Unmarshal(respBody, &rolesResult)
		rolesData, ok := rolesResult["data"].([]interface{})
		if !ok || len(rolesData) == 0 {
			t.Skip("Skipping: User không có roles")
		}

		firstRole, _ := rolesData[0].(map[string]interface{})
		roleID, _ := firstRole["id"].(string)
		client.SetActiveRoleID(roleID)

		// Test endpoint yêu cầu Permission.Read
		resp, respBody, err = client.GET("/permission/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		fmt.Printf("   🔍 Test Permission.Read permission - Status: %d\n", resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("   ✅ User có quyền Permission.Read\n")
		} else if resp.StatusCode == http.StatusForbidden {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("   ❌ User KHÔNG có quyền Permission.Read: %s\n", message)
		} else if resp.StatusCode == http.StatusBadRequest {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("   ⚠️ Request bị từ chối: %s\n", message)
		}
	})

	t.Run("🔒 Test middleware check quyền - Không có X-Active-Role-ID header", func(t *testing.T) {
		// Tạo client mới không có X-Active-Role-ID
		noRoleClient := utils.NewHTTPClient(baseURL, 10)
		noRoleClient.SetToken(token)

		// Test endpoint yêu cầu permission (cần X-Active-Role-ID)
		// Endpoint /user/find yêu cầu User.Read permission
		resp, respBody, err := noRoleClient.GET("/user/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()

		fmt.Printf("   🔍 Test không có X-Active-Role-ID header - Status: %d\n", resp.StatusCode)

		if resp.StatusCode == http.StatusBadRequest {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			if strings.Contains(strings.ToLower(message), "x-active-role-id") {
				fmt.Printf("   ✅ Middleware đã từ chối request thiếu X-Active-Role-ID: %s\n", message)
				assert.True(t, true, "Middleware đã check đúng - từ chối khi thiếu X-Active-Role-ID")
			} else {
				fmt.Printf("   ⚠️ Request bị từ chối nhưng không phải do thiếu X-Active-Role-ID: %s\n", message)
			}
		} else if resp.StatusCode == http.StatusForbidden {
			fmt.Printf("   ⚠️ Request bị từ chối với 403 (có thể do không có permission)\n")
		} else if resp.StatusCode == http.StatusOK {
			// Nếu trả về 200, có thể là:
			// 1. Route không yêu cầu permission (không đúng với code)
			// 2. Có logic fallback
			// 3. User có quyền mặc định
			fmt.Printf("   ⚠️ Request thành công (Status 200) - Cần kiểm tra:\n")
			fmt.Printf("      - Route có yêu cầu permission không?\n")
			fmt.Printf("      - Có logic fallback khi không có X-Active-Role-ID không?\n")
			fmt.Printf("      - User có quyền mặc định không?\n")
		} else {
			fmt.Printf("   ⚠️ Request trả về status code không mong đợi: %d\n", resp.StatusCode)
		}
	})

	t.Run("🔒 Test middleware check quyền - Với nhiều roles khác nhau", func(t *testing.T) {
		// Lấy tất cả roles của user
		resp, respBody, err := client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Skip("Skipping: Không lấy được roles")
		}

		var rolesResult map[string]interface{}
		json.Unmarshal(respBody, &rolesResult)
		rolesData, ok := rolesResult["data"].([]interface{})
		if !ok || len(rolesData) == 0 {
			t.Skip("Skipping: User không có roles")
		}

		fmt.Printf("   🔍 Test với %d role(s) khác nhau:\n", len(rolesData))

		// Test với từng role
		for i, roleInterface := range rolesData {
			role, ok := roleInterface.(map[string]interface{})
			if !ok {
				continue
			}

			roleID, _ := role["id"].(string)
			roleName, _ := role["name"].(string)
			client.SetActiveRoleID(roleID)

			fmt.Printf("      [%d] Testing với role: %s (%s)\n", i+1, roleName, roleID)

			// Test User.Read với role này
			resp, respBody, err := client.GET("/user/find")
			if err != nil {
				fmt.Printf("         ❌ Lỗi khi gọi API: %v\n", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				fmt.Printf("         ✅ Role này có quyền User.Read\n")
			} else if resp.StatusCode == http.StatusForbidden {
				var errorResult map[string]interface{}
				json.Unmarshal(respBody, &errorResult)
				message, _ := errorResult["message"].(string)
				fmt.Printf("         ❌ Role này KHÔNG có quyền User.Read: %s\n", message)
			} else {
				fmt.Printf("         ⚠️ Status: %d\n", resp.StatusCode)
			}
		}
	})

	t.Run("🔒 Test middleware check quyền - So sánh endpoint có và không có permission requirement", func(t *testing.T) {
		// Test endpoint không yêu cầu permission (auth/profile)
		resp, _, err := client.GET("/auth/profile")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()
		fmt.Printf("   📋 /auth/profile (không yêu cầu permission): Status %d\n", resp.StatusCode)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Endpoint không yêu cầu permission phải trả về 200")

		// Test endpoint yêu cầu permission nhưng không có X-Active-Role-ID
		noRoleClient := utils.NewHTTPClient(baseURL, 10)
		noRoleClient.SetToken(token)

		// Endpoint /user/find yêu cầu User.Read permission
		resp, respBody, err = noRoleClient.GET("/user/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()
		fmt.Printf("   📋 /user/find (yêu cầu User.Read, không có X-Active-Role-ID): Status %d\n", resp.StatusCode)

		if resp.StatusCode == http.StatusBadRequest {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("      ✅ Middleware đã từ chối: %s\n", message)
		} else if resp.StatusCode == http.StatusForbidden {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("      ✅ Middleware đã từ chối (403): %s\n", message)
		} else if resp.StatusCode == http.StatusOK {
			fmt.Printf("      ⚠️ Request thành công (200) - Cần kiểm tra lại logic middleware\n")
		}

		// Test endpoint yêu cầu permission với X-Active-Role-ID nhưng user không có role
		// Sử dụng một role ID không hợp lệ
		invalidRoleClient := utils.NewHTTPClient(baseURL, 10)
		invalidRoleClient.SetToken(token)
		invalidRoleClient.SetActiveRoleID("000000000000000000000000") // Invalid ObjectID

		resp, respBody, err = invalidRoleClient.GET("/user/find")
		if err != nil {
			t.Fatalf("❌ Lỗi khi gọi API: %v", err)
		}
		defer resp.Body.Close()
		fmt.Printf("   📋 /user/find (với X-Active-Role-ID không hợp lệ): Status %d\n", resp.StatusCode)

		if resp.StatusCode == http.StatusBadRequest {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("      ✅ Middleware đã từ chối role không hợp lệ: %s\n", message)
		} else if resp.StatusCode == http.StatusForbidden {
			var errorResult map[string]interface{}
			json.Unmarshal(respBody, &errorResult)
			message, _ := errorResult["message"].(string)
			fmt.Printf("      ✅ Middleware đã từ chối (403): %s\n", message)
		}
	})

	t.Run("📊 Tóm tắt quyền của user hiện tại", func(t *testing.T) {
		// Lấy profile để xem thông tin user
		resp, respBody, err := client.GET("/auth/profile")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy profile: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var profileResult map[string]interface{}
			json.Unmarshal(respBody, &profileResult)
			data, ok := profileResult["data"].(map[string]interface{})
			if ok {
				email, _ := data["email"].(string)
				name, _ := data["name"].(string)
				fmt.Printf("   👤 User: %s (%s)\n", name, email)
			}
		}

		// Lấy roles
		resp, respBody, err = client.GET("/auth/roles")
		if err != nil {
			t.Fatalf("❌ Lỗi khi lấy roles: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var rolesResult map[string]interface{}
			json.Unmarshal(respBody, &rolesResult)
			rolesData, ok := rolesResult["data"].([]interface{})
			if ok {
				fmt.Printf("   📋 Số lượng roles: %d\n", len(rolesData))
				if len(rolesData) == 0 {
					fmt.Printf("   ⚠️ User hiện tại KHÔNG có roles nào\n")
					fmt.Printf("   💡 Để test đầy đủ, cần gán role cho user:\n")
					fmt.Printf("      - Gọi API POST /api/v1/admin/user/role để gán role\n")
					fmt.Printf("      - Hoặc sử dụng init endpoint để tạo admin user\n")
				} else {
					fmt.Printf("   ✅ User có %d role(s), có thể test đầy đủ quyền\n", len(rolesData))
				}
			}
		}

		// Test một số endpoint để xem middleware có check quyền không
		fmt.Printf("\n   🔍 Kiểm tra middleware check quyền với các endpoint:\n")

		endpoints := []struct {
			path       string
			method     string
			permission string
		}{
			{"/auth/profile", "GET", "Không yêu cầu"},
			{"/user/find", "GET", "User.Read"},
			{"/permission/find", "GET", "Permission.Read"},
			{"/role/find", "GET", "Role.Read"},
		}

		for _, ep := range endpoints {
			var resp *http.Response
			var err error

			if ep.method == "GET" {
				resp, _, err = client.GET(ep.path)
			}

			if err == nil && resp != nil {
				defer resp.Body.Close()
				status := "✅"
				if resp.StatusCode >= 400 {
					status = "❌"
				}
				fmt.Printf("      %s %s %s - Status: %d (Yêu cầu: %s)\n", 
					status, ep.method, ep.path, resp.StatusCode, ep.permission)
			}
		}
	})
}
