package monitor

import (
	"log"
	"net/http"
	"time"

	"github.com/sparrowdb/db"

	"golang.org/x/net/websocket"
)

type WSServer struct {
	Config *db.SparrowConfig
}

func (wss *WSServer) webHandler(ws *websocket.Conn) {
	for {
		b := MetricToJSON()
		websocket.Message.Send(ws, string(b))
		time.Sleep(1 * time.Second)
	}
}

func (wss *WSServer) Start() {
	log.Printf("Starting WebSocket Server %s:%s", wss.Config.WSHost, wss.Config.WSPort)

	http.Handle("/", websocket.Handler(wss.webHandler))
	http.ListenAndServe(":"+wss.Config.WSPort, nil)
}

func (wss *WSServer) Stop() {
}

func NewWebSocketServer(config *db.SparrowConfig) WSServer {
	return WSServer{
		Config: config,
	}
}
