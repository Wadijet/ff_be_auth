package services

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService là cấu trúc chứa các phương thức liên quan đến người dùng
type UserService struct {
	*BaseServiceImpl[models.User]
	userRoleService *BaseServiceImpl[models.UserRole]
}

// NewUserService tạo mới UserService
func NewUserService(c *config.Configuration, db *mongo.Client) *UserService {
	userCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Users)).Collection(global.MongoDB_ColNames.Users)
	userRoleCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.UserRoles)).Collection(global.MongoDB_ColNames.UserRoles)

	return &UserService{
		BaseServiceImpl: NewBaseService[models.User](userCollection),
		userRoleService: NewBaseService[models.UserRole](userRoleCollection),
	}
}

// IsEmailExist kiểm tra email có tồn tại hay không
func (s *UserService) IsEmailExist(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	var user models.User
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create tạo mới một người dùng
func (s *UserService) Create(ctx context.Context, input *models.UserCreateInput) (*models.User, error) {
	// Kiểm tra email tồn tại
	exists, err := s.IsEmailExist(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utility.ErrDuplicate
	}

	// Validate email
	if err := utility.ValidateEmail(input.Email); err != nil {
		return nil, err
	}

	// Validate password
	if err := utility.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	// Hash password
	salt := uuid.New().String()
	passwordBytes := []byte(input.Password + salt)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Tạo user mới
	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Salt:     salt,
		IsBlock:  false,
	}

	// Lưu user
	createdUser, err := s.BaseServiceImpl.Create(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

// Login đăng nhập người dùng
func (s *UserService) Login(ctx context.Context, input *models.UserLoginInput) (*models.User, error) {
	// Tìm user theo email
	filter := bson.M{"email": input.Email}
	user, err := s.BaseServiceImpl.FindOneByFilter(ctx, filter, nil)
	if err != nil {
		if err == utility.ErrNotFound {
			return nil, utility.ErrInvalidCredentials
		}
		return nil, err
	}

	// Kiểm tra mật khẩu
	if err := user.ComparePassword(input.Password); err != nil {
		return nil, utility.ErrInvalidCredentials
	}

	// Tạo chuỗi random và curentTime để tạo token mới
	rdNumber := rand.Intn(100)
	currentTime := time.Now().Unix()

	tokenMap, err := utility.CreateToken(global.MongoDB_ServerConfig.JwtSecret, user.ID.Hex(), strconv.FormatInt(currentTime, 16), strconv.Itoa(rdNumber))
	if err != nil {
		return nil, err
	}

	// Cập nhật token mới
	user.Token = tokenMap["token"]

	var idTokenExist int = -1
	// duyệt qua tất cả các token, kiểm tra hwid đó đã có token chưa, nếu có thì idTokenExist = i
	for i, _token := range user.Tokens {
		if _token.Hwid == input.Hwid {
			idTokenExist = i
			break
		}
	}

	// Nếu không có token, thêm token mới
	if idTokenExist == -1 {
		user.Tokens = append(user.Tokens, models.Token{
			Hwid:     input.Hwid,
			JwtToken: tokenMap["token"],
		})
	} else {
		// Nếu có token, cập nhật token mới cho hwid đó
		user.Tokens[idTokenExist].JwtToken = tokenMap["token"]
	}

	// Cập nhật user
	updatedUser, err := s.BaseServiceImpl.Update(ctx, user.ID, user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// Logout đăng xuất người dùng
func (s *UserService) Logout(ctx context.Context, userID primitive.ObjectID, input *models.UserLogoutInput) error {
	// Tìm user
	user, err := s.BaseServiceImpl.FindOne(ctx, userID)
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
	_, err = s.BaseServiceImpl.Update(ctx, userID, user)
	return err
}

// ChangePassword thay đổi mật khẩu
func (s *UserService) ChangePassword(ctx context.Context, userID primitive.ObjectID, input *models.UserChangePasswordInput) error {
	// Tìm user
	user, err := s.BaseServiceImpl.FindOne(ctx, userID)
	if err != nil {
		return err
	}

	// Kiểm tra mật khẩu cũ
	if err := user.ComparePassword(input.OldPassword); err != nil {
		return errors.New("Invalid old password")
	}

	// Validate mật khẩu mới
	if err := utility.ValidatePassword(input.NewPassword); err != nil {
		return err
	}

	// Tạo salt mới và hash mật khẩu mới
	salt := uuid.New().String()
	passwordBytes := []byte(input.NewPassword + salt)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Cập nhật thông tin
	user.Password = string(hashedPassword)
	user.Salt = salt
	user.Tokens = nil // Xóa tất cả token
	user.UpdatedAt = time.Now().Unix()

	// Cập nhật user
	_, err = s.BaseServiceImpl.Update(ctx, userID, user)
	return err
}
