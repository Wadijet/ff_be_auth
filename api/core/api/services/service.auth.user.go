package services

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"meta_commerce/core/api/dto"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
	"meta_commerce/core/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService là cấu trúc chứa các phương thức liên quan đến người dùng
type UserService struct {
	*BaseServiceMongoImpl[models.User]
	userRoleService *BaseServiceMongoImpl[models.UserRole]
}

// NewUserService tạo mới UserService
func NewUserService() (*UserService, error) {
	// Lấy collections từ registry mới
	userCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.Users)
	if !exist {
		return nil, fmt.Errorf("failed to get users collection: %v", common.ErrNotFound)
	}

	userRoleCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.UserRoles)
	if !exist {
		return nil, fmt.Errorf("failed to get user_roles collection: %v", common.ErrNotFound)
	}

	return &UserService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.User](userCollection),
		userRoleService:      NewBaseServiceMongo[models.UserRole](userRoleCollection),
	}, nil
}

// Logout đăng xuất người dùng
func (s *UserService) Logout(ctx context.Context, userID primitive.ObjectID, input *dto.UserLogoutInput) error {
	// Tìm user
	user, err := s.BaseServiceMongoImpl.FindOneById(ctx, userID)
	if err != nil {
		return err
	}

	// Xóa token của hwid
	newTokens := make([]models.Token, 0)
	for _, t := range user.Tokens {
		if t.Hwid != input.Hwid {
			newTokens = append(newTokens, t)
		}
	}
	user.Tokens = newTokens
	user.Token = "" // Xóa token hiện tại
	user.UpdatedAt = time.Now().Unix()

	// Cập nhật user
	_, err = s.BaseServiceMongoImpl.UpdateById(ctx, userID, user)
	return err
}

// LoginWithFirebase đăng nhập bằng Firebase ID token
func (s *UserService) LoginWithFirebase(ctx context.Context, input *dto.FirebaseLoginInput) (*models.User, error) {
	// 1. Verify Firebase ID token
	token, err := utility.VerifyIDToken(ctx, input.IDToken)
	if err != nil {
		return nil, common.NewError(
			common.ErrCodeAuthCredentials,
			"Token không hợp lệ",
			common.StatusUnauthorized,
			err,
		)
	}

	// 2. Lấy thông tin user từ Firebase
	firebaseUser, err := utility.GetUserByUID(ctx, token.UID)
	if err != nil {
		return nil, err
	}

	// 3. Tìm user trong MongoDB theo Firebase UID
	filter := bson.M{"firebaseUid": token.UID}
	user, err := s.BaseServiceMongoImpl.FindOne(ctx, filter, nil)

	// 4. Nếu không tìm thấy, tạo user mới
	if err == common.ErrNotFound {
		newUser := &models.User{
			FirebaseUID:    token.UID,
			Email:          firebaseUser.Email,
			EmailVerified:  firebaseUser.EmailVerified,
			Phone:          firebaseUser.PhoneNumber,
			PhoneVerified:  firebaseUser.PhoneNumber != "",
			Name:           firebaseUser.DisplayName,
			AvatarURL:      firebaseUser.PhotoURL,
			IsBlock:        false,
			Tokens:         []models.Token{},
			CreatedAt:      time.Now().Unix(),
			UpdatedAt:      time.Now().Unix(),
		}

		user, err = s.BaseServiceMongoImpl.InsertOne(ctx, *newUser)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// 5. Nếu tìm thấy, sync thông tin từ Firebase (nếu có thay đổi)
		updated := false

		if user.Email != firebaseUser.Email {
			user.Email = firebaseUser.Email
			updated = true
		}

		if user.EmailVerified != firebaseUser.EmailVerified {
			user.EmailVerified = firebaseUser.EmailVerified
			updated = true
		}

		if user.Phone != firebaseUser.PhoneNumber {
			user.Phone = firebaseUser.PhoneNumber
			user.PhoneVerified = firebaseUser.PhoneNumber != ""
			updated = true
		}

		if user.Name != firebaseUser.DisplayName && firebaseUser.DisplayName != "" {
			user.Name = firebaseUser.DisplayName
			updated = true
		}

		if user.AvatarURL != firebaseUser.PhotoURL && firebaseUser.PhotoURL != "" {
			user.AvatarURL = firebaseUser.PhotoURL
			updated = true
		}

		if updated {
			user.UpdatedAt = time.Now().Unix()
			user, err = s.BaseServiceMongoImpl.UpdateById(ctx, user.ID, user)
			if err != nil {
				return nil, err
			}
		}
	}

	// 6. Kiểm tra user bị block
	if user.IsBlock {
		return nil, common.NewError(
			common.ErrCodeAuth,
			"Tài khoản đã bị khóa",
			common.StatusForbidden,
			nil,
		)
	}

	// 7. Tạo JWT token cho backend
	rdNumber := rand.Intn(100)
	currentTime := time.Now().Unix()

	tokenMap, err := utility.CreateToken(
		global.MongoDB_ServerConfig.JwtSecret,
		user.ID.Hex(),
		strconv.FormatInt(currentTime, 16),
		strconv.Itoa(rdNumber),
	)
	if err != nil {
		return nil, err
	}

	// 8. Cập nhật token vào user
	user.Token = tokenMap["token"]

	// Cập nhật hoặc thêm token vào tokens array (theo hwid)
	var idTokenExist int = -1
	for i, _token := range user.Tokens {
		if _token.Hwid == input.Hwid {
			idTokenExist = i
			break
		}
	}

	if idTokenExist == -1 {
		user.Tokens = append(user.Tokens, models.Token{
			Hwid:     input.Hwid,
			JwtToken: tokenMap["token"],
		})
	} else {
		user.Tokens[idTokenExist].JwtToken = tokenMap["token"]
	}

	// 9. Lưu user
	updatedUser, err := s.BaseServiceMongoImpl.UpdateById(ctx, user.ID, user)
	if err != nil {
		return nil, err
	}

	// 10. Nếu chưa có admin nào, tự động set user đầu tiên làm admin
	// Đây là phương án phổ biến: "First user becomes admin"
	initService, err := NewInitService()
	if err == nil {
		hasAdmin, err := initService.HasAnyAdministrator()
		if err == nil && !hasAdmin {
			// Chưa có admin, tự động set user này làm admin
			_, err = initService.SetAdministrator(updatedUser.ID)
			if err != nil && err != common.ErrUserAlreadyAdmin {
				// Log warning nhưng không fail login
				// User vẫn có thể login, chỉ là chưa có quyền admin
				// Có thể set admin sau bằng cách khác
			}
		}
	}

	return &updatedUser, nil
}