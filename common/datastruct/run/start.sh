# !/bin/zsh

# !/bin/bash

echo "args : $@"

go build 

nohup ./run >log/main.log &

ls -a

sleep 2s

ps -aux|grep run

