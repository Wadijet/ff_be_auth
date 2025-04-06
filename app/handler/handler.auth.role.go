package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// FiberRoleHandler xử lý các route liên quan đến vai trò cho Fiber
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
type FiberRoleHandler struct {
	FiberBaseHandler[models.Role, models.RoleCreateInput, models.Role]
	RoleService *services.RoleService
}

// NewFiberRoleHandler tạo một instance mới của FiberRoleHandler
// Returns:
//   - *FiberRoleHandler: Instance mới của FiberRoleHandler đã được khởi tạo với RoleService
func NewFiberRoleHandler() *FiberRoleHandler {
	handler := &FiberRoleHandler{
		RoleService: services.NewRoleService(),
	}
	handler.Service = handler.RoleService
	return handler
}

// Các hàm đặc thù của Role (nếu có) sẽ được thêm vào đây
