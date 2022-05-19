# !/bin/bash

bin/kafka-topics.sh  --create --topic DEPTH.BTC_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic DEPTH.ETH_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic DEPTH.ETH_BTC._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic DEPTH.USDT_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1

bin/kafka-topics.sh  --create --topic TRADE.BTC_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic TRADE.ETH_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic TRADE.ETH_BTC._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic TRADE.USDT_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1

bin/kafka-topics.sh  --create --topic KLINE.BTC_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic KLINE.ETH_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic KLINE.ETH_BTC._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic KLINE.USDT_USD._bcts_ --bootstrap-server 152.32.254.76:9117  --replication-factor 1 --partitions 1