package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"math/rand"
	"strconv"
	"time"

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

// Đăng nhập người dùng bằng email và mật khẩu
func (h *UserService) Login(ctx *fasthttp.RequestCtx, credential *models.UserLoginInput) (*models.User, error) {

	// Tìm người dùng theo email
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

	// So sánh mật khẩu
	if err = user.ComparePassword(credential.Password); err != nil {
		return nil, err
	}

	// Tạo chuỗi random và curentTime để tạo token mới
	rdNumber := rand.Intn(100)
	currentTime := time.Now().Unix()

	// Tạo token mới
	tokenMap, err := utility.CreateToken(global.MongoDB_ServerConfig.JwtSecret, user.ID.Hex(), strconv.FormatInt(currentTime, 16), strconv.Itoa(rdNumber))
	if err != nil {
		return nil, err
	}

	var idTokenExist int = -1

	// duyệt qua tất cả các token, kiểm tra hwid đó đã có token chưa, nếu có thì idTokenExist = i
	for i, _token := range user.Tokens {
		if _token.Hwid == credential.Hwid {
			idTokenExist = i
		}
	}

	// Cập nhật token nếu đã tồn tại
	if idTokenExist != -1 { // nếu idTokenExist != -1 thì cập nhật token mới
		user.Tokens[idTokenExist].JwtToken = tokenMap["token"]
	} else { // Thêm token mới nếu chưa tồn tại hwid
		// Thêm token mới nếu chưa tồn tại
		var newToken models.Token
		newToken.Hwid = credential.Hwid
		newToken.JwtToken = tokenMap["token"]
		newToken.RoleID = ""

		user.Tokens = append(user.Tokens, newToken)
	}

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	// Cập nhật thông tin người dùng trong cơ sở dữ liệu
	_, err = h.crudUser.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Chọn role để làm việc
// Sau khi đăng nhập thành công, token của người dùng sẽ để trống RoleID,
// Khi người dùng chọn RoleID để làm việc, token sẽ được gán RoleID
func (h *UserService) SetWorkingRole(ctx *fasthttp.RequestCtx, userID string, currentToken string, roleID string) (*models.User, error) {

	// Tìm user theo ID
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

	// Tìm token theo currentToken
	var idTokenExist int = -1
	for i, _token := range user.Tokens {
		if _token.JwtToken == currentToken {
			idTokenExist = i
		}
	}

	// Check xem roleID có tồn tại không
	result, err = h.crudRole.FindOneById(ctx, roleID, nil)
	if result == nil {
		return nil, err
	}

	// Cập nhật RoleID
	if idTokenExist != -1 {
		user.Tokens[idTokenExist].RoleID = roleID
	} else {
		return nil, err
	}

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	// Cập nhật thông tin người dùng trong cơ sở dữ liệu
	_, err = h.crudUser.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

// Xóa token tại vị trí chỉ định
func RemoveIndex(s []models.Token, index int) []models.Token {
	return append(s[:index], s[index+1:]...)
}

// Đăng xuất người dùng
func (h *UserService) Logout(ctx *fasthttp.RequestCtx, userID string, credential *models.UserLogoutInput) (LogoutResult interface{}, err error) {

	// Tìm người dùng theo userID
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

	// Tạo mảng mới, không chứa token có hwid trùng với hwid trong credential
	var newTokens []models.Token = []models.Token{}
	for _, _token := range user.Tokens {
		if _token.Hwid != credential.Hwid {
			newTokens = append(newTokens, _token)
		}
	}

	// Cập nhật danh sách token
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

	// Tìm người dùng theo userID
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

	// So sánh mật khẩu cũ
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
	user.Tokens = nil

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	return h.crudUser.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
}

// Thay đổi thông tin người dùng
func (h *UserService) ChangeInfo(ctx *fasthttp.RequestCtx, userID string, credential *models.UserChangeInfoInput) (ChangeInfoResult interface{}, err error) {

	// Tìm người dùng theo userID
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
