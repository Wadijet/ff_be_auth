package utility

import (
	"fmt"
	"time"

	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

// String2ObjectID chuyển đổi chuỗi thành ObjectID
// @params - chuỗi cần chuyển đổi
// @returns - ObjectID
func String2ObjectID(id string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return objectId
}

// ObjectID2String chuyển đổi ObjectID thành chuỗi
// @params - ObjectID cần chuyển đổi
// @returns - chuỗi ObjectID
func ObjectID2String(id primitive.ObjectID) string {
	stringObjectID := id.Hex()
	return stringObjectID
}
