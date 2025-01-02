package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
)

// RepositoryService là cấu trúc chứa thông tin kết nối đến MongoDB
type RepositoryService struct {
	mongoClient     *mongo.Client
	mongoCollection *mongo.Collection
	config          *config.Configuration
}

// GetDBName trả về tên cơ sở dữ liệu dựa trên tên collection
func GetDBName(c *config.Configuration, collectionName string) string {
	switch collectionName {
	// AUTH
	case global.MongoDB_ColNames.Users:
		return c.MongoDB_DBNameAuth
	case global.MongoDB_ColNames.Permissions:
		return c.MongoDB_DBNameAuth
	case global.MongoDB_ColNames.Roles:
		return c.MongoDB_DBNameAuth
	case global.MongoDB_ColNames.RolePermissions:
		return c.MongoDB_DBNameAuth
	case global.MongoDB_ColNames.UserRoles:
		return c.MongoDB_DBNameAuth
	case global.MongoDB_ColNames.Agents:
		return c.MongoDB_DBNameAuth
	// LOG

	default:
		return ""
	}
}

// Khởi tạo Repository
// trả về interface gắn với Repository
func NewRepository(c *config.Configuration, db *mongo.Client, collection_name string) *RepositoryService {
	dbName := GetDBName(c, collection_name)
	if dbName == "" {
		return nil
	} else {
		return &RepositoryService{mongoClient: db, mongoCollection: db.Database(dbName).Collection(collection_name), config: c}
	}
}

// Cài đặt collection để làm việc
func (service *RepositoryService) SetCollection(collection_name string) (ResultCollection *mongo.Collection) {
	dbName := GetDBName(service.config, collection_name)
	if dbName == "" {
		return nil
	} else {
		return service.mongoClient.Database(dbName).Collection(collection_name)
	}
}

// InsertOne chèn một tài liệu vào collection
// Params:	collection name (string)
// return: 	*mongo.Collection
func (service *RepositoryService) InsertOne(ctx context.Context, model interface{}) (InsertOneResult *mongo.InsertOneResult, err error) {

	// Thêm createdAt, updatedAt vào dữ liệu đầu vào
	myMap, err := utility.ToMap(model)
	if err != nil {
		return nil, errors.New("Input data is not a map")
	}
	myMap["createdAt"] = utility.CurrentTimeInMilli()
	myMap["updatedAt"] = utility.CurrentTimeInMilli()
	return service.mongoCollection.InsertOne(ctx, myMap)

}

// InsertMany chèn nhiều tài liệu vào collection
func (service *RepositoryService) InsertMany(ctx context.Context, models []interface{}) (InsertManyResult *mongo.InsertManyResult, err error) {

	var Maps []interface{}
	for _, model := range models {
		// Thêm createdAt, updatedAt vào dữ liệu đầu vào
		myMap, err := utility.ToMap(model)
		if err != nil {
			return nil, errors.New("Input data is not a map")
		}
		myMap["createdAt"] = utility.CurrentTimeInMilli()
		myMap["updatedAt"] = utility.CurrentTimeInMilli()

		Maps = append(Maps, myMap)
	}

	return service.mongoCollection.InsertMany(ctx, Maps)
}

// UpdateOneById cập nhật một tài liệu theo ID
func (service *RepositoryService) UpdateOneById(ctx context.Context, id string, change interface{}) (UpdateResult *mongo.UpdateResult, err error) {

	// Chuyển đổi dữ liệu thành map
	myMap, err := utility.ToMap(change)
	if err != nil {
		return nil, errors.New("Input data is not a map")
	}

	// Đảm bảo các trường trong `$set` không bị bỏ qua
	myChange, ok := myMap["$set"].(map[string]interface{})
	if !ok {
		myChange = make(map[string]interface{})
	}

	// Thêm `updatedAt` vào `$set`
	myChange["updatedAt"] = utility.CurrentTimeInMilli()

	// Đảm bảo các trường mảng rỗng không bị loại bỏ
	for key, value := range myChange {
		if array, ok := value.([]interface{}); ok && len(array) == 0 {
			myChange[key] = []interface{}{} // Giữ mảng rỗng
		}
	}

	myMap["$set"] = myChange

	// Tạo query
	query := bson.D{{Key: "_id", Value: utility.String2ObjectID(id)}}
	return service.mongoCollection.UpdateOne(ctx, query, myMap)

}

// UpdateMany cập nhật nhiều tài liệu theo query
func (service *RepositoryService) UpdateMany(ctx context.Context, query, change interface{}) (UpdateResult *mongo.UpdateResult, err error) {

	// Thêm updatedAt vào map thay đổi
	myMap, err := utility.ToMap(change)
	if err != nil {
		return nil, errors.New("Input data is not a map")
	}
	myChange, err := utility.ToMap(myMap["$set"])
	if err != nil {
		return nil, errors.New("Input data is not a map")
	}
	myChange["updatedAt"] = utility.CurrentTimeInMilli()
	myMap["$set"] = myChange

	return service.mongoCollection.UpdateMany(ctx, query, myMap)

}

// DeleteOneById xóa một tài liệu theo ID
func (service *RepositoryService) DeleteOneById(ctx context.Context, id string) (DeleteResult *mongo.DeleteResult, err error) {

	query := bson.D{{Key: "_id", Value: utility.String2ObjectID(id)}}
	result, err := service.mongoCollection.DeleteOne(ctx, query)
	return result, err

}

// DeleteMany xóa nhiều tài liệu theo query
func (service *RepositoryService) DeleteMany(ctx context.Context, query interface{}) (DeleteResult *mongo.DeleteResult, err error) {

	result, err := service.mongoCollection.DeleteMany(ctx, query)
	return result, err

}

// FindOneById tìm một tài liệu theo ID
func (service *RepositoryService) FindOneById(ctx context.Context, id string, opts *options.FindOneOptions) (FindOneResult bson.M, err error) {

	query := bson.D{{Key: "_id", Value: utility.String2ObjectID(id)}}
	var resultDecoded bson.M
	if opts != nil {
		result := service.mongoCollection.FindOne(ctx, query, opts)
		err = result.Decode(&resultDecoded)
	} else {
		result := service.mongoCollection.FindOne(ctx, query)
		err = result.Decode(&resultDecoded)
	}

	if err != nil {
		return nil, err
	}

	return resultDecoded, nil

}

// FindOne tìm một tài liệu theo query
func (service *RepositoryService) FindOne(ctx context.Context, query interface{}, opts *options.FindOneOptions) (FindOneResult bson.M, err error) {
	var resultDecoded bson.M
	if opts != nil {
		result := service.mongoCollection.FindOne(ctx, query, opts)
		err = result.Decode(&resultDecoded)
	} else {
		result := service.mongoCollection.FindOne(ctx, query)
		err = result.Decode(&resultDecoded)
	}

	if err != nil {
		return nil, err
	}

	return resultDecoded, nil
}

// CountAll đếm tất cả tài liệu theo filter
func (service *RepositoryService) CountAll(ctx context.Context, filter interface{}, limit int64) (CountResult *models.CountResult, err error) {

	count, err := service.mongoCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	countResult := new(models.CountResult)
	countResult.TotalCount = count
	countResult.Limit = limit
	countResult.TotalPage = count / limit

	return countResult, nil
}

// FindAll tìm tất cả tài liệu theo filter
func (service *RepositoryService) FindAll(ctx context.Context, filter interface{}, opts *options.FindOptions) (FindResult []bson.M, err error) {
	if filter == nil {
		filter = bson.D{}
	}

	var cursor *mongo.Cursor
	if opts != nil {
		cursor, err = service.mongoCollection.Find(ctx, filter, opts)
	} else {
		cursor, err = service.mongoCollection.Find(ctx, filter)
	}
	if err != nil {
		return nil, err
	}

	// lấy danh sách tất cả tài liệu trả về và in ra
	var items []bson.M
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, err
}

// FindAllWithPaginate tìm tất cả tài liệu với phân trang
func (service *RepositoryService) FindAllWithPaginate(ctx context.Context, filter interface{}, opts *options.FindOptions) (FindResult *models.PaginateResult, err error) {

	if filter == nil {
		filter = bson.D{}
	}

	cursor, err := service.mongoCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	// lấy danh sách tất cả tài liệu trả về và in ra
	var items []bson.M
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	var result = new(models.PaginateResult)
	result.ItemCount = int64(len(items))
	result.Limit = *opts.Limit
	result.Page = *opts.Skip / *opts.Limit
	result.Items = items

	return result, err

}
