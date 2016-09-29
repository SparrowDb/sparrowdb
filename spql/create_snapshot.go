package spql

import (
	"encoding/json"
	"errors"
)

// CreateSnapshotStmt holds database parsed arguments from query
type CreateSnapshotStmt struct {
	Name string `json:"name"`
}

// ParseCreateSnapshotStmt parse raw json to query.params
func ParseCreateSnapshotStmt(q *Query, raw *json.RawMessage) error {
	q.Method = "CreateSnapshot"

	stmt := &CreateSnapshotStmt{}
	json.Unmarshal(*raw, stmt)
	q.Params = stmt

	err := ValidateDatabaseName.MatchString(stmt.Name)
	if !err {
		return errors.New("Invalid database name")
	}

	return nil
}
