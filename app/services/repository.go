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

// Repository là cấu trúc chứa thông tin kết nối đến MongoDB
type Repository struct {
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
	case global.MongoDB_ColNames.RolePermission:
		return c.MongoDB_DBNameAuth
	case global.MongoDB_ColNames.UserRole:
		return c.MongoDB_DBNameAuth
	// LOG

	default:
		return ""
	}
}

// Khởi tạo Repository
// trả về interface gắn với Repository
func NewRepository(c *config.Configuration, db *mongo.Client, collection_name string) *Repository {
	dbName := GetDBName(c, collection_name)
	if dbName == "" {
		return nil
	} else {
		return &Repository{mongoClient: db, mongoCollection: db.Database(dbName).Collection(collection_name), config: c}
	}
}

// Cài đặt collection để làm việc
func (service *Repository) SetCollection(collection_name string) (ResultCollection *mongo.Collection) {
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
func (service *Repository) InsertOne(ctx context.Context, model interface{}) (InsertOneResult interface{}, err error) {

	// Thêm createdAt, updatedAt vào dữ liệu đầu vào
	myMap, err := utility.ToMap(model)
	if err != nil {
		return nil, errors.New("Input data is not a map")
	}
	myMap["createdAt"] = utility.CurrentTimeInMilli()
	//myMap["updatedAt"] = utility.CurrentTimeInMilli()
	return service.mongoCollection.InsertOne(ctx, myMap)

}

// InsertMany chèn nhiều tài liệu vào collection
func (service *Repository) InsertMany(ctx context.Context, models []interface{}) (InsertOneResult interface{}, err error) {

	var Maps []interface{}
	for _, model := range models {
		// Thêm createdAt, updatedAt vào dữ liệu đầu vào
		myMap, err := utility.ToMap(model)
		if err != nil {
			return nil, errors.New("Input data is not a map")
		}
		myMap["createdAt"] = utility.CurrentTimeInMilli()
		//myMap["updatedAt"] = utility.CurrentTimeInMilli()

		Maps = append(Maps, myMap)
	}

	return service.mongoCollection.InsertMany(ctx, Maps)
}

// FindOneById tìm một tài liệu theo ID
func (service *Repository) FindOneById(ctx context.Context, id string, opts *options.FindOneOptions) (FindOneResult interface{}, err error) {

	query := bson.D{{"_id", utility.String2ObjectID(id)}}
	var result bson.M
	if opts != nil {
		err = service.mongoCollection.FindOne(ctx, query, opts).Decode(&result)
	} else {
		err = service.mongoCollection.FindOne(ctx, query).Decode(&result)
	}

	if err != nil {
		return nil, err
	}

	return result, nil

}

// FindOne tìm một tài liệu theo query
func (service *Repository) FindOne(ctx context.Context, query interface{}, opts *options.FindOneOptions) (FindOneResult interface{}, err error) {
	var result bson.M
	if opts != nil {
		err = service.mongoCollection.FindOne(ctx, query, opts).Decode(&result)
	} else {
		err = service.mongoCollection.FindOne(ctx, query).Decode(&result)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

// CountAll đếm tất cả tài liệu theo filter
func (service *Repository) CountAll(ctx context.Context, filter interface{}, limit int64) (CountResult interface{}, err error) {

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
func (service *Repository) FindAll(ctx context.Context, filter interface{}, opts *options.FindOptions) (FindResult interface{}, err error) {
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
func (service *Repository) FindAllWithPaginate(ctx context.Context, filter interface{}, opts *options.FindOptions) (FindResult interface{}, err error) {
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

// UpdateOneById cập nhật một tài liệu theo ID
func (service *Repository) UpdateOneById(ctx context.Context, id string, change interface{}) (UpdateResult interface{}, err error) {

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

	// Tạo query
	query := bson.D{{"_id", utility.String2ObjectID(id)}}
	return service.mongoCollection.UpdateOne(ctx, query, myMap)

}

// UpdateMany cập nhật nhiều tài liệu theo query
func (service *Repository) UpdateMany(ctx context.Context, query, change interface{}) (UpdateResult interface{}, err error) {

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
func (service *Repository) DeleteOneById(ctx context.Context, id string) (DeleteResult interface{}, err error) {

	query := bson.D{{"_id", utility.String2ObjectID(id)}}
	result, err := service.mongoCollection.DeleteOne(ctx, query)
	return result, err

}

// DeleteMany xóa nhiều tài liệu theo query
func (service *Repository) DeleteMany(ctx context.Context, query interface{}) (DeleteResult interface{}, err error) {

	result, err := service.mongoCollection.DeleteMany(ctx, query)
	return result, err

}
