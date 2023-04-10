#!/bin/bash

cp ../.env ./server/.env
cd ui
yarn
yarn build
rm -rf ../server/dist
mv dist ../server/dist
cd ../server
yarn
yarn start
