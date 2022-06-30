# !/bin/zsh

# !/bin/bash

echo "args : $@"

go build

nohup ./front $1 >log/main.log &

sleep 2s

ps -aux|grep front

