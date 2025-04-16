package utility

import (
	"go.mongodb.org/mongo-driver/bson"
)

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
