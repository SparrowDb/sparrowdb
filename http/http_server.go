package http

import (
	"net"
	"net/http"
	"strings"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/monitor"
	"github.com/sparrowdb/slog"
	"github.com/sparrowdb/spql"
)

// HTTPServer holds HTTP server configuration and routes
type HTTPServer struct {
	Config        *db.SparrowConfig
	mux           *http.ServeMux
	dbManager     *db.DBManager
	routers       map[string]*controllerInfo
	queryExecutor *spql.QueryExecutor
	listener      net.Listener
}

type controllerInfo struct {
	route      string
	httpMethod string
	method     func(request *RequestData)
}

func (httpServer *HTTPServer) add(c *controllerInfo) {
	parts := strings.Split(c.route[1:], "/")
	httpServer.routers[parts[0]] = c
}

func (httpServer *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path[1:], "/")

	if parts[0] == "favicon.ico" {
		return
	}

	monitor.IncHTTPRequests()

	if controller, ok := httpServer.routers[parts[0]]; ok {
		parts := strings.Split(r.URL.Path[1:], "/")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		controller.method(&RequestData{responseWriter: w, request: r, params: parts[1:]})
	}
}

// Start starts HTTP server listener
func (httpServer *HTTPServer) Start() {
	var err error
	httpServer.listener, err = net.Listen("tcp", ":"+httpServer.Config.HTTPPort)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	slog.Infof("Starting HTTP Server %s:%s", httpServer.Config.HTTPHost, httpServer.Config.HTTPPort)

	handler := NewServeHandler(httpServer.dbManager, httpServer.queryExecutor)

	r, w, q := httpServer.Config.GetMode()
	if r == true {
		httpServer.add(&controllerInfo{route: "/g", httpMethod: "GET", method: handler.get})
	}
	if w == true {
		httpServer.add(&controllerInfo{route: "/upload", httpMethod: "POST", method: handler.upload})
	}
	if q == true {
		httpServer.add(&controllerInfo{route: "/query", httpMethod: "POST", method: handler.serveQuery})
	}

	httpServer.mux.Handle("/", httpServer)

	http.Serve(httpServer.listener, httpServer.mux)
}

// Stop stops HTTP server listener
func (httpServer *HTTPServer) Stop() {
	slog.Infof("Stopping HTTP Server")
	httpServer.listener.Close()
}

// NewHTTPServer returns new HTTPServer
func NewHTTPServer(config *db.SparrowConfig, dbm *db.DBManager) HTTPServer {
	return HTTPServer{
		Config:        config,
		dbManager:     dbm,
		queryExecutor: spql.NewQueryExecutor(dbm),
		mux:           http.NewServeMux(),
		routers:       make(map[string]*controllerInfo),
	}
}
