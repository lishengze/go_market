# !/bin/zsh

# !/bin/bash

echo "args : $@"


nohup ./rpc $1 >log/main.log &

go build 

ll

sleep 2s

ps -aux|grep rpc

