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

// ResponseType định nghĩa kiểu response
type ResponseType struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

// Payload tạo payload với trạng thái, dữ liệu và thông điệp
func Payload(isSuccess bool, data interface{}, message string, statusCode ...int) map[string]interface{} {
	response := ResponseType{
		Status:  "error",
		Data:    data,
		Message: message,
		Code:    StatusInternalServerError,
	}

	if isSuccess {
		response.Status = "success"
		response.Code = StatusOK
	}

	if len(statusCode) > 0 {
		response.Code = statusCode[0]
	}

	result := make(map[string]interface{})
	result["status"] = response.Status
	result["data"] = response.Data
	result["message"] = response.Message
	result["code"] = response.Code

	return result
}

// Convert2Struct chuyển đổi dữ liệu JSON thành struct
func Convert2Struct(data []byte, myStruct interface{}) map[string]interface{} {
	reader := bytes.NewReader(data)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	err := decoder.Decode(&myStruct)
	if err != nil {
		return Payload(false, NewError(ErrCodeValidationFormat, MsgInvalidFormat, StatusBadRequest, err), MsgInvalidFormat)
	}

	return nil
}

// ValidateStruct kiểm tra tính hợp lệ của struct
func ValidateStruct(myStruct interface{}) map[string]interface{} {
	err := global.Validate.Struct(myStruct)
	if err != nil {
		return Payload(false, NewError(ErrCodeValidationInput, MsgValidationError, StatusBadRequest, err), MsgValidationError)
	}

	return nil
}

// CreateChangeMap tạo bản đồ thay đổi từ struct
func CreateChangeMap(myStruct interface{}, myChange *map[string]interface{}) map[string]interface{} {
	CustomBson := &CustomBson{}
	change, err := CustomBson.Set(myStruct)
	if err != nil {
		return Payload(false, NewError(ErrCodeValidationInput, MsgValidationError, StatusBadRequest, err), MsgValidationError)
	}

	*myChange = change
	return nil
}

// FinalResponse tạo phản hồi cuối cùng dựa trên kết quả và lỗi
func FinalResponse(result interface{}, err error) map[string]interface{} {
	if err != nil {
		if customErr, ok := err.(*Error); ok {
			return Payload(false, customErr, customErr.Message, customErr.StatusCode)
		}
		return Payload(false, NewError(ErrCodeDatabaseConnection, MsgDatabaseError, StatusInternalServerError, err), MsgDatabaseError)
	} else {
		return Payload(true, result, MsgSuccess, StatusOK)
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
