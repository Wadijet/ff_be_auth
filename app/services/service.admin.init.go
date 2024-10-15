package services

import (
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"atk-go-server/app/models"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
)

// InitService định nghĩa các CRUD repository cho User, Permission và Role
type InitService struct {
	UserCRUD       Repository
	PermissionCRUD Repository
	RoleCRUD       Repository
}

// NewInitService khởi tạo các repository và trả về một đối tượng InitService
func NewInitService(c *config.Configuration, db *mongo.Client) *InitService {
	newService := new(InitService)
	newService.UserCRUD = *NewRepository(c, db, global.ColNames.Users)
	newService.PermissionCRUD = *NewRepository(c, db, global.ColNames.Permissions)
	newService.RoleCRUD = *NewRepository(c, db, global.ColNames.Roles)
	return newService
}

// SetAdmin thêm quyền Admin cho người dùng dựa trên email
func (service *InitService) SetAdmin(ctx *fasthttp.RequestCtx, email string) (InsertOneResult interface{}, err error) {

	// Tìm người dùng theo email
	userByEmail, err := service.UserCRUD.FindOne(ctx, bson.M{"email": email}, nil)
	if userByEmail == nil {
		return nil, err
	}

	// Chuyển đổi kết quả tìm kiếm thành model User
	var user models.User
	bsonBytes, err := bson.Marshal(userByEmail)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		return nil, err
	}

	// Tìm role Admin
	adminRole, err := service.RoleCRUD.FindOne(ctx, bson.M{"name": "Admin"}, nil)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi role Admin thành map để lấy ID
	adminRoleMap, err := utility.ToMap(adminRole)
	if err != nil {
		return nil, err
	}

	adminRoleID := adminRoleMap["_id"].(primitive.ObjectID)
	user.Role = adminRoleID

	// Tạo thay đổi cần thiết để cập nhật người dùng
	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	// Cập nhật người dùng với role Admin
	return service.UserCRUD.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
}
