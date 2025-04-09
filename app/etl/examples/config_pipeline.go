package examples

import (
	"log"
	"meta_commerce/app/etl"
	"meta_commerce/app/etl/scheduler"
	"os"
	"os/signal"
	"syscall"
)

// ExampleConfigPipeline demo việc chạy pipeline từ config
func ExampleConfigPipeline() {
	// 1. Khởi tạo pipeline loader
	loader := etl.NewPipelineLoader("app/etl/pipeline_config.yaml")

	// 2. Load pipeline từ config
	pipeline, err := loader.LoadPipeline("user_sync")
	if err != nil {
		log.Fatalf("Lỗi load pipeline: %v", err)
	}

	// 3. Khởi tạo scheduler
	sch := scheduler.NewScheduler()

	// 4. Thêm các jobs từ config
	jobConfigs := []scheduler.JobConfig{
		{
			Name:     "user_sync_5m",
			Schedule: "*/5 * * * *",
			Enabled:  true,
		},
		{
			Name:     "user_sync_daily",
			Schedule: "0 0 * * *",
			Enabled:  true,
		},
	}

	// Thêm các jobs vào scheduler
	for _, config := range jobConfigs {
		if err := sch.AddJob(config, pipeline); err != nil {
			log.Fatalf("Lỗi thêm job %s: %v", config.Name, err)
		}
	}

	// 5. Khởi động scheduler
	sch.Start()
	log.Println("Pipeline scheduler đã được khởi động")
	log.Printf("Đang chạy pipeline '%s' với các jobs:\n", "user_sync")
	for name, config := range sch.GetJobs() {
		log.Printf("- Job %s: %s (enabled: %v)\n", name, config.Schedule, config.Enabled)
	}

	// 6. Đợi signal để dừng gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 7. Dừng scheduler
	sch.Stop()
	log.Println("Pipeline scheduler đã dừng gracefully")
}

/*
Cách sử dụng:

1. Set environment variables:
export ENV_EXTERNAL_API_TOKEN=your_external_token
export ENV_INTERNAL_API_KEY=your_internal_key

2. Chạy example:
go run main.go

3. Pipeline sẽ:
- Load cấu hình từ file YAML
- Tự động xử lý các biến môi trường
- Validate cấu hình
- Khởi tạo các components
- Chạy theo lịch đã cấu hình
- Xử lý lỗi và retry
- Log thông tin chi tiết
*/
