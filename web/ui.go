package web

import (
	"fmt"
	"html/template"
	"net"
	_http "net/http"
	"os"

	"path/filepath"

	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/http"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/gin-gonic/gin"
)

func buildTemplate(files ...string) *template.Template {
	var pfiles []string
	pwd, _ := os.Getwd()
	for _, file := range files {
		pfiles = append(pfiles, filepath.Join(pwd, "web", "templates", file))
	}
	return template.Must(template.ParseFiles(pfiles...))
}

// UIServer holds HTTP server configuration and routes
type UIServer struct {
	Config   *db.SparrowConfig
	router   *gin.Engine
	listener net.Listener
}

// Start starts HTTP server listener
func (s *UIServer) Start() {
	var err error
	//hostAddr := fmt.Sprintf("127.0.0.1:%s", s.Config.HTTPPort)

	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", s.Config.AdminHost, s.Config.AdminPort))
	if err != nil {
		slog.Fatalf(err.Error())
	}

	s.router.Use(http.BasicMiddleware())

	pwd, _ := os.Getwd()
	//s.router.LoadHTMLGlob(filepath.Join(pwd, "web", "templates/*"))
	s.router.StaticFS("/", _http.Dir(filepath.Join(pwd, "web", "static")))
	s.router.OPTIONS("/*cors", func(c *gin.Context) {})

	/*s.router.GET("/", func(c *gin.Context) {
		s.router.SetHTMLTemplate(buildTemplate("base.html"))
		c.HTML(_http.StatusOK, "base", gin.H{
			"hostAddr": hostAddr,
		})
	})*/

	_http.Serve(s.listener, s.router)
}

// Stop stops HTTP server listener
func (s *UIServer) Stop() {
	slog.Infof("Stopping Admin Server")
}

// NewUIServer returns new UI server
func NewUIServer(config *db.SparrowConfig) UIServer {
	gin.SetMode(gin.ReleaseMode)
	return UIServer{
		Config: config,
		router: gin.New(),
	}
}
