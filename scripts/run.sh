#!/bin/sh
set -e

# Get the api key 
apiKey=$1

# Build the service
env GOOS=linux go install -ldflags "-w -X main.apiKey=$apiKey" ./src/app
chmod +x ./bin/linux_amd64/app

# Copy the Dockerfile
cp ./src/Dockerfile ./bin/

cd ./bin

# Docker build
docker build -t messagebird_proxy .

# Docker run!
docker run --expose 8080 -p 8080:8080 -t messagebird_proxy 
