# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep data_manager |grep -v grep | awk '{print $2}'`
kill -9 $PID

rm -fr data_manager

sleep 2s

ps -aux|grep data_manager