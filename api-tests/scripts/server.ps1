# Module quản lý server
# Chức năng: Khởi động, dừng, và kiểm tra trạng thái server

param(
    [string]$Action = "",  # start, stop, status (để trống khi import)
    [int]$Port = 8080,
    [string]$HealthEndpoint = "/api/v1/system/health"
)

$script:ServerProcess = $null
$script:ServerPort = $Port
$script:ServerHealthUrl = "http://localhost:$Port$HealthEndpoint"

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

# Khởi động server
function Start-TestServer {
    Write-Host "▶ Khởi động server trên port $script:ServerPort" -ForegroundColor Yellow
    
    # Kiểm tra xem server đã chạy chưa
    if (Test-Port -Port $script:ServerPort) {
        Write-Host "ℹ Port $script:ServerPort đã được sử dụng. Đang dừng process cũ..." -ForegroundColor Cyan
        Stop-TestServer
        Start-Sleep -Seconds 2
    }
    
    # Thiết lập biến môi trường
    $env:GO_ENV = "development"
    
    # Khởi động server trong background
    Push-Location $projectRoot
    try {
        $script:ServerProcess = Start-Process -FilePath "go" `
            -ArgumentList "run", ".\api\cmd\server\" `
            -PassThru `
            -WindowStyle Hidden
        
        Write-Host "ℹ Server process đã khởi động (PID: $($script:ServerProcess.Id))" -ForegroundColor Cyan
        
        # Đợi server sẵn sàng (tăng thời gian đợi lên 60 giây)
        Write-Host "▶ Đợi server sẵn sàng (có thể mất 20-40 giây)..." -ForegroundColor Yellow
        $ready = Wait-ForHttpEndpoint -Url $script:ServerHealthUrl -MaxWaitSeconds 60
        
        if ($ready) {
            Write-Host "✓ Server đã sẵn sàng!" -ForegroundColor Green
            return $script:ServerProcess
        } else {
            Write-Host "✗ Server không sẵn sàng sau 30 giây" -ForegroundColor Red
            Stop-TestServer
            return $null
        }
    }
    finally {
        Pop-Location
    }
}

# Dừng server
function Stop-TestServer {
    Write-Host "▶ Dừng server..." -ForegroundColor Yellow
    
    # Dừng process nếu có
    if ($script:ServerProcess -and -not $script:ServerProcess.HasExited) {
        try {
            Stop-Process -Id $script:ServerProcess.Id -Force -ErrorAction SilentlyContinue
            Write-Host "✓ Đã dừng server process (PID: $($script:ServerProcess.Id))" -ForegroundColor Green
        }
        catch {
            Write-Host "ℹ Process đã tự dừng hoặc không tồn tại" -ForegroundColor Cyan
        }
        $script:ServerProcess = $null
    }
    
    # Đợi một chút để port được giải phóng
    Start-Sleep -Seconds 1
    
    # Kiểm tra lại xem port đã được giải phóng chưa
    if (Test-Port -Port $script:ServerPort) {
        Write-Host "ℹ Port $script:ServerPort vẫn đang được sử dụng, đang tìm và dừng process..." -ForegroundColor Cyan
        try {
            $processes = Get-NetTCPConnection -LocalPort $script:ServerPort -ErrorAction SilentlyContinue | 
                Select-Object -ExpandProperty OwningProcess -Unique
            if ($processes) {
                foreach ($pid in $processes) {
                    try {
                        Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
                        Write-Host "  ℹ Đã dừng process $pid" -ForegroundColor Cyan
                    }
                    catch {
                        # Bỏ qua nếu không thể dừng
                    }
                }
                Start-Sleep -Seconds 1
            }
        }
        catch {
            # Bỏ qua lỗi
        }
    }
    
    Write-Host "✓ Server đã được dừng" -ForegroundColor Green
}

# Kiểm tra trạng thái server
function Get-TestServerStatus {
    if ($script:ServerProcess -and -not $script:ServerProcess.HasExited) {
        try {
            $response = Invoke-WebRequest -Uri $script:ServerHealthUrl -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
            if ($response.StatusCode -eq 200) {
                Write-Host "✓ Server đang chạy và sẵn sàng (PID: $($script:ServerProcess.Id))" -ForegroundColor Green
                return $true
            }
        }
        catch {
            Write-Host "✗ Server process đang chạy nhưng không phản hồi" -ForegroundColor Red
            return $false
        }
    } else {
        Write-Host "ℹ Server không đang chạy" -ForegroundColor Cyan
        return $false
    }
}

# Chỉ xử lý action khi được gọi trực tiếp và có tham số Action
if ($MyInvocation.InvocationName -ne '.' -and $Action -ne "") {
    switch ($Action.ToLower()) {
        "start" {
            Start-TestServer | Out-Null
        }
        "stop" {
            Stop-TestServer
        }
        "status" {
            Get-TestServerStatus | Out-Null
        }
        default {
            Write-Host "✗ Invalid action: $Action. Use: start, stop, status" -ForegroundColor Red
        }
    }
}

