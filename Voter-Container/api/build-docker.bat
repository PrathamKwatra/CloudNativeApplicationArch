@echo off

docker build --tag voter-api:v1  -f ./dockerfile.dockerfile .
