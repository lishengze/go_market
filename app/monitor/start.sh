# !/bin/zsh

# !/bin/bash

echo "args : $@"

go build

nohup ./monitor $1 >log/main.log &

sleep 2s

ps -aux|grep monitor

