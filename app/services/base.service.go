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
	Create(ctx context.Context, data T) (T, error)                                                                                  // Tạo mới một bản ghi
	CreateMany(ctx context.Context, data []T) ([]T, error)                                                                          // Tạo nhiều bản ghi
	FindOne(ctx context.Context, id primitive.ObjectID) (T, error)                                                                  // Tìm một document theo ObjectId
	FindOneByFilter(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (T, error)                               // Tìm một bản ghi theo filter
	FindAll(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]T, error)                                        // Tìm tất cả bản ghi theo filter
	FindAllWithPaginate(ctx context.Context, filter interface{}, page, limit int64) (*models.PaginateResult[T], error)              // Tìm tất cả bản ghi với phân trang
	Update(ctx context.Context, id primitive.ObjectID, data T) (T, error)                                                           // Cập nhật một document theo ObjectId
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (int64, error)                                          // Cập nhật nhiều bản ghi
	Delete(ctx context.Context, id primitive.ObjectID) error                                                                        // Xóa một document theo ObjectId
	DeleteMany(ctx context.Context, filter interface{}) (int64, error)                                                              // Xóa nhiều bản ghi theo điều kiện
	CountAll(ctx context.Context, filter interface{}) (int64, error)                                                                // Đếm số lượng bản ghi theo filter
	Upsert(ctx context.Context, filter interface{}, data T) (T, error)                                                              // Upsert thực hiện thao tác update nếu tồn tại, insert nếu chưa tồn tại
	UpsertMany(ctx context.Context, filter interface{}, data []T) ([]T, error)                                                      // UpsertMany thực hiện thao tác upsert cho nhiều document
	FindByIds(ctx context.Context, ids []primitive.ObjectID) ([]T, error)                                                           // Tìm nhiều document theo danh sách ID
	Exists(ctx context.Context, filter interface{}) (bool, error)                                                                   // Kiểm tra xem một document có tồn tại không
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts *options.FindOneAndUpdateOptions) (T, error) // Tìm và cập nhật một document
	FindOneAndDelete(ctx context.Context, filter interface{}, opts *options.FindOneAndDeleteOptions) (T, error)                     // Tìm và xóa một document
	Distinct(ctx context.Context, fieldName string, filter interface{}) ([]interface{}, error)                                      // Lấy danh sách các giá trị duy nhất của một trường
}

// BaseServiceImpl định nghĩa các phương thức cơ bản cho service
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

// FindOne tìm một document theo ObjectId
func (s *BaseServiceImpl[T]) FindOne(ctx context.Context, id primitive.ObjectID) (T, error) {
	var result T
	filter := bson.M{"_id": id}
	err := s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
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

// Update cập nhật một document theo ObjectId
func (s *BaseServiceImpl[T]) Update(ctx context.Context, id primitive.ObjectID, data T) (T, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": data}

	_, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return data, err
	}

	return s.FindOne(ctx, id)
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

// Delete xóa một document theo ObjectId
func (s *BaseServiceImpl[T]) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := s.collection.DeleteOne(ctx, filter)
	return err
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

// Upsert thực hiện thao tác update nếu tồn tại, insert nếu chưa tồn tại
func (s *BaseServiceImpl[T]) Upsert(ctx context.Context, filter interface{}, data T) (T, error) {
	var zero T

	// Chuyển data thành map để thêm timestamps
	dataMap, err := utility.ToMap(data)
	if err != nil {
		return zero, err
	}

	// Thêm timestamps
	now := time.Now().UnixMilli()
	dataMap["updatedAt"] = now

	// Tạo options cho upsert với sort để đảm bảo chỉ update một document
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After).
		SetSort(bson.D{{Key: "_id", Value: 1}}) // Sắp xếp theo _id để đảm bảo tính nhất quán

	// Thực hiện upsert và lấy document sau khi update
	var upserted T
	err = s.collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": dataMap}, opts).Decode(&upserted)
	if err != nil {
		return zero, err
	}

	return upserted, nil
}

// UpsertMany thực hiện thao tác upsert cho nhiều document
func (s *BaseServiceImpl[T]) UpsertMany(ctx context.Context, filter interface{}, data []T) ([]T, error) {
	if len(data) == 0 {
		return []T{}, nil
	}

	// Tạo các models cho bulk write
	var models []mongo.WriteModel
	now := time.Now().UnixMilli()

	for _, item := range data {
		// Chuyển data thành map để thêm timestamps
		dataMap, err := utility.ToMap(item)
		if err != nil {
			return nil, err
		}

		// Thêm timestamps
		dataMap["updatedAt"] = now

		// Tạo upsert model
		upsertModel := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(bson.M{"$set": dataMap}).
			SetUpsert(true)

		models = append(models, upsertModel)
	}

	// Thực hiện bulk write
	opts := options.BulkWrite().SetOrdered(false) // SetOrdered(false) để thực hiện song song
	result, err := s.collection.BulkWrite(ctx, models, opts)
	if err != nil {
		return nil, err
	}

	// Lấy lại các documents sau khi upsert
	var upserted []T
	if result.UpsertedCount > 0 {
		// Nếu có documents mới được tạo
		var upsertedIDs []primitive.ObjectID
		for _, id := range result.UpsertedIDs {
			if objectID, ok := id.(primitive.ObjectID); ok {
				upsertedIDs = append(upsertedIDs, objectID)
			}
		}

		if len(upsertedIDs) > 0 {
			cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": upsertedIDs}})
			if err != nil {
				return nil, err
			}
			defer cursor.Close(ctx)

			if err = cursor.All(ctx, &upserted); err != nil {
				return nil, err
			}
		}
	}

	// Lấy các documents đã được update
	if result.ModifiedCount > 0 {
		cursor, err := s.collection.Find(ctx, filter)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		var updated []T
		if err = cursor.All(ctx, &updated); err != nil {
			return nil, err
		}

		// Kết hợp cả documents mới và documents đã update
		upserted = append(upserted, updated...)
	}

	return upserted, nil
}

// FindByIds tìm nhiều document theo danh sách ID
func (s *BaseServiceImpl[T]) FindByIds(ctx context.Context, ids []primitive.ObjectID) ([]T, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := s.collection.Find(ctx, filter)
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

// Exists kiểm tra xem một document có tồn tại không
func (s *BaseServiceImpl[T]) Exists(ctx context.Context, filter interface{}) (bool, error) {
	if filter == nil {
		filter = bson.D{}
	}

	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindOneAndUpdate tìm và cập nhật một document
func (s *BaseServiceImpl[T]) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts *options.FindOneAndUpdateOptions) (T, error) {
	var zero T

	if filter == nil {
		filter = bson.D{}
	}

	if opts == nil {
		opts = options.FindOneAndUpdate()
	}

	// Thêm updatedAt vào update
	updateMap, err := utility.ToMap(update)
	if err != nil {
		return zero, err
	}

	if setMap, ok := updateMap["$set"].(map[string]interface{}); ok {
		setMap["updatedAt"] = time.Now().UnixMilli()
		updateMap["$set"] = setMap
	} else {
		updateMap["$set"] = bson.M{"updatedAt": time.Now().UnixMilli()}
	}

	var result T
	err = s.collection.FindOneAndUpdate(ctx, filter, updateMap, opts).Decode(&result)
	if err != nil {
		return zero, err
	}

	return result, nil
}

// FindOneAndDelete tìm và xóa một document
func (s *BaseServiceImpl[T]) FindOneAndDelete(ctx context.Context, filter interface{}, opts *options.FindOneAndDeleteOptions) (T, error) {
	var zero T

	if filter == nil {
		filter = bson.D{}
	}

	if opts == nil {
		opts = options.FindOneAndDelete()
	}

	var result T
	err := s.collection.FindOneAndDelete(ctx, filter, opts).Decode(&result)
	if err != nil {
		return zero, err
	}

	return result, nil
}

// Distinct lấy danh sách các giá trị duy nhất của một trường
func (s *BaseServiceImpl[T]) Distinct(ctx context.Context, fieldName string, filter interface{}) ([]interface{}, error) {
	if filter == nil {
		filter = bson.D{}
	}

	return s.collection.Distinct(ctx, fieldName, filter)
}
