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
	SYMBOL_SUB      = "sub_symbol"
	SYMBOL_SUB_OK   = "sub_symbol_ok"
	SYMBOL_UNSUB    = "unsub_symbol"
	SYMBOL_UNSUB_OK = "unsub_symbol_ok"

	DEPTH_SUB      = "depth_sub"
	DEPTH_SUB_OK   = "depth_sub_ok"
	DEPTH_UNSUB    = "depth_unsub"
	DEPTH_UNSUB_OK = "depth_unsub_ok"

	TRADE_SUB      = "trade_sub"
	TRADE_SUB_OK   = "trade_sub_ok"
	TRADE_UNSUB    = "trade_unsub"
	TRADE_UNSUB_OK = "trade_unsub_ok"

	KLINE_SUB       = "kline_sub"
	KLINE_SUB_OK    = "kline_sub_ok"
	KLINE_UNSUMB    = "kline_unsub"
	KLINE_UNSUMB_OK = "kline_unsub_ok"

	SYMBOL_UPDATE = "symbol_update"
	DEPTH_UPDATE  = "depth_update"
	TRADE_UPATE   = "trade_update"
	KLINE_UPATE   = "kline_update"

	KLINE_HIST = "History"
	KLINE_REAL = "RealTime"

	HEARTBEAT = "heartbeat"
	ERROR     = "error"
)
