# !/bin/zsh

# !/bin/bash

echo "args : $@"


nohup ./market_aggregate $1 >log/main.log &


go build

ll

sleep 2s

ps -aux|grep market_aggregate

