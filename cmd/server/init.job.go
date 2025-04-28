package main

import (
	"context"
	"fmt"

	"meta_commerce/core/scheduler"
	"meta_commerce/core/scheduler/jobs"
)

func InitJobs(s *scheduler.Scheduler) error {
	// Khởi tạo SyncFbPagesJob với schedule động
	schedule := "0 */5 * * * *" // Có thể thay đổi khi cần
	syncFbPagesJob, err := jobs.NewSyncFbPagesJob(schedule)
	if err != nil {
		return fmt.Errorf("failed to create sync_fb_pages job: %v", err)
	}

	// Đăng ký job vào scheduler
	err = s.AddJob(syncFbPagesJob.GetName(), syncFbPagesJob.GetSchedule(), func() {
		ctx := context.Background()
		if err := syncFbPagesJob.Execute(ctx); err != nil {
			jobs.GetJobLogger().WithError(err).Error("Error running sync_fb_pages job")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add sync_fb_pages job to scheduler: %v", err)
	}

	return nil
}
