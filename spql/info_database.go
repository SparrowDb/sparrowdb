package spql

import (
	"encoding/json"
	"errors"
)

// InfoDatabaseStmt holds database parsed arguments from query
type InfoDatabaseStmt struct {
	Name string `json:"name"`
}

// ParseInfoDatabasetStmt parse raw json to query.params
func ParseInfoDatabasetStmt(q *Query, raw *json.RawMessage) error {
	q.Method = "InfoDatabase"

	stmt := &InfoDatabaseStmt{}
	json.Unmarshal(*raw, stmt)
	q.Params = stmt

	err := ValidateDatabaseName.MatchString(stmt.Name)
	if !err {
		return errors.New("Invalid database name")
	}

	return nil
}
