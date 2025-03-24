package utility

import (
	"fmt"
	"log"
	"time"

	"encoding/json"
)

// GoProtect là một hàm bao bọc (wrapper) giúp bảo vệ một hàm khác khỏi bị panic.
// Nếu xảy ra panic trong hàm f(), GoProtect sẽ bắt lại và in ra lỗi thay vì làm chương trình dừng hẳn.
func GoProtect(f func()) {
	defer func() {
		// Sử dụng recover() để bắt lỗi panic nếu có
		if err := recover(); err != nil {
			fmt.Printf("Đã bắt lỗi panic: %v\n", err)
		}
	}()

	// Gọi hàm f() được truyền vào
	f()
}

// Describe mô tả kiểu và giá trị của interface
// @params - interface cần mô tả
func Describe(t interface{}) {
	fmt.Printf("Interface type %T value %v\n", t, t)
}

// PrettyPrint in đẹp một interface dưới dạng JSON
// @params - interface cần in đẹp
// @returns - chuỗi JSON đẹp
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// UnixMilli dùng để lấy mili giây của thời gian cho trước
// @params - thời gian
// @returns - mili giây của thời gian cho trước
func UnixMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// CurrentTimeInMilli dùng để lấy thời gian hiện tại tính bằng mili giây
// Hàm này sẽ được sử dụng khi cần timestamp hiện tại
// @returns - timestamp hiện tại (tính bằng mili giây)
func CurrentTimeInMilli() int64 {
	return UnixMilli(time.Now())
}

// LogWarning ghi log cảnh báo với các thông tin bổ sung
func LogWarning(msg string, args ...interface{}) {
	// Tạo chuỗi thông tin bổ sung
	details := ""
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			details += fmt.Sprintf(" %s=%v", args[i], args[i+1])
		}
	}
	log.Printf("WARNING: %s%s", msg, details)
}
