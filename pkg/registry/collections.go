package registry

import (
	"meta_commerce/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// Collections là instance toàn cục của Registry cho MongoDB collections
var Collections = NewRegistry[*mongo.Collection]()

// RegisterCollection đăng ký collection mới vào registry

// InitCollections khởi tạo và đăng ký các collections MongoDB
func InitCollections(client *mongo.Client, cfg *config.Configuration) error {
	db := client.Database(cfg.MongoDB_DBNameAuth)
	colNames := []string{"users", "permissions", "roles", "role_permissions", "user_roles",
		"agents", "access_tokens", "fb_pages", "fb_conversations", "fb_messages", "fb_posts", "pc_orders"}

	for _, name := range colNames {
		if err := Collections.Register(name, db.Collection(name)); err != nil {
			return err
		}
	}
	return nil
}
