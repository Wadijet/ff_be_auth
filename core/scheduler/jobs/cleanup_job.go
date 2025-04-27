/*
Package examples cung cấp các ví dụ về cách triển khai và sử dụng jobs.
File này minh họa cách tạo một job cụ thể bằng cách kế thừa từ BaseJob.

CleanupJob là một ví dụ về job định kỳ dọn dẹp dữ liệu cũ trong hệ thống:
- Chạy vào 00:00 mỗi ngày
- Xóa các dữ liệu cũ hơn X ngày
- Có thể tùy chỉnh số ngày lưu trữ thông qua tham số retentionDays
*/
package jobs

import (
	"context"
	"fmt"
	"time"

	"meta_commerce/core/scheduler"
)

// CleanupJob là một ví dụ về job dọn dẹp dữ liệu.
// Job này kế thừa từ BaseJob và thêm trường retentionDays để cấu hình
// thời gian lưu trữ dữ liệu.
type CleanupJob struct {
	// Nhúng BaseJob để kế thừa các chức năng cơ bản
	*scheduler.BaseJob
	// retentionDays: số ngày giữ lại dữ liệu, dữ liệu cũ hơn sẽ bị xóa
	retentionDays int
}

// NewCleanupJob tạo một instance mới của CleanupJob.
// Tham số:
// - retentionDays: số ngày giữ lại dữ liệu (vd: 30 ngày)
// Job được cấu hình để chạy vào 00:00 mỗi ngày thông qua cron expression "0 0 * * *"
func NewCleanupJob(retentionDays int) *CleanupJob {
	return &CleanupJob{
		// Khởi tạo BaseJob với tên "cleanup" và lịch chạy "0 0 * * *"
		BaseJob:       scheduler.NewBaseJob("cleanup", "0 0 * * *"),
		retentionDays: retentionDays,
	}
}

// Execute triển khai logic dọn dẹp dữ liệu.
// Phương thức này được gọi tự động bởi scheduler theo lịch đã định nghĩa.
// Tham số:
// - ctx: context để kiểm soát thời gian thực thi và hủy job
// Trả về error nếu có lỗi xảy ra trong quá trình dọn dẹp
func (j *CleanupJob) Execute(ctx context.Context) error {
	// In thông báo bắt đầu với số ngày cấu hình
	fmt.Printf("Bắt đầu dọn dẹp dữ liệu cũ hơn %d ngày...\n", j.retentionDays)

	// Tính thời điểm ngưỡng để xóa dữ liệu
	// Ví dụ: Nếu retentionDays = 30, threshold sẽ là thời điểm 30 ngày trước
	threshold := time.Now().AddDate(0, 0, -j.retentionDays)

	// TODO: Triển khai logic xóa dữ liệu cũ
	// Các loại dữ liệu cần xem xét:
	// 1. Log files:
	//    - Application logs
	//    - Access logs
	//    - Error logs
	// 2. Temporary files:
	//    - Upload temps
	//    - Cache files
	//    - Backup files
	// 3. Database records:
	//    - Audit logs
	//    - Session data
	//    - Temporary data
	// 4. Cache data:
	//    - Redis cache
	//    - Memory cache
	//    - File cache

	// In thông báo hoàn thành với thời điểm ngưỡng
	fmt.Printf("Đã hoàn thành dọn dẹp dữ liệu cũ hơn %s\n", threshold.Format("2006-01-02"))
	return nil
}
