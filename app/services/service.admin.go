package services

import (
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminService chứa các CRUD repository cho User, Permission và Role
type AdminService struct {
	UserCRUD       Repository
	PermissionCRUD Repository
	RoleCRUD       Repository
}

// Khởi tạo AdminService với các repository tương ứng
// trả về interface gắn với RepositoryImp
func NewAdminService(c *config.Configuration, db *mongo.Client) *AdminService {
	newService := new(AdminService)
	newService.UserCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.PermissionCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newService.RoleCRUD = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
	return newService
}

// SetRole gán Role cho User dựa trên Email và RoleID
func (service *AdminService) SetRole(ctx *fasthttp.RequestCtx, Email string, RoleID primitive.ObjectID) (SetRoleResult interface{}, err error) {

	// Tìm Role theo RoleID
	findRoleId, err := service.RoleCRUD.FindOneById(ctx, utility.ObjectID2String(RoleID), nil)
	if err != nil {
		return nil, err
	}

	var resultRole models.Role
	bsonBytes, err := bson.Marshal(findRoleId)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &resultRole)
	if err != nil {
		return nil, err
	}

	// Tìm User theo Email
	userByEmail, err := service.UserCRUD.FindOne(ctx, bson.M{"email": Email}, nil)
	if userByEmail == nil {
		return nil, err
	}

	var user models.User
	bsonBytes, err = bson.Marshal(userByEmail)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		return nil, err
	}

	// Gán Role cho User
	//user.Role = resultRole.ID

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	return service.UserCRUD.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
}

// BlockUser chặn hoặc bỏ chặn User dựa trên Email và trạng thái Block
func (service *AdminService) BlockUser(ctx *fasthttp.RequestCtx, Email string, Block bool, Note string) (AddBalanceResult interface{}, err error) {

	// Tìm User theo Email
	userByEmail, err := service.UserCRUD.FindOne(ctx, bson.M{"email": Email}, nil)
	if userByEmail == nil {
		return nil, err
	}

	var user models.User
	bsonBytes, err := bson.Marshal(userByEmail)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		return nil, err
	}

	// Cập nhật trạng thái Block và ghi chú
	user.IsBlock = Block
	user.BlockNote = Note

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(user)
	if err != nil {
		return nil, err
	}

	return service.UserCRUD.UpdateOneById(ctx, utility.ObjectID2String(user.ID), change)
}
