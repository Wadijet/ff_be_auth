package handler

// Package handler chứa các handler xử lý request HTTP trong ứng dụng.
// Package này cung cấp các chức năng CRUD cơ bản và các tiện ích để xử lý request/response.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"meta_commerce/core/api/services"
	"meta_commerce/core/common"
	"meta_commerce/core/global"
	"meta_commerce/core/utility"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	mongoopts "go.mongodb.org/mongo-driver/mongo/options"
)

// FilterOptions cấu hình cho việc validate filter
type FilterOptions struct {
	DeniedFields     []string // Các trường bị cấm filter
	AllowedOperators []string // Các operator MongoDB được phép
	MaxFields        int      // Số lượng field tối đa trong một filter
}

// BaseHandler là base handler cho các Fiber handler, cung cấp các chức năng CRUD cơ bản.
// Struct này sử dụng Generic Type để có thể tái sử dụng cho nhiều loại model khác nhau.
//
// Type parameters:
// - T: Kiểu dữ liệu của model
// - CreateInput: Kiểu dữ liệu của input khi tạo mới
// - UpdateInput: Kiểu dữ liệu của input khi cập nhật
type BaseHandler[T any, CreateInput any, UpdateInput any] struct {
	BaseService   services.BaseServiceMongo[T] // Service xử lý logic nghiệp vụ với MongoDB
	filterOptions FilterOptions                // Cấu hình validate filter
}

// NewBaseHandler tạo mới một BaseHandler với BaseService được cung cấp
func NewBaseHandler[T any, CreateInput any, UpdateInput any](baseService services.BaseServiceMongo[T]) *BaseHandler[T, CreateInput, UpdateInput] {
	return &BaseHandler[T, CreateInput, UpdateInput]{
		BaseService: baseService,
		filterOptions: FilterOptions{
			DeniedFields: []string{
				"password",
				"token",
				"secret",
				"key",
				"hash",
			},
			AllowedOperators: []string{
				"$eq",
				"$gt",
				"$gte",
				"$lt",
				"$lte",
				"$in",
				"$nin",
				"$exists",
			},
			MaxFields: 10,
		},
	}
}

// validateInput thực hiện validate chi tiết dữ liệu đầu vào
func (h *BaseHandler[T, CreateInput, UpdateInput]) validateInput(input interface{}) error {
	// Validate với validator từ global
	if err := global.Validate.Struct(input); err != nil {
		return common.NewError(common.ErrCodeValidationInput, common.MsgValidationError, common.StatusBadRequest, err)
	}

	// Kiểm tra các trường đặc biệt
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Chỉ xử lý nếu input là struct
	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Kiểm tra các trường string
		if field.Kind() == reflect.String {
			// Kiểm tra độ dài tối đa (nếu có tag maxLength)
			if maxTag := fieldType.Tag.Get("maxLength"); maxTag != "" {
				maxLen, err := strconv.Atoi(maxTag)
				if err == nil && len(field.String()) > maxLen {
					return common.NewError(
						common.ErrCodeValidationInput,
						fmt.Sprintf("Trường %s vượt quá độ dài cho phép (%d ký tự)", fieldType.Name, maxLen),
						common.StatusBadRequest,
						nil,
					)
				}
			}
		}

		// Kiểm tra các trường số
		if field.Kind() == reflect.Int || field.Kind() == reflect.Int64 {
			// Kiểm tra giá trị tối thiểu (nếu có tag min)
			if minTag := fieldType.Tag.Get("min"); minTag != "" {
				min, err := strconv.ParseInt(minTag, 10, 64)
				if err == nil && field.Int() < min {
					return common.NewError(
						common.ErrCodeValidationInput,
						fmt.Sprintf("Trường %s phải lớn hơn hoặc bằng %d", fieldType.Name, min),
						common.StatusBadRequest,
						nil,
					)
				}
			}

			// Kiểm tra giá trị tối đa (nếu có tag max)
			if maxTag := fieldType.Tag.Get("max"); maxTag != "" {
				max, err := strconv.ParseInt(maxTag, 10, 64)
				if err == nil && field.Int() > max {
					return common.NewError(
						common.ErrCodeValidationInput,
						fmt.Sprintf("Trường %s phải nhỏ hơn hoặc bằng %d", fieldType.Name, max),
						common.StatusBadRequest,
						nil,
					)
				}
			}
		}
	}

	return nil
}

// ParseRequestBody parse và validate dữ liệu từ request body.
// Sử dụng json.Decoder với UseNumber() để xử lý chính xác các số.
//
// Parameters:
// - c: Fiber context
// - input: Con trỏ tới struct sẽ chứa dữ liệu được parse
//
// Returns:
// - error: Lỗi nếu có trong quá trình parse hoặc validate
func (h *BaseHandler[T, CreateInput, UpdateInput]) ParseRequestBody(c fiber.Ctx, input interface{}) error {
	// Parse body thành struct T
	body := c.Body()
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(input); err != nil {
		return common.NewError(common.ErrCodeValidationFormat, common.MsgValidationError, common.StatusBadRequest, err)
	}

	// Validate chi tiết input
	if err := h.validateInput(input); err != nil {
		return err
	}

	return nil
}

// ParseRequestQuery parse và validate dữ liệu từ query string.
// Query string phải được encode dưới dạng JSON.
//
// Parameters:
// - c: Fiber context
// - input: Con trỏ tới struct sẽ chứa dữ liệu được parse
//
// Returns:
// - error: Lỗi nếu có trong quá trình parse hoặc validate
func (h *BaseHandler[T, CreateInput, UpdateInput]) ParseRequestQuery(c fiber.Ctx, input interface{}) error {
	query := c.Query("query", "")

	// Parse query
	reader := bytes.NewReader([]byte(query))
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	if err := decoder.Decode(input); err != nil {
		return common.NewError(common.ErrCodeValidationFormat, common.MsgValidationError, common.StatusBadRequest, err)
	}

	// Validate struct
	if err := global.Validate.Struct(input); err != nil {
		return common.NewError(common.ErrCodeValidationInput, common.MsgValidationError, common.StatusBadRequest, err)
	}

	return nil
}

// ParseRequestParams parse và validate các tham số từ URI.
// Sử dụng Fiber's URI binding để parse các tham số.
//
// Parameters:
// - c: Fiber context
// - input: Con trỏ tới struct sẽ chứa dữ liệu được parse
//
// Returns:
// - error: Lỗi nếu có trong quá trình parse hoặc validate
func (h *BaseHandler[T, CreateInput, UpdateInput]) ParseRequestParams(c fiber.Ctx, input interface{}) error {
	// Parse URI params
	if err := c.Bind().URI(input); err != nil {
		return common.NewError(common.ErrCodeValidationFormat, common.MsgValidationError, common.StatusBadRequest, err)
	}

	// Validate struct
	if err := global.Validate.Struct(input); err != nil {
		return common.NewError(common.ErrCodeValidationInput, common.MsgValidationError, common.StatusBadRequest, err)
	}

	return nil
}

// validateFilter kiểm tra tính hợp lệ của filter
func (h *BaseHandler[T, CreateInput, UpdateInput]) validateFilter(filter map[string]interface{}) error {
	// Kiểm tra số lượng field
	if len(filter) > h.filterOptions.MaxFields {
		return common.NewError(
			common.ErrCodeValidationFormat,
			fmt.Sprintf("Số lượng field filter không được vượt quá %d", h.filterOptions.MaxFields),
			common.StatusBadRequest,
			nil,
		)
	}

	// Kiểm tra từng field và operator
	for field, value := range filter {
		// Kiểm tra field có bị cấm không
		if utility.Contains(h.filterOptions.DeniedFields, field) {
			return common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Field không được phép filter: %s", field),
				common.StatusBadRequest,
				nil,
			)
		}

		// Kiểm tra operator nếu value là map
		if mapValue, ok := value.(map[string]interface{}); ok {
			for op := range mapValue {
				if strings.HasPrefix(op, "$") && !utility.Contains(h.filterOptions.AllowedOperators, op) {
					return common.NewError(
						common.ErrCodeValidationFormat,
						fmt.Sprintf("Operator không được phép sử dụng: %s", op),
						common.StatusBadRequest,
						nil,
					)
				}
			}
		}
	}

	return nil
}

// processFilter xử lý và validate filter từ request
func (h *BaseHandler[T, CreateInput, UpdateInput]) processFilter(c fiber.Ctx) (map[string]interface{}, error) {
	var filter map[string]interface{}

	// Parse filter từ query
	if err := json.Unmarshal([]byte(c.Query("filter", "{}")), &filter); err != nil {
		return nil, common.NewError(
			common.ErrCodeValidationFormat,
			"Filter không hợp lệ",
			common.StatusBadRequest,
			err,
		)
	}

	// Validate filter
	if err := h.validateFilter(filter); err != nil {
		return nil, err
	}

	return filter, nil
}

// validateMongoOptions kiểm tra tính hợp lệ của các options
func (h *BaseHandler[T, CreateInput, UpdateInput]) validateMongoOptions(options map[string]interface{}) error {
	// Danh sách các options được phép
	allowedOptions := map[string]bool{
		"projection": true,
		"sort":       true,
		"limit":      true,
		"skip":       true,
	}

	// Kiểm tra các options không hợp lệ
	for key := range options {
		if !allowedOptions[key] {
			return common.NewError(
				common.ErrCodeValidationFormat,
				fmt.Sprintf("Option không hợp lệ: %s", key),
				common.StatusBadRequest,
				nil,
			)
		}
	}

	// Validate projection
	if projection, ok := options["projection"].(map[string]interface{}); ok {
		for field := range projection {
			// Kiểm tra các trường bị cấm
			if utility.Contains(h.filterOptions.DeniedFields, field) {
				return common.NewError(
					common.ErrCodeValidationFormat,
					fmt.Sprintf("Field không được phép projection: %s", field),
					common.StatusBadRequest,
					nil,
				)
			}
		}
	}

	// Validate sort
	if sort, ok := options["sort"].(map[string]interface{}); ok {
		for field, value := range sort {
			// Kiểm tra các trường bị cấm
			if utility.Contains(h.filterOptions.DeniedFields, field) {
				return common.NewError(
					common.ErrCodeValidationFormat,
					fmt.Sprintf("Field không được phép sort: %s", field),
					common.StatusBadRequest,
					nil,
				)
			}
			// Kiểm tra giá trị sort (1 hoặc -1)
			if v, ok := value.(float64); !ok || (v != 1 && v != -1) {
				return common.NewError(
					common.ErrCodeValidationFormat,
					fmt.Sprintf("Giá trị sort không hợp lệ cho field %s: %v (phải là 1 hoặc -1)", field, value),
					common.StatusBadRequest,
					nil,
				)
			}
		}
	}

	// Validate limit
	if limit, ok := options["limit"].(float64); ok {
		if limit <= 0 {
			return common.NewError(
				common.ErrCodeValidationFormat,
				"Limit phải lớn hơn 0",
				common.StatusBadRequest,
				nil,
			)
		}
		if limit > 1000 { // Giới hạn tối đa 1000 documents
			return common.NewError(
				common.ErrCodeValidationFormat,
				"Limit không được vượt quá 1000",
				common.StatusBadRequest,
				nil,
			)
		}
	}

	// Validate skip
	if skip, ok := options["skip"].(float64); ok {
		if skip < 0 {
			return common.NewError(
				common.ErrCodeValidationFormat,
				"Skip không được âm",
				common.StatusBadRequest,
				nil,
			)
		}
	}

	return nil
}

// processMongoOptions xử lý options từ query string và chuyển đổi sang MongoDB options
func (h *BaseHandler[T, CreateInput, UpdateInput]) processMongoOptions(c fiber.Ctx, isFindOne bool) (interface{}, error) {
	var rawOptions map[string]interface{}

	// Parse options từ query string
	if err := json.Unmarshal([]byte(c.Query("options", "{}")), &rawOptions); err != nil {
		return nil, common.NewError(
			common.ErrCodeValidationFormat,
			"Options không hợp lệ",
			common.StatusBadRequest,
			err,
		)
	}

	// Validate options
	if err := h.validateMongoOptions(rawOptions); err != nil {
		return nil, err
	}

	// Chuyển đổi sang MongoDB options
	if isFindOne {
		opts := mongoopts.FindOne()
		if projection, ok := rawOptions["projection"].(map[string]interface{}); ok {
			opts.SetProjection(projection)
		}
		if sort, ok := rawOptions["sort"].(map[string]interface{}); ok {
			opts.SetSort(sort)
		}
		return opts, nil
	}

	opts := mongoopts.Find()
	if projection, ok := rawOptions["projection"].(map[string]interface{}); ok {
		opts.SetProjection(projection)
	}
	if sort, ok := rawOptions["sort"].(map[string]interface{}); ok {
		opts.SetSort(sort)
	}
	if limit, ok := rawOptions["limit"].(float64); ok {
		opts.SetLimit(int64(limit))
	}
	if skip, ok := rawOptions["skip"].(float64); ok {
		opts.SetSkip(int64(skip))
	}
	return opts, nil
}

// ParsePagination xử lý việc parse thông tin phân trang từ request.
// Hỗ trợ các tham số:
// - page: Số trang (mặc định: 1)
// - limit: Số lượng item trên một trang (mặc định: 10)
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - page: Số trang
// - limit: Số lượng item trên một trang
func (h *BaseHandler[T, CreateInput, UpdateInput]) ParsePagination(c fiber.Ctx) (int64, int64) {
	page := utility.P2Int64(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}

	limit := utility.P2Int64(c.Query("limit", "10"))
	if limit <= 0 {
		limit = 10
	}

	return page, limit
}

// GetIDFromContext lấy ID từ URI params của request.
//
// Parameters:
// - c: Fiber context
//
// Returns:
// - string: ID từ params
func (h *BaseHandler[T, CreateInput, UpdateInput]) GetIDFromContext(c fiber.Ctx) string {
	return c.Params("id")
}
