package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"meta_commerce/app/database/registry"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/utility"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminService chứa các service cho User, Permission và Role
type AdminService struct {
	UserService       *BaseServiceMongoImpl[models.User]
	PermissionService *BaseServiceMongoImpl[models.Permission]
	RoleService       *BaseServiceMongoImpl[models.Role]
}

// NewAdminService tạo mới AdminService với các service tương ứng
func NewAdminService() *AdminService {
	userCollection := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Users)
	permissionCollection := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Permissions)
	roleCollection := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Roles)

	return &AdminService{
		UserService:       NewBaseServiceMongo[models.User](userCollection),
		PermissionService: NewBaseServiceMongo[models.Permission](permissionCollection),
		RoleService:       NewBaseServiceMongo[models.Role](roleCollection),
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
			return nil, utility.ErrNotFound
		}
		return nil, utility.ConvertMongoError(err)
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
			return nil, utility.ErrNotFound
		}
		return nil, utility.ConvertMongoError(err)
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

// UnBlockUser mở khóa người dùng
func (s *AdminService) UnBlockUser(ctx context.Context, email string) (*models.User, error) {
	// Tìm User theo Email
	filter := bson.M{"email": email}
	var user models.User
	err := s.UserService.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utility.ErrNotFound
		}
		return nil, utility.ConvertMongoError(err)
	}

	// Cập nhật trạng thái Block và ghi chú
	user.IsBlock = false
	user.BlockNote = ""
	user.UpdatedAt = time.Now().Unix()

	// Cập nhật User
	updatedUser, err := s.UserService.UpdateById(ctx, user.ID, user)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
