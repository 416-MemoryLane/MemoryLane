cd .\ui
yarn build
Move-Item .\dist ..\server\dist
cd ..\server
yarn start