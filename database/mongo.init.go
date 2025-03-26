package database

import (
	"context"
	"fmt"
	"log"
	"meta_commerce/global"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureDatabaseAndCollections đảm bảo rằng cơ sở dữ liệu và các collection cần thiết tồn tại.
// Nếu cơ sở dữ liệu không tồn tại, nó sẽ được tạo ra. Nếu các collection không tồn tại, chúng sẽ được tạo ra bằng cách
// chèn một tài liệu dummy và sau đó xóa nó.
//
// Tham số:
// - client: Một đối tượng *mongo.Client kết nối tới MongoDB.
//
// Trả về:
// - error: Lỗi nếu có vấn đề xảy ra trong quá trình kiểm tra hoặc tạo cơ sở dữ liệu và collection.
func EnsureDatabaseAndCollections(client *mongo.Client) error {
	dbName := global.MongoDB_ServerConfig.MongoDB_DBNameAuth

	// Tạo 1 context tổng 30 giây để duyệt tất cả collections
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Kiểm tra database
	dbList, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list databases: %w", err)
	}

	dbExists := false
	for _, name := range dbList {
		if name == dbName {
			dbExists = true
			break
		}
	}
	if !dbExists {
		log.Printf("Database %s does not exist, will create automatically by creating collections", dbName)
	}

	// Tạo database nếu chưa tồn tại
	db := client.Database(dbName)
	collections := []string{}
	v := reflect.ValueOf(global.MongoDB_ColNames)
	for i := 0; i < v.NumField(); i++ {
		collections = append(collections, v.Field(i).String())
	}

	// Kiểm tra và tạo collections
	collList, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	for _, collectionName := range collections {
		// Kiểm tra collection có tồn tại hay không
		exists := false
		for _, existingColl := range collList {
			if existingColl == collectionName {
				exists = true
				break
			}
		}
		// Tạo collection nếu chưa tồn tại
		if !exists {
			log.Printf("Collection %s chưa tồn tại, tạo mới.", collectionName)
			if err := db.CreateCollection(ctx, collectionName); err != nil {
				return fmt.Errorf("failed to create collection %s: %w", collectionName, err)
			}
		}
	}

	log.Printf("Database and collections are ensured in database: %s", dbName)
	return nil
}

// Hàm parseOrder: Trích xuất thứ tự sắp xếp từ tag (1 hoặc -1)
func parseOrder(tag string) int {
	if strings.Contains(tag, "order:-1") {
		return -1 // Nếu tag chứa "order:-1", trả về -1 (giảm dần)
	}
	return 1 // Mặc định trả về 1 (tăng dần)
}

// Hàm parseIndexTag: Phân tách và phân tích tag index
func parseIndexTag(tag string) []map[string]string {
	parts := strings.Split(tag, ";") // Tách tag theo dấu ';'
	result := []map[string]string{}

	for _, part := range parts {
		subParts := strings.Split(part, ",") // Tách từng cấu hình theo dấu ','
		entry := map[string]string{}
		for _, subPart := range subParts {
			kv := strings.Split(subPart, ":") // Tách thành key và value (nếu có)
			if len(kv) == 2 {
				entry[kv[0]] = kv[1]
			} else {
				entry[kv[0]] = ""
			}
		}
		result = append(result, entry)
	}

	return result // Trả về danh sách các cấu hình index
}

func compareIndex(existingIndex bson.M, keys bson.D, options *options.IndexOptions) bool {
	existingKeys, ok := existingIndex["key"].(bson.M)
	if !ok {
		return false
	}

	// So sánh các khóa
	for _, key := range keys {
		existingValue, exists := existingKeys[key.Key]
		if !exists {
			return false
		}

		// Xử lý cho trường hợp 1 / -1
		newVal, isInt := key.Value.(int)
		if isInt {
			// convert existingValue về int (nếu có thể)
			switch ev := existingValue.(type) {
			case int32:
				if int(ev) != newVal {
					return false
				}
			case int64:
				if int(ev) != newVal {
					return false
				}
			case float64:
				if int(ev) != newVal {
					return false
				}
			default:
				return false
			}
		} else {
			// fallback so sánh kiểu cũ
			if existingValue != key.Value {
				return false
			}
		}
	}

	// So sánh các tùy chọn (unique)
	if unique, ok := existingIndex["unique"].(bool); ok && options.Unique != nil {
		if unique != *options.Unique {
			return false
		}
	} else if options.Unique != nil && *options.Unique {
		// index cũ không unique, index mới lại unique => mismatch
		return false
	}

	// So sánh TTL
	if ttl, ok := existingIndex["expireAfterSeconds"].(int32); ok && options.ExpireAfterSeconds != nil {
		if ttl != *options.ExpireAfterSeconds {
			return false
		}
	}

	return true
}

// checkAndReplaceIndex kiểm tra và thay thế index nếu cần thiết
func checkAndReplaceIndex(
	ctx context.Context,
	collection *mongo.Collection,
	existingIndexes map[string]bson.M,
	indexName string,
	keys bson.D,
	options *options.IndexOptions,
) error {
	// Kiểm tra nếu index đã tồn tại
	if existingIndex, exists := existingIndexes[indexName]; exists {
		// So sánh cấu hình index hiện tại với cấu hình mới
		if compareIndex(existingIndex, keys, options) {
			fmt.Printf("Index %s đã tồn tại và đúng cấu hình, bỏ qua...\n", indexName)
			return nil
		}
		// Xóa index nếu cấu hình không khớp
		if _, err := collection.Indexes().DropOne(ctx, indexName); err != nil {
			return fmt.Errorf("không thể xóa index %s: %w", indexName, err)
		}
		fmt.Printf("Đã xóa index cũ: %s\n", indexName)
	}

	// Tạo index mới
	if _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: options,
	}); err != nil {
		return fmt.Errorf("không thể tạo index %s: %w", indexName, err)
	}
	fmt.Printf("Đã tạo index: %s\n", indexName)
	return nil
}

func CreateIndexes(ctx context.Context, collection *mongo.Collection, model interface{}) error {

	fmt.Printf("Bắt đầu xử lý index cho collection: %s\n", collection.Name())

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	fmt.Printf("Lấy danh sách indexs hiện có.\n")
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return fmt.Errorf("không thể lấy danh sách index: %w", err)
	}
	defer cursor.Close(ctx)

	existingIndexes := map[string]bson.M{}
	for cursor.Next(ctx) {
		var indexInfo bson.M
		if err := cursor.Decode(&indexInfo); err != nil {
			return fmt.Errorf("không thể giải mã thông tin index: %w", err)
		}
		if name, ok := indexInfo["name"].(string); ok {
			existingIndexes[name] = indexInfo
		}
	}

	compoundGroups := map[string]bson.D{}
	compoundOptions := map[string]*options.IndexOptions{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag, ok := field.Tag.Lookup("index")
		if !ok {
			continue
		}

		bsonField := field.Tag.Get("bson")
		if bsonField == "" || bsonField == "-" {
			continue
		}

		indexConfigs := parseIndexTag(tag)
		for _, config := range indexConfigs {

			if _, ok := config["text"]; ok {
				keys := bson.D{{Key: bsonField, Value: "text"}} // Định nghĩa kiểu text index
				indexName := bsonField + "_text"
				options := options.Index().SetName(indexName)

				if err := checkAndReplaceIndex(ctx, collection, existingIndexes, indexName, keys, options); err != nil {
					return err
				}
			}

			if _, ok := config["single"]; ok {
				order := parseOrder(tag)
				keys := bson.D{{Key: bsonField, Value: order}}
				indexName := bsonField + "_single"
				options := options.Index().SetName(indexName)

				if err := checkAndReplaceIndex(ctx, collection, existingIndexes, indexName, keys, options); err != nil {
					return err
				}
			}

			if _, ok := config["unique"]; ok {
				keys := bson.D{{Key: bsonField, Value: 1}}
				indexName := bsonField + "_unique"
				options := options.Index().SetName(indexName).SetUnique(true)

				if err := checkAndReplaceIndex(ctx, collection, existingIndexes, indexName, keys, options); err != nil {
					return err
				}
			}

			if ttlValue, ok := config["ttl"]; ok {
				ttl, err := strconv.Atoi(ttlValue)
				if err != nil {
					return fmt.Errorf("TTL không hợp lệ: %w", err)
				}
				keys := bson.D{{Key: bsonField, Value: 1}}
				indexName := bsonField + "_ttl"
				options := options.Index().SetExpireAfterSeconds(int32(ttl)).SetName(indexName)

				if err := checkAndReplaceIndex(ctx, collection, existingIndexes, indexName, keys, options); err != nil {
					return err
				}
			}

			if groupName, ok := config["compound"]; ok {
				order := parseOrder(tag)
				compoundGroups[groupName] = append(compoundGroups[groupName], bson.E{Key: bsonField, Value: order})
				if _, exists := compoundOptions[groupName]; !exists {
					compoundOptions[groupName] = options.Index().SetName(groupName)
				}
			}
		}
	}

	// Tạo compound index
	for groupName, fields := range compoundGroups {
		if err := checkAndReplaceIndex(ctx, collection, existingIndexes, groupName, fields, compoundOptions[groupName]); err != nil {
			return err
		}
	}

	return nil
}
