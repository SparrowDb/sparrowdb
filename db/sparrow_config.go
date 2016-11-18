package db

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/SparrowDb/sparrowdb/errors"
	"github.com/SparrowDb/sparrowdb/slog"
)

const (
	// DefaultSparrowConfigFile is the default configuration file
	DefaultSparrowConfigFile = "sparrow.xml"
)

// SparrowConfig holds general configuration of SparrowDB
type SparrowConfig struct {
	NodeName             string  `xml:"node_name"`
	HTTPPort             string  `xml:"http_port"`
	HTTPHost             string  `xml:"http_host"`
	WSPort               string  `xml:"ws_port"`
	WSHost               string  `xml:"ws_host"`
	ReadOnly             bool    `xml:"read_only"`
	MaxDataLogSize       uint64  `xml:"max_datalog_size"`
	MaxCacheSize         uint64  `xml:"max_cache_size"`
	BloomFilterFp        float32 `xml:"bloomfilter_fpp"`
	CronExp              string  `xml:"dataholder_cron_compaction"`
	Path                 string  `xml:"data_file_directory"`
	SnapshotPath         string  `xml:"snapshot_path"`
	TokenActive          bool    `xml:"generate_token"`
	AuthenticationActive bool    `xml:"enable_authentication"`
	UserExpire           int     `xml:"user_expire"`
}

// NewSparrowConfig return configuration from file
func NewSparrowConfig(filePath string) *SparrowConfig {
	filePath = filePath + DefaultSparrowConfigFile

	xmlFile, err := os.Open(filePath)
	if err != nil {
		slog.Fatalf(errors.ErrFileNotFound.Error(), filePath)
	}

	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)

	cfg := SparrowConfig{}

	if err := xml.Unmarshal(data, &cfg); err != nil {
		slog.Fatalf(errors.ErrParseFile.Error(), filePath)
	}

	return &cfg
}
