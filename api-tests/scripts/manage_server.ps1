# Script quản lý server độc lập
# Sử dụng: .\api-tests\scripts\manage_server.ps1 start|stop|status

param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("start", "stop", "status")]
    [string]$Action
)

# Import modules
$scriptDir = if ($PSScriptRoot) { 
    $PSScriptRoot 
} else { 
    Split-Path -Parent $MyInvocation.MyCommand.Path 
}
. "$scriptDir\utils.ps1"
. "$scriptDir\server.ps1"

# Xử lý action
switch ($Action.ToLower()) {
    "start" {
        Start-TestServer
    }
    "stop" {
        Stop-TestServer
    }
    "status" {
        Get-TestServerStatus
    }
}

