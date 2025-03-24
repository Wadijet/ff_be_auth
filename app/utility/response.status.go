package utility

// HTTP Status Code Constants
const (
	// Success Codes (2xx)
	StatusOK        = 200 // Thành công
	StatusCreated   = 201 // Tạo mới thành công
	StatusAccepted  = 202 // Yêu cầu được chấp nhận
	StatusNoContent = 204 // Thành công nhưng không có nội dung trả về

	// Client Error Codes (4xx)
	StatusBadRequest   = 400 // Yêu cầu không hợp lệ
	StatusUnauthorized = 401 // Chưa xác thực
	StatusForbidden    = 403 // Không có quyền truy cập
	StatusNotFound     = 404 // Không tìm thấy tài nguyên
	StatusConflict     = 409 // Xung đột dữ liệu

	// Server Error Codes (5xx)
	StatusInternalServerError = 500 // Lỗi server
	StatusServiceUnavailable  = 503 // Dịch vụ không khả dụng
)

// Response Messages
const (
	MsgSuccess         = "Thao tác thành công"
	MsgCreated         = "Tạo mới thành công"
	MsgBadRequest      = "Yêu cầu không hợp lệ"
	MsgUnauthorized    = "Vui lòng đăng nhập"
	MsgForbidden       = "Không có quyền truy cập"
	MsgNotFound        = "Không tìm thấy tài nguyên"
	MsgInternalError   = "Lỗi hệ thống"
	MsgValidationError = "Dữ liệu không hợp lệ"
	MsgDatabaseError   = "Lỗi tương tác với cơ sở dữ liệu"
)

// Error định nghĩa cấu trúc lỗi
type Error struct {
	Message    string
	StatusCode int
}

// Error trả về message của lỗi
func (e *Error) Error() string {
	return e.Message
}

// NewError tạo một error mới với status code và message
func NewError(statusCode int, message string) error {
	return &Error{
		Message:    message,
		StatusCode: statusCode,
	}
}
