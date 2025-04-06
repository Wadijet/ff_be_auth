package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FiberRolePermissionHandler xử lý các route liên quan đến phân quyền cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
// Các phương thức của FiberBaseHandler đã có sẵn:
// - InsertOne: Thêm mới một role permission
// - InsertMany: Thêm nhiều role permission
// - FindOne: Tìm một role permission theo điều kiện
// - FindOneById: Tìm một role permission theo ID
// - FindManyByIds: Tìm nhiều role permission theo danh sách ID
// - FindWithPagination: Tìm role permission với phân trang
// - Find: Tìm nhiều role permission theo điều kiện
// - UpdateOne: Cập nhật một role permission theo điều kiện
// - UpdateMany: Cập nhật nhiều role permission theo điều kiện
// - UpdateById: Cập nhật một role permission theo ID
// - DeleteOne: Xóa một role permission theo điều kiện
// - DeleteMany: Xóa nhiều role permission theo điều kiện
// - DeleteById: Xóa một role permission theo ID
// - FindOneAndUpdate: Tìm và cập nhật một role permission
// - FindOneAndDelete: Tìm và xóa một role permission
// - CountDocuments: Đếm số lượng role permission theo điều kiện
// - Distinct: Lấy danh sách giá trị duy nhất của một trường
// - Upsert: Thêm mới hoặc cập nhật một role permission
// - UpsertMany: Thêm mới hoặc cập nhật nhiều role permission
// - DocumentExists: Kiểm tra role permission có tồn tại không
type FiberRolePermissionHandler struct {
	FiberBaseHandler[models.RolePermission, models.RolePermissionCreateInput, models.RolePermission]
	RolePermissionService *services.RolePermissionService
}

// NewFiberRolePermissionHandler tạo một instance mới của FiberRolePermissionHandler
// Returns:
//   - *FiberRolePermissionHandler: Instance mới của FiberRolePermissionHandler đã được khởi tạo với RolePermissionService
func NewFiberRolePermissionHandler() *FiberRolePermissionHandler {
	handler := &FiberRolePermissionHandler{
		RolePermissionService: services.NewRolePermissionService(),
	}
	handler.Service = handler.RolePermissionService
	return handler
}

// HandleUpdateRolePermissions xử lý cập nhật quyền cho vai trò
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - roleId: ID của vai trò cần cập nhật quyền
//   - permissionIds: Danh sách ID của các quyền
//
// Response:
//   - 200: Cập nhật quyền thành công
//     {
//     "message": "Thành công",
//     "data": [
//     {
//     "id": "...",
//     "roleId": "...",
//     "permissionId": "...",
//     "scope": 0,
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     ]
//     }
//   - 400: Dữ liệu không hợp lệ
//   - 500: Lỗi server
func (h *FiberRolePermissionHandler) HandleUpdateRolePermissions(c fiber.Ctx) error {
	// Parse input từ request body
	input := new(models.RolePermissionUpdateInput)
	if err := c.Bind().Body(input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": utility.MsgValidationError,
		})
	}

	// Chuyển đổi roleId từ string sang ObjectID
	roleId, err := primitive.ObjectIDFromHex(input.RoleId)
	if err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": "ID vai trò không hợp lệ",
		})
	}

	// Xóa tất cả role permission cũ của role
	filter := bson.M{"roleId": roleId}
	_, err = h.RolePermissionService.DeleteMany(c.Context(), filter)
	if err != nil {
		if customErr, ok := err.(*utility.Error); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"code":    customErr.Code,
				"message": customErr.Message,
				"details": customErr.Details,
			})
		}
		return c.Status(utility.StatusInternalServerError).JSON(fiber.Map{
			"code":    utility.ErrCodeDatabase,
			"message": err.Error(),
		})
	}

	// Tạo danh sách role permission mới
	var rolePermissions []models.RolePermission
	for _, permissionId := range input.PermissionIds {
		permissionIdObj, err := primitive.ObjectIDFromHex(permissionId)
		if err != nil {
			continue // Bỏ qua các permissionId không hợp lệ
		}
		rolePermission := models.RolePermission{
			ID:           primitive.NewObjectID(),
			RoleID:       roleId,
			PermissionID: permissionIdObj,
			Scope:        0,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		}
		rolePermissions = append(rolePermissions, rolePermission)
	}

	// Thêm các role permission mới
	for _, rolePermission := range rolePermissions {
		_, err = h.RolePermissionService.InsertOne(c.Context(), rolePermission)
		if err != nil {
			if customErr, ok := err.(*utility.Error); ok {
				return c.Status(customErr.StatusCode).JSON(fiber.Map{
					"code":    customErr.Code,
					"message": customErr.Message,
					"details": customErr.Details,
				})
			}
			return c.Status(utility.StatusInternalServerError).JSON(fiber.Map{
				"code":    utility.ErrCodeDatabase,
				"message": err.Error(),
			})
		}
	}

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    rolePermissions,
	})
}
