package services

import (
	"context"
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoopts "go.mongodb.org/mongo-driver/mongo/options"
)

// OrganizationService là cấu trúc chứa các phương thức liên quan đến tổ chức
type OrganizationService struct {
	*BaseServiceMongoImpl[models.Organization]
	roleService *RoleService
}

// NewOrganizationService tạo mới OrganizationService
func NewOrganizationService() (*OrganizationService, error) {
	organizationCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.Organizations)
	if !exist {
		return nil, fmt.Errorf("failed to get organizations collection: %v", common.ErrNotFound)
	}

	roleService, err := NewRoleService()
	if err != nil {
		return nil, fmt.Errorf("failed to create role service: %v", err)
	}

	return &OrganizationService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Organization](organizationCollection),
		roleService:          roleService,
	}, nil
}

// GetChildrenIDs lấy tất cả ID của organization con (dùng cho Scope = 1)
func (s *OrganizationService) GetChildrenIDs(ctx context.Context, parentID primitive.ObjectID) ([]primitive.ObjectID, error) {
	// Lấy organization cha
	parent, err := s.FindOneById(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// Query tất cả organization có Path bắt đầu với parent.Path
	filter := bson.M{
		"path":     bson.M{"$regex": "^" + parent.Path},
		"isActive": true,
	}

	orgs, err := s.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, 0, len(orgs))
	for _, org := range orgs {
		ids = append(ids, org.ID)
	}

	return ids, nil
}

// validateBeforeDelete kiểm tra các điều kiện trước khi xóa organization
// - Không cho phép xóa System organization
// - Không cho phép xóa nếu có tổ chức con
// - Không cho phép xóa nếu có role trực thuộc
func (s *OrganizationService) validateBeforeDelete(ctx context.Context, orgID primitive.ObjectID) error {
	// Lấy thông tin organization cần xóa
	org, err := s.FindOneById(ctx, orgID)
	if err != nil {
		return err
	}

	var modelOrg models.Organization
	bsonBytes, _ := bson.Marshal(org)
	err = bson.Unmarshal(bsonBytes, &modelOrg)
	if err != nil {
		return common.ErrInvalidFormat
	}

	// Kiểm tra 1: Nếu là System organization thì không cho phép xóa
	if modelOrg.Type == models.OrganizationTypeSystem && modelOrg.Code == "SYSTEM" && modelOrg.Level == -1 {
		return common.NewError(
			common.ErrCodeBusinessOperation,
			"Không thể xóa System organization. Đây là tổ chức cấp cao nhất chứa Administrator và không thể xóa.",
			common.StatusForbidden,
			nil,
		)
	}

	// Kiểm tra 2: Kiểm tra xem có tổ chức con không
	// Tìm các tổ chức có parentId = org.ID hoặc path bắt đầu với org.Path + "/"
	childrenFilter := bson.M{
		"$or": []bson.M{
			{"parentId": modelOrg.ID},
			{"path": bson.M{"$regex": "^" + modelOrg.Path + "/"}},
		},
	}
	children, err := s.Find(ctx, childrenFilter, nil)
	if err != nil && err != common.ErrNotFound {
		return err
	}
	if err == nil && len(children) > 0 {
		return common.NewError(
			common.ErrCodeBusinessOperation,
			fmt.Sprintf("Không thể xóa tổ chức '%s' vì có %d tổ chức con. Vui lòng xóa hoặc di chuyển các tổ chức con trước.", modelOrg.Name, len(children)),
			common.StatusConflict,
			nil,
		)
	}

	// Kiểm tra 3: Kiểm tra xem có role nào trực thuộc tổ chức này không
	rolesFilter := bson.M{
		"organizationId": modelOrg.ID,
	}
	roles, err := s.roleService.Find(ctx, rolesFilter, nil)
	if err != nil && err != common.ErrNotFound {
		return err
	}
	if err == nil && len(roles) > 0 {
		return common.NewError(
			common.ErrCodeBusinessOperation,
			fmt.Sprintf("Không thể xóa tổ chức '%s' vì có %d role trực thuộc. Vui lòng xóa hoặc di chuyển các role trước.", modelOrg.Name, len(roles)),
			common.StatusConflict,
			nil,
		)
	}

	return nil
}

// DeleteOne override method DeleteOne để kiểm tra trước khi xóa
func (s *OrganizationService) DeleteOne(ctx context.Context, filter interface{}) error {
	// Lấy thông tin organization cần xóa
	org, err := s.BaseServiceMongoImpl.FindOne(ctx, filter, nil)
	if err != nil {
		return err
	}

	var modelOrg models.Organization
	bsonBytes, _ := bson.Marshal(org)
	err = bson.Unmarshal(bsonBytes, &modelOrg)
	if err != nil {
		return common.ErrInvalidFormat
	}

	// Kiểm tra trước khi xóa
	if err := s.validateBeforeDelete(ctx, modelOrg.ID); err != nil {
		return err
	}

	// Thực hiện xóa nếu không có ràng buộc
	return s.BaseServiceMongoImpl.DeleteOne(ctx, filter)
}

// DeleteById override method DeleteById để kiểm tra trước khi xóa
func (s *OrganizationService) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	// Kiểm tra trước khi xóa
	if err := s.validateBeforeDelete(ctx, id); err != nil {
		return err
	}

	// Thực hiện xóa nếu không có ràng buộc
	return s.BaseServiceMongoImpl.DeleteById(ctx, id)
}

// DeleteMany override method DeleteMany để kiểm tra trước khi xóa
func (s *OrganizationService) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	// Lấy danh sách organizations sẽ bị xóa
	orgs, err := s.BaseServiceMongoImpl.Find(ctx, filter, nil)
	if err != nil && err != common.ErrNotFound {
		return 0, err
	}

	// Kiểm tra từng organization trước khi xóa
	for _, org := range orgs {
		var modelOrg models.Organization
		bsonBytes, _ := bson.Marshal(org)
		if err := bson.Unmarshal(bsonBytes, &modelOrg); err != nil {
			continue
		}

		// Kiểm tra trước khi xóa
		if err := s.validateBeforeDelete(ctx, modelOrg.ID); err != nil {
			return 0, err
		}
	}

	// Thực hiện xóa nếu không có ràng buộc
	return s.BaseServiceMongoImpl.DeleteMany(ctx, filter)
}

// FindOneAndDelete override method FindOneAndDelete để kiểm tra trước khi xóa
func (s *OrganizationService) FindOneAndDelete(ctx context.Context, filter interface{}, opts *mongoopts.FindOneAndDeleteOptions) (models.Organization, error) {
	var zero models.Organization

	// Lấy thông tin organization sẽ bị xóa
	org, err := s.BaseServiceMongoImpl.FindOne(ctx, filter, nil)
	if err != nil {
		return zero, err
	}

	var modelOrg models.Organization
	bsonBytes, _ := bson.Marshal(org)
	err = bson.Unmarshal(bsonBytes, &modelOrg)
	if err != nil {
		return zero, common.ErrInvalidFormat
	}

	// Kiểm tra trước khi xóa
	if err := s.validateBeforeDelete(ctx, modelOrg.ID); err != nil {
		return zero, err
	}

	// Thực hiện xóa nếu không có ràng buộc
	return s.BaseServiceMongoImpl.FindOneAndDelete(ctx, filter, opts)
}
