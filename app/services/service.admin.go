package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminService chứa các service cho User, Permission và Role
type AdminService struct {
	UserService       *BaseServiceImpl[models.User]
	PermissionService *BaseServiceImpl[models.Permission]
	RoleService       *BaseServiceImpl[models.Role]
}

// NewAdminService tạo mới AdminService với các service tương ứng
func NewAdminService(c *config.Configuration, db *mongo.Client) *AdminService {
	userCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Users)).Collection(global.MongoDB_ColNames.Users)
	permissionCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	roleCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)

	return &AdminService{
		UserService:       NewBaseService[models.User](userCollection),
		PermissionService: NewBaseService[models.Permission](permissionCollection),
		RoleService:       NewBaseService[models.Role](roleCollection),
	}
}

// SetRole gán Role cho User dựa trên Email và RoleID
func (s *AdminService) SetRole(ctx context.Context, email string, roleID primitive.ObjectID) (*models.User, error) {
	// Kiểm tra Role có tồn tại không
	_, err := s.RoleService.FindOneById(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Tìm User theo Email
	filter := bson.M{"email": email}
	var user models.User
	err = s.UserService.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Cập nhật Role cho User
	user.Token = roleID.Hex() // Sử dụng Token để lưu RoleID
	user.UpdatedAt = time.Now().Unix()

	// Cập nhật User
	updatedUser, err := s.UserService.UpdateById(ctx, user.ID, user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// BlockUser chặn hoặc bỏ chặn User dựa trên Email và trạng thái Block
func (s *AdminService) BlockUser(ctx context.Context, email string, block bool, note string) (*models.User, error) {
	// Tìm User theo Email
	filter := bson.M{"email": email}
	var user models.User
	err := s.UserService.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Cập nhật trạng thái Block và ghi chú
	user.IsBlock = block
	user.BlockNote = note
	user.UpdatedAt = time.Now().Unix()

	// Cập nhật User
	updatedUser, err := s.UserService.UpdateById(ctx, user.ID, user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
