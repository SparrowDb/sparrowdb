package spql

import (
	"encoding/json"
	"errors"
)

// SelectStmt holds database parsed arguments from query
type SelectStmt struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// ParseSelectStmt parse raw json to query.params
func ParseSelectStmt(q *Query, raw *json.RawMessage) error {
	q.Method = "Select"

	stmt := &SelectStmt{}
	json.Unmarshal(*raw, stmt)
	q.Params = stmt

	var err error

	if err := ValidateDatabaseName.MatchString(stmt.Name); !err {
		return errors.New("Invalid database name")
	}

	if err := ValidateDatabaseName.MatchString(stmt.Key); !err {
		return errors.New("Invalid key")
	}

	return err
}
