package global

import (
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/go-playground/validator.v9"
)

// InitValidator khởi tạo và đăng ký các custom validator
func InitValidator() {
	// Khởi tạo validator
	Validate = validator.New()

	// Đăng ký các custom validator
	_ = Validate.RegisterValidation("no_xss", validateNoXSS)
	_ = Validate.RegisterValidation("no_sql_injection", validateNoSQLInjection)
	_ = Validate.RegisterValidation("strong_password", validateStrongPassword)
}

// validateNoXSS kiểm tra XSS
func validateNoXSS(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	dangerousPatterns := []string{
		"<script",
		"javascript:",
		"onerror=",
		"onload=",
		"onclick=",
		"onmouseover=",
		"eval(",
		"document.cookie",
		"document.write",
		"innerHTML",
		"fromCharCode",
		"window.location",
		"<iframe",
		"<object",
		"<embed",
	}

	value = strings.ToLower(value)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(value, pattern) {
			return false
		}
	}
	return true
}

// validateNoSQLInjection kiểm tra SQL Injection
func validateNoSQLInjection(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	sqlPatterns := []string{
		"'",
		";",
		"--",
		"/*",
		"*/",
		"xp_",
		"SELECT",
		"DROP",
		"DELETE",
		"UPDATE",
		"INSERT",
		"UNION",
		"OR 1=1",
		"OR '1'='1",
		"OR 'a'='a",
		"OR 1 = 1",
		"WAITFOR",
		"DELAY",
		"BENCHMARK",
	}

	value = strings.ToUpper(value)
	for _, pattern := range sqlPatterns {
		if strings.Contains(value, strings.ToUpper(pattern)) {
			return false
		}
	}
	return true
}

// validateStrongPassword kiểm tra mật khẩu mạnh
// DEPRECATED: Không còn sử dụng - Firebase quản lý authentication và password
// Giữ lại để tương thích ngược, nhưng không nên sử dụng
func validateStrongPassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Kiểm tra độ dài tối thiểu
	if len(value) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Yêu cầu ít nhất 3 trong 4 điều kiện
	conditions := 0
	if hasUpper {
		conditions++
	}
	if hasLower {
		conditions++
	}
	if hasNumber {
		conditions++
	}
	if hasSpecial {
		conditions++
	}

	return conditions >= 3
}

// validateEmail kiểm tra định dạng email
func validateEmail(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(value)
}
