# 行情前置服务器

## 接口说明

### 服务类型

1. websoket 服务  
   a. websoket 的请求与推送都是以 json 字符串的形式完成. 

   b. 请求与推送的 json 中有固定的 type 属性说明请求与推送的数据格式.

   c. 心跳机制:  
    服务器会定时发送心跳数据包给客户端，type:"heartbeat", 客户端只需回复 {"type":"heartbeat"} 即可保持连接.


### 数据说明
行情发布的主要数据是 depth列表，k线, 实时trade，24小时涨跌幅等数据;

1. symbol 列表  
   symbol信息是所有行情请求的基础，也是由websocket server 提供。
   1) 初始化： 当websocket 建立时， 服务端会自动地将当前存储的 symbol 发送给客户端， 无需请求。 
        json: {"type": "symbol_update", "symbol": ["symbol1", "symbol2" ]}
   2) 更新:  
       当有symbol加入时，websocket 会推送更新json 过来
       json: {"type": "symbol_update", "symbol": ["symbol1", "symbol2" ]}

2. depth数据 
    1) 同时订阅  symbol 对应的 depth, trade数据，24小时涨跌幅数据; 
    {  
        "type":"sub_symbol",  
        "symbol":["symbolName"]  //需要订阅的symbole数组
    }  

    2) 单独对 depth 的订阅  
    {  
        "type":"depth_sub",  
        "symbol":["symbolName"]   
    }

    3) 取消订阅  
      1). 同时取消对 depth, trade 的订阅。  
        {  
            "type":"unsub_symbol",  
            "symbol":["symbolName"]   
        }  
      2). 单独取消对 depth 的订阅  
        {
            "type":"depth_unsub",  
            "symbol":["symbolName"]   
        }

    4) websocket 推送的depth数据:  
    {  
        "ask_length":numb,  
        "asks":[[price, volume, accumulatedVolume]...],  // 每个数组原子按顺序存储 价格，成交量，累积成交量等信息;  
        "bid_length":numb,  
        "bids":[[price, volume, accumulatedVolume]...],  // bids 是逆序
        "symbol":"symbolName",  
        "exchange":"",  
        "seqno":0,  // 当前行情序列号  
        "tick":0,   // 时间戳  
        "type":"depth_update"     // type 类型   
    }



3. Trade 数据与 24小时涨跌幅数据
    Trade数据和24小时涨跌幅数据是在订阅 depth 数据时一起订阅的。因为24小时涨跌幅数据是依赖trade 数据更新的，所以这两类数据是放在同一个json 数据结构中回报。  
    1)  回报的数据格式:
    {
        "type":"trade_update",
        "symbol":"",
        "price":"",
        "volume":"",
        "change":"",
        "change_rate":"",
        "high":"",
        "low":"",
    }

    2) 单独对 trade 的订阅  
        {  
            "type":"trade_sub",  
            "symbol":["symbolName"]   
        }

    3)  单独取消对trade的订阅  
    {  
        "type":"trade_unsub",  
        "symbol":["symbolName"]           
    }


4. K线数据
    1) 发送订阅请求，k线的请求订阅包含了请求的历史数据的元信息，请求的json字符串为:     
    {   
        "type":"kline_sub",    
        "symbol":"symbolName",  
        "start_time":start_time,    // 必须是秒级的UTC时间戳   字符串形式;
        "end_time":end_time,        // 必须是秒级的UTC时间戳  字符串形式;
        "count": "1000",            // 数量和时间，任选一个为准; 同时填写以 Count 为准;  字符串形式;
        "frequency":"60"            // 数据频率，以秒为单位，现在必须是60的整数倍.  字符串形式;
    }  
    历史数据的区间通过 [start_time, end_time] 或者 count 设置; 默认推荐的是按照数量来请求，初次订阅时默认展示的是 1000根1分频的k线数据。

    2) 取消订阅  
        取消对某个币对某个频率的K线订阅  
    {   
        "type":"kline_unsub",    
        "symbol":"symbolName",  
        "frequency":"60"            // 数据频率，以秒为单位，现在必须是60的整数倍.  
    }    

    3) websocket 推送的数据:  
    {  
        "data":[["open":,"high":,"low":,"close":,"volume":,"tick":,]...]  // 每个数组原子储存 open, high, low, close, volume, tick-时间戳 等信息;
        "symbol":"symbolName",    
        "start_time":"",    // 回复数据的开始时间，秒级的UTC时间戳   
        "end_time":0,       // 回复数据的结束时间，秒级的UTC时间戳   
        "frequency":0,      // 请求的时间频率  
        "type":"kline_update"     // type 类型   
        "data_count": ;     // 实际返回的k 线数目;  
        "data_type": ;      // "RealTime", "History"，返回的是实时还是历史数据;
        
    }    

5. 心跳数据
    服务端主动推送，客户端收到返回的机制  
    1) 服务端每隔一段时间，会发送心跳字段给客户端形式为   
    {"time":"2020-12-04 07:41:20.15205969","type":"heartbeat"}.  
    2) 客户端需要回应  {"type": "heartbeat"} 这样的字符串即可.  
    3) 服务端判断失活的依据是，规定时间未收到客户的请求信息-包括心跳回报，这个时间通常是 心跳发送时间的整数被.

6. 特殊说明  
    所有的时间都是以 UTC时间为准;

7. type 汇总说明  
	SYMBOL_SUB   = "sub_symbol"  
	SYMBOL_UNSUB = "unsub_symbol"  
	DEPTH_SUB    = "depth_sub"  
	DEPTH_UNSUB  = "depth_unsub"  
	TRADE_SUB    = "trade_sub"  
	TRADE_UNSUB  = "trade_unsub"  
	KLINE_SUB    = "kline_sub"  
	KLINE_UNSUMB = "kline_unsub"  
  
	SYMBOL_UPDATE = "symbol_update"  
	DEPTH_UPDATE  = "depth_update"  
	TRADE_UPATE   = "trade_update"  
	KLINE_UPATE   = "kline_update"  
  
	HEARTBEAT = "heartbeat"  

8. 错误信息说明  
    若是订阅出现错误，会发送错误数据信息  
    {   
        "type":"error",    
        "info":"err_msg",  
    }        
   
## 更新日志 2022.6.28  
    针对app 访问的特点，对数据的订阅，取消订阅的 type 都做了更新。  
    增加了单独针对 depth 和 trade 的订阅与取消订阅的结构。
    kline 数据的订阅与取消订阅的 type 也做了对应的统一。

## 更新日志 2022.6.29  
    将所有回报的数据的格式做了统一，一致设置为 dataType_update.

## 更新日志 2022.7.7  
    增加错误信息回报

    
