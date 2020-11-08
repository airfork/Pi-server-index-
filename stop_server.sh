#!/usr/bin/env bash

stopCode=$(cat stopCode.txt)

curl -s "localhost:8000/${stopCode}" >/dev/null 2>&1;
errorCode=$?
if [ $errorCode -eq 7 ]; then
  echo "Server is already stopped"
elif [ $errorCode -eq 52 ]; then
  echo "Successfully stopped the server"
else
  echo "Curl exited with error code: ${errorCode}"
fi