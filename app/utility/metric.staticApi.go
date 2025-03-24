package utility

import (
	"strings"

	"github.com/valyala/fasthttp"
)

// ApiInfo chứa thông tin về một API request
type ApiInfo struct {
	Time   int64  `json:"time" bson:"time"`     // Thời gian request
	Method string `json:"method" bson:"method"` // Phương thức HTTP (GET, POST, ...)
	Url    string `json:"url" bson:"url"`       // URL của request
	Result bool   `json:"result" bson:"result"` // Kết quả của request
}

// StackApiInfo là một stack chứa các ApiInfo
type StackApiInfo []ApiInfo

var stackApiInfo StackApiInfo

// IsEmpty: kiểm tra xem stack có rỗng hay không
func (s *StackApiInfo) IsEmpty() bool {
	return len(*s) == 0
}

// PushStackApiInfo: thêm một giá trị mới vào stack
func PushStackApiInfo(r *fasthttp.RequestCtx) {
	spilts := strings.Split(r.URI().String(), "?")

	str := new(ApiInfo)
	str.Time = CurrentTimeInMilli()                // Lấy thời gian hiện tại
	str.Url = spilts[0]                            // Lấy URL của request
	str.Method = string(r.Request.Header.Method()) // Lấy phương thức HTTP của request
	stackApiInfo = append(stackApiInfo, *str)      // Thêm giá trị mới vào cuối stack
}

// PopStackApiInfo: loại bỏ và trả về phần tử trên cùng của stack. Trả về false nếu stack rỗng.
func PopStackApiInfo() bool {
	if stackApiInfo.IsEmpty() {
		return false
	} else {
		index := len(stackApiInfo) - 1        // Lấy chỉ số của phần tử trên cùng.
		stackApiInfo = (stackApiInfo)[:index] // Loại bỏ phần tử đó khỏi stack bằng cách cắt bỏ.
		return true
	}
}

// RemoveStackApiInfo: loại bỏ các phần tử trong stack đã hết hạn
func RemoveStackApiInfo(expiryMinute int64) {
	expiryTime := CurrentTimeInMilli() - expiryMinute*60*1000 // Tính thời gian hết hạn

	for len(stackApiInfo) > 0 && (stackApiInfo[len(stackApiInfo)-1].Time < expiryTime) {
		PopStackApiInfo() // Loại bỏ phần tử đã hết hạn
	}
}

// ApiCount chứa thông tin về số lượng request tới một API cụ thể
type ApiCount struct {
	Url    string `json:"url" bson"url"`       // URL của API
	Method string `json:"method" bson"method"` // Phương thức HTTP
	Count  int64  `json:"count" bson"count"`   // Số lượng request
}

// GetApiStatic: lấy thống kê về các API request trong khoảng thời gian nhất định
func GetApiStatic(inSecond int64) []ApiCount {
	expiryTime := CurrentTimeInMilli() - inSecond*1000 // Tính thời gian hết hạn
	var apiCount []ApiCount
	for _, s := range stackApiInfo {
		if s.Time > expiryTime {
			isExist := false
			for i, k := range apiCount {
				if s.Url == k.Url && s.Method == k.Method {
					isExist = true
					apiCount[i].Count += 1 // Tăng số lượng request nếu API đã tồn tại trong danh sách
				}
			}

			if isExist == false {
				newApiCount := new(ApiCount)
				newApiCount.Url = s.Url
				newApiCount.Method = s.Method
				newApiCount.Count = 1

				apiCount = append(apiCount, *newApiCount) // Thêm API mới vào danh sách
			}
		}
	}

	return apiCount
}
