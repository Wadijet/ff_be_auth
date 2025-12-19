package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meta_commerce/config"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Script này phân tích cấu trúc dữ liệu thực tế trong MongoDB
// Đọc documents mẫu và vẽ lại cấu trúc đầy đủ, sâu xuống các tầng
func main() {
	fmt.Println("=== Script Phân Tích Cấu Trúc Dữ Liệu Thực Tế ===\n")

	// Đọc cấu hình từ file env
	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("Không thể đọc cấu hình từ file env")
	}

	// Kết nối với MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDB_ConnectionURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Không thể kết nối với MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Kiểm tra kết nối
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Không thể ping MongoDB: %v", err)
	}

	fmt.Printf("✓ Đã kết nối với MongoDB: %s\n", cfg.MongoDB_ConnectionURI)
	fmt.Println()

	db := client.Database(cfg.MongoDB_DBName_Auth)

	// Phân tích các collections quan trọng
	collections := []string{
		"customers",
		"pc_pos_orders",
		"pc_pos_products",
		"pc_pos_variations",
		"pc_pos_shops",
		"pc_pos_warehouses",
		"fb_conversations",
		"fb_messages",
		"fb_message_items",
		"fb_posts",
		"fb_pages",
	}

	for _, collName := range collections {
		analyzeCollectionStructure(ctx, db, collName)
		fmt.Println()
	}
}

// analyzeCollectionStructure phân tích cấu trúc thực tế của collection
func analyzeCollectionStructure(ctx context.Context, db *mongo.Database, collName string) {
	collection := db.Collection(collName)

	// Đếm số documents
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Printf("⚠ Không thể đếm documents trong %s: %v\n", collName, err)
		return
	}

	if count == 0 {
		fmt.Printf("📋 %s: Không có documents\n", collName)
		return
	}

	fmt.Printf("📋 PHÂN TÍCH: %s (%d documents)\n", collName, count)
	fmt.Println(strings.Repeat("=", 80))

	// Lấy 3 documents mẫu để phân tích
	var samples []bson.M
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(3))
	if err != nil {
		fmt.Printf("⚠ Không thể lấy documents: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &samples); err != nil {
		fmt.Printf("⚠ Không thể decode documents: %v\n", err)
		return
	}

	if len(samples) == 0 {
		fmt.Printf("⚠ Không có documents để phân tích\n")
		return
	}

	// Phân tích cấu trúc từ documents mẫu
	analyzeDocumentStructure(samples[0], collName, 0)

	// Nếu có nhiều documents, so sánh để tìm fields khác nhau
	if len(samples) > 1 {
		fmt.Println("\n📊 So sánh với documents khác:")
		compareDocuments(samples)
	}
}

// analyzeDocumentStructure phân tích cấu trúc của một document
func analyzeDocumentStructure(doc bson.M, path string, depth int) {
	if depth > 5 { // Giới hạn độ sâu để tránh quá dài
		return
	}

	indent := strings.Repeat("  ", depth)

	for key, value := range doc {
		if key == "_id" {
			continue // Bỏ qua _id
		}

		fullPath := path
		if fullPath != "" {
			fullPath += "."
		}
		fullPath += key

		switch v := value.(type) {
		case map[string]interface{}:
			fmt.Printf("%s%s: (object)\n", indent, key)
			analyzeDocumentStructure(bson.M(v), fullPath, depth+1)

		case bson.M:
			fmt.Printf("%s%s: (object)\n", indent, key)
			analyzeDocumentStructure(v, fullPath, depth+1)

		case []interface{}:
			fmt.Printf("%s%s: (array[%d])\n", indent, key, len(v))
			if len(v) > 0 {
				// Phân tích phần tử đầu tiên
				fmt.Printf("%s  └─ [0]: ", indent)
				switch elem := v[0].(type) {
				case map[string]interface{}:
					fmt.Println("(object)")
					analyzeDocumentStructure(bson.M(elem), fullPath+"[0]", depth+2)
				case bson.M:
					fmt.Println("(object)")
					analyzeDocumentStructure(elem, fullPath+"[0]", depth+2)
				default:
					fmt.Printf("%T\n", elem)
				}
			}

		case primitive.DateTime:
			fmt.Printf("%s%s: (datetime) %v\n", indent, key, time.Unix(int64(v)/1000, 0))

		case time.Time:
			fmt.Printf("%s%s: (datetime) %v\n", indent, key, v)

		default:
			// Hiển thị giá trị mẫu (truncate nếu quá dài)
			valueStr := fmt.Sprintf("%v", v)
			if len(valueStr) > 100 {
				valueStr = valueStr[:100] + "..."
			}
			fmt.Printf("%s%s: (%T) %s\n", indent, key, v, valueStr)
		}
	}
}

// compareDocuments so sánh các documents để tìm fields khác nhau
func compareDocuments(docs []bson.M) {
	if len(docs) < 2 {
		return
	}

	// Lấy tất cả keys từ tất cả documents
	allKeys := make(map[string]bool)
	for _, doc := range docs {
		extractKeys(doc, "", allKeys)
	}

	// Kiểm tra key nào có trong tất cả documents
	commonKeys := make(map[string]bool)
	for key := range allKeys {
		inAll := true
		for _, doc := range docs {
			if !hasKey(doc, key) {
				inAll = false
				break
			}
		}
		if inAll {
			commonKeys[key] = true
		}
	}

	// Kiểm tra key nào chỉ có trong một số documents
	optionalKeys := make(map[string]bool)
	for key := range allKeys {
		if !commonKeys[key] {
			optionalKeys[key] = true
		}
	}

	if len(commonKeys) > 0 {
		fmt.Println("  ✓ Fields luôn có (required):")
		for key := range commonKeys {
			fmt.Printf("    - %s\n", key)
		}
	}

	if len(optionalKeys) > 0 {
		fmt.Println("  ⚠ Fields tùy chọn (optional):")
		for key := range optionalKeys {
			fmt.Printf("    - %s\n", key)
		}
	}
}

// extractKeys trích xuất tất cả keys từ document (bao gồm nested)
func extractKeys(doc bson.M, prefix string, keys map[string]bool) {
	for key, value := range doc {
		if key == "_id" {
			continue
		}

		fullKey := prefix
		if fullKey != "" {
			fullKey += "."
		}
		fullKey += key
		keys[fullKey] = true

		switch v := value.(type) {
		case map[string]interface{}:
			extractKeys(bson.M(v), fullKey, keys)
		case bson.M:
			extractKeys(v, fullKey, keys)
		case []interface{}:
			if len(v) > 0 {
				if elem, ok := v[0].(map[string]interface{}); ok {
					extractKeys(bson.M(elem), fullKey+"[]", keys)
				} else if elem, ok := v[0].(bson.M); ok {
					extractKeys(elem, fullKey+"[]", keys)
				}
			}
		}
	}
}

// hasKey kiểm tra xem document có chứa key (có thể nested) không
func hasKey(doc bson.M, keyPath string) bool {
	parts := strings.Split(keyPath, ".")
	current := doc

	for i, part := range parts {
		// Xử lý array notation
		if strings.HasSuffix(part, "[]") {
			part = strings.TrimSuffix(part, "[]")
		}

		value, exists := current[part]
		if !exists {
			return false
		}

		// Nếu là phần tử cuối, return true
		if i == len(parts)-1 {
			return true
		}

		// Nếu không phải object, không thể tiếp tục
		switch v := value.(type) {
		case map[string]interface{}:
			current = bson.M(v)
		case bson.M:
			current = v
		case []interface{}:
			if len(v) > 0 {
				if elem, ok := v[0].(map[string]interface{}); ok {
					current = bson.M(elem)
				} else if elem, ok := v[0].(bson.M); ok {
					current = elem
				} else {
					return false
				}
			} else {
				return false
			}
		default:
			return false
		}
	}

	return true
}

// Export một document mẫu ra JSON để xem chi tiết
func exportSampleDocument(ctx context.Context, collection *mongo.Collection, collName string) {
	var sample bson.M
	err := collection.FindOne(ctx, bson.M{}).Decode(&sample)
	if err != nil {
		return
	}

	jsonData, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		return
	}

	fmt.Printf("\n📄 Document mẫu (JSON):\n")
	fmt.Println(string(jsonData))
}

