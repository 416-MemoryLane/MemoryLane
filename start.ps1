if ($args.Length -ne 2) {
    Write-Host "Usage: $0 <username> <password>"
    exit 1
}

$user = $args[0]
$pass = $args[1]

# get the path of the current script's directory
$script_dir = Split-Path -Parent $MyInvocation.MyCommand.Path
$envFilePath = "$script_dir\.env"

if (Test-Path $envFilePath) {
    $envContents = Get-Content $envFilePath
    $username = $envContents | Select-String "^ML_USERNAME=(.+)$" | ForEach-Object {$_.Matches.Groups[1].Value}
    $password = $envContents | Select-String "^ML_PASSWORD=(.+)$" | ForEach-Object {$_.Matches.Groups[1].Value}
    if ($user -eq $username -and $pass -eq $password) {
        # user match .env file
        Write-Host "Logged in as $user"
    } else {
        # user pass don't match .env file
        Write-Host "Incorrect username and password. Please login as $username. Try again."
        exit 1
    } 
} else {
    # Send login request to galactus
    $url = "https://memory-lane-381119.wl.r.appspot.com/login"
    $body = @{
        username = $user
        password = $pass
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri $url -Method POST -ContentType "application/json" -Body $body
    $message = $response.message
    
    if ($message -eq "$user successfully logged in" -or $message -eq "Account with username $user successfully created") {
        if ($message -eq "Account with username $user successfully created") {
            Write-Host $message
        }
        "ML_USERNAME=$user" | Out-File $envFilePath
        "ML_PASSWORD=$pass" | Out-File $envFilePath -Append
        Write-Host "Logged in as $user"
    } else {
        Write-Host "Incorrect username and password. Please try again."
        exit 1
    }
}

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
    Write-Host "Received termination signal, stopping jobs. This may take a minute..."
    Set-Location $script_dir
    Stop-Job $ui_pid
    Stop-Job $go_pid
    Wait-Job $ui_pid, $go_pid | Out-Null
}
