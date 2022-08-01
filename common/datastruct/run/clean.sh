# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep run |grep -v grep | awk '{print $2}'`
kill -9 $PID

rm -fr run

sleep 2s

ps -aux|grep run