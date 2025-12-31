# Hướng Dẫn Tạo User Có Quyền Để Test Full Các Case

## 🎯 Mục Tiêu

Tạo user có đầy đủ quyền (admin) để test tất cả các API endpoints và chức năng của hệ thống.

## 📋 Các Phương Án

### Phương Án 1: Sử Dụng Helper Function (Khuyến Nghị) ⭐

**File**: `api-tests/utils/test_helper.go`

**Cách sử dụng:**

```go
func TestMyFeature(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"
    
    // Setup test với admin user có full quyền
    fixtures, adminEmail, adminToken, client, err := utils.SetupTestWithAdminUser(t, baseURL)
    if err != nil {
        t.Fatalf("❌ Không thể setup test: %v", err)
    }
    
    // client đã được set token và active role
    // Có thể dùng ngay để test
    
    // Test các API
    resp, body, err := client.GET("/some/endpoint")
    // ...
}
```

**Ưu điểm:**
- ✅ Tự động setup tất cả: wait for health, init data, create admin user, set active role
- ✅ Code ngắn gọn, dễ sử dụng
- ✅ Tự động xử lý các edge cases

**Helper Functions có sẵn:**

1. **`SetupTestWithAdminUser()`** - Tạo admin user với full quyền
2. **`SetupTestWithRegularUser()`** - Tạo user thường (không có quyền admin)

### Phương Án 2: Sử Dụng TestFixtures Trực Tiếp

**Cách sử dụng:**

```go
func TestMyFeature(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"
    fixtures := utils.NewTestFixtures(baseURL)
    
    // Lấy Firebase ID token
    firebaseIDToken := utils.GetTestFirebaseIDToken()
    if firebaseIDToken == "" {
        t.Skip("TEST_FIREBASE_ID_TOKEN không được set")
    }
    
    // Tạo admin user
    adminEmail, _, adminToken, userID, err := fixtures.CreateAdminUser(firebaseIDToken)
    if err != nil {
        t.Fatalf("❌ Không thể tạo admin user: %v", err)
    }
    
    // Tạo client và set token
    client := utils.NewHTTPClient(baseURL, 10)
    client.SetToken(adminToken)
    
    // Set active role
    resp, body, err := client.GET("/auth/roles")
    // ... parse và set active role
}
```

### Phương Án 3: Sử Dụng Init API (Chỉ Khi Chưa Có Admin)

**Cách sử dụng:**

```go
func TestMyFeature(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"
    fixtures := utils.NewTestFixtures(baseURL)
    
    // 1. Khởi tạo dữ liệu mặc định (chỉ hoạt động khi chưa có admin)
    err := fixtures.InitData()
    if err != nil {
        t.Logf("ℹ️ Init data: %v (có thể đã init rồi)", err)
    }
    
    // 2. Tạo user thường
    firebaseIDToken := utils.GetTestFirebaseIDToken()
    email, _, token, err := fixtures.CreateTestUser(firebaseIDToken)
    
    // 3. Lấy userID từ profile
    client := utils.NewHTTPClient(baseURL, 10)
    client.SetToken(token)
    resp, body, err := client.GET("/auth/profile")
    // ... parse userID
    
    // 4. Set làm administrator (chỉ hoạt động khi chưa có admin)
    resp, _, err = client.POST(fmt.Sprintf("/init/set-administrator/%s", userID), nil)
    // ...
}
```

## 🔑 Yêu Cầu

### 1. Firebase ID Token

**Cách lấy Firebase ID Token:**

**Option A: Từ Environment Variable**
```bash
# Windows PowerShell
$env:TEST_FIREBASE_ID_TOKEN="your-firebase-id-token"

# Linux/Mac
export TEST_FIREBASE_ID_TOKEN="your-firebase-id-token"
```

**Option B: Sử dụng Script**
```bash
# Windows
.\scripts\get-firebase-token.ps1

# Linux/Mac
./scripts/get-firebase-token.sh
```

**Option C: Từ Firebase Console**
1. Đăng nhập Firebase Console
2. Authentication > Users
3. Tạo user test hoặc lấy token từ user hiện có
4. Copy ID token

### 2. Server Phải Đang Chạy

```bash
# Chạy server
cd folkgroup-backend/api
go run cmd/server/main.go
```

## 📝 Ví Dụ Đầy Đủ

### Ví Dụ 1: Test với Admin User

```go
package tests

import (
    "testing"
    "ff_be_auth_tests/utils"
)

func TestFullFeatureWithAdmin(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"
    
    // Setup với admin user
    fixtures, adminEmail, _, client, err := utils.SetupTestWithAdminUser(t, baseURL)
    if err != nil {
        t.Fatalf("❌ Setup failed: %v", err)
    }
    
    t.Logf("✅ Test với admin user: %s", adminEmail)
    
    // Test các API cần quyền admin
    t.Run("Test Admin API", func(t *testing.T) {
        resp, body, err := client.GET("/admin/users")
        // ... verify response
    })
    
    // Test CRUD operations
    t.Run("Test CRUD", func(t *testing.T) {
        // Create
        payload := map[string]interface{}{
            "name": "Test Resource",
        }
        resp, _, err := client.POST("/resource/insert-one", payload)
        // ... verify
        
        // Read, Update, Delete...
    })
}
```

### Ví Dụ 2: Test với Regular User

```go
func TestFeatureWithRegularUser(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"
    
    // Setup với regular user
    fixtures, userEmail, _, client, err := utils.SetupTestWithRegularUser(t, baseURL)
    if err != nil {
        t.Fatalf("❌ Setup failed: %v", err)
    }
    
    t.Logf("✅ Test với user: %s", userEmail)
    
    // Test các API không cần quyền admin
    t.Run("Test Public API", func(t *testing.T) {
        resp, body, err := client.GET("/public/endpoint")
        // ... verify
    })
}
```

## 🔍 Kiểm Tra User Có Quyền

### Kiểm Tra User Có Phải Admin Không

```go
// Lấy profile để xem roles
resp, body, err := client.GET("/auth/profile")
// Parse response và kiểm tra roles

// Hoặc kiểm tra có thể gọi admin API không
resp, _, err := client.GET("/admin/users")
if resp.StatusCode == http.StatusOK {
    // User có quyền admin
}
```

### Kiểm Tra Permissions

```go
// Lấy danh sách roles
resp, body, err := client.GET("/auth/roles")
// Parse và kiểm tra role có permissions gì

// Hoặc test trực tiếp với API cần permission
resp, _, err := client.POST("/some/endpoint", payload)
if resp.StatusCode == http.StatusForbidden {
    // User không có quyền
}
```

## ⚠️ Lưu Ý

1. **First User Becomes Admin**: User đầu tiên đăng nhập tự động trở thành admin (nếu chưa có admin)

2. **Init APIs Chỉ Hoạt Động Khi Chưa Có Admin**: 
   - `/init/all` - Chỉ hoạt động khi chưa có admin
   - `/init/set-administrator/:id` - Chỉ hoạt động khi chưa có admin
   - Khi đã có admin, các API này sẽ trả về 404

3. **Firebase ID Token Phải Hợp Lệ**: 
   - Token phải từ Firebase Authentication
   - Token phải chưa hết hạn
   - User phải tồn tại trong Firebase

4. **Database Phải Sạch (Nếu Cần)**:
   - Nếu test với database mới, init APIs sẽ hoạt động
   - Nếu database đã có admin, cần dùng `CreateAdminUser()` hoặc set admin thủ công

## 🚀 Quick Start

1. **Set Firebase ID Token:**
```bash
$env:TEST_FIREBASE_ID_TOKEN="your-token"
```

2. **Chạy Test:**
```bash
cd api-tests
go test -v ./cases/admin_full_test.go
```

3. **Test sẽ tự động:**
   - ✅ Đợi server sẵn sàng
   - ✅ Khởi tạo dữ liệu mặc định
   - ✅ Tạo admin user
   - ✅ Set active role
   - ✅ Sẵn sàng để test

## 📚 Tài Liệu Tham Khảo

- `api-tests/utils/test_helper.go` - Helper functions
- `api-tests/utils/test_fixtures.go` - Test fixtures
- `api-tests/cases/admin_full_test.go` - Ví dụ test với admin user
- `docs/01-getting-started/khoi-tao.md` - Hướng dẫn khởi tạo hệ thống
