
DROP TABLE IF EXISTS `exchange_event`;
CREATE TABLE `exchange_event` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_order` varchar(20) DEFAULT '' COMMENT '订单ID',
  `user_id` bigint(18) DEFAULT 0 COMMENT '用户ID',
  `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '挂单时间,取消时间,完成时间',
  `type` int(11) DEFAULT 0 COMMENT '事件类型',
  `status` tinyint(1) DEFAULT 0 COMMENT '订单状态',
  `name` varchar(20) DEFAULT '' COMMENT '事件名',
  `message` varchar(512) DEFAULT '' COMMENT '事件',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ;

DROP TABLE IF EXISTS `order`;
CREATE TABLE `order` (
    `id` bigint(10) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id` bigint(10) NOT NULL COMMENT '用户id',
    `order_local_id` varchar(16) NOT NULL COMMENT '本地报单编号',
    `order_type` tinyint(1) NOT NULL COMMENT '订单类型 1:otc订单 2:聚合交易订单',
    `order_mode` tinyint(1) DEFAULT 0 COMMENT '下单模式 0:按数量询价/下单 1:按金额询价/下单',
    `order_price_type` tinyint(1) NOT NULL COMMENT '报单价格条件 1:限价单 2:市价单',
    `symbol` varchar(16) NOT NULL COMMENT '币对 BTC_USDT',
    `base_currency` varchar(16) NOT NULL COMMENT '基础币种 USDT',
    `target_currency` varchar(16) NOT NULL COMMENT '目标币种 BTC',
    `direction` tinyint(1) NOT NULL COMMENT '0:买 1:卖',
    `price` decimal(18, 10) COMMENT '报单价格',
    `volume` decimal(18, 10) COMMENT '报单数量',
    `amount` decimal(18, 10) COMMENT '报单金额, 根据金额下单时用到',
    `order_maker` tinyint(1) DEFAULT 0 COMMENT '1:taker 2:maker',
    `trade_volume` decimal(18, 10) COMMENT '已成交数量',
    `trade_amount` decimal(18, 10) COMMENT '已成交金额(包括手续费)',
    `fee_kind` tinyint(1) NOT NULL COMMENT '订单创建时的手续费类型 1:百分比 2:绝对值',
    `fee_rate` decimal(18, 10) COMMENT '订单创建时的手续费类型',
    `order_status` tinyint(1) DEFAULT 0 COMMENT '订单状态 0:已发送 1:全部成交 2:部分成交在队列 3:部分成交已撤单 4:撤单',
    `order_create_time` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '订单创建时间',
    `order_modify_time` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '订单修改时间',
    PRIMARY KEY (`id`),
    INDEX index_user_id (`user_id`),
    INDEX index_order_local_id (`order_local_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '订单表';

DROP TABLE IF EXISTS `trade`;
CREATE TABLE `trade` (
    `id` bigint(10) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id` bigint(10) NOT NULL COMMENT '用户id',
    `order_local_id` varchar(16) NOT NULL COMMENT '本地报单编号',
    `trade_id` varchar(16) NOT NULL COMMENT '成交编号',
    `order_type` tinyint(1) NOT NULL COMMENT '订单类型 1:otc订单 2:聚合交易订单',
    `order_mode` tinyint(1) DEFAULT 0 COMMENT '下单模式 0:按数量询价/下单 1:按金额询价/下单',
    `order_price_type` tinyint(1) NOT NULL COMMENT '报单价格条件 1:限价单 2:市价单',
    `symbol` varchar(16) NOT NULL COMMENT '币对 BTC_USDT',
    `base_currency` varchar(16) NOT NULL COMMENT '基础币种 USDT',
    `target_currency` varchar(16) NOT NULL COMMENT '目标币种 BTC',
    `direction` tinyint(1) NOT NULL COMMENT '0:买 1:卖',
    `price` decimal(18, 10) COMMENT '报单价格',
    `volume` decimal(18, 10) COMMENT '报单数量，冗余order.volume',
    `amount` decimal(18, 10) COMMENT '报单金额, 根据金额下单时用到，冗余order.amount',
    `order_maker` tinyint(1) DEFAULT 0 COMMENT '1:taker 2:maker',
    `trade_volume` decimal(18, 10) COMMENT '已成交数量',
    `trade_amount` decimal(18, 10) COMMENT '已成交金额(包括手续费)',
    `fee_kind` tinyint(1) NOT NULL COMMENT '订单创建时的手续费类型 1:百分比 2:绝对值',
    `fee_rate` decimal(18, 10) COMMENT '订单创建时的手续费类型',
    `source` tinyint(1) NOT NULL COMMENT '成交来源 1:otc 2:聚合交易 3:线下交易录入',
    `offline_trade_id` varchar(16) DEFAULT '' COMMENT '线下交易录入id',
    `fee` decimal(18, 10)  COMMENT '成交手续费',
    `turn_over` decimal(18, 10) COMMENT '成交金额',
    `base_currency_price` decimal(18, 10) COMMENT 'pms, 基础币种成交时的价格',
    `target_currency_price` decimal(18, 10) COMMENT 'pms, 目标币种成交时的价格',
    `trade_time` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '成交时间',
    PRIMARY KEY (id),
    INDEX index_user_id (`user_id`),
    INDEX index_order_local_id (`order_local_id`),
    INDEX index_trade_id (`trade_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '成交表';


DROP TABLE IF EXISTS `trade_flow`;
CREATE TABLE `trade_flow` (
    `id` bigint(10) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id` bigint(10) NOT NULL COMMENT '用户id',
    `account_id` bigint(10) NOT NULL COMMENT '账户id',
    `currency`  varchar(16) NOT NULL COMMENT '币种',
    `source` tinyint(1) NOT NULL COMMENT '1:冻结 2:解冻 3:成交 4:线下交易 5:法币出入金  最小原子状态',
    `type` tinyint(1) NOT NULL COMMENT '1:交易 2:交易手续费 3:交易修正 4:交易手续费修正  前台筛选用',
    `order_local_id` varchar(16) NOT NULL COMMENT '本地报单编号',
    `trade_id` varchar(16) NOT NULL COMMENT '成交编号',
    `offline_trade_id` varchar(16) DEFAULT '' COMMENT '线下交易录入id',
    `amount` decimal(18,10)  COMMENT '金额',
    `flow_create_time` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT '流水创建时间',
    PRIMARY KEY (`id`),
    INDEX index_user_id (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '交易流水表';

CREATE TABLE `account`
(
    `id`          bigint      NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_id`     bigint      NOT NULL COMMENT '用户id',
    `currency`    varchar(32) NOT NULL COMMENT '币种',
    `currency_id` varchar(6)  NOT NULL COMMENT '币种id',
--     `available`   decimal(18, 10) COMMENT '可用',
    `frozen`      decimal(18, 10) COMMENT '冻结',
    `balance`     decimal(18, 10) COMMENT '总额',
    `cost_price`  varchar(16) DEFAULT '' COMMENT '成本价',
    PRIMARY KEY (`id`) USING BTREE,
    KEY           `index_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '本地账户表';

-- 以前的admin_member_group改为user_group
DROP TABLE IF EXISTS `user_group`;
CREATE TABLE `user_group`
(
    `group_id`   bigint(20) NOT NULL AUTO_INCREMENT,
    `group_code` varchar(255) NOT NULL,
    `name`       varchar(255) NOT NULL,
    `memo`       varchar(255) DEFAULT NULL,
    PRIMARY KEY (`group_id`),
    UNIQUE KEY `uix_member_group_group_code` (`group_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;