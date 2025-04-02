package services

import (
	"context"
	"fmt"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InitService là cấu trúc chứa các phương thức khởi tạo dữ liệu ban đầu
type InitService struct {
	userService           *UserService
	roleService           *RoleService
	permissionService     *PermissionService
	rolePermissionService *RolePermissionService
	userRoleService       *UserRoleService
}

// NewInitService tạo mới InitService
func NewInitService() *InitService {
	return &InitService{
		userService:           NewUserService(),
		roleService:           NewRoleService(),
		permissionService:     NewPermissionService(),
		rolePermissionService: NewRolePermissionService(),
		userRoleService:       NewUserRoleService(),
	}
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

// InitPermission để khởi tạo các quyền mặc định
func (h *InitService) InitPermission() {
	for _, permission := range InitialPermissions {
		filter := bson.M{"name": permission.Name}
		_, err := h.permissionService.FindOne(context.TODO(), filter, nil)
		if err != nil && err != utility.ErrNotFound {
			continue
		}

		if err == utility.ErrNotFound {
			h.permissionService.InsertOne(context.TODO(), permission)
		}
	}
}

// InitRole để khởi tạo các vai trò mặc định
func (h *InitService) InitRole() error {
	adminRole, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != utility.ErrNotFound {
		return err
	}
	if err == nil {
		return utility.ErrInvalidInput
	}

	newAdminRole := models.Role{
		Name:     "Administrator",
		Describe: "Vai trò quản trị hệ thống",
	}

	adminRole, err = h.roleService.InsertOne(context.TODO(), newAdminRole)
	if err != nil {
		return utility.ErrInvalidInput
	}

	permissions, err := h.permissionService.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return utility.ErrInvalidInput
	}

	for _, permission := range permissions {
		rolePermission := models.RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: permission.ID,
		}
		_, err = h.rolePermissionService.InsertOne(context.TODO(), rolePermission)
		if err != nil {
			continue
		}
	}
	return nil
}

// CheckPermissionForAdministrator kiểm tra và cập nhật quyền cho Administrator
func (h *InitService) CheckPermissionForAdministrator() (err error) {
	role, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != utility.ErrNotFound {
		return err
	}
	if err == utility.ErrNotFound {
		return h.InitRole()
	}

	var modelRole models.Role
	bsonBytes, _ := bson.Marshal(role)
	err = bson.Unmarshal(bsonBytes, &modelRole)
	if err != nil {
		return utility.ErrInvalidFormat
	}

	permissions, err := h.permissionService.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return utility.ErrInvalidInput
	}

	for _, permissionData := range permissions {
		var modelPermission models.Permission
		bsonBytes, _ := bson.Marshal(permissionData)
		err := bson.Unmarshal(bsonBytes, &modelPermission)
		if err != nil {
			fmt.Errorf("Failed to decode permission")
			continue
		}

		filter := bson.M{
			"roleId":       modelRole.ID,
			"permissionId": modelPermission.ID,
			"scope":        0,
		}

		_, err = h.rolePermissionService.FindOne(context.TODO(), filter, nil)
		if err != nil && err != utility.ErrNotFound {
			continue
		}
		if err == utility.ErrNotFound {
			rolePermission := models.RolePermission{
				RoleID:       modelRole.ID,
				PermissionID: modelPermission.ID,
				Scope:        0,
			}
			_, err = h.rolePermissionService.InsertOne(context.TODO(), rolePermission)
			if err != nil {
				fmt.Errorf("Failed to insert role permission: %v", err)
				continue
			}
		}
	}

	return nil
}

// SetAdministrator gán quyền admin cho user
func (h *InitService) SetAdministrator(userID primitive.ObjectID) (result interface{}, err error) {
	user, err := h.userService.FindOneById(context.TODO(), userID)
	if err != nil {
		return nil, err
	}

	role, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != utility.ErrNotFound {
		return nil, err
	}
	if err == utility.ErrNotFound {
		err = h.InitRole()
		if err != nil {
			return nil, err
		}
		role, err = h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
		if err != nil {
			return nil, err
		}
	}

	var modelRole models.Role
	bsonBytes, _ := bson.Marshal(role)
	err = bson.Unmarshal(bsonBytes, &modelRole)
	if err != nil {
		return nil, utility.ErrInvalidFormat
	}

	userRole := models.UserRole{
		UserID: user.ID,
		RoleID: modelRole.ID,
	}
	result, err = h.userRoleService.InsertOne(context.TODO(), userRole)
	if err != nil {
		return nil, err
	}

	return result, nil
}
