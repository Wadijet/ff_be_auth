package scheduler

import (
	"context"
	"fmt"
	"log"
	"meta_commerce/app/etl"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// JobConfig chứa cấu hình cho một pipeline job
type JobConfig struct {
	Name     string `json:"name"`     // Tên job để nhận dạng
	Schedule string `json:"schedule"` // Cron expression (ví dụ: "*/5 * * * *" chạy 5 phút/lần)
	Enabled  bool   `json:"enabled"`  // Bật/tắt job
}

// Job đại diện cho một pipeline job
type Job struct {
	Config   JobConfig
	Pipeline *etl.Pipeline
	EntryID  cron.EntryID // ID của job trong cron scheduler
}

// Scheduler quản lý các pipeline jobs
type Scheduler struct {
	cron     *cron.Cron
	jobs     map[string]*Job
	jobsLock sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewScheduler tạo một instance mới của Scheduler
func NewScheduler() *Scheduler {
	// Tạo cron với cấu hình:
	// - Cho phép chạy jobs song song
	// - Recover từ panic
	// - Chạy jobs với timezone local
	c := cron.New(
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow,
		)),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger), // Tự động recover từ panic
		),
	)

	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		cron:     c,
		jobs:     make(map[string]*Job),
		jobsLock: sync.RWMutex{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

// AddJob thêm một pipeline job mới vào scheduler
func (s *Scheduler) AddJob(config JobConfig, pipeline *etl.Pipeline) error {
	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()

	// Kiểm tra job đã tồn tại
	if _, exists := s.jobs[config.Name]; exists {
		return fmt.Errorf("job %s đã tồn tại", config.Name)
	}

	// Tạo job mới
	job := &Job{
		Config:   config,
		Pipeline: pipeline,
	}

	// Thêm job vào cron nếu được enable
	if config.Enabled {
		entryID, err := s.cron.AddFunc(config.Schedule, func() {
			// Tạo context mới cho mỗi lần chạy với timeout
			ctx, cancel := context.WithTimeout(s.ctx, 30*time.Minute)
			defer cancel()

			log.Printf("Bắt đầu chạy job %s\n", config.Name)
			if err := pipeline.Execute(ctx); err != nil {
				log.Printf("Lỗi chạy job %s: %v\n", config.Name, err)
				return
			}
			log.Printf("Hoàn thành job %s\n", config.Name)
		})

		if err != nil {
			return fmt.Errorf("lỗi thêm job vào cron: %v", err)
		}
		job.EntryID = entryID
	}

	// Lưu job vào map
	s.jobs[config.Name] = job
	return nil
}

// RemoveJob xóa một job khỏi scheduler
func (s *Scheduler) RemoveJob(name string) error {
	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()

	job, exists := s.jobs[name]
	if !exists {
		return fmt.Errorf("không tìm thấy job %s", name)
	}

	if job.EntryID != 0 {
		s.cron.Remove(job.EntryID)
	}

	delete(s.jobs, name)
	return nil
}

// EnableJob bật một job
func (s *Scheduler) EnableJob(name string) error {
	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()

	job, exists := s.jobs[name]
	if !exists {
		return fmt.Errorf("không tìm thấy job %s", name)
	}

	if job.Config.Enabled {
		return nil // Job đã được bật
	}

	// Thêm job vào cron
	entryID, err := s.cron.AddFunc(job.Config.Schedule, func() {
		ctx, cancel := context.WithTimeout(s.ctx, 30*time.Minute)
		defer cancel()

		log.Printf("Bắt đầu chạy job %s\n", job.Config.Name)
		if err := job.Pipeline.Execute(ctx); err != nil {
			log.Printf("Lỗi chạy job %s: %v\n", job.Config.Name, err)
			return
		}
		log.Printf("Hoàn thành job %s\n", job.Config.Name)
	})

	if err != nil {
		return fmt.Errorf("lỗi bật job: %v", err)
	}

	job.EntryID = entryID
	job.Config.Enabled = true
	return nil
}

// DisableJob tắt một job
func (s *Scheduler) DisableJob(name string) error {
	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()

	job, exists := s.jobs[name]
	if !exists {
		return fmt.Errorf("không tìm thấy job %s", name)
	}

	if !job.Config.Enabled {
		return nil // Job đã bị tắt
	}

	if job.EntryID != 0 {
		s.cron.Remove(job.EntryID)
		job.EntryID = 0
	}

	job.Config.Enabled = false
	return nil
}

// Start khởi động scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("ETL Scheduler đã được khởi động")
}

// Stop dừng scheduler
func (s *Scheduler) Stop() {
	s.cancel() // Hủy tất cả contexts
	s.cron.Stop()
	log.Println("ETL Scheduler đã dừng")
}

// GetJobs trả về danh sách các jobs
func (s *Scheduler) GetJobs() map[string]JobConfig {
	s.jobsLock.RLock()
	defer s.jobsLock.RUnlock()

	jobs := make(map[string]JobConfig)
	for name, job := range s.jobs {
		jobs[name] = job.Config
	}
	return jobs
}
