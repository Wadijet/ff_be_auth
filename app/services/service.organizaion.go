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
	crudOrgannization Repository
}

// Khởi tạo UserService với cấu hình và kết nối cơ sở dữ liệu
func NewOrgannizationService(c *config.Configuration, db *mongo.Client) *OrgannizationService {
	newService := new(OrgannizationService)
	newService.crudOrgannization = *NewRepository(c, db, global.MongoDB_ColNames.Organizations)
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

	// Nếu ParentId không tồn tại thì báo lỗi
	if !h.IsIdExist(ctx, input.ParentID) {
		return nil, errors.New("ParentId not found")
	}

	// Nếu ParentId có tồn tại thì thêm tổ chức
	return h.crudOrgannization.InsertOne(ctx, input)
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
