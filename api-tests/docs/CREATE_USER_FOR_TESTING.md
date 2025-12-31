# Hướng Dẫn Tạo User Cho Testing

## 🎯 Tổng Quan

Hệ thống sử dụng **Firebase Authentication** để xác thực người dùng. Để test, bạn có thể:

1. **Sử dụng Firebase ID Token** (Khuyến nghị) - Tạo user trong Firebase, lấy ID token
2. **First User Becomes Admin** - User đầu tiên đăng nhập tự động trở thành admin
3. **Tạo user trực tiếp trong database** (Chỉ cho test, bypass Firebase) - Không khuyến nghị

## 📋 Các Phương Án

### Phương Án 1: Sử Dụng Firebase ID Token (Khuyến Nghị) ⭐

**Cách hoạt động:**
- Tạo user trong Firebase Console hoặc qua Firebase SDK
- Lấy Firebase ID Token
- Dùng token để login và tạo user trong database

**Ưu điểm:**
- ✅ Giống với production flow
- ✅ Test đầy đủ authentication flow
- ✅ An toàn và đúng với kiến trúc hệ thống

**Cách sử dụng:**

```go
// 1. Set Firebase ID Token
$env:TEST_FIREBASE_ID_TOKEN="your-firebase-id-token"

// 2. Sử dụng helper function
fixtures, email, token, client, err := utils.SetupTestWithAdminUser(t, baseURL)
```

### Phương Án 2: First User Becomes Admin

**Cách hoạt động:**
- User đầu tiên đăng nhập tự động trở thành admin
- Không cần set admin thủ công

**Cách sử dụng:**

```go
// Tạo user đầu tiên (tự động trở thành admin)
fixtures := utils.NewTestFixtures(baseURL)
firebaseIDToken := utils.GetTestFirebaseIDToken()

email, _, token, err := fixtures.CreateTestUser(firebaseIDToken)
// User này tự động trở thành admin
```

### Phương Án 3: Tạo User Trực Tiếp (Không Khuyến Nghị)

**Lưu ý:** 
- User model không có password field (đã deprecated)
- User phải có FirebaseUID để link với Firebase
- Không thể bypass Firebase authentication trong production
- Chỉ nên dùng cho test database operations (không test auth flow)

**Cách sử dụng (Nếu cần):**

```go
// Tạo user trực tiếp trong database (bypass Firebase)
// CHỈ DÙNG CHO TEST, KHÔNG DÙNG CHO PRODUCTION
user := models.User{
    Email:       "test@example.com",
    Name:        "Test User",
    FirebaseUID: "test_fake_uid_123", // Fake UID cho test
    EmailVerified: true,
    IsBlock:     false,
}
// Insert vào database...
```

## 🚀 Helper Functions

### SetupTestWithAdminUser()

Tự động setup test với admin user có full quyền:

```go
fixtures, adminEmail, adminToken, client, err := utils.SetupTestWithAdminUser(t, baseURL)
```

**Tự động:**
- ✅ Wait for health
- ✅ Init data (organization, permissions, roles)
- ✅ Tạo admin user từ Firebase ID token
- ✅ Set active role
- ✅ Sẵn sàng để test

### SetupTestWithRegularUser()

Tự động setup test với user thường:

```go
fixtures, userEmail, userToken, client, err := utils.SetupTestWithRegularUser(t, baseURL)
```

### CreateTestUser()

Tạo user test từ Firebase ID token:

```go
email, firebaseUID, token, err := fixtures.CreateTestUser(firebaseIDToken)
```

### CreateAdminUser()

Tạo admin user từ Firebase ID token:

```go
email, firebaseUID, token, userID, err := fixtures.CreateAdminUser(firebaseIDToken)
```

## 🔑 Lấy Firebase ID Token

### Cách 1: Sử dụng Script

**Windows:**
```powershell
.\scripts\get-firebase-token.ps1 -Email "test@example.com" -Password "Test@123"
```

**Linux/Mac:**
```bash
./scripts/get-firebase-token.sh -e "test@example.com" -p "Test@123"
```

### Cách 2: Từ Firebase Console

1. Đăng nhập Firebase Console
2. Authentication > Users
3. Tạo user test hoặc lấy token từ user hiện có
4. Copy ID token

### Cách 3: Từ Frontend App

1. Đăng nhập vào app
2. Lấy Firebase ID token từ browser console:
```javascript
firebase.auth().currentUser.getIdToken().then(token => console.log(token))
```

## 📝 Ví Dụ Đầy Đủ

### Ví Dụ 1: Test với Admin User

```go
func TestAdminFeature(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"
    
    // Setup với admin user
    fixtures, adminEmail, adminToken, client, err := utils.SetupTestWithAdminUser(t, baseURL)
    if err != nil {
        t.Fatalf("❌ Setup failed: %v", err)
    }
    
    // Test admin API
    resp, body, err := client.GET("/admin/users")
    // ...
}
```

### Ví Dụ 2: Test với Regular User

```go
func TestUserFeature(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"
    
    // Setup với regular user
    fixtures, userEmail, userToken, client, err := utils.SetupTestWithRegularUser(t, baseURL)
    if err != nil {
        t.Fatalf("❌ Setup failed: %v", err)
    }
    
    // Test user API
    resp, body, err := client.GET("/auth/profile")
    // ...
}
```

### Ví Dụ 3: Tạo User Mới và Set Admin

```go
func TestCreateUserAndSetAdmin(t *testing.T) {
    baseURL := "http://localhost:8080/api/v1"
    fixtures := utils.NewTestFixtures(baseURL)
    
    // 1. Tạo user thường
    firebaseIDToken := utils.GetTestFirebaseIDToken()
    email, _, token, err := fixtures.CreateTestUser(firebaseIDToken)
    
    // 2. Lấy userID từ profile
    client := utils.NewHTTPClient(baseURL, 10)
    client.SetToken(token)
    resp, body, err := client.GET("/auth/profile")
    // Parse userID...
    
    // 3. Set làm admin (nếu chưa có admin trong hệ thống)
    resp, _, err = client.POST(fmt.Sprintf("/init/set-administrator/%s", userID), nil)
    // ...
}
```

## ⚠️ Lưu Ý

1. **Firebase ID Token Bắt Buộc**: 
   - Tất cả user phải được tạo qua Firebase
   - Không thể bypass Firebase authentication
   - User model không có password field (đã deprecated)

2. **First User Becomes Admin**:
   - User đầu tiên đăng nhập tự động trở thành admin
   - Chỉ hoạt động khi chưa có admin trong hệ thống

3. **Init APIs Chỉ Hoạt Động Khi Chưa Có Admin**:
   - `/init/all` - Chỉ hoạt động khi chưa có admin
   - `/init/set-administrator/:id` - Chỉ hoạt động khi chưa có admin
   - Khi đã có admin, các API này trả về 404

4. **Database Phải Sạch (Nếu Cần)**:
   - Nếu test với database mới, init APIs sẽ hoạt động
   - Nếu database đã có admin, cần dùng `CreateAdminUser()` hoặc set admin thủ công

## 🔍 Kiểm Tra User Có Quyền

```go
// Kiểm tra user có phải admin không
resp, _, err := client.GET("/admin/users")
if resp.StatusCode == http.StatusOK {
    // User có quyền admin
}

// Kiểm tra user có role gì
resp, body, err := client.GET("/auth/roles")
// Parse và kiểm tra roles
```

## 📚 Tài Liệu Tham Khảo

- `api-tests/utils/test_helper.go` - Helper functions
- `api-tests/utils/test_fixtures.go` - Test fixtures
- `api-tests/docs/SETUP_TEST_USER_WITH_FULL_PERMISSIONS.md` - Hướng dẫn setup user có quyền
- `docs/01-getting-started/khoi-tao.md` - Hướng dẫn khởi tạo hệ thống
