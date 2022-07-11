# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep front |grep -v grep | awk '{print $2}'`
kill -9 $PID

rm -fr front

sleep 2s

ps -aux|grep front