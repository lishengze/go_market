package monitor_market

import "github.com/gorilla/websocket"

type WSClient struct {
	config *WSConfig
	client *websocket.Conn
}
