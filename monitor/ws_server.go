package monitor

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/slog"

	"golang.org/x/net/websocket"
)

// WSServer holds WebSocket server configuration
type WSServer struct {
	Config   *db.SparrowConfig
	listener net.Listener
}

func (wss *WSServer) webHandler(ws *websocket.Conn) {
	for {
		b := MetricToJSON()
		websocket.Message.Send(ws, string(b))
		time.Sleep(1 * time.Second)
	}
}

// Start starts WebSocket server listener
func (wss *WSServer) Start() {
	var err error
	wss.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", wss.Config.WSHost, wss.Config.WSPort))
	if err != nil {
		slog.Fatalf(err.Error())
	}

	slog.Infof("Starting WebSocket Server %s:%s", wss.Config.WSHost, wss.Config.WSPort)

	http.Handle("/", websocket.Handler(wss.webHandler))

	http.Serve(wss.listener, nil)
}

// Stop stops WebSocket server listener
func (wss *WSServer) Stop() {
	slog.Infof("Stopping WebSocket Server")
	wss.listener.Close()
}

// NewWebSocketServer returns new WSServer
func NewWebSocketServer(config *db.SparrowConfig) WSServer {
	return WSServer{
		Config: config,
	}
}
