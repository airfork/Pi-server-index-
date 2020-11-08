#!/usr/bin/env bash

curl -s "localhost:8000" >/dev/null 2>&1;
errorCode=$?

if [ $errorCode -eq 0 ]; then
   echo "Server is already running"
   exit 1
fi

go build -o server
nohup ./server &> /dev/null &