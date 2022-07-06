# !/bin/zsh

# !/bin/bash

echo "args : $@"

go build 

nohup ./rpc $1 >log/main.log &

ll

sleep 2s

ps -aux|grep rpc

