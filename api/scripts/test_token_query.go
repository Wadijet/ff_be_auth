package main

import (
	"context"
	"fmt"
	"log"
	"meta_commerce/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("=== Test Token Query ===\n")

	// Token cần test
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI2OTRmNGFjYjVkZWJjMjEwZWVlMTdhODciLCJ0aW1lIjoiNjk1MGE4N2IiLCJyYW5kb21OdW1iZXIiOiIwIn0.RN1HmsEpHSeMPzaPUnKFNGMppnMdB2mg04aIrwBss7k"

	// Đọc cấu hình
	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("Không thể đọc cấu hình")
	}

	// Kết nối MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDB_ConnectionURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Không thể kết nối MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Không thể ping MongoDB: %v", err)
	}

	fmt.Printf("✓ Đã kết nối MongoDB: %s\n", cfg.MongoDB_ConnectionURI)
	fmt.Printf("✓ Database: %s\n", cfg.MongoDB_DBName_Auth)
	fmt.Println()

	// List tất cả collections
	collections, err := client.Database(cfg.MongoDB_DBName_Auth).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Printf("Không thể list collections: %v", err)
	} else {
		fmt.Printf("Collections trong database: %v\n", collections)
		fmt.Println()
	}

	db := client.Database(cfg.MongoDB_DBName_Auth)
	collection := db.Collection("auth_users")
	
	// Đếm số documents
	count, _ := collection.CountDocuments(ctx, bson.M{})
	fmt.Printf("Tổng số users trong collection 'auth_users': %d\n", count)
	fmt.Println()

	// Query 1: Tìm trong field "token"
	fmt.Println("=== Query 1: Tìm trong field 'token' ===")
	query1 := bson.M{"token": token}
	var result1 bson.M
	err1 := collection.FindOne(ctx, query1).Decode(&result1)
	if err1 != nil {
		fmt.Printf("❌ Không tìm thấy: %v\n", err1)
	} else {
		fmt.Printf("✅ Tìm thấy user!\n")
		if id, ok := result1["_id"]; ok {
			fmt.Printf("   User ID: %v\n", id)
		}
		if email, ok := result1["email"]; ok {
			fmt.Printf("   Email: %v\n", email)
		}
		if tokenField, ok := result1["token"]; ok {
			tokenStr := fmt.Sprintf("%v", tokenField)
			if len(tokenStr) > 50 {
				fmt.Printf("   Token (first 50 chars): %s...\n", tokenStr[:50])
			} else {
				fmt.Printf("   Token: %s\n", tokenStr)
			}
			// So sánh token
			if tokenStr == token {
				fmt.Printf("   ✅ Token khớp hoàn toàn!\n")
			} else {
				fmt.Printf("   ⚠️ Token không khớp!\n")
				fmt.Printf("   Expected length: %d\n", len(token))
				fmt.Printf("   Actual length: %d\n", len(tokenStr))
			}
		}
	}
	fmt.Println()

	// Query 2: Tìm trong array "tokens" với dot notation
	fmt.Println("=== Query 2: Tìm trong array 'tokens' với dot notation ===")
	query2 := bson.M{"tokens.jwtToken": token}
	var result2 bson.M
	err2 := collection.FindOne(ctx, query2).Decode(&result2)
	if err2 != nil {
		fmt.Printf("❌ Không tìm thấy: %v\n", err2)
	} else {
		fmt.Printf("✅ Tìm thấy user trong tokens array!\n")
		if id, ok := result2["_id"]; ok {
			fmt.Printf("   User ID: %v\n", id)
		}
	}
	fmt.Println()

	// Query 3: Tìm với $elemMatch
	fmt.Println("=== Query 3: Tìm với $elemMatch ===")
	query3 := bson.M{
		"tokens": bson.M{
			"$elemMatch": bson.M{
				"jwtToken": token,
			},
		},
	}
	var result3 bson.M
	err3 := collection.FindOne(ctx, query3).Decode(&result3)
	if err3 != nil {
		fmt.Printf("❌ Không tìm thấy: %v\n", err3)
	} else {
		fmt.Printf("✅ Tìm thấy user với $elemMatch!\n")
		if id, ok := result3["_id"]; ok {
			fmt.Printf("   User ID: %v\n", id)
		}
	}
	fmt.Println()

	// Query 4: Tìm user theo userID từ token để xem token thực tế trong DB
	fmt.Println("=== Query 4: Tìm user theo userID để xem token thực tế ===")
	userID := "694f4acb5debc210eee17a87"
	objectID, _ := primitive.ObjectIDFromHex(userID)
	query4 := bson.M{"_id": objectID}
	var result4 bson.M
	err4 := collection.FindOne(ctx, query4).Decode(&result4)
	if err4 != nil {
		fmt.Printf("❌ Không tìm thấy user: %v\n", err4)
	} else {
		fmt.Printf("✅ Tìm thấy user!\n")
		if id, ok := result4["_id"]; ok {
			fmt.Printf("   User ID: %v\n", id)
		}
		if email, ok := result4["email"]; ok {
			fmt.Printf("   Email: %v\n", email)
		}
		if tokenField, ok := result4["token"]; ok {
			tokenStr := fmt.Sprintf("%v", tokenField)
			fmt.Printf("   Token trong DB (length: %d): %s\n", len(tokenStr), tokenStr)
			fmt.Printf("   Token cần tìm (length: %d): %s\n", len(token), token)
			if tokenStr == token {
				fmt.Printf("   ✅ Token khớp hoàn toàn!\n")
			} else {
				fmt.Printf("   ⚠️ Token không khớp!\n")
				// So sánh từng ký tự
				minLen := len(token)
				if len(tokenStr) < minLen {
					minLen = len(tokenStr)
				}
				diffCount := 0
				for i := 0; i < minLen; i++ {
					if token[i] != tokenStr[i] {
						diffCount++
						if diffCount <= 5 {
							fmt.Printf("   Khác nhau tại vị trí %d: expected='%c' (%d), actual='%c' (%d)\n",
								i, token[i], token[i], tokenStr[i], tokenStr[i])
						}
					}
				}
				if diffCount > 5 {
					fmt.Printf("   ... và %d vị trí khác nữa\n", diffCount-5)
				}
			}
		} else {
			fmt.Printf("   ⚠️ Không có field 'token' trong document!\n")
		}
		if tokensArray, ok := result4["tokens"]; ok {
			fmt.Printf("   Tokens array: %v\n", tokensArray)
		}
	}
}
