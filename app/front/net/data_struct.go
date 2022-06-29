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

   if (js["type"].get<string>() == KLINE_SUB)
   {
       process_kline_req(ori_msg, socket_id, ws_safe);
   }

   if (js["type"].get<string>() == TRADE)
   {
       process_trade_req(ori_msg, socket_id, ws_safe);
   }
*/
const (
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
)
