package spql

import (
	"encoding/json"
	"errors"
)

// CreateDatabaseStmt holds database parsed arguments from query
type CreateDatabaseStmt struct {
	Name           string  `json:"name"`
	MaxDataLogSize uint64  `json:"max_datalog_size"`
	MaxCacheSize   uint64  `json:"max_cache_size"`
	BloomFilterFp  float32 `json:"bloomfilter_fpp"`
	CronExp        string  `json:"dataholder_cron_compaction"`
	Path           string  `json:"path"`
	SnapshotPath   string  `json:"snapshot_path"`
}

// ParseCreateDatabaseStmt parse raw json to query.params
func ParseCreateDatabaseStmt(q *Query, raw *json.RawMessage) error {
	stmt := &CreateDatabaseStmt{}
	json.Unmarshal(*raw, stmt)
	q.Params = stmt

	err := ValidateDatabaseName.MatchString(stmt.Name)
	if !err {
		return errors.New("Invalid database name")
	}

	return nil
}
