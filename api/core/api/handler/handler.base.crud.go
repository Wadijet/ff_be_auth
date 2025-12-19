package handler

// Package handler chứa các handler xử lý request HTTP trong ứng dụng.
// Package này cung cấp các chức năng CRUD cơ bản và các tiện ích để xử lý request/response.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/utility"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoopts "go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne thêm mới một document vào database.
// Dữ liệu được parse từ request body và validate trước khi thêm vào DB.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) InsertOne(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		// Parse request body thành struct T
		input := new(T)
		if err := h.ParseRequestBody(c, input); err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Dữ liệu gửi lên không đúng định dạng JSON hoặc không khớp với cấu trúc yêu cầu. Chi tiết: %v", err),
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		data, err := h.BaseService.InsertOne(c.Context(), *input)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// InsertMany thêm nhiều document vào database.
// Dữ liệu được parse từ request body dưới dạng mảng và validate trước khi thêm vào DB.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) InsertMany(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		var inputs []T
		if err := h.ParseRequestBody(c, &inputs); err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Dữ liệu gửi lên phải là một mảng JSON và các phần tử phải khớp với cấu trúc yêu cầu. Chi tiết: %v", err),
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		data, err := h.BaseService.InsertMany(c.Context(), inputs)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// FindOne tìm một document theo điều kiện filter.
// Filter và options được truyền qua query string dưới dạng JSON.
// Ví dụ options: {"projection": {"field": 1}, "sort": {"field": 1}}
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) FindOne(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		options, err := h.processMongoOptions(c, true)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		data, err := h.BaseService.FindOne(c.Context(), filter, options.(*mongoopts.FindOneOptions))
		h.HandleResponse(c, data, err)
		return nil
	})
}

// FindOneById tìm một document theo ID.
// ID được truyền qua URI params.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) FindOneById(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		id := c.Params("id")
		if id == "" {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				"ID không được để trống trong URL params",
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		if !primitive.IsValidObjectID(id) {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("ID '%s' không đúng định dạng MongoDB ObjectID (phải là chuỗi hex 24 ký tự)", id),
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		data, err := h.BaseService.FindOneById(c.Context(), utility.String2ObjectID(id))
		h.HandleResponse(c, data, err)
		return nil
	})
}

// FindManyByIds tìm nhiều document theo danh sách ID.
// Danh sách ID được truyền qua query string dưới dạng mảng JSON.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) FindManyByIds(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		var ids []string
		idsStr := c.Query("ids", "[]")
		if err := json.Unmarshal([]byte(idsStr), &ids); err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Danh sách ID phải là một mảng JSON. Giá trị nhận được: %s. Chi tiết lỗi: %v", idsStr, err),
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		// Validate từng ID
		objectIds := make([]primitive.ObjectID, len(ids))
		for i, id := range ids {
			if !primitive.IsValidObjectID(id) {
				h.HandleResponse(c, nil, common.NewError(
					common.ErrCodeValidationFormat,
					fmt.Sprintf("ID '%s' tại vị trí %d không đúng định dạng MongoDB ObjectID (phải là chuỗi hex 24 ký tự)", id, i),
					common.StatusBadRequest,
					nil,
				))
				return nil
			}
			objectIds[i] = utility.String2ObjectID(id)
		}

		data, err := h.BaseService.FindManyByIds(c.Context(), objectIds)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// FindWithPagination tìm nhiều document với phân trang.
// Hỗ trợ filter, options và phân trang với page và limit.
//
// Parameters:
// - c: Fiber context
// Query params:
// - filter: Điều kiện tìm kiếm (JSON)
// - options: Tùy chọn tìm kiếm (JSON). Ví dụ: {"projection": {"field": 1}, "sort": {"field": 1}}
// - page: Số trang (mặc định: 1)
// - limit: Số lượng item trên một trang (mặc định: 10)
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) FindWithPagination(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		// Sử dụng processFilter để có normalizeFilter và validate
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		options, err := h.processMongoOptions(c, false)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		// Parse page và limit từ query string
		page, err := strconv.ParseInt(c.Query("page", "1"), 10, 64)
		if err != nil {
			page = 1
		}
		// Đảm bảo page >= 1 để tránh skip âm
		if page < 1 {
			page = 1
		}

		limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
		if err != nil {
			limit = 10
		}
		// Đảm bảo limit > 0
		if limit <= 0 {
			limit = 10
		}

		// Không set limit và skip vào options ở đây
		// Service sẽ tự tính toán và set vào options để đảm bảo tính nhất quán
		findOptions := options.(*mongoopts.FindOptions)

		data, err := h.BaseService.FindWithPagination(c.Context(), filter, page, limit, findOptions)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// Find tìm nhiều document theo điều kiện filter.
// Filter và options được truyền qua query string dưới dạng JSON.
// Ví dụ options: {"projection": {"field": 1}, "sort": {"field": 1}, "limit": 10, "skip": 0}
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) Find(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		options, err := h.processMongoOptions(c, false)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		data, err := h.BaseService.Find(c.Context(), filter, options.(*mongoopts.FindOptions))
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		// Đảm bảo data không bao giờ là nil, luôn trả về mảng rỗng nếu không có kết quả
		if data == nil {
			data = []T{}
		}

		h.HandleResponse(c, data, nil)
		return nil
	})
}

// UpdateOne cập nhật một document theo điều kiện filter.
// Filter được truyền qua query string, dữ liệu cập nhật trong request body.
// Chỉ update các trường có trong input, giữ nguyên các trường khác.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) UpdateOne(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		// Parse input thành map để chỉ update các trường được chỉ định
		var updateData map[string]interface{}
		if err := json.NewDecoder(bytes.NewReader(c.Body())).Decode(&updateData); err != nil {
			h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Dữ liệu cập nhật không hợp lệ", common.StatusBadRequest, nil))
			return nil
		}

		// Tạo update data với $set operator
		update := &services.UpdateData{
			Set: updateData,
		}

		data, err := h.BaseService.UpdateOne(c.Context(), filter, update, nil)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// UpdateMany cập nhật nhiều document theo điều kiện filter.
// Filter được truyền qua query string, dữ liệu cập nhật trong request body.
// Chỉ update các trường có trong input, giữ nguyên các trường khác.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) UpdateMany(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		// Parse input thành map để chỉ update các trường được chỉ định
		var updateData map[string]interface{}
		if err := json.NewDecoder(bytes.NewReader(c.Body())).Decode(&updateData); err != nil {
			h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Dữ liệu cập nhật không hợp lệ", common.StatusBadRequest, nil))
			return nil
		}

		// Tạo update data với $set operator
		update := &services.UpdateData{
			Set: updateData,
		}

		count, err := h.BaseService.UpdateMany(c.Context(), filter, update, nil)
		h.HandleResponse(c, count, err)
		return nil
	})
}

// UpdateById cập nhật một document theo ID.
// ID được truyền qua URI params, dữ liệu cập nhật trong request body.
// Chỉ update các trường có trong input, giữ nguyên các trường khác.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) UpdateById(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		id := c.Params("id")
		if id == "" {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				"ID không được để trống trong URL params",
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		if !primitive.IsValidObjectID(id) {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("ID '%s' không đúng định dạng MongoDB ObjectID (phải là chuỗi hex 24 ký tự)", id),
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		// Parse input thành map để chỉ update các trường được chỉ định
		var updateData map[string]interface{}
		if err := json.NewDecoder(bytes.NewReader(c.Body())).Decode(&updateData); err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Dữ liệu cập nhật phải là một object JSON hợp lệ. Chi tiết lỗi: %v", err),
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		// Tạo update data với $set operator
		update := &services.UpdateData{
			Set: updateData,
		}

		data, err := h.BaseService.UpdateById(c.Context(), utility.String2ObjectID(id), update)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// DeleteOne xóa một document theo điều kiện filter.
// Filter được truyền qua query string dưới dạng JSON.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) DeleteOne(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		err = h.BaseService.DeleteOne(c.Context(), filter)
		h.HandleResponse(c, nil, err)
		return nil
	})
}

// DeleteMany xóa nhiều document theo điều kiện filter.
// Filter được truyền qua query string dưới dạng JSON.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có và số lượng document đã xóa
func (h *BaseHandler[T, CreateInput, UpdateInput]) DeleteMany(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		count, err := h.BaseService.DeleteMany(c.Context(), filter)
		h.HandleResponse(c, count, err)
		return nil
	})
}

// DeleteById xóa một document theo ID.
// ID được truyền qua URI params.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) DeleteById(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		id := c.Params("id")
		if id == "" {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				"ID không được để trống trong URL params",
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		if !primitive.IsValidObjectID(id) {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("ID '%s' không đúng định dạng MongoDB ObjectID (phải là chuỗi hex 24 ký tự)", id),
				common.StatusBadRequest,
				nil,
			))
			return nil
		}

		err := h.BaseService.DeleteById(c.Context(), utility.String2ObjectID(id))
		h.HandleResponse(c, nil, err)
		return nil
	})
}

// FindOneAndUpdate tìm và cập nhật một document.
// Filter được truyền qua query string, dữ liệu cập nhật trong request body.
// Trả về document sau khi cập nhật.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) FindOneAndUpdate(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		// Parse input thành map để chỉ update các trường được chỉ định
		var updateData map[string]interface{}
		if err := json.NewDecoder(bytes.NewReader(c.Body())).Decode(&updateData); err != nil {
			h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Dữ liệu cập nhật không hợp lệ", common.StatusBadRequest, nil))
			return nil
		}

		// Tạo update data với $set operator
		update := &services.UpdateData{
			Set: updateData,
		}

		data, err := h.BaseService.FindOneAndUpdate(c.Context(), filter, update, nil)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// FindOneAndDelete tìm và xóa một document.
// Filter được truyền qua query string.
// Trả về document đã xóa.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) FindOneAndDelete(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		data, err := h.BaseService.FindOneAndDelete(c.Context(), filter, nil)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// CountDocuments đếm số lượng document theo điều kiện filter.
// Filter được truyền qua query string dưới dạng JSON.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) CountDocuments(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		var filter map[string]interface{}
		// Lấy giá trị filter từ query string, mặc định là "{}" nếu không có
		filterStr := c.Query("filter", "{}")

		// Log giá trị filter để debug
		fmt.Printf("Filter string từ query: %s\n", filterStr)

		// Chuyển đổi chuỗi JSON thành map
		if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
			// Log lỗi để debug
			fmt.Printf("Lỗi khi parse filter: %v\n", err)

			// Trả về lỗi cho client
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				"Filter không hợp lệ",
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		// Log filter sau khi parse thành công
		fmt.Printf("Filter sau khi parse: %+v\n", filter)

		count, err := h.BaseService.CountDocuments(c.Context(), filter)
		h.HandleResponse(c, count, err)
		return nil
	})
}

// Distinct lấy danh sách giá trị duy nhất của một trường.
// Tên trường được truyền qua URI params, filter qua query string.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) Distinct(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		field := c.Params("field")
		if field == "" {
			h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Tên trường không hợp lệ", common.StatusBadRequest, nil))
			return nil
		}

		var filter map[string]interface{}
		if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
			h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Filter không hợp lệ", common.StatusBadRequest, nil))
			return nil
		}

		data, err := h.BaseService.Distinct(c.Context(), field, filter)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// Upsert thêm mới hoặc cập nhật một document.
// Filter được truyền qua query string, dữ liệu trong request body.
// Nếu không tìm thấy document thỏa mãn filter sẽ tạo mới, ngược lại sẽ cập nhật.
// Parse body thành struct T để struct tag `extract` có thể hoạt động tự động
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) Upsert(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		// Parse filter từ query string (sử dụng processFilter để có normalizeFilter và validate)
		filter, err := h.processFilter(c)
		if err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		// Parse request body thành struct T (model) để struct tag `extract` có thể hoạt động
		// Struct tag `extract` sẽ tự động extract dữ liệu từ PanCakeData, FacebookData, etc.
		input := new(T)
		if err := h.ParseRequestBody(c, input); err != nil {
			h.HandleResponse(c, nil, common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Dữ liệu gửi lên không đúng định dạng JSON hoặc không khớp với cấu trúc yêu cầu. Chi tiết: %v", err),
				common.StatusBadRequest,
				err,
			))
			return nil
		}

		// Gọi Upsert với struct T - extract sẽ tự động chạy trong ToMap() khi ToUpdateData() được gọi
		data, err := h.BaseService.Upsert(c.Context(), filter, *input)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// UpsertMany thêm mới hoặc cập nhật nhiều document.
// Filter được truyền qua query string, dữ liệu trong request body dưới dạng mảng.
// Với mỗi item trong mảng: nếu không tìm thấy document thỏa mãn filter sẽ tạo mới, ngược lại sẽ cập nhật.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) UpsertMany(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		var filter map[string]interface{}
		if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
			h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Filter không hợp lệ", common.StatusBadRequest, nil))
			return nil
		}

		var inputs []T
		if err := h.ParseRequestBody(c, &inputs); err != nil {
			h.HandleResponse(c, nil, err)
			return nil
		}

		data, err := h.BaseService.UpsertMany(c.Context(), filter, inputs)
		h.HandleResponse(c, data, err)
		return nil
	})
}

// DocumentExists kiểm tra document có tồn tại không.
// Filter được truyền qua query string dưới dạng JSON.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) DocumentExists(c fiber.Ctx) error {
	return h.SafeHandler(c, func() error {
		var filter map[string]interface{}
		if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
			h.HandleResponse(c, nil, common.NewError(common.ErrCodeValidationFormat, "Filter không hợp lệ", common.StatusBadRequest, nil))
			return nil
		}

		exists, err := h.BaseService.DocumentExists(c.Context(), filter)
		h.HandleResponse(c, exists, err)
		return nil
	})
}
