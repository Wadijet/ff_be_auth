package database

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connectToMongoDB kết nối đến cơ sở dữ liệu MongoDB và trả về client
func ConnectToMongoDB(uri string) (*mongo.Client, error) {

	if uri == "" {
		return nil, fmt.Errorf("URI kết nối cơ sở dữ liệu MongoDB trống")
	}

	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(50).                 // Giới hạn tối đa 50 connections
		SetMinPoolSize(10).                 // Giữ tối thiểu 10 connections trong pool
		SetConnectTimeout(5 * time.Second). // Timeout khi kết nối
		SetSocketTimeout(10 * time.Second)  // Timeout khi gửi nhận dữ liệu

	// Kết nối thử với MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tạo client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Kiểm tra kết nối
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()

	err = client.Ping(ctxPing, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

// EnsureDatabaseExists kiểm tra xem database đã tồn tại chưa, nếu chưa thì tạo mới
func EnsureDatabaseExists(client *mongo.Client, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databases, err := client.ListDatabaseNames(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list databases: %w", err)
	}

	for _, db := range databases {
		if db == dbName {
			log.Printf("Database %s already exists", dbName)
			return nil
		}
	}

	// Tạo mới database
	err = client.Database(dbName).CreateCollection(ctx, "dummyCollection")
	if err != nil {
		return fmt.Errorf("failed to create database %s: %w", dbName, err)
	}

	log.Printf("Database %s created successfully", dbName)
	return nil
}

// CloseInstance closes the MongoDB client connection.
func DisconnectFromMongoDB(client *mongo.Client) error {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("Failed to disconnect MongoDB client: %v", err)
		return err
	}
	log.Println("Successfully disconnected from MongoDB")
	return nil
}
