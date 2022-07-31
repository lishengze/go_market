# !/bin/zsh

# !/bin/bash

echo "args : $@"

go build 

nohup ./rpc $1 $2 >log/main.log &

ls -a

sleep 2s

ps -aux|grep rpc

