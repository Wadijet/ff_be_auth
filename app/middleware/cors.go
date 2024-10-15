package middleware

import (
	"github.com/valyala/fasthttp"
)

// Các biến cấu hình cho CORS
var (
	corsAllowHeaders     = "*"                                // Cho phép tất cả các header
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS" // Các phương thức HTTP được phép
	corsAllowOrigin      = "*"                                // Cho phép tất cả các nguồn gốc
	corsAllowCredentials = "true"                             // Cho phép gửi thông tin xác thực
)

// Hàm CORS là một middleware để xử lý CORS cho các yêu cầu HTTP
func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Thiết lập các header CORS cho phản hồi
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		// Gọi handler tiếp theo trong chuỗi middleware
		next(ctx)
	}
}
