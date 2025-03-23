package handler

import (
	"atk-go-server/app/utility"
	"net/http"
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

// HandleResponse xử lý response chung cho tất cả các handler
// @param ctx: Context của request
// @param data: Dữ liệu cần trả về
// @param err: Lỗi nếu có
// Hàm này sẽ tự động xử lý:
// - Nếu có lỗi: trả về response với status error và message lỗi
// - Nếu thành công: trả về response với status success và data
func (h *BaseHandler) HandleResponse(ctx *fasthttp.RequestCtx, data interface{}, err error) {
	var response map[string]interface{}
	if err != nil {
		response = utility.Payload(false, nil, err.Error())
	} else {
		response = utility.FinalResponse(data, nil)
	}
	utility.JSON(ctx, response)
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
	postValues := ctx.PostBody()
	response := utility.Convert2Struct(postValues, input)
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
// 1. Tạo response với status error
// 2. Thêm message lỗi vào response
// 3. Gửi response về client
func (h *BaseHandler) HandleError(ctx *fasthttp.RequestCtx, err error) {
	response := utility.Payload(false, nil, err.Error())
	utility.JSON(ctx, response)
}

// HandleSuccess xử lý response thành công
// @param ctx: Context của request
// @param data: Dữ liệu cần trả về
// Hàm này sẽ:
// 1. Tạo response với status success
// 2. Thêm data vào response
// 3. Gửi response về client
func (h *BaseHandler) HandleSuccess(ctx *fasthttp.RequestCtx, data interface{}) {
	response := utility.FinalResponse(data, nil)
	utility.JSON(ctx, response)
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
		if err == nil {
			ctx.SetStatusCode(http.StatusOK)
			h.HandleSuccess(ctx, result)
		} else {
			ctx.SetStatusCode(http.StatusInternalServerError)
			h.HandleError(ctx, err)
		}
	} else {
		ctx.SetStatusCode(http.StatusBadRequest)
		h.HandleError(ctx, nil)
	}
}
