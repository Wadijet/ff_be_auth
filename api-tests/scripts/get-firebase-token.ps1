# Script lấy Firebase ID token cho test
# Sử dụng: .\api-tests\scripts\get-firebase-token.ps1 -Email "test@example.com" -Password "Test@123"

param(
    [string]$Email = "test@example.com",
    [string]$Password = "Test@123",
    [string]$ApiKey = "",
    [string]$ProjectId = ""
)

# Kiểm tra API Key
if ($ApiKey -eq "") {
    # Thử lấy từ environment variable
    $ApiKey = $env:FIREBASE_API_KEY
    if ($ApiKey -eq "") {
        Write-Host "[ERROR] Firebase API Key không được cung cấp" -ForegroundColor Red
        Write-Host "Sử dụng: -ApiKey 'YOUR_API_KEY' hoặc set FIREBASE_API_KEY environment variable" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  LẤY FIREBASE ID TOKEN CHO TEST" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Tạo request body
$body = @{
    email             = $Email
    password          = $Password
    returnSecureToken = $true
} | ConvertTo-Json

# Gọi Firebase REST API
$url = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=$ApiKey"

try {
    Write-Host "Đang đăng nhập với Firebase..." -ForegroundColor Yellow
    $response = Invoke-RestMethod -Uri $url `
        -Method Post `
        -Body $body `
        -ContentType "application/json"

    $idToken = $response.idToken
    if ($idToken) {
        Write-Host "[OK] Đăng nhập thành công!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Firebase ID Token:" -ForegroundColor Cyan
        Write-Host $idToken -ForegroundColor White
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "  SET ENVIRONMENT VARIABLE" -ForegroundColor Cyan
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "PowerShell:" -ForegroundColor Yellow
        Write-Host "`$env:TEST_FIREBASE_ID_TOKEN=`"$idToken`"" -ForegroundColor White
        Write-Host ""
        Write-Host "Bash:" -ForegroundColor Yellow
        Write-Host "export TEST_FIREBASE_ID_TOKEN=`"$idToken`"" -ForegroundColor White
        Write-Host ""
        Write-Host "Hoặc chạy lệnh sau để set tự động:" -ForegroundColor Yellow
        Write-Host "`$env:TEST_FIREBASE_ID_TOKEN=`"$idToken`"" -ForegroundColor White
        
        # Hỏi có muốn set tự động không
        $setAuto = Read-Host "Bạn có muốn set token vào environment variable ngay bây giờ? (y/n)"
        if ($setAuto -eq "y" -or $setAuto -eq "Y") {
            $env:TEST_FIREBASE_ID_TOKEN = $idToken
            Write-Host "[OK] Đã set TEST_FIREBASE_ID_TOKEN vào environment variable" -ForegroundColor Green
            Write-Host "Lưu ý: Token chỉ có hiệu lực trong session hiện tại" -ForegroundColor Yellow
        }
    } else {
        Write-Host "[ERROR] Không nhận được ID token từ response" -ForegroundColor Red
        Write-Host "Response: $($response | ConvertTo-Json)" -ForegroundColor Yellow
        exit 1
    }
} catch {
    Write-Host "[ERROR] Lỗi khi đăng nhập:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    
    if ($_.ErrorDetails.Message) {
        $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
        Write-Host "Chi tiết: $($errorDetails.error.message)" -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Host "Gợi ý:" -ForegroundColor Yellow
    Write-Host "1. Kiểm tra email và password có đúng không" -ForegroundColor White
    Write-Host "2. Đảm bảo user đã được tạo trong Firebase Console" -ForegroundColor White
    Write-Host "3. Kiểm tra Firebase API Key có đúng không" -ForegroundColor White
    exit 1
}

