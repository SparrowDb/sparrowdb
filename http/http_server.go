package http

import (
	"fmt"
	"net"
	"net/http"

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

// Start starts HTTP server listener
func (httpServer *HTTPServer) Start() {
	var err error
	httpServer.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", httpServer.Config.HTTPHost, httpServer.Config.HTTPPort))
	if err != nil {
		slog.Fatalf(err.Error())
	}

	handler := NewServeHandler(httpServer.dbManager)

	// register basic middleware, for cors and server name
	httpServer.router.Use(BasicMiddleware())

	// auth group
	authorized := httpServer.router.Group("/")

	// Checks if auth is active, if true, register auth middleware
	if httpServer.Config.AuthenticationActive {
		authorized.Use(AuthMiddleware(func(c *gin.Context) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}))
	}
	httpServer.router.POST("/user/login", handler.userLogin)

	// register routes based on configuration file permission
	if !httpServer.Config.ReadOnly {
		// database create/delete
		authorized.PUT("/api/:dbname", handler.createDatabase)
		authorized.DELETE("/api/:dbname", handler.dropDatabase)

		// image insert/delete
		authorized.PUT("/api/:dbname/:key", handler.uploadData)
		authorized.DELETE("/api/:dbname/:key", handler.deleteData)

		// register script route
		// if :name is "_all" it will retrieve all scripts
		authorized.GET("/script/:name", getScriptList)
		authorized.POST("/script/:name", saveScript)
		authorized.DELETE("/script/:name", deleteScript)
	}

	// if :dbname is "_all" it will retrieve all databases or dbname
	// is a valid database name, it will retrive database information
	authorized.GET("/api/:dbname", handler.infoDatabase)

	// get image information by database/image_key
	authorized.GET("/api/:dbname/:key", handler.getDataInfo)

	// get image by database/image_key
	httpServer.router.GET("/g/:dbname/:key", handler.get)

	httpServer.router.GET("/ping", handler.ping)
	httpServer.router.OPTIONS("/*cors", func(c *gin.Context) {})

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
