# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep market_aggregate |grep -v grep | awk '{print $2}'`
kill -9 $PID

rm -fr market_aggregate

ls -a

sleep 2s

ps -aux|grep market_aggregate