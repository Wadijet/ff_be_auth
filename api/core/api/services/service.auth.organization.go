package services

import (
	"context"
	"fmt"
	models "meta_commerce/core/api/models/mongodb"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationService là cấu trúc chứa các phương thức liên quan đến tổ chức
type OrganizationService struct {
	*BaseServiceMongoImpl[models.Organization]
}

// NewOrganizationService tạo mới OrganizationService
func NewOrganizationService() (*OrganizationService, error) {
	organizationCollection, exist := global.RegistryCollections.Get(global.MongoDB_ColNames.Organizations)
	if !exist {
		return nil, fmt.Errorf("failed to get organizations collection: %v", common.ErrNotFound)
	}

	return &OrganizationService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Organization](organizationCollection),
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
		"path": bson.M{"$regex": "^" + parent.Path},
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

