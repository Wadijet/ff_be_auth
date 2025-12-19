package main

import (
	"context"
	"fmt"
	"log"
	"meta_commerce/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Script này xóa index pageId_unique khỏi collection fb_posts
// Vì một page có thể có nhiều posts, nên không nên có unique index trên pageId
func main() {
	fmt.Println("=== Script xóa index pageId_unique từ collection fb_posts ===")

	// Đọc cấu hình từ file env
	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("Không thể đọc cấu hình từ file env")
	}

	// Kết nối với MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	fmt.Printf("Đã kết nối với MongoDB: %s\n", cfg.MongoDB_ConnectionURI)
	fmt.Printf("Database: %s\n", cfg.MongoDB_DBName_Auth)
	fmt.Printf("Collection: fb_posts\n")

	// Lấy collection
	db := client.Database(cfg.MongoDB_DBName_Auth)
	collection := db.Collection("fb_posts")

	// Kiểm tra xem index có tồn tại không
	indexes, err := collection.Indexes().List(ctx)
	if err != nil {
		log.Fatalf("Không thể lấy danh sách index: %v", err)
	}

	indexExists := false
	for indexes.Next(ctx) {
		var indexInfo map[string]interface{}
		if err := indexes.Decode(&indexInfo); err != nil {
			continue
		}

		if name, ok := indexInfo["name"].(string); ok && name == "pageId_unique" {
			indexExists = true
			break
		}
	}
	indexes.Close(ctx)

	if !indexExists {
		fmt.Println("✓ Index 'pageId_unique' không tồn tại. Không cần xóa.")
		return
	}

	// Xóa index
	fmt.Println("Đang xóa index 'pageId_unique'...")
	_, err = collection.Indexes().DropOne(ctx, "pageId_unique")
	if err != nil {
		log.Fatalf("Không thể xóa index: %v", err)
	}

	fmt.Println("✓ Đã xóa index 'pageId_unique' thành công!")
	fmt.Println("\nLưu ý: Khi restart server, hàm CreateIndexes sẽ tự động xóa các index không còn được định nghĩa trong model.")
}
