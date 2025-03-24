# Hướng dẫn Test API

## Cấu trúc Test

- `api_test.go`: Chứa các test case cho API
- `setup_test.go`: Chứa các hàm helper cho việc setup test environment

## Cách chạy Test

### Chạy test locally

1. Đảm bảo MongoDB đã được cài đặt và đang chạy
2. Chạy lệnh sau để cài đặt dependencies:
```bash
go mod download
```

3. Chạy test:
```bash
# Chạy tất cả test
go test ./... -v

# Chạy test cho một package cụ thể
go test ./tests -v

# Chạy test cho một function cụ thể
go test ./tests -v -run TestStaticAPI
```

### Chạy test với Newman (Postman)

1. Cài đặt Newman:
```bash
npm install -g newman
```

2. Chạy collection test:
```bash
newman run ../doc/SS_key_server.postman_collection.json
```

## Cấu trúc Test Case

Mỗi test case nên bao gồm:

1. Setup test environment
2. Thực hiện request
3. Kiểm tra response
4. Cleanup (nếu cần)

## Best Practices

1. Sử dụng subtest để tổ chức các test case liên quan
2. Mock các service bên ngoài khi cần thiết
3. Cleanup dữ liệu test sau mỗi test case
4. Sử dụng các assertion từ package testify
5. Viết test case cho cả trường hợp thành công và thất bại

## CI/CD

Test được tự động chạy thông qua GitHub Actions khi:
- Push code lên branch main
- Tạo pull request vào branch main

Xem file `.github/workflows/test.yml` để biết thêm chi tiết về quy trình CI/CD. 