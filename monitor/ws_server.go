package monitor

import (
	"net"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/slog"

	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

// WSServer holds WebSocket server configuration
type WSServer struct {
	Config   *db.SparrowConfig
	listener net.Listener

	Notifier chan []byte
}

func (wss *WSServer) webHandler(ws *websocket.Conn) {
	var in []byte

	ws.Write(getJSON())

	go func() {
		for {
			if err := websocket.Message.Receive(ws, &in); err != nil {
				break
			}
		}
	}()

	for {
		ws.Write(<-wss.Notifier)
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
		Config:   config,
		Notifier: make(chan []byte),
	}
}
