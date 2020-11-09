#!/usr/bin/env bash

printf "Trying to stop server\n\n"
bash stop_server.sh
printf "Pulling from master\n\n"
git pull origin master
printf "Trying to start server\n\n"
bash start_server.sh
