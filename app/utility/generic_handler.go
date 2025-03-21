package utility

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// HandlerFunc định nghĩa kiểu hàm xử lý logic chính
type HandlerFunc func(ctx *fasthttp.RequestCtx, input interface{}) (interface{}, error)

// GenericHandler xử lý request theo một flow chung
func GenericHandler[T any](ctx *fasthttp.RequestCtx, handler HandlerFunc) {
	var response map[string]interface{} = nil

	// Lấy và validate input
	postValues := ctx.PostBody()
	inputStruct := new(T)

	response = Convert2Struct(postValues, inputStruct)
	if response == nil {
		response = ValidateStruct(inputStruct)
		if response == nil {
			// Thực thi logic chính
			result, err := handler(ctx, inputStruct)
			response = FinalResponse(result, err)
			if err == nil {
				ctx.SetStatusCode(http.StatusOK)
			} else {
				ctx.SetStatusCode(http.StatusInternalServerError)
			}
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
	} else {
		ctx.SetStatusCode(http.StatusBadRequest)
	}

	JSON(ctx, response)
}
