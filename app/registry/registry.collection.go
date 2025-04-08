// Package registry cung cấp implementation của registry pattern cho MongoDB collections.
package registry

import (
	"meta_commerce/app/global"
	"meta_commerce/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// Collections là registry dùng để quản lý MongoDB collections
var Collections = NewRegistry[*mongo.Collection]()

// InitCollections khởi tạo và đăng ký tất cả các collection cần thiết
// Parameters:
//   - db: Client MongoDB
//   - cfg: Cấu hình hệ thống
//
// Returns:
//   - error: Lỗi nếu có
func InitCollections(db *mongo.Client, cfg *config.Configuration) error {
	// Lấy database
	database := db.Database(cfg.MongoDB_DBNameAuth)

	// Định nghĩa danh sách các collection cần khởi tạo
	collections := []string{
		global.MongoDB_ColNames.Users,
		global.MongoDB_ColNames.Roles,
		global.MongoDB_ColNames.Permissions,
		global.MongoDB_ColNames.RolePermissions,
		global.MongoDB_ColNames.Agents,
		global.MongoDB_ColNames.UserRoles,
		global.MongoDB_ColNames.AccessTokens,
		global.MongoDB_ColNames.PcOrders,
		global.MongoDB_ColNames.FbPages,
		global.MongoDB_ColNames.FbPosts,
		global.MongoDB_ColNames.FbMessages,
		global.MongoDB_ColNames.FbConvesations,
	}

	// Đăng ký từng collection
	for _, name := range collections {
		collection := database.Collection(name)
		if err := Collections.Register(name, collection); err != nil {
			return err
		}
	}

	return nil
}

// GetCollectionNames trả về danh sách tên các collection đã được định nghĩa
// Returns:
//   - []string: Danh sách tên các collection
func GetCollectionNames() []string {
	return []string{
		global.MongoDB_ColNames.Users,
		global.MongoDB_ColNames.Roles,
		global.MongoDB_ColNames.Permissions,
		global.MongoDB_ColNames.RolePermissions,
		global.MongoDB_ColNames.Agents,
		global.MongoDB_ColNames.UserRoles,
		global.MongoDB_ColNames.AccessTokens,
		global.MongoDB_ColNames.PcOrders,
		global.MongoDB_ColNames.FbPages,
		global.MongoDB_ColNames.FbPosts,
		global.MongoDB_ColNames.FbMessages,
		global.MongoDB_ColNames.FbConvesations,
	}
}
