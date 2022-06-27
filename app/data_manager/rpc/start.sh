# !/bin/zsh

# !/bin/bash

echo "args : $@"


nohup ./data_manager $1 >log/main.log &

sleep 2s

ps -aux|grep data_manager

