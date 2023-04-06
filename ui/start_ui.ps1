
Copy-Item -Path "..\.env" -Destination ".\server\.env"
Set-Location -Path ".\ui"
yarn
yarn build
Remove-Item -Path "..\server\dist" -Recurse -Force
Move-Item -Path ".\dist" -Destination "..\server\dist"
Set-Location -Path "..\server"
yarn
yarn start
