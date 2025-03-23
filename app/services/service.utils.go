package services

import (
	"errors"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Custom errors
var (
	ErrNotFound     = errors.New("record not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrDuplicate    = errors.New("duplicate record")
	ErrInvalidEmail = errors.New("invalid email format")
	ErrWeakPassword = errors.New("password is too weak")
)

// ValidateEmail kiểm tra định dạng email
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword kiểm tra độ mạnh của mật khẩu
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

// CreatePaginationOptions tạo options cho phân trang
func CreatePaginationOptions(page, limit int64) *options.FindOptions {
	skip := (page - 1) * limit
	return options.Find().
		SetSkip(skip).
		SetLimit(limit)
}

// CreateSortOptions tạo options cho sắp xếp
func CreateSortOptions(sortField string, sortOrder int) *options.FindOptions {
	return options.Find().
		SetSort(bson.D{{Key: sortField, Value: sortOrder}})
}

// CreateTimeRangeFilter tạo filter cho khoảng thời gian
func CreateTimeRangeFilter(field string, startTime, endTime time.Time) bson.M {
	filter := bson.M{}
	if !startTime.IsZero() {
		filter[field] = bson.M{"$gte": startTime}
	}
	if !endTime.IsZero() {
		if _, exists := filter[field]; exists {
			filter[field].(bson.M)["$lte"] = endTime
		} else {
			filter[field] = bson.M{"$lte": endTime}
		}
	}
	return filter
}

// CreateSearchFilter tạo filter cho tìm kiếm
func CreateSearchFilter(fields []string, searchTerm string) bson.M {
	if searchTerm == "" {
		return bson.M{}
	}

	conditions := make([]bson.M, len(fields))
	for i, field := range fields {
		conditions[i] = bson.M{
			field: bson.M{
				"$regex":   searchTerm,
				"$options": "i",
			},
		}
	}

	return bson.M{"$or": conditions}
}
