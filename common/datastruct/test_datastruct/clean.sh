# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep test_datastruct |grep -v grep | awk '{print $2}'`
kill -9 $PID

rm -fr test_datastruct

sleep 2s

ps -aux|grep test_datastruct