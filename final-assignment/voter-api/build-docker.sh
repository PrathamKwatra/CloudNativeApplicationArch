#!/bin/bash
# docker buildx create --use 
# docker buildx build --platform linux/amd64,linux/arm64 -f ./dockerfile . -t polls-api:latest
docker build -f ./dockerfile . -t voters-api:latest
```