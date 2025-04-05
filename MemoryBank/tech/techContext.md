# Technical Context

## Stack công nghệ
- **Ngôn ngữ**: Go (Golang)
- **Framework Web**: FastHTTP với router tùy chỉnh
- **Cơ sở dữ liệu**: MongoDB
- **Middleware**: CORS, Rate limiting, Timeout, Recovery, Measure
- **Validation**: go-playground/validator

## Kiến trúc hệ thống
- **API RESTful**: Các endpoint cung cấp dịch vụ xác thực và phân quyền
- **Mô hình đa tầng**: Phân tách router, service, model và database layer
- **Stateless**: Sử dụng JWT hoặc token để xác thực giữa các request

## Phát triển và triển khai
- **Môi trường phát triển**: Phát triển trên máy Windows 10
- **Quản lý phụ thuộc**: Go Modules (go.mod, go.sum)
- **Quản lý mã nguồn**: Git

## Lưu ý kỹ thuật
- Sử dụng cấu trúc middleware FastHTTP khác với chuẩn net/http
- MongoDB yêu cầu cấu hình kết nối và indexes riêng
- Các xử lý async và concurrent được quản lý qua goroutines và channels 