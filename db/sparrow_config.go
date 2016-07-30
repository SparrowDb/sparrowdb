package db

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

const (
	// DefaultSparrowConfigFile is the default configuration file
	DefaultSparrowConfigFile = "sparrow.xml"
)

// SparrowConfig holds general configuration of SparrowDB
type SparrowConfig struct {
	NodeName string `xml:"node_name"`
	HTTPPort string `xml:"http_port"`
	HTTPHost string `xml:"http_host"`
	WSPort   string `xml:"ws_port"`
	WSHost   string `xml:"ws_host"`
	Mode     string `xml:"mode"`
}

func (sc *SparrowConfig) isValid() bool {
	reg := regexp.MustCompile("^(WR|RW|R|W)$")
	return reg.MatchString(sc.Mode)
}

// GetMode returns string describind SparrowDB
// storage mode
func (sc *SparrowConfig) GetMode() string {
	var r string
	switch sc.Mode {
	case "R":
		r = "Read"
	case "W":
		r = "Write"
	case "RW", "WR":
		r = "Read and Write"
	}
	return r
}

// NewSparrowConfig return configuration from file
func NewSparrowConfig(filePath string) *SparrowConfig {
	filePath = filePath + DefaultSparrowConfigFile

	xmlFile, _ := os.Open(filePath)

	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)

	cfg := SparrowConfig{}
	xml.Unmarshal(data, &cfg)

	if !cfg.isValid() {
		log.Fatalf("Not valid SparrowDB mode, it must be [R]ead, [W]write or [RW]read-write")
	}

	return &cfg
}
