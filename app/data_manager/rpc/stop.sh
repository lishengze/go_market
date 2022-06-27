# !/bin/bash

PID=`ps aux|grep data_manager |grep -v grep | awk '{print $2}'`
kill -9 $PID

sleep 2s

ps -aux|grep data_manager