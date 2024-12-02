package services

import (
	"atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService là cấu trúc chứa các phương thức liên quan đến người dùng
type UserService struct {
	crudUser Repository
	crudRole Repository
}

// Khởi tạo UserService với cấu hình và kết nối cơ sở dữ liệu
func NewUserService(c *config.Configuration, db *mongo.Client) *UserService {
	newService := new(UserService)
	newService.crudUser = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.crudRole = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
	return newService
}

// Kiểm tra email có tồn tại hay không
func (h *UserService) IsEmailExist(ctx *fasthttp.RequestCtx, email string) bool {
	filter := bson.M{"email": email}
	result, _ := h.crudUser.FindOne(ctx, filter, nil)
	if result == nil {
		return false
	} else {
		return true
	}
}

// Đăng nhập người dùng
func (h *UserService) Login(ctx *fasthttp.RequestCtx, credential *models.UserLoginInput) (*models.User, error) {
	query := bson.M{"email": credential.Email}
	result, err := h.crudUser.FindOne(ctx, query, nil)
	if result == nil {
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		return nil, err
	}

	if err = user.ComparePassword(credential.Password); err != nil {
		return nil, err
	}

	// Tạo chuỗi tạm thời để tạo token mới
	rdNumber := rand.Intn(100)
	currentTime := time.Now().Unix()

	tokenMap, err := utility.CreateToken(global.MongoDB_ServerConfig.JwtSecret, user.ID.Hex(), strconv.FormatInt(currentTime, 16), strconv.Itoa(rdNumber))
	if err != nil {
		return nil, err
	}

	var idTokenExist int = -1

	for i, _token := range user.Tokens {
		if _token.Hwid == credential.Hwid {
			idTokenExist = i
		}
	}

	if idTokenExist != -1 {
		user.Tokens[idTokenExist].Token = tokenMap["token"]
	} else {
		var newToken models.Token
		newToken.Hwid = credential.Hwid
		newToken.Token = tokenMap["token"]

		user.Tokens = append(user.Tokens, newToken)
	}

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	_, err = h.crudUser.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
	if err != nil {
		return nil, err
	}

	user.Token = tokenMap["token"]

	return &user, nil
}

// Xóa token tại vị trí chỉ định
func RemoveIndex(s []models.Token, index int) []models.Token {
	return append(s[:index], s[index+1:]...)
}

// Đăng xuất người dùng
func (h *UserService) Logout(ctx *fasthttp.RequestCtx, userID string, credential *models.UserLogoutInput) (LogoutResult interface{}, err error) {
	result, err := h.crudUser.FindOneById(ctx, userID, nil)
	if result == nil {
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		return nil, err
	}

	var newTokens []models.Token
	for _, _token := range user.Tokens {
		if _token.Hwid != credential.Hwid {
			newTokens = append(newTokens, _token)
		}
	}
	user.Tokens = newTokens

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	result, err = h.crudUser.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Thay đổi mật khẩu người dùng
func (h *UserService) ChangePassword(ctx *fasthttp.RequestCtx, userID string, credential *models.UserChangePasswordInput) (ChangePasswordResult interface{}, err error) {
	result, err := h.crudUser.FindOneById(ctx, userID, nil)
	if result == nil {
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		return nil, err
	}

	err = user.ComparePassword(credential.OldPassword)
	if err != nil {
		return nil, err
	}

	// Thay đổi mật khẩu
	user.Salt = uuid.New().String()
	passwordBytes := []byte(credential.NewPassword + user.Salt)

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash[:])

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	return h.crudUser.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
}

// Thay đổi thông tin người dùng
func (h *UserService) ChangeInfo(ctx *fasthttp.RequestCtx, userID string, credential *models.UserChangeInfoInput) (ChangeInfoResult interface{}, err error) {
	result, err := h.crudUser.FindOneById(ctx, userID, nil)
	if result == nil {
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		return nil, err
	}

	// Thay đổi thông tin
	user.Name = credential.Name

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	return h.crudUser.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
}

// Kiểm tra token người dùng
func (h *UserService) CheckToken(ctx *fasthttp.RequestCtx, JwtSecret string, tokenString string, requirePermissions []string) (CheckTokenResult interface{}, err error) {
	unauthError := errors.New("An unauthorized access!")
	userBlockedError := errors.New("User has been blocked!")
	notPermissionError := errors.New("You do not have permission to perform the action!")

	// Giải mã token
	t := models.JwtToken{}
	token, err := jwt.ParseWithClaims(tokenString, &t, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, unauthError
	} else {
		findUser, err := h.crudUser.FindOneById(context.TODO(), t.ID, nil)
		if err != nil {
			return nil, unauthError
		}

		var user models.User
		bsonBytes, err := bson.Marshal(findUser)
		if err != nil {
			return nil, unauthError
		}

		err = bson.Unmarshal(bsonBytes, &user)
		if err != nil {
			return nil, unauthError
		}

		if user.IsBlock {
			return nil, userBlockedError
		}

		isRightToken := false
		for _, _token := range user.Tokens {
			if _token.Token == tokenString {
				ctx.SetUserValue("userId", t.ID)                               // đặt ID người dùng đã đăng nhập vào context
				ctx.SetUserValue("roleId", utility.ObjectID2String(user.Role)) // đặt ID vai trò người dùng đã đăng nhập vào context
				isRightToken = true
				break
			}
		}

		if isRightToken == false {
			return nil, unauthError
		}

		if len(requirePermissions) == 0 {
			return user, nil
		}

		strRoleID := utility.ObjectID2String(user.Role)
		findRole, err := h.crudRole.FindOneById(context.TODO(), strRoleID, nil)
		if err != nil {
			return nil, notPermissionError
		}

		var result_findRole models.Role
		bsonBytes, err = bson.Marshal(findRole)
		if err != nil {
			return nil, notPermissionError
		}

		err = bson.Unmarshal(bsonBytes, &result_findRole)
		if err != nil {
			return nil, notPermissionError
		}

		totalCheck := true
		for _, requirePermisson := range requirePermissions {
			checkPermission := false
			for _, s := range result_findRole.Permissions {
				if requirePermisson == s.Name {
					checkPermission = true
					break
				}
			}

			if checkPermission == false {
				totalCheck = false
				break
			}
		}

		if totalCheck == true {
			return user, nil
		} else {
			return nil, notPermissionError
		}
	}
}
