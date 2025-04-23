# Script chạy test và tạo báo cáo
# Thiết lập UTF-8 cho terminal
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
[Console]::InputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

# Lấy đường dẫn tuyệt đối của thư mục hiện tại
$scriptPath = $PSScriptRoot
if (-not $scriptPath) {
    $scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Definition
}

# Lấy đường dẫn tới thư mục gốc của project
$projectRoot = Split-Path -Parent $scriptPath

# Tạo thư mục reports nếu chưa tồn tại
$reportsPath = Join-Path $scriptPath "reports"
if (-not (Test-Path $reportsPath)) {
    New-Item -ItemType Directory -Path $reportsPath
}

# Tạo tên file báo cáo với timestamp
$timestamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$reportFile = Join-Path $reportsPath "test_report_$timestamp.txt"

try {
    # Đọc template báo cáo
    $templatePath = Join-Path $scriptPath "templates/report_template.txt"
    $template = Get-Content -Path $templatePath -Raw -Encoding UTF8

    # Thay thế thời gian bắt đầu
    $startTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $report = $template.Replace("{{START_TIME}}", $startTime)

    # Khởi động server trong background
    Write-Host "Starting server..."
    $env:GO_ENV = "development"
    Push-Location $projectRoot
    $serverProcess = Start-Process -FilePath "go" -ArgumentList "run", ".\cmd\server\" -PassThru -WindowStyle Hidden
    Pop-Location

    # Đợi server khởi động (5 giây)
    Write-Host "Waiting for server..."
    Start-Sleep -Seconds 5

    # Chạy test và lưu kết quả
    Write-Host "Running test suite..."
    Push-Location $scriptPath
    $testOutput = & go test -v ./cases/... 2>&1
    Pop-Location

    # Dừng server
    Write-Host "Stopping server..."
    if ($serverProcess -and -not $serverProcess.HasExited) {
        Stop-Process -Id $serverProcess.Id -Force -ErrorAction SilentlyContinue
    }

    # Thay thế kết quả test
    $report = $report.Replace("{{TEST_OUTPUT}}", [string]::Join("`n", $testOutput))

    # Đếm số lượng test cases và test passed
    $totalTests = ($testOutput | Select-String -Pattern "=== RUN" | Measure-Object).Count
    $passedTests = ($testOutput | Select-String -Pattern "--- PASS:" | Measure-Object).Count
    $failedTests = ($testOutput | Select-String -Pattern "--- FAIL:" | Measure-Object).Count

    # Thay thế các placeholder còn lại
    $endTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $report = $report.Replace("{{TOTAL_TESTS}}", $totalTests)
    $report = $report.Replace("{{PASSED_TESTS}}", $passedTests)
    $report = $report.Replace("{{FAILED_TESTS}}", $failedTests)
    $report = $report.Replace("{{END_TIME}}", $endTime)

    # Ghi báo cáo với UTF-8 BOM
    $utf8Bom = New-Object System.Text.UTF8Encoding $true
    [System.IO.File]::WriteAllText($reportFile, $report, $utf8Bom)

    Write-Host "Test completed!"
    Write-Host "Report saved at: $reportFile"

    if ($failedTests -gt 0) {
        Write-Host "Warning: $failedTests test cases failed!" -ForegroundColor Red
        exit 1
    }
}
catch {
    # Đảm bảo server được dừng nếu có lỗi
    if ($serverProcess -and -not $serverProcess.HasExited) {
        Stop-Process -Id $serverProcess.Id -Force -ErrorAction SilentlyContinue
    }

    $errorMessage = "!!! TEST ERROR !!!`r`n$($_.Exception.Message)"
    [System.IO.File]::WriteAllText($reportFile, $errorMessage, $utf8Bom)
    Write-Host "Error occurred while running tests. See report for details." -ForegroundColor Red
    exit 1
} 