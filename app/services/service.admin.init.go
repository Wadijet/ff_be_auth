// Package services chứa các service xử lý logic nghiệp vụ của ứng dụng
package services

import (
	"context"
	"fmt"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InitService là cấu trúc chứa các phương thức khởi tạo dữ liệu ban đầu cho hệ thống
// Bao gồm khởi tạo người dùng, vai trò, quyền và các quan hệ giữa chúng
type InitService struct {
	userService           *UserService           // Service xử lý người dùng
	roleService           *RoleService           // Service xử lý vai trò
	permissionService     *PermissionService     // Service xử lý quyền
	rolePermissionService *RolePermissionService // Service xử lý quan hệ vai trò-quyền
	userRoleService       *UserRoleService       // Service xử lý quan hệ người dùng-vai trò
}

// NewInitService tạo mới một đối tượng InitService
// Khởi tạo các service con cần thiết để xử lý các tác vụ liên quan
func NewInitService() *InitService {
	return &InitService{
		userService:           NewUserService(),
		roleService:           NewRoleService(),
		permissionService:     NewPermissionService(),
		rolePermissionService: NewRolePermissionService(),
		userRoleService:       NewUserRoleService(),
	}
}

// InitialPermissions định nghĩa danh sách các quyền mặc định của hệ thống
// Được chia thành các module: Auth (Xác thực) và Pancake (Quản lý trang Facebook)
var InitialPermissions = []models.Permission{
	// ====================================  AUTH MODULE =============================================
	// Quản lý người dùng: Thêm, xem, sửa, xóa, khóa và phân quyền
	{Name: "User.Insert", Describe: "Quyền tạo người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Read", Describe: "Quyền xem danh sách người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Update", Describe: "Quyền cập nhật thông tin người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Delete", Describe: "Quyền xóa người dùng", Group: "Auth", Category: "User"},
	{Name: "User.Block", Describe: "Quyền khóa/mở khóa người dùng", Group: "Auth", Category: "User"},
	{Name: "User.SetRole", Describe: "Quyền phân quyền cho người dùng", Group: "Auth", Category: "User"},

	// Quản lý vai trò: Thêm, xem, sửa, xóa vai trò
	{Name: "Role.Insert", Describe: "Quyền tạo vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Read", Describe: "Quyền xem danh sách vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Update", Describe: "Quyền cập nhật vai trò", Group: "Auth", Category: "Role"},
	{Name: "Role.Delete", Describe: "Quyền xóa vai trò", Group: "Auth", Category: "Role"},

	// Quản lý quyền: Thêm, xem, sửa, xóa quyền
	{Name: "Permission.Insert", Describe: "Quyền tạo quyền", Group: "Auth", Category: "Permission"},
	{Name: "Permission.Read", Describe: "Quyền xem danh sách quyền", Group: "Auth", Category: "Permission"},
	{Name: "Permission.Update", Describe: "Quyền cập nhật quyền", Group: "Auth", Category: "Permission"},
	{Name: "Permission.Delete", Describe: "Quyền xóa quyền", Group: "Auth", Category: "Permission"},

	// Quản lý phân quyền cho vai trò: Thêm, xem, sửa, xóa phân quyền
	{Name: "RolePermission.Insert", Describe: "Quyền tạo phân quyền cho vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Read", Describe: "Quyền xem phân quyền của vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Update", Describe: "Quyền cập nhật phân quyền của vai trò", Group: "Auth", Category: "RolePermission"},
	{Name: "RolePermission.Delete", Describe: "Quyền xóa phân quyền của vai trò", Group: "Auth", Category: "RolePermission"},

	// Quản lý phân vai trò cho người dùng: Thêm, xem, sửa, xóa phân vai trò
	{Name: "UserRole.Insert", Describe: "Quyền phân công vai trò cho người dùng", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Read", Describe: "Quyền xem vai trò của người dùng", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Update", Describe: "Quyền cập nhật vai trò của người dùng", Group: "Auth", Category: "UserRole"},
	{Name: "UserRole.Delete", Describe: "Quyền xóa vai trò của người dùng", Group: "Auth", Category: "UserRole"},

	// Quản lý đại lý: Thêm, xem, sửa, xóa và kiểm tra trạng thái
	{Name: "Agent.Insert", Describe: "Quyền tạo đại lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Read", Describe: "Quyền xem danh sách đại lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Update", Describe: "Quyền cập nhật thông tin đại lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.Delete", Describe: "Quyền xóa đại lý", Group: "Auth", Category: "Agent"},
	{Name: "Agent.CheckIn", Describe: "Quyền kiểm tra trạng thái đại lý", Group: "Auth", Category: "Agent"},

	// ==================================== PANCAKE MODULE ===========================================
	// Quản lý token truy cập: Thêm, xem, sửa, xóa token
	{Name: "AccessToken.Insert", Describe: "Quyền tạo token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Read", Describe: "Quyền xem danh sách token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Update", Describe: "Quyền cập nhật token", Group: "Pancake", Category: "AccessToken"},
	{Name: "AccessToken.Delete", Describe: "Quyền xóa token", Group: "Pancake", Category: "AccessToken"},

	// Quản lý trang Facebook: Thêm, xem, sửa, xóa và cập nhật token
	{Name: "FbPage.Insert", Describe: "Quyền tạo trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Read", Describe: "Quyền xem danh sách trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Update", Describe: "Quyền cập nhật thông tin trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.Delete", Describe: "Quyền xóa trang Facebook", Group: "Pancake", Category: "FbPage"},
	{Name: "FbPage.UpdateToken", Describe: "Quyền cập nhật token trang Facebook", Group: "Pancake", Category: "FbPage"},

	// Quản lý cuộc trò chuyện Facebook: Thêm, xem, sửa, xóa
	{Name: "FbConversation.Insert", Describe: "Quyền tạo cuộc trò chuyện", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Read", Describe: "Quyền xem danh sách cuộc trò chuyện", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Update", Describe: "Quyền cập nhật cuộc trò chuyện", Group: "Pancake", Category: "FbConversation"},
	{Name: "FbConversation.Delete", Describe: "Quyền xóa cuộc trò chuyện", Group: "Pancake", Category: "FbConversation"},

	// Quản lý tin nhắn Facebook: Thêm, xem, sửa, xóa
	{Name: "FbMessage.Insert", Describe: "Quyền tạo tin nhắn", Group: "Pancake", Category: "FbMessage"},
	{Name: "FbMessage.Read", Describe: "Quyền xem danh sách tin nhắn", Group: "Pancake", Category: "FbMessage"},
	{Name: "FbMessage.Update", Describe: "Quyền cập nhật tin nhắn", Group: "Pancake", Category: "FbMessage"},
	{Name: "FbMessage.Delete", Describe: "Quyền xóa tin nhắn", Group: "Pancake", Category: "FbMessage"},

	// Quản lý bài viết Facebook: Thêm, xem, sửa, xóa
	{Name: "FbPost.Insert", Describe: "Quyền tạo bài viết", Group: "Pancake", Category: "FbPost"},
	{Name: "FbPost.Read", Describe: "Quyền xem danh sách bài viết", Group: "Pancake", Category: "FbPost"},
	{Name: "FbPost.Update", Describe: "Quyền cập nhật bài viết", Group: "Pancake", Category: "FbPost"},
	{Name: "FbPost.Delete", Describe: "Quyền xóa bài viết", Group: "Pancake", Category: "FbPost"},

	// Quản lý đơn hàng Pancake: Thêm, xem, sửa, xóa
	{Name: "PcOrder.Insert", Describe: "Quyền tạo đơn hàng", Group: "Pancake", Category: "PcOrder"},
	{Name: "PcOrder.Read", Describe: "Quyền xem danh sách đơn hàng", Group: "Pancake", Category: "PcOrder"},
	{Name: "PcOrder.Update", Describe: "Quyền cập nhật đơn hàng", Group: "Pancake", Category: "PcOrder"},
	{Name: "PcOrder.Delete", Describe: "Quyền xóa đơn hàng", Group: "Pancake", Category: "PcOrder"},
}

// InitPermission khởi tạo các quyền mặc định cho hệ thống
// Chỉ tạo mới các quyền chưa tồn tại trong database
func (h *InitService) InitPermission() {
	// Duyệt qua danh sách quyền mặc định
	for _, permission := range InitialPermissions {
		// Kiểm tra quyền đã tồn tại chưa
		filter := bson.M{"name": permission.Name}
		_, err := h.permissionService.FindOne(context.TODO(), filter, nil)

		// Bỏ qua nếu có lỗi khác ErrNotFound
		if err != nil && err != utility.ErrNotFound {
			continue
		}

		// Tạo mới quyền nếu chưa tồn tại
		if err == utility.ErrNotFound {
			h.permissionService.InsertOne(context.TODO(), permission)
		}
	}
}

// InitRole khởi tạo vai trò Administrator mặc định
// Tạo vai trò và gán tất cả các quyền cho vai trò này
func (h *InitService) InitRole() error {
	// Kiểm tra vai trò Administrator đã tồn tại chưa
	adminRole, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != utility.ErrNotFound {
		return err
	}
	if err == nil {
		return utility.ErrInvalidInput
	}

	// Tạo mới vai trò Administrator
	newAdminRole := models.Role{
		Name:     "Administrator",
		Describe: "Vai trò quản trị hệ thống",
	}

	// Lưu vai trò vào database
	adminRole, err = h.roleService.InsertOne(context.TODO(), newAdminRole)
	if err != nil {
		return utility.ErrInvalidInput
	}

	// Lấy danh sách tất cả các quyền
	permissions, err := h.permissionService.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return utility.ErrInvalidInput
	}

	// Gán tất cả quyền cho vai trò Administrator
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

// CheckPermissionForAdministrator kiểm tra và cập nhật quyền cho vai trò Administrator
// Đảm bảo vai trò Administrator có đầy đủ tất cả các quyền trong hệ thống
func (h *InitService) CheckPermissionForAdministrator() (err error) {
	// Kiểm tra vai trò Administrator có tồn tại không
	role, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != utility.ErrNotFound {
		return err
	}
	// Nếu chưa có vai trò Administrator, tạo mới
	if err == utility.ErrNotFound {
		return h.InitRole()
	}

	// Chuyển đổi dữ liệu sang model
	var modelRole models.Role
	bsonBytes, _ := bson.Marshal(role)
	err = bson.Unmarshal(bsonBytes, &modelRole)
	if err != nil {
		return utility.ErrInvalidFormat
	}

	// Lấy danh sách tất cả các quyền
	permissions, err := h.permissionService.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return utility.ErrInvalidInput
	}

	// Kiểm tra và cập nhật từng quyền cho vai trò Administrator
	for _, permissionData := range permissions {
		var modelPermission models.Permission
		bsonBytes, _ := bson.Marshal(permissionData)
		err := bson.Unmarshal(bsonBytes, &modelPermission)
		if err != nil {
			fmt.Errorf("Failed to decode permission")
			continue
		}

		// Kiểm tra quyền đã được gán chưa
		filter := bson.M{
			"roleId":       modelRole.ID,
			"permissionId": modelPermission.ID,
			"scope":        0,
		}

		_, err = h.rolePermissionService.FindOne(context.TODO(), filter, nil)
		if err != nil && err != utility.ErrNotFound {
			continue
		}

		// Nếu chưa có quyền, thêm mới
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

// SetAdministrator gán quyền Administrator cho một người dùng
// Trả về lỗi nếu người dùng không tồn tại hoặc đã có quyền Administrator
func (h *InitService) SetAdministrator(userID primitive.ObjectID) (result interface{}, err error) {
	// Kiểm tra user có tồn tại không
	user, err := h.userService.FindOneById(context.TODO(), userID)
	if err != nil {
		return nil, err
	}

	// Kiểm tra role Administrator có tồn tại không
	role, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != utility.ErrNotFound {
		return nil, err
	}

	// Nếu chưa có role Administrator, tạo mới
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

	// Kiểm tra userRole đã tồn tại chưa
	_, err = h.userRoleService.FindOne(context.TODO(), bson.M{"userId": user.ID, "roleId": role.ID}, nil)
	// Kiểm tra nếu userRole đã tồn tại
	if err == nil {
		// Nếu không có lỗi, tức là đã tìm thấy userRole, trả về lỗi đã định nghĩa
		return nil, utility.ErrUserAlreadyAdmin
	}

	// Xử lý các lỗi khác ngoài ErrNotFound
	if err != utility.ErrNotFound {
		return nil, err
	}

	// Nếu userRole chưa tồn tại (err == utility.ErrNotFound), tạo mới
	userRole := models.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}
	result, err = h.userRoleService.InsertOne(context.TODO(), userRole)
	if err != nil {
		return nil, err
	}

	return result, nil
}
