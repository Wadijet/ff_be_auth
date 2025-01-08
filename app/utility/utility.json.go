package utility

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"

	"atk-go-server/global"
)

// JSON thiết lập header và trả về dữ liệu JSON
func JSON(ctx *fasthttp.RequestCtx, data map[string]interface{}) {

	// Thiết lập Header
	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	// Chuyển đổi dữ liệu thành JSON
	res, err := json.Marshal(data)

	if err != nil {
		log.Println("Error Convert to JSON")
		data["error"] = err
	}

	// Ghi dữ liệu ra output
	ctx.Write(res)

	// Thiết lập mã trạng thái HTTP
	//ctx.SetStatusCode(statusCode)
}

// Payload tạo payload với trạng thái, dữ liệu và thông điệp
func Payload(isSuccess bool, data interface{}, message string, args ...int) map[string]interface{} {
	result := make(map[string]interface{})
	if isSuccess {
		result["status"] = "success"
	} else {
		result["status"] = "error"
	}
	result["data"] = data
	result["message"] = message
	// if len(args) == 0 {
	// 	if isSuccess {
	// 		result["code"] = 200
	// 	} else {
	// 		result["code"] = 400
	// 	}
	// } else {
	// 	result["code"] = args[0]
	// }
	return result
}

// Convert2Struct chuyển đổi dữ liệu JSON thành struct
func Convert2Struct(data []byte, myStruct interface{}) map[string]interface{} {
	reader := bytes.NewReader(data)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	err := decoder.Decode(&myStruct)
	if err != nil {
		return Payload(false, err.Error(), "Unable to decode input data!")
	}

	return nil
}

// ValidateStruct kiểm tra tính hợp lệ của struct
func ValidateStruct(myStruct interface{}) map[string]interface{} {
	err := global.Validate.Struct(myStruct)
	if err != nil {
		return Payload(false, err.Error(), "Input data is incorrect!")
	}

	return nil
}

// CreateChangeMap tạo bản đồ thay đổi từ struct
func CreateChangeMap(myStruct interface{}, myChange *map[string]interface{}) map[string]interface{} {

	CustomBson := &CustomBson{}
	change, err := CustomBson.Set(myStruct)
	if err != nil {
		return Payload(false, err.Error(), "Input data is incorrect!")
	}

	*myChange = change
	return nil
}

// FinalResponse tạo phản hồi cuối cùng dựa trên kết quả và lỗi
func FinalResponse(result interface{}, err error) map[string]interface{} {

	if err != nil {
		return Payload(false, err.Error(), "Database interaction error!")
	} else {
		return Payload(true, result, "Successful manipulation!")
	}
}

// ==========================================================================
// P2Float64 chuyển đổi interface thành float64
func P2Float64(input interface{}) float64 {
	jsonNumber, ok := input.(json.Number)
	if !ok {
		return 0
	}
	number, err := jsonNumber.Float64()
	if err != nil {
		return 0
	}

	return number
}

// P2Int64 chuyển đổi interface thành int64
func P2Int64(input interface{}) int64 {
	jsonNumber, ok := input.(json.Number)
	if !ok {
		return 0
	}
	result, err := jsonNumber.Int64()
	if err != nil {
		return 0
	}

	return result
}
