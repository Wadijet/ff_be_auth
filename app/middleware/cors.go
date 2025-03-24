package middleware

import (
	"strings"

	"github.com/valyala/fasthttp"
)

// CorsConfig chứa cấu hình cho CORS middleware
type CorsConfig struct {
	AllowHeaders     []string
	AllowMethods     []string
	AllowOrigins     []string
	AllowCredentials bool
}

// DefaultCorsConfig trả về cấu hình CORS mặc định
func DefaultCorsConfig() *CorsConfig {
	return &CorsConfig{
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"HEAD", "GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}
}

// CORS là middleware để xử lý CORS cho các yêu cầu HTTP
func CORS(config *CorsConfig) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			// Xử lý preflight request
			if string(ctx.Method()) == "OPTIONS" {
				ctx.SetStatusCode(fasthttp.StatusOK)
				return
			}

			// Lấy origin từ request
			origin := string(ctx.Request.Header.Peek("Origin"))

			// Kiểm tra origin có được phép không
			allowed := false
			for _, allowedOrigin := range config.AllowOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				// Thiết lập các header CORS
				ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
				ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
				ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ","))
				ctx.Response.Header.Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ","))
			}

			next(ctx)
		}
	}
}
