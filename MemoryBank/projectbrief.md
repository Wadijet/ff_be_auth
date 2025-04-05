# Mô tả dự án

## Tổng quan
Dự án FolkForm Backend Authentication (ff_be_auth) là một phần của hệ thống NextCommerce, phụ trách về xác thực và quản lý quyền người dùng.

## Công nghệ
- Ngôn ngữ: Go
- Framework: FastHTTP với Router tùy chỉnh
- Database: MongoDB
- Middleware: CORS, Rate limiting, Timeout, Recovery, Measure

## Cấu trúc dự án
- **app/**: Chứa mã nguồn chính của ứng dụng
  - **database/**: Tương tác với cơ sở dữ liệu
  - **global/**: Biến và cấu hình toàn cục
  - **middleware/**: Các middleware xử lý request
  - **models/**: Các model đại diện cho dữ liệu
  - **router/**: Định nghĩa các endpoint API
  - **services/**: Logic nghiệp vụ
  - **utility/**: Các hàm tiện ích
- **config/**: Cấu hình ứng dụng
- **database/**: Kết nối và quản lý cơ sở dữ liệu
- **main.go**: Điểm khởi chạy ứng dụng

## Chức năng chính
- Xác thực người dùng
- Quản lý phân quyền (Role-based Access Control)
- Quản lý token truy cập
- Tích hợp với các nền tảng xã hội (Facebook) 