#!/bin/bash
# docker buildx create --use 
# docker buildx build --platform linux/amd64,linux/arm64 -f ./dockerfile . -t votes-api:latest
docker build -f ./dockerfile . -t votes-api:latest
