:: Description: Start Redis with Cache Data mapped to data folder
:: For windows, the volume provided to docker should be absolute path (only for dockers set up using WSL2)
:: For docker set up using HyperV, I am not sure if this is required
@echo off
title Start Redis Volume

CALL :NORMALIZEPATH ".\cache-data"

ECHO "Volume Location (abs path): %RETVAL%"

@REM docker run -d --rm --name cnse-redis -e REDIS_ARGS="--save 20 1 --appendonly yes"  -p 6379:6379 -p 8001:8001 -v %RETVAL%:/data  redis/redis-stack:latest 

:NORMALIZEPATH
  SET RETVAL=%~f1
  EXIT /B
