package services

import (
	"context"
	"fmt"
	"meta_commerce/core/common"
	"meta_commerce/core/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RelationshipCheck định nghĩa một quan hệ cần kiểm tra
type RelationshipCheck struct {
	// CollectionName: Tên collection cần kiểm tra
	CollectionName string
	// FieldName: Tên field trong collection đó trỏ tới record hiện tại
	FieldName string
	// ErrorMessage: Thông báo lỗi khi tìm thấy quan hệ (có thể dùng %d để thay thế số lượng)
	ErrorMessage string
	// Optional: Nếu true, sẽ không trả về lỗi nếu không tìm thấy collection (dùng cho optional relationships)
	Optional bool
}

// CheckRelationshipExists kiểm tra xem có record nào trong collection khác đang trỏ tới record này không
//
// Parameters:
//   - ctx: Context
//   - recordID: ID của record cần kiểm tra
//   - checks: Danh sách các quan hệ cần kiểm tra
//
// Returns:
//   - error: Lỗi nếu tìm thấy quan hệ hoặc lỗi khác
func CheckRelationshipExists(ctx context.Context, recordID primitive.ObjectID, checks []RelationshipCheck) error {
	for _, check := range checks {
		// Lấy collection
		collection, exists := global.RegistryCollections.Get(check.CollectionName)
		if !exists {
			if check.Optional {
				// Nếu là optional và không tìm thấy collection, bỏ qua
				continue
			}
			return common.NewError(
				common.ErrCodeInternalServer,
				fmt.Sprintf("Không tìm thấy collection '%s' để kiểm tra quan hệ", check.CollectionName),
				common.StatusInternalServerError,
				nil,
			)
		}

		// Tạo filter để tìm các record có field trỏ tới recordID
		filter := bson.M{check.FieldName: recordID}

		// Đếm số lượng record
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return common.ConvertMongoError(err)
		}

		// Nếu tìm thấy record, trả về lỗi
		if count > 0 {
			errorMsg := check.ErrorMessage
			if errorMsg == "" {
				errorMsg = fmt.Sprintf("Không thể xóa record vì có %d record trong collection '%s' đang tham chiếu tới record này", count, check.CollectionName)
			} else {
				// Thay thế %d bằng số lượng nếu có
				errorMsg = fmt.Sprintf(check.ErrorMessage, count)
			}

			return common.NewError(
				common.ErrCodeBusinessOperation,
				errorMsg,
				common.StatusConflict,
				nil,
			)
		}
	}

	return nil
}

// CheckRelationshipExistsWithFilter kiểm tra quan hệ với filter tùy chỉnh
//
// Parameters:
//   - ctx: Context
//   - filter: Filter tùy chỉnh để kiểm tra
//   - checks: Danh sách các quan hệ cần kiểm tra (chỉ dùng CollectionName và ErrorMessage)
//
// Returns:
//   - error: Lỗi nếu tìm thấy quan hệ hoặc lỗi khác
func CheckRelationshipExistsWithFilter(ctx context.Context, filter bson.M, checks []RelationshipCheck) error {
	for _, check := range checks {
		// Lấy collection
		collection, exists := global.RegistryCollections.Get(check.CollectionName)
		if !exists {
			if check.Optional {
				continue
			}
			return common.NewError(
				common.ErrCodeInternalServer,
				fmt.Sprintf("Không tìm thấy collection '%s' để kiểm tra quan hệ", check.CollectionName),
				common.StatusInternalServerError,
				nil,
			)
		}

		// Đếm số lượng record với filter tùy chỉnh
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return common.ConvertMongoError(err)
		}

		// Nếu tìm thấy record, trả về lỗi
		if count > 0 {
			errorMsg := check.ErrorMessage
			if errorMsg == "" {
				errorMsg = fmt.Sprintf("Không thể xóa record vì có %d record trong collection '%s' đang tham chiếu tới record này", count, check.CollectionName)
			} else {
				errorMsg = fmt.Sprintf(check.ErrorMessage, count)
			}

			return common.NewError(
				common.ErrCodeBusinessOperation,
				errorMsg,
				common.StatusConflict,
				nil,
			)
		}
	}

	return nil
}

// GetRelationshipCount trả về số lượng record đang tham chiếu tới record này
//
// Parameters:
//   - ctx: Context
//   - recordID: ID của record cần kiểm tra
//   - collectionName: Tên collection cần kiểm tra
//   - fieldName: Tên field trong collection đó trỏ tới record hiện tại
//
// Returns:
//   - int64: Số lượng record đang tham chiếu
//   - error: Lỗi nếu có
func GetRelationshipCount(ctx context.Context, recordID primitive.ObjectID, collectionName, fieldName string) (int64, error) {
	collection, exists := global.RegistryCollections.Get(collectionName)
	if !exists {
		return 0, common.NewError(
			common.ErrCodeInternalServer,
			fmt.Sprintf("Không tìm thấy collection '%s'", collectionName),
			common.StatusInternalServerError,
			nil,
		)
	}

	filter := bson.M{fieldName: recordID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, common.ConvertMongoError(err)
	}

	return count, nil
}

// ValidateBeforeDeleteRole kiểm tra các quan hệ của Role trước khi xóa
// Đây là helper function chuyên dụng cho Role
func ValidateBeforeDeleteRole(ctx context.Context, roleID primitive.ObjectID) error {
	checks := []RelationshipCheck{
		{
			CollectionName: global.MongoDB_ColNames.UserRoles,
			FieldName:      "roleId",
			ErrorMessage:   "Không thể xóa role vì có %d user đang sử dụng role này. Vui lòng gỡ role khỏi các user trước.",
		},
		{
			CollectionName: global.MongoDB_ColNames.RolePermissions,
			FieldName:      "roleId",
			ErrorMessage:   "Không thể xóa role vì có %d permission đang được gán cho role này. Vui lòng gỡ các permission trước.",
		},
	}

	return CheckRelationshipExists(ctx, roleID, checks)
}

// ValidateBeforeDeleteOrganization kiểm tra các quan hệ của Organization trước khi xóa
// Đây là helper function chuyên dụng cho Organization
// Lưu ý: Function này chỉ kiểm tra quan hệ trực tiếp, các kiểm tra phức tạp hơn (như children)
// nên được thực hiện trong OrganizationService.validateBeforeDelete
func ValidateBeforeDeleteOrganization(ctx context.Context, orgID primitive.ObjectID) error {
	checks := []RelationshipCheck{
		{
			CollectionName: global.MongoDB_ColNames.Roles,
			FieldName:      "organizationId",
			ErrorMessage:   "Không thể xóa tổ chức vì có %d role trực thuộc. Vui lòng xóa hoặc di chuyển các role trước.",
		},
	}

	return CheckRelationshipExists(ctx, orgID, checks)
}

// ValidateBeforeDeletePermission kiểm tra các quan hệ của Permission trước khi xóa
func ValidateBeforeDeletePermission(ctx context.Context, permissionID primitive.ObjectID) error {
	checks := []RelationshipCheck{
		{
			CollectionName: global.MongoDB_ColNames.RolePermissions,
			FieldName:      "permissionId",
			ErrorMessage:   "Không thể xóa permission vì có %d role đang sử dụng permission này. Vui lòng gỡ permission khỏi các role trước.",
		},
	}

	return CheckRelationshipExists(ctx, permissionID, checks)
}

// ValidateBeforeDeleteUser kiểm tra các quan hệ của User trước khi xóa
func ValidateBeforeDeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	checks := []RelationshipCheck{
		{
			CollectionName: global.MongoDB_ColNames.UserRoles,
			FieldName:      "userId",
			ErrorMessage:   "Không thể xóa user vì có %d role đang được gán cho user này. Vui lòng gỡ các role trước.",
		},
	}

	return CheckRelationshipExists(ctx, userID, checks)
}
