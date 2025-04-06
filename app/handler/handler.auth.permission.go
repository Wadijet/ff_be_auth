package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"strconv"

	"github.com/gofiber/fiber/v3"
)

// FiberPermissionHandler xử lý các route liên quan đến permission cho Fiber
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
type FiberPermissionHandler struct {
	FiberBaseHandler[models.Permission, models.PermissionCreateInput, models.PermissionUpdateInput]
}

// NewFiberPermissionHandler tạo một instance mới của FiberPermissionHandler
// Returns:
//   - *FiberPermissionHandler: Instance mới của FiberPermissionHandler đã được khởi tạo với PermissionService
func NewFiberPermissionHandler() *FiberPermissionHandler {
	handler := &FiberPermissionHandler{}
	handler.Service = services.NewPermissionService()
	return handler
}

// HandleCreatePermission xử lý tạo mới permission
func (h *FiberPermissionHandler) HandleCreatePermission(c fiber.Ctx) error {
	input := new(models.PermissionCreateInput)
	if err := c.Bind().Body(input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": utility.MsgValidationError,
		})
	}

	// Chuyển đổi từ PermissionCreateInput sang Permission
	permission := models.Permission{
		Name:     input.Name,
		Describe: input.Describe,
		Category: input.Category,
		Group:    input.Group,
	}

	data, err := h.Service.InsertOne(c.Context(), permission)
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

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}

// HandleUpdatePermission xử lý cập nhật permission
func (h *FiberPermissionHandler) HandleUpdatePermission(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": "ID không hợp lệ",
		})
	}

	input := new(models.PermissionUpdateInput)
	if err := c.Bind().Body(input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": utility.MsgValidationError,
		})
	}

	// Chuyển đổi từ PermissionUpdateInput sang Permission
	permission := models.Permission{
		Name:     input.Name,
		Describe: input.Describe,
		Category: input.Category,
		Group:    input.Group,
	}

	data, err := h.Service.UpdateById(c.Context(), utility.String2ObjectID(id), permission)
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

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}

// HandleGetPermissionById xử lý lấy thông tin permission theo ID
func (h *FiberPermissionHandler) HandleGetPermissionById(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": "ID không hợp lệ",
		})
	}

	data, err := h.Service.FindOneById(c.Context(), utility.String2ObjectID(id))
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

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}

// HandleGetPermissions xử lý lấy danh sách permission với phân trang
func (h *FiberPermissionHandler) HandleGetPermissions(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := c.Bind().Query(&filter); err != nil {
		filter = make(map[string]interface{})
	}

	// Parse page và limit từ query string
	page, err := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	data, err := h.Service.FindWithPagination(c.Context(), filter, page, limit)
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

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}

// HandleDeletePermission xử lý xóa permission
func (h *FiberPermissionHandler) HandleDeletePermission(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": "ID không hợp lệ",
		})
	}

	err := h.Service.DeleteById(c.Context(), utility.String2ObjectID(id))
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

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
	})
}

// HandleGetPermissionsByCategory xử lý lấy danh sách permission theo category
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//   - category: Tên category cần tìm (lấy từ params)
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Thành công
//     {
//     "message": "Thành công",
//     "data": [
//     {
//     "id": "...",
//     "name": "...",
//     "describe": "...",
//     "category": "...",
//     "group": "...",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     ]
//     }
//   - 400: Category không hợp lệ
//   - 500: Lỗi server
func (h *FiberPermissionHandler) HandleGetPermissionsByCategory(c fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": "Category không hợp lệ",
		})
	}

	filter := map[string]interface{}{
		"category": category,
	}

	data, err := h.Service.Find(c.Context(), filter, nil)
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

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}

// HandleGetPermissionsByGroup xử lý lấy danh sách permission theo group
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//   - group: Tên group cần tìm (lấy từ params)
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Thành công
//     {
//     "message": "Thành công",
//     "data": [
//     {
//     "id": "...",
//     "name": "...",
//     "describe": "...",
//     "category": "...",
//     "group": "...",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     ]
//     }
//   - 400: Group không hợp lệ
//   - 500: Lỗi server
func (h *FiberPermissionHandler) HandleGetPermissionsByGroup(c fiber.Ctx) error {
	group := c.Params("group")
	if group == "" {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": "Group không hợp lệ",
		})
	}

	filter := map[string]interface{}{
		"group": group,
	}

	data, err := h.Service.Find(c.Context(), filter, nil)
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

	return c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}
