package http

import (
	"fmt"
	"net"
	"net/http"

	"github.com/SparrowDb/sparrowdb/auth"
	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/gin-gonic/gin"
)

// HTTPServer holds HTTP server configuration and routes
type HTTPServer struct {
	Config    *db.SparrowConfig
	router    *gin.Engine
	dbManager *db.DBManager
	listener  net.Listener
}

func (httpServer *HTTPServer) basicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Server", "SparrowDb")
		c.Next()
	}
}

func (httpServer *HTTPServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.ParseFromRequest(c.Request)

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

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

	handler := NewServeHandler(httpServer.dbManager)

	// register basic middleware, for cors and server name
	httpServer.router.Use(httpServer.basicMiddleware())

	// auth group
	authorized := httpServer.router.Group("/")

	// Checks if auth is active, if true, register auth middleware
	// and login route
	if httpServer.Config.AuthenticationActive {
		authorized.Use(httpServer.authMiddleware())
		httpServer.router.POST("/user/login", handler.userLogin)
	}

	// register routes based on configuration file permission
	if !httpServer.Config.ReadOnly {
		authorized.PUT("/api/:dbname", handler.createDatabase)
		authorized.DELETE("/api/:dbname", handler.dropDatabase)
		authorized.GET("/api/:dbname", handler.infoDatabase)

		authorized.PUT("/api/:dbname/:key", handler.uploadData)
		authorized.DELETE("/api/:dbname/:key", handler.deleteData)

		httpServer.router.GET("/api/:dbname/:key", handler.getDataInfo)
	}

	httpServer.router.GET("/g/:dbname/:key", handler.get)
	httpServer.router.GET("/ping", handler.ping)

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
		Config:    config,
		dbManager: dbm,
		router:    gin.New(),
	}
}
