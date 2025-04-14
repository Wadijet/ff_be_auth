package handler

import (
	"fmt"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// UserRoleHandler xử lý các route liên quan đến vai trò của người dùng cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
// Các phương thức của FiberBaseHandler đã có sẵn:
// - InsertOne: Thêm mới một user role
// - InsertMany: Thêm nhiều user role
// - FindOne: Tìm một user role theo điều kiện
// - FindOneById: Tìm một user role theo ID
// - FindManyByIds: Tìm nhiều user role theo danh sách ID
// - FindWithPagination: Tìm user role với phân trang
// - Find: Tìm nhiều user role theo điều kiện
// - UpdateOne: Cập nhật một user role theo điều kiện
// - UpdateMany: Cập nhật nhiều user role theo điều kiện
// - UpdateById: Cập nhật một user role theo ID
// - DeleteOne: Xóa một user role theo điều kiện
// - DeleteMany: Xóa nhiều user role theo điều kiện
// - DeleteById: Xóa một user role theo ID
// - FindOneAndUpdate: Tìm và cập nhật một user role
// - FindOneAndDelete: Tìm và xóa một user role
// - CountDocuments: Đếm số lượng user role theo điều kiện
// - Distinct: Lấy danh sách giá trị duy nhất của một trường
// - Upsert: Thêm mới hoặc cập nhật một user role
// - UpsertMany: Thêm mới hoặc cập nhật nhiều user role
// - DocumentExists: Kiểm tra user role có tồn tại không
type UserRoleHandler struct {
	BaseHandler[models.UserRole, models.UserRoleCreateInput, models.UserRoleCreateInput]
	UserRoleService *services.UserRoleService
}

// NewUserRoleHandler tạo một instance mới của FiberUserRoleHandler
// Returns:
//   - *FiberUserRoleHandler: Instance mới của FiberUserRoleHandler đã được khởi tạo với UserRoleService
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewUserRoleHandler() (*UserRoleHandler, error) {
	// Khởi tạo UserRoleService
	userRoleService, err := services.NewUserRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create user role service: %v", err)
	}

	handler := &UserRoleHandler{
		UserRoleService: userRoleService,
	}
	handler.BaseService = handler.UserRoleService
	return handler, nil
}

// Các hàm đặc thù của UserRole (nếu có) sẽ được thêm vào đây
