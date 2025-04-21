/*
Package scheduler cung cấp chức năng quản lý và thực thi các tác vụ định kỳ (cron jobs).
Package này sử dụng thư viện robfig/cron để quản lý việc lập lịch các tác vụ.

Các tính năng chính:
- Khởi tạo và quản lý scheduler
- Thêm/xóa/theo dõi các jobs
- Đồng bộ hóa truy cập vào scheduler thông qua mutex
- Hỗ trợ định dạng cron expression với độ chính xác đến giây
*/
package scheduler

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"
)

// Scheduler đại diện cho một scheduler quản lý các cron jobs.
// Struct này đảm bảo thread-safe thông qua việc sử dụng RWMutex.
type Scheduler struct {
	// cron là instance của cron scheduler từ thư viện robfig/cron
	cron *cron.Cron
	// jobs lưu trữ map giữa tên job và ID của nó trong cron scheduler
	jobs map[string]cron.EntryID
	// mu là mutex để đồng bộ hóa truy cập vào scheduler
	mu sync.RWMutex
}

// NewScheduler tạo một instance mới của Scheduler.
// Scheduler được khởi tạo với:
// - Cron scheduler có độ chính xác đến giây
// - Map rỗng để lưu trữ jobs
func NewScheduler() *Scheduler {
	return &Scheduler{
		// WithSeconds() cho phép định nghĩa cron expression với độ chính xác đến giây
		cron: cron.New(cron.WithSeconds()),
		jobs: make(map[string]cron.EntryID),
	}
}

// Start khởi động scheduler.
// Sau khi gọi Start, scheduler sẽ bắt đầu thực thi các jobs theo lịch đã định nghĩa.
// Các jobs mới có thể được thêm vào ngay cả khi scheduler đang chạy.
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop dừng scheduler một cách an toàn.
// - Dừng tất cả các jobs đang chạy
// - Đợi cho đến khi tất cả jobs hoàn thành
// - Trả về context để caller có thể theo dõi khi nào scheduler dừng hoàn toàn
func (s *Scheduler) Stop() context.Context {
	return s.cron.Stop()
}

// AddJob thêm một job mới vào scheduler.
// Tham số:
// - name: Tên định danh của job
// - spec: Biểu thức cron định nghĩa lịch chạy (vd: "0 0 * * *" - chạy lúc 00:00 mỗi ngày)
// - job: Hàm thực thi của job
// Trả về error nếu biểu thức cron không hợp lệ
func (s *Scheduler) AddJob(name string, spec string, job func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Nếu job đã tồn tại, xóa job cũ trước khi thêm job mới
	if _, exists := s.jobs[name]; exists {
		s.RemoveJob(name)
	}

	// Thêm job mới vào cron scheduler
	id, err := s.cron.AddFunc(spec, job)
	if err != nil {
		return err
	}

	// Lưu ID của job để có thể quản lý sau này
	s.jobs[name] = id
	return nil
}

// RemoveJob xóa một job khỏi scheduler dựa trên tên của job.
// Job sẽ không được lập lịch chạy nữa sau khi bị xóa.
// Nếu job không tồn tại, hàm này không làm gì cả.
func (s *Scheduler) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id, exists := s.jobs[name]; exists {
		s.cron.Remove(id)
		delete(s.jobs, name)
	}
}

// GetJobs trả về danh sách các jobs đang được quản lý bởi scheduler.
// Trả về một bản sao của map jobs để tránh data race.
// Key là tên job, value là ID của job trong cron scheduler.
func (s *Scheduler) GetJobs() map[string]cron.EntryID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make(map[string]cron.EntryID)
	for k, v := range s.jobs {
		jobs[k] = v
	}
	return jobs
}
