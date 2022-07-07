# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep monitor |grep -v grep | awk '{print $2}'`
kill -9 $PID

rm -fr monitor

sleep 2s

ps -aux|grep monitor