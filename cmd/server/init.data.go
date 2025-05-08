package main

import (
	"meta_commerce/core/api/services"

	"github.com/sirupsen/logrus"
)

func InitDefaultData() {
	initDefaultPermissions()

}

// Hàm khởi tạo các quyền mặc định
func initDefaultPermissions() {
	initService, err := services.NewInitService()
	if err != nil {
		logrus.Fatalf("Failed to initialize init service: %v", err) // Ghi log lỗi nếu khởi tạo init service thất bại
	}
	initService.InitPermission()
	initService.CheckPermissionForAdministrator()
}
