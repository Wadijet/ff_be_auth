package handler

import (
	"encoding/json"
	"meta_commerce/app/services"
	"meta_commerce/app/utility"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FiberBaseHandler là base handler cho các Fiber handler
type FiberBaseHandler[T any, CreateInput any, UpdateInput any] struct {
	Service services.BaseServiceMongo[T]
}

// HandleError xử lý lỗi và trả về response cho client
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) HandleError(c fiber.Ctx, err error) {
	// Kiểm tra nếu là Error từ utility
	if customErr, ok := err.(*utility.Error); ok {
		c.Status(customErr.StatusCode).JSON(fiber.Map{
			"code":    customErr.Code,
			"message": customErr.Message,
			"details": customErr.Details,
		})
		return
	}

	// Nếu không phải Error từ utility, trả về lỗi mặc định
	c.Status(utility.StatusInternalServerError).JSON(fiber.Map{
		"code":    utility.ErrCodeDatabase,
		"message": err.Error(),
	})
}

// HandleResponse xử lý response thành công và trả về cho client
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) HandleResponse(c fiber.Ctx, data interface{}, err error) {
	if err != nil {
		h.HandleError(c, err)
		return
	}

	c.Status(utility.StatusOK).JSON(fiber.Map{
		"message": utility.MsgSuccess,
		"data":    data,
	})
}

// ParseRequestBody parse request body thành struct
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) ParseRequestBody(c fiber.Ctx, input interface{}) error {
	if err := c.Bind().Body(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationFormat, utility.MsgValidationError, utility.StatusBadRequest, nil)
	}
	return nil
}

// ParseRequestQuery parse request query thành struct
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) ParseRequestQuery(c fiber.Ctx, input interface{}) error {
	if err := c.Bind().Query(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationFormat, "Dữ liệu không hợp lệ", utility.StatusBadRequest, nil)
	}
	return nil
}

// ParseRequestParams parse request params thành struct
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) ParseRequestParams(c fiber.Ctx, input interface{}) error {
	if err := c.Bind().URI(input); err != nil {
		return utility.NewError(utility.ErrCodeValidationFormat, "Dữ liệu không hợp lệ", utility.StatusBadRequest, nil)
	}
	return nil
}

// InsertOne thêm mới một document
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) InsertOne(c fiber.Ctx) error {
	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleError(c, err)
		return nil
	}

	data, err := h.Service.InsertOne(c.Context(), *input)
	h.HandleResponse(c, data, err)
	return nil
}

// InsertMany thêm mới nhiều document
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) InsertMany(c fiber.Ctx) error {
	var inputs []T
	if err := h.ParseRequestBody(c, &inputs); err != nil {
		h.HandleError(c, err)
		return nil
	}

	data, err := h.Service.InsertMany(c.Context(), inputs)
	h.HandleResponse(c, data, err)
	return nil
}

// FindOne tìm một document theo điều kiện
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) FindOne(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.Service.FindOne(c.Context(), filter, nil)
	h.HandleResponse(c, data, err)
	return nil
}

// FindOneById tìm một document theo ID
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) FindOneById(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.Service.FindOneById(c.Context(), utility.String2ObjectID(id))
	h.HandleResponse(c, data, err)
	return nil
}

// FindManyByIds tìm nhiều document theo danh sách ID
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) FindManyByIds(c fiber.Ctx) error {
	var ids []string
	if err := json.Unmarshal([]byte(c.Query("ids", "[]")), &ids); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Danh sách ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	objectIds := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectIds[i] = utility.String2ObjectID(id)
	}

	data, err := h.Service.FindManyByIds(c.Context(), objectIds)
	h.HandleResponse(c, data, err)
	return nil
}

// FindWithPagination tìm nhiều document với phân trang
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) FindWithPagination(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, utility.MsgValidationError, utility.StatusBadRequest, nil))
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

	data, err := h.Service.FindWithPagination(c.Context(), filter, page, limit)
	h.HandleResponse(c, data, err)
	return nil
}

// Find tìm nhiều document theo điều kiện
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) Find(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.Service.Find(c.Context(), filter, nil)
	h.HandleResponse(c, data, err)
	return nil
}

// UpdateOne cập nhật một document theo điều kiện
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) UpdateOne(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleError(c, err)
		return nil
	}

	data, err := h.Service.UpdateOne(c.Context(), filter, input, nil)
	h.HandleResponse(c, data, err)
	return nil
}

// UpdateMany cập nhật nhiều document theo điều kiện
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) UpdateMany(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleError(c, err)
		return nil
	}

	count, err := h.Service.UpdateMany(c.Context(), filter, input, nil)
	h.HandleResponse(c, count, err)
	return nil
}

// UpdateById cập nhật một document theo ID
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) UpdateById(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleError(c, err)
		return nil
	}

	data, err := h.Service.UpdateById(c.Context(), utility.String2ObjectID(id), *input)
	h.HandleResponse(c, data, err)
	return nil
}

// DeleteOne xóa một document theo điều kiện
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) DeleteOne(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	err := h.Service.DeleteOne(c.Context(), filter)
	h.HandleResponse(c, nil, err)
	return nil
}

// DeleteMany xóa nhiều document theo điều kiện
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) DeleteMany(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	count, err := h.Service.DeleteMany(c.Context(), filter)
	h.HandleResponse(c, count, err)
	return nil
}

// DeleteById xóa một document theo ID
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) DeleteById(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "ID không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	err := h.Service.DeleteById(c.Context(), utility.String2ObjectID(id))
	h.HandleResponse(c, nil, err)
	return nil
}

// FindOneAndUpdate tìm và cập nhật một document
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) FindOneAndUpdate(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleError(c, err)
		return nil
	}

	data, err := h.Service.FindOneAndUpdate(c.Context(), filter, input, nil)
	h.HandleResponse(c, data, err)
	return nil
}

// FindOneAndDelete tìm và xóa một document
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) FindOneAndDelete(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.Service.FindOneAndDelete(c.Context(), filter, nil)
	h.HandleResponse(c, data, err)
	return nil
}

// CountDocuments đếm số lượng document theo điều kiện
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) CountDocuments(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	count, err := h.Service.CountDocuments(c.Context(), filter)
	h.HandleResponse(c, count, err)
	return nil
}

// Distinct lấy danh sách giá trị duy nhất của một trường
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) Distinct(c fiber.Ctx) error {
	field := c.Params("field")
	if field == "" {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Tên trường không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	data, err := h.Service.Distinct(c.Context(), field, filter)
	h.HandleResponse(c, data, err)
	return nil
}

// Upsert thêm mới hoặc cập nhật một document
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) Upsert(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	input := new(T)
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleError(c, err)
		return nil
	}

	data, err := h.Service.Upsert(c.Context(), filter, *input)
	h.HandleResponse(c, data, err)
	return nil
}

// UpsertMany thêm mới hoặc cập nhật nhiều document
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) UpsertMany(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	var inputs []T
	if err := h.ParseRequestBody(c, &inputs); err != nil {
		h.HandleError(c, err)
		return nil
	}

	data, err := h.Service.UpsertMany(c.Context(), filter, inputs)
	h.HandleResponse(c, data, err)
	return nil
}

// DocumentExists kiểm tra document có tồn tại không
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) DocumentExists(c fiber.Ctx) error {
	var filter map[string]interface{}
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		h.HandleError(c, utility.NewError(utility.ErrCodeValidationFormat, "Filter không hợp lệ", utility.StatusBadRequest, nil))
		return nil
	}

	exists, err := h.Service.DocumentExists(c.Context(), filter)
	h.HandleResponse(c, exists, err)
	return nil
}

// ParsePagination xử lý việc parse thông tin phân trang từ request
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) ParsePagination(c fiber.Ctx) (int, int) {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	return page, limit
}

// GetIDFromContext lấy ID từ context của request
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) GetIDFromContext(c fiber.Ctx) string {
	return c.Params("id")
}

// GenericHandler xử lý request theo một flow chung
func (h *FiberBaseHandler[T, CreateInput, UpdateInput]) GenericHandler(c fiber.Ctx, input interface{}, handler func(c fiber.Ctx, input interface{}) (interface{}, error)) error {
	if err := h.ParseRequestBody(c, input); err != nil {
		h.HandleError(c, err)
		return nil
	}

	result, err := handler(c, input)
	h.HandleResponse(c, result, err)
	return nil
}
