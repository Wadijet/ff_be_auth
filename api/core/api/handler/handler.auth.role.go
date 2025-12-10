package handler

import (
	"fmt"
	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// RoleHandler xử lý các route liên quan đến vai trò cho Fiber
// Kế thừa từ FiberBaseHandler để có các chức năng CRUD cơ bản
// Các phương thức của FiberBaseHandler đã có sẵn:
// - InsertOne: Thêm mới một vai trò
// - InsertMany: Thêm nhiều vai trò
// - FindOne: Tìm một vai trò theo điều kiện
// - FindOneById: Tìm một vai trò theo ID
// - FindManyByIds: Tìm nhiều vai trò theo danh sách ID
// - FindWithPagination: Tìm vai trò với phân trang
// - Find: Tìm nhiều vai trò theo điều kiện
// - UpdateOne: Cập nhật một vai trò theo điều kiện
// - UpdateMany: Cập nhật nhiều vai trò theo điều kiện
// - UpdateById: Cập nhật một vai trò theo ID
// - DeleteOne: Xóa một vai trò theo điều kiện
// - DeleteMany: Xóa nhiều vai trò theo điều kiện
// - DeleteById: Xóa một vai trò theo ID
// - FindOneAndUpdate: Tìm và cập nhật một vai trò
// - FindOneAndDelete: Tìm và xóa một vai trò
// - CountDocuments: Đếm số lượng vai trò theo điều kiện
// - Distinct: Lấy danh sách giá trị duy nhất của một trường
// - Upsert: Thêm mới hoặc cập nhật một vai trò
// - UpsertMany: Thêm mới hoặc cập nhật nhiều vai trò
// - DocumentExists: Kiểm tra vai trò có tồn tại không
type RoleHandler struct {
	BaseHandler[models.Role, dto.RoleCreateInput, dto.RoleUpdateInput]
	RoleService *services.RoleService
}

// NewRoleHandler tạo một instance mới của FiberRoleHandler
// Returns:
//   - *FiberRoleHandler: Instance mới của FiberRoleHandler đã được khởi tạo với RoleService
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewRoleHandler() (*RoleHandler, error) {
	// Khởi tạo RoleService
	roleService, err := services.NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}

	handler := &RoleHandler{
		RoleService: roleService,
	}
	handler.BaseService = handler.RoleService
	return handler, nil
}

// Các hàm đặc thù của Role (nếu có) sẽ được thêm vào đây
