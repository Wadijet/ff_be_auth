package handler

// Package handler chứa các handler xử lý request HTTP trong ứng dụng.
// Package này cung cấp các chức năng CRUD cơ bản và các tiện ích để xử lý request/response.

import (
	"bytes"
	"encoding/json"
	"errors"
	"meta_commerce/core/api/services"
	"meta_commerce/core/global"
	"meta_commerce/core/utility"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseHandler là base handler cho các Fiber handler, cung cấp các chức năng CRUD cơ bản.
// Struct này sử dụng Generic Type để có thể tái sử dụng cho nhiều loại model khác nhau.
//
// Type parameters:
// - T: Kiểu dữ liệu của model
// - CreateInput: Kiểu dữ liệu của input khi tạo mới
// - UpdateInput: Kiểu dữ liệu của input khi cập nhật
type BaseHandler[T any, CreateInput any, UpdateInput any] struct {
	BaseService services.BaseServiceMongo[T] // Service xử lý logic nghiệp vụ với MongoDB
}

// HandleResponse xử lý và chuẩn hóa response trả về cho client.
// Phương thức này đảm bảo format response thống nhất trong toàn bộ ứng dụng.
//
// Parameters:
// - c: Fiber context
// - data: Dữ liệu trả về cho client (có thể là nil nếu chỉ trả về lỗi)
// - err: Lỗi nếu có (nil nếu không có lỗi)
func (h *BaseHandler[T, CreateInput, UpdateInput]) HandleResponse(c fiber.Ctx, data interface{}, err error) {
	if err != nil {
		var customErr *utility.Error
		if errors.As(err, &customErr) {
			c.Status(customErr.StatusCode).JSON(fiber.Map{
				"code":    customErr.Code.Code,
				"message": customErr.Message,
				"details": customErr.Details,
				"status":  "error",
			})
			return
		}
		// Nếu không phải custom error, trả về internal server error
		c.Status(utility.StatusInternalServerError).JSON(fiber.Map{
			"code":    utility.ErrCodeDatabase.Code,
			"message": err.Error(),
			"status":  "error",
		})
		return
	}

	// Trường hợp thành công
	c.Status(utility.StatusOK).JSON(fiber.Map{
		"code":    utility.StatusOK,
		"message": utility.MsgSuccess,
		"data":    data,
		"status":  "success",
	})
}

// ParseRequestBody parse và validate dữ liệu từ request body.
// Sử dụng json.Decoder với UseNumber() để xử lý chính xác các số.
//
// Parameters:
// - c: Fiber context
// - input: Con trỏ tới struct sẽ chứa dữ liệu được parse
//
// Returns:
// - error: Lỗi nếu có trong quá trình parse hoặc validate
func (h *BaseHandler[T, CreateInput, UpdateInput]) ParseRequestBody(c fiber.Ctx, input interface{}) error {
	// Parse body thành struct T
	body := c.Body()
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationFormat, utility.MsgValidationError, utility.StatusBadRequest, err)
	}

	// Validate struct input xem có hợp lệ không
	if err := global.Validate.Struct(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationInput, utility.MsgValidationError, utility.StatusBadRequest, err)
	}

	return nil
}

// ParseRequestQuery parse và validate dữ liệu từ query string.
// Query string phải được encode dưới dạng JSON.
//
// Parameters:
// - c: Fiber context
// - input: Con trỏ tới struct sẽ chứa dữ liệu được parse
//
// Returns:
// - error: Lỗi nếu có trong quá trình parse hoặc validate
func (h *BaseHandler[T, CreateInput, UpdateInput]) ParseRequestQuery(c fiber.Ctx, input interface{}) error {
	query := c.Query("query", "")

	// Parse query
	reader := bytes.NewReader([]byte(query))
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationFormat, utility.MsgValidationError, utility.StatusBadRequest, err)
	}

	// Validate struct
	if err := global.Validate.Struct(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationInput, utility.MsgValidationError, utility.StatusBadRequest, err)
	}

	return nil
}

// ParseRequestParams parse và validate các tham số từ URI.
// Sử dụng Fiber's URI binding để parse các tham số.
//
// Parameters:
// - c: Fiber context
// - input: Con trỏ tới struct sẽ chứa dữ liệu được parse
//
// Returns:
// - error: Lỗi nếu có trong quá trình parse hoặc validate
func (h *BaseHandler[T, CreateInput, UpdateInput]) ParseRequestParams(c fiber.Ctx, input interface{}) error {
	// Parse URI params
	if err := c.Bind().URI(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationFormat, utility.MsgValidationError, utility.StatusBadRequest, err)
	}

	// Validate struct
	if err := global.Validate.Struct(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationInput, utility.MsgValidationError, utility.StatusBadRequest, err)
	}

	return nil
}

// InsertOne thêm mới một document vào database.
// Dữ liệu được parse từ request body và validate trước khi thêm vào DB.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) InsertOne(c fiber.Ctx) error {

	// Parse request body thành struct T
	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	data, err := h.BaseService.InsertOne(c.Context(), *input)
	h.HandleResponse(c, data, err)
	return nil
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
	var inputs []T
	if err := h.ParseRequestBody(c, &inputs); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	data, err := h.BaseService.InsertMany(c.Context(), inputs)
	h.HandleResponse(c, data, err)
	return nil
}

// FindOne tìm một document theo điều kiện filter.
// Filter được truyền qua query string dưới dạng JSON.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) FindOne(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.BaseService.FindOne(c.Context(), filter, nil)
	h.HandleResponse(c, data, err)
	return nil
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
	id := c.Params("id")
	if id == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.BaseService.FindOneById(c.Context(), utility.String2ObjectID(id))
	h.HandleResponse(c, data, err)
	return nil
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
	var ids []string
	if err := json.Unmarshal([]byte(c.Query("ids", "[]")), &ids); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Danh sách ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	objectIds := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectIds[i] = utility.String2ObjectID(id)
	}

	data, err := h.BaseService.FindManyByIds(c.Context(), objectIds)
	h.HandleResponse(c, data, err)
	return nil
}

// FindWithPagination tìm nhiều document với phân trang.
// Hỗ trợ filter và phân trang với page và limit.
//
// Parameters:
// - c: Fiber context
// Query params:
// - filter: Điều kiện tìm kiếm (JSON)
// - page: Số trang (mặc định: 1)
// - limit: Số lượng item trên một trang (mặc định: 10)
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) FindWithPagination(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, utility.MsgValidationError, utility.StatusBadRequest, nil))
		return nil
	}

	// Parse page và limit từ query string
	page, err := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
	if err != nil {
		limit = 10
	}

	data, err := h.BaseService.FindWithPagination(c.Context(), filter, page, limit)
	h.HandleResponse(c, data, err)
	return nil
}

// Find tìm nhiều document theo điều kiện filter.
// Filter được truyền qua query string dưới dạng JSON.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) Find(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.BaseService.Find(c.Context(), filter, nil)
	h.HandleResponse(c, data, err)
	return nil
}

// UpdateOne cập nhật một document theo điều kiện filter.
// Filter được truyền qua query string, dữ liệu cập nhật trong request body.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) UpdateOne(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	data, err := h.BaseService.UpdateOne(c.Context(), filter, input, nil)
	h.HandleResponse(c, data, err)
	return nil
}

// UpdateMany cập nhật nhiều document theo điều kiện filter.
// Filter được truyền qua query string, dữ liệu cập nhật trong request body.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) UpdateMany(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	count, err := h.BaseService.UpdateMany(c.Context(), filter, input, nil)
	h.HandleResponse(c, count, err)
	return nil
}

// UpdateById cập nhật một document theo ID.
// ID được truyền qua URI params, dữ liệu cập nhật trong request body.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) UpdateById(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	data, err := h.BaseService.UpdateById(c.Context(), utility.String2ObjectID(id), *input)
	h.HandleResponse(c, data, err)
	return nil
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
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	err := h.BaseService.DeleteOne(c.Context(), filter)
	h.HandleResponse(c, nil, err)
	return nil
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
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	count, err := h.BaseService.DeleteMany(c.Context(), filter)
	h.HandleResponse(c, count, err)
	return nil
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
	id := c.Params("id")
	if id == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	err := h.BaseService.DeleteById(c.Context(), utility.String2ObjectID(id))
	h.HandleResponse(c, nil, err)
	return nil
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
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	data, err := h.BaseService.FindOneAndUpdate(c.Context(), filter, input, nil)
	h.HandleResponse(c, data, err)
	return nil
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
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.BaseService.FindOneAndDelete(c.Context(), filter, nil)
	h.HandleResponse(c, data, err)
	return nil
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
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	count, err := h.BaseService.CountDocuments(c.Context(), filter)
	h.HandleResponse(c, count, err)
	return nil
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
	field := c.Params("field")
	if field == "" {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Tên trường không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.BaseService.Distinct(c.Context(), field, filter)
	h.HandleResponse(c, data, err)
	return nil
}

// Upsert thêm mới hoặc cập nhật một document.
// Filter được truyền qua query string, dữ liệu trong request body.
// Nếu không tìm thấy document thỏa mãn filter sẽ tạo mới, ngược lại sẽ cập nhật.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) Upsert(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	data, err := h.BaseService.Upsert(c.Context(), filter, *input)
	h.HandleResponse(c, data, err)
	return nil
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
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
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
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleResponse(c, nil, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	exists, err := h.BaseService.DocumentExists(c.Context(), filter)
	h.HandleResponse(c, exists, err)
	return nil
}

// ParsePagination xử lý việc parse thông tin phân trang từ request.
// Hỗ trợ các tham số:
// - page: Số trang (mặc định: 1)
// - limit: Số lượng item trên một trang (mặc định: 10)
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - page: Số trang
// - limit: Số lượng item trên một trang
func (h *BaseHandler[T, CreateInput, UpdateInput]) ParsePagination(c fiber.Ctx) (int64, int64) {
	page := utility.P2Int64(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}

	limit := utility.P2Int64(c.Query("limit", "10"))
	if limit <= 0 {
		limit = 10
	}

	return page, limit
}

// GetIDFromContext lấy ID từ URI params của request.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - string: ID từ params
func (h *BaseHandler[T, CreateInput, UpdateInput]) GetIDFromContext(c fiber.Ctx) string {
	return c.Params("id")
}

// GenericHandler xử lý request theo một flow chung.
// Flow bao gồm:
// 1. Parse và validate input từ request body
// 2. Gọi handler function với input đã parse
// 3. Xử lý response
//
// Parameters:
// - c: Fiber context
// - input: Struct chứa dữ liệu input
// - handler: Function xử lý logic chính
//
// Returns:
// - error: Lỗi nếu có
func (h *BaseHandler[T, CreateInput, UpdateInput]) GenericHandler(c fiber.Ctx, input interface{}, handler func(c fiber.Ctx, input interface{}) (interface{}, error)) error {
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleResponse(c, nil, err)
		return nil
	}

	result, err := handler(c, input)
	h.HandleResponse(c, result, err)
	return nil
}
