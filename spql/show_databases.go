package spql

import "encoding/json"

// ShowDatabaseStmt holds database parsed arguments from query
type ShowDatabaseStmt struct {
}

// ParseShowDatabaseStmt parse raw json to query.params
func ParseShowDatabaseStmt(q *Query, raw *json.RawMessage) error {
	q.Method = "ShowDatabases"
	return nil
}
