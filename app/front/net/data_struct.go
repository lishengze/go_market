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

   if (js["type"].get<string>() == KLINE_UPDATE)
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
	SYMBOL_LIST        = "symbol_list"
	SYMBOL_UPDATE      = "symbol_update"
	KLINE_RSP          = "kline_rsp"
	KLINE_UPDATE       = "kline_update"
	HEARTBEAT          = "heartbeat"
	TRADE              = "trade"
)
