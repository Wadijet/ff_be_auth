package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/config"
	"atk-go-server/global"
	"errors"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService là cấu trúc chứa các phương thức liên quan đến người dùng
type OrgannizationService struct {
	crudOrgannization   Repository
	crudRole            Repository
	crudPermission      Repository
	crudRolePermissions Repository
	crudUser            Repository
	crudUserRoles       Repository
}

// Khởi tạo UserService với cấu hình và kết nối cơ sở dữ liệu
func NewOrgannizationService(c *config.Configuration, db *mongo.Client) *OrgannizationService {
	newService := new(OrgannizationService)
	newService.crudOrgannization = *NewRepository(c, db, global.MongoDB_ColNames.Organizations)
	newService.crudRole = *NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newService.crudPermission = *NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newService.crudRolePermissions = *NewRepository(c, db, global.MongoDB_ColNames.RolePermissions)
	newService.crudUser = *NewRepository(c, db, global.MongoDB_ColNames.Users)
	newService.crudUserRoles = *NewRepository(c, db, global.MongoDB_ColNames.UserRoles)
	return newService
}

// Kiểm tra Id có tồn tại hay không
func (h *OrgannizationService) IsIdExist(ctx *fasthttp.RequestCtx, id string) bool {

	// Tạo filter để tìm kiếm theo Id
	filter := bson.M{"id": id}
	result, _ := h.crudOrgannization.FindOne(ctx, filter, nil)
	if result == nil {
		return false
	} else {
		return true
	}
}

// Kiểm tra tổ chức có tổ chức con hay không
func (h *OrgannizationService) IsHaveChild(ctx *fasthttp.RequestCtx, id string) bool {
	// Kiểm tra tổ chức có tổ chức con hay
	filter := bson.M{"parent_id": id}
	result, _ := h.crudOrgannization.FindOne(ctx, filter, nil)
	if result == nil {
		return false
	} else {
		return true
	}
}

// Tạo mới một tổ chức, kiểm tra ParentId có tồn tại không trước khi thêm
func (h *OrgannizationService) Create(ctx *fasthttp.RequestCtx, input *models.OrganizationCreateInput) (InsertOneResult interface{}, err error) {

	// Nếu ParentId để trắng, thì là tổ chức gốc. Level = 1; ParentId = ""
	// Thêm biến chứa model tổ chức
	if input.ParentID == "" {

		// Tạo mới một tổ chức gốc
		var newOrgan = new(models.Organization)
		newOrgan.Name = input.Name
		newOrgan.Describe = input.Describe
		newOrgan.ParentID = ""
		newOrgan.Level = 1
		resultCreateOrgan, err := h.crudOrgannization.InsertOne(ctx, newOrgan)
		if err != nil {
			return nil, err
		}

		// Tạo mới role Administrator mặc định cho tổ chức này
		if resultCreateOrgan != nil {

			insertedOrganID := resultCreateOrgan.InsertedID

			var createdOrgan models.Organization
			bsonBytes, err := bson.Marshal(resultCreateOrgan)
			if err != nil {
				return nil, err
			}

			err = bson.Unmarshal(bsonBytes, &createdOrgan)
			if err != nil {
				return nil, err
			}

			// Tạo mới một role mặc định cho tổ chức
			var newRole = new(models.Role)
			newRole.Name = "Administrator"
			newRole.Describe = "Quản trị viên của " + input.Name
			newRole.OrganizationId = createdOrgan.ID

			resultCreateRole, err := h.crudRole.InsertOne(ctx, newRole)
			if err != nil {
				return nil, err
			}

			// Thêm cho role Administrator tất cả các quyền của tổ chức bao gồm tổ chức con, phạm vi là Subtree
			if resultCreateRole != nil {

			}

			// Gán user tạo tổ chức là Administator của tổ chức
		}
	} else {
		// Nếu ParentId không tồn tại thì báo lỗi
		// Tạo filter để tìm kiếm theo Id
		filter := bson.M{"id": input.ParentID}
		result, err := h.crudOrgannization.FindOne(ctx, filter, nil)
		if err != nil {
			return nil, err
		}

		// Nếu ParentId không tồn tại thì báo lỗi
		if result == nil {
			return nil, errors.New("ParentId not found")
		}

		if result != nil {

			var parentOrgan models.Organization
			bsonBytes, err := bson.Marshal(result)
			if err != nil {
				return nil, err
			}

			err = bson.Unmarshal(bsonBytes, &parentOrgan)
			if err != nil {
				return nil, err
			}

			var newOrgan = new(models.Organization)
			newOrgan.Name = input.Name
			newOrgan.Describe = input.Describe
			newOrgan.ParentID = input.ParentID
			newOrgan.Level = parentOrgan.Level + 1

			return h.crudOrgannization.InsertOne(ctx, newOrgan)
		}
	}

	return nil, errors.New("unexpected error")
}

// Sửa một tổ chức, kiểm tra ParentId có tồn tại không trước khi sửa
func (h *OrgannizationService) Update(ctx *fasthttp.RequestCtx, id string, input *models.OrganizationUpdateInput) (UpdateResult interface{}, err error) {

	// Nếu ParentId không tồn tại thì báo lỗi
	if !h.IsIdExist(ctx, input.ParentID) {
		return nil, errors.New("ParentId not found")
	}

	// Nếu ParentId có tồn tại thì sửa tổ chức
	return h.crudOrgannization.UpdateOneById(ctx, id, input)
}
