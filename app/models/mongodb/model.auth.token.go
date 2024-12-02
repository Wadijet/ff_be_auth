package models

import "github.com/dgrijalva/jwt-go"

// JwtToken ,contains data that will enrypted in JWT token
// When jwt token will decrypt, token model will returns
// Need this model to authenticate and validate resources access by loggedIn user
// JwtToken đại diện cho cấu trúc của một JSON Web Token (JWT).
//
// Các trường:
// - ID: ID của người dùng.
// - Time: Thời gian liên quan đến token.
// - RandomNum: Số ngẫu nhiên để tăng tính bảo mật.
// - StandardClaims: Các yêu cầu tiêu chuẩn của JWT.
type JwtToken struct {
	ID        string `json:"id"`        // User id
	Time      string `json:"time"`      // User id
	RandomNum string `json:"randomNum"` // User id
	jwt.StandardClaims
}

type Token struct {
	Hwid  string `json:"hwid" bson:"hwid,omitempty"`
	Token string `json:"token,omitempty" bson:"token,omitempty"`
}
