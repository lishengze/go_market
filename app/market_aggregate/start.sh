# !/bin/zsh

# !/bin/bash

echo "args : $@"


go build


nohup ./market_aggregate $1 >log/main.log &


ll

sleep 2s

ps -aux|grep market_aggregate

