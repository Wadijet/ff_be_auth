package examples

import (
	"context"
	"log"
	"meta_commerce/app/etl"
	"meta_commerce/app/etl/dest"
	"meta_commerce/app/etl/source"
	"meta_commerce/app/etl/transform"
	"time"
)

// ExampleUserSync demo việc đồng bộ user từ external API vào internal system
func ExampleUserSync() {
	// 1. Cấu hình source - lấy dữ liệu từ external API
	sourceConfig := source.RestConfig{
		URL: "https://api.example.com/users",
		Headers: map[string]string{
			"Authorization": "Bearer external-api-token",
		},
		Timeout: 30 * time.Second,
	}

	src, err := source.NewRestApiSource(sourceConfig)
	if err != nil {
		log.Fatalf("Lỗi khởi tạo source: %v", err)
	}

	// 2. Cấu hình transformer - map các trường dữ liệu
	transformConfig := transform.MapperConfig{
		Mappings: []transform.FieldMapping{
			{
				Source: "id",          // id từ external API
				Target: "external_id", // lưu dưới dạng external_id
				Type:   "string",
			},
			{
				Source: "name",      // name từ external API
				Target: "full_name", // lưu dưới dạng full_name
				Type:   "string",
			},
			{
				Source: "email",
				Target: "email",
				Type:   "string",
			},
			{
				Source: "age", // age có thể là string từ API
				Target: "age", // chuyển thành number
				Type:   "number",
			},
			{
				Source: "is_active",
				Target: "active",
				Type:   "boolean",
			},
		},
	}

	transformer, err := transform.NewFieldMapper(transformConfig)
	if err != nil {
		log.Fatalf("Lỗi khởi tạo transformer: %v", err)
	}

	// 3. Cấu hình destination - lưu vào internal API
	destConfig := dest.InternalAPIConfig{
		URL: "http://localhost:8080/api/users/sync",
		Headers: map[string]string{
			"X-API-Key":    "internal-system-key",
			"Content-Type": "application/json",
		},
		Timeout: 30 * time.Second,
	}

	destination, err := dest.NewInternalAPIDestination(destConfig)
	if err != nil {
		log.Fatalf("Lỗi khởi tạo destination: %v", err)
	}

	// 4. Tạo pipeline
	pipeline, err := etl.NewPipeline(src, transformer, destination)
	if err != nil {
		log.Fatalf("Lỗi khởi tạo pipeline: %v", err)
	}

	// 5. Thực thi pipeline
	ctx := context.Background()
	if err := pipeline.Execute(ctx); err != nil {
		log.Fatalf("Lỗi thực thi pipeline: %v", err)
	}

	// 6. In thông tin về các components đã sử dụng
	log.Printf("Thông tin Pipeline:\n%+v\n", pipeline.GetComponents())
}

/*
Ví dụ dữ liệu từ external API:
{
    "id": "ext123",
    "name": "John Doe",
    "email": "john@example.com",
    "age": "30",
    "is_active": true
}

Dữ liệu sau khi transform và gửi đến internal API:
{
    "external_id": "ext123",
    "full_name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "active": true
}
*/
