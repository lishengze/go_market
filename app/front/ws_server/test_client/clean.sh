# !/bin/bash

rm -fr log/*.log

PID=`ps aux|grep test_client |grep -v grep | awk '{print $2}'`
kill -9 $PID

rm -fr test_client

sleep 2s

ps -aux|grep test_client