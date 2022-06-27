# !/bin/zsh

# !/bin/bash

echo "args : $@"


nohup ./front $1 >log/main.log &

sleep 2s

ps -aux|grep front

