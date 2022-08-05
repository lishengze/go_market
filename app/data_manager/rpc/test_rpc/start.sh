# !/bin/zsh

# !/bin/bash

echo "args : $@"

go build 

nohup ./test_rpc >log/main.log &

ls -a

sleep 2s

ps -aux|grep test_rpc

