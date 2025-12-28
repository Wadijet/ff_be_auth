package main

import (
	"meta_commerce/core/api/services"
	"meta_commerce/core/global"

	"github.com/sirupsen/logrus"
)

func InitDefaultData() {
	initService, err := services.NewInitService()
	if err != nil {
		logrus.Fatalf("Failed to initialize init service: %v", err)
	}

	// 1. Khởi tạo Organization Root (PHẢI LÀM TRƯỚC)
	if err := initService.InitRootOrganization(); err != nil {
		logrus.Fatalf("Failed to initialize root organization: %v", err)
	}

	// 2. Khởi tạo Permissions (tạo các quyền mới nếu chưa có, bao gồm Customer, FbMessageItem, ...)
	if err := initService.InitPermission(); err != nil {
		logrus.Fatalf("Failed to initialize permissions: %v", err)
	}
	logrus.Info("Permissions initialized/updated successfully")

	// 3. Tạo Role Administrator (nếu chưa có) + Đảm bảo đầy đủ Permission cho Administrator
	// Tự động gán tất cả quyền trong hệ thống (bao gồm quyền mới) cho role Administrator
	if err := initService.CheckPermissionForAdministrator(); err != nil {
		logrus.Warnf("Failed to check permissions for administrator: %v", err)
	} else {
		logrus.Info("Administrator role permissions synchronized successfully")
	}

	// 4. Tạo user admin tự động từ Firebase UID (nếu có config) - Tùy chọn
	// Lưu ý: User phải đã tồn tại trong Firebase Authentication
	// Nếu không có FIREBASE_ADMIN_UID, user đầu tiên login sẽ tự động trở thành admin
	if global.MongoDB_ServerConfig.FirebaseAdminUID != "" {
		if err := initService.InitAdminUser(global.MongoDB_ServerConfig.FirebaseAdminUID); err != nil {
			logrus.Warnf("Failed to initialize admin user from Firebase UID: %v", err)
			logrus.Info("User đầu tiên login sẽ tự động trở thành admin")
		} else {
			logrus.Info("Admin user initialized successfully from Firebase UID")
		}
	} else {
		logrus.Info("FIREBASE_ADMIN_UID not set")
		logrus.Info("User đầu tiên login sẽ tự động trở thành admin (First user becomes admin)")
	}

	// 5. Khởi tạo dữ liệu mặc định cho hệ thống notification
	// Tạo các sender và template mặc định (global), các thông tin như token/password sẽ để trống để admin bổ sung sau
	if err := initService.InitNotificationData(); err != nil {
		logrus.Warnf("Failed to initialize notification data: %v", err)
	} else {
		logrus.Info("Notification data initialized successfully")
	}
}
