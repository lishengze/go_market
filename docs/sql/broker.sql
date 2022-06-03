
DROP TABLE IF EXISTS `account_settlements`;
CREATE TABLE `account_settlements` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` varchar(16) DEFAULT NULL COMMENT '会员ID, 系统管理的ID为admin',
  `currency` varchar(16) DEFAULT NULL COMMENT '交易币种',
  `bank_balance` decimal(26,16) DEFAULT NULL COMMENT '当日银行/钱包余额',
  `bank_withdraw` decimal(26,16) DEFAULT NULL COMMENT '当日银行/钱包充值',
  `bank_deposit` decimal(26,16) DEFAULT NULL COMMENT '当日银行/钱包提现',
  `cold_wallet_amount` decimal(26,16) DEFAULT NULL COMMENT '冷钱包余额',
  `total_member_amount` decimal(26,16) DEFAULT NULL COMMENT '内部账户余额',
  `recharge_amount` decimal(26,16) DEFAULT NULL COMMENT '充值金额',
  `withdraw_amount` decimal(26,16) DEFAULT NULL COMMENT '提现金额',
  `withdraw_fee` decimal(26,16) DEFAULT NULL COMMENT '提现手续费',
  `trade_fee` decimal(26,16) DEFAULT NULL COMMENT '交易手续费',
  `status` tinyint(1) DEFAULT '0' COMMENT '结算状态 0 待结算, 1 结算中, 2 结算完成',
  `settle_date` date DEFAULT NULL COMMENT '结算日期',
  `operator` varchar(16) DEFAULT NULL COMMENT '操作员',
  `create_time` datetime DEFAULT NULL COMMENT '记录创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '记录修改时间',
  `bank_deposite` decimal(26,16) DEFAULT NULL COMMENT '''当日银行/钱包 deposite''',
  `wallet_adjustment` decimal(26,16) DEFAULT NULL COMMENT '当日银行/钱包调整',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8221 DEFAULT CHARSET=utf8;

#
# Structure for table "account_snapshot"
#

DROP TABLE IF EXISTS `account_snapshot`;
CREATE TABLE `account_snapshot` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `account_id` bigint(18) DEFAULT NULL COMMENT '账号ID',
  `member_id` bigint(18) DEFAULT NULL COMMENT '会员ID',
  `coin_id` bigint(18) DEFAULT NULL COMMENT '币种ID',
  `coin_unit` varchar(16) DEFAULT NULL COMMENT '币种',
  `balance` decimal(26,16) DEFAULT NULL COMMENT '可用余额',
  `frozen_balance` decimal(26,16) DEFAULT NULL COMMENT '冻结余额',
  `today_deposit` decimal(26,16) DEFAULT NULL COMMENT '当日充值金额',
  `today_withdraw` decimal(26,16) DEFAULT NULL COMMENT '当日提现金额',
  `today_deal` decimal(26,16) DEFAULT NULL COMMENT '当日交易金额',
  `today_date` date DEFAULT NULL COMMENT '日期',
  `valuation` decimal(26,16) DEFAULT NULL COMMENT '估值',
  `valuation_currency` decimal(26,16) DEFAULT NULL COMMENT '估值币价',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `total_balance` decimal(26,16) DEFAULT NULL COMMENT '总额',
  `today_withdrawfee` decimal(26,16) DEFAULT NULL COMMENT '当日提现手续费',
  `today_dealfee` decimal(26,16) DEFAULT NULL COMMENT '当日交易手续费',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_member_id` (`member_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4211261 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='账户资产快照表';

#
# Structure for table "account_snapshotplat"
#

DROP TABLE IF EXISTS `account_snapshotplat`;
CREATE TABLE `account_snapshotplat` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `account_id` bigint(18) DEFAULT NULL COMMENT '账号ID',
  `member_id` bigint(18) DEFAULT NULL COMMENT '会员ID',
  `coin_id` bigint(18) DEFAULT NULL COMMENT '币种ID',
  `coin_unit` varchar(16) DEFAULT NULL COMMENT '币种',
  `balance` decimal(26,16) DEFAULT NULL COMMENT '可用余额',
  `frozen_balance` decimal(26,16) DEFAULT NULL COMMENT '冻结余额',
  `today_deposit` decimal(26,16) DEFAULT NULL COMMENT '当日充值金额',
  `today_withdraw` decimal(26,16) DEFAULT NULL COMMENT '当日提现金额',
  `today_deal` decimal(26,16) DEFAULT NULL COMMENT '当日交易金额',
  `today_date` date DEFAULT NULL COMMENT '日期',
  `valuation` decimal(26,16) DEFAULT NULL COMMENT '估值',
  `valuation_currency` decimal(26,16) DEFAULT NULL COMMENT '估值币价',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `total_balance` decimal(26,16) DEFAULT NULL COMMENT '总额',
  `today_withdrawfee` decimal(26,16) DEFAULT NULL COMMENT '当日提现手续费',
  `today_dealfee` decimal(26,16) DEFAULT NULL COMMENT '当日交易手续费',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_member_id` (`member_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2878619 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='账户资产快照表';

#
# Structure for table "accounts"
#

DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` bigint(18) DEFAULT NULL COMMENT '会员ID',
  `coin_id` bigint(18) DEFAULT NULL COMMENT '币种ID',
  `coin_unit` varchar(16) DEFAULT NULL COMMENT '币种',
  `balance` decimal(26,16) DEFAULT NULL COMMENT '可用余额',
  `frozen_balance` decimal(26,16) DEFAULT NULL COMMENT '冻结余额',
  `to_released` decimal(26,16) DEFAULT NULL COMMENT '待释放总量',
  `address` varchar(64) DEFAULT NULL COMMENT '充值地址',
  `version` int(2) DEFAULT NULL COMMENT '版本',
  `is_lock` tinyint(1) DEFAULT NULL COMMENT '账户是否锁定，0否，1是',
  `average_price` decimal(26,16) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_member_coin` (`member_id`,`coin_unit`)
) ENGINE=InnoDB AUTO_INCREMENT=164037 DEFAULT CHARSET=utf8;

#
# Structure for table "admin_member_group"
#

DROP TABLE IF EXISTS `admin_member_group`;
CREATE TABLE `admin_member_group` (
  `group_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `group_code` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `memo` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`group_id`),
  UNIQUE KEY `uix_member_group_group_code` (`group_code`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8;

#
# Structure for table "admin_member_news"
#

DROP TABLE IF EXISTS `admin_member_news`;
CREATE TABLE `admin_member_news` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` varchar(255) NOT NULL,
  `news_id` int(11) NOT NULL,
  `is_read` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uix_member_news` (`member_id`,`news_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3127 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

#
# Structure for table "admin_news"
#

DROP TABLE IF EXISTS `admin_news`;
CREATE TABLE `admin_news` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `content` varchar(1024) DEFAULT NULL,
  `publisher` varchar(255) NOT NULL,
  `is_to_all` tinyint(1) NOT NULL,
  `to_member_ids` varchar(2048) DEFAULT NULL,
  `fail_member_ids` varchar(2048) DEFAULT NULL,
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

#
# Structure for table "admin_user"
#

DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user` (
  `user_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role` int(11) DEFAULT NULL,
  `is_google_authenticated` tinyint(1) NOT NULL,
  `name` varchar(255) NOT NULL,
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `department` varchar(50) NOT NULL DEFAULT '' COMMENT '部门',
  `login_password` varchar(255) NOT NULL,
  `google_auth_secret` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '0-待生效，1-生效中，2-冻结',
  `last_operator` bigint(20) NOT NULL DEFAULT '0' COMMENT ' 最后一次操作人,0-表示由程序自动化创建',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`) USING BTREE,
  UNIQUE KEY `uix_user_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

#
# Structure for table "admin_withdraw_approve"
#

DROP TABLE IF EXISTS `admin_withdraw_approve`;
CREATE TABLE `admin_withdraw_approve` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `withdraw_no` bigint(20) NOT NULL,
  `action` int(11) NOT NULL,
  `approve_note` varchar(255) DEFAULT NULL,
  `approve_user` varchar(255) DEFAULT NULL,
  `updated` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uix_withdraw_approve_withdraw_no` (`withdraw_no`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

#
# Structure for table "assets"
#

DROP TABLE IF EXISTS `assets`;
CREATE TABLE `assets` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `asset_type` tinyint(1) DEFAULT NULL COMMENT '资产状态',
  `request_id` varchar(32) DEFAULT NULL COMMENT '请求编号',
  `member_id` bigint(18) DEFAULT NULL COMMENT '会员ID',
  `coin_id` bigint(18) DEFAULT NULL COMMENT '币种ID',
  `coin_unit` varchar(16) DEFAULT NULL COMMENT '币种',
  `coin_spec` varchar(16) DEFAULT NULL COMMENT '具体币种',
  `to` varchar(64) DEFAULT NULL COMMENT '提现接收地址',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '提现数量',
  `memo` varchar(16) DEFAULT NULL COMMENT '备忘录',
  `status` varchar(15) DEFAULT NULL COMMENT '状态',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `fee` decimal(26,16) DEFAULT NULL COMMENT '提现手续费',
  PRIMARY KEY (`id`),
  KEY `idx_coin` (`coin_unit`),
  KEY `idx_member_id` (`member_id`),
  KEY `idx_coin_id` (`coin_id`)
) ENGINE=InnoDB AUTO_INCREMENT=250592532505165825 DEFAULT CHARSET=utf8;

#
# Structure for table "coin_prices"
#

DROP TABLE IF EXISTS `coin_prices`;
CREATE TABLE `coin_prices` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `coin_unit` varchar(255) DEFAULT NULL COMMENT '币种',
  `price` decimal(26,16) DEFAULT NULL COMMENT '日出价格',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `date_time` varchar(255) CHARACTER SET eucjpms DEFAULT NULL COMMENT '日期',
  PRIMARY KEY (`id`),
  KEY `coin_date` (`date_time`,`coin_unit`)
) ENGINE=InnoDB AUTO_INCREMENT=3723 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "coins"
#

DROP TABLE IF EXISTS `coins`;
CREATE TABLE `coins` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(16) DEFAULT NULL COMMENT '货币',
  `name_cn` varchar(32) DEFAULT NULL COMMENT '中文名称',
  `unit` varchar(20) DEFAULT NULL COMMENT '单位',
  `status` varchar(8) DEFAULT NULL COMMENT '状态',
  `min_tx_fee` decimal(26,16) DEFAULT NULL COMMENT '最小提币手续费',
  `cny_rate` decimal(26,16) DEFAULT '0.0000000000000000' COMMENT '人民币汇率',
  `max_tx_fee` decimal(26,16) DEFAULT NULL COMMENT '最大提币手续费',
  `usd_rate` decimal(26,16) DEFAULT '0.0000000000000000' COMMENT '对美元汇率',
  `sort` tinyint(1) DEFAULT NULL COMMENT '排序',
  `can_withdraw` tinyint(1) DEFAULT NULL COMMENT '是否能提币',
  `can_recharge` tinyint(1) DEFAULT NULL COMMENT '是否能充币',
  `can_transfer` tinyint(1) DEFAULT NULL COMMENT '是否能转账',
  `can_auto_withdraw` tinyint(1) DEFAULT NULL COMMENT '是否能自动提币',
  `withdraw_threshold` decimal(26,16) DEFAULT NULL COMMENT '提币阈值',
  `min_withdraw_amount` decimal(26,16) DEFAULT NULL COMMENT '最小提币数量',
  `max_withdraw_amount` decimal(26,16) DEFAULT NULL COMMENT '最大提币数量',
  `min_recharge_amount` decimal(26,16) DEFAULT NULL COMMENT '最小充值数量',
  `is_legal` tinyint(1) DEFAULT NULL COMMENT '是否是合法币种',
  `all_balance` decimal(26,16) DEFAULT NULL COMMENT '总余额',
  `cold_wallet_address` varchar(64) DEFAULT NULL COMMENT '冷钱包地址',
  `hot_all_balance` decimal(26,16) DEFAULT NULL,
  `block_height` bigint(20) DEFAULT NULL COMMENT '块高',
  `miner_fee` decimal(26,16) DEFAULT NULL COMMENT '转账时付给矿工的手续费',
  `withdraw_scale` int(2) DEFAULT '4' COMMENT '提币精度',
  `info_link` varchar(64) DEFAULT NULL COMMENT '币种资料链接',
  `description` text COMMENT '币种简介',
  `account_type` tinyint(1) DEFAULT NULL COMMENT '账户类型',
  `deposit_address` varchar(64) DEFAULT NULL COMMENT '充值地址',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_coins_name` (`name`),
  UNIQUE KEY `idx_coins_unit` (`unit`)
) ENGINE=InnoDB AUTO_INCREMENT=54 DEFAULT CHARSET=utf8;

#
# Structure for table "content_kv"
#

DROP TABLE IF EXISTS `content_kv`;
CREATE TABLE `content_kv` (
  `id` bigint(18) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` bigint(18) DEFAULT NULL COMMENT '用户ID',
  `key` varchar(120) DEFAULT NULL COMMENT '键名称',
  `content` varchar(1024) DEFAULT NULL COMMENT '内容',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8;

#
# Structure for table "exchange_coins"
#

DROP TABLE IF EXISTS `exchange_coins`;
CREATE TABLE `exchange_coins` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `symbol` varchar(16) DEFAULT NULL COMMENT '交易币种名称',
  `symbol_no` varchar(6) DEFAULT NULL COMMENT '交易品种编号',
  `coin_symbol` varchar(16) DEFAULT NULL COMMENT '交易币种符号',
  `base_symbol` varchar(16) DEFAULT NULL COMMENT '结算币种符号',
  `enable` int(2) DEFAULT NULL COMMENT '状态，1：启用，2：禁止',
  `maker_fee` decimal(26,16) DEFAULT NULL COMMENT 'Maker交易手续费',
  `taker_fee` decimal(26,16) DEFAULT NULL COMMENT 'Taker交易手续费',
  `spread` decimal(26,16) DEFAULT NULL COMMENT '点差',
  `sort` tinyint(1) DEFAULT NULL COMMENT '排序',
  `coin_scale` int(6) DEFAULT NULL COMMENT '交易币小数精度',
  `base_coin_scale` int(6) DEFAULT NULL COMMENT '基币小数精度',
  `min_sell_price` decimal(26,16) DEFAULT NULL COMMENT '卖单最低价格',
  `max_buy_price` decimal(26,16) DEFAULT NULL COMMENT '最高买单价',
  `enable_market_sell` tinyint(1) DEFAULT NULL COMMENT '是否启用市价卖',
  `enable_market_buy` tinyint(1) DEFAULT NULL COMMENT '是否启用市价买',
  `max_trading_time` bigint(18) DEFAULT NULL COMMENT '最大交易时间',
  `max_trading_order` int(8) DEFAULT NULL COMMENT '最大在交易中的委托数量',
  `robot_type` tinyint(1) DEFAULT NULL COMMENT '机器人类型',
  `flag` tinyint(1) DEFAULT NULL COMMENT '标签位',
  `min_turnover` decimal(26,16) DEFAULT NULL COMMENT '最小成交额',
  `zone` int(8) DEFAULT NULL COMMENT '交易区域',
  `min_volume` decimal(26,16) DEFAULT NULL COMMENT '最小下单量',
  `max_volume` decimal(26,16) DEFAULT NULL COMMENT '最大下单量',
  `publish_type` int(11) DEFAULT '1' COMMENT '发行活动类型',
  `start_time` varchar(30) DEFAULT '2000-01-01 01:00:00' COMMENT '活动开始时间',
  `end_time` varchar(30) DEFAULT '2000-01-01 01:00:00' COMMENT '活动结束时间',
  `clear_time` varchar(30) DEFAULT '2000-01-01 01:00:00' COMMENT '活动清盘时间',
  `publish_price` decimal(26,16) DEFAULT '0.0000000000000000' COMMENT '分摊发行价格',
  `publish_amount` decimal(26,16) DEFAULT '0.0000000000000000' COMMENT '活动发行数量',
  `visible` int(11) DEFAULT '1' COMMENT '前台可见状态，1：可见，2：不可见',
  `exchangeable` int(11) DEFAULT '1' COMMENT '是否可交易，1：可交易，2：不可交易',
  `current_time` bigint(20) DEFAULT NULL COMMENT '服务器当前市价戳',
  `engine_status` tinyint(1) DEFAULT NULL COMMENT '交易引擎状态（0：不可用，1：可用)',
  `market_engine_status` tinyint(1) DEFAULT NULL COMMENT '行情引擎状态（0：不可用，1：可用',
  `ex_engine_status` tinyint(1) DEFAULT NULL COMMENT '交易机器人状态（0：非运行中，1：运行中）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_exchange_coins_symbol` (`symbol`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

#
# Structure for table "exchange_hedges"
#

DROP TABLE IF EXISTS `exchange_hedges`;
CREATE TABLE `exchange_hedges` (
  `id` bigint(18) NOT NULL COMMENT '主键ID',
  `parent_id` bigint(18) DEFAULT NULL COMMENT '对冲母单ID',
  `trade_id` bigint(18) DEFAULT NULL COMMENT '交易ID',
  `exchange_id` varchar(16) DEFAULT NULL COMMENT '对冲平台',
  `account_name` varchar(32) DEFAULT NULL COMMENT '交易所账户',
  `symbol` varchar(12) DEFAULT NULL COMMENT '交易币种名称，格式：BTC_USDT',
  `coin_symbol` varchar(16) DEFAULT NULL COMMENT '交易币单位',
  `base_symbol` varchar(16) DEFAULT NULL COMMENT '结算单位',
  `price` decimal(26,16) DEFAULT NULL COMMENT '内部交易价格/委托价格',
  `market_price` decimal(26,16) DEFAULT NULL COMMENT '行情价格',
  `match_price` decimal(26,16) DEFAULT NULL COMMENT '交易所成交价格',
  `volume` decimal(26,16) DEFAULT NULL COMMENT '发起对冲量',
  `match_volume` decimal(26,16) DEFAULT NULL COMMENT '交易所成交量',
  `match_value` decimal(26,16) DEFAULT NULL COMMENT '交易所成交额',
  `trade_no` varchar(32) DEFAULT NULL COMMENT '交易所成交编号',
  `fee` decimal(18,10) DEFAULT NULL COMMENT '交易所成交手续费',
  `fee_currency` varchar(16) DEFAULT NULL COMMENT '交易所产生手续费的币种',
  `buy_turnover` decimal(26,16) DEFAULT NULL COMMENT '买入成交额',
  `sell_turnover` decimal(26,16) DEFAULT NULL COMMENT '卖出成交额',
  `direction` varchar(5) DEFAULT NULL COMMENT '交易方向',
  `asset_type` varchar(16) DEFAULT NULL COMMENT '交易所资产类别',
  `hedge_type` varchar(16) DEFAULT NULL COMMENT '对冲交易类型',
  `status` varchar(8) DEFAULT NULL COMMENT '委托状态',
  `trade_channel` varchar(16) DEFAULT NULL COMMENT '交易通道',
  `time` bigint(18) DEFAULT NULL COMMENT '成交时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '交易发起时间',
  `platform_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '交易更新时间',
  `trade_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '交易所成交时间',
  `rsp_local_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '成交回到柜台时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '交易更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_exchange_hedges_parent_id` (`parent_id`),
  KEY `idx_exchange_hedges_trade_id` (`trade_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#
# Structure for table "exchange_order_details"
#

DROP TABLE IF EXISTS `exchange_order_details`;
CREATE TABLE `exchange_order_details` (
  `id` bigint(18) NOT NULL COMMENT '主键ID',
  `order_id` varchar(20) DEFAULT NULL COMMENT '订单ID',
  `price` decimal(26,16) DEFAULT NULL COMMENT '挂单价格',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '成交量',
  `turnover` decimal(26,16) DEFAULT NULL COMMENT '成交额',
  `fee` decimal(26,16) DEFAULT NULL COMMENT '手续费',
  `time` bigint(18) DEFAULT NULL COMMENT '时间戳',
  `detail_type` tinyint(1) DEFAULT NULL COMMENT '订单明细类型',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `profit` decimal(26,16) DEFAULT NULL COMMENT '收益',
  `member_id` bigint(20) DEFAULT NULL COMMENT '会员id',
  `symbol` varchar(50) DEFAULT NULL COMMENT '投资品平仓收益',
  KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#
# Structure for table "exchange_orders"
#

DROP TABLE IF EXISTS `exchange_orders`;
CREATE TABLE `exchange_orders` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_id` varchar(20) DEFAULT NULL COMMENT '订单ID',
  `member_id` bigint(18) DEFAULT NULL COMMENT '会员ID',
  `exchange_order_type` tinyint(1) DEFAULT NULL COMMENT '挂单类型',
  `order_type` varchar(16) DEFAULT NULL COMMENT '订单类型',
  `match_id` varchar(20) DEFAULT NULL COMMENT '匹配订单ID',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '买入或卖出量',
  `symbol` varchar(16) DEFAULT NULL COMMENT '交易对符号',
  `traded_amount` decimal(26,16) DEFAULT NULL COMMENT '成交量',
  `turnover` decimal(26,16) DEFAULT NULL COMMENT '成交额',
  `coin_symbol` varchar(16) DEFAULT NULL COMMENT '交易币单位',
  `base_symbol` varchar(16) DEFAULT NULL COMMENT '结算单位',
  `status` tinyint(1) DEFAULT NULL COMMENT '订单状态',
  `direction` tinyint(1) DEFAULT NULL COMMENT '订单方向',
  `price` decimal(26,16) DEFAULT NULL COMMENT '挂单价格',
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '挂单时间',
  `completed_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '交易完成时间',
  `canceled_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '取消时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `use_discount` tinyint(1) DEFAULT NULL COMMENT '是否使用折扣 0 不使用 1使用',
  `estimate_fee` decimal(26,16) DEFAULT NULL COMMENT '预估手续费',
  `fee_kind` tinyint(1) DEFAULT NULL COMMENT '费用类型',
  `fee_rate` decimal(26,16) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_exchange_orders_order_id` (`order_id`),
  KEY `idx_member` (`member_id`),
  KEY `idx_match` (`match_id`)
) ENGINE=InnoDB AUTO_INCREMENT=347 DEFAULT CHARSET=utf8;

#
# Structure for table "exchange_trades"
#

DROP TABLE IF EXISTS `exchange_trades`;
CREATE TABLE `exchange_trades` (
  `id` bigint(18) NOT NULL COMMENT '主键ID',
  `symbol` varchar(12) DEFAULT NULL COMMENT '交易币种名称，格式：BTC_USDT',
  `coin_symbol` varchar(16) DEFAULT NULL COMMENT '交易币单位',
  `base_symbol` varchar(16) DEFAULT NULL COMMENT '结算单位',
  `trade_type` tinyint(1) DEFAULT NULL COMMENT '成交类型',
  `offline_trade_id` varchar(32) DEFAULT NULL COMMENT '线下交易ID',
  `price` decimal(26,16) DEFAULT NULL COMMENT '交易价格',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '成交量',
  `buy_turnover` decimal(26,16) DEFAULT NULL COMMENT '买入成交额',
  `buy_fee` decimal(26,16) DEFAULT NULL COMMENT '买入手续费',
  `sell_turnover` decimal(26,16) DEFAULT NULL COMMENT '卖出成交额',
  `sell_fee` decimal(26,16) DEFAULT NULL COMMENT '卖出手续费',
  `direction` varchar(5) DEFAULT NULL COMMENT '交易方向',
  `buy_order_id` varchar(32) DEFAULT NULL COMMENT '买入订单ID',
  `sell_order_id` varchar(32) DEFAULT NULL COMMENT '卖出订单ID',
  `status` tinyint(1) DEFAULT NULL COMMENT '交易状态',
  `time` bigint(18) DEFAULT NULL COMMENT '成交时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '交易更新时间',
  `need_hedge_amount` decimal(26,16) DEFAULT NULL COMMENT '待对冲量',
  `member_id` bigint(18) DEFAULT NULL COMMENT '用户ID',
  `coin_price` decimal(26,16) DEFAULT NULL COMMENT 'coin当前价格',
  `base_price` decimal(26,16) DEFAULT NULL COMMENT 'base当前价格',
  `buy_time` timestamp NULL DEFAULT NULL COMMENT '买订单时间',
  `sell_time` timestamp NULL DEFAULT NULL COMMENT '卖订单时间',
  PRIMARY KEY (`id`),
  KEY `idx_symbol` (`symbol`),
  KEY `idx_exchange_trades_sell_order_id` (`sell_order_id`),
  KEY `idx_exchange_trades_buy_order_id` (`buy_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#
# Structure for table "fiat_recharge"
#

DROP TABLE IF EXISTS `fiat_recharge`;
CREATE TABLE `fiat_recharge` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `recharge_no` varchar(32) NOT NULL COMMENT '申请编号',
  `uid` varchar(32) DEFAULT NULL COMMENT 'uid',
  `email` varchar(32) DEFAULT NULL COMMENT '客户邮箱',
  `name` varchar(32) DEFAULT NULL COMMENT '姓名',
  `account_type` tinyint(1) DEFAULT NULL COMMENT '1个人 2机构',
  `fiat_type` tinyint(1) DEFAULT NULL COMMENT '1入金 2出金',
  `currency` varchar(32) DEFAULT NULL COMMENT '出入金币种',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '出入金数量',
  `fee` decimal(26,16) DEFAULT NULL COMMENT '出入金手续费',
  `bank_name` varchar(32) DEFAULT NULL COMMENT '出入金银行名',
  `bank_address` varchar(64) DEFAULT NULL COMMENT '银行地址',
  `beneficiary_bank_account_no` varchar(64) DEFAULT NULL COMMENT '受益人银行账号',
  `beneficiary_name_on_account` varchar(64) DEFAULT NULL COMMENT '受益人账户名称',
  `swift` varchar(16) DEFAULT NULL COMMENT 'swift',
  `route` varchar(64) DEFAULT NULL COMMENT '路由码',
  `beneficiary_address_on_account` varchar(64) DEFAULT NULL COMMENT '受益人银行地址',
  `remarks` varchar(64) DEFAULT NULL COMMENT '出入金备注',
  `voucher_photo` varchar(1280) DEFAULT NULL COMMENT '出入金凭证',
  `operator` varchar(32) DEFAULT NULL COMMENT '操作员',
  `approver` varchar(32) DEFAULT NULL COMMENT '审核员',
  `approve_note` varchar(512) DEFAULT NULL COMMENT '审核意见',
  `approve_status` tinyint(1) DEFAULT NULL COMMENT '审核状态 1待审核 2审核通过 3审核拒绝',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `txid` varchar(128) DEFAULT '' COMMENT '提现交易编号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_fiat_recharge_recharge_no` (`recharge_no`)
) ENGINE=InnoDB AUTO_INCREMENT=632 DEFAULT CHARSET=utf8;

#
# Structure for table "fiat_settlement_info"
#

DROP TABLE IF EXISTS `fiat_settlement_info`;
CREATE TABLE `fiat_settlement_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `fiat_settle_no` varchar(32) NOT NULL COMMENT '申请编号',
  `uid` bigint(18) DEFAULT NULL COMMENT 'uid',
  `user_name` varchar(32) DEFAULT NULL,
  `is_fiat` tinyint(1) DEFAULT NULL,
  `bank_name` varchar(255) DEFAULT NULL,
  `bank_address` varchar(255) DEFAULT NULL,
  `swift` varchar(255) DEFAULT NULL,
  `route_code` varchar(64) DEFAULT NULL COMMENT '路由码',
  `beneficiary_name_on_account` varchar(255) DEFAULT NULL,
  `beneficiary_bank_account_no` varchar(255) DEFAULT NULL,
  `beneficiary_address_on_account` varchar(255) DEFAULT NULL,
  `remarks` varchar(255) DEFAULT NULL,
  `source` tinyint(1) DEFAULT NULL COMMENT '是否kyc同步',
  `approve_status` tinyint(1) DEFAULT NULL,
  `approve_note` varchar(512) DEFAULT NULL COMMENT '审核意见',
  `operator` varchar(32) DEFAULT NULL COMMENT '操作员',
  `approver` varchar(32) DEFAULT NULL COMMENT '审核员',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_fiat_settlement_info_fiat_settle_no` (`fiat_settle_no`)
) ENGINE=InnoDB AUTO_INCREMENT=147 DEFAULT CHARSET=utf8;

#
# Structure for table "fiat_settlement_info_test"
#

DROP TABLE IF EXISTS `fiat_settlement_info_test`;
CREATE TABLE `fiat_settlement_info_test` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `fiat_settle_no` varchar(32) NOT NULL COMMENT '申请编号',
  `uid` bigint(18) DEFAULT NULL COMMENT 'uid',
  `user_name` varchar(32) DEFAULT NULL,
  `is_fiat` tinyint(1) DEFAULT NULL,
  `bank_name` varchar(255) DEFAULT NULL,
  `bank_address` varchar(255) DEFAULT NULL,
  `swift` varchar(255) DEFAULT NULL,
  `route_code` varchar(64) DEFAULT NULL COMMENT '路由码',
  `beneficiary_name_on_account` varchar(255) DEFAULT NULL,
  `beneficiary_bank_account_no` varchar(255) DEFAULT NULL,
  `beneficiary_address_on_account` varchar(255) DEFAULT NULL,
  `remarks` varchar(255) DEFAULT NULL,
  `source` tinyint(1) DEFAULT NULL COMMENT '是否kyc同步',
  `approve_status` tinyint(1) DEFAULT NULL,
  `approve_note` varchar(512) DEFAULT NULL COMMENT '审核意见',
  `operator` varchar(32) DEFAULT NULL COMMENT '操作员',
  `approver` varchar(32) DEFAULT NULL COMMENT '审核员',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_fiat_settlement_info_fiat_settle_no` (`fiat_settle_no`)
) ENGINE=InnoDB AUTO_INCREMENT=108 DEFAULT CHARSET=utf8;

#
# Structure for table "hedge_instructs"
#

DROP TABLE IF EXISTS `hedge_instructs`;
CREATE TABLE `hedge_instructs` (
  `id` bigint(18) NOT NULL AUTO_INCREMENT,
  `exchange_name` varchar(16) DEFAULT NULL,
  `symbol` varchar(16) DEFAULT NULL,
  `direction` int(2) DEFAULT NULL,
  `price` decimal(26,16) DEFAULT NULL,
  `amount` decimal(26,16) DEFAULT NULL,
  `hedge_batch_id` bigint(18) DEFAULT NULL,
  `hedge_id` bigint(18) DEFAULT NULL,
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `hedge_id` (`hedge_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7495 DEFAULT CHARSET=utf8;

#
# Structure for table "hedge_matches"
#

DROP TABLE IF EXISTS `hedge_matches`;
CREATE TABLE `hedge_matches` (
  `id` bigint(18) NOT NULL AUTO_INCREMENT,
  `match_id` varchar(100) DEFAULT NULL,
  `order_id` varchar(100) DEFAULT NULL,
  `exchange_name` varchar(16) DEFAULT NULL,
  `account_name` varchar(32) DEFAULT NULL,
  `symbol` varchar(16) DEFAULT NULL,
  `direction` int(2) DEFAULT NULL,
  `trade_price` decimal(26,16) DEFAULT NULL,
  `trade_amount` decimal(26,16) DEFAULT NULL,
  `trade_value` decimal(26,16) DEFAULT NULL,
  `hedge_id` bigint(18) DEFAULT NULL,
  `match_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `amount` decimal(26,16) DEFAULT NULL,
  `fee_currency` varchar(16) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3579947139 DEFAULT CHARSET=utf8;

#
# Structure for table "hedge_orders"
#

DROP TABLE IF EXISTS `hedge_orders`;
CREATE TABLE `hedge_orders` (
  `id` bigint(18) NOT NULL AUTO_INCREMENT,
  `order_id` varchar(32) DEFAULT NULL,
  `exchange_name` varchar(16) DEFAULT NULL,
  `account_name` varchar(32) DEFAULT NULL,
  `symbol` varchar(16) DEFAULT NULL,
  `direction` int(2) DEFAULT NULL,
  `price` decimal(26,16) DEFAULT NULL,
  `amount` decimal(26,16) DEFAULT NULL,
  `trade_price` decimal(26,16) DEFAULT NULL,
  `trade_amount` decimal(26,16) DEFAULT NULL,
  `remain_amount` decimal(26,16) DEFAULT NULL,
  `hedge_id` bigint(18) DEFAULT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `order_status` varchar(16) DEFAULT NULL,
  `fee` decimal(26,16) DEFAULT NULL,
  `fee_currency` varchar(16) DEFAULT NULL,
  `error_id` int(8) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7480 DEFAULT CHARSET=utf8;

#
# Structure for table "hedge_to_centers"
#

DROP TABLE IF EXISTS `hedge_to_centers`;
CREATE TABLE `hedge_to_centers` (
  `id` bigint(18) NOT NULL AUTO_INCREMENT,
  `hedge_match_id` bigint(18) DEFAULT NULL,
  `central_match_id` bigint(18) DEFAULT NULL,
  `hedge_amount` decimal(26,16) DEFAULT NULL,
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=209 DEFAULT CHARSET=utf8;

#
# Structure for table "investment_snapshots"
#

DROP TABLE IF EXISTS `investment_snapshots`;
CREATE TABLE `investment_snapshots` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) DEFAULT NULL COMMENT '会员id',
  `coin_unit` varchar(255) DEFAULT NULL COMMENT '币种',
  `position` decimal(26,16) DEFAULT NULL COMMENT '投资品当前持有数量',
  `average_price` decimal(26,16) DEFAULT NULL COMMENT '均价',
  `cost_basis` decimal(26,16) DEFAULT NULL COMMENT '成本',
  `last` decimal(26,16) DEFAULT NULL COMMENT '市价',
  `market_value` decimal(26,16) DEFAULT NULL COMMENT '市值',
  `market_value_percent` decimal(26,16) DEFAULT NULL COMMENT '市值占比',
  `daily_change` decimal(26,16) DEFAULT NULL COMMENT '日波动',
  `daily_profit` decimal(26,16) DEFAULT NULL COMMENT '日盈亏',
  `unrealized_profit` decimal(26,16) DEFAULT NULL COMMENT '未实现盈亏',
  `unrealized_profit_percent` decimal(26,16) DEFAULT NULL COMMENT '未实现盈亏率',
  `earlier_price` decimal(26,16) DEFAULT NULL COMMENT '日初价格',
  `date_time` date DEFAULT NULL COMMENT '日期',
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `daily_profit_loss` decimal(26,16) DEFAULT NULL COMMENT '日盈亏dailyP&L',
  `realized_profit_loss` decimal(26,16) DEFAULT NULL COMMENT '累计已实现盈亏Realized P&L',
  PRIMARY KEY (`id`),
  KEY `mid_coin` (`member_id`,`coin_unit`,`date_time`)
) ENGINE=InnoDB AUTO_INCREMENT=224175 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "kyc_organ"
#

DROP TABLE IF EXISTS `kyc_organ`;
CREATE TABLE `kyc_organ` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` bigint(18) DEFAULT NULL COMMENT '会员ID 对应 Member表的ID',
  `last_update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  `name` varchar(128) DEFAULT NULL COMMENT '公司法定全称',
  `form` varchar(128) DEFAULT NULL COMMENT '法律注册形式',
  `country` varchar(32) DEFAULT NULL COMMENT '注册国家（地区）',
  `register_date` varchar(32) DEFAULT NULL COMMENT '注册成立日期',
  `company_number` varchar(32) DEFAULT NULL COMMENT '公司注册编号（或同等）/商业登记证（香港）',
  `address` varchar(256) DEFAULT NULL COMMENT '公司注册地址',
  `business_country` varchar(32) DEFAULT NULL COMMENT '主要经营国家（地区）',
  `business_nature` varchar(128) DEFAULT NULL COMMENT '业务性质',
  `email` varchar(32) DEFAULT NULL COMMENT '邮箱',
  `mobile_pre` varchar(16) DEFAULT NULL COMMENT '手机号',
  `mobile` varchar(16) DEFAULT NULL COMMENT '手机号',
  `fund_source` tinyint(1) DEFAULT NULL COMMENT '资金来源',
  `fund_source_str` varchar(256) DEFAULT NULL COMMENT '资金来源',
  `asset_source` tinyint(1) DEFAULT NULL COMMENT '资产来源',
  `asset_source_str` varchar(256) DEFAULT NULL COMMENT '资产来源',
  `website` varchar(128) DEFAULT NULL COMMENT '政府企业/公司注册网页链接',
  `licence_code` varchar(128) DEFAULT NULL COMMENT '该实体持有的监管许可信息、办法此许可的监管机构、适用的司法管辖区和许可证号',
  `name_hashkey` varchar(32) DEFAULT NULL COMMENT '您在Hashkey的联系人',
  `representative_name` varchar(32) DEFAULT NULL COMMENT '公司代表姓名',
  `representative_job` varchar(32) DEFAULT NULL COMMENT '公司代表职务',
  `representative_email` varchar(32) DEFAULT NULL COMMENT '公司代表邮箱',
  `representative_mobile_pre` varchar(16) DEFAULT NULL COMMENT '公司代表手机号',
  `representative_mobile` varchar(16) DEFAULT NULL COMMENT '公司代表手机号',
  `reasion` varchar(32) DEFAULT NULL COMMENT '注册本平台的原因',
  `reasion_str` varchar(128) DEFAULT NULL COMMENT '注册本平台的原因',
  `trade_year` int(8) DEFAULT NULL COMMENT '交易数字货币几年',
  `what_coins` varchar(128) DEFAULT NULL COMMENT '希望交易哪些数字货币',
  `what_platforms` varchar(128) DEFAULT NULL COMMENT '目前持有的数字货币存放在哪些平台',
  `is_fiat` tinyint(1) DEFAULT NULL COMMENT '是否打算交易法币币对',
  `bank_name` varchar(32) DEFAULT NULL COMMENT '银行名称',
  `bank_address` varchar(256) DEFAULT NULL COMMENT '银行地址',
  `swift` varchar(32) DEFAULT NULL COMMENT 'SWIFT',
  `route_code` varchar(64) DEFAULT NULL COMMENT '路由码',
  `beneficiary_name_on_account` varchar(64) DEFAULT NULL COMMENT '受益人账户名称',
  `beneficiary_bank_account_no` varchar(32) DEFAULT NULL COMMENT '受益人银行账号',
  `beneficiary_address_on_account` varchar(256) DEFAULT NULL COMMENT '受益人银行账号地址',
  `remarks` varchar(256) DEFAULT NULL COMMENT '备注',
  `regist_photo` varchar(512) DEFAULT NULL COMMENT '公司注册证书',
  `rule_photo` varchar(512) DEFAULT NULL COMMENT '备忘录和公司章程',
  `power_photo` varchar(512) DEFAULT NULL COMMENT '董事职权证明书',
  `subsist_photo` varchar(512) DEFAULT NULL COMMENT '公司良好续存证明书',
  `book_photo` varchar(512) DEFAULT NULL COMMENT '董事/高级职员登记',
  `share_book_photo` varchar(512) DEFAULT NULL COMMENT '股东登记册(如果股权结构不止一层，请提供:1、结构图 2、任何股权占25%或以上最终受益人的股东登记册)',
  `reg_add_photo` varchar(512) DEFAULT NULL COMMENT '现有注册地址证明（以及注册地址不同的通讯地址）',
  `director_name` varchar(32) DEFAULT NULL COMMENT '法定名称',
  `dir_live_country` varchar(32) DEFAULT NULL COMMENT '居住国家',
  `dir_other_job` varchar(32) DEFAULT NULL COMMENT '在公司担任的其他职务/职责',
  `dir_pass_photo` varchar(128) DEFAULT NULL COMMENT '有效护照照片页（清楚显示名字、姓氏、出生日期、性别、国籍）',
  `dir_add_photo` varchar(512) DEFAULT NULL COMMENT '现有注册地址证明(必须是3个月内签发的银行或公用事业对账单，清楚地显示个人姓名、地址及签发日期)',
  `dir_hand_photo` varchar(128) DEFAULT NULL COMMENT '董事本人手持护照和“Hashkey + 今天日期”便条的照片',
  `sha_name` varchar(32) DEFAULT NULL,
  `percentage` varchar(32) DEFAULT NULL COMMENT '法定名称',
  `live_country` varchar(32) DEFAULT NULL COMMENT '居住国家',
  `other_job` varchar(32) DEFAULT NULL COMMENT '在公司担任的其他职务/职责',
  `sha_pass_photo` varchar(128) DEFAULT NULL COMMENT '有效护照照片页（清楚显示名字、姓氏、出生日期、性别、国籍）',
  `sha_add_photo` varchar(256) DEFAULT NULL COMMENT '现有注册地址证明(必须是3个月内签发的银行或公用事业对账单，清楚地显示个人姓名、地址及签发日期)',
  `trader_name` varchar(32) DEFAULT NULL COMMENT '法定名称',
  `trader_mobile_pre` varchar(16) DEFAULT NULL COMMENT '手机号',
  `trader_mobile` varchar(16) DEFAULT NULL COMMENT '手机号',
  `msn_type` varchar(32) DEFAULT NULL COMMENT '即时通讯工具(WeChat...)',
  `msn_number` varchar(32) DEFAULT NULL COMMENT '即时通讯工具联系方式',
  `trader_email` varchar(32) DEFAULT NULL COMMENT '邮箱',
  `trader_pass_photo` varchar(128) DEFAULT NULL COMMENT '有效护照照片页（清楚显示名字、姓氏、出生日期、性别、国籍）',
  `trader_add_photo` varchar(256) DEFAULT NULL COMMENT '现有注册地址证明(必须是3个月内签发的银行或公用事业对账单，清楚地显示个人姓名、地址及签发日期)',
  `other_files` varchar(512) DEFAULT NULL COMMENT '其他文件',
  `director_info` varchar(2048) DEFAULT NULL COMMENT '董事信息',
  `shareholder_info` varchar(2048) DEFAULT NULL COMMENT '股东信息',
  `auth_trader_info` varchar(2048) DEFAULT NULL COMMENT '授权交易员',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=61 DEFAULT CHARSET=utf8;

#
# Structure for table "kyc_person"
#

DROP TABLE IF EXISTS `kyc_person`;
CREATE TABLE `kyc_person` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(18) DEFAULT NULL,
  `last_update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `name` varchar(32) DEFAULT NULL,
  `sur_name` varchar(32) DEFAULT NULL,
  `birth_date` varchar(32) DEFAULT NULL,
  `country` varchar(32) DEFAULT NULL,
  `post_code` varchar(32) DEFAULT NULL,
  `address` varchar(1280) DEFAULT NULL,
  `tax_code` varchar(32) DEFAULT NULL,
  `fund_source` tinyint(1) DEFAULT NULL,
  `fund_source_str` varchar(256) DEFAULT NULL,
  `occupation` varchar(256) DEFAULT NULL,
  `employer_nature` varchar(128) DEFAULT NULL,
  `asset_source` tinyint(1) DEFAULT NULL,
  `asset_source_str` varchar(256) DEFAULT NULL,
  `name_hashkey` varchar(32) DEFAULT NULL,
  `reasion` varchar(32) DEFAULT NULL,
  `reasion_str` varchar(128) DEFAULT NULL,
  `trade_year` int(8) DEFAULT NULL,
  `what_coins` varchar(128) DEFAULT NULL,
  `what_platforms` varchar(128) DEFAULT NULL,
  `is_fiat` tinyint(1) DEFAULT NULL,
  `bank_name` varchar(32) DEFAULT NULL,
  `bank_address` varchar(1280) DEFAULT NULL,
  `swift` varchar(32) DEFAULT NULL,
  `route_code` varchar(64) DEFAULT NULL COMMENT '路由码',
  `beneficiary_name_on_account` varchar(32) DEFAULT NULL,
  `beneficiary_bank_account_no` varchar(32) DEFAULT NULL,
  `beneficiary_address_on_account` varchar(1280) DEFAULT NULL,
  `remarks` varchar(256) DEFAULT NULL,
  `certificates_country` varchar(10) DEFAULT NULL,
  `id_type` tinyint(1) DEFAULT NULL,
  `id_number` varchar(18) DEFAULT NULL,
  `passport_photo` varchar(128) DEFAULT NULL,
  `id_photo_1` varchar(128) DEFAULT NULL,
  `id_photo_2` varchar(128) DEFAULT NULL,
  `hand_photo` varchar(128) DEFAULT NULL,
  `address_photo` varchar(12800) DEFAULT NULL,
  `other_files` varchar(1280) DEFAULT NULL COMMENT '其他文件',
  `company` varchar(256) DEFAULT '' COMMENT '当前职业',
  `msn_type` varchar(32) DEFAULT '' COMMENT '即时通讯工具(WeChat...)',
  `msn_number` varchar(32) DEFAULT '' COMMENT '即时通讯工具联系方式',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=218 DEFAULT CHARSET=utf8;

#
# Structure for table "member_admins"
#

DROP TABLE IF EXISTS `member_admins`;
CREATE TABLE `member_admins` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` bigint(18) DEFAULT NULL COMMENT '会员ID',
  `id_photo_1` varchar(128) DEFAULT NULL COMMENT '证件照片1',
  `id_photo_2` varchar(128) DEFAULT NULL COMMENT '证件照片2',
  `id_photo_3` varchar(128) DEFAULT NULL COMMENT '证件照片3',
  `approve_status` tinyint(1) DEFAULT NULL COMMENT '审核状态',
  `approve_note(d)` varchar(16) DEFAULT NULL COMMENT '审核意见',
  `approve_user` varchar(16) DEFAULT NULL COMMENT '审核人',
  `approve_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '审核时间',
  `member_status` tinyint(1) DEFAULT NULL COMMENT '状态',
  `approve_status_basic` tinyint(1) DEFAULT NULL,
  `approve_node_of_basic` varchar(512) DEFAULT NULL,
  `approve_status_trade` tinyint(1) DEFAULT NULL,
  `approve_note_of_trade` varchar(512) DEFAULT NULL,
  `approve_status_settle` tinyint(1) DEFAULT NULL,
  `approve_note_of_settle` varchar(512) DEFAULT NULL,
  `approve_status_file` tinyint(1) DEFAULT NULL,
  `approve_note_of_file` varchar(512) DEFAULT NULL,
  `operate_user` varchar(16) DEFAULT NULL,
  `account_type` tinyint(1) DEFAULT NULL,
  `approve_status_basic_client` tinyint(1) DEFAULT NULL,
  `approve_status_trade_client` tinyint(1) DEFAULT NULL,
  `approve_status_settle_client` tinyint(1) DEFAULT NULL,
  `approve_status_file_client` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_member_admins_uid` (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=596 DEFAULT CHARSET=utf8;

#
# Structure for table "member_cash_flow"
#

DROP TABLE IF EXISTS `member_cash_flow`;
CREATE TABLE `member_cash_flow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) DEFAULT NULL COMMENT '会员id',
  `cash_flow` decimal(26,16) DEFAULT NULL COMMENT '现金流',
  `market_value` decimal(26,16) DEFAULT NULL COMMENT '日终市值',
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `date_time` date DEFAULT NULL COMMENT '日期',
  `market_value_start` decimal(26,16) DEFAULT NULL COMMENT '日初市价',
  PRIMARY KEY (`id`),
  KEY `member_date` (`member_id`,`date_time`)
) ENGINE=InnoDB AUTO_INCREMENT=1256 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "member_transactions"
#

DROP TABLE IF EXISTS `member_transactions`;
CREATE TABLE `member_transactions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` bigint(18) DEFAULT NULL COMMENT '会员ID',
  `account_id` bigint(18) DEFAULT NULL COMMENT '账户ID',
  `detail_id` bigint(18) DEFAULT NULL COMMENT '订单明细ID',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '交易金额',
  `transaction_type` tinyint(1) DEFAULT NULL COMMENT '交易类型',
  `symbol` varchar(16) DEFAULT NULL COMMENT '币种名称',
  `tx_id` varchar(128) DEFAULT NULL COMMENT '提现交易编号',
  `address` varchar(64) DEFAULT NULL COMMENT '充值或提现地址、或转账地址',
  `fee` decimal(26,16) DEFAULT NULL COMMENT '交易手续费',
  `flag` tinyint(1) DEFAULT NULL COMMENT '标识位',
  `real_fee` decimal(26,16) DEFAULT NULL COMMENT '实收手续费',
  `discount_fee` decimal(26,16) DEFAULT NULL COMMENT '折扣手续费',
  `pre_balance` decimal(26,16) DEFAULT NULL COMMENT '处理前余额',
  `balance` decimal(26,16) DEFAULT NULL COMMENT '处理前余额',
  `pre_frozen_bal` decimal(26,16) DEFAULT NULL COMMENT '处理前冻结余额',
  `frozen_bal` decimal(26,16) DEFAULT NULL COMMENT '处理后冻结余额',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `coin_price` decimal(26,16) DEFAULT NULL COMMENT 'base当前价格',
  `average_price` decimal(26,16) DEFAULT NULL COMMENT '当前均价',
  `profit` decimal(26,16) DEFAULT NULL COMMENT '收益',
  PRIMARY KEY (`id`),
  KEY `idx_member` (`member_id`),
  KEY `idx_account` (`account_id`),
  KEY `idx_detail` (`detail_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10758 DEFAULT CHARSET=utf8;

#
# Structure for table "members"
#

DROP TABLE IF EXISTS `members`;
CREATE TABLE `members` (
  `id` bigint(20) NOT NULL COMMENT '主键ID',
  `member_id` varchar(16) DEFAULT NULL COMMENT '会员ID',
  `group_id` bigint(18) DEFAULT NULL COMMENT '会员组ID',
  `salt` varchar(32) DEFAULT NULL COMMENT '盐值',
  `username` varchar(32) DEFAULT NULL COMMENT '用户名',
  `password` varchar(128) DEFAULT NULL COMMENT '密码',
  `gender` tinyint(1) DEFAULT NULL COMMENT '性别',
  `nickname` varchar(32) DEFAULT NULL COMMENT '昵称',
  `is_google_auth` tinyint(1) DEFAULT NULL COMMENT '是否谷歌认证',
  `google_auth_secret` varchar(32) DEFAULT NULL COMMENT '谷歌密钥',
  `margin` tinyint(1) DEFAULT NULL COMMENT '是否缴纳保证金',
  `jy_password` varchar(32) DEFAULT NULL COMMENT '交易密码',
  `real_name` varchar(32) DEFAULT NULL COMMENT '会员真实姓名',
  `id_type` tinyint(1) DEFAULT NULL COMMENT '证件类型',
  `id_number` varchar(18) DEFAULT NULL COMMENT '身份证号码',
  `email` varchar(32) DEFAULT NULL COMMENT '邮箱',
  `mobile` varchar(16) DEFAULT NULL COMMENT '手机号',
  `location` varchar(20) DEFAULT NULL COMMENT '定位',
  `member_type` tinyint(1) DEFAULT NULL COMMENT '会员类型',
  `member_level` varchar(10) DEFAULT NULL COMMENT '会员等级',
  `status` int(2) DEFAULT NULL COMMENT '状态',
  `registration_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `last_login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后登录时间',
  `token` varchar(64) DEFAULT NULL COMMENT 'token',
  `super_partner` tinyint(1) DEFAULT NULL COMMENT '超级合伙人',
  `transactions` int(10) DEFAULT NULL COMMENT '交易次数',
  `inviter_id` bigint(18) DEFAULT NULL COMMENT '邀请者ID',
  `promotion_code` varchar(20) DEFAULT NULL COMMENT '推广码',
  `real_name_status` tinyint(1) DEFAULT NULL COMMENT '实名认证状态',
  `login_count` int(8) DEFAULT NULL COMMENT '登录次数',
  `country` varchar(10) DEFAULT NULL COMMENT '国家',
  `province` varchar(40) DEFAULT NULL COMMENT '最后登录的省',
  `city` varchar(40) DEFAULT NULL COMMENT '最后登录的市',
  `device_id` varchar(256) DEFAULT NULL,
  `ip` varchar(32) DEFAULT NULL COMMENT '最后登录的ip',
  `token_expire_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'token预计过期时间',
  `transaction_status` tinyint(1) DEFAULT NULL COMMENT '交易状态',
  `register_type` tinyint(1) DEFAULT NULL,
  `area_code` varchar(16) DEFAULT NULL,
  `modify_username_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_members_member_id` (`member_id`) USING BTREE,
  KEY `idx_group` (`group_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

#
# Structure for table "menu"
#

DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增长id',
  `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT 'parent_id',
  `menu_name_cn` varchar(50) NOT NULL,
  `menu_name_en` varchar(50) NOT NULL,
  `menu_name_desc` varchar(50) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '菜单名',
  `front_key` varchar(50) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '前端用的唯一key参数',
  `menu_order` int(11) NOT NULL DEFAULT '0' COMMENT '菜单顺序，数值越小优先级越高',
  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '全局唯一的路径',
  `last_operator` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后一次操作人,0-表示由程序自动化创建',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `menu_order` (`menu_order`) USING BTREE,
  KEY `create_time` (`created`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=127 DEFAULT CHARSET=utf8 COMMENT='菜单表';

#
# Structure for table "official_recharge_bank"
#

DROP TABLE IF EXISTS `official_recharge_bank`;
CREATE TABLE `official_recharge_bank` (
  `id` bigint(18) NOT NULL COMMENT '主键ID',
  `bank_name` varchar(32) DEFAULT NULL COMMENT '银行名称',
  `bank_address` varchar(1280) DEFAULT NULL COMMENT '银行地址',
  `swift` varchar(32) DEFAULT NULL COMMENT 'SWIFT',
  `route_code` varchar(32) DEFAULT NULL COMMENT '路由传输号码 / ABA路由号',
  `beneficiary_name_on_account` varchar(32) DEFAULT NULL COMMENT '银行账号姓名',
  `beneficiary_bank_account_no` varchar(32) DEFAULT NULL COMMENT '银行账号',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#
# Structure for table "offline_trade_input"
#

DROP TABLE IF EXISTS `offline_trade_input`;
CREATE TABLE `offline_trade_input` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `trade_business_no` varchar(32) NOT NULL COMMENT '业务编号',
  `trade_time` date DEFAULT NULL COMMENT '交易时间',
  `offline_trade_id` varchar(32) NOT NULL COMMENT '线下交易ID',
  `uid` bigint(18) DEFAULT NULL COMMENT 'uid',
  `user_name` varchar(32) DEFAULT NULL COMMENT '客户账号',
  `name` varchar(32) DEFAULT NULL COMMENT '姓名',
  `symbol` varchar(32) DEFAULT NULL COMMENT '品种',
  `direction` tinyint(1) DEFAULT NULL COMMENT '0买 1卖',
  `quantity` decimal(26,16) DEFAULT NULL COMMENT '数量',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '金额',
  `price` decimal(26,16) DEFAULT NULL COMMENT '价格',
  `fee` decimal(26,16) DEFAULT NULL COMMENT '手续费',
  `trade_id` varchar(32) DEFAULT NULL COMMENT '内部成交ID',
  `remarks` varchar(512) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT NULL COMMENT '业务状态 1一次录入 2二次录入',
  `operator1` varchar(32) DEFAULT NULL COMMENT '操作员1',
  `operator2` varchar(32) DEFAULT NULL COMMENT '操作员2',
  `revise_operator` varchar(32) DEFAULT NULL COMMENT '修正操作员',
  `operate_time` datetime DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_offline_trade_input_trade_business_no` (`trade_business_no`),
  UNIQUE KEY `offline_trade_id_UNIQUE` (`offline_trade_id`)
) ENGINE=InnoDB AUTO_INCREMENT=190 DEFAULT CHARSET=utf8;

#
# Structure for table "offline_trade_input_history"
#

DROP TABLE IF EXISTS `offline_trade_input_history`;
CREATE TABLE `offline_trade_input_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `trade_business_no` varchar(32) NOT NULL COMMENT '业务编号',
  `trade_time` date DEFAULT NULL,
  `offline_trade_id` varchar(32) NOT NULL COMMENT '线下交易ID',
  `uid` bigint(18) DEFAULT NULL COMMENT 'uid',
  `user_name` varchar(32) DEFAULT NULL COMMENT '客户账号',
  `name` varchar(32) DEFAULT NULL COMMENT '姓名',
  `symbol` varchar(32) DEFAULT NULL COMMENT '品种',
  `direction` tinyint(1) DEFAULT NULL COMMENT '0买 1卖',
  `quantity` decimal(26,16) DEFAULT NULL COMMENT '数量',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '金额',
  `price` decimal(26,16) DEFAULT NULL COMMENT '价格',
  `fee` decimal(26,16) DEFAULT NULL COMMENT '手续费',
  `trade_id` varchar(32) DEFAULT NULL COMMENT '内部成交ID',
  `remarks` varchar(512) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT NULL COMMENT '业务状态 1一次录入 2二次录入',
  `operator1` varchar(32) DEFAULT NULL COMMENT '操作员1',
  `operator2` varchar(32) DEFAULT NULL COMMENT '操作员2',
  `revise_operator` varchar(32) DEFAULT NULL COMMENT '修正操作员',
  `operate_time` datetime DEFAULT NULL COMMENT '操作时间',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=228 DEFAULT CHARSET=utf8;

#
# Structure for table "portfolio_group_items"
#

DROP TABLE IF EXISTS `portfolio_group_items`;
CREATE TABLE `portfolio_group_items` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `gid` bigint(20) DEFAULT NULL COMMENT '资产组ID',
  `symbol` varchar(128) DEFAULT NULL COMMENT '分组中包含的币种',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `gid_symbol` (`gid`,`symbol`)
) ENGINE=InnoDB AUTO_INCREMENT=161 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "portfolio_groups"
#

DROP TABLE IF EXISTS `portfolio_groups`;
CREATE TABLE `portfolio_groups` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(128) DEFAULT NULL COMMENT '资产分组名称',
  `member_id` bigint(20) DEFAULT NULL COMMENT '会员id',
  `remark` varchar(256) DEFAULT NULL COMMENT '备注',
  `status` tinyint(1) DEFAULT '1' COMMENT '1:默认分组，2:私有分组',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `last_operator` bigint(20) DEFAULT '0' COMMENT '最后一次操作人,0-表示由程序自动创建',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "portfolio_performance_details"
#

DROP TABLE IF EXISTS `portfolio_performance_details`;
CREATE TABLE `portfolio_performance_details` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` bigint(20) DEFAULT NULL COMMENT '会员id',
  `detail_id` bigint(20) DEFAULT NULL COMMENT '订单明细ID',
  `record_type` int(11) DEFAULT NULL COMMENT '记录类型',
  `market_price` varchar(5000) DEFAULT NULL COMMENT '持有币种的市价',
  `coin` varchar(5000) DEFAULT NULL COMMENT '持有币种以及量',
  `date_time` date DEFAULT NULL COMMENT '日期',
  `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `market_value_before` decimal(26,16) DEFAULT NULL COMMENT '本次资金变动之前的市值',
  `market_value_after` decimal(26,16) DEFAULT NULL COMMENT '本次资金变动之后的市值',
  PRIMARY KEY (`id`),
  KEY `idx_date_time` (`date_time`),
  KEY `idx_member` (`member_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14806 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "profits"
#

DROP TABLE IF EXISTS `profits`;
CREATE TABLE `profits` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` bigint(20) DEFAULT NULL COMMENT '会员id',
  `detail_id` bigint(20) DEFAULT NULL COMMENT '订单明细ID',
  `symbol` varchar(255) DEFAULT NULL COMMENT '币种',
  `transaction_type` tinyint(1) DEFAULT NULL COMMENT '交易类型',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `profit` decimal(26,16) DEFAULT NULL COMMENT '收益',
  `profit_time` date DEFAULT NULL COMMENT '交易时间',
  `direction` tinyint(1) DEFAULT NULL COMMENT '交易方向,0买入 1卖出',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态,1有效 2 无效',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_detail` (`detail_id`),
  KEY `idx_profit` (`member_id`,`profit_time`,`symbol`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1046 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "role"
#

DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增长id',
  `role_name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
  `role_desc` varchar(50) NOT NULL DEFAULT '' COMMENT '角色描述',
  `last_operator` bigint(20) NOT NULL DEFAULT '0' COMMENT ' 最后一次操作人,0-表示由程序自动化创建',
  `deleted_timestamp` bigint(20) NOT NULL DEFAULT '0' COMMENT '删除时间戳',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_name` (`deleted_timestamp`,`role_name`),
  KEY `deleted_at` (`deleted_timestamp`,`updated`),
  KEY `create_time` (`created`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

#
# Structure for table "role_operation"
#

DROP TABLE IF EXISTS `role_operation`;
CREATE TABLE `role_operation` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增长id',
  `role_id` bigint(20) NOT NULL DEFAULT '0' COMMENT 'role_id',
  `operation_id` bigint(20) NOT NULL DEFAULT '0' COMMENT 'operation_id（=menu_id），目前operation也存放在menu表',
  `last_operator` bigint(20) NOT NULL DEFAULT '0' COMMENT ' 最后一次操作人,0-表示由程序自动化创建',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id__operation_id` (`role_id`,`operation_id`),
  KEY `operation_id` (`operation_id`),
  KEY `create_time` (`created`)
) ENGINE=InnoDB AUTO_INCREMENT=1426 DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

#
# Structure for table "security_log"
#

DROP TABLE IF EXISTS `security_log`;
CREATE TABLE `security_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增长id',
  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '操作内容',
  `operator` bigint(20) NOT NULL DEFAULT '0' COMMENT '操作人,0-表示由程序自动化创建',
  `ip` varchar(100) NOT NULL DEFAULT '',
  `content` varchar(100) NOT NULL DEFAULT '' COMMENT '操作内容',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `operator` (`operator`) USING BTREE,
  KEY `create_time` (`created`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3070 DEFAULT CHARSET=utf8 COMMENT='安全日志表';

#
# Structure for table "trading_hedging_alert_logs"
#

DROP TABLE IF EXISTS `trading_hedging_alert_logs`;
CREATE TABLE `trading_hedging_alert_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `type` varchar(50) NOT NULL COMMENT '类型',
  `message` text COMMENT '详情',
  `send_status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0-待发送 2-已发送',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
  PRIMARY KEY (`id`),
  KEY `type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=7962 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "trading_hedging_alert_setting"
#

DROP TABLE IF EXISTS `trading_hedging_alert_setting`;
CREATE TABLE `trading_hedging_alert_setting` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `key` varchar(50) NOT NULL COMMENT 'key',
  `value` text COMMENT 'value',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `type` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "trading_limit_setting"
#

DROP TABLE IF EXISTS `trading_limit_setting`;
CREATE TABLE `trading_limit_setting` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `symbol` varchar(16) NOT NULL COMMENT '币对',
  `daily_buying_limit_weekday` decimal(26,16) DEFAULT NULL COMMENT '日买限额weekday',
  `daily_buying_limit_weekend` decimal(26,16) DEFAULT NULL COMMENT '日买限额weekend',
  `daily_selling_limit_weekday` decimal(26,16) DEFAULT NULL COMMENT '日卖限额weekday',
  `daily_selling_limit_weekend` decimal(26,16) DEFAULT NULL COMMENT '日卖限额weekend',
  `daily_buying_limit_per_client_weekday` decimal(26,16) DEFAULT NULL COMMENT '单人日买限额weekday',
  `daily_buying_limit_per_client_weekend` decimal(26,16) DEFAULT NULL COMMENT '单人日买限额weekend',
  `daily_selling_limit_per_client_weekday` decimal(26,16) DEFAULT NULL COMMENT '单人日卖限额weekday',
  `daily_selling_limit_per_client_weekend` decimal(26,16) DEFAULT NULL COMMENT '单人日卖限额weekend',
  `status` int(11) DEFAULT NULL COMMENT '0-审核中，1-生效，2-拒绝，3-作废',
  `pending_status_setting_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '关联的pending状态的记录id，方便搜索使用',
  `obsoleted_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '作废或者拒绝的时间戳，单位秒',
  `create_operator` bigint(20) DEFAULT NULL COMMENT '创建操作人',
  `review_operator` bigint(20) DEFAULT NULL COMMENT '审核操作人',
  `rejected_reason` varchar(200) NOT NULL COMMENT '审核拒绝时的备注',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_symbol_status` (`symbol`,`status`,`obsoleted_time`),
  KEY `symbol_updated` (`symbol`,`updated`),
  KEY `symbol_status` (`symbol`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4;

#
# Structure for table "user_role"
#

DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增长id',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT 'user_id',
  `role_id` bigint(20) NOT NULL DEFAULT '0' COMMENT 'role_id',
  `last_operator` bigint(20) NOT NULL DEFAULT '0' COMMENT ' 最后一次操作人,0-表示由程序自动化创建',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id__role_id` (`user_id`,`role_id`),
  KEY `role_id` (`role_id`),
  KEY `create_time` (`created`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

#
# Structure for table "withdraws"
#

DROP TABLE IF EXISTS `withdraws`;
CREATE TABLE `withdraws` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` bigint(18) DEFAULT NULL COMMENT '会员ID',
  `coin_id` bigint(18) DEFAULT NULL COMMENT '币种ID',
  `coin_unit` varchar(16) DEFAULT NULL COMMENT '币种',
  `to` varchar(64) DEFAULT NULL COMMENT '提现接收地址',
  `amount` decimal(26,16) DEFAULT NULL COMMENT '提现数量',
  `memo` varchar(16) DEFAULT NULL COMMENT '备忘录',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `request_id` varchar(32) DEFAULT NULL COMMENT '请求编号',
  PRIMARY KEY (`id`),
  KEY `idx_member_id` (`member_id`),
  KEY `idx_coin_id` (`coin_id`),
  KEY `idx_coin` (`coin_unit`)
) ENGINE=InnoDB AUTO_INCREMENT=82861926749962241 DEFAULT CHARSET=utf8;

#
# Procedure "idata"
#

DROP PROCEDURE IF EXISTS `idata`;
CREATE PROCEDURE `idata`()
begin  
declare i int;  
set i=0;  while i<2 do    
insert into member_return_rates (myday, member_id, rate_of_return) values (date_sub(curdate(),interval -i day), 323072, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "idata0"
#

DROP PROCEDURE IF EXISTS `idata0`;
CREATE PROCEDURE `idata0`()
begin  
declare i int;  
set i=0;  while i<366 do    
insert into member_return_rates (date_time, member_id, rate_of_return) values (date_sub(curdate(),interval i day), 45886720, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "idata1"
#

DROP PROCEDURE IF EXISTS `idata1`;
CREATE PROCEDURE `idata1`()
begin  
declare i int;  
set i=0;  while i<2 do    
insert into member_return_rates (date_time, member_id, rate_of_return) values (date_sub(curdate(),interval -i day), 323072, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "idata11"
#

DROP PROCEDURE IF EXISTS `idata11`;
CREATE PROCEDURE `idata11`()
begin  
declare i int;  
set i=0;  while i<366 do    
insert into member_return_rates (date_time, member_id, rate_of_return) values (date_sub(curdate(),interval i day), 76276480, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "idata12"
#

DROP PROCEDURE IF EXISTS `idata12`;
CREATE PROCEDURE `idata12`()
begin  
declare i int;  
set i=0;  while i<366 do    
insert into member_return_rates (date_time, member_id, rate_of_return) values (date_sub(curdate(),interval i day), 76276480, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "idata13"
#

DROP PROCEDURE IF EXISTS `idata13`;
CREATE PROCEDURE `idata13`()
begin  
declare i int;  
set i=0;  while i<366 do    
insert into member_return_rates (date_time, member_id, rate_of_return) values (date_sub(curdate(),interval i day), 45886720, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "idata2"
#

DROP PROCEDURE IF EXISTS `idata2`;
CREATE PROCEDURE `idata2`()
begin  
declare i int;  
set i=0;  while i<5 do    
insert into member_return_rates (date_time, member_id, rate_of_return) values (date_sub(curdate(),interval -i day), 323072, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "idata3"
#

DROP PROCEDURE IF EXISTS `idata3`;
CREATE PROCEDURE `idata3`()
begin  
declare i int;  
set i=0;  while i<366 do    
insert into member_return_rates (date_time, member_id, rate_of_return) values (date_sub(curdate(),interval i day), 323072, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "idata9"
#

DROP PROCEDURE IF EXISTS `idata9`;
CREATE PROCEDURE `idata9`()
begin  
declare i int;  
set i=0;  while i<366 do    
insert into member_return_rates (date_time, member_id, rate_of_return) values (date_sub(curdate(),interval i day), 76276480, RAND());
set i=i+1; 
end while;
end;

#
# Procedure "wk"
#

DROP PROCEDURE IF EXISTS `wk`;
CREATE PROCEDURE `wk`()
begin
declare i int;
set i = 1;
while i < 2 do
insert into member_return_rates (myday, member_id, rate_of_return) values (date_sub(curdate(),interval -i day), 323072, RAND());
-- select i;
set i = i +1;
end while;
end;
