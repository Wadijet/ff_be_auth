package handler

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
	"meta_commerce/core/utility"

	"github.com/gofiber/fiber/v3"
)

// PermissionHandler xử lý các route liên quan đến permission cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
// Các phương thức của FiberBaseHandler đã có sẵn:
// - InsertOne: Thêm mới một permission
// - InsertMany: Thêm nhiều permission
// - FindOne: Tìm một permission theo điều kiện
// - FindOneById: Tìm một permission theo ID
// - FindManyByIds: Tìm nhiều permission theo danh sách ID
// - FindWithPagination: Tìm permission với phân trang
// - Find: Tìm nhiều permission theo điều kiện
// - UpdateOne: Cập nhật một permission theo điều kiện
// - UpdateMany: Cập nhật nhiều permission theo điều kiện
// - UpdateById: Cập nhật một permission theo ID
// - DeleteOne: Xóa một permission theo điều kiện
// - DeleteMany: Xóa nhiều permission theo điều kiện
// - DeleteById: Xóa một permission theo ID
// - FindOneAndUpdate: Tìm và cập nhật một permission
// - FindOneAndDelete: Tìm và xóa một permission
// - CountDocuments: Đếm số lượng permission theo điều kiện
// - Distinct: Lấy danh sách giá trị duy nhất của một trường
// - Upsert: Thêm mới hoặc cập nhật một permission
// - UpsertMany: Thêm mới hoặc cập nhật nhiều permission
// - DocumentExists: Kiểm tra permission có tồn tại không
type PermissionHandler struct {
	BaseHandler[models.Permission, models.PermissionCreateInput, models.PermissionUpdateInput]
}

// NewPermissionHandler tạo một instance mới của FiberPermissionHandler
// Returns:
//   - *FiberPermissionHandler: Instance mới của FiberPermissionHandler đã được khởi tạo với PermissionService
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewPermissionHandler() (*PermissionHandler, error) {
	handler := &PermissionHandler{}

	// Khởi tạo PermissionService
	permissionService, err := services.NewPermissionService()
	if err != nil {
		return nil, fmt.Errorf("failed to create permission service: %v", err)
	}

	handler.BaseService = permissionService
	return handler, nil
}

// HandleCreatePermission xử lý tạo mới permission
func (h *PermissionHandler) HandleCreatePermission(c fiber.Ctx) error {
	input := new(models.PermissionCreateInput)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	// Chuyển đổi từ PermissionCreateInput sang Permission
	permission := models.Permission{
		Name:     input.Name,
		Describe: input.Describe,
		Category: input.Category,
		Group:    input.Group,
	}

	data, err := h.BaseService.InsertOne(c.Context(), permission)
	h.HandleResponse(c, data, err)
	return nil
}

// HandleUpdatePermission xử lý cập nhật permission
func (h *PermissionHandler) HandleUpdatePermission(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	if id == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(models.PermissionUpdateInput)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	// Chuyển đổi từ PermissionUpdateInput sang Permission
	permission := models.Permission{
		Name:     input.Name,
		Describe: input.Describe,
		Category: input.Category,
		Group:    input.Group,
	}

	data, err := h.BaseService.UpdateById(c.Context(), utility.String2ObjectID(id), permission)
	h.HandleResponse(c, data, err)
	return nil
}

// HandleGetPermissionById xử lý lấy thông tin permission theo ID
func (h *PermissionHandler) HandleGetPermissionById(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	if id == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.BaseService.FindOneById(c.Context(), utility.String2ObjectID(id))
	h.HandleResponse(c, data, err)
	return nil
}

// HandleGetPermissions xử lý lấy danh sách permission với phân trang
func (h *PermissionHandler) HandleGetPermissions(c fiber.Ctx) error {
	// Parse filter từ query params
	var filter map[string]interface{}
	if err := c.Bind().Query(&filter); err != nil {
		filter = make(map[string]interface{})
	}

	// Lấy thông tin phân trang
	page, limit := h.ParsePagination(c)

	data, err := h.BaseService.FindWithPagination(c.Context(), filter, page, limit)
	h.HandleResponse(c, data, err)
	return nil
}

// HandleDeletePermission xử lý xóa permission
func (h *PermissionHandler) HandleDeletePermission(c fiber.Ctx) error {
	id := h.GetIDFromContext(c)
	if id == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	err := h.BaseService.DeleteById(c.Context(), utility.String2ObjectID(id))
	h.HandleResponse(c, nil, err)
	return nil
}

// HandleGetPermissionsByCategory xử lý lấy danh sách permission theo category
func (h *PermissionHandler) HandleGetPermissionsByCategory(c fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Category không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	filter := map[string]interface{}{
		"category": category,
	}

	data, err := h.BaseService.Find(c.Context(), filter, nil)
	h.HandleResponse(c, data, err)
	return nil
}

// HandleGetPermissionsByGroup xử lý lấy danh sách permission theo group
func (h *PermissionHandler) HandleGetPermissionsByGroup(c fiber.Ctx) error {
	group := c.Params("group")
	if group == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Group không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	filter := map[string]interface{}{
		"group": group,
	}

	data, err := h.BaseService.Find(c.Context(), filter, nil)
	h.HandleResponse(c, data, err)
	return nil
}
