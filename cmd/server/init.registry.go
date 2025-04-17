package main

import (
	"meta_commerce/config"
	"meta_commerce/core/global"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitRegistry() {

	logrus.Info("Initialized registry") // Ghi log thông báo đã khởi tạo registry

	// Khởi tạo registry và đăng ký các collections
	err := InitCollections(global.MongoDB_Session, global.MongoDB_ServerConfig)
	if err != nil {
		logrus.Fatalf("Failed to initialize collections: %v", err)
	}
	logrus.Info("Initialized collection registry")
}

// InitCollections khởi tạo và đăng ký các collections MongoDB
func InitCollections(client *mongo.Client, cfg *config.Configuration) error {
	db := client.Database(cfg.MongoDB_DBName_Auth)
	colNames := []string{"users", "permissions", "roles", "role_permissions", "user_roles",
		"agents", "access_tokens", "fb_pages", "fb_conversations", "fb_messages", "fb_posts", "pc_orders"}

	for _, name := range colNames {
		if err := global.RegistryCollections.Register(name, db.Collection(name)); err != nil {
			return err
		}
	}
	return nil
}
