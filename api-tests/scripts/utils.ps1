# Module tiện ích chung
# Các hàm helper dùng chung cho các script khác

function Write-Step {
    param(
        [string]$Message,
        [string]$Color = "Yellow"
    )
    Write-Host "▶ $Message" -ForegroundColor $Color
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-ErrorMsg {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ $Message" -ForegroundColor Cyan
}

# Lấy đường dẫn gốc của project
function Get-ProjectRoot {
    if ($PSScriptRoot) {
        return Split-Path -Parent $PSScriptRoot
    }
    return $PWD
}

# Kiểm tra xem port có đang được sử dụng không
function Test-Port {
    param(
        [int]$Port
    )
    try {
        $connection = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
        return $null -ne $connection
    }
    catch {
        return $false
    }
}

# Đợi port sẵn sàng (có process đang listen)
function Wait-ForPort {
    param(
        [int]$Port,
        [int]$MaxWaitSeconds = 30,
        [int]$IntervalSeconds = 1
    )
    
    $elapsed = 0
    while ($elapsed -lt $MaxWaitSeconds) {
        if (Test-Port -Port $Port) {
            return $true
        }
        Start-Sleep -Seconds $IntervalSeconds
        $elapsed += $IntervalSeconds
        Write-Host "." -NoNewline
    }
    Write-Host ""
    return $false
}

# Đợi HTTP endpoint sẵn sàng
function Wait-ForHttpEndpoint {
    param(
        [string]$Url,
        [int]$MaxWaitSeconds = 30,
        [int]$IntervalSeconds = 1
    )
    
    $elapsed = 0
    while ($elapsed -lt $MaxWaitSeconds) {
        try {
            $response = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
            if ($response.StatusCode -eq 200) {
                return $true
            }
        }
        catch {
            # Endpoint chưa sẵn sàng, tiếp tục đợi
        }
        
        Start-Sleep -Seconds $IntervalSeconds
        $elapsed += $IntervalSeconds
        if ($elapsed % 5 -eq 0) {
            Write-Host "  Đang đợi... ($elapsed/$MaxWaitSeconds giây)" -ForegroundColor Gray
        } else {
            Write-Host "." -NoNewline
        }
    }
    Write-Host ""
    return $false
}

