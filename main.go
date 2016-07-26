package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/sparrowdb/db"
	"github.com/sparrowdb/http"
	"github.com/sparrowdb/monitor"
	"github.com/sparrowdb/util"
)

const (
	// Version SparrowDb version
	Version = "0.0.1"
)

var (
	//configPathFlag = flag.String("config", "/home/mauricio/Sources/sparrow/sparrow_conf/conf1/", "Description")
	configPathFlag = flag.String("config", "./config/", "Description")
)

func checkAndCreateDefaultDirs() {
	dirs := []string{"log", "config", "data", "plugin"}
	for _, val := range dirs {
		if _, err := os.Stat(val); os.IsNotExist(err) {
			util.CreateDir(val)
		}
	}
}

func init() {
	createPIDfile()

	// Configure signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go handleSignal(c)
}

func handleSignal(c chan os.Signal) {
	<-c
	log.Printf("Quinting SparrowDB")
	os.Exit(1)
}

func createPIDfile() {
	p := strconv.Itoa(os.Getpid())
	ioutil.WriteFile("sparrow.pid", []byte(p), 0644)
}

func main() {
	flag.Parse()

	checkAndCreateDefaultDirs()

	log.Printf("%s v%s", "SparrowDB", Version)
	log.Printf("Cores: %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	sparrowConfig := db.NewSparrowConfig(*configPathFlag)
	databaseConfig := db.NewDatabaseConfig(*configPathFlag)

	dbManager := db.NewDBManager(sparrowConfig, databaseConfig)
	dbManager.LoadDatabases()

	monitor.StartMonitor()

	httpServer := http.NewHTTPServer(sparrowConfig, dbManager)
	wsServer := monitor.NewWebSocketServer(sparrowConfig)

	serviceManager := db.NewServiceManager()
	serviceManager.AddService("wsServer", &wsServer)
	serviceManager.AddService("httpServer", &httpServer)
	serviceManager.StartAll()
}
