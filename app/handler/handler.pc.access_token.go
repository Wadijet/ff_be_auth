package handler

import (
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/services"
	"meta_commerce/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// AccessTokenHandler là cấu trúc xử lý các yêu cầu liên quan đến Access Token
// Kế thừa từ BaseHandler với các type parameter:
// - Model: models.AccessToken
// - CreateInput: models.AccessTokenCreateInput
// - UpdateInput: models.AccessTokenUpdateInput
type AccessTokenHandler struct {
	BaseHandler[models.AccessToken, models.AccessTokenCreateInput, models.AccessTokenUpdateInput]
	AccessTokenService *services.AccessTokenService
}

// NewAccessTokenHandler khởi tạo một AccessTokenHandler mới
func NewAccessTokenHandler(c *config.Configuration, db *mongo.Client) *AccessTokenHandler {
	handler := &AccessTokenHandler{}
	handler.AccessTokenService = services.NewAccessTokenService(c, db)
	handler.Service = handler.AccessTokenService
	return handler
}

// Các hàm đặc thù của AccessToken (nếu có) sẽ được thêm vào đây
