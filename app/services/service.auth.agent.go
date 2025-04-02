package services

import (
	"context"
	"time"

	"meta_commerce/app/database/registry"
	"meta_commerce/app/global"
	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/utility"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AgentService là cấu trúc chứa các phương thức liên quan đến trợ lý
type AgentService struct {
	*BaseServiceMongoImpl[models.Agent]
}

// NewAgentService tạo mới AgentService
func NewAgentService() *AgentService {
	agentCollection := registry.GetRegistry().MustGetCollection(global.MongoDB_ColNames.Agents)
	return &AgentService{
		BaseServiceMongoImpl: NewBaseServiceMongo[models.Agent](agentCollection),
	}
}

// CheckOnlineStatus kiểm tra tình trạng Online của tất cả các trợ lý
func (s *AgentService) CheckOnlineStatus(ctx context.Context) error {
	// Lấy tất cả các agent
	opts := options.Find()
	agents, err := s.BaseServiceMongoImpl.Find(ctx, bson.M{}, opts)
	if err != nil {
		return utility.ConvertMongoError(err)
	}

	// Duyệt qua tất cả các agent
	for _, agent := range agents {
		// Kiểm tra tình trạng Online của Agent
		if agent.Status == 1 && ((utility.CurrentTimeInMilli() - agent.UpdatedAt) > 300) {
			// Cập nhật tình trạng Online của Agent
			agent.Status = 0
			agent.UpdatedAt = time.Now().Unix()

			_, err := s.BaseServiceMongoImpl.UpdateById(ctx, agent.ID, agent)
			if err != nil {
				return utility.ConvertMongoError(err)
			}
		}
	}

	return nil
}

// CheckIn điểm danh cho một trợ lý
func (s *AgentService) CheckIn(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	// Kiểm tra agent tồn tại
	agent, err := s.BaseServiceMongoImpl.FindOneById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cập nhật tình trạng Online của Agent
	agent.Status = 1
	agent.UpdatedAt = time.Now().Unix()

	// Cập nhật agent
	updatedAgent, err := s.BaseServiceMongoImpl.UpdateById(ctx, id, agent)
	if err != nil {
		return nil, err
	}

	return &updatedAgent, nil
}

// CheckOut điểm danh cho một trợ lý
func (s *AgentService) CheckOut(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	// Kiểm tra agent tồn tại
	agent, err := s.BaseServiceMongoImpl.FindOneById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cập nhật tình trạng Online của Agent
	agent.Status = 0
	agent.UpdatedAt = time.Now().Unix()

	// Cập nhật agent
	updatedAgent, err := s.BaseServiceMongoImpl.UpdateById(ctx, id, agent)
	if err != nil {
		return nil, err
	}

	return &updatedAgent, nil
}
