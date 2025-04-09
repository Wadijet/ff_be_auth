package examples

import (
	"log"
	"meta_commerce/app/etl"
	"meta_commerce/app/etl/dest"
	"meta_commerce/app/etl/scheduler"
	"meta_commerce/app/etl/source"
	"meta_commerce/app/etl/transform"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ExampleScheduledSync demo việc chạy pipeline theo lịch
func ExampleScheduledSync() {
	// 1. Tạo pipeline như ví dụ trước
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

	transformConfig := transform.MapperConfig{
		Mappings: []transform.FieldMapping{
			{
				Source: "id",
				Target: "external_id",
				Type:   "string",
			},
			{
				Source: "name",
				Target: "full_name",
				Type:   "string",
			},
		},
	}

	transformer, err := transform.NewFieldMapper(transformConfig)
	if err != nil {
		log.Fatalf("Lỗi khởi tạo transformer: %v", err)
	}

	destConfig := dest.InternalAPIConfig{
		URL: "http://localhost:8080/api/users/sync",
		Headers: map[string]string{
			"X-API-Key": "internal-system-key",
		},
		Timeout: 30 * time.Second,
	}

	destination, err := dest.NewInternalAPIDestination(destConfig)
	if err != nil {
		log.Fatalf("Lỗi khởi tạo destination: %v", err)
	}

	pipeline, err := etl.NewPipeline(src, transformer, destination)
	if err != nil {
		log.Fatalf("Lỗi khởi tạo pipeline: %v", err)
	}

	// 2. Tạo và cấu hình scheduler
	sch := scheduler.NewScheduler()

	// 3. Thêm các jobs vào scheduler

	// Job 1: Chạy mỗi 5 phút
	job1Config := scheduler.JobConfig{
		Name:     "user_sync_5m",
		Schedule: "*/5 * * * *", // Cron expression: chạy mỗi 5 phút
		Enabled:  true,
	}
	if err := sch.AddJob(job1Config, pipeline); err != nil {
		log.Fatalf("Lỗi thêm job 1: %v", err)
	}

	// Job 2: Chạy vào 00:00 mỗi ngày
	job2Config := scheduler.JobConfig{
		Name:     "user_sync_daily",
		Schedule: "0 0 * * *", // Cron expression: chạy lúc 00:00 mỗi ngày
		Enabled:  true,
	}
	if err := sch.AddJob(job2Config, pipeline); err != nil {
		log.Fatalf("Lỗi thêm job 2: %v", err)
	}

	// 4. Khởi động scheduler
	sch.Start()
	log.Println("Scheduler đã được khởi động với các jobs:")
	for name, config := range sch.GetJobs() {
		log.Printf("- Job %s: %s (enabled: %v)\n", name, config.Schedule, config.Enabled)
	}

	// 5. Đợi signal để dừng gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 6. Dừng scheduler
	sch.Stop()
	log.Println("Scheduler đã dừng gracefully")
}

/*
Cách sử dụng trong main.go:

func main() {
	// Khởi động server
	go server.Start()

	// Khởi động ETL scheduler
	go examples.ExampleScheduledSync()

	// Đợi signal để dừng
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Cleanup và thoát
	log.Println("Shutting down gracefully...")
}
*/
