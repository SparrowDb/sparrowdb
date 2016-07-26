package spql

import (
	"encoding/json"
	"errors"
)

// DropDatabaseStmt holds database parsed arguments from query
type DropDatabaseStmt struct {
	Name string `json:"name"`
}

// ParseDropDatabaseStmt parse raw json to query.params
func ParseDropDatabaseStmt(q *Query, raw *json.RawMessage) error {
	stmt := &DropDatabaseStmt{}
	json.Unmarshal(*raw, stmt)
	q.Params = stmt

	err := ValidateDatabaseName.MatchString(stmt.Name)
	if !err {
		return errors.New("Invalid database name")
	}

	return nil
}
