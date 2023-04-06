cd .\ui
yarn build
Remove-Item -Recurse -Force ..\server\dist
Move-Item .\dist ..\server\dist
cd ..\server
yarn start
