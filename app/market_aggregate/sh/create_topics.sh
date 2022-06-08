# !/bin/bash

bin/kafka-topics.sh  --create --topic DEPTH.BTC_USDT._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic DEPTH.ETH_USDT._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic DEPTH.BTC_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic DEPTH.ETH_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic DEPTH.ETH_BTC._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic DEPTH.USDT_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1

bin/kafka-topics.sh  --create --topic TRADE.BTC_USDT._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic TRADE.ETH_USDT._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic TRADE.BTC_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic TRADE.ETH_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic TRADE.ETH_BTC._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic TRADE.USDT_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1

bin/kafka-topics.sh  --create --topic KLINE.BTC_USDT._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic KLINE.ETH_USDT._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic KLINE.BTC_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic KLINE.ETH_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic KLINE.ETH_BTC._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1
bin/kafka-topics.sh  --create --topic KLINE.USDT_USD._bcts_ --bootstrap-server 10.10.1.45:9092  --replication-factor 1 --partitions 1