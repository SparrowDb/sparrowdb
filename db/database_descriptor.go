package db

import (
	"encoding/json"
	"encoding/xml"
)

// XMLDatabaseList holds root node and DatabaseDescriptor
// list
type XMLDatabaseList struct {
	XMLName   xml.Name             `xml:"databases"`
	Databases []DatabaseDescriptor `xml:"database"`
}

// DatabaseDescriptor holds database configuration
type DatabaseDescriptor struct {
	XMLName        xml.Name `xml:"database"`
	Name           string   `xml:"name"`
	MaxDataLogSize uint64   `xml:"max_datalog_size"`
	MaxCacheSize   uint64   `xml:"max_cache_size"`
	BloomFilterFp  float32  `xml:"bloomfilter_fpp"`
	CronExp        string   `xml:"dataholder_cron_compaction"`
	Path           string   `xml:"path"`
	TokenActive    bool     `xml:"generate_token"`
	Mode           string   `xml:"mode"`
}

// ToJSON returns DatabaseDescriptor as JSON
func (dd *DatabaseDescriptor) ToJSON() []byte {
	b, _ := json.Marshal(dd)
	return b
}
