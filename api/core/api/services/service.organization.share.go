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

// OrganizationShareService là service quản lý sharing giữa các organizations
type OrganizationShareService struct {
	*BaseServiceMongoImpl[models.OrganizationShare]
}

// NewOrganizationShareService tạo mới OrganizationShareService
func NewOrganizationShareService() (*OrganizationShareService, error) {
	collectionName := "auth_organization_shares"
	collection, exist := global.RegistryCollections.Get(collectionName)
	if !exist {
		// Nếu chưa có, tạo mới collection từ MongoDB database
		if global.MongoDB_Session == nil {
			return nil, fmt.Errorf("MongoDB session chưa được khởi tạo")
		}
		if global.MongoDB_ServerConfig == nil {
			return nil, fmt.Errorf("MongoDB config chưa được khởi tạo")
		}

		// Lấy database
		db := global.MongoDB_Session.Database(global.MongoDB_ServerConfig.MongoDB_DBName_Auth)
		// Tạo collection
		newCollection := db.Collection(collectionName)
		// Đăng ký vào registry
		_, err := global.RegistryCollections.Register(collectionName, newCollection)
		if err != nil {
			return nil, fmt.Errorf("failed to register collection: %v", err)
		}
		collection = newCollection
	}

	return &OrganizationShareService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.OrganizationShare](collection),
	}, nil
}

// GetSharedOrganizationIDs lấy organizations được share với user's organizations
// userOrgIDs: Danh sách organization IDs của user (từ scope)
// permissionName: Permission name cụ thể (nếu rỗng = tất cả permissions)
func GetSharedOrganizationIDs(ctx context.Context, userOrgIDs []primitive.ObjectID, permissionName string) ([]primitive.ObjectID, error) {
	shareService, err := NewOrganizationShareService()
	if err != nil {
		return nil, err
	}

	if len(userOrgIDs) == 0 {
		return []primitive.ObjectID{}, nil
	}

	// Query: toOrgId trong userOrgIDs (dùng tên field trong bson tag)
	filter := bson.M{
		"toOrgId": bson.M{"$in": userOrgIDs},
	}

	// Nếu có permissionName, filter thêm
	if permissionName != "" {
		// Share nếu:
		// 1. PermissionNames rỗng/nil (share tất cả permissions)
		// 2. PermissionNames chứa permissionName cụ thể
		filter["$or"] = []bson.M{
			{"permissionNames": bson.M{"$exists": false}},                // Không có field
			{"permissionNames": bson.M{"$size": 0}},                      // Array rỗng
			{"permissionNames": bson.M{"$in": []string{permissionName}}}, // Chứa permissionName
		}
	}

	shares, err := shareService.Find(ctx, filter, nil)
	if err != nil {
		if err == common.ErrNotFound {
			return []primitive.ObjectID{}, nil
		}
		return nil, err
	}

	// Lấy fromOrgIDs (organizations share data với user)
	sharedOrgIDsMap := make(map[primitive.ObjectID]bool)
	for _, share := range shares {
		// Nếu có permissionName, kiểm tra kỹ hơn
		if permissionName != "" {
			// Nếu PermissionNames không rỗng và không chứa permissionName → skip
			if len(share.PermissionNames) > 0 {
				hasPermission := false
				for _, pn := range share.PermissionNames {
					if pn == permissionName {
						hasPermission = true
						break
					}
				}
				if !hasPermission {
					continue // Skip share này
				}
			}
		}

		sharedOrgIDsMap[share.OwnerOrganizationID] = true
	}

	// Convert to slice
	result := make([]primitive.ObjectID, 0, len(sharedOrgIDsMap))
	for orgID := range sharedOrgIDsMap {
		result = append(result, orgID)
	}

	return result, nil
}
