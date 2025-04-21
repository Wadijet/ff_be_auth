# FolkForm Authentication Backend (ff_be_auth)

## Tổng Quan Dự Án
FolkForm Authentication Backend là một phần của hệ thống NextCommerce, cung cấp các dịch vụ xác thực và phân quyền cho nền tảng FolkForm.

## Cấu Trúc Hệ Thống
Dự án được tổ chức theo mô hình Clean Architecture với các thành phần chính:

### Core Components
- `cmd/`: Entry points của ứng dụng (server, worker)
- `core/`: Logic nghiệp vụ chính
  - `api/`: REST API handlers và middleware
  - `database/`: Database interfaces và implementations
  - `services/`: Business logic services
  - `models/`: Domain models
- `config/`: Cấu hình ứng dụng
- `custom_modes/`: Các mode tùy chỉnh cho phát triển

### Công Nghệ Sử Dụng
- Ngôn ngữ: Go
- Framework: (cần xác định)
- Database: MongoDB
- Authentication: (cần xác định)

## Trạng Thái Hiện Tại
- Dự án đang trong giai đoạn phát triển
- Đã có cấu trúc cơ bản của hệ thống
- Đang sử dụng Memory Bank system mới với các mode phát triển chuyên biệt

## Mục Tiêu
1. Xây dựng hệ thống xác thực an toàn và hiệu quả
2. Tích hợp với các thành phần khác của NextCommerce
3. Đảm bảo khả năng mở rộng và bảo trì

## Các Mode Phát Triển
- VAN: Khởi tạo và xác định độ phức tạp
- PLAN: Lập kế hoạch chi tiết
- CREATIVE: Khám phá giải pháp thiết kế
- IMPLEMENT: Triển khai code
- QA: Kiểm tra và đảm bảo chất lượng 