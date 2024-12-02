package handler

import (
	"atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"strconv"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler là struct chứa các dịch vụ và repository cần thiết để xử lý người dùng
type UserHandler struct {
	UserCRUD    services.Repository
	RoleCRUD    services.Repository
	UserService services.UserService
}

// NewUserHandler khởi tạo một UserHandler mới
func NewUserHandler(c *config.Configuration, db *mongo.Client) *UserHandler {
	newHandler := new(UserHandler)
	newHandler.UserCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Users)
	newHandler.RoleCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newHandler.UserService = *services.NewUserService(c, db)

	return newHandler
}

// CRUD functions ======================================================

// FindOneById tìm một người dùng theo ID
func (h *UserHandler) FindOneById(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	// GET ID
	id := ctx.UserValue("id").(string)
	// Cài đặt
	opts := new(options.FindOneOptions)
	opts.SetProjection(bson.D{{"salt", 0}, {"password", 0}})

	response = utility.FinalResponse(h.UserCRUD.FindOneById(ctx, id, opts))

	utility.JSON(ctx, response)
}

// FindAllWithFilter tìm tất cả người dùng với bộ lọc
func (h *UserHandler) FindAllWithFilter(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	postValues := ctx.PostBody()
	inputStruct := new(models.UserFilterInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm xử lý logic

			// Lấy dữ liệu
			buf := string(ctx.FormValue("limit"))
			limit, err := strconv.ParseInt(buf, 10, 64)
			if err != nil {
				limit = 10
			}

			buf = string(ctx.FormValue("page"))
			page, err := strconv.ParseInt(buf, 10, 64)
			if err != nil {
				page = 0
			}

			// Cài đặt
			opts := new(options.FindOptions)
			opts.SetLimit(limit)
			opts.SetSkip(page * limit)
			opts.SetSort(bson.D{{"updatedAt", 1}})
			opts.SetProjection(bson.D{{"salt", 0}, {"password", 0}})

			filterMap := make(map[string]interface{})
			if len(inputStruct.Emails) > 0 {
				filterMap["email"] = bson.M{"$in": inputStruct.Emails}
			}
			if len(inputStruct.RoleIDs) > 0 {
				filterMap["role"] = bson.M{"$in": inputStruct.RoleIDs}
			}

			var filter bson.M
			data, err := bson.Marshal(filterMap)
			if err != nil {
				return
			}

			err = bson.Unmarshal(data, &filter)
			if err != nil {
				return
			}

			response = utility.FinalResponse(h.UserCRUD.FindAllWithPaginate(ctx, filter, opts))
		}
	}

	utility.JSON(ctx, response)
}

// OTHER functions =======================================================

// Registry đăng ký người dùng mới
func (h *UserHandler) Registry(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	postValues := ctx.PostBody()
	inputStruct := new(models.UserCreateInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm tạo json changes

			if h.UserService.IsEmailExist(ctx, inputStruct.Email) == true {
				response = utility.Payload(false, nil, "User already exists!")
			} else {

				// Tìm Role User
				userRole, err := h.RoleCRUD.FindOne(ctx, bson.M{"name": "User"}, nil)
				if userRole != nil {
					userRoleMap, err := utility.ToMap(userRole)
					if err != nil {
						response = utility.Payload(false, err, "Can not create user!")
					} else {
						userRoleID := userRoleMap["_id"].(primitive.ObjectID)

						newUser := new(models.User)
						newUser.Name = inputStruct.Name
						newUser.Email = inputStruct.Email
						newUser.Role = userRoleID

						newUser.Salt = uuid.New().String()
						passwordBytes := []byte(inputStruct.Password + newUser.Salt)

						hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
						if err != nil {
							response = utility.Payload(false, err.Error(), "Can not create hash password!")
						} else {
							newUser.Password = string(hash[:])
							response = utility.FinalResponse(h.UserCRUD.InsertOne(ctx, newUser))
						}
					}
				} else {
					response = utility.Payload(false, err, "Can not create user!")
				}

			}
		}
	}
	utility.JSON(ctx, response)
}

// Login đăng nhập người dùng
func (h *UserHandler) Login(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	postValues := ctx.PostBody()
	inputStruct := new(models.UserLoginInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { // Gọi hàm tạo json changes
			user, err := h.UserService.Login(ctx, inputStruct)
			if user == nil {
				response = utility.Payload(false, err, "Login information is incorrect!")
			} else {

				response = utility.Payload(true, user, "Logged in successfully.")
			}
		}
	}
	utility.JSON(ctx, response)
}

// Logout đăng xuất người dùng
func (h *UserHandler) Logout(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	if ctx.UserValue("userId") != nil {
		strMyID := ctx.UserValue("userId").(string)

		// Lấy dữ liệu
		postValues := ctx.PostBody()
		inputStruct := new(models.UserLogoutInput)
		response = utility.Convert2Struct(postValues, inputStruct)
		if response == nil { // Kiểm tra dữ liệu đầu vào
			response = utility.ValidateStruct(inputStruct)
			if response == nil { // Gọi hàm tạo json changes
				response = utility.FinalResponse(h.UserService.Logout(ctx, strMyID, inputStruct))
			}
		}
	} else {
		response = utility.Payload(true, nil, "An unauthorized access!")
	}

	utility.JSON(ctx, response)
}

// GetMyInfo lấy thông tin của người dùng hiện tại
func (h *UserHandler) GetMyInfo(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	if ctx.UserValue("userId") != nil {
		strMyID := ctx.UserValue("userId").(string)
		// Cài đặt
		opts := new(options.FindOneOptions)
		opts.SetProjection(bson.D{{"salt", 0}, {"password", 0}})
		response = utility.FinalResponse(h.UserCRUD.FindOneById(ctx, strMyID, opts))
	} else {
		response = utility.Payload(true, nil, "An unauthorized access!")
	}

	utility.JSON(ctx, response)
}

// ChangePassword thay đổi mật khẩu người dùng
func (h *UserHandler) ChangePassword(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	postValues := ctx.PostBody()
	inputStruct := new(models.UserChangePasswordInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { //
			if ctx.UserValue("userId") != nil {
				strMyID := ctx.UserValue("userId").(string)
				response = utility.FinalResponse(h.UserService.ChangePassword(ctx, strMyID, inputStruct))
			} else {
				response = utility.Payload(true, nil, "An unauthorized access!")
			}
		}
	}

	utility.JSON(ctx, response)
}

// ChangeInfo thay đổi thông tin người dùng
func (h *UserHandler) ChangeInfo(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	postValues := ctx.PostBody()
	inputStruct := new(models.UserChangeInfoInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { //
			if ctx.UserValue("userId") != nil {
				strMyID := ctx.UserValue("userId").(string)
				response = utility.FinalResponse(h.UserService.ChangeInfo(ctx, strMyID, inputStruct))
			} else {
				response = utility.Payload(true, nil, "An unauthorized access!")
			}

		}
	}

	utility.JSON(ctx, response)
}

// CheckToken kiểm tra token của người dùng
func (h *UserHandler) CheckToken(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	postValues := ctx.PostBody()
	inputStruct := new(models.UserCheckTokenInput)
	response = utility.Convert2Struct(postValues, inputStruct)
	if response == nil { // Kiểm tra dữ liệu đầu vào
		response = utility.ValidateStruct(inputStruct)
		if response == nil { //
			response = utility.FinalResponse(h.UserService.CheckToken(ctx, global.MongoDB_ServerConfig.JwtSecret, inputStruct.Token, inputStruct.Permissions))
		}
	}

	utility.JSON(ctx, response)
}
