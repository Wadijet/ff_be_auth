package middleware

import (
	"meta_commerce/app/utility"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
)

// MeasureConfig chứa cấu hình cho Measure middleware
type MeasureConfig struct {
	WarningThreshold time.Duration // Ngưỡng cảnh báo cho thời gian xử lý
	LogSlowRequests  bool          // Có log các request chậm không
}

// DefaultMeasureConfig trả về cấu hình mặc định cho Measure middleware
func DefaultMeasureConfig() *MeasureConfig {
	return &MeasureConfig{
		WarningThreshold: 500 * time.Millisecond,
		LogSlowRequests:  true,
	}
}

// Measure là middleware dùng để đo lường thời gian xử lý và tài nguyên sử dụng của một request
func Measure(config *MeasureConfig) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			start := time.Now()

			// Lấy thông tin memory trước khi xử lý
			var memBefore runtime.MemStats
			runtime.ReadMemStats(&memBefore)

			// Loại bỏ thông tin API cũ từ stack
			utility.RemoveStackApiInfo(15)
			// Thêm thông tin API mới vào stack
			utility.PushStackApiInfo(ctx)

			// Gọi handler tiếp theo
			next(ctx)

			// Tính toán thời gian xử lý
			duration := time.Since(start)

			// Lấy thông tin memory sau khi xử lý
			var memAfter runtime.MemStats
			runtime.ReadMemStats(&memAfter)

			// Tính toán memory sử dụng
			memUsed := memAfter.Alloc - memBefore.Alloc

			// Log thông tin nếu request xử lý chậm
			if config.LogSlowRequests && duration > config.WarningThreshold {
				utility.LogWarning("Slow request detected",
					"path", string(ctx.Path()),
					"method", string(ctx.Method()),
					"duration", duration,
					"memory_used", memUsed,
				)
			}

			// Thêm thông tin metrics vào response header
			ctx.Response.Header.Set("X-Request-Duration", duration.String())
			ctx.Response.Header.Set("X-Memory-Used", utility.FormatBytes(memUsed))
		}
	}
}
