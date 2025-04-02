// Package services cung cấp các service cơ bản cho việc tương tác với MongoDB
package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	models "meta_commerce/app/models/mongodb"
	"meta_commerce/app/utility"
)

// ====================================
// INTERFACE VÀ STRUCT
// ====================================

// BaseServiceMongo định nghĩa interface chứa các phương thức cơ bản cho việc tương tác với MongoDB
// Type Parameters:
//   - Model: Kiểu dữ liệu của model
type BaseServiceMongo[Model any] interface {
	// NHÓM 1: CÁC HÀM CHUẨN MONGODB DRIVER
	// ====================================

	// 1.1 Thao tác Insert
	InsertOne(ctx context.Context, data Model) (Model, error)
	InsertMany(ctx context.Context, data []Model) ([]Model, error)

	// 1.2 Thao tác Find
	FindOne(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (Model, error)
	Find(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]Model, error)

	// 1.3 Thao tác Update
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *options.UpdateOptions) (Model, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *options.UpdateOptions) (int64, error)

	// 1.4 Thao tác Delete
	DeleteOne(ctx context.Context, filter interface{}) error
	DeleteMany(ctx context.Context, filter interface{}) (int64, error)

	// 1.5 Thao tác Atomic
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts *options.FindOneAndUpdateOptions) (Model, error)
	FindOneAndDelete(ctx context.Context, filter interface{}, opts *options.FindOneAndDeleteOptions) (Model, error)

	// 1.6 Các thao tác khác
	CountDocuments(ctx context.Context, filter interface{}) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter interface{}) ([]interface{}, error)

	// NHÓM 2: CÁC HÀM TIỆN ÍCH MỞ RỘNG
	// ================================

	// 2.1 Các hàm Find mở rộng
	FindOneById(ctx context.Context, id primitive.ObjectID) (Model, error)
	FindManyByIds(ctx context.Context, ids []primitive.ObjectID) ([]Model, error)
	FindWithPagination(ctx context.Context, filter interface{}, page, limit int64) (*models.PaginateResult[Model], error)

	// 2.2 Các hàm Update/Delete mở rộng
	UpdateById(ctx context.Context, id primitive.ObjectID, data Model) (Model, error)
	DeleteById(ctx context.Context, id primitive.ObjectID) error

	// 2.3 Các hàm Upsert tiện ích
	Upsert(ctx context.Context, filter interface{}, data Model) (Model, error)
	UpsertMany(ctx context.Context, filter interface{}, data []Model) ([]Model, error)

	// 2.4 Các hàm kiểm tra
	DocumentExists(ctx context.Context, filter interface{}) (bool, error)
}

// BaseServiceMongoImpl định nghĩa struct triển khai các phương thức cơ bản cho service
// Type Parameters:
//   - Model: Kiểu dữ liệu của model
type BaseServiceMongoImpl[T any] struct {
	collection *mongo.Collection // Collection MongoDB
}

// NewBaseServiceMongo tạo mới một BaseServiceImpl
// Parameters:
//   - collection: Collection MongoDB
//
// Returns:
//   - *BaseServiceImpl[T]: Instance mới của BaseServiceImpl
func NewBaseServiceMongo[T any](collection *mongo.Collection) *BaseServiceMongoImpl[T] {
	return &BaseServiceMongoImpl[T]{
		collection: collection,
	}
}

// ====================================
// NHÓM 1: CÁC HÀM CHUẨN MONGODB DRIVER
// ====================================

// 1.1 Thao tác Insert
// -------------------

// InsertOne tạo mới một bản ghi trong database
func (s *BaseServiceMongoImpl[T]) InsertOne(ctx context.Context, data T) (T, error) {
	var zero T

	// Chuyển data thành map để thêm timestamps
	dataMap, err := utility.ToMap(data)
	if err != nil {
		return zero, utility.ErrInvalidFormat
	}

	// Thêm timestamps
	now := time.Now().UnixMilli()
	dataMap["createdAt"] = now
	dataMap["updatedAt"] = now

	result, err := s.collection.InsertOne(ctx, dataMap)
	if err != nil {
		return zero, utility.ConvertMongoError(err)
	}

	// Lấy lại document vừa tạo
	var created T
	err = s.collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&created)
	if err != nil {
		return zero, utility.ConvertMongoError(err)
	}

	return created, nil
}

// InsertMany tạo nhiều bản ghi trong database
func (s *BaseServiceMongoImpl[T]) InsertMany(ctx context.Context, data []T) ([]T, error) {
	var documents []interface{}
	now := time.Now().UnixMilli()

	for _, item := range data {
		dataMap, err := utility.ToMap(item)
		if err != nil {
			return nil, utility.ErrInvalidFormat
		}
		dataMap["createdAt"] = now
		dataMap["updatedAt"] = now
		documents = append(documents, dataMap)
	}

	result, err := s.collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	// Lấy lại các documents vừa tạo
	var created []T
	filter := bson.M{"_id": bson.M{"$in": result.InsertedIDs}}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	err = cursor.All(ctx, &created)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return created, nil
}

// 1.2 Thao tác Find
// ----------------

// FindOne tìm một document theo điều kiện lọc
func (s *BaseServiceMongoImpl[T]) FindOne(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (T, error) {
	var zero T
	var result T

	if filter == nil {
		filter = bson.D{}
	}

	if opts == nil {
		opts = options.FindOne()
	}

	findResult := s.collection.FindOne(ctx, filter, opts)
	if err := findResult.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return zero, utility.ErrNotFound
		}
		return zero, utility.ConvertMongoError(err)
	}

	if err := findResult.Decode(&result); err != nil {
		return zero, utility.ConvertMongoError(err)
	}

	return result, nil
}

// Find tìm tất cả bản ghi theo điều kiện lọc
func (s *BaseServiceMongoImpl[T]) Find(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]T, error) {
	if filter == nil {
		filter = bson.D{}
	}

	if opts == nil {
		opts = options.Find()
	}

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return results, nil
}

// 1.3 Thao tác Update
// ------------------

// UpdateOne cập nhật một document
func (s *BaseServiceMongoImpl[T]) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts *options.UpdateOptions) (T, error) {
	var zero T

	if filter == nil {
		filter = bson.D{}
	}

	if opts == nil {
		opts = options.Update()
	}

	// Thêm updatedAt vào update
	updateMap, err := utility.ToMap(update)
	if err != nil {
		return zero, utility.ErrInvalidFormat
	}

	if setMap, ok := updateMap["$set"].(map[string]interface{}); ok {
		setMap["updatedAt"] = time.Now().UnixMilli()
		updateMap["$set"] = setMap
	} else {
		updateMap["$set"] = bson.M{"updatedAt": time.Now().UnixMilli()}
	}

	result, err := s.collection.UpdateOne(ctx, filter, updateMap, opts)
	if err != nil {
		return zero, utility.ConvertMongoError(err)
	}

	if result.ModifiedCount == 0 {
		return zero, utility.ErrNotFound
	}

	// Lấy lại document đã update
	var updated T
	err = s.collection.FindOne(ctx, filter).Decode(&updated)
	if err != nil {
		return zero, utility.ConvertMongoError(err)
	}

	return updated, nil
}

// UpdateMany cập nhật nhiều document
func (s *BaseServiceMongoImpl[T]) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts *options.UpdateOptions) (int64, error) {
	if filter == nil {
		filter = bson.D{}
	}

	if opts == nil {
		opts = options.Update()
	}

	// Thêm updatedAt vào update
	updateMap, err := utility.ToMap(update)
	if err != nil {
		return 0, utility.ErrInvalidFormat
	}

	if setMap, ok := updateMap["$set"].(map[string]interface{}); ok {
		setMap["updatedAt"] = time.Now().UnixMilli()
		updateMap["$set"] = setMap
	} else {
		updateMap["$set"] = bson.M{"updatedAt": time.Now().UnixMilli()}
	}

	result, err := s.collection.UpdateMany(ctx, filter, updateMap, opts)
	if err != nil {
		return 0, utility.ConvertMongoError(err)
	}

	return result.ModifiedCount, nil
}

// 1.4 Thao tác Delete
// ------------------

// DeleteOne xóa một document
func (s *BaseServiceMongoImpl[T]) DeleteOne(ctx context.Context, filter interface{}) error {
	if filter == nil {
		filter = bson.D{}
	}

	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return utility.ConvertMongoError(err)
	}

	if result.DeletedCount == 0 {
		return utility.ErrNotFound
	}

	return nil
}

// DeleteMany xóa nhiều document
func (s *BaseServiceMongoImpl[T]) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	if filter == nil {
		filter = bson.D{}
	}

	result, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, utility.ConvertMongoError(err)
	}

	return result.DeletedCount, nil
}

// 1.5 Thao tác Atomic
// ------------------

// FindOneAndUpdate tìm và cập nhật một document
func (s *BaseServiceMongoImpl[T]) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts *options.FindOneAndUpdateOptions) (T, error) {
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
		return zero, utility.ErrInvalidFormat
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
		return zero, utility.ConvertMongoError(err)
	}

	return result, nil
}

// FindOneAndDelete tìm và xóa một document
func (s *BaseServiceMongoImpl[T]) FindOneAndDelete(ctx context.Context, filter interface{}, opts *options.FindOneAndDeleteOptions) (T, error) {
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
		return zero, utility.ConvertMongoError(err)
	}

	return result, nil
}

// 1.6 Các thao tác khác
// --------------------

// CountDocuments đếm số lượng document
func (s *BaseServiceMongoImpl[T]) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	if filter == nil {
		filter = bson.D{}
	}

	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, utility.ConvertMongoError(err)
	}

	return count, nil
}

// Distinct lấy danh sách các giá trị duy nhất
func (s *BaseServiceMongoImpl[T]) Distinct(ctx context.Context, fieldName string, filter interface{}) ([]interface{}, error) {
	if filter == nil {
		filter = bson.D{}
	}

	values, err := s.collection.Distinct(ctx, fieldName, filter)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return values, nil
}

// ====================================
// NHÓM 2: CÁC HÀM TIỆN ÍCH MỞ RỘNG
// ====================================

// 2.1 Các hàm Find mở rộng
// -----------------------

// FindOneById tìm một document theo ObjectId
func (s *BaseServiceMongoImpl[T]) FindOneById(ctx context.Context, id primitive.ObjectID) (T, error) {
	var zero T
	filter := bson.M{"_id": id}
	err := s.collection.FindOne(ctx, filter).Decode(&zero)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return zero, utility.ErrNotFound
		}
		return zero, utility.ConvertMongoError(err)
	}
	return zero, nil
}

// FindManyByIds tìm nhiều document theo danh sách ID
func (s *BaseServiceMongoImpl[T]) FindManyByIds(ctx context.Context, ids []primitive.ObjectID) ([]T, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, utility.ConvertMongoError(err)
	}

	return results, nil
}

// FindWithPagination tìm tất cả bản ghi với phân trang
func (s *BaseServiceMongoImpl[T]) FindWithPagination(ctx context.Context, filter interface{}, page, limit int64) (*models.PaginateResult[T], error) {
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
		return nil, utility.ConvertMongoError(err)
	}

	// Lấy dữ liệu theo trang
	var items []T
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, utility.ConvertMongoError(err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &items); err != nil {
		return nil, utility.ConvertMongoError(err)
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

// 2.2 Các hàm Update/Delete mở rộng
// --------------------------------

// UpdateById cập nhật một document theo ObjectId
// Parameters:
//   - ctx: Context cho việc hủy bỏ hoặc timeout
//   - id: ObjectId của document cần cập nhật
//   - data: Dữ liệu cần cập nhật
//
// Returns:
//   - T: Document đã được cập nhật
//   - error: Lỗi nếu có
func (s *BaseServiceMongoImpl[T]) UpdateById(ctx context.Context, id primitive.ObjectID, data T) (T, error) {
	var zero T
	filter := bson.M{"_id": id}

	// Chuyển data thành map để thêm timestamps
	dataMap, err := utility.ToMap(data)
	if err != nil {
		return zero, utility.ErrInvalidFormat
	}

	// Thêm timestamps
	now := time.Now().UnixMilli()
	dataMap["updatedAt"] = now

	// Tạo options cho update
	opts := options.Update().SetUpsert(false)

	// Thực hiện update
	result, err := s.collection.UpdateOne(ctx, filter, bson.M{"$set": dataMap}, opts)
	if err != nil {
		return zero, utility.ConvertMongoError(err)
	}

	if result.ModifiedCount == 0 {
		return zero, utility.ErrNotFound
	}

	// Lấy lại document đã update
	var updated T
	err = s.collection.FindOne(ctx, filter).Decode(&updated)
	if err != nil {
		return zero, utility.ConvertMongoError(err)
	}

	return updated, nil
}

// DeleteById xóa một document theo ObjectId
// Parameters:
//   - ctx: Context cho việc hủy bỏ hoặc timeout
//   - id: ObjectId của document cần xóa
//
// Returns:
//   - error: Lỗi nếu có
func (s *BaseServiceMongoImpl[T]) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return utility.ConvertMongoError(err)
	}

	if result.DeletedCount == 0 {
		return utility.ErrNotFound
	}

	return nil
}

// 2.3 Các hàm Upsert tiện ích
// --------------------------

// Upsert thực hiện thao tác update nếu tồn tại, insert nếu chưa tồn tại
func (s *BaseServiceMongoImpl[T]) Upsert(ctx context.Context, filter interface{}, data T) (T, error) {
	var zero T

	// Chuyển data thành map để thêm timestamps
	dataMap, err := utility.ToMap(data)
	if err != nil {
		return zero, utility.ErrInvalidFormat
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
		return zero, utility.ConvertMongoError(err)
	}

	return upserted, nil
}

// UpsertMany thực hiện thao tác upsert cho nhiều document
func (s *BaseServiceMongoImpl[T]) UpsertMany(ctx context.Context, filter interface{}, data []T) ([]T, error) {
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
			return nil, utility.ErrInvalidFormat
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
		return nil, utility.ConvertMongoError(err)
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
				return nil, utility.ConvertMongoError(err)
			}
			defer cursor.Close(ctx)

			if err = cursor.All(ctx, &upserted); err != nil {
				return nil, utility.ConvertMongoError(err)
			}
		}
	}

	// Lấy các documents đã được update
	if result.ModifiedCount > 0 {
		cursor, err := s.collection.Find(ctx, filter)
		if err != nil {
			return nil, utility.ConvertMongoError(err)
		}
		defer cursor.Close(ctx)

		var updated []T
		if err = cursor.All(ctx, &updated); err != nil {
			return nil, utility.ConvertMongoError(err)
		}

		// Kết hợp cả documents mới và documents đã update
		upserted = append(upserted, updated...)
	}

	return upserted, nil
}

// 2.4 Các hàm kiểm tra
// -------------------

// DocumentExists kiểm tra xem một document có tồn tại không
func (s *BaseServiceMongoImpl[T]) DocumentExists(ctx context.Context, filter interface{}) (bool, error) {
	if filter == nil {
		filter = bson.D{}
	}

	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, utility.ConvertMongoError(err)
	}

	return count > 0, nil
}
