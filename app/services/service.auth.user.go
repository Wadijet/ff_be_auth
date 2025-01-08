package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// UserService là cấu trúc chứa các phương thức liên quan đến người dùng
type UserService struct {
	crudUser     RepositoryService
	crudUserRole RepositoryService
}

// Khởi tạo UserService với cấu hình và kết nối cơ sở dữ liệu
func NewUserService(c *config.Configuration, db *mongo.Client) *UserService {
	newService := new(UserService)
	newService.crudUser = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.crudUserRole = *NewRepository(c, db, global.MongoDB_ColNames.UserRoles)
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

// Tạo mới một người dùng
func (h *UserService) Create(ctx *fasthttp.RequestCtx, credential *models.UserCreateInput) (CreateResult interface{}, err error) {
	// Kiểm tra email của người dùng đã tồn tại chưa bằng cách gọi hàm IsEmailExist
	if h.IsEmailExist(ctx, credential.Email) {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		return nil, errors.New("Email already exists")
	}

	// Tạo mới người dùng
	newUser := new(models.User)
	newUser.Name = credential.Name
	newUser.Email = credential.Email

	newUser.Salt = uuid.New().String()
	passwordBytes := []byte(credential.Password + newUser.Salt)

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}
	newUser.Password = string(hash[:])

	// Thêm người dùng vào cơ sở dữ liệu
	result, err := h.crudUser.InsertOne(ctx, newUser)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	return result, nil
}

// Tìm một người dùng theo ID
func (h *UserService) FindOneById(ctx *fasthttp.RequestCtx, userID string) (result *models.User, err error) {
	findOneResult, err := h.crudUser.FindOneById(ctx, utility.String2ObjectID(userID), nil)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(findOneResult)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return &user, nil
}

// Tìm tất cả các người dùng với phân trang
func (h *UserService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (results interface{}, err error) {
	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})

	findAllResult, err := h.crudUser.FindAllWithPaginate(ctx, bson.D{}, opts)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	var bsonMUsers []models.User
	for _, bsonM := range findAllResult.Items {
		var user models.User
		bsonBytes, err := bson.Marshal(bsonM)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return nil, err
		}

		err = bson.Unmarshal(bsonBytes, &user)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return nil, err
		}

		bsonMUsers = append(bsonMUsers, user)
	}

	var findAllResult2 models.UserPaginateResult
	findAllResult2.Page = findAllResult.Page
	findAllResult2.Limit = findAllResult.Limit
	findAllResult2.Items = bsonMUsers

	ctx.SetStatusCode(fasthttp.StatusOK)
	return findAllResult2, nil
}

// Đăng nhập người dùng bằng email và mật khẩu
func (h *UserService) Login(ctx *fasthttp.RequestCtx, credential *models.UserLoginInput) (*models.User, error) {

	// Tìm người dùng theo email
	query := bson.M{"email": credential.Email}
	result, err := h.crudUser.FindOne(ctx, query, nil)
	if result == nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return nil, errors.New("Invalid email or password")
	}

	var user models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	// So sánh mật khẩu
	if err = user.ComparePassword(credential.Password); err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return nil, errors.New("Invalid email or password")
	}

	// Tạo chuỗi random và curentTime để tạo token mới
	rdNumber := rand.Intn(100)
	currentTime := time.Now().Unix()

	// Tạo token mới
	tokenMap, err := utility.CreateToken(global.MongoDB_ServerConfig.JwtSecret, user.ID.Hex(), strconv.FormatInt(currentTime, 16), strconv.Itoa(rdNumber))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
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
	user.Token = tokenMap["token"]

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	// Cập nhật thông tin người dùng trong cơ sở dữ liệu
	_, err = h.crudUser.UpdateOneById(ctx, user.ID, change)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return &user, nil
}

// Đăng xuất người dùng
func (h *UserService) Logout(ctx *fasthttp.RequestCtx, userID string, credential *models.UserLogoutInput) (LogoutResult interface{}, err error) {

	// Tìm người dùng theo userID
	result, err := h.crudUser.FindOneById(ctx, utility.String2ObjectID(userID), nil)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
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
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	updateResult, err := h.crudUser.UpdateOneById(ctx, user.ID, change)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return updateResult, nil
}

// Thay đổi mật khẩu người dùng
func (h *UserService) ChangePassword(ctx *fasthttp.RequestCtx, userID string, credential *models.UserChangePasswordInput) (ChangePasswordResult interface{}, err error) {

	// Tìm người dùng theo userID
	result, err := h.crudUser.FindOneById(ctx, utility.String2ObjectID(userID), nil)
	if result == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	// So sánh mật khẩu cũ
	err = user.ComparePassword(credential.OldPassword)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return nil, errors.New("Invalid old password")
	}

	// Thay đổi mật khẩu
	user.Salt = uuid.New().String()
	passwordBytes := []byte(credential.NewPassword + user.Salt)

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}
	user.Password = string(hash[:])
	user.Tokens = nil

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	updateResult, err := h.crudUser.UpdateOneById(ctx, user.ID, change)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return updateResult, nil
}

// Thay đổi thông tin người dùng
func (h *UserService) ChangeInfo(ctx *fasthttp.RequestCtx, userID string, credential *models.UserChangeInfoInput) (ChangeInfoResult interface{}, err error) {

	// Tìm người dùng theo userID
	result, err := h.crudUser.FindOneById(ctx, utility.String2ObjectID(userID), nil)
	if result == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(result)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	// Thay đổi thông tin
	user.Name = credential.Name

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	updateResult, err := h.crudUser.UpdateOneById(ctx, user.ID, change)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return updateResult, nil
}

// Lấy tất cả các Role của người dùng
func (h *UserService) GetRoles(ctx *fasthttp.RequestCtx, strUserID string) (GetRolesResult interface{}, err error) {

	userID := utility.String2ObjectID(strUserID)

	// Tìm tất cả các Role của người dùng
	// Cài đặt bộ lọc tìm kiếm
	filter := bson.D{{Key: "userId", Value: userID}}

	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetSort(bson.D{{Key: "updatedAt", Value: 1}})

	// Tìm tất cả các Role của người dùng
	findAllResult, err := h.crudUserRole.FindAll(ctx, filter, opts)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil, err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return findAllResult, nil
}
