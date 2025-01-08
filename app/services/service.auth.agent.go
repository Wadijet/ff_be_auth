package services

import (
	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"errors"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AgentService là cấu trúc chứa các phương thức liên quan đến người dùng
type AgentService struct {
	crudAgent RepositoryService
}

// Khởi tạo UserService với cấu hình và kết nối cơ sở dữ liệu
func NewAgentService(c *config.Configuration, db *mongo.Client) *AgentService {
	newService := new(AgentService)
	newService.crudAgent = *NewRepository(c, db, global.MongoDB_ColNames.Agents)
	return newService
}

// Tìm một Agent theo ID
func (h *AgentService) FindOneById(ctx *fasthttp.RequestCtx, id string) (FindResult interface{}, err error) {
	return h.crudAgent.FindOneById(ctx, utility.String2ObjectID(id), nil)
}

// Tìm tất cả các Agent với phân trang
func (h *AgentService) FindAll(ctx *fasthttp.RequestCtx, page int64, limit int64) (FindResult interface{}, err error) {
	// Cài đặt tùy chọn tìm kiếm
	opts := new(options.FindOptions)
	opts.SetLimit(limit)
	opts.SetSkip(page * limit)
	opts.SetSort(bson.D{{"updatedAt", 1}})

	return h.crudAgent.FindAllWithPaginate(ctx, bson.D{}, opts)
}

// Tạo mới một Agent
func (h *AgentService) Create(ctx *fasthttp.RequestCtx, credential *models.AgentCreateInput) (CreateResult interface{}, err error) {
	// Kiểm tra tên của Agent đã tồn tại chưa
	filter := bson.M{"name": credential.Name}
	checkResult, _ := h.crudAgent.FindOne(ctx, filter, nil)
	if checkResult != nil {
		return nil, errors.New("Agent already exists")
	}

	// Chuyển credential.AssignedUsers từ dạng []string sang dạng []primitive.ObjectID
	assignedUsers := make([]primitive.ObjectID, 0)
	for _, userID := range credential.AssignedUsers {
		assignedUsers = append(assignedUsers, utility.String2ObjectID(userID))
	}

	newAgent := new(models.Agent)
	newAgent.Name = credential.Name
	newAgent.Describe = credential.Describe
	newAgent.Status = 0
	newAgent.Command = 0
	newAgent.AssignedUsers = assignedUsers

	// Thêm Agent vào cơ sở dữ liệu
	return h.crudAgent.InsertOne(ctx, newAgent)
}

// Cập nhật một Agent theo ID
func (h *AgentService) Update(ctx *fasthttp.RequestCtx, id string, credential *models.AgentUpdateInput) (UpdateResult interface{}, err error) {
	// Kiểm tra Agent đã tồn tại chưa
	filter := bson.M{"_id": utility.String2ObjectID(id)}
	checkResult, _ := h.crudAgent.FindOne(ctx, filter, nil)
	if checkResult == nil {
		return nil, errors.New("Agent not found")
	}

	var agent models.Agent
	bsonBytes, err := bson.Marshal(checkResult)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonBytes, &agent)
	if err != nil {
		return nil, err
	}

	// Chuyển credential.AssignedUsers từ dạng []string sang dạng []primitive.ObjectID
	assignedUsers := make([]primitive.ObjectID, 0)
	for _, userID := range credential.AssignedUsers {
		assignedUsers = append(assignedUsers, utility.String2ObjectID(userID))
	}

	agent.Name = credential.Name
	agent.Describe = credential.Describe
	agent.Status = credential.Status
	agent.Command = credential.Command
	agent.AssignedUsers = assignedUsers

	CustomBson := &utility.CustomBson{}
	change, err := CustomBson.Set(agent)
	if err != nil {
		return nil, err
	}

	return h.crudAgent.UpdateOneById(ctx, utility.String2ObjectID(id), change)
}

// Xóa một Agent theo ID
func (h *AgentService) Delete(ctx *fasthttp.RequestCtx, id string) (DeleteResult interface{}, err error) {
	return h.crudAgent.DeleteOneById(ctx, utility.String2ObjectID(id))
}

// Hàm kiểm tra tình trạng Online của tất cả các Agent
// Duyệt qua tất cả các Agent, nếu Status = 1 và UpdateAt > 5 phút thì trả về Status = 0
func (h *AgentService) CheckOnlineStatus(ctx *fasthttp.RequestCtx) {
	// Lấy tất cả các Agent
	agents, _ := h.crudAgent.FindAll(ctx, nil, nil)

	// Duyệt qua tất cả các Agent
	for _, agent := range agents {
		// Chuyển đổi agent từ interface{} sang models.Agent
		var agentData models.Agent
		bsonBytes, err := bson.Marshal(agent)
		if err != nil {
			continue
		}

		err = bson.Unmarshal(bsonBytes, &agentData)
		if err != nil {
			continue
		}

		// Kiểm tra tình trạng Online của Agent
		if agentData.Status == 1 && ((utility.CurrentTimeInMilli() - agentData.UpdatedAt) > 300) {
			// Cập nhật tình trạng Online của Agent
			agentData.Status = 0
			CustomBson := &utility.CustomBson{}
			change, err := CustomBson.Set(agentData)
			if err != nil {
				continue
			}

			h.crudAgent.UpdateOneById(ctx, agentData.ID, change)
		}
	}
}
