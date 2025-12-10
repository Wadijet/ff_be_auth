# Viết Test Case

Hướng dẫn cách viết test case cho hệ thống.

## 📋 Tổng Quan

Test cases được viết bằng Go và nằm trong thư mục `api-tests/cases/`.

## 🏗️ Cấu Trúc Test

### Basic Test Structure

```go
package cases

import (
    "testing"
    "api-tests/utils"
)

func TestExample(t *testing.T) {
    // Setup
    client := utils.NewHTTPClient("http://localhost:8080/api/v1")
    
    // Test
    resp, err := client.Get("/system/health")
    if err != nil {
        t.Fatalf("Failed to get health: %v", err)
    }
    
    // Assert
    if resp.StatusCode != 200 {
        t.Errorf("Expected status 200, got %d", resp.StatusCode)
    }
}
```

## 📝 Test Utilities

### HTTP Client

**Vị trí:** `api-tests/utils/http_client.go`

```go
client := utils.NewHTTPClient("http://localhost:8080/api/v1")

// GET request
resp, err := client.Get("/endpoint")

// POST request
resp, err := client.Post("/endpoint", body)

// With authentication
client.SetToken("jwt-token")
resp, err := client.Get("/protected-endpoint")
```

### Test Fixtures

**Vị trí:** `api-tests/utils/test_fixtures.go`

```go
// Tạo test user
user := utils.CreateTestUser()

// Tạo test role
role := utils.CreateTestRole()

// Lấy Firebase token
token := utils.GetFirebaseToken("firebase-uid")
```

## ✅ Best Practices

1. **Tên Test**: Mô tả rõ ràng test case
2. **Setup/Teardown**: Cleanup sau mỗi test
3. **Assertions**: Kiểm tra kỹ kết quả
4. **Isolation**: Mỗi test độc lập
5. **Error Handling**: Xử lý lỗi đúng cách

## 📚 Tài Liệu Liên Quan

- [Tổng Quan Testing](tong-quan.md)
- [Chạy Test Suite](chay-test.md)
- [Báo Cáo Test](bao-cao-test.md)

