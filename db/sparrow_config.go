package db

import (
	"encoding/xml"
	"io/ioutil"
	"os"
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
}

// NewSparrowConfig return configuration from file
func NewSparrowConfig(filePath string) *SparrowConfig {
	filePath = filePath + DefaultSparrowConfigFile

	xmlFile, _ := os.Open(filePath)

	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)

	cfg := SparrowConfig{}
	xml.Unmarshal(data, &cfg)

	return &cfg
}
