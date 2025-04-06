package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
)

// FiberAccessTokenHandler xử lý các route liên quan đến Access Token cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.AccessToken - Model chính của AccessToken
// - CreateInput: models.AccessTokenCreateInput - Struct đầu vào cho việc tạo mới
// - UpdateInput: models.AccessTokenUpdateInput - Struct đầu vào cho việc cập nhật
type FiberAccessTokenHandler struct {
	FiberBaseHandler[models.AccessToken, models.AccessTokenCreateInput, models.AccessTokenUpdateInput]
	AccessTokenService *services.AccessTokenService
}

// NewFiberAccessTokenHandler tạo một instance mới của FiberAccessTokenHandler
// Returns:
//   - *FiberAccessTokenHandler: Instance mới của FiberAccessTokenHandler đã được khởi tạo với các service cần thiết
func NewFiberAccessTokenHandler() *FiberAccessTokenHandler {
	handler := &FiberAccessTokenHandler{}
	handler.AccessTokenService = services.NewAccessTokenService()
	handler.Service = handler.AccessTokenService // Gán service cho BaseHandler
	return handler
}

// Các hàm đặc thù của AccessToken (nếu có) sẽ được thêm vào đây
