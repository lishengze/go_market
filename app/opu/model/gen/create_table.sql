CREATE TABLE IF NOT EXISTS `account`
(
    `id`               VARCHAR(64)  NOT NULL COMMENT 'id',
    `alias`            VARCHAR(64)  NOT NULL COMMENT 'alias',
    `key`              VARCHAR(128) NOT NULL COMMENT 'key',
    `secret`           VARCHAR(128) NOT NULL COMMENT 'secret',
    `passphrase`       VARCHAR(128) NOT NULL COMMENT 'passphrase',
    `sub_account_name` VARCHAR(128) NOT NULL COMMENT 'sub_account_name',
    `exchange`         VARCHAR(32)  NOT NULL COMMENT 'exchange',
    `create_time`      DATETIME(6)  NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `update_time`      DATETIME(6)  NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `alias` (`alias`),
    UNIQUE KEY `alias-key` (`alias`, `key`)

) ENGINE = InnoDB,
  CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT 'account';

CREATE TABLE IF NOT EXISTS `order`
(
    `id`              VARCHAR(64) NOT NULL COMMENT 'id, 同时是报给交易所的 order id',
    `account_id`      VARCHAR(64) NOT NULL COMMENT 'account_id',
    `account_alias`   VARCHAR(64) NOT NULL COMMENT 'account_alias',
    `client_order_id` VARCHAR(64) NOT NULL COMMENT 'client_order_id',
    `ex_order_id`     VARCHAR(64) NOT NULL COMMENT 'ex_order_id',
    `api_type`        VARCHAR(32) NOT NULL COMMENT 'api_type',
    `side`            VARCHAR(16) NOT NULL COMMENT 'side',
    `status`          VARCHAR(16) NOT NULL COMMENT 'status',
    `volume`          VARCHAR(32) NOT NULL COMMENT 'volume',
    `filled_volume`   VARCHAR(32) NOT NULL COMMENT 'filled_volume',
    `price`           VARCHAR(32) NOT NULL COMMENT 'price',
    `tp`              VARCHAR(16) NOT NULL COMMENT 'type',
    `std_symbol`      VARCHAR(32) NOT NULL COMMENT 'std_symbol',
    `ex_symbol`       VARCHAR(32) NOT NULL COMMENT 'ex_symbol',
    `exchange`        VARCHAR(32) NOT NULL COMMENT 'exchange',
    `reject_reason`   TEXT        NOT NULL COMMENT 'reject_reason',
    `send_flag`       VARCHAR(16) NOT NULL DEFAULT 'UNSENT' COMMENT 'send_flag,表示订单是否发送至交易所 UNSENT|SENT',
    `cancel_flag`     VARCHAR(16) NOT NULL DEFAULT 'UNSET' COMMENT 'cancel_flag,表示客户是否下达撤单指令 UNSET|SET',
    `last_sync_time`  DATETIME(6) NULL     DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'last_sync_time,上次和交易所同步状态的时间',
    `create_time`     DATETIME(6) NULL     DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `update_time`     DATETIME(6) NULL     DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `account_id-client_order_id` (`account_id`, `client_order_id`)
) ENGINE = InnoDB,
  CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT 'order';


CREATE TABLE IF NOT EXISTS `trade`
(
    `id`           VARCHAR(64) NOT NULL COMMENT 'id',
    `order_id`     VARCHAR(64) NOT NULL COMMENT 'order_id',
    `ex_trade_id`  VARCHAR(64) NOT NULL COMMENT 'ex_trade_id',
    `exchange`     VARCHAR(32) NOT NULL COMMENT 'exchange',
    `std_symbol`   VARCHAR(32) NOT NULL COMMENT 'std_symbol',
    `liquidity`    VARCHAR(32) NOT NULL COMMENT 'liquidity',
    `side`         VARCHAR(16) NOT NULL COMMENT 'side',
    `volume`       VARCHAR(32) NOT NULL COMMENT 'volume',
    `price`        VARCHAR(32) NOT NULL COMMENT 'price',
    `fee`          VARCHAR(32) NOT NULL COMMENT 'fee',
    `fee_currency` VARCHAR(32) NOT NULL COMMENT 'fee_currency',
    `trade_time`   DATETIME(6) NOT NULL COMMENT 'trade_time',
    `create_time`  DATETIME(6) NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `update_time`  DATETIME(6) NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `order_id-ex_trade_id` (`order_id`, `ex_trade_id`)
) ENGINE = InnoDB,
  CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT 'trade';


CREATE TABLE IF NOT EXISTS `symbol`
(
    `id`             VARCHAR(64) NOT NULL COMMENT 'id',
    `type`           VARCHAR(32) NOT NULL COMMENT 'type',
    `api_type`       VARCHAR(32) NOT NULL COMMENT 'api_type',
    `std_symbol`     VARCHAR(64) NOT NULL COMMENT 'std_symbol',
    `ex_format`      VARCHAR(64) NOT NULL COMMENT 'std_symbol',
    `base_currency`  VARCHAR(64) NOT NULL COMMENT 'base_currency',
    `quote_currency` VARCHAR(64) NOT NULL COMMENT 'quote_currency',
    `exchange`       VARCHAR(32) NOT NULL COMMENT 'exchange',
    `volume_scale`   VARCHAR(64) NOT NULL COMMENT 'volume_scale',
    `price_scale`    VARCHAR(64) NOT NULL COMMENT 'price_scale',
    `min_volume`     VARCHAR(64) NOT NULL COMMENT 'min_volume',
    `contract_size`  VARCHAR(64) NOT NULL COMMENT 'contract_size',
    `create_time`    DATETIME(6) NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
    `update_time`    DATETIME(6) NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `std_symbol-exchange` (`std_symbol`, `exchange`),
    UNIQUE KEY `ex_format-api_type-exchange` (`ex_format`, `api_type`, `exchange`)

) ENGINE = InnoDB,
  CHARSET = utf8mb4
  COLLATE = utf8mb4_bin COMMENT 'symbol';
