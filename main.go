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
	totalProcs      = runtime.NumCPU()
	configPathFlag  = flag.String("config", "./config/", "Description")
	configProcsFlag = flag.Int("j", totalProcs, "Description")
	instance        *Instance
)

// Instance holds SparrowDb instance configuration
type Instance struct {
	sparrowConfig  *db.SparrowConfig
	databaseConfig *db.DatabaseConfig
	dbManager      *db.DBManager
	httpServer     http.HTTPServer
	wsServer       monitor.WSServer
	serviceManager db.ServiceManager
}

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
	instance.serviceManager.StopAll()
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

	// validate flag processors
	if *configProcsFlag > totalProcs || *configProcsFlag < 0 {
		log.Fatalf("Invalid number of processors: %d, max of %d", *configProcsFlag, totalProcs)
	}

	log.Printf("%s v%s", "SparrowDB", Version)
	log.Printf("Cores: %d", *configProcsFlag)
	runtime.GOMAXPROCS(*configProcsFlag)

	instance = &Instance{}
	instance.sparrowConfig = db.NewSparrowConfig(*configPathFlag)
	instance.databaseConfig = db.NewDatabaseConfig(*configPathFlag)
	log.Printf("Database Mode: %s", instance.sparrowConfig.GetStringMode())

	instance.dbManager = db.NewDBManager(instance.sparrowConfig, instance.databaseConfig)
	instance.dbManager.LoadDatabases()

	monitor.StartMonitor()

	instance.httpServer = http.NewHTTPServer(instance.sparrowConfig, instance.dbManager)
	instance.wsServer = monitor.NewWebSocketServer(instance.sparrowConfig)

	instance.serviceManager = db.NewServiceManager()
	instance.serviceManager.AddService("wsServer", &instance.wsServer)
	instance.serviceManager.AddService("httpServer", &instance.httpServer)
	instance.serviceManager.StartAll()
}
