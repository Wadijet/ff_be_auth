package middleware

import (
	"context"
	"time"

	"github.com/valyala/fasthttp"
)

// Timeout middleware để giới hạn thời gian xử lý request
func Timeout(timeout time.Duration) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			// Tạo context với timeout
			ctxWithTimeout, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// Tạo channel để theo dõi việc hoàn thành request
			done := make(chan struct{})
			go func() {
				next(ctx)
				close(done)
			}()

			// Đợi request hoàn thành hoặc timeout
			select {
			case <-done:
				// Request hoàn thành bình thường
			case <-ctxWithTimeout.Done():
				// Request timeout
				ctx.SetStatusCode(fasthttp.StatusRequestTimeout)
				ctx.SetBody([]byte("Request Timeout"))
			}
		}
	}
}
