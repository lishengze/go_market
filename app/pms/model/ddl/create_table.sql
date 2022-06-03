DROP TABLE IF EXISTS `portfolio_investments`;
CREATE TABLE `portfolio_investments` (
                                    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                    `name` varchar(128) DEFAULT NOT NULL COMMENT '资产分组名称',
                                    `user_id` bigint(20) DEFAULT  NOT NULL COMMENT '用户id',
                                    `remark` varchar(256) DEFAULT  NOT NULL COMMENT '备注',
                                    `status` tinyint(1) DEFAULT '1' COMMENT '1:默认分组，2:私有分组',
                                    `last_operator` bigint(20) DEFAULT '0' COMMENT '最后一次操作人,0-表示由程序自动创建',
                                    `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
                                    PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

#
# Structure for table "portfolio_investment_items"
#

DROP TABLE IF EXISTS `portfolio_investment_items`;
CREATE TABLE `portfolio_investment_items` (
                                         `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                         `gid` bigint(20) DEFAULT  NOT NULL COMMENT '资产组ID',
                                         `investment` varchar(128) DEFAULT  NOT NULL COMMENT '分组中包含的投资品',
                                         `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                         `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
                                         PRIMARY KEY (`id`),
                                         UNIQUE KEY `gid_investment` (`gid`,`investment`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `investment_snapshots`;
CREATE TABLE `investment_snapshots` (
                                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                        `user_id` bigint(20) DEFAULT  NOT NULL COMMENT '用户id',
                                        `investment` varchar(255)  COMMENT '投资品',
                                        `position` decimal(26,16)  COMMENT '投资品当前持有数量',
                                        `average_price` decimal(26,16)  COMMENT '均价',
                                        `cost_basis` decimal(26,16)  COMMENT '成本',
                                        `last` decimal(26,16)  COMMENT '市价',
                                        `market_value` decimal(26,16) COMMENT '市值',
                                        `market_value_percent` decimal(26,16)  COMMENT '市值占比',
                                        `daily_change` decimal(26,16) COMMENT '日波动',
                                        `daily_profit` decimal(26,16)  COMMENT '日盈亏',
                                        `unrealized_profit` decimal(26,16) COMMENT '未实现盈亏',
                                        `unrealized_profit_percent` decimal(26,16)  COMMENT '未实现盈亏率',
                                        `daily_profit_add_loss` decimal(26,16)  COMMENT '日盈亏dailyP&L',
                                        `realized_profit_add_loss` decimal(26,16)  COMMENT '累计已实现盈亏Realized P&L',
                                        `date_time` date DEFAULT NOT NULL COMMENT '日期',
                                        `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                        `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
                                        PRIMARY KEY (`id`),
                                        KEY `mid_coin` (`member_id`,`date_time`,`investment`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

#
# Structure for table "member_cash_flow"
#

DROP TABLE IF EXISTS `member_cash_flow`;
CREATE TABLE `member_cash_flow` (
                                    `id` bigint(20) NOT NULL AUTO_INCREMENT,
                                    `user_id` bigint(20) DEFAULT  NOT NULL COMMENT '用户id',
                                    `cash_flow` decimal(26,16)  COMMENT '现金流',
                                    `market_value_day_start` decimal(26,16)  COMMENT '日初市值',
                                    `market_value_day_end` decimal(26,16)  COMMENT '日终市值',
                                    `date_time` date DEFAULT  NOT NULL COMMENT '日期',
                                    `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
                                    PRIMARY KEY (`id`),
                                    KEY `member_date` (`member_id`,`date_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

#
# Structure for table "coin_prices"
#

DROP TABLE IF EXISTS `coin_prices`;
CREATE TABLE `coin_prices` (
                               `id` bigint(20) NOT NULL AUTO_INCREMENT,
                               `coin_unit` varchar(255) DEFAULT NOT NULL COMMENT '币种',
                               `price` decimal(26,16)  COMMENT '日出价格',
                               `date_time` date DEFAULT  NOT NULL COMMENT '日期',
                               `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最近更新时间',
                               PRIMARY KEY (`id`),
                               KEY `coin_date` (`date_time`,`coin_unit`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
