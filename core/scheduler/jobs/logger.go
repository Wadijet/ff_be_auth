package jobs

import (
	"meta_commerce/core/logger"

	"github.com/sirupsen/logrus"
)

// GetJobLogger trả về logger chuyên dụng cho jobs
func GetJobLogger() *logrus.Logger {
	return logger.GetLogger("jobs")
}
