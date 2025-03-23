package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
)

// GetDBName lấy tên database từ cấu hình
func GetDBName(c *config.Configuration, collectionName string) string {
	return c.MongoDB_DBNameAuth
}

// BaseService định nghĩa các phương thức cơ bản mà tất cả service cần có
type BaseService[T any] interface {
	Create(ctx context.Context, data T) (T, error)
	CreateMany(ctx context.Context, data []T) ([]T, error)
	FindOne(ctx context.Context, id string) (T, error)
	FindOneByFilter(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (T, error)
	FindAll(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]T, error)
	FindAllWithPaginate(ctx context.Context, filter interface{}, page, limit int64) (*models.PaginateResult[T], error)
	Update(ctx context.Context, id string, data T) (T, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (int64, error)
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, filter interface{}) (int64, error)
	CountAll(ctx context.Context, filter interface{}) (int64, error)
}

// BaseServiceImpl là implementation mặc định cho BaseService
type BaseServiceImpl[T any] struct {
	collection *mongo.Collection
}

// NewBaseService tạo mới một BaseServiceImpl
func NewBaseService[T any](collection *mongo.Collection) *BaseServiceImpl[T] {
	return &BaseServiceImpl[T]{
		collection: collection,
	}
}

// Create tạo mới một bản ghi
func (s *BaseServiceImpl[T]) Create(ctx context.Context, data T) (T, error) {
	var zero T

	// Chuyển data thành map để thêm timestamps
	dataMap, err := utility.ToMap(data)
	if err != nil {
		return zero, err
	}

	// Thêm timestamps
	now := time.Now().UnixMilli()
	dataMap["createdAt"] = now
	dataMap["updatedAt"] = now

	result, err := s.collection.InsertOne(ctx, dataMap)
	if err != nil {
		return zero, err
	}

	// Lấy lại document vừa tạo
	var created T
	err = s.collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&created)
	if err != nil {
		return zero, err
	}

	return created, nil
}

// CreateMany tạo nhiều bản ghi
func (s *BaseServiceImpl[T]) CreateMany(ctx context.Context, data []T) ([]T, error) {
	var documents []interface{}
	now := time.Now().UnixMilli()

	for _, item := range data {
		dataMap, err := utility.ToMap(item)
		if err != nil {
			return nil, err
		}
		dataMap["createdAt"] = now
		dataMap["updatedAt"] = now
		documents = append(documents, dataMap)
	}

	result, err := s.collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}

	// Lấy lại các documents vừa tạo
	var created []T
	filter := bson.M{"_id": bson.M{"$in": result.InsertedIDs}}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &created)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// FindOne tìm một bản ghi theo ID
func (s *BaseServiceImpl[T]) FindOne(ctx context.Context, id string) (T, error) {
	var zero T
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return zero, err
	}

	// Khởi tạo options nếu nil
	opts := options.FindOne()

	var result T
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return zero, errors.New("record not found")
		}
		return zero, err
	}

	return result, nil
}

// FindOneByFilter tìm một bản ghi theo filter
func (s *BaseServiceImpl[T]) FindOneByFilter(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (T, error) {
	var zero T
	var result T

	if filter == nil {
		filter = bson.D{}
	}

	// Khởi tạo options nếu nil
	if opts == nil {
		opts = options.FindOne()
	}

	// Thực hiện FindOne và lưu kết quả
	findResult := s.collection.FindOne(ctx, filter, opts)

	// Kiểm tra lỗi từ FindOne
	if err := findResult.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return zero, errors.New("record not found")
		}
		return zero, err
	}

	// Thực hiện Decode
	if err := findResult.Decode(&result); err != nil {
		return zero, err
	}

	return result, nil
}

// FindAll tìm tất cả bản ghi theo filter
func (s *BaseServiceImpl[T]) FindAll(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]T, error) {
	if filter == nil {
		filter = bson.D{}
	}

	// Khởi tạo options nếu nil
	if opts == nil {
		opts = options.Find()
	}

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// FindAllWithPaginate tìm tất cả bản ghi với phân trang
func (s *BaseServiceImpl[T]) FindAllWithPaginate(ctx context.Context, filter interface{}, page, limit int64) (*models.PaginateResult[T], error) {
	if filter == nil {
		filter = bson.D{}
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit)

	// Lấy tổng số bản ghi
	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Lấy dữ liệu theo trang
	var items []T
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return &models.PaginateResult[T]{
		Items:     items,
		Page:      page,
		Limit:     limit,
		ItemCount: int64(len(items)),
		Total:     total,
		TotalPage: (total + limit - 1) / limit,
	}, nil
}

// Update cập nhật một bản ghi
func (s *BaseServiceImpl[T]) Update(ctx context.Context, id string, data T) (T, error) {
	var zero T
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return zero, err
	}

	// Chuyển data thành map để cập nhật
	dataMap, err := utility.ToMap(data)
	if err != nil {
		return zero, err
	}

	// Cập nhật timestamp
	dataMap["updatedAt"] = time.Now().UnixMilli()

	update := bson.M{"$set": dataMap}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return zero, err
	}

	// Khởi tạo options nếu nil
	opts := options.FindOne()

	// Lấy lại document đã cập nhật
	var updated T
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}, opts).Decode(&updated)
	if err != nil {
		return zero, err
	}

	return updated, nil
}

// UpdateMany cập nhật nhiều bản ghi
func (s *BaseServiceImpl[T]) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	// Thêm updatedAt vào update
	updateMap, err := utility.ToMap(update)
	if err != nil {
		return 0, err
	}

	if setMap, ok := updateMap["$set"].(map[string]interface{}); ok {
		setMap["updatedAt"] = time.Now().UnixMilli()
		updateMap["$set"] = setMap
	} else {
		updateMap["$set"] = bson.M{"updatedAt": time.Now().UnixMilli()}
	}

	result, err := s.collection.UpdateMany(ctx, filter, updateMap)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

// Delete xóa một bản ghi
func (s *BaseServiceImpl[T]) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("record not found")
	}

	return nil
}

// DeleteMany xóa nhiều bản ghi theo điều kiện
func (s *BaseServiceImpl[T]) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	result, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

// CountAll đếm số lượng bản ghi theo filter
func (s *BaseServiceImpl[T]) CountAll(ctx context.Context, filter interface{}) (int64, error) {
	return s.collection.CountDocuments(ctx, filter)
}
