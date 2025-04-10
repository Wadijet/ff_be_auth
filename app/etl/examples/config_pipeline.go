package examples

import (
	"log"
	"meta_commerce/app/etl"
	"meta_commerce/app/etl/scheduler"
	"os"
	"os/signal"
	"syscall"
)

// ExampleConfigPipeline demo việc chạy pipeline từ code
func ExampleConfigPipeline() {
	// 1. Khởi tạo pipeline builder
	builder := etl.NewPipelineBuilder()

	// 2. Build tất cả pipelines
	pipelines, err := builder.BuildAllPipelines()
	if err != nil {
		log.Fatalf("Lỗi build pipelines: %v", err)
	}

	// 3. Khởi tạo scheduler
	sch := scheduler.NewScheduler()

	// 4. Thêm các jobs cho từng pipeline
	jobConfigs := map[string][]scheduler.JobConfig{
		"user_sync": {
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
		},
		"order_sync": {
			{
				Name:     "order_sync_15m",
				Schedule: "*/15 * * * *",
				Enabled:  true,
			},
		},
		"product_sync": {
			{
				Name:     "product_sync_hourly",
				Schedule: "0 * * * *",
				Enabled:  true,
			},
		},
	}

	// Thêm các jobs vào scheduler
	for i, pipeline := range pipelines {
		var pipelineType string
		switch i {
		case 0:
			pipelineType = "user_sync"
		case 1:
			pipelineType = "order_sync"
		case 2:
			pipelineType = "product_sync"
		}

		configs := jobConfigs[pipelineType]
		for _, config := range configs {
			if err := sch.AddJob(config, pipeline); err != nil {
				log.Fatalf("Lỗi thêm job %s: %v", config.Name, err)
			}
		}
	}

	// 5. Khởi động scheduler
	sch.Start()
	log.Println("Pipeline scheduler đã được khởi động")
	log.Println("Danh sách các jobs đang chạy:")
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
- Build các pipelines từ code
- Khởi tạo các components
- Chạy theo lịch đã cấu hình
- Xử lý lỗi và retry
- Log thông tin chi tiết
*/
