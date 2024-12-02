package database

import (
	"atk-go-server/config"
	"atk-go-server/global"
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

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

// EnsureDatabaseAndCollections checks if the database and collections exist, and creates them if they don't.
//
// Parameters:
// - client: The MongoDB client object.
//
// Returns:
// - error: An error object if there is an issue during the process.
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
