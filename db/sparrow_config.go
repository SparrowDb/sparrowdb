package db

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/sparrowdb/slog"
)

const (
	// DefaultSparrowConfigFile is the default configuration file
	DefaultSparrowConfigFile = "sparrow.xml"
)

// SparrowConfig holds general configuration of SparrowDB
type SparrowConfig struct {
	NodeName       string  `xml:"node_name"`
	HTTPPort       string  `xml:"http_port"`
	HTTPHost       string  `xml:"http_host"`
	WSPort         string  `xml:"ws_port"`
	WSHost         string  `xml:"ws_host"`
	Mode           string  `xml:"mode"`
	MaxDataLogSize uint64  `xml:"max_datalog_size"`
	MaxCacheSize   uint64  `xml:"max_cache_size"`
	BloomFilterFp  float32 `xml:"bloomfilter_fpp"`
	CronExp        string  `xml:"dataholder_cron_compaction"`
	Path           string  `xml:"data_file_directory"`
	TokenActive    bool    `xml:"generate_token"`
}

func (sc *SparrowConfig) isValid() bool {
	reg := regexp.MustCompile("^([Q])|([W])|([R])$")
	return reg.MatchString(sc.Mode)
}

// GetMode returns bool for each SparrowDB mode
func (sc *SparrowConfig) GetMode() (read bool, write bool, query bool) {
	if strings.Contains(sc.Mode, "R") {
		read = true
	}

	if strings.Contains(sc.Mode, "W") {
		write = true
	}

	if strings.Contains(sc.Mode, "Q") {
		query = true
	}
	return
}

// GetStringMode returns string describind SparrowDB
// storage mode
func (sc *SparrowConfig) GetStringMode() string {
	var r string
	if strings.Contains(sc.Mode, "R") {
		r += "Read "
	}

	if strings.Contains(sc.Mode, "W") {
		r += "Write "
	}

	if strings.Contains(sc.Mode, "Q") {
		r += "Query "
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
		slog.Fatalf("Not valid SparrowDB mode, it must be [R]ead, [W]write or [RW]read-write")
	}

	return &cfg
}
