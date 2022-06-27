# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep front |grep -v grep | awk '{print $2}'`
kill -9 $PID

sleep 2s

ps -aux|grep front