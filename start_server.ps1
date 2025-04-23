# Đặt biến môi trường GO_ENV
$env:GO_ENV = "development"

# Đường dẫn tới thư mục gốc của project
$projectRoot = $PSScriptRoot

# Chạy server
Write-Host "Starting server from $projectRoot\cmd\server..."
Set-Location $projectRoot
go run .\cmd\server\ 