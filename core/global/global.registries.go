package global

import (
	"meta_commerce/core/registry"

	"go.mongodb.org/mongo-driver/mongo"
)

// Các Registry
var RegistryCollections = registry.NewRegistry[*mongo.Collection]() // Registry cho các collection
var RegistryDatabase = registry.NewRegistry[*mongo.Database]()      // Registry cho các database
