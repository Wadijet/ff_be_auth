package handler

import (
	"context"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/utility"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BaseHandler là struct cơ sở cho tất cả các handler trong hệ thống
// Type Parameters:
//   - Model: Kiểu dữ liệu của model
//   - CreateInput: Kiểu dữ liệu cho input tạo mới
//   - UpdateInput: Kiểu dữ liệu cho input cập nhật
type BaseHandler[Model any, CreateInput any, UpdateInput any] struct {
	// Service tương ứng với model
	Service interface {
		// NHÓM 1: CÁC HÀM CHUẨN MONGODB DRIVER
		// ====================================

		// 1.1 Thao tác Insert
		InsertOne(ctx context.Context, data Model) (Model, error)
		InsertMany(ctx context.Context, data []Model) ([]Model, error)

		// 1.2 Thao tác Find
		FindOne(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (Model, error)
		Find(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]Model, error)

		// 1.3 Thao tác Update
		UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *options.UpdateOptions) (Model, error)
		UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *options.UpdateOptions) (int64, error)

		// 1.4 Thao tác Delete
		DeleteOne(ctx context.Context, filter interface{}) error
		DeleteMany(ctx context.Context, filter interface{}) (int64, error)

		// 1.5 Thao tác Atomic
		FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts *options.FindOneAndUpdateOptions) (Model, error)
		FindOneAndDelete(ctx context.Context, filter interface{}, opts *options.FindOneAndDeleteOptions) (Model, error)

		// 1.6 Các thao tác khác
		CountDocuments(ctx context.Context, filter interface{}) (int64, error)
		Distinct(ctx context.Context, fieldName string, filter interface{}) ([]interface{}, error)

		// NHÓM 2: CÁC HÀM TIỆN ÍCH MỞ RỘNG
		// ================================

		// 2.1 Các hàm Find mở rộng
		FindOneById(ctx context.Context, id primitive.ObjectID) (Model, error)
		FindManyByIds(ctx context.Context, ids []primitive.ObjectID) ([]Model, error)
		FindWithPagination(ctx context.Context, filter interface{}, page, limit int64) (*models.PaginateResult[Model], error)

		// 2.2 Các hàm Update/Delete mở rộng
		UpdateById(ctx context.Context, id primitive.ObjectID, data Model) (Model, error)
		DeleteById(ctx context.Context, id primitive.ObjectID) error

		// 2.3 Các hàm Upsert tiện ích
		Upsert(ctx context.Context, filter interface{}, data Model) (Model, error)
		UpsertMany(ctx context.Context, filter interface{}, data []Model) ([]Model, error)

		// 2.4 Các hàm kiểm tra
		DocumentExists(ctx context.Context, filter interface{}) (bool, error)
	}
}

// HandlerFunc định nghĩa kiểu hàm xử lý logic chính
type HandlerFunc func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error)

// ResponseError là cấu trúc dữ liệu cho response lỗi
type ResponseError struct {
	Status    string            `json:"status"`
	ErrorCode utility.ErrorCode `json:"error_code"`
	Message   string            `json:"message"`
}

// ResponseSuccess là cấu trúc dữ liệu cho response thành công
type ResponseSuccess struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// ====================================
// CÁC HÀM TIỆN ÍCH CHUNG
// ====================================

// HandleResponse xử lý response chung cho tất cả các handler
func (h *BaseHandler[Model, CreateInput, UpdateInput]) HandleResponse(ctx *fasthttp.RequestCtx, data interface{}, err error) {
	if err != nil {
		h.HandleError(ctx, err)
	} else {
		h.HandleSuccess(ctx, data)
	}
}

// ParseRequestBody xử lý việc parse request body thành struct
func (h *BaseHandler[Model, CreateInput, UpdateInput]) ParseRequestBody(ctx *fasthttp.RequestCtx, input interface{}) map[string]interface{} {
	var body []byte
	contentType := string(ctx.Request.Header.ContentType())
	if contentType == "application/json" {
		body = ctx.Request.Body()
	} else {
		body = ctx.PostBody()
	}
	response := utility.Convert2Struct(body, input)
	if response == nil {
		response = utility.ValidateStruct(input)
	}
	return response
}

// ParsePagination xử lý việc parse thông tin phân trang từ request
func (h *BaseHandler[Model, CreateInput, UpdateInput]) ParsePagination(ctx *fasthttp.RequestCtx) (int64, int64) {
	limit, _ := strconv.ParseInt(string(ctx.FormValue("limit")), 10, 64)
	if limit <= 0 {
		limit = 10
	}
	page, _ := strconv.ParseInt(string(ctx.FormValue("page")), 10, 64)
	if page < 0 {
		page = 0
	}
	return page, limit
}

// GetIDFromContext lấy ID từ context của request
func (h *BaseHandler[Model, CreateInput, UpdateInput]) GetIDFromContext(ctx *fasthttp.RequestCtx) string {
	return ctx.UserValue("id").(string)
}

// HandleError xử lý lỗi chung và trả về response lỗi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) HandleError(ctx *fasthttp.RequestCtx, err error) {
	var statusCode int = utility.StatusInternalServerError
	var message string
	var errorCode utility.ErrorCode = utility.ErrCodeDatabaseConnection

	if customErr, ok := err.(*utility.Error); ok {
		statusCode = customErr.StatusCode
		message = customErr.Message
		errorCode = customErr.Code
	} else {
		message = err.Error()
	}

	ctx.SetStatusCode(statusCode)
	response := ResponseError{
		Status:    "error",
		ErrorCode: errorCode,
		Message:   message,
	}
	responseMap := map[string]interface{}{
		"status":     response.Status,
		"error_code": response.ErrorCode,
		"message":    response.Message,
	}
	utility.JSON(ctx, responseMap)
}

// HandleSuccess xử lý response thành công
func (h *BaseHandler[Model, CreateInput, UpdateInput]) HandleSuccess(ctx *fasthttp.RequestCtx, data interface{}) {
	ctx.SetStatusCode(utility.StatusOK)
	response := ResponseSuccess{
		Status: "success",
		Data:   data,
	}
	responseMap := map[string]interface{}{
		"status": response.Status,
		"data":   response.Data,
	}
	utility.JSON(ctx, responseMap)
}

// GenericHandler xử lý request theo một flow chung
func (h *BaseHandler[Model, CreateInput, UpdateInput]) GenericHandler(ctx *fasthttp.RequestCtx, input interface{}, handler HandlerFunc) {
	response := h.ParseRequestBody(ctx, input)
	if response == nil {
		result, err := handler(ctx, input)
		h.HandleResponse(ctx, result, err)
	} else {
		h.HandleError(ctx, utility.NewError(utility.ErrCodeValidationInput, utility.MsgValidationError, utility.StatusBadRequest, nil))
	}
}

// ====================================
// NHÓM 1: CÁC HÀM CHUẨN MONGODB DRIVER
// ====================================

// 1.1 Thao tác Insert
// -------------------

// InsertOne tạo mới một bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) InsertOne(ctx *fasthttp.RequestCtx) {
	input := new(CreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		createInput := input.(*CreateInput)
		model, err := utility.ConvertStruct(createInput, new(Model))
		if err != nil {
			return nil, err
		}
		return h.Service.InsertOne(context.Background(), model.(Model))
	})
}

// InsertMany tạo nhiều bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) InsertMany(ctx *fasthttp.RequestCtx) {
	input := new([]CreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		createInputs := input.(*[]CreateInput)
		models, err := utility.ConvertStruct(createInputs, new([]Model))
		if err != nil {
			return nil, err
		}
		return h.Service.InsertMany(context.Background(), models.([]Model))
	})
}

// 1.2 Thao tác Find
// ----------------

// FindOne tìm một bản ghi theo điều kiện
func (h *BaseHandler[Model, CreateInput, UpdateInput]) FindOne(ctx *fasthttp.RequestCtx) {
	filter := bson.M{} // Có thể parse filter từ query params
	opts := options.FindOne()
	data, err := h.Service.FindOne(context.Background(), filter, opts)
	h.HandleResponse(ctx, data, err)
}

// Find tìm tất cả bản ghi theo điều kiện
func (h *BaseHandler[Model, CreateInput, UpdateInput]) Find(ctx *fasthttp.RequestCtx) {
	filter := bson.M{} // Có thể parse filter từ query params
	opts := options.Find()
	data, err := h.Service.Find(context.Background(), filter, opts)
	h.HandleResponse(ctx, data, err)
}

// 1.3 Thao tác Update
// ------------------

// UpdateOne cập nhật một bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) UpdateOne(ctx *fasthttp.RequestCtx) {
	input := new(UpdateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		updateInput := input.(*UpdateInput)
		model, err := utility.ConvertStruct(updateInput, new(Model))
		if err != nil {
			return nil, err
		}
		filter := bson.M{} // Có thể parse filter từ query params
		opts := options.Update()
		return h.Service.UpdateOne(context.Background(), filter, model, opts)
	})
}

// UpdateMany cập nhật nhiều bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) UpdateMany(ctx *fasthttp.RequestCtx) {
	input := new(UpdateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		updateInput := input.(*UpdateInput)
		model, err := utility.ConvertStruct(updateInput, new(Model))
		if err != nil {
			return nil, err
		}
		filter := bson.M{} // Có thể parse filter từ query params
		opts := options.Update()
		return h.Service.UpdateMany(context.Background(), filter, model, opts)
	})
}

// 1.4 Thao tác Delete
// ------------------

// DeleteOne xóa một bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) DeleteOne(ctx *fasthttp.RequestCtx) {
	filter := bson.M{} // Có thể parse filter từ query params
	err := h.Service.DeleteOne(context.Background(), filter)
	h.HandleResponse(ctx, nil, err)
}

// DeleteMany xóa nhiều bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) DeleteMany(ctx *fasthttp.RequestCtx) {
	filter := bson.M{} // Có thể parse filter từ query params
	count, err := h.Service.DeleteMany(context.Background(), filter)
	h.HandleResponse(ctx, count, err)
}

// 1.5 Thao tác Atomic
// ------------------

// FindOneAndUpdate tìm và cập nhật một bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) FindOneAndUpdate(ctx *fasthttp.RequestCtx) {
	input := new(UpdateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		updateInput := input.(*UpdateInput)
		model, err := utility.ConvertStruct(updateInput, new(Model))
		if err != nil {
			return nil, err
		}
		filter := bson.M{} // Có thể parse filter từ query params
		opts := options.FindOneAndUpdate()
		return h.Service.FindOneAndUpdate(context.Background(), filter, model, opts)
	})
}

// FindOneAndDelete tìm và xóa một bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) FindOneAndDelete(ctx *fasthttp.RequestCtx) {
	filter := bson.M{} // Có thể parse filter từ query params
	opts := options.FindOneAndDelete()
	data, err := h.Service.FindOneAndDelete(context.Background(), filter, opts)
	h.HandleResponse(ctx, data, err)
}

// 1.6 Các thao tác khác
// --------------------

// CountDocuments đếm số lượng bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) CountDocuments(ctx *fasthttp.RequestCtx) {
	filter := bson.M{} // Có thể parse filter từ query params
	count, err := h.Service.CountDocuments(context.Background(), filter)
	h.HandleResponse(ctx, count, err)
}

// Distinct lấy danh sách các giá trị duy nhất
func (h *BaseHandler[Model, CreateInput, UpdateInput]) Distinct(ctx *fasthttp.RequestCtx) {
	fieldName := string(ctx.FormValue("field"))
	filter := bson.M{} // Có thể parse filter từ query params
	data, err := h.Service.Distinct(context.Background(), fieldName, filter)
	h.HandleResponse(ctx, data, err)
}

// ====================================
// NHÓM 2: CÁC HÀM TIỆN ÍCH MỞ RỘNG
// ====================================

// 2.1 Các hàm Find mở rộng
// -----------------------

// FindOneById tìm một bản ghi theo ID
func (h *BaseHandler[Model, CreateInput, UpdateInput]) FindOneById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	data, err := h.Service.FindOneById(context.Background(), utility.String2ObjectID(id))
	h.HandleResponse(ctx, data, err)
}

// FindManyByIds tìm nhiều bản ghi theo danh sách ID
func (h *BaseHandler[Model, CreateInput, UpdateInput]) FindManyByIds(ctx *fasthttp.RequestCtx) {
	// Parse IDs từ query params
	ids := []primitive.ObjectID{} // Cần implement logic parse IDs
	data, err := h.Service.FindManyByIds(context.Background(), ids)
	h.HandleResponse(ctx, data, err)
}

// FindWithPagination tìm tất cả bản ghi với phân trang
func (h *BaseHandler[Model, CreateInput, UpdateInput]) FindWithPagination(ctx *fasthttp.RequestCtx) {
	page, limit := h.ParsePagination(ctx)
	filter := bson.M{} // Có thể parse filter từ query params
	data, err := h.Service.FindWithPagination(context.Background(), filter, page, limit)
	h.HandleResponse(ctx, data, err)
}

// 2.2 Các hàm Update/Delete mở rộng
// --------------------------------

// UpdateById cập nhật một bản ghi theo ID
func (h *BaseHandler[Model, CreateInput, UpdateInput]) UpdateById(ctx *fasthttp.RequestCtx) {
	input := new(UpdateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		updateInput := input.(*UpdateInput)
		model, err := utility.ConvertStruct(updateInput, new(Model))
		if err != nil {
			return nil, err
		}
		id := h.GetIDFromContext(ctx)
		return h.Service.UpdateById(context.Background(), utility.String2ObjectID(id), model.(Model))
	})
}

// DeleteById xóa một bản ghi theo ID
func (h *BaseHandler[Model, CreateInput, UpdateInput]) DeleteById(ctx *fasthttp.RequestCtx) {
	id := h.GetIDFromContext(ctx)
	err := h.Service.DeleteById(context.Background(), utility.String2ObjectID(id))
	h.HandleResponse(ctx, nil, err)
}

// 2.3 Các hàm Upsert tiện ích
// --------------------------

// Upsert thực hiện thao tác update nếu tồn tại, insert nếu chưa tồn tại
func (h *BaseHandler[Model, CreateInput, UpdateInput]) Upsert(ctx *fasthttp.RequestCtx) {
	input := new(CreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		createInput := input.(*CreateInput)
		model, err := utility.ConvertStruct(createInput, new(Model))
		if err != nil {
			return nil, err
		}
		filter := bson.M{} // Có thể parse filter từ query params
		return h.Service.Upsert(context.Background(), filter, model.(Model))
	})
}

// UpsertMany thực hiện thao tác upsert cho nhiều bản ghi
func (h *BaseHandler[Model, CreateInput, UpdateInput]) UpsertMany(ctx *fasthttp.RequestCtx) {
	input := new([]CreateInput)
	h.GenericHandler(ctx, input, func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error) {
		createInputs := input.(*[]CreateInput)
		models, err := utility.ConvertStruct(createInputs, new([]Model))
		if err != nil {
			return nil, err
		}
		filter := bson.M{} // Có thể parse filter từ query params
		return h.Service.UpsertMany(context.Background(), filter, models.([]Model))
	})
}

// 2.4 Các hàm kiểm tra
// -------------------

// DocumentExists kiểm tra xem một bản ghi có tồn tại không
func (h *BaseHandler[Model, CreateInput, UpdateInput]) DocumentExists(ctx *fasthttp.RequestCtx) {
	filter := bson.M{} // Có thể parse filter từ query params
	exists, err := h.Service.DocumentExists(context.Background(), filter)
	h.HandleResponse(ctx, exists, err)
}
