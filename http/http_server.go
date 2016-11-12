package http

import (
	"fmt"
	"net"
	"net/http"

	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/SparrowDb/sparrowdb/spql"
	"github.com/gin-gonic/gin"
)

// HTTPServer holds HTTP server configuration and routes
type HTTPServer struct {
	Config        *db.SparrowConfig
	router        *gin.Engine
	dbManager     *db.DBManager
	queryExecutor *spql.QueryExecutor
	listener      net.Listener
}

func (httpServer *HTTPServer) basicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Server", "SparrowDb")
		c.Next()
	}
}

// Start starts HTTP server listener
func (httpServer *HTTPServer) Start() {
	var err error
	httpServer.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", httpServer.Config.HTTPHost, httpServer.Config.HTTPPort))
	if err != nil {
		slog.Fatalf(err.Error())
	}

	handler := NewServeHandler(httpServer.dbManager, httpServer.queryExecutor)

	httpServer.router.Use(httpServer.basicMiddleware())

	r, w, q := httpServer.Config.GetMode()

	httpServer.router.GET("/ping", handler.ping)

	if r == true {
		httpServer.router.GET("/g/:dbname/:key", handler.get)
	}
	if w == true {
		httpServer.router.POST("/upload", handler.upload)
	}
	if q == true {
		httpServer.router.POST("/query", handler.serveQuery)
	}

	http.Serve(httpServer.listener, httpServer.router)
}

// Stop stops HTTP server listener
func (httpServer *HTTPServer) Stop() {
	slog.Infof("Stopping HTTP Server")
	httpServer.listener.Close()
}

// NewHTTPServer returns new HTTPServer
func NewHTTPServer(config *db.SparrowConfig, dbm *db.DBManager) HTTPServer {
	gin.SetMode(gin.ReleaseMode)

	return HTTPServer{
		Config:        config,
		dbManager:     dbm,
		queryExecutor: spql.NewQueryExecutor(dbm),
		router:        gin.New(),
	}
}
