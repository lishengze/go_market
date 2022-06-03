DROP TABLE IF EXISTS `exchange_orders`;
CREATE TABLE `exchange_orders` (
   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
   `order_id` varchar(20) DEFAULT NULL COMMENT '订单ID',
   `user_id` bigint(18) DEFAULT NULL COMMENT '用户ID',
   `symbol` varchar(16) DEFAULT NULL COMMENT '交易币对',
   `price` decimal(18,8) DEFAULT NULL COMMENT '挂单价格',
   `status` tinyint(1) DEFAULT NULL COMMENT '订单状态',
   `direction` tinyint(1) DEFAULT NULL COMMENT '订单方向',
   `traded_volume` decimal(26,16) DEFAULT NULL COMMENT '成交量',
   `turnover` decimal(26,16) DEFAULT NULL COMMENT '成交额',
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `exchange_event`;
CREATE TABLE `exchange_event` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_order` varchar(20) DEFAULT NULL COMMENT '订单ID',
  `user_id` bigint(18) DEFAULT NULL COMMENT '用户ID',
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '挂单时间,取消时间,完成时间',
  `type` int(11) DEFAULT NULL COMMENT '事件类型',
  `status` tinyint(1) DEFAULT NULL COMMENT '订单状态',
  `name` varchar(20) DEFAULT NULL COMMENT '事件名',
  `message` varchar(512) DEFAULT NULL COMMENT '事件',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `exchange_trades`;
CREATE TABLE `exchange_trades` (
   `id` bigint(18) NOT NULL COMMENT '主键ID',
   `symbol` varchar(12) DEFAULT NULL COMMENT '交易币种名称，格式：BTC_USDT',
   `trade_type` tinyint(1) DEFAULT NULL COMMENT '成交类型',
   `offline_trade_id` varchar(32) DEFAULT NULL COMMENT '线下交易ID',
   `price` decimal(18,8) DEFAULT NULL COMMENT '交易价格',
   `volume` decimal(26,16) DEFAULT NULL COMMENT '成交量',
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
   PRIMARY KEY (`id`),
   KEY `idx_symbol` (`symbol`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `hedge_matches`;
CREATE TABLE `hedge_matches` (
 `id` bigint(18) NOT NULL AUTO_INCREMENT,
 `match_id` varchar(100) DEFAULT NULL,
 `order_id` varchar(100) DEFAULT NULL,
 `exchange_name` varchar(16) DEFAULT NULL,
 `account_name` varchar(32) DEFAULT NULL,
 `symbol` varchar(16) DEFAULT NULL,
 `direction` tinyint(1)  DEFAULT NULL,
 `trade_price` decimal(26,16) DEFAULT NULL,
 `trade_amount` decimal(26,16) DEFAULT NULL,
 `trade_value` decimal(26,16) DEFAULT NULL,
 `hedge_id` bigint(18) DEFAULT NULL,
 `match_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `amount` decimal(26,16) DEFAULT NULL,
 `fee_currency` varchar(16) DEFAULT NULL,
 PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `hedge_to_centers`;
CREATE TABLE `hedge_to_centers` (
    `id` bigint(18) NOT NULL AUTO_INCREMENT,
    `hedge_match_id` bigint(18) DEFAULT NULL,
    `central_match_id` bigint(18) DEFAULT NULL,
    `hedge_amount` decimal(26,16) DEFAULT NULL,
    `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
