package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"

	"go.mongodb.org/mongo-driver/mongo"
)

// InitService định nghĩa các CRUD repository cho User, Permission và Role
type InitService struct {
	UserCRUD       Repository
	PermissionCRUD Repository
	RoleCRUD       Repository
}

// NewInitService khởi tạo các repository và trả về một đối tượng InitService
func NewInitService(c *config.Configuration, db *mongo.Client) *InitService {
	newService := new(InitService)
	newService.UserCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.PermissionCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newService.RoleCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
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
		result, _ := h.PermissionCRUD.FindOne(nil, filter, nil)
		// Nếu quyền chưa tồn tại thì thêm quyền vào collection
		if result == nil {
			h.PermissionCRUD.InsertOne(nil, permission)
		}
	}
}

// Viết hàm InitRole để khởi tạo các vai trò mặc định theo nguyên tắc sau:
// Kiểm tra vai trò Administrator đã tồn tại chưa
// Nếu chưa tồn tại thì thêm vai trò Administrator vào collection Roles
// Sau đó, gán tất cả các quyền cho vai trò Administrator
func (h *InitService) InitRole() (err error) {
	// Tạo filter để tìm kiếm vai trò theo tên
	filter := map[string]interface{}{"name": "Administrator"}
	// Tìm vai trò theo filter
	result, _ := h.RoleCRUD.FindOne(nil, filter, nil)
	// Nếu vai trò chưa tồn tại thì thêm vai trò vào collection
	if result != mongo.ErrNoDocuments {
		return errors.New("Role Administrator is already existed")

	}

	// Tạo vai trò Administrator
	adminRole := models.Role{
		Name:     "Administrator",
		Describe: "Vai trò quản trị hệ thống",
	}
	// Thêm vai trò vào collection
	resultInsertRole, err := h.RoleCRUD.InsertOne(nil, adminRole)
	if err != nil {
		return errors.New("Failed to insert role Administrator")
	}
	
	insertedRoldID := resultInsertRole.InsertedID

	// Lấy tất cả quyền
	permissions, err := h.PermissionCRUD.FindAll(nil, nil, nil)
	if err != nil {
		return errors.New("Failed to get all permissions")
	}

	// Gán tất cả quyền cho vai trò Administrator
	for _, permission := range permissions {
		rolePermission := models.RolePermission{
			RoleID:      insertedRoldID,
			PermissionID: permission.ID,
		}
		_, err = h.RoleCRUD.InsertOne(nil, rolePermission)
	return
}
