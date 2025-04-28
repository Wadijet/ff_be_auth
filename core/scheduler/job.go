/*
Package scheduler định nghĩa các interface và model cần thiết cho việc quản lý jobs.
File này cung cấp các thành phần cơ bản để xây dựng một job:
- Interface Job định nghĩa các phương thức cần thiết
- Struct JobMetadata lưu trữ thông tin về một lần chạy job
- Struct BaseJob cung cấp triển khai cơ bản của interface Job
*/
package scheduler

import (
	"context"
	"time"
)

// ================== INTERFACE ĐỊNH NGHĨA JOB ==================

// Job là interface chuẩn cho mọi job trong hệ thống.
type Job interface {
	// Execute thực thi logic chính của job
	// ctx: context để kiểm soát thời gian thực thi và hủy job
	// Trả về error nếu có lỗi xảy ra trong quá trình thực thi
	Execute(ctx context.Context) error

	// GetName trả về tên định danh của job
	// Tên này được sử dụng để đăng ký và quản lý job trong scheduler
	GetName() string

	// GetSchedule trả về biểu thức cron định nghĩa lịch chạy của job
	// Ví dụ: "0 0 * * *" - chạy lúc 00:00 mỗi ngày
	GetSchedule() string
}

// ================== BASE JOB ==================

// BaseJob cung cấp sẵn name, schedule và các hàm mặc định.
// Các job cụ thể chỉ cần nhúng *BaseJob và implement Execute.
type BaseJob struct {
	name     string
	schedule string
}

// NewBaseJob khởi tạo BaseJob với tên và lịch chạy.
func NewBaseJob(name, schedule string) *BaseJob {
	return &BaseJob{name: name, schedule: schedule}
}

func (j *BaseJob) GetName() string     { return j.name }
func (j *BaseJob) GetSchedule() string { return j.schedule }
func (j *BaseJob) Execute(ctx context.Context) error {
	// Mặc định không làm gì, job con phải override
	return nil
}

// ================== TRẠNG THÁI & METADATA ==================

// JobStatus là enum trạng thái job.
type JobStatus string

const (
	// JobStatusPending: job đã được lập lịch nhưng chưa bắt đầu chạy
	JobStatusPending JobStatus = "pending"
	// JobStatusRunning: job đang trong quá trình thực thi
	JobStatusRunning JobStatus = "running"
	// JobStatusCompleted: job đã hoàn thành thành công
	JobStatusCompleted JobStatus = "completed"
	// JobStatusFailed: job thực thi thất bại, có thể cần retry
	JobStatusFailed JobStatus = "failed"
)

// JobMetadata lưu thông tin về từng lần chạy job.
type JobMetadata struct {
	// Name: tên định danh của job
	Name string `json:"name" bson:"name"`
	// Schedule: biểu thức cron định nghĩa lịch chạy
	Schedule string `json:"schedule" bson:"schedule"`
	// Status: trạng thái hiện tại của job
	Status JobStatus `json:"status" bson:"status"`
	// LastRun: thời điểm job chạy lần cuối
	LastRun time.Time `json:"last_run" bson:"last_run"`
	// NextRun: thời điểm dự kiến job sẽ chạy lần tiếp theo
	NextRun time.Time `json:"next_run" bson:"next_run"`
	// Duration: thời gian thực thi của lần chạy cuối (tính bằng giây)
	Duration float64 `json:"duration" bson:"duration"`
	// Error: thông tin lỗi nếu job thất bại
	Error string `json:"error,omitempty" bson:"error,omitempty"`
	// RetryCount: số lần đã retry
	RetryCount int `json:"retry_count" bson:"retry_count"`
	// MaxRetries: số lần retry tối đa cho phép
	MaxRetries int `json:"max_retries" bson:"max_retries"`
	// CreatedAt: thời điểm job được tạo
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	// UpdatedAt: thời điểm cập nhật thông tin gần nhất
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
