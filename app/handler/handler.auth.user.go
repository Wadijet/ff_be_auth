package handler

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserHandler là struct chứa các dịch vụ và repository cần thiết để xử lý người dùng
type UserHandler struct {
	RoleService services.RoleService
	UserService services.UserService
}

// NewUserHandler khởi tạo một UserHandler mới
func NewUserHandler(c *config.Configuration, db *mongo.Client) *UserHandler {
	newHandler := new(UserHandler)
	newHandler.UserService = *services.NewUserService(c, db)
	newHandler.RoleService = *services.NewRoleService(c, db)

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
	//opts := new(options.FindOneOptions)
	//opts.SetProjection(bson.D{{"salt", 0}, {"password", 0}})

	response = utility.FinalResponse(h.UserService.FindOneById(ctx, id))

	utility.JSON(ctx, response)
}

// FindAllWithFilter tìm tất cả người dùng với bộ lọc
func (h *UserHandler) FindAllWithFilter(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu phân trang từ request
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

	response = utility.FinalResponse(h.UserService.FindAll(ctx, page, limit))

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

			response = utility.FinalResponse(h.UserService.Create(ctx, inputStruct))

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
		//opts := new(options.FindOneOptions)
		//opts.SetProjection(bson.D{{Key: "salt", Value: 0}, {"password", 0}})
		response = utility.FinalResponse(h.UserService.FindOneById(ctx, strMyID))
	} else {
		response = utility.Payload(true, nil, "An unauthorized access!")
	}

	utility.JSON(ctx, response)
}

// GetMyInfo lấy thông tin của người dùng hiện tại
func (h *UserHandler) GetMyRoles(ctx *fasthttp.RequestCtx) {
	var response map[string]interface{} = nil

	// Lấy dữ liệu
	if ctx.UserValue("userId") != nil {
		strMyID := ctx.UserValue("userId").(string)

		response = utility.FinalResponse(h.UserService.GetRoles(ctx, strMyID))
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

//
