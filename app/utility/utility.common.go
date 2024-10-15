package utility

import (
	"fmt"
	"time"

	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GoProtect bảo vệ hàm khỏi panic
// @params - hàm cần bảo vệ
func GoProtect(f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recovered: %v", err)
		}
	}()

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

// ****************************************************  Bson *******************************************
// Các thao tác Bson tùy chỉnh

// CustomBson dùng để thực hiện các thao tác bson tùy chỉnh
// như set, push, unset, v.v. bằng cách sử dụng các struct
// Điều này rất hữu ích khi cần tạo bản đồ bson từ struct
type CustomBson struct{}

// BsonWrapper chứa các thao tác bson cơ bản
// như $set, $push, $addToSet
// Nó rất hữu ích để chuyển đổi struct thành bson
type BsonWrapper struct {

	// Set sẽ đặt dữ liệu trong db
	// ví dụ - nếu cần đặt "name":"Jack", thì cần tạo một struct chứa trường name và gán struct đó vào trường này.
	// Sau khi mã hóa thành bson, nó sẽ như { $set : {name : "Jack"}} và điều này sẽ hữu ích trong truy vấn mongo
	Set interface{} `json:"$set,omitempty" bson:"$set,omitempty"`

	// Toán tử Unset xóa một trường cụ thể.
	// Nếu trường không tồn tại, thì Unset không làm gì cả
	// Nếu cần unset trường name thì chỉ cần tạo một struct chứa trường name và gán "" cho name.
	// Bây giờ để unset, gán struct đó vào trường Unset. Sau khi mã hóa, nó sẽ trở thành { $unset: { name: "" } }
	Unset interface{} `json:"$unset,omitempty" bson:"$unset,omitempty"`

	// Toán tử Push thêm một giá trị cụ thể vào một mảng.
	// Nếu trường không có trong tài liệu để cập nhật,
	// Push thêm trường mảng với giá trị là phần tử của nó.
	// Nếu trường không phải là một mảng, thao tác sẽ thất bại.
	Push interface{} `json:"$push,omitempty" bson:"$push,omitempty"`

	// Toán tử AddToSet thêm một giá trị vào một mảng trừ khi giá trị đã có, trong trường hợp đó AddToSet không làm gì với mảng đó.
	// Nếu sử dụng AddToSet trên một trường không có trong tài liệu để cập nhật,
	// AddToSet tạo trường mảng với giá trị cụ thể là phần tử của nó.
	AddToSet interface{} `json:"$addToSet,omitempty" bson:"$addToSet,omitempty"`
}

// ToMap chuyển đổi interface thành bản đồ.
// Nó nhận interface làm tham số và trả về bản đồ và lỗi nếu có
func ToMap(s interface{}) (map[string]interface{}, error) {
	var stringInterfaceMap map[string]interface{}
	itr, err := bson.Marshal(s)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(itr, &stringInterfaceMap)
	return stringInterfaceMap, err
}

// Set tạo truy vấn để thay thế giá trị của một trường bằng giá trị cụ thể
// @params - dữ liệu cần đặt
// @returns - bản đồ truy vấn và lỗi nếu có
func (customBson *CustomBson) Set(data interface{}) (map[string]interface{}, error) {
	s := BsonWrapper{Set: data}
	return ToMap(s)
}

// Push tạo truy vấn để thêm một giá trị cụ thể vào một trường mảng
// @params - dữ liệu cần thêm
// @returns - bản đồ truy vấn và lỗi nếu có
func (customBson *CustomBson) Push(data interface{}) (map[string]interface{}, error) {
	s := BsonWrapper{Push: data}
	return ToMap(s)
}

// Unset tạo truy vấn để xóa một trường cụ thể
// @params - dữ liệu cần unset
// @returns - bản đồ truy vấn và lỗi nếu có
func (customBson *CustomBson) Unset(data interface{}) (map[string]interface{}, error) {
	s := BsonWrapper{Unset: data}
	return ToMap(s)
}

// AddToSet tạo truy vấn để thêm một giá trị vào một mảng trừ khi giá trị đã có.
// @params - dữ liệu cần thêm vào set
// @returns - bản đồ truy vấn và lỗi nếu có
func (customBson *CustomBson) AddToSet(data interface{}) (map[string]interface{}, error) {
	s := BsonWrapper{AddToSet: data}
	return ToMap(s)
}

// ****************************************************  Bson End  *******************************************

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
