package ws_server

import (
	"encoding/json"
	"market_server/app/front/net"
	"market_server/common/datastruct"
	"market_server/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

func GetHeartbeatMsg() []byte {

	heartbeat_map := map[string]interface{}{
		"time": util.UTCNanoTime(),
		"type": net.HEARTBEAT,
	}
	rst, err := json.Marshal(heartbeat_map)

	if err != nil {
		logx.Errorf("GetHeartbeatMsg: %+v", err)
		return nil
	}
	return rst
}

func GetSubTradeRspMsg(req_trade *datastruct.ReqTrade, symbol_list []string) []byte {
	heartbeat_map := map[string]interface{}{
		"time":            util.UTCNanoTime(),
		"symbol":          symbol_list,
		"req_arrive_time": util.TimeStrFromInt(req_trade.ReqArriveTime),
		"type":            net.HEARTBEAT,
	}
	rst, err := json.Marshal(heartbeat_map)

	if err != nil {
		logx.Errorf("GetSubTradeRspMsg: %+v", err)
		return nil
	}
	return rst
}

func GetSubSymbolRspMsg(req_depth *datastruct.ReqDepth, symbol_list []string) []byte {
	heartbeat_map := map[string]interface{}{
		"time":            util.UTCNanoTime(),
		"symbol":          req_depth.Symbol,
		"req_arrive_time": util.TimeStrFromInt(req_depth.ReqArriveTime),
		"type":            net.HEARTBEAT,
	}
	rst, err := json.Marshal(heartbeat_map)

	if err != nil {
		logx.Errorf("GetSubTradeRspMsg: %+v", err)
		return nil
	}
	return rst
}
