# Script tạo self-signed certificate cho development
# Chạy script này để tạo certificate và key file cho HTTPS

$certDir = Join-Path $PSScriptRoot "..\config\tls"
$certFile = Join-Path $certDir "server.crt"
$keyFile = Join-Path $certDir "server.key"

# Tạo thư mục nếu chưa tồn tại
if (-not (Test-Path $certDir)) {
    New-Item -ItemType Directory -Path $certDir -Force | Out-Null
    Write-Host "Đã tạo thư mục: $certDir" -ForegroundColor Green
}

# Kiểm tra xem đã có certificate chưa
if ((Test-Path $certFile) -and (Test-Path $keyFile)) {
    Write-Host "Certificate đã tồn tại tại:" -ForegroundColor Yellow
    Write-Host "  Cert: $certFile" -ForegroundColor Yellow
    Write-Host "  Key:  $keyFile" -ForegroundColor Yellow
    $overwrite = Read-Host "Bạn có muốn tạo lại không? (y/N)"
    if ($overwrite -ne "y" -and $overwrite -ne "Y") {
        Write-Host "Bỏ qua tạo certificate mới." -ForegroundColor Gray
        exit 0
    }
}

Write-Host "Đang tạo self-signed certificate..." -ForegroundColor Cyan

# Tạo certificate và key bằng OpenSSL
try {
    # Thử dùng OpenSSL nếu có
    $opensslPath = Get-Command openssl -ErrorAction SilentlyContinue
    if ($opensslPath) {
        Write-Host "Sử dụng OpenSSL..." -ForegroundColor Gray
        
        # Chạy OpenSSL command
        $opensslArgs = @(
            "req",
            "-x509",
            "-newkey", "rsa:4096",
            "-keyout", $keyFile,
            "-out", $certFile,
            "-days", "365",
            "-nodes",
            "-subj", "/C=VN/ST=HCM/L=HoChiMinh/O=Development/CN=localhost"
        )
        
        & openssl $opensslArgs
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Đã tạo certificate thành công!" -ForegroundColor Green
            Write-Host "  Cert: $certFile" -ForegroundColor Green
            Write-Host "  Key:  $keyFile" -ForegroundColor Green
            Write-Host ""
            Write-Host "⚠️  Lưu ý: Đây là self-signed certificate, trình duyệt sẽ cảnh báo." -ForegroundColor Yellow
            Write-Host "   Chấp nhận cảnh báo để tiếp tục sử dụng." -ForegroundColor Yellow
            exit 0
        } else {
            Write-Host "❌ Lỗi khi tạo certificate với OpenSSL" -ForegroundColor Red
            exit 1
        }
    }
} catch {
    Write-Host "Không tìm thấy OpenSSL, đang thử cách khác..." -ForegroundColor Yellow
}

# Nếu không có OpenSSL, hướng dẫn cài đặt
Write-Host ""
Write-Host "❌ Không tìm thấy OpenSSL!" -ForegroundColor Red
Write-Host ""
Write-Host "Có 2 cách để tạo certificate:" -ForegroundColor Yellow
Write-Host ""
Write-Host "Cách 1: Cài OpenSSL và chạy lại script này" -ForegroundColor Cyan
Write-Host "  - Windows: choco install openssl" -ForegroundColor Gray
Write-Host "  - Hoặc tải từ: https://slproweb.com/products/Win32OpenSSL.html" -ForegroundColor Gray
Write-Host ""
Write-Host "Cách 2: Tạo thủ công bằng lệnh sau:" -ForegroundColor Cyan
$subj = "/C=VN/ST=HCM/L=HoChiMinh/O=Development/CN=localhost"
Write-Host "  openssl req -x509 -newkey rsa:4096 -keyout `"$keyFile`" -out `"$certFile`" -days 365 -nodes -subj $subj" -ForegroundColor Gray
Write-Host ""
exit 1
