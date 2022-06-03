[TOC]

## 修订记录

| 版本 | 修订日期   | 修订说明 |
| ---- | ---------- | -------- |
| v0.1 | 2022-04-11 | 发布初版 |

## 文档说明

用户询价接口约定  Demo

## 基本约定

### 报文交换格式

**Json**

### 通讯协议

#### 协议

**http**

## oct_quote 用户询价接口说明

#### 地址：/trading/api/otcQuote

#### 方法：POST

#### 请求

| 字段       | 是否必须 | 类型   | 说明     |
| ---------- | -------- | ------ | -------- |
| direction     | M        | string | 方向 |
| symbol | M        | string | 币对   |
| type | M        | string | 类型   |
| turnover | M        | string | 交易额   |

#### 返回

| 字段       | 是否必须 | 类型   | 说明         |
| ---------- | -------- | ------ | ------------ |
| price     | M        | string | 价格 |
| symbol | M        | string | 币对   |
| otc_number | M        | string | 数量   |
| otc_turnover | M        | string | 交易额   |




