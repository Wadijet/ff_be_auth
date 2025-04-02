package tests

import (
	"fmt"
	"testing"

	"meta_commerce/tests/utils"

	"github.com/stretchr/testify/assert"
)

// TestUserRegistration test API đăng ký user mới
func TestUserRegistration(t *testing.T) {
	// Gọi API đăng ký
	resp, err := utils.RegisterUser(utils.TestUserData)
	if err != nil {
		t.Fatalf("Failed to call register API: %v", err)
	}
	fmt.Printf("Register response: %+v\n", resp)

	// Kiểm tra response
	status, ok := resp["status"].(string)
	if !ok {
		t.Fatalf("Response status is not string: %v", resp["status"])
	}

	// Nếu user đã tồn tại, đăng nhập để lấy token
	if status == "error" {
		// Gọi API đăng nhập
		resp, err = utils.LoginUser(utils.TestLoginData)
		if err != nil {
			t.Fatalf("Failed to call login API: %v", err)
		}
		fmt.Printf("Login response: %+v\n", resp)

		// Kiểm tra response
		status, ok = resp["status"].(string)
		if !ok {
			t.Fatalf("Response status is not string: %v", resp["status"])
		}
		assert.Equal(t, "success", status)

		data, ok := resp["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Response data is not map: %v", resp["data"])
		}

		// Lưu user ID để dùng cho test tiếp theo
		userId, ok := data["id"].(string)
		if !ok {
			t.Fatalf("User ID is not string: %v", data["id"])
		}
		utils.TestAdminData["userId"] = userId

		// Lưu token để dùng cho test tiếp theo
		token, ok := data["token"].(string)
		if !ok {
			t.Fatalf("Token is not string: %v", data["token"])
		}
		utils.TestAdminData["token"] = token

		// Bỏ qua test đăng nhập vì đã đăng nhập ở đây
		t.Skip("Skip login test because already logged in")
	} else {
		assert.Equal(t, "success", status)

		data, ok := resp["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Response data is not map: %v", resp["data"])
		}
		assert.NotEmpty(t, data)

		// Lưu user ID để dùng cho test tiếp theo
		userId, ok := data["id"].(string)
		if !ok {
			t.Fatalf("User ID is not string: %v", data["id"])
		}
		utils.TestAdminData["userId"] = userId
	}
}

// TestUserLogin test API đăng nhập
func TestUserLogin(t *testing.T) {
	// Gọi API đăng nhập
	resp, err := utils.LoginUser(utils.TestLoginData)
	if err != nil {
		t.Fatalf("Failed to call login API: %v", err)
	}
	fmt.Printf("Login response: %+v\n", resp)

	// Kiểm tra response
	status, ok := resp["status"].(string)
	if !ok {
		t.Fatalf("Response status is not string: %v", resp["status"])
	}
	assert.Equal(t, "success", status)

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Response data is not map: %v", resp["data"])
	}

	token, ok := data["token"].(string)
	if !ok {
		t.Fatalf("Token is not string: %v", data["token"])
	}
	assert.NotEmpty(t, token)

	// Lưu token để dùng cho test tiếp theo
	utils.TestAdminData["token"] = token
}

// TestSetAdministrator test API set quyền Administrator
func TestSetAdministrator(t *testing.T) {
	// Lấy user ID và token từ test trước
	userId, ok := utils.TestAdminData["userId"].(string)
	if !ok {
		t.Fatal("User ID not found in test data")
	}
	token, ok := utils.TestAdminData["token"].(string)
	if !ok {
		t.Fatal("Token not found in test data")
	}

	// Gọi API set quyền admin
	resp, err := utils.SetAdministrator(userId, token)
	if err != nil {
		t.Fatalf("Failed to call set admin API: %v", err)
	}
	fmt.Printf("Set admin response: %+v\n", resp)

	// Kiểm tra response
	status, ok := resp["status"].(string)
	if !ok {
		t.Fatalf("Response status is not string: %v", resp["status"])
	}
	assert.Equal(t, "success", status)
}

// TestUserLogout test API đăng xuất
func TestUserLogout(t *testing.T) {
	// Lấy token từ test trước
	token, ok := utils.TestAdminData["token"].(string)
	if !ok {
		t.Fatal("Token not found in test data")
	}

	// Gọi API đăng xuất
	resp, err := utils.LogoutUser(token)
	if err != nil {
		t.Fatalf("Failed to call logout API: %v", err)
	}
	fmt.Printf("Logout response: %+v\n", resp)

	// Kiểm tra response
	status, ok := resp["status"].(string)
	if !ok {
		t.Fatalf("Response status is not string: %v", resp["status"])
	}
	assert.Equal(t, "success", status)
}
