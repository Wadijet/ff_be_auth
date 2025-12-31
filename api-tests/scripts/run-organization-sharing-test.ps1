# Script chạy test Organization Sharing
# Sử dụng: .\api-tests\scripts\run-organization-sharing-test.ps1

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  TEST ORGANIZATION SHARING" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Kiểm tra Firebase token
$firebaseToken = $env:TEST_FIREBASE_ID_TOKEN
if (-not $firebaseToken) {
    Write-Host "[WARNING] TEST_FIREBASE_ID_TOKEN chưa được set" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Cách lấy Firebase token:" -ForegroundColor Cyan
    Write-Host "1. Lấy Firebase API Key từ config hoặc Firebase Console" -ForegroundColor White
    Write-Host "2. Chạy script:" -ForegroundColor White
    Write-Host "   .\api-tests\scripts\get-firebase-token.ps1 -Email 'test@example.com' -Password 'Test@123' -ApiKey 'YOUR_API_KEY'" -ForegroundColor Gray
    Write-Host "3. Set environment variable:" -ForegroundColor White
    Write-Host "   `$env:TEST_FIREBASE_ID_TOKEN='your_token_here'" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Hoặc chạy test với token:" -ForegroundColor Cyan
    Write-Host "   `$env:TEST_FIREBASE_ID_TOKEN='token'; go test -v ./api-tests/cases -run TestOrganizationSharing" -ForegroundColor Gray
    Write-Host ""
    
    # Hỏi user có muốn tiếp tục không
    $continue = Read-Host "Bạn có muốn tiếp tục chạy test không? (y/n)"
    if ($continue -ne "y" -and $continue -ne "Y") {
        exit 0
    }
}

# Kiểm tra server
Write-Host "[1/2] Kiểm tra server..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/system/health" -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
    if ($response.StatusCode -eq 200) {
        Write-Host "  [OK] Server đang chạy" -ForegroundColor Green
    }
} catch {
    Write-Host "  [ERROR] Server chưa chạy hoặc không thể kết nối" -ForegroundColor Red
    Write-Host "  Vui lòng khởi động server trước khi chạy test" -ForegroundColor Yellow
    exit 1
}

# Chạy test
Write-Host "[2/2] Chạy test..." -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan

$projectRoot = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $projectRoot

go test -v ./api-tests/cases -run TestOrganizationSharing

$exitCode = $LASTEXITCODE

Write-Host "========================================" -ForegroundColor Cyan

if ($exitCode -eq 0) {
    Write-Host ""
    Write-Host "✅ TEST PASSED!" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ TEST FAILED!" -ForegroundColor Red
}

exit $exitCode
