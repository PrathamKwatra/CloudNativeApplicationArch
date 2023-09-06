#!/bin/bash
docker buildx create --use 
docker buildx build --platform linux/amd64,linux/arm64 -f ./dockerfile . -t ninjaversionfive0/voter-api:latest --push
