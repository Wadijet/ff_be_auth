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
	colNames := []string{"auth_users", "auth_permissions", "auth_roles", "auth_role_permissions", "auth_user_roles",
		"agents", "access_tokens", "fb_pages", "fb_conversations", "fb_messages", "fb_posts", "pc_orders"}

	for _, name := range colNames {
		registered, err := global.RegistryCollections.Register(name, db.Collection(name))
		if err != nil {
			logrus.Errorf("Failed to register collection %s: %v", name, err)
			return err
		}

		if registered {
			logrus.Infof("Collection %s registered successfully", name)
		} else {
			logrus.Errorf("Collection %s already registered", name)
		}

	}

	return nil
}
