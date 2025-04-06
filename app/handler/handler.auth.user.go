package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"

	"github.com/gofiber/fiber/v3"
)

// FiberAuthUserHandler xử lý các route liên quan đến xác thực người dùng cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
// Các phương thức của FiberBaseHandler đã có sẵn:
// - InsertOne: Thêm mới một user
// - InsertMany: Thêm nhiều user
// - FindOne: Tìm một user theo điều kiện
// - FindOneById: Tìm một user theo ID
// - FindManyByIds: Tìm nhiều user theo danh sách ID
// - FindWithPagination: Tìm user với phân trang
// - Find: Tìm nhiều user theo điều kiện
// - UpdateOne: Cập nhật một user theo điều kiện
// - UpdateMany: Cập nhật nhiều user theo điều kiện
// - UpdateById: Cập nhật một user theo ID
// - DeleteOne: Xóa một user theo điều kiện
// - DeleteMany: Xóa nhiều user theo điều kiện
// - DeleteById: Xóa một user theo ID
// - FindOneAndUpdate: Tìm và cập nhật một user
// - FindOneAndDelete: Tìm và xóa một user
// - CountDocuments: Đếm số lượng user theo điều kiện
// - Distinct: Lấy danh sách giá trị duy nhất của một trường
// - Upsert: Thêm mới hoặc cập nhật một user
// - UpsertMany: Thêm mới hoặc cập nhật nhiều user
// - DocumentExists: Kiểm tra user có tồn tại không
type FiberAuthUserHandler struct {
	FiberBaseHandler[models.User, models.UserCreateInput, models.User]
	UserService *services.UserService
	RoleService *services.RoleService
}

// NewFiberAuthUserHandler tạo một instance mới của FiberAuthUserHandler
// Returns:
//   - *FiberAuthUserHandler: Instance mới của FiberAuthUserHandler đã được khởi tạo với UserService và RoleService
func NewFiberAuthUserHandler() *FiberAuthUserHandler {
	handler := &FiberAuthUserHandler{
		UserService: services.NewUserService(),
		RoleService: services.NewRoleService(),
	}
	handler.Service = handler.UserService
	return handler
}

// HandleLogin xử lý đăng nhập
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - email: Email đăng nhập
//   - password: Mật khẩu
//
// Response:
//   - 200: Đăng nhập thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "token": "...",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 400: Dữ liệu không hợp lệ
//   - 401: Email hoặc mật khẩu không đúng
//   - 500: Lỗi server
func (h *FiberAuthUserHandler) HandleLogin(c fiber.Ctx) error {
	input := new(models.UserLoginInput)
	if err := c.Bind().Body(input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": utility.MsgValidationError,
		})
	}

	data, err := h.UserService.Login(c.Context(), input)
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

// HandleRegister xử lý đăng ký tài khoản mới
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - email: Email đăng ký
//   - password: Mật khẩu
//   - name: Tên người dùng
//
// Response:
//   - 200: Đăng ký thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 400: Dữ liệu không hợp lệ
//   - 409: Email đã tồn tại
//   - 500: Lỗi server
func (h *FiberAuthUserHandler) HandleRegister(c fiber.Ctx) error {
	input := new(models.UserCreateInput)
	if err := c.Bind().Body(input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": utility.MsgValidationError,
		})
	}

	data, err := h.UserService.Registry(c.Context(), input)
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

// HandleLogout xử lý đăng xuất
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - deviceId: ID của thiết bị đăng xuất (tùy chọn)
//
// Response:
//   - 200: Đăng xuất thành công
//     {
//     "message": "Thành công"
//     }
//   - 400: Dữ liệu không hợp lệ
//   - 401: Chưa đăng nhập
//   - 500: Lỗi server
func (h *FiberAuthUserHandler) HandleLogout(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(utility.StatusUnauthorized).JSON(fiber.Map{
			"code":    utility.ErrCodeAuthCredentials,
			"message": utility.MsgUnauthorized,
		})
	}

	input := new(models.UserLogoutInput)
	if err := c.Bind().Body(input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": utility.MsgValidationError,
		})
	}

	err := h.UserService.Logout(c.Context(), utility.String2ObjectID(userID.(string)), input)
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

// HandleGetMyInfo xử lý lấy thông tin cá nhân của người dùng đang đăng nhập
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Lấy thông tin thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 401: Chưa đăng nhập
//   - 404: Không tìm thấy thông tin người dùng
//   - 500: Lỗi server
func (h *FiberAuthUserHandler) HandleGetMyInfo(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(utility.StatusUnauthorized).JSON(fiber.Map{
			"code":    utility.ErrCodeAuthCredentials,
			"message": utility.MsgUnauthorized,
		})
	}

	data, err := h.UserService.FindOneById(c.Context(), utility.String2ObjectID(userID.(string)))
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

// HandleChangePassword xử lý thay đổi mật khẩu
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - oldPassword: Mật khẩu cũ
//   - newPassword: Mật khẩu mới
//
// Response:
//   - 200: Thay đổi mật khẩu thành công
//     {
//     "message": "Thành công"
//     }
//   - 400: Dữ liệu không hợp lệ
//   - 401: Chưa đăng nhập hoặc mật khẩu cũ không đúng
//   - 500: Lỗi server
func (h *FiberAuthUserHandler) HandleChangePassword(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(utility.StatusUnauthorized).JSON(fiber.Map{
			"code":    utility.ErrCodeAuthCredentials,
			"message": utility.MsgUnauthorized,
		})
	}

	input := new(models.UserChangePasswordInput)
	if err := c.Bind().Body(input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": utility.MsgValidationError,
		})
	}

	err := h.UserService.ChangePassword(c.Context(), utility.String2ObjectID(userID.(string)), input)
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

// HandleGetMyRoles xử lý lấy danh sách vai trò của người dùng đang đăng nhập
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Response:
//   - 200: Lấy danh sách vai trò thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "name": "...",
//     "permissions": ["..."],
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 401: Chưa đăng nhập
//   - 404: Không tìm thấy thông tin người dùng hoặc vai trò
//   - 500: Lỗi server
func (h *FiberAuthUserHandler) HandleGetMyRoles(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(utility.StatusUnauthorized).JSON(fiber.Map{
			"code":    utility.ErrCodeAuthCredentials,
			"message": utility.MsgUnauthorized,
		})
	}

	user, err := h.UserService.FindOneById(c.Context(), utility.String2ObjectID(userID.(string)))
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

	if user.Token == "" {
		return c.Status(utility.StatusOK).JSON(fiber.Map{
			"message": utility.MsgSuccess,
		})
	}

	role, err := h.RoleService.FindOneById(c.Context(), utility.String2ObjectID(user.Token))
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
		"data":    role,
	})
}

// HandleChangeInfo xử lý thay đổi thông tin cá nhân
// Parameters:
//   - c: Context của Fiber chứa thông tin request
//
// Returns:
//   - error: Lỗi nếu có
//
// Request Body:
//   - name: Tên mới của người dùng
//
// Response:
//   - 200: Thay đổi thông tin thành công
//     {
//     "message": "Thành công",
//     "data": {
//     "id": "...",
//     "email": "...",
//     "name": "...",
//     "createdAt": 123,
//     "updatedAt": 123
//     }
//     }
//   - 400: Dữ liệu không hợp lệ
//   - 401: Chưa đăng nhập
//   - 500: Lỗi server
func (h *FiberAuthUserHandler) HandleChangeInfo(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Status(utility.StatusUnauthorized).JSON(fiber.Map{
			"code":    utility.ErrCodeAuthCredentials,
			"message": utility.MsgUnauthorized,
		})
	}

	input := new(models.UserChangeInfoInput)
	if err := c.Bind().Body(input); err != nil {
		return c.Status(utility.StatusBadRequest).JSON(fiber.Map{
			"code":    utility.ErrCodeValidationFormat,
			"message": utility.MsgValidationError,
		})
	}

	user := models.User{
		Name: input.Name,
	}
	data, err := h.UserService.UpdateById(c.Context(), utility.String2ObjectID(userID.(string)), user)
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
