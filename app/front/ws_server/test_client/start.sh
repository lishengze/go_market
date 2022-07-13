# !/bin/zsh

# !/bin/bash

go build

nohup ./test_client >log/main.log &

sleep 2s

ps -aux|grep test_client
