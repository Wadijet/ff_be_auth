package services

import (
	"context"
	"time"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AgentService là cấu trúc chứa các phương thức liên quan đến trợ lý
type AgentService struct {
	*BaseServiceImpl[models.Agent]
}

// NewAgentService tạo mới AgentService
func NewAgentService(c *config.Configuration, db *mongo.Client) *AgentService {
	agentCollection := db.Database(GetDBName(c, global.MongoDB_ColNames.Agents)).Collection(global.MongoDB_ColNames.Agents)
	return &AgentService{
		BaseServiceImpl: NewBaseService[models.Agent](agentCollection),
	}
}

// IsNameExist kiểm tra tên trợ lý có tồn tại hay không
func (s *AgentService) IsNameExist(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}
	var agent models.Agent
	err := s.BaseServiceImpl.collection.FindOne(ctx, filter).Decode(&agent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create tạo mới một trợ lý
func (s *AgentService) Create(ctx context.Context, input *models.AgentCreateInput) (*models.Agent, error) {
	// Kiểm tra tên tồn tại
	exists, err := s.IsNameExist(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Agent already exists")
	}

	// Chuyển input.AssignedUsers từ dạng []string sang dạng []primitive.ObjectID
	assignedUsers := make([]primitive.ObjectID, 0)
	for _, userID := range input.AssignedUsers {
		assignedUsers = append(assignedUsers, utility.String2ObjectID(userID))
	}

	// Tạo agent mới
	agent := &models.Agent{
		ID:            primitive.NewObjectID(),
		Name:          input.Name,
		Describe:      input.Describe,
		Status:        0,
		Command:       0,
		AssignedUsers: assignedUsers,
		ConfigData:    input.ConfigData,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	// Lưu agent
	createdAgent, err := s.BaseServiceImpl.Create(ctx, *agent)
	if err != nil {
		return nil, err
	}

	return &createdAgent, nil
}

// Update cập nhật thông tin trợ lý
func (s *AgentService) Update(ctx context.Context, id string, input *models.AgentUpdateInput) (*models.Agent, error) {
	// Kiểm tra agent tồn tại
	agent, err := s.BaseServiceImpl.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Nếu có thay đổi tên, kiểm tra tên mới
	if input.Name != "" && input.Name != agent.Name {
		exists, err := s.IsNameExist(ctx, input.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("Agent name already exists")
		}
		agent.Name = input.Name
	}

	// Cập nhật thông tin khác
	if input.Describe != "" {
		agent.Describe = input.Describe
	}
	if input.Status != 0 {
		agent.Status = input.Status
	}
	if input.Command != 0 {
		agent.Command = input.Command
	}
	if len(input.AssignedUsers) > 0 {
		assignedUsers := make([]primitive.ObjectID, 0)
		for _, userID := range input.AssignedUsers {
			assignedUsers = append(assignedUsers, utility.String2ObjectID(userID))
		}
		agent.AssignedUsers = assignedUsers
	}
	if input.ConfigData != nil {
		agent.ConfigData = input.ConfigData
	}
	agent.UpdatedAt = time.Now().Unix()

	// Cập nhật agent
	updatedAgent, err := s.BaseServiceImpl.Update(ctx, id, agent)
	if err != nil {
		return nil, err
	}

	return &updatedAgent, nil
}

// Delete xóa trợ lý
func (s *AgentService) Delete(ctx context.Context, id string) error {
	return s.BaseServiceImpl.Delete(ctx, id)
}

// CheckOnlineStatus kiểm tra tình trạng Online của tất cả các trợ lý
func (s *AgentService) CheckOnlineStatus(ctx context.Context) error {
	// Lấy tất cả các agent
	opts := options.Find()
	agents, err := s.BaseServiceImpl.FindAll(ctx, bson.M{}, opts)
	if err != nil {
		return err
	}

	// Duyệt qua tất cả các agent
	for _, agent := range agents {
		// Kiểm tra tình trạng Online của Agent
		if agent.Status == 1 && ((utility.CurrentTimeInMilli() - agent.UpdatedAt) > 300) {
			// Cập nhật tình trạng Online của Agent
			agent.Status = 0
			agent.UpdatedAt = time.Now().Unix()

			_, err := s.BaseServiceImpl.Update(ctx, agent.ID.Hex(), agent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckIn điểm danh cho một trợ lý
func (s *AgentService) CheckIn(ctx context.Context, id string) (*models.Agent, error) {
	// Kiểm tra agent tồn tại
	agent, err := s.BaseServiceImpl.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cập nhật tình trạng Online của Agent
	agent.Status = 1
	agent.UpdatedAt = time.Now().Unix()

	// Cập nhật agent
	updatedAgent, err := s.BaseServiceImpl.Update(ctx, id, agent)
	if err != nil {
		return nil, err
	}

	return &updatedAgent, nil
}
