package entity

// MetadataDatabase là struct cho metadata của database
type MetadataDatabase struct {
	Name           string `json:"name" bson:"name"`                       // Tên database
	Description    string `json:"description" bson:"description"`         // Mô tả database
	URI            string `json:"uri" bson:"uri"`                         // URI kết nối database
	MaxPoolSize    int    `json:"max_pool_size" bson:"max_pool_size"`     // Số lượng connection tối đa
	MinPoolSize    int    `json:"min_pool_size" bson:"min_pool_size"`     // Số lượng connection tối thiểu
	ConnectTimeout int    `json:"connect_timeout" bson:"connect_timeout"` // Timeout khi kết nối
	SocketTimeout  int    `json:"socket_timeout" bson:"socket_timeout"`   // Timeout khi gửi nhận dữ liệu
}

// MetadataCollection là struct cho metadata của collection
type MetadataCollection struct {
	Name        string          `json:"name" bson:"name"`               // Tên collection
	Description string          `json:"description" bson:"description"` // Mô tả collection
	Database    string          `json:"database" bson:"database"`       // Tên database
	Fields      []MetadataField `json:"fields" bson:"fields"`           // Các field của collection
	Indexes     []MetadataIndex `json:"indexes" bson:"indexes"`         // Các index của collection
}

// MetadataField là struct cho metadata của field
type MetadataField struct {
	Name        string `json:"name" bson:"name"`               // Tên field
	Description string `json:"description" bson:"description"` // Mô tả field
	Type        string `json:"type" bson:"type"`               // Kiểu dữ liệu field
}

// MetadataIndex là struct cho metadata của index
type MetadataIndex struct {
	Name        string `json:"name" bson:"name"`               // Tên index
	Description string `json:"description" bson:"description"` // Mô tả index
	Type        string `json:"type" bson:"type"`               // Kiểu index
}
