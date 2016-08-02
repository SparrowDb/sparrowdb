package spql

import (
	"encoding/json"
	"errors"
)

// DeleteStmt holds database parsed arguments from query
type DeleteStmt struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// ParseDeleteStmt parse raw json to query.params
func ParseDeleteStmt(q *Query, raw *json.RawMessage) error {
	stmt := &DeleteStmt{}
	json.Unmarshal(*raw, stmt)
	q.Params = stmt

	if err := ValidateDatabaseName.MatchString(stmt.Name); !err {
		return errors.New("Invalid database name")
	}

	if err := ValidateDatabaseName.MatchString(stmt.Key); !err {
		return errors.New("Invalid key")
	}

	return nil
}
