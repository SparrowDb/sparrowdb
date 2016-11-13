package model

// CreateDatabase holds database parsed arguments from http request
type CreateDatabase struct {
	MaxDataLogSize uint64  `json:"max_datalog_size"`
	MaxCacheSize   uint64  `json:"max_cache_size"`
	BloomFilterFp  float32 `json:"bloomfilter_fpp"`
	CronExp        string  `json:"dataholder_cron_compaction"`
	Path           string  `json:"path"`
	SnapshotPath   string  `json:"snapshot_path"`
}
