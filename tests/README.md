# Hướng Dẫn Chạy Test

## Cấu trúc thư mục
```
tests/
├── cases/              # Chứa các test case
├── reports/            # Chứa báo cáo test
├── templates/          # Templates cho báo cáo
├── utils/             # Utilities cho test
└── run_tests.ps1      # Script chạy test tự động
```

## Quy trình chạy test đơn giản

### 1. Chạy test và tạo báo cáo
```powershell
# Di chuyển vào thư mục tests
cd tests

# Chạy script test tự động
powershell -File run_tests.ps1
```

Script sẽ tự động:
- Khởi động server trong background
- Đợi server sẵn sàng (5 giây)
- Chạy toàn bộ test cases
- Tự động dừng server sau khi test xong
- Tạo báo cáo chi tiết trong thư mục reports/
- Hiển thị kết quả tổng quan

### 2. Xem kết quả test
Báo cáo test được lưu tại `reports/test_report_YYYY-MM-DD_HH-mm-ss.txt` với nội dung:
- Thời gian bắt đầu và kết thúc
- Danh sách test cases đã chạy
- Kết quả từng test (PASS/FAIL)
- Thống kê tổng hợp:
  + Tổng số test cases
  + Số test passed
  + Số test failed

### 3. Phân tích kết quả
- ✅ PASS: Test case chạy thành công
- ❌ FAIL: Test case thất bại, xem chi tiết lỗi trong báo cáo
- 🕒 SKIP: Test case bị bỏ qua

## Cơ chế khởi động server tự động

### 1. Quy trình khởi động server
Script `run_tests.ps1` tự động quản lý vòng đời của server:

1. **Khởi động server:**
   ```powershell
   # Thiết lập môi trường
   $env:GO_ENV = "development"
   
   # Khởi động server trong background
   Push-Location $projectRoot
   $serverProcess = Start-Process -FilePath "go" -ArgumentList "run", ".\cmd\server\" -PassThru -WindowStyle Hidden
   Pop-Location
   ```

2. **Đợi server sẵn sàng:**
   ```powershell
   # Đợi 5 giây để server khởi động hoàn tất
   Write-Host "Đợi server khởi động..."
   Start-Sleep -Seconds 5
   ```

3. **Dừng server sau khi test:**
   ```powershell
   # Dừng server
   Stop-Process -Id $serverProcess.Id -Force
   ```

4. **Xử lý lỗi:**
   ```powershell
   # Đảm bảo server được dừng khi có lỗi
   if ($serverProcess -and -not $serverProcess.HasExited) {
       Stop-Process -Id $serverProcess.Id -Force -ErrorAction SilentlyContinue
   }
   ```

### 2. Ưu điểm của cơ chế tự động
1. **Tự động hóa hoàn toàn:**
   - Không cần khởi động server thủ công
   - Tránh quên tắt server sau khi test
   - Giảm thiểu sai sót do con người

2. **Quản lý tài nguyên tốt:**
   - Server chỉ chạy trong thời gian test
   - Tự động giải phóng port sau khi test
   - Xử lý graceful shutdown

3. **Xử lý lỗi tốt:**
   - Tự động dừng server khi test fail
   - Đảm bảo không có process zombie
   - Log rõ ràng về trạng thái server

### 3. Lưu ý quan trọng
1. **Port conflict:**
   - Script sẽ fail nếu port 8080 đang bị sử dụng
   - Cần đảm bảo không có instance server nào đang chạy

2. **Thời gian chờ:**
   - Mặc định đợi 5 giây cho server khởi động
   - Có thể điều chỉnh thời gian này nếu cần

3. **Debug mode:**
   - Server chạy ở chế độ background
   - Xem log trong file báo cáo để debug

## Lưu ý quan trọng
- Đảm bảo server đang chạy trước khi test
- Server chạy mặc định ở port 8080
- Báo cáo được tự động lưu trong thư mục reports/
- Mỗi báo cáo có timestamp riêng để dễ theo dõi

## Chạy test nâng cao
```powershell
# Chạy test cho module cụ thể
go test -v ./cases/your_module/...

# Chạy test có pattern nhất định
go test -v -run "TestPattern" ./cases/...

# Chạy test và hiển thị coverage
go test -v -cover ./cases/...
```

## Lưu Ý Quan Trọng Về PowerShell

> ⚠️ **Chú ý khi chạy lệnh trong PowerShell:**
> 1. PowerShell KHÔNG hỗ trợ toán tử `&&` như Unix/Linux. Thay vào đó:
>    ```powershell
>    # KHÔNG dùng
>    cd thư_mục && chạy_lệnh    # ❌ Sẽ gây lỗi
>    
>    # Nên dùng
>    cd thư_mục                  # ✅ Chạy từng lệnh riêng biệt
>    chạy_lệnh
>    
>    # Hoặc dùng dấu chấm phẩy
>    cd thư_mục; chạy_lệnh      # ✅ Nhiều lệnh trên một dòng
>    ```
> 
> 2. Để chạy nhiều lệnh:
>    ```powershell
>    # Sử dụng dấu chấm phẩy
>    cd thư_mục; go test        # ✅ Cách 1
>    
>    # Hoặc dùng pipe
>    cd thư_mục | go test       # ✅ Cách 2
>    ```
>


## Cấu trúc thư mục
```
tests/
├── cases/              # Chứa các test case
│   └── health_test.go  # Test health check API
├── utils/              # Utilities cho test
│   └── http_client.go  # HTTP client wrapper
├── main_test.go        # File test chính
└── README.md           # File này
```

## Quy trình chạy test

### 1. Khởi động server

#### Cách 1: Chạy bằng VS Code (KHUYẾN NGHỊ)
- Mở VS Code
- Nhấn F5 để chạy server với launch config có sẵn
- Kiểm tra trong Debug Console để đảm bảo server đã chạy thành công
- Có thể đặt breakpoint để debug khi cần

> ⚠️ **Lưu ý khi chạy bằng VS Code:**
> 1. Đảm bảo đang mở đúng workspace của project
> 2. Kiểm tra server đã khởi động thành công khi thấy thông báo "Server is running on :8080"
> 3. Nếu thấy lỗi "address already in use" -> Port 8080 đang bị chiếm, cần tắt process đang dùng port này

#### Cách 2: Chạy bằng command line (Chỉ dùng khi cần thiết)
```powershell
# Mở terminal, cd vào thư mục gốc
cd C:\Projects\DMD\NextCommerce\FolkForm\ff_be_auth

# Chạy server trực tiếp để dễ theo dõi lỗi
go run cmd/server/main.go
```

> ⚠️ **Lưu ý khi chạy bằng command line:**
> 1. Chỉ sử dụng cách này khi không thể chạy bằng VS Code
> 2. Có thể gặp một số lỗi về import hoặc init functions
> 3. Nếu gặp lỗi, nên quay lại sử dụng VS Code với launch config có sẵn

### Debug Process Server

#### Theo dõi lỗi khi chạy test
1. **Debug với VS Code (Khuyến nghị):**
   - Mở VS Code, nhấn F5 để chạy server trong chế độ debug
   - Đặt breakpoint tại các endpoint đang test
   - Theo dõi giá trị biến và stack trace trong Debug Console
   - Xem log chi tiết trong Debug Console khi test fail

2. **Debug với command line:**
   ```powershell
   # Chạy server với log chi tiết
   go run cmd/server/main.go -v debug
   
   # Mở terminal khác để theo dõi log
   Get-Content -Path "logs/server.log" -Wait
   ```

#### Cách phân tích lỗi
1. **Khi test fail:**
   - Kiểm tra log server để xem request/response
   - Xem stack trace trong báo cáo test
   - So sánh expected vs actual response
   - Kiểm tra các điều kiện tiên quyết (DB, cache, etc)

2. **Lỗi kết nối:**
   - Kiểm tra server có đang chạy không
   - Xác nhận port 8080 không bị chiếm
   - Kiểm tra firewall và network settings

3. **Lỗi nghiệp vụ:**
   - Đặt breakpoint tại handler function
   - Theo dõi luồng xử lý request
   - Kiểm tra giá trị các biến middleware
   - Xem log validation và business logic

#### Công cụ hỗ trợ debug
1. **VS Code:**
   - Debug Console: Xem log và theo dõi biến
   - Call Stack: Xem luồng thực thi
   - Watch: Theo dõi giá trị biến
   - Breakpoints: Dừng tại điểm cần kiểm tra

2. **Delve debugger:**
   ```powershell
   # Cài đặt delve
   go install github.com/go-delve/delve/cmd/dlv@latest
   
   # Debug với delve
   dlv debug cmd/server/main.go
   ```

3. **Log viewer:**
   - Tail logs trong real-time
   - Lọc log theo severity
   - Tìm kiếm theo pattern

#### Best Practices
1. **Logging:**
   - Log đầy đủ request/response
   - Thêm request ID để trace
   - Log chi tiết các lỗi nghiệp vụ

2. **Testing:**
   - Test từng endpoint riêng biệt
   - Chuẩn bị test data đầy đủ
   - Clean up sau mỗi test case

3. **Monitoring:**
   - Theo dõi memory usage
   - Kiểm tra goroutine leaks
   - Monitor response time

### 2. Chạy test suite

#### Cách 1: Chạy test và tạo báo cáo đẹp
```powershell
# Di chuyển vào thư mục tests
cd tests

# Kiểm tra server đã chạy chưa
$serverCheck = try { 
    Invoke-WebRequest -Uri "http://localhost:8080/api/v1/system/health" -UseBasicParsing -TimeoutSec 1 
} catch { $null }

if (-not $serverCheck) {
    Write-Host "🚀 Server chưa chạy, đang khởi động..."
    # Khởi động VS Code với debug session
    Start-Process -FilePath "code" -ArgumentList "--folder", ".." -Wait
    # Gửi F5 để start debug
    Add-Type -AssemblyName System.Windows.Forms
    [System.Windows.Forms.SendKeys]::SendWait("{F5}")
    
    # Đợi server khởi động
    Write-Host "⏳ Đang đợi server khởi động..."
    do {
        Start-Sleep -Seconds 1
        $serverCheck = try { 
            Invoke-WebRequest -Uri "http://localhost:8080/api/v1/system/health" -UseBasicParsing -TimeoutSec 1 
        } catch { $null }
    } while (-not $serverCheck)
    Write-Host "✅ Server đã khởi động thành công!"
}

Write-Host "🧪 Bắt đầu chạy test..."

# Tạo thư mục reports nếu chưa có
if (-not (Test-Path reports)) { New-Item -ItemType Directory -Path reports }

# Chạy test và tạo báo cáo đẹp
$timestamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$reportFile = "reports/test_report_$timestamp.txt"

# Tạo header báo cáo
@"
===========================================
BÁO CÁO KẾT QUẢ TEST
Thời gian bắt đầu: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
===========================================

"@ | Out-File -FilePath $reportFile -Encoding UTF8

# Chạy test và lưu kết quả
$testOutput = go test -v -count=1 ./cases/... 2>&1
$testOutput | Tee-Object -Append -FilePath $reportFile

# Tạo footer báo cáo
$totalTests = ($testOutput | Select-String -Pattern "=== RUN" | Measure-Object).Count
$passedTests = ($testOutput | Select-String -Pattern "--- PASS:" | Measure-Object).Count
$failedTests = ($testOutput | Select-String -Pattern "--- FAIL:" | Measure-Object).Count

@"

===========================================
TỔNG KẾT BÁO CÁO
- Tổng số test cases: $totalTests
- Số test passed: $passedTests
- Số test failed: $failedTests
- Thời gian kết thúc: $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")
===========================================
"@ | Add-Content -Path $reportFile -Encoding UTF8

# Tạo symlink cho báo cáo mới nhất
Copy-Item -Path $reportFile -Destination "reports/latest_report.txt" -Force

Write-Host "`n✨ Báo cáo đã được lưu tại: $reportFile"
```

> 💡 **Hoàn toàn tự động:**
> 1. Tự động kiểm tra server
> 2. Tự động khởi động VS Code và debug session
> 3. Tự động đợi server sẵn sàng
> 4. Tự động tạo và format báo cáo
> 5. Tự động tính toán kết quả test

#### Cách 2: Chạy test riêng từng module
```powershell
# Chạy test cho health check
cd tests
go test -v ./cases/health_test.go

# Chạy test cho một package cụ thể
go test -v ./cases/your_package/...

# Chạy test có pattern cụ thể
go test -v -run "TestHealth" ./cases/...
```

#### Cách 3: Chạy test với các tùy chọn hữu ích
```powershell
# Chạy test và hiển thị coverage
go test -v -cover ./cases/...

# Chạy test và xuất coverage report
go test -v -coverprofile=coverage.out ./cases/...
go tool cover -html=coverage.out -o coverage.html

# Chạy test với timeout cụ thể
go test -v -timeout 30s ./cases/...

# Chạy test ở chế độ verbose và hiển thị log
go test -v -test.v ./cases/...
```

> 💡 **Mẹo hay khi chạy test:**
> 1. Sử dụng `-v` để xem chi tiết từng test case
> 2. Dùng `-run` để chạy test case cụ thể
> 3. Thêm `-count=1` để disable test cache
> 4. Sử dụng `| Tee-Object` để vừa xem kết quả vừa lưu file
> 5. Kết hợp với `grep` để lọc kết quả: `go test -v ./... | grep -E "FAIL|PASS"`

### 3. Xem kết quả test

#### Xem kết quả trực tiếp
- ✅ PASS: Test case chạy thành công
- ❌ FAIL: Test case thất bại, xem chi tiết lỗi bên dưới
- 🕒 SKIP: Test case bị bỏ qua (có thể do điều kiện không phù hợp)

#### Phân tích lỗi test
Khi test fail, bạn sẽ thấy:
1. Tên test case bị fail
2. File và line number gây lỗi
3. Expected vs Actual value
4. Stack trace (nếu có panic)

Ví dụ về test fail:
```
--- FAIL: TestHealthCheck (0.02s)
    health_test.go:25: 
        Error Trace:    health_test.go:25
        Error:          Not equal: 
                       expected: 200
                       actual  : 500
        Test:          TestHealthCheck
```

#### Các bước debug test fail
1. Xem log của server để kiểm tra request/response
2. Check điều kiện tiên quyết (server running, database, etc)
3. Đặt breakpoint tại vị trí fail trong VS Code
4. Chạy lại test cụ thể với flag -v để xem chi tiết

## Lưu ý
- Server phải được khởi động trước khi chạy test
- Server chạy mặc định ở port 8080
- Mỗi test case nên được đặt trong thư mục `cases/`
- Các utility function nên được đặt trong thư mục `utils/`
- Báo cáo test được lưu tự động trong thư mục `reports/` với định dạng tên: `test_report_YYYY-MM-DD_HH-mm-ss.txt`

## Cách thêm test case mới

1. Tạo file test mới trong thư mục `cases/`:
```go
package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewFeature(t *testing.T) {
    // Code test ở đây
}
```

2. Sử dụng HTTP client có sẵn từ utils:
```go
client := utils.NewHTTPClient("http://localhost:8080", 5)
resp, body, err := client.GET("/api/v1/your-endpoint")
``` 

## Báo Cáo Test

### Cấu trúc báo cáo
Báo cáo test được tự động tạo trong thư mục `reports/` với cấu trúc như sau:
```
reports/
├── test_report_2024-03-20_14-30-00.txt   # Báo cáo theo ngày giờ
├── test_report_2024-03-20_15-45-00.txt
└── latest_report.txt                      # Báo cáo mới nhất
```

### Nội dung báo cáo
Mỗi báo cáo test bao gồm:
- Thời gian bắt đầu và kết thúc test
- Tổng số test case đã chạy
- Kết quả từng test case (Pass/Fail)
- Chi tiết lỗi nếu test case thất bại
- Thống kê tổng hợp (% pass/fail)
- Thời gian chạy của từng test case

### Cách đọc báo cáo
1. **Xem báo cáo mới nhất:**
```powershell
type reports\latest_report.txt
```

2. **Xem báo cáo theo ngày:**
```powershell
type reports\test_report_YYYY-MM-DD_HH-mm-ss.txt
```

3. **Tìm kiếm test case thất bại:**
- Tìm dòng bắt đầu bằng "FAIL" trong báo cáo
- Xem chi tiết lỗi ở phần stack trace bên dưới

### Lưu trữ báo cáo
- Báo cáo được tự động lưu trữ trong thư mục `reports/`
- Định dạng tên file: `test_report_YYYY-MM-DD_HH-mm-ss.txt`
- File `latest_report.txt` luôn chứa kết quả của lần chạy test gần nhất
- Nên giữ lại báo cáo để theo dõi lịch sử và phát hiện regression 