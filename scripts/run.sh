#!/bin/sh
set -e

# Build the service
env GOOS=linux go install ./src/app
chmod +x ./bin/linux_amd64/app

# Copy the Dockerfile
cp ./src/Dockerfile ./bin/
cd ./bin

# Docker build
docker build -t messagebird_proxy .

# Docker run
docker run -e API_KEY=Wxsljyqzf0kbikO96mtpyY2xw --expose 8080 -p 8080:8080 -t messagebird_proxy  