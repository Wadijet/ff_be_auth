# Product Context

## Tổng quan sản phẩm
FolkForm Backend Authentication (ff_be_auth) là một phần của hệ thống NextCommerce, chịu trách nhiệm xác thực và quản lý phân quyền người dùng. Nó cung cấp các API để đăng nhập, đăng ký, quản lý người dùng và quyền truy cập.

## Yêu cầu nghiệp vụ
- Xác thực người dùng (đăng nhập, đăng ký, quên mật khẩu)
- Quản lý phân quyền dựa trên vai trò (RBAC)
- Tích hợp xác thực với nền tảng xã hội (Facebook)
- Quản lý token truy cập

## Người dùng mục tiêu
- Các ứng dụng frontend của hệ thống NextCommerce
- Quản trị viên hệ thống
- Nhà phát triển tích hợp

## Các yếu tố độc đáo
- Sử dụng FastHTTP thay vì net/http chuẩn để tối ưu hiệu suất
- Tích hợp sẵn với MongoDB
- Kiến trúc middleware linh hoạt 