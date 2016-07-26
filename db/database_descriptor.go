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
	MaxDataLogSize string   `xml:"max_datalog_size"`
	MaxCacheSize   string   `xml:"max_cache_size"`
	BloomFilterFp  string   `xml:"bloomfilter_fpp"`
	CronExp        string   `xml:"dataholder_cron_compaction"`
	Path           string   `xml:"path"`
}

// ToJSON returns DatabaseDescriptor as JSON
func (dd *DatabaseDescriptor) ToJSON() []byte {
	b, _ := json.Marshal(dd)
	return b
}
