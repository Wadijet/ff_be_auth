package handler

import (
	"atk-go-server/app/utility"
	"strconv"

	"github.com/valyala/fasthttp"
)

// BaseHandler là struct cơ sở cho tất cả các handler trong hệ thống
// Cung cấp các phương thức tiện ích chung để xử lý request và response
type BaseHandler struct {
	// Có thể thêm các trường chung ở đây
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

// HandleResponse xử lý response chung cho tất cả các handler
// @param ctx: Context của request
// @param data: Dữ liệu cần trả về
// @param err: Lỗi nếu có
// Hàm này sẽ tự động xử lý:
// - Nếu có lỗi: trả về response với status error và message lỗi
// - Nếu thành công: trả về response với status success và data
func (h *BaseHandler) HandleResponse(ctx *fasthttp.RequestCtx, data interface{}, err error) {
	if err != nil {
		h.HandleError(ctx, err)
	} else {
		h.HandleSuccess(ctx, data)
	}
}

// ParseRequestBody xử lý việc parse request body thành struct
// @param ctx: Context của request
// @param input: Con trỏ đến struct cần parse dữ liệu vào
// @return: Response chứa lỗi nếu có, nil nếu thành công
// Hàm này sẽ:
// 1. Lấy body từ request
// 2. Chuyển đổi JSON thành struct
// 3. Validate dữ liệu của struct
func (h *BaseHandler) ParseRequestBody(ctx *fasthttp.RequestCtx, input interface{}) map[string]interface{} {
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
// @param ctx: Context của request
// @return: (page, limit) - Số trang và số lượng item mỗi trang
// Hàm này sẽ:
// 1. Lấy thông tin limit và page từ query parameters
// 2. Chuyển đổi sang kiểu int64
// 3. Xử lý giá trị mặc định:
//   - limit <= 0: mặc định là 10
//   - page < 0: mặc định là 0
func (h *BaseHandler) ParsePagination(ctx *fasthttp.RequestCtx) (int64, int64) {
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
// @param ctx: Context của request
// @return: ID dưới dạng string
// Hàm này được sử dụng khi cần lấy ID từ URL parameter
// Ví dụ: /api/users/{id} -> lấy giá trị của {id}
func (h *BaseHandler) GetIDFromContext(ctx *fasthttp.RequestCtx) string {
	return ctx.UserValue("id").(string)
}

// HandleError xử lý lỗi chung và trả về response lỗi
// @param ctx: Context của request
// @param err: Lỗi cần xử lý
// Hàm này sẽ:
// 1. Tạo response với status "error"
// 2. Thêm error code và message lỗi vào response
// 3. Gửi response về client
func (h *BaseHandler) HandleError(ctx *fasthttp.RequestCtx, err error) {
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
	// Chuyển đổi response thành map
	responseMap := map[string]interface{}{
		"status":     response.Status,
		"error_code": response.ErrorCode,
		"message":    response.Message,
	}
	utility.JSON(ctx, responseMap)
}

// HandleSuccess xử lý response thành công
// @param ctx: Context của request
// @param data: Dữ liệu cần trả về
// Hàm này sẽ:
// 1. Tạo response với status "success"
// 2. Thêm data vào response
// 3. Gửi response về client
func (h *BaseHandler) HandleSuccess(ctx *fasthttp.RequestCtx, data interface{}) {
	ctx.SetStatusCode(utility.StatusOK)
	response := ResponseSuccess{
		Status: "success",
		Data:   data,
	}
	// Chuyển đổi response thành map
	responseMap := map[string]interface{}{
		"status": response.Status,
		"data":   response.Data,
	}
	utility.JSON(ctx, responseMap)
}

// GenericHandler xử lý request theo một flow chung
// @param ctx: Context của request
// @param input: Con trỏ đến struct cần parse dữ liệu vào
// @param handler: Hàm xử lý logic chính
// Hàm này sẽ:
// 1. Parse và validate input từ request body
// 2. Thực thi logic chính thông qua handler function
// 3. Xử lý response và error
func (h *BaseHandler) GenericHandler(ctx *fasthttp.RequestCtx, input interface{}, handler HandlerFunc) {
	response := h.ParseRequestBody(ctx, input)
	if response == nil {
		// Thực thi logic chính
		result, err := handler(ctx, input)
		h.HandleResponse(ctx, result, err)
	} else {
		h.HandleError(ctx, utility.NewError(utility.ErrCodeValidationInput, utility.MsgValidationError, utility.StatusBadRequest, nil))
	}
}
