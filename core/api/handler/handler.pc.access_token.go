package handler

import (
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/api/services"
)

// AccessTokenHandler xử lý các route liên quan đến Access Token cho Fiber
// Kế thừa từ FiberBaseHandler với các type parameter:
// - Model: models.AccessToken - Model chính của AccessToken
// - CreateInput: models.AccessTokenCreateInput - Struct đầu vào cho việc tạo mới
// - UpdateInput: models.AccessTokenUpdateInput - Struct đầu vào cho việc cập nhật
type AccessTokenHandler struct {
	BaseHandler[models.AccessToken, models.AccessTokenCreateInput, models.AccessTokenUpdateInput]
	AccessTokenService *services.AccessTokenService
}

// NewAccessTokenHandler tạo một instance mới của FiberAccessTokenHandler
// Returns:
//   - *FiberAccessTokenHandler: Instance mới của FiberAccessTokenHandler đã được khởi tạo với các service cần thiết
//   - error: Lỗi nếu có trong quá trình khởi tạo
func NewAccessTokenHandler() (*AccessTokenHandler, error) {
	handler := &AccessTokenHandler{}

	// Khởi tạo AccessTokenService và xử lý error
	service, err := services.NewAccessTokenService()
	if err != nil {
		return nil, fmt.Errorf("failed to create access token service: %v", err)
	}
	handler.AccessTokenService = service
	handler.BaseService = handler.AccessTokenService // Gán service cho BaseHandler

	return handler, nil
}

// Các hàm đặc thù của AccessToken (nếu có) sẽ được thêm vào đây
