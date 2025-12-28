// Package services chứa các service xử lý logic nghiệp vụ của ứng dụng
package services

import (
	"context"
	"fmt"
	"strings"
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
	userService                *UserService                // Service xử lý người dùng
	roleService             *RoleService                  // Service xử lý vai trò
	permissionService          *PermissionService          // Service xử lý quyền
	rolePermissionService      *RolePermissionService      // Service xử lý quan hệ vai trò-quyền
	userRoleService            *UserRoleService            // Service xử lý quan hệ người dùng-vai trò
	organizationService            *OrganizationService            // Service xử lý tổ chức
	notificationSenderService      *NotificationSenderService      // Service xử lý notification sender
	notificationTemplateService    *NotificationTemplateService    // Service xử lý notification template
	notificationChannelService     *NotificationChannelService     // Service xử lý notification channel
	notificationRoutingService     *NotificationRoutingService     // Service xử lý notification routing
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

	notificationSenderService, err := NewNotificationSenderService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification sender service: %v", err)
	}

	notificationTemplateService, err := NewNotificationTemplateService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification template service: %v", err)
	}

	notificationChannelService, err := NewNotificationChannelService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification channel service: %v", err)
	}

	notificationRoutingService, err := NewNotificationRoutingService()
	if err != nil {
		return nil, fmt.Errorf("failed to create notification routing service: %v", err)
	}

	return &InitService{
		userService:                 userService,
		roleService:                  roleService,
		permissionService:            permissionService,
		rolePermissionService:       rolePermissionService,
		userRoleService:              userRoleService,
	organizationService:            organizationService,
	notificationSenderService:      notificationSenderService,
	notificationTemplateService:    notificationTemplateService,
	notificationChannelService:      notificationChannelService,
	notificationRoutingService:     notificationRoutingService,
}, nil
}

// InitDefaultNotificationTeam khởi tạo team mặc định cho hệ thống notification
// Tạo team "Tech Team" thuộc System Organization và channel mặc định
// Returns:
//   - *models.Organization: Team mặc định đã tạo
//   - error: Lỗi nếu có trong quá trình khởi tạo
func (h *InitService) InitDefaultNotificationTeam() (*models.Organization, error) {
	// Sử dụng context cho phép insert system data trong quá trình init
	// Lưu ý: withSystemDataInsertAllowed là unexported, chỉ có thể gọi từ trong package services
	ctx := withSystemDataInsertAllowed(context.TODO())
	currentTime := time.Now().Unix()

	// Lấy System Organization
	systemOrg, err := h.GetRootOrganization()
	if err != nil {
		return nil, fmt.Errorf("failed to get system organization: %v", err)
	}

	// Kiểm tra team mặc định đã tồn tại chưa
	teamFilter := bson.M{
		"code":     "TECH_TEAM",
		"parentId": systemOrg.ID,
	}
	existingTeam, err := h.organizationService.FindOne(ctx, teamFilter, nil)
	if err != nil && err != common.ErrNotFound {
		return nil, fmt.Errorf("failed to check existing tech team: %v", err)
	}

	var techTeam *models.Organization
	if err == common.ErrNotFound {
		// Tạo mới Tech Team
		techTeamModel := models.Organization{
			Name:      "Tech Team",
			Code:      "TECH_TEAM",
			Type:      models.OrganizationTypeTeam,
			ParentID:  &systemOrg.ID,
			Path:      systemOrg.Path + "/TECH_TEAM",
			Level:     systemOrg.Level + 1, // Level = 0 (vì System là -1)
			IsActive:  true,
			IsSystem:  true, // Đánh dấu là dữ liệu hệ thống
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}

		createdTeam, err := h.organizationService.InsertOne(ctx, techTeamModel)
		if err != nil {
			return nil, fmt.Errorf("failed to create tech team: %v", err)
		}

		var modelTeam models.Organization
		bsonBytes, _ := bson.Marshal(createdTeam)
		if err := bson.Unmarshal(bsonBytes, &modelTeam); err != nil {
			return nil, fmt.Errorf("failed to decode tech team: %v", err)
		}
		techTeam = &modelTeam
	} else {
		// Team đã tồn tại
		var modelTeam models.Organization
		bsonBytes, _ := bson.Marshal(existingTeam)
		if err := bson.Unmarshal(bsonBytes, &modelTeam); err != nil {
			return nil, fmt.Errorf("failed to decode existing tech team: %v", err)
		}
		techTeam = &modelTeam
	}

	return techTeam, nil
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

	// Quản lý tin nhắn Facebook Item: Thêm, xem, sửa, xóa
	{Name: "FbMessageItem.Insert", Describe: "Quyền tạo tin nhắn item", Group: "Pancake", Category: "FbMessageItem"},
	{Name: "FbMessageItem.Read", Describe: "Quyền xem danh sách tin nhắn item", Group: "Pancake", Category: "FbMessageItem"},
	{Name: "FbMessageItem.Update", Describe: "Quyền cập nhật tin nhắn item", Group: "Pancake", Category: "FbMessageItem"},
	{Name: "FbMessageItem.Delete", Describe: "Quyền xóa tin nhắn item", Group: "Pancake", Category: "FbMessageItem"},

	// Quản lý khách hàng: Thêm, xem, sửa, xóa
	{Name: "Customer.Insert", Describe: "Quyền tạo khách hàng", Group: "Pancake", Category: "Customer"},
	{Name: "Customer.Read", Describe: "Quyền xem danh sách khách hàng", Group: "Pancake", Category: "Customer"},
	{Name: "Customer.Update", Describe: "Quyền cập nhật thông tin khách hàng", Group: "Pancake", Category: "Customer"},
	{Name: "Customer.Delete", Describe: "Quyền xóa khách hàng", Group: "Pancake", Category: "Customer"},

	// Quản lý khách hàng Facebook: Thêm, xem, sửa, xóa
	{Name: "FbCustomer.Insert", Describe: "Quyền tạo khách hàng Facebook", Group: "Pancake", Category: "FbCustomer"},
	{Name: "FbCustomer.Read", Describe: "Quyền xem danh sách khách hàng Facebook", Group: "Pancake", Category: "FbCustomer"},
	{Name: "FbCustomer.Update", Describe: "Quyền cập nhật thông tin khách hàng Facebook", Group: "Pancake", Category: "FbCustomer"},
	{Name: "FbCustomer.Delete", Describe: "Quyền xóa khách hàng Facebook", Group: "Pancake", Category: "FbCustomer"},

	// Quản lý khách hàng POS: Thêm, xem, sửa, xóa
	{Name: "PcPosCustomer.Insert", Describe: "Quyền tạo khách hàng POS", Group: "Pancake", Category: "PcPosCustomer"},
	{Name: "PcPosCustomer.Read", Describe: "Quyền xem danh sách khách hàng POS", Group: "Pancake", Category: "PcPosCustomer"},
	{Name: "PcPosCustomer.Update", Describe: "Quyền cập nhật thông tin khách hàng POS", Group: "Pancake", Category: "PcPosCustomer"},
	{Name: "PcPosCustomer.Delete", Describe: "Quyền xóa khách hàng POS", Group: "Pancake", Category: "PcPosCustomer"},

	// Quản lý cửa hàng Pancake POS: Thêm, xem, sửa, xóa
	{Name: "PcPosShop.Insert", Describe: "Quyền tạo cửa hàng từ Pancake POS", Group: "Pancake", Category: "PcPosShop"},
	{Name: "PcPosShop.Read", Describe: "Quyền xem danh sách cửa hàng từ Pancake POS", Group: "Pancake", Category: "PcPosShop"},
	{Name: "PcPosShop.Update", Describe: "Quyền cập nhật thông tin cửa hàng từ Pancake POS", Group: "Pancake", Category: "PcPosShop"},
	{Name: "PcPosShop.Delete", Describe: "Quyền xóa cửa hàng từ Pancake POS", Group: "Pancake", Category: "PcPosShop"},

	// Quản lý kho hàng Pancake POS: Thêm, xem, sửa, xóa
	{Name: "PcPosWarehouse.Insert", Describe: "Quyền tạo kho hàng từ Pancake POS", Group: "Pancake", Category: "PcPosWarehouse"},
	{Name: "PcPosWarehouse.Read", Describe: "Quyền xem danh sách kho hàng từ Pancake POS", Group: "Pancake", Category: "PcPosWarehouse"},
	{Name: "PcPosWarehouse.Update", Describe: "Quyền cập nhật thông tin kho hàng từ Pancake POS", Group: "Pancake", Category: "PcPosWarehouse"},
	{Name: "PcPosWarehouse.Delete", Describe: "Quyền xóa kho hàng từ Pancake POS", Group: "Pancake", Category: "PcPosWarehouse"},

	// Quản lý sản phẩm Pancake POS: Thêm, xem, sửa, xóa
	{Name: "PcPosProduct.Insert", Describe: "Quyền tạo sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosProduct"},
	{Name: "PcPosProduct.Read", Describe: "Quyền xem danh sách sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosProduct"},
	{Name: "PcPosProduct.Update", Describe: "Quyền cập nhật thông tin sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosProduct"},
	{Name: "PcPosProduct.Delete", Describe: "Quyền xóa sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosProduct"},

	// Quản lý biến thể sản phẩm Pancake POS: Thêm, xem, sửa, xóa
	{Name: "PcPosVariation.Insert", Describe: "Quyền tạo biến thể sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosVariation"},
	{Name: "PcPosVariation.Read", Describe: "Quyền xem danh sách biến thể sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosVariation"},
	{Name: "PcPosVariation.Update", Describe: "Quyền cập nhật thông tin biến thể sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosVariation"},
	{Name: "PcPosVariation.Delete", Describe: "Quyền xóa biến thể sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosVariation"},

	// Quản lý danh mục sản phẩm Pancake POS: Thêm, xem, sửa, xóa
	{Name: "PcPosCategory.Insert", Describe: "Quyền tạo danh mục sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosCategory"},
	{Name: "PcPosCategory.Read", Describe: "Quyền xem danh sách danh mục sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosCategory"},
	{Name: "PcPosCategory.Update", Describe: "Quyền cập nhật thông tin danh mục sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosCategory"},
	{Name: "PcPosCategory.Delete", Describe: "Quyền xóa danh mục sản phẩm từ Pancake POS", Group: "Pancake", Category: "PcPosCategory"},

	// Quản lý đơn hàng Pancake POS: Thêm, xem, sửa, xóa
	{Name: "PcPosOrder.Insert", Describe: "Quyền tạo đơn hàng từ Pancake POS", Group: "Pancake", Category: "PcPosOrder"},
	{Name: "PcPosOrder.Read", Describe: "Quyền xem danh sách đơn hàng từ Pancake POS", Group: "Pancake", Category: "PcPosOrder"},
	{Name: "PcPosOrder.Update", Describe: "Quyền cập nhật thông tin đơn hàng từ Pancake POS", Group: "Pancake", Category: "PcPosOrder"},
	{Name: "PcPosOrder.Delete", Describe: "Quyền xóa đơn hàng từ Pancake POS", Group: "Pancake", Category: "PcPosOrder"},

	// ==================================== NOTIFICATION MODULE ===========================================
	// Quản lý Notification Sender: Thêm, xem, sửa, xóa
	{Name: "NotificationSender.Insert", Describe: "Quyền tạo cấu hình sender thông báo", Group: "Notification", Category: "NotificationSender"},
	{Name: "NotificationSender.Read", Describe: "Quyền xem danh sách cấu hình sender thông báo", Group: "Notification", Category: "NotificationSender"},
	{Name: "NotificationSender.Update", Describe: "Quyền cập nhật cấu hình sender thông báo", Group: "Notification", Category: "NotificationSender"},
	{Name: "NotificationSender.Delete", Describe: "Quyền xóa cấu hình sender thông báo", Group: "Notification", Category: "NotificationSender"},

	// Quản lý Notification Channel: Thêm, xem, sửa, xóa
	{Name: "NotificationChannel.Insert", Describe: "Quyền tạo kênh thông báo cho team", Group: "Notification", Category: "NotificationChannel"},
	{Name: "NotificationChannel.Read", Describe: "Quyền xem danh sách kênh thông báo", Group: "Notification", Category: "NotificationChannel"},
	{Name: "NotificationChannel.Update", Describe: "Quyền cập nhật kênh thông báo", Group: "Notification", Category: "NotificationChannel"},
	{Name: "NotificationChannel.Delete", Describe: "Quyền xóa kênh thông báo", Group: "Notification", Category: "NotificationChannel"},

	// Quản lý Notification Template: Thêm, xem, sửa, xóa
	{Name: "NotificationTemplate.Insert", Describe: "Quyền tạo template thông báo", Group: "Notification", Category: "NotificationTemplate"},
	{Name: "NotificationTemplate.Read", Describe: "Quyền xem danh sách template thông báo", Group: "Notification", Category: "NotificationTemplate"},
	{Name: "NotificationTemplate.Update", Describe: "Quyền cập nhật template thông báo", Group: "Notification", Category: "NotificationTemplate"},
	{Name: "NotificationTemplate.Delete", Describe: "Quyền xóa template thông báo", Group: "Notification", Category: "NotificationTemplate"},

	// Quản lý Notification Routing Rule: Thêm, xem, sửa, xóa
	{Name: "NotificationRouting.Insert", Describe: "Quyền tạo routing rule thông báo", Group: "Notification", Category: "NotificationRouting"},
	{Name: "NotificationRouting.Read", Describe: "Quyền xem danh sách routing rule thông báo", Group: "Notification", Category: "NotificationRouting"},
	{Name: "NotificationRouting.Update", Describe: "Quyền cập nhật routing rule thông báo", Group: "Notification", Category: "NotificationRouting"},
	{Name: "NotificationRouting.Delete", Describe: "Quyền xóa routing rule thông báo", Group: "Notification", Category: "NotificationRouting"},

	// Quản lý Notification History: Chỉ xem
	{Name: "NotificationHistory.Read", Describe: "Quyền xem lịch sử thông báo", Group: "Notification", Category: "NotificationHistory"},

	// Trigger Notification: Gửi thông báo
	{Name: "Notification.Trigger", Describe: "Quyền trigger/gửi thông báo", Group: "Notification", Category: "Notification"},
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
			// Set IsSystem = true cho tất cả permissions được tạo trong init
			permission.IsSystem = true
			// Sử dụng context cho phép insert system data trong quá trình init
			// Lưu ý: withSystemDataInsertAllowed là unexported, chỉ có thể gọi từ trong package services
			initCtx := withSystemDataInsertAllowed(context.TODO())
			_, err = h.permissionService.InsertOne(initCtx, permission)
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
		IsSystem: true, // Đánh dấu là dữ liệu hệ thống
	}

	// Sử dụng context cho phép insert system data trong quá trình init
	// Lưu ý: withSystemDataInsertAllowed là unexported, chỉ có thể gọi từ trong package services
	initCtx := withSystemDataInsertAllowed(context.TODO())
	_, err = h.organizationService.InsertOne(initCtx, systemOrgModel)
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
			IsSystem:       true,        // Đánh dấu là dữ liệu hệ thống
		}

		// Lưu vai trò vào database
		// Sử dụng context cho phép insert system data trong quá trình init
		// Lưu ý: withSystemDataInsertAllowed là unexported, chỉ có thể gọi từ trong package services
		initCtx := withSystemDataInsertAllowed(context.TODO())
		adminRole, err = h.roleService.InsertOne(initCtx, newAdminRole)
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

	// Đảm bảo role Administrator có đầy đủ tất cả các quyền trong hệ thống
	// Gọi CheckPermissionForAdministrator để cập nhật quyền cho role Administrator
	err = h.CheckPermissionForAdministrator()
	if err != nil {
		// Log lỗi nhưng không fail việc set administrator
		// Vì role đã được gán, chỉ là quyền có thể chưa được cập nhật đầy đủ
		_ = fmt.Errorf("failed to check permissions for administrator: %v", err)
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

// InitNotificationData khởi tạo dữ liệu mặc định cho hệ thống notification
// Tạo các sender và template mặc định (global), các thông tin như token/password sẽ để trống để admin bổ sung sau
// Returns:
//   - error: Lỗi nếu có trong quá trình khởi tạo
func (h *InitService) InitNotificationData() error {
	// Sử dụng context cho phép insert system data trong quá trình init
	// Lưu ý: withSystemDataInsertAllowed là unexported, chỉ có thể gọi từ trong package services
	ctx := withSystemDataInsertAllowed(context.TODO())
	currentTime := time.Now().Unix()
	var err error

	// ==================================== 0. KHỞI TẠO TEAM MẶC ĐỊNH CHO NOTIFICATION =============================================
	// Tạo Tech Team thuộc System Organization để có thể tạo channel và routing rule mặc định
	techTeam, err := h.InitDefaultNotificationTeam()
	if err != nil {
		return fmt.Errorf("failed to initialize default notification team: %v", err)
	}

	// ==================================== 1. KHỞI TẠO NOTIFICATION SENDERS (GLOBAL) =============================================
	// Sender cho Email
	emailSenderFilter := bson.M{
		"organizationId": nil,
		"channelType":    "email",
		"name":           "Email Sender Mặc Định",
	}
	_, err = h.notificationSenderService.FindOne(ctx, emailSenderFilter, nil)
	if err != nil && err != common.ErrNotFound {
		return fmt.Errorf("failed to check existing email sender: %v", err)
	}
	if err == common.ErrNotFound {
		emailSender := models.NotificationChannelSender{
			OrganizationID: nil, // Global sender
			ChannelType:    "email",
			Name:           "Email Sender Mặc Định",
			IsActive:       false, // Tắt mặc định, admin cần cấu hình token/password trước khi bật
			IsSystem:       true,  // Đánh dấu là dữ liệu hệ thống, không thể xóa
			SMTPHost:       "",   // Admin cần bổ sung
			SMTPPort:       587,  // Port mặc định
			SMTPUsername:   "",   // Admin cần bổ sung
			SMTPPassword:   "",   // Admin cần bổ sung
			FromEmail:      "",   // Admin cần bổ sung
			FromName:       "",   // Admin cần bổ sung
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
		}
		_, err = h.notificationSenderService.InsertOne(ctx, emailSender)
		if err != nil {
			return fmt.Errorf("failed to create email sender: %v", err)
		}
	}

	// Sender cho Telegram
	telegramSenderFilter := bson.M{
		"organizationId": nil,
		"channelType":    "telegram",
		"name":           "Telegram Bot Mặc Định",
	}
	_, err = h.notificationSenderService.FindOne(ctx, telegramSenderFilter, nil)
	if err != nil && err != common.ErrNotFound {
		return fmt.Errorf("failed to check existing telegram sender: %v", err)
	}
	if err == common.ErrNotFound {
		telegramSender := models.NotificationChannelSender{
			OrganizationID: nil, // Global sender
			ChannelType:    "telegram",
			Name:           "Telegram Bot Mặc Định",
			IsActive:       false, // Tắt mặc định, admin cần cấu hình bot token trước khi bật
			IsSystem:       true,  // Đánh dấu là dữ liệu hệ thống, không thể xóa
			BotToken:       "",   // Admin cần bổ sung
			BotUsername:    "",   // Admin cần bổ sung
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
		}
		_, err = h.notificationSenderService.InsertOne(ctx, telegramSender)
		if err != nil {
			return fmt.Errorf("failed to create telegram sender: %v", err)
		}
	}

	// Sender cho Webhook
	webhookSenderFilter := bson.M{
		"organizationId": nil,
		"channelType":    "webhook",
		"name":           "Webhook Sender Mặc Định",
	}
	_, err = h.notificationSenderService.FindOne(ctx, webhookSenderFilter, nil)
	if err != nil && err != common.ErrNotFound {
		return fmt.Errorf("failed to check existing webhook sender: %v", err)
	}
	if err == common.ErrNotFound {
		webhookSender := models.NotificationChannelSender{
			OrganizationID: nil, // Global sender
			ChannelType:    "webhook",
			Name:           "Webhook Sender Mặc Định",
			IsActive:       false, // Tắt mặc định, admin cần cấu hình trước khi bật
			IsSystem:       true,  // Đánh dấu là dữ liệu hệ thống, không thể xóa
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
		}
		_, err = h.notificationSenderService.InsertOne(ctx, webhookSender)
		if err != nil {
			return fmt.Errorf("failed to create webhook sender: %v", err)
		}
	}

	// ==================================== 2. KHỞI TẠO NOTIFICATION TEMPLATES (GLOBAL) =============================================
	// Template cho event conversation_unreplied - Email
	convUnrepliedEmailFilter := bson.M{
		"organizationId": nil,
		"eventType":      "conversation_unreplied",
		"channelType":    "email",
	}
	_, err = h.notificationTemplateService.FindOne(ctx, convUnrepliedEmailFilter, nil)
	if err == common.ErrNotFound {
		template := models.NotificationTemplate{
			OrganizationID: nil, // Global template
			EventType:      "conversation_unreplied",
			ChannelType:    "email",
			Subject:        "Cảnh báo: Cuộc trò chuyện chưa được trả lời",
			Content: `Xin chào,

Bạn có một cuộc trò chuyện chưa được trả lời trong {{minutes}} phút.

Thông tin cuộc trò chuyện:
- ID: {{conversationId}}
- Khách hàng: {{customerName}}
- Thời gian: {{lastMessageAt}}

Vui lòng kiểm tra và phản hồi sớm nhất có thể.

Trân trọng,
Hệ thống thông báo`,
			Variables: []string{"conversationId", "minutes", "customerName", "lastMessageAt"},
			CTAs: []models.NotificationCTA{
				{
					Label:  "Xem cuộc trò chuyện",
					Action: "{{baseUrl}}/conversations/{{conversationId}}",
					Style:  "primary",
				},
			},
			IsActive:  true,
			IsSystem:  true, // Đánh dấu là dữ liệu hệ thống, không thể xóa
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		_, err = h.notificationTemplateService.InsertOne(ctx, template)
		if err != nil {
			return fmt.Errorf("failed to create conversation_unreplied email template: %v", err)
		}
	}

	// Template cho event conversation_unreplied - Telegram
	convUnrepliedTelegramFilter := bson.M{
		"organizationId": nil,
		"eventType":      "conversation_unreplied",
		"channelType":    "telegram",
	}
	_, err = h.notificationTemplateService.FindOne(ctx, convUnrepliedTelegramFilter, nil)
	if err == common.ErrNotFound {
		template := models.NotificationTemplate{
			OrganizationID: nil, // Global template
			EventType:      "conversation_unreplied",
			ChannelType:    "telegram",
			Subject:        "", // Telegram không có subject
			Content: `🚨 *Cảnh báo: Cuộc trò chuyện chưa được trả lời*

Bạn có một cuộc trò chuyện chưa được trả lời trong *{{minutes}}* phút.

*Thông tin:*
• ID: ` + "`{{conversationId}}`" + `
• Khách hàng: {{customerName}}
• Thời gian: {{lastMessageAt}}

Vui lòng kiểm tra và phản hồi sớm nhất có thể.`,
			Variables: []string{"conversationId", "minutes", "customerName", "lastMessageAt"},
			CTAs: []models.NotificationCTA{
				{
					Label:  "Xem cuộc trò chuyện",
					Action: "{{baseUrl}}/conversations/{{conversationId}}",
					Style:  "primary",
				},
			},
			IsActive:  true,
			IsSystem:  true, // Đánh dấu là dữ liệu hệ thống, không thể xóa
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		_, err = h.notificationTemplateService.InsertOne(ctx, template)
		if err != nil {
			return fmt.Errorf("failed to create conversation_unreplied telegram template: %v", err)
		}
	}

	// Template cho event conversation_unreplied - Webhook
	convUnrepliedWebhookFilter := bson.M{
		"organizationId": nil,
		"eventType":      "conversation_unreplied",
		"channelType":    "webhook",
	}
	_, err = h.notificationTemplateService.FindOne(ctx, convUnrepliedWebhookFilter, nil)
	if err == common.ErrNotFound {
		template := models.NotificationTemplate{
			OrganizationID: nil, // Global template
			EventType:      "conversation_unreplied",
			ChannelType:    "webhook",
			Subject:        "", // Webhook không có subject
			Content: `{"eventType":"conversation_unreplied","conversationId":"{{conversationId}}","minutes":{{minutes}},"customerName":"{{customerName}}","lastMessageAt":"{{lastMessageAt}}","baseUrl":"{{baseUrl}}"}`,
			Variables: []string{"conversationId", "minutes", "customerName", "lastMessageAt", "baseUrl"},
			IsActive:  true,
			IsSystem:  true, // Đánh dấu là dữ liệu hệ thống, không thể xóa
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		_, err = h.notificationTemplateService.InsertOne(ctx, template)
		if err != nil {
			return fmt.Errorf("failed to create conversation_unreplied webhook template: %v", err)
		}
	}

	// ==================================== 3. KHỞI TẠO NOTIFICATION CHANNELS CHO TECH TEAM =============================================
	// Channel Email mặc định cho Tech Team
	emailChannelFilter := bson.M{
		"organizationId": techTeam.ID,
		"channelType":    "email",
		"name":           "Email Channel Mặc Định",
	}
	_, err = h.notificationChannelService.FindOne(ctx, emailChannelFilter, nil)
	if err == common.ErrNotFound {
		emailChannel := models.NotificationChannel{
			OrganizationID: techTeam.ID,
			ChannelType:    "email",
			Name:           "Email Channel Mặc Định",
			IsActive:       false, // Tắt mặc định, admin cần cấu hình recipients trước khi bật
			IsSystem:       true,  // Đánh dấu là dữ liệu hệ thống, không thể xóa
			Recipients:     []string{}, // Admin cần bổ sung email addresses
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
		}
		_, err = h.notificationChannelService.InsertOne(ctx, emailChannel)
		if err != nil {
			return fmt.Errorf("failed to create email channel for tech team: %v", err)
		}
	}

	// Channel Telegram mặc định cho Tech Team
	telegramChannelFilter := bson.M{
		"organizationId": techTeam.ID,
		"channelType":    "telegram",
		"name":           "Telegram Channel Mặc Định",
	}
	_, err = h.notificationChannelService.FindOne(ctx, telegramChannelFilter, nil)
	if err == common.ErrNotFound {
		telegramChannel := models.NotificationChannel{
			OrganizationID: techTeam.ID,
			ChannelType:    "telegram",
			Name:           "Telegram Channel Mặc Định",
			IsActive:       false, // Tắt mặc định, admin cần cấu hình chat IDs trước khi bật
			IsSystem:       true,  // Đánh dấu là dữ liệu hệ thống, không thể xóa
			ChatIDs:        []string{}, // Admin cần bổ sung Telegram chat IDs
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
		}
		_, err = h.notificationChannelService.InsertOne(ctx, telegramChannel)
		if err != nil {
			return fmt.Errorf("failed to create telegram channel for tech team: %v", err)
		}
	}

	// Channel Webhook mặc định cho Tech Team
	webhookChannelFilter := bson.M{
		"organizationId": techTeam.ID,
		"channelType":    "webhook",
		"name":           "Webhook Channel Mặc Định",
	}
	_, err = h.notificationChannelService.FindOne(ctx, webhookChannelFilter, nil)
	if err == common.ErrNotFound {
		webhookChannel := models.NotificationChannel{
			OrganizationID: techTeam.ID,
			ChannelType:    "webhook",
			Name:           "Webhook Channel Mặc Định",
			IsActive:       false, // Tắt mặc định, admin cần cấu hình webhook URL trước khi bật
			IsSystem:       true,  // Đánh dấu là dữ liệu hệ thống, không thể xóa
			WebhookURL:     "",   // Admin cần bổ sung webhook URL
			WebhookHeaders: map[string]string{}, // Admin có thể bổ sung headers nếu cần
			CreatedAt:      currentTime,
			UpdatedAt:      currentTime,
		}
		_, err = h.notificationChannelService.InsertOne(ctx, webhookChannel)
		if err != nil {
			return fmt.Errorf("failed to create webhook channel for tech team: %v", err)
		}
	}

	// ==================================== 4. KHỞI TẠO TEMPLATES CHO CÁC EVENT CẤP HỆ THỐNG =============================================
	systemEvents := []struct {
		eventType string
		subject   string
		content   string
		variables []string
	}{
		{
			eventType: "system_startup",
			subject:   "Hệ thống đã khởi động",
			content: `Xin chào,

Hệ thống đã được khởi động thành công.

Thông tin:
- Thời gian: {{timestamp}}
- Phiên bản: {{version}}
- Môi trường: {{environment}}

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "version", "environment"},
		},
		{
			eventType: "system_shutdown",
			subject:   "Cảnh báo: Hệ thống đang tắt",
			content: `Xin chào,

Hệ thống đang được tắt.

Thông tin:
- Thời gian: {{timestamp}}
- Lý do: {{reason}}

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "reason"},
		},
		{
			eventType: "system_error",
			subject:   "🚨 Lỗi hệ thống nghiêm trọng",
			content: `Xin chào,

Hệ thống đã gặp lỗi nghiêm trọng.

Thông tin lỗi:
- Thời gian: {{timestamp}}
- Loại lỗi: {{errorType}}
- Mô tả: {{errorMessage}}
- Chi tiết: {{errorDetails}}

Vui lòng kiểm tra và xử lý ngay lập tức.

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "errorType", "errorMessage", "errorDetails"},
		},
		{
			eventType: "system_warning",
			subject:   "⚠️ Cảnh báo hệ thống",
			content: `Xin chào,

Hệ thống có cảnh báo cần chú ý.

Thông tin:
- Thời gian: {{timestamp}}
- Loại cảnh báo: {{warningType}}
- Mô tả: {{warningMessage}}

Vui lòng kiểm tra và xử lý.

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "warningType", "warningMessage"},
		},
		{
			eventType: "database_error",
			subject:   "🚨 Lỗi kết nối Database",
			content: `Xin chào,

Hệ thống gặp lỗi khi kết nối với Database.

Thông tin lỗi:
- Thời gian: {{timestamp}}
- Database: {{databaseName}}
- Lỗi: {{errorMessage}}

Vui lòng kiểm tra kết nối database ngay lập tức.

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "databaseName", "errorMessage"},
		},
		{
			eventType: "api_error",
			subject:   "⚠️ Lỗi API",
			content: `Xin chào,

Hệ thống gặp lỗi khi xử lý API request.

Thông tin:
- Thời gian: {{timestamp}}
- Endpoint: {{endpoint}}
- Method: {{method}}
- Lỗi: {{errorMessage}}
- Status Code: {{statusCode}}

Vui lòng kiểm tra và xử lý.

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "endpoint", "method", "errorMessage", "statusCode"},
		},
		{
			eventType: "backup_completed",
			subject:   "✅ Backup hoàn tất",
			content: `Xin chào,

Quá trình backup đã hoàn tất thành công.

Thông tin:
- Thời gian: {{timestamp}}
- Loại backup: {{backupType}}
- Kích thước: {{backupSize}}
- Vị trí: {{backupLocation}}

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "backupType", "backupSize", "backupLocation"},
		},
		{
			eventType: "backup_failed",
			subject:   "❌ Backup thất bại",
			content: `Xin chào,

Quá trình backup đã thất bại.

Thông tin:
- Thời gian: {{timestamp}}
- Loại backup: {{backupType}}
- Lỗi: {{errorMessage}}

Vui lòng kiểm tra và thử lại.

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "backupType", "errorMessage"},
		},
		{
			eventType: "rate_limit_exceeded",
			subject:   "⚠️ Vượt quá Rate Limit",
			content: `Xin chào,

Hệ thống đã vượt quá rate limit.

Thông tin:
- Thời gian: {{timestamp}}
- Endpoint: {{endpoint}}
- IP: {{ipAddress}}
- Số request: {{requestCount}}
- Giới hạn: {{rateLimit}}

Vui lòng kiểm tra và điều chỉnh.

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "endpoint", "ipAddress", "requestCount", "rateLimit"},
		},
		{
			eventType: "security_alert",
			subject:   "🚨 Cảnh báo bảo mật",
			content: `Xin chào,

Hệ thống phát hiện hoạt động đáng ngờ hoặc vi phạm bảo mật.

Thông tin:
- Thời gian: {{timestamp}}
- Loại cảnh báo: {{alertType}}
- Mô tả: {{alertMessage}}
- IP: {{ipAddress}}
- User: {{username}}

Vui lòng kiểm tra và xử lý ngay lập tức.

Trân trọng,
Hệ thống thông báo`,
			variables: []string{"timestamp", "alertType", "alertMessage", "ipAddress", "username"},
		},
	}

	// Tạo templates cho mỗi system event (Email, Telegram, Webhook)
	for _, event := range systemEvents {
		// Email template
		emailFilter := bson.M{
			"organizationId": nil,
			"eventType":      event.eventType,
			"channelType":    "email",
		}
		_, err = h.notificationTemplateService.FindOne(ctx, emailFilter, nil)
		if err == common.ErrNotFound {
			template := models.NotificationTemplate{
				OrganizationID: nil,
				EventType:      event.eventType,
				ChannelType:    "email",
				Subject:        event.subject,
				Content:        event.content,
				Variables:      event.variables,
				IsActive:       true,
				IsSystem:       true, // Đánh dấu là dữ liệu hệ thống, không thể xóa
				CreatedAt:      currentTime,
				UpdatedAt:      currentTime,
			}
			_, err = h.notificationTemplateService.InsertOne(ctx, template)
			if err != nil {
				return fmt.Errorf("failed to create %s email template: %v", event.eventType, err)
			}
		}

		// Telegram template
		telegramFilter := bson.M{
			"organizationId": nil,
			"eventType":      event.eventType,
			"channelType":    "telegram",
		}
		_, err = h.notificationTemplateService.FindOne(ctx, telegramFilter, nil)
		if err == common.ErrNotFound {
			// Convert content to Telegram format (Markdown)
			telegramContent := event.content
			telegramContent = fmt.Sprintf("*%s*\n\n%s", event.subject, telegramContent)
			// Replace bullet points with Telegram format
			telegramContent = strings.ReplaceAll(telegramContent, "- ", "• ")

			template := models.NotificationTemplate{
				OrganizationID: nil,
				EventType:      event.eventType,
				ChannelType:    "telegram",
				Subject:        "",
				Content:        telegramContent,
				Variables:      event.variables,
				IsActive:       true,
				IsSystem:       true, // Đánh dấu là dữ liệu hệ thống, không thể xóa
				CreatedAt:      currentTime,
				UpdatedAt:      currentTime,
			}
			_, err = h.notificationTemplateService.InsertOne(ctx, template)
			if err != nil {
				return fmt.Errorf("failed to create %s telegram template: %v", event.eventType, err)
			}
		}

		// Webhook template (JSON format)
		webhookFilter := bson.M{
			"organizationId": nil,
			"eventType":      event.eventType,
			"channelType":    "webhook",
		}
		_, err = h.notificationTemplateService.FindOne(ctx, webhookFilter, nil)
		if err == common.ErrNotFound {
			// Create JSON template with all variables
			jsonVars := make([]string, 0)
			for _, v := range event.variables {
				jsonVars = append(jsonVars, fmt.Sprintf(`"%s":"{{%s}}"`, v, v))
			}
			jsonContent := fmt.Sprintf(`{"eventType":"%s",%s}`, event.eventType, strings.Join(jsonVars, ","))

			template := models.NotificationTemplate{
				OrganizationID: nil,
				EventType:      event.eventType,
				ChannelType:    "webhook",
				Subject:        "",
				Content:        jsonContent,
				Variables:      event.variables,
				IsActive:       true,
				IsSystem:       true, // Đánh dấu là dữ liệu hệ thống, không thể xóa
				CreatedAt:      currentTime,
				UpdatedAt:      currentTime,
			}
			_, err = h.notificationTemplateService.InsertOne(ctx, template)
			if err != nil {
				return fmt.Errorf("failed to create %s webhook template: %v", event.eventType, err)
			}
		}
	}

	// ==================================== 5. KHỞI TẠO ROUTING RULES MẶC ĐỊNH CHO SYSTEM EVENTS =============================================
	// Tạo routing rules để kết nối system events với Tech Team
	for _, event := range systemEvents {
		routingFilter := bson.M{
			"eventType": event.eventType,
		}
		_, err = h.notificationRoutingService.FindOne(ctx, routingFilter, nil)
		if err == common.ErrNotFound {
			routingRule := models.NotificationRoutingRule{
				EventType:       event.eventType,
				OrganizationIDs: []primitive.ObjectID{techTeam.ID},
				ChannelTypes:    []string{"email", "telegram", "webhook"}, // Tất cả channel types
				IsActive:       false, // Tắt mặc định, admin cần bật sau khi cấu hình channels
				IsSystem:       true,  // Đánh dấu là dữ liệu hệ thống, không thể xóa
				CreatedAt:      currentTime,
				UpdatedAt:      currentTime,
			}
			_, err = h.notificationRoutingService.InsertOne(ctx, routingRule)
			if err != nil {
				return fmt.Errorf("failed to create routing rule for %s: %v", event.eventType, err)
			}
		}
	}

	return nil
}
