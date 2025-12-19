package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meta_commerce/config"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Script này export documents mẫu ra JSON để phân tích cấu trúc
func main() {
	fmt.Println("=== Export Documents Mẫu để Phân Tích ===\n")

	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("Không thể đọc cấu hình từ file env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDB_ConnectionURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Không thể kết nối với MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Không thể ping MongoDB: %v", err)
	}

	fmt.Printf("✓ Đã kết nối với MongoDB\n\n")

	db := client.Database(cfg.MongoDB_DBName_Auth)

	// Tạo thư mục output (từ thư mục gốc dự án)
	currentDir, _ := os.Getwd()
	// Nếu đang ở thư mục api, đi lên 1 cấp
	if filepath.Base(currentDir) == "api" {
		currentDir = filepath.Dir(currentDir)
	}
	outputDir := filepath.Join(currentDir, "docs", "09-ai-context", "sample-data")
	os.MkdirAll(outputDir, 0755)
	fmt.Printf("Output directory: %s\n\n", outputDir)

	collections := []string{
		"customers",
		"pc_pos_orders",
		"pc_pos_products",
		"pc_pos_variations",
		"pc_pos_shops",
		"pc_pos_warehouses",
		"pc_pos_categories",
		"fb_conversations",
		"fb_messages",
		"fb_message_items",
		"fb_posts",
		"fb_pages",
	}

	for _, collName := range collections {
		exportCollection(ctx, db, collName, outputDir)
	}

	fmt.Println("\n✓ Hoàn thành export documents mẫu")
}

func exportCollection(ctx context.Context, db *mongo.Database, collName string, outputDir string) {
	collection := db.Collection(collName)

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Printf("⚠ %s: Không thể đếm documents\n", collName)
		return
	}

	if count == 0 {
		fmt.Printf("⏭ %s: Không có documents\n", collName)
		return
	}

	// Lấy 2 documents mẫu
	var samples []bson.M
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(2))
	if err != nil {
		fmt.Printf("⚠ %s: Không thể lấy documents\n", collName)
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &samples); err != nil {
		fmt.Printf("⚠ %s: Không thể decode documents\n", collName)
		return
	}

	if len(samples) == 0 {
		return
	}

	// Export ra file JSON
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s-sample.json", collName))
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("⚠ %s: Không thể tạo file\n", collName)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(samples); err != nil {
		fmt.Printf("⚠ %s: Không thể encode JSON\n", collName)
		return
	}

	fmt.Printf("✓ %s: %d documents → %s\n", collName, count, outputFile)
}

