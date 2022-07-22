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

func GetSubTradeRspMsg(req_arrive_time int64, symbol_list []string) []byte {
	heartbeat_map := map[string]interface{}{
		"symbol":          symbol_list,
		"req_arrive_time": util.TimeStrFromInt(req_arrive_time),
		"type":            net.TRADE_SUB_OK,
	}
	rst, err := json.Marshal(heartbeat_map)

	if err != nil {
		logx.Errorf("GetSubTradeRspMsg: %+v", err)
		return nil
	}
	return rst
}

func GetSubDepthRspMsg(req_arrive_time int64, symbol_list []string) []byte {
	heartbeat_map := map[string]interface{}{
		"symbol":          symbol_list,
		"req_arrive_time": util.TimeStrFromInt(req_arrive_time),
		"type":            net.DEPTH_SUB_OK,
	}
	rst, err := json.Marshal(heartbeat_map)

	if err != nil {
		logx.Errorf("GetSubSymbolRspMsg: %+v", err)
		return nil
	}
	return rst
}

func GetSubKlineRspMsg(req_kline *datastruct.ReqHistKline) []byte {
	heartbeat_map := map[string]interface{}{
		"symbol":          req_kline.Symbol,
		"resolution":      req_kline.Frequency,
		"count":           req_kline.Count,
		"req_arrive_time": util.TimeStrFromInt(req_kline.ReqArriveTime),
		"type":            net.KLINE_SUB_OK,
	}
	rst, err := json.Marshal(heartbeat_map)

	if err != nil {
		logx.Errorf("GetSubKlineRspMsg: %+v", err)
		return nil
	}
	return rst
}
