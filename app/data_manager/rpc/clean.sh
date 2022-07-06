# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep rpc |grep -v grep | awk '{print $2}'`
kill -9 $PID

rm -fr rpc

sleep 2s

ps -aux|grep rpc