cd .\ui
yarn
yarn build
Remove-Item -Recurse -Force ..\server\dist
Move-Item .\dist ..\server\dist
cd ..\server
yarn
yarn start
