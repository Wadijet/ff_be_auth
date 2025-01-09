package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitService định nghĩa các CRUD repository cho User, Permission và Role
type InitService struct {
	UserCRUD           RepositoryService
	PermissionCRUD     RepositoryService
	RoleCRUD           RepositoryService
	RolePermissionCRUD RepositoryService
	UserRoleCRUD       RepositoryService
}

// NewInitService khởi tạo các repository và trả về một đối tượng InitService
func NewInitService(c *config.Configuration, db *mongo.Client) *InitService {
	newService := new(InitService)
	newService.UserCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.PermissionCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newService.RoleCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newService.RolePermissionCRUD = *NewRepository(c, db, global.MongoDB_ColNames.RolePermissions)
	newService.UserRoleCRUD = *NewRepository(c, db, global.MongoDB_ColNames.UserRoles)
	return newService
}

var InitialPermissions = []models.Permission{

	{Name: "User.Read", Describe: "Quyền xem người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Block", Describe: "Quyền khóa người dùng", Group: "Auth", Category: "User"},
	{Name: "Permission.Read", Describe: "Quyền xem các quyền", Group: "Auth", Category: "Permission"},
	{Name: "Role.Create", Describe: "Quyền tạo vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Read", Describe: "Quyền xem vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Update", Describe: "Quyền cập nhật vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Delete", Describe: "Quyền xóa vai trò", Group: "Auth", Category: "Role"},
	{Name: "RolePermission.Create", Describe: "Quyền tạo phân quyền cho vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Read", Describe: "Quyền xem phân quyền cho vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Update", Describe: "Quyền cập nhật phân quyền cho vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Delete", Describe: "Quyền xóa phân quyền cho vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "UserRole.Create", Describe: "Quyền tạo phân công vai trò", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Read", Describe: "Quyền xem phân công vai trò", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Update", Describe: "Quyền cập nhật phân công vai trò", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Delete", Describe: "Quyền xóa phân công vai trò", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Delete", Describe: "Quyền xóa phân công vai trò", Group: "Auth", Category: "UserRole"},
	{Name: "Agent.Read", Describe: "Quyền xem trợ lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Create", Describe: "Quyền tạo trợ lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Update", Describe: "Quyền cập nhật trợ lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Delete", Describe: "Quyền xóa trợ lý", Group: "Auth", Category: "Agent"},
	{Name: "AccessToken.Read", Describe: "Quyền xem Access token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Create", Describe: "Quyền tạo Access token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Update", Describe: "Quyền cập nhật Access token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Delete", Describe: "Quyền xóa Access token", Group: "Pancake", Category: "AccessToken"},
	{Name: "FbPage.Read", Describe: "Quyền xem trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Create", Describe: "Quyền tạo trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Update", Describe: "Quyền cập nhật trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Delete", Describe: "Quyền xóa trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbConversation.Read", Describe: "Quyền xem cuộc trò chuyện trên Facebook", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Create", Describe: "Quyền tạo cuộc trò chuyện trên Facebook", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Update", Describe: "Quyền cập nhật cuộc trò chuyện trên Facebook", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Delete", Describe: "Quyền xóa cuộc trò chuyện trên Facebook", Group: "Pancake", Category: "FbConversation"},
}

// Viết hàm InitPermission để khởi tạo các quyền mặc định theo nguyên tắc sau:
// Duyệt tất cả các quyền trong mảng InitialPermissions
// Kiểm tra quyền đã tồn tại trong collection Permissions chưa
// Nếu chưa tồn tại thì thêm quyền đó vào collection Permissions
func (h *InitService) InitPermission() {
	for _, permission := range InitialPermissions {
		// Tạo filter để tìm kiếm quyền theo tên
		filter := map[string]interface{}{"name": permission.Name}
		// Tìm quyền theo filter
		result, _ := h.PermissionCRUD.FindOne(context.TODO(), filter, nil)
		// Nếu quyền chưa tồn tại thì thêm quyền vào collection
		if result == nil {
			h.PermissionCRUD.InsertOne(context.TODO(), permission)
		}
	}
}

// Viết hàm InitRole để khởi tạo các vai trò mặc định theo nguyên tắc sau:
// Kiểm tra vai trò Administrator đã tồn tại chưa
// Nếu chưa tồn tại thì thêm vai trò Administrator vào collection Roles
// Sau đó, gán tất cả các quyền cho vai trò Administrator
func (h *InitService) InitRole() (err error) {

	// Kiểm tra vai trò Administrator đã tồn tại chưa
	filter := map[string]interface{}{"name": "Administrator"}
	findResult, err := h.RoleCRUD.FindOne(context.TODO(), filter, nil)

	if findResult != nil {
		return errors.New("Role Administrator is already existed")
	}

	// Nếu vai trò chưa tồn tại thì thêm vai trò vào collection
	adminRole := models.Role{
		Name:     "Administrator",
		Describe: "Vai trò quản trị hệ thống",
	}

	// Thêm vai trò vào collection
	resultInsertRole, err := h.RoleCRUD.InsertOne(nil, adminRole)
	if err != nil {
		return errors.New("Failed to insert role Administrator")
	}
	insertedRoleID := resultInsertRole.InsertedID.(primitive.ObjectID)

	// Lấy tất cả quyền
	permissions, err := h.PermissionCRUD.FindAll(nil, nil, nil)
	if err != nil {
		return errors.New("Failed to get all permissions")
	}

	// Gán tất cả quyền cho vai trò Administrator
	for _, permissionData := range permissions {
		// decode permission từ bson.M về models.Permission
		var modelPermission models.Permission
		bsonBytes, _ := bson.Marshal(permissionData)
		err := bson.Unmarshal(bsonBytes, &modelPermission)
		if err != nil {
			fmt.Errorf("Failed to decode permission")
			continue
		}

		rolePermission := models.RolePermission{
			RoleID:       insertedRoleID,
			PermissionID: modelPermission.ID,
		}
		_, err = h.RolePermissionCRUD.InsertOne(context.TODO(), rolePermission)
		if err != nil {
			fmt.Errorf("Failed to insert role permission: %v", err)
			continue
		}
	}
	return nil
}

// Viết hàm kiểm tra các quyền của role Administrator, nếu thiếu quyền nào thì thêm vào
func (h *InitService) CheckPermissionForAdministrator() (err error) {
	// Tìm role theo tên
	filter := map[string]interface{}{"name": "Administrator"}
	role, err := h.RoleCRUD.FindOne(context.TODO(), filter, nil)
	if role == nil {
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
	permissions, err := h.PermissionCRUD.FindAll(nil, nil, nil)
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
		filter := bson.D{
			{Key: "roleId", Value: modelRole.ID},
			{Key: "permissionId", Value: modelPermission.ID},
			{Key: "scope", Value: 0},
		}

		rolePermission, err := h.RolePermissionCRUD.FindOne(context.TODO(), filter, nil)
		if rolePermission == nil {
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
func (h *InitService) SetAdministrator(userID string) (result interface{}, err error) {
	// Tìm user theo ID
	user, err := h.UserCRUD.FindOneById(context.TODO(), utility.String2ObjectID(userID), nil)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User not found")
	}

	// Tìm role theo tên
	role, err := h.RoleCRUD.FindOne(context.TODO(), map[string]interface{}{"name": "Administrator"}, nil)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("Role not found")
	}

	// Gán role cho user
	roleID := role["_id"].(primitive.ObjectID)
	userRole := models.UserRole{
		UserID: utility.String2ObjectID(userID),
		RoleID: roleID,
	}
	insertResult, err := h.UserRoleCRUD.InsertOne(context.TODO(), userRole)
	if err != nil {
		return nil, err
	}

	return insertResult, nil
}
