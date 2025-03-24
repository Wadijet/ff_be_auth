package middleware

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// RateLimiter để theo dõi số lượng request
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter tạo một rate limiter mới
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit middleware để giới hạn số lượng request
func RateLimit(limit int, window time.Duration) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	limiter := NewRateLimiter(limit, window)
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			ip := ctx.RemoteIP().String()

			limiter.mu.Lock()
			now := time.Now()

			// Xóa các request cũ
			requests := limiter.requests[ip]
			valid := requests[:0]
			for _, t := range requests {
				if now.Sub(t) <= limiter.window {
					valid = append(valid, t)
				}
			}

			// Kiểm tra giới hạn
			if len(valid) >= limiter.limit {
				limiter.mu.Unlock()
				ctx.SetStatusCode(fasthttp.StatusTooManyRequests)
				ctx.SetBody([]byte("Too Many Requests"))
				return
			}

			// Thêm request mới
			valid = append(valid, now)
			limiter.requests[ip] = valid
			limiter.mu.Unlock()

			next(ctx)
		}
	}
}
