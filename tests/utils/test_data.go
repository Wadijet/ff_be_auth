package utils

import (
	"fmt"
)

// TestUserData chứa dữ liệu test cho user
var TestUserData = map[string]interface{}{
	"name":     "Admin User",
	"email":    "admin_test@example.com",
	"password": "admin123",
}

// TestLoginData chứa dữ liệu test cho đăng nhập
var TestLoginData = map[string]interface{}{
	"email":    "admin_test@example.com",
	"password": "admin123",
	"hwid":     "test_hwid",
}

// TestAdminData chứa dữ liệu test cho set quyền admin
var TestAdminData = map[string]interface{}{
	"userId": "", // Sẽ được set sau khi đăng ký thành công
	"token":  "", // Sẽ được set sau khi đăng nhập thành công
}

// TestUserData2 chứa dữ liệu test cho user thứ 2
var TestUserData2 = map[string]interface{}{
	"name":     "User 2",
	"email":    "user2_test@example.com",
	"password": "user123",
}

// TestLoginData2 chứa dữ liệu test cho đăng nhập user thứ 2
var TestLoginData2 = map[string]interface{}{
	"email":    "user2_test@example.com",
	"password": "user123",
	"hwid":     "test-hwid-2",
}

// TestLogoutData chứa dữ liệu test cho đăng xuất
var TestLogoutData = map[string]interface{}{
	"hwid": "test_hwid",
}

// DeleteTestUser xóa user test cũ
func DeleteTestUser() (map[string]interface{}, error) {
	return CallAPI("DELETE", fmt.Sprintf("/users/by_email/%s", TestUserData["email"]), nil, "")
}
