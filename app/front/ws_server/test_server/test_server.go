package main

import (
	"flag"
	"fmt"
	"log"
	"market_server/app/front/net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	// "golang.org/x/net/websocket"
)

func basic_test() {
	// websocket.Server()
}

var addr = flag.String("addr", "127.0.0.1:8114", "http service address")

var upgrader = websocket.Upgrader{} // use default options

// type WS_INFO struct {
// 	Conn *websocket.Conn
// 	ID   int64
// }

func echo(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%+v, Echo: %+v \n", time.Now(), *r)

	c, err := upgrader.Upgrade(w, r, nil)

	ws := net.NewWSInfo(c)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.Conn.ReadMessage()
		if err != nil {
			log.Printf("read, type: %d, err: %+v\n", mt, err)
			break
		}
		log.Printf("ws: %d, recv type: %d, msg: %s", ws.ID, mt, message)
		err = ws.Conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/trading/marketws", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func TestWS() {
	logx.Info("TestWS")
	fmt.Println("TestWS")
}
