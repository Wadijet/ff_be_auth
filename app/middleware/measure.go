package middleware

import (
	"atk-go-server/app/utility"

	"github.com/valyala/fasthttp"
)

// Measure là middleware dùng để đo lường thời gian xử lý của một request.
// Nó sẽ loại bỏ thông tin API cũ và thêm thông tin API mới vào stack trước khi gọi handler tiếp theo.
func Measure(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		// Loại bỏ thông tin API cũ từ stack
		utility.RemoveStackApiInfo(15)
		// Thêm thông tin API mới vào stack
		utility.PushStackApiInfo(ctx)

		// Gọi handler tiếp theo
		next(ctx)
	}
}
