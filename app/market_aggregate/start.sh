# !/bin/zsh

# !/bin/bash

echo "args : $@"


go build


nohup ./market_aggregate $1 $2>log/main.log &


ls -a

sleep 2s

ps -aux|grep market_aggregate

