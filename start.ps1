if ($args.Count -ne 2) {
    Write-Host "Usage: script.ps1 <username> <password>"
    exit 1
}

$user = $args[0]
$pass = $args[1]

# run the Go app in the background
Start-Job -ScriptBlock {
    & go mod download
    & go run app.go --username $using:user --password $using:pass
}

# save the PID of the last background process
$go_pid = (Get-Job | Select-Object -Last 1).Id

# run the UI in the background
cd ui
Start-Job -ScriptBlock { & ./start_ui.sh }
$ui_pid = (Get-Job | Select-Object -Last 1).Id

try {
    # tail the UI output in real-time
    Receive-Job -Id $ui_pid -Wait
} finally {
    Write-Host "Received termination signal, stopping jobs."
    cd ..
    Stop-Job $ui_pid -NoWait
    Stop-Job $go_pid -NoWait
    
}
