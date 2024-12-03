package services

import (
	"go.mongodb.org/mongo-driver/mongo"
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"
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
	{Name: "User.Read", 				Describe: "Quyền xem người dùng", 				Category: "User"},
	//{Name: "Permission.Create", 				Describe: "Quyền tạo vai trò", 					Category: "Permission"},
	{Name: "Permission.Read", 				Describe: "Quyền xem các quyền", 					Category: "Permission"},
	//{Name: "Permission.Update", 				Describe: "Quyền cập nhật vai trò", 			Category: "Permission"},
	//{Name: "Permission.Delete", 				Describe: "Quyền xóa vai trò", 					Category: "Permission"},
	{Name: "Organization.Create", 		Describe: "Quyền tạo tổ chức", 					Category: "Organization"},
	{Name: "Organization.Read", 		Describe: "Quyền xem tổ chức", 					Category: "Organization"},
	{Name: "Organization.Update", 		Describe: "Quyền cập nhật tổ chức", 			Category: "Organization"},
	{Name: "Organization.Delete", 		Describe: "Quyền xóa tổ chức", 					Category: "Organization"},
	{Name: "Role.Create", 				Describe: "Quyền tạo vai trò", 					Category: "Role"},
	{Name: "Role.Read", 				Describe: "Quyền xem vai trò", 					Category: "Role"},
	{Name: "Role.Update", 				Describe: "Quyền cập nhật vai trò", 			Category: "Role"},
	{Name: "Role.Delete", 				Describe: "Quyền xóa vai trò", 					Category: "Role"},
	{Name: "UserRole.Create", 				Describe: "Quyền tạo phân công vai trò", 					Category: "UserRole"},
	{Name: "UserRole.Read", 				Describe: "Quyền xem phân công vai trò", 					Category: "UserRole"},
	{Name: "UserRole.Update", 				Describe: "Quyền cập nhật phân công vai trò", 			Category: "UserRole"},
	{Name: "UserRole.Delete", 				Describe: "Quyền xóa phân công vai trò", 					Category: "UserRole"},
	{Name: "RolePermission.Create", 				Describe: "Quyền tạo phân quyền cho vai trò", 					Category: "RolePermission"},
	{Name: "RolePermission.Read", 					Describe: "Quyền xem phân quyền cho vai trò", 						Category: "RolePermission"},
	{Name: "RolePermission.Update", 				Describe: "Quyền cập nhật phân quyền cho vai trò", 			Category: "RolePermission"},
	{Name: "RolePermission.Delete", 				Describe: "Quyền xóa phân quyền cho vai trò", 					Category: "RolePermission"},
}
