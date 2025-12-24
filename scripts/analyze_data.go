package main

import (
	"context"
	"fmt"
	"log"
	"meta_commerce/config"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Script này phân tích dữ liệu trong MongoDB
// Bao gồm: đếm số documents, thống kê collections, phân tích cấu trúc dữ liệu
func main() {
	fmt.Println("=== Script Phân Tích Dữ Liệu MongoDB ===\n")

	// Đọc cấu hình từ file env
	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("Không thể đọc cấu hình từ file env")
	}

	// Kết nối với MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	// Phân tích các databases
	databases := []string{
		cfg.MongoDB_DBName_Auth,
		cfg.MongoDB_DBName_Staging,
		cfg.MongoDB_DBName_Data,
	}

	for _, dbName := range databases {
		if dbName == "" {
			continue
		}
		analyzeDatabase(ctx, client, dbName)
		fmt.Println()
	}
}

// analyzeDatabase phân tích một database cụ thể
func analyzeDatabase(ctx context.Context, client *mongo.Client, dbName string) {
	fmt.Printf("📊 PHÂN TÍCH DATABASE: %s\n", dbName)
	fmt.Println(strings.Repeat("=", 60))

	db := client.Database(dbName)

	// Kiểm tra database có tồn tại không
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		fmt.Printf("⚠ Không thể truy cập database %s: %v\n", dbName, err)
		return
	}

	if len(collections) == 0 {
		fmt.Printf("⚠ Database %s không có collection nào\n", dbName)
		return
	}

	// Thống kê tổng quan
	totalDocs := int64(0)
	collectionStats := make(map[string]int64)

	for _, collName := range collections {
		collection := db.Collection(collName)
		count, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Printf("⚠ Không thể đếm documents trong %s: %v\n", collName, err)
			continue
		}
		collectionStats[collName] = count
		totalDocs += count
	}

	fmt.Printf("📈 Tổng số collections: %d\n", len(collections))
	fmt.Printf("📈 Tổng số documents: %d\n", totalDocs)
	fmt.Println()

	// Hiển thị chi tiết từng collection
	fmt.Println("📋 CHI TIẾT CÁC COLLECTIONS:")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-40s %15s\n", "Collection", "Số Documents")
	fmt.Println(strings.Repeat("-", 60))

	// Sắp xếp theo số documents giảm dần
	type collInfo struct {
		name  string
		count int64
	}
	var sortedColls []collInfo
	for name, count := range collectionStats {
		sortedColls = append(sortedColls, collInfo{name, count})
	}

	// Sắp xếp
	for i := 0; i < len(sortedColls)-1; i++ {
		for j := i + 1; j < len(sortedColls); j++ {
			if sortedColls[i].count < sortedColls[j].count {
				sortedColls[i], sortedColls[j] = sortedColls[j], sortedColls[i]
			}
		}
	}

	for _, coll := range sortedColls {
		fmt.Printf("%-40s %15d\n", coll.name, coll.count)
	}
	fmt.Println()

	// Phân tích chi tiết một số collections quan trọng
	importantCollections := []string{
		"customers",
		"pc_pos_orders",
		"pc_pos_products",
		"fb_messages",
		"fb_conversations",
		"auth_users",
	}

	for _, collName := range importantCollections {
		if count, exists := collectionStats[collName]; exists && count > 0 {
			analyzeCollection(ctx, db, collName)
		}
	}
}

// analyzeCollection phân tích chi tiết một collection
func analyzeCollection(ctx context.Context, db *mongo.Database, collName string) {
	collection := db.Collection(collName)

	fmt.Printf("🔍 PHÂN TÍCH CHI TIẾT: %s\n", collName)
	fmt.Println(strings.Repeat("-", 60))

	// Đếm tổng số documents
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Printf("⚠ Không thể đếm documents: %v\n", err)
		return
	}
	fmt.Printf("Tổng số documents: %d\n", totalCount)

	// Lấy một sample document để xem cấu trúc
	var sample bson.M
	err = collection.FindOne(ctx, bson.M{}).Decode(&sample)
	if err == nil {
		fmt.Println("\n📄 Cấu trúc document mẫu:")
		fmt.Println("Các trường chính:")
		for key := range sample {
			if key != "_id" {
				fmt.Printf("  - %s\n", key)
			}
		}
	}

	// Phân tích theo loại collection
	switch collName {
	case "customers":
		analyzeCustomers(ctx, collection)
	case "pc_pos_orders":
		analyzeOrders(ctx, collection)
	case "pc_pos_products":
		analyzeProducts(ctx, collection)
	case "fb_messages":
		analyzeFbMessages(ctx, collection)
	case "fb_conversations":
		analyzeFbConversations(ctx, collection)
	case "auth_users":
		analyzeUsers(ctx, collection)
	}

	fmt.Println()
}

// analyzeCustomers phân tích collection customers
func analyzeCustomers(ctx context.Context, collection *mongo.Collection) {
	// Đếm customers theo source
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$source",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		fmt.Println("\n📊 Phân bố theo source:")
		for cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err == nil {
				source := result["_id"]
				if source == nil {
					source = "null"
				}
				count := result["count"]
				fmt.Printf("  - %v: %v\n", source, count)
			}
		}
	}
}

// analyzeOrders phân tích collection orders
func analyzeOrders(ctx context.Context, collection *mongo.Collection) {
	// Đếm orders theo status
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		fmt.Println("\n📊 Phân bố theo status:")
		for cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err == nil {
				status := result["_id"]
				if status == nil {
					status = "null"
				}
				count := result["count"]
				fmt.Printf("  - %v: %v\n", status, count)
			}
		}
	}

	// Đếm orders theo shop
	pipeline = []bson.M{
		{"$group": bson.M{
			"_id":   "$shopId",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	cursor, err = collection.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		fmt.Println("\n📊 Top 10 shops theo số orders:")
		for cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err == nil {
				shopId := result["_id"]
				if shopId == nil {
					shopId = "null"
				}
				count := result["count"]
				fmt.Printf("  - Shop %v: %v orders\n", shopId, count)
			}
		}
	}
}

// analyzeProducts phân tích collection products
func analyzeProducts(ctx context.Context, collection *mongo.Collection) {
	// Đếm products theo shop
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$shopId",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		fmt.Println("\n📊 Top 10 shops theo số products:")
		for cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err == nil {
				shopId := result["_id"]
				if shopId == nil {
					shopId = "null"
				}
				count := result["count"]
				fmt.Printf("  - Shop %v: %v products\n", shopId, count)
			}
		}
	}
}

// analyzeFbMessages phân tích collection fb_messages
func analyzeFbMessages(ctx context.Context, collection *mongo.Collection) {
	// Đếm messages theo page
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$pageId",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		fmt.Println("\n📊 Top 10 pages theo số messages:")
		for cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err == nil {
				pageId := result["_id"]
				if pageId == nil {
					pageId = "null"
				}
				count := result["count"]
				fmt.Printf("  - Page %v: %v messages\n", pageId, count)
			}
		}
	}
}

// analyzeFbConversations phân tích collection fb_conversations
func analyzeFbConversations(ctx context.Context, collection *mongo.Collection) {
	// Đếm conversations theo page
	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$pageId",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		fmt.Println("\n📊 Top 10 pages theo số conversations:")
		for cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err == nil {
				pageId := result["_id"]
				if pageId == nil {
					pageId = "null"
				}
				count := result["count"]
				fmt.Printf("  - Page %v: %v conversations\n", pageId, count)
			}
		}
	}
}

// analyzeUsers phân tích collection users
func analyzeUsers(ctx context.Context, collection *mongo.Collection) {
	// Đếm users có email
	emailCount, _ := collection.CountDocuments(ctx, bson.M{"email": bson.M{"$exists": true, "$ne": ""}})
	fmt.Printf("\n📊 Users có email: %d\n", emailCount)

	// Đếm users có phone
	phoneCount, _ := collection.CountDocuments(ctx, bson.M{"phone": bson.M{"$exists": true, "$ne": ""}})
	fmt.Printf("📊 Users có phone: %d\n", phoneCount)

	// Đếm users verified
	emailVerified, _ := collection.CountDocuments(ctx, bson.M{"emailVerified": true})
	fmt.Printf("📊 Users đã verify email: %d\n", emailVerified)

	phoneVerified, _ := collection.CountDocuments(ctx, bson.M{"phoneVerified": true})
	fmt.Printf("📊 Users đã verify phone: %d\n", phoneVerified)
}



