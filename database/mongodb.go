package database

import (
	"atk-go-server/config"
	"atk-go-server/global"
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetInstance initializes and returns a *mongo.Client object.
// This function uses the database connection URL from the provided configuration.
//
// Parameters:
// - c: Pointer to the config.Configuration object containing configuration information.
//
// Returns:
// - *mongo.Client: The connected MongoDB client object.
//
// Notes:
// - This function will log and return an error if there is an issue during connection or connection check.
func GetInstance(c *config.Configuration) (*mongo.Client, error) {
	if c.MongoDB_ConnectionURL == "" {
		return nil, fmt.Errorf("database connection URL is empty")
	}

	clientOptions := options.Client().ApplyURI(c.MongoDB_ConnectionURL).
		SetConnectTimeout(10 * time.Second) // Set a connection timeout

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()

	err = client.Ping(ctxPing, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

// CloseInstance closes the MongoDB client connection.
func CloseInstance(client *mongo.Client) error {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("Failed to disconnect MongoDB client: %v", err)
		return err
	}
	log.Println("Successfully disconnected from MongoDB")
	return nil
}

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
	// Check if the database exists
	dbList, err := client.ListDatabaseNames(context.Background(), map[string]interface{}{})
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
		log.Printf("Database %s does not exist, creating it", dbName)
	}

	db := client.Database(dbName)

	collections := []string{}
	v := reflect.ValueOf(global.MongoDB_ColNames)
	for i := 0; i < v.NumField(); i++ {
		collections = append(collections, v.Field(i).String())
	}

	for _, collectionName := range collections {
		collection := db.Collection(collectionName)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if the collection exists by running a simple command
		err := collection.FindOne(ctx, map[string]interface{}{}).Err()
		if err == mongo.ErrNoDocuments {
			// Collection does not exist, create it by inserting a dummy document and then deleting it
			_, err := collection.InsertOne(ctx, map[string]interface{}{"dummy": "dummy"})
			if err != nil {
				return fmt.Errorf("failed to create collection %s: %w", collectionName, err)
			}
			_, err = collection.DeleteOne(ctx, map[string]interface{}{"dummy": "dummy"})
			if err != nil {
				return fmt.Errorf("failed to delete dummy document from collection %s: %w", collectionName, err)
			}
			log.Printf("Created collection: %s", collectionName)
		} else if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("failed to check collection %s: %w", collectionName, err)
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
		if existingValue, exists := existingKeys[key.Key]; !exists || existingValue != key.Value {
			return false
		}
	}

	// So sánh các tùy chọn (e.g., unique)
	if unique, ok := existingIndex["unique"].(bool); ok && options.Unique != nil {
		if unique != *options.Unique {
			return false
		}
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
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Lấy danh sách index hiện có
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
