# Module chạy test
# Chức năng: Chạy test suite và báo cáo kết quả

param(
    [string]$TestPath = "./cases/...",
    [switch]$Verbose = $true
)

# Lấy đường dẫn script directory
$scriptDir = if ($PSScriptRoot) { 
    $PSScriptRoot 
} else { 
    Split-Path -Parent $MyInvocation.MyCommand.Path 
}

# Lấy đường dẫn project root (2 cấp trên scripts vì scripts nằm trong api-tests)
$projectRoot = Split-Path -Parent (Split-Path -Parent $scriptDir)

# Import utils
. "$scriptDir\utils.ps1"

# Chạy test suite
function Run-TestSuite {
    param(
        [string]$Path = "./cases/...",
        [switch]$Verbose = $true
    )
    
    Write-Step "Chạy test suite..."
    Write-Host "========================================" -ForegroundColor Cyan
    
    Push-Location $projectRoot
    try {
        $testArgs = @("test")
        if ($Verbose) {
            $testArgs += "-v"
        }
        # Đảm bảo đường dẫn test là từ project root
        if ($Path -notlike "./api-tests/*") {
            $Path = "./api-tests/$Path"
        }
        $testArgs += $Path
        
        $testOutput = & go $testArgs 2>&1
        $exitCode = $LASTEXITCODE
        
        Write-Host "========================================" -ForegroundColor Cyan
        
        # Phân tích kết quả
        $totalTests = ($testOutput | Select-String -Pattern "=== RUN" | Measure-Object).Count
        $passedTests = ($testOutput | Select-String -Pattern "--- PASS:" | Measure-Object).Count
        $failedTests = ($testOutput | Select-String -Pattern "--- FAIL:" | Measure-Object).Count
        
        # Hiển thị kết quả
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "  KẾT QUẢ TEST" -ForegroundColor Cyan
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "  Tổng số test: $totalTests" -ForegroundColor White
        Write-Host "  Passed: $passedTests" -ForegroundColor Green
        Write-Host "  Failed: $failedTests" -ForegroundColor $(if ($failedTests -gt 0) { "Red" } else { "Green" })
        
        if ($totalTests -gt 0) {
            $passRate = [math]::Round(($passedTests / $totalTests) * 100, 1)
            Write-Host "  Pass Rate: $passRate%" -ForegroundColor $(if ($passRate -eq 100) { "Green" } else { "Yellow" })
        }
        Write-Host "========================================" -ForegroundColor Cyan
        
        return @{
            ExitCode = $exitCode
            TotalTests = $totalTests
            PassedTests = $passedTests
            FailedTests = $failedTests
            Output = $testOutput
            PassRate = $passRate
        }
    }
    finally {
        Pop-Location
    }
}

# Xử lý khi được gọi trực tiếp
if ($MyInvocation.InvocationName -ne '.') {
    $result = Run-TestSuite -Path $TestPath -Verbose:$Verbose
    exit $result.ExitCode
}

