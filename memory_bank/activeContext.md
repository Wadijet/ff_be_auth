# Active Context

## Context hiện tại
- Dự án: FolkForm Backend Authentication (ff_be_auth)
- Giai đoạn: Khởi tạo Memory Bank

## Các nhiệm vụ ưu tiên
1. Hoàn thiện cấu trúc Memory Bank
2. Xem xét và hiểu cấu trúc dự án hiện tại

## Các file quan trọng
- `main.go`: File chính của ứng dụng
- `go.mod` & `go.sum`: Quản lý dependency

## Lưu ý
- Dự án viết bằng Go, sử dụng MongoDB làm cơ sở dữ liệu
- Sử dụng FastHTTP thay vì chuẩn net/http của Go
- Có sử dụng các middleware cho xử lý request
- Tuân thủ nguyên tắc phát triển:
  - Luôn lựa chọn giải pháp tinh, gọn để nhanh chóng đưa dự án vào hoạt động
  - Luôn đề cao việc tái sử dụng code, tối thiểu hóa lượng code 