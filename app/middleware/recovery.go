package middleware

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// Recovery middleware để bắt và xử lý panic
func Recovery(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("Panic recovered: %v", err)
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				ctx.SetBody([]byte("Internal Server Error"))
			}
		}()
		next(ctx)
	}
}
