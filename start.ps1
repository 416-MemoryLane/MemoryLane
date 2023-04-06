if ($args.Count -ne 2) {
    Write-Host "Usage: script.ps1 <username> <password>"
    exit 1
}

$user = $args[0]
$pass = $args[1]

# get the path of the current script's directory
$script_dir = Split-Path -Parent $MyInvocation.MyCommand.Path

# run the Go app in the background
Start-Job -ScriptBlock {
    & go mod download
    & go run app.go --username $using:user --password $using:pass
} | Out-Null

# save the PID of the last background process
$go_pid = (Get-Job | Select-Object -Last 1).Id

# run the UI in the background
Start-Job -ScriptBlock { 
    Set-Location $using:script_dir\ui
    & .\start_ui.ps1 
} | Out-Null
$ui_pid = (Get-Job | Select-Object -Last 1).Id

try {
    # tail the Go app output in real-time
    Receive-Job -Id $go_pid -Wait
} finally {
    Write-Host "Received termination signal, stopping jobs."
    Set-Location $script_dir
    Stop-Job $ui_pid
    Stop-Job $go_pid
    Wait-Job $ui_pid, $go_pid | Out-Null
}
