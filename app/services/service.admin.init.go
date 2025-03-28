package services

import (
	"context"
	"errors"
	"fmt"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/config"
	"meta_commerce/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitService định nghĩa các CRUD repository cho User, Permission và Role
type InitService struct {
	UserCRUD           BaseServiceMongo[models.User]
	PermissionCRUD     BaseServiceMongo[models.Permission]
	RoleCRUD           BaseServiceMongo[models.Role]
	RolePermissionCRUD BaseServiceMongo[models.RolePermission]
	UserRoleCRUD       BaseServiceMongo[models.UserRole]
}

// NewInitService khởi tạo các repository và trả về một đối tượng InitService
func NewInitService(c *config.Configuration, db *mongo.Client) *InitService {
	newService := new(InitService)

	// Khởi tạo các collection
	userCol := GetCollectionFromName(db, GetDBNameFromCollectionName(c, global.MongoDB_ColNames.Users), global.MongoDB_ColNames.Users)
	permissionCol := GetCollectionFromName(db, GetDBNameFromCollectionName(c, global.MongoDB_ColNames.Permissions), global.MongoDB_ColNames.Permissions)
	roleCol := GetCollectionFromName(db, GetDBNameFromCollectionName(c, global.MongoDB_ColNames.Roles), global.MongoDB_ColNames.Roles)
	rolePermissionCol := GetCollectionFromName(db, GetDBNameFromCollectionName(c, global.MongoDB_ColNames.RolePermissions), global.MongoDB_ColNames.RolePermissions)
	userRoleCol := GetCollectionFromName(db, GetDBNameFromCollectionName(c, global.MongoDB_ColNames.UserRoles), global.MongoDB_ColNames.UserRoles)

	// Khởi tạo các service với BaseService
	newService.UserCRUD = NewBaseServiceMongo[models.User](userCol)
	newService.PermissionCRUD = NewBaseServiceMongo[models.Permission](permissionCol)
	newService.RoleCRUD = NewBaseServiceMongo[models.Role](roleCol)
	newService.RolePermissionCRUD = NewBaseServiceMongo[models.RolePermission](rolePermissionCol)
	newService.UserRoleCRUD = NewBaseServiceMongo[models.UserRole](userRoleCol)

	return newService
}

var InitialPermissions = []models.Permission{
	// ====================================  AUTH MODULE =============================================
	// User Management
	{Name: "User.Insert", Describe: "Quyền tạo người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Read", Describe: "Quyền xem danh sách người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Update", Describe: "Quyền cập nhật thông tin người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Delete", Describe: "Quyền xóa người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Block", Describe: "Quyền khóa/mở khóa người dùng", Group: "Auth", Category: "User"},
	{Name: "User.SetRole", Describe: "Quyền phân quyền cho người dùng", Group: "Auth", Category: "User"},

	// Role Management
	{Name: "Role.Insert", Describe: "Quyền tạo vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Read", Describe: "Quyền xem danh sách vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Update", Describe: "Quyền cập nhật vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Delete", Describe: "Quyền xóa vai trò", Group: "Auth", Category: "Role"},

	// Permission Management
	{Name: "Permission.Insert", Describe: "Quyền tạo quyền", Group: "Auth", Category: "Permission"},
	{Name: "Permission.Read", Describe: "Quyền xem danh sách quyền", Group: "Auth", Category: "Permission"},
	{Name: "Permission.Update", Describe: "Quyền cập nhật quyền", Group: "Auth", Category: "Permission"},
	{Name: "Permission.Delete", Describe: "Quyền xóa quyền", Group: "Auth", Category: "Permission"},

	// Role-Permission Management
	{Name: "RolePermission.Insert", Describe: "Quyền tạo phân quyền cho vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Read", Describe: "Quyền xem phân quyền của vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Update", Describe: "Quyền cập nhật phân quyền của vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Delete", Describe: "Quyền xóa phân quyền của vai trò", Group: "Auth", Category: "RolePermission"},

	// User-Role Management
	{Name: "UserRole.Insert", Describe: "Quyền phân công vai trò cho người dùng", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Read", Describe: "Quyền xem vai trò của người dùng", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Update", Describe: "Quyền cập nhật vai trò của người dùng", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Delete", Describe: "Quyền xóa vai trò của người dùng", Group: "Auth", Category: "UserRole"},

	// Agent Management
	{Name: "Agent.Insert", Describe: "Quyền tạo đại lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Read", Describe: "Quyền xem danh sách đại lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Update", Describe: "Quyền cập nhật thông tin đại lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Delete", Describe: "Quyền xóa đại lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.CheckIn", Describe: "Quyền kiểm tra trạng thái đại lý", Group: "Auth", Category: "Agent"},

	// ====================================  PANCAKE MODULE ===========================================
	// Access Token Management
	{Name: "AccessToken.Insert", Describe: "Quyền tạo token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Read", Describe: "Quyền xem danh sách token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Update", Describe: "Quyền cập nhật token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Delete", Describe: "Quyền xóa token", Group: "Pancake", Category: "AccessToken"},

	// Facebook Page Management
	{Name: "FbPage.Insert", Describe: "Quyền tạo trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Read", Describe: "Quyền xem danh sách trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Update", Describe: "Quyền cập nhật thông tin trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Delete", Describe: "Quyền xóa trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.UpdateToken", Describe: "Quyền cập nhật token trang Facebook", Group: "Pancake", Category: "FbPage"},

	// Facebook Conversation Management
	{Name: "FbConversation.Insert", Describe: "Quyền tạo cuộc trò chuyện", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Read", Describe: "Quyền xem danh sách cuộc trò chuyện", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Update", Describe: "Quyền cập nhật cuộc trò chuyện", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Delete", Describe: "Quyền xóa cuộc trò chuyện", Group: "Pancake", Category: "FbConversation"},

	// Facebook Message Management
	{Name: "FbMessage.Insert", Describe: "Quyền tạo tin nhắn", Group: "Pancake", Category: "FbMessage"},
	{Name: "FbMessage.Read", Describe: "Quyền xem danh sách tin nhắn", Group: "Pancake", Category: "FbMessage"},
	{Name: "FbMessage.Update", Describe: "Quyền cập nhật tin nhắn", Group: "Pancake", Category: "FbMessage"},
	{Name: "FbMessage.Delete", Describe: "Quyền xóa tin nhắn", Group: "Pancake", Category: "FbMessage"},

	// Facebook Post Management
	{Name: "FbPost.Insert", Describe: "Quyền tạo bài viết", Group: "Pancake", Category: "FbPost"},
	{Name: "FbPost.Read", Describe: "Quyền xem danh sách bài viết", Group: "Pancake", Category: "FbPost"},
	{Name: "FbPost.Update", Describe: "Quyền cập nhật bài viết", Group: "Pancake", Category: "FbPost"},
	{Name: "FbPost.Delete", Describe: "Quyền xóa bài viết", Group: "Pancake", Category: "FbPost"},

	// Pancake Order Management
	{Name: "PcOrder.Insert", Describe: "Quyền tạo đơn hàng", Group: "Pancake", Category: "PcOrder"},
	{Name: "PcOrder.Read", Describe: "Quyền xem danh sách đơn hàng", Group: "Pancake", Category: "PcOrder"},
	{Name: "PcOrder.Update", Describe: "Quyền cập nhật đơn hàng", Group: "Pancake", Category: "PcOrder"},
	{Name: "PcOrder.Delete", Describe: "Quyền xóa đơn hàng", Group: "Pancake", Category: "PcOrder"},
}

// Viết hàm InitPermission để khởi tạo các quyền mặc định theo nguyên tắc sau:
// Duyệt tất cả các quyền trong mảng InitialPermissions
// Kiểm tra quyền đã tồn tại trong collection Permissions chưa
// Nếu chưa tồn tại thì thêm quyền đó vào collection Permissions
func (h *InitService) InitPermission() {
	for _, permission := range InitialPermissions {
		// Tìm quyền theo filter

		filter := bson.M{"name": permission.Name}
		_, err := h.PermissionCRUD.FindOne(context.TODO(), filter, nil)
		if err != nil && err != mongo.ErrNoDocuments {
			continue
		}

		// Nếu quyền chưa tồn tại thì thêm quyền vào collection
		if err == mongo.ErrNoDocuments {
			h.PermissionCRUD.InsertOne(context.TODO(), permission)
		}
	}
}

// Viết hàm InitRole để khởi tạo các vai trò mặc định theo nguyên tắc sau:
// Kiểm tra vai trò Administrator đã tồn tại chưa
// Nếu chưa tồn tại thì thêm vai trò Administrator vào collection Roles
// Sau đó, gán tất cả các quyền cho vai trò Administrator
func (h *InitService) InitRole() error {
	// Kiểm tra vai trò Administrator đã tồn tại chưa
	adminRole, err := h.RoleCRUD.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == nil {
		return errors.New("Role Administrator is already existed")
	}

	// Tạo vai trò Administrator
	newAdminRole := models.Role{
		Name:     "Administrator",
		Describe: "Vai trò quản trị hệ thống",
	}

	// Thêm vai trò vào collection
	adminRole, err = h.RoleCRUD.InsertOne(context.TODO(), newAdminRole)
	if err != nil {
		return errors.New("Failed to insert role Administrator")
	}

	// Lấy tất cả quyền
	permissions, err := h.PermissionCRUD.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return errors.New("Failed to get all permissions")
	}

	// Gán tất cả quyền cho vai trò Administrator
	for _, permission := range permissions {
		rolePermission := models.RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: permission.ID,
		}
		_, err = h.RolePermissionCRUD.InsertOne(context.TODO(), rolePermission)
		if err != nil {
			continue
		}
	}
	return nil
}

// Viết hàm kiểm tra các quyền của role Administrator, nếu thiếu quyền nào thì thêm vào
func (h *InitService) CheckPermissionForAdministrator() (err error) {
	// Tìm role theo tên
	role, err := h.RoleCRUD.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return h.InitRole()
	}

	// Chuyển đổi role từ bson.M về models.Role
	var modelRole models.Role
	bsonBytes, _ := bson.Marshal(role)
	err = bson.Unmarshal(bsonBytes, &modelRole)
	if err != nil {
		return errors.New("Failed to decode role")
	}

	// Lấy tất cả quyền
	permissions, err := h.PermissionCRUD.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return errors.New("Failed to get all permissions")
	}

	// duyệt qua danh sách các quyền
	for _, permissionData := range permissions {
		// decode permission từ bson.M về models.Permission
		var modelPermission models.Permission
		bsonBytes, _ := bson.Marshal(permissionData)
		err := bson.Unmarshal(bsonBytes, &modelPermission)
		if err != nil {
			fmt.Errorf("Failed to decode permission")
			continue
		}

		// Tìm quyền của role Administrator
		filter := bson.M{
			"roleId":       modelRole.ID,
			"permissionId": modelPermission.ID,
			"scope":        0,
		}

		// Tìm quyền của role Administrator
		_, err = h.RolePermissionCRUD.FindOne(context.TODO(), filter, nil)
		if err != nil && err != mongo.ErrNoDocuments {
			continue
		}
		if err == mongo.ErrNoDocuments {
			rolePermission := models.RolePermission{
				RoleID:       modelRole.ID,
				PermissionID: modelPermission.ID,
				Scope:        0,
			}
			_, err = h.RolePermissionCRUD.InsertOne(context.TODO(), rolePermission)

			if err != nil {
				fmt.Errorf("Failed to insert role permission: %v", err)
				continue
			}
		}
	}

	return nil
}

// Viết hàm set administator để gán quyền admin cho user
func (h *InitService) SetAdministrator(userID primitive.ObjectID) (result interface{}, err error) {
	// Tìm user theo ID
	user, err := h.UserCRUD.FindOneById(context.TODO(), userID)
	if err != nil {
		return nil, err
	}

	// Tìm role theo tên
	role, err := h.RoleCRUD.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("Role not found")
	}

	// Gán role cho user
	userRole := models.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}
	result, err = h.UserRoleCRUD.InsertOne(context.TODO(), userRole)
	if err != nil {
		return nil, err
	}

	return result, nil
}
