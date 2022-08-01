# !/bin/zsh

# !/bin/bash

echo "args : $@"

go build 

nohup ./test_datastruct >log/main.log &

ls -a

sleep 2s

ps -aux|grep test_datastruct

