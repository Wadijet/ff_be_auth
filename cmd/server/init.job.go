package main

import (
	"fmt"

	"meta_commerce/core/scheduler"
	"meta_commerce/core/scheduler/jobs"
)

func InitJobs(s *scheduler.Scheduler) error {
	// Khởi tạo SyncFbPagesJob
	syncFbPagesJob, err := jobs.NewSyncFbPagesJob()
	if err != nil {
		return fmt.Errorf("failed to create sync_fb_pages job: %v", err)
	}

	// Đăng ký job vào scheduler
	err = s.AddJob(syncFbPagesJob.Name(), syncFbPagesJob.Schedule(), func() {
		if err := syncFbPagesJob.Run(); err != nil {
			jobs.GetJobLogger().WithError(err).Error("Error running sync_fb_pages job")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add sync_fb_pages job to scheduler: %v", err)
	}

	return nil
}
