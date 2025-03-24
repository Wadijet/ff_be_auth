package tests

import (
	"atk-go-server/app/handler"
	"atk-go-server/app/middleware"
	"atk-go-server/config"
	"atk-go-server/database"
	"atk-go-server/global"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	testClient *mongo.Client
	testConfig *config.Configuration
)

func init() {
	// Khởi tạo config test
	testConfig = &config.Configuration{
		MongoDB_ConnectionURL: "mongodb://localhost:27017",
		MongoDB_DBNameAuth:    "test_db",
		MongoDb_Uri_Auth:      "mongodb://localhost:27017",
		MongoDb_Uri_Data:      "mongodb://localhost:27017",
		JwtSecret:             "test_secret",
		Metadata_Path:         "./metadata",
	}

	// Khởi tạo kết nối MongoDB
	var err error
	testClient, err = database.GetInstance(testConfig)
	if err != nil {
		panic(err)
	}

	// Khởi tạo validator
	global.Validate = validator.New()

	// Khởi tạo logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Khởi tạo config trong global
	global.MongoDB_ServerConfig = testConfig

	// Khởi tạo tên collection
	global.MongoDB_ColNames.Users = "users"
	global.MongoDB_ColNames.Permissions = "permissions"
	global.MongoDB_ColNames.Roles = "roles"
	global.MongoDB_ColNames.RolePermissions = "role_permissions"
	global.MongoDB_ColNames.UserRoles = "user_roles"
	global.MongoDB_ColNames.Agents = "agents"
	global.MongoDB_ColNames.AccessTokens = "access_tokens"
	global.MongoDB_ColNames.FbPages = "fb_pages"
	global.MongoDB_ColNames.FbConvesations = "fb_conversations"
	global.MongoDB_ColNames.FbMessages = "fb_messages"
	global.MongoDB_ColNames.FbPosts = "fb_posts"
	global.MongoDB_ColNames.PcOrders = "pc_orders"

	// Xóa dữ liệu test cũ
	testClient.Database(testConfig.MongoDB_DBNameAuth).Collection(global.MongoDB_ColNames.Users).DeleteMany(context.Background(), bson.M{})

	// Tạo user test
	userHandler := handler.NewUserHandler(testConfig, testClient)
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody([]byte(`{
		"name": "Test User",
		"email": "test@example.com",
		"password": "Test@123"
	}`))
	userHandler.Registry(ctx)

	// Kiểm tra kết quả tạo user
	var response TestResponse
	err = json.Unmarshal(ctx.Response.Body(), &response)
	if err != nil {
		panic(err)
	}
	if response.Status != "success" {
		panic("Failed to create test user: " + response.Message)
	}
	logrus.Info("Created test user successfully")
}

// TestResponse cấu trúc response chuẩn
type TestResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// setupTestRouter khởi tạo router cho test
func setupTestRouter() *router.Router {
	// Khởi tạo router
	r := router.New()

	// Khởi tạo middleware
	middle := middleware.NewJwtToken(testConfig, testClient)

	// Khởi tạo handlers
	staticHandler := handler.NewStaticHandler()
	userHandler := handler.NewUserHandler(testConfig, testClient)

	// Đăng ký routes
	r.GET("/api/v1/static/test", staticHandler.TestApi)
	r.GET("/api/v1/static/system", middle.CheckUserAuth("", staticHandler.GetSystemStatic))
	r.POST("/api/v1/users/login", userHandler.Login)

	return r
}

// TestStaticAPI test các API static
func TestStaticAPI(t *testing.T) {
	r := setupTestRouter()

	// Test case cho API /static/test
	t.Run("Test Static Test API", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/api/v1/static/test")
		ctx.Request.Header.SetMethod("GET")

		r.Handler(ctx)

		// Kiểm tra status code
		assert.Equal(t, 200, ctx.Response.StatusCode())

		// Parse response
		var response TestResponse
		err := json.Unmarshal(ctx.Response.Body(), &response)
		assert.NoError(t, err)

		// Kiểm tra response
		assert.Equal(t, "success", response.Status)
		assert.NotEmpty(t, response.Message)
	})

	// Test case cho API /static/system khi chưa đăng nhập
	t.Run("Test Static System API Without Auth", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/api/v1/static/system")
		ctx.Request.Header.SetMethod("GET")

		r.Handler(ctx)

		// Kiểm tra status code
		assert.Equal(t, 401, ctx.Response.StatusCode())

		// Parse response
		var response TestResponse
		err := json.Unmarshal(ctx.Response.Body(), &response)
		assert.NoError(t, err)

		// Kiểm tra response
		assert.Equal(t, "error", response.Status)
		assert.NotEmpty(t, response.Message)
	})
}

// TestAuthAPI test các API xác thực
func TestAuthAPI(t *testing.T) {
	r := setupTestRouter()

	// Test case cho API login với dữ liệu không hợp lệ
	t.Run("Test Login API With Invalid Data", func(t *testing.T) {
		// Tạo request body
		loginBody := map[string]string{
			"email":    "invalid_email",
			"password": "password123",
			"hwid":     "test_hwid",
		}
		body, _ := json.Marshal(loginBody)

		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/api/v1/users/login")
		ctx.Request.Header.SetMethod("POST")
		ctx.Request.Header.SetContentType("application/json")
		ctx.Request.SetBody(body)

		r.Handler(ctx)

		// Kiểm tra status code
		assert.Equal(t, 400, ctx.Response.StatusCode())

		// Parse response
		var response TestResponse
		err := json.Unmarshal(ctx.Response.Body(), &response)
		assert.NoError(t, err)

		// Kiểm tra response
		assert.Equal(t, "error", response.Status)
		assert.NotEmpty(t, response.Message)
	})

	// Test case cho API login với dữ liệu hợp lệ
	t.Run("Test Login API With Valid Data", func(t *testing.T) {
		// Tạo request body
		loginBody := map[string]string{
			"email":    "test@example.com",
			"password": "Test@123",
			"hwid":     "test_hwid",
		}
		body, _ := json.Marshal(loginBody)

		ctx := &fasthttp.RequestCtx{}
		ctx.Request.SetRequestURI("/api/v1/users/login")
		ctx.Request.Header.SetMethod("POST")
		ctx.Request.Header.SetContentType("application/json")
		ctx.Request.SetBody(body)

		r.Handler(ctx)

		// Kiểm tra status code
		assert.Equal(t, 200, ctx.Response.StatusCode())

		// Parse response
		var response TestResponse
		err := json.Unmarshal(ctx.Response.Body(), &response)
		assert.NoError(t, err)

		// Kiểm tra response
		assert.Equal(t, "success", response.Status)
		assert.NotEmpty(t, response.Message)
		assert.NotNil(t, response.Data)
	})
}

func TestMain(m *testing.M) {
	// Chạy tests
	code := m.Run()

	// Cleanup
	if testClient != nil {
		database.CloseInstance(testClient)
	}

	// Exit
	os.Exit(code)
}
