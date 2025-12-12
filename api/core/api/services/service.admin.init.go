// Package services chứa các service xử lý logic nghiệp vụ của ứng dụng
package services

import (
	"context"
	"fmt"
	"time"

	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/utility"

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
	organizationService   *OrganizationService   // Service xử lý tổ chức
}

// NewInitService tạo mới một đối tượng InitService
// Khởi tạo các service con cần thiết để xử lý các tác vụ liên quan
// Returns:
//   - *InitService: Instance mới của InitService
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewInitService() (*InitService, error) {
	// Khởi tạo các services
	userService, err := NewUserService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user service: %v", err)
	}

	roleService, err := NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}

	permissionService, err := NewPermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create permission service: %v", err)
	}

	rolePermissionService, err := NewRolePermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role permission service: %v", err)
	}

	userRoleService, err := NewUserRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user role service: %v", err)
	}

	organizationService, err := NewOrganizationService()
	if err != nil {
		return nil, fmt.Errorf("failed to create organization service: %v", err)
	}

	return &InitService{
		userService:           userService,
		roleService:           roleService,
		permissionService:     permissionService,
		rolePermissionService: rolePermissionService,
		userRoleService:       userRoleService,
		organizationService:   organizationService,
	}, nil
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

	// Quản lý tổ chức: Thêm, xem, sửa, xóa
	{Name: "Organization.Insert", Describe: "Quyền tạo tổ chức", Group: "Auth", Category: "Organization"},
	{Name: "Organization.Read", Describe: "Quyền xem danh sách tổ chức", Group: "Auth", Category: "Organization"},
	{Name: "Organization.Update", Describe: "Quyền cập nhật tổ chức", Group: "Auth", Category: "Organization"},
	{Name: "Organization.Delete", Describe: "Quyền xóa tổ chức", Group: "Auth", Category: "Organization"},

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
	{Name: "Agent.CheckOut", Describe: "Quyền kiểm tra trạng thái đại lý", Group: "Auth", Category: "Agent"},

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
// Returns:
//   - error: Lỗi nếu có trong quá trình khởi tạo
func (h *InitService) InitPermission() error {
	// Duyệt qua danh sách quyền mặc định
	for _, permission := range InitialPermissions {
		// Kiểm tra quyền đã tồn tại chưa
		filter := bson.M{"name": permission.Name}
		_, err := h.permissionService.FindOne(context.TODO(), filter, nil)

		// Bỏ qua nếu có lỗi khác ErrNotFound
		if err != nil && err != common.ErrNotFound {
			continue
		}

		// Tạo mới quyền nếu chưa tồn tại
		if err == common.ErrNotFound {
			_, err = h.permissionService.InsertOne(context.TODO(), permission)
			if err != nil {
				return fmt.Errorf("failed to insert permission %s: %v", permission.Name, err)
			}
		}
	}
	return nil
}

// InitRootOrganization khởi tạo Organization System (Level -1)
// System organization là tổ chức cấp cao nhất, chứa Administrator, không có parent, không thể xóa
// System thay thế ROOT_GROUP cũ
// Returns:
//   - error: Lỗi nếu có trong quá trình khởi tạo
func (h *InitService) InitRootOrganization() error {
	// Kiểm tra System Organization đã tồn tại chưa
	systemFilter := bson.M{
		"type":  models.OrganizationTypeSystem,
		"level": -1,
		"code":  "SYSTEM",
	}
	_, err := h.organizationService.FindOne(context.TODO(), systemFilter, nil)
	if err != nil && err != common.ErrNotFound {
		return fmt.Errorf("failed to check system organization: %v", err)
	}

	// Nếu đã tồn tại, không cần tạo mới
	if err == nil {
		return nil
	}

	// Tạo mới System Organization
	systemOrgModel := models.Organization{
		Name:     "Hệ Thống",
		Code:     "SYSTEM",
		Type:     models.OrganizationTypeSystem,
		ParentID: nil, // System không có parent
		Path:     "/system",
		Level:    -1,
		IsActive: true,
	}

	_, err = h.organizationService.InsertOne(context.TODO(), systemOrgModel)
	if err != nil {
		return fmt.Errorf("failed to create system organization: %v", err)
	}

	return nil
}

// GetRootOrganization lấy System Organization (Level -1) - tổ chức cấp cao nhất
// Returns:
//   - *models.Organization: System Organization
//   - error: Lỗi nếu có
func (h *InitService) GetRootOrganization() (*models.Organization, error) {
	filter := bson.M{
		"type":  models.OrganizationTypeSystem,
		"level": -1,
		"code":  "SYSTEM",
	}
	org, err := h.organizationService.FindOne(context.TODO(), filter, nil)
	if err != nil {
		return nil, fmt.Errorf("system organization not found: %v", err)
	}

	var modelOrg models.Organization
	bsonBytes, _ := bson.Marshal(org)
	err = bson.Unmarshal(bsonBytes, &modelOrg)
	if err != nil {
		return nil, common.ErrInvalidFormat
	}

	return &modelOrg, nil
}

// InitRole khởi tạo vai trò Administrator mặc định
// Tạo vai trò và gán tất cả các quyền cho vai trò này
// Role Administrator phải thuộc System Organization (Level -1)
func (h *InitService) InitRole() error {
	// Lấy System Organization (Level -1)
	rootOrg, err := h.GetRootOrganization()
	if err != nil {
		return fmt.Errorf("failed to get system organization: %v", err)
	}

	// Kiểm tra vai trò Administrator đã tồn tại chưa
	adminRole, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != common.ErrNotFound {
		return err
	}

	var modelRole models.Role
	roleExists := false

	if err == nil {
		// Nếu đã tồn tại, kiểm tra OrganizationID
		bsonBytes, _ := bson.Marshal(adminRole)
		err = bson.Unmarshal(bsonBytes, &modelRole)
		if err == nil {
			roleExists = true
			// Nếu chưa có OrganizationID, cập nhật
			if modelRole.OrganizationID.IsZero() {
				updateData := bson.M{
					"$set": bson.M{
						"organizationId": rootOrg.ID,
					},
				}
				_, err = h.roleService.UpdateOne(context.TODO(), bson.M{"_id": modelRole.ID}, updateData, nil)
				if err != nil {
					return fmt.Errorf("failed to update administrator role with organization: %v", err)
				}
			}
		}
	}

	// Nếu chưa tồn tại, tạo mới vai trò Administrator với OrganizationID
	if !roleExists {
		newAdminRole := models.Role{
			Name:           "Administrator",
			Describe:       "Vai trò quản trị hệ thống",
			OrganizationID: rootOrg.ID, // Gán vào Organization Root
		}

		// Lưu vai trò vào database
		adminRole, err = h.roleService.InsertOne(context.TODO(), newAdminRole)
		if err != nil {
			return fmt.Errorf("failed to create administrator role: %v", err)
		}

		// Chuyển đổi sang model để sử dụng
		bsonBytes, _ := bson.Marshal(adminRole)
		err = bson.Unmarshal(bsonBytes, &modelRole)
		if err != nil {
			return fmt.Errorf("failed to decode administrator role: %v", err)
		}
	}

	// Đảm bảo role Administrator có đầy đủ tất cả permissions
	// Lấy danh sách tất cả các quyền
	permissions, err := h.permissionService.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return fmt.Errorf("failed to get permissions: %v", err)
	}

	// Gán tất cả quyền cho vai trò Administrator với Scope = 1 (Tổ chức đó và tất cả các tổ chức con)
	for _, permissionData := range permissions {
		var modelPermission models.Permission
		bsonBytes, _ := bson.Marshal(permissionData)
		err := bson.Unmarshal(bsonBytes, &modelPermission)
		if err != nil {
			continue // Bỏ qua permission không decode được
		}

		// Kiểm tra quyền đã được gán chưa
		filter := bson.M{
			"roleId":       modelRole.ID,
			"permissionId": modelPermission.ID,
		}

		existingRP, err := h.rolePermissionService.FindOne(context.TODO(), filter, nil)
		if err != nil && err != common.ErrNotFound {
			continue // Bỏ qua nếu có lỗi khác ErrNotFound
		}

		// Nếu chưa có quyền, thêm mới với Scope = 1 (Tổ chức đó và tất cả các tổ chức con)
		if err == common.ErrNotFound {
			rolePermission := models.RolePermission{
				RoleID:       modelRole.ID,
				PermissionID: modelPermission.ID,
				Scope:        1, // Scope = 1: Tổ chức đó và tất cả các tổ chức con - Vì thuộc Root, sẽ xem tất cả
			}
			_, err = h.rolePermissionService.InsertOne(context.TODO(), rolePermission)
			if err != nil {
				continue // Bỏ qua nếu insert thất bại
			}
		} else {
			// Nếu đã có, kiểm tra scope - nếu là 0 thì cập nhật thành 1 (để admin có quyền xem tất cả)
			var existingModelRP models.RolePermission
			bsonBytes, _ := bson.Marshal(existingRP)
			err = bson.Unmarshal(bsonBytes, &existingModelRP)
			if err == nil && existingModelRP.Scope == 0 {
				// Cập nhật scope từ 0 → 1 (chỉ tổ chức → tổ chức + các tổ chức con)
				updateData := bson.M{
					"$set": bson.M{
						"scope": 1,
					},
				}
				_, err = h.rolePermissionService.UpdateOne(context.TODO(), bson.M{"_id": existingModelRP.ID}, updateData, nil)
				if err != nil {
					// Log error nhưng tiếp tục với permission tiếp theo
					continue
				}
			}
		}
	}

	return nil
}

// CheckPermissionForAdministrator kiểm tra và cập nhật quyền cho vai trò Administrator
// Đảm bảo vai trò Administrator có đầy đủ tất cả các quyền trong hệ thống
func (h *InitService) CheckPermissionForAdministrator() (err error) {
	// Kiểm tra vai trò Administrator có tồn tại không
	role, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil && err != common.ErrNotFound {
		return err
	}
	// Nếu chưa có vai trò Administrator, tạo mới
	if err == common.ErrNotFound {
		return h.InitRole()
	}

	// Chuyển đổi dữ liệu sang model
	var modelRole models.Role
	bsonBytes, _ := bson.Marshal(role)
	err = bson.Unmarshal(bsonBytes, &modelRole)
	if err != nil {
		return common.ErrInvalidFormat
	}

	// Lấy danh sách tất cả các quyền
	permissions, err := h.permissionService.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return common.ErrInvalidInput
	}

	// Kiểm tra và cập nhật từng quyền cho vai trò Administrator
	for _, permissionData := range permissions {
		var modelPermission models.Permission
		bsonBytes, _ := bson.Marshal(permissionData)
		err := bson.Unmarshal(bsonBytes, &modelPermission)
		if err != nil {
			// Log error nhưng tiếp tục với permission tiếp theo
			_ = fmt.Errorf("failed to decode permission: %v", err)
			continue
		}

		// Kiểm tra quyền đã được gán chưa (không filter scope)
		filter := bson.M{
			"roleId":       modelRole.ID,
			"permissionId": modelPermission.ID,
		}

		existingRP, err := h.rolePermissionService.FindOne(context.TODO(), filter, nil)
		if err != nil && err != common.ErrNotFound {
			continue
		}

		// Nếu chưa có quyền, thêm mới với Scope = 1 (Tổ chức đó và tất cả các tổ chức con)
		if err == common.ErrNotFound {
			rolePermission := models.RolePermission{
				RoleID:       modelRole.ID,
				PermissionID: modelPermission.ID,
				Scope:        1, // Scope = 1: Tổ chức đó và tất cả các tổ chức con - Vì thuộc Root, sẽ xem tất cả
			}
			_, err = h.rolePermissionService.InsertOne(context.TODO(), rolePermission)
			if err != nil {
				// Log error nhưng tiếp tục với permission tiếp theo
				_ = fmt.Errorf("failed to insert role permission: %v", err)
				continue
			}
		} else {
			// Nếu đã có, kiểm tra scope - nếu là 0 thì cập nhật thành 1 (để admin có quyền xem tất cả)
			var existingModelRP models.RolePermission
			bsonBytes, _ := bson.Marshal(existingRP)
			err = bson.Unmarshal(bsonBytes, &existingModelRP)
			if err == nil && existingModelRP.Scope == 0 {
				// Cập nhật scope từ 0 → 1 (chỉ tổ chức → tổ chức + các tổ chức con)
				updateData := bson.M{
					"$set": bson.M{
						"scope": 1,
					},
				}
				_, err = h.rolePermissionService.UpdateOne(context.TODO(), bson.M{"_id": existingModelRP.ID}, updateData, nil)
				if err != nil {
					// Log error nhưng tiếp tục với permission tiếp theo
					_ = fmt.Errorf("failed to update role permission scope: %v", err)
				}
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
	if err != nil && err != common.ErrNotFound {
		return nil, err
	}

	// Nếu chưa có role Administrator, tạo mới
	if err == common.ErrNotFound {
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
		return nil, common.ErrUserAlreadyAdmin
	}

	// Xử lý các lỗi khác ngoài ErrNotFound
	if err != common.ErrNotFound {
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

// InitAdminUser tạo user admin tự động từ Firebase UID (nếu có config)
// Sử dụng khi có FIREBASE_ADMIN_UID trong config
// User sẽ được tạo từ Firebase và tự động gán role Administrator
func (h *InitService) InitAdminUser(firebaseUID string) error {
	if firebaseUID == "" {
		return nil // Không có config, bỏ qua
	}

	// Kiểm tra user đã tồn tại chưa
	filter := bson.M{"firebaseUid": firebaseUID}
	existingUser, err := h.userService.FindOne(context.TODO(), filter, nil)
	if err != nil && err != common.ErrNotFound {
		return fmt.Errorf("failed to check existing admin user: %v", err)
	}

	var userID primitive.ObjectID

	// Nếu user chưa tồn tại, tạo từ Firebase
	if err == common.ErrNotFound {
		// Lấy thông tin user từ Firebase
		firebaseUser, err := utility.GetUserByUID(context.TODO(), firebaseUID)
		if err != nil {
			return fmt.Errorf("failed to get user from Firebase: %v", err)
		}

		// Tạo user mới
		currentTime := time.Now().Unix()
		newUser := &models.User{
			FirebaseUID:   firebaseUID,
			Email:         firebaseUser.Email,
			EmailVerified: firebaseUser.EmailVerified,
			Phone:         firebaseUser.PhoneNumber,
			PhoneVerified: firebaseUser.PhoneNumber != "",
			Name:          firebaseUser.DisplayName,
			AvatarURL:     firebaseUser.PhotoURL,
			IsBlock:       false,
			Tokens:        []models.Token{},
			CreatedAt:     currentTime,
			UpdatedAt:     currentTime,
		}

		createdUser, err := h.userService.InsertOne(context.TODO(), *newUser)
		if err != nil {
			return fmt.Errorf("failed to create admin user: %v", err)
		}

		userID = createdUser.ID
	} else {
		// User đã tồn tại
		userID = existingUser.ID
	}

	// Gán role Administrator cho user
	_, err = h.SetAdministrator(userID)
	if err != nil && err != common.ErrUserAlreadyAdmin {
		return fmt.Errorf("failed to set administrator role: %v", err)
	}

	return nil
}

// GetInitStatus kiểm tra trạng thái khởi tạo hệ thống
// Trả về thông tin về các đơn vị cơ bản đã được khởi tạo chưa
func (h *InitService) GetInitStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	// Kiểm tra Organization Root
	_, err := h.GetRootOrganization()
	status["organization"] = map[string]interface{}{
		"initialized": err == nil,
		"error": func() string {
			if err != nil {
				return err.Error()
			} else {
				return ""
			}
		}(),
	}

	// Kiểm tra Permissions
	permissions, err := h.permissionService.Find(context.TODO(), bson.M{}, nil)
	permissionCount := 0
	if err == nil {
		permissionCount = len(permissions)
	}
	status["permissions"] = map[string]interface{}{
		"initialized": err == nil && permissionCount > 0,
		"count":       permissionCount,
		"error": func() string {
			if err != nil {
				return err.Error()
			} else {
				return ""
			}
		}(),
	}

	// Kiểm tra Role Administrator và admin users
	adminRole, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	status["roles"] = map[string]interface{}{
		"initialized": err == nil,
		"error": func() string {
			if err != nil && err != common.ErrNotFound {
				return err.Error()
			} else {
				return ""
			}
		}(),
	}
	adminUserCount := 0
	if err == nil {
		var modelRole models.Role
		bsonBytes, _ := bson.Marshal(adminRole)
		if err := bson.Unmarshal(bsonBytes, &modelRole); err == nil {
			userRoles, err := h.userRoleService.Find(context.TODO(), bson.M{"roleId": modelRole.ID}, nil)
			if err == nil {
				adminUserCount = len(userRoles)
			}
		}
	}
	status["adminUsers"] = map[string]interface{}{
		"count":    adminUserCount,
		"hasAdmin": adminUserCount > 0,
	}

	return status, nil
}

// HasAnyAdministrator kiểm tra xem hệ thống đã có administrator chưa
// Returns:
//   - bool: true nếu đã có ít nhất một administrator
//   - error: Lỗi nếu có
func (h *InitService) HasAnyAdministrator() (bool, error) {
	// Kiểm tra role Administrator có tồn tại không
	adminRole, err := h.roleService.FindOne(context.TODO(), bson.M{"name": "Administrator"}, nil)
	if err != nil {
		if err == common.ErrNotFound {
			return false, nil // Chưa có role Administrator
		}
		return false, err
	}

	// Chuyển đổi sang model
	var modelRole models.Role
	bsonBytes, _ := bson.Marshal(adminRole)
	if err := bson.Unmarshal(bsonBytes, &modelRole); err != nil {
		return false, err
	}

	// Kiểm tra có user nào có role Administrator không
	userRoles, err := h.userRoleService.Find(context.TODO(), bson.M{"roleId": modelRole.ID}, nil)
	if err != nil {
		return false, err
	}

	return len(userRoles) > 0, nil
}
