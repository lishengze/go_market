package net

/*
   if (js["type"].get<string>() == "sub_symbol")
   {
       process_depth_req(ori_msg, socket_id, ws_safe);

       process_trade_req(ori_msg, socket_id, ws_safe);
   }

   if (js["type"].get<string>() == HEARTBEAT)
   {
       process_heartbeat(socket_id, ws_safe);
   }

   if (js["type"].get<string>() == KLINE_UPDATE_SUB)
   {
       process_kline_req(ori_msg, socket_id, ws_safe);
   }

   if (js["type"].get<string>() == TRADE)
   {
       process_trade_req(ori_msg, socket_id, ws_safe);
   }
*/
const (
	MARKET_DATA_UPDATE = "market_data_update"
	SYMBOL_SUB         = "sub_symbol"
	SYMBOL_UNSUB       = "unsub_symbol"
	SYMBOL_LIST        = "symbol_list"
	SYMBOL_UPDATE      = "symbol_update"

	DEPTH_SUB           = "depth_sub"
	DEPTH_UNSUB         = "depth_unsub"
	TRADE_SUB           = "trade_sub"
	TRADE_UNSUB         = "trade_unsub"
	KLINE_UPDATE_SUB    = "kline_sub"
	KLINE_UPDATE_UNSUMB = "kline_unsub"

	KLINE_RSP = "kline_rsp"

	HEARTBEAT = "heartbeat"
	TRADE     = "trade"
)
